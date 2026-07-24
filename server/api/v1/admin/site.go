package admin

import (
	"io"

	"oneclickvirt/model/common"
	productModel "oneclickvirt/model/product"
	"oneclickvirt/service/product"

	"github.com/gin-gonic/gin"
)

// GetAdminSiteConfig 获取完整站点配置
// @Summary 获取完整站点配置
// @Description 管理员获取完整的站点配置信息
// @Tags 管理员站点配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 401 {object} common.Response "未授权"
// @Failure 500 {object} common.Response "获取失败"
// @Router /admin/site-config [get]
func GetAdminSiteConfig(c *gin.Context) {
	service := product.NewSiteService()
	config, err := service.GetAdminSiteConfig()
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, config)
}

// UpdateSiteConfig 更新站点配置
// @Summary 更新站点配置
// @Description 管理员更新站点配置信息
// @Tags 管理员站点配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body productModel.SiteConfig true "站点配置参数"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 500 {object} common.Response "更新失败"
// @Router /admin/site-config [put]
func UpdateSiteConfig(c *gin.Context) {
	var config productModel.SiteConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := product.NewSiteService()
	if err := service.UpdateSiteConfig(&config); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "站点配置更新成功")
}

// UploadSiteImage 上传站点图片
// @Summary 上传站点图片
// @Description 管理员上传站点图片（Logo、Favicon等）
// @Tags 管理员站点配置
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param image formData file true "图片文件"
// @Success 200 {object} common.Response{data=object} "上传成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 500 {object} common.Response "上传失败"
// @Router /admin/site-config/upload-image [post]
func UploadSiteImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "请选择要上传的图片"))
		return
	}
	defer file.Close()

	// 读取文件内容
	imageData, err := io.ReadAll(file)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "读取图片失败"))
		return
	}

	service := product.NewSiteService()
	imageURL, err := service.UploadSiteImage(imageData, header.Filename)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, gin.H{
		"url": imageURL,
	}, "图片上传成功")
}
