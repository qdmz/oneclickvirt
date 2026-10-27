package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders 为所有响应设置基础安全响应头：
//   - X-Content-Type-Options: nosniff —— 阻止浏览器 MIME 嗅探（配合上传文件的扩展名白名单降低伪装文件风险）
//   - X-Frame-Options: SAMEORIGIN —— 阻止第三方站点 iframe 嵌入，防点击劫持（与控制台代理既有取值一致，见 api/v1/admin/console_proxy.go）
//   - Referrer-Policy: strict-origin-when-cross-origin —— 跨站请求只发送来源站点，不泄漏完整 URL 路径
//
// 暂不设置全局 Content-Security-Policy：控制台代理等页面已有各自的 CSP，
// 全局 CSP 需要先盘点内联脚本/样式后再单独引入，避免误伤现有页面。
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "SAMEORIGIN")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}
