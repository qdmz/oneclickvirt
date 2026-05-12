import request from '@/utils/request'

export function getTicketList(params) {
  return request({
    url: '/v1/user/tickets',
    method: 'get',
    params
  })
}

export function createTicket(data) {
  return request({
    url: '/v1/user/tickets',
    method: 'post',
    data
  })
}

export function getTicketDetail(id) {
  return request({
    url: `/v1/user/tickets/${id}`,
    method: 'get'
  })
}

export function replyTicket(id, data) {
  return request({
    url: `/v1/user/tickets/${id}/reply`,
    method: 'post',
    data
  })
}

// Admin APIs
export function getAdminTicketList(params) {
  return request({
    url: '/v1/admin/tickets',
    method: 'get',
    params
  })
}

export function getAdminTicketDetail(id) {
  return request({
    url: `/v1/admin/tickets/${id}`,
    method: 'get'
  })
}

export function updateTicket(id, data) {
  return request({
    url: `/v1/admin/tickets/${id}`,
    method: 'put',
    data
  })
}

export function adminReplyTicket(id, data) {
  return request({
    url: `/v1/admin/tickets/${id}/reply`,
    method: 'post',
    data
  })
}

export function closeTicket(id, data) {
  return request({
    url: `/v1/admin/tickets/${id}/close`,
    method: 'post',
    data
  })
}
