<template>
  <el-form
    ref="formRef"
    :model="modelValue"
    :rules="rules"
    label-width="120px"
    class="server-form"
  >
    <el-form-item
      :label="$t('admin.providers.serverName')"
      prop="name"
    >
      <el-input
        v-model="modelValue.name"
        :placeholder="$t('admin.providers.serverNamePlaceholder')"
        maxlength="7"
        show-word-limit
      />
    </el-form-item>
    <el-form-item
      :label="$t('admin.providers.serverType')"
      prop="type"
    >
      <el-select
        v-model="modelValue.type"
        :placeholder="$t('admin.providers.serverTypePlaceholder')"
      >
        <el-option
          v-for="option in PROVIDER_TYPE_OPTIONS"
          :key="option.value"
          :label="$t(option.labelKey)"
          :value="option.value"
        />
      </el-select>
    </el-form-item>

    <!-- SSH模式：主机地址（Agent/本机模式下不需要SSH IP/端口） -->
    <template v-if="!isAgentMode && !isLocalMode">
      <el-form-item
        :label="$t('admin.providers.hostAddress')"
        prop="host"
      >
        <el-input
          v-model="modelValue.host"
          :placeholder="$t('admin.providers.hostPlaceholder')"
        />
      </el-form-item>
      <div
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.hostTip') }}
        </el-text>
      </div>
    </template>

    <el-form-item
      :label="$t('admin.providers.portIP')"
      :prop="isAgentMode || isLocalMode ? '' : 'portIP'"
    >
      <el-input
        v-model="modelValue.portIP"
        :placeholder="$t('admin.providers.portIPPlaceholder')"
      />
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ isAgentMode ? $t('admin.providers.portIPTipAgent') : (isLocalMode ? $t('admin.providers.portIPTipLocal') : $t('admin.providers.portIPTip')) }}
      </el-text>
    </div>

    <!-- SSH模式：SSH端口（Agent/本机模式下不需要） -->
    <template v-if="!isAgentMode && !isLocalMode">
      <el-form-item
        :label="$t('admin.providers.port')"
        prop="port"
      >
        <el-input-number
          v-model="modelValue.port"
          :min="1"
          :max="65535"
          :controls="false"
        />
      </el-form-item>
      <div
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.portTip') }}
        </el-text>
      </div>
    </template>

    <!-- 节点模式选择 -->
    <el-form-item
      :label="$t('admin.providers.nodeMode')"
      prop="discoverMode"
    >
      <el-radio-group v-model="modelValue.discoverMode">
        <el-radio :label="false">
          {{ $t('admin.providers.cleanNode') }}
        </el-radio>
        <el-radio :label="true">
          {{ $t('admin.providers.nodeWithInstances') }}
        </el-radio>
      </el-radio-group>
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.nodeModeTip') }}
      </el-text>
    </div>

    <!-- 发现模式配置 - 仅在选择"有实例的节点"时显示 -->
    <template v-if="modelValue.discoverMode">
      <el-form-item
        :label="$t('admin.providers.autoImport')"
        prop="autoImport"
      >
        <el-switch
          v-model="modelValue.autoImport"
          :active-text="$t('common.enabled')"
          :inactive-text="$t('common.disabled')"
        />
      </el-form-item>
      <div
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.autoImportTip') }}
        </el-text>
      </div>

      <el-form-item
        v-if="modelValue.autoImport"
        :label="$t('admin.providers.autoAdjustQuota')"
        prop="autoAdjustQuota"
      >
        <el-switch
          v-model="modelValue.autoAdjustQuota"
          :active-text="$t('common.enabled')"
          :inactive-text="$t('common.disabled')"
        />
      </el-form-item>
      <div
        v-if="modelValue.autoImport"
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.autoAdjustQuotaTip') }}
        </el-text>
      </div>

      <el-form-item
        v-if="modelValue.autoImport"
        :label="$t('admin.providers.importedInstanceOwner')"
        prop="importedInstanceOwner"
      >
        <el-input
          v-model="modelValue.importedInstanceOwner"
          :placeholder="$t('admin.providers.importedInstanceOwnerPlaceholder')"
        />
      </el-form-item>
      <div
        v-if="modelValue.autoImport"
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.importedInstanceOwnerTip') }}
        </el-text>
      </div>
    </template>

    <el-form-item
      :label="$t('common.status')"
      prop="status"
    >
      <el-select
        v-model="modelValue.status"
        :placeholder="$t('admin.providers.statusPlaceholder')"
      >
        <el-option
          :label="$t('common.enabled')"
          value="active"
        />
        <el-option
          :label="$t('common.disabled')"
          value="inactive"
        />
        <el-option
          :label="$t('common.partial')"
          value="partial"
        />
      </el-select>
    </el-form-item>
    <el-form-item
      :label="$t('admin.providers.architecture')"
      prop="architecture"
    >
      <el-select
        v-model="modelValue.architecture"
        :placeholder="$t('admin.providers.architecturePlaceholder')"
      >
        <el-option
          label="amd64 (x86_64)"
          value="amd64"
        />
        <el-option
          label="arm64 (aarch64)"
          value="arm64"
        />
        <el-option
          label="s390x (IBM Z)"
          value="s390x"
        />
      </el-select>
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.architectureTip') }}
      </el-text>
    </div>
    <el-form-item
      :label="$t('common.description')"
      prop="description"
    >
      <el-input
        v-model="modelValue.description"
        type="textarea"
        :rows="3"
        :placeholder="$t('admin.providers.descriptionPlaceholder')"
      />
    </el-form-item>
  </el-form>
</template>

<script setup>
import { ref, computed } from 'vue'
import { PROVIDER_TYPE_OPTIONS } from '@/utils/providerTypes'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true
  },
  rules: {
    type: Object,
    required: true
  }
})

const isAgentMode = computed(() => props.modelValue.connectionType === 'agent')
const isLocalMode = computed(() => props.modelValue.connectionType === 'local')

// 暴露表单引用供父组件使用
const formRef = ref()
defineExpose({
  formRef
})
</script>

<style scoped>
.server-form {
  max-height: 500px;
  overflow-y: auto;
  padding-right: 10px;
}

.form-tip {
  margin-top: 5px;
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
