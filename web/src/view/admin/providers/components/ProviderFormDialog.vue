<template>
  <el-dialog 
    v-model="dialogVisible" 
    :title="isEditing ? $t('admin.providers.editServer') : $t('admin.providers.addServer')" 
    width="1000px"
    :close-on-click-modal="false"
    :before-close="handleBeforeClose"
  >
    <!-- 配置分类标签页 -->
    <el-tabs
      v-model="activeTab"
      type="border-card"
      class="server-config-tabs"
      :lazy="false"
    >
      <!-- 基本信息 -->
      <el-tab-pane
        :label="$t('admin.providers.basicInfo')"
        name="basic"
      >
        <BasicInfoTab
          ref="basicInfoTabRef"
          v-model="formData"
          :rules="rules"
        />
      </el-tab-pane>

      <!-- 连接配置 -->
      <el-tab-pane
        :label="$t('admin.providers.connectionConfig')"
        name="connection"
      >
        <ConnectionTab
          v-model="formData"
          :is-editing="isEditing"
          :testing-connection="testingConnection"
          :connection-test-result="connectionTestResult"
          :generating-secret="generatingSecret"
          :agent-connect-cmd="agentConnectCmd"
          :agent-connect-cmd-github="agentConnectCmdGithub"
          :exec-loading="execLoading"
          :exec-result="execResult"
          :checking-agent-status="checkingAgentStatus"
          @test-connection="handleTestConnection"
          @apply-timeout="handleApplyTimeout"
          @auth-method-change="handleAuthMethodChange"
          @generate-agent-secret="handleGenerateAgentSecret"
          @check-agent-status="handleCheckAgentStatus"
          @exec-command="handleExecCommand"
          @clear-exec-result="execResult = null"
        />
      </el-tab-pane>

      <!-- 地理位置 -->
      <el-tab-pane
        :label="$t('admin.providers.location')"
        name="location"
      >
        <LocationTab
          v-model="formData"
          :grouped-countries="groupedCountries"
        />
      </el-tab-pane>

      <!-- 虚拟化配置 -->
      <el-tab-pane
        :label="$t('admin.providers.virtualizationConfig')"
        name="virtualization"
      >
        <VirtualizationTab
          v-model="formData"
        />
      </el-tab-pane>

      <!-- IP映射配置 -->
      <el-tab-pane
        :label="$t('admin.providers.ipMappingConfig')"
        name="mapping"
      >
        <MappingTab
          v-model="formData"
          @provider-updated="handleProviderUpdated"
        />
      </el-tab-pane>

      <!-- 带宽配置 -->
      <el-tab-pane
        :label="$t('admin.providers.bandwidthConfig')"
        name="bandwidth"
      >
        <BandwidthTab
          v-model="formData"
        />
      </el-tab-pane>

      <!-- 等级限制配置 -->
      <el-tab-pane
        :label="$t('admin.providers.levelLimits')"
        name="levelLimits"
      >
        <LevelLimitsTab
          v-model="formData"
          @reset-defaults="handleResetLevelLimits"
        />
      </el-tab-pane>

      <!-- 过期设置 -->
      <el-tab-pane
        :label="$t('admin.providers.expirySettings')"
        name="expiry"
      >
        <ExpirySettingsTab
          v-model="formData"
          :provider-id="formData.id"
          :is-editing="isEditing"
        />
      </el-tab-pane>

      <!-- 高级设置 -->
      <el-tab-pane
        :label="$t('admin.providers.advancedSettings')"
        name="advanced"
      >
        <AdvancedTab
          v-model="formData"
        />
      </el-tab-pane>

      <!-- 硬件配置（LXD/Incus 容器和虚拟机） -->
      <el-tab-pane
        v-if="showHardwareConfigTab"
        :label="$t('admin.providers.hardwareConfig')"
        name="hardwareConfig"
      >
        <HardwareConfigTab
          v-model="formData"
        />
      </el-tab-pane>
    </el-tabs>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">{{ $t('common.cancel') }}</el-button>
        <el-button
          type="primary"
          :loading="loading"
          @click="handleSubmit"
        >{{ $t('common.save') }}</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
// 导入子标签页组件
import BasicInfoTab from './formTabs/BasicInfoTab.vue'
import ConnectionTab from './formTabs/ConnectionTab.vue'
import LocationTab from './formTabs/LocationTab.vue'
import VirtualizationTab from './formTabs/VirtualizationTab.vue'
import MappingTab from './formTabs/MappingTab.vue'
import BandwidthTab from './formTabs/BandwidthTab.vue'
import LevelLimitsTab from './formTabs/LevelLimitsTab.vue'
import ExpirySettingsTab from './formTabs/ExpirySettingsTab.vue'
import AdvancedTab from './formTabs/AdvancedTab.vue'
import HardwareConfigTab from './formTabs/HardwareConfigTab.vue'
import { useProviderForm } from './composables/useProviderForm'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  isEditing: {
    type: Boolean,
    default: false
  },
  providerData: {
    type: Object,
    default: () => ({})
  },
  groupedCountries: {
    type: Object,
    default: () => ({})
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'submit', 'cancel', 'reset-level-limits', 'provider-updated'])

const {
  dialogVisible,
  activeTab,
  basicInfoTabRef,
  formData,
  rules,
  groupedCountries,
  showHardwareConfigTab,
  testingConnection,
  connectionTestResult,
  generatingSecret,
  checkingAgentStatus,
  agentConnectCmd,
  agentConnectCmdGithub,
  execLoading,
  execResult,
  handleTestConnection,
  handleApplyTimeout,
  handleAuthMethodChange,
  handleGenerateAgentSecret,
  handleExecCommand,
  handleCheckAgentStatus,
  acknowledgePersistedFields,
  handleResetLevelLimits,
  handleSubmit,
  handleBeforeClose,
  handleClose
} = useProviderForm(props, emit)

const handleProviderUpdated = (updates) => {
  acknowledgePersistedFields(updates)
  emit('provider-updated', updates)
}
</script>

<style scoped>
.server-config-tabs {
  margin-bottom: 20px;
}

:deep(.server-form) {
  max-height: 500px;
  overflow-y: auto;
  padding-right: 10px;
}

:deep(.form-tip) {
  margin-top: 5px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

@media (max-width: 768px) {
  .server-config-tabs {
    margin-bottom: 12px;
  }

  :deep(.server-config-tabs .el-tabs__content) {
    padding: 12px;
  }

  :deep(.server-form) {
    max-height: 62vh;
    padding-right: 0;
  }

  :deep(.server-config-tabs [style*="margin-left: 120px"]),
  :deep(.server-config-tabs [style*="margin-left: 160px"]),
  :deep(.server-config-tabs [style*="margin-left: 180px"]) {
    margin-left: 0 !important;
    margin-top: 4px !important;
  }

  :deep(.server-config-tabs [style*="width: 300px"]),
  :deep(.server-config-tabs [style*="width: 400px"]) {
    width: 100% !important;
    max-width: 100%;
  }

  .dialog-footer {
    flex-wrap: wrap;
  }

  .dialog-footer .el-button {
    flex: 1 1 120px;
    margin-left: 0;
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
