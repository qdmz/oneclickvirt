package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APIKey API 密钥模型
type APIKey struct {
	ID        uint           `json:"id" gorm:"primarykey;type:int(11)"`
	UUID      string         `json:"uuid" gorm:"uniqueIndex;not null;size:36"`
	UserID    uint           `json:"userId" gorm:"not null;index:idx_user_id;type:int(11)"`
	Name      string         `json:"name" gorm:"size:100;not null"` // 密钥名称，便于识别
	Key       string         `json:"key" gorm:"uniqueIndex;not null;size:64;column:api_key"` // API 密钥
	ExpiresAt *time.Time     `json:"expiresAt"` // 过期时间，nil 表示永不过期
	LastUsedAt *time.Time    `json:"lastUsedAt"` // 最后使用时间
	Status    string         `json:"status" gorm:"default:active;size:20"` // 状态：active, disabled
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联用户
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (a *APIKey) BeforeCreate(tx *gorm.DB) error {
	a.UUID = uuid.New().String()
	return nil
}
