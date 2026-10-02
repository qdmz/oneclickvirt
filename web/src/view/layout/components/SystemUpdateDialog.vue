<template>
  <el-dialog
    v-model="visible"
    :title="t('home.footer.updateDialogTitle')"
    width="min(760px, calc(100vw - 24px))"
    destroy-on-close
    class="system-update-dialog"
    @closed="stopPolling"
  >
    <div
      v-loading="loading"
      class="update-dialog-body"
    >
      <div class="update-summary">
        <div class="summary-item">
          <span class="summary-label">{{ t('home.footer.currentVersion') }}</span>
          <code>{{ info.currentVersion || '-' }}</code>
        </div>
        <div class="summary-item">
          <span class="summary-label">{{ t('home.footer.deploymentMode') }}</span>
          <el-tag
            size="small"
            effect="plain"
          >
            {{ info.capability?.mode || 'unknown' }}
          </el-tag>
        </div>
        <div class="summary-item">
          <span class="summary-label">{{ t('home.footer.deploymentFlavor') }}</span>
          <el-tag
            size="small"
            type="info"
            effect="plain"
          >
            {{ info.capability?.flavor || '-' }}
          </el-tag>
        </div>
      </div>

      <el-alert
        v-if="info.capability?.reason"
        :title="info.capability.reason"
        type="info"
        :closable="false"
        show-icon
        class="update-reason"
      />
      <el-alert
        v-if="info.error"
        :title="info.error"
        type="warning"
        :closable="false"
        show-icon
        class="update-reason"
      />

      <el-tabs v-model="activeTab">
        <el-tab-pane
          :label="t('home.footer.updateTab')"
          name="update"
        >
          <div class="version-row">
            <span>{{ t('home.footer.latestVersion') }}</span>
            <a
              v-if="info.releaseUrl"
              :href="info.releaseUrl"
              target="_blank"
              rel="noopener noreferrer"
            >{{ info.latestVersion || '-' }}</a>
            <span v-else>{{ info.latestVersion || '-' }}</span>
          </div>
          <el-select
            v-model="selectedUpdateVersion"
            :placeholder="t('home.footer.selectVersion')"
            clearable
            filterable
            class="version-select"
          >
            <el-option
              v-for="release in updateReleases"
              :key="release.tag"
              :label="release.tag"
              :value="release.tag"
              :disabled="!release.canUpdate"
            >
              <div class="release-option">
                <span>{{ release.tag }}</span>
                <el-tag
                  v-if="!release.canUpdate"
                  size="small"
                  type="warning"
                >
                  {{ t('home.footer.assetUnavailable') }}
                </el-tag>
              </div>
            </el-option>
          </el-select>
          <div class="action-row">
            <el-button
              type="primary"
              :disabled="!info.capability?.canUpdate || !selectedUpdateRelease?.canUpdate"
              :loading="actionLoading"
              @click="submitUpdate"
            >
              {{ t('home.footer.updateNow') }}
            </el-button>
            <el-button
              :disabled="!info.capability?.canRestart"
              :loading="actionLoading"
              @click="submitRestart"
            >
              {{ t('home.footer.restartNow') }}
            </el-button>
          </div>
          <p class="update-note">
            {{ t('home.footer.rollbackDatabaseNote') }}
          </p>
        </el-tab-pane>

        <el-tab-pane
          :label="t('home.footer.rollbackTab')"
          name="rollback"
        >
          <el-select
            v-model="selectedRollback"
            :placeholder="t('home.footer.selectRollbackVersion')"
            value-key="key"
            filterable
            class="version-select"
          >
            <el-option
              v-for="item in rollbackOptions"
              :key="item.key"
              :label="item.label"
              :value="item"
              :disabled="!item.canApply"
            >
              <div class="release-option">
                <span>{{ item.label }}</span>
                <el-tag
                  v-if="item.local"
                  size="small"
                  type="success"
                >
                  {{ t('home.footer.localBackup') }}
                </el-tag>
              </div>
            </el-option>
          </el-select>
          <div class="action-row">
            <el-button
              type="warning"
              :disabled="!info.capability?.canRollback || !selectedRollback?.canApply"
              :loading="actionLoading"
              @click="submitRollback"
            >
              {{ t('home.footer.rollbackNow') }}
            </el-button>
          </div>
          <p class="update-note">
            {{ t('home.footer.rollbackWarning') }}
          </p>
        </el-tab-pane>

        <el-tab-pane
          :label="t('home.footer.commandsTab')"
          name="commands"
        >
          <div
            v-if="!info.capability?.commands?.length"
            class="empty-state"
          >
            {{ t('home.footer.noCommands') }}
          </div>
          <div
            v-for="command in info.capability?.commands || []"
            :key="command.key"
            class="command-item"
          >
            <div class="command-heading">
              <span>{{ command.label }}</span>
              <el-tag
                v-if="command.destructive"
                size="small"
                type="warning"
              >
                {{ t('home.footer.destructiveCommand') }}
              </el-tag>
            </div>
            <p
              v-if="command.description"
              class="command-description"
            >
              {{ command.description }}
            </p>
            <div class="command-line">
              <el-input
                :model-value="resolvedCommand(command)"
                type="textarea"
                :rows="2"
                readonly
              />
              <el-button
                class="copy-command"
                :title="t('home.footer.copyCommand')"
                :aria-label="t('home.footer.copyCommand')"
                @click="copyCommand(resolvedCommand(command))"
              >
                <el-icon><DocumentCopy /></el-icon>
              </el-button>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>

      <el-alert
        v-if="operation && isOperationActive"
        :title="operation.message || t('home.footer.operationRunning')"
        type="info"
        :closable="false"
        show-icon
        class="operation-alert"
      >
        <template #default>
          <span>{{ operation.status }}</span>
          <span v-if="reconnecting"> · {{ t('home.footer.reconnecting') }}</span>
        </template>
      </el-alert>
      <el-alert
        v-if="operation?.status === 'failed'"
        :title="operation.error || t('home.footer.operationFailed')"
        type="error"
        :closable="false"
        show-icon
        class="operation-alert"
      />
    </div>

    <template #footer>
      <el-button @click="visible = false">
        {{ t('common.close') }}
      </el-button>
      <el-button
        :loading="loading"
        @click="loadInfo"
      >
        <el-icon><Refresh /></el-icon>
        {{ t('common.refresh') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { DocumentCopy, Refresh } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import {
  getRollbackVersions,
  getSystemUpdateStatus,
  getUpdateInfo,
  restartSystem,
  startSystemRollback,
  startSystemUpdate
} from '@/api/admin'

const props = defineProps({
  modelValue: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue'])
const { t } = useI18n()
const visible = computed({
  get: () => props.modelValue,
  set: value => emit('update:modelValue', value)
})

const loading = ref(false)
const actionLoading = ref(false)
const activeTab = ref('update')
const info = ref({ capability: { commands: [] }, releases: [], rollbackVersions: [] })
const rollbackReleases = ref([])
const selectedUpdateVersion = ref('')
const selectedRollback = ref(null)
const operation = ref(null)
const reconnecting = ref(false)
let pollTimer = null

const updateReleases = computed(() => (info.value.releases || []).filter(release => release.tag))
const rollbackOptions = computed(() => {
  const options = []
  for (const backup of info.value.rollbackVersions || []) {
    options.push({
      key: `backup:${backup.id}`,
      label: `${backup.version} (${t('home.footer.localBackup')})`,
      version: backup.version,
      backupId: backup.id,
      local: true,
      canApply: Boolean(backup.id)
    })
  }
  for (const release of rollbackReleases.value) {
    if (release.tag && !options.some(option => option.version === release.tag)) {
      options.push({
        key: `release:${release.tag}`,
        label: release.tag,
        version: release.tag,
        local: false,
        canApply: Boolean(release.canRollback)
      })
    }
  }
  return options
})
const isOperationActive = computed(() => ['scheduled', 'staging', 'applying'].includes(operation.value?.status))
const selectedUpdateRelease = computed(() => updateReleases.value.find(release => release.tag === selectedUpdateVersion.value))

const resolvedCommand = (command) => {
  const version = command.key === 'script-rollback'
    ? selectedRollback.value?.version
    : selectedUpdateVersion.value
  if (!version) return command.command
  return command.command.split('<版本号>').join(version)
}

const loadInfo = async () => {
  loading.value = true
  try {
    const response = await getUpdateInfo()
    if (response?.data) {
      info.value = response.data
      operation.value = response.data.operation || operation.value
    }
    if (!selectedUpdateVersion.value && info.value.latestVersion) {
      const latest = updateReleases.value.find(item => item.tag === info.value.latestVersion && item.canUpdate)
      selectedUpdateVersion.value = latest?.tag || ''
    }
    const rollbackResponse = await getRollbackVersions()
    if (rollbackResponse?.data) {
      rollbackReleases.value = rollbackResponse.data.releases || []
      info.value = {
        ...info.value,
        rollbackVersions: rollbackResponse.data.rollbackVersions || info.value.rollbackVersions || [],
        error: rollbackResponse.data.error || info.value.error
      }
    }
    if (visible.value && isOperationActive.value) startPolling()
  } catch (error) {
    ElMessage.error(error?.userMessage || error?.message || t('home.footer.updateLoadFailed'))
  } finally {
    loading.value = false
  }
}

const requireConfirmation = async (message) => {
  try {
    await ElMessageBox.confirm(message, t('common.warning'), {
      type: 'warning',
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel')
    })
    return true
  } catch {
    return false
  }
}

const submitUpdate = async () => {
  if (!selectedUpdateRelease.value?.canUpdate) return
  if (!await requireConfirmation(t('home.footer.updateConfirm', { version: selectedUpdateVersion.value }))) return
  actionLoading.value = true
  try {
    const response = await startSystemUpdate(selectedUpdateVersion.value)
    operation.value = response?.data || null
    ElMessage.success(t('home.footer.operationSubmitted'))
    startPolling()
  } catch (error) {
    ElMessage.error(error?.userMessage || error?.message || t('home.footer.operationFailed'))
  } finally {
    actionLoading.value = false
  }
}

const submitRollback = async () => {
  if (!selectedRollback.value?.canApply) return
  if (!await requireConfirmation(t('home.footer.rollbackConfirm', { version: selectedRollback.value.version }))) return
  actionLoading.value = true
  try {
    const response = await startSystemRollback(selectedRollback.value.version, selectedRollback.value.backupId)
    operation.value = response?.data || null
    ElMessage.success(t('home.footer.operationSubmitted'))
    startPolling()
  } catch (error) {
    ElMessage.error(error?.userMessage || error?.message || t('home.footer.operationFailed'))
  } finally {
    actionLoading.value = false
  }
}

const submitRestart = async () => {
  if (!await requireConfirmation(t('home.footer.restartConfirm'))) return
  actionLoading.value = true
  try {
    const response = await restartSystem()
    operation.value = response?.data || null
    ElMessage.success(t('home.footer.operationSubmitted'))
    startPolling()
  } catch (error) {
    ElMessage.error(error?.userMessage || error?.message || t('home.footer.operationFailed'))
  } finally {
    actionLoading.value = false
  }
}

const startPolling = () => {
  stopPolling()
  pollTimer = window.setInterval(async () => {
    try {
      const response = await getSystemUpdateStatus()
      if (response?.data) operation.value = response.data
      reconnecting.value = false
      if (!isOperationActive.value) {
        stopPolling()
        if (operation.value?.status === 'succeeded') await loadInfo()
      }
    } catch {
      reconnecting.value = true
    }
  }, 2000)
}

const stopPolling = () => {
  if (pollTimer) {
    window.clearInterval(pollTimer)
    pollTimer = null
  }
}

const copyCommand = async (command) => {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(command)
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = command
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      textarea.remove()
    }
    ElMessage.success(t('common.copySuccess'))
  } catch {
    ElMessage.error(t('common.copyFailed'))
  }
}

watch(() => props.modelValue, value => {
  if (value) loadInfo()
  else stopPolling()
})

onBeforeUnmount(stopPolling)
</script>

<style scoped>
.update-dialog-body {
  min-height: 220px;
}

.update-summary {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 24px;
  padding: 2px 0 14px;
  border-bottom: 1px solid var(--border-color);
}

.summary-item,
.version-row,
.command-heading {
  display: flex;
  align-items: center;
  gap: 8px;
}

.command-description {
  margin: 0 0 6px;
  color: var(--text-color-secondary);
  font-size: 12px;
  line-height: 1.5;
}

.summary-label {
  color: var(--text-color-secondary);
  font-size: 12px;
}

.update-reason,
.operation-alert {
  margin-top: 12px;
}

.version-row {
  margin: 4px 0 14px;
  color: var(--text-color-secondary);
}

.version-row a {
  color: var(--primary-color);
}

.version-select {
  width: min(100%, 420px);
}

.release-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.action-row {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 16px;
}

.update-note,
.empty-state {
  color: var(--text-color-secondary);
  font-size: 12px;
  line-height: 1.6;
}

.command-item {
  padding: 10px 0;
  border-bottom: 1px solid var(--border-color-light);
}

.command-heading {
  justify-content: space-between;
  margin-bottom: 6px;
  font-size: 13px;
}

.command-line {
  position: relative;
}

.command-line :deep(.el-textarea__inner) {
  padding-right: 42px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  line-height: 1.45;
  overflow-wrap: anywhere;
}

.copy-command {
  position: absolute;
  top: 6px;
  right: 6px;
  z-index: 1;
}

@media (max-width: 600px) {
  .update-summary {
    display: grid;
    grid-template-columns: 1fr 1fr;
  }

  .version-select {
    width: 100%;
  }
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
