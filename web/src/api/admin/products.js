import request from '@/utils/request'

// 管理员：产品列表（含下架）
export const getAdminProducts = (params) => {
  return request({
    url: '/v1/admin/products',
    method: 'get',
    params
  })
}

// 管理员：创建产品
export const createAdminProduct = (data) => {
  return request({
    url: '/v1/admin/products',
    method: 'post',
    data
  })
}

// 管理员：更新产品
export const updateAdminProduct = (id, data) => {
  return request({
    url: `/v1/admin/products/${id}`,
    method: 'put',
    data
  })
}

// 管理员：删除产品
export const deleteAdminProduct = (id) => {
  return request({
    url: `/v1/admin/products/${id}`,
    method: 'delete'
  })
}
