<template>
  <div class="performance-monitor">
    <el-card
      class="header-card"
      shadow="never"
    >
      <div class="header-content">
        <div class="title-section">
          <h2>
            <el-icon><Monitor /></el-icon>
            {{ $t('admin.performance.title') }}
          </h2>
          <p class="subtitle">
            {{ $t('admin.performance.subtitle') }}
          </p>
        </div>
      </div>
    </el-card>

    <!-- 关键指标卡片 -->
    <el-row
      :gutter="20"
      class="metrics-cards"
    >
      <el-col
        :xs="24"
        :sm="12"
        :md="6"
      >
        <el-card
          shadow="hover"
          class="metric-card"
        >
          <div class="metric-icon goroutine">
            <el-icon><Connection /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">
              {{ $t('admin.performance.goroutineCount') }}
            </div>
            <div class="metric-value">
              {{ metrics.goroutine_count || 0 }}
            </div>
            <div :class="['metric-status', getGoroutineStatus()]">
              {{ getGoroutineStatusText() }}
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col
        :xs="24"
        :sm="12"
        :md="6"
      >
        <el-card
          shadow="hover"
          class="metric-card"
        >
          <div class="metric-icon memory">
            <el-icon><Memo /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">
              {{ $t('admin.performance.memoryUsage') }}
            </div>
            <div class="metric-value">
              {{ metrics.memory_alloc || 0 }} MB
            </div>
            <div :class="['metric-status', getMemoryStatus()]">
              {{ getMemoryStatusText() }}
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col
        :xs="24"
        :sm="12"
        :md="6"
      >
        <el-card
          shadow="hover"
          class="metric-card"
        >
          <div class="metric-icon gc">
            <el-icon><DeleteFilled /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">
              {{ $t('admin.performance.gcCount') }}
            </div>
            <div class="metric-value">
              {{ metrics.gc_count || 0 }}
            </div>
            <div class="metric-status normal">
              {{ $t('admin.performance.averagePause') }}: {{ formatDuration(metrics.gc_pause_avg) }}
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col
        :xs="24"
        :sm="12"
        :md="6"
      >
        <el-card
          shadow="hover"
          class="metric-card"
        >
          <div class="metric-icon database">
            <el-icon><Coin /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">
              {{ $t('admin.performance.databaseConnections') }}
            </div>
            <div class="metric-value">
              {{ dbStats.in_use || 0 }} / {{ dbStats.max_open_connections || 0 }}
            </div>
            <div :class="['metric-status', getDBStatus()]">
              {{ $t('admin.performance.utilization') }}: {{ getDBUtilization() }}%
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 详细信息 -->
    <el-row
      :gutter="20"
      class="detail-section"
    >
      <!-- 内存详情 -->
      <el-col
        :xs="24"
        :md="12"
      >
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>{{ $t('admin.performance.memoryDetails') }}</span>
              <el-tag
                type="info"
                size="small"
              >
                {{ $t('admin.performance.unit') }}: MB
              </el-tag>
            </div>
          </template>
          <el-descriptions
            :column="2"
            border
          >
            <el-descriptions-item :label="$t('admin.performance.currentAlloc')">
              {{ metrics.memory_alloc || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.totalAlloc')">
              {{ metrics.memory_total_alloc || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.systemMemory')">
              {{ metrics.memory_sys || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.heapMemory')">
              {{ metrics.memory_heap_alloc || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.heapSystem')">
              {{ metrics.memory_heap_sys || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.stackUsage')">
              {{ metrics.memory_stack_inuse || 0 }}
            </el-descriptions-item>
          </el-descriptions>
          
          <!-- 内存使用趋势图 -->
          <div
            ref="memoryChartRef"
            class="chart-container"
          />
        </el-card>
      </el-col>

      <!-- GC 详情 -->
      <el-col
        :xs="24"
        :md="12"
      >
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>{{ $t('admin.performance.gcDetails') }}</span>
            </div>
          </template>
          <el-descriptions
            :column="2"
            border
          >
            <el-descriptions-item :label="$t('admin.performance.gcCount')">
              {{ metrics.gc_count || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.totalPauseTime')">
              {{ formatDuration(metrics.gc_pause_total) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.averagePause')">
              {{ formatDuration(metrics.gc_pause_avg) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.lastPause')">
              {{ formatDuration(metrics.gc_last_pause) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.nextGC')">
              {{ metrics.next_gc || 0 }} MB
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.cpuCores')">
              {{ metrics.cpu_count || 0 }}
            </el-descriptions-item>
          </el-descriptions>

          <!-- GC 频率趋势图 -->
          <div
            ref="gcChartRef"
            class="chart-container"
          />
        </el-card>
      </el-col>
    </el-row>

    <!-- 数据库和连接池 -->
    <el-row
      :gutter="20"
      class="detail-section"
    >
      <el-col
        :xs="24"
        :md="12"
      >
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>{{ $t('admin.performance.databasePool') }}</span>
              <el-tag
                v-if="dbManagerStats && dbManagerStats.connected" 
                :type="dbManagerStats.reconnecting ? 'warning' : 'success'" 
                size="small"
              >
                {{ dbManagerStats.reconnecting ? $t('admin.performance.reconnecting') : $t('admin.performance.connected') }}
              </el-tag>
            </div>
          </template>
          
          <!-- 数据库管理器状态 -->
          <el-alert 
            v-if="dbManagerStats && dbManagerStats.heartbeat_active"
            :title="$t('admin.performance.dbManagerStatus')"
            type="success"
            :closable="false"
            show-icon
            style="margin-bottom: 16px"
          >
            <template #default>
              <div style="font-size: 12px; line-height: 1.8;">
                <div>{{ $t('admin.performance.heartbeatActive') }}: {{ dbManagerStats.heartbeat_active ? $t('common.yes') : $t('common.no') }}</div>
                <div>{{ $t('admin.performance.maxReconnectRetry') }}: {{ dbManagerStats.max_reconnect_retry }}</div>
                <div>{{ $t('admin.performance.reconnectInterval') }}: {{ dbManagerStats.reconnect_interval }}</div>
              </div>
            </template>
          </el-alert>

          <el-descriptions
            v-if="dbStats"
            :column="2"
            border
          >
            <el-descriptions-item :label="$t('admin.performance.maxConnections')">
              {{ dbStats.max_open_connections || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.currentConnections')">
              {{ dbStats.open_connections || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.inUse')">
              {{ dbStats.in_use || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.idle')">
              {{ dbStats.idle || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.waitCount')">
              {{ dbStats.wait_count || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.waitDuration')">
              {{ formatDuration(dbStats.wait_duration) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.maxIdleClosed')">
              {{ dbStats.max_idle_closed || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.maxLifetimeClosed')">
              {{ dbStats.max_lifetime_closed || 0 }}
            </el-descriptions-item>
          </el-descriptions>
          <el-empty
            v-else
            :description="$t('admin.performance.noData')"
          />
        </el-card>
      </el-col>

      <el-col
        :xs="24"
        :md="12"
      >
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>{{ $t('admin.performance.sshPool') }}</span>
              <el-tag
                v-if="sshPoolStats && sshPoolStats.utilization !== undefined" 
                :type="getSSHPoolUtilizationType(sshPoolStats.utilization)" 
                size="small"
              >
                {{ $t('admin.performance.utilization') }}: {{ sshPoolStats.utilization?.toFixed(1) || 0 }}%
              </el-tag>
            </div>
          </template>
          <el-descriptions
            v-if="sshPoolStats && sshPoolStats.total_connections !== undefined"
            :column="2"
            border
          >
            <el-descriptions-item :label="$t('admin.performance.totalConnections')">
              {{ sshPoolStats.total_connections || 0 }} / {{ sshPoolStats.max_connections || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.healthyConnections')">
              <el-tag
                :type="sshPoolStats.healthy_connections === sshPoolStats.total_connections ? 'success' : 'warning'"
                size="small"
              >
                {{ sshPoolStats.healthy_connections || 0 }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.unhealthyConnections')">
              <el-tag
                :type="sshPoolStats.unhealthy_connections > 0 ? 'danger' : 'info'"
                size="small"
              >
                {{ sshPoolStats.unhealthy_connections || 0 }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.activeConnections')">
              <el-tag
                type="success"
                size="small"
              >
                {{ sshPoolStats.active_connections || 0 }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.idleConnections')">
              <el-tag
                type="info"
                size="small"
              >
                {{ sshPoolStats.idle_connections || 0 }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.avgConnectionAge')">
              {{ formatDuration(sshPoolStats.avg_connection_age) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.oldestConnectionAge')">
              {{ formatDuration(sshPoolStats.oldest_connection_age) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.performance.maxIdleTime')">
              {{ formatDuration(sshPoolStats.max_idle_time) }}
            </el-descriptions-item>
          </el-descriptions>
          <el-empty
            v-else
            :description="$t('admin.performance.noSSHData')"
          />
        </el-card>
      </el-col>
    </el-row>

    <!-- Goroutine 趋势图 -->
    <el-card
      shadow="hover"
      class="chart-card"
    >
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.performance.goroutineTrend') }}</span>
          <el-radio-group
            v-model="timeRange"
            size="small"
            @change="fetchHistory"
          >
            <el-radio-button label="5m">
              {{ $t('admin.performance.timeRange.5m') }}
            </el-radio-button>
            <el-radio-button label="15m">
              {{ $t('admin.performance.timeRange.15m') }}
            </el-radio-button>
            <el-radio-button label="1h">
              {{ $t('admin.performance.timeRange.1h') }}
            </el-radio-button>
            <el-radio-button label="6h">
              {{ $t('admin.performance.timeRange.6h') }}
            </el-radio-button>
          </el-radio-group>
        </div>
      </template>
      <div
        ref="goroutineChartRef"
        class="chart-container-large"
      />
    </el-card>
  </div>
</template>

<script setup>
import { 
  Monitor, 
  Refresh, 
  Connection, 
  Memo, 
  DeleteFilled,
  Coin 
} from '@element-plus/icons-vue'
import { usePerformanceMetrics } from './composables/usePerformanceMetrics'

const {
  t,
  loading,
  timeRange,
  metrics,
  dbStats,
  dbManagerStats,
  sshPoolStats,
  memoryChartRef,
  gcChartRef,
  goroutineChartRef,
  fetchMetrics,
  fetchHistory,
  formatDuration,
  getGoroutineStatus,
  getGoroutineStatusText,
  getMemoryStatus,
  getMemoryStatusText,
  getDBStatus,
  getDBUtilization,
  getSSHPoolUtilizationType,
} = usePerformanceMetrics()
</script>
<script>
export default {
  name: 'PerformanceMonitor'
}
</script>

<style scoped lang="scss">
.performance-monitor {
  padding: 20px;

  .header-card {
    margin-bottom: 20px;
    
    .header-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .title-section {
        h2 {
          margin: 0;
          font-size: 24px;
          font-weight: 600;
          display: flex;
          align-items: center;
          gap: 10px;
        }
        
        .subtitle {
          margin: 5px 0 0;
          color: var(--text-color-secondary);
          font-size: 14px;
        }
      }
    }
  }

  .metrics-cards {
    margin-bottom: 20px;
    
    .metric-card {
      margin-bottom: 20px;
      cursor: pointer;
      transition: all 0.3s;
      
      &:hover {
        transform: translateY(-5px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }
      
      :deep(.el-card__body) {
        display: flex;
        align-items: center;
        padding: 20px;
      }
      
      .metric-icon {
        width: 60px;
        height: 60px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 28px;
        color: white;
        margin-right: 15px;
        
        &.goroutine { background: linear-gradient(135deg, #16a34a 0%, #22c55e 100%); }
        &.memory { background: linear-gradient(135deg, #059669 0%, #10b981 100%); }
        &.gc { background: linear-gradient(135deg, #0f766e 0%, #14b8a6 100%); }
        &.database { background: linear-gradient(135deg, #65a30d 0%, #84cc16 100%); }
      }
      
      .metric-content {
        flex: 1;
        
        .metric-label {
          font-size: 14px;
          color: var(--text-color-secondary);
          margin-bottom: 5px;
        }
        
        .metric-value {
          font-size: 28px;
          font-weight: 600;
          margin-bottom: 5px;
        }
        
        .metric-status {
          font-size: 12px;
          padding: 2px 8px;
          border-radius: 4px;
          display: inline-block;
          
          &.normal {
            background: var(--success-bg);
            color: #16a34a;
          }
          
          &.warning {
            background: var(--warning-bg);
            color: #e6a23c;
          }
          
          &.critical {
            background: var(--error-bg);
            color: #f56c6c;
          }
        }
      }
    }
  }

  .detail-section {
    margin-bottom: 20px;
  }

  .chart-card {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
  }

  .chart-container {
    height: 300px;
    margin-top: 20px;
  }

  .chart-container-large {
    height: 400px;
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
