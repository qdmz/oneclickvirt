package router

import (
	"oneclickvirt/api/v1/public"
	"oneclickvirt/api/v1/system"
	"oneclickvirt/api/v1/user"

	"github.com/gin-gonic/gin"
)

// InitPublicRouter 公开路由（需要数据库连接）
func InitPublicRouter(Router *gin.RouterGroup) {
	PublicRouter := Router.Group("v1/public")
	{
		PublicRouter.GET("announcements", system.GetAnnouncement)
		PublicRouter.GET("stats", public.GetDashboardStats)
		PublicRouter.GET("system-images/available", system.GetAvailableSystemImages)
		PublicRouter.GET("products", public.GetPublicProducts)
		PublicRouter.GET("payment-config", public.GetPaymentConfig)

		// KYC callback
		PublicRouter.GET("kyc/callback", user.HandleKYCCallback)
	}
}
