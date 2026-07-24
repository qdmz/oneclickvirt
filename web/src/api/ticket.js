import request from '@/utils/request'

// ========== 工单相关 API ==========

/**
 * 创建工单
 * @param {Object} data - 工单数据 {subject, content, category}
 * @returns {Promise}
 */
export const createTicket = (data) => {
  return request({
    url: '/v1/user/tickets',
    method: 'post',
    data
  })
}

/**
 * 获取工单列表
 * @param {Object} params - 查询参数 {page, pageSize, status, category}
 * @returns {Promise}
 */
export const getTicketList = (params) => {
  return request({
    url: '/v1/user/tickets',
    method: 'get',
    params
  })
}

/**
 * 获取工单详情
 * @param {number} id - 工单ID
 * @returns {Promise}
 */
export const getTicketDetail = (id) => {
  return request({
    url: `/v1/user/tickets/${id}`,
    method: 'get'
  })
}

/**
 * 回复工单
 * @param {number} id - 工单ID
 * @param {Object} data - 回复内容 {content}
 * @returns {Promise}
 */
export const replyTicket = (id, data) => {
  return request({
    url: `/v1/user/tickets/${id}/reply`,
    method: 'post',
    data
  })
}

/**
 * 关闭工单
 * @param {number} id - 工单ID
 * @returns {Promise}
 */
export const closeTicket = (id) => {
  return request({
    url: `/v1/user/tickets/${id}/close`,
    method: 'post'
  })
}
