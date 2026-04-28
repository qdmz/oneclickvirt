<template>
  <div class="api-keys-container">
    <!-- 加载状态 -->
    <div
      v-if="loading"
      class="loading-container"
    >
      <el-loading-directive />
      <div class="loading-text">
        {{ t('user.apiKeys.loading') }}
      </div>
    </div>

    <!-- 主要内容 -->
    <div v-else>
      <!-- 统计卡片 -->
      <el-row
        :gutter="20"
        class="stats-row"
      >
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon total">
                <el-icon><Key /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">
                  {{ stats.total }}
                </div>
                <div class="stat-label">
                  {{ t('user.apiKeys.total') }}
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon active">
                <el-icon><CircleCheck /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">
                  {{ stats.active }}
                </div>
                <div class="stat-label">
                  {{ t('user.apiKeys.active') }}
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon disabled">
                <el-icon><CircleClose /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">
                  {{ stats.disabled }}
                </div>
                <div class="stat-label">
                  {{ t('user.apiKeys.disabled') }}
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon expired">
                <el-icon><Warning /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">
                  {{ stats.expired }}
                </div>
                <div class="stat-label">
                  {{ t('user.apiKeys.expired') }}
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- API 密钥列表 -->
      <el-card class="api-keys-card">
        <template #header>
          <div class="card-header">
            <span>{{ t('user.apiKeys.title') }}</span>
            <el-button
              type="primary"
              :icon="Plus"
              @click="showCreateDialog = true"
            >
              {{ t('user.apiKeys.create') }}
            </el-button>
          </div>
        </template>

        <!-- API 密钥表格 -->
        <el-table
          :data="apiKeys"
          style="width: 100%"
          v-loading="tableLoading"
        >
          <el-table-column
            prop="id"
            :label="t('user.apiKeys.id')"
            width="80"
          />
          <el-table-column
            prop="name"
            :label="t('user.apiKeys.name')"
            min-width="150"
          />
          <el-table-column
            prop="key"
            :label="t('user.apiKeys.key')"
            min-width="200"
          >
            <template #default="{ row }">
              <div class="key-display">
                <el-input
                  v-model="row.key"
                  readonly
                  size="small"
                >
                  <template #append>
                    <el-button
                      :icon="CopyDocument"
                      @click="copyKey(row.key)"
                    />
                  </template>
                </el-input>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="status"
            :label="t('user.apiKeys.status')"
            width="100"
          >
            <template #default="{ row }">
              <el-tag :type="getStatusTagType(row.status)">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            prop="expiresAt"
            :label="t('user.apiKeys.expiresAt')"
            width="180"
          >
            <template #default="{ row }">
              <span v-if="row.expiresAt">
                {{ formatDate(row.expiresAt) }}
              </span>
              <span v-else>
                {{ t('user.apiKeys.neverExpire') }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            prop="lastUsedAt"
            :label="t('user.apiKeys.lastUsedAt')"
            width="180"
          >
            <template #default="{ row }">
              <span v-if="row.lastUsedAt">
                {{ formatDate(row.lastUsedAt) }}
              </span>
              <span v-else>
                {{ t('user.apiKeys.neverUsed') }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            prop="createdAt"
            :label="t('user.apiKeys.createdAt')"
            width="180"
          >
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column
            :label="t('common.actions')"
            width="200"
            fixed="right"
          >
            <template #default="{ row }">
              <el-button
                v-if="row.status === 'active'"
                type="warning"
                size="small"
                @click="handleRevoke(row)"
              >
                {{ t('user.apiKeys.revoke') }}
              </el-button>
              <el-button
                v-if="row.status === 'disabled'"
                type="success"
                size="small"
                @click="handleEnable(row)"
              >
                {{ t('user.apiKeys.enable') }}
              </el-button>
              <el-button
                type="danger"
                size="small"
                @click="handleDelete(row)"
              >
                {{ t('common.delete') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <!-- 创建 API 密钥对话框 -->
      <el-dialog
        v-model="showCreateDialog"
        :title="t('user.apiKeys.create')"
        width="500px"
      >
        <el-form
          ref="createFormRef"
          :model="createForm"
          :rules="createRules"
          label-width="100px"
        >
          <el-form-item
            :label="t('user.apiKeys.name')"
            prop="name"
          >
            <el-input
              v-model="createForm.name"
              :placeholder="t('user.apiKeys.namePlaceholder')"
              clearable
            />
          </el-form-item>
          <el-form-item
            :label="t('user.apiKeys.expiresAt')"
            prop="expiresAt"
          >
            <el-date-picker
              v-model="createForm.expiresAt"
              type="datetime"
              :placeholder="t('user.apiKeys.expiresAtPlaceholder')"
              :disabled-date="disabledDate"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
              style="width: 100%;"
            />
            <div class="form-tip">
              {{ t('user.apiKeys.expiresAtTip') }}
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

      <!-- 显示新创建的 API 密钥 -->
      <el-dialog
        v-model="showNewKeyDialog"
        :title="t('user.apiKeys.createSuccess')"
        width="600px"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
      >
        <el-alert
          :title="t('user.apiKeys.saveKeyWarning')"
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 20px;"
        />
        <div class="new-key-content">
          <el-text
            type="info"
            style="display: block; margin-bottom: 10px;"
          >
            {{ t('user.apiKeys.newKey') }}：
          </el-text>
          <el-input
            v-model="newKey"
            readonly
            style="width: 100%; font-family: monospace; font-size: 14px;"
          >
            <template #append>
              <el-button
                :icon="CopyDocument"
                @click="copyKey(newKey)"
              >
                {{ t('common.copy') }}
              </el-button>
            </template>
          </el-input>
        </div>
        <template #footer>
          <el-button
            type="primary"
            @click="closeNewKeyDialog"
          >
            {{ t('user.apiKeys.iHaveSaved') }}
          </el-button>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, CopyDocument, Key, CircleCheck, CircleClose, Warning } from '@element-plus/icons-vue'
import {
  getAPIKeys,
  createAPIKey,
  updateAPIKey,
  deleteAPIKey,
  revokeAPIKey,
  getAPIKeyStats
} from '@/api/api-key'

const { t } = useI18n()

// 加载状态
const loading = ref(true)
const tableLoading = ref(false)
const creating = ref(false)

// 统计数据
const stats = reactive({
  total: 0,
  active: 0,
  disabled: 0,
  expired: 0
})

// API 密钥列表
const apiKeys = ref([])

// 创建对话框
const showCreateDialog = ref(false)
const showNewKeyDialog = ref(false)
const newKey = ref('')

// 创建表单
const createFormRef = ref()
const createForm = reactive({
  name: '',
  expiresAt: null
})

const createRules = reactive({
  name: [
    { required: true, message: t('user.apiKeys.nameRequired'), trigger: 'blur' },
    { min: 1, max: 100, message: t('user.apiKeys.nameLength'), trigger: 'blur' }
  ]
})

// 获取统计数据
const fetchStats = async () => {
  try {
    const response = await getAPIKeyStats()
    if (response.code === 0 || response.code === 200) {
      Object.assign(stats, response.data)
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 获取 API 密钥列表
const fetchAPIKeys = async () => {
  tableLoading.value = true
  try {
    const response = await getAPIKeys()
    if (response.code === 0 || response.code === 200) {
      apiKeys.value = response.data || []
    }
  } catch (error) {
    console.error('获取 API 密钥列表失败:', error)
    ElMessage.error(t('user.apiKeys.fetchFailed'))
  } finally {
    tableLoading.value = false
  }
}

// 创建 API 密钥
const handleCreate = async () => {
  if (!createFormRef.value) return

  await createFormRef.value.validate(async (valid) => {
    if (!valid) return

    creating.value = true
    try {
      const response = await createAPIKey({
        name: createForm.name,
        expiresAt: createForm.expiresAt
      })

      if (response.code === 0 || response.code === 200) {
        newKey.value = response.data.key
        showNewKeyDialog.value = true
        showCreateDialog.value = false

        // 重置表单
        createForm.name = ''
        createForm.expiresAt = null

        // 刷新列表
        await fetchAPIKeys()
        await fetchStats()
      } else {
        ElMessage.error(response.msg || t('user.apiKeys.createFailed'))
      }
    } catch (error) {
      console.error('创建 API 密钥失败:', error)
      ElMessage.error(t('user.apiKeys.createFailed'))
    } finally {
      creating.value = false
    }
  })
}

// 撤销 API 密钥
const handleRevoke = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('user.apiKeys.revokeConfirm'),
      t('user.apiKeys.revokeTitle'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )

    const response = await revokeAPIKey(row.id)
    if (response.code === 0 || response.code === 200) {
      ElMessage.success(t('user.apiKeys.revokeSuccess'))
      await fetchAPIKeys()
      await fetchStats()
    } else {
      ElMessage.error(response.msg || t('user.apiKeys.revokeFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('撤销 API 密钥失败:', error)
      ElMessage.error(t('user.apiKeys.revokeFailed'))
    }
  }
}

// 启用 API 密钥
const handleEnable = async (row) => {
  try {
    const response = await updateAPIKey(row.id, { status: 'active' })
    if (response.code === 0 || response.code === 200) {
      ElMessage.success(t('user.apiKeys.enableSuccess'))
      await fetchAPIKeys()
      await fetchStats()
    } else {
      ElMessage.error(response.msg || t('user.apiKeys.enableFailed'))
    }
  } catch (error) {
    console.error('启用 API 密钥失败:', error)
    ElMessage.error(t('user.apiKeys.enableFailed'))
  }
}

// 删除 API 密钥
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('user.apiKeys.deleteConfirm'),
      t('user.apiKeys.deleteTitle'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )

    const response = await deleteAPIKey(row.id)
    if (response.code === 0 || response.code === 200) {
      ElMessage.success(t('user.apiKeys.deleteSuccess'))
      await fetchAPIKeys()
      await fetchStats()
    } else {
      ElMessage.error(response.msg || t('user.apiKeys.deleteFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除 API 密钥失败:', error)
      ElMessage.error(t('user.apiKeys.deleteFailed'))
    }
  }
}

// 复制密钥
const copyKey = async (key) => {
  try {
    // 优先使用 Clipboard API
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(key)
      ElMessage.success(t('user.apiKeys.copySuccess'))
      return
    }

    // 降级方案：使用传统的 document.execCommand
    const textArea = document.createElement('textarea')
    textArea.value = key
    textArea.style.position = 'fixed'
    textArea.style.left = '-999999px'
    textArea.style.top = '-999999px'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()

    try {
      // @ts-ignore - execCommand 已废弃但作为降级方案仍需使用
      const successful = document.execCommand('copy')
      if (successful) {
        ElMessage.success(t('user.apiKeys.copySuccess'))
      } else {
        throw new Error('execCommand failed')
      }
    } finally {
      document.body.removeChild(textArea)
    }
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error(t('user.apiKeys.copyFailed'))
  }
}

// 关闭新密钥对话框
const closeNewKeyDialog = () => {
  showNewKeyDialog.value = false
  newKey.value = ''
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 禁用过去的日期
const disabledDate = (time) => {
  return time.getTime() < Date.now() - 8.64e7
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  switch (status) {
    case 'active':
      return 'success'
    case 'disabled':
      return 'info'
    default:
      return 'info'
  }
}

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 'active':
      return t('user.apiKeys.active')
    case 'disabled':
      return t('user.apiKeys.disabled')
    default:
      return status
  }
}

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([fetchStats(), fetchAPIKeys()])
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  color: #666;
}

.loading-text {
  margin-top: 16px;
  font-size: 14px;
}

.api-keys-container {
  padding: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 10px;
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
  font-size: 24px;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-icon.active {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.stat-icon.disabled {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.stat-icon.expired {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.api-keys-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.key-display {
  width: 100%;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.new-key-content {
  padding: 10px 0;
}
</style>
