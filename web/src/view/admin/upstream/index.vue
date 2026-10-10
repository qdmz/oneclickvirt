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
</style>
