<template>
  <div class="upstream-container">
    <!-- 顶部操作栏 -->
    <el-card class="header-card">
      <div class="header-row">
        <div class="header-left">
          <h3 class="page-title">{{ $t('admin.upstream.title') }}</h3>
          <p class="page-desc">{{ $t('admin.upstream.desc') }}</p>
        </div>
        <div class="header-right">
          <el-button type="primary" :icon="Plus" @click="openCreateDialog">
            {{ $t('admin.upstream.addProvider') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 上游节点列表 -->
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.upstream.providerList') }}</span>
        </div>
      </template>

      <el-empty
        v-if="providers.length === 0 && !loading"
        :description="$t('admin.upstream.empty')"
      />

      <el-table v-else :data="providers" style="width: 100%">
        <el-table-column prop="name" :label="$t('admin.upstream.name')" min-width="140" />
        <el-table-column :label="$t('admin.upstream.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? $t('admin.upstream.active') : $t('admin.upstream.inactive') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('admin.upstream.region')" width="120">
          <template #default="{ row }">{{ row.region || '-' }}</template>
        </el-table-column>
        <el-table-column :label="$t('admin.upstream.actions')" width="360" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openTestDialog(row)">
              {{ $t('admin.upstream.test') }}
            </el-button>
            <el-button size="small" type="primary" :loading="syncingId === row.id" @click="handleSync(row)">
              {{ $t('admin.upstream.syncProducts') }}
            </el-button>
            <el-button size="small" type="warning" @click="openEditDialog(row)">
              {{ $t('admin.upstream.edit') }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">
              {{ $t('admin.upstream.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 已同步产品列表 -->
    <el-card v-loading="productsLoading" class="products-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.upstream.syncedProducts') }}</span>
          <el-button size="small" @click="loadSyncedProducts">
            {{ $t('admin.upstream.refresh') }}
          </el-button>
        </div>
      </template>

      <el-empty
        v-if="syncedProducts.length === 0 && !productsLoading"
        :description="$t('admin.upstream.noProducts')"
      />

      <el-table v-else :data="syncedProducts" style="width: 100%">
        <el-table-column prop="name" :label="$t('admin.upstream.productName')" min-width="160" />
        <el-table-column :label="$t('admin.upstream.productSpec')" min-width="180">
          <template #default="{ row }">
            <span>{{ row.cpu }}C / {{ row.memory }}MB / {{ row.disk }}MB</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('admin.upstream.bandwidth')" width="120">
          <template #default="{ row }">{{ row.bandwidth }}Mbps</template>
        </el-table-column>
        <el-table-column :label="$t('admin.upstream.price')" width="120">
          <template #default="{ row }">¥{{ row.price }}</template>
        </el-table-column>
        <el-table-column :label="$t('admin.upstream.period')" width="120">
          <template #default="{ row }">{{ row.periodType }} / {{ row.periodValue }}</template>
        </el-table-column>
        <el-table-column :label="$t('admin.upstream.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? $t('admin.upstream.onSale') : $t('admin.upstream.offSale') }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingId ? $t('admin.upstream.editProvider') : $t('admin.upstream.addProvider')"
      width="620px"
      destroy-on-close
    >
      <el-alert
        v-if="editingId"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 16px"
      >
        {{ $t('admin.upstream.editHint') }}
      </el-alert>
      <el-form ref="formRef" :model="form" :rules="activeFormRules" label-width="120px">
        <el-form-item :label="$t('admin.upstream.name')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('admin.upstream.namePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('admin.upstream.region')">
          <el-input v-model="form.region" :placeholder="$t('admin.upstream.regionPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('admin.upstream.authType')" prop="authConfig.authType">
          <el-select v-model="form.authConfig.authType" style="width: 100%">
            <el-option label="api_client（API客户端签名）" value="api_client" />
            <el-option label="module（用户名密码）" value="module" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('admin.upstream.baseUrl')" prop="authConfig.baseUrl">
          <el-input v-model="form.authConfig.baseUrl" :placeholder="$t('admin.upstream.baseUrlPlaceholder')" />
        </el-form-item>

        <template v-if="form.authConfig.authType === 'api_client'">
          <el-form-item :label="$t('admin.upstream.apiId')" prop="authConfig.apiId">
            <el-input v-model="form.authConfig.apiId" />
          </el-form-item>
          <el-form-item :label="$t('admin.upstream.apiKey')" prop="authConfig.apiKey">
            <el-input v-model="form.authConfig.apiKey" type="password" show-password />
          </el-form-item>
        </template>
        <template v-else>
          <el-form-item :label="$t('admin.upstream.username')" prop="authConfig.username">
            <el-input v-model="form.authConfig.username" />
          </el-form-item>
          <el-form-item :label="$t('admin.upstream.password')" prop="authConfig.password">
            <el-input v-model="form.authConfig.password" type="password" show-password />
          </el-form-item>
        </template>

        <el-form-item :label="$t('admin.upstream.signMethod')">
          <el-select v-model="form.authConfig.signMethod" style="width: 100%">
            <el-option label="md5" value="md5" />
            <el-option label="sha1" value="sha1" />
            <el-option label="sha256" value="sha256" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('admin.upstream.timeout')">
          <el-input-number v-model="form.authConfig.timeout" :min="5" :max="120" /> <span>s</span>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('admin.upstream.cancel') }}</el-button>
        <el-button @click="handleTestBeforeSave">{{ $t('admin.upstream.testAndSave') }}</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">
          {{ $t('admin.upstream.save') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 测试连接对话框 -->
    <el-dialog
      v-model="testDialogVisible"
      :title="$t('admin.upstream.testConnection')"
      width="480px"
    >
      <div v-loading="testing" class="test-result">
        <el-alert
          v-if="testResult !== null"
          :type="testResult.ok ? 'success' : 'error'"
          :title="testResult.ok ? $t('admin.upstream.testSuccess') : $t('admin.upstream.testFailed')"
          :description="testResult.msg"
          :closable="false"
        />
        <el-empty
          v-else-if="!testing"
          :description="$t('admin.upstream.testHint')"
        />
      </div>
      <template #footer>
        <el-button @click="testDialogVisible = false">{{ $t('admin.upstream.close') }}</el-button>
        <el-button type="primary" :loading="testing" @click="runTest">
          {{ $t('admin.upstream.testNow') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getUpstreamProviderList,
  createUpstreamProvider,
  updateUpstreamProvider,
  deleteUpstreamProvider,
  testUpstreamConnection,
  syncUpstreamProducts
} from '@/api/admin/upstream'
import { getAdminProductList } from '@/api/admin'

const loading = ref(false)
const providers = ref([])
const syncingId = ref(null)

// 已同步产品列表
const productsLoading = ref(false)
const syncedProducts = ref([])

// 对话框状态
const dialogVisible = ref(false)
const editingId = ref(null)
const saving = ref(false)
const formRef = ref(null)

const defaultAuthConfig = () => ({
  authType: 'api_client',
  baseUrl: '',
  apiId: '',
  apiKey: '',
  username: '',
  password: '',
  signMethod: 'md5',
  timeout: 30
})

const form = reactive({
  name: '',
  region: '',
  authConfig: defaultAuthConfig()
})

const formRules = {
  name: [{ required: true, message: '请输入节点名称', trigger: 'blur' }],
  'authConfig.baseUrl': [{ required: true, message: '请输入API地址', trigger: 'blur' }]
}

// 编辑模式：baseUrl 留空则保留原有配置，不校验必填
const activeFormRules = computed(() => {
  if (editingId.value) {
    return {
      name: [{ required: true, message: '请输入节点名称', trigger: 'blur' }]
    }
  }
  return formRules
})

// 测试连接对话框
const testDialogVisible = ref(false)
const testing = ref(false)
const testResult = ref(null)
const testProviderId = ref(null)

const loadProviders = async () => {
  loading.value = true
  try {
    const res = await getUpstreamProviderList()
    providers.value = res.data || []
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

const loadSyncedProducts = async () => {
  productsLoading.value = true
  try {
    const res = await getAdminProductList({ page: 1, pageSize: 100 })
    const list = res.data?.list || res.data?.items || res.data || []
    syncedProducts.value = list.filter((p) => p.upstreamType === 'idcsmart')
  } catch (e) {
    ElMessage.error(e.message || '加载产品失败')
  } finally {
    productsLoading.value = false
  }
}

const openCreateDialog = () => {
  editingId.value = null
  Object.assign(form, { name: '', region: '', authConfig: defaultAuthConfig() })
  dialogVisible.value = true
}

const openEditDialog = (row) => {
  editingId.value = row.id
  // AuthConfig 不返回前端，编辑时仅保留基础信息，需重新填写 API 配置
  Object.assign(form, {
    name: row.name,
    region: row.region,
    authConfig: defaultAuthConfig()
  })
  dialogVisible.value = true
}

const handleSave = async () => {
  await formRef.value.validate()
  saving.value = true
  try {
    const payload = {
      name: form.name,
      region: form.region,
      // 编辑模式下 baseUrl 留空 = 保留原有 API 配置；有值 = 更新配置
      authConfig: (editingId.value && !form.authConfig.baseUrl) ? undefined : form.authConfig
    }
    if (editingId.value) {
      await updateUpstreamProvider(editingId.value, payload)
      ElMessage.success('更新成功')
    } else {
      await createUpstreamProvider(payload)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    await loadProviders()
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const handleTestBeforeSave = async () => {
  await formRef.value.validate()
  try {
    const res = await testUpstreamConnection({ authConfig: form.authConfig })
    ElMessage.success(res.msg || '连接成功')
  } catch (e) {
    ElMessage.error(e.message || '连接失败')
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除上游节点「${row.name}」？`, '提示', {
      type: 'warning'
    })
    await deleteUpstreamProvider(row.id)
    ElMessage.success('删除成功')
    await loadProviders()
  } catch (e) {
    if (e !== 'cancel' && e?.message) {
      ElMessage.error(e.message)
    }
  }
}

const handleSync = async (row) => {
  syncingId.value = row.id
  try {
    const res = await syncUpstreamProducts(row.id)
    ElMessage.success(`同步完成：新增 ${res.data?.synced ?? 0} 个，跳过 ${res.data?.skipped ?? 0} 个`)
    await loadSyncedProducts()
  } catch (e) {
    ElMessage.error(e.message || '同步失败')
  } finally {
    syncingId.value = null
  }
}

const openTestDialog = (row) => {
  testProviderId.value = row.id
  testResult.value = null
  testDialogVisible.value = true
}

const runTest = async () => {
  testing.value = true
  testResult.value = null
  try {
    await testUpstreamConnection({ providerId: testProviderId.value })
    testResult.value = { ok: true, msg: '连接成功' }
  } catch (e) {
    testResult.value = { ok: false, msg: e.message || '连接失败' }
  } finally {
    testing.value = false
  }
}

onMounted(() => {
  loadProviders()
  loadSyncedProducts()
})
</script>

<style scoped>
.upstream-container {
  padding: 8px;
}
.header-card {
  margin-bottom: 16px;
}
.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.page-title {
  margin: 0 0 4px 0;
  font-size: 18px;
  font-weight: 600;
}
.page-desc {
  margin: 0;
  color: #909399;
  font-size: 13px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.products-card {
  margin-top: 16px;
}
.test-result {
  min-height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
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
