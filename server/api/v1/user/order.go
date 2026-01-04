package user

import (
	"encoding/json"
	"fmt"
	"oneclickvirt/global"
	orderModel "oneclickvirt/model/order"
	productModel "oneclickvirt/model/product"
	instanceModel "oneclickvirt/model/provider"
	userModel "oneclickvirt/model/user"
	walletModel "oneclickvirt/model/wallet"
	"oneclickvirt/service/cache"

	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// fillProductResourcesFromLevelLimit 从用户等级限制中自动填充产品资源配置
func fillProductResourcesFromLevelLimit(product *productModel.Product) {
	// 获取全局配置中的等级限制
	levelLimits := global.APP_CONFIG.Quota.LevelLimits
	level := product.Level

	if levelInfo, exists := levelLimits[level]; exists {
		// 填充最大实例数
		product.MaxInstances = levelInfo.MaxInstances
		// 填充最大流量（直接使用int64类型，无需转换）
		product.Traffic = levelInfo.MaxTraffic
		// 填充资源配置
		maxResources := levelInfo.MaxResources
		if maxResources != nil {
			// 填充CPU核心数
			if cpu, ok := maxResources["cpu"].(float64); ok {
				product.CPU = int(cpu)
			} else if cpu, ok := maxResources["cpu"].(int); ok {
				product.CPU = cpu
			}

			// 填充内存(MB)
			if memory, ok := maxResources["memory"].(float64); ok {
				product.Memory = int(memory)
			} else if memory, ok := maxResources["memory"].(int); ok {
				product.Memory = memory
			}

			// 填充磁盘(MB)
			if disk, ok := maxResources["disk"].(float64); ok {
				product.Disk = int(disk)
			} else if disk, ok := maxResources["disk"].(int); ok {
				product.Disk = disk
			}

			// 填充带宽(Mbps)
			if bandwidth, ok := maxResources["bandwidth"].(float64); ok {
				product.Bandwidth = int(bandwidth)
			} else if bandwidth, ok := maxResources["bandwidth"].(int); ok {
				product.Bandwidth = bandwidth
			}
		}

		global.APP_LOG.Info(fmt.Sprintf("产品 %s (等级 %d) 资源配置已从等级限制自动填充", product.Name, level))
	} else {
		global.APP_LOG.Warn(fmt.Sprintf("未找到等级 %d 的限制配置，无法自动填充产品资源配置", level))
		// 记录可用的等级限制键
		availableLevels := make([]int, 0, len(levelLimits))
		for k := range levelLimits {
			availableLevels = append(availableLevels, k)
		}
		global.APP_LOG.Warn(fmt.Sprintf("可用的等级限制键: %v", availableLevels))
	}
}

// GetUserOrders 获取我的订单列表
// @Summary 获取我的订单列表
// @Description 获取当前登录用户的订单列表
// @Tags 用户/订单
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param status query string false "订单状态"
// @Success 200 {object} common.Response
// @Router /v1/user/orders [get]
func GetUserOrders(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未授权"})
		return
	}

	// 分页参数
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if parsed, err := parseUint(p); err == nil && parsed > 0 {
			page = int(parsed)
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if parsed, err := parseUint(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = int(parsed)
		}
	}

	var orders []orderModel.Order
	var total int64

	query := global.APP_DB.Model(&orderModel.Order{}).Where("user_id = ?", userID)

	// 订单状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 查询总数
	query.Count(&total)

	// 查询列表
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Preload("Product").
		Limit(pageSize).
		Offset(offset).
		Find(&orders).Error; err != nil {
		global.APP_LOG.Error("获取订单列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取订单列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":     orders,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetUserOrder 获取订单详情
// @Summary 获取订单详情
// @Description 获取当前登录用户的订单详情
// @Tags 用户/订单
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 200 {object} common.Response
// @Router /v1/user/orders/{id} [get]
func GetUserOrder(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未授权"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "订单ID不能为空"})
		return
	}

	var order orderModel.Order
	if err := global.APP_DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("Product").
		First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"code": 404, "message": "订单不存在"})
			return
		}
		global.APP_LOG.Error("获取订单详情失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取订单详情失败"})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    order,
	})
}

// CancelOrder 取消订单
// @Summary 取消订单
// @Description 取消待支付的订单
// @Tags 用户/订单
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 200 {object} common.Response
// @Router /v1/user/orders/{id}/cancel [post]
func CancelOrder(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未授权"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "订单ID不能为空"})
		return
	}

	var order orderModel.Order
	if err := global.APP_DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"code": 404, "message": "订单不存在"})
			return
		}
		global.APP_LOG.Error("查询订单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "查询订单失败"})
		return
	}

	// 检查订单状态
	if order.Status != orderModel.OrderStatusPending {
		c.JSON(400, gin.H{"code": 400, "message": "只能取消待支付订单"})
		return
	}

	// 更新订单状态
	order.Status = orderModel.OrderStatusCancelled
	if err := global.APP_DB.Save(&order).Error; err != nil {
		global.APP_LOG.Error("取消订单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "取消订单失败"})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "订单已取消",
	})
}

// GetUserProducts 获取可用产品列表
// @Summary 获取可用产品列表
// @Description 获取当前可用购买的产品列表
// @Tags 用户/订单
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /v1/user/products [get]
func GetUserProducts(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未授权"})
		return
	}

	// 获取所有启用的产品
	var products []productModel.Product
	if err := global.APP_DB.Where("is_enabled = ?", 1).
		Order("sort_order ASC, id ASC").
		Find(&products).Error; err != nil {
		global.APP_LOG.Error("获取产品列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取产品列表失败"})
		return
	}

	// 从用户等级限制中自动填充产品资源配置，确保产品配置与最新的等级配额一致
	for i := range products {
		fillProductResourcesFromLevelLimit(&products[i])
	}

	// 查询用户已经购买过的产品ID列表
	var purchasedProductIDs []uint
	if err := global.APP_DB.Model(&productModel.ProductPurchase{}).
		Where("user_id = ?", userID).
		Pluck("product_id", &purchasedProductIDs).Error; err != nil {
		global.APP_LOG.Error("查询用户购买记录失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取产品列表失败"})
		return
	}

	// 将购买记录转换为map，方便查询
	purchasedMap := make(map[uint]bool)
	for _, productID := range purchasedProductIDs {
		purchasedMap[productID] = true
	}

	// 过滤产品列表，只返回用户可以购买的产品
	var availableProducts []productModel.Product
	for _, product := range products {
		// 如果产品允许重复购买，或者用户没有购买过该产品，则显示该产品
		if product.AllowRepeat == 1 || !purchasedMap[product.ID] {
			availableProducts = append(availableProducts, product)
		}
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    availableProducts,
	})
}

// PurchaseProduct 购买产品
// @Summary 购买产品
// @Description 用户购买产品
// @Tags 用户/订单
// @Accept json
// @Produce json
// @Param id path uint true "产品ID"
// @Param data body map[string]interface{} true "购买信息: paymentMethod(支付方式)"
// @Success 200 {object} common.Response
// @Router /v1/user/products/{id}/purchase [post]
func PurchaseProduct(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未授权"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "产品ID不能为空"})
		return
	}

	var params struct {
		PaymentMethod string `json:"paymentMethod" binding:"required"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	tx := global.APP_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取产品信息
	var product productModel.Product
	if err := tx.First(&product, id).Error; err != nil {
		tx.Rollback()
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	// 检查产品是否启用
	if product.IsEnabled != 1 {
		tx.Rollback()
		c.JSON(400, gin.H{"code": 400, "message": "产品已下架"})
		return
	}

	// 检查是否允许重复购买
	if product.AllowRepeat == 0 {
		// 查询用户是否已经购买过该产品
		var purchaseCount int64
		if err := tx.Model(&productModel.ProductPurchase{}).
			Where("user_id = ? AND product_id = ?", userID, product.ID).
			Count(&purchaseCount).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"code": 500, "message": "查询购买记录失败"})
			return
		}

		if purchaseCount > 0 {
			tx.Rollback()
			c.JSON(400, gin.H{"code": 400, "message": "该产品不允许重复购买，您已经购买过该产品"})
			return
		}
	}

	// 获取用户信息
	var user userModel.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"code": 500, "message": "查询用户失败"})
		return
	}

	// 检查用户等级：只能购买等于或大于自己当前等级的产品
	if user.Level > product.Level {
		tx.Rollback()
		c.JSON(400, gin.H{"code": 400, "message": "此产品不符合该用户等级购买请返回重新选择其它等级产品"})
		return
	}

	// 生成订单号
	orderNo := generateOrderNo()

	// 序列化产品数据
	productData := marshalProductData(product)

	// 计算有效期
	var expireTime *time.Time
	if product.Period > 0 {
		t := time.Now().AddDate(0, 0, product.Period)
		expireTime = &t
	}

	// 创建订单
	order := orderModel.Order{
		OrderNo:       orderNo,
		UserID:        userID.(uint),
		ProductID:     &product.ID,
		Amount:        product.Price,
		Status:        orderModel.OrderStatusPending,
		PaymentMethod: params.PaymentMethod,
		ProductData:   productData,
		ExpireAt:      time.Now().Add(30 * time.Minute), // 30分钟过期
	}

	// 如果是余额支付,立即处理
	if params.PaymentMethod == orderModel.PaymentMethodBalance {
		// 获取用户钱包
		var wallet walletModel.UserWallet
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			// 钱包不存在则创建
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
		}

		// 检查余额是否足够
		if wallet.Balance < product.Price {
			tx.Rollback()
			c.JSON(400, gin.H{"code": 400, "message": "余额不足"})
			return
		}

		// 扣除余额
		wallet.Balance -= product.Price
		wallet.TotalExpense += product.Price
		if err := tx.Save(&wallet).Error; err != nil {
			tx.Rollback()
			global.APP_LOG.Error("更新钱包失败", zap.Error(err))
			c.JSON(500, gin.H{"code": 500, "message": "支付失败"})
			return
		}

		// 创建交易记录
		transaction := walletModel.WalletTransaction{
			UserID:      userID.(uint),
			Type:        walletModel.TransactionTypeConsume,
			Amount:      -product.Price,
			Balance:     wallet.Balance,
			Description: "购买产品: " + product.Name,
			OrderID:     &order.ID,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"code": 500, "message": "创建交易记录失败"})
			return
		}

		// 更新订单状态
		order.Status = orderModel.OrderStatusPaid
		now := time.Now()
		order.PaymentTime = &now
		order.PaidAmount = product.Price
	}

	// 创建订单
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		global.APP_LOG.Error("创建订单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "创建订单失败"})
		return
	}

	// 如果已支付,升级用户等级并更新资源配额
	if order.Status == orderModel.OrderStatusPaid {
		if err := upgradeUserLevel(tx, userID.(uint), product.Level, expireTime, product); err != nil {
			tx.Rollback()
			global.APP_LOG.Error("升级用户等级失败", zap.Error(err))
			c.JSON(500, gin.H{"code": 500, "message": "升级用户等级失败"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		global.APP_LOG.Error("提交事务失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "购买失败"})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "购买成功",
		"data":    order,
	})
}

// upgradeUserLevel 升级用户等级
func upgradeUserLevel(tx *gorm.DB, userID uint, newLevel int, expireTime *time.Time, product productModel.Product) error {
	var user userModel.User
	if err := tx.First(&user, userID).Error; err != nil {
		return err
	}

	// 只更新用户等级为所购产品等级
	user.Level = newLevel

	// 更新有效期：当前到期日期+产品有效期
	if expireTime != nil {
		now := time.Now()
		if user.LevelExpireAt != nil && user.LevelExpireAt.After(now) {
			// 如果用户已有有效的到期时间，在其基础上延长
			newExpire := user.LevelExpireAt.AddDate(0, 0, product.Period)
			user.LevelExpireAt = &newExpire
			global.APP_LOG.Info(fmt.Sprintf("用户已有有效到期时间，在基础上延长: %v + %d天 = %v", user.LevelExpireAt, product.Period, newExpire))
		} else {
			// 如果没有有效的到期时间，从当前时间开始计算
			user.LevelExpireAt = expireTime
			global.APP_LOG.Info(fmt.Sprintf("用户没有有效到期时间，从当前时间开始计算: %v", user.LevelExpireAt))
		}
	} else {
		// 永久产品,设置为9999年后
		farFuture := time.Now().AddDate(9999, 0, 0)
		user.LevelExpireAt = &farFuture
		global.APP_LOG.Info(fmt.Sprintf("永久产品，设置到期时间为: %v", user.LevelExpireAt))
	}

	global.APP_LOG.Info(fmt.Sprintf("用户 %d 升级到等级 %d，到期时间: %v", userID, newLevel, user.LevelExpireAt))

	// 更新用户名下所有实例的到期时间
	if err := updateUserInstancesExpireTime(tx, userID, user.LevelExpireAt); err != nil {
		return err
	}

	// 保存用户信息
	if err := tx.Save(&user).Error; err != nil {
		return err
	}

	// 创建产品购买记录，用于记录用户购买的产品
	now := time.Now()
	productPurchase := productModel.ProductPurchase{
		UserID:    userID,
		ProductID: product.ID,
		Level:     newLevel,
		StartDate: now,
		EndDate:   user.LevelExpireAt,
		IsActive:  true,
	}
	if err := tx.Create(&productPurchase).Error; err != nil {
		global.APP_LOG.Error("创建产品购买记录失败", zap.Error(err))
		return err
	}

	// 清除用户Dashboard缓存，确保用户下次访问个人中心时能看到最新的用户信息
	cacheService := cache.GetUserCacheService()
	cacheService.InvalidateUserCache(userID)
	global.APP_LOG.Info(fmt.Sprintf("已清除用户 %d 的Dashboard缓存", userID))

	return nil
}

// updateUserInstancesExpireTime 更新用户名下所有实例的到期时间
func updateUserInstancesExpireTime(tx *gorm.DB, userID uint, expireAt *time.Time) error {
	// 更新用户名下所有实例的到期时间
	return tx.Model(&instanceModel.Instance{}).Where("user_id = ?", userID).Update("expired_at", expireAt).Error
}

// GetPurchaseOrderStatus 获取购买订单状态
// @Summary 获取购买订单状态
// @Description 查询购买产品订单的支付状态
// @Tags 用户/订单
// @Accept json
// @Produce json
// @Param orderNo path string true "订单号"
// @Success 200 {object} common.Response
// @Router /v1/user/orders/status/{orderNo} [get]
func GetPurchaseOrderStatus(c *gin.Context) {
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
		"code":    200,
		"message": "success",
		"data": gin.H{
			"orderNo":  order.OrderNo,
			"status":   order.Status,
			"amount":   order.Amount,
			"paidTime": order.PaymentTime,
		},
	})
}

// marshalProductData 序列化产品数据
func marshalProductData(product productModel.Product) string {
	data := map[string]interface{}{
		"id":           product.ID,
		"name":         product.Name,
		"description":  product.Description,
		"level":        product.Level,
		"price":        product.Price,
		"period":       product.Period,
		"cpu":          product.CPU,
		"memory":       product.Memory,
		"disk":         product.Disk,
		"bandwidth":    product.Bandwidth,
		"traffic":      product.Traffic,
		"maxInstances": product.MaxInstances,
	}

	if dataBytes, err := json.Marshal(data); err == nil {
		return string(dataBytes)
	}
	return "{}"
}
