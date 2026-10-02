<template>
  <div class="user-dashboard">
    <!-- 加载状态 -->
    <div
      v-if="loading"
      class="loading-container"
    >
      <el-loading-directive />
      <div class="loading-text">
        {{ t('common.loading') }}
      </div>
    </div>
    
    <!-- 主要内容 -->
    <div v-else>
      <div class="dashboard-header">
        <h1>
          {{ t('user.dashboard.welcome', { name: userInfo?.nickname || userInfo?.username || t('common.user') }) }}
          <el-tag
            :type="getLevelTagType(userLimits.level)"
            size="large"
            effect="plain"
            class="level-tag"
          >
            {{ t('user.dashboard.levelTag', { level: userLimits.level, text: getLevelText(userLimits.level) }) }}
          </el-tag>
        </h1>
        <p>{{ t('user.dashboard.subtitle') }}</p>
      </div>

      <!-- 资源限制信息 -->
      <div class="resource-limits-section">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ t('user.dashboard.resourceQuota') }}</span>
              <el-button
                @click="loadUserLimits"
              >
                {{ t('common.refresh') }}
              </el-button>
            </div>
          </template>
        
          <div class="limits-grid">
            <!-- 实例数量限制 -->
            <div class="limit-item">
              <div class="limit-header">
                <span class="limit-title">{{ t('user.dashboard.instanceCount') }}</span>
                <span class="limit-usage">{{ userLimits.usedInstances }} / {{ userLimits.maxInstances }}</span>
              </div>
              <el-progress 
                :percentage="getUsagePercentage(userLimits.usedInstances, userLimits.maxInstances)"
                :color="getProgressColor(userLimits.usedInstances, userLimits.maxInstances)"
                :stroke-width="8"
              />
              <div class="limit-description">
                {{ t('user.dashboard.instanceCountDesc') }}
                <span
                  v-if="userLimits.containerCount !== undefined || userLimits.vmCount !== undefined"
                  style="display: block; margin-top: 4px; color: #909399; font-size: 12px;"
                >
                  {{ t('user.dashboard.containerCount') }}: {{ userLimits.containerCount || 0 }} / {{ t('user.dashboard.vmCount') }}: {{ userLimits.vmCount || 0 }}
                </span>
              </div>
            </div>

            <!-- CPU核心限制 -->
            <div class="limit-item">
              <div class="limit-header">
                <span class="limit-title">{{ t('user.dashboard.cpuCores') }}</span>
                <span class="limit-usage">{{ userLimits.usedCpu }} / {{ userLimits.maxCpu }}{{ t('user.dashboard.cores') }}</span>
              </div>
              <el-progress 
                :percentage="getUsagePercentage(userLimits.usedCpu, userLimits.maxCpu)"
                :color="getProgressColor(userLimits.usedCpu, userLimits.maxCpu)"
                :stroke-width="8"
              />
              <div class="limit-description">
                {{ t('user.dashboard.cpuCoresDesc') }}
              </div>
            </div>

            <!-- 内存限制 -->
            <div class="limit-item">
              <div class="limit-header">
                <span class="limit-title">{{ t('user.dashboard.memorySize') }}</span>
                <span class="limit-usage">{{ formatMemory(userLimits.usedMemory) }} / {{ formatMemory(userLimits.maxMemory) }}</span>
              </div>
              <el-progress 
                :percentage="getUsagePercentage(userLimits.usedMemory, userLimits.maxMemory)"
                :color="getProgressColor(userLimits.usedMemory, userLimits.maxMemory)"
                :stroke-width="8"
              />
              <div class="limit-description">
                {{ t('user.dashboard.memorySizeDesc') }}
              </div>
            </div>

            <!-- 存储空间限制 -->
            <div class="limit-item">
              <div class="limit-header">
                <span class="limit-title">{{ t('user.dashboard.storageSpace') }}</span>
                <span class="limit-usage">{{ formatStorage(userLimits.usedDisk) }} / {{ formatStorage(userLimits.maxDisk) }}</span>
              </div>
              <el-progress 
                :percentage="getUsagePercentage(userLimits.usedDisk, userLimits.maxDisk)"
                :color="getProgressColor(userLimits.usedDisk, userLimits.maxDisk)"
                :stroke-width="8"
              />
              <div class="limit-description">
                {{ t('user.dashboard.storageSpaceDesc') }}
              </div>
            </div>

            <!-- 流量限制 -->
            <div class="limit-item">
              <div class="limit-header">
                <span class="limit-title">{{ t('user.dashboard.trafficLimit') }}</span>
                <span class="limit-usage">
                  {{ userLimits.maxTraffic > 0 ? `${formatTraffic(userLimits.usedTraffic)} / ${formatTraffic(userLimits.maxTraffic)}` : t('user.dashboard.unlimited') }}
                </span>
              </div>
              <el-progress 
                v-if="userLimits.maxTraffic > 0"
                :percentage="getUsagePercentage(userLimits.usedTraffic, userLimits.maxTraffic)"
                :color="getProgressColor(userLimits.usedTraffic, userLimits.maxTraffic)"
                :stroke-width="8"
              />
              <div
                v-else
                class="unlimited-badge"
              >
                <el-tag
                  type="success"
                  size="small"
                >
                  {{ t('user.dashboard.unlimitedTraffic') }}
                </el-tag>
              </div>
              <div class="limit-description">
                {{ userLimits.maxTraffic > 0 ? t('user.dashboard.trafficLimitDesc') : t('user.dashboard.unlimitedTrafficDesc') }}
              </div>
            </div>

            <!-- 快照配额 -->
            <div class="limit-item">
              <div class="limit-header">
                <span class="limit-title">{{ t('user.dashboard.snapshotQuota') }}</span>
                <span class="limit-usage">{{ userLimits.usedSnapshots }} / {{ userLimits.maxSnapshots }}</span>
              </div>
              <el-progress
                :percentage="getUsagePercentage(userLimits.usedSnapshots, userLimits.maxSnapshots)"
                :color="getProgressColor(userLimits.usedSnapshots, userLimits.maxSnapshots)"
                :stroke-width="8"
              />
              <div class="limit-description">
                {{ t('user.dashboard.snapshotQuotaDesc', {
                  remaining: userLimits.remainingSnapshots,
                  perInstance: userLimits.maxSnapshotsPerInstance
                }) }}
              </div>
            </div>
          </div>
        </el-card>
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
        class="announcements"
      >
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ t('user.dashboard.systemAnnouncements') }}</span>
            </div>
          </template>
        
          <div class="announcements-list">
            <div 
              v-for="announcement in announcements" 
              :key="announcement.id"
              class="announcement-item"
            >
              <div class="announcement-title">
                {{ announcement.title }}
              </div>
              <div class="announcement-content">
                {{ announcement.content }}
              </div>
              <div class="announcement-date">
                {{ formatDate(announcement.createdAt) }}
              </div>
            </div>
          </div>
        </el-card>
      </div>
    </div> <!-- 结束主要内容区域 -->
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onActivated, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getUserLimits } from '@/api/user'
import { getAnnouncements } from '@/api/public'
import { useUserStore } from '@/pinia/modules/user'
import { formatMemorySize, formatDiskSize, formatBandwidthSpeed } from '@/utils/unit-formatter'
import TrafficOverview from '@/components/TrafficOverview.vue'
import TrafficHistoryChart from '@/components/TrafficHistoryChart.vue'

const { t, locale } = useI18n()
const userStore = useUserStore()
const userInfo = userStore.user || {}
const loading = ref(true)

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

// 获取用户限制信息
const loadUserLimits = async () => {
  const loadingMsg = ElMessage({
    message: t('user.dashboard.refreshingQuota'),
    type: 'info',
    duration: 0, // 不自动关闭
    showClose: false
  })
  
  try {
    const response = await getUserLimits()
    Object.assign(userLimits, response.data)
    loadingMsg.close()
    ElMessage.success(t('user.dashboard.quotaRefreshed'))
  } catch (error) {
    console.error(t('user.dashboard.getUserLimitsFailed'), error)
    loadingMsg.close()
    ElMessage.error(error?.details || error?.message || t('user.dashboard.loadQuotaFailed'))
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

// 获取等级标签类型
const getLevelTagType = (level) => {
  const levelMap = {
    1: '',
    2: 'warning',
    3: 'success', 
    4: 'danger'
  }
  return levelMap[level] || ''
}

// 获取等级文本
const getLevelText = () => {
  return t('user.dashboard.normalUser')
}

// 获取等级描述
const getLevelDescription = () => {
  return t('user.dashboard.enjoyIt')
}

// 获取使用百分比
const getUsagePercentage = (used, total) => {
  if (!total) return 0
  return Math.round((used / total) * 100)
}

// 获取进度条颜色
const getProgressColor = (used, total) => {
  const percentage = getUsagePercentage(used, total)
  if (percentage >= 90) return '#f56c6c'
  if (percentage >= 70) return '#e6a23c'
  return '#67c23a'
}

// 格式化内存显示
const formatMemory = (memory) => {
  return formatMemorySize(memory)
}

// 格式化存储显示
const formatStorage = (disk) => {
  return formatDiskSize(disk)
}

// 格式化带宽显示
const formatBandwidth = (bandwidth) => {
  return formatBandwidthSpeed(bandwidth)
}

// 格式化流量显示
const formatTraffic = (traffic) => {
  return formatDiskSize(traffic) // 流量和磁盘都是以MB为单位，使用相同的格式化
}

// 格式化日期
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString(locale.value === 'en-US' ? 'en-US' : 'zh-CN')
}

onMounted(async () => {
  // 强制页面刷新监听器
  window.addEventListener('force-page-refresh', handleForceRefresh)
  
  loading.value = true
  try {
    await Promise.all([
      loadUserLimits(),
      loadAnnouncements()
    ])
  } finally {
    loading.value = false
  }
})

// 使用 onActivated 确保每次页面激活时都重新加载数据
onActivated(async () => {
  loading.value = true
  try {
    await Promise.all([
      loadUserLimits(),
      loadAnnouncements()
    ])
  } finally {
    loading.value = false
  }
})

// 处理强制刷新事件
const handleForceRefresh = async (event) => {
  if (event.detail && event.detail.path === '/user/dashboard') {
    loading.value = true
    try {
      await Promise.all([
        loadUserLimits(),
        loadAnnouncements()
      ])
    } finally {
      loading.value = false
    }
  }
}

onUnmounted(() => {
  // 清理事件监听器
  window.removeEventListener('force-page-refresh', handleForceRefresh)
})
</script>
<style scoped>
.user-dashboard {
  padding: 24px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  color: #666;
}

.loading-text {
  margin-top: 16px;
  font-size: 14px;
}

.dashboard-header {
  margin-bottom: 24px;
  animation: fadeInUp 0.5s cubic-bezier(0.4, 0, 0.2, 1) both;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.dashboard-header h1 {
  margin: 0 0 8px 0;
  color: var(--text-color-primary);
  font-size: 28px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  letter-spacing: -0.5px;
}

.dashboard-header p {
  margin: 0;
  color: var(--text-color-secondary);
  font-size: 15px;
}

.level-tag {
  font-size: 14px;
  font-weight: 600;
  border-radius: 999px !important;
  padding: 2px 14px !important;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 700;
  color: var(--text-color-primary);
  font-size: 15px;
}

/* 资源限制信息 */
.resource-limits-section {
  margin-bottom: 24px;
  animation: fadeInUp 0.55s cubic-bezier(0.4, 0, 0.2, 1) 0.05s both;
}

.resource-limits-section :deep(.el-card) {
  border-radius: 20px !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.06) !important;
  border: 1px solid var(--border-color) !important;
  transition: box-shadow 0.35s cubic-bezier(0.4, 0, 0.2, 1);
}

.resource-limits-section :deep(.el-card:hover) {
  box-shadow: 0 12px 40px rgba(99, 102, 241, 0.1) !important;
}

.limits-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.limit-item {
  padding: 20px;
  background: var(--neutral-bg);
  border-radius: 16px;
  border: 1px solid var(--border-color);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1),
              box-shadow 0.3s cubic-bezier(0.4, 0, 0.2, 1),
              border-color 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.limit-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--primary-color), var(--primary-color-light));
  opacity: 0;
  transition: opacity 0.3s ease;
}

.limit-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.08);
  border-color: var(--border-color-hover);
}

.limit-item:hover::before {
  opacity: 1;
}

.limit-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.limit-title {
  font-weight: 700;
  color: var(--text-color-primary);
  font-size: 14px;
}

.limit-usage {
  font-weight: 700;
  color: var(--text-color-secondary);
  font-size: 14px;
  font-variant-numeric: tabular-nums;
}

.limit-item :deep(.el-progress-bar__outer) {
  border-radius: 999px !important;
  background-color: rgba(148, 163, 184, 0.15) !important;
}

.limit-item :deep(.el-progress-bar__inner) {
  border-radius: 999px !important;
  transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

.limit-description {
  margin-top: 10px;
  font-size: 12px;
  color: #9ca3af;
  line-height: 1.6;
}

/* 公告 */
.announcements {
  margin-bottom: 24px;
  animation: fadeInUp 0.6s cubic-bezier(0.4, 0, 0.2, 1) 0.1s both;
}

.announcements :deep(.el-card) {
  border-radius: 20px !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.06) !important;
}

/* 流量历史图表 */
.traffic-history-section {
  margin-bottom: 24px;
  animation: fadeInUp 0.6s cubic-bezier(0.4, 0, 0.2, 1) 0.08s both;
}

.traffic-history-section :deep(.el-card) {
  border-radius: 20px !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.06) !important;
}

.announcements-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.announcement-item {
  padding: 18px 20px;
  background: var(--neutral-bg);
  border-radius: 14px;
  border-left: 4px solid #10b981;
  transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1),
              box-shadow 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.announcement-item:hover {
  transform: translateX(4px);
  box-shadow: 0 4px 20px rgba(16, 185, 129, 0.08);
}

.announcement-title {
  font-weight: 700;
  color: var(--text-color-primary);
  margin-bottom: 8px;
  font-size: 15px;
}

.announcement-content {
  color: #4b5563;
  margin-bottom: 8px;
  line-height: 1.6;
  font-size: 14px;
}

.announcement-date {
  font-size: 12px;
  color: #9ca3af;
}

.unlimited-badge {
  margin: 8px 0;
  text-align: center;
}

.unlimited-badge :deep(.el-tag) {
  border-radius: 999px !important;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .limits-grid {
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
    gap: 16px;
  }
  
  .limit-item {
    padding: 16px;
  }
}

@media (max-width: 768px) {
  .user-dashboard {
    padding: 16px;
  }
  
  .dashboard-header h1 {
    font-size: 22px;
  }
  
  .dashboard-header p {
    font-size: 13px;
  }
  
  .limits-grid {
    grid-template-columns: 1fr;
    gap: 14px;
  }
  
  .limit-item {
    padding: 14px;
    border-radius: 12px;
  }
  
  .limit-header {
    margin-bottom: 10px;
  }
  
  .limit-title {
    font-size: 13px;
  }
  
  .limit-usage {
    font-size: 13px;
  }
  
  .announcement-item {
    padding: 14px;
    border-radius: 12px;
  }
}

@media (max-width: 480px) {
  .user-dashboard {
    padding: 12px;
  }
  
  .dashboard-header h1 {
    font-size: 20px;
    gap: 8px;
  }
  
  .level-tag {
    font-size: 12px;
    padding: 2px 10px !important;
  }
  
  .limit-item {
    padding: 12px;
  }
  
  .limit-title {
    font-size: 12px;
  }
  
  .limit-usage {
    font-size: 12px;
  }
  
  .limit-description {
    font-size: 11px;
  }
  
  .announcement-item {
    padding: 12px;
  }
  
  .announcement-title {
    font-size: 14px;
  }
  
  .announcement-content {
    font-size: 13px;
  }
}
</style>
