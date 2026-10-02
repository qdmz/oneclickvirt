<template>
  <div class="redemption-codes-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('admin.redemptionCodes.title') }}</span>
          <el-button
            type="primary"
            @click="openCreateDialog"
          >
            {{ t('admin.redemptionCodes.batchCreate') }}
          </el-button>
        </div>
      </template>

      <!-- 过滤栏 -->
      <el-row
        :gutter="12"
        class="filter-bar"
      >
        <el-col :span="6">
          <el-input
            v-model="filterCode"
            :placeholder="t('admin.redemptionCodes.searchCode')"
            clearable
            @change="handleFilterChange"
          />
        </el-col>
        <el-col :span="5">
          <el-select
            v-model="filterStatus"
            :placeholder="t('admin.redemptionCodes.filterStatus')"
            clearable
            @change="handleFilterChange"
          >
            <el-option
              value=""
              :label="t('admin.redemptionCodes.allStatus')"
            />
            <el-option
              value="pending_create"
              :label="t('admin.redemptionCodes.statusPendingCreate')"
            />
            <el-option
              value="creating"
              :label="t('admin.redemptionCodes.statusCreating')"
            />
            <el-option
              value="pending_use"
              :label="t('admin.redemptionCodes.statusPendingUse')"
            />
            <el-option
              value="used"
              :label="t('admin.redemptionCodes.statusUsed')"
            />
            <el-option
              value="deleting"
              :label="t('admin.redemptionCodes.statusDeleting')"
            />
          </el-select>
        </el-col>
        <el-col :span="5">
          <el-select
            v-model="filterProvider"
            :placeholder="t('admin.redemptionCodes.filterProvider')"
            clearable
            @change="handleFilterChange"
          >
            <el-option
              value=""
              :label="t('admin.redemptionCodes.allProviders')"
            />
            <el-option
              v-for="p in allProviders"
              :key="p.id"
              :value="p.id"
              :label="p.name"
            />
          </el-select>
        </el-col>
      </el-row>

      <!-- 批量操作栏 -->
      <div
        v-if="selectedRows.length > 0"
        class="batch-actions"
      >
        <span style="margin-right: 12px">{{ selectedRows.length }} {{ t('common.selected') }}&nbsp;</span>
        <el-button
          type="primary"
          size="small"
          @click="handleExport"
        >
          {{ t('admin.redemptionCodes.export') }}
        </el-button>
        <el-button
          type="danger"
          size="small"
          @click="handleBatchDelete"
        >
          {{ t('admin.redemptionCodes.batchDelete') }}
        </el-button>
      </div>

      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        :empty-text="t('admin.redemptionCodes.emptyHint')"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="50"
        />
        <el-table-column
          prop="code"
          :label="t('admin.redemptionCodes.colCode')"
          min-width="160"
        />
        <el-table-column
          :label="t('admin.redemptionCodes.colStatus')"
          width="110"
        >
          <template #default="scope">
            <el-tag
              :type="statusTagType(scope.row.status)"
              size="small"
            >
              {{ statusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="providerName"
          :label="t('admin.redemptionCodes.colProvider')"
          width="120"
        />
        <el-table-column
          :label="t('admin.redemptionCodes.colInstanceType')"
          min-width="150"
        >
          <template #default="scope">
            {{ scope.row.instanceType === 'container' ? t('admin.redemptionCodes.container') : t('admin.redemptionCodes.vm') }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.redemptionCodes.colCreationMode')"
          min-width="150"
        >
          <template #default="scope">
            <el-tag
              v-if="scope.row.creationMode === 'copy'"
              type="warning"
              size="small"
            >
              {{ t('admin.redemptionCodes.modeCopy') }}
            </el-tag>
            <el-tag
              v-else
              type="info"
              size="small"
            >
              {{ t('admin.redemptionCodes.modeStandard') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.redemptionCodes.colSpecs')"
          min-width="200"
        >
          <template #default="scope">
            <span v-if="scope.row.cpuName || scope.row.memoryName">
              CPU: {{ scope.row.cpuName || scope.row.cpuId }} / {{ t('admin.redemptionCodes.memory') }}: {{ scope.row.memoryName || scope.row.memoryId }}
              <span v-if="scope.row.diskName || scope.row.diskId"> / {{ t('admin.redemptionCodes.disk') }}: {{ scope.row.diskName || scope.row.diskId }}</span>
              <span v-if="scope.row.bandwidthName || scope.row.bandwidthId"> / {{ t('admin.redemptionCodes.bandwidth') }}: {{ scope.row.bandwidthName || scope.row.bandwidthId }}</span>
            </span>
          </template>
        </el-table-column>
        <el-table-column
          prop="createdByUser"
          :label="t('admin.redemptionCodes.colCreatedBy')"
          min-width="120"
        />
        <el-table-column
          prop="instanceName"
          :label="t('admin.redemptionCodes.colInstanceName')"
          min-width="150"
        >
          <template #default="scope">
            {{ scope.row.instanceName || '-' }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.redemptionCodes.colCreatedAt')"
          width="160"
        >
          <template #default="scope">
            {{ scope.row.createdAt ? new Date(scope.row.createdAt).toLocaleString() : '' }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.redemptionCodes.colRedeemedAt')"
          width="160"
        >
          <template #default="scope">
            {{ scope.row.redeemedAt ? new Date(scope.row.redeemedAt).toLocaleString() : '-' }}
          </template>
        </el-table-column>
        <el-table-column
          prop="remark"
          :label="t('admin.redemptionCodes.colRemark')"
          min-width="120"
        />
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
    </el-card>

    <!-- 批量创建对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="t('admin.redemptionCodes.createDialogTitle')"
      width="min(560px, calc(100vw - 32px))"
      :before-close="handleCreateDialogClose"
    >
      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-width="110px"
      >
        <el-form-item
          :label="t('admin.redemptionCodes.colProvider')"
          prop="providerId"
        >
          <el-select
            v-model="createForm.providerId"
            :placeholder="t('admin.redemptionCodes.providerPlaceholder')"
            style="width: 100%"
            @change="onProviderChange"
          >
            <el-option
              v-for="p in allProviders"
              :key="p.id"
              :value="p.id"
              :label="p.name"
            />
          </el-select>
        </el-form-item>
        <!-- 创建模式（支持复制的容器节点显示） -->
        <el-form-item
          v-if="isLxdIncusProvider"
          :label="t('admin.redemptionCodes.creationMode')"
        >
          <el-radio-group v-model="createForm.creationMode">
            <el-radio value="standard">
              {{ t('admin.redemptionCodes.modeStandard') }}
            </el-radio>
            <el-radio value="copy">
              {{ t('admin.redemptionCodes.modeCopy') }}
            </el-radio>
          </el-radio-group>
          <div style="font-size: 12px; color: #909399; margin-top: 4px;">
            {{ t('admin.redemptionCodes.copyModeTip') }}
          </div>
        </el-form-item>
        <!-- 源容器（仅复制模式显示） -->
        <el-form-item
          v-if="isLxdIncusProvider && createForm.creationMode === 'copy'"
          :label="t('admin.redemptionCodes.sourceContainer')"
          prop="sourceContainer"
        >
          <el-select
            :key="'src-cnt-' + stoppedContainerOptions.length"
            v-model="createForm.sourceContainer"
            :placeholder="stoppedContainersLoading ? t('admin.redemptionCodes.loadingContainers') : t('admin.redemptionCodes.sourceContainerPlaceholder')"
            style="width: 100%"
            :loading="stoppedContainersLoading"
            class="source-container-select"
            popper-class="source-container-popper"
          >
            <el-option
              v-if="stoppedContainers.length === 0 && !stoppedContainersLoading"
              value=""
              :label="t('admin.redemptionCodes.noStoppedContainers')"
              disabled
            />
            <el-option
              v-for="c in stoppedContainerOptions"
              :key="c.name"
              :value="c.name"
              :label="c.label"
            >
              <div
                class="source-container-option"
                :title="c.label"
              >
                {{ c.label }}
              </div>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="createForm.creationMode !== 'copy'"
          :label="t('admin.redemptionCodes.colInstanceType')"
          prop="instanceType"
        >
          <el-select
            v-model="createForm.instanceType"
            :placeholder="t('admin.redemptionCodes.instanceTypePlaceholder')"
            style="width: 100%"
            :disabled="!createForm.providerId"
            @change="onInstanceTypeChange"
          >
            <el-option
              v-if="providerCaps.containerEnabled"
              value="container"
              :label="t('admin.redemptionCodes.container')"
            />
            <el-option
              v-if="providerCaps.vmEnabled"
              value="vm"
              :label="t('admin.redemptionCodes.vm')"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="createForm.creationMode !== 'copy'"
          :label="t('admin.redemptionCodes.colSpecs') + ' - ' + t('admin.redemptionCodes.image')"
          prop="imageId"
        >
          <el-select
            v-model="createForm.imageId"
            :placeholder="t('admin.redemptionCodes.imagePlaceholder')"
            style="width: 100%"
            :disabled="!createForm.instanceType"
          >
            <el-option
              v-for="img in availableImages"
              :key="img.id"
              :value="img.id"
              :label="img.displayName || img.name"
            />
          </el-select>
        </el-form-item>
        <el-row
          v-if="createForm.creationMode !== 'copy'"
          :gutter="12"
        >
          <el-col :span="12">
            <el-form-item
              label="CPU"
              prop="cpuId"
            >
              <el-select
                v-model="createForm.cpuId"
                :placeholder="t('admin.redemptionCodes.cpuPlaceholder')"
                style="width: 100%"
                :disabled="!createForm.instanceType"
              >
                <el-option
                  v-for="spec in cpuSpecs"
                  :key="spec.id"
                  :value="spec.id"
                  :label="spec.name || spec.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="t('admin.redemptionCodes.memory')"
              prop="memoryId"
            >
              <el-select
                v-model="createForm.memoryId"
                :placeholder="t('admin.redemptionCodes.memoryPlaceholder')"
                style="width: 100%"
                :disabled="!createForm.instanceType"
              >
                <el-option
                  v-for="spec in memorySpecs"
                  :key="spec.id"
                  :value="spec.id"
                  :label="spec.name || spec.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row
          v-if="createForm.creationMode !== 'copy'"
          :gutter="12"
        >
          <el-col :span="12">
            <el-form-item
              :label="t('admin.redemptionCodes.disk')"
              prop="diskId"
            >
              <el-select
                v-model="createForm.diskId"
                :placeholder="t('admin.redemptionCodes.diskPlaceholder')"
                style="width: 100%"
                :disabled="!createForm.instanceType"
              >
                <el-option
                  v-for="spec in diskSpecs"
                  :key="spec.id"
                  :value="spec.id"
                  :label="spec.name || spec.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="t('admin.redemptionCodes.bandwidth')"
              prop="bandwidthId"
            >
              <el-select
                v-model="createForm.bandwidthId"
                :placeholder="t('admin.redemptionCodes.bandwidthPlaceholder')"
                style="width: 100%"
                :disabled="!createForm.instanceType"
              >
                <el-option
                  v-for="spec in bandwidthSpecs"
                  :key="spec.id"
                  :value="spec.id"
                  :label="spec.name || spec.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <!-- GPU 直通（支持容器 GPU 的节点，含复制模式） -->
        <el-form-item
          v-if="canConfigureGpuPassthrough"
          :label="t('admin.redemptionCodes.gpuPassthrough')"
        >
          <div class="gpu-config-wrap">
            <el-checkbox
              v-model="createForm.gpuEnabled"
              class="gpu-enable-checkbox"
            >
              {{ t('admin.redemptionCodes.gpuEnabled') }}
            </el-checkbox>
            <div
              v-if="createForm.gpuEnabled"
              class="gpu-config-panel"
            >
              <div class="gpu-actions-row">
                <span class="gpu-actions-label">{{ t('admin.redemptionCodes.gpuDeviceIds') }}:</span>
                <el-button
                  size="small"
                  :loading="gpuDetecting"
                  @click="detectGpus"
                >
                  {{ gpuDetecting ? t('admin.redemptionCodes.gpuDetecting') : t('admin.redemptionCodes.gpuDetect') }}
                </el-button>
              </div>
              <!-- 检测结果列表 -->
              <div
                v-if="detectedGpus.length > 0"
                class="gpu-options-wrap"
              >
                <el-checkbox-group
                  v-model="selectedGpuIndices"
                  class="gpu-options"
                >
                  <el-checkbox
                    v-for="(gpu, idx) in detectedGpus"
                    :key="idx"
                    :value="idx"
                    :label="idx"
                    class="gpu-option-item"
                  >
                    {{ t('admin.redemptionCodes.gpuDeviceIdx', { idx }) }} — {{ gpu.name || gpu.vendor || '' }}
                  </el-checkbox>
                </el-checkbox-group>
                <div class="gpu-batch-actions">
                  <el-button
                    size="small"
                    text
                    @click="selectAllGpus"
                  >
                    {{ t('admin.redemptionCodes.gpuSelectAll') }}
                  </el-button>
                  <el-button
                    size="small"
                    text
                    @click="deselectAllGpus"
                  >
                    {{ t('admin.redemptionCodes.gpuDeselectAll') }}
                  </el-button>
                </div>
              </div>
              <div
                v-else-if="gpuChecked && !gpuDetecting"
                class="gpu-empty-tip"
              >
                {{ t('admin.redemptionCodes.gpuNoneFound') }}
              </div>
            </div>
          </div>
        </el-form-item>
        <el-form-item
          :label="t('admin.redemptionCodes.countLabel')"
          prop="count"
        >
          <el-input-number
            v-model="createForm.count"
            :min="1"
            :max="100"
            :controls="false"
            style="width: 140px"
          />
        </el-form-item>
        <el-form-item
          :label="t('admin.redemptionCodes.remarkLabel')"
          prop="remark"
        >
          <el-input
            v-model="createForm.remark"
            type="textarea"
            :rows="2"
            :placeholder="t('admin.redemptionCodes.remarkPlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelCreate">{{ t('common.cancel') }}</el-button>
          <el-button
            type="primary"
            :loading="createLoading"
            @click="submitCreate"
          >
            {{ t('common.create') }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 导出对话框 -->
    <el-dialog
      v-model="showExportDialog"
      :title="t('admin.redemptionCodes.exportDialogTitle')"
      width="600px"
    >
      <!-- 步骤1：选择导出字段 -->
      <div v-if="!exportResult">
        <p style="margin-bottom: 12px; font-weight: 600;">
          {{ t('admin.redemptionCodes.selectExportFields') }}
        </p>
        <el-checkbox-group
          v-model="exportFields"
          style="margin-bottom: 16px;"
        >
          <el-checkbox
            v-for="field in allExportFields"
            :key="field.value"
            :value="field.value"
            style="margin-bottom: 6px; width: 180px;"
          >
            {{ field.label }}
          </el-checkbox>
        </el-checkbox-group>
        <div style="margin-bottom: 12px;">
          <el-button
            size="small"
            @click="exportFields = allExportFields.map(f => f.value)"
          >
            {{ t('admin.redemptionCodes.selectAll') }}
          </el-button>
          <el-button
            size="small"
            @click="exportFields = []"
          >
            {{ t('admin.redemptionCodes.deselectAll') }}
          </el-button>
        </div>
      </div>
      <!-- 步骤2：显示导出结果 -->
      <div v-else>
        <p style="margin-bottom: 8px">
          {{ t('admin.redemptionCodes.exportedCodes') }}
        </p>
        <el-input
          v-model="exportedCodesText"
          type="textarea"
          :rows="16"
          readonly
        />
      </div>
      <template #footer>
        <span
          v-if="!exportResult"
          class="dialog-footer"
        >
          <el-button @click="showExportDialog = false">{{ t('common.cancel') }}</el-button>
          <el-button
            type="primary"
            :loading="exportLoading"
            :disabled="exportFields.length === 0"
            @click="doExport"
          >
            {{ t('admin.redemptionCodes.export') }}
          </el-button>
        </span>
        <span
          v-else
          class="dialog-footer"
        >
          <el-button @click="resetExportDialog">{{ t('common.back') }}</el-button>
          <el-button @click="showExportDialog = false; exportResult = null">{{ t('common.close') }}</el-button>
          <el-button
            type="primary"
            @click="copyExportedCodes"
          >
            {{ t('admin.redemptionCodes.copyAll') }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import useRedemptionCodes from './useRedemptionCodes'

const { t } = useI18n()

const {
  loading, tableData, total, currentPage, pageSize,
  filterCode, filterStatus, filterProvider,
  selectedRows, allProviders,
  showCreateDialog, createLoading, createFormRef,
  createForm, createRules,
  providerCaps, availableImages,
  cpuSpecs, memorySpecs, diskSpecs, bandwidthSpecs,
  stoppedContainers, stoppedContainerOptions, stoppedContainersLoading,
  isLxdIncusProvider, canConfigureGpuPassthrough,
  gpuDetecting, gpuChecked, detectedGpus, selectedGpuIndices,
  showExportDialog, exportedCodesText, exportResult, exportLoading, exportFields,
  allExportFields,
  statusTagType, statusLabel,
  handleFilterChange, handleSelectionChange,
  openCreateDialog, cancelCreate, handleCreateDialogClose,
  onProviderChange, onInstanceTypeChange,
  detectGpus, selectAllGpus, deselectAllGpus,
  submitCreate,
  handleExport, resetExportDialog, doExport, copyExportedCodes,
  handleBatchDelete,
  handleSizeChange, handleCurrentChange
} = useRedemptionCodes()
</script>


<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.card-header > span {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color-primary);
}
.filter-bar {
  margin-bottom: 16px;
}
.batch-actions {
  margin-bottom: 12px;
  padding: 10px 12px;
  background-color: var(--neutral-bg);
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.gpu-config-wrap {
  width: 100%;
}

.gpu-enable-checkbox {
  margin-bottom: 8px;
}

.gpu-config-panel {
  margin-top: 8px;
}

.gpu-actions-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.gpu-actions-label {
  font-size: 13px;
  color: #606266;
}

.gpu-options-wrap {
  margin-bottom: 8px;
  max-height: 220px;
  overflow: auto;
  padding-right: 4px;
}

.gpu-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

:deep(.gpu-option-item.el-checkbox) {
  display: flex;
  align-items: flex-start;
  margin-right: 0;
  min-width: 0;
  height: auto;
}

:deep(.gpu-option-item .el-checkbox__label) {
  white-space: normal;
  line-height: 1.4;
  word-break: break-word;
  overflow-wrap: anywhere;
}

.gpu-batch-actions {
  margin-top: 4px;
  font-size: 12px;
  color: #909399;
}

.gpu-empty-tip {
  font-size: 12px;
  color: #909399;
}

:deep(.source-container-select .el-select__selected-item) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

:deep(.source-container-popper .el-select-dropdown__item) {
  height: auto;
  min-height: 34px;
  line-height: 1.35;
  padding-top: 7px;
  padding-bottom: 7px;
}

.source-container-option {
  white-space: normal;
  word-break: break-word;
  overflow-wrap: anywhere;
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
