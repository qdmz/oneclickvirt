<template>
  <div class="invite-codes-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.inviteCodes.title') }}</span>
          <div>
            <el-button
              type="success"
              @click="showCreateDialog = true"
            >
              {{ $t('admin.inviteCodes.createCustomCode') }}
            </el-button>
            <el-button
              type="primary"
              @click="showGenerateDialog = true"
            >
              {{ $t('admin.inviteCodes.batchGenerate') }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-form :inline="true">
          <el-form-item :label="$t('admin.inviteCodes.usageStatus')">
            <el-select
              v-model="filterForm.isUsed"
              :placeholder="$t('common.all')"
              clearable
              style="width: 120px"
              @change="handleFilterChange"
            >
              <el-option
                :label="$t('common.all')"
                :value="null"
              />
              <el-option
                :label="$t('admin.inviteCodes.unused')"
                :value="false"
              />
              <el-option
                :label="$t('admin.inviteCodes.used')"
                :value="true"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('common.status')">
            <el-select
              v-model="filterForm.status"
              :placeholder="$t('common.all')"
              clearable
              style="width: 120px"
              @change="handleFilterChange"
            >
              <el-option
                :label="$t('common.all')"
                :value="0"
              />
              <el-option
                :label="$t('admin.inviteCodes.available')"
                :value="1"
              />
            </el-select>
          </el-form-item>
        </el-form>
      </div>

      <!-- 批量操作按钮 -->
      <div
        v-if="selectedCodes.length > 0"
        class="batch-actions"
      >
        <el-button
          type="primary"
          @click="handleBatchExport"
        >
          {{ $t('admin.inviteCodes.exportSelected') }} ({{ selectedCodes.length }})
        </el-button>
        <el-button
          type="danger"
          @click="handleBatchDelete"
        >
          {{ $t('admin.inviteCodes.deleteSelected') }} ({{ selectedCodes.length }})
        </el-button>
      </div>
      
      <el-table
        v-loading="loading"
        :data="inviteCodes"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column
          prop="id"
          label="ID"
          width="60"
        />
        <el-table-column
          prop="code"
          :label="$t('admin.inviteCodes.code')"
          min-width="130"
        />
        <el-table-column
          prop="maxUses"
          :label="$t('admin.inviteCodes.maxUses')"
          min-width="130"
        >
          <template #default="scope">
            {{ scope.row.maxUses === 0 ? $t('admin.inviteCodes.unlimited') : scope.row.maxUses }}
          </template>
        </el-table-column>
        <el-table-column
          prop="usedCount"
          :label="$t('admin.inviteCodes.usedCount')"
          width="120"
        />
        <el-table-column
          prop="status"
          :label="$t('common.status')"
          width="100"
        >
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? $t('admin.inviteCodes.available') : $t('admin.inviteCodes.expired') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="expiresAt"
          :label="$t('admin.inviteCodes.expiryDate')"
          width="160"
        >
          <template #default="scope">
            {{ scope.row.expiresAt ? new Date(scope.row.expiresAt).toLocaleString() : $t('admin.inviteCodes.neverExpires') }}
          </template>
        </el-table-column>
        <el-table-column
          prop="createdAt"
          :label="$t('common.createTime')"
          width="160"
        >
          <template #default="scope">
            {{ new Date(scope.row.createdAt).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.actions')"
          width="120"
        >
          <template #default="scope">
            <el-button
              size="small"
              type="danger"
              @click="deleteCode(scope.row.id)"
            >
              {{ $t('common.delete') }}
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
    </el-card>

    <!-- 创建自定义邀请码对话框 -->
    <el-dialog 
      v-model="showCreateDialog" 
      :title="$t('admin.inviteCodes.createCustomCode')" 
      width="500px"
      :before-close="handleCreateDialogClose"
    >
      <el-form 
        ref="createFormRef" 
        :model="createForm" 
        :rules="createRules" 
        label-width="120px"
      >
        <el-form-item
          :label="$t('admin.inviteCodes.code')"
          prop="code"
        >
          <el-input 
            v-model="createForm.code" 
            :placeholder="$t('admin.inviteCodes.codeInputPlaceholder')"
            maxlength="50"
            show-word-limit
          />
          <div class="form-tip">
            {{ $t('admin.inviteCodes.codeFormatTip') }}
          </div>
        </el-form-item>
        <el-form-item
          :label="$t('admin.inviteCodes.maxUses')"
          prop="maxUses"
        >
          <el-input-number
            v-model="createForm.maxUses"
            :min="0"
            :controls="false"
          />
          <div class="form-tip">
            {{ $t('admin.inviteCodes.maxUsesTip') }}
          </div>
        </el-form-item>
        <el-form-item
          :label="$t('admin.inviteCodes.expiryDate')"
          prop="expiresAt"
        >
          <el-date-picker
            v-model="createForm.expiresAt"
            type="datetime"
            :placeholder="$t('admin.inviteCodes.selectExpiryDate')"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
          <div class="form-tip">
            {{ $t('admin.inviteCodes.expiryDateTip') }}
          </div>
        </el-form-item>
        <el-form-item
          :label="$t('common.description')"
          prop="description"
        >
          <el-input 
            v-model="createForm.description" 
            type="textarea" 
            :rows="3"
            :placeholder="$t('admin.inviteCodes.descriptionPlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelCreate">{{ $t('common.cancel') }}</el-button>
          <el-button
            type="primary"
            :loading="createLoading"
            @click="submitCreate"
          >{{ $t('common.create') }}</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 生成邀请码对话框 -->
    <el-dialog 
      v-model="showGenerateDialog" 
      :title="$t('admin.inviteCodes.batchGenerate')" 
      width="500px"
    >
      <el-form 
        ref="generateFormRef" 
        :model="generateForm" 
        :rules="generateRules" 
        label-width="120px"
      >
        <el-form-item
          :label="$t('admin.inviteCodes.generateCount')"
          prop="count"
        >
          <el-input-number
            v-model="generateForm.count"
            :min="1"
            :max="100"
            :controls="false"
          />
        </el-form-item>
        <el-form-item
          :label="$t('admin.inviteCodes.maxUses')"
          prop="maxUses"
        >
          <el-input-number
            v-model="generateForm.maxUses"
            :min="0"
            :controls="false"
          />
          <div class="form-tip">
            {{ $t('admin.inviteCodes.maxUsesTip') }}
          </div>
        </el-form-item>
        <el-form-item
          :label="$t('admin.inviteCodes.expiryDate')"
          prop="expiresAt"
        >
          <el-date-picker
            v-model="generateForm.expiresAt"
            type="datetime"
            :placeholder="$t('admin.inviteCodes.selectExpiryDate')"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
          <div class="form-tip">
            {{ $t('admin.inviteCodes.expiryDateTip') }}
          </div>
        </el-form-item>
        <el-form-item
          :label="$t('common.description')"
          prop="description"
        >
          <el-input 
            v-model="generateForm.description" 
            type="textarea" 
            :rows="3"
            :placeholder="$t('admin.inviteCodes.descriptionPlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelGenerate">{{ $t('common.cancel') }}</el-button>
          <el-button
            type="primary"
            :loading="generateLoading"
            @click="submitGenerate"
          >{{ $t('admin.inviteCodes.generate') }}</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 导出邀请码对话框 -->
    <el-dialog
      v-model="showExportDialog"
      :title="$t('admin.inviteCodes.exportCodes')"
      width="600px"
    >
      <div class="export-content">
        <el-input
          v-model="exportedCodes"
          type="textarea"
          :rows="15"
          readonly
        />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showExportDialog = false">{{ $t('common.close') }}</el-button>
          <el-button
            type="primary"
            @click="copyExportedCodes"
          >{{ $t('common.copy') }}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { useInviteCodeManagement } from './composables/useInviteCodeManagement.js'

const {
  inviteCodes,
  loading,
  showCreateDialog,
  showGenerateDialog,
  showExportDialog,
  createLoading,
  generateLoading,
  createFormRef,
  generateFormRef,
  selectedCodes,
  exportedCodes,
  filterForm,
  currentPage,
  pageSize,
  total,
  createForm,
  createRules,
  generateForm,
  generateRules,
  handleFilterChange,
  handleSelectionChange,
  handleBatchExport,
  handleBatchDelete,
  copyExportedCodes,
  cancelCreate,
  handleCreateDialogClose,
  submitCreate,
  cancelGenerate,
  submitGenerate,
  deleteCode,
  handleSizeChange,
  handleCurrentChange
} = useInviteCodeManagement()
</script>

<style scoped>
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

.filter-bar {
  margin-bottom: 20px;
}

.batch-actions {
  margin-bottom: 15px;
  padding: 10px;
  background-color: var(--neutral-bg);
  border-radius: 4px;
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

.form-tip {
  font-size: 12px;
  color: var(--text-color-secondary);
  margin-top: 4px;
}

.export-content {
  margin: 20px 0;
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
