<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('admin.providers.taskLog')"
    width="80%"
    style="max-width: 1000px;"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div
      v-if="loading"
      style="text-align: center; padding: 40px;"
    >
      <el-icon
        class="is-loading"
        style="font-size: 32px;"
      >
        <Loading />
      </el-icon>
      <p style="margin-top: 16px;">
        {{ $t('admin.providers.loadingTaskLog') }}
      </p>
    </div>
    <div
      v-else-if="error"
      style="text-align: center; padding: 40px;"
    >
      <el-alert 
        type="error" 
        :title="error" 
        show-icon 
        :closable="false"
      />
    </div>
    <div v-else>
      <!-- 任务基本信息 -->
      <el-card
        v-if="task"
        style="margin-bottom: 20px;"
      >
        <template #header>
          <span>{{ $t('admin.providers.taskInfo') }}</span>
        </template>
        <el-descriptions
          :column="2"
          border
        >
          <el-descriptions-item :label="$t('admin.providers.taskID')">
            {{ task.id }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.providers.provider')">
            {{ task.providerName }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.providers.taskType')">
            {{ task.taskType }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.providers.status')">
            <el-tag :type="getTaskStatusType(task.status)">
              {{ getTaskStatusText(task.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.providers.executor')">
            {{ task.executorName }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.providers.duration')">
            {{ task.duration }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.providers.startTime')">
            {{ task.startedAt ? new Date(task.startedAt).toLocaleString() : '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.providers.completionTime')">
            {{ task.completedAt ? new Date(task.completedAt).toLocaleString() : '-' }}
          </el-descriptions-item>
        </el-descriptions>
        <div
          v-if="task.errorMessage"
          style="margin-top: 16px;"
        >
          <el-alert 
            type="error" 
            :title="task.errorMessage" 
            show-icon 
            :closable="false"
          />
        </div>
      </el-card>

      <!-- 日志内容 -->
      <el-card>
        <template #header>
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <span>{{ $t('admin.providers.executionLog') }}</span>
            <el-button 
              v-if="task && task.logOutput" 
              size="small"
              @click="handleCopyLog"
            >
              {{ $t('admin.providers.copyLog') }}
            </el-button>
          </div>
        </template>
        <div 
          class="task-log-content"
          :style="{
            height: '400px',
            overflow: 'auto',
            backgroundColor: '#1e1e1e',
            color: '#ffffff',
            padding: '16px',
            fontFamily: 'Monaco, Consolas, monospace',
            fontSize: '13px',
            lineHeight: '1.5',
            borderRadius: '4px'
          }"
        >
          <pre v-if="task && task.logOutput">{{ task.logOutput }}</pre>
          <div
            v-else
            style="color: #999; text-align: center; padding: 40px;"
          >
            {{ t('common.noLogContent') }}
          </div>
        </div>
      </el-card>
    </div>

    <template #footer>
      <div style="text-align: center;">
        <el-button @click="handleClose">
          {{ t('common.close') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed } from 'vue'
import { Loading } from '@element-plus/icons-vue'
import { copyToClipboard } from '@/utils/clipboard'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  visible: {
    type: Boolean,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: null
  },
  task: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:visible', 'close'])

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const handleClose = () => {
  emit('close')
}

const handleCopyLog = () => {
  if (props.task && props.task.logOutput) {
    copyToClipboard(props.task.logOutput)
    ElMessage.success(t('admin.providers.logCopied'))
  }
}

const getTaskStatusType = (status) => {
  const statusMap = {
    'pending': 'info',
    'running': 'primary',
    'completed': 'success',
    'failed': 'danger',
    'cancelled': 'warning'
  }
  return statusMap[status] || 'info'
}

const getTaskStatusText = (status) => {
  const statusTextMap = {
    'pending': t('admin.providers.taskStatusPending'),
    'running': t('admin.providers.taskStatusRunning'),
    'completed': t('admin.providers.taskStatusCompleted'),
    'failed': t('admin.providers.taskStatusFailed'),
    'cancelled': t('admin.providers.taskStatusCancelled')
  }
  return statusTextMap[status] || status
}
</script>

<style scoped>
.task-log-content::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.task-log-content::-webkit-scrollbar-track {
  background: #2d2d2d;
  border-radius: 4px;
}

.task-log-content::-webkit-scrollbar-thumb {
  background: #555;
  border-radius: 4px;
}

.task-log-content::-webkit-scrollbar-thumb:hover {
  background: #777;
}

.task-log-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
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
