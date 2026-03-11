package agent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"oneclickvirt/global"
	"oneclickvirt/model/agent"

	"gorm.io/gorm"
)

type AgentService struct {
	DB *gorm.DB
}

func NewAgentService() *AgentService {
	return &AgentService{DB: global.APP_DB}
}

func (s *AgentService) CreateAgent(userID uint, req agent.CreateAgentRequest) (*agent.Agent, error) {
	code, err := s.generateCode()
	if err != nil {
		return nil, err
	}
	a := &agent.Agent{
		UserID:       userID,
		Code:         code,
		Name:         req.Name,
		ContactName:  req.ContactName,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		Status:       0,
	}
	if err := s.DB.Create(a).Error; err != nil {
		return nil, fmt.Errorf("创建代理商失败: %v", err)
	}
	return a, nil
}

func (s *AgentService) GetAgentByUserID(userID uint) (*agent.Agent, error) {
	var a agent.Agent
	if err := s.DB.Where("user_id = ?", userID).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AgentService) GetAgentByCode(code string) (*agent.Agent, error) {
	var a agent.Agent
	if err := s.DB.Where("code = ? AND status = 1", code).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AgentService) GetAgentByID(id uint) (*agent.Agent, error) {
	var a agent.Agent
	if err := s.DB.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AgentService) UpdateAgent(agentID uint, updates map[string]interface{}) error {
	return s.DB.Model(&agent.Agent{}).Where("id = ?", agentID).Updates(updates).Error
}

func (s *AgentService) ListSubUsers(agentID uint, page, pageSize int, keyword string, status *int) ([]map[string]interface{}, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	q := s.DB.Model(&agent.SubUserRelation{}).
		Where("agent_id = ?", agentID).
		Select("sub_user_relations.id, sub_user_relations.agent_id, sub_user_relations.user_id, sub_user_relations.created_at, users.username, users.nickname, users.email, users.status as user_status, users.created_at as user_created_at").
		Joins("LEFT JOIN users ON users.id = sub_user_relations.user_id")

	if keyword != "" {
		q = q.Where("users.username LIKE ? OR users.email LIKE ? OR users.nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		q = q.Where("users.status = ?", *status)
	}

	var total int64
	q.Count(&total)

	var results []map[string]interface{}
	err := q.Order("sub_user_relations.created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&results).Error
	return results, total, err
}

func (s *AgentService) DeleteSubUser(agentID uint, userID uint) error {
	return s.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).Delete(&agent.SubUserRelation{}).Error
}

func (s *AgentService) BatchUpdateSubUserStatus(agentID uint, userIDs []uint, status int) error {
	return s.DB.Exec("UPDATE users SET status = ? WHERE id IN ? AND id IN (SELECT user_id FROM sub_user_relations WHERE agent_id = ?)", status, userIDs, agentID).Error
}

func (s *AgentService) BatchDeleteSubUsers(agentID uint, userIDs []uint) error {
	return s.DB.Where("agent_id = ? AND user_id IN ?", agentID, userIDs).Delete(&agent.SubUserRelation{}).Error
}

func (s *AgentService) GetAgentStatistics(agentID uint) (*agent.AgentStatistics, error) {
	stats := &agent.AgentStatistics{}

	s.DB.Model(&agent.SubUserRelation{}).Where("agent_id = ?", agentID).Count(&stats.SubUserCount)

	s.DB.Model(&agent.SubUserRelation{}).
		Joins("LEFT JOIN users ON users.id = sub_user_relations.user_id").
		Where("agent_id = ? AND users.status = 1", agentID).
		Count(&stats.ActiveUserCount)

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	s.DB.Model(&agent.Commission{}).
		Where("agent_id = ? AND status = 1 AND created_at >= ?", agentID, monthStart).
		Select("COALESCE(SUM(amount), 0)").Scan(&stats.MonthCommission)

	s.DB.Model(&agent.Commission{}).
		Where("agent_id = ? AND status = 1", agentID).
		Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalCommission)

	s.DB.Model(&agent.Commission{}).
		Where("agent_id = ? AND status = 1 AND amount < 0", agentID).
		Select("COALESCE(SUM(ABS(amount)), 0)").Scan(&stats.TotalWithdrawn)

	var a agent.Agent
	s.DB.Select("balance").Where("id = ?", agentID).First(&a)
	stats.Balance = a.Balance

	return stats, nil
}

func (s *AgentService) ListCommissions(agentID uint, page, pageSize int, status *int) ([]agent.Commission, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.DB.Model(&agent.Commission{}).Where("agent_id = ?", agentID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var total int64
	query.Count(&total)

	var commissions []agent.Commission
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&commissions).Error
	return commissions, total, err
}

func (s *AgentService) SettleCommission(commissionID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var comm agent.Commission
		if err := tx.First(&comm, commissionID).Error; err != nil {
			return fmt.Errorf("佣金记录不存在")
		}
		if comm.Status != 0 {
			return fmt.Errorf("佣金状态不允许结算")
		}

		now := time.Now()
		comm.Status = 1
		comm.SettledAt = &now
		if err := tx.Save(&comm).Error; err != nil {
			return err
		}

		return tx.Model(&agent.Agent{}).Where("id = ?", comm.AgentID).
			Update("balance", gorm.Expr("balance + ?", comm.Amount)).Error
	})
}

func (s *AgentService) Withdraw(agentID uint, amount int64) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var a agent.Agent
		if err := tx.First(&a, agentID).Error; err != nil {
			return fmt.Errorf("代理商不存在")
		}
		if a.Balance < amount {
			return fmt.Errorf("余额不足")
		}

		now := time.Now()
		if err := tx.Model(&a).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		comm := agent.Commission{
			AgentID:     agentID,
			Amount:      -amount,
			Rate:        0,
			Status:      1,
			Description: fmt.Sprintf("提现 %.2f 元", float64(amount)/100),
			SettledAt:   &now,
		}
		return tx.Create(&comm).Error
	})
}

func (s *AgentService) ListAgents(page, pageSize int, keyword string, status *int) ([]agent.Agent, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.DB.Model(&agent.Agent{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR contact_email LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var total int64
	query.Count(&total)

	var agents []agent.Agent
	err := query.Preload("User").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&agents).Error
	return agents, total, err
}

func (s *AgentService) ApproveAgent(agentID uint) error {
	return s.DB.Model(&agent.Agent{}).Where("id = ? AND status = 0", agentID).Update("status", 1).Error
}

func (s *AgentService) DisableAgent(agentID uint) error {
	return s.DB.Model(&agent.Agent{}).Where("id = ?", agentID).Update("status", 2).Error
}

func (s *AgentService) AdjustCommission(agentID uint, rate float64) error {
	return s.DB.Model(&agent.Agent{}).Where("id = ?", agentID).Update("commission_rate", rate).Error
}

func (s *AgentService) CreateSubUserRelation(agentID uint, userID uint) error {
	var existing agent.SubUserRelation
	if err := s.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).First(&existing).Error; err == nil {
		return nil
	}
	relation := agent.SubUserRelation{
		AgentID: agentID,
		UserID:  userID,
	}
	return s.DB.Create(&relation).Error
}

func (s *AgentService) GetAgentSubUsers(agentID uint, page, pageSize int) ([]map[string]interface{}, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.DB.Model(&agent.SubUserRelation{}).
		Where("agent_id = ?", agentID).
		Select("sub_user_relations.id, sub_user_relations.agent_id, sub_user_relations.user_id, sub_user_relations.created_at, users.username, users.nickname, users.email, users.status as user_status").
		Joins("LEFT JOIN users ON users.id = sub_user_relations.user_id")

	var total int64
	query.Count(&total)

	var results []map[string]interface{}
	err := query.Order("sub_user_relations.created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&results).Error
	return results, total, err
}

func (s *AgentService) GetWalletTransactions(agentID uint, page, pageSize int) ([]agent.Commission, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.DB.Model(&agent.Commission{}).Where("agent_id = ?", agentID)
	var total int64
	query.Count(&total)

	var commissions []agent.Commission
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&commissions).Error
	return commissions, total, err
}

func (s *AgentService) generateCode() (string, error) {
	b := make([]byte, 8)
	for i := 0; i < 3; i++ {
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		code := "AG" + hex.EncodeToString(b)[:12]
		var count int64
		s.DB.Model(&agent.Agent{}).Where("code = ?", code).Count(&count)
		if count == 0 {
			return code, nil
		}
	}
	return "", fmt.Errorf("生成推广码失败")
}
