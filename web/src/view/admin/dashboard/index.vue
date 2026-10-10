<template>
  <div class="oh-page admin-dashboard">
    <!-- 页头 -->
    <div class="oh-page-header">
      <div>
        <div class="oh-page-title">
          <span>{{ greeting }}</span>
          <span class="oh-gradient-text">{{ $t('admin.dashboard.title') }}</span>
        </div>
        <div class="oh-page-subtitle">
          {{ $t('ui.dashboard.adminSubtitle') }}
          <span
            v-if="lastUpdated"
            class="updated-at"
          >· {{ $t('ui.dashboard.lastUpdated', { time: lastUpdated }) }}</span>
        </div>
      </div>
      <div class="oh-page-actions">
        <el-button
          :loading="loading"
          @click="handleManualRefresh"
        >
          <el-icon><Refresh /></el-icon>
          {{ $t('common.refresh') }}
        </el-button>
      </div>
    </div>

    <!-- 首次加载骨架屏 -->
    <template v-if="loading && !loaded">
      <div class="oh-grid oh-grid-4 section">
        <div
          v-for="n in 4"
          :key="n"
          class="oh-surface skeleton-card"
        >
          <el-skeleton
            :rows="2"
            animated
          />
        </div>
      </div>
      <div class="oh-grid oh-grid-2 section">
        <div
          v-for="n in 2"
          :key="n"
          class="oh-surface skeleton-panel"
        >
          <el-skeleton
            :rows="6"
            animated
          />
        </div>
      </div>
    </template>

    <template v-else>
      <!-- 关键指标 -->
      <div class="oh-grid oh-grid-4 section">
        <div
          v-for="card in statCards"
          :key="card.key"
          class="oh-surface oh-stat stat-card"
          :class="`oh-tone-${card.tone}`"
        >
          <div class="oh-stat-top">
            <span class="oh-stat-label">{{ card.label }}</span>
            <span
              class="stat-icon-gradient"
              :class="`grad-${card.tone}`"
            >
              <i :class="card.icon" />
            </span>
          </div>
          <div class="oh-stat-value">
            {{ card.value }}
          </div>
          <div class="oh-stat-foot">
            <span
              v-if="card.dot"
              class="status-dot"
              :class="card.dot"
            />
            {{ card.foot }}
          </div>
        </div>
      </div>

      <!-- 服务器状态 + 实例分布 -->
      <div class="oh-grid oh-grid-2 section">
        <div class="oh-surface panel">
          <div class="panel-header">
            <span class="oh-section-title">
              <i class="fas fa-heartbeat title-icon" />
              {{ $t('ui.dashboard.serverStatus') }}
            </span>
            <el-tag
              :type="healthTag.type"
              effect="light"
              round
              size="small"
            >
              {{ healthTag.text }}
            </el-tag>
          </div>
          <div class="rings">
            <div
              v-for="ring in systemRings"
              :key="ring.key"
              class="ring-item"
            >
              <el-progress
                type="circle"
                :percentage="ring.value"
                :width="112"
                :stroke-width="9"
                :color="usageColor(ring.value)"
              >
                <template #default>
                  <div class="ring-inner">
                    <span class="ring-value">{{ ring.value.toFixed(1) }}<small>%</small></span>
                  </div>
                </template>
              </el-progress>
              <div class="ring-label">
                <i :class="ring.icon" />
                {{ ring.label }}
              </div>
            </div>
          </div>
          <div class="uptime-bar">
            <i class="far fa-clock" />
            <span>{{ $t('ui.dashboard.uptime') }}</span>
            <strong>{{ uptimeText }}</strong>
          </div>
        </div>

        <div class="oh-surface panel">
          <div class="panel-header">
            <span class="oh-section-title">
              <i class="fas fa-chart-pie title-icon" />
              {{ $t('ui.dashboard.instanceDistribution') }}
            </span>
            <span class="panel-meta">
              {{ $t('ui.dashboard.runningOfTotal', { running: dashboardData.runningInstances, total: totalInstances }) }}
            </span>
          </div>
          <div
            v-if="totalInstances > 0"
            ref="distributionChartRef"
            class="chart-box"
          />
          <div
            v-else
            class="oh-empty chart-empty"
          >
            <div class="oh-empty-icon">
              <i class="fas fa-cubes" />
            </div>
            {{ $t('ui.dashboard.noData') }}
          </div>
        </div>
      </div>

      <!-- 资源分配率 -->
      <div class="oh-surface panel section">
        <div class="panel-header">
          <span class="oh-section-title">
            <i class="fas fa-layer-group title-icon" />
            {{ $t('ui.dashboard.resourceAllocation') }}
          </span>
          <span class="panel-meta">{{ $t('admin.dashboard.resourceOverview') }}</span>
        </div>
        <div class="allocation">
          <div
            ref="allocationChartRef"
            class="chart-box allocation-chart"
          />
          <div class="allocation-list">
            <div
              v-for="item in allocationItems"
              :key="item.key"
              class="allocation-item"
              :class="`oh-tone-${item.tone}`"
            >
              <div class="allocation-head">
                <span class="allocation-icon"><i :class="item.icon" /></span>
                <span class="allocation-title">{{ item.label }}</span>
                <span class="allocation-percent">{{ item.percent }}%</span>
              </div>
              <el-progress
                :percentage="item.percent"
                :stroke-width="6"
                :show-text="false"
                :color="usageColor(item.percent)"
              />
              <div class="allocation-numbers">
                <span>{{ $t('ui.dashboard.allocated') }} <b>{{ item.used }}</b></span>
                <span>{{ $t('admin.dashboard.total') }} <b>{{ item.total }}</b></span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 快捷入口 -->
      <div class="section">
        <div class="oh-section-title quick-title">
          <i class="fas fa-bolt title-icon" />
          {{ $t('ui.dashboard.quickEntry') }}
        </div>
        <div class="quick-grid">
          <router-link
            v-for="link in quickLinks"
            :key="link.to"
            :to="link.to"
            class="oh-surface quick-link"
            :class="`oh-tone-${link.tone}`"
          >
            <span class="quick-icon"><i :class="link.icon" /></span>
            <span class="quick-label">{{ link.label }}</span>
            <i class="fas fa-arrow-right quick-arrow" />
          </router-link>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import * as echarts from 'echarts'
import { getAdminDashboard } from '@/api/admin'
import { withChartTheme, useChartTheme, getChartTokens } from '@/utils/chartTheme'

const { t, locale } = useI18n()

const loading = ref(false)
const loaded = ref(false)
const lastUpdated = ref('')

const dashboardData = reactive({
  totalUsers: 0,
  totalProviders: 0,
  totalVMs: 0,
  totalContainers: 0,
  activeUsers: 0,
  totalInstances: 0,
  runningInstances: 0,
  activeProviders: 0
})

const resourceUsage = reactive({
  totalCpuCores: 0,
  usedCpuCores: 0,
  totalMemoryMB: 0,
  usedMemoryMB: 0,
  totalDiskMB: 0,
  usedDiskMB: 0
})

const systemStatus = reactive({
  cpuUsage: 0,
  memoryUsage: 0,
  diskUsage: 0,
  uptime: ''
})

// ---------- 格式化 ----------
const formatMB = (mb) => {
  if (!mb || mb <= 0) return '0 MB'
  const GB = 1024
  const TB = 1024 * 1024
  if (mb >= TB) return (mb / TB).toFixed(2) + ' TB'
  if (mb >= GB) return (mb / GB).toFixed(2) + ' GB'
  return mb + ' MB'
}

const formatNumber = (n) => Number(n || 0).toLocaleString(locale.value === 'en-US' ? 'en-US' : 'zh-CN')

const clampPercent = (v) => {
  const n = Number(v)
  if (!Number.isFinite(n) || n < 0) return 0
  return Math.min(Math.round(n * 10) / 10, 100)
}

const percentOf = (used, total) => {
  if (!total || total <= 0) return 0
  return Math.min(Math.round((used / total) * 100), 100)
}

// Go time.Duration 字符串（如 "72h3m4.5s"）→ 本地化时长
const parseGoDuration = (str) => {
  if (!str || typeof str !== 'string') return null
  const re = /([\d.]+)(h|m(?!s)|s|ms|µs|us|ns)/g
  let seconds = 0
  let matched = false
  let m
  while ((m = re.exec(str)) !== null) {
    matched = true
    const v = parseFloat(m[1])
    switch (m[2]) {
      case 'h': seconds += v * 3600; break
      case 'm': seconds += v * 60; break
      case 's': seconds += v; break
      default: break
    }
  }
  return matched ? Math.floor(seconds) : null
}

const uptimeText = computed(() => {
  const secs = parseGoDuration(systemStatus.uptime)
  if (secs === null) return systemStatus.uptime || '—'
  const d = Math.floor(secs / 86400)
  const h = Math.floor((secs % 86400) / 3600)
  const min = Math.floor((secs % 3600) / 60)
  if (d > 0) return t('ui.dashboard.uptimeDHM', { d, h, m: min })
  if (h > 0) return t('ui.dashboard.uptimeHM', { h, m: min })
  return t('ui.dashboard.uptimeM', { m: Math.max(min, 0) })
})

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return t('ui.dashboard.goodMorning')
  if (hour < 18) return t('ui.dashboard.goodAfternoon')
  return t('ui.dashboard.goodEvening')
})

// 阈值颜色（CSS 变量，随主题切换）
const usageColor = (p) => {
  if (p >= 90) return 'var(--error-color)'
  if (p >= 70) return 'var(--warning-color)'
  return 'var(--primary-color)'
}

const totalInstances = computed(() =>
  dashboardData.totalInstances || (dashboardData.totalVMs + dashboardData.totalContainers)
)

// ---------- 卡片数据 ----------
const statCards = computed(() => [
  {
    key: 'activeUsers',
    tone: 'indigo',
    icon: 'fas fa-user-check',
    label: t('ui.dashboard.activeUsers'),
    value: formatNumber(dashboardData.activeUsers),
    foot: `${t('admin.dashboard.totalUsers')} ${formatNumber(dashboardData.totalUsers)}`
  },
  {
    key: 'runningInstances',
    tone: 'emerald',
    icon: 'fas fa-play-circle',
    label: t('admin.dashboard.runningInstances'),
    value: formatNumber(dashboardData.runningInstances),
    foot: `${t('ui.dashboard.totalInstances')} ${formatNumber(totalInstances.value)}`,
    dot: dashboardData.runningInstances > 0 ? 'dot-ok' : 'dot-idle'
  },
  {
    key: 'activeProviders',
    tone: 'sky',
    icon: 'fas fa-server',
    label: t('ui.dashboard.activeProviders'),
    value: formatNumber(dashboardData.activeProviders),
    foot: `${t('admin.dashboard.totalProviders')} ${formatNumber(dashboardData.totalProviders)}`,
    dot: dashboardData.activeProviders > 0 ? 'dot-ok' : 'dot-idle'
  },
  {
    key: 'uptime',
    tone: 'violet',
    icon: 'fas fa-stopwatch',
    label: t('ui.dashboard.uptime'),
    value: uptimeText.value,
    foot: t('ui.dashboard.serverStatus'),
    dot: healthLevel.value === 'ok' ? 'dot-ok' : (healthLevel.value === 'warn' ? 'dot-warn' : 'dot-bad')
  }
])

const systemRings = computed(() => [
  { key: 'cpu', icon: 'fas fa-microchip', label: t('ui.dashboard.cpuUsage'), value: clampPercent(systemStatus.cpuUsage) },
  { key: 'mem', icon: 'fas fa-memory', label: t('ui.dashboard.memoryUsage'), value: clampPercent(systemStatus.memoryUsage) },
  { key: 'disk', icon: 'fas fa-hdd', label: t('ui.dashboard.diskUsage'), value: clampPercent(systemStatus.diskUsage) }
])

const healthLevel = computed(() => {
  const max = Math.max(...systemRings.value.map((r) => r.value))
  if (max >= 90) return 'bad'
  if (max >= 70) return 'warn'
  return 'ok'
})

const healthTag = computed(() => ({
  ok: { type: 'success', text: t('ui.dashboard.healthOk') },
  warn: { type: 'warning', text: t('ui.dashboard.healthWarn') },
  bad: { type: 'danger', text: t('ui.dashboard.healthBad') }
}[healthLevel.value]))

const allocationItems = computed(() => [
  {
    key: 'cpu',
    tone: 'indigo',
    icon: 'fas fa-microchip',
    label: t('admin.dashboard.cpuCores'),
    used: formatNumber(resourceUsage.usedCpuCores),
    total: formatNumber(resourceUsage.totalCpuCores),
    percent: percentOf(resourceUsage.usedCpuCores, resourceUsage.totalCpuCores)
  },
  {
    key: 'memory',
    tone: 'violet',
    icon: 'fas fa-memory',
    label: t('admin.dashboard.memory'),
    used: formatMB(resourceUsage.usedMemoryMB),
    total: formatMB(resourceUsage.totalMemoryMB),
    percent: percentOf(resourceUsage.usedMemoryMB, resourceUsage.totalMemoryMB)
  },
  {
    key: 'disk',
    tone: 'amber',
    icon: 'fas fa-hdd',
    label: t('admin.dashboard.disk'),
    used: formatMB(resourceUsage.usedDiskMB),
    total: formatMB(resourceUsage.totalDiskMB),
    percent: percentOf(resourceUsage.usedDiskMB, resourceUsage.totalDiskMB)
  }
])

const quickLinks = computed(() => [
  { to: '/admin/users', icon: 'fas fa-users', tone: 'indigo', label: t('sidebar.userManagement') },
  { to: '/admin/providers', icon: 'fas fa-server', tone: 'sky', label: t('sidebar.providerManagement') },
  { to: '/admin/instances', icon: 'fas fa-cubes', tone: 'emerald', label: t('sidebar.instanceManagement') },
  { to: '/admin/orders', icon: 'fas fa-receipt', tone: 'amber', label: t('sidebar.orderManagement') },
  { to: '/admin/tickets', icon: 'fas fa-headset', tone: 'violet', label: t('sidebar.ticketManagement') },
  { to: '/admin/config', icon: 'fas fa-sliders-h', tone: 'rose', label: t('sidebar.systemConfiguration') }
])

// ---------- 图表 ----------
const distributionChartRef = ref(null)
const allocationChartRef = ref(null)
let distributionChart = null
let allocationChart = null
let resizeObserver = null

const ensureChart = (instance, el) => {
  if (!el) return null
  if (instance && !instance.isDisposed() && instance.getDom() === el) return instance
  if (instance && !instance.isDisposed()) instance.dispose()
  return echarts.init(el)
}

const renderDistributionChart = () => {
  if (!distributionChartRef.value) {
    if (distributionChart && !distributionChart.isDisposed()) distributionChart.dispose()
    distributionChart = null
    return
  }
  distributionChart = ensureChart(distributionChart, distributionChartRef.value)
  const tk = getChartTokens()
  const running = dashboardData.runningInstances || 0
  const total = totalInstances.value
  distributionChart.setOption(withChartTheme({
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { bottom: 0, left: 'center' },
    graphic: [
      {
        type: 'text',
        left: 'center',
        top: '36%',
        style: { text: formatNumber(total), fill: tk.text, fontSize: 26, fontWeight: 700, textAlign: 'center' }
      },
      {
        type: 'text',
        left: 'center',
        top: '50%',
        style: { text: t('ui.dashboard.totalInstances'), fill: tk.textMuted, fontSize: 12, textAlign: 'center' }
      }
    ],
    series: [
      {
        name: t('ui.dashboard.instanceDistribution'),
        type: 'pie',
        radius: ['58%', '78%'],
        center: ['50%', '44%'],
        avoidLabelOverlap: true,
        padAngle: 2,
        itemStyle: { borderRadius: 6, borderColor: tk.card, borderWidth: 2 },
        label: { show: false },
        emphasis: { scale: true, scaleSize: 4 },
        data: [
          { name: t('admin.dashboard.totalVMs'), value: dashboardData.totalVMs || 0, itemStyle: { color: tk.primary } },
          { name: t('admin.dashboard.totalContainers'), value: dashboardData.totalContainers || 0, itemStyle: { color: tk.accent } }
        ]
      },
      {
        name: t('ui.dashboard.running'),
        type: 'pie',
        radius: ['44%', '50%'],
        center: ['50%', '44%'],
        silent: false,
        label: { show: false },
        itemStyle: { borderColor: tk.card, borderWidth: 1 },
        data: [
          { name: t('ui.dashboard.running'), value: running, itemStyle: { color: tk.success } },
          { name: t('ui.dashboard.notRunning'), value: Math.max(total - running, 0), itemStyle: { color: tk.borderLight } }
        ]
      }
    ]
  }), true)
}

const renderAllocationChart = () => {
  if (!allocationChartRef.value) return
  allocationChart = ensureChart(allocationChart, allocationChartRef.value)
  const tk = getChartTokens()
  const items = allocationItems.value
  allocationChart.setOption(withChartTheme({
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      formatter: (params) => {
        const idx = params[0]?.dataIndex ?? 0
        const it = items[idx]
        return `<b>${it.label}</b><br/>${t('ui.dashboard.allocated')}: ${it.used} (${it.percent}%)<br/>${t('admin.dashboard.total')}: ${it.total}`
      }
    },
    legend: { top: 0, right: 0, data: [t('ui.dashboard.allocated'), t('ui.dashboard.available')] },
    grid: { top: 36, left: 4, right: 24, bottom: 4, containLabel: true },
    xAxis: { type: 'value', max: 100, axisLabel: { formatter: '{value}%' } },
    yAxis: { type: 'category', data: items.map((i) => i.label), inverse: true },
    series: [
      {
        name: t('ui.dashboard.allocated'),
        type: 'bar',
        stack: 'alloc',
        barWidth: 16,
        data: items.map((i, idx) => ({
          value: i.percent,
          itemStyle: { color: [tk.primary, tk.violet, tk.warning][idx], borderRadius: [8, 0, 0, 8] }
        }))
      },
      {
        name: t('ui.dashboard.available'),
        type: 'bar',
        stack: 'alloc',
        barWidth: 16,
        itemStyle: { color: tk.borderLight, borderRadius: [0, 8, 8, 0] },
        data: items.map((i) => 100 - i.percent)
      }
    ]
  }), true)
}

const renderCharts = async () => {
  await nextTick()
  renderDistributionChart()
  renderAllocationChart()
  attachResizeObserver()
}

const handleResize = () => {
  distributionChart?.resize()
  allocationChart?.resize()
}

const attachResizeObserver = () => {
  if (typeof ResizeObserver === 'undefined') return
  if (!resizeObserver) resizeObserver = new ResizeObserver(() => handleResize())
  resizeObserver.disconnect()
  if (distributionChartRef.value) resizeObserver.observe(distributionChartRef.value)
  if (allocationChartRef.value) resizeObserver.observe(allocationChartRef.value)
}

// 主题 / 语言切换时重绘
useChartTheme(() => {
  if (loaded.value) renderCharts()
})
watch(locale, () => {
  if (loaded.value) renderCharts()
})

// ---------- 数据 ----------
const fetchDashboardData = async () => {
  loading.value = true
  try {
    const response = await getAdminDashboard()
    if (response.code === 200) {
      if (response.data && response.data.statistics) {
        Object.assign(dashboardData, response.data.statistics)
      } else {
        Object.assign(dashboardData, response.data)
      }
      if (response.data && response.data.resourceUsage) {
        Object.assign(resourceUsage, response.data.resourceUsage)
      }
      if (response.data && response.data.systemStatus) {
        Object.assign(systemStatus, response.data.systemStatus)
      }
      lastUpdated.value = new Date().toLocaleTimeString(locale.value === 'en-US' ? 'en-US' : 'zh-CN', {
        hour: '2-digit',
        minute: '2-digit'
      })
    }
    return true
  } catch (error) {
    ElMessage.error(t('admin.dashboard.loadDataFailed'))
    console.error('Dashboard data fetch error:', error)
    return false
  } finally {
    loading.value = false
    loaded.value = true
    renderCharts()
  }
}

const handleManualRefresh = async () => {
  const ok = await fetchDashboardData()
  if (ok) ElMessage.success(t('ui.dashboard.refreshed'))
}

onMounted(async () => {
  window.addEventListener('resize', handleResize)
  await fetchDashboardData()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  resizeObserver?.disconnect()
  distributionChart?.dispose()
  allocationChart?.dispose()
  distributionChart = null
  allocationChart = null
})
</script>

<style scoped lang="scss">
.section {
  margin-bottom: 16px;
}

.updated-at {
  color: var(--text-color-tertiary);
}

.skeleton-card {
  padding: 18px;
  min-height: 120px;
}

.skeleton-panel {
  padding: 20px;
  min-height: 300px;
}

/* 指标卡 */
.stat-card .oh-stat-value {
  font-size: 26px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.stat-icon-gradient {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 16px;
  flex-shrink: 0;
  box-shadow: 0 8px 18px -8px var(--oh-tone, var(--primary-color));
}

.grad-indigo { background: linear-gradient(135deg, #818cf8 0%, #4f46e5 100%); }
.grad-emerald { background: linear-gradient(135deg, #34d399 0%, #059669 100%); }
.grad-sky { background: linear-gradient(135deg, #38bdf8 0%, #0284c7 100%); }
.grad-violet { background: linear-gradient(135deg, #a78bfa 0%, #7c3aed 100%); }
.grad-amber { background: linear-gradient(135deg, #fbbf24 0%, #d97706 100%); }

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  display: inline-block;
  background: var(--text-color-tertiary);
}

.dot-ok {
  background: var(--success-color);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--success-color) 22%, transparent);
}

.dot-warn {
  background: var(--warning-color);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--warning-color) 22%, transparent);
}

.dot-bad {
  background: var(--error-color);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--error-color) 22%, transparent);
}

/* 面板 */
.panel {
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
  flex-wrap: wrap;
}

.panel-meta {
  font-size: 12.5px;
  color: var(--text-color-tertiary);
}

.title-icon {
  color: var(--primary-color);
  font-size: 14px;
}

/* 服务器状态环 */
.rings {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  padding: 8px 0 4px;
  flex: 1;
  align-items: center;
}

.ring-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.ring-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.ring-value {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--text-color-primary);
  font-variant-numeric: tabular-nums;

  small {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-color-secondary);
    margin-left: 1px;
  }
}

.ring-label {
  font-size: 13px;
  color: var(--text-color-secondary);
  display: inline-flex;
  align-items: center;
  gap: 6px;

  i {
    color: var(--text-color-tertiary);
  }
}

.uptime-bar {
  margin-top: 16px;
  padding: 10px 14px;
  border-radius: var(--radius-md);
  background: var(--brand-gradient-soft);
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-color-secondary);

  i {
    color: var(--primary-color);
  }

  strong {
    margin-left: auto;
    color: var(--text-color-primary);
    font-variant-numeric: tabular-nums;
  }
}

/* 图表 */
.chart-box {
  width: 100%;
  height: 260px;
}

.chart-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.allocation {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(0, 1fr);
  gap: 24px;
  align-items: center;
}

.allocation-chart {
  height: 220px;
}

.allocation-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.allocation-item {
  padding: 12px 14px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color-light);
  background: var(--surface-muted);
}

.allocation-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.allocation-icon {
  width: 26px;
  height: 26px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: var(--oh-tone);
  background: color-mix(in srgb, var(--oh-tone) 14%, transparent);
}

.allocation-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-color-primary);
}

.allocation-percent {
  margin-left: auto;
  font-size: 13px;
  font-weight: 700;
  color: var(--text-color-primary);
  font-variant-numeric: tabular-nums;
}

.allocation-numbers {
  margin-top: 6px;
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-color-tertiary);

  b {
    color: var(--text-color-secondary);
    font-weight: 600;
    font-variant-numeric: tabular-nums;
  }
}

/* 快捷入口 */
.quick-title {
  margin-bottom: 12px;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 12px;
}

.quick-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px;
  text-decoration: none;
  color: var(--text-color-primary);
  transition: transform var(--duration-normal, 0.2s) var(--ease-out, ease), border-color 0.2s, box-shadow 0.2s;

  &:hover {
    transform: translateY(-2px);
    border-color: var(--border-color-hover);
    box-shadow: var(--shadow-md);

    .quick-arrow {
      opacity: 1;
      transform: translateX(0);
    }
  }
}

.quick-icon {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--oh-tone);
  background: color-mix(in srgb, var(--oh-tone) 12%, transparent);
}

.quick-label {
  font-size: 13px;
  font-weight: 500;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.quick-arrow {
  margin-left: auto;
  font-size: 11px;
  color: var(--text-color-tertiary);
  opacity: 0;
  transform: translateX(-4px);
  transition: all 0.2s;
}

@media (max-width: 1400px) {
  .quick-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 1024px) {
  .allocation {
    grid-template-columns: minmax(0, 1fr);
  }
}

@media (max-width: 768px) {
  .rings {
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 4px;
  }

  .quick-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 480px) {
  .rings {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
