<template>
  <div class="snapshot-page">
    <el-card class="summary-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('admin.snapshots.title') }}</span>
          <el-button
            type="primary"
            :loading="loading"
            @click="loadAll"
          >
            {{ t('admin.snapshots.refresh') }}
          </el-button>
        </div>
      </template>
      <el-row :gutter="16">
        <el-col
          :xs="24"
          :sm="8"
          :md="6"
        >
          <el-statistic
            :title="t('admin.snapshots.totalSnapshots')"
            :value="overview.total || 0"
          />
        </el-col>
        <el-col
          :xs="24"
          :sm="8"
          :md="6"
        >
          <el-statistic
            :title="t('admin.snapshots.availableSnapshots')"
            :value="overview.available || 0"
          />
        </el-col>
        <el-col
          :xs="24"
          :sm="8"
          :md="6"
        >
          <el-statistic
            :title="t('admin.snapshots.failedSnapshots')"
            :value="overview.failed || 0"
          />
        </el-col>
        <el-col
          :xs="24"
          :sm="8"
          :md="6"
        >
          <el-statistic
            :title="t('admin.snapshots.schedules')"
            :value="overview.schedules || 0"
          />
        </el-col>
      </el-row>
    </el-card>

    <el-card class="content-card">
      <el-tabs
        v-model="activeTab"
        @tab-change="handleTabChange"
      >
        <el-tab-pane
          :label="t('admin.snapshots.overview')"
          name="snapshots"
        >
          <div class="toolbar">
            <el-select
              v-model="snapshotFilter.instanceId"
              :placeholder="t('admin.snapshots.searchInstancePlaceholder')"
              class="toolbar-input wide"
              filterable
              remote
              clearable
              :remote-method="searchInstances"
              :loading="instanceSearchLoading"
              @focus="searchInstances('')"
            >
              <el-option
                v-for="item in instanceOptions"
                :key="item.id"
                :label="formatInstanceOption(item)"
                :value="item.id"
              />
            </el-select>
            <el-select
              v-model="snapshotFilter.status"
              :placeholder="t('admin.snapshots.status')"
              clearable
              class="toolbar-input"
            >
              <el-option
                :label="t('admin.snapshots.creating')"
                value="creating"
              />
              <el-option
                :label="t('admin.snapshots.available')"
                value="available"
              />
              <el-option
                :label="t('admin.snapshots.failed')"
                value="failed"
              />
            </el-select>
            <el-select
              v-model="snapshotFilter.providerType"
              :placeholder="t('admin.snapshots.provider')"
              clearable
              class="toolbar-input"
            >
              <el-option
                label="Proxmox"
                value="proxmox"
              />
              <el-option
                label="LXD"
                value="lxd"
              />
              <el-option
                label="Incus"
                value="incus"
              />
              <el-option
                label="QEMU/Libvirt"
                value="qemu"
              />
              <el-option
                label="KubeVirt"
                value="kubevirt"
              />
              <el-option
                label="Docker"
                value="docker"
              />
              <el-option
                label="Podman"
                value="podman"
              />
            </el-select>
            <el-button
              type="primary"
              @click="loadSnapshots"
            >
              {{ t('admin.snapshots.query') }}
            </el-button>
            <el-button @click="resetSnapshotFilter">
              {{ t('admin.snapshots.reset') }}
            </el-button>
            <el-button
              type="success"
              @click="openCreateDialog"
            >
              {{ t('admin.snapshots.createSnapshot') }}
            </el-button>
          </div>
          <el-table
            v-loading="loading"
            :data="snapshots"
            border
          >
            <el-table-column
              prop="id"
              label="ID"
              width="80"
            />
            <el-table-column
              prop="name"
              :label="t('admin.snapshots.snapshotName')"
              min-width="160"
            />
            <el-table-column
              prop="instanceName"
              :label="t('admin.snapshots.instance')"
              min-width="160"
            />
            <el-table-column
              prop="providerType"
              :label="t('admin.snapshots.provider')"
              width="120"
            />
            <el-table-column
              prop="instanceType"
              :label="t('admin.snapshots.type')"
              width="100"
            />
            <el-table-column
              prop="source"
              :label="t('admin.snapshots.source')"
              width="100"
            />
            <el-table-column
              prop="status"
              :label="t('admin.snapshots.status')"
              width="110"
            >
              <template #default="{ row }">
                <el-tag :type="snapshotStatusType(row.status)">
                  {{ translateStatus(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              prop="createdAt"
              :label="t('admin.snapshots.createTime')"
              width="180"
            >
              <template #default="{ row }">
                {{ formatDate(row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.snapshots.actions')"
              width="290"
              fixed="right"
            >
              <template #default="{ row }">
                <el-button
                  size="small"
                  type="warning"
                  :disabled="row.status !== 'available'"
                  @click="restoreSnapshot(row)"
                >
                  {{ t('admin.snapshots.restore') }}
                </el-button>
                <el-button
                  size="small"
                  @click="downloadSnapshot(row)"
                >
                  {{ t('admin.snapshots.download') }}
                </el-button>
                <el-button
                  size="small"
                  type="danger"
                  @click="deleteSnapshot(row)"
                >
                  {{ t('admin.snapshots.delete') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-model:current-page="snapshotPagination.page"
            v-model:page-size="snapshotPagination.pageSize"
            :total="snapshotPagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            class="pagination"
            @size-change="loadSnapshots"
            @current-change="loadSnapshots"
          />
        </el-tab-pane>

        <el-tab-pane
          :label="t('admin.snapshots.scheduleSnapshots')"
          name="schedules"
        >
          <div class="toolbar">
            <el-button
              type="success"
              @click="openScheduleDialog"
            >
              {{ t('admin.snapshots.newSchedule') }}
            </el-button>
            <el-button @click="loadSchedules">
              {{ t('admin.snapshots.refresh') }}
            </el-button>
          </div>
          <el-table
            v-loading="loading"
            :data="schedules"
            border
          >
            <el-table-column
              prop="id"
              label="ID"
              width="80"
            />
            <el-table-column
              prop="name"
              :label="t('admin.snapshots.scheduleName')"
              min-width="160"
            />
            <el-table-column
              prop="instanceName"
              :label="t('admin.snapshots.instance')"
              min-width="160"
            />
            <el-table-column
              prop="intervalHours"
              :label="t('admin.snapshots.intervalHours')"
              min-width="170"
            />
            <el-table-column
              prop="retentionDays"
              :label="t('admin.snapshots.retentionDays')"
              min-width="170"
            />
            <el-table-column
              prop="maxSnapshots"
              :label="t('admin.snapshots.maxSnapshots')"
              min-width="150"
            />
            <el-table-column
              prop="enabled"
              :label="t('admin.snapshots.enabled')"
              min-width="100"
            >
              <template #default="{ row }">
                <el-switch
                  v-model="row.enabled"
                  @change="toggleSchedule(row)"
                />
              </template>
            </el-table-column>
            <el-table-column
              prop="nextRunAt"
              :label="t('admin.snapshots.nextRunAt')"
              width="180"
            >
              <template #default="{ row }">
                {{ formatDate(row.nextRunAt) }}
              </template>
            </el-table-column>
            <el-table-column
              prop="lastError"
              :label="t('admin.snapshots.lastError')"
              min-width="180"
              show-overflow-tooltip
            />
            <el-table-column
              :label="t('admin.snapshots.actions')"
              width="110"
              fixed="right"
            >
              <template #default="{ row }">
                <el-button
                  size="small"
                  type="danger"
                  @click="deleteSchedule(row)"
                >
                  {{ t('admin.snapshots.delete') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-model:current-page="schedulePagination.page"
            v-model:page-size="schedulePagination.pageSize"
            :total="schedulePagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            class="pagination"
            @size-change="loadSchedules"
            @current-change="loadSchedules"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog
      v-model="createDialogVisible"
      :title="t('admin.snapshots.createSnapshot')"
      width="620px"
    >
      <el-form
        :model="createForm"
        label-width="120px"
      >
        <el-form-item
          :label="t('admin.snapshots.selectInstances')"
          required
        >
          <el-select
            v-model="createForm.instanceIds"
            multiple
            filterable
            remote
            reserve-keyword
            clearable
            collapse-tags
            collapse-tags-tooltip
            :placeholder="t('admin.snapshots.searchInstancePlaceholder')"
            :remote-method="searchInstances"
            :loading="instanceSearchLoading"
            style="width: 100%"
            @focus="searchInstances('')"
          >
            <el-option
              v-for="item in instanceOptions"
              :key="item.id"
              :label="formatInstanceOption(item)"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('admin.snapshots.snapshotName')">
          <el-input
            v-model="createForm.name"
            :placeholder="t('admin.snapshots.autoNamePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('admin.snapshots.description')">
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

    <el-dialog
      v-model="scheduleDialogVisible"
      :title="t('admin.snapshots.addSchedule')"
      width="560px"
    >
      <el-form
        :model="scheduleForm"
        label-width="120px"
      >
        <el-form-item
          :label="t('admin.snapshots.selectInstance')"
          required
        >
          <el-select
            v-model="scheduleForm.instanceId"
            filterable
            remote
            clearable
            :placeholder="t('admin.snapshots.searchInstancePlaceholder')"
            :remote-method="searchInstances"
            :loading="instanceSearchLoading"
            style="width: 100%"
            @focus="searchInstances('')"
          >
            <el-option
              v-for="item in instanceOptions"
              :key="item.id"
              :label="formatInstanceOption(item)"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="t('admin.snapshots.scheduleName')"
          required
        >
          <el-input v-model="scheduleForm.name" />
        </el-form-item>
        <el-form-item :label="t('admin.snapshots.intervalHours')">
          <el-input-number
            v-model="scheduleForm.intervalHours"
            :min="1"
            :max="720"
            controls-position="right"
          />
        </el-form-item>
        <el-form-item :label="t('admin.snapshots.retentionDays')">
          <el-input-number
            v-model="scheduleForm.retentionDays"
            :min="1"
            :max="365"
            controls-position="right"
          />
        </el-form-item>
        <el-form-item :label="t('admin.snapshots.maxSnapshots')">
          <el-input-number
            v-model="scheduleForm.maxSnapshots"
            :min="1"
            :max="100"
            controls-position="right"
          />
        </el-form-item>
        <el-form-item :label="t('admin.snapshots.enabled')">
          <el-switch v-model="scheduleForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="scheduleDialogVisible = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="createSchedule"
        >
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { snapshotApi, getAllInstances } from '@/api/admin'

const { t, locale } = useI18n()

const activeTab = ref('snapshots')
const loading = ref(false)
const submitting = ref(false)
const overview = reactive({ total: 0, available: 0, failed: 0, schedules: 0 })
const snapshots = ref([])
const schedules = ref([])
const snapshotPagination = reactive({ page: 1, pageSize: 20, total: 0 })
const schedulePagination = reactive({ page: 1, pageSize: 20, total: 0 })
const snapshotFilter = reactive({ instanceId: null, providerType: '', status: '' })
const createDialogVisible = ref(false)
const scheduleDialogVisible = ref(false)
const createForm = reactive({ instanceIds: [], name: '', description: '' })
const scheduleForm = reactive({ instanceId: null, name: '', intervalHours: 24, retentionDays: 7, maxSnapshots: 3, enabled: true })
const instanceOptions = ref([])
const instanceSearchLoading = ref(false)

const loadOverview = async () => {
  const res = await snapshotApi.overview()
  Object.assign(overview, res.data || {})
}

const loadSnapshots = async () => {
  loading.value = true
  try {
    const params = {
      page: snapshotPagination.page,
      pageSize: snapshotPagination.pageSize,
      status: snapshotFilter.status || undefined,
      providerType: snapshotFilter.providerType || undefined,
      instanceId: snapshotFilter.instanceId || undefined
    }
    const res = await snapshotApi.list(params)
    snapshots.value = res.data?.list || []
    snapshotPagination.total = res.data?.total || 0
  } catch (error) {
    ElMessage.error(error.message || t('admin.snapshots.loadSnapshotsFailed'))
  } finally {
    loading.value = false
  }
}

const loadSchedules = async () => {
  loading.value = true
  try {
    const res = await snapshotApi.schedules({ page: schedulePagination.page, pageSize: schedulePagination.pageSize })
    schedules.value = res.data?.list || []
    schedulePagination.total = res.data?.total || 0
  } catch (error) {
    ElMessage.error(error.message || t('admin.snapshots.loadSchedulesFailed'))
  } finally {
    loading.value = false
  }
}

const loadAll = async () => {
  loading.value = true
  try {
    await Promise.all([loadOverview(), loadSnapshots(), loadSchedules(), searchInstances('')])
  } finally {
    loading.value = false
  }
}

const searchInstances = async (keyword = '') => {
  instanceSearchLoading.value = true
  try {
    const res = await getAllInstances({ page: 1, pageSize: 30, name: keyword || undefined })
    instanceOptions.value = res.data?.list || []
  } catch (error) {
    ElMessage.error(error.message || t('admin.snapshots.searchInstancesFailed'))
  } finally {
    instanceSearchLoading.value = false
  }
}

const formatInstanceOption = (item) => {
  if (!item) return ''
  return `${item.name || '-'} (#${item.id}) ${item.provider ? ` / ${item.provider}` : ''}`
}

const handleTabChange = (tab) => {
  if (tab === 'schedules') loadSchedules()
  else loadSnapshots()
}

const resetSnapshotFilter = () => {
  Object.assign(snapshotFilter, { instanceId: null, providerType: '', status: '' })
  snapshotPagination.page = 1
  loadSnapshots()
}

const openCreateDialog = async () => {
  createForm.instanceIds = snapshotFilter.instanceId ? [snapshotFilter.instanceId] : []
  createForm.name = ''
  createForm.description = ''
  createDialogVisible.value = true
  await searchInstances('')
}

const createSnapshot = async () => {
  if (!createForm.instanceIds.length) {
    ElMessage.warning(t('admin.snapshots.fillInstances'))
    return
  }
  submitting.value = true
  try {
    const payload = { instanceIds: createForm.instanceIds, name: createForm.name, description: createForm.description }
    const res = await snapshotApi.batchCreate(payload)
    const failed = res.data?.failed || 0
    if (failed > 0) {
      ElMessage.warning(t('admin.snapshots.batchCreatePartial', { failed }))
    } else {
      ElMessage.success(t('admin.snapshots.createSnapshotSuccess'))
    }
    createDialogVisible.value = false
    await Promise.all([loadOverview(), loadSnapshots()])
  } catch (error) {
    ElMessage.error(error.message || t('admin.snapshots.createSnapshotFailed'))
  } finally {
    submitting.value = false
  }
}

const restoreSnapshot = async (row) => {
  try {
    await ElMessageBox.confirm(t('admin.snapshots.restoreConfirm', { name: row.name }), t('admin.snapshots.restoreTitle'), { type: 'warning' })
    await snapshotApi.restore(row.id)
    ElMessage.success(t('admin.snapshots.restoreSubmitted'))
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(error.message || t('admin.snapshots.restoreFailed'))
  }
}

const deleteSnapshot = async (row) => {
  try {
    await ElMessageBox.confirm(t('admin.snapshots.deleteSnapshotConfirm', { name: row.name }), t('admin.snapshots.deleteTitle'), { type: 'warning' })
    await snapshotApi.delete(row.id)
    ElMessage.success(t('admin.snapshots.deleteSuccess'))
    await Promise.all([loadOverview(), loadSnapshots()])
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(error.message || t('admin.snapshots.deleteFailed'))
  }
}

const downloadSnapshot = async (row) => {
  try {
    const res = await snapshotApi.download(row.id)
    saveBlob(res.data, filenameFromResponse(res) || `snapshot-${row.instanceName}-${row.name}.json`)
  } catch (error) {
    ElMessage.error(error.message || t('admin.snapshots.downloadFailed'))
  }
}

const openScheduleDialog = async () => {
  Object.assign(scheduleForm, { instanceId: null, name: '', intervalHours: 24, retentionDays: 7, maxSnapshots: 3, enabled: true })
  scheduleDialogVisible.value = true
  await searchInstances('')
}

const createSchedule = async () => {
  if (!scheduleForm.instanceId || !scheduleForm.name) {
    ElMessage.warning(t('admin.snapshots.fillScheduleRequired'))
    return
  }
  submitting.value = true
  try {
    await snapshotApi.createSchedule({ ...scheduleForm })
    ElMessage.success(t('admin.snapshots.createScheduleSuccess'))
    scheduleDialogVisible.value = false
    await Promise.all([loadOverview(), loadSchedules()])
  } catch (error) {
    ElMessage.error(error.message || t('admin.snapshots.createScheduleFailed'))
  } finally {
    submitting.value = false
  }
}

const toggleSchedule = async (row) => {
  try {
    await snapshotApi.updateSchedule(row.id, { enabled: row.enabled })
    ElMessage.success(t('admin.snapshots.updateSuccess'))
  } catch (error) {
    row.enabled = !row.enabled
    ElMessage.error(error.message || t('admin.snapshots.updateFailed'))
  }
}

const deleteSchedule = async (row) => {
  try {
    await ElMessageBox.confirm(t('admin.snapshots.deleteScheduleConfirm', { name: row.name }), t('admin.snapshots.deleteTitle'), { type: 'warning' })
    await snapshotApi.deleteSchedule(row.id)
    ElMessage.success(t('admin.snapshots.deleteSuccess'))
    await Promise.all([loadOverview(), loadSchedules()])
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(error.message || t('admin.snapshots.deleteFailed'))
  }
}

const snapshotStatusType = (status) => {
  if (status === 'available') return 'success'
  if (status === 'failed') return 'danger'
  return 'warning'
}

const translateStatus = (status) => {
  const keyMap = { creating: 'creating', available: 'available', failed: 'failed' }
  return keyMap[status] ? t(`admin.snapshots.${keyMap[status]}`) : status
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

onMounted(loadAll)
</script>

<style scoped>
.snapshot-page {
  padding: 20px;
}
.summary-card,
.content-card {
  margin-bottom: 16px;
}
.card-header,
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}
.toolbar {
  justify-content: flex-start;
  margin-bottom: 16px;
}
.toolbar-input {
  width: 180px;
}
.toolbar-input.wide {
  width: 260px;
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
