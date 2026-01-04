package product

import (
	"time"
	userModel "oneclickvirt/model/user"
)

// Product 产品配置表
type Product struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null;comment:产品名称"`
	Description  string    `json:"description" gorm:"type:text;comment:产品描述"`
	Level        int       `json:"level" gorm:"not null;comment:对应等级"`
	Price        int64     `json:"price" gorm:"not null;comment:价格(分)"`
	Period       int       `json:"period" gorm:"not null;comment:有效期(天), 0表示永久"`
	CPU          int       `json:"cpu" gorm:"column:cpu;not null;comment:CPU核心数"`
	Memory       int       `json:"memory" gorm:"column:memory;not null;comment:内存(MB)"`
	Disk         int       `json:"disk" gorm:"column:disk;not null;comment:磁盘(MB)"`
	Bandwidth    int       `json:"bandwidth" gorm:"column:bandwidth;not null;comment:带宽(Mbps)"`
	Traffic      int64     `json:"traffic" gorm:"column:traffic;not null;comment:流量配额(MB)"`
	MaxInstances int       `json:"maxInstances" gorm:"column:max_instances;not null;comment:最大实例数"`
	IsEnabled    int       `json:"isEnabled" gorm:"column:is_enabled;default:1;comment:是否启用(1:启用, 0:禁用)"`
	SortOrder    int       `json:"sortOrder" gorm:"column:sort_order;default:0;comment:排序"`
	Features     string    `json:"features" gorm:"column:features;type:text;comment:特性(JSON数组)"`
	AllowRepeat  int       `json:"allowRepeat" gorm:"column:allow_repeat;default:1;comment:是否允许重复购买(1:允许, 0:不允许)"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// ProductPurchase 产品购买记录表
type ProductPurchase struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index:idx_user_purchase;not null;comment:用户ID"`
	ProductID uint      `gorm:"index;not null;comment:产品ID"`
	OrderID   uint      `gorm:"index;comment:订单ID"`
	Level     int       `gorm:"not null;comment:购买后的等级"`
	StartDate time.Time `gorm:"comment:开始时间"`
	EndDate   *time.Time `gorm:"comment:结束时间, 为空表示永久"`
	IsActive  bool      `gorm:"default:true;comment:是否激活"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User    userModel.User `gorm:"foreignKey:UserID"`
	Product Product        `gorm:"foreignKey:ProductID"`
}

// TableName 指定表名
func (ProductPurchase) TableName() string {
	return "product_purchases"
}
