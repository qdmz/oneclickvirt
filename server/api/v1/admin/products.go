package admin

import (
	"strconv"

	"oneclickvirt/middleware"
	"oneclickvirt/model/common"
	"oneclickvirt/model/product"
	productSvc "oneclickvirt/service/product"

	"github.com/gin-gonic/gin"
)

// ListProductsAdmin 管理端产品列表
// @Router /admin/products [get]
func ListProductsAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	svc := productSvc.NewService()
	list, total, err := svc.AdminList(page, pageSize)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccessWithPagination(c, list, total, page, pageSize)
}

// CreateProductAdmin 创建产品
// @Router /admin/products [post]
func CreateProductAdmin(c *gin.Context) {
	// 仅超管可管理产品
	if !middleware.IsSuperAdmin(c) {
		common.ResponseWithError(c, common.NewError(common.CodeForbidden, "需要超管权限"))
		return
	}
	var p product.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误: "+err.Error()))
		return
	}
	svc := productSvc.NewService()
	created, err := svc.AdminCreate(&p)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, created, "产品已创建")
}

// UpdateProductAdmin 更新产品
// @Router /admin/products/:id [put]
func UpdateProductAdmin(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		common.ResponseWithError(c, common.NewError(common.CodeForbidden, "需要超管权限"))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的产品ID"))
		return
	}
	var upd map[string]interface{}
	if err := c.ShouldBindJSON(&upd); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误: "+err.Error()))
		return
	}
	svc := productSvc.NewService()
	if err := svc.AdminUpdate(uint(id), upd); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, nil, "产品已更新")
}

// DeleteProductAdmin 删除产品
// @Router /admin/products/:id [delete]
func DeleteProductAdmin(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		common.ResponseWithError(c, common.NewError(common.CodeForbidden, "需要超管权限"))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的产品ID"))
		return
	}
	svc := productSvc.NewService()
	if err := svc.AdminDelete(uint(id)); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, nil, "产品已删除")
}
