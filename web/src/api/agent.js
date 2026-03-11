import request from '@/utils/request'

// 代理商端 API

// 创建代理商申请
export const createAgent = (data) => {
  return request({
    url: '/v1/agent/apply',
    method: 'post',
    data
  })
}

// 获取代理商资料
export const getAgentProfile = () => {
  return request({
    url: '/v1/agent/profile',
    method: 'get'
  })
}

// 更新代理商资料
export const updateAgentProfile = (data) => {
  return request({
    url: '/v1/agent/profile',
    method: 'put',
    data
  })
}

// 获取子用户列表
export const getSubUsers = (params) => {
  return request({
    url: '/v1/agent/sub-users',
    method: 'get',
    params
  })
}

// 删除子用户
export const deleteSubUser = (userId) => {
  return request({
    url: `/v1/agent/sub-users/${userId}`,
    method: 'delete'
  })
}

// 批量更新子用户状态
export const batchUpdateSubUserStatus = (userIds, status) => {
  return request({
    url: '/v1/agent/sub-users/batch-status',
    method: 'put',
    data: { userIds, status }
  })
}

// 批量删除子用户
export const batchDeleteSubUsers = (userIds) => {
  return request({
    url: '/v1/agent/sub-users/batch-delete',
    method: 'post',
    data: { userIds }
  })
}

// 获取代理商统计
export const getAgentStatistics = () => {
  return request({
    url: '/v1/agent/statistics',
    method: 'get'
  })
}

// 获取佣金记录
export const getCommissions = (params) => {
  return request({
    url: '/v1/agent/commissions',
    method: 'get',
    params
  })
}

// 获取钱包信息
export const getAgentWallet = () => {
  return request({
    url: '/v1/agent/wallet',
    method: 'get'
  })
}

// 获取钱包交易记录
export const getWalletTransactions = (params) => {
  return request({
    url: '/v1/agent/wallet/transactions',
    method: 'get',
    params
  })
}

// 提现
export const withdraw = (amount) => {
  return request({
    url: '/v1/agent/wallet/withdraw',
    method: 'post',
    data: { amount }
  })
}

// 管理员端 API

// 获取代理商列表
export const getAgentList = (params) => {
  return request({
    url: '/v1/admin/agents',
    method: 'get',
    params
  })
}

// 管理员创建代理商
export const createAgentByAdmin = (data) => {
  return request({
    url: '/v1/admin/agents',
    method: 'post',
    data
  })
}

// 管理员更新代理商
export const updateAgentByAdmin = (id, data) => {
  return request({
    url: `/v1/admin/agents/${id}`,
    method: 'put',
    data
  })
}

// 删除代理商
export const deleteAgent = (id) => {
  return request({
    url: `/v1/admin/agents/${id}`,
    method: 'delete'
  })
}

// 审核通过
export const approveAgent = (id) => {
  return request({
    url: `/v1/admin/agents/${id}/approve`,
    method: 'post'
  })
}

// 更新代理商状态
export const updateAgentStatus = (id, status) => {
  return request({
    url: `/v1/admin/agents/${id}/status`,
    method: 'put',
    data: { status }
  })
}

// 调整佣金比例
export const adjustCommission = (id, commissionRate) => {
  return request({
    url: `/v1/admin/agents/${id}/commission`,
    method: 'put',
    data: { commissionRate }
  })
}

// 获取代理商详情
export const getAgentDetail = (id) => {
  return request({
    url: `/v1/admin/agents/${id}/detail`,
    method: 'get'
  })
}

// 获取代理商子用户
export const getAgentSubUsers = (id, params) => {
  return request({
    url: `/v1/admin/agents/${id}/sub-users`,
    method: 'get',
    params
  })
}

// 结算佣金
export const settleCommission = (id) => {
  return request({
    url: `/v1/admin/commissions/${id}/settle`,
    method: 'post'
  })
}
