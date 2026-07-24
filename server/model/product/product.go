package product

import (
	"time"

	"gorm.io/gorm"
)

// Product 虚拟机产品定义表
type Product struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Name        string  `json:"name" gorm:"not null;size:128"`           // 产品名称
	Description string  `json:"description" gorm:"type:text"`              // 产品描述
	Type        string  `json:"type" gorm:"not null;size:32;index"`       // 虚拟化类型: lxd/incus/docker/qemu/kubevirt
	Category    string  `json:"category" gorm:"default:vm;size:16"`       // 类别: vm/container

	// 资源配置
	CPU       int `json:"cpu" gorm:"not null;default:1"`       // CPU核心数
	Memory    int `json:"memory" gorm:"not null;default:512"`  // 内存(MB)
	Disk      int `json:"disk" gorm:"not null;default:10240"`  // 磁盘(MB)
	Bandwidth int `json:"bandwidth" gorm:"default:100"`        // 带宽(Mbps)
	Traffic   int `json:"traffic" gorm:"default:102400"`       // 流量配额(MB)

	// 计费配置
	Price        float64 `json:"price" gorm:"not null;default:0"`       // 单价(元/周期)
	PeriodType   string  `json:"periodType" gorm:"default:month;size:16"` // 周期类型: hour/day/month/year
	PeriodValue  int     `json:"periodValue" gorm:"default:1"`           // 周期值

	// 限制配置
	MaxSnapshots int `json:"maxSnapshots" gorm:"default:1"`  // 最大快照数
	MaxPorts     int `json:"maxPorts" gorm:"default:0"`      // 最大端口映射数(0=不限)

	// 状态
	Status      int    `json:"status" gorm:"default:1;index"`  // 0=下架 1=上架
	SortOrder   int    `json:"sortOrder" gorm:"default:0"`      // 排序
	Icon        string `json:"icon" gorm:"size:256"`            // 产品图标URL

	// 关联镜像(逗号分隔的镜像ID列表,空则使用所有可用镜像)
	ImageIDs string `json:"imageIds" gorm:"type:text"`

	// 节点限制(逗号分隔的Provider ID列表,空则所有节点可用)
	ProviderIDs string `json:"providerIds" gorm:"type:text"`
}

// ProductOrder 产品订单表
type ProductOrder struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	OrderNo   string `json:"orderNo" gorm:"uniqueIndex;not null;size:64"` // 订单号
	UserID    uint   `json:"userId" gorm:"not null;index"`                // 用户ID
	ProductID uint   `json:"productId" gorm:"not null;index"`             // 产品ID

	// 产品快照(防止产品修改后订单信息变化)
	ProductName  string  `json:"productName" gorm:"size:128"`
	ProductType  string  `json:"productType" gorm:"size:32"`
	CPU          int     `json:"cpu"`
	Memory       int     `json:"memory"`
	Disk         int     `json:"disk"`
	Bandwidth    int     `json:"bandwidth"`
	Traffic      int     `json:"traffic"`
	PeriodType   string  `json:"periodType" gorm:"size:16"`
	PeriodValue  int     `json:"periodValue"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity" gorm:"default:1"`  // 购买数量(周期数)
	TotalAmount  float64 `json:"totalAmount"`                 // 总金额

	// 支付信息
	PaymentMethod string  `json:"paymentMethod" gorm:"size:32;index"` // balance/yipay/alipay
	PaymentStatus int     `json:"paymentStatus" gorm:"default:0;index"` // 0=未支付 1=已支付 2=支付失败 3=已退款
	PaidAt        *time.Time `json:"paidAt"`
	TradeNo       string  `json:"tradeNo" gorm:"size:128;index"` // 第三方支付流水号

	// 开通信息
	ProvisionStatus int        `json:"provisionStatus" gorm:"default:0;index"` // 0=待开通 1=开通中 2=已开通 3=开通失败
	InstanceID      uint       `json:"instanceId" gorm:"index"`                // 关联的实例ID
	ProvisionedAt   *time.Time `json:"provisionedAt"`
	ExpireAt        *time.Time `json:"expireAt"` // 到期时间

	// 续费信息
	IsRenewal     bool   `json:"isRenewal" gorm:"default:false"` // 是否续费订单
	RenewOrderID  uint   `json:"renewOrderId" gorm:"index"`      // 原订单ID(续费时)

	// 镜像选择
	ImageID   uint   `json:"imageId"`   // 选择的镜像ID
	ImageName string `json:"imageName" gorm:"size:128"` // 镜像名称快照

	// 备注
	Remark string `json:"remark" gorm:"type:text"`
}

// Ticket 用户工单表
type Ticket struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	TicketNo  string `json:"ticketNo" gorm:"uniqueIndex;not null;size:64"` // 工单编号
	UserID    uint   `json:"userId" gorm:"not null;index"`                 // 用户ID
	Title     string `json:"title" gorm:"not null;size:256"`               // 标题
	Content   string `json:"content" gorm:"type:text"`                     // 内容

	// 分类和优先级
	Category string `json:"category" gorm:"default:general;size:32;index"` // general/technical/billing/complaint
	Priority string `json:"priority" gorm:"default:normal;size:16"`        // low/normal/high/urgent

	// 状态
	Status    int    `json:"status" gorm:"default:0;index"` // 0=待处理 1=处理中 2=已解决 3=已关闭
	ClosedAt  *time.Time `json:"closedAt"`
	SolvedAt  *time.Time `json:"solvedAt"`

	// 关联资源
	InstanceID uint `json:"instanceId" gorm:"index"` // 关联实例ID(可选)
	OrderID    uint `json:"orderId" gorm:"index"`    // 关联订单ID(可选)

	// 最后回复信息
	LastReplyAt     *time.Time `json:"lastReplyAt"`
	LastReplyBy     string     `json:"lastReplyBy" gorm:"size:16"` // user/admin
	LastReplyUserID uint       `json:"lastReplyUserId"`
}

// TicketReply 工单回复表
type TicketReply struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`

	TicketID uint   `json:"ticketId" gorm:"not null;index"`
	UserID   uint   `json:"userId" gorm:"not null"`       // 回复者ID
	UserType string `json:"userType" gorm:"size:16"`      // user/admin
	Content  string `json:"content" gorm:"type:text;not null"`

	// 附件(JSON格式)
	Attachments string `json:"attachments" gorm:"type:text"`

	// 是否内部备注(仅管理员可见)
	IsInternal bool `json:"isInternal" gorm:"default:false"`
}

// SiteConfig 站点前端配置表
type SiteConfig struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 基础信息
	SiteName        string `json:"site_name" gorm:"size:128"`      // 站点名称
	SiteDescription string `json:"site_description" gorm:"type:text"` // 站点描述
	SiteKeywords    string `json:"site_keywords" gorm:"size:512"`     // SEO关键词

	// Logo和图标
	LogoURL        string `json:"logo_url" gorm:"size:512"`
	FaviconURL     string `json:"favicon_url" gorm:"size:512"`
	DarkLogoURL    string `json:"dark_logo_url" gorm:"size:512"`    // 暗色模式Logo

	// 页眉配置
	HeaderHTML     string `json:"custom_header" gorm:"type:text"`      // 自定义页眉HTML
	HeaderEnabled  bool   `json:"header_enabled" gorm:"default:false"` // 是否启用自定义页眉
	ShowNav        bool   `json:"show_nav" gorm:"default:true"`        // 是否显示导航

	// 页脚配置
	FooterHTML     string `json:"custom_footer" gorm:"type:text"`      // 自定义页脚HTML
	FooterEnabled  bool   `json:"footer_enabled" gorm:"default:false"` // 是否启用自定义页脚
	CopyrightText  string `json:"copyright_text" gorm:"size:512"`      // 版权文字
	ICPNumber      string `json:"icp_number" gorm:"size:128"`          // ICP备案号

	// 首页配置
	HomeTitle      string `json:"home_title" gorm:"size:256"`        // 首页大标题
	HomeSubtitle   string `json:"home_subtitle" gorm:"size:512"`     // 首页副标题
	HomeBackground string `json:"home_background" gorm:"size:512"`   // 首页背景图
	ShowHomeStats  bool   `json:"show_home_stats" gorm:"default:true"` // 是否显示首页统计

	// 主题配置
	PrimaryColor   string `json:"primary_color" gorm:"size:16;default:#409EFF"` // 主题色
	ThemeMode      string `json:"theme_mode" gorm:"default:auto;size:16"`       // light/dark/auto
	CustomCSS      string `json:"custom_css" gorm:"type:text"`                  // 自定义CSS

	// 联系信息
	ContactEmail    string `json:"contact_email" gorm:"size:128"`
	ContactPhone    string `json:"contact_phone" gorm:"size:32"`
	ContactQQ       string `json:"contact_qq" gorm:"size:32"`
	ContactTelegram string `json:"contact_telegram" gorm:"size:128"`

	// 支付配置(前端展示用)
	ShowBalance     bool   `json:"show_balance" gorm:"default:true"`  // 显示余额支付
	ShowYiPay       bool   `json:"show_yipay" gorm:"default:false"`   // 显示易支付

	// 功能开关
	EnableRegistration bool `json:"enable_registration" gorm:"default:true"`
	EnableTicket       bool `json:"enable_ticket" gorm:"default:true"`
	EnableProductStore bool `json:"enable_product_store" gorm:"default:true"`

	// 公告栏
	AnnouncementBar     string `json:"announcement_bar" gorm:"type:text"`      // 顶部公告栏内容
	AnnouncementEnabled bool   `json:"announcement_enabled" gorm:"default:false"`
}

// UserBalanceLog 用户余额变动记录
type UserBalanceLog struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`

	UserID    uint    `json:"userId" gorm:"not null;index"`
	Type      string  `json:"type" gorm:"size:32;index"`     // recharge/consume/refund/bonus
	Amount    float64 `json:"amount"`                          // 变动金额(正数增加,负数减少)
	BalanceBefore float64 `json:"balanceBefore"`               // 变动前余额
	BalanceAfter  float64 `json:"balanceAfter"`                // 变动后余额

	OrderID   uint   `json:"orderId" gorm:"index"`           // 关联订单
	Remark    string `json:"remark" gorm:"size:512"`         // 备注
	TradeNo   string `json:"tradeNo" gorm:"size:128;index"`  // 交易流水号
}

// YiPayConfig 易支付配置表
type YiPayConfig struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Enabled     bool   `json:"enabled" gorm:"default:false"`
	Name        string `json:"name" gorm:"size:64;default:易支付"`
	ApiURL      string `json:"apiUrl" gorm:"size:512"`           // 支付网关地址
	Pid         string `json:"pid" gorm:"size:64"`               // 商户ID
	Key         string `json:"key" gorm:"size:256"`              // 商户密钥
	NotifyURL   string `json:"notifyUrl" gorm:"size:512"`        // 异步通知地址
	ReturnURL   string `json:"returnUrl" gorm:"size:512"`        // 同步跳转地址
	PayType     string `json:"payType" gorm:"size:32;default:alipay"` // alipay/wxpay/qqpay

	// 费率配置
	FeePercent  float64 `json:"feePercent" gorm:"default:0"`     // 手续费百分比
	MinAmount   float64 `json:"minAmount" gorm:"default:0.01"`   // 最小支付金额
	MaxAmount   float64 `json:"maxAmount" gorm:"default:100000"` // 最大支付金额
}
