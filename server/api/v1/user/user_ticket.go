package user

import (
	"fmt"
	"net/http"
	ticketModel "oneclickvirt/model/ticket"
	ticketService "oneclickvirt/service/ticket"

	"github.com/gin-gonic/gin"
)

type CreateTicketRequest struct {
	Title       string                     `json:"title" binding:"required"`
	Description string                     `json:"description" binding:"required"`
	Type        ticketModel.TicketType     `json:"type" binding:"required,oneof=question issue feature complaint other"`
	Priority    ticketModel.TicketPriority `json:"priority" binding:"required,oneof=low medium high urgent"`
	InstanceID  *uint                      `json:"instanceId"`
}

type ReplyRequest struct {
	Content string `json:"content" binding:"required"`
}

func CreateTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	svc := ticketService.NewService()
	t, err := svc.CreateTicket(userID.(uint), req.Title, req.Description, req.Type, req.Priority, req.InstanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": t})
}

func GetUserTickets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	status := c.Query("status")
	page := 1
	pageSize := 10
	c.Query("page")
	c.Query("pageSize")

	svc := ticketService.NewService()
	tickets, total, err := svc.GetUserTickets(userID.(uint), ticketModel.TicketStatus(status), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets, "total": total})
}

func GetTicketDetail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	ticketID := c.Param("id")
	var tid uint
	if _, err := fmt.Sscanf(ticketID, "%d", &tid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID"})
		return
	}

	svc := ticketService.NewService()
	t, err := svc.GetTicketByID(tid, userID.(uint), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	replies, err := svc.GetReplies(tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": t, "replies": replies})
}

func ReplyTicket(c *gin.Context) {
	var req ReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	ticketID := c.Param("id")
	var tid uint
	if _, err := fmt.Sscanf(ticketID, "%d", &tid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID"})
		return
	}

	svc := ticketService.NewService()
	reply, err := svc.AddReply(tid, userID.(uint), req.Content, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reply})
}