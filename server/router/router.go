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
		PaymentGroup.POST("/mapay/notify", payment.MapayNotify)
	}
}
