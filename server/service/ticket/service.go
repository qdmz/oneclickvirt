package ticket

import (
	"errors"
	"oneclickvirt/global"
	"oneclickvirt/model/ticket"
	"time"

	"gorm.io/gorm"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateTicket(userID uint, title, description string, ticketType ticket.TicketType, priority ticket.TicketPriority, instanceID *uint) (*ticket.Ticket, error) {
	if title == "" || description == "" {
		return nil, errors.New("标题和描述不能为空")
	}

	t := &ticket.Ticket{
		UserID:      userID,
		Title:       title,
		Description: description,
		Type:        ticketType,
		Priority:    priority,
		Status:      ticket.TicketStatusOpen,
		InstanceID:  instanceID,
	}

	if err := global.APP_DB.Create(t).Error; err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) GetTicketByID(ticketID, userID uint, isAdmin bool) (*ticket.Ticket, error) {
	var t ticket.Ticket
	query := global.APP_DB.Where("id = ?", ticketID)
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工单不存在")
		}
		return nil, err
	}

	return &t, nil
}

func (s *Service) GetUserTickets(userID uint, status ticket.TicketStatus, page, pageSize int) ([]ticket.Ticket, int64, error) {
	var tickets []ticket.Ticket
	var total int64

	query := global.APP_DB.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tickets).Error; err != nil {
		return nil, 0, err
	}

	return tickets, total, nil
}

func (s *Service) GetAllTickets(status ticket.TicketStatus, page, pageSize int) ([]ticket.Ticket, int64, error) {
	var tickets []ticket.Ticket
	var total int64

	query := global.APP_DB.Model(&ticket.Ticket{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tickets).Error; err != nil {
		return nil, 0, err
	}

	return tickets, total, nil
}

func (s *Service) UpdateTicket(ticketID uint, updates map[string]interface{}) error {
	if err := global.APP_DB.Model(&ticket.Ticket{}).Where("id = ?", ticketID).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) CloseTicket(ticketID uint, resolutionNotes string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":          ticket.TicketStatusClosed,
		"closed_at":       &now,
		"resolution_notes": resolutionNotes,
	}
	return s.UpdateTicket(ticketID, updates)
}

func (s *Service) AddReply(ticketID, userID uint, content string, isAdmin bool) (*ticket.TicketReply, error) {
	if content == "" {
		return nil, errors.New("回复内容不能为空")
	}

	reply := &ticket.TicketReply{
		TicketID: ticketID,
		UserID:   userID,
		Content:  content,
		IsAdmin:  isAdmin,
	}

	if err := global.APP_DB.Create(reply).Error; err != nil {
		return nil, err
	}

	if !isAdmin {
		global.APP_DB.Model(&ticket.Ticket{}).Where("id = ?", ticketID).Update("status", ticket.TicketStatusPending)
	}

	return reply, nil
}

func (s *Service) GetReplies(ticketID uint) ([]ticket.TicketReply, error) {
	var replies []ticket.TicketReply
	if err := global.APP_DB.Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&replies).Error; err != nil {
		return nil, err
	}
	return replies, nil
}