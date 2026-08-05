package admin

import (
	"strconv"

	"oneclickvirt/middleware"
	"oneclickvirt/model/common"
	ticketModel "oneclickvirt/model/ticket"
	ticketService "oneclickvirt/service/ticket"

	"github.com/gin-gonic/gin"
)

// GetTickets 管理员获取全部工单列表
func GetTickets(c *gin.Context) {
	var req ticketModel.TicketListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	svc := ticketService.NewService()
	list, total, err := svc.ListAllTickets(req)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccessWithPagination(c, list, total, req.Page, req.PageSize)
}

// GetTicketDetail 管理员获取工单详情
func GetTicketDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的工单ID"))
		return
	}
	svc := ticketService.NewService()
	t, err := svc.GetTicketDetail(uint(id), 0, true)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, t, "获取成功")
}

// ReplyTicket 管理员回复工单
func ReplyTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的工单ID"))
		return
	}
	var req ticketModel.ReplyTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	uid, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, "未授权"))
		return
	}
	svc := ticketService.NewService()
	if err := svc.ReplyTicket(uint(id), uid, true, req.Content); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, nil, "回复成功")
}

// UpdateTicketStatus 管理员更新工单状态
func UpdateTicketStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的工单ID"))
		return
	}
	var req ticketModel.UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	svc := ticketService.NewService()
	if err := svc.UpdateStatus(uint(id), req.Status); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, nil, "状态更新成功")
}

// CloseTicket 管理员关闭工单
func CloseTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的工单ID"))
		return
	}
	svc := ticketService.NewService()
	if err := svc.CloseTicket(uint(id), 0, true); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, nil, "工单已关闭")
}

// DeleteTicket 管理员删除工单
func DeleteTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的工单ID"))
		return
	}
	svc := ticketService.NewService()
	if err := svc.DeleteTicket(uint(id)); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, nil, "工单已删除")
}

// GetTicketStats 管理员获取工单状态统计
func GetTicketStats(c *gin.Context) {
	svc := ticketService.NewService()
	stats, err := svc.CountByStatus()
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, stats, "获取成功")
}
