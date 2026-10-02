<template>
  <!-- Agent 说明 -->
  <el-alert
    :title="$t('admin.providers.agentMonitoringDescTitle')"
    type="success"
    :closable="false"
    show-icon
    style="margin-bottom: 16px;"
  >
    <template #default>
      <p style="margin: 4px 0 0;">
        {{ $t('admin.providers.agentMonitoringDesc') }}
      </p>
    </template>
  </el-alert>

  <!-- Agent 状态 -->
  <div class="agent-status-section">
    <el-descriptions
      :column="2"
      border
      size="small"
    >
      <el-descriptions-item :label="$t('admin.providers.monitoringMode')">
        <el-tag
          :type="config.monitoring_mode === 'agent' ? 'success' : 'info'"
          size="small"
        >
          {{ config.monitoring_mode === 'agent' ? 'Agent' : 'PMAcct' }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item :label="$t('admin.providers.agentStatus')">
        <el-tag
          :type="agentStatusType"
          size="small"
        >
          {{ agentStatusText }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item
        v-if="config.agent_version"
        :label="$t('admin.providers.agentVersion')"
      >
        {{ config.agent_version }}
      </el-descriptions-item>
      <el-descriptions-item :label="$t('admin.providers.agentPort')">
        {{ config.agent_port || 23782 }}
      </el-descriptions-item>
      <el-descriptions-item :label="$t('admin.providers.collectInterval')">
        {{ config.collect_interval || 5 }}s
      </el-descriptions-item>
      <el-descriptions-item :label="$t('admin.providers.resourceCollectInterval')">
        {{ config.resource_collect_interval || 30 }}s
      </el-descriptions-item>
    </el-descriptions>

    <!-- Agent Token 展示与复制 -->
    <div
      v-if="config.agent_token"
      class="token-section"
    >
      <el-descriptions
        :column="1"
        border
        size="small"
        style="margin-top: 12px;"
      >
        <el-descriptions-item :label="$t('admin.providers.agentToken')">
          <div style="display: flex; align-items: center; gap: 8px;">
            <el-input
              :model-value="showToken ? config.agent_token : '••••••••••••••••'"
              readonly
              size="small"
              style="flex: 1; max-width: 320px;"
            />
            <el-button
              size="small"
              @click="$emit('toggle-show-token')"
            >
              {{ showToken ? $t('admin.providers.hideToken') : $t('admin.providers.showToken') }}
            </el-button>
            <el-button
              size="small"
              type="primary"
              @click="$emit('copy-token')"
            >
              {{ $t('admin.providers.copyToken') }}
            </el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item
          v-if="!isAgentProvider"
          :label="$t('admin.providers.agentTestUrl')"
        >
          <div>
            <div style="display: flex; align-items: center; gap: 8px;">
              <el-input
                :model-value="agentSwaggerUrl"
                readonly
                size="small"
                style="flex: 1; max-width: 400px;"
              />
              <el-button
                size="small"
                @click="$emit('copy-url', agentSwaggerUrl)"
              >
                {{ $t('admin.providers.copyUrl') }}
              </el-button>
            </div>
            <div
              v-if="config.agent_version"
              style="margin-top: 6px;"
            >
              <el-text
                size="small"
                type="info"
              >
                {{ $t('admin.providers.agentVersion') }}: {{ config.agent_version }}
              </el-text>
            </div>
          </div>
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <!-- 操作按钮 -->
    <div class="action-buttons">
      <el-button
        v-if="!isAgentProvider"
        type="success"
        :loading="deployLoading"
        @click="$emit('deploy-agent')"
      >
        {{ config.agent_installed ? $t('admin.providers.redeployAgent') : $t('admin.providers.deployAgent') }}
      </el-button>
      <el-button
        v-if="!isAgentProvider"
        type="danger"
        :loading="uninstallLoading"
        :disabled="!config.agent_installed"
        @click="$emit('uninstall-agent')"
      >
        {{ $t('admin.providers.uninstallAgent') }}
      </el-button>
      <el-button
        type="primary"
        :loading="statusLoading"
        @click="$emit('check-status')"
      >
        {{ $t('admin.providers.checkAgentStatus') }}
      </el-button>
      <el-button
        type="warning"
        :loading="syncLoading"
        :disabled="!config.agent_installed"
        @click="$emit('sync-monitors')"
      >
        {{ $t('admin.providers.syncMonitors') }}
      </el-button>
      <el-button
        type="danger"
        :loading="clearMonitorsLoading"
        :disabled="!config.agent_installed"
        @click="$emit('clear-monitors')"
      >
        {{ $t('admin.providers.clearMonitors') }}
      </el-button>
      <el-button
        @click="$emit('toggle-config-editor')"
      >
        {{ $t('admin.providers.editConfig') }}
      </el-button>
    </div>

    <!-- 配置编辑器 -->
    <el-card
      v-if="showConfigEditor"
      shadow="never"
      style="margin-top: 16px;"
    >
      <template #header>
        <span>{{ $t('admin.providers.monitoringConfig') }}</span>
      </template>
      <el-form
        :model="editConfig"
        label-width="180px"
        size="small"
      >
        <el-form-item :label="$t('admin.providers.monitoringMode')">
          <el-select
            v-model="editConfig.monitoring_mode"
            style="width: 160px;"
          >
            <el-option
              label="Agent"
              value="agent"
            />
            <el-option
              label="PMAcct"
              value="pmacct"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('admin.providers.trafficCollectMethod')">
          <el-select
            v-model="editConfig.traffic_collect_method"
            style="width: 160px;"
          >
            <el-option
              label="nftables (NFT)"
              value="nft"
            />
            <el-option
              label="iptables (IPT)"
              value="ipt"
            />
          </el-select>
          <el-text
            type="info"
            size="small"
            style="margin-left: 8px;"
          >
            {{ $t('admin.providers.trafficCollectMethodHint') }}
          </el-text>
        </el-form-item>
        <el-form-item :label="$t('admin.providers.agentPort')">
          <el-input-number
            v-model="editConfig.agent_port"
            :min="1024"
            :max="65535"
          />
        </el-form-item>
        <el-form-item :label="$t('admin.providers.collectInterval')">
          <el-input-number
            v-model="editConfig.collect_interval"
            :min="1"
            :max="300"
          />
          <span style="margin-left: 8px; color: #909399;">s</span>
          <el-text
            type="info"
            size="small"
            style="margin-left: 8px;"
          >
            {{ $t('admin.providers.collectIntervalHint') }}
          </el-text>
        </el-form-item>
        <el-form-item :label="$t('admin.providers.resourceCollectInterval')">
          <el-input-number
            v-model="editConfig.resource_collect_interval"
            :min="10"
            :max="3600"
          />
          <span style="margin-left: 8px; color: #909399;">s</span>
          <el-text
            type="info"
            size="small"
            style="margin-left: 8px;"
          >
            {{ $t('admin.providers.resourceCollectIntervalHint') }}
          </el-text>
        </el-form-item>
        <el-form-item :label="$t('admin.providers.extraExcludeCIDRsV4')">
          <el-input
            v-model="editConfig.extra_exclude_cidrs_v4"
            type="textarea"
            :rows="2"
            :placeholder="$t('admin.providers.extraExcludeCIDRsPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="$t('admin.providers.extraExcludeCIDRsV6')">
          <el-input
            v-model="editConfig.extra_exclude_cidrs_v6"
            type="textarea"
            :rows="2"
            :placeholder="$t('admin.providers.extraExcludeCIDRsV6Placeholder')"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="saveConfigLoading"
            @click="$emit('save-config')"
          >
            {{ $t('common.save') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 部署输出 -->
    <div
      v-if="deployOutput"
      class="deploy-output"
    >
      <h4>{{ $t('admin.providers.deployOutput') }}</h4>
      <div class="output-content">
        <pre>{{ deployOutput }}</pre>
      </div>
    </div>
  </div>

  <!-- 监控列表 -->
  <div style="margin-top: 20px;">
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
      <h4 style="margin: 0;">
        {{ $t('admin.providers.instanceMonitors') }} ({{ monitorsPagination.total }})
      </h4>
      <el-button
        v-if="config.agent_installed && agentIsOnline"
        size="small"
        :loading="listAgentLoading"
        @click="$emit('list-agent-monitors')"
      >
        {{ $t('admin.providers.viewAgentMonitors') }}
      </el-button>
    </div>
    <el-table
      v-loading="monitorsLoading"
      :data="monitors"
      size="small"
      max-height="300"
    >
      <el-table-column
        prop="instance_name"
        :label="$t('admin.providers.instanceName')"
        width="150"
      >
        <template #default="{ row }">
          <span>{{ row.instance_name || '-' }}</span>
          <el-tag
            v-if="row.instance_deleted"
            type="danger"
            size="small"
            style="margin-left: 4px;"
          >
            {{ $t('admin.providers.deleted') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        prop="interfaces"
        :label="$t('admin.providers.interfaces')"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          {{ row.interfaces || '-' }}
        </template>
      </el-table-column>
      <el-table-column
        prop="agent_monitor_id"
        :label="$t('admin.providers.agentId')"
        min-width="160"
      />
      <el-table-column
        :label="$t('admin.providers.trafficIn')"
        min-width="120"
      >
        <template #default="{ row }">
          {{ formatBytes(row.last_traffic_bytes_in || 0) }}
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.trafficOut')"
        min-width="130"
      >
        <template #default="{ row }">
          {{ formatBytes(row.last_traffic_bytes_out || 0) }}
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.status')"
        min-width="90"
      >
        <template #default="{ row }">
          <el-tag
            :type="row.is_enabled ? 'success' : 'info'"
            size="small"
          >
            {{ row.is_enabled ? $t('common.enabled') : $t('common.disabled') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.monitorHealth')"
        min-width="120"
      >
        <template #default="{ row }">
          <el-tooltip
            :disabled="!row.health_error"
            :content="row.health_error"
            placement="top"
          >
            <el-tag
              :type="healthTagType(row.health_status)"
              size="small"
            >
              {{ healthStatusText(row.health_status) }}
            </el-tag>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.lastSync')"
        width="160"
      >
        <template #default="{ row }">
          {{ row.last_sync_at ? formatDateTime(row.last_sync_at) : '-' }}
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-if="monitorsPagination.total > monitorsPagination.pageSize"
      v-model:current-page="monitorsPagination.page"
      v-model:page-size="monitorsPagination.pageSize"
      :page-sizes="[10, 20, 50]"
      :total="monitorsPagination.total"
      layout="total, sizes, prev, pager, next"
      size="small"
      style="margin-top: 8px; justify-content: center;"
      @current-change="$emit('load-monitors')"
      @size-change="() => { monitorsPagination.page = 1; $emit('load-monitors') }"
    />
  </div>

  <!-- Agent端监控列表弹窗 -->
  <el-dialog
    v-model="localShowAgentMonitors"
    :title="$t('admin.providers.agentMonitorsList')"
    width="1000px"
    append-to-body
  >
    <el-table
      :data="agentMonitors"
      size="small"
      max-height="400"
    >
      <el-table-column
        prop="id"
        label="ID"
        width="70"
      />
      <el-table-column
        :label="$t('admin.providers.interfaces')"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          {{ (row.interface || []).join(', ') || '-' }}
        </template>
      </el-table-column>
      <el-table-column
        prop="instance_name"
        :label="$t('admin.providers.instanceName')"
        width="150"
      >
        <template #default="{ row }">
          <span>{{ row.instance_name || '-' }}</span>
          <el-tag
            v-if="row.instance_deleted"
            type="danger"
            size="small"
            style="margin-left: 4px;"
          >
            {{ $t('admin.providers.deleted') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        prop="provider_kind"
        :label="$t('admin.providers.provider')"
        min-width="160"
      >
        <template #default="{ row }">
          {{ row.provider_kind || '-' }}
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.monitorHealth')"
        min-width="130"
      >
        <template #default="{ row }">
          <el-tooltip
            :disabled="!monitorHealthDetail(row)"
            :content="monitorHealthDetail(row)"
            placement="top"
          >
            <el-tag
              :type="healthTagType(row.health_status)"
              size="small"
            >
              {{ healthStatusText(row.health_status) }}
            </el-tag>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.trafficIn')"
        min-width="120"
      >
        <template #default="{ row }">
          {{ formatBytes(row.total_bytes_in || 0) }}
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.trafficOut')"
        min-width="130"
      >
        <template #default="{ row }">
          {{ formatBytes(row.total_bytes_out || 0) }}
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('admin.providers.totalTraffic')"
        min-width="150"
      >
        <template #default="{ row }">
          {{ formatBytes(row.total_bytes || 0) }}
        </template>
      </el-table-column>
    </el-table>
    <template #footer>
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <el-text
          type="info"
          size="small"
        >
          {{ $t('admin.providers.agentMonitorsTotal') }}: {{ agentMonitorsPagination.total }}
        </el-text>
        <el-pagination
          v-if="agentMonitorsPagination.total > agentMonitorsPagination.pageSize"
          v-model:current-page="agentMonitorsPagination.page"
          v-model:page-size="agentMonitorsPagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="agentMonitorsPagination.total"
          layout="total, sizes, prev, pager, next"
          size="small"
          @current-change="$emit('list-agent-monitors')"
          @size-change="() => { agentMonitorsPagination.page = 1; $emit('list-agent-monitors') }"
        />
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

const props = defineProps({
  config: { type: Object, required: true },
  editConfig: { type: Object, required: true },
  showConfigEditor: { type: Boolean, default: false },
  showToken: { type: Boolean, default: false },
  deployLoading: { type: Boolean, default: false },
  uninstallLoading: { type: Boolean, default: false },
  statusLoading: { type: Boolean, default: false },
  syncLoading: { type: Boolean, default: false },
  clearMonitorsLoading: { type: Boolean, default: false },
  saveConfigLoading: { type: Boolean, default: false },
  listAgentLoading: { type: Boolean, default: false },
  monitorsLoading: { type: Boolean, default: false },
  deployOutput: { type: String, default: '' },
  monitors: { type: Array, default: () => [] },
  agentIsOnline: { type: Boolean, default: false },
  showAgentMonitors: { type: Boolean, default: false },
  agentMonitors: { type: Array, default: () => [] },
  monitorsPagination: { type: Object, required: true },
  agentMonitorsPagination: { type: Object, required: true },
  agentSwaggerUrl: { type: String, default: '' },
  agentStatusType: { type: String, default: 'info' },
  agentStatusText: { type: String, default: '' },
  isAgentProvider: { type: Boolean, default: false }
})

const emit = defineEmits([
  'toggle-show-token', 'copy-token', 'copy-url',
  'deploy-agent', 'uninstall-agent', 'check-status',
  'sync-monitors', 'clear-monitors', 'list-agent-monitors',
  'save-config', 'toggle-config-editor', 'load-monitors',
  'update:showAgentMonitors'
])

const localShowAgentMonitors = computed({
  get: () => props.showAgentMonitors,
  set: (val) => emit('update:showAgentMonitors', val)
})

function formatDateTime(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

function formatBytes(bytes) {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function healthTagType(status) {
  if (status === 'healthy') return 'success'
  if (status === 'unhealthy') return 'danger'
  return 'warning'
}

function healthStatusText(status) {
  if (status === 'healthy') return t('admin.providers.monitorHealthy')
  if (status === 'unhealthy') return t('admin.providers.monitorUnhealthy')
  return t('admin.providers.monitorUnknown')
}

function monitorHealthDetail(row) {
  if (row.health_error) return row.health_error
  if (Array.isArray(row.missing_interfaces) && row.missing_interfaces.length > 0) {
    return `${t('admin.providers.missingInterfaces')}: ${row.missing_interfaces.join(', ')}`
  }
  return ''
}
</script>

<style scoped>
.agent-status-section { padding: 0 4px; }
.token-section { margin-top: 12px; }

.action-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 16px;
}

.deploy-output { margin-top: 20px; }
.deploy-output h4 { margin: 0 0 12px; font-size: 14px; font-weight: 600; }

.output-content {
  background: var(--neutral-bg);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 12px;
  max-height: 400px;
  overflow-y: auto;
}

.output-content pre {
  margin: 0;
  font-family: 'Courier New', Courier, monospace;
  font-size: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-wrap: break-word;
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
