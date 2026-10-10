<template>
  <div class="oh-page user-dashboard">
    <!-- 加载骨架屏 -->
    <template v-if="loading && !loaded">
      <div class="skeleton-header">
        <el-skeleton
          animated
          :rows="1"
        />
      </div>
      <div class="oh-grid oh-grid-4 section">
        <div
          v-for="n in 4"
          :key="`qa-${n}`"
          class="oh-surface skeleton-card"
        >
          <el-skeleton
            animated
            :rows="1"
          />
        </div>
      </div>
      <div class="oh-surface skeleton-panel section">
        <el-skeleton
          animated
          :rows="6"
        />
      </div>
    </template>

    <!-- 主要内容 -->
    <template v-else>
      <div class="oh-page-header">
        <div>
          <div class="oh-page-title">
            <span>{{ greeting }}</span>
            <span class="oh-gradient-text">{{ userInfo?.nickname || userInfo?.username || t('common.user') }}</span>
            <el-tag
              :type="getLevelTagType(userLimits.level)"
              effect="light"
              round
              class="level-tag"
            >
              <i class="fas fa-crown" />
              {{ t('user.dashboard.levelTag', { level: userLimits.level, text: getLevelText(userLimits.level) }) }}
            </el-tag>
          </div>
          <div class="oh-page-subtitle">
            {{ t('user.dashboard.subtitle') }}
          </div>
        </div>
      </div>

      <!-- 快捷操作 -->
      <div class="oh-grid oh-grid-4 section">
        <router-link
          v-for="action in quickActions"
          :key="action.to"
          :to="action.to"
          class="oh-surface quick-card"
          :class="`oh-tone-${action.tone}`"
        >
          <span
            class="quick-icon"
            :class="`grad-${action.tone}`"
          ><i :class="action.icon" /></span>
          <span class="quick-text">
            <span class="quick-label">{{ action.label }}</span>
            <span class="quick-desc">{{ action.desc }}</span>
          </span>
          <i class="fas fa-chevron-right quick-arrow" />
        </router-link>
      </div>

      <!-- 资源配额 -->
      <div class="oh-surface panel section">
        <div class="panel-header">
          <span class="oh-section-title">
            <i class="fas fa-tachometer-alt title-icon" />
            {{ t('user.dashboard.resourceQuota') }}
          </span>
          <el-button
            size="small"
            :loading="refreshing"
            @click="handleManualRefresh"
          >
            <el-icon><Refresh /></el-icon>
            {{ t('common.refresh') }}
          </el-button>
        </div>

        <div class="quota-grid">
          <div
            v-for="q in quotaCards"
            :key="q.key"
            class="quota-card"
            :class="`oh-tone-${q.tone}`"
          >
            <div class="quota-head">
              <span class="quota-icon"><i :class="q.icon" /></span>
              <span class="quota-title">{{ q.title }}</span>
              <span
                v-if="q.limited"
                class="quota-percent"
                :class="percentClass(q.percent)"
              >{{ q.percent }}%</span>
            </div>
            <div class="quota-usage">
              <template v-if="q.limited">
                <span class="quota-used">{{ q.used }}</span>
                <span class="quota-sep">/</span>
                <span class="quota-max">{{ q.max }}</span>
              </template>
              <template v-else>
                <span class="quota-used">{{ t('user.dashboard.unlimited') }}</span>
                <el-tag
                  type="success"
                  size="small"
                  effect="light"
                  round
                >
                  {{ t('user.dashboard.unlimitedTraffic') }}
                </el-tag>
              </template>
            </div>
            <el-progress
              v-if="q.limited"
              :percentage="q.percent"
              :color="getProgressColor(q.percent)"
              :stroke-width="6"
              :show-text="false"
            />
            <div
              v-else
              class="quota-unlimited-bar"
            />
            <div class="quota-desc">
              {{ q.desc }}
              <span
                v-if="q.extra"
                class="quota-extra"
              >{{ q.extra }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 流量使用统计 -->
      <TrafficOverview />

      <!-- 流量历史趋势图 -->
      <div class="traffic-history-section">
        <TrafficHistoryChart
          type="user"
          :title="t('user.dashboard.trafficHistoryChart')"
          :auto-refresh="0"
        />
      </div>

      <!-- 系统公告 -->
      <div
        v-if="announcements.length > 0"
        class="oh-surface panel section"
      >
        <div class="panel-header">
          <span class="oh-section-title">
            <i class="fas fa-bullhorn title-icon" />
            {{ t('user.dashboard.systemAnnouncements') }}
          </span>
        </div>
        <div class="announcements-list">
          <div
            v-for="announcement in announcements"
            :key="announcement.id"
            class="announcement-item"
          >
            <div class="announcement-top">
              <span class="announcement-title">{{ announcement.title }}</span>
              <span class="announcement-date">{{ formatDate(announcement.createdAt) }}</span>
            </div>
            <div class="announcement-content">
              {{ announcement.content }}
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onActivated, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { getUserLimits } from '@/api/user'
import { getAnnouncements } from '@/api/public'
import { useUserStore } from '@/pinia/modules/user'
import { formatMemorySize, formatDiskSize } from '@/utils/unit-formatter'
import TrafficOverview from '@/components/TrafficOverview.vue'
import TrafficHistoryChart from '@/components/TrafficHistoryChart.vue'

const { t, locale } = useI18n()
const userStore = useUserStore()
const userInfo = userStore.user || {}
const loading = ref(true)
const loaded = ref(false)
const refreshing = ref(false)

const userLimits = reactive({
  level: 1,
  maxInstances: 0,
  usedInstances: 0,
  containerCount: 0,
  vmCount: 0,
  maxCpu: 0,
  usedCpu: 0,
  maxMemory: 0,
  usedMemory: 0,
  maxDisk: 0,
  usedDisk: 0,
  maxBandwidth: 0,
  usedBandwidth: 0,
  maxTraffic: 0,
  usedTraffic: 0,
  maxSnapshots: 0,
  usedSnapshots: 0,
  remainingSnapshots: 0,
  maxSnapshotsPerInstance: 0
})

const announcements = ref([])

/**
 * 获取用户限制信息
 * @param {boolean} notify 是否显示"正在刷新 / 已刷新"提示（仅手动刷新时为 true，页面加载时静默）
 */
const loadUserLimits = async (notify = false) => {
  const loadingMsg = notify
    ? ElMessage({
      message: t('user.dashboard.refreshingQuota'),
      type: 'info',
      duration: 0, // 不自动关闭
      showClose: false
    })
    : null

  try {
    const response = await getUserLimits()
    Object.assign(userLimits, response.data)
    loadingMsg?.close()
    if (notify) ElMessage.success(t('user.dashboard.quotaRefreshed'))
  } catch (error) {
    console.error(t('user.dashboard.getUserLimitsFailed'), error)
    loadingMsg?.close()
    ElMessage.error(error?.details || error?.message || t('user.dashboard.loadQuotaFailed'))
  }
}

// 手动刷新（配额卡片右上角按钮）
const handleManualRefresh = async () => {
  refreshing.value = true
  try {
    await loadUserLimits(true)
  } finally {
    refreshing.value = false
  }
}

// 获取公告信息
const loadAnnouncements = async () => {
  try {
    const response = await getAnnouncements({ page: 1, pageSize: 3 })
    const data = response.data
    announcements.value = Array.isArray(data) ? data : (data?.list || data?.items || [])
  } catch (error) {
    console.error(t('user.dashboard.getAnnouncementsFailed'), error)
  }
}

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return t('ui.dashboard.goodMorning')
  if (hour < 18) return t('ui.dashboard.goodAfternoon')
  return t('ui.dashboard.goodEvening')
})

// 获取等级标签类型
const getLevelTagType = (level) => {
  const levelMap = {
    1: 'primary',
    2: 'warning',
    3: 'success',
    4: 'danger'
  }
  return levelMap[level] || 'primary'
}

// 获取等级文本
const getLevelText = () => {
  return t('user.dashboard.normalUser')
}

// 获取使用百分比
const getUsagePercentage = (used, total) => {
  if (!total) return 0
  return Math.min(Math.round((used / total) * 100), 100)
}

// 进度条颜色（设计令牌，随亮暗主题切换）
const getProgressColor = (percentage) => {
  if (percentage >= 90) return 'var(--error-color)'
  if (percentage >= 70) return 'var(--warning-color)'
  return 'var(--primary-color)'
}

const percentClass = (p) => (p >= 90 ? 'is-danger' : p >= 70 ? 'is-warning' : '')

const formatMemory = (memory) => formatMemorySize(memory)
const formatStorage = (disk) => formatDiskSize(disk)
// 流量和磁盘都是以MB为单位，使用相同的格式化
const formatTraffic = (traffic) => formatDiskSize(traffic)

// 格式化日期
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString(locale.value === 'en-US' ? 'en-US' : 'zh-CN')
}

// 快捷操作（路由均存在于 router/index.js）
const quickActions = computed(() => [
  { to: '/user/apply', icon: 'fas fa-plus', tone: 'indigo', label: t('sidebar.apply'), desc: t('ui.dashboard.quickApplyDesc') },
  { to: '/user/instances', icon: 'fas fa-server', tone: 'sky', label: t('sidebar.myInstances'), desc: t('ui.dashboard.quickInstancesDesc') },
  { to: '/user/orders', icon: 'fas fa-receipt', tone: 'amber', label: t('sidebar.orders'), desc: t('ui.dashboard.quickOrdersDesc') },
  { to: '/user/tickets', icon: 'fas fa-headset', tone: 'violet', label: t('sidebar.tickets'), desc: t('ui.dashboard.quickTicketsDesc') }
])

// 配额卡片
const quotaCards = computed(() => {
  const L = userLimits
  const hasTypeBreakdown = L.containerCount !== undefined || L.vmCount !== undefined
  return [
    {
      key: 'instances',
      tone: 'indigo',
      icon: 'fas fa-cubes',
      title: t('user.dashboard.instanceCount'),
      used: L.usedInstances,
      max: L.maxInstances,
      limited: true,
      percent: getUsagePercentage(L.usedInstances, L.maxInstances),
      desc: t('user.dashboard.instanceCountDesc'),
      extra: hasTypeBreakdown
        ? `${t('user.dashboard.containerCount')}: ${L.containerCount || 0} · ${t('user.dashboard.vmCount')}: ${L.vmCount || 0}`
        : ''
    },
    {
      key: 'cpu',
      tone: 'sky',
      icon: 'fas fa-microchip',
      title: t('user.dashboard.cpuCores'),
      used: L.usedCpu,
      max: `${L.maxCpu}${t('user.dashboard.cores')}`,
      limited: true,
      percent: getUsagePercentage(L.usedCpu, L.maxCpu),
      desc: t('user.dashboard.cpuCoresDesc')
    },
    {
      key: 'memory',
      tone: 'violet',
      icon: 'fas fa-memory',
      title: t('user.dashboard.memorySize'),
      used: formatMemory(L.usedMemory),
      max: formatMemory(L.maxMemory),
      limited: true,
      percent: getUsagePercentage(L.usedMemory, L.maxMemory),
      desc: t('user.dashboard.memorySizeDesc')
    },
    {
      key: 'disk',
      tone: 'amber',
      icon: 'fas fa-hdd',
      title: t('user.dashboard.storageSpace'),
      used: formatStorage(L.usedDisk),
      max: formatStorage(L.maxDisk),
      limited: true,
      percent: getUsagePercentage(L.usedDisk, L.maxDisk),
      desc: t('user.dashboard.storageSpaceDesc')
    },
    {
      key: 'traffic',
      tone: 'emerald',
      icon: 'fas fa-exchange-alt',
      title: t('user.dashboard.trafficLimit'),
      used: formatTraffic(L.usedTraffic),
      max: formatTraffic(L.maxTraffic),
      limited: L.maxTraffic > 0,
      percent: getUsagePercentage(L.usedTraffic, L.maxTraffic),
      desc: L.maxTraffic > 0 ? t('user.dashboard.trafficLimitDesc') : t('user.dashboard.unlimitedTrafficDesc')
    },
    {
      key: 'snapshots',
      tone: 'rose',
      icon: 'fas fa-camera',
      title: t('user.dashboard.snapshotQuota'),
      used: L.usedSnapshots,
      max: L.maxSnapshots,
      limited: true,
      percent: getUsagePercentage(L.usedSnapshots, L.maxSnapshots),
      desc: t('user.dashboard.snapshotQuotaDesc', {
        remaining: L.remainingSnapshots,
        perInstance: L.maxSnapshotsPerInstance
      })
    }
  ]
})

// 静默加载全部数据；首次加载显示骨架屏，之后后台刷新不闪烁
const loadAll = async () => {
  if (!loaded.value) loading.value = true
  try {
    await Promise.all([
      loadUserLimits(false),
      loadAnnouncements()
    ])
  } finally {
    loading.value = false
    loaded.value = true
  }
}

// keep-alive 下首次挂载时 onActivated 也会紧接着触发，避免重复请求
let skipNextActivation = false

onMounted(async () => {
  // 强制页面刷新监听器
  window.addEventListener('force-page-refresh', handleForceRefresh)
  skipNextActivation = true
  await loadAll()
})

// 使用 onActivated 确保每次页面激活时都重新加载数据
onActivated(async () => {
  if (skipNextActivation) {
    skipNextActivation = false
    return
  }
  await loadAll()
})

// 处理强制刷新事件
const handleForceRefresh = async (event) => {
  if (event.detail && event.detail.path === '/user/dashboard') {
    await loadAll()
  }
}

onUnmounted(() => {
  // 清理事件监听器
  window.removeEventListener('force-page-refresh', handleForceRefresh)
})
</script>

<style scoped lang="scss">
.section {
  margin-bottom: 16px;
}

.skeleton-header {
  max-width: 420px;
  margin-bottom: 20px;
}

.skeleton-card {
  padding: 16px;
  min-height: 76px;
}

.skeleton-panel {
  padding: 20px;
}

.level-tag {
  font-size: 12.5px;
  font-weight: 500;

  i {
    margin-right: 4px;
    font-size: 11px;
  }
}

/* 快捷操作 */
.quick-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  text-decoration: none;
  color: var(--text-color-primary);
  transition: transform var(--duration-normal, 0.2s) var(--ease-out, ease), border-color 0.2s, box-shadow 0.2s;

  &:hover {
    transform: translateY(-2px);
    border-color: var(--border-color-hover);
    box-shadow: var(--shadow-md);

    .quick-arrow {
      color: var(--oh-tone);
      transform: translateX(2px);
    }
  }
}

.quick-icon {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 16px;
  box-shadow: 0 8px 18px -8px var(--oh-tone, var(--primary-color));
}

.grad-indigo { background: linear-gradient(135deg, #818cf8 0%, #4f46e5 100%); }
.grad-sky { background: linear-gradient(135deg, #38bdf8 0%, #0284c7 100%); }
.grad-amber { background: linear-gradient(135deg, #fbbf24 0%, #d97706 100%); }
.grad-violet { background: linear-gradient(135deg, #a78bfa 0%, #7c3aed 100%); }

.quick-text {
  display: flex;
  flex-direction: column;
  min-width: 0;
  gap: 2px;
}

.quick-label {
  font-size: 14px;
  font-weight: 600;
}

.quick-desc {
  font-size: 12px;
  color: var(--text-color-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.quick-arrow {
  margin-left: auto;
  font-size: 11px;
  color: var(--text-color-tertiary);
  transition: all 0.2s;
}

/* 面板 */
.panel {
  padding: 18px 20px;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.title-icon {
  color: var(--primary-color);
  font-size: 14px;
}

/* 配额 */
.quota-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

.quota-card {
  position: relative;
  padding: 14px 16px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color-light);
  background: var(--surface-muted);
  overflow: hidden;
  transition: border-color 0.2s;

  &:hover {
    border-color: color-mix(in srgb, var(--oh-tone) 35%, transparent);
  }
}

.quota-head {
  display: flex;
  align-items: center;
  gap: 8px;
}

.quota-icon {
  width: 30px;
  height: 30px;
  border-radius: 9px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  color: var(--oh-tone);
  background: color-mix(in srgb, var(--oh-tone) 13%, transparent);
}

.quota-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-color-secondary);
}

.quota-percent {
  margin-left: auto;
  font-size: 12px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 999px;
  color: var(--primary-color);
  background: var(--primary-color-bg);
  font-variant-numeric: tabular-nums;

  &.is-warning {
    color: var(--warning-color);
    background: var(--warning-color-bg);
  }

  &.is-danger {
    color: var(--error-color);
    background: var(--danger-color-bg);
  }
}

.quota-usage {
  display: flex;
  align-items: baseline;
  flex-wrap: wrap;
  gap: 6px;
  margin: 12px 0 10px;
  font-variant-numeric: tabular-nums;
}

.quota-used {
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--text-color-primary);
}

.quota-sep,
.quota-max {
  font-size: 13.5px;
  color: var(--text-color-tertiary);
}

.quota-unlimited-bar {
  height: 6px;
  border-radius: 999px;
  background: linear-gradient(90deg, color-mix(in srgb, var(--success-color) 55%, transparent), color-mix(in srgb, var(--success-color) 12%, transparent));
}

.quota-desc {
  margin-top: 10px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--text-color-tertiary);
}

.quota-extra {
  display: block;
  margin-top: 2px;
  color: var(--text-color-secondary);
}

/* 流量 */
.traffic-history-section {
  margin-bottom: 16px;

  :deep(.traffic-history-chart) {
    margin-top: 0;
  }
}

/* 公告 */
.announcements-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.announcement-item {
  padding: 14px 16px;
  border-radius: var(--radius-md);
  background: var(--surface-muted);
  border: 1px solid var(--border-color-light);
  border-left: 3px solid var(--primary-color);
}

.announcement-top {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 6px;
}

.announcement-title {
  font-weight: 600;
  color: var(--text-color-primary);
}

.announcement-content {
  color: var(--text-color-secondary);
  line-height: 1.6;
  font-size: 13.5px;
}

.announcement-date {
  flex-shrink: 0;
  font-size: 12px;
  color: var(--text-color-tertiary);
}

@media (max-width: 1200px) {
  .quota-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .quota-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .panel {
    padding: 16px;
  }
}
</style>
