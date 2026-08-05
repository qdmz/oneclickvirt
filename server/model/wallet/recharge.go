package wallet

import (
	"time"

	"gorm.io/gorm"
)

// RechargeStatus 充值订单状态
const (
	RechargeStatusPending = "pending" // 待支付
	RechargeStatusPaid    = "paid"    // 已支付
	RechargeStatusFailed  = "failed"  // 支付失败/已关闭
)

// RechargeMethod 充值方式
const (
	RechargeMethodEpay  = "epay"  // 易支付
	RechargeMethodAdmin = "admin" // 管理员手动调整
)

// RechargeOrder 余额充值订单
type RechargeOrder struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	UserID     uint           `json:"userId" gorm:"index;not null"`
	TradeNo    string         `json:"tradeNo" gorm:"uniqueIndex;size:64;not null"` // 系统内部订单号（out_trade_no）
	Amount     int64          `json:"amount" gorm:"not null"`                      // 充值金额（分）
	Method     string         `json:"method" gorm:"size:16;not null;default:epay"` // 充值方式
	Status     string         `json:"status" gorm:"size:16;not null;default:pending;index"`
	EpayTradeNo string        `json:"epayTradeNo" gorm:"size:64"`                 // 易支付返回的交易号
	Remark     string         `json:"remark" gorm:"size:255"`
	PaidAt     *time.Time     `json:"paidAt"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (RechargeOrder) TableName() string {
	return "recharge_orders"
}
