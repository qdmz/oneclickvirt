package wallet

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"oneclickvirt/global"
	"oneclickvirt/model/common"
	"oneclickvirt/model/order"
	"oneclickvirt/model/user"
	"oneclickvirt/model/wallet"
	paymentSvc "oneclickvirt/service/payment"
	"oneclickvirt/service/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Service 钱包/充值服务
type Service struct{}

func NewService() *Service { return &Service{} }

// CreateRechargeOrder 创建余额充值订单。
// 返回订单与易支付跳转参数（若易支付未启用或配置缺失，返回错误）。
func (s *Service) CreateRechargeOrder(userID uint, amountCents int64, clientBaseURL string) (*wallet.RechargeOrder, map[string]string, error) {
	if amountCents <= 0 {
		return nil, nil, common.NewError(common.CodeValidationError, "充值金额必须大于 0")
	}
	if amountCents > 1000000000 { // 上限 1000 万元
		return nil, nil, common.NewError(common.CodeValidationError, "充值金额过大")
	}
	cfg := global.GetAppConfig().Payment
	if !cfg.EpayEnabled {
		return nil, nil, common.NewError(common.CodeUnavailable, "在线充值未开启，请联系管理员")
	}
	if cfg.EpayPID == "" || cfg.EpayKey == "" || cfg.EpayApiURL == "" {
		return nil, nil, common.NewError(common.CodeConfigError, "易支付参数未配置完整（PID / KEY / 网关地址）")
	}

	tradeNo := paymentSvc.GenerateTradeNo("RC")
	order := &wallet.RechargeOrder{
		UserID:  userID,
		TradeNo: tradeNo,
		Amount:  amountCents,
		Method:  wallet.RechargeMethodEpay,
		Status:  wallet.RechargeStatusPending,
	}

	if err := global.APP_DB.Create(order).Error; err != nil {
		return nil, nil, fmt.Errorf("创建充值订单失败: %w", err)
	}

	notifyURL := cfg.EpayNotifyURL
	if notifyURL == "" {
		notifyURL = joinURL(clientBaseURL, "/api/v1/payment/epay/notify")
	}
	returnURL := cfg.EpayReturnURL
	if returnURL == "" {
		returnURL = joinURL(clientBaseURL, "/api/v1/payment/epay/return")
	}

	params, err := paymentSvc.BuildEpaySubmitParams(tradeNo, "余额充值", paymentSvc.FormatMoney(amountCents), notifyURL, returnURL)
	if err != nil {
		return nil, nil, err
	}
	return order, params, nil
}

// HandleEpayNotify 处理易支付异步回调。
func (s *Service) HandleEpayNotify(params map[string]string) error {
	cfg := global.GetAppConfig().Payment
	if !paymentSvc.EpayVerify(params, cfg.EpayKey) {
		return errors.New("易支付签名校验失败")
	}
	outTradeNo := params["out_trade_no"]
	if outTradeNo == "" {
		return errors.New("缺少订单号")
	}
	tradeStatus := params["trade_status"]
	if tradeStatus != "TRADE_SUCCESS" && tradeStatus != "SUCCESS" {
		_ = global.APP_DB.Model(&wallet.RechargeOrder{}).
			Where("trade_no = ?", outTradeNo).
			Update("status", wallet.RechargeStatusFailed).Error
		return fmt.Errorf("交易状态非成功: %s", tradeStatus)
	}

	var order wallet.RechargeOrder
	if err := global.APP_DB.Where("trade_no = ?", outTradeNo).First(&order).Error; err != nil {
		return fmt.Errorf("充值订单不存在: %w", err)
	}
	if order.Status == wallet.RechargeStatusPaid {
		return nil // 幂等：已处理
	}

	now := time.Now()
	err := database.GetDatabaseService().ExecuteTransaction(context.Background(), func(tx *gorm.DB) error {
		if err := tx.Model(&user.User{}).Where("id = ?", order.UserID).
			Update("balance", gorm.Expr("balance + ?", order.Amount)).Error; err != nil {
			return fmt.Errorf("更新余额失败: %w", err)
		}
		if err := tx.Model(&wallet.RechargeOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
			"status":        wallet.RechargeStatusPaid,
			"epay_trade_no": params["trade_no"],
			"paid_at":       now,
		}).Error; err != nil {
			return fmt.Errorf("更新订单状态失败: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	global.APP_LOG.Info("易支付充值成功入账",
		zap.Uint("userID", order.UserID),
		zap.String("tradeNo", outTradeNo),
		zap.Int64("amount", order.Amount))
	return nil
}

// AdjustBalance 管理员手动调整用户余额（deltaCents 可正可负），
// 同时写入一条 method=admin 的余额变动记录，便于审计。
func (s *Service) AdjustBalance(userID uint, deltaCents int64, remark string) (int64, error) {
	if deltaCents == 0 {
		return 0, common.NewError(common.CodeValidationError, "调整金额不能为 0")
	}
	var u user.User
	if err := global.APP_DB.Select("id", "balance").First(&u, userID).Error; err != nil {
		return 0, common.NewError(common.CodeNotFound, "用户不存在")
	}
	if u.Balance+deltaCents < 0 {
		return 0, common.NewError(common.CodeValidationError,
			fmt.Sprintf("扣减后余额为负：当前 ¥%s", paymentSvc.FormatMoney(u.Balance)))
	}

	now := time.Now()
	err := database.GetDatabaseService().ExecuteTransaction(context.Background(), func(tx *gorm.DB) error {
		if err := tx.Model(&user.User{}).Where("id = ?", userID).
			Update("balance", gorm.Expr("balance + ?", deltaCents)).Error; err != nil {
			return fmt.Errorf("更新余额失败: %w", err)
		}
		rec := &wallet.RechargeOrder{
			UserID:  userID,
			TradeNo: paymentSvc.GenerateTradeNo("AD"),
			Amount:  deltaCents,
			Method:  wallet.RechargeMethodAdmin,
			Status:  wallet.RechargeStatusPaid,
			Remark:  remark,
			PaidAt:  &now,
		}
		return tx.Create(rec).Error
	})
	if err != nil {
		return 0, err
	}
	global.APP_LOG.Info("管理员调整用户余额",
		zap.Uint("userID", userID), zap.Int64("delta", deltaCents), zap.String("remark", remark))
	return u.Balance + deltaCents, nil
}

// GetBalance 获取用户余额（分）。
func (s *Service) GetBalance(userID uint) (int64, error) {
	var u user.User
	if err := global.APP_DB.Select("balance").First(&u, userID).Error; err != nil {
		return 0, err
	}
	return u.Balance, nil
}

// ListUserRechargeOrders 分页查询用户充值订单。
func (s *Service) ListUserRechargeOrders(userID uint, page, pageSize int) ([]wallet.RechargeOrder, int64, error) {
	var list []wallet.RechargeOrder
	var total int64
	db := global.APP_DB.Model(&wallet.RechargeOrder{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// AdminListRechargeOrders 管理端余额流水（userID=0 表示全部）。
func (s *Service) AdminListRechargeOrders(userID uint, page, pageSize int) ([]wallet.RechargeOrder, int64, error) {
	var list []wallet.RechargeOrder
	var total int64
	db := global.APP_DB.Model(&wallet.RechargeOrder{})
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 20
	}
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// AdminListOrders 管理端商品订单列表（userID=0 表示全部）。
func (s *Service) AdminListOrders(userID uint, status string, page, pageSize int) ([]order.Order, int64, error) {
	var list []order.Order
	var total int64
	db := global.APP_DB.Model(&order.Order{})
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 20
	}
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func joinURL(base, path string) string {
	if base == "" {
		return path
	}
	return strings.TrimRight(base, "/") + path
}
