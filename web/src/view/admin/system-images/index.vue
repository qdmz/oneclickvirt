<template>
  <div class="system-images-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.systemImages.title') }}</span>
          <div class="header-buttons">
            <el-button
              type="success"
              :loading="syncing"
              @click="handleSync"
            >
              {{ $t('admin.systemImages.syncImages') }}
            </el-button>
            <el-button @click="handleReset">
              {{ $t('common.reset') }}
            </el-button>
            <el-button
              type="primary"
              @click="handleSearch"
            >
              {{ $t('common.search') }}
            </el-button>
            <el-button
              type="primary"
              @click="handleCreate"
            >
              {{ $t('admin.systemImages.addImage') }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索过滤 -->
      <div class="filter-container">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-input
              v-model="searchForm.search"
              :placeholder="$t('admin.systemImages.searchPlaceholder')"
              clearable
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </el-col>
          <el-col :span="4">
            <el-select
              v-model="searchForm.providerType"
              :placeholder="$t('admin.systemImages.providerType')"
              clearable
              style="width: 100%;"
            >
              <el-option
                v-for="option in PROVIDER_TYPE_OPTIONS"
                :key="option.value"
                :label="$t(option.labelKey)"
                :value="option.value"
              />
            </el-select>
          </el-col>
          <el-col :span="3">
            <el-select
              v-model="searchForm.instanceType"
              :placeholder="$t('admin.systemImages.instanceType')"
              clearable
              style="width: 100%;"
            >
              <el-option
                :label="$t('admin.systemImages.vm')"
                value="vm"
              />
              <el-option
                :label="$t('admin.systemImages.container')"
                value="container"
              />
            </el-select>
          </el-col>
          <el-col :span="3">
            <el-select
              v-model="searchForm.architecture"
              :placeholder="$t('admin.systemImages.architecture')"
              clearable
              style="width: 100%;"
            >
              <el-option
                label="amd64"
                value="amd64"
              />
              <el-option
                label="arm64"
                value="arm64"
              />
              <el-option
                label="s390x"
                value="s390x"
              />
            </el-select>
          </el-col>
          <el-col :span="3">
            <el-select
              v-model="searchForm.osType"
              :placeholder="$t('admin.systemImages.osType')"
              clearable
              style="width: 100%;"
            >
              <el-option-group
                v-for="(osList, category) in groupedOperatingSystems"
                :key="category"
                :label="category"
              >
                <el-option
                  v-for="os in osList"
                  :key="os.name"
                  :label="os.displayName"
                  :value="os.name"
                />
              </el-option-group>
            </el-select>
          </el-col>
          <el-col :span="3">
            <el-select
              v-model="searchForm.status"
              :placeholder="$t('common.status')"
              clearable
              style="width: 100%;"
            >
              <el-option
                :label="$t('admin.systemImages.active')"
                value="active"
              />
              <el-option
                :label="$t('admin.systemImages.inactive')"
                value="inactive"
              />
            </el-select>
          </el-col>
        </el-row>
      </div>

      <!-- 批量操作 -->
      <div
        v-if="selectedRows.length > 0"
        class="batch-actions"
      >
        <el-alert
          :title="$t('admin.systemImages.selectedCount', { count: selectedRows.length })"
          type="info"
          show-icon
          :closable="false"
        >
          <template #default>
            <el-button
              type="success"
              size="small"
              @click="handleBatchStatus('active')"
            >
              {{ $t('admin.systemImages.batchActivate') }}
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="handleBatchStatus('inactive')"
            >
              {{ $t('admin.systemImages.batchDisable') }}
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleBatchDelete"
            >
              {{ $t('admin.systemImages.batchDelete') }}
            </el-button>
          </template>
        </el-alert>
      </div>

      <!-- 数据表格 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        class="system-images-table"
        :cell-style="{ padding: '12px 0' }"
        :header-cell-style="{ background: '#f5f7fa', padding: '14px 0', fontWeight: '600' }"
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
          prop="name"
          :label="$t('admin.systemImages.imageName')"
          min-width="140"
          show-overflow-tooltip
        />
        <el-table-column
          :label="$t('admin.systemImages.providerType')"
          min-width="190"
          align="center"
        >
          <template #default="scope">
            <el-tag :type="getProviderTypeColor(scope.row.providerType)">
              {{ getProviderTypeName(scope.row.providerType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('admin.systemImages.instanceType')"
          width="110"
          align="center"
        >
          <template #default="scope">
            <el-tag :type="scope.row.instanceType === 'vm' ? 'primary' : 'success'">
              {{ scope.row.instanceType === 'vm' ? $t('admin.systemImages.vm') : $t('admin.systemImages.container') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="architecture"
          :label="$t('admin.systemImages.architecture')"
          min-width="140"
          align="center"
          show-overflow-tooltip
        />
        <el-table-column
          :label="$t('admin.systemImages.osType')"
          min-width="170"
          show-overflow-tooltip
        >
          <template #default="scope">
            <span style="display: inline-flex; align-items: center; gap: 6px;">
              <OsIcon
                :name="scope.row.osType || scope.row.name"
                :size="20"
              />
              {{ getDisplayName(scope.row.osType) || scope.row.osType || '-' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          prop="osVersion"
          :label="$t('admin.systemImages.version')"
          width="120"
          show-overflow-tooltip
        />
        <el-table-column
          label="URL"
          min-width="200"
          show-overflow-tooltip
        >
          <template #default="scope">
            <span class="url-text">{{ scope.row.url }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('admin.systemImages.size')"
          width="100"
          align="center"
        >
          <template #default="scope">
            {{ formatFileSize(scope.row.size) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.status')"
          width="100"
          align="center"
        >
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'">
              {{ scope.row.status === 'active' ? $t('admin.systemImages.active') : $t('admin.systemImages.inactive') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.createTime')"
          width="180"
          align="center"
        >
          <template #default="scope">
            {{ formatDateTime(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.actions')"
          width="240"
          fixed="right"
          align="center"
        >
          <template #default="scope">
            <div class="action-buttons">
              <el-button
                type="primary"
                size="small"
                @click="handleEdit(scope.row)"
              >
                {{ $t('common.edit') }}
              </el-button>
              <el-button
                :type="scope.row.status === 'active' ? 'warning' : 'success'"
                size="small"
                @click="handleToggleStatus(scope.row)"
              >
                {{ scope.row.status === 'active' ? $t('common.disable') : $t('admin.systemImages.activate') }}
              </el-button>
              <el-button
                type="danger"
                size="small"
                @click="handleDelete(scope.row)"
              >
                {{ $t('common.delete') }}
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
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

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="800px"
      :before-close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              :label="$t('admin.systemImages.imageName')"
              prop="name"
            >
              <el-input
                v-model="form.name"
                :placeholder="$t('admin.systemImages.imageNamePlaceholder')"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="$t('admin.systemImages.providerType')"
              prop="providerType"
            >
              <el-select
                v-model="form.providerType"
                :placeholder="$t('admin.systemImages.selectProviderType')"
                @change="handleProviderTypeChange"
              >
                <el-option
                  v-for="option in PROVIDER_TYPE_OPTIONS"
                  :key="option.value"
                  :label="$t(option.labelKey)"
                  :value="option.value"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              :label="$t('admin.systemImages.instanceType')"
              prop="instanceType"
            >
              <el-select
                v-model="form.instanceType"
                :placeholder="$t('admin.systemImages.selectInstanceType')"
                @change="handleInstanceTypeChange"
              >
                <el-option
                  :label="$t('admin.systemImages.vm')"
                  value="vm"
                />
                <el-option
                  :label="$t('admin.systemImages.container')"
                  value="container"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="$t('admin.systemImages.architecture')"
              prop="architecture"
            >
              <el-select
                v-model="form.architecture"
                :placeholder="$t('admin.systemImages.selectArchitecture')"
              >
                <el-option
                  label="amd64"
                  value="amd64"
                />
                <el-option
                  label="arm64"
                  value="arm64"
                />
                <el-option
                  label="s390x"
                  value="s390x"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item
          :label="$t('admin.systemImages.imageUrl')"
          prop="url"
        >
          <el-input
            v-model="form.url"
            :placeholder="$t('admin.systemImages.imageUrlPlaceholder')"
          />
          <div class="form-hint">
            <template v-if="getUrlHint()">
              {{ getUrlHint() }}
            </template>
          </div>
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              :label="$t('admin.systemImages.osType')"
              prop="osType"
            >
              <el-select 
                v-model="form.osType" 
                :placeholder="$t('admin.systemImages.selectOsType')"
                filterable
                @change="handleOsTypeChange"
              >
                <el-option-group
                  v-for="(osList, category) in groupedOperatingSystems"
                  :key="category"
                  :label="category"
                >
                  <el-option
                    v-for="os in osList"
                    :key="os.name"
                    :label="os.displayName"
                    :value="os.name"
                  />
                </el-option-group>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="$t('admin.systemImages.osVersion')"
              prop="osVersion"
            >
              <el-input 
                v-model="form.osVersion" 
                :placeholder="$t('admin.systemImages.selectOsVersion')"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('admin.systemImages.fileSize')">
              <el-input
                v-model.number="form.size"
                type="number"
                :placeholder="$t('admin.systemImages.optional')"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('admin.systemImages.checksum')">
              <el-input
                v-model="form.checksum"
                :placeholder="$t('admin.systemImages.checksumPlaceholder')"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              :label="$t('admin.systemImages.minMemoryMB')"
              prop="minMemoryMB"
            >
              <el-input
                v-model.number="form.minMemoryMB"
                type="number"
                :placeholder="$t('admin.systemImages.minMemoryPlaceholder')"
              />
              <div class="form-hint">
                {{ $t('admin.systemImages.minMemoryHint') }}
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="$t('admin.systemImages.minDiskMB')"
              prop="minDiskMB"
            >
              <el-input
                v-model.number="form.minDiskMB"
                type="number"
                :placeholder="$t('admin.systemImages.minDiskPlaceholder')"
              />
              <div class="form-hint">
                {{ $t('admin.systemImages.minDiskHint') }}
              </div>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('admin.systemImages.useCdn')">
          <el-switch
            v-model="form.useCdn"
            :active-text="$t('admin.systemImages.useCdnActive')"
            :inactive-text="$t('admin.systemImages.useCdnInactive')"
          />
          <div class="form-hint">
            {{ $t('admin.systemImages.useCdnHint') }}
          </div>
        </el-form-item>
        <el-form-item :label="$t('admin.systemImages.tags')">
          <el-input
            v-model="form.tags"
            :placeholder="$t('admin.systemImages.tagsPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="$t('common.description')">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            :placeholder="$t('admin.systemImages.descriptionPlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleDialogClose">{{ $t('common.cancel') }}</el-button>
          <el-button
            type="primary"
            :loading="submitting"
            @click="handleSubmit"
          >
            {{ $t('common.confirm') }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>


<script setup>
import { onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import OsIcon from '@/components/OsIcon.vue'
import { useSystemImageManagement } from './composables/useSystemImageManagement'
import { PROVIDER_TYPE_OPTIONS } from '@/utils/providerTypes'

const {
  loading, submitting, syncing, dialogVisible, selectedRows, tableData,
  searchForm, pagination, form, formRef, isEdit, editId,
  groupedOperatingSystems, dialogTitle, rules,
  fetchData, handleSearch, handleReset, handleSync, handleSelectionChange,
  handleCreate, handleEdit, handleSubmit, handleDelete,
  handleToggleStatus, handleBatchDelete, handleBatchStatus,
  handleSizeChange, handleCurrentChange, handleDialogClose,
  handleProviderTypeChange, handleInstanceTypeChange, handleOsTypeChange,
  getUrlHint, getProviderTypeName, getProviderTypeColor,
  truncateUrl, formatFileSize, formatDateTime,
  getDisplayName,
  t
} = useSystemImageManagement()

onMounted(() => {
  fetchData()
})
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

.header-buttons {
  display: flex;
  gap: 8px;
  align-items: center;
}

.system-images-table {
  width: 100%;
  
  .action-buttons {
    display: flex;
    gap: 10px;
    justify-content: center;
    align-items: center;
    flex-wrap: wrap;
    padding: 4px 0;
    
    .el-button {
      margin: 0 !important;
    }
  }
  
  :deep(.el-table__cell:not(.el-table-fixed-column--right)) {
    .cell {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

.filter-container {
  margin-bottom: 20px;
}

.batch-actions {
  margin-bottom: 16px;
}

.pagination-container {
  margin-top: 20px;
  text-align: center;
}

.url-text {
  cursor: pointer;
  color: #16a34a;
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.is-default {
  color: #f56c6c;
}

.form-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.dialog-footer {
  text-align: right;
}

:deep(.el-table) {
  margin-bottom: 0;
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
