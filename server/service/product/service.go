package product

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"oneclickvirt/global"
	"oneclickvirt/model/common"
	"oneclickvirt/model/order"
	"oneclickvirt/model/product"
	userModel "oneclickvirt/model/user"
	"oneclickvirt/model/wallet"
	paymentSvc "oneclickvirt/service/payment"
	"oneclickvirt/service/resources"
	userService "oneclickvirt/service/user"
	"oneclickvirt/service/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Service 产品/订单服务
type Service struct{}

func NewService() *Service { return &Service{} }

// ListProducts 列出已上架产品（用户端）。
func (s *Service) ListProducts(page, pageSize int, instanceType string) ([]product.Product, int64, error) {
	var list []product.Product
	var total int64
	db := global.APP_DB.Model(&product.Product{}).Where("status = ?", product.ProductStatusOn)
	if instanceType != "" {
		db = db.Where("type = ?", instanceType)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}
	if err := db.Order("sort DESC, id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// AdminList 管理端列表（含下架）。
func (s *Service) AdminList(page, pageSize int) ([]product.Product, int64, error) {
	var list []product.Product
	var total int64
	db := global.APP_DB.Model(&product.Product{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 20
	}
	if err := db.Order("sort DESC, id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Get 获取单个产品（管理端）。
func (s *Service) Get(id uint) (*product.Product, error) {
	var p product.Product
	if err := global.APP_DB.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// AdminCreate 创建产品。
func (s *Service) AdminCreate(p *product.Product) (*product.Product, error) {
	if p.Name == "" {
		return nil, common.NewError(common.CodeValidationError, "产品名称不能为空")
	}
	if p.Price < 0 {
		return nil, common.NewError(common.CodeValidationError, "价格不能为负")
	}
	if err := global.APP_DB.Create(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

// AdminUpdate 更新产品。
func (s *Service) AdminUpdate(id uint, upd map[string]interface{}) error {
	return global.APP_DB.Model(&product.Product{}).Where("id = ?", id).Updates(upd).Error
}

// AdminDelete 删除产品（软删）。
func (s *Service) AdminDelete(id uint) error {
	return global.APP_DB.Delete(&product.Product{}, id).Error
}

// BuyRequest 购买请求
type BuyRequest struct {
	Method      string `json:"method" binding:"required,oneof=balance epay"`
	NodeID      uint   `json:"nodeId"`
	ImageID     uint   `json:"imageId"`
	PeriodDays  int    `json:"periodDays"`
	RenewalOf   *uint  `json:"renewalOf"`
	ClientBaseURL string `json:"-"`
}

// BuyResult 购买结果
type BuyResult struct {
	Order      *order.Order
	EpayParams map[string]string
}

// Buy 购买产品。
func (s *Service) Buy(userID uint, productID uint, req BuyRequest) (*BuyResult, error) {
	var p product.Product
	if err := global.APP_DB.First(&p, productID).Error; err != nil {
		return nil, common.NewError(common.CodeNotFound, "产品不存在")
	}
	if p.Status != product.ProductStatusOn {
		return nil, common.NewError(common.CodeValidationError, "产品已下架")
	}
	if p.Stock == 0 {
		return nil, common.NewError(common.CodeValidationError, "产品已售罄")
	}

	// 节点解析：产品可绑定固定节点，否则由用户下单时选择。
	// 若两者均未提供，则视为「仅购买资源额度」，支付后额度立即到账，
	// 用户可稍后在「创建实例」页自行选择节点开通。
	nodeID := p.NodeID
	if nodeID == 0 {
		nodeID = req.NodeID
	}
	imageID := req.ImageID
	if req.RenewalOf == nil && nodeID != 0 && imageID == 0 {
		return nil, common.NewError(common.CodeValidationError, "已选择开通节点，请同时选择操作系统镜像")
	}

	period := req.PeriodDays
	if period <= 0 {
		period = p.PeriodDays
	}

	tradeNo := paymentSvc.GenerateTradeNo("PD")

	if req.Method == order.OrderMethodBalance {
		var u userModel.User
		if err := global.APP_DB.Select("balance").First(&u, userID).Error; err != nil {
			return nil, err
		}
		if u.Balance < p.Price {
			return nil, common.NewError(common.CodeValidationError,
				fmt.Sprintf("余额不足：当前 ¥%s，需支付 ¥%s，请先充值",
					paymentSvc.FormatMoney(u.Balance), paymentSvc.FormatMoney(p.Price)))
		}
	}

	ord := &order.Order{
		UserID:     userID,
		ProductID:  productID,
		TradeNo:    tradeNo,
		Amount:     p.Price,
		Method:     req.Method,
		Status:     order.OrderStatusPending,
		PeriodDays: period,
		RenewalOf:  req.RenewalOf,
	}

	err := database.GetDatabaseService().ExecuteTransaction(context.Background(), func(tx *gorm.DB) error {
		if err := tx.Create(ord).Error; err != nil {
			return err
		}
		if req.Method == order.OrderMethodBalance {
			if err := tx.Model(&userModel.User{}).Where("id = ?", userID).
				Update("balance", gorm.Expr("balance - ?", p.Price)).Error; err != nil {
				return fmt.Errorf("扣减余额失败: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if req.Method == order.OrderMethodBalance {
		if err := s.Fulfill(ord, &p, userID, nodeID, imageID); err != nil {
			global.APP_LOG.Error("余额购买履约失败", zap.Uint("orderID", ord.ID), zap.Error(err))
			return nil, fmt.Errorf("支付成功但开通失败: %w", err)
		}
		return &BuyResult{Order: ord}, nil
	}

	notifyURL := global.GetAppConfig().Payment.EpayNotifyURL
	if notifyURL == "" {
		notifyURL = joinURL(req.ClientBaseURL, "/api/v1/payment/epay/notify")
	}
	returnURL := global.GetAppConfig().Payment.EpayReturnURL
	if returnURL == "" {
		returnURL = joinURL(req.ClientBaseURL, "/api/v1/payment/epay/return")
	}
	params, err := paymentSvc.BuildEpaySubmitParams(tradeNo, p.Name, paymentSvc.FormatMoney(p.Price), notifyURL, returnURL)
	if err != nil {
		return nil, err
	}
	return &BuyResult{Order: ord, EpayParams: params}, nil
}

// Fulfill 履约：授予用户资源并开通实例（best-effort）。
func (s *Service) Fulfill(ord *order.Order, p *product.Product, userID, nodeID, imageID uint) error {
	now := time.Now()
	expiresAt := now.AddDate(0, 0, ord.PeriodDays)

	err := database.GetDatabaseService().ExecuteTransaction(context.Background(), func(tx *gorm.DB) error {
		if err := grantUserResources(tx, userID, p); err != nil {
			return err
		}
		if err := tx.Model(&order.Order{}).Where("id = ?", ord.ID).Updates(map[string]interface{}{
			"status":     order.OrderStatusPaid,
			"paid_at":    now,
			"expires_at": expiresAt,
		}).Error; err != nil {
			return err
		}
		if err := tx.Model(&product.Product{}).Where("id = ?", p.ID).
			Updates(map[string]interface{}{
				"sold":  gorm.Expr("sold + 1"),
				"stock": gorm.Expr("CASE WHEN stock > 0 THEN stock - 1 ELSE stock END"),
			}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	if nodeID == 0 || imageID == 0 {
		// 续费 / 仅购买额度：资源已到账，跳过实例开通
		_ = global.APP_DB.Model(&order.Order{}).Where("id = ?", ord.ID).
			Update("remark", "资源额度已到账，可前往「创建实例」自行开通").Error
		return nil
	}

	req := userModel.CreateInstanceRequest{
		ProviderId:  nodeID,
		ImageId:     imageID,
		CPUId:       fmt.Sprintf("cpu-%d", p.CPU),
		MemoryId:    fmt.Sprintf("memory-%dmb", p.MemoryMB),
		DiskId:      fmt.Sprintf("disk-%dmb", p.DiskMB),
		BandwidthId: fmt.Sprintf("bandwidth-%dmb", p.BandwidthMbps),
		Description: fmt.Sprintf("购买产品：%s", p.Name),
	}
	task, err := userService.NewService().CreateUserInstance(userID, req)
	if err != nil {
		global.APP_LOG.Warn("购买后开通实例失败（资源已授予）",
			zap.Uint("userID", userID), zap.Uint("productID", p.ID), zap.Error(err))
		_ = global.APP_DB.Model(&order.Order{}).Where("id = ?", ord.ID).
			Update("remark", "资源已授予，实例开通待节点就绪: "+err.Error()).Error
		return nil
	}
	_ = global.APP_DB.Model(&order.Order{}).Where("id = ?", ord.ID).
		Update("instance_id", task.ID).Error
	return nil
}

// grantUserResources 提升用户配额上限以容纳本次购买（购买即得资源）。
func grantUserResources(tx *gorm.DB, userID uint, p *product.Product) error {
	var u userModel.User
	if err := tx.Select("max_instances", "max_cpu", "max_memory", "max_disk", "max_bandwidth").
		First(&u, userID).Error; err != nil {
		return err
	}
	upd := map[string]interface{}{}
	if p.CPU > u.MaxCPU {
		upd["max_cpu"] = p.CPU
	}
	if p.MemoryMB > u.MaxMemory {
		upd["max_memory"] = p.MemoryMB
	}
	if p.DiskMB > u.MaxDisk {
		upd["max_disk"] = p.DiskMB
	}
	if p.BandwidthMbps > u.MaxBandwidth {
		upd["max_bandwidth"] = p.BandwidthMbps
	}
	if u.MaxInstances < 1 {
		upd["max_instances"] = 1
	}
	if len(upd) > 0 {
		if err := tx.Model(&userModel.User{}).Where("id = ?", userID).Updates(upd).Error; err != nil {
			return err
		}
	}
	quotaService := resources.NewQuotaService()
	return quotaService.RecalculateUserQuotaInTx(tx, userID)
}

// HandleEpayNotify 处理购买订单的易支付回调并履约。
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
		_ = global.APP_DB.Model(&order.Order{}).Where("trade_no = ?", outTradeNo).
			Update("status", order.OrderStatusFailed).Error
		return fmt.Errorf("交易状态非成功: %s", tradeStatus)
	}

	var ord order.Order
	if err := global.APP_DB.Where("trade_no = ?", outTradeNo).First(&ord).Error; err != nil {
		return fmt.Errorf("订单不存在: %w", err)
	}
	if ord.Status == order.OrderStatusPaid {
		return nil
	}
	var p product.Product
	if err := global.APP_DB.First(&p, ord.ProductID).Error; err != nil {
		return err
	}
	nodeID := p.NodeID
	if nodeID == 0 {
		now := time.Now()
		_ = grantUserResources(global.APP_DB, ord.UserID, &p)
		_ = global.APP_DB.Model(&order.Order{}).Where("id = ?", ord.ID).Updates(map[string]interface{}{
			"status":     order.OrderStatusPaid,
			"paid_at":    now,
			"expires_at": now.AddDate(0, 0, ord.PeriodDays),
			"remark":     "已支付，待选择节点开通",
		}).Error
		return nil
	}
	return s.Fulfill(&ord, &p, ord.UserID, nodeID, 0)
}

// MyOrders 用户订单列表。
func (s *Service) MyOrders(userID uint, page, pageSize int) ([]order.Order, int64, error) {
	var list []order.Order
	var total int64
	db := global.APP_DB.Model(&order.Order{}).Where("user_id = ?", userID)
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

func joinURL(base, path string) string {
	if base == "" {
		return path
	}
	return strings.TrimRight(base, "/") + path
}

var _ = wallet.RechargeOrder{}
