import request from '@/utils/request'

// 管理员：调整用户余额（amount 单位：元，可为负数表示扣减）
export const adjustUserBalance = (userId, data) => {
  return request({
    url: `/v1/admin/users/${userId}/balance`,
    method: 'post',
    data
  })
}

// 管理员：充值/余额变动流水
export const getAdminRechargeOrders = (params) => {
  return request({
    url: '/v1/admin/recharge-orders',
    method: 'get',
    params
  })
}

// 管理员：商品购买订单列表
export const getAdminOrders = (params) => {
  return request({
    url: '/v1/admin/orders',
    method: 'get',
    params
  })
}
