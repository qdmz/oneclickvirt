package user

import (
	"encoding/json"
	"fmt"
	"oneclickvirt/global"
	orderModel "oneclickvirt/model/order"
	redemptionModel "oneclickvirt/model/redemption"
	userModel "oneclickvirt/model/user"
	walletModel "oneclickvirt/model/wallet"

	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateRechargeOrder 创建充值订单
// @Summary 创建充值订单
// @Description 用户创建充值订单
// @Tags 用户/充值
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "充值信息: amount(金额,分), paymentMethod(支付方式)"
// @Success 200 {object} common.Response
// @Router /v1/user/recharge/create-order [post]
func CreateRechargeOrder(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var params struct {
		Amount        int64  `json:"amount" binding:"required"`
		PaymentMethod string `json:"paymentMethod" binding:"required"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 验证金额
	if params.Amount <= 0 {
		c.JSON(400, gin.H{"code": 400, "message": "金额必须大于0"})
		return
	}

	// 验证支付方式
	if params.PaymentMethod != orderModel.PaymentMethodAlipay &&
		params.PaymentMethod != orderModel.PaymentMethodWechat &&
		params.PaymentMethod != orderModel.PaymentMethodBalance {
		c.JSON(400, gin.H{"code": 400, "message": "不支持的支付方式"})
		return
	}

	// 生成订单号
	orderNo := generateOrderNo()

	// 创建订单
	order := orderModel.Order{
		OrderNo:       orderNo,
		UserID:        userID.(uint),
		Amount:        params.Amount,
		Status:        orderModel.OrderStatusPending,
		PaymentMethod: params.PaymentMethod,
		PaidAmount:    0,
		ProductData:   "{}", // 充值订单设置默认空JSON对象，避免约束失败
		ExpireAt:      time.Now().Add(30 * time.Minute), // 30分钟过期
	}

	if err := global.APP_DB.Create(&order).Error; err != nil {
		global.APP_LOG.Error("创建充值订单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "创建订单失败: " + err.Error()})
		return
	}

	// 返回订单信息
	c.JSON(200, gin.H{
		"code": 200,
		"message": "创建订单成功",
		"data": order,
	})
}

// GetRechargeAlipayQR 获取支付宝支付二维码
// @Summary 获取支付宝支付二维码
// @Description 获取支付宝支付二维码URL
// @Tags 用户/充值
// @Accept json
// @Produce json
// @Param orderNo path string true "订单号"
// @Success 200 {object} common.Response
// @Router /v1/user/recharge/alipay-qr/{orderNo} [get]
func GetRechargeAlipayQR(c *gin.Context) {
	orderNo := c.Param("orderNo")
	if orderNo == "" {
		c.JSON(400, gin.H{"code": 400, "message": "订单号不能为空"})
		return
	}

	var order orderModel.Order
	if err := global.APP_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "订单不存在"})
		return
	}

	// 验证订单状态
	if order.Status != orderModel.OrderStatusPending {
		c.JSON(400, gin.H{"code": 400, "message": "订单状态异常"})
		return
	}

	// 检查订单是否过期
	if time.Now().After(order.ExpireAt) {
		order.Status = orderModel.OrderStatusExpired
		global.APP_DB.Save(&order)
		c.JSON(400, gin.H{"code": 400, "message": "订单已过期"})
		return
	}

	// 这里应该是调用支付宝SDK生成支付二维码URL
	// 由于需要配置支付宝商户信息，这里返回模拟数据
	qrCode := "https://qr.alipay.com/" + orderNo

	c.JSON(200, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"qrCode":  qrCode,
			"orderNo": order.OrderNo,
			"amount":  order.Amount,
			"expireAt": order.ExpireAt,
		},
	})
}

// GetRechargeWechatQR 获取微信支付二维码
// @Summary 获取微信支付二维码
// @Description 获取微信支付二维码URL
// @Tags 用户/充值
// @Accept json
// @Produce json
// @Param orderNo path string true "订单号"
// @Success 200 {object} common.Response
// @Router /v1/user/recharge/wechat-qr/{orderNo} [get]
func GetRechargeWechatQR(c *gin.Context) {
	orderNo := c.Param("orderNo")
	if orderNo == "" {
		c.JSON(400, gin.H{"code": 400, "message": "订单号不能为空"})
		return
	}

	var order orderModel.Order
	if err := global.APP_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "订单不存在"})
		return
	}

	// 验证订单状态
	if order.Status != orderModel.OrderStatusPending {
		c.JSON(400, gin.H{"code": 400, "message": "订单状态异常"})
		return
	}

	// 检查订单是否过期
	if time.Now().After(order.ExpireAt) {
		order.Status = orderModel.OrderStatusExpired
		global.APP_DB.Save(&order)
		c.JSON(400, gin.H{"code": 400, "message": "订单已过期"})
		return
	}

	// 这里应该是调用微信支付SDK生成支付二维码URL
	// 由于需要配置微信支付商户信息,这里返回模拟数据
	qrCode := "weixin://wxpay/bizpayurl?pr=" + orderNo

	c.JSON(200, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"qrCode":  qrCode,
			"orderNo": order.OrderNo,
			"amount":  order.Amount,
			"expireAt": order.ExpireAt,
		},
	})
}

// ExchangeRedemptionCode 使用兑换码
// @Summary 使用兑换码
// @Description 用户使用兑换码充值
// @Tags 用户/充值
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "兑换码"
// @Success 200 {object} common.Response
// @Router /v1/user/recharge/exchange-code [post]
func ExchangeRedemptionCode(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var params struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if params.Code == "" {
		c.JSON(400, gin.H{"code": 400, "message": "兑换码不能为空"})
		return
	}

	tx := global.APP_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// 确保事务在函数结束时正确关闭
	var committed bool
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 查询兑换码
	var redemptionCode redemptionModel.RedemptionCode
	if err := tx.Where("code = ?", params.Code).First(&redemptionCode).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"code": 404, "message": "兑换码不存在"})
			return
		}
		c.JSON(500, gin.H{"code": 500, "message": "查询兑换码失败"})
		return
	}

	// 检查兑换码是否启用
	if !redemptionCode.IsEnabled {
		tx.Rollback()
		c.JSON(400, gin.H{"code": 400, "message": "兑换码已禁用"})
		return
	}

	// 检查兑换码是否过期
	if redemptionCode.ExpireAt != nil && time.Now().After(*redemptionCode.ExpireAt) {
		tx.Rollback()
		c.JSON(400, gin.H{"code": 400, "message": "兑换码已过期"})
		return
	}

	// 检查使用次数
	if redemptionCode.UsedCount >= redemptionCode.MaxUses {
		tx.Rollback()
		c.JSON(400, gin.H{"code": 400, "message": "兑换码使用次数已用完"})
		return
	}

	// 检查用户是否已经使用过
	var existingUsage redemptionModel.RedemptionCodeUsage
	if err := tx.Where("code_id = ? AND user_id = ?", redemptionCode.ID, userID).First(&existingUsage).Error; err == nil {
		tx.Rollback()
		c.JSON(400, gin.H{"code": 400, "message": "您已经使用过此兑换码"})
		return
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		c.JSON(500, gin.H{"code": 500, "message": "查询兑换码使用记录失败"})
		return
	}

	// 获取或创建用户钱包
	var wallet walletModel.UserWallet
	if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建钱包
			wallet = walletModel.UserWallet{
				UserID:        userID.(uint),
				Balance:       0,
				Frozen:        0,
				TotalRecharge: 0,
				TotalExpense:  0,
			}
			if err := tx.Create(&wallet).Error; err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{"code": 500, "message": "创建钱包失败"})
				return
			}
		} else {
			tx.Rollback()
			c.JSON(500, gin.H{"code": 500, "message": "查询钱包失败"})
			return
		}
	}

	// 根据兑换码类型处理奖励
	var reward map[string]interface{}
	switch redemptionCode.Type {
	case redemptionModel.RedemptionTypeBalance:
		// 余额充值
		wallet.Balance += redemptionCode.Amount
		wallet.TotalRecharge += redemptionCode.Amount
		if err := tx.Save(&wallet).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"code": 500, "message": "更新钱包失败"})
			return
		}

		// 创建交易记录
		transaction := walletModel.WalletTransaction{
			UserID:      userID.(uint),
			Type:        walletModel.TransactionTypeExchange,
			Amount:      redemptionCode.Amount,
			Balance:     wallet.Balance,
			Description: "兑换码充值",
			RelatedID:   &redemptionCode.ID,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"code": 500, "message": "创建交易记录失败"})
			return
		}

		reward = map[string]interface{}{
			"type":    "balance",
			"amount":  redemptionCode.Amount,
			"balance": wallet.Balance,
		}

	case redemptionModel.RedemptionTypeLevel:
		// 等级提升
		var user userModel.User
		if err := tx.First(&user, userID).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"code": 500, "message": "查询用户失败"})
			return
		}
		user.Level = int(redemptionCode.Amount)
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"code": 500, "message": "更新用户等级失败"})
			return
		}

		reward = map[string]interface{}{
			"type":  "level",
			"level": redemptionCode.Amount,
		}

	case redemptionModel.RedemptionTypeProduct:
		// 产品奖励
		if redemptionCode.ProductID == nil {
			tx.Rollback()
			c.JSON(400, gin.H{"code": 400, "message": "产品类型兑换码缺少产品ID"})
			return
		}
		// TODO: 实现产品奖励逻辑，比如创建产品订单或更新用户产品权限
		reward = map[string]interface{}{
			"type":      "product",
			"productId": redemptionCode.ProductID,
		}
	default:
		tx.Rollback()
		c.JSON(400, gin.H{"code": 400, "message": "不支持的兑换码类型"})
		return
	}

	// 更新兑换码使用次数
	redemptionCode.UsedCount++
	if err := tx.Save(&redemptionCode).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"code": 500, "message": "更新兑换码失败"})
		return
	}

	// 记录使用记录
	usage := redemptionModel.RedemptionCodeUsage{
		CodeID: redemptionCode.ID,
		UserID: userID.(uint),
		Reward:  marshalReward(reward),
	}
	if err := tx.Create(&usage).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"code": 500, "message": "记录使用失败"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		global.APP_LOG.Error("提交事务失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "兑换失败"})
		return
	}

	committed = true

	c.JSON(200, gin.H{
		"code": 200,
		"message": "兑换成功",
		"data": reward,
	})
}

// GetRechargeOrderStatus 获取充值订单状态
// @Summary 获取充值订单状态
// @Description 查询充值订单的支付状态
// @Tags 用户/充值
// @Accept json
// @Produce json
// @Param orderNo path string true "订单号"
// @Success 200 {object} common.Response
// @Router /v1/user/recharge/order-status/{orderNo} [get]
func GetRechargeOrderStatus(c *gin.Context) {
	orderNo := c.Param("orderNo")
	if orderNo == "" {
		c.JSON(400, gin.H{"code": 400, "message": "订单号不能为空"})
		return
	}

	var order orderModel.Order
	if err := global.APP_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "订单不存在"})
		return
	}

	// 检查订单是否过期
	if order.Status == orderModel.OrderStatusPending && time.Now().After(order.ExpireAt) {
		order.Status = orderModel.OrderStatusExpired
		global.APP_DB.Save(&order)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"orderNo":  order.OrderNo,
			"status":   order.Status,
			"amount":   order.Amount,
			"paidTime": order.PaymentTime,
		},
	})
}

// generateOrderNo 生成订单号
func generateOrderNo() string {
	return fmt.Sprintf("%d%06d", time.Now().Unix(), time.Now().Nanosecond()%1000000)
}

// marshalReward 序列化奖励
func marshalReward(reward map[string]interface{}) string {
	if data, err := json.Marshal(reward); err == nil {
		return string(data)
	}
	return "{}"
}
