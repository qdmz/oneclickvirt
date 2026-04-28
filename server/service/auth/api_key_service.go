package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"oneclickvirt/global"
	"oneclickvirt/model/user"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// APIKeyService API 密钥服务
type APIKeyService struct{}

// NewAPIKeyService 创建 API 密钥服务
func NewAPIKeyService() *APIKeyService {
	return &APIKeyService{}
}

// CreateAPIKeyRequest 创建 API 密钥请求
type CreateAPIKeyRequest struct {
	Name      string     `json:"name" binding:"required,min=1,max=100"`
	ExpiresAt *time.Time `json:"expiresAt"`
}

// UpdateAPIKeyRequest 更新 API 密钥请求
type UpdateAPIKeyRequest struct {
	Name   string `json:"name" binding:"omitempty,min=1,max=100"`
	Status string `json:"status" binding:"omitempty,oneof=active disabled"`
}

// CreateAPIKey 创建 API 密钥
func (s *APIKeyService) CreateAPIKey(userID uint, req CreateAPIKeyRequest) (*user.APIKey, error) {
	// 生成随机密钥
	key, err := s.generateRandomKey()
	if err != nil {
		return nil, fmt.Errorf("生成密钥失败: %w", err)
	}

	// 检查密钥是否已存在
	var existingKey user.APIKey
	err = global.APP_DB.Where("api_key = ?", key).First(&existingKey).Error
	if err == nil {
		return nil, errors.New("密钥已存在，请重试")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("检查密钥失败: %w", err)
	}

	// 创建 API 密钥
	apiKey := &user.APIKey{
		UserID:    userID,
		Name:      req.Name,
		Key:       key,
		ExpiresAt: req.ExpiresAt,
		Status:    "active",
	}

	if err := global.APP_DB.Create(apiKey).Error; err != nil {
		return nil, fmt.Errorf("创建密钥失败: %w", err)
	}

	return apiKey, nil
}

// GetAPIKeys 获取用户的 API 密钥列表
func (s *APIKeyService) GetAPIKeys(userID uint) ([]user.APIKey, error) {
	var apiKeys []user.APIKey

	err := global.APP_DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&apiKeys).Error

	if err != nil {
		return nil, fmt.Errorf("获取密钥列表失败: %w", err)
	}

	return apiKeys, nil
}

// GetAPIKeyByID 根据 ID 获取 API 密钥
func (s *APIKeyService) GetAPIKeyByID(id uint, userID uint) (*user.APIKey, error) {
	var apiKey user.APIKey

	err := global.APP_DB.Where("id = ? AND user_id = ?", id, userID).
		First(&apiKey).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("密钥不存在")
		}
		return nil, fmt.Errorf("获取密钥失败: %w", err)
	}

	return &apiKey, nil
}

// UpdateAPIKey 更新 API 密钥
func (s *APIKeyService) UpdateAPIKey(id uint, userID uint, req UpdateAPIKeyRequest) (*user.APIKey, error) {
	apiKey, err := s.GetAPIKeyByID(id, userID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if len(updates) == 0 {
		return apiKey, nil
	}

	if err := global.APP_DB.Model(apiKey).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新密钥失败: %w", err)
	}

	// 重新获取更新后的数据
	if err := global.APP_DB.First(apiKey, id).Error; err != nil {
		return nil, fmt.Errorf("获取更新后的密钥失败: %w", err)
	}

	return apiKey, nil
}

// DeleteAPIKey 删除 API 密钥
func (s *APIKeyService) DeleteAPIKey(id uint, userID uint) error {
	result := global.APP_DB.Where("id = ? AND user_id = ?", id, userID).Delete(&user.APIKey{})
	if result.Error != nil {
		return fmt.Errorf("删除密钥失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("密钥不存在")
	}
	return nil
}

// ValidateAPIKey 验证 API 密钥
func (s *APIKeyService) ValidateAPIKey(key string) (*user.APIKey, error) {
	var apiKey user.APIKey

	err := global.APP_DB.Where("api_key = ? AND status = ?", key, "active").
		First(&apiKey).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("密钥无效或已禁用")
		}
		return nil, fmt.Errorf("验证密钥失败: %w", err)
	}

	// 检查是否过期
	if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
		return nil, errors.New("密钥已过期")
	}

	// 更新最后使用时间
	now := time.Now()
	apiKey.LastUsedAt = &now
	if err := global.APP_DB.Model(&apiKey).Update("last_used_at", now).Error; err != nil {
		// 更新失败不影响验证结果，只记录日志
		global.APP_LOG.Error("更新 API 密钥最后使用时间失败", zap.Error(err))
	}

	return &apiKey, nil
}

// RevokeAPIKey 撤销 API 密钥
func (s *APIKeyService) RevokeAPIKey(id uint, userID uint) error {
	_, err := s.UpdateAPIKey(id, userID, UpdateAPIKeyRequest{
		Status: "disabled",
	})
	return err
}

// generateRandomKey 生成随机密钥
func (s *APIKeyService) generateRandomKey() (string, error) {
	// 生成 32 字节的随机数据
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// 使用 Base64 编码，并移除填充字符
	key := base64.URLEncoding.EncodeToString(bytes)

	// 添加前缀以便识别
	return "ocv_" + key, nil
}
