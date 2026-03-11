package agent

import (
	"time"

	"oneclickvirt/model/user"

	"gorm.io/gorm"
)

// Agent 代理商
type Agent struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"uniqueIndex;not null" json:"userId"`
	Code             string         `gorm:"uniqueIndex;size:32;not null" json:"code"`
	Name             string         `gorm:"size:100" json:"name"`
	ContactName      string         `gorm:"size:64" json:"contactName"`
	ContactEmail     string         `gorm:"size:128" json:"contactEmail"`
	ContactPhone     string         `gorm:"size:32" json:"contactPhone"`
	CommissionRate   float64        `gorm:"type:decimal(5,2);default:0" json:"commissionRate"`
	MaxSubUsers      int            `gorm:"default:0" json:"maxSubUsers"`
	MaxDomainsPerUser int           `gorm:"default:3" json:"maxDomainsPerUser"`
	Status           int            `gorm:"default:0;index" json:"status"` // 0=待审核 1=正常 2=禁用
	Balance          int64          `gorm:"default:0" json:"balance"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	User user.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// SubUserRelation 子用户关系
type SubUserRelation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   uint      `gorm:"index;not null" json:"agentId"`
	UserID    uint      `gorm:"index;not null" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`

	Agent Agent     `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	User  user.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Commission 佣金记录
type Commission struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	AgentID     uint       `gorm:"index;not null" json:"agentId"`
	OrderID     uint       `gorm:"index" json:"orderId"`
	SubUserID   uint       `gorm:"index" json:"subUserId"`
	Amount      int64      `gorm:"not null" json:"amount"`
	Rate        float64    `gorm:"type:decimal(5,2)" json:"rate"`
	Status      int        `gorm:"default:0" json:"status"` // 0=待结算 1=已结算 2=已取消
	Description string     `gorm:"size:255" json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	SettledAt   *time.Time `json:"settledAt"`
}

// 请求/响应结构体
type CreateAgentRequest struct {
	Name         string  `json:"name" binding:"required"`
	ContactName  string  `json:"contactName"`
	ContactEmail string  `json:"contactEmail" binding:"required,email"`
	ContactPhone string  `json:"contactPhone"`
}

type UpdateAgentRequest struct {
	Name         *string  `json:"name"`
	ContactName  *string  `json:"contactName"`
	ContactEmail *string  `json:"contactEmail"`
	ContactPhone *string  `json:"contactPhone"`
}

type AdminUpdateAgentRequest struct {
	Name               *string  `json:"name"`
	ContactName        *string  `json:"contactName"`
	ContactEmail       *string  `json:"contactEmail"`
	ContactPhone       *string  `json:"contactPhone"`
	CommissionRate     *float64 `json:"commissionRate"`
	MaxSubUsers        *int     `json:"maxSubUsers"`
	MaxDomainsPerUser  *int     `json:"maxDomainsPerUser"`
}

type WithdrawRequest struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

type AdjustCommissionRequest struct {
	CommissionRate float64 `json:"commissionRate" binding:"required,gte=0,lte=100"`
}

type ListSubUsersRequest struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Keyword  string `form:"keyword" json:"keyword"`
	Status   *int   `form:"status" json:"status"`
}

type ListAgentsRequest struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Keyword  string `form:"keyword" json:"keyword"`
	Status   *int   `form:"status" json:"status"`
}

type AgentStatistics struct {
	SubUserCount     int64   `json:"subUserCount"`
	ActiveUserCount  int64   `json:"activeUserCount"`
	MonthCommission  int64   `json:"monthCommission"`
	TotalCommission  int64   `json:"totalCommission"`
	TotalWithdrawn   int64   `json:"totalWithdrawn"`
	Balance          int64   `json:"balance"`
}

type WalletTransaction struct {
	ID          uint       `json:"id"`
	Amount      int64      `json:"amount"`
	Type        string     `json:"type"` // commission, withdraw
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
}
