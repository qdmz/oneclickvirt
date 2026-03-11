import request from '@/utils/request'

export function submitKYC(data) {
  return request({
    url: '/v1/user/kyc/submit',
    method: 'post',
    data
  })
}

export function getKYCStatus() {
  return request({
    url: '/v1/user/kyc/status',
    method: 'get'
  })
}

export function queryKYC() {
  return request({
    url: '/v1/user/kyc/query',
    method: 'post'
  })
}

// Admin APIs
export function getKYCRecords(params) {
  return request({
    url: '/v1/admin/kyc/records',
    method: 'get',
    params
  })
}

export function updateKYCStatus(id, data) {
  return request({
    url: `/v1/admin/kyc/${id}/status`,
    method: 'put',
    data
  })
}

export function getKYCStats() {
  return request({
    url: '/v1/admin/kyc/stats',
    method: 'get'
  })
}
