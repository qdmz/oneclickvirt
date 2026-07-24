package product

import (
	"strconv"

	"oneclickvirt/middleware"
	"oneclickvirt/model/common"
	"oneclickvirt/service/product"

	"github.com/gin-gonic/gin"
)

// CreateTicketRequest 创建工单请求
type CreateTicketRequest struct {
	Title      string `json:"title" binding:"required,max=256"`
	Content    string `json:"content" binding:"required"`
	Category   string `json:"category"`
	Priority   string `json:"priority"`
	InstanceID uint   `json:"instanceId"`
	OrderID    uint   `json:"orderId"`
}

// ReplyTicketRequest 回复工单请求
type ReplyTicketRequest struct {
	Content string `json:"content" binding:"required"`
}

// getUserID 从上下文获取用户ID
func getUserID(c *gin.Context) (uint, error) {
	return middleware.GetUserIDFromContext(c)
}

// CreateTicket 创建工单
// @Summary 创建工单
// @Description 用户创建新的支持工单
// @Tags 用户工单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTicketRequest true "创建工单请求参数"
// @Success 200 {object} common.Response{data=object} "创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "创建失败"
// @Router /user/tickets [post]
func CreateTicket(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := product.NewTicketService()
	ticket, err := service.CreateTicket(userID, req.Title, req.Content, req.Category, req.Priority, req.InstanceID, req.OrderID)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, ticket, "工单创建成功")
}

// GetTicketList 获取用户工单列表
// @Summary 获取用户工单列表
// @Description 获取当前登录用户的工单列表，支持分页和状态筛选
// @Tags 用户工单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param status query int false "状态筛选 0=待处理 1=处理中 2=已解决 3=已关闭"
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "获取失败"
// @Router /user/tickets [get]
func GetTicketList(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var pageInfo common.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	pageInfo.Normalize(common.DefaultPageSize)

	// 状态筛选
	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil && s >= 0 && s <= 3 {
			status = &s
		}
	}

	service := product.NewTicketService()
	tickets, total, err := service.GetUserTicketList(userID, status, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccessWithPagination(c, tickets, total, pageInfo.Page, pageInfo.PageSize)
}

// GetTicketDetail 获取工单详情
// @Summary 获取工单详情
// @Description 获取指定工单的详细信息，包含回复列表
// @Tags 用户工单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单ID"
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 403 {object} common.Response "无权查看"
// @Failure 404 {object} common.Response "工单不存在"
// @Router /user/tickets/{id} [get]
func GetTicketDetail(c *gin.Context) {
	userID, err := getUserID(c)
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

	service := product.NewTicketService()
	ticket, replies, err := service.GetTicketDetail(userID, uint(ticketID))
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, gin.H{
		"ticket":  ticket,
		"replies": replies,
	})
}

// ReplyTicket 回复工单
// @Summary 回复工单
// @Description 用户对指定工单进行回复
// @Tags 用户工单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单ID"
// @Param request body ReplyTicketRequest true "回复请求参数"
// @Success 200 {object} common.Response{data=object} "回复成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 403 {object} common.Response "无权回复"
// @Failure 500 {object} common.Response "回复失败"
// @Router /user/tickets/{id}/reply [post]
func ReplyTicket(c *gin.Context) {
	userID, err := getUserID(c)
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

	var req ReplyTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := product.NewTicketService()
	reply, err := service.ReplyTicket(userID, uint(ticketID), req.Content)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, reply, "回复成功")
}

// CloseTicket 关闭工单
// @Summary 关闭工单
// @Description 用户关闭指定工单
// @Tags 用户工单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单ID"
// @Success 200 {object} common.Response "关闭成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 403 {object} common.Response "无权关闭"
// @Failure 500 {object} common.Response "关闭失败"
// @Router /user/tickets/{id}/close [post]
func CloseTicket(c *gin.Context) {
	userID, err := getUserID(c)
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

	service := product.NewTicketService()
	if err := service.CloseTicket(userID, uint(ticketID)); err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, nil, "工单已关闭")
}
