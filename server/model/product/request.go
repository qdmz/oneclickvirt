package product

import "oneclickvirt/model/common"

// ========== 产品相关请求 ==========

// ProductListRequest 产品列表请求
type ProductListRequest struct {
	common.PageInfo
	Category string `json:"category" form:"category"` // 类别筛选
	Type     string `json:"type" form:"type"`         // 虚拟化类型筛选
}

// ========== 订单相关请求 ==========

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	ProductID uint `json:"productId" binding:"required"` // 产品ID
	ImageID   uint `json:"imageId" binding:"required"`   // 镜像ID
	Quantity  int  `json:"quantity" binding:"required,min=1,max=36"` // 购买周期数
}

// OrderListRequest 订单列表请求
type OrderListRequest struct {
	common.PageInfo
	PaymentStatus   int    `json:"paymentStatus" form:"paymentStatus"`   // 支付状态筛选
	ProvisionStatus int    `json:"provisionStatus" form:"provisionStatus"` // 开通状态筛选
}

// PayOrderRequest 余额支付请求
type PayOrderRequest struct {
	OrderID uint `json:"orderId" binding:"required"` // 订单ID
}

// RenewOrderRequest 续费订单请求
type RenewOrderRequest struct {
	OrderID  uint `json:"orderId" binding:"required"` // 原订单ID
	Quantity int  `json:"quantity" binding:"required,min=1,max=36"` // 续费周期数
}

// CancelOrderRequest 取消订单请求
type CancelOrderRequest struct {
	OrderID uint `json:"orderId" binding:"required"` // 订单ID
}

// ========== 易支付相关请求 ==========

// CreateYiPayOrderRequest 创建易支付订单请求
type CreateYiPayOrderRequest struct {
	Amount  float64 `json:"amount" binding:"required,min=0.01"` // 充值金额
	PayType string  `json:"payType" binding:"required"`         // 支付方式 alipay/wxpay/qqpay
}

// YiPayNotifyRequest 易支付异步通知请求
type YiPayNotifyRequest struct {
	Pid         string  `json:"pid" form:"pid"`                   // 商户ID
	TradeNo     string  `json:"trade_no" form:"trade_no"`         // 平台订单号
	OutTradeNo  string  `json:"out_trade_no" form:"out_trade_no"` // 商户订单号
	Type        string  `json:"type" form:"type"`                 // 支付方式
	Name        string  `json:"name" form:"name"`                 // 商品名称
	Money       string  `json:"money" form:"money"`               // 订单金额
	TradeStatus string  `json:"trade_status" form:"trade_status"` // 交易状态
	Sign        string  `json:"sign" form:"sign"`                 // 签名
	SignType    string  `json:"sign_type" form:"sign_type"`       // 签名类型
}

// ========== 管理员相关请求 ==========

// AdminProductListRequest 管理员产品列表请求
type AdminProductListRequest struct {
	common.PageInfo
	Status   int    `json:"status" form:"status"`     // 状态筛选
	Category string `json:"category" form:"category"` // 类别筛选
	Type     string `json:"type" form:"type"`         // 虚拟化类型筛选
}

// CreateProductRequest 创建产品请求
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,max=128"`      // 产品名称
	Description string  `json:"description"`                            // 产品描述
	Type        string  `json:"type" binding:"required,max=32"`       // 虚拟化类型
	Category    string  `json:"category" binding:"max=16"`            // 类别
	CPU         int     `json:"cpu" binding:"required,min=1"`         // CPU核心数
	Memory      int     `json:"memory" binding:"required,min=128"`    // 内存(MB)
	Disk        int     `json:"disk" binding:"required,min=1024"`     // 磁盘(MB)
	Bandwidth   int     `json:"bandwidth" binding:"min=1"`            // 带宽(Mbps)
	Traffic     int     `json:"traffic" binding:"min=0"`              // 流量配额(MB)
	Price       float64 `json:"price" binding:"required,min=0"`       // 单价
	PeriodType  string  `json:"periodType" binding:"required,max=16"` // 周期类型
	PeriodValue int     `json:"periodValue" binding:"required,min=1"` // 周期值
	MaxSnapshots int    `json:"maxSnapshots" binding:"min=0"`         // 最大快照数
	MaxPorts    int     `json:"maxPorts" binding:"min=0"`             // 最大端口映射数
	Status      int     `json:"status" binding:"oneof=0 1"`           // 状态
	SortOrder   int     `json:"sortOrder"`                              // 排序
	Icon        string  `json:"icon" binding:"max=256"`               // 图标URL
	ImageIDs    string  `json:"imageIds"`                               // 关联镜像ID列表
	ProviderIDs string  `json:"providerIds"`                            // 节点限制
}

// UpdateProductRequest 更新产品请求
type UpdateProductRequest struct {
	Name         string  `json:"name" binding:"required,max=128"`
	Description  string  `json:"description"`
	Type         string  `json:"type" binding:"required,max=32"`
	Category     string  `json:"category" binding:"max=16"`
	CPU          int     `json:"cpu" binding:"required,min=1"`
	Memory       int     `json:"memory" binding:"required,min=128"`
	Disk         int     `json:"disk" binding:"required,min=1024"`
	Bandwidth    int     `json:"bandwidth" binding:"min=1"`
	Traffic      int     `json:"traffic" binding:"min=0"`
	Price        float64 `json:"price" binding:"required,min=0"`
	PeriodType   string  `json:"periodType" binding:"required,max=16"`
	PeriodValue  int     `json:"periodValue" binding:"required,min=1"`
	MaxSnapshots int     `json:"maxSnapshots" binding:"min=0"`
	MaxPorts     int     `json:"maxPorts" binding:"min=0"`
	Status       int     `json:"status" binding:"oneof=0 1"`
	SortOrder    int     `json:"sortOrder"`
	Icon         string  `json:"icon" binding:"max=256"`
	ImageIDs     string  `json:"imageIds"`
	ProviderIDs  string  `json:"providerIds"`
}

// UpdateProductStatusRequest 更新产品状态请求
type UpdateProductStatusRequest struct {
	Status int `json:"status" binding:"oneof=0 1"` // 0=下架 1=上架
}

// AdminOrderListRequest 管理员订单列表请求
type AdminOrderListRequest struct {
	common.PageInfo
	UserID          uint   `json:"userId" form:"userId"`                   // 用户ID筛选
	PaymentStatus   int    `json:"paymentStatus" form:"paymentStatus"`     // 支付状态筛选
	ProvisionStatus int    `json:"provisionStatus" form:"provisionStatus"` // 开通状态筛选
	OrderNo         string `json:"orderNo" form:"orderNo"`                 // 订单号搜索
}

// UpdateOrderStatusRequest 更新订单状态请求
type UpdateOrderStatusRequest struct {
	OrderID         uint   `json:"orderId" binding:"required"`
	PaymentStatus   int    `json:"paymentStatus" binding:"oneof=0 1 2 3"`    // 0=未支付 1=已支付 2=支付失败 3=已退款
	ProvisionStatus int    `json:"provisionStatus" binding:"oneof=0 1 2 3"`  // 0=待开通 1=开通中 2=已开通 3=开通失败
	Remark          string `json:"remark" binding:"max=512"`
}
