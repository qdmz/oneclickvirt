package agent

import (
	"strconv"

	"oneclickvirt/global"
	agentModel "oneclickvirt/model/agent"
	agentService "oneclickvirt/service/agent"
	"oneclickvirt/model/common"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetProfile 获取代理商资料
func GetProfile(c *gin.Context) {
	userID, _ := getUserID(c)
	svc := agentService.NewAgentService()

	a, err := svc.GetAgentByUserID(userID)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeNotFound, "代理商信息不存在"))
		return
	}
	common.ResponseSuccess(c, a)
}

// CreateAgent 创建代理商申请
func CreateAgent(c *gin.Context) {
	userID, _ := getUserID(c)
	var req agentModel.CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	svc := agentService.NewAgentService()

	// 检查是否已是代理商
	existing, _ := svc.GetAgentByUserID(userID)
	if existing != nil {
		common.ResponseWithError(c, common.NewError(common.CodeConflict, "您已提交过代理商申请"))
		return
	}

	a, err := svc.CreateAgent(userID, req)
	if err != nil {
		global.APP_LOG.Error("创建代理商失败", zap.Error(err))
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, a, "代理商申请已提交，请等待审核")
}

// UpdateProfile 更新代理商资料
func UpdateProfile(c *gin.Context) {
	agentID, _ := getAgentID(c)
	var req agentModel.UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	svc := agentService.NewAgentService()
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.ContactName != nil {
		updates["contact_name"] = *req.ContactName
	}
	if req.ContactEmail != nil {
		updates["contact_email"] = *req.ContactEmail
	}
	if req.ContactPhone != nil {
		updates["contact_phone"] = *req.ContactPhone
	}

	if err := svc.UpdateAgent(agentID, updates); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "更新成功")
}

// GetSubUsers 获取子用户列表
func GetSubUsers(c *gin.Context) {
	agentID, _ := getAgentID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	svc := agentService.NewAgentService()
	results, total, err := svc.ListSubUsers(agentID, page, pageSize, keyword, nil)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccessWithPagination(c, results, total, page, pageSize)
}

// DeleteSubUser 删除子用户
func DeleteSubUser(c *gin.Context) {
	agentID, _ := getAgentID(c)
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "无效的用户ID"))
		return
	}

	svc := agentService.NewAgentService()
	if err := svc.DeleteSubUser(agentID, uint(userID)); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "已移除子用户")
}

// BatchUpdateSubUserStatus 批量更新子用户状态
func BatchUpdateSubUserStatus(c *gin.Context) {
	agentID, _ := getAgentID(c)
	var req struct {
		UserIDs []uint `json:"userIds" binding:"required"`
		Status  int    `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	svc := agentService.NewAgentService()
	if err := svc.BatchUpdateSubUserStatus(agentID, req.UserIDs, req.Status); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "操作成功")
}

// BatchDeleteSubUsers 批量删除子用户
func BatchDeleteSubUsers(c *gin.Context) {
	agentID, _ := getAgentID(c)
	var req struct {
		UserIDs []uint `json:"userIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	svc := agentService.NewAgentService()
	if err := svc.BatchDeleteSubUsers(agentID, req.UserIDs); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "批量删除成功")
}

// GetStatistics 获取代理商统计
func GetStatistics(c *gin.Context) {
	agentID, _ := getAgentID(c)
	svc := agentService.NewAgentService()
	stats, err := svc.GetAgentStatistics(agentID)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}
	common.ResponseSuccess(c, stats)
}

// GetCommissions 获取佣金记录
func GetCommissions(c *gin.Context) {
	agentID, _ := getAgentID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	var statusPtr *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		statusPtr = &v
	}

	svc := agentService.NewAgentService()
	commissions, total, err := svc.ListCommissions(agentID, page, pageSize, statusPtr)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccessWithPagination(c, commissions, total, page, pageSize)
}

// Withdraw 提现
func Withdraw(c *gin.Context) {
	agentID, _ := getAgentID(c)
	var req agentModel.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	svc := agentService.NewAgentService()
	if err := svc.Withdraw(agentID, req.Amount); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "提现申请已提交")
}

// GetWallet 获取钱包信息
func GetWallet(c *gin.Context) {
	agentID, _ := getAgentID(c)
	svc := agentService.NewAgentService()
	stats, err := svc.GetAgentStatistics(agentID)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}
	common.ResponseSuccess(c, gin.H{
		"balance":         stats.Balance,
		"totalCommission": stats.TotalCommission,
		"totalWithdrawn":  stats.TotalWithdrawn,
	})
}

// GetWalletTransactions 获取钱包交易记录
func GetWalletTransactions(c *gin.Context) {
	agentID, _ := getAgentID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	svc := agentService.NewAgentService()
	commissions, total, err := svc.GetWalletTransactions(agentID, page, pageSize)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccessWithPagination(c, commissions, total, page, pageSize)
}

// helpers
func getUserID(c *gin.Context) (uint, error) {
	return getUintFromContext(c, "user_id")
}

func getAgentID(c *gin.Context) (uint, error) {
	return getUintFromContext(c, "agent_id")
}

func getUintFromContext(c *gin.Context, key string) (uint, error) {
	if val, exists := c.Get(key); exists {
		if id, ok := val.(uint); ok {
			return id, nil
		}
	}
	return 0, common.NewError(common.CodeUnauthorized, "认证信息缺失")
}

func init() {
	global.APP_LOG.Debug("agent API handlers loaded")
}
