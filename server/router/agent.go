package router

import (
	"oneclickvirt/api/v1/agent"
	"oneclickvirt/middleware"
	authModel "oneclickvirt/model/auth"

	"github.com/gin-gonic/gin"
)

// InitAgentRouter 代理商路由
func InitAgentRouter(Router *gin.RouterGroup) {
	AgentGroup := Router.Group("/v1/agent")
	AgentGroup.Use(middleware.RequireAuth(authModel.AuthLevelUser))
	{
		// 代理商申请（已认证用户即可调用）
		AgentGroup.POST("/apply", agent.CreateAgent)
		AgentGroup.GET("/profile", agent.GetProfile)

		// 以下接口需要代理商权限
		agentAuth := AgentGroup.Group("")
		agentAuth.Use(middleware.RequireAgent())
		{
			agentAuth.PUT("/profile", agent.UpdateProfile)

			// 子用户管理
			agentAuth.GET("/sub-users", agent.GetSubUsers)
			agentAuth.DELETE("/sub-users/:userId", agent.DeleteSubUser)
			agentAuth.PUT("/sub-users/batch-status", agent.BatchUpdateSubUserStatus)
			agentAuth.POST("/sub-users/batch-delete", agent.BatchDeleteSubUsers)

			// 统计
			agentAuth.GET("/statistics", agent.GetStatistics)

			// 佣金
			agentAuth.GET("/commissions", agent.GetCommissions)

			// 钱包
			agentAuth.GET("/wallet", agent.GetWallet)
			agentAuth.GET("/wallet/transactions", agent.GetWalletTransactions)
			agentAuth.POST("/wallet/withdraw", agent.Withdraw)
		}
	}
}
