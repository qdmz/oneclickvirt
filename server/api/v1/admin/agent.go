package admin

import (
	"strconv"

	agentModel "oneclickvirt/model/agent"
	agentService "oneclickvirt/service/agent"
	"oneclickvirt/model/common"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetAgentList 获取代理商列表
func GetAgentList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	var statusPtr *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		statusPtr = &v
	}

	svc := agentService.NewAgentService()
	agents, total, err := svc.ListAgents(page, pageSize, keyword, statusPtr)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccessWithPagination(c, agents, total, page, pageSize)
}

// CreateAgentByAdmin 管理员创建代理商
func CreateAgentByAdmin(c *gin.Context) {
	var req struct {
		UserID         uint    `json:"userId" binding:"required"`
		Name           string  `json:"name" binding:"required"`
		ContactName    string  `json:"contactName"`
		ContactEmail   string  `json:"contactEmail"`
		ContactPhone   string  `json:"contactPhone"`
		CommissionRate float64 `json:"commissionRate"`
		MaxSubUsers    int     `json:"maxSubUsers"`
		Status         int     `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	svc := agentService.NewAgentService()
	a, err := svc.CreateAgent(req.UserID, agentModel.CreateAgentRequest{
		Name:         req.Name,
		ContactName:  req.ContactName,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
	})
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	if req.CommissionRate > 0 {
		svc.AdjustCommission(a.ID, req.CommissionRate)
	}
	if req.MaxSubUsers > 0 {
		svc.UpdateAgent(a.ID, map[string]interface{}{"max_sub_users": req.MaxSubUsers})
	}
	if req.Status > 0 {
		svc.UpdateAgent(a.ID, map[string]interface{}{"status": req.Status})
	}

	common.ResponseSuccess(c, a, "代理商创建成功")
}

// UpdateAgentByAdmin 管理员更新代理商
func UpdateAgentByAdmin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req agentModel.AdminUpdateAgentRequest
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
	if req.CommissionRate != nil {
		updates["commission_rate"] = *req.CommissionRate
	}
	if req.MaxSubUsers != nil {
		updates["max_sub_users"] = *req.MaxSubUsers
	}
	if req.MaxDomainsPerUser != nil {
		updates["max_domains_per_user"] = *req.MaxDomainsPerUser
	}

	if err := svc.UpdateAgent(uint(id), updates); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "更新成功")
}

// DeleteAgent 删除代理商
func DeleteAgent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	svc := agentService.NewAgentService()
	if err := svc.UpdateAgent(uint(id), map[string]interface{}{"status": 2}); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}
	common.ResponseSuccess(c, nil, "已禁用代理商")
}

// ApproveAgent 审核通过
func ApproveAgent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	svc := agentService.NewAgentService()
	if err := svc.ApproveAgent(uint(id)); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}
	common.ResponseSuccess(c, nil, "审核通过")
}

// UpdateAgentStatus 更新代理商状态
func UpdateAgentStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	svc := agentService.NewAgentService()
	if err := svc.UpdateAgent(uint(id), map[string]interface{}{"status": req.Status}); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}
	common.ResponseSuccess(c, nil, "状态更新成功")
}

// AdjustCommission 调整佣金比例
func AdjustCommission(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req agentModel.AdjustCommissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	svc := agentService.NewAgentService()
	if err := svc.AdjustCommission(uint(id), req.CommissionRate); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}
	common.ResponseSuccess(c, nil, "佣金比例已调整")
}

// GetAgentDetail 获取代理商详情
func GetAgentDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	svc := agentService.NewAgentService()
	a, err := svc.GetAgentByID(uint(id))
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeNotFound, "代理商不存在"))
		return
	}

	stats, _ := svc.GetAgentStatistics(uint(id))

	common.ResponseSuccess(c, gin.H{
		"agent":      a,
		"statistics": stats,
	})
}

// GetAgentSubUsers 获取代理商子用户列表
func GetAgentSubUsers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	svc := agentService.NewAgentService()
	results, total, err := svc.GetAgentSubUsers(uint(id), page, pageSize)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccessWithPagination(c, results, total, page, pageSize)
}

// SettleCommission 结算佣金
func SettleCommission(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	svc := agentService.NewAgentService()
	if err := svc.SettleCommission(uint(id)); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeError, err.Error()))
		return
	}
	common.ResponseSuccess(c, nil, "佣金已结算")
}

func init() {
	zap.L().Debug("admin agent handlers loaded")
}
