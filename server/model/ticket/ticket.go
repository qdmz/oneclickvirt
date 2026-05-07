package ticket

import (
	"time"

	"gorm.io/gorm"
)

type TicketStatus string
type TicketPriority string
type TicketType string

const (
	TicketStatusOpen     TicketStatus = "open"
	TicketStatusPending  TicketStatus = "pending"
	TicketStatusResolved TicketStatus = "resolved"
	TicketStatusClosed   TicketStatus = "closed"

	TicketPriorityLow    TicketPriority = "low"
	TicketPriorityMedium TicketPriority = "medium"
	TicketPriorityHigh   TicketPriority = "high"
	TicketPriorityUrgent TicketPriority = "urgent"

	TicketTypeQuestion    TicketType = "question"
	TicketTypeIssue       TicketType = "issue"
	TicketTypeFeature     TicketType = "feature"
	TicketTypeComplaint   TicketType = "complaint"
	TicketTypeOther       TicketType = "other"
)

type Ticket struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"userId" gorm:"not null;index"`
	Title           string         `json:"title" gorm:"not null;size:255"`
	Description     string         `json:"description" gorm:"type:text"`
	Type            TicketType     `json:"type" gorm:"not null;default:'question'"`
	Priority        TicketPriority `json:"priority" gorm:"not null;default:'medium'"`
	Status          TicketStatus   `json:"status" gorm:"not null;default:'open';index"`
	AssignedTo      *uint          `json:"assignedTo" gorm:"index"`
	InstanceID      *uint          `json:"instanceId" gorm:"index"`
	Tags            string         `json:"tags" gorm:"size:255"`
	ResolutionNotes string         `json:"resolutionNotes" gorm:"type:text"`
	CreatedAt       time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	ClosedAt        *time.Time     `json:"closedAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Ticket) TableName() string {
	return "tickets"
}

type TicketReply struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	TicketID  uint           `json:"ticketId" gorm:"not null;index"`
	UserID    uint           `json:"userId" gorm:"not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	IsAdmin   bool           `json:"isAdmin" gorm:"default:false"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (TicketReply) TableName() string {
	return "ticket_replies"
}