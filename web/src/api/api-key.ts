import request from '@/utils/request'

/**
 * API 密钥接口
 */

// 创建 API 密钥
export function createAPIKey(data: {
  name: string
  expiresAt?: string
}) {
  return request({
    url: '/api/v1/user/api-keys',
    method: 'post',
    data
  })
}

// 获取 API 密钥列表
export function getAPIKeys() {
  return request({
    url: '/api/v1/user/api-keys',
    method: 'get'
  })
}

// 获取 API 密钥详情
export function getAPIKey(id: number) {
  return request({
    url: `/api/v1/user/api-keys/${id}`,
    method: 'get'
  })
}

// 更新 API 密钥
export function updateAPIKey(id: number, data: {
  name?: string
  status?: string
}) {
  return request({
    url: `/api/v1/user/api-keys/${id}`,
    method: 'put',
    data
  })
}

// 删除 API 密钥
export function deleteAPIKey(id: number) {
  return request({
    url: `/api/v1/user/api-keys/${id}`,
    method: 'delete'
  })
}

// 撤销 API 密钥
export function revokeAPIKey(id: number) {
  return request({
    url: `/api/v1/user/api-keys/${id}/revoke`,
    method: 'put'
  })
}

// 获取 API 密钥统计信息
export function getAPIKeyStats() {
  return request({
    url: '/api/v1/user/api-keys/stats',
    method: 'get'
  })
}
