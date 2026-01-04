package admin

import (
	"oneclickvirt/global"
	orderModel "oneclickvirt/model/order"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetOrders 获取所有订单列表
// @Summary 获取所有订单列表
// @Description 管理员获取所有订单列表(支持分页)
// @Tags 管理员/订单管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Param status query string false "订单状态"
// @Param orderNo query string false "订单号"
// @Param username query string false "用户名"
// @Success 200 {object} common.Response
// @Router /v1/admin/orders [get]
func GetOrders(c *gin.Context) {
	global.APP_LOG.Info("开始获取订单列表")

	// 分页参数
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if num, err := strconv.Atoi(p); err == nil && num > 0 {
			page = num
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if num, err := strconv.Atoi(ps); err == nil && num > 0 {
			pageSize = num
		}
	}

	// 构建查询
	query := global.APP_DB.Model(&orderModel.Order{}).Preload("User").Preload("Product")

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 订单号筛选
	if orderNo := c.Query("orderNo"); orderNo != "" {
		query = query.Where("order_no = ?", orderNo)
	}

	// 用户名筛选
	if username := c.Query("username"); username != "" {
		query = query.Joins("JOIN users ON users.id = orders.user_id").Where("users.username LIKE ?", "%"+username+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.APP_LOG.Error("获取订单总数失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取订单列表失败: " + err.Error()})
		return
	}
	global.APP_LOG.Info("订单总数", zap.Int64("total", total))

	// 查询订单列表
	global.APP_LOG.Info("开始查询订单列表",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	// 使用结构体嵌套查询，获取订单和用户名
	type OrderWithUser struct {
		orderModel.Order
		Username string `json:"username"`
	}

	var ordersWithUser []OrderWithUser
	offset := (page - 1) * pageSize

	// 使用JOIN查询获取订单和关联的用户名
	if err := global.APP_DB.Table("orders").Select("orders.*, users.username").
		Joins("LEFT JOIN users ON orders.user_id = users.id").
		Where(query.Session(&gorm.Session{})).
		Order("orders.created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&ordersWithUser).Error; err != nil {
		global.APP_LOG.Error("获取订单列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取订单列表失败: " + err.Error()})
		return
	}

	global.APP_LOG.Info("订单列表查询成功",
		zap.Int("count", len(ordersWithUser)),
		zap.Int64("total", total))

	c.JSON(200, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"list":  ordersWithUser,
			"total": total,
		},
	})
}

// GetOrder 获取订单详情
// @Summary 获取订单详情
// @Description 管理员获取订单详情
// @Tags 管理员/订单管理
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 200 {object} common.Response
// @Router /v1/admin/orders/{id} [get]
func GetOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "订单ID不能为空"})
		return
	}

	// 使用结构体嵌套查询，获取订单和用户名
	type OrderWithUser struct {
		orderModel.Order
		Username string `json:"username"`
	}

	var orderWithUser OrderWithUser
	if err := global.APP_DB.Table("orders").Select("orders.*, users.username").
		Joins("LEFT JOIN users ON orders.user_id = users.id").
		Preload("Product").
		Where("orders.id = ?", id).
		First(&orderWithUser).Error; err != nil {
		global.APP_LOG.Error("获取订单详情失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取订单详情失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "success",
		"data": orderWithUser,
	})
}

// DeleteOrder 删除订单
// @Summary 删除订单
// @Description 管理员删除订单
// @Tags 管理员/订单管理
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 200 {object} common.Response
// @Router /v1/admin/orders/{id} [delete]
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "订单ID不能为空"})
		return
	}

	if err := global.APP_DB.Delete(&orderModel.Order{}, id).Error; err != nil {
		global.APP_LOG.Error("删除订单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "删除订单失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "删除成功",
	})
}

// CancelOrder 取消订单
// @Summary 取消订单
// @Description 管理员取消订单
// @Tags 管理员/订单管理
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 200 {object} common.Response
// @Router /v1/admin/orders/{id}/cancel [post]
func CancelOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "订单ID不能为空"})
		return
	}

	if err := global.APP_DB.Model(&orderModel.Order{}).
		Where("id = ?", id).
		Where("status = ?", "pending").
		Update("status", "cancelled").Error; err != nil {
		global.APP_LOG.Error("取消订单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "取消订单失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "订单已取消",
	})
}

// RefundOrder 退款订单
// @Summary 退款订单
// @Description 管理员退款订单
// @Tags 管理员/订单管理
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Param request body object true "退款信息"
// @Success 200 {object} common.Response
// @Router /v1/admin/orders/{id}/refund [post]
func RefundOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "订单ID不能为空"})
		return
	}

	var req struct {
		Amount int    `json:"amount"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 更新订单状态为已退款
	updates := map[string]interface{}{
		"status": "refunded",
	}
	if req.Reason != "" {
		updates["remark"] = req.Reason
	}

	if err := global.APP_DB.Model(&orderModel.Order{}).
		Where("id = ?", id).
		Where("status = ?", "paid").
		Updates(updates).Error; err != nil {
		global.APP_LOG.Error("退款失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "退款失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "退款成功",
	})
}
