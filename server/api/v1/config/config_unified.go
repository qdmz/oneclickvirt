package config

import (
	"fmt"
	"net/http"
	"oneclickvirt/service/auth"
	"oneclickvirt/service/email"
	"strings"

	"oneclickvirt/config"
	"oneclickvirt/global"
	"oneclickvirt/middleware"
	authModel "oneclickvirt/model/auth"
	"oneclickvirt/model/common"
	configModel "oneclickvirt/model/config"

	"github.com/gin-gonic/gin"
)

// GetUnifiedConfig 获取统一配置接口
// @Summary 获取系统配置
// @Description 根据用户权限返回相应的配置信息
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param scope query string false "配置范围" Enums(public,user,admin) default(user)
// @Success 200 {object} common.Response{data=interface{}} "获取成功"
// @Failure 401 {object} common.Response "认证失败"
// @Failure 403 {object} common.Response "权限不足"
// @Failure 500 {object} common.Response "获取失败"
// @Router /config [get]
func GetUnifiedConfig(c *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.Response{
			Code: 401,
			Msg:  "用户未认证",
		})
		return
	}

	// 根据请求路径自动判断 scope
	scope := c.DefaultQuery("scope", "")
	if scope == "" {
		// 如果没有提供 scope 参数，根据路径判断
		if strings.Contains(c.Request.URL.Path, "/admin/") {
			scope = "admin"
		} else if strings.Contains(c.Request.URL.Path, "/public/") {
			scope = "public"
		} else {
			scope = "user"
		}
	}

	// 根据用户权限和请求范围决定返回的配置
	configManager := config.GetConfigManager()
	if configManager == nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Code: 500,
			Msg:  "配置管理器未初始化",
		})
		return
	}

	var result map[string]interface{}

	switch scope {
	case "public":
		// 公开配置，所有用户都可以访问
		result = getPublicConfig(configManager)
	case "user":
		// 用户配置，普通用户可以访问的配置
		result = getUserConfig(configManager, authCtx)
	case "admin", "global":
		// 管理员配置和全局配置，只有管理员可以访问
		permissionService := auth.PermissionService{}
		hasAdminPermission := permissionService.HasPermission(authCtx.UserID, "admin")
		if !hasAdminPermission {
			c.JSON(http.StatusForbidden, common.Response{
				Code: 403,
				Msg:  "权限不足",
			})
			return
		}
		result = getAdminConfig(configManager)
	default:
		c.JSON(http.StatusBadRequest, common.Response{
			Code: 400,
			Msg:  "无效的配置范围",
		})
		return
	}

	common.ResponseSuccess(c, result)
}

// UpdateUnifiedConfig 更新统一配置接口
// @Summary 更新系统配置
// @Description 根据用户权限更新相应的配置信息
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body configModel.UnifiedConfigRequest true "配置更新请求"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "认证失败"
// @Failure 403 {object} common.Response "权限不足"
// @Failure 500 {object} common.Response "更新失败"
// @Router /config [put]
func UpdateUnifiedConfig(c *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.Response{
			Code: 401,
			Msg:  "用户未认证",
		})
		return
	}

	// 解析请求体
	var rawData map[string]interface{}
	if err := c.ShouldBindJSON(&rawData); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	var req configModel.UnifiedConfigRequest

	// 检查是否是新的统一格式
	if scope, exists := rawData["scope"]; exists {
		if config, configExists := rawData["config"]; configExists {
			req.Scope = scope.(string)
			req.Config = config.(map[string]interface{})
		} else {
			common.ResponseWithError(c, common.NewError(common.CodeValidationError, "统一格式缺少config字段"))
			return
		}
	} else {
		// 向后兼容：直接配置数据，根据路径判断 scope
		if strings.Contains(c.Request.URL.Path, "/admin/") {
			req.Scope = "admin"
		} else {
			req.Scope = "user"
		}
		req.Config = rawData
	}

	// 验证权限
	if !hasConfigUpdatePermission(authCtx, req.Scope) {
		c.JSON(http.StatusForbidden, common.Response{
			Code: 403,
			Msg:  "权限不足",
		})
		return
	}

	configManager := config.GetConfigManager()
	if configManager == nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Code: 500,
			Msg:  "配置管理器未初始化",
		})
		return
	}

	// 根据范围过滤配置项
	filteredConfig := filterConfigByScope(req.Config, req.Scope, authCtx)

	// 更新配置
	// UpdateConfig 会自动：
	// 1. 将配置保存到数据库（自动转换为 kebab-case 格式）
	// 2. 通过已注册的回调函数同步到 global.APP_CONFIG
	// 3. 写回到 YAML 文件
	if err := configManager.UpdateConfig(filteredConfig); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeConfigError, err.Error()))
		return
	}

	// ConfigManager.UpdateConfig 已经通过回调机制自动同步到全局配置
	// 回调函数在 initialize/config_manager.go 的 syncConfigToGlobal 中定义
	// 它会正确处理 kebab-case 和 camelCase 两种格式的键名

	common.ResponseSuccess(c, nil, "配置更新成功")
}

// getPublicConfig 获取公开配置
func getPublicConfig(cm *config.ConfigManager) map[string]interface{} {
	allConfig := cm.GetAllConfig()
	publicConfig := make(map[string]interface{})

	// 只返回公开的配置项
	publicKeys := []string{
		"app.name",
		"app.version",
		"app.description",
		"auth.enablePublicRegistration",
	}

	for _, key := range publicKeys {
		if value, exists := allConfig[key]; exists {
			publicConfig[key] = value
		}
	}

	// 将扁平化配置转换为嵌套结构
	return unflattenConfig(publicConfig)
}

// getUserConfig 获取用户配置（使用服务端权限验证）
func getUserConfig(cm *config.ConfigManager, authCtx *authModel.AuthContext) map[string]interface{} {
	result := make(map[string]interface{})
	permissionService := auth.PermissionService{}

	// 基础配置 - 所有用户可见
	result["auth"] = map[string]interface{}{
		"enablePublicRegistration": global.APP_CONFIG.Auth.EnablePublicRegistration,
	}

	// 配额配置 - 从 global.APP_CONFIG 获取完整配置
	levelLimits := make(map[string]interface{})
	for level, limitInfo := range global.APP_CONFIG.Quota.LevelLimits {
		levelKey := fmt.Sprintf("%d", level)
		levelLimits[levelKey] = map[string]interface{}{
			"max-instances": limitInfo.MaxInstances,
			"max-resources": limitInfo.MaxResources,
			"max-traffic":   limitInfo.MaxTraffic,
		}
	}

	result["quota"] = map[string]interface{}{
		"defaultLevel": global.APP_CONFIG.Quota.DefaultLevel,
		"levelLimits":  levelLimits,
	}

	// 管理员可以看到更多配置
	hasAdminPermission := permissionService.HasPermission(authCtx.UserID, "admin")
	if hasAdminPermission {
		authConfig := result["auth"].(map[string]interface{})
		authConfig["enableEmail"] = global.APP_CONFIG.Auth.EnableEmail
		authConfig["enableTelegram"] = global.APP_CONFIG.Auth.EnableTelegram
		authConfig["enableQQ"] = global.APP_CONFIG.Auth.EnableQQ
	}

	return result
}

// getAdminConfig 获取管理员配置
func getAdminConfig(cm *config.ConfigManager) map[string]interface{} {
	// 直接从 global.APP_CONFIG 构建完整配置返回
	// 确保返回所有配置项（包括默认值）
	result := make(map[string]interface{})

	// 认证配置
	result["auth"] = map[string]interface{}{
		"enableEmail":              global.APP_CONFIG.Auth.EnableEmail,
		"enableTelegram":           global.APP_CONFIG.Auth.EnableTelegram,
		"enableQQ":                 global.APP_CONFIG.Auth.EnableQQ,
		"enableOAuth2":             global.APP_CONFIG.Auth.EnableOAuth2,
		"enablePublicRegistration": global.APP_CONFIG.Auth.EnablePublicRegistration,
		"emailSMTPHost":            global.APP_CONFIG.Auth.EmailSMTPHost,
		"emailSMTPPort":            global.APP_CONFIG.Auth.EmailSMTPPort,
		"emailUsername":            global.APP_CONFIG.Auth.EmailUsername,
		"emailPassword":            global.APP_CONFIG.Auth.EmailPassword,
		"telegramBotToken":         global.APP_CONFIG.Auth.TelegramBotToken,
		"qqAppID":                  global.APP_CONFIG.Auth.QQAppID,
		"qqAppKey":                 global.APP_CONFIG.Auth.QQAppKey,
	}

	// 邀请码配置
	result["inviteCode"] = map[string]interface{}{
		"enabled":  global.APP_CONFIG.InviteCode.Enabled,
		"required": global.APP_CONFIG.InviteCode.Required,
	}

	// 配额配置 - 从 global.APP_CONFIG 获取完整的等级限制
	levelLimits := make(map[string]interface{})
	for level, limitInfo := range global.APP_CONFIG.Quota.LevelLimits {
		levelKey := fmt.Sprintf("%d", level)
		levelLimits[levelKey] = map[string]interface{}{
			"max-instances": limitInfo.MaxInstances,
			"max-resources": limitInfo.MaxResources,
			"max-traffic":   limitInfo.MaxTraffic,
		}
	}

	result["quota"] = map[string]interface{}{
		"defaultLevel": global.APP_CONFIG.Quota.DefaultLevel,
		"levelLimits":  levelLimits,
		"instanceTypePermissions": map[string]interface{}{
			"minLevelForContainer":       global.APP_CONFIG.Quota.InstanceTypePermissions.MinLevelForContainer,
			"minLevelForVM":              global.APP_CONFIG.Quota.InstanceTypePermissions.MinLevelForVM,
			"minLevelForDeleteContainer": global.APP_CONFIG.Quota.InstanceTypePermissions.MinLevelForDeleteContainer,
			"minLevelForDeleteVM":        global.APP_CONFIG.Quota.InstanceTypePermissions.MinLevelForDeleteVM,
			"minLevelForResetContainer":  global.APP_CONFIG.Quota.InstanceTypePermissions.MinLevelForResetContainer,
			"minLevelForResetVM":         global.APP_CONFIG.Quota.InstanceTypePermissions.MinLevelForResetVM,
		},
	}

	// 其他配置
	result["other"] = map[string]interface{}{
		"maxAvatarSize":   global.APP_CONFIG.Other.MaxAvatarSize,
		"defaultLanguage": global.APP_CONFIG.Other.DefaultLanguage,
	}

	// 系统配置
	result["system"] = map[string]interface{}{
		"addr":                     global.APP_CONFIG.System.Addr,
		"dbType":                   global.APP_CONFIG.System.DbType,
		"env":                      global.APP_CONFIG.System.Env,
		"frontendURL":              global.APP_CONFIG.System.FrontendURL,
		"ipLimitCount":             global.APP_CONFIG.System.LimitCountIP,
		"ipLimitTime":              global.APP_CONFIG.System.LimitTimeIP,
		"oauth2StateTokenMinutes":  global.APP_CONFIG.System.OAuth2StateTokenMinutes,
		"ossType":                  global.APP_CONFIG.System.OssType,
		"providerInactiveHours":    global.APP_CONFIG.System.ProviderInactiveHours,
		"useMultipoint":            global.APP_CONFIG.System.UseMultipoint,
		"useRedis":                 global.APP_CONFIG.System.UseRedis,
	}

	// 验证码配置
	result["captcha"] = map[string]interface{}{
		"enabled":    global.APP_CONFIG.Captcha.Enabled,
		"expireTime": global.APP_CONFIG.Captcha.ExpireTime,
		"height":     global.APP_CONFIG.Captcha.Height,
		"length":     global.APP_CONFIG.Captcha.Length,
		"width":      global.APP_CONFIG.Captcha.Width,
	}

	// 任务配置
	result["task"] = map[string]interface{}{
		"deleteRetryCount": global.APP_CONFIG.Task.DeleteRetryCount,
		"deleteRetryDelay": global.APP_CONFIG.Task.DeleteRetryDelay,
	}

	// 站点配置
	result["site"] = map[string]interface{}{
		"name":        "OneClickVirt",
		"description": "虚拟化管理平台",
		"keywords":    "虚拟化,Docker,LXD,Incus,Proxmox",
	}

	// 系统名称和描述
	result["systemName"] = "虚拟化管理平台"
	result["systemDescription"] = "支持多种虚拟化技术的管理平台"

	// 支付接口配置
	result["payment"] = map[string]interface{}{
		"alipayAppId":       global.APP_CONFIG.Payment.AlipayAppID,
		"alipayGateway":     global.APP_CONFIG.Payment.AlipayGateway,
		"alipayPrivateKey":  global.APP_CONFIG.Payment.AlipayPrivateKey,
		"alipayPublicKey":   global.APP_CONFIG.Payment.AlipayPublicKey,
		"balanceDailyLimit": global.APP_CONFIG.Payment.BalanceDailyLimit,
		"balanceMinAmount":  global.APP_CONFIG.Payment.BalanceMinAmount,
		"enableAlipay":      global.APP_CONFIG.Payment.EnableAlipay,
		"enableBalance":     global.APP_CONFIG.Payment.EnableBalance,
		"enableWechat":      global.APP_CONFIG.Payment.EnableWechat,
		"wechatApiKey":      global.APP_CONFIG.Payment.WechatAPIKey,
		"wechatApiV3Key":    global.APP_CONFIG.Payment.WechatAPIV3Key,
		"wechatAppId":       global.APP_CONFIG.Payment.WechatAppID,
		"wechatMchId":       global.APP_CONFIG.Payment.WechatMchID,
		"wechatSerialNo":    global.APP_CONFIG.Payment.WechatSerialNo,
		"wechatType":        global.APP_CONFIG.Payment.WechatType,
		// 易支付配置
		"enableEpay":     global.APP_CONFIG.Payment.EnableEpay,
		"epayAPIURL":     global.APP_CONFIG.Payment.EpayAPIURL,
		"epayPID":        global.APP_CONFIG.Payment.EpayPID,
		"epayKey":        global.APP_CONFIG.Payment.EpayKey,
		"epayReturnURL":  global.APP_CONFIG.Payment.EpayReturnURL,
		"epayNotifyURL":  global.APP_CONFIG.Payment.EpayNotifyURL,
		// 码支付配置
		"enableMapay":    global.APP_CONFIG.Payment.EnableMapay,
		"mapayAPIURL":    global.APP_CONFIG.Payment.MapayAPIURL,
		"mapayID":        global.APP_CONFIG.Payment.MapayID,
		"mapayKey":       global.APP_CONFIG.Payment.MapayKey,
		"mapayReturnURL": global.APP_CONFIG.Payment.MapayReturnURL,
		"mapayNotifyURL": global.APP_CONFIG.Payment.MapayNotifyURL,
	}

	return result
} // unflattenConfig 将扁平化的配置（如 quota.defaultLevel）转换为嵌套结构（如 quota: { defaultLevel: 1 }）
func unflattenConfig(flatConfig map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range flatConfig {
		setNestedValue(result, key, value)
	}

	return result
}

// setNestedValue 将点分隔的键设置为嵌套结构
func setNestedValue(target map[string]interface{}, key string, value interface{}) {
	keys := strings.Split(key, ".")
	current := target

	for i := 0; i < len(keys)-1; i++ {
		k := keys[i]
		if _, exists := current[k]; !exists {
			current[k] = make(map[string]interface{})
		}
		if nested, ok := current[k].(map[string]interface{}); ok {
			current = nested
		}
	}

	current[keys[len(keys)-1]] = value
}

// hasConfigUpdatePermission 检查配置更新权限（使用服务端权限验证）
func hasConfigUpdatePermission(authCtx *authModel.AuthContext, scope string) bool {
	// 使用权限服务进行服务端权限验证
	permissionService := auth.PermissionService{}

	switch scope {
	case "public":
		// 公开配置不允许更新
		return false
	case "user":
		// 普通用户配置，管理员可以更新
		// 使用权限服务验证，而不是依赖客户端传入的userType
		hasAdminPermission := permissionService.HasPermission(authCtx.UserID, "admin")
		return hasAdminPermission
	case "admin", "global":
		// 管理员配置和全局配置，只有管理员可以更新
		hasAdminPermission := permissionService.HasPermission(authCtx.UserID, "admin")
		return hasAdminPermission
	default:
		return false
	}
}

// filterConfigByScope 根据范围过滤配置（使用服务端权限验证）
func filterConfigByScope(config map[string]interface{}, scope string, authCtx *authModel.AuthContext) map[string]interface{} {
	filtered := make(map[string]interface{})
	permissionService := auth.PermissionService{}

	switch scope {
	case "user":
		// 只允许更新用户级别的配置
		allowedKeys := map[string]bool{
			"quota.defaultLevel": true,
			"quota.levelLimits":  true,
		}

		// 使用权限服务验证，而不是依赖JWT中的userType
		hasAdminPermission := permissionService.HasPermission(authCtx.UserID, "admin")
		if hasAdminPermission {
			allowedKeys["auth.enablePublicRegistration"] = true
			// 允许邮箱相关配置
			allowedKeys["auth.enableEmail"] = true
			allowedKeys["auth.emailSMTPPort"] = true
			allowedKeys["auth.emailSMTPHost"] = true
			allowedKeys["auth.emailUsername"] = true
			allowedKeys["auth.emailPassword"] = true
			// 允许其他认证相关配置
			allowedKeys["auth.enableTelegram"] = true
			allowedKeys["auth.enableQQ"] = true
			allowedKeys["auth.enableOAuth2"] = true
			allowedKeys["auth.telegramBotToken"] = true
			allowedKeys["auth.qqAppID"] = true
			allowedKeys["auth.qqAppKey"] = true
		}

		for key, value := range config {
			if allowedKeys[key] {
				filtered[key] = value
			}
		}
	case "admin":
		// 管理员可以更新所有非系统级配置
		hasAdminPermission := permissionService.HasPermission(authCtx.UserID, "admin")
		if hasAdminPermission {
			for key, value := range config {
				// 过滤掉系统级配置
				if !isSystemLevelConfig(key) {
					filtered[key] = value
				}
			}
		}
	case "global":
		// 全局配置，只有管理员可以更新，且排除系统级配置
		hasAdminPermission := permissionService.HasPermission(authCtx.UserID, "admin")
		if hasAdminPermission {
			for key, value := range config {
				// 过滤掉系统级配置
				if !isSystemLevelConfig(key) {
					filtered[key] = value
				}
			}
		}
	}

	return filtered
}

// isSystemLevelConfig 检查是否为系统级配置（启动必需，必须来自YAML）
func isSystemLevelConfig(key string) bool {
	systemLevelConfigKeys := map[string]bool{
		// System 配置（所有 system.* 都是系统级配置）
		"system.addr":                       true,
		"system.db-type":                    true,
		"system.env":                        true,
		"system.frontend-url":               true,
		"system.iplimit-count":              true,
		"system.iplimit-time":               true,
		"system.oauth2-state-token-minutes": true,
		"system.oss-type":                   true,
		"system.provider-inactive-hours":    true,
		"system.use-multipoint":             true,
		"system.use-redis":                  true,

		// MySQL 配置（数据库连接信息，必须在连接数据库前读取）
		"mysql.path":           true,
		"mysql.port":           true,
		"mysql.config":         true,
		"mysql.db-name":        true,
		"mysql.username":       true,
		"mysql.password":       true,
		"mysql.prefix":         true,
		"mysql.singular":       true,
		"mysql.engine":         true,
		"mysql.max-idle-conns": true,
		"mysql.max-open-conns": true,
		"mysql.max-lifetime":   true,
		"mysql.log-mode":       true,
		"mysql.log-zap":        true,
		"mysql.auto-create":    true,

		// Redis 配置（如果启用Redis，也是启动必需）
		"redis.addr":     true,
		"redis.password": true,
		"redis.db":       true,

		// Zap 日志配置（日志系统启动必需）
		"zap.level":              true,
		"zap.format":             true,
		"zap.prefix":             true,
		"zap.director":           true,
		"zap.encode-level":       true,
		"zap.stacktrace-key":     true,
		"zap.max-file-size":      true,
		"zap.max-backups":        true,
		"zap.max-log-length":     true,
		"zap.retention-day":      true,
		"zap.show-line":          true,
		"zap.log-in-console":     true,
		"zap.max-string-length":  true,
		"zap.max-array-elements": true,
	}
	return systemLevelConfigKeys[key]
}

// TestEmailSend 测试邮件发送接口
// @Summary 测试邮件发送
// @Description 测试SMTP配置是否能正确发送邮件
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "测试邮件请求" schema:{"recipient":"test@example.com"}
// @Success 200 {object} common.Response "发送成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "认证失败"
// @Failure 403 {object} common.Response "权限不足"
// @Failure 500 {object} common.Response "发送失败"
// @Router /admin/config/test-email [post]
func TestEmailSend(c *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.Response{
			Code: 401,
			Msg:  "用户未认证",
		})
		return
	}

	// 验证权限
	permissionService := auth.PermissionService{}
	hasAdminPermission := permissionService.HasPermission(authCtx.UserID, "admin")
	if !hasAdminPermission {
		c.JSON(http.StatusForbidden, common.Response{
			Code: 403,
			Msg:  "权限不足",
		})
		return
	}

	// 解析请求体
	var req struct {
		Recipient string `json:"recipient" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误: 请输入有效的邮箱地址"))
		return
	}

	// 检查邮箱配置
	config := global.APP_CONFIG.Auth
	if !config.EnableEmail {
		common.ResponseWithError(c, common.NewError(common.CodeConfigError, "邮箱功能未启用"))
		return
	}

	if config.EmailSMTPHost == "" || config.EmailSMTPPort == 0 || config.EmailUsername == "" || config.EmailPassword == "" {
		common.ResponseWithError(c, common.NewError(common.CodeConfigError, "SMTP配置不完整"))
		return
	}

	// 发送测试邮件
	emailService := email.NewEmailService()
	err := emailService.SendWelcomeEmail(req.Recipient, "测试用户")
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeConfigError, fmt.Sprintf("邮件发送失败: %v", err)))
		return
	}

	common.ResponseSuccess(c, nil, "测试邮件发送成功，请查收")
}
