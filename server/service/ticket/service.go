package ticket

import (
	"context"
	"errors"
	"time"

	"oneclickvirt/global"
	ticketModel "oneclickvirt/model/ticket"
	userModel "oneclickvirt/model/user"
	"oneclickvirt/service/database"
	"oneclickvirt/utils"

	"gorm.io/gorm"
)

// Service 工单服务
type Service struct{}

// NewService 创建工单服务
func NewService() *Service {
	return &Service{}
}

// normalizePriority 校验并归一化优先级
func normalizePriority(p string) string {
	switch p {
	case ticketModel.PriorityLow, ticketModel.PriorityNormal, ticketModel.PriorityHigh, ticketModel.PriorityUrgent:
		return p
	default:
		return ticketModel.PriorityNormal
	}
}

// CreateTicket 用户创建工单（同时写入首条留言）
func (s *Service) CreateTicket(userID uint, req ticketModel.CreateTicketRequest) (*ticketModel.Ticket, error) {
	category := req.Category
	if category == "" {
		category = "general"
	}
	now := time.Now()
	t := ticketModel.Ticket{
		UserID:      userID,
		Title:       req.Title,
		Category:    category,
		Priority:    normalizePriority(req.Priority),
		Status:      ticketModel.StatusOpen,
		ReplyCount:  1,
		LastReplyAt: &now,
		LastReplyBy: "user",
	}

	dbService := database.GetDatabaseService()
	err := dbService.ExecuteTransaction(context.Background(), func(tx *gorm.DB) error {
		if err := tx.Create(&t).Error; err != nil {
			return err
		}
		reply := ticketModel.TicketReply{
			TicketID: t.ID,
			UserID:   userID,
			IsAdmin:  false,
			Content:  utils.SanitizeHTML(req.Content),
		}
		return tx.Create(&reply).Error
	})
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// listTickets 内部通用列表查询；userID>0 时按用户过滤（用户侧），=0 时不过滤（管理员侧）
func (s *Service) listTickets(userID uint, req ticketModel.TicketListRequest) ([]ticketModel.Ticket, int64, error) {
	var tickets []ticketModel.Ticket
	var total int64

	query := global.APP_DB.Model(&ticketModel.Ticket{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Priority != "" {
		query = query.Where("priority = ?", req.Priority)
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page, pageSize := req.Page, req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	// 排序：未关闭优先、最后回复时间倒序
	if err := query.Order("CASE WHEN status='closed' THEN 1 ELSE 0 END ASC").
		Order("last_reply_at DESC").
		Offset(offset).Limit(pageSize).Find(&tickets).Error; err != nil {
		return nil, 0, err
	}

	s.fillUsernames(tickets)
	return tickets, total, nil
}

// fillUsernames 批量填充工单创建者用户名（避免 N+1）
func (s *Service) fillUsernames(tickets []ticketModel.Ticket) {
	if len(tickets) == 0 {
		return
	}
	idSet := make(map[uint]struct{})
	for _, t := range tickets {
		idSet[t.UserID] = struct{}{}
	}
	ids := make([]uint, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	userMap := make(map[uint]string)
	var users []userModel.User
	if err := global.APP_DB.Select("id", "username").Where("id IN ?", ids).Find(&users).Error; err == nil {
		for _, u := range users {
			userMap[u.ID] = u.Username
		}
	}
	for i := range tickets {
		tickets[i].Username = userMap[tickets[i].UserID]
	}
}

// ListUserTickets 用户查看自己的工单列表
func (s *Service) ListUserTickets(userID uint, req ticketModel.TicketListRequest) ([]ticketModel.Ticket, int64, error) {
	return s.listTickets(userID, req)
}

// ListAllTickets 管理员查看所有工单列表
func (s *Service) ListAllTickets(req ticketModel.TicketListRequest) ([]ticketModel.Ticket, int64, error) {
	return s.listTickets(0, req)
}

// GetTicketDetail 获取工单详情（含回复）；requesterID>0 且非管理员时校验归属
func (s *Service) GetTicketDetail(ticketID uint, requesterID uint, isAdmin bool) (*ticketModel.Ticket, error) {
	var t ticketModel.Ticket
	if err := global.APP_DB.Where("id = ?", ticketID).First(&t).Error; err != nil {
		return nil, err
	}
	if !isAdmin && t.UserID != requesterID {
		return nil, errors.New("无权访问该工单")
	}

	var replies []ticketModel.TicketReply
	if err := global.APP_DB.Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&replies).Error; err != nil {
		return nil, err
	}

	// 填充回复者用户名
	idSet := make(map[uint]struct{})
	idSet[t.UserID] = struct{}{}
	for _, r := range replies {
		idSet[r.UserID] = struct{}{}
	}
	ids := make([]uint, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	userMap := make(map[uint]string)
	var users []userModel.User
	if err := global.APP_DB.Select("id", "username").Where("id IN ?", ids).Find(&users).Error; err == nil {
		for _, u := range users {
			userMap[u.ID] = u.Username
		}
	}
	t.Username = userMap[t.UserID]
	for i := range replies {
		replies[i].Username = userMap[replies[i].UserID]
	}
	t.Replies = replies
	return &t, nil
}

// ReplyTicket 回复工单；isAdmin 决定回复方与新状态
func (s *Service) ReplyTicket(ticketID uint, requesterID uint, isAdmin bool, content string) error {
	var t ticketModel.Ticket
	if err := global.APP_DB.Where("id = ?", ticketID).First(&t).Error; err != nil {
		return err
	}
	if !isAdmin && t.UserID != requesterID {
		return errors.New("无权回复该工单")
	}
	if t.Status == ticketModel.StatusClosed {
		return errors.New("工单已关闭，无法回复")
	}

	now := time.Now()
	newStatus := ticketModel.StatusOpen
	lastReplyBy := "user"
	if isAdmin {
		newStatus = ticketModel.StatusAnswered
		lastReplyBy = "admin"
	}

	dbService := database.GetDatabaseService()
	return dbService.ExecuteTransaction(context.Background(), func(tx *gorm.DB) error {
		reply := ticketModel.TicketReply{
			TicketID: ticketID,
			UserID:   requesterID,
			IsAdmin:  isAdmin,
			Content:  utils.SanitizeHTML(content),
		}
		if err := tx.Create(&reply).Error; err != nil {
			return err
		}
		return tx.Model(&ticketModel.Ticket{}).Where("id = ?", ticketID).Updates(map[string]interface{}{
			"status":        newStatus,
			"last_reply_at": &now,
			"last_reply_by": lastReplyBy,
			"reply_count":   gorm.Expr("reply_count + 1"),
		}).Error
	})
}

// CloseTicket 关闭工单；用户仅能关闭自己的工单，管理员可关闭任意
func (s *Service) CloseTicket(ticketID uint, requesterID uint, isAdmin bool) error {
	var t ticketModel.Ticket
	if err := global.APP_DB.Where("id = ?", ticketID).First(&t).Error; err != nil {
		return err
	}
	if !isAdmin && t.UserID != requesterID {
		return errors.New("无权操作该工单")
	}
	now := time.Now()
	return global.APP_DB.Model(&ticketModel.Ticket{}).Where("id = ?", ticketID).Updates(map[string]interface{}{
		"status":    ticketModel.StatusClosed,
		"closed_at": &now,
	}).Error
}

// UpdateStatus 管理员更新工单状态
func (s *Service) UpdateStatus(ticketID uint, status string) error {
	updates := map[string]interface{}{"status": status}
	if status == ticketModel.StatusClosed {
		now := time.Now()
		updates["closed_at"] = &now
	}
	return global.APP_DB.Model(&ticketModel.Ticket{}).Where("id = ?", ticketID).Updates(updates).Error
}

// DeleteTicket 管理员删除工单（软删除，连带回复）
func (s *Service) DeleteTicket(ticketID uint) error {
	dbService := database.GetDatabaseService()
	return dbService.ExecuteTransaction(context.Background(), func(tx *gorm.DB) error {
		if err := tx.Where("ticket_id = ?", ticketID).Delete(&ticketModel.TicketReply{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", ticketID).Delete(&ticketModel.Ticket{}).Error
	})
}

// CountByStatus 统计各状态工单数量（管理员概览用）
func (s *Service) CountByStatus() (map[string]int64, error) {
	type row struct {
		Status string
		Cnt    int64
	}
	var rows []row
	if err := global.APP_DB.Model(&ticketModel.Ticket{}).
		Select("status, count(*) as cnt").Group("status").Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := map[string]int64{
		ticketModel.StatusOpen:     0,
		ticketModel.StatusAnswered: 0,
		ticketModel.StatusClosed:   0,
	}
	for _, r := range rows {
		result[r.Status] = r.Cnt
	}
	return result, nil
}
