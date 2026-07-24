package product

import (
	"oneclickvirt/model/common"
	"oneclickvirt/service/product"

	"github.com/gin-gonic/gin"
)

// GetSiteConfig 获取站点配置
// @Summary 获取站点配置
// @Description 获取站点前端配置信息（前端使用）
// @Tags 站点配置
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 500 {object} common.Response "获取失败"
// @Router /public/site-config [get]
func GetSiteConfig(c *gin.Context) {
	service := product.NewSiteService()
	config, err := service.GetPublicSiteInfo()
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, config)
}

// GetPublicSiteInfo 获取公开站点信息
// @Summary 获取公开站点信息
// @Description 获取站点的公开基本信息
// @Tags 站点配置
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 500 {object} common.Response "获取失败"
// @Router /public/site-info [get]
func GetPublicSiteInfo(c *gin.Context) {
	service := product.NewSiteService()
	info, err := service.GetPublicSiteInfo()
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, info)
}
