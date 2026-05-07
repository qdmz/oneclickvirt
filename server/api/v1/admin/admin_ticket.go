package admin

import (
	"fmt"
	"net/http"
	ticketModel "oneclickvirt/model/ticket"
	ticketService "oneclickvirt/service/ticket"

	"github.com/gin-gonic/gin"
)

type AdminReplyRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateTicketRequest struct {
	Status     ticketModel.TicketStatus   `json:"status"`
	Priority   ticketModel.TicketPriority `json:"priority"`
	AssignedTo *uint                      `json:"assignedTo"`
	Tags       string                     `json:"tags"`
}

func GetAllTickets(c *gin.Context) {
	status := c.Query("status")
	page := 1
	pageSize := 10
	c.Query("page")
	c.Query("pageSize")

	svc := ticketService.NewService()
	tickets, total, err := svc.GetAllTickets(ticketModel.TicketStatus(status), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets, "total": total})
}

func GetTicketDetail(c *gin.Context) {
	ticketID := c.Param("id")
	var tid uint
	if _, err := fmt.Sscanf(ticketID, "%d", &tid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID"})
		return
	}

	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "管理员未登录"})
		return
	}

	svc := ticketService.NewService()
	t, err := svc.GetTicketByID(tid, adminID.(uint), true)
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

func UpdateTicket(c *gin.Context) {
	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticketID := c.Param("id")
	var tid uint
	if _, err := fmt.Sscanf(ticketID, "%d", &tid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID"})
		return
	}

	updates := make(map[string]interface{})
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Priority != "" {
		updates["priority"] = req.Priority
	}
	if req.AssignedTo != nil {
		updates["assigned_to"] = req.AssignedTo
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}

	svc := ticketService.NewService()
	if err := svc.UpdateTicket(tid, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工单更新成功"})
}

func AdminReplyTicket(c *gin.Context) {
	var req AdminReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticketID := c.Param("id")
	var tid uint
	if _, err := fmt.Sscanf(ticketID, "%d", &tid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID"})
		return
	}

	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "管理员未登录"})
		return
	}

	svc := ticketService.NewService()
	reply, err := svc.AddReply(tid, adminID.(uint), req.Content, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reply})
}

func CloseTicket(c *gin.Context) {
	var req struct {
		ResolutionNotes string `json:"resolutionNotes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticketID := c.Param("id")
	var tid uint
	if _, err := fmt.Sscanf(ticketID, "%d", &tid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID"})
		return
	}

	svc := ticketService.NewService()
	if err := svc.CloseTicket(tid, req.ResolutionNotes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工单已关闭"})
}