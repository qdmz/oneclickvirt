<template>
  <el-dialog
    :model-value="visible"
    :title="dialogTitle"
    width="900px"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    @close="handleClose"
  >
    <div v-if="provider">
      <!-- 历史记录视图 -->
      <div v-if="showHistory">
        <el-alert
          :title="$t('admin.providers.trafficMonitorHistory')"
          type="info"
          :closable="false"
          show-icon
          style="margin-bottom: 20px;"
        >
          <template #default>
            <p>{{ $t('admin.providers.trafficMonitorHistoryMessage') }}</p>
          </template>
        </el-alert>

        <!-- 正在运行的任务 -->
        <div
          v-if="runningTask"
          style="margin-bottom: 20px;"
        >
          <el-alert
            :title="$t('admin.providers.runningTrafficMonitorTask')"
            type="warning"
            :closable="false"
            show-icon
          >
            <template #default>
              <p>{{ $t('admin.providers.taskID') }}: {{ runningTask.id }}</p>
              <p>{{ $t('admin.providers.taskType') }}: {{ getTaskTypeLabel(runningTask.taskType) }}</p>
              <p>{{ $t('admin.providers.startTime') }}: {{ formatDateTime(runningTask.startedAt) }}</p>
              <p>{{ $t('admin.providers.progress') }}: {{ runningTask.progress }}%</p>
            </template>
          </el-alert>
        </div>

        <!-- 历史任务列表 -->
        <div v-if="historyTasks.length > 0">
          <h4>{{ $t('admin.providers.trafficMonitorHistoryRecords') }}</h4>
          <el-table
            :data="historyTasks"
            size="small"
            style="margin-bottom: 20px;"
          >
            <el-table-column
              prop="id"
              :label="$t('admin.providers.taskID')"
              min-width="100"
            />
            <el-table-column
              :label="$t('admin.providers.taskType')"
              width="120"
            >
              <template #default="{ row }">
                <el-tag 
                  :type="getTaskTypeTagType(row.taskType)"
                  size="small"
                >
                  {{ getTaskTypeLabel(row.taskType) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('admin.providers.status')"
              min-width="90"
            >
              <template #default="{ row }">
                <el-tag 
                  :type="getTaskStatusTagType(row.status)"
                  size="small"
                >
                  {{ getTaskStatusLabel(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('admin.providers.executionTime')"
              min-width="160"
            >
              <template #default="{ row }">
                {{ formatDateTime(row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('admin.providers.progress')"
              min-width="110"
            >
              <template #default="{ row }">
                <el-progress
                  :percentage="row.progress"
                  :status="row.status === 'failed' ? 'exception' : row.status === 'completed' ? 'success' : undefined"
                />
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('admin.providers.result')"
              show-overflow-tooltip
            >
              <template #default="{ row }">
                <span
                  v-if="row.status === 'completed'"
                  style="color: #67C23A;"
                >
                  ✅ {{ $t('common.success') }}: {{ row.successCount }}/{{ row.totalCount }}
                </span>
                <span
                  v-else-if="row.status === 'failed'"
                  style="color: #F56C6C;"
                >
                  ❌ {{ row.errorMsg || $t('common.failed') }}
                </span>
                <span v-else>{{ row.message || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('common.actions')"
              width="100"
            >
              <template #default="{ row }">
                <el-button 
                  type="primary" 
                  size="small"
                  @click="handleViewTaskLog(row.id)"
                >
                  {{ $t('admin.providers.viewLog') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <!-- 分页组件 -->
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[5, 10, 20, 50]"
            :small="false"
            :background="true"
            layout="total, sizes, prev, pager, next, jumper"
            :total="pagination.total"
            style="justify-content: center;"
            @size-change="handlePageSizeChange"
            @current-change="handlePageChange"
          />
        </div>

        <!-- 操作按钮 -->
        <div style="text-align: center; margin-top: 20px;">
          <el-button 
            v-if="runningTask"
            type="primary"
            @click="handleViewRunningTask"
          >
            {{ $t('admin.providers.viewRunningTaskLog') }}
          </el-button>
          <el-button 
            type="success"
            @click="handleExecuteOperation('enable')"
          >
            {{ $t('admin.providers.enableTrafficMonitor') }}
          </el-button>
          <el-button 
            type="warning"
            @click="handleExecuteOperation('disable')"
          >
            {{ $t('admin.providers.disableTrafficMonitor') }}
          </el-button>
          <el-button 
            type="info"
            @click="handleExecuteOperation('detect')"
          >
            {{ $t('admin.providers.detectTrafficMonitor') }}
          </el-button>
          <el-button @click="handleClose">
            {{ $t('common.close') }}
          </el-button>
        </div>
      </div>

      <!-- 任务执行视图 -->
      <div
        v-else-if="task"
        class="task-container"
      >
        <!-- 任务基本信息 -->
        <el-descriptions
          :column="2"
          border
        >
          <el-descriptions-item :label="$t('admin.providers.trafficMonitorTaskType')">
            <el-tag :type="getTaskTypeTagType(task.taskType)">
              {{ getTaskTypeLabel(task.taskType) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.providers.trafficMonitorTaskStatus')">
            <el-tag :type="getTaskStatusTagType(task.status)">
              {{ getTaskStatusLabel(task.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item
            :label="$t('admin.providers.trafficMonitorTaskProgress')"
            :span="2"
          >
            <div class="progress-container">
              <el-progress 
                :percentage="task.progress" 
                :status="task.status === 'failed' ? 'exception' : task.status === 'completed' ? 'success' : undefined"
              />
              <div class="progress-details">
                <span>{{ $t('common.total') }}: {{ task.totalCount }}</span>
                <span class="success-count">{{ $t('common.success') }}: {{ task.successCount }}</span>
                <span class="failed-count">{{ $t('common.failed') }}: {{ task.failedCount }}</span>
              </div>
            </div>
          </el-descriptions-item>
        </el-descriptions>

        <!-- 任务输出日志 -->
        <div class="output-section">
          <div class="section-header">
            <h4>{{ $t('admin.providers.trafficMonitorTaskOutput') }}</h4>
            <el-button
              v-if="task.status === 'running'"
              type="primary"
              size="small"
              :icon="Refresh"
              :loading="loading"
              @click="$emit('refresh')"
            >
              {{ $t('common.refresh') }}
            </el-button>
          </div>
          <div class="output-content">
            <pre v-if="task.output">{{ task.output }}</pre>
            <el-empty
              v-else
              :description="task.status === 'pending' ? $t('admin.providers.taskExecuting') : $t('common.noData')"
              :image-size="80"
            />
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <el-button @click="handleClose">
        {{ $t('common.close') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  provider: {
    type: Object,
    default: null
  },
  showHistory: {
    type: Boolean,
    default: false
  },
  task: {
    type: Object,
    default: null
  },
  runningTask: {
    type: Object,
    default: null
  },
  historyTasks: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  pagination: {
    type: Object,
    default: () => ({
      page: 1,
      pageSize: 10,
      total: 0
    })
  }
})

const emit = defineEmits(['update:visible', 'close', 'refresh', 'viewTaskLog', 'viewRunningTask', 'executeOperation', 'pageChange', 'pageSizeChange'])

const dialogTitle = computed(() => {
  if (props.showHistory) {
    return t('admin.providers.trafficMonitorManagement')
  }
  return t('admin.providers.trafficMonitorTaskTitle')
})

const handleClose = () => {
  emit('update:visible', false)
  emit('close')
}

const handleViewTaskLog = (taskId) => {
  emit('viewTaskLog', taskId)
}

const handleViewRunningTask = () => {
  emit('viewRunningTask')
}

const handleExecuteOperation = (operation) => {
  emit('executeOperation', operation)
}

const handlePageChange = (page) => {
  emit('pageChange', page)
}

const handlePageSizeChange = (pageSize) => {
  emit('pageSizeChange', pageSize)
}

const formatDateTime = (dateTime) => {
  if (!dateTime) return '-'
  return new Date(dateTime).toLocaleString()
}

const getTaskTypeLabel = (taskType) => {
  const labels = {
    'enable_all': t('admin.providers.trafficMonitorTaskTypeEnableAll'),
    'disable_all': t('admin.providers.trafficMonitorTaskTypeDisableAll'),
    'detect_all': t('admin.providers.trafficMonitorTaskTypeDetectAll')
  }
  return labels[taskType] || taskType
}

const getTaskTypeTagType = (taskType) => {
  const types = {
    'enable_all': 'success',
    'disable_all': 'danger',
    'detect_all': 'info'
  }
  return types[taskType] || 'info'
}

const getTaskStatusLabel = (status) => {
  const labels = {
    'pending': t('admin.providers.trafficMonitorTaskStatusPending'),
    'running': t('admin.providers.trafficMonitorTaskStatusRunning'),
    'completed': t('admin.providers.trafficMonitorTaskStatusCompleted'),
    'failed': t('admin.providers.trafficMonitorTaskStatusFailed')
  }
  return labels[status] || status
}

const getTaskStatusTagType = (status) => {
  const types = {
    'pending': 'info',
    'running': 'warning',
    'completed': 'success',
    'failed': 'danger'
  }
  return types[status] || 'info'
}
</script>

<style scoped>
.task-container {
  max-height: 600px;
  overflow-y: auto;
}

.progress-container {
  width: 100%;
}

.progress-details {
  display: flex;
  gap: 20px;
  margin-top: 8px;
  font-size: 13px;
  color: var(--text-color-secondary);
}

.success-count {
  color: #67c23a;
}

.failed-count {
  color: #f56c6c;
}

.output-section {
  margin-top: 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.section-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color-primary);
}

.output-content {
  background: var(--neutral-bg);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 12px;
  max-height: 400px;
  overflow-y: auto;
}

.output-content pre {
  margin: 0;
  font-family: 'Courier New', Courier, monospace;
  font-size: 12px;
  line-height: 1.6;
  color: var(--text-color-primary);
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 自定义滚动条 */
.task-container::-webkit-scrollbar,
.output-content::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.task-container::-webkit-scrollbar-track,
.output-content::-webkit-scrollbar-track {
  background: var(--neutral-bg);
  border-radius: 4px;
}

.task-container::-webkit-scrollbar-thumb,
.output-content::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 4px;
}

.task-container::-webkit-scrollbar-thumb:hover,
.output-content::-webkit-scrollbar-thumb:hover {
  background: #909399;
}


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
