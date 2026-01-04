package admin

import (
	"oneclickvirt/global"
	siteModel "oneclickvirt/model/site"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetPublicSiteConfigs 获取站点配置(公开API)
// @Summary 获取站点配置(公开API)
// @Description 获取公开的站点配置信息
// @Tags 公开/站点配置
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /api/v1/public/site-configs [get]
func GetPublicSiteConfigs(c *gin.Context) {
	var configs []siteModel.SiteConfig
	// 只获取公开的配置项
	publicKeys := []string{
		siteModel.SiteName,
		siteModel.SiteURL,
		siteModel.SiteIconURL,
		siteModel.SiteHeader,
		siteModel.SiteFooter,
		siteModel.ContactEmail,
		siteModel.ContactPhone,
		siteModel.CompanyName,
		siteModel.ICPNumber,
		siteModel.AnalyticsCode, // 添加统计代码到公开配置项
	}

	if err := global.APP_DB.Where("`key` IN ?", publicKeys).Find(&configs).Error; err != nil {
		global.APP_LOG.Error("获取站点配置失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取站点配置失败"})
		return
	}

	// 转换为map格式
	configMap := make(map[string]interface{})
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	// 确保contact_phone字段存在
	if _, ok := configMap["contact_phone"]; !ok {
		configMap["contact_phone"] = "888-888-8888"
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "success",
		"data": configMap,
	})
}

// GetSiteConfigs 获取站点配置列表
// @Summary 获取站点配置列表
// @Description 管理员获取所有站点配置
// @Tags 管理员/站点配置
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /v1/admin/site-configs [get]
func GetSiteConfigs(c *gin.Context) {
	var configs []siteModel.SiteConfig
	if err := global.APP_DB.Find(&configs).Error; err != nil {
		global.APP_LOG.Error("获取站点配置失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取站点配置失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "success",
		"data": configs,
	})
}

// GetSiteConfig 获取单个站点配置
// @Summary 获取单个站点配置
// @Description 根据配置键获取站点配置
// @Tags 管理员/站点配置
// @Accept json
// @Produce json
// @Param key path string true "配置键"
// @Success 200 {object} common.Response
// @Router /v1/admin/site-configs/{key} [get]
func GetSiteConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(400, gin.H{"code": 400, "message": "配置键不能为空"})
		return
	}

	var config siteModel.SiteConfig
	if err := global.APP_DB.Where("`key` = ?", key).First(&config).Error; err != nil {
		global.APP_LOG.Error("获取站点配置失败", zap.Error(err))
		c.JSON(404, gin.H{"code": 404, "message": "配置不存在"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "success",
		"data": config,
	})
}

// UpdateSiteConfigs 批量更新站点配置
// @Summary 批量更新站点配置
// @Description 管理员批量更新站点配置
// @Tags 管理员/站点配置
// @Accept json
// @Produce json
// @Param configs body []siteModel.SiteConfig true "配置列表"
// @Success 200 {object} common.Response
// @Router /v1/admin/site-configs [put]
func UpdateSiteConfigs(c *gin.Context) {
	var configs []siteModel.SiteConfig
	if err := c.ShouldBindJSON(&configs); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 使用事务更新
	tx := global.APP_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, config := range configs {
		// 检查配置是否存在
		var existingConfig siteModel.SiteConfig
		if err := tx.Where("`key` = ?", config.Key).First(&existingConfig).Error; err != nil {
			// 不存在则创建
			if err := tx.Create(&config).Error; err != nil {
				tx.Rollback()
				global.APP_LOG.Error("创建站点配置失败", zap.Error(err))
				c.JSON(500, gin.H{"code": 500, "message": "创建配置失败"})
				return
			}
		} else {
			// 存在则更新
			if err := tx.Model(&existingConfig).Updates(map[string]interface{}{
				"value":        config.Value,
				"type":         config.Type,
				"description":  config.Description,
			}).Error; err != nil {
				tx.Rollback()
				global.APP_LOG.Error("更新站点配置失败", zap.Error(err))
				c.JSON(500, gin.H{"code": 500, "message": "更新配置失败"})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		global.APP_LOG.Error("提交事务失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "更新配置失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "更新成功",
	})
}

// InitializeSiteConfigs 初始化默认站点配置
// @Summary 初始化默认站点配置
// @Description 初始化默认站点配置
// @Tags 管理员/站点配置
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /v1/admin/site-configs/initialize [post]
func InitializeSiteConfigs(c *gin.Context) {
	defaultConfigs := []siteModel.SiteConfig{
		{Key: siteModel.SiteName, Value: "OneClickVirt", Type: siteModel.ConfigTypeString, Description: "网站名称"},
		{Key: siteModel.SiteURL, Value: "https://example.com", Type: siteModel.ConfigTypeString, Description: "网站URL"},
		{Key: siteModel.SiteIconURL, Value: "/favicon.ico", Type: siteModel.ConfigTypeImage, Description: "网站图标URL"},
		{Key: siteModel.SiteHeader, Value: "欢迎使用 OneClickVirt", Type: siteModel.ConfigTypeText, Description: "页眉内容"},
		{Key: siteModel.SiteFooter, Value: "© 2025 OneClickVirt. All rights reserved.", Type: siteModel.ConfigTypeText, Description: "页脚内容"},
		{Key: siteModel.ContactEmail, Value: "admin@example.com", Type: siteModel.ConfigTypeString, Description: "联系邮箱"},
		{Key: siteModel.ContactPhone, Value: "", Type: siteModel.ConfigTypeString, Description: "联系电话"},
		{Key: siteModel.CompanyName, Value: "", Type: siteModel.ConfigTypeString, Description: "公司名称"},
		{Key: siteModel.ICPNumber, Value: "", Type: siteModel.ConfigTypeString, Description: "ICP备案号"},
		{Key: siteModel.CustomCSS, Value: "", Type: siteModel.ConfigTypeText, Description: "自定义CSS"},
		{Key: siteModel.AnalyticsCode, Value: "", Type: siteModel.ConfigTypeText, Description: "统计代码"},
	}

	tx := global.APP_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, config := range defaultConfigs {
		var existingConfig siteModel.SiteConfig
		if err := tx.Where("`key` = ?", config.Key).First(&existingConfig).Error; err != nil {
			// 不存在则创建
			if err := tx.Create(&config).Error; err != nil {
				tx.Rollback()
				global.APP_LOG.Error("初始化站点配置失败", zap.Error(err))
				c.JSON(500, gin.H{"code": 500, "message": "初始化配置失败"})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		global.APP_LOG.Error("提交事务失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "初始化配置失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "初始化成功",
	})
}
