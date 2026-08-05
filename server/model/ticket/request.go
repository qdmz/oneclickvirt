package ticket

import "oneclickvirt/model/common"

// CreateTicketRequest 用户创建工单请求
type CreateTicketRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"` // 首条留言内容
	Category string `json:"category"`                   // 分类，可空
	Priority string `json:"priority"`                   // 优先级，可空
}

// ReplyTicketRequest 回复工单请求（用户与管理员通用）
type ReplyTicketRequest struct {
	Content string `json:"content" binding:"required"`
}

// TicketListRequest 工单列表查询请求
type TicketListRequest struct {
	common.PageInfo
	Status   string `json:"status" form:"status"`     // 状态过滤，空=全部
	Priority string `json:"priority" form:"priority"` // 优先级过滤，空=全部
	Category string `json:"category" form:"category"` // 分类过滤，空=全部
	Keyword  string `json:"keyword" form:"keyword"`   // 标题关键词
}

// UpdateTicketStatusRequest 管理员更新工单状态请求
type UpdateTicketStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=open answered closed"`
}
