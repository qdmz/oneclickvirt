import request from '@/utils/request'

export function getDashboardStats() {
  return request({
    url: '/v1/dashboard/stats',
    method: 'get'
  })
}

export function getRecentActivities() {
  return request({
    url: '/v1/dashboard/activities',
    method: 'get'
  })
}

export function getResourceTrends(timeRange = '24h') {
  return request({
    url: '/v1/dashboard/resource-trends',
    method: 'get',
    params: { timeRange }
  })
}

export function getSystemStatus() {
  return request({
    url: '/v1/dashboard/system-status',
    method: 'get'
  })
}