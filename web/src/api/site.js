import request from '@/utils/request'

// ========== 站点配置 API ==========

/**
 * 获取站点配置（用户端公开）
 * @returns {Promise}
 */
export const getSiteConfig = () => {
  return request({
    url: '/v1/public/site-config',
    method: 'get'
  })
}

/**
 * 获取站点配置（管理员端）
 * @returns {Promise}
 */
export const getAdminSiteConfig = () => {
  return request({
    url: '/v1/admin/site-config',
    method: 'get'
  })
}

/**
 * 更新站点配置（管理员端）
 * @param {Object} data - 配置数据
 * @returns {Promise}
 */
export const updateAdminSiteConfig = (data) => {
  return request({
    url: '/v1/admin/site-config',
    method: 'put',
    data
  })
}

/**
 * 上传Logo
 * @param {FormData} formData - 包含logo文件的表单数据
 * @returns {Promise}
 */
export const uploadLogo = (formData) => {
  return request({
    url: '/v1/admin/site-config/logo',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
