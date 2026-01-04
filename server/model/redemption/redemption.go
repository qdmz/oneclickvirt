package redemption

import (
	"time"
	productModel "oneclickvirt/model/product"
	userModel "oneclickvirt/model/user"
)

// RedemptionCode 兑换码表
type RedemptionCode struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Code      string     `json:"code" gorm:"type:varchar(32);uniqueIndex;not null;comment:兑换码"`
	Type      string     `json:"type" gorm:"type:varchar(20);not null;comment:兑换码类型"`
	Amount    int64      `json:"amount" gorm:"default:0;comment:金额(分)或等级数"`
	ProductID *uint      `json:"productId" gorm:"index;comment:产品ID"`
	MaxUses   int        `json:"maxUses" gorm:"default:1;comment:最大使用次数"`
	UsedCount int        `json:"usedCount" gorm:"default:0;comment:已使用次数"`
	IsEnabled bool       `json:"isEnabled" gorm:"default:true;comment:是否启用"`
	ExpireAt  *time.Time `json:"expireAt" gorm:"comment:过期时间"`
	Remark    string     `json:"remark" gorm:"type:varchar(255);comment:备注"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`

	Product *productModel.Product `json:"product" gorm:"foreignKey:ProductID"`
}

// TableName 指定表名
func (RedemptionCode) TableName() string {
	return "redemption_codes"
}

// 兑换码类型常量
const (
	RedemptionTypeBalance = "balance" // 余额
	RedemptionTypeLevel   = "level"   // 等级
	RedemptionTypeProduct = "product" // 产品
)

// RedemptionCodeUsage 兑换码使用记录表
type RedemptionCodeUsage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CodeID    uint      `json:"codeId" gorm:"index;not null;comment:兑换码ID"`
	UserID    uint      `json:"userId" gorm:"index;not null;comment:用户ID"`
	Reward    string    `json:"reward" gorm:"type:json;comment:奖励详情"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`

	Code RedemptionCode `json:"code" gorm:"foreignKey:CodeID"`
	User userModel.User  `json:"user" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (RedemptionCodeUsage) TableName() string {
	return "redemption_code_usages"
}
