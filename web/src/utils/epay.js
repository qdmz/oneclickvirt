/**
 * 易支付跳转工具
 * 后端返回 { epayUrl, epay: { pid, out_trade_no, name, money, sign, ... } }
 * 这里构造一个隐藏表单并自动 POST 提交，跳转到支付网关。
 */

/**
 * 提交到易支付网关
 * @param {string} epayUrl 网关地址（submit.php）
 * @param {Object} params 签名后的参数
 * @param {string} target 打开方式，默认当前页跳转；传 '_blank' 则新窗口
 */
export function submitToEpay(epayUrl, params, target = '_self') {
  if (!epayUrl || !params) return false

  const form = document.createElement('form')
  form.method = 'POST'
  form.action = epayUrl
  form.target = target
  form.style.display = 'none'

  Object.keys(params).forEach((key) => {
    const val = params[key]
    if (val === undefined || val === null) return
    const input = document.createElement('input')
    input.type = 'hidden'
    input.name = key
    input.value = String(val)
    form.appendChild(input)
  })

  document.body.appendChild(form)
  form.submit()
  // 提交后清理（跳转通常已发生，保险起见延迟移除）
  setTimeout(() => {
    if (form.parentNode) form.parentNode.removeChild(form)
  }, 1000)
  return true
}

/**
 * 获取当前站点基础地址，用于后端拼接支付回调 return_url
 */
export function currentBaseURL() {
  return `${window.location.protocol}//${window.location.host}`
}

/**
 * 分 -> 元 显示
 */
export function fenToYuan(cents) {
  const n = Number(cents || 0)
  return (n / 100).toFixed(2)
}
