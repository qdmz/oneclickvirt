package kyc

import (
	"time"

	"gorm.io/gorm"
)

// KYCRecord real name verification record
type KYCRecord struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"uniqueIndex;not null;index" json:"userId"`
	CertifyID   string         `gorm:"size:64;index" json:"certifyId"`
	RealName    string         `gorm:"size:64" json:"realName,omitempty"`
	IDCardNumber string        `gorm:"size:128" json:"-"`
	IDCardHash  string         `gorm:"size:64;index" json:"idCardHash"`
	Status      int            `gorm:"default:0;index" json:"status"` // 0=pending 1=certified 2=rejected 3=expired
	CertifiedAt *time.Time     `json:"certifiedAt"`
	IDType      string         `gorm:"size:10;default:IDENTITY_CARD" json:"idType"`
	Gender      string         `gorm:"size:4" json:"gender,omitempty"`
	Province    string         `gorm:"size:32" json:"province,omitempty"`
	City        string         `gorm:"size:32" json:"city,omitempty"`
	Remark      string         `gorm:"size:255" json:"remark,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (KYCRecord) TableName() string {
	return "kyc_records"
}

// KYC status constants
const (
	KYCStatusPending   = 0
	KYCStatusCertified = 1
	KYCStatusRejected  = 2
	KYCStatusExpired   = 3
)

// MaskRealName masks a name, e.g. "Zhang San" -> "Zhang S*"
func MaskRealName(name string) string {
	runes := []rune(name)
	if len(runes) <= 1 {
		return name
	}
	if len(runes) == 2 {
		return string(runes[0]) + "*"
	}
	result := string(runes[0])
	for i := 1; i < len(runes)-1; i++ {
		result += "*"
	}
	result += string(runes[len(runes)-1])
	return result
}
