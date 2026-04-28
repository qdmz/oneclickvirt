package middleware

import (
	"errors"
	"net/http"
	"oneclickvirt/global"
	authModel "oneclickvirt/model/auth"
	userModel "oneclickvirt/model/user"
	authService "oneclickvirt/service/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth API 密钥认证中间件
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "缺少认证信息",
			})
			c.Abort()
			return
		}

		// 检查是否是 API Key
		if strings.HasPrefix(authHeader, "ApiKey ") {
			// API Key 认证
			apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
			authCtx, err := validateAPIKey(apiKey)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "API Key 无效: " + err.Error(),
				})
				c.Abort()
				return
			}
			c.Set("auth_context", authCtx)
			c.Next()
			return
		}

		// 如果不是 API Key，返回错误
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "不支持的认证方式",
		})
		c.Abort()
	}
}

// APIKeyOnly 仅 API 密钥认证中间件
func APIKeyOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "缺少认证信息",
			})
			c.Abort()
			return
		}

		// 检查是否是 API Key
		if !strings.HasPrefix(authHeader, "ApiKey ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "仅支持 API Key 认证",
			})
			c.Abort()
			return
		}

		// API Key 认证
		apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
		authCtx, err := validateAPIKey(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "API Key 无效: " + err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("auth_context", authCtx)
		c.Next()
	}
}

// validateAPIKey 验证 API 密钥
func validateAPIKey(apiKey string) (*authModel.AuthContext, error) {
	service := authService.NewAPIKeyService()
	keyModel, err := service.ValidateAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	var user userModel.User
	if err := global.APP_DB.Select("id, username, user_type, status, level").First(&user, keyModel.UserID).Error; err != nil {
		return nil, err
	}

	// 严格检查用户状态
	if user.Status != 1 {
		return nil, errors.New("账户已被禁用")
	}

	// 使用权限服务获取用户有效权限
	permissionService := authService.PermissionService{}
	effectivePermission, err := permissionService.GetUserEffectivePermission(user.ID)
	if err != nil {
		return nil, errors.New("权限验证失败")
	}

	// 构建认证上下文
	authCtx := &authModel.AuthContext{
		UserID:       user.ID,
		Username:     user.Username,
		UserType:     effectivePermission.EffectiveType,
		Level:        effectivePermission.EffectiveLevel,
		BaseUserType: user.UserType,
		AllUserTypes: effectivePermission.AllTypes,
		IsEffective:  true,
	}

	return authCtx, nil
}
