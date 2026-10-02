<template>
  <el-form
    :model="modelValue"
    label-width="180px"
    class="server-form"
  >
    <el-alert
      :title="$t('admin.providers.hardwareConfigTip')"
      type="info"
      :closable="false"
      show-icon
      style="margin-bottom: 20px;"
    />

    <!-- 通用配置（容器和虚拟机都支持） -->
    <el-divider content-position="left">
      <el-text
        type="primary"
        size="large"
      >
        {{ $t('admin.providers.commonConfig') }}
      </el-text>
    </el-divider>

    <!-- 内存交换（容器和虚拟机都支持） -->
    <el-form-item
      :label="$t('admin.providers.containerMemorySwap')"
      prop="containerMemorySwap"
    >
      <el-switch
        v-model="modelValue.containerMemorySwap"
        :active-text="$t('common.enable')"
        :inactive-text="$t('common.disable')"
      />
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 180px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.containerMemorySwapTip') }}
      </el-text>
    </div>

    <!-- 容器专用配置 -->
    <template v-if="modelValue.containerEnabled">
      <el-divider content-position="left">
        <el-text
          type="warning"
          size="large"
        >
          {{ $t('admin.providers.containerOnlyConfig') }}
        </el-text>
      </el-divider>

      <!-- 特权模式 -->
      <el-form-item
        :label="$t('admin.providers.containerPrivileged')"
        prop="containerPrivileged"
      >
        <el-switch
          v-model="modelValue.containerPrivileged"
          :active-text="$t('common.enable')"
          :inactive-text="$t('common.disable')"
        />
      </el-form-item>
      <div
        v-if="modelValue.containerEnabled"
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 180px;"
      >
        <el-text
          size="small"
          type="warning"
        >
          {{ $t('admin.providers.containerPrivilegedTip') }}
        </el-text>
      </div>

      <!-- 容器嵌套 -->
      <el-form-item
        :label="$t('admin.providers.containerAllowNesting')"
        prop="containerAllowNesting"
      >
        <el-switch
          v-model="modelValue.containerAllowNesting"
          :active-text="$t('common.enable')"
          :inactive-text="$t('common.disable')"
        />
      </el-form-item>
      <div
        v-if="modelValue.containerEnabled"
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 180px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.containerAllowNestingTip') }}
        </el-text>
      </div>

      <!-- CPU限制（容器专用：与limits.cpu互斥） -->
      <el-form-item
        :label="$t('admin.providers.containerCpuAllowance')"
        prop="containerCpuAllowance"
      >
        <el-input
          v-model="modelValue.containerCpuAllowance"
          :placeholder="$t('admin.providers.containerCpuAllowancePlaceholder')"
          style="width: 200px"
        />
      </el-form-item>
      <div
        v-if="modelValue.containerEnabled"
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 180px;"
      >
        <el-text
          size="small"
          type="warning"
        >
          {{ $t('admin.providers.containerCpuAllowanceTip') }}
        </el-text>
      </div>

      <!-- 最大进程数 -->
      <el-form-item
        :label="$t('admin.providers.containerMaxProcesses')"
        prop="containerMaxProcesses"
      >
        <el-input-number
          v-model="modelValue.containerMaxProcesses"
          :min="0"
          :max="100000"
          :step="100"
          :controls="false"
          :placeholder="$t('admin.providers.containerMaxProcessesPlaceholder')"
          style="width: 200px"
        />
      </el-form-item>
      <div
        v-if="modelValue.containerEnabled"
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 180px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.containerMaxProcessesTip') }}
        </el-text>
      </div>

      <!-- 磁盘IO限制 -->
      <el-form-item
        :label="$t('admin.providers.containerDiskIoLimit')"
        prop="containerDiskIoLimit"
      >
        <el-input
          v-model="modelValue.containerDiskIoLimit"
          :placeholder="$t('admin.providers.containerDiskIoLimitPlaceholder')"
          style="width: 200px"
        />
      </el-form-item>
      <div
        v-if="modelValue.containerEnabled"
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 180px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.containerDiskIoLimitTip') }}
        </el-text>
      </div>
    </template>

    <!-- 虚拟机配置提示 -->
    <template v-if="modelValue.vmEnabled && !modelValue.containerEnabled">
      <el-divider content-position="left">
        <el-text
          type="info"
          size="large"
        >
          {{ $t('admin.providers.vmConfigNote') }}
        </el-text>
      </el-divider>
      <el-alert
        :title="$t('admin.providers.vmHardwareConfigTip')"
        type="info"
        :closable="false"
        show-icon
      />
    </template>

    <!-- GPU 直通配置（仅 LXD / Incus 节点） -->
    <template v-if="modelValue.type === 'lxd' || modelValue.type === 'incus'">
      <el-divider content-position="left">
        <el-text
          type="warning"
          size="large"
        >
          {{ $t('admin.providers.gpuPassthrough') }}
        </el-text>
      </el-divider>

      <el-alert
        :title="$t('admin.providers.gpuDriverWarning')"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 16px;"
      />

      <el-form-item
        :label="$t('admin.providers.gpuEnabled')"
        prop="gpuEnabled"
      >
        <el-switch
          v-model="modelValue.gpuEnabled"
          :active-text="$t('common.enable')"
          :inactive-text="$t('common.disable')"
        />
      </el-form-item>

      <template v-if="modelValue.gpuEnabled">
        <el-form-item
          :label="$t('admin.providers.gpuDeviceIds')"
          prop="gpuDeviceIds"
        >
          <el-input
            v-model="modelValue.gpuDeviceIds"
            :placeholder="$t('admin.providers.gpuDeviceIdsPlaceholder')"
            style="width: 300px"
          />
          <el-button
            v-if="modelValue.id"
            type="primary"
            plain
            :loading="detectingGpus"
            style="margin-left: 8px;"
            @click="handleDetectGPUs"
          >
            {{ $t('admin.providers.gpuDetect') }}
          </el-button>
        </el-form-item>
        <div
          class="form-tip"
          style="margin-top: -10px; margin-bottom: 15px; margin-left: 180px;"
        >
          <el-text
            size="small"
            type="info"
          >
            {{ $t('admin.providers.gpuDeviceIdsTip') }}
          </el-text>
        </div>

        <!-- 检测到的 GPU 列表 -->
        <div
          v-if="detectedGpus.length > 0"
          style="margin-left: 180px; margin-bottom: 16px;"
        >
          <el-text
            size="small"
            type="success"
            style="display:block; margin-bottom: 8px;"
          >
            {{ $t('admin.providers.gpuDetectedList') }}
          </el-text>
          <el-tag
            v-for="(gpu, idx) in detectedGpus"
            :key="idx"
            style="margin-right: 6px; margin-bottom: 6px; cursor: pointer;"
            @click="selectGpuId(gpu)"
          >
            {{ formatDeviceLabel(gpu, idx) }}
          </el-tag>
        </div>
      </template>
    </template>
  </el-form>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { detectProviderGPUs } from '@/api/admin'

const { t } = useI18n()

const props = defineProps({
  modelValue: {
    type: Object,
    required: true
  }
})

const detectingGpus = ref(false)
const detectedGpus = ref([])

async function handleDetectGPUs() {
  const providerId = props.modelValue?.id
  if (!providerId) return

  detectingGpus.value = true
  try {
    const res = await detectProviderGPUs(providerId)
    const gpuList = res?.data?.gpus || []
    const npuList = res?.data?.npus || []
    detectedGpus.value = res?.data?.accelerators || [...gpuList, ...npuList]
    if (detectedGpus.value.length === 0) {
      ElMessage.info(t('admin.providers.gpuNoneFound'))
    } else {
      ElMessage.success(t('admin.providers.gpuDetectSuccess', {
        count: detectedGpus.value.length,
        gpuCount: gpuList.length,
        npuCount: npuList.length
      }))
    }
  } catch (e) {
    ElMessage.error(t('admin.providers.gpuDetectFailed'))
  } finally {
    detectingGpus.value = false
  }
}

function selectGpuId(gpu) {
  if (gpu?.kind === 'npu') {
    ElMessage.info(t('admin.providers.npuIdHint'))
    return
  }
  const id = gpu.id
  if (!id) return
  const current = props.modelValue.gpuDeviceIds || ''
  const ids = current ? current.split(',').map(s => s.trim()).filter(Boolean) : []
  if (!ids.includes(String(id))) {
    ids.push(String(id))
    props.modelValue.gpuDeviceIds = ids.join(',')
  }
}

function formatDeviceLabel(gpu, idx) {
  const typeLabel = gpu?.kind === 'npu'
    ? t('admin.providers.acceleratorNpu')
    : t('admin.providers.acceleratorGpu')
  const idPart = gpu?.id ? `${gpu.id} - ` : `${idx} - `
  const namePart = gpu?.product || gpu?.name || gpu?.card || t('admin.providers.gpuUnknown')
  return `[${typeLabel}] ${idPart}${namePart}`
}
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
