import request from '@/utils/request'

// ========== 管理员产品管理 API ==========

/**
 * 获取管理员产品列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export const getAdminProductList = (params) => {
  return request({ url: '/v1/admin/products', method: 'get', params })
}

/**
 * 创建产品
 * @param {Object} data - 产品数据
 * @returns {Promise}
 */
export const createAdminProduct = (data) => {
  return request({ url: '/v1/admin/products', method: 'post', data })
}

/**
 * 获取产品详情
 * @param {number} id - 产品ID
 * @returns {Promise}
 */
export const getAdminProductDetail = (id) => {
  return request({ url: `/v1/admin/products/${id}`, method: 'get' })
}

/**
 * 更新产品
 * @param {number} id - 产品ID
 * @param {Object} data - 产品数据
 * @returns {Promise}
 */
export const updateAdminProduct = (id, data) => {
  return request({ url: `/v1/admin/products/${id}`, method: 'put', data })
}

/**
 * 删除产品
 * @param {number} id - 产品ID
 * @returns {Promise}
 */
export const deleteAdminProduct = (id) => {
  return request({ url: `/v1/admin/products/${id}`, method: 'delete' })
}

/**
 * 更新产品状态
 * @param {number} id - 产品ID
 * @param {Object} data - 状态数据
 * @returns {Promise}
 */
export const updateAdminProductStatus = (id, data) => {
  return request({ url: `/v1/admin/products/${id}/status`, method: 'put', data })
}

// ========== 管理员订单管理 API ==========

/**
 * 获取管理员订单列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export const getAdminOrderList = (params) => {
  return request({ url: '/v1/admin/orders', method: 'get', params })
}

/**
 * 获取订单详情
 * @param {number} id - 订单ID
 * @returns {Promise}
 */
export const getAdminOrderDetail = (id) => {
  return request({ url: `/v1/admin/orders/${id}`, method: 'get' })
}

/**
 * 更新订单状态
 * @param {number} id - 订单ID
 * @param {Object} data - 状态数据
 * @returns {Promise}
 */
export const updateAdminOrderStatus = (id, data) => {
  return request({ url: `/v1/admin/orders/${id}/status`, method: 'put', data })
}

/**
 * 手动开通订单
 * @param {number} id - 订单ID
 * @returns {Promise}
 */
export const manualProvision = (id) => {
  return request({ url: `/v1/admin/orders/${id}/provision`, method: 'post' })
}

// ========== 管理员工单管理 API ==========

/**
 * 获取管理员工单列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export const getAdminTicketList = (params) => {
  return request({ url: '/v1/admin/tickets', method: 'get', params })
}

/**
 * 获取工单详情
 * @param {number} id - 工单ID
 * @returns {Promise}
 */
export const getAdminTicketDetail = (id) => {
  return request({ url: `/v1/admin/tickets/${id}`, method: 'get' })
}

/**
 * 管理员回复工单
 * @param {number} id - 工单ID
 * @param {Object} data - 回复数据
 * @returns {Promise}
 */
export const adminReplyTicket = (id, data) => {
  return request({ url: `/v1/admin/tickets/${id}/reply`, method: 'post', data })
}

/**
 * 更新工单状态
 * @param {number} id - 工单ID
 * @param {Object} data - 状态数据
 * @returns {Promise}
 */
export const updateAdminTicketStatus = (id, data) => {
  return request({ url: `/v1/admin/tickets/${id}/status`, method: 'put', data })
}

/**
 * 分配工单
 * @param {number} id - 工单ID
 * @param {Object} data - 分配数据
 * @returns {Promise}
 */
export const assignAdminTicket = (id, data) => {
  return request({ url: `/v1/admin/tickets/${id}/assign`, method: 'post', data })
}

/**
 * 获取工单统计数据
 * @returns {Promise}
 */
export const getAdminTicketStats = () => {
  return request({ url: '/v1/admin/tickets/stats', method: 'get' })
}

// ========== 管理员站点配置 API ==========

/**
 * 获取站点配置（管理员端）
 * @returns {Promise}
 */
export const getAdminSiteConfig = () => {
  return request({ url: '/v1/admin/site-config', method: 'get' })
}

/**
 * 更新站点配置（管理员端）
 * @param {Object} data - 配置数据
 * @returns {Promise}
 */
export const updateAdminSiteConfig = (data) => {
  return request({ url: '/v1/admin/site-config', method: 'put', data })
}

/**
 * 上传站点图片（Logo/Favicon）
 * @param {FormData} formData - 包含图片文件和类型的表单数据
 * @returns {Promise}
 */
export const uploadSiteImage = (formData) => {
  return request({
    url: '/v1/admin/site-config/upload-image',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// ========== 管理员易支付配置 API ==========

/**
 * 获取易支付配置
 * @returns {Promise}
 */
export const getYiPayConfig = () => {
  return request({ url: '/v1/admin/yipay-config', method: 'get' })
}

/**
 * 更新易支付配置
 * @param {Object} data - 配置数据
 * @returns {Promise}
 */
export const updateYiPayConfig = (data) => {
  return request({ url: '/v1/admin/yipay-config', method: 'put', data })
}

/**
 * 测试易支付连通性/密钥自检
 * 用当前配置的 key 向网关发起订单查询，返回网关原始响应，用于排查「MD5签名校验失败」(key 不一致)
 * @returns {Promise}
 */
export const testYiPayConfig = () => {
  return request({ url: '/v1/admin/yipay-config/test', method: 'get' })
}

// ========== 管理员站点链接管理 API ==========

/**
 * 获取站点链接列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export const getAdminSiteLinkList = (params) => {
  return request({ url: '/v1/admin/site-links', method: 'get', params })
}

/**
 * 创建站点链接
 * @param {Object} data - 站点链接数据
 * @returns {Promise}
 */
export const createAdminSiteLink = (data) => {
  return request({ url: '/v1/admin/site-links', method: 'post', data })
}

/**
 * 获取站点链接详情
 * @param {number} id - 站点链接ID
 * @returns {Promise}
 */
export const getAdminSiteLinkDetail = (id) => {
  return request({ url: `/v1/admin/site-links/${id}`, method: 'get' })
}

/**
 * 更新站点链接
 * @param {number} id - 站点链接ID
 * @param {Object} data - 站点链接数据
 * @returns {Promise}
 */
export const updateAdminSiteLink = (id, data) => {
  return request({ url: `/v1/admin/site-links/${id}`, method: 'put', data })
}

/**
 * 删除站点链接
 * @param {number} id - 站点链接ID
 * @returns {Promise}
 */
export const deleteAdminSiteLink = (id) => {
  return request({ url: `/v1/admin/site-links/${id}`, method: 'delete' })
}
