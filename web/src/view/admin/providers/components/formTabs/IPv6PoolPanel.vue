<template>
  <!-- IPv6 地址池/节点地址文件同步。范围按需分配，不会在前端展开。 -->
  <template v-if="managedIPv6NAT">
    <el-divider
      content-position="left"
      style="margin-top: 24px;"
    >
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.ipv6Pool.managedNatTitle') }}</span>
    </el-divider>
    <el-alert
      type="info"
      :closable="false"
      :title="$t('admin.providers.ipv6Pool.managedNatTitle')"
      :description="$t('admin.providers.ipv6Pool.managedNatTip')"
      show-icon
      style="margin-bottom: 16px;"
    />
  </template>
  <template v-else-if="usesStaticIPv6Pool">
    <el-divider
      content-position="left"
      style="margin-top: 24px;"
    >
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.ipv6Pool.management') }}</span>
    </el-divider>
    <el-alert
      v-if="!modelValue.id"
      type="info"
      :closable="false"
      :title="$t('admin.providers.ipv6Pool.newProviderNote')"
      style="margin-bottom: 16px;"
    />
    <template v-else>
      <el-alert
        v-if="!supportsStaticIPv6"
        type="warning"
        :title="$t('admin.providers.ipv6Pool.staticAllocationUnsupported')"
        :description="$t('admin.providers.ipv6Pool.staticAllocationUnsupportedTip', { type: modelValue.type || '-' })"
        :closable="false"
        show-icon
        style="margin-bottom: 16px;"
      />
      <el-alert
        v-else-if="requiresRoutedStaticIPv6"
        type="info"
        :title="$t('admin.providers.ipv6Pool.routedAllocationRequired')"
        :description="$t('admin.providers.ipv6Pool.routedAllocationRequiredTip', { type: modelValue.type || '-' })"
        :closable="false"
        show-icon
        style="margin-bottom: 16px;"
      />
      <el-form-item :label="$t('admin.providers.ipv6Pool.filePath')">
        <div class="ipv6-file-config">
          <el-input
            v-model="modelValue.ipv6AddressFilePath"
            :placeholder="$t('admin.providers.ipv6Pool.filePathPlaceholder')"
            :disabled="!canManageStaticIPv6Pool"
            clearable
          />
          <div class="ipv6-file-actions">
            <el-button
              :icon="DocumentChecked"
              :loading="ipv6FileSaving"
              :disabled="!canManageStaticIPv6Pool || ipv6FileSyncing"
              @click="saveIPv6FilePath"
            >
              {{ $t('admin.providers.ipv6Pool.fileSaveBtn') }}
            </el-button>
            <el-popconfirm
              :title="$t('admin.providers.ipv6Pool.fileClearConfirm')"
              @confirm="clearIPv6FilePath"
            >
              <template #reference>
                <el-button
                  type="danger"
                  plain
                  :icon="Delete"
                  :loading="ipv6FileSaving"
                  :disabled="!hasIPv6FilePath || ipv6FileSyncing"
                >
                  {{ $t('admin.providers.ipv6Pool.fileClearBtn') }}
                </el-button>
              </template>
            </el-popconfirm>
            <el-button
              type="primary"
              :icon="Refresh"
              :loading="ipv6FileSyncing"
              :disabled="!canManageStaticIPv6Pool || !hasIPv6FilePath || ipv6FileSaving"
              @click="syncIPv6File"
            >
              {{ $t('admin.providers.ipv6Pool.syncBtn') }}
            </el-button>
          </div>
          <div class="ipv6-file-state">
            <el-tag
              size="small"
              :type="hasIPv6FilePath ? 'success' : 'info'"
            >
              {{ hasIPv6FilePath ? $t('admin.providers.ipv6Pool.fileConfigured') : $t('admin.providers.ipv6Pool.autoDetection') }}
            </el-tag>
            <el-tag
              size="small"
              :type="supportsStaticIPv6 ? 'success' : 'warning'"
            >
              {{ supportsStaticIPv6 ? $t('admin.providers.ipv6Pool.staticAllocationSupported') : $t('admin.providers.ipv6Pool.staticAllocationUnavailable') }}
            </el-tag>
            <el-text
              size="small"
              type="info"
            >
              {{ hasIPv6FilePath ? $t('admin.providers.ipv6Pool.filePathTip') : $t('admin.providers.ipv6Pool.autoDetectionTip') }}
            </el-text>
          </div>
        </div>
      </el-form-item>

      <el-alert
        v-if="modelValue.ipv6AddressFileSyncError"
        type="error"
        :title="$t('admin.providers.ipv6Pool.lastSyncError')"
        :description="modelValue.ipv6AddressFileSyncError"
        :closable="false"
        show-icon
        style="margin-bottom: 16px;"
      />

      <el-descriptions
        v-if="ipv6SyncResult || ipv6LastSyncedAt"
        :title="$t('admin.providers.ipv6Pool.syncStatus')"
        :column="3"
        border
        size="small"
        class="ipv6-sync-result"
      >
        <el-descriptions-item :label="$t('admin.providers.ipv6Pool.lastSyncedAt')">
          {{ ipv6LastSyncedAt || '-' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('admin.providers.ipv6Pool.parsedCount')">
          {{ ipv6SyncResult?.parsedCount ?? '-' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('admin.providers.ipv6Pool.remoteReadCount')">
          {{ ipv6SyncResult?.remoteReadCount ?? '-' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('admin.providers.ipv6Pool.addedCount')">
          {{ syncItemCount(ipv6SyncResult?.added) }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('admin.providers.ipv6Pool.removedCount')">
          {{ syncItemCount(ipv6SyncResult?.removed) }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('admin.providers.ipv6Pool.preservedCount')">
          {{ syncItemCount(ipv6SyncResult?.preservedAllocated) }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('admin.providers.ipv6Pool.invalidLines')">
          <el-tooltip
            v-if="ipv6SyncResult?.invalidLines?.length"
            :content="ipv6SyncResult.invalidLines.join(', ')"
            placement="top"
          >
            <el-tag
              type="warning"
              size="small"
            >
              {{ ipv6SyncResult.invalidLines.length }}
            </el-tag>
          </el-tooltip>
          <span v-else>{{ ipv6SyncResult ? 0 : '-' }}</span>
        </el-descriptions-item>
      </el-descriptions>

      <el-row
        :gutter="16"
        style="margin-bottom: 16px;"
      >
        <el-col :span="8">
          <el-statistic
            :title="$t('admin.providers.ipv6Pool.total')"
            :value="ipv6PoolStats.total"
          />
        </el-col>
        <el-col :span="8">
          <el-statistic
            :title="$t('admin.providers.ipv6Pool.allocated')"
            :value="ipv6PoolStats.allocated"
          />
        </el-col>
        <el-col :span="8">
          <el-statistic
            :title="$t('admin.providers.ipv6Pool.available')"
            :value="ipv6PoolStats.available"
          />
        </el-col>
      </el-row>
      <el-form-item :label="$t('admin.providers.ipv6Pool.addresses')">
        <div class="ipv6-pool-editor">
          <el-input
            v-model="newIPv6Addresses"
            type="textarea"
            :rows="4"
            :placeholder="$t('admin.providers.ipv6Pool.addressesPlaceholder')"
            :disabled="!canManageStaticIPv6Pool"
          />
          <el-space wrap>
            <el-button
              type="primary"
              :loading="ipv6PoolSaving"
              :disabled="!canManageStaticIPv6Pool"
              @click="addIPv6ToPool"
            >
              {{ $t('admin.providers.ipv6Pool.addBtn') }}
            </el-button>
            <el-popconfirm
              :title="$t('admin.providers.ipv6Pool.clearConfirm')"
              @confirm="clearIPv6Pool"
            >
              <template #reference>
                <el-button
                  type="danger"
                  plain
                >
                  {{ $t('admin.providers.ipv6Pool.clearBtn') }}
                </el-button>
              </template>
            </el-popconfirm>
          </el-space>
        </div>
      </el-form-item>
      <el-form-item :label="$t('admin.providers.ipv6Pool.list')">
        <el-table
          v-loading="ipv6PoolLoading"
          :data="ipv6PoolEntries"
          style="width: 100%"
          size="small"
          max-height="240"
        >
          <el-table-column
            :label="$t('admin.providers.ipv6Pool.address')"
            prop="address"
            min-width="220"
            show-overflow-tooltip
          />
          <el-table-column
            :label="$t('admin.providers.ipv6Pool.status')"
            min-width="100"
          >
            <template #default="{ row }">
              <el-tag
                :type="row.is_range || row.is_reserved ? 'info' : (row.is_allocated ? 'warning' : 'success')"
                size="small"
              >
                {{ row.is_range ? $t('admin.providers.ipv6Pool.statusRange') : (row.is_reserved ? $t('admin.providers.ipv6Pool.statusReserved') : (row.is_allocated ? $t('admin.providers.ipv6Pool.statusAllocated') : $t('admin.providers.ipv6Pool.statusFree'))) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            :label="$t('admin.providers.ipv6Pool.instance')"
            prop="instance_id"
            min-width="100"
          >
            <template #default="{ row }">
              <span>{{ row.instance_id || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column
            width="80"
            align="center"
          >
            <template #default="{ row }">
              <el-popconfirm
                v-if="!row.is_allocated && !row.is_range && !row.is_reserved"
                :title="$t('admin.providers.ipv6Pool.deleteConfirm')"
                @confirm="deleteIPv6Entry(row.id)"
              >
                <template #reference>
                  <el-button
                    type="danger"
                    link
                    size="small"
                  >
                    {{ $t('common.delete') }}
                  </el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
    </template>
  </template>
</template>

<script setup>
import { computed } from 'vue'
import { Delete, DocumentChecked, Refresh } from '@element-plus/icons-vue'
import { useIPv6Pool } from './composables/useIPv6Pool'
import { requiresRoutedStaticIPv6Provider, supportsStaticIPv6Provider } from '@/utils/ipv6Capabilities'
import { usesControllerIPv6Pool, usesManagedIPv6NAT } from '@/utils/networkType'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true
  }
})
const emit = defineEmits(['provider-updated'])

const hasIPv6FilePath = computed(() => Boolean(String(props.modelValue.ipv6AddressFilePath || '').trim()))
const managedIPv6NAT = computed(() => usesManagedIPv6NAT(
  props.modelValue.type,
  props.modelValue.networkType,
  props.modelValue.ipv6PortMappingMethod
))
const usesStaticIPv6Pool = computed(() => usesControllerIPv6Pool(
  props.modelValue.type,
  props.modelValue.networkType,
  props.modelValue.ipv6PortMappingMethod
))
const supportsStaticIPv6 = computed(() => supportsStaticIPv6Provider(props.modelValue.type))
const requiresRoutedStaticIPv6 = computed(() => requiresRoutedStaticIPv6Provider(props.modelValue.type))
const canManageStaticIPv6Pool = computed(() => supportsStaticIPv6.value && !requiresRoutedStaticIPv6.value)
const syncItemCount = value => Array.isArray(value) ? value.length : '-'

const {
  ipv6PoolEntries,
  ipv6PoolStats,
  ipv6PoolLoading,
  newIPv6Addresses,
  ipv6PoolSaving,
  ipv6FileSaving,
  ipv6FileSyncing,
  ipv6SyncResult,
  ipv6LastSyncedAt,
  addIPv6ToPool,
  clearIPv6Pool,
  deleteIPv6Entry,
  saveIPv6FilePath,
  clearIPv6FilePath,
  syncIPv6File,
} = useIPv6Pool(props, updates => emit('provider-updated', updates))
</script>

<style scoped>
.ipv6-file-config,
.ipv6-pool-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.ipv6-file-actions,
.ipv6-file-state {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.ipv6-sync-result {
  margin-bottom: 16px;
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
