package admin

import (
	"strconv"

	"oneclickvirt/middleware"
	"oneclickvirt/model/common"
	"oneclickvirt/service/product"

	"github.com/gin-gonic/gin"
)

// AdminReplyTicketRequest 管理员回复工单请求
type AdminReplyTicketRequest struct {
	Content    string `json:"content" binding:"required"`
	IsInternal bool   `json:"isInternal"`
}

// UpdateTicketStatusRequest 更新工单状态请求
type UpdateTicketStatusRequest struct {
	Status int `json:"status" binding:"required,min=0,max=3"`
}

// AssignTicketRequest 分配工单请求
type AssignTicketRequest struct {
	AdminID uint `json:"adminId" binding:"required"`
}

// getAdminID 从上下文获取管理员ID
func getAdminID(c *gin.Context) (uint, error) {
	return middleware.GetUserIDFromContext(c)
}

// GetAdminTicketList 获取所有工单列表
// @Summary 获取所有工单列表
// @Description 管理员获取所有工单列表，支持按状态/分类/优先级筛选
// @Tags 管理员工单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param status query int false "状态 0=待处理 1=处理中 2=已解决 3=已关闭"
// @Param category query string false "分类 general/technical/billing/complaint"
// @Param priority query string false "优先级 low/normal/high/urgent"
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 401 {object} common.Response "未授权"
// @Failure 500 {object} common.Response "获取失败"
// @Router /admin/tickets [get]
func GetAdminTicketList(c *gin.Context) {
	var pageInfo common.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	pageInfo.Normalize(common.DefaultPageSize)

	// 筛选参数
	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil && s >= 0 && s <= 3 {
			status = &s
		}
	}
	category := c.Query("category")
	priority := c.Query("priority")

	service := product.NewTicketService()
	tickets, total, err := service.GetAdminTicketListWithUser(status, category, priority, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccessWithPagination(c, tickets, total, pageInfo.Page, pageInfo.PageSize)
}

// GetAdminTicketDetail 获取工单详情（管理员）
// @Summary 获取工单详情
// @Description 管理员获取指定工单的详细信息，包含所有回复（包括内部备注）
// @Tags 管理员工单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单ID"
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "工单不存在"
// @Router /admin/tickets/{id} [get]
func GetAdminTicketDetail(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "无效的工单ID"))
		return
	}

	service := product.NewTicketService()
	ticket, replies, user, err := service.GetAdminTicketDetail(uint(ticketID))
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, gin.H{
		"ticket":  ticket,
		"replies": replies,
		"user":    user,
	})
}

// AdminReplyTicket 管理员回复工单
// @Summary 管理员回复工单
// @Description 管理员对指定工单进行回复，支持内部备注
// @Tags 管理员工单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单ID"
// @Param request body AdminReplyTicketRequest true "回复请求参数"
// @Success 200 {object} common.Response{data=object} "回复成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "工单不存在"
// @Router /admin/tickets/{id}/reply [post]
func AdminReplyTicket(c *gin.Context) {
	adminID, err := getAdminID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "无效的工单ID"))
		return
	}

	var req AdminReplyTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := product.NewTicketService()
	reply, err := service.AdminReplyTicket(adminID, uint(ticketID), req.Content, req.IsInternal)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, reply, "回复成功")
}

// UpdateTicketStatus 更新工单状态
// @Summary 更新工单状态
// @Description 管理员更新指定工单的状态
// @Tags 管理员工单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单ID"
// @Param request body UpdateTicketStatusRequest true "更新状态请求参数"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "工单不存在"
// @Router /admin/tickets/{id}/status [put]
func UpdateTicketStatus(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "无效的工单ID"))
		return
	}

	var req UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := product.NewTicketService()
	if err := service.UpdateTicketStatus(uint(ticketID), req.Status); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "状态更新成功")
}

// AssignTicket 分配工单给管理员
// @Summary 分配工单
// @Description 将指定工单分配给某个管理员
// @Tags 管理员工单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单ID"
// @Param request body AssignTicketRequest true "分配请求参数"
// @Success 200 {object} common.Response "分配成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "工单或管理员不存在"
// @Router /admin/tickets/{id}/assign [post]
func AssignTicket(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "无效的工单ID"))
		return
	}

	var req AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := product.NewTicketService()
	if err := service.AssignTicket(uint(ticketID), req.AdminID); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "工单分配成功")
}

// GetTicketStats 获取工单统计
// @Summary 获取工单统计
// @Description 获取系统工单统计数据
// @Tags 管理员工单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 401 {object} common.Response "未授权"
// @Failure 500 {object} common.Response "获取失败"
// @Router /admin/tickets/stats [get]
func GetTicketStats(c *gin.Context) {
	service := product.NewTicketService()
	stats, err := service.GetTicketStats()
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, stats)
}
