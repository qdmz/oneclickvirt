package middleware

import (
	"net/http"

	"oneclickvirt/global"
	"oneclickvirt/model/agent"
	"oneclickvirt/model/common"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequireAgent 检查当前用户是否是状态正常的代理商
func RequireAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, "用户未认证"))
			c.Abort()
			return
		}

		var a agent.Agent
		if err := global.APP_DB.Where("user_id = ? AND status = 1", userID).First(&a).Error; err != nil {
			global.APP_LOG.Debug("用户不是有效代理商", zap.Uint("userID", userID), zap.Error(err))
			c.JSON(http.StatusForbidden, gin.H{
				"code":    common.CodeForbidden,
				"message": "您不是有效代理商",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Set("agent_id", a.ID)
		c.Set("agent", &a)
		c.Next()
	}
}
