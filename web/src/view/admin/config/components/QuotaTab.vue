<template>
  <el-form
    v-loading="loading"
    :model="config"
    label-width="140px"
    class="config-form"
  >
    <el-alert
      :title="$t('admin.config.userLevelDesc')"
      type="info"
      :closable="false"
      show-icon
      style="margin-bottom: 20px;"
    >
      <div>{{ $t('admin.config.userLevelHint') }}</div>
      <div style="margin-top: 8px; color: #67C23A;">
        <i class="el-icon-check" />
        {{ $t('admin.config.autoSyncHint') }}
      </div>
      <div style="margin-top: 8px; color: #E6A23C;">
        <i class="el-icon-warning" />
        {{ $t('admin.config.resourceLimitWarning') }}
      </div>
    </el-alert>

    <el-form-item :label="$t('admin.config.newUserDefaultLevel')">
      <el-select
        v-model="config.quota.defaultLevel"
        :placeholder="$t('admin.config.selectDefaultLevel')"
        style="width: 200px"
      >
        <el-option
          v-for="level in levelKeys"
          :key="level"
          :label="$t('admin.config.levelN', { level })"
          :value="level"
        />
      </el-select>
    </el-form-item>

    <el-divider content-position="left">
      <div class="divider-title">
        <span>{{ $t('admin.config.levelLimitsConfig') }}</span>
        <el-button
          type="primary"
          size="small"
          @click="addLevel"
        >
          {{ $t('admin.config.addLevel') }}
        </el-button>
      </div>
    </el-divider>

    <!-- 等级限制配置 -->
    <el-row :gutter="15">
      <el-col
        v-for="level in levelKeys"
        :key="level"
        :span="24"
        style="margin-bottom: 15px;"
      >
        <el-card
          class="level-card"
          :class="{ 'default-level': config.quota.defaultLevel === level }"
          shadow="hover"
        >
          <template #header>
            <div class="level-header">
              <span class="level-title">{{ $t('admin.config.levelNLimits', { level }) }}</span>
              <div class="level-actions">
                <el-tag
                  v-if="config.quota.defaultLevel === level"
                  type="success"
                  size="small"
                >
                  {{ $t('admin.config.defaultLevel') }}
                </el-tag>
                <el-button
                  size="small"
                  type="danger"
                  text
                  :disabled="levelKeys.length <= 1 || config.quota.defaultLevel === level"
                  @click="removeLevel(level)"
                >
                  {{ $t('common.delete') }}
                </el-button>
              </div>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="6">
              <el-form-item :label="$t('admin.config.maxInstances')">
                <el-input-number
                  v-model="config.quota.levelLimits[level]['maxInstances']"
                  :min="1"
                  :max="1000"
                  :controls="false"
                  :step="1"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('admin.config.maxCPU')">
                <el-input-number
                  v-model="config.quota.levelLimits[level]['maxResources']['cpu']"
                  :min="1"
                  :max="10240"
                  :controls="false"
                  :step="1"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('admin.config.maxMemoryMB')">
                <el-input-number
                  v-model="config.quota.levelLimits[level]['maxResources']['memory']"
                  :min="128"
                  :max="10485760"
                  :controls="false"
                  :step="128"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('admin.config.maxDiskMB')">
                <el-input-number
                  v-model="config.quota.levelLimits[level]['maxResources']['disk']"
                  :min="512"
                  :max="1024000000"
                  :controls="false"
                  :step="512"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="6">
              <el-form-item :label="$t('admin.config.maxBandwidthMbps')">
                <el-input-number
                  v-model="config.quota.levelLimits[level]['maxResources']['bandwidth']"
                  :min="1"
                  :max="1000000"
                  :controls="false"
                  :step="1"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('admin.config.trafficLimitMB')">
                <el-input-number
                  v-model="config.quota.levelLimits[level]['maxTraffic']"
                  :min="1024"
                  :max="1024000000"
                  :controls="false"
                  :step="1024"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('admin.config.maxSnapshots')">
                <el-input-number
                  v-model="config.quota.levelLimits[level]['maxSnapshots']"
                  :min="0"
                  :max="1000"
                  :controls="false"
                  :step="1"
                  style="width: 100%"
                />
                <div class="form-item-hint">
                  {{ $t('admin.config.maxSnapshotsHint') }}
                </div>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('admin.config.expiryDays')">
                <el-input-number
                  v-model="config.quota.levelLimits[level]['expiryDays']"
                  :min="0"
                  :max="36500"
                  :controls="false"
                  :step="1"
                  style="width: 100%"
                />
                <div class="form-item-hint">
                  {{ $t('admin.config.expiryDaysHint') }}
                </div>
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
  </el-form>
</template>

<script setup>
import { computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { DEFAULT_QUOTA_LEVEL_LIMITS, buildDefaultLevelLimit, getSortedLevelKeys } from '@/utils/levels'

const { t } = useI18n()

const props = defineProps({
  config: { type: Object, required: true },
  loading: { type: Boolean, default: false }
})

const levelKeys = computed(() => getSortedLevelKeys(props.config.quota.levelLimits))

const addLevel = () => {
  const nextLevel = (levelKeys.value.at(-1) || 0) + 1
  const previousLevel = levelKeys.value.at(-1)
  props.config.quota.levelLimits[nextLevel] = buildDefaultLevelLimit(nextLevel, props.config.quota.levelLimits[previousLevel], DEFAULT_QUOTA_LEVEL_LIMITS)
  ElMessage.success(t('admin.config.levelAdded', { level: nextLevel }))
}

const removeLevel = (level) => {
  if (props.config.quota.defaultLevel === level || levelKeys.value.length <= 1) return
  delete props.config.quota.levelLimits[level]
  ElMessage.success(t('admin.config.levelRemoved', { level }))
}
</script>

<style scoped>
.config-form {
  max-height: 600px;
  overflow-y: auto;
}
.level-card {
  border: 1px solid var(--border-color);
}
.level-card.default-level {
  border-color: var(--el-color-success);
}
.level-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.level-title {
  font-weight: 500;
}
.divider-title,
.level-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}
.divider-title {
  justify-content: space-between;
  width: 100%;
}
.form-item-hint {
  font-size: 12px;
  color: var(--text-color-tertiary);
  margin-top: 4px;
  line-height: 1.4;
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
