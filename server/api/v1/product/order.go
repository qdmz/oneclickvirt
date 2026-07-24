package product

import (
	"strconv"

	"oneclickvirt/middleware"
	"oneclickvirt/model/common"
	productModel "oneclickvirt/model/product"
	productService "oneclickvirt/service/product"

	"github.com/gin-gonic/gin"
)

// CreateOrder 创建订单
// @Summary 创建订单
// @Description 用户选择产品、镜像和周期创建订单
// @Tags 产品商城
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body product.CreateOrderRequest true "创建订单请求参数"
// @Success 200 {object} common.Response{data=product.ProductOrder} "创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req productModel.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := productService.NewService()
	order, err := service.CreateOrder(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, order, "订单创建成功")
}

// GetOrderList 获取用户订单列表
// @Summary 获取用户订单列表
// @Description 获取当前登录用户的订单列表
// @Tags 产品商城
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param paymentStatus query int false "支付状态筛选"
// @Param provisionStatus query int false "开通状态筛选"
// @Success 200 {object} common.Response{data=common.PageResult} "获取成功"
// @Failure 401 {object} common.Response "用户未登录"
// @Router /orders [get]
func GetOrderList(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req productModel.OrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	req.Normalize(common.DefaultPageSize)

	service := productService.NewService()
	orders, total, err := service.GetOrderList(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccessWithPagination(c, orders, total, req.Page, req.PageSize)
}

// GetOrderDetail 获取订单详情
// @Summary 获取订单详情
// @Description 获取指定订单的详细信息
// @Tags 产品商城
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} common.Response{data=product.ProductOrder} "获取成功"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 403 {object} common.Response "无权查看"
// @Failure 404 {object} common.Response "订单不存在"
// @Router /orders/{id} [get]
func GetOrderDetail(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的订单ID"))
		return
	}

	service := productService.NewService()
	order, err := service.GetOrderDetail(userID, uint(orderID))
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, order)
}

// PayOrderWithBalance 余额支付
// @Summary 余额支付订单
// @Description 使用账户余额支付指定订单
// @Tags 产品商城
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body product.PayOrderRequest true "支付请求参数"
// @Success 200 {object} common.Response "支付成功"
// @Failure 400 {object} common.Response "参数错误或余额不足"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 404 {object} common.Response "订单不存在"
// @Router /orders/pay [post]
func PayOrderWithBalance(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req productModel.PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := productService.NewService()
	if err := service.PayOrderWithBalance(userID, req.OrderID); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "支付成功")
}

// GetOrderPayStatus 查询支付状态
// @Summary 查询订单支付状态
// @Description 查询指定订单的支付和开通状态
// @Tags 产品商城
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 404 {object} common.Response "订单不存在"
// @Router /orders/{id}/pay-status [get]
func GetOrderPayStatus(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的订单ID"))
		return
	}

	service := productService.NewService()
	order, err := service.GetOrderDetail(userID, uint(orderID))
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	status := map[string]interface{}{
		"orderId":         order.ID,
		"orderNo":         order.OrderNo,
		"paymentStatus":   order.PaymentStatus,
		"provisionStatus": order.ProvisionStatus,
		"instanceId":      order.InstanceID,
		"paidAt":          order.PaidAt,
		"expireAt":        order.ExpireAt,
	}

	common.ResponseSuccess(c, status)
}

// RenewOrder 续费订单
// @Summary 续费订单
// @Description 为已开通的订单创建续费订单
// @Tags 产品商城
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body product.RenewOrderRequest true "续费请求参数"
// @Success 200 {object} common.Response{data=product.ProductOrder} "创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Router /orders/renew [post]
func RenewOrder(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req productModel.RenewOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := productService.NewService()
	order, err := service.RenewOrder(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, order, "续费订单创建成功")
}

// CancelOrder 取消未支付订单
// @Summary 取消未支付订单
// @Description 取消指定未支付订单
// @Tags 产品商城
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body product.CancelOrderRequest true "取消订单请求参数"
// @Success 200 {object} common.Response "取消成功"
// @Failure 400 {object} common.Response "参数错误或订单状态不支持取消"
// @Failure 401 {object} common.Response "用户未登录"
// @Router /orders/cancel [post]
func CancelOrder(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req productModel.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := productService.NewService()
	if err := service.CancelOrder(userID, req.OrderID); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "订单取消成功")
}

// CompleteRenewOrder 完成续费支付（内部调用，也可以由前端在支付成功后调用）
// @Summary 完成续费
// @Description 支付续费订单后调用，延长原订单和实例到期时间
// @Tags 产品商城
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "续费订单ID"
// @Success 200 {object} common.Response "续费成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Router /orders/{id}/complete-renew [post]
func CompleteRenewOrder(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的订单ID"))
		return
	}

	service := productService.NewService()
	order, err := service.GetOrderDetail(userID, uint(orderID))
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	// 检查是否是续费订单
	if !order.IsRenewal {
		common.ResponseWithError(c, common.NewError(common.CodeBadRequest, "不是续费订单"))
		return
	}

	if err := service.CompleteRenewal(order); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "续费成功")
}
