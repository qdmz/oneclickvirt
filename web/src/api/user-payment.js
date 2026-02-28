import request from '@/utils/request'

// ========== 钱包管理 ==========

// 获取钱包信息
export const getWallet = () => {
  return request({
    url: '/v1/user/wallet',
    method: 'get'
  })
}

// 获取交易记录
export const getWalletTransactions = (params) => {
  return request({
    url: '/v1/user/wallet/transactions',
    method: 'get',
    params
  })
}

// ========== 充值管理 ==========

// 创建充值订单
export const createRechargeOrder = (data) => {
  return request({
    url: '/v1/user/recharge/create-order',
    method: 'post',
    data
  })
}

// 获取支付宝支付二维码
export const getAlipayQR = (orderNo) => {
  return request({
    url: `/v1/user/recharge/alipay-qr/${orderNo}`,
    method: 'get'
  })
}

// 获取微信支付二维码
export const getWechatQR = (orderNo) => {
  return request({
    url: `/v1/user/recharge/wechat-qr/${orderNo}`,
    method: 'get'
  })
}

// 获取易支付二维码
export const getEpayQR = (orderNo) => {
  return request({
    url: `/v1/user/recharge/epay-qr/${orderNo}`,
    method: 'get'
  })
}

// 获取码支付二维码
export const getMapayQR = (orderNo) => {
  return request({
    url: `/v1/user/recharge/mapay-qr/${orderNo}`,
    method: 'get'
  })
}

// 获取产品购买易支付二维码
export const getPurchaseEpayQR = (orderNo) => {
  return request({
    url: `/v1/user/orders/epay-qr/${orderNo}`,
    method: 'get'
  })
}

// 获取产品购买码支付二维码
export const getPurchaseMapayQR = (orderNo) => {
  return request({
    url: `/v1/user/orders/mapay-qr/${orderNo}`,
    method: 'get'
  })
}

// 查询充值订单状态
export const getRechargeOrderStatus = (orderNo) => {
  return request({
    url: `/v1/user/recharge/order-status/${orderNo}`,
    method: 'get'
  })
}

// 使用兑换码
export const exchangeRedemptionCode = (data) => {
  return request({
    url: '/v1/user/recharge/exchange-code',
    method: 'post',
    data
  })
}

// ========== 订单管理 ==========

// 获取我的订单列表
export const getUserOrders = (params) => {
  return request({
    url: '/v1/user/orders',
    method: 'get',
    params
  })
}

// 获取订单详情
export const getUserOrderDetail = (id) => {
  return request({
    url: `/v1/user/orders/${id}`,
    method: 'get'
  })
}

// 取消订单
export const cancelOrder = (id) => {
  return request({
    url: `/v1/user/orders/${id}/cancel`,
    method: 'post'
  })
}

// ========== 产品管理 ==========

// 获取可用产品列表
export const getUserProducts = () => {
  return request({
    url: '/v1/user/products',
    method: 'get'
  })
}

// 购买产品
export const purchaseProduct = (id, data) => {
  return request({
    url: `/v1/user/products/${id}/purchase`,
    method: 'post',
    data
  })
}

// 查询购买订单状态
export const getPurchaseOrderStatus = (orderNo) => {
  return request({
    url: `/v1/user/orders/status/${orderNo}`,
    method: 'get'
  })
}
