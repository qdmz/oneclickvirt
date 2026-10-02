<template>
  <div>
    <el-table
      v-loading="loading"
      :data="providers"
      style="width: 100%"
      @selection-change="$emit('selection-change', $event)"
    >
      <el-table-column
        type="selection"
        width="55"
        fixed="left"
      />
      <el-table-column
        prop="name"
        :label="$t('common.name')"
        width="100"
        fixed="left"
      />
      <el-table-column
        prop="type"
        :label="$t('admin.providers.providerType')"
        min-width="110"
      />
      <el-table-column
        prop="version"
        :label="$t('admin.providers.version')"
        width="100"
      >
        <template #default="scope">
          <div class="version-tags">
            <el-tag
              v-if="scope.row.version && scope.row.version !== ''"
              size="small"
              type="info"
            >
              {{ scope.row.version }}
            </el-tag>
            <el-text
              v-else
              size="small"
              type="info"
            >
              -
            </el-text>
            <el-tag
              v-if="scope.row.type === 'proxmox' && scope.row.pveKvmAvailable === true"
              size="small"
              type="success"
            >
              {{ $t('admin.providers.pveKvmAvailableTrue') }}
            </el-tag>
            <el-tag
              v-else-if="scope.row.type === 'proxmox' && scope.row.pveKvmAvailable === false"
              size="small"
              type="warning"
            >
              {{ $t('admin.providers.pveKvmAvailableFalse') }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.location')"
        min-width="110"
      >
        <template #default="scope">
          <div class="location-cell-vertical">
            <div
              v-if="scope.row.countryCode"
              class="location-flag"
            >
              {{ getFlagEmoji(scope.row.countryCode) }}
            </div>
            <div
              v-if="scope.row.country"
              class="location-country"
            >
              {{ scope.row.country }}
            </div>
            <div
              v-if="scope.row.city"
              class="location-city"
            >
              {{ scope.row.city }}
            </div>
            <div
              v-if="!scope.row.country && !scope.row.city"
              class="location-empty"
            >
              -
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.apiEndpoint')"
        width="140"
      >
        <template #default="scope">
          {{ scope.row.connectionType === 'agent'
            ? (scope.row.agentRemoteIP || '-')
            : (scope.row.connectionType === 'local' ? $t('admin.providers.localConnection') : (scope.row.endpoint ? extractEndpointHost(scope.row.endpoint) : '-')) }}
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.sshPort')"
        width="120"
      >
        <template #default="scope">
          {{ scope.row.connectionType === 'agent' || scope.row.connectionType === 'local' ? '-' : (scope.row.sshPort || 22) }}
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.supportTypes')"
        min-width="150"
      >
        <template #default="scope">
          <div class="support-types">
            <el-tag
              v-if="scope.row.container_enabled"
              size="small"
              type="primary"
            >
              {{ $t('admin.providers.container') }}
            </el-tag>
            <el-tag
              v-if="scope.row.vm_enabled"
              size="small"
              type="success"
            >
              {{ $t('admin.providers.vm') }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        prop="architecture"
        :label="$t('admin.providers.architecture')"
        min-width="140"
      >
        <template #default="scope">
          <el-tag
            size="small"
            type="info"
          >
            {{ scope.row.architecture || 'amd64' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.storagePool')"
        min-width="140"
      >
        <template #default="scope">
          <el-tag
            v-if="scope.row.storagePool"
            size="small"
            type="warning"
          >
            <el-icon style="margin-right: 4px;">
              <FolderOpened />
            </el-icon>
            {{ scope.row.storagePool }}
          </el-tag>
          <el-text
            v-else-if="scope.row.type === 'proxmox'"
            size="small"
            type="info"
          >
            {{ $t('admin.providers.notConfigured') }}
          </el-text>
          <el-text
            v-else
            size="small"
            type="info"
          >
            -
          </el-text>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.connectionStatus')"
        min-width="120"
      >
        <template #default="scope">
          <div class="connection-status">
            <template v-if="scope.row.connectionType === 'agent'">
              <el-tag
                v-if="(scope.row.agentRuntimeStatus || scope.row.agentStatus) === 'online'"
                size="small"
                type="success"
              >
                {{ $t('admin.providers.agentOnlineShort') }}
              </el-tag>
              <el-tag
                v-else
                size="small"
                type="danger"
              >
                {{ $t('admin.providers.agentStatusOffline') }}
              </el-tag>
            </template>
            <template v-else>
              <div class="connection-status-row">
                <span class="status-prefix">A</span>
                <el-tag
                  size="small"
                  :type="getStatusType(scope.row.apiStatus)"
                >
                  {{ getStatusText(scope.row.apiStatus) }}
                </el-tag>
              </div>
              <div class="connection-status-row">
                <span class="status-prefix">S</span>
                <el-tag
                  size="small"
                  :type="getStatusType(scope.row.sshStatus)"
                >
                  {{ getStatusText(scope.row.sshStatus) }}
                </el-tag>
              </div>
            </template>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.cpuResource')"
        width="140"
      >
        <template #default="scope">
          <div
            v-if="scope.row.resourceSynced"
            class="resource-info"
          >
            <div class="resource-usage">
              <span>{{ scope.row.allocatedCpuCores || 0 }}</span>
              <span class="separator">/</span>
              <span>{{ scope.row.nodeCpuCores || 0 }} {{ $t('admin.providers.cores') }}</span>
            </div>
            <div class="resource-progress">
              <el-progress
                :percentage="getResourcePercentage(scope.row.allocatedCpuCores, scope.row.nodeCpuCores)"
                :status="getResourceProgressStatus(scope.row.allocatedCpuCores, scope.row.nodeCpuCores)"
                :stroke-width="6"
                :show-text="false"
              />
            </div>
          </div>
          <div
            v-else
            class="resource-placeholder"
          >
            <el-text
              size="small"
              type="info"
            >
              <el-icon><Loading /></el-icon>
              {{ $t('admin.providers.notSynced') }}
            </el-text>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.memoryResource')"
        width="140"
      >
        <template #default="scope">
          <div
            v-if="scope.row.resourceSynced"
            class="resource-info"
          >
            <div class="resource-usage">
              <span>{{ formatMemorySize(scope.row.allocatedMemory) }}</span>
              <span class="separator">/</span>
              <span>{{ formatMemorySize(scope.row.nodeMemoryTotal) }}</span>
            </div>
            <div class="resource-progress">
              <el-progress
                :percentage="getResourcePercentage(scope.row.allocatedMemory, scope.row.nodeMemoryTotal)"
                :status="getResourceProgressStatus(scope.row.allocatedMemory, scope.row.nodeMemoryTotal)"
                :stroke-width="6"
                :show-text="false"
              />
            </div>
          </div>
          <div
            v-else
            class="resource-placeholder"
          >
            <el-text
              size="small"
              type="info"
            >
              <el-icon><Loading /></el-icon>
              {{ $t('admin.providers.notSynced') }}
            </el-text>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.diskResource')"
        min-width="150"
      >
        <template #default="scope">
          <div
            v-if="scope.row.resourceSynced"
            class="resource-info"
          >
            <div class="resource-usage">
              <span>{{ formatDiskSize(scope.row.allocatedDisk) }}</span>
              <span class="separator">/</span>
              <span>{{ formatDiskSize(scope.row.nodeDiskTotal) }}</span>
            </div>
            <div class="resource-progress">
              <el-progress
                :percentage="getResourcePercentage(scope.row.allocatedDisk, scope.row.nodeDiskTotal)"
                :status="getResourceProgressStatus(scope.row.allocatedDisk, scope.row.nodeDiskTotal)"
                :stroke-width="6"
                :show-text="false"
              />
            </div>
          </div>
          <div
            v-else
            class="resource-placeholder"
          >
            <el-text
              size="small"
              type="info"
            >
              <el-icon><Loading /></el-icon>
              {{ $t('admin.providers.notSynced') }}
            </el-text>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.trafficUsage')"
        min-width="150"
      >
        <template #default="scope">
          <div
            v-if="scope.row.enableTrafficControl"
            class="traffic-info"
          >
            <div class="traffic-usage">
              <span>{{ formatTraffic(scope.row.usedTraffic) }}</span>
              <span class="separator">/</span>
              <span>{{ formatTraffic(scope.row.maxTraffic) }}</span>
            </div>
            <div class="traffic-progress">
              <el-progress
                :percentage="getTrafficPercentage(scope.row.usedTraffic, scope.row.maxTraffic)"
                :status="scope.row.trafficLimited ? 'exception' : getTrafficProgressStatus(scope.row.usedTraffic, scope.row.maxTraffic)"
                :stroke-width="6"
                :show-text="false"
              />
            </div>
            <div
              v-if="scope.row.trafficLimited"
              class="traffic-status"
            >
              <el-tag
                type="danger"
                size="small"
              >
                {{ $t('admin.providers.trafficExceeded') }}
              </el-tag>
            </div>
          </div>
          <div
            v-else
            class="traffic-disabled"
          >
            <el-text
              size="small"
              type="info"
            >
              {{ $t('admin.providers.trafficDisabled') }}
            </el-text>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.instanceQuota')"
        width="160"
      >
        <template #default="scope">
          <div class="instance-quota-info">
            <div
              v-if="scope.row.container_enabled"
              class="quota-item"
            >
              <el-tag
                size="small"
                type="primary"
              >
                {{ $t('admin.providers.container') }}
              </el-tag>
              <span class="quota-text">
                {{ scope.row.currentContainerCount || 0 }} / {{ scope.row.maxContainerInstances === 0 ? '∞' : scope.row.maxContainerInstances }}
              </span>
              <el-progress
                v-if="scope.row.maxContainerInstances > 0"
                :percentage="getQuotaPercentage(scope.row.currentContainerCount, scope.row.maxContainerInstances)"
                :status="getQuotaProgressStatus(scope.row.currentContainerCount, scope.row.maxContainerInstances)"
                :stroke-width="4"
                :show-text="false"
              />
            </div>
            <div
              v-if="scope.row.vm_enabled"
              class="quota-item"
            >
              <el-tag
                size="small"
                type="success"
              >
                {{ $t('admin.providers.vm') }}
              </el-tag>
              <span class="quota-text">
                {{ scope.row.currentVMCount || 0 }} / {{ scope.row.maxVMInstances === 0 ? '∞' : scope.row.maxVMInstances }}
              </span>
              <el-progress
                v-if="scope.row.maxVMInstances > 0"
                :percentage="getQuotaPercentage(scope.row.currentVMCount, scope.row.maxVMInstances)"
                :status="getQuotaProgressStatus(scope.row.currentVMCount, scope.row.maxVMInstances)"
                :stroke-width="4"
                :show-text="false"
              />
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('common.status')"
        min-width="90"
      >
        <template #default="scope">
          <el-tag
            v-if="scope.row.isFrozen"
            type="danger"
            size="small"
          >
            {{ $t('admin.providers.frozen') }}
          </el-tag>
          <el-tag
            v-else-if="isExpired(scope.row.expiresAt)"
            type="warning"
            size="small"
          >
            {{ $t('admin.providers.expired') }}
          </el-tag>
          <el-tag
            v-else
            type="success"
            size="small"
          >
            {{ $t('common.normal') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.expiryTime')"
        width="130"
      >
        <template #default="scope">
          <div v-if="scope.row.expiresAt">
            <el-tag
              :type="isExpired(scope.row.expiresAt) ? 'danger' : isNearExpiry(scope.row.expiresAt) ? 'warning' : 'success'"
              size="small"
            >
              {{ formatDateTime(scope.row.expiresAt) }}
            </el-tag>
          </div>
          <el-text
            v-else
            size="small"
            type="info"
          >
            {{ $t('admin.providers.neverExpires') }}
          </el-text>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('common.actions')"
        width="260"
        fixed="right"
      >
        <template #default="scope">
          <div class="action-buttons">
            <el-button
              size="small"
              type="primary"
              @click="$emit('edit', scope.row)"
            >
              {{ $t('common.edit') }}
            </el-button>
            <el-button
              size="small"
              type="primary"
              @click="$emit('show-actions', scope.row)"
            >
              {{ $t('common.actions') }}
            </el-button>
            <el-button
              size="small"
              type="danger"
              @click="$emit('delete', scope.row)"
            >
              {{ $t('common.delete') }}
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="$emit('size-change', $event)"
        @current-change="$emit('page-change', $event)"
      />
    </div>
  </div>
</template>

<script setup>
import {
  formatMemorySize,
  formatDiskSize,
  formatTraffic,
  getTrafficPercentage,
  getTrafficProgressStatus,
  getResourcePercentage,
  getResourceProgressStatus,
  getQuotaPercentage,
  getQuotaProgressStatus,
  formatDateTime,
  isExpired,
  isNearExpiry,
  getStatusType,
  getStatusText,
  getFlagEmoji
} from '../composables/useProviderUtils'
import { extractEndpointHost } from '@/utils/endpoint'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

defineProps({
  loading: { type: Boolean, default: false },
  providers: { type: Array, default: () => [] },
  currentPage: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
  total: { type: Number, default: 0 }
})

defineEmits(['selection-change', 'edit', 'show-actions', 'delete', 'size-change', 'page-change'])
</script>

<style scoped>
.location-cell-vertical {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  min-height: 75px;
  justify-content: center;
}

.location-flag { font-size: 20px; }
.location-country, .location-city { font-size: 12px; color: #606266; }
.location-empty { color: #c0c4cc; }
.version-tags { display: flex; flex-direction: column; align-items: flex-start; gap: 4px; }
.support-types { display: flex; flex-direction: column; gap: 4px; }
.connection-status { display: flex; flex-direction: column; gap: 4px; }
.connection-status-row { display: flex; align-items: center; gap: 4px; }
.status-prefix { width: 10px; font-size: 11px; color: #909399; line-height: 1; }
.resource-info, .traffic-info { display: flex; flex-direction: column; gap: 4px; }
.resource-usage, .traffic-usage { font-size: 12px; text-align: center; }
.separator { margin: 0 4px; color: #909399; }
.resource-placeholder { text-align: center; }
.traffic-status { text-align: center; }
.traffic-disabled { display: flex; align-items: center; justify-content: center; height: 60px; text-align: center; }
.instance-quota-info { display: flex; flex-direction: column; gap: 8px; }
.quota-item { display: flex; flex-direction: column; gap: 4px; align-items: flex-start; }
.quota-text { font-size: 12px; font-weight: 500; color: #606266; margin-left: 4px; }
.pagination-wrapper { margin-top: 20px; display: flex; justify-content: center; }
.action-buttons { display: flex; gap: 5px; flex-wrap: wrap; align-items: center; }
.action-buttons .el-button { margin: 0; }


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
