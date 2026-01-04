package middleware

import (
	"runtime/debug"

	"oneclickvirt/global"
	"oneclickvirt/model/common"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic信息
				global.APP_LOG.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("ip", c.ClientIP()),
				)

				// 返回500错误
				common.ResponseWithError(c, common.NewError(common.CodeInternalError, "系统遇到了意外错误，请稍后重试"))
				c.Abort()
			}
		}()

		c.Next()

		// 处理业务错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// 记录错误日志
			global.APP_LOG.Error("Request error",
				zap.Error(err.Err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("ip", c.ClientIP()),
			)

			// 根据错误类型返回不同的HTTP状态码
			switch err.Type {
			case gin.ErrorTypeBind:
				common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, err.Error()))
			case gin.ErrorTypePublic:
				common.ResponseWithError(c, common.NewError(common.CodeError, err.Error()))
			default:
				common.ResponseWithError(c, common.NewError(common.CodeInternalError, "请稍后重试"))
			}
		}
	}
}
