<template>
  <div class="instances-container">
    <el-card>
      <template #header>
        <div class="header-row">
          <span>{{ $t('admin.instances.title') }}</span>
          <div class="header-actions">
            <el-button
              v-if="selectedInstances.length > 0"
              type="success"
              @click="batchStartInstances"
            >
              {{ $t('admin.instances.batchStart') }} ({{ selectedInstances.length }})
            </el-button>
            <el-button
              v-if="selectedInstances.length > 0"
              type="warning"
              @click="batchStopInstances"
            >
              {{ $t('admin.instances.batchStop') }} ({{ selectedInstances.length }})
            </el-button>
            <el-button
              v-if="selectedInstances.length > 0"
              type="danger"
              @click="batchDeleteInstances"
            >
              {{ $t('admin.instances.batchDelete') }} ({{ selectedInstances.length }})
            </el-button>
            <el-button
              type="primary"
              :loading="loading"
              @click="loadInstances"
            >
              {{ $t('common.refresh') }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- 筛选条件 -->
      <div class="filter-row">
        <el-input
          v-model="filters.instanceName"
          :placeholder="$t('admin.instances.searchByInstanceName')"
          style="width: 200px; margin-right: 10px;"
          clearable
        />
        <el-input
          v-model="filters.providerName"
          :placeholder="$t('admin.instances.searchByProviderName')"
          style="width: 200px; margin-right: 10px;"
          clearable
        />
        <el-input
          v-model="filters.ownerName"
          :placeholder="$t('admin.instances.searchByOwner')"
          style="width: 200px; margin-right: 10px;"
          clearable
        />
        <el-select
          v-model="filters.status"
          :placeholder="$t('admin.instances.filterByStatus')"
          style="width: 120px; margin-right: 10px;"
          clearable
        >
          <el-option
            :label="$t('admin.instances.statusRunning')"
            value="running"
          />
          <el-option
            :label="$t('admin.instances.statusStopped')"
            value="stopped"
          />
          <el-option
            :label="$t('admin.instances.statusCreating')"
            value="creating"
          />
          <el-option
            :label="$t('admin.instances.statusStarting')"
            value="starting"
          />
          <el-option
            :label="$t('admin.instances.statusStopping')"
            value="stopping"
          />
          <el-option
            :label="$t('admin.instances.statusRestarting')"
            value="restarting"
          />
          <el-option
            :label="$t('admin.instances.statusRebuilding')"
            value="rebuilding"
          />
          <el-option
            :label="$t('admin.instances.statusResetting')"
            value="resetting"
          />
          <el-option
            :label="$t('admin.instances.statusError')"
            value="error"
          />
          <el-option
            :label="$t('admin.instances.statusFailed')"
            value="failed"
          />
          <el-option
            :label="$t('admin.instances.statusDeleting')"
            value="deleting"
          />
          <el-option
            :label="$t('admin.instances.statusDeleted')"
            value="deleted"
          />
        </el-select>
        <el-select
          v-model="filters.instanceType"
          :placeholder="$t('admin.instances.filterByType')"
          style="width: 120px; margin-right: 10px;"
          clearable
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
        <el-button
          type="primary"
          @click="handleSearch"
        >
          {{ $t('common.search') }}
        </el-button>
        <el-button
          @click="handleReset"
        >
          {{ $t('common.reset') }}
        </el-button>
      </div>

      <el-table
        ref="tableRef"
        v-loading="loading"
        :data="instances"
        style="width: 100%"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column
          prop="name"
          :label="$t('admin.instances.instanceName')"
          min-width="150"
          show-overflow-tooltip
          fixed="left"
        />
        <el-table-column
          prop="userName"
          :label="$t('admin.instances.owner')"
          width="100"
        />
        <el-table-column
          prop="providerName"
          :label="$t('admin.instances.provider')"
          min-width="160"
          show-overflow-tooltip
        />
        <el-table-column
          prop="instance_type"
          :label="$t('admin.instances.instanceType')"
          min-width="100"
        >
          <template #default="scope">
            <el-tag
              :type="scope.row.instance_type === 'container' ? 'primary' : 'success'"
              size="small"
            >
              {{ scope.row.instance_type === 'container' ? $t('admin.instances.typeContainer') : $t('admin.instances.typeVM') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('admin.instances.acceleratorConfig')"
          min-width="190"
        >
          <template #default="scope">
            <template v-if="scope.row.gpuEnabled || scope.row.npuEnabled">
              <el-tag
                v-if="scope.row.gpuEnabled"
                type="warning"
                size="small"
                style="margin-right: 4px;"
              >
                {{ scope.row.gpuDeviceIds ? 'GPU:' + scope.row.gpuDeviceIds : $t('admin.instances.gpuEnabled') }}
              </el-tag>
              <el-tag
                v-if="scope.row.npuEnabled"
                type="danger"
                size="small"
              >
                {{ scope.row.npuDeviceIds ? 'NPU:' + scope.row.npuDeviceIds : $t('admin.instances.gpuEnabled') }}
              </el-tag>
            </template>
            <span
              v-else
              style="color: #c0c4cc; font-size: 12px;"
            >—</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="sshPort"
          :label="$t('admin.instances.sshPort')"
          min-width="120"
        />
        <el-table-column
          prop="osType"
          :label="$t('admin.instances.system')"
          min-width="90"
        />
        <el-table-column
          :label="$t('admin.instances.instanceStatus')"
          width="100"
        >
          <template #default="scope">
            <el-tag
              :type="getStatusType(scope.row.status)"
              size="small"
            >
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('admin.instances.trafficStatus')"
          width="100"
        >
          <template #default="scope">
            <el-tag
              v-if="scope.row.trafficLimited"
              type="danger"
              size="small"
            >
              {{ $t('admin.instances.limited') }}
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
          prop="createdAt"
          :label="$t('common.createTime')"
          width="140"
        >
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="expiresAt"
          :label="$t('admin.instances.expiryTime')"
          width="140"
        >
          <template #default="scope">
            <span :class="{ 'expired': isExpired(scope.row.expiresAt), 'expiring-soon': isExpiringSoon(scope.row.expiresAt) }">
              {{ formatDate(scope.row.expiresAt) }}
            </span>
            <div
              v-if="scope.row.isManualExpiry"
              style="margin-top: 4px;"
            >
              <el-tag
                size="small"
                type="info"
              >
                {{ $t('admin.instances.manualExpiry') }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          prop="isFrozen"
          :label="$t('admin.instances.freezeStatus')"
          min-width="150"
          align="center"
        >
          <template #default="scope">
            <el-tag
              :type="scope.row.isFrozen ? 'danger' : 'success'"
              size="small"
            >
              {{ scope.row.isFrozen ? $t('admin.instances.frozen') : $t('admin.instances.normal') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.actions')"
          width="430"
          fixed="right"
        >
          <template #default="scope">
            <div class="action-buttons">
              <el-button
                size="small"
                type="info"
                :disabled="!canOpenInstanceDetail(scope.row)"
                @click="viewInstanceDetail(scope.row)"
              >
                {{ $t('admin.instances.viewDetail') }}
              </el-button>
              <el-button
                size="small"
                type="primary"
                :disabled="isInstanceBusy(scope.row)"
                @click="showActionDialog(scope.row)"
              >
                {{ $t('admin.instances.actions') }}
              </el-button>
              <el-button
                size="small"
                type="success"
                :disabled="isInstanceBusy(scope.row) || scope.row.status !== 'running' || (!scope.row.hasSshMapping && scope.row.networkType === 'no_port_mapping')"
                :title="(!scope.row.hasSshMapping && scope.row.networkType === 'no_port_mapping') ? $t('admin.instances.sshNoPortMapping') : ''"
                @click="openSSHTerminal(scope.row)"
              >
                {{ $t('admin.instances.connect') }}
              </el-button>
              <el-tooltip :content="$t('user.instanceDetail.webConsole')">
                <el-button
                  size="small"
                  type="info"
                  :aria-label="$t('user.instanceDetail.webConsole')"
                  :disabled="isInstanceBusy(scope.row) || scope.row.status !== 'running'"
                  @click="openConsole(scope.row)"
                >
                  <el-icon><Monitor /></el-icon>
                </el-button>
              </el-tooltip>
              <el-button
                size="small"
                type="warning"
                :disabled="!canOpenInstanceDetail(scope.row)"
                @click="showTransferDialog(scope.row)"
              >
                {{ $t('admin.instances.transfer') }}
              </el-button>
              <el-button
                size="small"
                type="success"
                :disabled="!canOpenInstanceDetail(scope.row)"
                @click="createShareLink(scope.row)"
              >
                <el-icon><Link /></el-icon>
                {{ $t('admin.instances.share') }}
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-row">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 实例详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="$t('admin.instances.instanceDetails')"
      width="60%"
    >
      <div
        v-if="selectedInstance"
        class="instance-detail"
      >
        <el-descriptions
          :column="2"
          border
        >
          <el-descriptions-item :label="$t('admin.instances.instanceName')">
            {{ selectedInstance.name }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.uuid')">
            {{ selectedInstance.uuid }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.owner')">
            {{ selectedInstance.userName }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.provider')">
            {{ selectedInstance.providerName }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.instanceType')">
            <el-tag :type="selectedInstance.instance_type === 'container' ? 'primary' : 'success'">
              {{ selectedInstance.instance_type === 'container' ? $t('admin.instances.typeContainer') : $t('admin.instances.typeVM') }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.instanceStatus')">
            <el-tag :type="getStatusType(selectedInstance.status)">
              {{ getStatusText(selectedInstance.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.image')">
            {{ selectedInstance.image }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.operatingSystem')">
            {{ selectedInstance.osType }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.cpu')">
            {{ selectedInstance.cpu }}{{ $t('admin.instances.cores') }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.memory')">
            {{ formatMemory(selectedInstance.memory) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.disk')">
            {{ formatDisk(selectedInstance.disk) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.bandwidth')">
            {{ selectedInstance.bandwidth }}Mbps
          </el-descriptions-item>
          <el-descriptions-item
            v-if="selectedInstance.gpuEnabled || selectedInstance.npuEnabled"
            :label="$t('admin.instances.acceleratorConfig')"
          >
            <template v-if="selectedInstance.gpuEnabled">
              <el-tag
                type="warning"
                size="small"
              >
                {{ $t('admin.instances.gpu') }} {{ selectedInstance.gpuDeviceIds || $t('admin.instances.gpuEnabled') }}
              </el-tag>&nbsp;
            </template>
            <template v-if="selectedInstance.npuEnabled">
              <el-tag
                type="danger"
                size="small"
              >
                {{ $t('admin.instances.npu') }} {{ selectedInstance.npuDeviceIds || $t('admin.instances.gpuEnabled') }}
              </el-tag>
            </template>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.publicIPv4')">
            {{ selectedInstance.publicIP || $t('admin.instances.unassigned') }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.privateIPv4')">
            {{ selectedInstance.privateIP || $t('admin.instances.unassigned') }}
          </el-descriptions-item>
          <el-descriptions-item
            v-if="selectedInstance.ipv6Address"
            :label="$t('admin.instances.privateIPv6')"
          >
            {{ selectedInstance.ipv6Address }}
          </el-descriptions-item>
          <el-descriptions-item
            v-if="selectedInstance.publicIPv6"
            :label="$t('admin.instances.publicIPv6')"
          >
            {{ selectedInstance.publicIPv6 }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.sshPort')">
            {{ selectedInstance.sshPort }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.username')">
            {{ selectedInstance.username }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.password')">
            <span v-if="showPassword">{{ selectedInstance.password }}</span>
            <span v-else>••••••••</span>
            <el-button
              link
              @click="showPassword = !showPassword"
            >
              {{ showPassword ? $t('admin.instances.hide') : $t('admin.instances.show') }}
            </el-button>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.trafficLimit')">
            <el-tag
              v-if="selectedInstance.trafficLimited"
              type="danger"
            >
              {{ $t('admin.instances.limited') }}
            </el-tag>
            <el-tag
              v-else
              type="success"
            >
              {{ $t('common.normal') }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.networkInterfaceV4')">
            {{ selectedInstance.pmacctInterfaceV4 || $t('admin.instances.notSet') }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.networkInterfaceV6')">
            {{ selectedInstance.pmacctInterfaceV6 || $t('admin.instances.notSet') }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('common.createTime')">
            {{ formatDate(selectedInstance.createdAt) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('common.updatedAt')">
            {{ formatDate(selectedInstance.updatedAt) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.expiryTime')">
            <span :class="{ 'expired': isExpired(selectedInstance.expiresAt), 'expiring-soon': isExpiringSoon(selectedInstance.expiresAt) }">
              {{ formatDate(selectedInstance.expiresAt) }}
            </span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('admin.instances.healthStatus')">
            <el-tag :type="selectedInstance.healthStatus === 'healthy' ? 'success' : 'danger'">
              {{ selectedInstance.healthStatus === 'healthy' ? $t('admin.instances.healthy') : $t('admin.instances.unhealthy') }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <div
          class="traffic-info"
          style="margin-top: 20px;"
        >
          <h4>{{ $t('admin.instances.historicalTraffic') }}</h4>
          <el-descriptions
            :column="2"
            border
          >
            <el-descriptions-item :label="$t('admin.instances.inboundTraffic')">
              {{ formatTraffic(selectedInstance.usedTrafficIn) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('admin.instances.outboundTraffic')">
              {{ formatTraffic(selectedInstance.usedTrafficOut) }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
      <template #footer>
        <el-button
          type="info"
          :disabled="selectedInstance?.status !== 'running' || isInstanceBusy(selectedInstance)"
          @click="openConsole(selectedInstance)"
        >
          <el-icon><Monitor /></el-icon>
          {{ $t('user.instanceDetail.webConsole') }}
        </el-button>
        <el-button
          type="primary"
          :loading="accessLoading"
          @click="showInstanceAccessDialog(selectedInstance)"
        >
          <el-icon><EditPen /></el-icon>
          {{ $t('admin.instances.editAccess') }}
        </el-button>
        <el-button @click="detailDialogVisible = false">
          {{ $t('common.close') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 实例操作对话框 -->
    <el-dialog
      v-model="actionDialogVisible"
      :title="$t('admin.instances.instanceActions')"
      width="400px"
    >
      <div
        v-if="actionInstance"
        class="action-dialog-content"
      >
        <el-button
          type="success"
          :disabled="isInstanceBusy(actionInstance) || actionInstance.status === 'running' || actionInstance.status === 'starting'"
          :loading="actionLoading"
          style="width: 100%; margin-bottom: 10px;"
          @click="performAction('start')"
        >
          <el-icon><VideoPlay /></el-icon>
          {{ $t('common.start') }}
        </el-button>
        <el-button
          type="warning"
          :disabled="isInstanceBusy(actionInstance) || actionInstance.status === 'stopped' || actionInstance.status === 'stopping'"
          :loading="actionLoading"
          style="width: 100%; margin-bottom: 10px;"
          @click="performAction('stop')"
        >
          <el-icon><VideoPause /></el-icon>
          {{ $t('common.stop') }}
        </el-button>
        <el-button
          type="primary"
          :disabled="isInstanceBusy(actionInstance) || actionInstance.status !== 'running'"
          :loading="actionLoading"
          style="width: 100%; margin-bottom: 10px;"
          @click="performAction('restart')"
        >
          <el-icon><Refresh /></el-icon>
          {{ $t('common.restart') }}
        </el-button>
        <el-button
          type="info"
          :disabled="isInstanceBusy(actionInstance) || actionInstance.status !== 'running'"
          :loading="actionLoading"
          style="width: 100%; margin-bottom: 10px;"
          @click="performAction('resetPassword')"
        >
          <el-icon><Lock /></el-icon>
          {{ $t('admin.instances.resetPassword') }}
        </el-button>
        <el-button
          type="warning"
          :disabled="isInstanceBusy(actionInstance) || actionInstance.status !== 'running'"
          :loading="actionLoading"
          style="width: 100%; margin-bottom: 10px;"
          @click="performAction('reset')"
        >
          <el-icon><RefreshRight /></el-icon>
          {{ $t('admin.instances.resetSystem') }}
        </el-button>
        <el-button
          type="primary"
          :disabled="isInstanceBusy(actionInstance)"
          style="width: 100%; margin-bottom: 10px;"
          @click="showEgressDialog(actionInstance)"
        >
          <el-icon><Connection /></el-icon>
          {{ $t('admin.instances.egressAction') }}
        </el-button>
        <el-button
          type="success"
          :disabled="isInstanceBusy(actionInstance) || actionInstance.status !== 'running'"
          :loading="actionLoading"
          style="width: 100%; margin-bottom: 10px;"
          @click="performAction('vnc')"
        >
          <el-icon><Monitor /></el-icon>
          {{ $t('admin.instances.vncConsole') || 'VNC控制台' }}
        </el-button>
        <el-divider />
        <el-button
          type="info"
          :disabled="isInstanceBusy(actionInstance)"
          style="width: 100%; margin-bottom: 10px;"
          @click="performAction('setExpiry')"
        >
          {{ $t('admin.instances.setExpiry') }}
        </el-button>
        <el-button
          type="primary"
          :disabled="isInstanceBusy(actionInstance)"
          style="width: 100%; margin-bottom: 10px;"
          @click="performAction('renew')"
        >
          {{ $t('admin.instances.renew') || '续期' }}
        </el-button>
        <el-button
          v-if="!actionInstance.isFrozen"
          type="warning"
          :disabled="isInstanceBusy(actionInstance)"
          style="width: 100%; margin-bottom: 10px;"
          @click="performAction('freeze')"
        >
          {{ $t('admin.instances.freeze') }}
        </el-button>
        <el-button
          v-else
          type="success"
          :disabled="isInstanceBusy(actionInstance)"
          style="width: 100%; margin-bottom: 10px;"
          @click="performAction('unfreeze')"
        >
          {{ $t('admin.instances.unfreeze') }}
        </el-button>
        <el-divider />
        <el-button
          type="danger"
          :loading="actionLoading"
          style="width: 100%;"
          @click="performAction('delete')"
        >
          <el-icon><Delete /></el-icon>
          {{ $t('common.delete') }}
        </el-button>
      </div>
    </el-dialog>

    <EgressDialog
      v-model="egressDialogVisible"
      :instance="egressInstance"
      @updated="loadInstances"
    />

    <InstanceAccessDialog
      v-model="accessDialogVisible"
      :instance="accessInstance"
      @updated="refreshInstanceAccess"
    />

    <VNCDialog
      v-if="consoleInstance"
      v-model="consoleDialogVisible"
      :instance-id="consoleInstance.id"
      :instance-name="consoleInstance.name"
      scope="admin"
    />

    <!-- 转移实例对话框 -->
    <el-dialog
      v-model="transferDialogVisible"
      :title="$t('admin.instances.transferInstance')"
      width="400px"
    >
      <el-form label-width="100px">
        <el-form-item :label="$t('admin.instances.instanceName')">
          <el-input
            :model-value="transferForm.instanceName"
            disabled
          />
        </el-form-item>
        <el-form-item :label="$t('admin.instances.targetUserId')">
          <el-select
            v-model="transferForm.targetUserId"
            filterable
            remote
            :remote-method="searchUsers"
            :loading="searchingUsers"
            :placeholder="$t('admin.instances.searchUserPlaceholder')"
            style="width: 100%"
          >
            <el-option
              v-for="user in userOptions"
              :key="user.id"
              :label="`${user.nickname || user.username} (ID: ${user.id})`"
              :value="user.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="transferDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="transferLoading"
          @click="confirmTransfer"
        >
          {{ $t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 重装系统镜像选择对话框 -->
    <el-dialog
      v-model="resetImageDialogVisible"
      :title="$t('admin.instances.selectImage') || '选择系统镜像'"
      width="500px"
    >
      <div v-loading="loadingResetImages">
        <el-select
          v-model="selectedResetImage"
          :placeholder="$t('admin.instances.selectImage') || '请选择系统镜像'"
          style="width: 100%"
        >
          <el-option
            v-for="img in resetImageList"
            :key="img.id"
            :label="img.name"
            :value="img.url || img.name"
          />
        </el-select>
      </div>
      <template #footer>
        <el-button @click="resetImageDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="actionLoading"
          @click="confirmResetWithImage"
        >
          {{ $t('admin.instances.confirmReset') || '确认重装' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- VNC 控制台对话框 -->
    <VNCDialog
      v-model="vncDialogVisible"
      :instance-id="vncInstanceId"
      :instance-name="vncInstanceName"
      admin-mode
    />
    <!-- 设置过期时间（日期选择器） -->
    <SetExpiryDialog
      :visible="showSetExpiryDialog"
      :form="expiryForm"
      :loading="expiryLoading"
      @confirm="confirmSetInstanceExpiry"
      @update:visible="showSetExpiryDialog = $event"
    />
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { 
  VideoPlay, 
  VideoPause, 
  Refresh, 
  RefreshRight, 
  Lock, 
  Delete,
  Link,
  Connection,
  EditPen,
  Monitor
} from '@element-plus/icons-vue'
import { useInstanceManagement } from './composables/useInstanceManagement'
import EgressDialog from './components/EgressDialog.vue'
import InstanceAccessDialog from './components/InstanceAccessDialog.vue'
import VNCDialog from '@/components/VNCDialog.vue'
import SetExpiryDialog from './components/SetExpiryDialog.vue'

const {
  instances, loading, detailDialogVisible, actionDialogVisible, egressDialogVisible, accessDialogVisible, consoleDialogVisible, accessLoading,
  selectedInstance, actionInstance, egressInstance, accessInstance, consoleInstance, actionLoading, showPassword,
  selectedInstances, transferDialogVisible, transferLoading, transferForm, tableRef,
  filters, pagination,
  loadInstances, handleSearch, handleReset, handleSizeChange, handleCurrentChange,
  viewInstanceDetail, showActionDialog, showEgressDialog, showInstanceAccessDialog, refreshInstanceAccess, performAction,
  getStatusType, getStatusText, formatDate, formatMemory, formatDisk, formatTraffic,
  isExpired, isExpiringSoon, openSSHTerminal, openConsole,
  handleSelectionChange, batchDeleteInstances, batchStartInstances, batchStopInstances,
  showTransferDialog, confirmTransfer, handleWindowResize,
  searchUsers, searchingUsers, userOptions, canOpenInstanceDetail, isInstanceBusy, createShareLink,
  resetImageDialogVisible, resetImageList, selectedResetImage, loadingResetImages, confirmResetWithImage,
  vncDialogVisible, vncInstanceId, vncInstanceName,
  showSetExpiryDialog, expiryForm, expiryLoading, confirmSetInstanceExpiry
} = useInstanceManagement()

onMounted(() => {
  loadInstances()
  window.addEventListener('resize', handleWindowResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleWindowResize)
})
</script>

<style scoped>
.instances-container {
  width: 100%;
  height: 100%;
}

.header-row {
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
  gap: 10px;
  align-items: center;
}

.filter-row {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

.action-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.action-buttons .el-button {
  margin: 0;
}

.pagination-row {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.instance-detail {
  max-height: 70vh;
  overflow-y: auto;
}

.traffic-info {
  border-top: 1px solid #ebeef5;
  padding-top: 20px;
}

.expired {
  color: #f56c6c;
  font-weight: bold;
}

.expiring-soon {
  color: #e6a23c;
  font-weight: bold;
}

.action-dialog-content {
  padding: 10px 0;
}

.action-dialog-content .el-button {
  margin: 0;
}

/* 确认表格在窗口调整时的显示问题 */
:deep(.el-table) {
  overflow: visible !important;
}

:deep(.el-table__body-wrapper) {
  overflow-x: auto;
}

/* 确保fixed列正确渲染 */
:deep(.el-table__fixed),
:deep(.el-table__fixed-right) {
  height: auto !important;
}

:deep(.el-table__fixed-body-wrapper),
:deep(.el-table__fixed-right .el-table__fixed-body-wrapper) {
  height: auto !important;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .action-buttons {
    flex-direction: column;
  }
  
  .action-buttons .el-button {
    width: 100%;
    margin-bottom: 2px;
  }
}

@media (max-width: 768px) {
  .filter-row {
    flex-direction: column;
    align-items: stretch;
  }
  
  .filter-row > * {
    width: 100% !important;
    margin-bottom: 10px;
  }
  
  .header-actions {
    flex-wrap: wrap;
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
