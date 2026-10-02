<template>
  <div class="admin-tasks">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.tasks.title') }}</span>
          <p class="header-subtitle">
            {{ $t('admin.tasks.subtitle') }}
          </p>
        </div>
      </template>

      <!-- 任务池全局开关 -->
      <el-alert
        class="task-pool-alert"
        :type="poolStatus.enabled ? 'success' : (poolStatus.drainComplete ? 'warning' : 'error')"
        :closable="false"
        show-icon
      >
        <template #title>
          <div class="task-pool-header">
            <div>
              <strong>
                {{ poolStatus.enabled ? $t('admin.tasks.poolEnabled') : $t('admin.tasks.poolDisabled') }}
              </strong>
              <span class="task-pool-desc">
                {{ poolStatus.enabled ? $t('admin.tasks.poolEnabledDesc') : (poolStatus.drainComplete ? $t('admin.tasks.poolMaintenanceReady') : $t('admin.tasks.poolDraining')) }}
              </span>
            </div>
            <el-button
              v-if="isSuperAdmin"
              class="task-pool-action"
              :type="poolStatus.enabled ? 'danger' : 'success'"
              :loading="poolLoading"
              size="small"
              @click="toggleTaskPool(!poolStatus.enabled)"
            >
              {{ poolStatus.enabled ? $t('admin.tasks.disableTaskPool') : $t('admin.tasks.enableTaskPool') }}
            </el-button>
          </div>
        </template>
        <div class="task-pool-meta">
          <span>{{ $t('admin.tasks.pendingTasks') }}: {{ poolStatus.pendingTasks }}</span>
          <span>{{ $t('admin.tasks.runningTasks') }}: {{ poolStatus.runningTasks }}</span>
          <span>{{ $t('admin.tasks.configurationPendingTasks') }}: {{ poolStatus.configurationPendingTasks }}</span>
          <span>{{ $t('admin.tasks.configurationRunningTasks') }}: {{ poolStatus.configurationRunningTasks }}</span>
        </div>
      </el-alert>

      <!-- 统计卡片 -->
      <div class="stats-cards">
        <el-row :gutter="20">
          <el-col :span="4">
            <el-card class="stats-card">
              <div class="stat-item">
                <div class="stat-number">
                  {{ stats.totalTasks }}
                </div>
                <div class="stat-label">
                  {{ $t('admin.tasks.totalTasks') }}
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="4">
            <el-card class="stats-card pending">
              <div class="stat-item">
                <div class="stat-number">
                  {{ stats.pendingTasks }}
                </div>
                <div class="stat-label">
                  {{ $t('admin.tasks.pendingTasks') }}
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="4">
            <el-card class="stats-card running">
              <div class="stat-item">
                <div class="stat-number">
                  {{ stats.runningTasks }}
                </div>
                <div class="stat-label">
                  {{ $t('admin.tasks.runningTasks') }}
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="4">
            <el-card class="stats-card completed">
              <div class="stat-item">
                <div class="stat-number">
                  {{ stats.completedTasks }}
                </div>
                <div class="stat-label">
                  {{ $t('admin.tasks.completedTasks') }}
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="4">
            <el-card class="stats-card failed">
              <div class="stat-item">
                <div class="stat-number">
                  {{ stats.failedTasks }}
                </div>
                <div class="stat-label">
                  {{ $t('admin.tasks.failedTasks') }}
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="4">
            <el-card class="stats-card timeout">
              <div class="stat-item">
                <div class="stat-number">
                  {{ stats.timeoutTasks }}
                </div>
                <div class="stat-label">
                  {{ $t('admin.tasks.timeoutTasks') }}
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 筛选器 -->
      <div class="filter-section">
        <el-form
          :inline="true"
          :model="filterForm"
          class="filter-form"
        >
          <el-form-item>
            <el-input
              v-model="filterForm.username"
              :placeholder="$t('admin.tasks.enterUsername')"
              clearable
              style="width: 120px"
            />
          </el-form-item>
          <el-form-item>
            <el-select
              v-model="filterForm.providerId"
              :placeholder="$t('admin.tasks.selectProvider')"
              clearable
              style="width: 150px"
            >
              <el-option
                v-for="provider in providers"
                :key="provider.id"
                :label="provider.name"
                :value="provider.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-select
              v-model="filterForm.taskType"
              :placeholder="$t('admin.tasks.selectTaskType')"
              clearable
              style="width: 120px"
            >
              <el-option
                :label="$t('admin.tasks.taskTypeCreate')"
                value="create"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeCreateInstance')"
                value="create_instance"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeStart')"
                value="start"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeStop')"
                value="stop"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeRestart')"
                value="restart"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeReset')"
                value="reset"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeRebuild')"
                value="rebuild"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeDelete')"
                value="delete"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeResetPassword')"
                value="reset-password"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeCreatePortMapping')"
                value="create-port-mapping"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeDeletePortMapping')"
                value="delete-port-mapping"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeSyncPortMappings')"
                value="sync-port-mappings"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeRepairPortMappings')"
                value="repair-port-mappings"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeCreateRedemptionInstance')"
                value="create_redemption_instance"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeSnapshotCreate')"
                value="snapshot-create"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeSnapshotDelete')"
                value="snapshot-delete"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeSnapshotRestore')"
                value="snapshot-restore"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeMonitorSync')"
                value="monitor-sync"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeAgentDeploy')"
                value="agent-deploy"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeAgentUninstall')"
                value="agent-uninstall"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeTrafficMonitorEnable')"
                value="traffic-monitor-enable"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeTrafficMonitorDisable')"
                value="traffic-monitor-disable"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeTrafficMonitorDetect')"
                value="traffic-monitor-detect"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeProviderImageCleanup')"
                value="provider-image-cleanup"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeProviderInstanceSync')"
                value="provider-instance-sync"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeProviderOrphanCleanup')"
                value="provider-orphan-cleanup"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeProviderHealthCheck')"
                value="provider-health-check"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeProviderIOLimitSync')"
                value="provider-io-limit-sync"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeProviderRuntimeReload')"
                value="provider-runtime-reload"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeProviderRecoverySync')"
                value="provider-recovery-sync"
              />
              <el-option
                :label="$t('admin.tasks.taskTypeProviderDelete')"
                value="provider-delete"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-select
              v-model="filterForm.status"
              :placeholder="$t('admin.tasks.selectStatus')"
              clearable
              style="width: 120px"
            >
              <el-option
                :label="$t('admin.tasks.statusPending')"
                value="pending"
              />
              <el-option
                :label="$t('admin.tasks.statusProcessing')"
                value="processing"
              />
              <el-option
                :label="$t('admin.tasks.statusRunning')"
                value="running"
              />
              <el-option
                :label="$t('admin.tasks.statusCompleted')"
                value="completed"
              />
              <el-option
                :label="$t('admin.tasks.statusFailed')"
                value="failed"
              />
              <el-option
                :label="$t('admin.tasks.statusCancelled')"
                value="cancelled"
              />
              <el-option
                :label="$t('admin.tasks.statusCancelling')"
                value="cancelling"
              />
              <el-option
                :label="$t('admin.tasks.statusTimeout')"
                value="timeout"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-select
              v-model="filterForm.instanceType"
              :placeholder="$t('admin.tasks.selectInstanceType')"
              clearable
              style="width: 120px"
            >
              <el-option
                :label="$t('admin.instances.typeContainer')"
                value="container"
              />
              <el-option
                :label="$t('admin.instances.typeVM')"
                value="vm"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="loadTasks"
            >
              {{ $t('common.filter') }}
            </el-button>
            <el-button @click="resetFilter">
              {{ $t('common.reset') }}
            </el-button>
            <el-button 
              :loading="loading"
              @click="loadTasks"
            >
              <el-icon><Refresh /></el-icon>
              {{ $t('common.refresh') }}
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 任务列表 -->
      <el-card class="tasks-card">
        <el-table
          v-loading="loading"
          :data="tasks"
          class="tasks-table"
          :cell-style="{ padding: '12px 0' }"
          :header-cell-style="{ background: '#f5f7fa', padding: '14px 0', fontWeight: '600' }"
          :default-sort="{prop: 'createdAt', order: 'descending'}"
        >
          <el-table-column
            prop="id"
            label="ID"
            width="80"
            align="center"
            sortable
          />
          <el-table-column
            prop="userName"
            :label="$t('common.user')"
            width="140"
            show-overflow-tooltip
          />
          <el-table-column
            prop="taskType"
            :label="$t('admin.tasks.taskType')"
            width="120"
            align="center"
          >
            <template #default="{ row }">
              <el-tag size="small">
                {{ getTaskTypeText(row.taskType) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            prop="status"
            :label="$t('common.status')"
            width="110"
            align="center"
          >
            <template #default="{ row }">
              <el-tag
                :type="getTaskStatusType(row.status)"
                size="small"
              >
                {{ getTaskStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            prop="progress"
            :label="$t('admin.tasks.progress')"
            width="140"
            align="center"
          >
            <template #default="{ row }">
              <el-progress
                v-if="row.status === 'running' || row.status === 'processing'"
                :percentage="row.progress"
                :status="row.status === 'failed' ? 'exception' : (row.status === 'completed' ? 'success' : undefined)"
                size="small"
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="providerName"
            :label="$t('admin.tasks.provider')"
            width="140"
            show-overflow-tooltip
          />
          <el-table-column
            prop="instanceName"
            :label="$t('admin.tasks.instance')"
            min-width="180"
          >
            <template #default="{ row }">
              <div
                v-if="row.instanceName"
                class="instance-info"
              >
                <div class="instance-name">
                  {{ row.instanceName }}
                </div>
                <el-tag
                  v-if="row.instanceType"
                  size="small"
                  :type="row.instanceType === 'vm' ? 'warning' : 'info'"
                >
                  {{ row.instanceType === 'vm' ? $t('admin.instances.typeVM') : $t('admin.instances.typeContainer') }}
                </el-tag>
              </div>
              <span
                v-else
                class="text-gray"
              >-</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="createdAt"
            :label="$t('common.createTime')"
            width="180"
            align="center"
            sortable
          >
            <template #default="{ row }">
              {{ formatDateTime(row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column
            prop="remainingTime"
            :label="$t('admin.tasks.remainingTime')"
            min-width="160"
            align="center"
          >
            <template #default="{ row }">
              <span v-if="row.status === 'running' && row.remainingTime > 0">
                {{ formatDuration(row.remainingTime) }}
              </span>
              <span
                v-else
                class="text-gray"
              >-</span>
            </template>
          </el-table-column>
          <el-table-column
            :label="$t('common.actions')"
            width="220"
            fixed="right"
            align="center"
          >
            <template #default="{ row }">
              <div class="action-buttons">
                <el-button
                  v-if="row.canForceStop"
                  type="danger"
                  size="small"
                  @click="showForceStopDialog(row)"
                >
                  {{ $t('admin.tasks.forceStop') }}
                </el-button>
                <el-button
                  v-if="row.status === 'pending'"
                  type="warning"
                  size="small"
                  @click="cancelTask(row)"
                >
                  {{ $t('admin.tasks.cancelTask') }}
                </el-button>
                <el-button
                  size="small"
                  @click="viewTaskDetail(row)"
                >
                  {{ $t('common.details') }}
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadTasks"
            @current-change="loadTasks"
          />
        </div>
      </el-card>

      <!-- 强制停止任务对话框 -->
      <el-dialog
        v-model="forceStopDialog.visible"
        :title="$t('admin.tasks.forceStopTask')"
        width="500px"
      >
        <el-form
          :model="forceStopDialog.form"
          label-width="80px"
        >
          <el-form-item :label="$t('admin.tasks.taskInfo')">
            <div class="task-info">
              <p><strong>ID:</strong> {{ forceStopDialog.task?.id }}</p>
              <p><strong>{{ $t('admin.tasks.taskType') }}:</strong> {{ getTaskTypeText(forceStopDialog.task?.taskType) }}</p>
              <p><strong>{{ $t('common.user') }}:</strong> {{ forceStopDialog.task?.userName }}</p>
              <p><strong>{{ $t('admin.tasks.instance') }}:</strong> {{ forceStopDialog.task?.instanceName || '-' }}</p>
            </div>
          </el-form-item>
          <el-form-item 
            :label="$t('admin.tasks.stopReason')"
          >
            <el-input
              v-model="forceStopDialog.form.reason"
              type="textarea"
              :rows="3"
              :placeholder="$t('admin.tasks.enterStopReason')"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="forceStopDialog.visible = false">
              {{ $t('common.cancel') }}
            </el-button>
            <el-button
              type="danger"
              :loading="forceStopDialog.loading"
              @click="confirmForceStop"
            >
              {{ $t('admin.tasks.forceStop') }}
            </el-button>
          </span>
        </template>
      </el-dialog>

      <!-- 任务详情对话框 -->
      <el-dialog
        v-model="detailDialog.visible"
        :title="$t('admin.tasks.taskDetails')"
        width="600px"
      >
        <div
          v-if="detailDialog.task"
          class="task-detail"
        >
          <el-descriptions
            :column="2"
            border
          >
            <el-descriptions-item :label="$t('admin.tasks.taskId')">
              {{ detailDialog.task.id }}
            </el-descriptions-item>
            <el-descriptions-item label="UUID">
              {{ detailDialog.task.uuid }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.tasks.taskType')">
              {{ getTaskTypeText(detailDialog.task.taskType) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('common.status')">
              <el-tag :type="getTaskStatusType(detailDialog.task.status)">
                {{ getTaskStatusText(detailDialog.task.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('common.user')">
              {{ detailDialog.task.userName }} (ID: {{ detailDialog.task.userId }})
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.tasks.provider')">
              {{ detailDialog.task.providerName }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.tasks.instance')">
              <div v-if="detailDialog.task.instanceName">
                {{ detailDialog.task.instanceName }}
                <el-tag
                  v-if="detailDialog.task.instanceType"
                  size="mini"
                  :type="detailDialog.task.instanceType === 'vm' ? 'warning' : 'info'"
                >
                  {{ detailDialog.task.instanceType === 'vm' ? $t('admin.instances.typeVM') : $t('admin.instances.typeContainer') }}
                </el-tag>
              </div>
              <span v-else>-</span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.tasks.progress')">
              <el-progress
                v-if="detailDialog.task.progress !== undefined && detailDialog.task.progress !== null"
                :percentage="detailDialog.task.progress"
                :status="detailDialog.task.status === 'failed' ? 'exception' : (detailDialog.task.status === 'completed' ? 'success' : undefined)"
              />
              <span v-else>-</span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.tasks.timeoutDuration')">
              {{ formatDuration(detailDialog.task.timeoutDuration) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.tasks.remainingTime')">
              <span v-if="detailDialog.task.status === 'running' && detailDialog.task.remainingTime > 0">
                {{ formatDuration(detailDialog.task.remainingTime) }}
              </span>
              <span v-else>-</span>
            </el-descriptions-item>
            <el-descriptions-item
              :label="$t('common.createTime')"
              :span="2"
            >
              {{ formatDateTime(detailDialog.task.createdAt) }}
            </el-descriptions-item>
            <el-descriptions-item
              :label="$t('admin.tasks.startTime')"
              :span="2"
            >
              {{ detailDialog.task.startedAt ? formatDateTime(detailDialog.task.startedAt) : '-' }}
            </el-descriptions-item>
            <el-descriptions-item
              :label="$t('admin.tasks.completionTime')"
              :span="2"
            >
              {{ detailDialog.task.completedAt ? formatDateTime(detailDialog.task.completedAt) : '-' }}
            </el-descriptions-item>
            <el-descriptions-item
              v-if="detailDialog.task.errorMessage"
              :label="$t('admin.tasks.errorMessage')"
              :span="2"
            >
              <el-text type="danger">
                {{ detailDialog.task.errorMessage }}
              </el-text>
            </el-descriptions-item>
            <el-descriptions-item
              v-if="detailDialog.task.cancelReason"
              :label="$t('admin.tasks.cancelReason')"
              :span="2"
            >
              <el-text type="warning">
                {{ detailDialog.task.cancelReason }}
              </el-text>
            </el-descriptions-item>
            <el-descriptions-item
              v-if="shouldShowPreallocatedConfig(detailDialog.task)"
              :label="$t('admin.tasks.preallocatedConfig')"
              :span="2"
            >
              <template v-if="detailDialog.task.preallocatedCpu && detailDialog.task.preallocatedCpu > 0">
                <el-tag
                  size="small"
                  type="info"
                >
                  CPU: {{ detailDialog.task.preallocatedCpu }} {{ $t('common.core') }}
                </el-tag>
                <el-tag
                  size="small"
                  type="info"
                  style="margin-left: 8px;"
                >
                  {{ $t('admin.tasks.memory') }}: {{ (detailDialog.task.preallocatedMemory / 1024).toFixed(1) }} GB
                </el-tag>
                <el-tag
                  size="small"
                  type="info"
                  style="margin-left: 8px;"
                >
                  {{ $t('admin.tasks.disk') }}: {{ (detailDialog.task.preallocatedDisk / 1024).toFixed(1) }} GB
                </el-tag>
                <el-tag
                  size="small"
                  type="info"
                  style="margin-left: 8px;"
                >
                  {{ $t('admin.tasks.bandwidth') }}: {{ detailDialog.task.preallocatedBandwidth }} Mbps
                </el-tag>
              </template>
              <template v-else>
                <el-text type="info">
                  {{ $t('admin.tasks.noPreallocatedConfig') }}
                </el-text>
              </template>
            </el-descriptions-item>
            <el-descriptions-item
              v-if="detailDialog.task.statusMessage"
              :label="$t('admin.tasks.statusMessage')"
              :span="2"
            >
              {{ translateStepMsg(detailDialog.task.statusMessage) }}
            </el-descriptions-item>
            <el-descriptions-item
              :label="$t('admin.tasks.progressLogs')"
              :span="2"
            >
              <div v-if="detailDialog.logsLoading">
                <el-text type="info">
                  {{ $t('common.loading') }}
                </el-text>
              </div>
              <div v-else-if="detailDialog.task.progressLogs">
                <task-steps-panel
                  :task-type="detailDialog.task.taskType"
                  :progress-logs="detailDialog.task.progressLogs"
                  :task-status="detailDialog.task.status"
                  style="margin-bottom: 10px;"
                />
                <el-button
                  link
                  size="small"
                  @click="toggleProgressLogs(detailDialog.task.id)"
                >
                  {{ expandedLogTaskIds.has(detailDialog.task.id) ? $t('admin.tasks.hideProgressLogs') : $t('admin.tasks.showProgressLogs') }}
                </el-button>
                <div
                  v-if="expandedLogTaskIds.has(detailDialog.task.id)"
                  class="progress-log-list"
                >
                  <div
                    v-for="(entry, idx) in parseProgressLogs(detailDialog.task.progressLogs)"
                    :key="idx"
                    class="progress-log-entry"
                  >
                    <span class="log-time">{{ entry.t }}</span>
                    <el-tag
                      size="small"
                      type="info"
                      style="margin: 0 6px;"
                    >
                      {{ entry.p }}%
                    </el-tag>
                    <div class="log-content">
                      <span class="log-msg">{{ translateStepMsg(entry.m) }}</span>
                      <div
                        v-if="entry.command || entry.output || entry.error"
                        class="log-detail"
                      >
                        <div v-if="entry.command">
                          <strong>Command</strong>
                          <pre>{{ entry.command }}</pre>
                        </div>
                        <div v-if="entry.output">
                          <strong>Output</strong>
                          <pre>{{ entry.output }}</pre>
                        </div>
                        <div v-if="entry.error">
                          <strong>Error</strong>
                          <pre>{{ entry.error }}</pre>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <el-text
                v-else
                type="info"
              >
                {{ $t('admin.tasks.noProgressLogs') }}
              </el-text>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup>
import { Refresh } from '@element-plus/icons-vue'
import { useTaskManagement } from './composables/useTaskManagement'
import TaskStepsPanel from '@/components/TaskStepsPanel.vue'

const {
  loading, poolLoading, tasks, providers, total, stats, poolStatus, isSuperAdmin,
  filterForm, pagination,
  forceStopDialog, detailDialog, expandedLogTaskIds,
  loadTasks, resetFilter, loadTaskPoolStatus, toggleTaskPool,
  showForceStopDialog, confirmForceStop,
  cancelTask, viewTaskDetail,
  parseProgressLogs, translateStepMsg, toggleProgressLogs,
  shouldShowPreallocatedConfig,
  getTaskTypeText, getTaskStatusType, getTaskStatusText,
  formatDateTime, formatDuration,
  t
} = useTaskManagement()
</script>

<style scoped lang="scss" src="./index.scss">

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
