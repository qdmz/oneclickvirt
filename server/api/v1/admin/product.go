package admin

import (
	"strconv"

	"oneclickvirt/model/common"
	productModel "oneclickvirt/model/product"
	productService "oneclickvirt/service/product"

	"github.com/gin-gonic/gin"
)

// GetAdminProductList 获取管理员产品列表
// @Summary 获取管理员产品列表
// @Description 获取所有产品列表（包含下架产品）
// @Tags 产品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param status query int false "状态筛选"
// @Param category query string false "类别筛选"
// @Param type query string false "虚拟化类型筛选"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} common.Response{data=common.PageResult} "获取成功"
// @Router /admin/products [get]
func GetAdminProductList(c *gin.Context) {
	var req productModel.AdminProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	req.Normalize(common.DefaultPageSize)

	service := productService.NewService()
	products, total, err := service.GetAdminProductList(req)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccessWithPagination(c, products, total, req.Page, req.PageSize)
}

// CreateAdminProduct 管理员创建产品
// @Summary 管理员创建产品
// @Description 创建新的虚拟机/容器产品
// @Tags 产品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body product.CreateProductRequest true "创建产品请求参数"
// @Success 200 {object} common.Response{data=product.Product} "创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Router /admin/products [post]
func CreateAdminProduct(c *gin.Context) {
	var req productModel.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := productService.NewService()
	product, err := service.CreateAdminProduct(req)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, product, "产品创建成功")
}

// UpdateAdminProduct 管理员更新产品
// @Summary 管理员更新产品
// @Description 更新指定产品的信息
// @Tags 产品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Param request body product.UpdateProductRequest true "更新产品请求参数"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 404 {object} common.Response "产品不存在"
// @Router /admin/products/{id} [put]
func UpdateAdminProduct(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的产品ID"))
		return
	}

	var req productModel.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := productService.NewService()
	if err := service.UpdateAdminProduct(uint(productID), req); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "产品更新成功")
}

// DeleteAdminProduct 管理员删除产品
// @Summary 管理员删除产品
// @Description 删除指定产品
// @Tags 产品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 404 {object} common.Response "产品不存在"
// @Router /admin/products/{id} [delete]
func DeleteAdminProduct(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的产品ID"))
		return
	}

	service := productService.NewService()
	if err := service.DeleteAdminProduct(uint(productID)); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "产品删除成功")
}

// GetAdminProductDetail 获取管理员产品详情
// @Summary 获取管理员产品详情
// @Description 获取指定产品的详细信息（管理端）
// @Tags 产品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Success 200 {object} common.Response{data=product.Product} "获取成功"
// @Failure 404 {object} common.Response "产品不存在"
// @Router /admin/products/{id} [get]
func GetAdminProductDetail(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的产品ID"))
		return
	}

	service := productService.NewService()
	product, err := service.GetProductDetail(uint(productID))
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, product)
}

// UpdateProductStatus 更新产品上下架状态
// @Summary 更新产品上下架状态
// @Description 快速更新产品的上下架状态
// @Tags 产品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Param request body product.UpdateProductStatusRequest true "更新状态请求参数"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Router /admin/products/{id}/status [put]
func UpdateProductStatus(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的产品ID"))
		return
	}

	var req productModel.UpdateProductStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := productService.NewService()
	// 复用 UpdateAdminProduct 的更新逻辑，只更新状态字段
	updateReq := productModel.UpdateProductRequest{
		Status: req.Status,
	}
	// 先获取原产品信息填充必填字段
	product, err := service.GetProductDetail(uint(productID))
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	updateReq.Name = product.Name
	updateReq.Type = product.Type
	updateReq.CPU = product.CPU
	updateReq.Memory = product.Memory
	updateReq.Disk = product.Disk
	updateReq.Price = product.Price
	updateReq.PeriodType = product.PeriodType
	updateReq.PeriodValue = product.PeriodValue

	if err := service.UpdateAdminProduct(uint(productID), updateReq); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "状态更新成功")
}
