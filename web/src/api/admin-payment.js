import request from '@/utils/request'

// ========== 站点配置管理 ==========

// 获取站点配置列表
export const getSiteConfigs = () => {
  return request({
    url: '/v1/admin/site-configs',
    method: 'get'
  })
}

// 获取单个站点配置
export const getSiteConfig = (key) => {
  return request({
    url: `/v1/admin/site-configs/${key}`,
    method: 'get'
  })
}

// 批量更新站点配置
export const updateSiteConfigs = (data) => {
  return request({
    url: '/v1/admin/site-configs',
    method: 'put',
    data
  })
}

// 初始化默认站点配置
export const initializeSiteConfigs = () => {
  return request({
    url: '/v1/admin/site-configs/initialize',
    method: 'post'
  })
}

// ========== 产品管理 ==========

// 获取产品列表
export const getProducts = (params) => {
  return request({
    url: '/v1/admin/products',
    method: 'get',
    params
  })
}

// 创建产品
export const createProduct = (data) => {
  return request({
    url: '/v1/admin/products',
    method: 'post',
    data
  })
}

// 更新产品
export const updateProduct = (id, data) => {
  return request({
    url: `/v1/admin/products/${id}`,
    method: 'put',
    data
  })
}

// 删除产品
export const deleteProduct = (id) => {
  return request({
    url: `/v1/admin/products/${id}`,
    method: 'delete'
  })
}

// 启用/禁用产品
export const toggleProduct = (id) => {
  return request({
    url: `/v1/admin/products/${id}/toggle`,
    method: 'put'
  })
}

// ========== 兑换码管理 ==========

// 获取兑换码列表
export const getRedemptionCodes = (params) => {
  return request({
    url: '/v1/admin/redemption-codes',
    method: 'get',
    params
  })
}

// 创建兑换码
export const createRedemptionCode = (data) => {
  return request({
    url: '/v1/admin/redemption-codes',
    method: 'post',
    data
  })
}

// 批量生成兑换码
export const generateRedemptionCodes = (data) => {
  return request({
    url: '/v1/admin/redemption-codes/generate',
    method: 'post',
    data
  })
}

// 删除兑换码
export const deleteRedemptionCode = (id) => {
  return request({
    url: `/v1/admin/redemption-codes/${id}`,
    method: 'delete'
  })
}

// 获取兑换码使用记录
export const getRedemptionCodeUsages = (id) => {
  return request({
    url: `/v1/admin/redemption-codes/${id}/usages`,
    method: 'get'
  })
}

// 启用/禁用兑换码
export const toggleRedemptionCode = (id) => {
  return request({
    url: `/v1/admin/redemption-codes/${id}/toggle`,
    method: 'put'
  })
}

// ========== 订单管理 ==========

// 获取所有订单列表
export const getAllOrders = (params) => {
  return request({
    url: '/v1/admin/orders',
    method: 'get',
    params
  })
}

// 获取订单详情
export const getOrderDetail = (id) => {
  return request({
    url: `/v1/admin/orders/${id}`,
    method: 'get'
  })
}

// 删除订单
export const deleteOrder = (id) => {
  return request({
    url: `/v1/admin/orders/${id}`,
    method: 'delete'
  })
}

// 取消订单
export const cancelOrder = (id) => {
  return request({
    url: `/v1/admin/orders/${id}/cancel`,
    method: 'post'
  })
}

// 退款订单
export const refundOrder = (id, data) => {
  return request({
    url: `/v1/admin/orders/${id}/refund`,
    method: 'post',
    data
  })
}
