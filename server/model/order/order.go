package order

import (
	"time"
	productModel "oneclickvirt/model/product"
	userModel "oneclickvirt/model/user"
)

// Order 订单表
type Order struct {
	ID                    uint         `json:"id" gorm:"primaryKey"`
	OrderNo               string       `json:"orderNo" gorm:"type:varchar(50);uniqueIndex;not null;comment:订单号"`
	UserID                uint         `json:"userId" gorm:"index:idx_user_order;not null;comment:用户ID"`
	ProductID             uint         `json:"productId" gorm:"index;not null;comment:产品ID"`
	Amount                float64      `json:"amount" gorm:"type:decimal(10,2);not null;comment:订单金额"`
	Status                int          `json:"status" gorm:"index;default:0;comment:订单状态"`
	PaymentMethod         string       `json:"paymentMethod" gorm:"type:varchar(20);comment:支付方式"`
	PaymentTransactionID  string       `json:"paymentTransactionId" gorm:"type:varchar(100);comment:支付交易ID"`
	PaymentTime           *time.Time   `json:"paymentTime" gorm:"comment:支付时间"`
	PaidAmount            float64      `json:"paidAmount" gorm:"type:decimal(10,2);default:0;comment:实付金额"`
	ProductData           string       `json:"productData" gorm:"type:json;comment:产品快照"`
	Remark                string       `json:"remark" gorm:"type:varchar(255);comment:备注"`
	ExpireAt              time.Time    `json:"expireAt" gorm:"comment:订单过期时间"`
	CreatedAt             time.Time    `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt             time.Time    `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt             *time.Time   `json:"deletedAt" gorm:"index"`

	User    *userModel.User      `json:"user" gorm:"foreignKey:UserID"`
	Product *productModel.Product `json:"product" gorm:"foreignKey:ProductID"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// 订单状态常量
const (
	OrderStatusPending   = 0 // 待支付
	OrderStatusPaid      = 1 // 已支付
	OrderStatusCancelled = 2 // 已取消
	OrderStatusRefunded  = 3 // 已退款
	OrderStatusExpired   = 4 // 已过期
)

// 支付方式常量
const (
	PaymentMethodAlipay   = "alipay"   // 支付宝
	PaymentMethodWechat   = "wechat"   // 微信支付
	PaymentMethodBalance  = "balance"  // 余额支付
	PaymentMethodExchange = "exchange" // 兑换码
	PaymentMethodEpay     = "epay"     // 易支付
	PaymentMethodMapay    = "mapay"    // 码支付
)

// PaymentRecord 支付记录表
type PaymentRecord struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OrderID        uint      `gorm:"index:idx_order_payment;not null;comment:订单ID" json:"orderId"`
	UserID         uint      `gorm:"index;not null;comment:用户ID" json:"userId"`
	Type           string    `gorm:"type:varchar(20);not null;comment:支付类型" json:"type"`
	TransactionID  string    `gorm:"type:varchar(64);uniqueIndex;comment:第三方交易号" json:"transactionId"`
	Amount         int64     `gorm:"not null;comment:支付金额(分)" json:"amount"`
	Status         string    `gorm:"type:varchar(20);not null;comment:支付状态" json:"status"`
	NotifyData     string    `gorm:"type:text;comment:回调数据" json:"notifyData"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Order Order       `gorm:"foreignKey:OrderID" json:"-"`
	User  userModel.User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (PaymentRecord) TableName() string {
	return "payment_records"
}

// 支付状态常量
const (
	PaymentStatusProcessing = "processing" // 处理中
	PaymentStatusSuccess    = "success"    // 成功
	PaymentStatusFailed     = "failed"     // 失败
	PaymentStatusCancelled  = "cancelled"  // 已取消
)
