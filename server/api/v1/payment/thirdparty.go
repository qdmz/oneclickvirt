package payment

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"oneclickvirt/global"
	"oneclickvirt/model/order"
	"oneclickvirt/service/cache"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 易支付相关配置
type EpayConfig struct {
	APIURL    string `json:"api_url"`
	PID       string `json:"pid"`
	Key       string `json:"key"`
	ReturnURL string `json:"return_url"`
	NotifyURL string `json:"notify_url"`
	Type      string `json:"type"` // alipay, wechat, qqpay
}

// 码支付相关配置
type MapayConfig struct {
	APIURL    string `json:"api_url"`
	ID        string `json:"id"`
	Key       string `json:"key"`
	ReturnURL string `json:"return_url"`
	NotifyURL string `json:"notify_url"`
	Type      string `json:"type"` // alipay, wechat, qqpay
}

// EpayNotify 易支付回调
// @Summary 易支付回调
// @Description 易支付异步通知
// @Tags 支付
// @Accept json
// @Produce json
// @Router /v1/payment/epay/notify [post]
func EpayNotify(c *gin.Context) {
	// 解析表单数据
	if err := c.Request.ParseForm(); err != nil {
		global.APP_LOG.Error("解析易支付回调数据失败", zap.Error(err))
		c.String(http.StatusBadRequest, "fail")
		return
	}

	params := c.Request.Form
	sign := params.Get("sign")
	
	// 验证签名
	if !verifyEpaySign(params, sign) {
		global.APP_LOG.Error("易支付签名验证失败")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 获取订单号
	orderNo := params.Get("out_trade_no")
	if orderNo == "" {
		global.APP_LOG.Error("易支付订单号为空")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 处理支付成功
	if err := processPaymentSuccess(orderNo, order.PaymentMethodEpay, params); err != nil {
		global.APP_LOG.Error("处理易支付成功失败", zap.Error(err))
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// MapayNotify 码支付回调
// @Summary 码支付回调
// @Description 码支付异步通知
// @Tags 支付
// @Accept json
// @Produce json
// @Router /v1/payment/mapay/notify [post]
func MapayNotify(c *gin.Context) {
	// 解析表单数据
	if err := c.Request.ParseForm(); err != nil {
		global.APP_LOG.Error("解析码支付回调数据失败", zap.Error(err))
		c.String(http.StatusBadRequest, "fail")
		return
	}

	params := c.Request.Form
	sign := params.Get("sign")
	
	// 验证签名
	if !verifyMapaySign(params, sign) {
		global.APP_LOG.Error("码支付签名验证失败")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 获取订单号
	orderNo := params.Get("out_trade_no")
	if orderNo == "" {
		global.APP_LOG.Error("码支付订单号为空")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 处理支付成功
	if err := processPaymentSuccess(orderNo, order.PaymentMethodMapay, params); err != nil {
		global.APP_LOG.Error("处理码支付成功失败", zap.Error(err))
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// verifyEpaySign 验证易支付签名
func verifyEpaySign(params url.Values, sign string) bool {
	// 构建待签名字符串
	var keys []string
	for key := range params {
		if key != "sign" && key != "sign_type" && params.Get(key) != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var signStr strings.Builder
	for i, key := range keys {
		if i > 0 {
			signStr.WriteString("&")
		}
		signStr.WriteString(key)
		signStr.WriteString("=")
		signStr.WriteString(params.Get(key))
	}
	signStr.WriteString(global.APP_CONFIG.Payment.EpayKey)

	// MD5签名
	h := md5.New()
	h.Write([]byte(signStr.String()))
	calculatedSign := hex.EncodeToString(h.Sum(nil))

	return strings.ToLower(calculatedSign) == strings.ToLower(sign)
}

// verifyMapaySign 验证码支付签名
func verifyMapaySign(params url.Values, sign string) bool {
	// 构建待签名字符串
	var keys []string
	for key := range params {
		if key != "sign" && key != "sign_type" && params.Get(key) != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var signStr strings.Builder
	for i, key := range keys {
		if i > 0 {
			signStr.WriteString("&")
		}
		signStr.WriteString(key)
		signStr.WriteString("=")
		signStr.WriteString(params.Get(key))
	}
	signStr.WriteString(global.APP_CONFIG.Payment.MapayKey)

	// MD5签名
	h := md5.New()
	h.Write([]byte(signStr.String()))
	calculatedSign := hex.EncodeToString(h.Sum(nil))

	return strings.ToLower(calculatedSign) == strings.ToLower(sign)
}

// processPaymentSuccess 处理支付成功
func processPaymentSuccess(orderNo string, paymentMethod string, notifyData url.Values) error {
	tx := global.APP_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询订单
	var order order.Order
	if err := tx.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return err
	}

	// 幂等性检查:如果订单已支付,直接返回成功
	if order.Status == order.OrderStatusPaid {
		tx.Rollback()
		return nil
	}

	// 检查订单状态
	if order.Status != order.OrderStatusPending {
		tx.Rollback()
		return nil
	}

	// 更新订单状态
	order.Status = order.OrderStatusPaid
	now := time.Now()
	order.PaymentTime = &now
	order.PaidAmount = order.Amount

	if err := tx.Save(&order).Error; err != nil {
		return err
	}

	// 创建支付记录
	var notifyDataStr string
	if data, err := json.Marshal(notifyData); err == nil {
		notifyDataStr = string(data)
	}

	paymentRecord := order.PaymentRecord{
		OrderID:       order.ID,
		UserID:        order.UserID,
		Type:          paymentMethod,
		TransactionID: orderNo, // 这里简化处理,实际应该是第三方交易号
		Amount:        order.Amount,
		Status:        order.PaymentStatusSuccess,
		NotifyData:    notifyDataStr,
	}

	if err := tx.Create(&paymentRecord).Error; err != nil {
		return err
	}

	// 处理订单类型
	if order.ProductID != nil {
		// 产品购买:提升用户等级
		var user userModel.User
		if err := tx.First(&user, order.UserID).Error; err != nil {
			return err
		}

		// 解析产品数据获取等级和有效期
		var productData map[string]interface{}
		if err := json.Unmarshal([]byte(order.ProductData), &productData); err == nil {
			if level, ok := productData["level"].(float64); ok {
				newLevel := int(level)
				user.Level = newLevel

				// 获取产品有效期
				var period int
				if p, ok := productData["period"].(float64); ok && p > 0 {
					period = int(p)
				}

				// 计算有效期
				var expireTime *time.Time
				if period > 0 {
					t := time.Now().AddDate(0, 0, period)
					expireTime = &t
				}

				// 更新有效期：当前到期日期+产品有效期
				if expireTime != nil {
					now := time.Now()
					var newExpire time.Time
					if user.LevelExpireAt != nil && user.LevelExpireAt.After(now) {
						// 如果用户已有有效的到期时间，在其基础上延长
						newExpire = user.LevelExpireAt.AddDate(0, 0, period)
						global.APP_LOG.Info(fmt.Sprintf("用户已有有效到期时间，在基础上延长: %v + %d天 = %v", user.LevelExpireAt, period, newExpire))
					} else {
						// 如果没有有效的到期时间，从当前时间开始计算
						newExpire = time.Now().AddDate(0, 0, period)
						global.APP_LOG.Info(fmt.Sprintf("用户没有有效到期时间，从当前时间开始计算: %v + %d天 = %v", now, period, newExpire))
					}
					user.LevelExpireAt = &newExpire
				} else {
					// 永久产品,设置为9999年后
					farFuture := time.Now().AddDate(9999, 0, 0)
					user.LevelExpireAt = &farFuture
					global.APP_LOG.Info(fmt.Sprintf("永久产品，设置到期时间为: %v", farFuture))
				}

				// 更新用户名下所有实例的到期时间
				if err := tx.Model(&instanceModel.Instance{}).Where("user_id = ?", order.UserID).Update("expired_at", user.LevelExpireAt).Error; err != nil {
					return err
				}

				if err := tx.Save(&user).Error; err != nil {
					return err
				}

				// 清除用户Dashboard缓存，确保用户下次访问个人中心时能看到最新的用户信息
				cacheService := cache.GetUserCacheService()
				cacheService.InvalidateUserCache(order.UserID)
				global.APP_LOG.Info(fmt.Sprintf("已清除用户 %d 的Dashboard缓存", order.UserID))
			}
		}
	} else {
		// 充值订单:增加钱包余额
		var wallet walletModel.UserWallet
		if err := tx.Where("user_id = ?", order.UserID).First(&wallet).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 创建钱包
				wallet = walletModel.UserWallet{
					UserID:        order.UserID,
					Balance:       0,
					Frozen:        0,
					TotalRecharge: 0,
					TotalExpense:  0,
				}
				if err := tx.Create(&wallet).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// 增加余额
		wallet.Balance += order.Amount
		wallet.TotalRecharge += order.Amount
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		// 创建交易记录
		transaction := walletModel.WalletTransaction{
			UserID:      order.UserID,
			Type:        walletModel.TransactionTypeRecharge,
			Amount:      order.Amount,
			Balance:     wallet.Balance,
			Description: "在线充值",
			OrderID:     &order.ID,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	global.APP_LOG.Info("订单支付成功",
		zap.String("orderNo", orderNo),
		zap.Uint("userId", order.UserID),
		zap.Int64("amount", order.Amount),
	)

	return nil
}