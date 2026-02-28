package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

// InputValidator 输入验证中间件
func InputValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过支付回调路径的验证
		path := c.Request.URL.Path
		if isPaymentCallbackPath(path) {
			c.Next()
			return
		}

		// 检查SQL注入
		if containsSQLInjection(c.Request.URL.RawQuery) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "检测到潜在的SQL注入攻击",
			})
			c.Abort()
			return
		}

		// 检查XSS攻击
		if containsXSS(c.Request.URL.RawQuery) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "检测到潜在的XSS攻击",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isPaymentCallbackPath 检查是否是支付回调路径
func isPaymentCallbackPath(path string) bool {
	paymentPaths := []string{
		"/api/v1/payment/alipay/notify",
		"/api/v1/payment/wechat/notify",
		"/api/v1/payment/epay/notify",
		"/api/v1/payment/mapay/notify",
	}
	for _, pp := range paymentPaths {
		if path == pp {
			return true
		}
	}
	return false
}

// containsSQLInjection 检查是否包含SQL注入模式
func containsSQLInjection(input string) bool {
	sqlPatterns := []string{
		`(?i)(union\s+select)`,
		`(?i)(drop\s+table)`,
		`(?i)(delete\s+from)`,
		`(?i)(insert\s+into)`,
		`(?i)(update\s+set)`,
		`(?i)(exec\s*\()`,
		`(?i)(script\s*>)`,
		`(?i)(\'\s*or\s*\'\s*=\s*\')`,
		`(?i)(\'\s*or\s*1\s*=\s*1)`,
		`(?i)(--\s)`,
		`(?i)(/\*.*\*/)`,
	}

	for _, pattern := range sqlPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			return true
		}
	}
	return false
}

// containsXSS 检查是否包含XSS攻击模式
func containsXSS(input string) bool {
	xssPatterns := []string{
		`(?i)(<script[^>]*>)`,
		`(?i)(</script>)`,
		`(?i)(javascript:)`,
		`(?i)(on\w+\s*=)`,
		`(?i)(<iframe[^>]*>)`,
		`(?i)(<object[^>]*>)`,
		`(?i)(<embed[^>]*>)`,
		`(?i)(<link[^>]*>)`,
		`(?i)(<meta[^>]*>)`,
	}

	for _, pattern := range xssPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			return true
		}
	}
	return false
}
