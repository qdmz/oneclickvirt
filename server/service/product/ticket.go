package product

import (
	"errors"
	"fmt"
	"time"

	"oneclickvirt/global"
	"oneclickvirt/model/product"
	userModel "oneclickvirt/model/user"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TicketService 工单服务
type TicketService struct{}

// NewTicketService 创建工单服务实例
func NewTicketService() *TicketService {
	return &TicketService{}
}

// GenerateTicketNo 生成工单编号
func (s *TicketService) GenerateTicketNo() string {
	return fmt.Sprintf("TK%s%s", time.Now().Format("20060102"), uuid.New().String()[:8])
}

// CreateTicket 创建工单
func (s *TicketService) CreateTicket(userID uint, title, content, category, priority string, instanceID, orderID uint) (*product.Ticket, error) {
	if title == "" || content == "" {
		return nil, errors.New("标题和内容不能为空")
	}

	// 验证分类
	validCategories := map[string]bool{"general": true, "technical": true, "billing": true, "complaint": true}
	if category == "" {
		category = "general"
	}
	if !validCategories[category] {
		return nil, errors.New("无效的工单分类")
	}

	// 验证优先级
	validPriorities := map[string]bool{"low": true, "normal": true, "high": true, "urgent": true}
	if priority == "" {
		priority = "normal"
	}
	if !validPriorities[priority] {
		return nil, errors.New("无效的优先级")
	}

	ticket := &product.Ticket{
		TicketNo:   s.GenerateTicketNo(),
		UserID:     userID,
		Title:      title,
		Content:    content,
		Category:   category,
		Priority:   priority,
		Status:     0, // 待处理
		InstanceID: instanceID,
		OrderID:    orderID,
	}

	if err := global.APP_DB.Create(ticket).Error; err != nil {
		global.APP_LOG.Error("创建工单失败", zap.Uint("userID", userID), zap.Error(err))
		return nil, err
	}

	global.APP_LOG.Info("创建工单成功",
		zap.Uint("ticketID", ticket.ID),
		zap.String("ticketNo", ticket.TicketNo),
		zap.Uint("userID", userID))
	return ticket, nil
}

// GetUserTicketList 获取用户工单列表
func (s *TicketService) GetUserTicketList(userID uint, status *int, page, pageSize int) ([]product.Ticket, int64, error) {
	var tickets []product.Ticket
	var total int64

	query := global.APP_DB.Model(&product.Ticket{}).Where("user_id = ?", userID)
	if status != nil {
		query = query.Where("status = ?", *status)
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

// GetTicketDetail 获取工单详情（用户视角，仅允许查看自己的工单）
func (s *TicketService) GetTicketDetail(userID, ticketID uint) (*product.Ticket, []product.TicketReply, error) {
	var ticket product.Ticket
	if err := global.APP_DB.First(&ticket, ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("工单不存在")
		}
		return nil, nil, err
	}

	// 权限检查：只能查看自己的工单
	if ticket.UserID != userID {
		return nil, nil, errors.New("无权查看该工单")
	}

	var replies []product.TicketReply
	if err := global.APP_DB.Where("ticket_id = ? AND is_internal = ?", ticketID, false).
		Order("created_at ASC").Find(&replies).Error; err != nil {
		return nil, nil, err
	}

	return &ticket, replies, nil
}

// ReplyTicket 用户回复工单
func (s *TicketService) ReplyTicket(userID, ticketID uint, content string) (*product.TicketReply, error) {
	if content == "" {
		return nil, errors.New("回复内容不能为空")
	}

	var ticket product.Ticket
	if err := global.APP_DB.First(&ticket, ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工单不存在")
		}
		return nil, err
	}

	// 权限检查
	if ticket.UserID != userID {
		return nil, errors.New("无权回复该工单")
	}

	// 检查工单是否已关闭
	if ticket.Status == 3 {
		return nil, errors.New("工单已关闭，无法回复")
	}

	reply := &product.TicketReply{
		TicketID: ticketID,
		UserID:   userID,
		UserType: "user",
		Content:  content,
	}

	// 在事务中创建回复并更新工单最后回复信息
	if err := global.APP_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(reply).Error; err != nil {
			return err
		}
		now := time.Now()
		if err := tx.Model(&ticket).Updates(map[string]interface{}{
			"last_reply_at":      now,
			"last_reply_by":      "user",
			"last_reply_user_id": userID,
			"status":             0, // 回复后状态重置为待处理
		}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		global.APP_LOG.Error("回复工单失败", zap.Uint("ticketID", ticketID), zap.Uint("userID", userID), zap.Error(err))
		return nil, err
	}

	global.APP_LOG.Info("用户回复工单成功",
		zap.Uint("ticketID", ticketID),
		zap.Uint("replyID", reply.ID),
		zap.Uint("userID", userID))
	return reply, nil
}

// CloseTicket 关闭工单
func (s *TicketService) CloseTicket(userID, ticketID uint) error {
	var ticket product.Ticket
	if err := global.APP_DB.First(&ticket, ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("工单不存在")
		}
		return err
	}

	// 权限检查
	if ticket.UserID != userID {
		return errors.New("无权关闭该工单")
	}

	if ticket.Status == 3 {
		return errors.New("工单已关闭")
	}

	now := time.Now()
	if err := global.APP_DB.Model(&ticket).Updates(map[string]interface{}{
		"status":    3,
		"closed_at": now,
	}).Error; err != nil {
		return err
	}

	global.APP_LOG.Info("用户关闭工单成功",
		zap.Uint("ticketID", ticketID),
		zap.Uint("userID", userID))
	return nil
}

// ── 管理员接口 ──

// GetAdminTicketList 获取所有工单列表（管理员）
func (s *TicketService) GetAdminTicketList(status *int, category, priority string, page, pageSize int) ([]product.Ticket, int64, error) {
	var tickets []product.Ticket
	var total int64

	query := global.APP_DB.Model(&product.Ticket{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
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

// GetAdminTicketDetail 获取工单详情（管理员视角，包含内部备注）
func (s *TicketService) GetAdminTicketDetail(ticketID uint) (*product.Ticket, []product.TicketReply, *userModel.User, error) {
	var ticket product.Ticket
	if err := global.APP_DB.First(&ticket, ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, errors.New("工单不存在")
		}
		return nil, nil, nil, err
	}

	var replies []product.TicketReply
	if err := global.APP_DB.Where("ticket_id = ?", ticketID).
		Order("created_at ASC").Find(&replies).Error; err != nil {
		return nil, nil, nil, err
	}

	// 获取用户信息
	var user userModel.User
	if err := global.APP_DB.First(&user, ticket.UserID).Error; err != nil {
		user = userModel.User{ID: ticket.UserID, Username: "未知用户"}
	}

	return &ticket, replies, &user, nil
}

// AdminReplyTicket 管理员回复工单
func (s *TicketService) AdminReplyTicket(adminID, ticketID uint, content string, isInternal bool) (*product.TicketReply, error) {
	if content == "" {
		return nil, errors.New("回复内容不能为空")
	}

	var ticket product.Ticket
	if err := global.APP_DB.First(&ticket, ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工单不存在")
		}
		return nil, err
	}

	reply := &product.TicketReply{
		TicketID:   ticketID,
		UserID:     adminID,
		UserType:   "admin",
		Content:    content,
		IsInternal: isInternal,
	}

	// 在事务中创建回复并更新工单
	if err := global.APP_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(reply).Error; err != nil {
			return err
		}
		now := time.Now()
		updates := map[string]interface{}{
			"last_reply_at":      now,
			"last_reply_by":      "admin",
			"last_reply_user_id": adminID,
		}
		// 如果不是内部备注，将工单状态改为处理中
		if !isInternal && ticket.Status == 0 {
			updates["status"] = 1 // 处理中
		}
		if err := tx.Model(&ticket).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		global.APP_LOG.Error("管理员回复工单失败", zap.Uint("ticketID", ticketID), zap.Uint("adminID", adminID), zap.Error(err))
		return nil, err
	}

	global.APP_LOG.Info("管理员回复工单成功",
		zap.Uint("ticketID", ticketID),
		zap.Uint("replyID", reply.ID),
		zap.Uint("adminID", adminID),
		zap.Bool("isInternal", isInternal))
	return reply, nil
}

// UpdateTicketStatus 更新工单状态
func (s *TicketService) UpdateTicketStatus(ticketID uint, status int) error {
	if status < 0 || status > 3 {
		return errors.New("无效的工单状态")
	}

	var ticket product.Ticket
	if err := global.APP_DB.First(&ticket, ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("工单不存在")
		}
		return err
	}

	updates := map[string]interface{}{
		"status": status,
	}
	now := time.Now()
	switch status {
	case 2: // 已解决
		updates["solved_at"] = now
	case 3: // 已关闭
		updates["closed_at"] = now
	}

	if err := global.APP_DB.Model(&ticket).Updates(updates).Error; err != nil {
		return err
	}

	global.APP_LOG.Info("更新工单状态成功",
		zap.Uint("ticketID", ticketID),
		zap.Int("status", status))
	return nil
}

// AssignTicket 分配工单给管理员
func (s *TicketService) AssignTicket(ticketID, adminID uint) error {
	var ticket product.Ticket
	if err := global.APP_DB.First(&ticket, ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("工单不存在")
		}
		return err
	}

	// 验证管理员是否存在
	var admin userModel.User
	if err := global.APP_DB.First(&admin, adminID).Error; err != nil {
		return errors.New("管理员不存在")
	}
	if admin.UserType != "admin" && admin.UserType != "normal_admin" {
		return errors.New("指定用户不是管理员")
	}

	// 更新工单分配信息（通过 last_reply_user_id 记录分配的管理员）
	if err := global.APP_DB.Model(&ticket).Update("last_reply_user_id", adminID).Error; err != nil {
		return err
	}

	global.APP_LOG.Info("分配工单成功",
		zap.Uint("ticketID", ticketID),
		zap.Uint("adminID", adminID))
	return nil
}

// GetTicketStats 获取工单统计
func (s *TicketService) GetTicketStats() (map[string]interface{}, error) {
	var total, pending, processing, solved, closed int64

	if err := global.APP_DB.Model(&product.Ticket{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if err := global.APP_DB.Model(&product.Ticket{}).Where("status = ?", 0).Count(&pending).Error; err != nil {
		return nil, err
	}
	if err := global.APP_DB.Model(&product.Ticket{}).Where("status = ?", 1).Count(&processing).Error; err != nil {
		return nil, err
	}
	if err := global.APP_DB.Model(&product.Ticket{}).Where("status = ?", 2).Count(&solved).Error; err != nil {
		return nil, err
	}
	if err := global.APP_DB.Model(&product.Ticket{}).Where("status = ?", 3).Count(&closed).Error; err != nil {
		return nil, err
	}

	// 今日新增
	today := time.Now().Format("2006-01-02")
	var todayCount int64
	if err := global.APP_DB.Model(&product.Ticket{}).Where("DATE(created_at) = ?", today).Count(&todayCount).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":      total,
		"pending":    pending,
		"processing": processing,
		"solved":     solved,
		"closed":     closed,
		"today":      todayCount,
	}, nil
}

// TicketDetailResponse 工单详情响应（包含用户信息）
type TicketDetailResponse struct {
	Ticket  *product.Ticket        `json:"ticket"`
	Replies []product.TicketReply  `json:"replies"`
	User    *userModel.User        `json:"user,omitempty"`
}

// AdminTicketListResponse 管理员工单列表响应
type AdminTicketListResponse struct {
	product.Ticket
	Username string `json:"username"`
}

// GetAdminTicketListWithUser 获取管理员工单列表（包含用户名）
func (s *TicketService) GetAdminTicketListWithUser(status *int, category, priority string, page, pageSize int) ([]AdminTicketListResponse, int64, error) {
	var results []AdminTicketListResponse
	var total int64

	query := global.APP_DB.Model(&product.Ticket{}).
		Select("tickets.*, users.username as username").
		Joins("LEFT JOIN users ON tickets.user_id = users.id")

	if status != nil {
		query = query.Where("tickets.status = ?", *status)
	}
	if category != "" {
		query = query.Where("tickets.category = ?", category)
	}
	if priority != "" {
		query = query.Where("tickets.priority = ?", priority)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("tickets.created_at DESC").Offset(offset).Limit(pageSize).Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
