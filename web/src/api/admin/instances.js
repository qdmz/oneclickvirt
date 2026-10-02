import request from '@/utils/request'
import { createLongTimeoutRequest } from '@/utils/longTimeoutRequest'

// 创建实例专用请求实例（120秒超时）
const instanceOperationRequest = createLongTimeoutRequest(120000, {
  requestPrefix: 'instance'
})

export const getAllInstances = (params) => {
  return request({
    url: '/v1/admin/instances',
    method: 'get',
    params
  })
}

export const getAdminInstance = (id) => {
  return request({
    url: `/v1/admin/instances/${id}`,
    method: 'get'
  })
}

export const createInstance = (data) => {
  return instanceOperationRequest({
    url: '/v1/admin/instances',
    method: 'post',
    data
  })
}

export const updateInstance = (id, data) => {
  return request({
    url: `/v1/admin/instances/${id}`,
    method: 'put',
    data
  })
}

export const deleteInstance = (id) => {
  return instanceOperationRequest({
    url: `/v1/admin/instances/${id}`,
    method: 'delete'
  })
}

export const adminInstanceAction = (id, action, extraData = {}) => {
  return instanceOperationRequest({
    url: `/v1/admin/instances/${id}/action`,
    method: 'post',
    data: { action, ...extraData }
  })
}

export const adminBatchInstanceAction = (instanceIds, action) => {
  return instanceOperationRequest({
    url: '/v1/admin/instances/batch-action',
    method: 'post',
    data: { instanceIds, action }
  })
}

export const getInstanceEgress = (id) => {
  return request({
    url: `/v1/admin/instances/${id}/egress`,
    method: 'get'
  })
}

export const bindInstanceEgress = (id, data) => {
  return instanceOperationRequest({
    url: `/v1/admin/instances/${id}/egress`,
    method: 'put',
    data
  })
}

export const unbindInstanceEgress = (id, apply = true) => {
  return instanceOperationRequest({
    url: `/v1/admin/instances/${id}/egress`,
    method: 'delete',
    params: { apply }
  })
}

export const reconcileInstanceEgress = (id, apply = true) => {
  return instanceOperationRequest({
    url: `/v1/admin/instances/${id}/egress/reconcile`,
    method: 'post',
    data: { apply }
  })
}

export const ensureInstanceEgressDependencies = (id, packageSet = 'wireguard') => {
  return instanceOperationRequest({
    url: `/v1/admin/instances/${id}/egress/dependencies`,
    method: 'post',
    data: { package_set: packageSet }
  })
}

export const createAdminInstanceShare = (id, data) => {
  return request({
    url: `/v1/admin/instances/${id}/share-links`,
    method: 'post',
    data
  })
}

export const resetInstancePassword = (id) => {
  return request({
    url: `/v1/admin/instances/${id}/reset-password`,
    method: 'put',
    data: {} // 空数据，后端会自动生成密码
  })
}

export const getAdminInstanceNewPassword = (instanceId, taskId) => {
  return request({
    url: `/v1/admin/instances/${instanceId}/password/${taskId}`,
    method: 'get'
  })
}

// 获取实例类型权限配置
export const getInstanceTypePermissions = () => {
  return request({
    url: '/v1/admin/instance-type-permissions',
    method: 'get'
  })
}

// 更新实例类型权限配置
export const updateInstanceTypePermissions = (data) => {
  return request({
    url: '/v1/admin/instance-type-permissions',
    method: 'put',
    data
  })
}

// 实例冻结管理
export const setInstanceExpiry = (data) => {
  return request({
    url: '/v1/admin/instances/set-expiry',
    method: 'post',
    data
  })
}

export const freezeInstance = (data) => {
  return request({
    url: '/v1/admin/instances/freeze',
    method: 'post',
    data
  })
}

export const unfreezeInstance = (data) => {
  return request({
    url: '/v1/admin/instances/unfreeze',
    method: 'post',
    data
  })
}

// 管理员 VNC 控制台
export const getAdminInstanceVNCInfo = (instanceId) => request({
  url: `/v1/admin/instances/${instanceId}/vnc`,
  method: 'get'
})

export const getAdminInstanceVNCWsUrl = (instanceId) => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  let host = window.location.host
  if (import.meta.env.MODE === 'development' && import.meta.env.VITE_SERVER_PORT) {
    host = `${window.location.hostname}:${import.meta.env.VITE_SERVER_PORT}`
  }
  const token = sessionStorage.getItem('token') || ''
  const query = token ? `?token=${encodeURIComponent(token)}` : ''
  return `${protocol}//${host}/api/v1/admin/instances/${instanceId}/vnc/ws${query}`
}
