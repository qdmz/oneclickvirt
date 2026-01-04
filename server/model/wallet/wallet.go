package wallet

import (
	"time"
	userModel "oneclickvirt/model/user"
)

// UserWallet 用户钱包表
type UserWallet struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index:idx_user_wallet;uniqueIndex;not null;comment:用户ID" json:"userId"`
	Balance       int64     `gorm:"default:0;not null;comment:余额(分)" json:"balance"`
	Frozen        int64     `gorm:"default:0;not null;comment:冻结金额" json:"frozen"`
	TotalRecharge int64     `gorm:"default:0;not null;comment:累计充值" json:"totalRecharge"`
	TotalExpense  int64     `gorm:"default:0;not null;comment:累计消费" json:"totalExpense"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	User userModel.User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (UserWallet) TableName() string {
	return "user_wallets"
}

// WalletTransaction 钱包交易记录表
type WalletTransaction struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"index:idx_user_txn;not null;comment:用户ID" json:"userId"`
	Type        string     `gorm:"type:varchar(20);not null;comment:交易类型" json:"type"`
	Amount      int64      `gorm:"not null;comment:金额(分)" json:"amount"`
	Balance     int64      `gorm:"not null;comment:交易后余额" json:"balance"`
	Description string     `gorm:"type:varchar(255);comment:交易说明" json:"description"`
	OrderID     *uint      `gorm:"index;comment:关联订单ID" json:"orderId"`
	RelatedID   *uint      `gorm:"index;comment:关联记录ID" json:"relatedId"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`

	User userModel.User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}

// 交易类型常量
const (
	TransactionTypeRecharge = "recharge" // 充值
	TransactionTypeConsume  = "consume"  // 消费
	TransactionTypeRefund   = "refund"   // 退款
	TransactionTypeWithdraw = "withdraw" // 提现
	TransactionTypeExchange = "exchange" // 兑换码
	TransactionTypeSystem   = "system"   // 系统调整
)
