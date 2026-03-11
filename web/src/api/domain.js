import request from '@/utils/request'

// 用户域名管理
export function getDomains(params) {
  return request({
    url: '/v1/user/domains',
    method: 'get',
    params
  })
}

export function createDomain(data) {
  return request({
    url: '/v1/user/domains',
    method: 'post',
    data
  })
}

export function updateDomain(id, data) {
  return request({
    url: `/v1/user/domains/${id}`,
    method: 'put',
    data
  })
}

export function deleteDomain(id) {
  return request({
    url: `/v1/user/domains/${id}`,
    method: 'delete'
  })
}

export function getDomainQuota() {
  return request({
    url: '/v1/user/domains/quota',
    method: 'get'
  })
}

// 管理员域名管理
export function adminGetDomains(params) {
  return request({
    url: '/v1/admin/domains',
    method: 'get',
    params
  })
}

export function adminDeleteDomain(id) {
  return request({
    url: `/v1/admin/domains/${id}`,
    method: 'delete'
  })
}

export function getDomainConfig() {
  return request({
    url: '/v1/admin/domain-config',
    method: 'get'
  })
}

export function updateDomainConfig(data) {
  return request({
    url: '/v1/admin/domain-config',
    method: 'put',
    data
  })
}

export function syncDNS() {
  return request({
    url: '/v1/admin/domains/sync-dns',
    method: 'post'
  })
}
