<template>
  <div class="snapshots-tab">
    <div class="toolbar">
      <el-button v-if="!isReadOnly" type="primary" :disabled="snapshotLimitReached" :loading="submitting" @click="openCreateDialog">
        {{ t('user.instanceDetail.createSnapshot') }}
      </el-button>
      <el-button v-if="!isReadOnly" :disabled="snapshotLimitReached" :loading="uploading" @click="triggerUpload">
        {{ t('user.instanceDetail.uploadSnapshot') }}
      </el-button>
      <el-button
        :loading="loading"
        @click="loadSnapshots"
      >
        {{ t('user.instanceDetail.refresh') }}
      </el-button>
      <span class="snapshot-quota">
        {{ t('user.instanceDetail.snapshotQuota', { used: pagination.total, max: limits.maxSnapshotsPerInstance > 0 ? limits.maxSnapshotsPerInstance : 0 }) }}
      </span>
      <input ref="uploadInput" class="hidden-file-input" type="file" accept="application/json,.json" @change="handleUpload" />
    </div>

    <el-table
      v-loading="loading"
      :data="snapshots"
      border
    >
      <el-table-column
        prop="name"
        :label="t('user.instanceDetail.snapshotName')"
        min-width="160"
      />
      <el-table-column
        prop="description"
        :label="t('user.instanceDetail.description')"
        min-width="180"
        show-overflow-tooltip
      />
      <el-table-column
        prop="status"
        :label="t('user.instanceDetail.status')"
        width="110"
      >
        <template #default="{ row }">
          <el-tag :type="snapshotStatusType(row.status)">
            {{ translateStatus(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        prop="source"
        :label="t('user.instanceDetail.source')"
        width="100"
      />
      <el-table-column
        prop="createdAt"
        :label="t('user.instanceDetail.createdAt')"
        width="180"
      >
        <template #default="{ row }">
          {{ formatDate(row.createdAt) }}
        </template>
      </el-table-column>
      <el-table-column
        :label="t('user.instanceDetail.actions')"
        :width="isReadOnly ? 130 : 290"
        fixed="right"
      >
        <template #default="{ row }">
          <el-button
            v-if="!isReadOnly"
            size="small"
            type="warning"
            :disabled="row.status !== 'available'"
            @click="restoreSnapshot(row)"
          >
            {{ t('user.instanceDetail.restoreSnapshot') }}
          </el-button>
          <el-button
            size="small"
            :disabled="row.status !== 'available'"
            @click="downloadSnapshot(row)"
          >
            {{ t('user.instanceDetail.downloadSnapshot') }}
          </el-button>
          <el-button
            v-if="!isReadOnly"
            size="small"
            type="danger"
            @click="deleteSnapshot(row)"
          >
            {{ t('user.instanceDetail.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-empty
      v-if="!loading && !snapshots.length"
      :description="t('user.instanceDetail.noSnapshots')"
    />

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      :page-sizes="[10, 20, 50]"
      layout="total, sizes, prev, pager, next"
      class="pagination"
      @size-change="loadSnapshots"
      @current-change="loadSnapshots"
    />

    <el-dialog
      v-model="createDialogVisible"
      :title="t('user.instanceDetail.createSnapshot')"
      width="520px"
      append-to-body
      destroy-on-close
      :lock-scroll="false"
    >
      <el-form
        :model="createForm"
        label-width="110px"
      >
        <el-form-item :label="t('user.instanceDetail.snapshotName')">
          <el-input
            v-model="createForm.name"
            :placeholder="t('user.instanceDetail.autoSnapshotName')"
          />
        </el-form-item>
        <el-form-item :label="t('user.instanceDetail.description')">
          <el-input
            v-model="createForm.description"
            type="textarea"
            :rows="3"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="createSnapshot"
        >
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createUserInstanceSnapshot,
  deleteUserSnapshot,
  downloadSharedSnapshot,
  downloadUserSnapshot,
  getSharedInstanceSnapshots,
  getUserInstanceSnapshots,
  getUserLimits,
  restoreUserSnapshot,
  uploadUserSnapshot
} from '@/api/user'

const props = defineProps({
  instanceId: {
    type: [Number, String],
    required: true
  },
  shareToken: {
    type: String,
    default: ''
  },
  readonly: {
    type: Boolean,
    default: false
  }
})

const { t, locale } = useI18n()
const loading = ref(false)
const submitting = ref(false)
const uploading = ref(false)
const snapshots = ref([])
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })
const createDialogVisible = ref(false)
const createForm = reactive({ name: '', description: '' })
const isReadOnly = computed(() => props.readonly || Boolean(props.shareToken))
const uploadInput = ref(null)

// 快照额度：每实例允许数量 与 已用数量
const limits = ref({ maxSnapshotsPerInstance: 0, maxSnapshots: 0 })
const loadLimits = async () => {
  try {
    const res = await getUserLimits()
    limits.value = res.data || { maxSnapshotsPerInstance: 0, maxSnapshots: 0 }
  } catch (e) {
    limits.value = { maxSnapshotsPerInstance: 0, maxSnapshots: 0 }
  }
}
// 是否已达到每实例快照上限（0 表示不允许）
const snapshotLimitReached = computed(() => {
  if (isReadOnly.value) return true
  const cap = limits.value.maxSnapshotsPerInstance
  if (cap <= 0) return true
  return pagination.total >= cap
})

const errorMessage = (error, fallback) => error?.details || error?.message || fallback

const loadSnapshots = async () => {
  if (!props.instanceId) return
  loading.value = true
  try {
    const params = { page: pagination.page, pageSize: pagination.pageSize }
    const res = props.shareToken
      ? await getSharedInstanceSnapshots(props.shareToken, params)
      : await getUserInstanceSnapshots(props.instanceId, params)
    snapshots.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (error) {
    ElMessage.error(errorMessage(error, t('user.instanceDetail.loadSnapshotsFailed')))
  } finally {
    loading.value = false
  }
}

const openCreateDialog = () => {
  if (isReadOnly.value) return
  createForm.name = ''
  createForm.description = ''
  createDialogVisible.value = true
}

const createSnapshot = async () => {
  submitting.value = true
  try {
    await createUserInstanceSnapshot(props.instanceId, { ...createForm })
    ElMessage.success(t('user.instanceDetail.createSnapshotSubmitted'))
    createDialogVisible.value = false
    await loadSnapshots()
  } catch (error) {
    ElMessage.error(errorMessage(error, t('user.instanceDetail.createSnapshotFailed')))
  } finally {
    submitting.value = false
  }
}

const triggerUpload = () => {
  if (isReadOnly.value) return
  uploadInput.value?.click()
}

const handleUpload = async (event) => {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file || isReadOnly.value) return
  uploading.value = true
  try {
    await uploadUserSnapshot(props.instanceId, file)
    ElMessage.success(t('user.instanceDetail.uploadSnapshotSuccess'))
    await loadSnapshots()
  } catch (error) {
    ElMessage.error(errorMessage(error, t('user.instanceDetail.uploadSnapshotFailed')))
  } finally {
    uploading.value = false
  }
}

const restoreSnapshot = async (row) => {
  try {
    await ElMessageBox.confirm(t('user.instanceDetail.restoreSnapshotConfirm', { name: row.name }), t('user.instanceDetail.confirmOperation'), { type: 'warning' })
    await restoreUserSnapshot(row.id)
    ElMessage.success(t('user.instanceDetail.restoreSnapshotSubmitted'))
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(errorMessage(error, t('user.instanceDetail.restoreSnapshotFailed')))
  }
}

const deleteSnapshot = async (row) => {
  try {
    await ElMessageBox.confirm(t('user.instanceDetail.deleteSnapshotConfirm', { name: row.name }), t('user.instanceDetail.confirmOperation'), { type: 'warning' })
    await deleteUserSnapshot(row.id)
    ElMessage.success(t('user.instanceDetail.deleteSuccess'))
    await loadSnapshots()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(errorMessage(error, t('user.instanceDetail.deleteFailed')))
  }
}

const downloadSnapshot = async (row) => {
  try {
    const res = props.shareToken
      ? await downloadSharedSnapshot(props.shareToken, row.id)
      : await downloadUserSnapshot(row.id)
    saveBlob(res.data, filenameFromResponse(res) || `snapshot-${row.instanceName || props.instanceId}-${row.name}.json`)
  } catch (error) {
    ElMessage.error(errorMessage(error, t('user.instanceDetail.downloadSnapshotFailed')))
  }
}

const snapshotStatusType = (status) => {
  if (status === 'available') return 'success'
  if (status === 'failed') return 'danger'
  return 'warning'
}

const translateStatus = (status) => {
  const keyMap = { creating: 'snapshotCreating', available: 'snapshotAvailable', failed: 'snapshotFailed' }
  return keyMap[status] ? t(`user.instanceDetail.${keyMap[status]}`) : status
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString(locale.value)
}

const filenameFromResponse = (response) => {
  const disposition = response?.headers?.['content-disposition'] || ''
  const match = disposition.match(/filename="?([^";]+)"?/i)
  return match ? decodeURIComponent(match[1]) : ''
}

const saveBlob = (blob, filename) => {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

watch(() => [props.instanceId, props.shareToken], () => {
  pagination.page = 1
  loadSnapshots()
})

onMounted(() => {
  loadSnapshots()
  loadLimits()
})
</script>

<style scoped>
.snapshots-tab {
  min-height: 220px;
}
.toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.hidden-file-input {
  display: none;
}
.snapshot-quota {
  font-size: 13px;
  color: #909399;
  margin-left: 4px;
}
.pagination {
  margin-top: 16px;
  justify-content: flex-end;
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
