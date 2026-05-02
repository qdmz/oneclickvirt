// 购物车API模块
import request from '@/utils/request'

/**
 * 购物车管理API
 */

// 获取购物车列表
export const getCartItems = (params) => {
  return request({
    url: '/v1/user/cart',
    method: 'get',
    params
  })
}

// 添加商品到购物车
export const addToCart = (productId, quantity = 1) => {
  return request({
    url: '/v1/user/cart/items',
    method: 'post',
    data: { productId, quantity }
  })
}

// 从购物车移除商品
export const removeFromCart = (itemId) => {
  return request({
    url: `/v1/user/cart/items/${itemId}`,
    method: 'delete'
  })
}

// 更新购物车商品数量
export const updateCartItem = (itemId, quantity) => {
  return request({
    url: `/v1/user/cart/items/${itemId}`,
    method: 'put',
    data: { quantity }
  })
}

// 清空购物车
export const clearCart = () => {
  return request({
    url: '/v1/user/cart/clear',
    method: 'delete'
  })
}

// 获取购物车总金额
export const getCartTotal = () => {
  return request({
    url: '/v1/user/cart/total',
    method: 'get'
  })
}

// 检查商品是否在购物车中
export const checkItemInCart = (productId) => {
  return request({
    url: `/v1/user/cart/check/${productId}`,
    method: 'get'
  })
}
