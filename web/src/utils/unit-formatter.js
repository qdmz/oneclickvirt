/**
 * 单位格式化工具函数
 * 统一处理MB、Mbps等单位的显示格式
 */

/**
 * 格式化内存大小显示
 * @param {number} sizeMB - 内存大小，单位MB
 * @returns {string} 格式化后的显示文本
 */
export function formatMemorySize(sizeMB) {
  if (sizeMB < 1024) {
    return `${sizeMB}MB`
  }
  const sizeGB = sizeMB / 1024
  if (sizeGB === Math.floor(sizeGB)) {
    return `${Math.floor(sizeGB)}GB`
  }
  return `${sizeGB.toFixed(1)}GB`
}

/**
 * 格式化磁盘大小显示
 * @param {number} sizeMB - 磁盘大小，单位MB
 * @returns {string} 格式化后的显示文本
 */
export function formatDiskSize(sizeMB) {
  if (sizeMB < 1024) {
    return `${sizeMB}MB`
  }
  const sizeGB = sizeMB / 1024
  if (sizeGB === Math.floor(sizeGB)) {
    return `${Math.floor(sizeGB)}GB`
  }
  return `${sizeGB.toFixed(1)}GB`
}

/**
 * 格式化带宽大小显示
 * @param {number} speedMbps - 带宽速度，单位Mbps
 * @returns {string} 格式化后的显示文本
 */
export function formatBandwidthSpeed(speedMbps) {
  if (speedMbps < 1000) {
    return `${speedMbps}Mbps`
  }
  const speedGbps = speedMbps / 1000
  if (speedGbps === Math.floor(speedGbps)) {
    return `${Math.floor(speedGbps)}Gbps`
  }
  return `${speedGbps.toFixed(1)}Gbps`
}

/**
 * 格式化资源使用量显示
 * @param {number} used - 已使用量
 * @param {number} total - 总量
 * @param {string} unit - 单位类型 'memory'|'disk'|'bandwidth'
 * @returns {string} 格式化后的显示文本
 */
export function formatResourceUsage(used, total, unit) {
  let formattedUsed, formattedTotal
  
  switch (unit) {
    case 'memory':
      formattedUsed = formatMemorySize(used)
      formattedTotal = formatMemorySize(total)
      break
    case 'disk':
      formattedUsed = formatDiskSize(used)
      formattedTotal = formatDiskSize(total)
      break
    case 'bandwidth':
      formattedUsed = formatBandwidthSpeed(used)
      formattedTotal = formatBandwidthSpeed(total)
      break
    default:
      formattedUsed = used
      formattedTotal = total
  }
  
  return `${formattedUsed} / ${formattedTotal}`
}
