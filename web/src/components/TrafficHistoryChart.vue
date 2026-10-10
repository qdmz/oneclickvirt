<template>
  <div class="traffic-history-chart">
    <el-card>
      <template #header>
        <div class="chart-header">
          <span v-if="title">{{ title }}</span>
          <div class="chart-controls">
            <slot name="extra-actions" />
            <span class="control-label">{{ $t('user.traffic.historyChart.timeRange') }}:</span>
            <el-select
              v-model="selectedPeriod"
              size="small"
              class="history-select"
              @change="loadData"
            >
              <el-option
                :label="$t('user.traffic.historyChart.period15m')"
                value="15m"
              />
              <el-option
                :label="$t('user.traffic.historyChart.period30m')"
                value="30m"
              />
              <el-option
                :label="$t('user.traffic.historyChart.period1h')"
                value="1h"
              />
              <el-option
                :label="$t('user.traffic.historyChart.period6h')"
                value="6h"
              />
              <el-option
                :label="$t('user.traffic.historyChart.period12h')"
                value="12h"
              />
              <el-option
                :label="$t('user.traffic.historyChart.period24h')"
                value="24h"
              />
            </el-select>
            <span class="control-label">{{ $t('user.traffic.historyChart.dataInterval') }}:</span>
            <el-select
              v-model="selectedInterval"
              size="small"
              class="history-select"
              @change="loadData"
            >
              <el-option
                :label="$t('user.traffic.historyChart.interval5m')"
                :value="5"
              />
              <el-option
                :label="$t('user.traffic.historyChart.interval10m')"
                :value="10"
              />
              <el-option
                :label="$t('user.traffic.historyChart.interval15m')"
                :value="15"
              />
              <el-option
                :label="$t('user.traffic.historyChart.interval30m')"
                :value="30"
              />
            </el-select>
          </div>
        </div>
      </template>

      <div
        v-show="loading"
        class="chart-loading"
      >
        <el-skeleton
          :rows="8"
          animated
        />
      </div>
      
      <div
        v-show="!loading"
        ref="chartRef"
        class="chart-container"
        style="width: 100%; height: 400px;"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'
import request from '@/utils/request'
import { getUserTrafficHistory, getInstanceTrafficHistory } from '@/api/user'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/pinia/modules/user'
import { withChartTheme, useChartTheme, areaGradient, getChartTokens } from '@/utils/chartTheme'

const { t, locale } = useI18n()
const userStore = useUserStore()

const props = defineProps({
  // 'instance', 'provider', 'user'
  type: {
    type: String,
    required: true,
    validator: (value) => ['instance', 'provider', 'user'].includes(value)
  },
  // 资源ID (instance_id, provider_id, user_id)
  resourceId: {
    type: [Number, String],
    default: null
  },
  // 图表标题
  title: {
    type: String,
    default: ''
  },
  // 自动刷新间隔（秒），0表示不自动刷新
  autoRefresh: {
    type: Number,
    default: 0
  }
})

const chartRef = ref(null)
const loading = ref(false)
const error = ref('')
const selectedPeriod = ref('1h')   // 默认1小时
const selectedInterval = ref(5)    // 默认5分钟间隔
const chartInstance = ref(null)
const refreshTimer = ref(null)
const chartData = ref([])          // 存储图表数据，用于语言切换时重新渲染

const normalizeHistoryData = (payload) => {
  if (Array.isArray(payload)) return payload
  if (Array.isArray(payload?.list)) return payload.list
  if (Array.isArray(payload?.data)) return payload.data
  if (Array.isArray(payload?.history)) return payload.history
  return []
}

// 格式化流量单位
const formatTraffic = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${units[i]}`
}

// 格式化时间标签（支持分钟级精度）
const formatTimeLabel = (record) => {
  // 优先使用 record_time 字段
  if (record.record_time) {
    const date = new Date(record.record_time)
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hour = String(date.getHours()).padStart(2, '0')
    const minute = String(date.getMinutes()).padStart(2, '0')
    
    // 根据时间范围调整显示格式
    if (selectedPeriod.value === '5m' || selectedPeriod.value === '10m' || selectedPeriod.value === '15m') {
      // 短时间范围：只显示 时:分
      return `${hour}:${minute}`
    } else if (selectedPeriod.value === '30m' || selectedPeriod.value === '45m' || selectedPeriod.value === '1h') {
      // 中等时间范围：显示 月-日 时:分
      return `${month}-${day} ${hour}:${minute}`
    } else {
      // 长时间范围：显示 月-日 时:分
      return `${month}-${day} ${hour}:${minute}`
    }
  }
  
  // 回退到使用分散的字段
  const month = String(record.month).padStart(2, '0')
  const day = String(record.day).padStart(2, '0')
  const hour = String(record.hour).padStart(2, '0')
  const minute = String(record.minute || 0).padStart(2, '0')
  
  // 根据时间范围调整显示格式
  if (selectedPeriod.value === '5m' || selectedPeriod.value === '10m' || selectedPeriod.value === '15m') {
    // 短时间范围：只显示 时:分
    return `${hour}:${minute}`
  } else if (selectedPeriod.value === '30m' || selectedPeriod.value === '45m' || selectedPeriod.value === '1h') {
    // 中等时间范围：显示 月-日 时:分
    return `${month}-${day} ${hour}:${minute}`
  } else {
    // 长时间范围：显示 月-日 时:分
    return `${month}-${day} ${hour}:${minute}`
  }
}

// 加载流量历史数据
const loadData = async () => {
  if (loading.value) return
  
  // 检查用户是否已登录
  if (!userStore.isLoggedIn) {
    console.warn('Traffic history: User not logged in, skipping data load')
    loading.value = false
    await nextTick()
    renderChart([])
    return
  }
  
  loading.value = true
  error.value = ''
  
  try {
    let response
    const params = {
      period: selectedPeriod.value,
      interval: selectedInterval.value
    }
    
    switch (props.type) {
      case 'instance':
        if (!props.resourceId) {
          throw new Error('Instance ID is required')
        }
        response = await getInstanceTrafficHistory(props.resourceId, params)
        break
      case 'provider':
        if (!props.resourceId) {
          throw new Error('Provider ID is required')
        }
        // Provider 使用admin API，仍然需要直接调用
        response = await request({
          url: `/v1/admin/providers/${props.resourceId}/traffic/history`,
          method: 'get',
          params
        })
        break
      case 'user':
        response = await getUserTrafficHistory(params)
        break
      default:
        throw new Error('Invalid type')
    }
    if (response && (response.code === 200)) {
      loading.value = false
      // 存储数据供语言切换时使用
      chartData.value = normalizeHistoryData(response.data)
      // 等待DOM更新后再渲染图表
      await nextTick()
      renderChart(chartData.value)
    } else {
      throw new Error(response?.message || response?.msg || 'Failed to load data')
    }
  } catch (err) {
    console.error('Load traffic history failed:', err)
    loading.value = false
    // 加载失败时显示全零曲线图，避免显示错误信息
    chartData.value = []
    await nextTick()
    renderChart([])
  }
}

// 生成零填充数据点（当无数据时显示全零曲线图）
const generateZeroFilledData = () => {
  const now = new Date()
  const periodMinutes = {
    '15m': 15, '30m': 30, '1h': 60, '6h': 360, '12h': 720, '24h': 1440
  }
  const minutes = periodMinutes[selectedPeriod.value] || 60
  const interval = selectedInterval.value || 5
  const points = []
  
  for (let i = minutes; i >= 0; i -= interval) {
    const time = new Date(now.getTime() - i * 60 * 1000)
    points.push({
      record_time: time.toISOString(),
      traffic_in: 0,
      traffic_out: 0,
      total_used: 0
    })
  }
  
  return points
}

// 渲染图表
const renderChart = (data) => {
  if (!chartRef.value) {
    console.warn('TrafficHistoryChart - chartRef is null')
    return
  }
  data = normalizeHistoryData(data)
  
  // 如果无数据，生成零填充数据点（显示全零曲线图）
  if (!data || data.length === 0) {
    data = generateZeroFilledData()
  }
  // 清除错误状态
  error.value = ''
  
  // 初始化或重新初始化图表实例
  // 如果图表实例已存在但DOM已销毁，需要重新初始化
  if (!chartInstance.value || chartInstance.value.isDisposed()) {
    if (chartInstance.value) {
      chartInstance.value.dispose()
    }
    chartInstance.value = echarts.init(chartRef.value)
  }
  
  // 准备数据
  const timeLabels = data.map(item => formatTimeLabel(item))
  // 支持多种数据格式：
  // 1. traffic_mb: MB单位（新格式）
  // 2. traffic_in/traffic_out/total_used: MB单位（实例详细数据）
  // 3. rx_bytes/tx_bytes/total_bytes: 字节单位（原始pmacct数据）
  const trafficIn = data.map(item => {
    if (item.traffic_in !== undefined) return parseFloat(item.traffic_in || 0).toFixed(2) // 已经是MB
    if (item.rx_bytes !== undefined) return ((item.rx_bytes || 0) / 1024 / 1024).toFixed(2) // 字节转MB
    return '0.00'
  })
  const trafficOut = data.map(item => {
    if (item.traffic_out !== undefined) return parseFloat(item.traffic_out || 0).toFixed(2)
    if (item.tx_bytes !== undefined) return ((item.tx_bytes || 0) / 1024 / 1024).toFixed(2)
    return '0.00'
  })
  const totalUsed = data.map(item => {
    if (item.total_used !== undefined) return parseFloat(item.total_used || 0).toFixed(2)
    if (item.traffic_mb !== undefined) return parseFloat(item.traffic_mb || 0).toFixed(2) // 历史数据
    if (item.total_bytes !== undefined) return ((item.total_bytes || 0) / 1024 / 1024).toFixed(2)
    return '0.00'
  })
  // 配置图表选项（颜色 / 坐标轴 / 提示框样式来自 chartTheme，随亮暗主题切换）
  const tk = getChartTokens()
  const seriesDefs = [
    { name: t('user.traffic.historyChart.inbound'), data: trafficIn, color: tk.accent },
    { name: t('user.traffic.historyChart.outbound'), data: trafficOut, color: tk.violet },
    { name: t('user.traffic.historyChart.total'), data: totalUsed, color: tk.primary }
  ]
  const option = withChartTheme({
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'line' },
      formatter: (params) => {
        let result = `<div style="margin-bottom:4px;opacity:.75">${params[0].axisValue}</div>`
        params.forEach(item => {
          const value = parseFloat(item.value)
          result += `${item.marker} ${item.seriesName}: <b>${formatTraffic(value * 1024 * 1024)}</b><br/>`
        })
        return result
      }
    },
    legend: {
      data: seriesDefs.map(s => s.name),
      top: 4,
      left: 'center'
    },
    grid: {
      top: 44,
      left: 8,
      right: 16,
      bottom: 8,
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: timeLabels,
      axisLabel: {
        rotate: timeLabels.length > 12 ? 30 : 0,
        interval: 'auto',
        hideOverlap: true
      }
    },
    yAxis: {
      type: 'value',
      name: 'MB',
      axisLabel: {
        formatter: (value) => {
          if (value >= 1024) {
            return `${(value / 1024).toFixed(1)} GB`
          }
          return `${value} MB`
        }
      }
    },
    series: seriesDefs.map((def, idx) => ({
      name: def.name,
      type: 'line',
      smooth: true,
      showSymbol: false,
      symbolSize: 6,
      data: def.data,
      emphasis: { focus: 'series' },
      itemStyle: { color: def.color },
      lineStyle: { width: idx === 2 ? 2.5 : 2, color: def.color },
      areaStyle: { color: areaGradient(def.color, idx === 2 ? 0.3 : 0.18, 0.01) }
    }))
  })
  chartInstance.value.setOption(option, true)
}

// 亮/暗主题切换时重新渲染（chartData 为空时 renderChart 会生成零填充数据）
useChartTheme(() => {
  if (!loading.value && chartInstance.value) renderChart(chartData.value)
})

// 窗口大小改变时重新渲染
const handleResize = () => {
  if (chartInstance.value) {
    chartInstance.value.resize()
  }
}

// 设置自动刷新
const setupAutoRefresh = () => {
  if (props.autoRefresh > 0) {
    refreshTimer.value = setInterval(() => {
      loadData()
    }, props.autoRefresh * 1000)
  }
}

// 清除自动刷新
const clearAutoRefresh = () => {
  if (refreshTimer.value) {
    clearInterval(refreshTimer.value)
    refreshTimer.value = null
  }
}

// 监听资源ID变化
watch(() => props.resourceId, () => {
  if (props.resourceId) {
    loadData()
  }
})

// 监听自动刷新配置变化
watch(() => props.autoRefresh, () => {
  clearAutoRefresh()
  setupAutoRefresh()
})

// 监听语言变化，重新渲染图表
watch(() => locale.value, () => {
  // 当语言切换时，如果图表已经有数据，重新渲染
  if (chartData.value && chartData.value.length > 0) {
    renderChart(chartData.value)
  }
})

onMounted(async () => {
  await nextTick()
  window.addEventListener('resize', handleResize)
  if (props.resourceId || props.type === 'user') {
    loadData()
  }
  setupAutoRefresh()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  clearAutoRefresh()
  if (chartInstance.value) {
    chartInstance.value.dispose()
    chartInstance.value = null
  }
})

defineExpose({
  refresh: loadData
})
</script>

<style scoped lang="scss">
.traffic-history-chart {
  margin-top: 20px;

  .chart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
    min-width: 0;

    > span {
      min-width: 0;
      overflow-wrap: anywhere;
    }

    .chart-controls {
      display: flex;
      align-items: center;
      justify-content: flex-end;
      flex-wrap: wrap;
      gap: 8px;
      min-width: 0;

      .control-label {
        font-size: 14px;
        white-space: nowrap;
      }

      .history-select {
        width: 120px;
      }
    }
  }

  .chart-loading {
    height: 400px;
    padding: 12px 4px;
  }

  .chart-container {
    min-height: 400px;
  }

  .control-label {
    color: var(--text-color-tertiary);
    font-size: 12.5px !important;
  }

  @media (max-width: 768px) {
    .chart-header {
      align-items: flex-start;
      flex-direction: column;

      .chart-controls {
        width: 100%;
        justify-content: flex-start;

        .history-select {
          flex: 1 1 120px;
          width: auto;
          max-width: 100%;
        }
      }
    }
  }
}
</style>
