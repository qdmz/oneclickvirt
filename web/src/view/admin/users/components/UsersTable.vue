<template>
  <el-table
    v-loading="loading"
    :data="users"
    class="users-table"
    :cell-style="{ padding: '12px 0' }"
    :header-cell-style="{ background: '#f5f7fa', padding: '14px 0', fontWeight: '600' }"
    @selection-change="$emit('selection-change', $event)"
  >
    <el-table-column
      type="selection"
      width="55"
      align="center"
    />
    <el-table-column
      prop="id"
      label="ID"
      width="80"
      align="center"
    />
    <el-table-column
      prop="username"
      :label="$t('admin.users.username')"
      min-width="140"
      show-overflow-tooltip
    />
    <el-table-column
      prop="email"
      :label="$t('admin.users.email')"
      min-width="180"
      show-overflow-tooltip
    />
    <el-table-column
      prop="nickname"
      :label="$t('admin.users.nickname')"
      min-width="140"
      show-overflow-tooltip
    />
    <el-table-column
      prop="level"
      :label="$t('admin.users.level')"
      width="100"
      align="center"
    >
      <template #default="scope">
        <el-tag :type="getLevelTagType(scope.row.level)">
          {{ $t('admin.users.levelTag', { level: scope.row.level }) }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column
      prop="balance"
      :label="$t('admin.users.balance')"
      width="120"
      align="center"
    >
      <template #default="scope">
        <span class="balance-cell">¥{{ Number(scope.row.balance || 0).toFixed(2) }}</span>
      </template>
    </el-table-column>
    <el-table-column
      prop="userType"
      :label="$t('admin.users.userType')"
      width="120"
      align="center"
    >
      <template #default="scope">
        <el-tag :type="getUserTypeTagType(scope.row.userType)">
          {{ getUserTypeLabel(scope.row.userType) }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column
      prop="status"
      :label="$t('common.status')"
      width="100"
      align="center"
    >
      <template #default="scope">
        <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
          {{ scope.row.status === 1 ? $t('admin.users.active') : $t('admin.users.disabled') }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column
      prop="expiresAt"
      :label="$t('admin.users.expiresAt')"
      width="180"
      align="center"
    >
      <template #default="scope">
        <div v-if="scope.row.expiresAt">
          <el-tag
            :type="isExpired(scope.row.expiresAt) ? 'danger' : 'success'"
            size="small"
          >
            {{ formatDateTime(scope.row.expiresAt) }}
          </el-tag>
          <div
            v-if="scope.row.isManualExpiry"
            style="margin-top: 4px;"
          >
            <el-tag
              size="small"
              type="info"
            >
              {{ $t('admin.users.manualExpiry') }}
            </el-tag>
          </div>
        </div>
        <span v-else>-</span>
      </template>
    </el-table-column>
    <el-table-column
      :label="$t('common.actions')"
      width="520"
      fixed="right"
      align="center"
    >
      <template #default="scope">
        <div class="action-buttons">
          <el-button
            size="small"
            @click="$emit('edit', scope.row)"
          >
            {{ $t('common.edit') }}
          </el-button>
          <el-dropdown @command="(level) => $emit('set-user-level', scope.row, level)">
            <el-button
              size="small"
              type="primary"
            >
              {{ $t('admin.users.levelSetting') }}<el-icon class="el-icon--right">
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
          <el-button
            size="small"
            type="warning"
            @click="$emit('set-expiry', scope.row)"
          >
            {{ $t('admin.users.setExpiry') }}
          </el-button>
          <el-button
            size="small"
            :type="scope.row.status === 1 ? 'danger' : 'success'"
            @click="$emit('toggle-status', scope.row)"
          >
            {{ scope.row.status === 1 ? $t('admin.users.disable') : $t('admin.users.enable') }}
          </el-button>
          <el-button
            size="small"
            type="warning"
            @click="$emit('reset-password', scope.row)"
          >
            {{ $t('admin.users.resetPassword') }}
          </el-button>
          <el-button
            size="small"
            type="success"
            @click="$emit('adjust-balance', scope.row)"
          >
            {{ $t('admin.users.adjustBalance') }}
          </el-button>
          <el-button
            v-if="scope.row.userType !== 'admin'"
            size="small"
            type="info"
            @click="$emit('login-as', scope.row)"
          >
            {{ $t('admin.users.loginAs') }}
          </el-button>
        </div>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { getLevelTagType } from '@/utils/levels'

defineProps({
  users: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  availableLevels: { type: Array, default: () => [1, 2, 3, 4, 5] }
})

defineEmits(['selection-change', 'edit', 'set-user-level', 'set-expiry', 'toggle-status', 'reset-password', 'login-as', 'adjust-balance'])

const { t, locale } = useI18n()

const getUserTypeLabel = (userType) => {
  const labelMap = {
    'user': t('admin.users.normalUser'),
    'normal_admin': t('admin.users.normalAdmin'),
    'admin': t('admin.users.adminUser')
  }
  return labelMap[userType] || t('common.unknown')
}

const getUserTypeTagType = (userType) => {
  const typeMap = { 'user': '', 'normal_admin': 'warning', 'admin': 'danger' }
  return typeMap[userType] || ''
}

const formatDateTime = (dateTimeStr) => {
  if (!dateTimeStr) return '-'
  const date = new Date(dateTimeStr)
  return date.toLocaleString(locale.value, {
    year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit'
  })
}

const isExpired = (dateTimeStr) => {
  if (!dateTimeStr) return false
  return new Date(dateTimeStr) < new Date()
}
</script>

<style scoped>
.balance-cell {
  color: var(--el-color-danger);
  font-weight: 600;
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
