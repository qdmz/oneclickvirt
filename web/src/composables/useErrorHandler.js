import { ref, readonly } from 'vue'
import { errorHandler } from '@/utils/errorHandler'

/**
 * 全局错误处理composable
 * 为Vue组件提供统一的错误处理能力
 */
export function useErrorHandler() {
  const loading = ref(false)
  const error = ref(null)

  /**
   * 包装async函数调用，自动处理loading和错误
   * @param {Function} asyncFn - 异步函数
   * @param {Object} options - 配置选项
   * @returns {Promise} 执行结果
   */
  const executeAsync = async (asyncFn, options = {}) => {
    const {
      showLoading = true,
      showError = true,
      successMessage = null,
      errorOptions = {}
    } = options

    if (showLoading) {
      loading.value = true
    }
    error.value = null

    try {
      const result = await asyncFn()
      
      // 显示成功消息
      if (successMessage) {
        errorHandler.showErrorMessage(successMessage, 'success')
      }
      
      return { success: true, data: result, error: null }
    } catch (err) {
      const errorInfo = errorHandler.handleApiError(err, {
        showMessage: showError,
        ...errorOptions
      })
      error.value = errorInfo
      return { success: false, data: null, error: errorInfo }
    } finally {
      if (showLoading) {
        loading.value = false
      }
    }
  }

  /**
   * 处理表单提交
   * @param {Function} submitFn - 提交函数
   * @param {Object} options - 配置选项
   * @returns {Promise} 提交结果
   */
  const handleSubmit = async (submitFn, options = {}) => {
    const {
      successMessage = '操作成功',
      confirmMessage = null,
      ...otherOptions
    } = options

    // 如果需要确认，先显示确认对话框
    if (confirmMessage) {
      const confirmed = await errorHandler.showConfirmDialog(confirmMessage)
      if (!confirmed) {
        return { success: false, cancelled: true }
      }
    }

    return executeAsync(submitFn, {
      successMessage,
      ...otherOptions
    })
  }

  /**
   * 处理删除操作
   * @param {Function} deleteFn - 删除函数
   * @param {Object} options - 配置选项
   * @returns {Promise} 删除结果
   */
  const handleDelete = async (deleteFn, options = {}) => {
    const {
      confirmMessage = '确定要删除吗？此操作不可撤销。',
      successMessage = '删除成功',
      ...otherOptions
    } = options

    return handleSubmit(deleteFn, {
      confirmMessage,
      successMessage,
      ...otherOptions
    })
  }

  /**
   * 处理批量操作
   * @param {Array} items - 要处理的项目列表
   * @param {Function} itemProcessor - 处理单个项目的函数
   * @param {Object} options - 配置选项
   * @returns {Promise} 批量处理结果
   */
  const handleBatchOperation = async (items, itemProcessor, options = {}) => {
    const {
      confirmMessage = `确定要处理 ${items.length} 个项目吗？`,
      successMessage = '批量操作完成',
      continueOnError = false,
      ...otherOptions
    } = options

    // 确认批量操作
    if (confirmMessage) {
      const confirmed = await errorHandler.showConfirmDialog(confirmMessage)
      if (!confirmed) {
        return { success: false, cancelled: true }
      }
    }

    return executeAsync(async () => {
      const results = []
      const errors = []

      for (let i = 0; i < items.length; i++) {
        try {
          const result = await itemProcessor(items[i], i)
          results.push({ item: items[i], result, success: true })
        } catch (err) {
          const errorInfo = errorHandler.handleApiError(err, { showMessage: false })
          errors.push({ item: items[i], error: errorInfo, success: false })
          
          if (!continueOnError) {
            throw err
          }
        }
      }

      return {
        results,
        errors,
        successCount: results.length,
        errorCount: errors.length,
        total: items.length
      }
    }, {
      successMessage,
      ...otherOptions
    })
  }

  /**
   * 清除错误状态
   */
  const clearError = () => {
    error.value = null
  }

  /**
   * 设置loading状态
   * @param {boolean} value - loading值
   */
  const setLoading = (value) => {
    loading.value = value
  }

  return {
    // 响应式状态
    loading: readonly(loading),
    error: readonly(error),
    
    // 方法
    executeAsync,
    handleSubmit,
    handleDelete,
    handleBatchOperation,
    clearError,
    setLoading,
    
    // 直接暴露errorHandler的方法
    showErrorMessage: errorHandler.showErrorMessage,
    showConfirmDialog: errorHandler.showConfirmDialog,
    handleValidationErrors: errorHandler.handleValidationErrors
  }
}

// 默认导出
export default useErrorHandler
