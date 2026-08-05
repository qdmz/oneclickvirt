package user

import (
	"strconv"

	"oneclickvirt/middleware"
	"oneclickvirt/model/common"
	ticketModel "oneclickvirt/model/ticket"
	ticketService "oneclickvirt/service/ticket"

	"github.com/gin-gonic/gin"
)

// CreateTicket 用户创建工单
func CreateTicket(c *gin.Context) {
	var req ticketModel.CreateTicketRequest
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
	t, err := svc.CreateTicket(uid, req)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, t, "工单创建成功")
}

// ListMyTickets 用户查看自己的工单列表
func ListMyTickets(c *gin.Context) {
	var req ticketModel.TicketListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	uid, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, "未授权"))
		return
	}
	svc := ticketService.NewService()
	list, total, err := svc.ListUserTickets(uid, req)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccessWithPagination(c, list, total, req.Page, req.PageSize)
}

// GetMyTicket 用户查看自己的工单详情
func GetMyTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的工单ID"))
		return
	}
	uid, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, "未授权"))
		return
	}
	svc := ticketService.NewService()
	t, err := svc.GetTicketDetail(uint(id), uid, false)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, t, "获取成功")
}

// ReplyMyTicket 用户回复自己的工单
func ReplyMyTicket(c *gin.Context) {
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
	if err := svc.ReplyTicket(uint(id), uid, false, req.Content); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, nil, "回复成功")
}

// CloseMyTicket 用户关闭自己的工单
func CloseMyTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的工单ID"))
		return
	}
	uid, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, "未授权"))
		return
	}
	svc := ticketService.NewService()
	if err := svc.CloseTicket(uint(id), uid, false); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, nil, "工单已关闭")
}
