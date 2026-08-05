package ticket

import (
	"time"

	"gorm.io/gorm"
)

// 工单状态常量
const (
	StatusOpen     = "open"     // 待处理（用户新建或用户回复后，等待管理员）
	StatusAnswered = "answered" // 已回复（管理员回复后，等待用户）
	StatusClosed   = "closed"   // 已关闭
)

// 工单优先级常量
const (
	PriorityLow    = "low"
	PriorityNormal = "normal"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"
)

// Ticket 工单模型
type Ticket struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	UserID      uint           `json:"userId" gorm:"index;not null"`               // 创建工单的用户
	Username    string         `json:"username" gorm:"-"`                           // 查询时填充，不入库
	Title       string         `json:"title" gorm:"not null;size:255"`             // 工单标题
	Category    string         `json:"category" gorm:"size:64;default:general"`    // 分类：general/billing/technical/other
	Priority    string         `json:"priority" gorm:"size:16;default:normal"`     // 优先级
	Status      string         `json:"status" gorm:"size:16;default:open;index"`   // 状态
	ReplyCount  int            `json:"replyCount" gorm:"default:0"`               // 回复总数
	LastReplyAt *time.Time     `json:"lastReplyAt"`                                // 最后回复时间
	LastReplyBy string         `json:"lastReplyBy" gorm:"size:16"`                // 最后回复方：user / admin
	ClosedAt    *time.Time     `json:"closedAt"`                                   // 关闭时间
	Replies     []TicketReply  `json:"replies,omitempty" gorm:"foreignKey:TicketID"`
}

// TableName 指定表名
func (Ticket) TableName() string {
	return "tickets"
}

// TicketReply 工单回复模型
type TicketReply struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	TicketID  uint           `json:"ticketId" gorm:"index;not null"` // 所属工单
	UserID    uint           `json:"userId" gorm:"not null"`         // 回复者用户ID
	Username  string         `json:"username" gorm:"-"`               // 查询时填充，不入库
	IsAdmin   bool           `json:"isAdmin" gorm:"default:false"`   // 是否管理员回复
	Content   string         `json:"content" gorm:"type:longtext"`   // 回复内容
}

// TableName 指定表名
func (TicketReply) TableName() string {
	return "ticket_replies"
}
