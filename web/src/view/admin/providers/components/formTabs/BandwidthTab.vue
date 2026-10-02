<template>
  <el-form
    :model="modelValue"
    label-width="120px"
    class="server-form"
  >
    <el-divider content-position="left">
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.bandwidthLimits') }}</span>
    </el-divider>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item
          :label="$t('admin.providers.defaultInboundBandwidth')"
          prop="defaultInboundBandwidth"
        >
          <el-input-number
            v-model="modelValue.defaultInboundBandwidth"
            :min="1"
            :max="1000000"
            :step="50"
            :controls="false"
            placeholder="300"
            style="width: 100%"
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
            {{ $t('admin.providers.defaultInboundBandwidthTip') }}
          </el-text>
        </div>
      </el-col>
      <el-col :span="12">
        <el-form-item
          :label="$t('admin.providers.defaultOutboundBandwidth')"
          prop="defaultOutboundBandwidth"
        >
          <el-input-number
            v-model="modelValue.defaultOutboundBandwidth"
            :min="1"
            :max="1000000"
            :step="50"
            :controls="false"
            placeholder="300"
            style="width: 100%"
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
            {{ $t('admin.providers.defaultOutboundBandwidthTip') }}
          </el-text>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item
          :label="$t('admin.providers.maxInboundBandwidth')"
          prop="maxInboundBandwidth"
        >
          <el-input-number
            v-model="modelValue.maxInboundBandwidth"
            :min="1"
            :max="1000000"
            :step="50"
            :controls="false"
            placeholder="1000"
            style="width: 100%"
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
            {{ $t('admin.providers.maxInboundBandwidthTip') }}
          </el-text>
        </div>
      </el-col>
      <el-col :span="12">
        <el-form-item
          :label="$t('admin.providers.maxOutboundBandwidth')"
          prop="maxOutboundBandwidth"
        >
          <el-input-number
            v-model="modelValue.maxOutboundBandwidth"
            :min="1"
            :max="1000000"
            :step="50"
            :controls="false"
            placeholder="1000"
            style="width: 100%"
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
            {{ $t('admin.providers.maxOutboundBandwidthTip') }}
          </el-text>
        </div>
      </el-col>
    </el-row>

    <el-divider content-position="left">
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.trafficStatistics') }}</span>
    </el-divider>

    <el-form-item
      :label="$t('admin.providers.enableTrafficControl')"
      prop="enableTrafficControl"
    >
      <el-switch
        v-model="modelValue.enableTrafficControl"
        :active-text="$t('admin.providers.enabled')"
        :inactive-text="$t('admin.providers.disabled')"
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
        {{ $t('admin.providers.enableTrafficControlTip') }}
      </el-text>
    </div>

    <el-form-item
      v-show="modelValue.enableTrafficControl"
      :label="$t('admin.providers.maxTraffic')"
      prop="maxTraffic"
    >
      <el-input-number
        v-model="maxTrafficTB"
        :min="0.001"
        :max="1000"
        :step="0.1"
        :precision="3"
        :controls="false"
        placeholder="1"
        style="width: 100%"
      />
    </el-form-item>
    <div
      v-show="modelValue.enableTrafficControl"
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.maxTrafficTip') }}
      </el-text>
    </div>

    <el-form-item
      v-show="modelValue.enableTrafficControl"
      :label="$t('admin.providers.trafficCountMode')"
      prop="trafficCountMode"
    >
      <el-select
        v-model="modelValue.trafficCountMode"
        :placeholder="$t('admin.providers.selectTrafficCountMode')"
        style="width: 100%"
      >
        <el-option
          :label="$t('admin.providers.trafficCountModeBoth')"
          value="both"
        />
        <el-option
          :label="$t('admin.providers.trafficCountModeOut')"
          value="out"
        />
        <el-option
          :label="$t('admin.providers.trafficCountModeIn')"
          value="in"
        />
      </el-select>
    </el-form-item>
    <div
      v-show="modelValue.enableTrafficControl"
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.trafficCountModeTip') }}
      </el-text>
    </div>

    <el-form-item
      v-show="modelValue.enableTrafficControl"
      :label="$t('admin.providers.trafficMultiplier')"
      prop="trafficMultiplier"
    >
      <el-input-number
        v-model="modelValue.trafficMultiplier"
        :min="0.1"
        :max="10"
        :step="0.1"
        :precision="2"
        :controls="false"
        placeholder="1.0"
        style="width: 100%"
      />
    </el-form-item>
    <div
      v-show="modelValue.enableTrafficControl"
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.trafficMultiplierTip') }}
      </el-text>
    </div>

    <el-form-item
      v-show="modelValue.enableTrafficControl"
      :label="$t('admin.providers.trafficSyncMethod')"
      prop="trafficSyncMethod"
    >
      <el-select
        v-model="modelValue.trafficSyncMethod"
        :placeholder="$t('admin.providers.selectTrafficSyncMethod')"
        style="width: 100%"
      >
        <el-option
          :label="$t('admin.providers.trafficSyncMethodPmacct')"
          value="pmacct"
        />
        <el-option
          :label="$t('admin.providers.trafficSyncMethodAgent')"
          value="agent"
        />
      </el-select>
    </el-form-item>
    <div
      v-show="modelValue.enableTrafficControl"
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.trafficSyncMethodTip') }}
      </el-text>
    </div>

    <el-divider
      v-show="modelValue.enableTrafficControl"
      content-position="left"
    >
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.trafficStatsConfig') }}</span>
    </el-divider>

    <el-form-item
      v-show="modelValue.enableTrafficControl"
      :label="$t('admin.providers.trafficStatsMode')"
      prop="trafficStatsMode"
    >
      <el-select
        v-model="modelValue.trafficStatsMode"
        :placeholder="$t('admin.providers.selectTrafficStatsMode')"
        style="width: 100%"
        @change="handlePresetChange"
      >
        <el-option
          :label="$t('admin.providers.trafficStatsModeHigh')"
          value="high"
        />
        <el-option
          :label="$t('admin.providers.trafficStatsModeStandard')"
          value="standard"
        />
        <el-option
          :label="$t('admin.providers.trafficStatsModeLight')"
          value="light"
        />
        <el-option
          :label="$t('admin.providers.trafficStatsModeMinimal')"
          value="minimal"
        />
        <el-option
          :label="$t('admin.providers.trafficStatsModeCustom')"
          value="custom"
        />
      </el-select>
    </el-form-item>
    <div
      v-show="modelValue.enableTrafficControl"
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.trafficStatsModeTip') }}
      </el-text>
    </div>

    <!-- 流量统计详细配置 - 始终显示，但非自定义模式为只读 -->
    <el-row
      v-show="modelValue.enableTrafficControl"
      :gutter="20"
    >
      <el-col :span="12">
        <el-form-item
          :label="$t('admin.providers.trafficCollectInterval')"
          prop="trafficCollectInterval"
        >
          <el-input-number
            v-model="modelValue.trafficCollectInterval"
            :min="30"
            :max="300"
            :step="30"
            :controls="false"
            :disabled="modelValue.trafficStatsMode !== 'custom'"
            placeholder="300"
            style="width: 100%"
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
            {{ $t('admin.providers.trafficCollectIntervalTip') }}{{ modelValue.trafficStatsMode !== 'custom' ? '（' + $t('common.presetValue') + '）' : '' }}
          </el-text>
        </div>
      </el-col>
      <el-col :span="12">
        <el-form-item
          :label="$t('admin.providers.trafficCollectBatchSize')"
          prop="trafficCollectBatchSize"
        >
          <el-input-number
            v-model="modelValue.trafficCollectBatchSize"
            :min="1"
            :max="100"
            :step="5"
            :controls="false"
            :disabled="modelValue.trafficStatsMode !== 'custom'"
            placeholder="10"
            style="width: 100%"
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
            {{ $t('admin.providers.trafficCollectBatchSizeTip') }}{{ modelValue.trafficStatsMode !== 'custom' ? '（' + $t('common.presetValue') + '）' : '' }}
          </el-text>
        </div>
      </el-col>
    </el-row>

    <el-row
      v-show="modelValue.enableTrafficControl"
      :gutter="20"
    >
      <el-col :span="12">
        <el-form-item
          :label="$t('admin.providers.trafficLimitCheckInterval')"
          prop="trafficLimitCheckInterval"
        >
          <el-input-number
            v-model="modelValue.trafficLimitCheckInterval"
            :min="60"
            :max="3600"
            :step="30"
            :controls="false"
            :disabled="modelValue.trafficStatsMode !== 'custom'"
            placeholder="600"
            style="width: 100%"
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
            {{ $t('admin.providers.trafficLimitCheckIntervalTip') }}{{ modelValue.trafficStatsMode !== 'custom' ? '（' + $t('common.presetValue') + '）' : '' }}
          </el-text>
        </div>
      </el-col>
      <el-col :span="12">
        <el-form-item
          :label="$t('admin.providers.trafficLimitCheckBatchSize')"
          prop="trafficLimitCheckBatchSize"
        >
          <el-input-number
            v-model="modelValue.trafficLimitCheckBatchSize"
            :min="1"
            :max="100"
            :step="5"
            :controls="false"
            :disabled="modelValue.trafficStatsMode !== 'custom'"
            placeholder="10"
            style="width: 100%"
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
            {{ $t('admin.providers.trafficLimitCheckBatchSizeTip') }}{{ modelValue.trafficStatsMode !== 'custom' ? '（' + $t('common.presetValue') + '）' : '' }}
          </el-text>
        </div>
      </el-col>
    </el-row>
  </el-form>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true
  }
})

// 流量单位转换：TB 转 MB (1TB = 1024 * 1024 MB = 1048576 MB)
const TB_TO_MB = 1048576

// 计算属性：maxTraffic 的 TB 单位显示
const maxTrafficTB = computed({
  get: () => {
    // 从 MB 转换为 TB
    return Number((props.modelValue.maxTraffic / TB_TO_MB).toFixed(3))
  },
  set: (value) => {
    // 从 TB 转换为 MB
    props.modelValue.maxTraffic = Math.round(value * TB_TO_MB)
  }
})

// 预设配置（与后端保持一致）- 简化版本，只保留实际使用的字段
const presets = {
  high: {
    trafficCollectInterval: 30,  // 0.5分钟采集+统计
    trafficCollectBatchSize: 20,
    trafficLimitCheckInterval: 30,  // 30秒检测
    trafficLimitCheckBatchSize: 20,
    trafficAutoResetInterval: 600,  // 10分钟检查
    trafficAutoResetBatchSize: 20
  },
  standard: {
    trafficCollectInterval: 60,  // 1分钟采集+统计
    trafficCollectBatchSize: 15,
    trafficLimitCheckInterval: 60,  // 1分钟检测
    trafficLimitCheckBatchSize: 15,
    trafficAutoResetInterval: 900,  // 15分钟检查
    trafficAutoResetBatchSize: 15
  },
  light: {
    trafficCollectInterval: 90,   // 1.5分钟采集+统计
    trafficCollectBatchSize: 10,
    trafficLimitCheckInterval: 90,   // 1.5分钟检测
    trafficLimitCheckBatchSize: 10,
    trafficAutoResetInterval: 1800,  // 30分钟检查
    trafficAutoResetBatchSize: 10
  },
  minimal: {
    trafficCollectInterval: 120,  // 2分钟采集+统计
    trafficCollectBatchSize: 5,
    trafficLimitCheckInterval: 120,  // 2分钟检测
    trafficLimitCheckBatchSize: 5,
    trafficAutoResetInterval: 3600,  // 60分钟检查
    trafficAutoResetBatchSize: 5
  }
}

// 处理预设模式变更
const handlePresetChange = (mode) => {
  if (mode !== 'custom' && presets[mode]) {
    Object.assign(props.modelValue, presets[mode])
  }
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
