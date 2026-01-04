package resource

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ResourceReservation 资源预留记录模型 - 彻底简化设计
type ResourceReservation struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	UUID      string         `json:"uuid" gorm:"uniqueIndex;not null;size:36"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 核心标识字段
	UserID     uint   `json:"userId" gorm:"not null;index"`                  // 用户ID
	ProviderID uint   `json:"providerId" gorm:"not null;index"`              // Provider ID
	SessionID  string `json:"sessionId" gorm:"not null;size:64;uniqueIndex"` // 会话标识，唯一主键

	// 实例规格
	InstanceType string `json:"instanceType" gorm:"not null;size:16"` // container 或 vm
	CPU          int    `json:"cpu" gorm:"not null"`                  // 预留的CPU核心数
	Memory       int64  `json:"memory" gorm:"not null"`               // 预留的内存(MB)
	Disk         int64  `json:"disk" gorm:"not null"`                 // 预留的磁盘(MB)
	Bandwidth    int    `json:"bandwidth" gorm:"not null"`            // 预留的带宽(Mbps)

	// TTL管理
	ExpiresAt time.Time `json:"expiresAt" gorm:"index;column:expires_at"` // 预留过期时间，自动清理
}

func (r *ResourceReservation) BeforeCreate(tx *gorm.DB) error {
	r.UUID = uuid.New().String()
	return nil
}

// IsExpired 检查预留是否已过期
func (r *ResourceReservation) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// IsActive 检查预留是否处于活跃状态（未过期）
func (r *ResourceReservation) IsActive() bool {
	return !r.IsExpired()
}

// TableName 设置表名
func (ResourceReservation) TableName() string {
	return "resource_reservations"
}
