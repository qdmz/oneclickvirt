package order

import (
	"time"

	"gorm.io/gorm"
)

// OrderStatus 订单状态
const (
	OrderStatusPending   = "pending"   // 待支付
	OrderStatusPaid      = "paid"      // 已支付（已开通）
	OrderStatusFailed    = "failed"    // 支付失败
	OrderStatusCancelled = "cancelled" // 已取消
)

// OrderMethod 支付方式
const (
	OrderMethodBalance = "balance" // 余额支付
	OrderMethodEpay    = "epay"    // 易支付
)

// Order 产品购买订单（在线购买 VM/容器产品）
type Order struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"userId" gorm:"index;not null"`
	ProductID   uint           `json:"productId" gorm:"index;not null"`
	TradeNo     string         `json:"tradeNo" gorm:"uniqueIndex;size:64;not null"` // 系统内部订单号
	Amount      int64          `json:"amount" gorm:"not null"`                      // 实付金额（分）
	Method      string         `json:"method" gorm:"size:16;not null"`             // balance / epay
	Status      string         `json:"status" gorm:"size:16;not null;default:pending;index"`
	EpayTradeNo string         `json:"epayTradeNo" gorm:"size:64"`
	// InstanceID 开通成功后关联的实例 ID
	InstanceID *uint `json:"instanceId"`
	// RenewalOf 续费所针对的已有实例/订单 ID（0 表示新购）
	RenewalOf *uint   `json:"renewalOf"`
	PeriodDays int    `json:"periodDays" gorm:"not null;default:30"`
	Remark     string `json:"remark" gorm:"size:255"`
	PaidAt     *time.Time `json:"paidAt"`
	// ExpiresAt 资源到期时间（购买成功后计算）
	ExpiresAt  *time.Time     `json:"expiresAt"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Order) TableName() string {
	return "orders"
}
