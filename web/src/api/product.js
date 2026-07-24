import request from '@/utils/request'

// ========== 产品相关 API ==========

/**
 * 获取产品列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export const getProductList = (params) => {
  return request({
    url: '/v1/user/products',
    method: 'get',
    params
  })
}

/**
 * 获取产品详情
 * @param {number} id - 产品ID
 * @returns {Promise}
 */
export const getProductDetail = (id) => {
  return request({
    url: `/v1/user/products/${id}`,
    method: 'get'
  })
}

/**
 * 获取产品图片列表
 * @param {number} productId - 产品ID
 * @returns {Promise}
 */
export const getProductImages = (productId) => {
  return request({
    url: `/v1/user/products/${productId}/images`,
    method: 'get'
  })
}

// ========== 订单相关 API ==========

/**
 * 创建订单
 * @param {Object} data - 订单数据
 * @returns {Promise}
 */
export const createOrder = (data) => {
  return request({
    url: '/v1/user/orders',
    method: 'post',
    data
  })
}

/**
 * 获取订单列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export const getOrderList = (params) => {
  return request({
    url: '/v1/user/orders',
    method: 'get',
    params
  })
}

/**
 * 获取订单详情
 * @param {number} id - 订单ID
 * @returns {Promise}
 */
export const getOrderDetail = (id) => {
  return request({
    url: `/v1/user/orders/${id}`,
    method: 'get'
  })
}

/**
 * 使用余额支付订单
 * @param {number} orderId - 订单ID
 * @returns {Promise}
 */
export const payWithBalance = (orderId) => {
  return request({
    url: `/v1/user/orders/${orderId}/pay`,
    method: 'post',
    data: { payment_method: 'balance' }
  })
}

/**
 * 取消未支付订单
 * @param {number} orderId - 订单ID
 * @returns {Promise}
 */
export const cancelOrder = (orderId) => {
  return request({
    url: `/v1/user/orders/${orderId}/cancel`,
    method: 'post'
  })
}

/**
 * 续费订单
 * @param {number} orderId - 订单ID
 * @param {Object} data - 续费数据
 * @returns {Promise}
 */
export const renewOrder = (orderId, data) => {
  return request({
    url: `/v1/user/orders/${orderId}/renew`,
    method: 'post',
    data
  })
}

// ========== 支付相关 API ==========

/**
 * 创建易支付订单
 * @param {Object} data - 支付数据
 * @returns {Promise}
 */
export const createYiPayOrder = (data) => {
  return request({
    url: '/v1/user/payments/yipay',
    method: 'post',
    data
  })
}

/**
 * 获取充值记录列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export const getRechargeList = (params) => {
  return request({
    url: '/v1/user/recharges',
    method: 'get',
    params
  })
}

// ========== 余额相关 API ==========

/**
 * 获取用户余额
 * @returns {Promise}
 */
export const getUserBalance = () => {
  return request({
    url: '/v1/user/balance',
    method: 'get'
  })
}

/**
 * 获取余额变动记录
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export const getBalanceLogs = (params) => {
  return request({
    url: '/v1/user/balance/logs',
    method: 'get',
    params
  })
}
