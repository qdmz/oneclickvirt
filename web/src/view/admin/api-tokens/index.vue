<template>
  <div class="api-tokens-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('admin.apiTokens.title') }}</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="toolbar">
        <el-input
          v-model="keyword"
          :placeholder="t('admin.apiTokens.searchPlaceholder')"
          style="width: 280px;"
          clearable
          @keyup.enter="handleSearch"
          @clear="handleSearch"
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
          {{ t('common.search') }}
        </el-button>
        <el-button
          style="margin-left: 10px;"
          @click="resetFilters"
        >
          {{ t('common.reset') }}
        </el-button>
        <el-popconfirm
          v-if="selectedIds.length > 0"
          :title="t('admin.apiTokens.batchDeleteConfirm', { count: selectedIds.length })"
          @confirm="handleBatchDelete"
        >
          <template #reference>
            <el-button
              type="danger"
              style="margin-left: 10px;"
            >
              {{ t('admin.apiTokens.batchDelete') }}（{{ selectedIds.length }}）
            </el-button>
          </template>
        </el-popconfirm>
      </div>

      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="tokens"
        style="width: 100%; margin-top: 16px;"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="50"
        />
        <el-table-column
          :label="t('admin.apiTokens.userId')"
          prop="userId"
          min-width="100"
        />
        <el-table-column
          :label="t('admin.apiTokens.username')"
          prop="username"
          min-width="120"
        />
        <el-table-column
          :label="t('admin.apiTokens.userType')"
          prop="userType"
          width="120"
        >
          <template #default="{ row }">
            <el-tag
              :type="row.userType === 'admin' ? 'danger' : row.userType === 'normal_admin' ? 'warning' : 'info'"
              size="small"
            >
              {{ formatUserType(row.userType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.apiTokens.tokenName')"
          prop="name"
          min-width="150"
        />
        <el-table-column
          :label="t('admin.apiTokens.tokenPrefix')"
          prop="tokenPrefix"
          min-width="150"
        >
          <template #default="{ row }">
            <el-tag
              type="info"
              size="small"
            >
              {{ row.tokenPrefix }}...
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.apiTokens.useCount')"
          prop="useCount"
          min-width="110"
        />
        <el-table-column
          :label="t('admin.apiTokens.createdAt')"
          width="160"
        >
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.apiTokens.expiresAt')"
          width="160"
        >
          <template #default="{ row }">
            <span v-if="!row.expiresAt">{{ t('admin.apiTokens.never') }}</span>
            <span v-else>{{ formatDate(row.expiresAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.apiTokens.lastUsedAt')"
          width="160"
        >
          <template #default="{ row }">
            <span v-if="!row.lastUsedAt">{{ t('admin.apiTokens.unused') }}</span>
            <span v-else>{{ formatDate(row.lastUsedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('common.actions')"
          width="100"
          fixed="right"
        >
          <template #default="{ row }">
            <el-popconfirm
              :title="t('admin.apiTokens.deleteConfirm')"
              @confirm="handleDelete(row.id)"
            >
              <template #reference>
                <el-button
                  type="danger"
                  size="small"
                  :icon="Delete"
                />
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Delete, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { adminGetApiTokenList, adminDeleteApiToken, adminBatchDeleteApiTokens } from '@/api/admin'

const { t } = useI18n()

const tokens = ref([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const keyword = ref('')
const selectedIds = ref([])

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

const formatUserType = (userType) => {
  if (userType === 'admin') return t('admin.apiTokens.userTypeAdmin')
  if (userType === 'normal_admin') return t('admin.apiTokens.userTypeNormalAdmin')
  return t('admin.apiTokens.userTypeUser')
}

const loadTokens = async () => {
  loading.value = true
  try {
    const res = await adminGetApiTokenList({
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: keyword.value
    })
    tokens.value = res.data?.items || []
    total.value = res.data?.total || 0
  } catch {
    ElMessage.error(t('common.operationFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadTokens()
}

const resetFilters = () => {
  keyword.value = ''
  currentPage.value = 1
  loadTokens()
}

const handleSelectionChange = (selection) => {
  selectedIds.value = selection.map((row) => row.id)
}

const handleDelete = async (id) => {
  try {
    await adminDeleteApiToken(id)
    ElMessage.success(t('admin.apiTokens.deleteSuccess'))
    await loadTokens()
  } catch {
    ElMessage.error(t('admin.apiTokens.deleteFailed'))
  }
}

const handleBatchDelete = async () => {
  if (selectedIds.value.length === 0) return
  try {
    await adminBatchDeleteApiTokens(selectedIds.value)
    ElMessage.success(t('admin.apiTokens.batchDeleteSuccess'))
    selectedIds.value = []
    await loadTokens()
  } catch {
    ElMessage.error(t('admin.apiTokens.batchDeleteFailed'))
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadTokens()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadTokens()
}

onMounted(() => {
  loadTokens()
})
</script>

<style scoped>
.api-tokens-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.toolbar {
  display: flex;
  align-items: center;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
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
