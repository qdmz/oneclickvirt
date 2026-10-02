<template>
  <div class="api-tokens-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ t('user.apiTokens.title') }}</span>
          <el-button
            type="primary"
            @click="showCreateDialog = true"
          >
            {{ t('user.apiTokens.createToken') }}
          </el-button>
        </div>
      </template>

      <el-alert
        :title="t('user.apiTokens.usageNote')"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 16px"
      />
      <el-alert
        :title="t('user.apiTokens.maxTokensWarning')"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 16px"
      />

      <div
        v-if="loading"
        class="loading-container"
      >
        <el-skeleton
          :rows="3"
          animated
        />
      </div>

      <el-empty
        v-else-if="tokens.length === 0"
        :description="t('user.apiTokens.noTokens')"
      >
        <el-button
          type="primary"
          @click="showCreateDialog = true"
        >
          {{ t('user.apiTokens.createNewToken') }}
        </el-button>
      </el-empty>

      <div v-else>
        <div class="tokens-summary">
          {{ t('user.apiTokens.totalTokens', { count: tokens.length }) }}
        </div>
        <el-table
          :data="tokens"
          style="width: 100%"
          stripe
        >
          <el-table-column
            :label="t('user.apiTokens.tokenName')"
            prop="name"
            min-width="150"
          />
          <el-table-column
            :label="t('user.apiTokens.tokenPrefix')"
            prop="tokenPrefix"
            width="150"
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
            :label="t('user.apiTokens.useCount')"
            prop="useCount"
            min-width="110"
          />
          <el-table-column
            :label="t('user.apiTokens.createdAt')"
            width="160"
          >
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column
            :label="t('user.apiTokens.expiresAt')"
            width="160"
          >
            <template #default="{ row }">
              <span v-if="!row.expiresAt">{{ t('user.apiTokens.never') }}</span>
              <span v-else>{{ formatDate(row.expiresAt) }}</span>
            </template>
          </el-table-column>
          <el-table-column
            :label="t('user.apiTokens.lastUsedAt')"
            width="160"
          >
            <template #default="{ row }">
              <span v-if="!row.lastUsedAt">{{ t('user.apiTokens.unused') }}</span>
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
                :title="t('user.apiTokens.deleteConfirm')"
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
      </div>
    </el-card>

    <!-- 创建 Token 对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="t('user.apiTokens.createToken')"
      width="480px"
      @close="resetCreateForm"
    >
      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-width="100px"
      >
        <el-form-item
          :label="t('user.apiTokens.tokenName')"
          prop="name"
        >
          <el-input
            v-model="createForm.name"
            :placeholder="t('user.apiTokens.tokenNamePlaceholder')"
            maxlength="64"
            show-word-limit
          />
        </el-form-item>
        <el-form-item
          :label="t('user.apiTokens.expireDays')"
          prop="expireDays"
        >
          <el-input-number
            v-model="createForm.expireDays"
            :min="0"
            :max="3650"
            style="width: 100%"
          />
          <div class="form-hint">
            {{ t('user.apiTokens.expireDaysHint') }}
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="creating"
          @click="handleCreate"
        >
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 新建 Token 展示对话框 -->
    <el-dialog
      v-model="showTokenDialog"
      :title="t('user.apiTokens.createSuccess')"
      width="560px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-alert
        :title="t('user.apiTokens.copyTokenWarning')"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 16px"
      />
      <div class="token-display">
        <el-input
          :value="newToken"
          readonly
          type="textarea"
          :rows="3"
          resize="none"
        />
        <el-button
          type="primary"
          :icon="CopyDocument"
          style="margin-top: 8px; width: 100%"
          @click="copyToken"
        >
          {{ t('common.copy') }}
        </el-button>
      </div>
      <template #footer>
        <el-button
          type="primary"
          @click="showTokenDialog = false; loadTokens()"
        >
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Delete, CopyDocument } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { createApiToken, getApiTokenList, deleteApiToken } from '@/api/auth'

const { t } = useI18n()

const tokens = ref([])
const loading = ref(false)
const creating = ref(false)
const showCreateDialog = ref(false)
const showTokenDialog = ref(false)
const newToken = ref('')
const createFormRef = ref(null)

const createForm = ref({
  name: '',
  expireDays: 0
})

const createRules = {
  name: [
    { required: true, message: () => t('user.apiTokens.tokenNameRequired'), trigger: 'blur' }
  ]
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

const loadTokens = async () => {
  loading.value = true
  try {
    const res = await getApiTokenList()
    tokens.value = res.data?.items || []
  } catch {
    ElMessage.error(t('common.operationFailed'))
  } finally {
    loading.value = false
  }
}

const resetCreateForm = () => {
  createForm.value = { name: '', expireDays: 0 }
  createFormRef.value?.clearValidate()
}

const handleCreate = async () => {
  if (!createFormRef.value) return
  const valid = await createFormRef.value.validate().catch(() => false)
  if (!valid) return

  creating.value = true
  try {
    const res = await createApiToken({
      name: createForm.value.name,
      expireDays: createForm.value.expireDays || 0
    })
    const token = res.data?.token || res.data?.Token || ''
    newToken.value = token
    showCreateDialog.value = false
    showTokenDialog.value = true
    resetCreateForm()
  } catch {
    ElMessage.error(t('user.apiTokens.createFailed'))
  } finally {
    creating.value = false
  }
}

const handleDelete = async (id) => {
  try {
    await deleteApiToken(id)
    ElMessage.success(t('user.apiTokens.deleteSuccess'))
    await loadTokens()
  } catch {
    ElMessage.error(t('user.apiTokens.deleteFailed'))
  }
}

const copyToken = async () => {
  try {
    await navigator.clipboard.writeText(newToken.value)
    ElMessage.success(t('user.apiTokens.copiedSuccess'))
  } catch {
    // fallback
    const el = document.createElement('textarea')
    el.value = newToken.value
    document.body.appendChild(el)
    el.select()
    document.execCommand('copy')
    document.body.removeChild(el)
    ElMessage.success(t('user.apiTokens.copiedSuccess'))
  }
}

onMounted(() => {
  loadTokens()
})
</script>

<style scoped>
.api-tokens-page {
  padding: 16px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.loading-container {
  padding: 20px;
}
.tokens-summary {
  margin-bottom: 12px;
  color: #606266;
  font-size: 14px;
}
.token-display {
  padding: 8px 0;
}
.form-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
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
