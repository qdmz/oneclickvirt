package product

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"oneclickvirt/constant"
	"oneclickvirt/global"
	adminModel "oneclickvirt/model/admin"
	"oneclickvirt/model/common"
	providerModel "oneclickvirt/model/provider"
	productModel "oneclickvirt/model/product"
	systemModel "oneclickvirt/model/system"
	userModel "oneclickvirt/model/user"
	userService "oneclickvirt/service/user"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Service 产品服务层
type Service struct {
	yiPayService *YiPayService
}

// NewService 创建产品服务实例
func NewService() *Service {
	return &Service{
		yiPayService: NewYiPayService(),
	}
}

// ========== 订单号生成 ==========

// generateOrderNo 生成订单号
// 格式: ORD + 年月日 + 6位随机数
func (s *Service) generateOrderNo() string {
	now := time.Now()
	randomNum := fmt.Sprintf("%06d", randInt(0, 999999))
	return "ORD" + now.Format("20060102") + randomNum
}

// ========== 产品相关方法 ==========

// GetProductList 获取上架产品列表
func (s *Service) GetProductList(req productModel.ProductListRequest) ([]productModel.Product, int64, error) {
	req.Normalize(common.DefaultPageSize)

	var total int64
	var products []productModel.Product

	db := global.APP_DB.Model(&productModel.Product{}).Where("status = ?", 1)
	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	}
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, common.NewError(common.CodeDatabaseError, err.Error())
	}

	if err := db.Order("sort_order DESC, id ASC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&products).Error; err != nil {
		return nil, 0, common.NewError(common.CodeDatabaseError, err.Error())
	}

	return products, total, nil
}

// GetProductDetail 获取产品详情
func (s *Service) GetProductDetail(productID uint) (*productModel.Product, error) {
	var product productModel.Product
	if err := global.APP_DB.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewError(common.CodeNotFound, "产品不存在")
		}
		return nil, common.NewError(common.CodeDatabaseError, err.Error())
	}
	return &product, nil
}

// GetProductImages 获取产品可用镜像
func (s *Service) GetProductImages(productID uint) ([]systemModel.SystemImage, error) {
	product, err := s.GetProductDetail(productID)
	if err != nil {
		return nil, err
	}

	db := global.APP_DB.Model(&systemModel.SystemImage{}).Where("status = ?", "active")

	// 如果产品配置了关联镜像，则只返回关联镜像
	if product.ImageIDs != "" {
		imageIDList := strings.Split(product.ImageIDs, ",")
		var ids []uint
		for _, idStr := range imageIDList {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
				ids = append(ids, uint(id))
			}
		}
		if len(ids) > 0 {
			db = db.Where("id IN ?", ids)
		}
	}

	// 根据产品类型过滤镜像
	if product.Type != "" {
		db = db.Where("provider_type = ?", product.Type)
	}

	var images []systemModel.SystemImage
	if err := db.Find(&images).Error; err != nil {
		return nil, common.NewError(common.CodeDatabaseError, err.Error())
	}

	return images, nil
}

// ========== 订单相关方法 ==========

// CreateOrder 创建订单
func (s *Service) CreateOrder(userID uint, req productModel.CreateOrderRequest) (*productModel.ProductOrder, error) {
	// 获取产品信息
	var product productModel.Product
	if err := global.APP_DB.First(&product, req.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewError(common.CodeNotFound, "产品不存在")
		}
		return nil, common.NewError(common.CodeDatabaseError, err.Error())
	}

	if product.Status != 1 {
		return nil, common.NewError(common.CodeBadRequest, "该产品已下架")
	}

	// 获取镜像信息
	var image systemModel.SystemImage
	if err := global.APP_DB.First(&image, req.ImageID).Error; err != nil {
		return nil, common.NewError(common.CodeNotFound, "镜像不存在")
	}
	if image.Status != "active" {
		return nil, common.NewError(common.CodeBadRequest, "所选镜像不可用")
	}

	// 计算总金额
	totalAmount := product.Price * float64(req.Quantity)

	// 生成订单号
	orderNo := s.generateOrderNo()

	// 计算到期时间
	expireAt := s.calculateExpireAt(product.PeriodType, product.PeriodValue, req.Quantity, time.Now())

	order := productModel.ProductOrder{
		OrderNo:       orderNo,
		UserID:        userID,
		ProductID:     req.ProductID,
		ProductName:   product.Name,
		ProductType:   product.Type,
		CPU:           product.CPU,
		Memory:        product.Memory,
		Disk:          product.Disk,
		Bandwidth:     product.Bandwidth,
		Traffic:       product.Traffic,
		PeriodType:    product.PeriodType,
		PeriodValue:   product.PeriodValue,
		Price:         product.Price,
		Quantity:      req.Quantity,
		TotalAmount:   totalAmount,
		PaymentStatus: 0, // 未支付
		ProvisionStatus: 0, // 待开通
		ImageID:       req.ImageID,
		ImageName:     image.Name,
		ExpireAt:      &expireAt,
	}

	if err := global.APP_DB.Create(&order).Error; err != nil {
		return nil, common.NewError(common.CodeDatabaseError, err.Error())
	}

	global.APP_LOG.Info("创建订单成功",
		zap.Uint("userID", userID),
		zap.String("orderNo", orderNo),
		zap.Float64("totalAmount", totalAmount))

	return &order, nil
}

// GetOrderList 获取用户订单列表
func (s *Service) GetOrderList(userID uint, req productModel.OrderListRequest) ([]productModel.ProductOrder, int64, error) {
	req.Normalize(common.DefaultPageSize)

	var total int64
	var orders []productModel.ProductOrder

	db := global.APP_DB.Model(&productModel.ProductOrder{}).Where("user_id = ?", userID)
	if req.PaymentStatus >= 0 {
		db = db.Where("payment_status = ?", req.PaymentStatus)
	}
	if req.ProvisionStatus >= 0 {
		db = db.Where("provision_status = ?", req.ProvisionStatus)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, common.NewError(common.CodeDatabaseError, err.Error())
	}

	if err := db.Order("created_at DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, common.NewError(common.CodeDatabaseError, err.Error())
	}

	return orders, total, nil
}

// GetOrderDetail 获取订单详情
func (s *Service) GetOrderDetail(userID uint, orderID uint) (*productModel.ProductOrder, error) {
	var order productModel.ProductOrder
	db := global.APP_DB.First(&order, orderID)
	if err := db.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewError(common.CodeNotFound, "订单不存在")
		}
		return nil, common.NewError(common.CodeDatabaseError, err.Error())
	}

	// 用户只能查看自己的订单
	if order.UserID != userID {
		return nil, common.NewError(common.CodeForbidden, "无权查看该订单")
	}

	return &order, nil
}

// GetOrderByOrderNo 根据订单号获取订单
func (s *Service) GetOrderByOrderNo(orderNo string) (*productModel.ProductOrder, error) {
	var order productModel.ProductOrder
	if err := global.APP_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewError(common.CodeNotFound, "订单不存在")
		}
		return nil, common.NewError(common.CodeDatabaseError, err.Error())
	}
	return &order, nil
}

// PayOrderWithBalance 余额支付
func (s *Service) PayOrderWithBalance(userID uint, orderID uint) error {
	var order productModel.ProductOrder
	if err := global.APP_DB.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewError(common.CodeNotFound, "订单不存在")
		}
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	if order.UserID != userID {
		return common.NewError(common.CodeForbidden, "无权操作该订单")
	}

	if order.PaymentStatus != 0 {
		return common.NewError(common.CodeBadRequest, "订单状态不支持支付")
	}

	// 查询用户余额
	var user userModel.User
	if err := global.APP_DB.First(&user, userID).Error; err != nil {
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	if user.Balance < order.TotalAmount {
		return common.NewError(common.CodeBadRequest, "余额不足")
	}

	now := time.Now()

	// 使用事务处理支付和余额扣除
	err := global.APP_DB.Transaction(func(tx *gorm.DB) error {
		// 扣除余额
		balanceBefore := user.Balance
		balanceAfter := user.Balance - order.TotalAmount
		if err := tx.Model(&userModel.User{}).Where("id = ?", userID).
			Update("balance", balanceAfter).Error; err != nil {
			return err
		}

		// 更新订单状态
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"payment_status": 1,
			"payment_method": "balance",
			"paid_at":        &now,
		}).Error; err != nil {
			return err
		}

		// 记录余额变动日志
		log := productModel.UserBalanceLog{
			UserID:        userID,
			Type:          "consume",
			Amount:        -order.TotalAmount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			OrderID:       order.ID,
			Remark:        fmt.Sprintf("订单支付: %s", order.OrderNo),
			TradeNo:       order.OrderNo,
		}
		if err := tx.Create(&log).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	global.APP_LOG.Info("余额支付成功",
		zap.Uint("userID", userID),
		zap.Uint("orderID", orderID),
		zap.String("orderNo", order.OrderNo),
		zap.Float64("amount", order.TotalAmount))

	// 支付成功后自动开通实例
	go s.autoProvisionOrder(order.ID)

	return nil
}

// RenewOrder 续费订单
func (s *Service) RenewOrder(userID uint, req productModel.RenewOrderRequest) (*productModel.ProductOrder, error) {
	var originalOrder productModel.ProductOrder
	if err := global.APP_DB.First(&originalOrder, req.OrderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewError(common.CodeNotFound, "原订单不存在")
		}
		return nil, common.NewError(common.CodeDatabaseError, err.Error())
	}

	if originalOrder.UserID != userID {
		return nil, common.NewError(common.CodeForbidden, "无权操作该订单")
	}

	if originalOrder.PaymentStatus != 1 || originalOrder.ProvisionStatus != 2 {
		return nil, common.NewError(common.CodeBadRequest, "该订单状态不支持续费")
	}

	// 创建续费订单
	totalAmount := originalOrder.Price * float64(req.Quantity)
	orderNo := s.generateOrderNo()

	renewOrder := productModel.ProductOrder{
		OrderNo:         orderNo,
		UserID:          userID,
		ProductID:       originalOrder.ProductID,
		ProductName:     originalOrder.ProductName,
		ProductType:     originalOrder.ProductType,
		CPU:             originalOrder.CPU,
		Memory:          originalOrder.Memory,
		Disk:            originalOrder.Disk,
		Bandwidth:       originalOrder.Bandwidth,
		Traffic:         originalOrder.Traffic,
		PeriodType:      originalOrder.PeriodType,
		PeriodValue:     originalOrder.PeriodValue,
		Price:           originalOrder.Price,
		Quantity:        req.Quantity,
		TotalAmount:     totalAmount,
		PaymentStatus:   0,
		ProvisionStatus: 0,
		IsRenewal:       true,
		RenewOrderID:    originalOrder.ID,
		ImageID:         originalOrder.ImageID,
		ImageName:       originalOrder.ImageName,
	}

	if err := global.APP_DB.Create(&renewOrder).Error; err != nil {
		return nil, common.NewError(common.CodeDatabaseError, err.Error())
	}

	global.APP_LOG.Info("创建续费订单成功",
		zap.Uint("userID", userID),
		zap.Uint("originalOrderID", originalOrder.ID),
		zap.String("orderNo", orderNo))

	return &renewOrder, nil
}

// CompleteRenewal 完成续费（支付成功后调用）
func (s *Service) CompleteRenewal(order *productModel.ProductOrder) error {
	if !order.IsRenewal || order.RenewOrderID == 0 {
		return common.NewError(common.CodeBadRequest, "不是续费订单")
	}

	var originalOrder productModel.ProductOrder
	if err := global.APP_DB.First(&originalOrder, order.RenewOrderID).Error; err != nil {
		return common.NewError(common.CodeNotFound, "原订单不存在")
	}

	// 计算新的到期时间
	var newExpireAt time.Time
	if originalOrder.ExpireAt != nil && originalOrder.ExpireAt.After(time.Now()) {
		// 在原到期时间基础上延长
		newExpireAt = s.calculateExpireAt(originalOrder.PeriodType, originalOrder.PeriodValue, order.Quantity, *originalOrder.ExpireAt)
	} else {
		// 从当前时间开始计算
		newExpireAt = s.calculateExpireAt(originalOrder.PeriodType, originalOrder.PeriodValue, order.Quantity, time.Now())
	}

	// 更新原订单到期时间和关联实例到期时间
	err := global.APP_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&originalOrder).Update("expire_at", &newExpireAt).Error; err != nil {
			return err
		}

		// 更新实例到期时间
		if originalOrder.InstanceID > 0 {
			if err := tx.Model(&providerModel.Instance{}).Where("id = ?", originalOrder.InstanceID).
				Update("expires_at", &newExpireAt).Error; err != nil {
				return err
			}
		}

		// 更新续费订单状态
		now := time.Now()
		if err := tx.Model(order).Updates(map[string]interface{}{
			"payment_status":   1,
			"payment_method":   "balance",
			"paid_at":          &now,
			"provision_status": 2,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	global.APP_LOG.Info("续费成功",
		zap.Uint("orderID", order.ID),
		zap.Uint("originalOrderID", originalOrder.ID),
		zap.Time("newExpireAt", newExpireAt))

	return nil
}

// CancelOrder 取消未支付订单
func (s *Service) CancelOrder(userID uint, orderID uint) error {
	var order productModel.ProductOrder
	if err := global.APP_DB.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewError(common.CodeNotFound, "订单不存在")
		}
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	if order.UserID != userID {
		return common.NewError(common.CodeForbidden, "无权操作该订单")
	}

	if order.PaymentStatus != 0 {
		return common.NewError(common.CodeBadRequest, "只有未支付订单可以取消")
	}

	if err := global.APP_DB.Delete(&order).Error; err != nil {
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	return nil
}

// ========== 支付成功后的自动开通逻辑 ==========

// autoProvisionOrder 自动开通订单对应的实例
func (s *Service) autoProvisionOrder(orderID uint) {
	// 延迟执行，给数据库事务足够时间提交
	time.Sleep(1 * time.Second)

	var order productModel.ProductOrder
	if err := global.APP_DB.First(&order, orderID).Error; err != nil {
		global.APP_LOG.Error("自动开通: 订单不存在", zap.Uint("orderID", orderID), zap.Error(err))
		return
	}

	if order.PaymentStatus != 1 {
		global.APP_LOG.Warn("自动开通: 订单未支付", zap.Uint("orderID", orderID))
		return
	}

	// 更新为开通中
	if err := global.APP_DB.Model(&order).Update("provision_status", 1).Error; err != nil {
		global.APP_LOG.Error("自动开通: 更新状态失败", zap.Uint("orderID", orderID), zap.Error(err))
		return
	}

	// 执行开通
	if err := s.provisionInstance(&order); err != nil {
		global.APP_LOG.Error("自动开通失败",
			zap.Uint("orderID", orderID),
			zap.String("orderNo", order.OrderNo),
			zap.Error(err))
		global.APP_DB.Model(&order).Update("provision_status", 3) // 开通失败
		return
	}
}

// provisionInstance 为订单开通实例
func (s *Service) provisionInstance(order *productModel.ProductOrder) error {
	// 获取产品信息以确定可用节点
	var product productModel.Product
	if err := global.APP_DB.First(&product, order.ProductID).Error; err != nil {
		return fmt.Errorf("获取产品信息失败: %w", err)
	}

	// 选择节点
	providerID, err := s.selectProvider(&product)
	if err != nil {
		return err
	}

	// 构建规格ID
	cpuID := fmt.Sprintf("cpu-%d", order.CPU)
	memoryID := fmt.Sprintf("mem-%dmb", order.Memory)
	diskID := fmt.Sprintf("disk-%dmb", order.Disk)
	bandwidthID := fmt.Sprintf("bw-%dmbps", order.Bandwidth)

	// 验证规格ID是否有效
	if _, err := constant.GetCPUSpecByID(cpuID); err != nil {
		return fmt.Errorf("CPU规格无效: %w", err)
	}
	if _, err := constant.GetMemorySpecByID(memoryID); err != nil {
		return fmt.Errorf("内存规格无效: %w", err)
	}
	if _, err := constant.GetDiskSpecByID(diskID); err != nil {
		return fmt.Errorf("磁盘规格无效: %w", err)
	}
	if _, err := constant.GetBandwidthSpecByID(bandwidthID); err != nil {
		return fmt.Errorf("带宽规格无效: %w", err)
	}

	// 构建创建实例请求
	createReq := userModel.CreateInstanceRequest{
		ProviderId:  providerID,
		ImageId:     order.ImageID,
		CPUId:       cpuID,
		MemoryId:    memoryID,
		DiskId:      diskID,
		BandwidthId: bandwidthID,
		Description: fmt.Sprintf("产品订单开通: %s", order.OrderNo),
		OrderID:     order.ID,
	}

	global.APP_LOG.Info("开始为订单开通实例",
		zap.Uint("orderID", order.ID),
		zap.String("orderNo", order.OrderNo),
		zap.Uint("providerID", providerID),
		zap.Uint("imageID", order.ImageID))

	// 调用已有实例创建接口
	us := userService.NewService()
	task, err := us.CreateUserInstance(order.UserID, createReq)
	if err != nil {
		return fmt.Errorf("创建实例任务失败: %w", err)
	}

	global.APP_LOG.Info("实例创建任务已提交",
		zap.Uint("orderID", order.ID),
		zap.Uint("taskID", task.ID))

	// 后台跟踪任务状态
	go s.trackProvisionTask(order.ID, task.ID)

	return nil
}

// trackProvisionTask 跟踪开通任务状态
func (s *Service) trackProvisionTask(orderID uint, taskID uint) {
	maxRetries := 120 // 最多跟踪10分钟（每5秒检查一次）
	interval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		time.Sleep(interval)

		var task adminModel.Task
		if err := global.APP_DB.First(&task, taskID).Error; err != nil {
			global.APP_LOG.Error("跟踪任务: 查询任务失败",
				zap.Uint("taskID", taskID),
				zap.Error(err))
			continue
		}

		if task.Status == "completed" {
			now := time.Now()
			updates := map[string]interface{}{
				"provision_status": 2,
				"provisioned_at":   &now,
			}
			if task.InstanceID != nil && *task.InstanceID > 0 {
				updates["instance_id"] = *task.InstanceID

				// 同时更新实例到期时间
				var order productModel.ProductOrder
				if err := global.APP_DB.First(&order, orderID).Error; err == nil && order.ExpireAt != nil {
					global.APP_DB.Model(&providerModel.Instance{}).Where("id = ?", *task.InstanceID).
						Update("expires_at", order.ExpireAt)
				}
			}

			if err := global.APP_DB.Model(&productModel.ProductOrder{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
				global.APP_LOG.Error("跟踪任务: 更新订单状态失败",
					zap.Uint("orderID", orderID),
					zap.Error(err))
			} else {
				global.APP_LOG.Info("实例开通成功",
					zap.Uint("orderID", orderID),
					zap.Uint("taskID", taskID),
					zap.Uintp("instanceID", task.InstanceID))
			}
			return
		}

		if task.Status == "failed" || task.Status == "timeout" {
			if err := global.APP_DB.Model(&productModel.ProductOrder{}).Where("id = ?", orderID).
				Update("provision_status", 3).Error; err != nil {
				global.APP_LOG.Error("跟踪任务: 更新订单失败状态失败",
					zap.Uint("orderID", orderID),
					zap.Error(err))
			}
			global.APP_LOG.Warn("实例开通失败",
				zap.Uint("orderID", orderID),
				zap.Uint("taskID", taskID),
				zap.String("taskStatus", task.Status),
				zap.String("errorMessage", task.ErrorMessage))
			return
		}
	}

	global.APP_LOG.Warn("跟踪任务超时",
		zap.Uint("orderID", orderID),
		zap.Uint("taskID", taskID))
}

// selectProvider 从产品配置的节点中选择一个可用节点
func (s *Service) selectProvider(product *productModel.Product) (uint, error) {
	if product.ProviderIDs == "" {
		return 0, errors.New("该产品未配置可用节点")
	}

	providerIDList := strings.Split(product.ProviderIDs, ",")
	for _, idStr := range providerIDList {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		providerID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			continue
		}

		// 检查节点是否可用
		var provider providerModel.Provider
		if err := global.APP_DB.First(&provider, providerID).Error; err != nil {
			continue
		}

		providerAvailable := (provider.ConnectionType == "agent" && provider.AgentStatus == "online") ||
			(provider.ConnectionType != "agent" && (provider.Status == "active" || provider.Status == "partial"))

		if providerAvailable && !provider.IsFrozen && provider.AllowClaim {
			return uint(providerID), nil
		}
	}

	return 0, errors.New("暂无可用的节点，请联系管理员")
}

// calculateExpireAt 计算到期时间
func (s *Service) calculateExpireAt(periodType string, periodValue, quantity int, start time.Time) time.Time {
	duration := time.Duration(periodValue * quantity)
	switch periodType {
	case "hour":
		return start.Add(duration * time.Hour)
	case "day":
		return start.Add(duration * 24 * time.Hour)
	case "month":
		return start.AddDate(0, int(duration), 0)
	case "year":
		return start.AddDate(int(duration), 0, 0)
	default:
		return start.AddDate(0, int(duration), 0)
	}
}

// ========== 易支付充值相关方法 ==========

// CreateYiPayOrder 创建易支付充值订单
func (s *Service) CreateYiPayOrder(userID uint, req productModel.CreateYiPayOrderRequest, clientIP string) (map[string]interface{}, error) {
	config, err := s.yiPayService.GetActiveConfig()
	if err != nil {
		return nil, err
	}

	if req.Amount < config.MinAmount || req.Amount > config.MaxAmount {
		return nil, common.NewError(common.CodeBadRequest,
			fmt.Sprintf("充值金额必须在 %.2f 到 %.2f 之间", config.MinAmount, config.MaxAmount))
	}

	// 生成充值订单号
	rechargeNo := s.yiPayService.GenerateOrderNo()

	// 创建余额变动记录（初始状态为未支付）
	log := productModel.UserBalanceLog{
		UserID:  userID,
		Type:    "recharge",
		Amount:  req.Amount,
		TradeNo: rechargeNo,
		Remark:  fmt.Sprintf("易支付充值: %s", config.Name),
	}
	if err := global.APP_DB.Create(&log).Error; err != nil {
		return nil, common.NewError(common.CodeDatabaseError, err.Error())
	}

	// 构建支付参数
	params := s.yiPayService.BuildPayParams(config, rechargeNo, req.Amount, req.PayType, clientIP)
	payURL := s.yiPayService.GeneratePayURL(config, params)

	return map[string]interface{}{
		"rechargeNo": rechargeNo,
		"payURL":     payURL,
		"amount":     req.Amount,
	}, nil
}

// ProcessYiPayNotify 处理易支付异步通知
func (s *Service) ProcessYiPayNotify(params map[string]string) error {
	config, err := s.yiPayService.GetActiveConfig()
	if err != nil {
		return err
	}

	// 验证签名
	if !s.yiPayService.VerifyNotify(params, config) {
		return common.NewError(common.CodeBadRequest, "签名验证失败")
	}

	// 检查交易状态
	if params["trade_status"] != "TRADE_SUCCESS" {
		global.APP_LOG.Warn("易支付通知: 交易未成功",
			zap.String("tradeStatus", params["trade_status"]))
		return nil
	}

	rechargeNo := params["out_trade_no"]
	money, _ := strconv.ParseFloat(params["money"], 64)
	tradeNo := params["trade_no"]

	// 查询充值记录
	var log productModel.UserBalanceLog
	if err := global.APP_DB.Where("trade_no = ? AND type = ?", rechargeNo, "recharge").First(&log).Error; err != nil {
		return common.NewError(common.CodeNotFound, "充值记录不存在")
	}

	// 已处理则直接返回
	if log.BalanceAfter > 0 || log.ID == 0 {
		// 检查是否已经处理过（通过查询用户余额是否有变动判断）
		var existingLog productModel.UserBalanceLog
		if err := global.APP_DB.Where("trade_no = ? AND balance_after > balance_before", rechargeNo).First(&existingLog).Error; err == nil {
			global.APP_LOG.Info("易支付通知: 该订单已处理", zap.String("rechargeNo", rechargeNo))
			return nil
		}
	}

	// 验证金额
	if money <= 0 {
		return common.NewError(common.CodeBadRequest, "无效的支付金额")
	}

	// 使用事务处理充值
	err = global.APP_DB.Transaction(func(tx *gorm.DB) error {
		// 查询用户当前余额
		var user userModel.User
		if err := tx.First(&user, log.UserID).Error; err != nil {
			return err
		}

		balanceBefore := user.Balance
		balanceAfter := user.Balance + money

		// 增加余额
		if err := tx.Model(&userModel.User{}).Where("id = ?", log.UserID).
			Update("balance", balanceAfter).Error; err != nil {
			return err
		}

		// 更新充值记录
		if err := tx.Model(&log).Updates(map[string]interface{}{
			"amount":         money,
			"balance_before": balanceBefore,
			"balance_after":  balanceAfter,
			"remark":         fmt.Sprintf("易支付充值成功: %s", tradeNo),
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	global.APP_LOG.Info("易支付充值成功",
		zap.String("rechargeNo", rechargeNo),
		zap.String("tradeNo", tradeNo),
		zap.Float64("amount", money),
		zap.Uint("userID", log.UserID))

	return nil
}

// GetRechargeList 获取用户充值记录
func (s *Service) GetRechargeList(userID uint, page, pageSize int) ([]productModel.UserBalanceLog, int64, error) {
	page, pageSize = common.NormalizePagination(page, pageSize, common.DefaultPageSize)

	var total int64
	var logs []productModel.UserBalanceLog

	db := global.APP_DB.Model(&productModel.UserBalanceLog{}).
		Where("user_id = ? AND type = ?", userID, "recharge")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, common.NewError(common.CodeDatabaseError, err.Error())
	}

	if err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		return nil, 0, common.NewError(common.CodeDatabaseError, err.Error())
	}

	return logs, total, nil
}

// ========== 管理员相关方法 ==========

// GetAdminProductList 获取管理员产品列表
func (s *Service) GetAdminProductList(req productModel.AdminProductListRequest) ([]productModel.Product, int64, error) {
	req.Normalize(common.DefaultPageSize)

	var total int64
	var products []productModel.Product

	db := global.APP_DB.Model(&productModel.Product{})
	if req.Status >= 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	}
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}
	if req.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, common.NewError(common.CodeDatabaseError, err.Error())
	}

	if err := db.Order("sort_order DESC, id ASC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&products).Error; err != nil {
		return nil, 0, common.NewError(common.CodeDatabaseError, err.Error())
	}

	return products, total, nil
}

// CreateAdminProduct 管理员创建产品
func (s *Service) CreateAdminProduct(req productModel.CreateProductRequest) (*productModel.Product, error) {
	product := productModel.Product{
		Name:         req.Name,
		Description:  req.Description,
		Type:         req.Type,
		Category:     req.Category,
		CPU:          req.CPU,
		Memory:       req.Memory,
		Disk:         req.Disk,
		Bandwidth:    req.Bandwidth,
		Traffic:      req.Traffic,
		Price:        req.Price,
		PeriodType:   req.PeriodType,
		PeriodValue:  req.PeriodValue,
		MaxSnapshots: req.MaxSnapshots,
		MaxPorts:     req.MaxPorts,
		Status:       req.Status,
		SortOrder:    req.SortOrder,
		Icon:         req.Icon,
		ImageIDs:     req.ImageIDs,
		ProviderIDs:  req.ProviderIDs,
	}

	if err := global.APP_DB.Create(&product).Error; err != nil {
		return nil, common.NewError(common.CodeDatabaseError, err.Error())
	}

	return &product, nil
}

// UpdateAdminProduct 管理员更新产品
func (s *Service) UpdateAdminProduct(productID uint, req productModel.UpdateProductRequest) error {
	var product productModel.Product
	if err := global.APP_DB.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewError(common.CodeNotFound, "产品不存在")
		}
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	updates := map[string]interface{}{
		"name":          req.Name,
		"description":   req.Description,
		"type":          req.Type,
		"category":      req.Category,
		"cpu":           req.CPU,
		"memory":        req.Memory,
		"disk":          req.Disk,
		"bandwidth":     req.Bandwidth,
		"traffic":       req.Traffic,
		"price":         req.Price,
		"period_type":   req.PeriodType,
		"period_value":  req.PeriodValue,
		"max_snapshots": req.MaxSnapshots,
		"max_ports":     req.MaxPorts,
		"status":        req.Status,
		"sort_order":    req.SortOrder,
		"icon":          req.Icon,
		"image_ids":     req.ImageIDs,
		"provider_ids":  req.ProviderIDs,
	}

	if err := global.APP_DB.Model(&product).Updates(updates).Error; err != nil {
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	return nil
}

// DeleteAdminProduct 管理员删除产品
func (s *Service) DeleteAdminProduct(productID uint) error {
	var product productModel.Product
	if err := global.APP_DB.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewError(common.CodeNotFound, "产品不存在")
		}
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	if err := global.APP_DB.Delete(&product).Error; err != nil {
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	return nil
}

// GetAdminOrderList 获取管理员订单列表
func (s *Service) GetAdminOrderList(req productModel.AdminOrderListRequest) ([]productModel.ProductOrder, int64, error) {
	req.Normalize(common.DefaultPageSize)

	var total int64
	var orders []productModel.ProductOrder

	db := global.APP_DB.Model(&productModel.ProductOrder{})
	if req.UserID > 0 {
		db = db.Where("user_id = ?", req.UserID)
	}
	if req.PaymentStatus >= 0 {
		db = db.Where("payment_status = ?", req.PaymentStatus)
	}
	if req.ProvisionStatus >= 0 {
		db = db.Where("provision_status = ?", req.ProvisionStatus)
	}
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.Keyword != "" {
		db = db.Where("product_name LIKE ? OR order_no LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, common.NewError(common.CodeDatabaseError, err.Error())
	}

	if err := db.Order("created_at DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, common.NewError(common.CodeDatabaseError, err.Error())
	}

	return orders, total, nil
}

// GetAdminOrderDetail 获取管理员订单详情
func (s *Service) GetAdminOrderDetail(orderID uint) (*productModel.ProductOrder, error) {
	var order productModel.ProductOrder
	if err := global.APP_DB.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewError(common.CodeNotFound, "订单不存在")
		}
		return nil, common.NewError(common.CodeDatabaseError, err.Error())
	}
	return &order, nil
}

// UpdateOrderStatus 管理员更新订单状态
func (s *Service) UpdateOrderStatus(orderID uint, req productModel.UpdateOrderStatusRequest) error {
	var order productModel.ProductOrder
	if err := global.APP_DB.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewError(common.CodeNotFound, "订单不存在")
		}
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	updates := map[string]interface{}{}
	if req.PaymentStatus >= 0 {
		updates["payment_status"] = req.PaymentStatus
	}
	if req.ProvisionStatus >= 0 {
		updates["provision_status"] = req.ProvisionStatus
	}
	if req.Remark != "" {
		updates["remark"] = gorm.Expr("CONCAT(COALESCE(remark, ''), '\n管理员备注: ', ?)", req.Remark)
	}

	if len(updates) == 0 {
		return common.NewError(common.CodeBadRequest, "没有需要更新的字段")
	}

	if err := global.APP_DB.Model(&order).Updates(updates).Error; err != nil {
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	return nil
}

// ManualProvision 管理员手动开通实例
func (s *Service) ManualProvision(orderID uint) error {
	var order productModel.ProductOrder
	if err := global.APP_DB.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewError(common.CodeNotFound, "订单不存在")
		}
		return common.NewError(common.CodeDatabaseError, err.Error())
	}

	if order.PaymentStatus != 1 {
		return common.NewError(common.CodeBadRequest, "订单未支付")
	}

	if order.ProvisionStatus == 2 {
		return common.NewError(common.CodeBadRequest, "实例已开通")
	}

	if err := s.provisionInstance(&order); err != nil {
		return common.NewError(common.CodeInternalError, err.Error())
	}

	return nil
}
