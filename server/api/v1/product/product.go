package product

import (
	"strconv"

	"oneclickvirt/model/common"
	productModel "oneclickvirt/model/product"
	productService "oneclickvirt/service/product"

	"github.com/gin-gonic/gin"
)

// GetProductList 获取上架产品列表
// @Summary 获取上架产品列表
// @Description 获取所有状态为上架的产品列表
// @Tags 产品商城
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param category query string false "类别筛选"
// @Param type query string false "虚拟化类型筛选"
// @Success 200 {object} common.Response{data=common.PageResult} "获取成功"
// @Router /products [get]
func GetProductList(c *gin.Context) {
	var req productModel.ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	req.Normalize(common.DefaultPageSize)

	service := productService.NewService()
	products, total, err := service.GetProductList(req)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccessWithPagination(c, products, total, req.Page, req.PageSize)
}

// GetProductDetail 获取产品详情
// @Summary 获取产品详情
// @Description 获取指定产品的详细信息
// @Tags 产品商城
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} common.Response{data=product.Product} "获取成功"
// @Failure 404 {object} common.Response "产品不存在"
// @Router /products/{id} [get]
func GetProductDetail(c *gin.Context) {
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

// GetPublicProductList 获取公开产品列表(无需登录)
func GetPublicProductList(c *gin.Context) {
	GetProductList(c)
}

// GetPublicProductDetail 获取公开产品详情(无需登录)
func GetPublicProductDetail(c *gin.Context) {
	GetProductDetail(c)
}

// GetProductImages 获取产品可用镜像
// @Summary 获取产品可用镜像
// @Description 获取指定产品配置的可用的系统镜像列表
// @Tags 产品商城
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} common.Response{data=[]system.SystemImage} "获取成功"
// @Failure 404 {object} common.Response "产品不存在"
// @Router /products/{id}/images [get]
func GetProductImages(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的产品ID"))
		return
	}

	service := productService.NewService()
	images, err := service.GetProductImages(uint(productID))
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, images)
}
