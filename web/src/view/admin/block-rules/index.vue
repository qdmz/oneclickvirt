<template>
  <div class="block-rules-container">
    <!-- Rules Tab -->
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('admin.blockRules.rules') }}</span>
          <div>
            <el-button
              v-if="isSuperAdmin"
              type="primary"
              @click="handleCreateRule"
            >
              {{ t('admin.blockRules.addRule') }}
            </el-button>
            <el-button
              type="success"
              :disabled="selectedRules.length === 0"
              @click="showApplyDialog = true"
            >
              {{ t('admin.blockRules.applyRules') }}
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loadingRules"
        :data="rules"
        stripe
        @selection-change="handleRuleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column
          prop="name"
          :label="t('admin.blockRules.ruleName')"
          min-width="150"
        />
        <el-table-column
          prop="category"
          :label="t('admin.blockRules.category')"
          width="120"
        >
          <template #default="{ row }">
            <el-tag
              :type="categoryTagType(row.category)"
              size="small"
            >
              {{ t(`admin.blockRules.categories.${row.category}`) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="description"
          :label="t('admin.blockRules.description')"
          min-width="200"
          show-overflow-tooltip
        />
        <el-table-column
          :label="t('admin.blockRules.strings')"
          min-width="150"
        >
          <template #default="{ row }">
            {{ parseStrings(row.strings).length }} {{ t('admin.blockRules.strings') }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.blockRules.enabled')"
          width="100"
        >
          <template #default="{ row }">
            <el-switch
              v-if="isSuperAdmin"
              :model-value="row.enabled"
              @change="(val) => handleToggleEnabled(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.blockRules.builtin')"
          min-width="110"
        >
          <template #default="{ row }">
            <el-tag
              v-if="row.is_builtin"
              type="info"
              size="small"
            >
              {{ t('admin.blockRules.builtin') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          v-if="isSuperAdmin"
          :label="t('admin.blockRules.actions')"
          width="160"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              @click="handleEditRule(row)"
            >
              <el-icon><Edit /></el-icon>
            </el-button>
            <el-button
              link
              type="danger"
              @click="handleDeleteRule(row)"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Applications Card -->
    <el-card style="margin-top: 16px;">
      <template #header>
        <div class="card-header">
          <span>{{ t('admin.blockRules.applications') }}</span>
          <el-button
            type="danger"
            :disabled="selectedApps.length === 0"
            @click="handleRemoveApplications"
          >
            {{ t('admin.blockRules.removeApplications') }}
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loadingApps"
        :data="applications"
        stripe
        @selection-change="handleAppSelectionChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column
          prop="rule_id"
          :label="t('admin.blockRules.ruleId')"
          min-width="100"
        />
        <el-table-column
          prop="scope"
          :label="t('admin.blockRules.scope')"
          width="100"
        >
          <template #default="{ row }">
            <el-tag size="small">
              {{ t(`admin.blockRules.scopes.${row.scope}`) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="target_name"
          :label="t('admin.blockRules.targetName')"
          min-width="150"
        />
        <el-table-column
          prop="status"
          :label="t('admin.blockRules.status')"
          width="100"
        >
          <template #default="{ row }">
            <el-tag
              :type="statusTagType(row.status)"
              size="small"
            >
              {{ t(`admin.blockRules.statuses.${row.status}`) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="ip_version"
          :label="t('admin.blockRules.ipVersion')"
          width="120"
        >
          <template #default="{ row }">
            <el-tag
              size="small"
              type="info"
            >
              {{ t(`admin.blockRules.ipVersions.${row.ip_version || 'both'}`) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="created_at"
          :label="t('admin.blockRules.createdAt')"
          width="180"
        >
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Create/Edit Rule Dialog -->
    <el-dialog
      v-model="showRuleDialog"
      :title="isEdit ? t('admin.blockRules.editRule') : t('admin.blockRules.addRule')"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="ruleFormRef"
        :model="ruleForm"
        :rules="ruleFormRules"
        label-width="120px"
      >
        <el-form-item
          :label="t('admin.blockRules.ruleName')"
          prop="name"
        >
          <el-input v-model="ruleForm.name" />
        </el-form-item>
        <el-form-item
          :label="t('admin.blockRules.category')"
          prop="category"
        >
          <el-select
            v-model="ruleForm.category"
            style="width: 100%;"
          >
            <el-option
              v-for="cat in categories"
              :key="cat"
              :label="t(`admin.blockRules.categories.${cat}`)"
              :value="cat"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('admin.blockRules.description')">
          <el-input
            v-model="ruleForm.description"
            type="textarea"
            :rows="2"
          />
        </el-form-item>
        <el-form-item
          :label="t('admin.blockRules.strings')"
          prop="stringsText"
        >
          <el-input
            v-model="ruleForm.stringsText"
            type="textarea"
            :rows="8"
            :placeholder="t('admin.blockRules.stringsPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('admin.blockRules.enabled')">
          <el-switch v-model="ruleForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRuleDialog = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="handleSubmitRule"
        >
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Apply Rules Dialog -->
    <el-dialog
      v-model="showApplyDialog"
      :title="t('admin.blockRules.applyRules')"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="applyFormRef"
        :model="applyForm"
        :rules="applyFormRules"
        label-width="120px"
      >
        <el-form-item :label="t('admin.blockRules.selectRules')">
          <div>
            <el-tag
              v-for="rule in selectedRules"
              :key="rule.id"
              style="margin: 2px;"
              size="small"
            >
              {{ rule.name }}
            </el-tag>
          </div>
        </el-form-item>
        <el-form-item
          :label="t('admin.blockRules.scope')"
          prop="scope"
        >
          <el-select
            v-model="applyForm.scope"
            style="width: 100%;"
            @change="handleScopeChange"
          >
            <el-option
              v-for="s in scopeOptions"
              :key="s"
              :label="t(`admin.blockRules.scopes.${s}`)"
              :value="s"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="applyForm.scope === 'provider'"
          :label="t('admin.blockRules.selectTargets')"
          prop="target_ids"
        >
          <el-select
            v-model="applyForm.target_ids"
            multiple
            filterable
            style="width: 100%;"
          >
            <el-option
              v-for="p in providerOptions"
              :key="p.id"
              :label="p.name || `Provider #${p.id}`"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="applyForm.scope === 'instance'"
          :label="t('admin.blockRules.selectTargets')"
          prop="target_ids"
        >
          <div style="width: 100%;">
            <el-select
              v-model="instanceProviderFilter"
              :placeholder="t('admin.blockRules.filterByProvider')"
              clearable
              style="width: 100%; margin-bottom: 8px;"
              @change="fetchInstancesForProvider"
            >
              <el-option
                v-for="p in providerOptions"
                :key="p.id"
                :label="p.name || `Provider #${p.id}`"
                :value="p.id"
              />
            </el-select>
            <el-select
              v-model="applyForm.target_ids"
              multiple
              filterable
              :loading="loadingInstances"
              style="width: 100%;"
            >
              <el-option
                v-for="inst in instanceOptions"
                :key="inst.id"
                :label="`${inst.name} (${inst.status})`"
                :value="inst.id"
              />
            </el-select>
          </div>
        </el-form-item>
        <el-form-item
          :label="t('admin.blockRules.ipVersion')"
          prop="ip_version"
        >
          <el-select
            v-model="applyForm.ip_version"
            style="width: 100%;"
          >
            <el-option
              v-for="v in ipVersionOptions"
              :key="v"
              :label="t(`admin.blockRules.ipVersions.${v}`)"
              :value="v"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showApplyDialog = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="handleApplyRules"
        >
          {{ t('admin.blockRules.applyRules') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { Edit, Delete } from '@element-plus/icons-vue'
import { useBlockRuleManagement } from './composables/useBlockRuleManagement.js'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

const {
  rules,
  applications,
  isSuperAdmin,
  providerOptions,
  instanceOptions,
  selectedRules,
  selectedApps,
  loadingRules,
  loadingApps,
  loadingInstances,
  submitting,
  showRuleDialog,
  showApplyDialog,
  isEdit,
  ruleFormRef,
  applyFormRef,
  instanceProviderFilter,
  categories,
  scopeOptions,
  ipVersionOptions,
  ruleForm,
  applyForm,
  ruleFormRules,
  applyFormRules,
  parseStrings,
  formatDate,
  categoryTagType,
  statusTagType,
  fetchInstancesForProvider,
  handleScopeChange,
  handleRuleSelectionChange,
  handleAppSelectionChange,
  handleCreateRule,
  handleEditRule,
  handleSubmitRule,
  handleDeleteRule,
  handleToggleEnabled,
  handleApplyRules,
  handleRemoveApplications
} = useBlockRuleManagement()
</script>

<style scoped>
.block-rules-container {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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
