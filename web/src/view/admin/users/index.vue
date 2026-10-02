<template>
  <div class="users-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.users.title') }}</span>
          <div class="header-actions">
            <el-button
              type="primary"
              @click="handleAddUser"
            >
              {{ $t('admin.users.addUser') }}
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 搜索和批量操作 -->
      <div class="toolbar">
        <div class="search-section">
          <el-input
            v-model="searchUsername"
            :placeholder="$t('admin.users.searchByUsername')"
            style="width: 200px;"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select
            v-model="searchStatus"
            :placeholder="$t('admin.users.selectStatus')"
            style="width: 150px; margin-left: 10px;"
            clearable
          >
            <el-option
              :label="$t('admin.users.all')"
              :value="null"
            />
            <el-option
              :label="$t('admin.users.active')"
              :value="1"
            />
            <el-option
              :label="$t('admin.users.disabled')"
              :value="0"
            />
          </el-select>
          <el-select
            v-model="searchUserType"
            :placeholder="$t('admin.users.selectUserType')"
            style="width: 180px; margin-left: 10px;"
            clearable
          >
            <el-option
              :label="$t('admin.users.all')"
              value=""
            />
            <el-option
              :label="$t('admin.users.normalUser')"
              value="user"
            />
            <el-option
              :label="$t('admin.users.normalAdmin')"
              value="normal_admin"
            />
            <el-option
              :label="$t('admin.users.adminUser')"
              value="admin"
            />
          </el-select>
          <el-button 
            type="primary" 
            style="margin-left: 10px;"
            @click="handleSearch"
          >
            {{ $t('admin.users.query') }}
          </el-button>
          <el-button 
            type="default" 
            style="margin-left: 10px;"
            @click="resetFilters"
          >
            {{ $t('admin.users.resetFilters') }}
          </el-button>
        </div>
        
        <div
          v-if="multipleSelection.length > 0"
          class="batch-actions"
        >
          <span class="selection-info">{{ $t('admin.users.selected') }} {{ multipleSelection.length }} {{ $t('admin.users.users') }}</span>
          <el-button
            size="small"
            type="danger"
            @click="handleBatchDelete"
          >
            {{ $t('admin.users.batchDelete') }}
          </el-button>
          <el-button
            size="small"
            type="warning"
            @click="handleBatchEnable"
          >
            {{ $t('admin.users.batchEnable') }}
          </el-button>
          <el-button
            size="small"
            type="info"
            @click="handleBatchDisable"
          >
            {{ $t('admin.users.batchDisable') }}
          </el-button>
          <el-dropdown @command="handleBatchLevelCommand">
            <el-button
              size="small"
              type="primary"
            >
              {{ $t('admin.users.batchSetLevel') }}<el-icon class="el-icon--right">
                <arrow-down />
              </el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item
                  v-for="level in availableLevels"
                  :key="level"
                  :command="level"
                >
                  {{ $t('admin.users.setToLevel', { level }) }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
      
      <UsersTable
        :users="users"
        :loading="loading"
        :available-levels="availableLevels"
        @selection-change="handleSelectionChange"
        @edit="editUser"
        @set-user-level="handleSetUserLevel"
        @set-expiry="handleSetExpiry"
        @toggle-status="handleToggleUserStatus"
        @reset-password="handleResetPassword"
        @login-as="handleLoginAsUser"
        @adjust-balance="handleAdjustBalance"
      />

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

    <AddEditUserDialog
      ref="addUserFormRef"
      :visible="showAddDialog"
      :is-editing="isEditing"
      :add-user-form="addUserForm"
      :add-user-rules="addUserRules"
      :available-levels="availableLevels"
      :loading="addUserLoading"
      @update:visible="showAddDialog = $event"
      @cancel="cancelAddUser"
      @submit="submitAddUser"
    />

    <ResetPasswordDialog
      :visible="showResetPasswordDialog"
      :reset-password-form="resetPasswordForm"
      :generated-password="generatedPassword"
      :loading="resetPasswordLoading"
      @update:visible="showResetPasswordDialog = $event"
      @cancel="cancelResetPassword"
      @confirm="confirmResetPassword"
      @copy-password="copyPassword"
    />

    <SetExpiryDialog
      :visible="showSetExpiryDialog"
      :freeze-form="freezeForm"
      :loading="freezeLoading"
      @update:visible="showSetExpiryDialog = $event"
      @confirm="confirmSetExpiry"
    />

    <AdjustBalanceDialog
      :visible="showAdjustBalanceDialog"
      :form="balanceForm"
      :loading="balanceLoading"
      @update:visible="showAdjustBalanceDialog = $event"
      @confirm="confirmAdjustBalance"
    />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { Search, ArrowDown } from '@element-plus/icons-vue'
import { useUserManagement } from './composables/useUserManagement'
import UsersTable from './components/UsersTable.vue'
import AddEditUserDialog from './components/AddEditUserDialog.vue'
import ResetPasswordDialog from './components/ResetPasswordDialog.vue'
import SetExpiryDialog from './components/SetExpiryDialog.vue'
import AdjustBalanceDialog from './components/AdjustBalanceDialog.vue'

const {
  users, loading, showAddDialog, addUserLoading, addUserFormRef, isEditing,
  showResetPasswordDialog, resetPasswordForm, resetPasswordLoading, generatedPassword,
  showSetExpiryDialog, freezeLoading, freezeForm,
  showAdjustBalanceDialog, balanceLoading, balanceForm,
  searchUsername, searchStatus, searchUserType,
  multipleSelection, currentPage, pageSize, total, availableLevels,
  addUserForm, addUserRules,
  loadUsers, loadLevelOptions, handleSearch, resetFilters,
  handleSelectionChange, handleBatchDelete, handleBatchEnable, handleBatchDisable,
  handleBatchLevelCommand, handleSetUserLevel,
  getLevelTagType, getUserTypeLabel, getUserTypeTagType,
  handleAddUser, editUser, cancelAddUser, submitAddUser,
  handleToggleUserStatus,
  handleResetPassword, confirmResetPassword, cancelResetPassword,
  handleLoginAsUser, copyPassword,
  handleSetExpiry, confirmSetExpiry,
  handleAdjustBalance, confirmAdjustBalance,
  formatDateTime, isExpired,
  handleSizeChange, handleCurrentChange,
  t
} = useUserManagement()

onMounted(() => {
  loadLevelOptions()
  loadUsers()
})
</script>

<style scoped lang="scss">
.users-container {
  .el-card {
    :deep(.el-card__header) {
      padding: 20px 24px;
      border-bottom: 1px solid #ebeef5;
    }
    
    :deep(.el-card__body) {
      padding: 24px;
    }
  }
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
  
  .header-actions {
    .el-button {
      padding: 10px 20px;
    }
  }
}

.users-table {
  width: 100%;
  
  .action-buttons {
    display: flex;
    gap: 12px;
    justify-content: center;
    align-items: center;
    flex-wrap: wrap;
    padding: 4px 0;
    
    .el-button {
      margin: 0 !important;
      padding: 8px 16px;
    }
    
    .el-dropdown {
      .el-button {
        margin: 0 !important;
        padding: 8px 16px;
      }
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

.toolbar {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  
  .el-button {
    padding: 10px 20px;
  }
}

.search-section {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  
  .el-button {
    padding: 10px 20px;
  }
}

.batch-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background-color: var(--neutral-bg);
  border-radius: 4px;
  
  .el-button {
    padding: 8px 16px;
  }
}

.selection-info {
  color: #16a34a;
  font-weight: 500;
}

.role-tag {
  margin-right: 5px;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.password-hint {
  margin-top: 5px;
  font-size: 12px;
  line-height: 1.4;
  color: var(--text-color-secondary);
}

.dialog-footer {
  text-align: right;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  
  .el-button {
    padding: 10px 24px;
    margin: 0 !important;
  }
}

:deep(.el-dialog) {
  .el-dialog__body {
    padding: 24px 24px 10px;
  }
  
  .el-form {
    .el-form-item {
      margin-bottom: 24px;
    }
    
    .el-row {
      margin-bottom: 8px;
    }
    
    .el-input {
      .el-input__inner {
        padding: 8px 12px;
      }
    }
    
    .el-select {
      .el-input__inner {
        padding: 8px 12px;
      }
    }
  }
  
  .el-input-group__append {
    .el-button {
      padding: 8px 16px;
    }
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
