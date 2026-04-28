package order

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 订单状态常量
const (
	OrderStatusPending   = "pending"   // 待支付
	OrderStatusPaid      = "paid"      // 已支付
	OrderStatusCancelled = "cancelled" // 已取消
	OrderStatusExpired   = "expired"   // 已过期
	OrderStatusActive    = "active"    // 已开通
)

// 支付方式常量
const (
	PaymentMethodAlipay   = "alipay"   // 支付宝
	PaymentMethodWechat   = "wechat"   // 微信
	PaymentMethodBalance  = "balance"  // 余额
	PaymentMethodYipay    = "yipay"    // 易支付
	PaymentMethodEpay     = "epay"     // 易支付（旧）
	PaymentMethodMapay   = "mapay"   // 马支付
)

// Order 订单模型
type Order struct {
	ID           uint           `json:"id" gorm:"primarykey;type:int(11)"`
	UUID         string         `json:"uuid" gorm:"uniqueIndex;not null;size:36"`
	UserID       uint           `json:"userId" gorm:"not null;index:idx_user_id;type:int(11)"`
	ProductID    uint           `json:"productId" gorm:"not null;index:idx_product_id;type:int(11)"`
	InstanceID   *uint          `json:"instanceId" gorm:"index:idx_instance_id;type:int(11)"`
	OrderNo      string         `json:"orderNo" gorm:"uniqueIndex;not null;size:64"` // 订单号
	Amount       float64        `json:"amount" gorm:"not null;type:decimal(10,2)"`    // 订单金额
	Status       string         `json:"status" gorm:"default:pending;size:20;index"`  // 订单状态
	PaymentMethod string        `json:"paymentMethod" gorm:"size:20"`                 // 支付方式
	PaymentID    string         `json:"paymentId" gorm:"size:64"`                     // 支付ID
	PaidAt       *time.Time     `json:"paidAt"`                                       // 支付时间
	ProvisionedAt *time.Time    `json:"provisionedAt"`                                // 开通时间
	ExpiredAt    *time.Time     `json:"expiredAt"`                                    // 过期时间
	Remark       string         `json:"remark" gorm:"size:500"`                      // 备注
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 兼容字段
	ProductData  string         `json:"productData" gorm:"type:text"`  // 产品数据（JSON）
	ExpireAt     *time.Time     `json:"expireAt"`                   // 过期时间（兼容）
	PaymentTime  *time.Time     `json:"paymentTime"`                // 支付时间（兼容）
	PaidAmount   float64        `json:"paidAmount" gorm:"type:decimal(10,2)"` // 已支付金额
}

// 支付状态常量
const (
	PaymentStatusSuccess = "success" // 支付成功
	PaymentStatusFailed  = "failed"  // 支付失败
	PaymentStatusPending = "pending" // 支付中
)

// PaymentRecord 支付记录模型
type PaymentRecord struct {
	ID           uint           `json:"id" gorm:"primarykey;type:int(11)"`
	UUID         string         `json:"uuid" gorm:"uniqueIndex;not null;size:36"`
	OrderID      uint           `json:"orderId" gorm:"not null;index:idx_order_id;type:int(11)"`
	UserID       uint           `json:"userId" gorm:"not null;index:idx_user_id;type:int(11)"`
	PaymentNo    string         `json:"paymentNo" gorm:"uniqueIndex;not null;size:64"` // 支付号
	Amount       float64        `json:"amount" gorm:"not null;type:decimal(10,2)"`       // 支付金额
	PaymentMethod string        `json:"paymentMethod" gorm:"size:20"`                   // 支付方式
	PaymentStatus string        `json:"paymentStatus" gorm:"size:20;index"`             // 支付状态
	PaymentTime  *time.Time     `json:"paymentTime"`                                   // 支付时间
	NotifyData   string         `json:"notifyData" gorm:"type:text"`                    // 回调数据
	Remark       string         `json:"remark" gorm:"size:500"`                        // 备注
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 兼容字段
	Type          string  `json:"type" gorm:"size:20"`          // 支付类型
	TransactionID string  `json:"transactionId" gorm:"size:64"` // 交易ID
	Status        string  `json:"status" gorm:"size:20"`        // 状态（兼容）
}

// BeforeCreate 创建前钩子
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	o.UUID = uuid.New().String()
	return nil
}

func (p *PaymentRecord) BeforeCreate(tx *gorm.DB) error {
	p.UUID = uuid.New().String()
	return nil
}
