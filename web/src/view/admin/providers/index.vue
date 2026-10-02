<template>
  <div class="providers-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.providers.title') }}</span>
          <div class="header-actions">
            <!-- 批量操作按钮组 - 仅在选中时显示 -->
            <template v-if="selectedProviders.length > 0">
              <el-button
                type="success"
                :icon="CircleCheck"
                :loading="batchHealthSubmitting"
                @click="handleBatchHealthCheck"
              >
                {{ $t('admin.providers.batchHealthCheck') }} ({{ selectedProviders.length }})
              </el-button>
              <el-button
                type="danger"
                :icon="Delete"
                @click="handleBatchDelete"
              >
                {{ $t('admin.providers.batchDelete') }} ({{ selectedProviders.length }})
              </el-button>
              <el-button
                type="warning"
                :icon="Lock"
                @click="handleBatchFreeze"
              >
                {{ $t('admin.providers.batchFreeze') }} ({{ selectedProviders.length }})
              </el-button>
            </template>
            <el-button
              :icon="Download"
              @click="handleExportCSV"
            >
              {{ $t('admin.providers.exportCsv') }}
            </el-button>
            <el-button
              :icon="Upload"
              @click="triggerImportCSV"
            >
              {{ $t('admin.providers.importCsv') }}
            </el-button>
            <input
              ref="importCsvInput"
              type="file"
              accept=".csv,text/csv"
              style="display: none"
              @change="handleImportCsvFileChange"
            >
            <!-- 添加服务器按钮 -->
            <el-button
              type="primary"
              @click="handleAddProvider"
            >
              {{ $t('admin.providers.addProvider') }}
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 搜索过滤 -->
      <SearchFilter 
        :search-form="searchForm"
        @search="handleSearch"
        @reset="handleReset"
      />
      
      <!-- Provider列表表格 -->
      <ProviderTable
        :loading="loading"
        :providers="providers"
        :current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        @selection-change="handleSelectionChange"
        @edit="editProvider"
        @auto-configure="autoConfigureAPI"
        @traffic-monitor="handleEnableTrafficMonitor"
        @health-check="checkHealth"
        @sync-instances="syncInstances"
        @force-recovery-sync="forceRecoverySync"
        @set-expiry="handleSetProviderExpiry"
        @freeze="freezeServer"
        @unfreeze="unfreezeServer"
        @cleanup-orphans="cleanupOrphans"
        @delete="handleDeleteProvider"
        @size-change="handleSizeChange"
        @page-change="handleCurrentChange"
      />
    </el-card>

    <!-- 模式选择对话框 -->
    <ProviderModeSelectDialog
      v-model:visible="showModeSelectDialog"
      @confirm="handleModeConfirm"
    />

    <!-- 添加/编辑服务器对话框 -->
    <ProviderFormDialog
      v-model:visible="showAddDialog"
      :is-editing="isEditing"
      :provider-data="addProviderForm"
      :grouped-countries="groupedCountries"
      :loading="addProviderLoading"
      @submit="handleProviderFormSubmit"
      @cancel="cancelAddServer"
      @reset-level-limits="resetLevelLimitsToDefault"
      @provider-updated="handleProviderUpdated"
    />

    <!-- 自动配置结果对话框 -->
    <ConfigDialog
      v-model:visible="configDialog.visible"
      :provider="configDialog.provider"
      :show-history="configDialog.showHistory"
      :running-task="configDialog.runningTask"
      :history-tasks="configDialog.historyTasks"
      :pagination="configDialog.pagination"
      @close="configDialog.visible = false"
      @view-task-log="viewTaskLog"
      @view-running-task="viewRunningTask"
      @rerun-configuration="rerunConfiguration"
      @page-change="handleConfigPageChange"
      @page-size-change="handleConfigPageSizeChange"
    />

    <!-- 任务日志查看对话框 -->
    <TaskLogDialog
      v-model:visible="taskLogDialog.visible"
      :loading="taskLogDialog.loading"
      :error="taskLogDialog.error"
      :task="taskLogDialog.task"
      @close="taskLogDialog.visible = false"
    />

    <!-- 监控管理对话框 -->
    <MonitoringManagementDialog
      v-model:visible="trafficMonitorDialog.visible"
      :provider="trafficMonitorDialog.provider"
      :show-history="trafficMonitorDialog.showHistory"
      :task="trafficMonitorDialog.task"
      :running-task="trafficMonitorDialog.runningTask"
      :history-tasks="trafficMonitorDialog.historyTasks"
      :loading="trafficMonitorDialog.loading"
      :pagination="trafficMonitorDialog.pagination"
      @close="resetTrafficMonitorDialog()"
      @refresh="refreshTrafficMonitorTask"
      @view-task-log="viewTrafficMonitorTaskLog"
      @view-running-task="viewRunningTrafficMonitorTask"
      @execute-operation="executeTrafficMonitorOperation"
      @page-change="handleTrafficMonitorPageChange"
      @page-size-change="handleTrafficMonitorPageSizeChange"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { Search, Delete, Lock, Upload, Download, CircleCheck } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import SearchFilter from './components/SearchFilter.vue'
import ConfigDialog from './components/ConfigDialog.vue'
import TaskLogDialog from './components/TaskLogDialog.vue'
import TrafficMonitorTaskDialog from './components/TrafficMonitorTaskDialog.vue'
import MonitoringManagementDialog from './components/MonitoringManagementDialog.vue'
import ProviderTable from './components/ProviderTable.vue'
import ProviderFormDialog from './components/ProviderFormDialog.vue'
import ProviderModeSelectDialog from './components/ProviderModeSelectDialog.vue'
import { useProviderCRUD } from './composables/useProviderCRUD'
import { useProviderForm } from './composables/useProviderForm'
import { useProviderDialogs } from './composables/useProviderDialogs'
import { CONTAINER_ONLY_PROVIDER_TYPES, VM_ONLY_PROVIDER_TYPES } from '@/utils/providerTypes'

const { t } = useI18n()

const {
  providers, selectedProviders, loading, batchHealthSubmitting,
  currentPage, pageSize, total, searchForm,
  loadProviders, handleSearch, handleReset,
  handleSizeChange, handleCurrentChange, handleSelectionChange,
  handleDeleteProvider, handleBatchDelete, handleBatchFreeze, handleBatchHealthCheck,
  handleSetProviderExpiry, freezeServer, unfreezeServer, checkHealth,
  handleExportCSV, handleImportCSV, cleanupOrphans, syncInstances, forceRecoverySync
} = useProviderCRUD()

const importCsvInput = ref(null)

const triggerImportCSV = () => {
  importCsvInput.value?.click()
}

const handleImportCsvFileChange = async (event) => {
  const file = event.target?.files?.[0]
  await handleImportCSV(file)
  if (event.target) {
    event.target.value = ''
  }
}

const {
  showAddDialog, addProviderLoading, isEditing, addProviderForm,
  maxTrafficTB, groupedCountries, getLevelTagType,
  resetLevelLimitsToDefault, cancelAddServer,
  editProvider, submitAddServer
} = useProviderForm(loadProviders)

const handleProviderUpdated = async (updates = {}) => {
  Object.assign(addProviderForm, updates)
  await loadProviders()
}

const showModeSelectDialog = ref(false)

// 点击添加服务器：先弹模式选择对话框
const handleAddProvider = () => {
  showModeSelectDialog.value = true
}

// 模式选择确认后打开表单对话框
const handleModeConfirm = (mode) => {
  isEditing.value = false
  cancelAddServer()
  addProviderForm.connectionType = mode  // 'ssh', 'agent' or 'local'
  if (mode === 'local') {
    addProviderForm.name = '本机'
    addProviderForm.type = 'qemu'
    addProviderForm.host = ''
    addProviderForm.port = 0
    addProviderForm.username = 'root'
    addProviderForm.containerEnabled = true
    addProviderForm.vmEnabled = true
    addProviderForm.networkType = 'nat_ipv4'
    addProviderForm.ipv4PortMappingMethod = 'iptables'
    addProviderForm.ipv6PortMappingMethod = 'native'
  }
  if (mode === 'agent') {
    addProviderForm.enableTrafficControl = true
    addProviderForm.enableResourceMonitoring = true
    addProviderForm.trafficSyncMethod = 'agent'
    addProviderForm.networkType = 'no_port_mapping'
  }
  showAddDialog.value = true
}

const {
  configDialog, taskLogDialog, trafficMonitorDialog,
  viewTaskLog, copyTaskLog, autoConfigureAPI,
  startNewConfiguration, rerunConfiguration, viewRunningTask,
  handleConfigPageChange, handleConfigPageSizeChange,
  handleEnableTrafficMonitor, loadTrafficMonitorHistory,
  openTrafficMonitorDialog, handleTrafficMonitorPageChange,
  handleTrafficMonitorPageSizeChange, executeTrafficMonitorOperation,
  viewTrafficMonitorTaskLog, viewRunningTrafficMonitorTask,
  refreshTrafficMonitorTask, resetTrafficMonitorDialog, debugAuthStatus
} = useProviderDialogs(loadProviders)


// 处理 ProviderFormDialog 提交 — 支持 agent 模式新增后留在对话框
const handleProviderFormSubmit = async (formData) => {
  const result = await submitAddServer(formData)
  if (result?.agentMode && result?.newId) {
    // agent 模式新增成功：切换到编辑模式，停留在对话框，触发生成密钥
    isEditing.value = true
    Object.assign(addProviderForm, formData)
    addProviderForm.id = result.newId
    addProviderForm.agentStatus = 'offline'
    // 通知对话框切换到连接页并自动生成密钥
    // 通过更新 showAddDialog 保持对话框打开（它本来就开着）
    // ProviderFormDialog 会监听到 isEditing 变化
  }
  // 其他情况已在 submitAddServer 内处理（关闭对话框等）
}

// 监听provider类型变化，自动设置虚拟化类型支持和端口映射方式
watch(() => addProviderForm.type, (newType) => {
  // 编辑模式下不自动修改虚拟化类型设置，保持用户已保存的配置
  if (isEditing.value) {
    return
  }
  if (CONTAINER_ONLY_PROVIDER_TYPES.includes(newType)) {
    // Docker/Podman/Containerd/Orbstack只支持容器，使用原生端口映射
    addProviderForm.containerEnabled = true
    addProviderForm.vmEnabled = false
    addProviderForm.ipv4PortMappingMethod = 'native'
    addProviderForm.ipv6PortMappingMethod = 'native'
  } else if (VM_ONLY_PROVIDER_TYPES.includes(newType)) {
    // 本地虚拟化类型仅支持虚拟机，使用iptables端口映射
    addProviderForm.containerEnabled = false
    addProviderForm.vmEnabled = true
    addProviderForm.ipv4PortMappingMethod = 'iptables'
    addProviderForm.ipv6PortMappingMethod = 'iptables'
  } else if (newType === 'qemu' || newType === 'kubevirt') {
    // QEMU/KubeVirt 同时支持容器和虚拟机，端口默认走节点侧映射
    addProviderForm.containerEnabled = true
    addProviderForm.vmEnabled = true
    addProviderForm.ipv4PortMappingMethod = 'iptables'
    addProviderForm.ipv6PortMappingMethod = 'iptables'
  } else if (newType === 'proxmox') {
    // Proxmox支持容器和虚拟机
    addProviderForm.containerEnabled = true
    addProviderForm.vmEnabled = true
    // IPv4: NAT模式下默认iptables，独立IP模式下默认native
    const isNATMode = addProviderForm.networkType === 'nat_ipv4' || addProviderForm.networkType === 'nat_ipv4_ipv6'
    addProviderForm.ipv4PortMappingMethod = isNATMode ? 'iptables' : 'native'
    // IPv6: 默认native
    addProviderForm.ipv6PortMappingMethod = 'native'
  } else if (['lxd', 'incus'].includes(newType)) {
    // LXD/Incus支持容器和虚拟机，默认使用device_proxy
    addProviderForm.containerEnabled = true
    addProviderForm.vmEnabled = true
    addProviderForm.ipv4PortMappingMethod = 'device_proxy'
    addProviderForm.ipv6PortMappingMethod = 'device_proxy'
  } else {
    // 其他类型保持默认设置
    addProviderForm.containerEnabled = true
    addProviderForm.vmEnabled = false
    addProviderForm.ipv4PortMappingMethod = 'device_proxy'
    addProviderForm.ipv6PortMappingMethod = 'device_proxy'
  }
})

// 监听网络类型变化，当Proxmox从NAT改为独立IP时，自动调整端口映射方法
watch(() => [addProviderForm.type, addProviderForm.networkType], ([type, networkType]) => {
  // 编辑模式下不自动修改虚拟化类型设置，但仍需处理端口映射方式的联动
  // 端口映射方式的联动由 MappingTab.vue 组件处理
  if (isEditing.value) {
    return
  }
  
  if (type === 'proxmox') {
    const isNATMode = networkType === 'nat_ipv4' || networkType === 'nat_ipv4_ipv6'
    if (isNATMode) {
      // NAT模式只能使用iptables
      addProviderForm.ipv4PortMappingMethod = 'iptables'
    } else {
      // 独立IP模式默认使用native，但也可以选择iptables
      if (addProviderForm.ipv4PortMappingMethod === 'iptables') {
        // 如果当前是iptables，保持不变
      } else {
        addProviderForm.ipv4PortMappingMethod = 'native'
      }
    }
  }
})

onMounted(() => {
  // 在开发环境下输出调试信息
  if (import.meta.env.DEV) {
    debugAuthStatus()
  }
  loadProviders()
})

</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  min-width: 0;
  
  > span {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-color-primary);
    min-width: 0;
    overflow-wrap: anywhere;
  }
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  justify-content: flex-end;
  flex-wrap: wrap;
  max-width: 100%;
}

.header-actions .el-button {
  margin-left: 0;
}

.filter-container {
  margin-bottom: 20px;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.support-type-group {
  display: flex;
  gap: 15px;
}

.form-tip {
  margin-top: 5px;
}

/* 服务器配置标签页样式 */
.server-config-tabs {
  margin-bottom: 20px;
}

.server-config-tabs .el-tab-pane {
  padding: 20px 0;
}

.server-form {
  max-height: 400px;
  overflow-y: auto;
  padding-right: 10px;
}

.location-cell {
  display: flex;
  align-items: center;
  gap: 5px;
}

.location-cell-vertical {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  font-size: 12px;
}

.location-flag {
  font-size: 20px;
  line-height: 1;
}

.location-country {
  font-weight: 500;
  color: var(--text-color-primary);
  text-align: center;
}

.location-city {
  font-size: 11px;
  color: var(--text-color-secondary);
  text-align: center;
}

.location-empty {
  color: #c0c4cc;
}

.flag-icon {
  font-size: 16px;
}

.support-types {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.el-select .el-input {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.connection-status {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.resource-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 12px;
}

.resource-usage {
  display: flex;
  align-items: center;
  gap: 2px;
  font-weight: 500;
}

.resource-usage .separator {
  color: #c0c4cc;
  margin: 0 2px;
}

.resource-progress {
  width: 100%;
}

.resource-item {
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.resource-item .el-icon {
  font-size: 14px;
  color: var(--text-color-secondary);
}

.resource-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 60px;
  color: #c0c4cc;
}

.sync-time {
  margin-top: 2px;
  text-align: center;
}

.traffic-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 12px;
}

.traffic-usage {
  display: flex;
  align-items: center;
  gap: 2px;
  font-weight: 500;
}

.traffic-usage .separator {
  color: #c0c4cc;
  margin: 0 2px;
}

.traffic-progress {
  width: 100%;
}

.traffic-status {
  text-align: center;
}

/* 资源限制配置样式 */
.resource-limit-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 20px;
  background: var(--neutral-bg);
  border-radius: 8px;
  transition: all 0.3s;
}

.resource-limit-item:hover {
  background: var(--bg-color-hover);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.resource-limit-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color-primary);
}

.resource-limit-label .el-icon {
  font-size: 18px;
  color: #16a34a;
}

.resource-limit-tip {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--text-color-secondary);
  text-align: center;
}

.resource-limit-tip .el-icon {
  color: #16a34a;
}

/* 等级限制配置样式 */
.level-limits-container {
  padding: 10px;
  max-height: 450px;
  overflow-y: auto;
}

/* 自定义滚动条样式 */
.level-limits-container::-webkit-scrollbar {
  width: 8px;
}

.level-limits-container::-webkit-scrollbar-track {
  background: var(--neutral-bg);
  border-radius: 4px;
}

.level-limits-container::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 4px;
}

.level-limits-container::-webkit-scrollbar-thumb:hover {
  background: #909399;
}

.level-config-card {
  margin-bottom: 16px;
  padding: 16px;
  background: var(--neutral-bg);
  border-radius: 6px;
  border: 1px solid var(--border-color);
  transition: all 0.3s;
}

@media (max-width: 768px) {
  .card-header {
    align-items: stretch;
    flex-direction: column;
  }

  .header-actions {
    justify-content: flex-start;
  }

  .header-actions .el-button {
    flex: 1 1 140px;
  }
}

.level-config-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border-color: #c0c4cc;
}

.level-header {
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 2px solid #e4e7ed;
}

.level-form {
  margin-top: 8px;
}

.level-form .el-form-item {
  margin-bottom: 12px;
}

.level-form .el-divider {
  margin: 12px 0;
}

.level-form .form-tip {
  margin-top: 2px;
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
