package admin

import (
	"strconv"

	"oneclickvirt/model/common"
	productModel "oneclickvirt/model/product"
	productService "oneclickvirt/service/product"

	"github.com/gin-gonic/gin"
)

// GetAdminOrderList 获取管理员订单列表
// @Summary 获取管理员订单列表
// @Description 获取所有用户的订单列表
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param userId query int false "用户ID筛选"
// @Param paymentStatus query int false "支付状态筛选"
// @Param provisionStatus query int false "开通状态筛选"
// @Param orderNo query string false "订单号搜索"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} common.Response{data=common.PageResult} "获取成功"
// @Router /admin/orders [get]
func GetAdminOrderList(c *gin.Context) {
	var req productModel.AdminOrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	req.Normalize(common.DefaultPageSize)

	service := productService.NewService()
	orders, total, err := service.GetAdminOrderList(req)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccessWithPagination(c, orders, total, req.Page, req.PageSize)
}

// GetAdminOrderDetail 获取管理员订单详情
// @Summary 获取管理员订单详情
// @Description 获取指定订单的详细信息（管理端）
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} common.Response{data=product.ProductOrder} "获取成功"
// @Failure 404 {object} common.Response "订单不存在"
// @Router /admin/orders/{id} [get]
func GetAdminOrderDetail(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的订单ID"))
		return
	}

	service := productService.NewService()
	order, err := service.GetAdminOrderDetail(uint(orderID))
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, order)
}

// UpdateOrderStatus 更新订单状态
// @Summary 更新订单状态
// @Description 管理员手动更新订单的支付或开通状态
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Param request body product.UpdateOrderStatusRequest true "更新订单状态请求参数"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 404 {object} common.Response "订单不存在"
// @Router /admin/orders/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的订单ID"))
		return
	}

	var req productModel.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := productService.NewService()
	if err := service.UpdateOrderStatus(uint(orderID), req); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "订单状态更新成功")
}

// ManualProvision 手动开通实例
// @Summary 手动开通实例
// @Description 管理员手动为已支付订单开通实例
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} common.Response "开通任务已提交"
// @Failure 400 {object} common.Response "订单状态不支持开通"
// @Failure 404 {object} common.Response "订单不存在"
// @Router /admin/orders/{id}/provision [post]
func ManualProvision(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的订单ID"))
		return
	}

	service := productService.NewService()
	if err := service.ManualProvision(uint(orderID)); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "实例开通任务已提交")
}
