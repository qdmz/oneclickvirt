<template>
  <div class="admin-traffic">
    <div class="page-header">
      <h1>{{ $t('admin.traffic.title') }}</h1>
      <p>{{ $t('admin.traffic.subtitle') }}</p>
    </div>

    <!-- 系统流量概览 -->
    <div class="system-overview">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>{{ $t('admin.traffic.systemOverview') }}</span>
            <div class="header-actions">
              <el-button
                :loading="overviewLoading"
                @click="loadSystemOverview"
              >
                {{ $t('common.refresh') }}
              </el-button>
              <el-button
                v-if="isSuperAdmin"
                type="primary"
                :loading="syncingAllTraffic"
                @click="syncAllTrafficData"
              >
                {{ $t('admin.traffic.syncAllTraffic') }}
              </el-button>
            </div>
          </div>
        </template>

        <div
          v-if="overviewLoading"
          class="loading-container"
        >
          <el-skeleton
            :rows="3"
            animated
          />
        </div>

        <div
          v-else-if="systemOverview"
          class="overview-content"
        >
          <el-row :gutter="20">
            <el-col :span="6">
              <div class="stat-card">
                <div class="stat-title">
                  {{ $t('admin.traffic.monthlyTotalTraffic') }}
                </div>
                <div class="stat-value">
                  {{ systemOverview.traffic?.formatted?.total_bytes || '0 B' }}
                </div>
                <div class="stat-subtitle">
                  {{ $t('admin.traffic.uplink') }}: {{ systemOverview.traffic?.formatted?.total_tx || '0 B' }} / 
                  {{ $t('admin.traffic.downlink') }}: {{ systemOverview.traffic?.formatted?.total_rx || '0 B' }}
                </div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-card">
                <div class="stat-title">
                  {{ $t('admin.traffic.userStats') }}
                </div>
                <div class="stat-value">
                  {{ systemOverview.users?.total || 0 }}
                </div>
                <div class="stat-subtitle">
                  {{ $t('admin.traffic.limited') }}: {{ systemOverview.users?.limited || 0 }} 
                  ({{ (systemOverview.users?.limited_percent || 0).toFixed(1) }}%)
                </div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-card">
                <div class="stat-title">
                  {{ $t('admin.traffic.providerStats') }}
                </div>
                <div class="stat-value">
                  {{ systemOverview.providers?.total || 0 }}
                </div>
                <div class="stat-subtitle">
                  {{ $t('admin.traffic.limited') }}: {{ systemOverview.providers?.limited || 0 }} 
                  ({{ (systemOverview.providers?.limited_percent || 0).toFixed(1) }}%)
                </div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-card">
                <div class="stat-title">
                  {{ $t('admin.traffic.totalInstances') }}
                </div>
                <div class="stat-value">
                  {{ systemOverview.instances || 0 }}
                </div>
                <div class="stat-subtitle">
                  {{ $t('admin.traffic.activeInstanceStats') }}
                </div>
              </div>
            </el-col>
          </el-row>

          <div class="period-info">
            <el-text
              type="info"
              size="small"
            >
              <el-icon><Calendar /></el-icon>
              {{ $t('admin.traffic.statsPeriod') }}:
              {{ systemOverview.period_type === 'current_cycle' ? $t('admin.traffic.currentTrafficCycle') : systemOverview.period }}
              <span v-if="systemOverview.period">({{ systemOverview.period }})</span>
            </el-text>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 流量排行榜 -->
    <div class="traffic-ranking">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>{{ $t('admin.traffic.trafficRanking') }}</span>
          </div>
        </template>

        <!-- 搜索和批量操作工具栏 -->
        <div class="toolbar">
          <div class="search-section">
            <el-input
              v-model="searchParams.username"
              :placeholder="$t('admin.traffic.searchByUsername')"
              style="width: 200px;"
              clearable
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-input
              v-model="searchParams.nickname"
              :placeholder="$t('admin.traffic.searchByNickname')"
              style="width: 200px; margin-left: 10px;"
              clearable
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button 
              type="primary" 
              style="margin-left: 10px;"
              @click="handleSearch"
            >
              {{ $t('common.search') }}
            </el-button>
            <el-button 
              @click="resetSearch"
            >
              {{ $t('common.reset') }}
            </el-button>
            <el-button
              size="default"
              :loading="rankingLoading"
              @click="loadTrafficRanking"
            >
              <el-icon><Refresh /></el-icon>
              {{ $t('common.refresh') }}
            </el-button>
          </div>

          <!-- 批量操作 -->
          <div
            v-if="isSuperAdmin && selectedUsers.length > 0"
            class="batch-actions"
          >
            <span class="selection-info">
              {{ $t('admin.traffic.selected') }} {{ selectedUsers.length }} {{ $t('admin.traffic.users') }}
            </span>
            <el-button
              size="small"
              type="primary"
              @click="handleBatchSync"
            >
              {{ $t('admin.traffic.batchSync') }}
            </el-button>
            <el-button
              size="small"
              type="warning"
              @click="handleBatchLimit"
            >
              {{ $t('admin.traffic.batchLimit') }}
            </el-button>
            <el-button
              size="small"
              type="success"
              @click="handleBatchUnlimit"
            >
              {{ $t('admin.traffic.batchUnlimit') }}
            </el-button>
          </div>
        </div>

        <div
          v-if="rankingLoading"
          class="loading-container"
        >
          <el-skeleton
            :rows="5"
            animated
          />
        </div>

        <div v-else-if="trafficRanking && trafficRanking.length > 0">
          <el-table
            :data="trafficRanking"
            stripe
            border
            @selection-change="handleSelectionChange"
          >
            <el-table-column
              type="selection"
              width="55"
              align="center"
            />
            <el-table-column
              :label="$t('admin.traffic.rank')"
              width="80"
              align="center"
            >
              <template #default="{ row }">
                <el-tag 
                  :type="getRankTagType(row.rank)"
                  effect="dark"
                  size="small"
                >
                  #{{ row.rank }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              prop="username"
              :label="$t('admin.traffic.username')"
              width="150"
            />
            <el-table-column
              prop="nickname"
              :label="$t('admin.traffic.nickname')"
              width="150"
            />
            <el-table-column
              :label="$t('admin.traffic.monthlyUsage')"
              width="120"
            >
              <template #default="{ row }">
                {{ row.formatted?.month_usage || formatTrafficMB(row.month_usage) }}
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('admin.traffic.totalLimit')"
              min-width="130"
            >
              <template #default="{ row }">
                {{ row.formatted?.total_limit || formatTrafficMB(row.total_limit) }}
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('admin.traffic.usageRate')"
              width="120"
              align="center"
            >
              <template #default="{ row }">
                <el-progress
                  :percentage="Math.min(row.usage_percent || 0, 100)"
                  :color="getUsageColor(row.usage_percent || 0)"
                  :stroke-width="8"
                  :show-text="false"
                />
                <div style="margin-top: 4px; font-size: 12px;">
                  {{ (row.usage_percent || 0).toFixed(1) }}%
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('common.status')"
              width="100"
              align="center"
            >
              <template #default="{ row }">
                <el-tag 
                  :type="row.is_limited ? 'danger' : 'success'"
                  size="small"
                >
                  {{ row.is_limited ? $t('admin.traffic.limitedStatus') : $t('common.normal') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              v-if="isSuperAdmin"
              :label="$t('common.actions')"
              width="450"
              align="center"
            >
              <template #default="{ row }">
                <el-button
                  size="small"
                  @click="viewUserTraffic(row.user_id)"
                >
                  {{ $t('admin.traffic.viewDetails') }}
                </el-button>
                <el-button
                  size="small"
                  type="primary"
                  :loading="syncingUsers.includes(row.user_id)"
                  @click="syncUserTrafficData(row.user_id)"
                >
                  {{ $t('admin.traffic.syncTraffic') }}
                </el-button>
                <el-button
                  v-if="!row.is_limited"
                  size="small"
                  type="warning"
                  @click="limitUser(row)"
                >
                  {{ $t('admin.traffic.limitTraffic') }}
                </el-button>
                <el-button
                  v-else
                  size="small"
                  type="success"
                  @click="unlimitUser(row)"
                >
                  {{ $t('admin.traffic.removeLimit') }}
                </el-button>
                <el-button
                  size="small"
                  type="danger"
                  @click="clearUserTraffic(row)"
                >
                  {{ $t('admin.traffic.clearTraffic') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <!-- 分页 -->
          <div class="pagination-wrapper">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>

        <div
          v-else
          class="empty-state"
        >
          <el-empty :description="$t('admin.traffic.noTrafficData')" />
        </div>
      </el-card>
    </div>

    <!-- 用户流量详情对话框 -->
    <el-dialog
      v-model="userTrafficDialogVisible"
      :title="$t('admin.traffic.userTrafficDetails')"
      width="600px"
    >
      <div
        v-if="userTrafficLoading"
        class="loading-container"
      >
        <el-skeleton
          :rows="4"
          animated
        />
      </div>

      <div
        v-else-if="selectedUserTraffic"
        class="user-traffic-detail"
      >
        <el-descriptions
          :column="2"
          border
        >
          <el-descriptions-item :label="$t('admin.traffic.userId')">
            {{ selectedUserTraffic.user_id }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.traffic.dataSource')">
            <el-tag :type="dataSourceTagType(selectedUserTraffic.data_source)">
              {{ dataSourceLabel(selectedUserTraffic.data_source) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.traffic.currentCycleUsage')">
            {{ selectedUserTraffic.formatted?.current_usage || formatTrafficMB(selectedUserTraffic.current_month_usage) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.traffic.totalLimit')">
            {{ selectedUserTraffic.formatted?.total_limit || formatTrafficMB(selectedUserTraffic.total_limit) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.traffic.usageRate')">
            {{ (selectedUserTraffic.usage_percent || 0).toFixed(2) }}%
          </el-descriptions-item>
          <el-descriptions-item :label="$t('common.status')">
            <el-tag :type="selectedUserTraffic.is_limited ? 'danger' : 'success'">
              {{ selectedUserTraffic.is_limited ? $t('admin.traffic.limitedStatus') : $t('common.normal') }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <div
          v-if="selectedUserTraffic.reset_time"
          style="margin-top: 15px;"
        >
          <el-text
            type="info"
            size="small"
          >
            <el-icon><Clock /></el-icon>
            {{ $t('admin.traffic.trafficResetTime') }}: {{ formatDate(selectedUserTraffic.reset_time) }}
          </el-text>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button 
            type="primary"
            :loading="syncingUserDetail"
            @click="syncUserTrafficFromDetail"
          >
            {{ $t('admin.traffic.syncNow') }}
          </el-button>
          <el-button @click="userTrafficDialogVisible = false">{{ $t('common.close') }}</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 流量限制对话框 -->
    <el-dialog
      v-model="limitDialogVisible"
      :title="limitAction === 'limit' ? $t('admin.traffic.limitUserTraffic') : $t('admin.traffic.removeLimitTitle')"
      width="400px"
    >
      <el-form
        ref="limitFormRef"
        :model="limitForm"
        :rules="limitFormRules"
        label-width="80px"
      >
        <el-form-item :label="$t('common.user')">
          <el-text>{{ selectedUser?.username }} ({{ selectedUser?.email }})</el-text>
        </el-form-item>
        <el-form-item
          v-if="limitAction === 'limit'"
          :label="$t('admin.traffic.limitReason')"
          prop="reason"
        >
          <el-input
            v-model="limitForm.reason"
            type="textarea"
            :rows="3"
            :placeholder="$t('admin.traffic.enterLimitReason')"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="limitDialogVisible = false">{{ $t('common.cancel') }}</el-button>
          <el-button
            type="primary"
            :loading="limitSubmitting"
            @click="submitLimitAction"
          >
            {{ $t('common.confirm') }}{{ limitAction === 'limit' ? $t('admin.traffic.limit') : $t('admin.traffic.remove') }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { Refresh, Calendar, Clock, Search } from '@element-plus/icons-vue'
import { useTrafficManagement } from './composables/useTrafficManagement'

const {
  overviewLoading, systemOverview, syncingAllTraffic, isSuperAdmin,
  rankingLoading, trafficRanking, currentPage, pageSize, total, selectedUsers,
  searchParams,
  userTrafficDialogVisible, userTrafficLoading, selectedUserTraffic, syncingUserDetail,
  limitDialogVisible, limitSubmitting, limitAction, selectedUser, syncingUsers,
  limitForm, limitFormRules,
  loadSystemOverview, loadTrafficRanking,
  handleSearch, resetSearch, handleSizeChange, handleCurrentChange, handleSelectionChange,
  handleBatchSync, handleBatchLimit, handleBatchUnlimit,
  viewUserTraffic, limitUser, unlimitUser, submitLimitAction,
  syncUserTrafficData, syncUserTrafficFromDetail, syncAllTrafficData,
  clearUserTraffic,
  formatBytes, formatTrafficMB, formatDate, getRankTagType, getUsageColor, dataSourceLabel, dataSourceTagType,
  t
} = useTrafficManagement()

onMounted(() => {
  loadSystemOverview()
  loadTrafficRanking()
})
</script>
<style scoped>
.admin-traffic {
  margin: -24px 0 -24px -24px;
  padding: 24px 0 24px 24px;
  width: auto;
  min-width: 0;
  max-width: 100%;
}

.page-header {
  margin-bottom: 20px;
  padding-right: 24px;
}

.page-header h1 {
  margin: 0 0 8px 0;
  color: var(--el-text-color-primary);
}

.page-header p {
  margin: 0;
  color: var(--el-text-color-regular);
}

.system-overview {
  margin-bottom: 20px;
  padding-right: 24px;
}

.traffic-ranking {
  padding-right: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  > span {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-color-primary);
  }
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.loading-container {
  padding: 20px;
}

.overview-content {
  padding: 10px 0;
}

.stat-card {
  text-align: center;
  padding: 20px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
  border: 1px solid var(--el-border-color-light);
}

.stat-title {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin-bottom: 10px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 8px;
  font-family: monospace;
}

.stat-subtitle {
  font-size: 12px;
  color: var(--el-text-color-regular);
}

.period-info {
  text-align: center;
  margin-top: 20px;
}

.traffic-ranking {
  margin-bottom: 20px;
  padding-right: 0;
}

.traffic-ranking :deep(.el-card) {
  border-radius: 0;
  margin-right: 0;
}

.empty-state {
  padding: 40px;
  text-align: center;
}

.user-traffic-detail {
  padding: 10px 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.toolbar {
  margin-bottom: 16px;
}

.search-section {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 12px;
}

.batch-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px;
  background-color: var(--info-bg);
  border: 1px solid rgba(22, 163, 74, 0.2);
  border-radius: 4px;
}

.selection-info {
  font-size: 14px;
  color: #16a34a;
  font-weight: 500;
  margin-right: 10px;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 1024px) {
  .admin-traffic {
    margin: 0;
    padding: 0;
    width: 100%;
  }

  .page-header,
  .system-overview {
    padding-right: 0;
  }

  .card-header {
    align-items: stretch;
    flex-direction: column;
  }

  .header-actions,
  .search-section,
  .batch-actions {
    justify-content: flex-start;
  }

  .search-section :deep(.el-input),
  .search-section :deep(.el-button),
  .header-actions .el-button,
  .batch-actions .el-button {
    margin-left: 0 !important;
  }

  .search-section :deep(.el-input) {
    flex: 1 1 180px;
    max-width: 100%;
  }

  .pagination-wrapper {
    justify-content: center;
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
