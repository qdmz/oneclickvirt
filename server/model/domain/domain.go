package domain

import (
	"time"

	"gorm.io/gorm"
)

type Domain struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"userId"`
	InstanceID   uint           `gorm:"index" json:"instanceId"`
	Domain       string         `gorm:"size:255;not null" json:"domain"`
	Protocol     string         `gorm:"size:10;default:http" json:"protocol"`
	InternalIP   string         `gorm:"size:45;not null" json:"internalIp"`
	InternalPort int            `gorm:"not null" json:"internalPort"`
	ExternalPort int            `gorm:"default:0" json:"externalPort"`
	SSL          bool           `gorm:"default:false" json:"ssl"`
	Status       int            `gorm:"default:0;index" json:"status"`
	AgentID      *uint          `gorm:"index" json:"agentId,omitempty"`
	ExpiresAt    *time.Time     `json:"expiresAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type DomainConfig struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	MaxDomainsPerUser     int       `gorm:"default:3" json:"maxDomainsPerUser"`
	MaxDomainsPerAgentUser int      `gorm:"default:5" json:"maxDomainsPerAgentUser"`
	DefaultTTL            int       `gorm:"default:300" json:"defaultTTL"`
	AutoSSL               bool      `gorm:"default:false" json:"autoSSL"`
	AllowedSuffixes       string    `gorm:"type:text" json:"allowedSuffixes"`
	DNSType               string    `gorm:"size:20;default:dnsmasq" json:"dnsType"`
	DNSConfigPath         string    `gorm:"size:255;default:/etc/dnsmasq.d/oneclickvirt-hosts.conf" json:"dnsConfigPath"`
	NginxConfigPath       string    `gorm:"size:255;default:/etc/nginx/conf.d/oneclickvirt-domains" json:"nginxConfigPath"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}
