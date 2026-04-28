package router

import (
	"oneclickvirt/api/v1/auth"
	"oneclickvirt/api/v1/payment"
	"oneclickvirt/middleware"
	authModel "oneclickvirt/model/auth"

	"github.com/gin-gonic/gin"
)

// InitAuthRouter 认证路由
func InitAuthRouter(Router *gin.RouterGroup) {
	AuthRouter := Router.Group("v1/auth")
	{
		AuthRouter.POST("login", auth.Login)
		AuthRouter.POST("register", auth.Register)
		AuthRouter.GET("captcha", auth.GetCaptcha)
		AuthRouter.POST("send-verify-code", auth.SendVerifyCode) // 发送登录验证码
		AuthRouter.POST("forgot-password", auth.ForgotPassword)
		AuthRouter.POST("reset-password", auth.ResetPassword)
		AuthRouter.POST("resend-verification", auth.ResendVerification)   // 重发激活邮件
		AuthRouter.GET("verify-email", auth.VerifyEmail)                   // 验证邮箱激活
		AuthRouter.GET("verify-reset-token", auth.VerifyResetToken)        // 验证密码重置token
		AuthRouter.POST("logout", middleware.RequireAuth(authModel.AuthLevelUser), auth.Logout)
	}
}

// InitPaymentRouter 支付回调路由
func InitPaymentRouter(Router *gin.RouterGroup) {
	PaymentGroup := Router.Group("v1/payment")
	{
		// 支付回调接口(不需要认证)
		PaymentGroup.POST("/alipay/notify", payment.AlipayNotify)
		PaymentGroup.POST("/wechat/notify", payment.WechatNotify)
		PaymentGroup.POST("/epay/notify", payment.EpayNotify)
		PaymentGroup.GET("/epay/notify", payment.EpayNotify) // 支持GET请求
		PaymentGroup.POST("/mapay/notify", payment.MapayNotify)

		// 易支付接口
		PaymentGroup.POST("/yipay/notify", payment.YipayNotify)
		PaymentGroup.GET("/yipay/notify", payment.YipayNotify) // 支持GET请求
		PaymentGroup.GET("/yipay/return", payment.YipayReturn)
	}

	// 易支付管理接口(需要认证)
	YipayAdminGroup := Router.Group("v1/admin/payment/yipay")
	YipayAdminGroup.Use(middleware.RequireAuth(authModel.AuthLevelAdmin))
	{
		YipayAdminGroup.GET("/config", payment.GetYipayConfig)
		YipayAdminGroup.PUT("/config", payment.UpdateYipayConfig)
		YipayAdminGroup.GET("/stats", payment.GetYipayStats)
		YipayAdminGroup.GET("/orders", payment.GetYipayOrders)
		YipayAdminGroup.POST("/test", payment.TestYipay)
		YipayAdminGroup.POST("/create", payment.CreateYipayOrder)
		YipayAdminGroup.GET("/query", payment.QueryYipayOrder)
	}
}
