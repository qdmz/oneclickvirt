import request from '@/utils/request'

// ========== 管理员代金券管理 API ==========

/**
 * 获取代金券列表
 * @param {Object} params - { page, pageSize, status, batchNo, code }
 */
export const getVoucherList = (params) => {
  return request({
    url: '/v1/admin/vouchers',
    method: 'get',
    params
  })
}

/**
 * 代金券统计
 */
export const getVoucherStats = () => {
  return request({
    url: '/v1/admin/vouchers/stats',
    method: 'get'
  })
}

/**
 * 批量生成代金券
 * @param {Object} data - { amount, count, prefix, expireAt, remark }
 */
export const createVouchers = (data) => {
  return request({
    url: '/v1/admin/vouchers',
    method: 'post',
    data
  })
}

/**
 * 作废代金券
 * @param {number} id
 */
export const voidVoucher = (id) => {
  return request({
    url: `/v1/admin/vouchers/${id}/void`,
    method: 'put'
  })
}

/**
 * 删除代金券
 * @param {number} id
 */
export const deleteVoucher = (id) => {
  return request({
    url: `/v1/admin/vouchers/${id}`,
    method: 'delete'
  })
}

/**
 * 批量删除代金券（按 ID 列表或批次号）
 * @param {Object} data - { ids?: number[], batchNo?: string }
 */
export const batchDeleteVouchers = (data) => {
  return request({
    url: '/v1/admin/vouchers/batch-delete',
    method: 'post',
    data
  })
}

// ========== 管理员用户余额调整 API ==========

/**
 * 调整指定用户余额
 * @param {number} userId
 * @param {Object} data - { mode: 'add'|'set', amount, remark }
 */
export const adjustUserBalance = (userId, data) => {
  return request({
    url: `/v1/admin/users/${userId}/balance`,
    method: 'put',
    data
  })
}

/**
 * 获取指定用户余额变动记录
 * @param {number} userId
 * @param {Object} params
 */
export const getUserBalanceLogs = (userId, params) => {
  return request({
    url: `/v1/admin/users/${userId}/balance-logs`,
    method: 'get',
    params
  })
}
