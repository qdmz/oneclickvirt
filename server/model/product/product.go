package product

import (
	"time"

	"gorm.io/gorm"
)

// ProductType 产品/实例类型
const (
	ProductTypeVM      = "vm"      // KVM 虚拟机
	ProductTypeLXC     = "lxc"     // LXC 容器
	ProductTypeContainer = "container" // 容器
)

// ProductStatus 产品状态
const (
	ProductStatusOn  = 1 // 上架
	ProductStatusOff = 0 // 下架
)

// Product VM/容器 产品（套餐）
type Product struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"size:128;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Type        string         `json:"type" gorm:"size:16;not null;default:vm"` // vm/lxc/container
	CPU         int            `json:"cpu" gorm:"not null;default:1"`          // CPU 核数
	MemoryMB    int            `json:"memoryMB" gorm:"not null;default:1024"`  // 内存 MB
	DiskMB      int            `json:"diskMB" gorm:"not null;default:10240"`   // 磁盘 MB
	BandwidthMbps int           `json:"bandwidthMbps" gorm:"not null;default:100"` // 带宽 Mbps
	Price       int64          `json:"price" gorm:"not null"`                 // 价格（分）
	PeriodDays  int            `json:"periodDays" gorm:"not null;default:30"` // 计费周期（天）
	// ImageFilter 镜像过滤条件：按虚拟化类型/实例类型过滤可选镜像，例如 "vm" / "lxc"
	ImageFilter string `json:"imageFilter" gorm:"size:32"`
	// NodeID 指定开通节点（0 表示用户下单时自选节点）
	NodeID     uint   `json:"nodeId" gorm:"default:0;index"`
	Stock      int    `json:"stock" gorm:"default:-1"` // 库存（-1 表示不限）
	Sold       int    `json:"sold" gorm:"default:0"`   // 已售数量
	Status     int    `json:"status" gorm:"default:1;index"`
	Sort       int    `json:"sort" gorm:"default:0"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}
