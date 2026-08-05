import request from '@/utils/request'

// 管理员：获取全部工单列表
export const getAdminTickets = (params) => {
  return request({
    url: '/v1/admin/tickets',
    method: 'get',
    params
  })
}

// 管理员：工单状态统计
export const getAdminTicketStats = () => {
  return request({
    url: '/v1/admin/tickets/stats',
    method: 'get'
  })
}

// 管理员：工单详情
export const getAdminTicketDetail = (id) => {
  return request({
    url: `/v1/admin/tickets/${id}`,
    method: 'get'
  })
}

// 管理员：回复工单
export const replyAdminTicket = (id, content) => {
  return request({
    url: `/v1/admin/tickets/${id}/replies`,
    method: 'post',
    data: { content }
  })
}

// 管理员：更新工单状态
export const updateAdminTicketStatus = (id, status) => {
  return request({
    url: `/v1/admin/tickets/${id}/status`,
    method: 'put',
    data: { status }
  })
}

// 管理员：关闭工单
export const closeAdminTicket = (id) => {
  return request({
    url: `/v1/admin/tickets/${id}/close`,
    method: 'post'
  })
}

// 管理员：删除工单
export const deleteAdminTicket = (id) => {
  return request({
    url: `/v1/admin/tickets/${id}`,
    method: 'delete'
  })
}
