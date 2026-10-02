<template>
  <div class="instance-detail">
    <!-- 页面头部 -->
    <div class="page-header">
      <el-button 
        type="text" 
        class="back-btn"
        @click="handleBack"
      >
        <el-icon><ArrowLeft /></el-icon>
        {{ t('user.instanceDetail.backToList') }}
      </el-button>
    </div>

    <!-- 实例概览卡片 -->
    <InstanceOverviewCard
      :instance="instance"
      :monitoring="monitoring"
      :action-loading="actionLoading"
      :instance-type-permissions="instanceTypePermissions"
      :share-mode="isShareMode"
      @perform-action="performAction"
      @reset-password="showResetPasswordDialog"
      @open-ssh="openSSHTerminal"
      @open-vnc="openConsole('')"
      @view-task="viewTaskDetail"
      @create-share="createShareLink"
    />

    <!-- 标签页内容 -->
    <el-card class="tabs-card">
      <el-tabs
        v-model="activeTab"
        type="border-card"
      >
        <!-- 概览标签页 -->
        <el-tab-pane
          :label="t('user.instanceDetail.overview')"
          name="overview"
        >
          <OverviewTab
            :instance="instance"
            :show-password="showPassword"
            @toggle-password="togglePassword"
            @copy="copyToClipboard"
          />
        </el-tab-pane>

        <!-- 端口映射标签页 -->
        <el-tab-pane
          :label="t('user.instanceDetail.portMapping')"
          name="ports"
        >
          <PortMappingsTab
            :instance="instance"
            :port-mappings="portMappings"
            @refresh="refreshPortMappings"
            @copy="copyToClipboard"
          />
        </el-tab-pane>

        <!-- 统计标签页 -->
        <el-tab-pane
          :label="t('user.instanceDetail.statistics')"
          name="stats"
        >
          <StatisticsTab
            :instance="instance"
            :monitoring="monitoring"
            :instance-id="currentInstanceId"
            :share-token="shareToken"
            @refresh="refreshMonitoring"
            @show-traffic-detail="showTrafficDetail = true"
          />
        </el-tab-pane>
        <!-- 资源监控标签页 -->
        <el-tab-pane
          :label="t('user.instanceDetail.resourceMonitoring')"
          name="resources"
        >
          <ResourceMonitorChart
            ref="resourceChartRef"
            :instance-id="currentInstanceId"
            :share-token="shareToken"
            :auto-refresh="300000"
          />
        </el-tab-pane>

        <!-- 快照标签页 -->
        <el-tab-pane
          :label="t('user.instanceDetail.snapshots')"
          name="snapshots"
        >
          <SnapshotsTab
            :instance-id="currentInstanceId"
            :share-token="shareToken"
            :readonly="isShareMode"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- PMAcct 流量详情对话框 -->
    <InstanceTrafficDetail
      v-model="showTrafficDetail"
      :instance-id="currentInstanceId"
      :share-token="shareToken"
      :instance-name="instance.name"
    />

    <VNCDialog
      v-model="showVNCDialog"
      :instance-id="currentInstanceId"
      :instance-name="instance.name"
      :initial-protocol="selectedConsoleProtocol"
      :scope="isShareMode ? 'share' : 'user'"
      :share-token="shareToken"
    />

    <!-- 重置系统镜像选择对话框 -->
    <el-dialog
      v-model="showResetImageDialog"
      :title="t('user.instanceDetail.selectResetImage')"
      width="500px"
      destroy-on-close
    >
      <div v-loading="loadingResetImages">
        <p style="margin-bottom: 12px; color: var(--el-text-color-secondary);">
          {{ t('user.instanceDetail.selectResetImageTip') }}
        </p>
        <el-radio-group
          v-model="selectedResetImage"
          style="display: flex; flex-direction: column; gap: 8px;"
        >
          <el-radio
            v-for="img in resetImages"
            :key="img.name || img.id"
            :value="img.name"
            border
            style="margin: 0; width: 100%;"
          >
            <span style="display: inline-flex; align-items: center; gap: 6px;">
              <OsIcon
                :name="img.name"
                :size="20"
              />
              {{ img.display_name || img.name }}
            </span>
          </el-radio>
        </el-radio-group>
      </div>
      <template #footer>
        <el-button @click="showResetImageDialog = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :disabled="!selectedResetImage"
          @click="confirmResetWithImage"
        >
          {{ t('user.instanceDetail.confirmReset') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import InstanceTrafficDetail from '@/components/InstanceTrafficDetail.vue'
import ResourceMonitorChart from '@/components/ResourceMonitorChart.vue'
import VNCDialog from '@/components/VNCDialog.vue'
import OsIcon from '@/components/OsIcon.vue'
import { useInstanceDetail } from './composables/useInstanceDetail'
import { useInstanceActions } from './composables/useInstanceActions'
import InstanceOverviewCard from './components/InstanceOverviewCard.vue'
import OverviewTab from './components/OverviewTab.vue'
import PortMappingsTab from './components/PortMappingsTab.vue'
import StatisticsTab from './components/StatisticsTab.vue'
import SnapshotsTab from './components/SnapshotsTab.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const activeTab = ref('overview')
const resourceChartRef = ref(null)
const showVNCDialog = ref(false)
const selectedConsoleProtocol = ref('')
const shareToken = computed(() => route.params.token ? String(route.params.token) : '')
const isShareMode = computed(() => Boolean(shareToken.value))

const {
  portMappings,
  instanceTypePermissions,
  instance,
  monitoring,
  updateInstancePermissions,
  loadInstanceDetail,
  refreshPortMappings,
  refreshMonitoring,
  loadInstanceTypePermissions,
  startShareTokenMonitor
} = useInstanceDetail(shareToken)
const currentInstanceId = computed(() => instance.value?.id || route.params.id || '')

const {
  actionLoading,
  showPassword,
  showTrafficDetail,
  showResetImageDialog,
  resetImages,
  selectedResetImage,
  loadingResetImages,
  confirmResetWithImage,
  viewTaskDetail,
  performAction,
  openSSHTerminal,
  createShareLink,
  showResetPasswordDialog,
  togglePassword,
  copyToClipboard
} = useInstanceActions(instance, monitoring, loadInstanceDetail, shareToken)

const openConsole = (protocol = '') => {
  selectedConsoleProtocol.value = protocol || ''
  showVNCDialog.value = true
}

// 标志位，防止 watch 循环触发
let isUpdatingFromRoute = false

const handleBack = () => {
  if (isShareMode.value) {
    router.push('/home')
    return
  }
  router.back()
}

watch(() => [route.params.id, route.params.token], async ([newId, newToken], [oldId, oldToken] = []) => {
  if ((newToken && newToken !== oldToken) || (newId && newId !== oldId && newId !== 'undefined')) {
    try {
      const [detailSuccess, permissionsSuccess] = await Promise.all([
        loadInstanceDetail(true),
        loadInstanceTypePermissions()
      ])
      if (detailSuccess && permissionsSuccess) {
        updateInstancePermissions()
        refreshMonitoring()
        refreshPortMappings()
      }
    } catch (error) {
      console.error('路由切换时加载数据失败:', error)
    }
  }
})

watch(() => route.query.tab, (newTab, oldTab) => {
  if (newTab === oldTab) return
  if (newTab && ['overview', 'ports', 'stats', 'resources', 'snapshots'].includes(newTab)) {
    if (activeTab.value === newTab) return
    isUpdatingFromRoute = true
    activeTab.value = newTab
    nextTick(() => { isUpdatingFromRoute = false })
  } else {
    if (activeTab.value !== 'overview') {
      isUpdatingFromRoute = true
      activeTab.value = 'overview'
      nextTick(() => { isUpdatingFromRoute = false })
    }
  }
}, { immediate: true })

watch(activeTab, (newTab, oldTab) => {
  if (newTab === oldTab || isUpdatingFromRoute) return
  if (newTab && route.query.tab !== newTab) {
    router.replace({ query: { ...route.query, tab: newTab } })
  }
})

let monitoringTimer = null

onMounted(async () => {
  await nextTick()
  try {
    const [detailSuccess, permissionsSuccess] = await Promise.all([
      loadInstanceDetail(true),
      loadInstanceTypePermissions()
    ])
    if (detailSuccess && permissionsSuccess) {
      updateInstancePermissions()
      refreshMonitoring()
      refreshPortMappings()
      startShareTokenMonitor() // 定时检测分享令牌有效期
      monitoringTimer = setInterval(refreshMonitoring, 30000)
    }
  } catch (error) {
    console.error('页面初始化失败:', error)
  }
})

onUnmounted(() => {
  if (monitoringTimer) {
    clearInterval(monitoringTimer)
    monitoringTimer = null
  }
})
</script>

<style src="./detail.css" scoped>

/* ── 现代 UI 增强 ── */

/* 卡片 */
:deep(.el-card),
:deep(.el-card.is-always-shadow),
:deep(.el-card.is-hover-shadow) {
  border-radius: 16px !important;
  border: 1px solid var(--border-color) !important;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05) !important;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1),
              box-shadow 0.3s cubic-bezier(0.4, 0, 0.2, 1),
              border-color 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.el-card:hover) {
  transform: translateY(-3px) !important;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.09) !important;
}

/* 统计卡片 */
.stat-card :deep(.el-card__body) {
  padding: 20px !important;
}

.stat-content {
  display: flex !important;
  align-items: center !important;
  gap: 14px !important;
}

.stat-icon {
  width: 48px !important;
  height: 48px !important;
  border-radius: 14px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  font-size: 22px !important;
  flex-shrink: 0 !important;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.08) !important;
}

.stat-icon.user-icon { background: linear-gradient(135deg, #6366f1, #818cf8); color: #fff; }
.stat-icon.server-icon { background: linear-gradient(135deg, #10b981, #34d399); color: #fff; }
.stat-icon.vm-icon { background: linear-gradient(135deg, #f59e0b, #fbbf24); color: #fff; }
.stat-icon.container-icon { background: linear-gradient(135deg, #ec4899, #f472b6); color: #fff; }
.stat-icon.cpu-icon { background: linear-gradient(135deg, #6366f1, #818cf8); color: #fff; }
.stat-icon.memory-icon { background: linear-gradient(135deg, #10b981, #34d399); color: #fff; }
.stat-icon.disk-icon { background: linear-gradient(135deg, #f59e0b, #fbbf24); color: #fff; }

.stat-number {
  font-size: 28px !important;
  font-weight: 800 !important;
  color: var(--text-color-primary) !important;
  line-height: 1.1 !important;
  font-variant-numeric: tabular-nums !important;
}

.stat-label {
  font-size: 13px !important;
  color: var(--text-color-secondary) !important;
  font-weight: 600 !important;
  margin-top: 2px !important;
}

/* 资源卡片 */
.resource-card :deep(.el-card__body) {
  padding: 16px 20px !important;
}

.resource-header {
  display: flex !important;
  align-items: center !important;
  gap: 10px !important;
  margin-bottom: 12px !important;
}

.resource-icon {
  width: 36px !important;
  height: 36px !important;
  border-radius: 10px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  font-size: 16px !important;
  flex-shrink: 0 !important;
}

.resource-title {
  font-weight: 700 !important;
  font-size: 14px !important;
  color: var(--text-color-primary) !important;
}

.resource-number-value {
  font-size: 22px !important;
  font-weight: 800 !important;
  color: var(--primary-color) !important;
  font-variant-numeric: tabular-nums !important;
}

.resource-number-label {
  font-size: 11px !important;
  color: var(--text-color-secondary) !important;
  font-weight: 600 !important;
}

.resource-number-separator {
  color: var(--text-color-tertiary) !important;
  font-weight: 700 !important;
  font-size: 16px !important;
  margin: 0 8px !important;
}

/* 按钮 */
:deep(.el-button--primary) {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-light)) !important;
  border: none !important;
  border-radius: 12px !important;
  font-weight: 600 !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1) !important;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.2) !important;
}

:deep(.el-button--primary:hover) {
  transform: translateY(-2px) !important;
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.3) !important;
}

:deep(.el-button--success) {
  border-radius: 12px !important;
  font-weight: 600 !important;
}

:deep(.el-button--warning) {
  border-radius: 12px !important;
  font-weight: 600 !important;
}

:deep(.el-button--danger) {
  border-radius: 12px !important;
  font-weight: 600 !important;
}

:deep(.el-button--default) {
  border-radius: 12px !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.el-button--default:hover) {
  transform: translateY(-1px) !important;
}

/* 表格 */
:deep(.el-table) {
  --el-table-border-color: var(--border-color) !important;
}

:deep(.el-table th) {
  background: rgba(99, 102, 241, 0.04) !important;
  color: var(--text-color-primary) !important;
  font-weight: 700 !important;
  font-size: 13px !important;
  letter-spacing: 0.3px !important;
}

:deep(.el-table td) {
  font-size: 13px !important;
  color: var(--text-color-primary) !important;
}

:deep(.el-table__row) {
  transition: all 0.2s ease !important;
}

:deep(.el-table__row:hover) {
  background: rgba(99, 102, 241, 0.04) !important;
  transform: scale(1.005) !important;
}

/* 徽章 */
:deep(.el-tag) {
  border-radius: 999px !important;
  font-weight: 600 !important;
  font-size: 12px !important;
  padding: 2px 12px !important;
  border: none !important;
}

:deep(.el-tag--success) {
  background: rgba(16, 185, 129, 0.12) !important;
  color: #10b981 !important;
}

:deep(.el-tag--warning) {
  background: rgba(245, 158, 11, 0.12) !important;
  color: #f59e0b !important;
}

:deep(.el-tag--danger) {
  background: rgba(239, 68, 68, 0.12) !important;
  color: #ef4444 !important;
}

:deep(.el-tag--info) {
  background: rgba(107, 114, 128, 0.1) !important;
  color: #6b7280 !important;
}

/* 输入框 */
:deep(.el-input__wrapper) {
  border-radius: 12px !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2) !important;
}

:deep(.el-select .el-input__wrapper),
:deep(.el-cascader .el-input__wrapper) {
  border-radius: 12px !important;
}

/* 分页 */
:deep(.el-pagination) {
  margin-top: 24px !important;
}

:deep(.el-pager li),
:deep(.el-pager li.btn-quicknext),
:deep(.el-pager li.btn-quickprev),
:deep(.el-pagination .btn-prev),
:deep(.el-pagination .btn-next) {
  border-radius: 10px !important;
  font-weight: 600 !important;
  transition: all 0.2s ease !important;
}

:deep(.el-pager li.is-active) {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-light)) !important;
  border-color: transparent !important;
}

/* 加载态 */
.el-loading-mask {
  backdrop-filter: blur(4px) !important;
  -webkit-backdrop-filter: blur(4px) !important;
  background: rgba(255, 255, 255, 0.7) !important;
}

.el-dark .el-loading-mask {
  background: rgba(15, 23, 42, 0.7) !important;
}

/* 标题 */
h1, h2, h3, h4 {
  color: var(--text-color-primary) !important;
}

/* 进场动画 */
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(18px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 响应式断点 */
@media (max-width: 768px) {
  .stat-content {
    flex-direction: column !important;
    align-items: flex-start !important;
  }

  .stat-number {
    font-size: 22px !important;
  }

  .stat-icon {
    width: 40px !important;
    height: 40px !important;
    font-size: 18px !important;
  }

  .resource-header {
    flex-direction: column !important;
    align-items: flex-start !important;
  }

  .resource-numbers {
    display: flex !important;
    align-items: center !important;
    gap: 4px !important;
    flex-wrap: wrap !important;
  }
}

@media (max-width: 480px) {
  .stat-number {
    font-size: 20px !important;
  }

  .stat-label {
    font-size: 12px !important;
  }

  .resource-number-value {
    font-size: 18px !important;
  }
}
</style>
