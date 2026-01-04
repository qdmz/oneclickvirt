package payment

import (
	"encoding/json"
	"fmt"
	"oneclickvirt/global"
	orderModel "oneclickvirt/model/order"
	instanceModel "oneclickvirt/model/provider"
	userModel "oneclickvirt/model/user"
	walletModel "oneclickvirt/model/wallet"
	"oneclickvirt/service/cache"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AlipayNotify 支付宝支付回调
// @Summary 支付宝支付回调
// @Description 支付宝支付异步通知
// @Tags 支付
// @Accept json
// @Produce json
// @Router /v1/payment/alipay/notify [post]
func AlipayNotify(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		global.APP_LOG.Error("解析支付宝回调数据失败", zap.Error(err))
		c.JSON(400, gin.H{"code": 400, "message": "数据格式错误"})
		return
	}

	// TODO: 验证支付宝签名

	// 获取订单号
	orderNo, ok := data["trade_no"].(string)
	if !ok || orderNo == "" {
		global.APP_LOG.Error("订单号为空")
		c.JSON(400, gin.H{"code": 400, "message": "订单号为空"})
		return
	}

	// 处理支付成功
	if err := processPaymentSuccess(orderNo, orderModel.PaymentMethodAlipay, data); err != nil {
		global.APP_LOG.Error("处理支付成功失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "处理失败"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "success"})
}

// WechatNotify 微信支付回调
// @Summary 微信支付回调
// @Description 微信支付异步通知
// @Tags 支付
// @Accept json
// @Produce json
// @Router /v1/payment/wechat/notify [post]
func WechatNotify(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		global.APP_LOG.Error("解析微信回调数据失败", zap.Error(err))
		c.JSON(400, gin.H{"code": 400, "message": "数据格式错误"})
		return
	}

	// TODO: 验证微信签名

	// 获取订单号
	orderNo, ok := data["out_trade_no"].(string)
	if !ok || orderNo == "" {
		global.APP_LOG.Error("订单号为空")
		c.JSON(400, gin.H{"code": 400, "message": "订单号为空"})
		return
	}

	// 处理支付成功
	if err := processPaymentSuccess(orderNo, orderModel.PaymentMethodWechat, data); err != nil {
		global.APP_LOG.Error("处理支付成功失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "处理失败"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "success"})
}

// processPaymentSuccess 处理支付成功
func processPaymentSuccess(orderNo string, paymentMethod string, notifyData map[string]interface{}) error {
	tx := global.APP_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询订单
	var order orderModel.Order
	if err := tx.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return err
	}

	// 幂等性检查:如果订单已支付,直接返回成功
	if order.Status == orderModel.OrderStatusPaid {
		tx.Rollback()
		return nil
	}

	// 检查订单状态
	if order.Status != orderModel.OrderStatusPending {
		tx.Rollback()
		return nil
	}

	// 更新订单状态
	order.Status = orderModel.OrderStatusPaid
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

	paymentRecord := orderModel.PaymentRecord{
		OrderID:       order.ID,
		UserID:        order.UserID,
		Type:          paymentMethod,
		TransactionID: orderNo, // 这里简化处理,实际应该是第三方交易号
		Amount:        order.Amount,
		Status:        orderModel.PaymentStatusSuccess,
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
