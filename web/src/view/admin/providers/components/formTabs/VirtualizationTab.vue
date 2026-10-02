<template>
  <el-form
    :model="modelValue"
    label-width="120px"
    class="server-form"
  >
    <!-- 虚拟化配置 -->
    <el-divider content-position="left">
      <el-icon><Monitor /></el-icon>
      <span style="margin-left: 8px;">{{ $t('admin.providers.virtualizationConfig') }}</span>
    </el-divider>

    <el-row
      :gutter="20"
      style="margin-bottom: 20px;"
    >
      <el-col :span="12">
        <el-card
          shadow="hover"
          style="height: 100%;"
        >
          <template #header>
            <div style="display: flex; align-items: center; font-weight: 600;">
              <el-icon
                size="18"
                style="margin-right: 8px;"
              >
                <Box />
              </el-icon>
              <span>{{ $t('admin.providers.supportTypes') }}</span>
            </div>
          </template>
          <div
            class="support-type-group"
            style="padding: 10px 0;"
          >
            <el-checkbox
              v-model="modelValue.containerEnabled"
              :disabled="isFixedSupportProvider"
              style="margin-right: 30px;"
            >
              <span style="font-size: 14px;">{{ $t('admin.providers.supportContainer') }}</span>
              <el-tooltip
                :content="$t('admin.providers.containerTech')"
                placement="top"
              >
                <el-icon style="margin-left: 5px;">
                  <InfoFilled />
                </el-icon>
              </el-tooltip>
            </el-checkbox>
            <el-checkbox 
              v-model="modelValue.vmEnabled"
              :disabled="isFixedSupportProvider"
            >
              <span style="font-size: 14px;">{{ $t('admin.providers.supportVM') }}</span>
              <el-tooltip
                :content="$t('admin.providers.vmTech')"
                placement="top"
              >
                <el-icon style="margin-left: 5px;">
                  <InfoFilled />
                </el-icon>
              </el-tooltip>
            </el-checkbox>
          </div>
          <div
            class="form-tip"
            style="margin-top: 10px;"
          >
            <el-text
              size="small"
              type="info"
            >
              {{ $t(supportTypeTip) }}
            </el-text>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card
          shadow="hover"
          style="height: 100%;"
        >
          <template #header>
            <div style="display: flex; align-items: center; font-weight: 600;">
              <el-icon
                size="18"
                style="margin-right: 8px;"
              >
                <DocumentCopy />
              </el-icon>
              <span>{{ $t('admin.providers.instanceLimits') }}</span>
            </div>
          </template>
          <div style="padding: 5px 0;">
            <el-form-item
              :label="$t('admin.providers.maxContainers')"
              label-width="100px"
              style="margin-bottom: 15px;"
            >
              <el-input-number
                v-model="modelValue.maxContainerInstances"
                :min="0"
                :max="1000"
                :step="1"
                :controls="false"
                :placeholder="$t('admin.providers.zeroUnlimited')"
                size="small"
                style="width: 100%"
              />
              <div
                class="form-tip"
                style="margin-top: 5px;"
              >
                <el-text
                  size="small"
                  type="info"
                >
                  {{ $t('admin.providers.maxContainersTip') }}
                </el-text>
              </div>
            </el-form-item>
            
            <el-form-item
              :label="$t('admin.providers.maxVMs')"
              label-width="100px"
              style="margin-bottom: 0;"
            >
              <el-input-number
                v-model="modelValue.maxVMInstances"
                :min="0"
                :max="1000"
                :step="1"
                :controls="false"
                :placeholder="$t('admin.providers.zeroUnlimited')"
                size="small"
                style="width: 100%"
              />
              <div
                class="form-tip"
                style="margin-top: 5px;"
              >
                <el-text
                  size="small"
                  type="info"
                >
                  {{ $t('admin.providers.maxVMsTip') }}
                </el-text>
              </div>
            </el-form-item>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 容器资源限制配置 -->
    <div style="margin-top: 20px;">
      <el-card shadow="hover">
        <template #header>
          <div style="display: flex; align-items: center; justify-content: space-between;">
            <div style="display: flex; align-items: center; font-weight: 600;">
              <el-icon
                size="18"
                style="margin-right: 8px;"
              >
                <Box />
              </el-icon>
              <span>{{ $t('admin.providers.containerResourceLimits') }}</span>
            </div>
            <el-tag
              size="small"
              type="info"
            >
              Container
            </el-tag>
          </div>
        </template>
        <el-alert
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 20px;"
        >
          <template #title>
            <span style="font-size: 13px;">{{ $t('admin.providers.configDescription') }}</span>
          </template>
          <div style="font-size: 12px; line-height: 1.8;">
            {{ $t('admin.providers.enableLimitTip') }}<br>
            {{ $t('admin.providers.noLimitTip') }}
          </div>
        </el-alert>
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="resource-limit-item">
              <div class="resource-limit-label">
                <el-icon><Cpu /></el-icon>
                <span>{{ $t('admin.providers.limitCPU') }}</span>
              </div>
              <el-switch
                v-model="modelValue.containerLimitCpu"
                :active-text="$t('admin.providers.limited')"
                :inactive-text="$t('admin.providers.unlimited')"
                inline-prompt
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949;"
              />
              <div class="resource-limit-tip">
                <el-icon size="12">
                  <InfoFilled />
                </el-icon>
                <span>{{ $t('admin.providers.defaultNoLimitCPU') }}</span>
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="resource-limit-item">
              <div class="resource-limit-label">
                <el-icon><Memo /></el-icon>
                <span>{{ $t('admin.providers.limitMemory') }}</span>
              </div>
              <el-switch
                v-model="modelValue.containerLimitMemory"
                :active-text="$t('admin.providers.limited')"
                :inactive-text="$t('admin.providers.unlimited')"
                inline-prompt
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949;"
              />
              <div class="resource-limit-tip">
                <el-icon size="12">
                  <InfoFilled />
                </el-icon>
                <span>{{ $t('admin.providers.defaultNoLimitMemory') }}</span>
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="resource-limit-item">
              <div class="resource-limit-label">
                <el-icon><Coin /></el-icon>
                <span>{{ $t('admin.providers.limitDisk') }}</span>
              </div>
              <el-switch
                v-model="modelValue.containerLimitDisk"
                :active-text="$t('admin.providers.limited')"
                :inactive-text="$t('admin.providers.unlimited')"
                inline-prompt
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949;"
              />
              <div class="resource-limit-tip">
                <el-icon size="12">
                  <InfoFilled />
                </el-icon>
                <span>{{ $t('admin.providers.defaultLimitDisk') }}</span>
              </div>
            </div>
          </el-col>
        </el-row>
      </el-card>
    </div>

    <!-- 虚拟机资源限制配置 -->
    <div style="margin-top: 20px;">
      <el-card shadow="hover">
        <template #header>
          <div style="display: flex; align-items: center; justify-content: space-between;">
            <div style="display: flex; align-items: center; font-weight: 600;">
              <el-icon
                size="18"
                style="margin-right: 8px;"
              >
                <Monitor />
              </el-icon>
              <span>{{ $t('admin.providers.vmResourceLimits') }}</span>
            </div>
            <el-tag
              size="small"
              type="success"
            >
              Virtual Machine
            </el-tag>
          </div>
        </template>
        <el-alert
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 20px;"
        >
          <template #title>
            <span style="font-size: 13px;">{{ $t('admin.providers.configDescription') }}</span>
          </template>
          <div style="font-size: 12px; line-height: 1.8;">
            {{ $t('admin.providers.enableLimitTip') }}<br>
            {{ $t('admin.providers.noLimitTip') }}
          </div>
        </el-alert>
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="resource-limit-item">
              <div class="resource-limit-label">
                <el-icon><Cpu /></el-icon>
                <span>{{ $t('admin.providers.limitCPU') }}</span>
              </div>
              <el-switch
                v-model="modelValue.vmLimitCpu"
                :active-text="$t('admin.providers.limited')"
                :inactive-text="$t('admin.providers.unlimited')"
                inline-prompt
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949;"
              />
              <div class="resource-limit-tip">
                <el-icon size="12">
                  <InfoFilled />
                </el-icon>
                <span>{{ $t('admin.providers.defaultLimitCPU') }}</span>
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="resource-limit-item">
              <div class="resource-limit-label">
                <el-icon><Memo /></el-icon>
                <span>{{ $t('admin.providers.limitMemory') }}</span>
              </div>
              <el-switch
                v-model="modelValue.vmLimitMemory"
                :active-text="$t('admin.providers.limited')"
                :inactive-text="$t('admin.providers.unlimited')"
                inline-prompt
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949;"
              />
              <div class="resource-limit-tip">
                <el-icon size="12">
                  <InfoFilled />
                </el-icon>
                <span>{{ $t('admin.providers.defaultLimitMemory') }}</span>
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="resource-limit-item">
              <div class="resource-limit-label">
                <el-icon><Coin /></el-icon>
                <span>{{ $t('admin.providers.limitDisk') }}</span>
              </div>
              <el-switch
                v-model="modelValue.vmLimitDisk"
                :active-text="$t('admin.providers.limited')"
                :inactive-text="$t('admin.providers.unlimited')"
                inline-prompt
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949;"
              />
              <div class="resource-limit-tip">
                <el-icon size="12">
                  <InfoFilled />
                </el-icon>
                <span>{{ $t('admin.providers.defaultLimitDisk') }}</span>
              </div>
            </div>
          </el-col>
        </el-row>
      </el-card>
    </div>

    <!-- 存储配置 -->
    <el-form-item
      v-if="showStoragePool"
      :label="$t('admin.providers.storagePool')"
      prop="storagePool"
      style="margin-top: 20px;"
    >
      <el-input
        v-model="modelValue.storagePool"
        :placeholder="$t('admin.providers.storagePoolPlaceholder')"
        maxlength="64"
        show-word-limit
      >
        <template #prepend>
          <el-icon><FolderOpened /></el-icon>
        </template>
      </el-input>
    </el-form-item>
    <div
      v-if="modelValue.type === 'proxmox'"
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.proxmoxStorageTip') }}
      </el-text>
    </div>

    <el-divider content-position="left">
      <el-icon><Memo /></el-icon>
      <span style="margin-left: 8px;">{{ $t('admin.providers.ioRateLimits') }}</span>
    </el-divider>
    <el-alert
      :title="$t('admin.providers.ioRateLimitsTip')"
      type="info"
      :closable="false"
      show-icon
      style="margin-bottom: 16px;"
    />
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <span>{{ $t('admin.providers.containerIoRateLimits') }}</span>
          </template>
          <el-form-item
            :label="$t('admin.providers.readIoLimit')"
            prop="containerReadIoLimit"
          >
            <el-input
              v-model="modelValue.containerReadIoLimit"
              :placeholder="$t('admin.providers.ioRateLimitPlaceholder')"
              clearable
            />
          </el-form-item>
          <el-form-item
            :label="$t('admin.providers.writeIoLimit')"
            prop="containerWriteIoLimit"
          >
            <el-input
              v-model="modelValue.containerWriteIoLimit"
              :placeholder="$t('admin.providers.ioRateLimitPlaceholder')"
              clearable
            />
          </el-form-item>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <span>{{ $t('admin.providers.vmIoRateLimits') }}</span>
          </template>
          <el-form-item
            :label="$t('admin.providers.readIoLimit')"
            prop="vmReadIoLimit"
          >
            <el-input
              v-model="modelValue.vmReadIoLimit"
              :placeholder="$t('admin.providers.ioRateLimitPlaceholder')"
              clearable
            />
          </el-form-item>
          <el-form-item
            :label="$t('admin.providers.writeIoLimit')"
            prop="vmWriteIoLimit"
          >
            <el-input
              v-model="modelValue.vmWriteIoLimit"
              :placeholder="$t('admin.providers.ioRateLimitPlaceholder')"
              clearable
            />
          </el-form-item>
        </el-card>
      </el-col>
    </el-row>
  </el-form>
</template>

<script setup>
import { computed } from 'vue'
import { Monitor, Box, DocumentCopy, InfoFilled, Cpu, Memo, Coin, FolderOpened } from '@element-plus/icons-vue'
import { CONTAINER_ONLY_PROVIDER_TYPES, STORAGE_POOL_PROVIDER_TYPES, VM_ONLY_PROVIDER_TYPES } from '@/utils/providerTypes'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true
  }
})

const containerOnlyTypes = CONTAINER_ONLY_PROVIDER_TYPES
const vmOnlyTypes = VM_ONLY_PROVIDER_TYPES

const isFixedSupportProvider = computed(() => {
  return containerOnlyTypes.includes(props.modelValue.type) || vmOnlyTypes.includes(props.modelValue.type)
})

const showStoragePool = computed(() => {
  return STORAGE_POOL_PROVIDER_TYPES.includes(props.modelValue.type)
})

const supportTypeTip = computed(() => {
  if (containerOnlyTypes.includes(props.modelValue.type)) {
    return 'admin.providers.dockerOnlyContainer'
  }
  if (props.modelValue.type === 'qemu') {
    return 'admin.providers.qemuSupportsLXCAndVM'
  }
  if (vmOnlyTypes.includes(props.modelValue.type)) {
    return 'admin.providers.vmOnlyProviders'
  }
  return 'admin.providers.selectVirtualizationType'
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

.resource-limit-item {
  text-align: center;
  padding: 15px 10px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  transition: all 0.3s;
}

.resource-limit-item:hover {
  border-color: #16a34a;
  background-color: var(--success-bg);
}

.resource-limit-label {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 10px;
  font-weight: 500;
  font-size: 14px;
}

.resource-limit-label .el-icon {
  margin-right: 5px;
}

.resource-limit-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 10px;
  font-size: 12px;
  color: var(--text-color-secondary);
}

.resource-limit-tip .el-icon {
  margin-right: 3px;
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
