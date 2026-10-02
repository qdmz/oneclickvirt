<template>
  <div class="domain-mgmt-container">
    <el-tabs v-model="activeTab">
      <el-tab-pane
        :label="t('admin.domain.instanceDomainBindings')"
        name="domains"
      >
        <div class="domain-panel">
          <div class="filter-bar">
            <el-input
              v-model="filters.keyword"
              :placeholder="t('admin.domain.keywordPlaceholder')"
              clearable
              class="filter-input"
              @keyup.enter="handleSearch"
              @clear="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-select
              v-model="filters.status"
              :placeholder="t('admin.domain.status')"
              clearable
              class="filter-select"
              @change="handleSearch"
              @clear="handleSearch"
            >
              <el-option
                :label="t('admin.domain.statusActive')"
                value="active"
              />
              <el-option
                :label="t('admin.domain.statusPending')"
                value="pending"
              />
              <el-option
                :label="t('admin.domain.statusError')"
                value="error"
              />
            </el-select>
            <el-select
              v-model="filters.providerId"
              :placeholder="t('admin.domain.providerName')"
              clearable
              filterable
              class="filter-select"
              @change="handleSearch"
              @clear="handleSearch"
            >
              <el-option
                v-for="provider in providers"
                :key="provider.id"
                :label="provider.name"
                :value="provider.id"
              />
            </el-select>
            <el-input-number
              v-model="filters.userId"
              :placeholder="t('admin.domain.userId')"
              :min="1"
              :controls="false"
              class="filter-number"
              @change="handleSearch"
            />
            <el-input-number
              v-model="filters.instanceId"
              :placeholder="t('admin.domain.instanceId')"
              :min="1"
              :controls="false"
              class="filter-number"
              @change="handleSearch"
            />
            <el-button
              type="primary"
              @click="handleSearch"
            >
              <el-icon><Search /></el-icon>
              <span>{{ t('common.search') }}</span>
            </el-button>
            <el-button @click="resetFilters">
              {{ t('common.reset') }}
            </el-button>
            <el-button
              :loading="syncingProxies"
              @click="handleSyncProxies"
            >
              <el-icon><Refresh /></el-icon>
              <span>{{ t('admin.domain.syncProxies') }}</span>
            </el-button>
          </div>

          <el-table
            v-loading="loading"
            :data="domains"
            stripe
            row-key="id"
          >
            <el-table-column
              prop="id"
              label="ID"
              width="72"
            />
            <el-table-column
              :label="t('admin.domain.domainName')"
              min-width="220"
            >
              <template #default="{ row }">
                <div class="domain-main">
                  <span class="domain-name">{{ row.domainName }}</span>
                  <el-tag
                    v-if="row.enableSSL"
                    size="small"
                    :type="row.hasCert ? 'success' : 'warning'"
                  >
                    {{ row.hasCert ? 'SSL' : t('admin.domain.certMissing') }}
                  </el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.domain.owner')"
              min-width="150"
            >
              <template #default="{ row }">
                <div>{{ row.username || '-' }}</div>
                <div class="sub-text">
                  ID: {{ row.userId }}
                  <span v-if="row.userNickname"> / {{ row.userNickname }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.domain.instance')"
              min-width="170"
            >
              <template #default="{ row }">
                <div>{{ row.instanceName || '-' }}</div>
                <div class="sub-text">
                  ID: {{ row.instanceId }}
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.domain.provider')"
              min-width="170"
            >
              <template #default="{ row }">
                <div>{{ row.providerName || '-' }}</div>
                <div class="sub-text">
                  {{ row.providerType || '-' }} / ID: {{ row.providerId }}
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.domain.target')"
              min-width="170"
            >
              <template #default="{ row }">
                <el-tag
                  size="small"
                  type="info"
                >
                  {{ row.protocol || 'http' }}
                </el-tag>
                <span class="target-text">{{ formatEndpointHostPort(row.internalIP, row.internalPort) }}</span>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.domain.status')"
              width="130"
            >
              <template #default="{ row }">
                <div class="status-cell">
                  <el-tag
                    :type="getStatusTagType(row.status)"
                    size="small"
                  >
                    {{ getStatusText(row.status) }}
                  </el-tag>
                  <el-tooltip
                    v-if="row.errorMsg"
                    :content="row.errorMsg"
                    placement="top"
                  >
                    <el-icon class="error-icon">
                      <Warning />
                    </el-icon>
                  </el-tooltip>
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('common.createdAt')"
              width="170"
            >
              <template #default="{ row }">
                {{ formatDateTime(row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.domain.actions')"
              width="230"
              fixed="right"
            >
              <template #default="{ row }">
                <el-button
                  size="small"
                  @click="handleEditDomain(row)"
                >
                  <el-icon><Edit /></el-icon>
                  <span>{{ t('common.edit') }}</span>
                </el-button>
                <el-button
                  size="small"
                  :loading="syncingDomainId === row.id"
                  @click="handleSyncDomain(row)"
                >
                  <el-icon><Connection /></el-icon>
                  <span>{{ t('admin.domain.retrySync') }}</span>
                </el-button>
                <el-button
                  size="small"
                  type="danger"
                  @click="handleDelete(row)"
                >
                  <el-icon><Delete /></el-icon>
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-empty
            v-if="!loading && domains.length === 0"
            :description="t('admin.domain.noDomains')"
          />

          <el-pagination
            v-if="pagination.total > 0"
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            class="pagination"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handlePageSizeChange"
            @current-change="fetchData"
          />
        </div>
      </el-tab-pane>

      <el-tab-pane
        :label="t('admin.domain.providerConfig')"
        name="config"
      >
        <div class="domain-panel">
          <el-table
            v-loading="configLoading"
            :data="providers"
            stripe
          >
            <el-table-column
              prop="id"
              label="ID"
              width="72"
            />
            <el-table-column
              prop="name"
              :label="t('admin.domain.providerName')"
              min-width="180"
            />
            <el-table-column
              prop="type"
              :label="t('admin.domain.providerType')"
              width="120"
            />
            <el-table-column
              :label="t('admin.domain.domainBindingEnabled')"
              min-width="160"
            >
              <template #default="{ row }">
                <el-tag
                  :type="getProviderConfig(row.id)?.enabled ? 'success' : 'info'"
                  size="small"
                >
                  {{ getProviderConfig(row.id)?.enabled ? t('common.enabled') : t('common.disabled') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.domain.maxDomainsPerUser')"
              min-width="200"
            >
              <template #default="{ row }">
                {{ getProviderConfig(row.id)?.maxDomainsPerUser ?? 3 }}
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.domain.allowedSuffixes')"
              min-width="220"
            >
              <template #default="{ row }">
                <span>{{ getProviderConfig(row.id)?.allowedSuffixes || t('admin.domain.noSuffixLimit') }}</span>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.domain.actions')"
              width="110"
              fixed="right"
            >
              <template #default="{ row }">
                <el-button
                  size="small"
                  @click="handleEditConfig(row)"
                >
                  <el-icon><Edit /></el-icon>
                  <span>{{ t('common.edit') }}</span>
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-dialog
      v-model="showDomainDialog"
      :title="t('admin.domain.editDomain')"
      width="680px"
      destroy-on-close
      @closed="resetDomainForm"
    >
      <el-form
        ref="domainFormRef"
        :model="domainForm"
        :rules="domainRules"
        label-width="140px"
      >
        <el-form-item
          :label="t('admin.domain.domainName')"
          prop="domainName"
        >
          <el-input
            v-model.trim="domainForm.domainName"
            placeholder="app.example.com"
          />
        </el-form-item>
        <el-form-item
          :label="t('admin.domain.instance')"
          prop="instanceId"
        >
          <el-select
            v-model="domainForm.instanceId"
            :loading="instanceLoading"
            :placeholder="t('admin.domain.selectInstance')"
            filterable
            style="width: 100%;"
          >
            <el-option
              v-for="instance in instanceOptions"
              :key="instance.id"
              :label="buildInstanceLabel(instance)"
              :value="instance.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="t('admin.domain.internalIP')"
          prop="internalIP"
        >
          <el-input
            v-model.trim="domainForm.internalIP"
            placeholder="127.0.0.1"
          />
        </el-form-item>
        <el-form-item
          :label="t('admin.domain.internalPort')"
          prop="internalPort"
        >
          <el-input-number
            v-model="domainForm.internalPort"
            :min="1"
            :max="65535"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item :label="t('admin.domain.protocol')">
          <el-radio-group v-model="domainForm.protocol">
            <el-radio-button value="http">
              HTTP
            </el-radio-button>
            <el-radio-button value="https">
              HTTPS
            </el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('admin.domain.enableSsl')">
          <el-switch v-model="domainForm.enableSSL" />
        </el-form-item>
        <el-form-item
          v-if="domainForm.enableSSL"
          :label="t('admin.domain.sslCert')"
        >
          <el-input
            v-model="domainForm.sslCertContent"
            type="textarea"
            :rows="4"
            :placeholder="t('admin.domain.sslCertPlaceholder')"
          />
        </el-form-item>
        <el-form-item
          v-if="domainForm.enableSSL"
          :label="t('admin.domain.sslKey')"
        >
          <el-input
            v-model="domainForm.sslKeyContent"
            type="textarea"
            :rows="4"
            :placeholder="t('admin.domain.sslKeyPlaceholder')"
          />
        </el-form-item>
        <el-form-item
          v-if="domainForm.hasCert"
          :label="t('admin.domain.clearCert')"
        >
          <el-switch v-model="domainForm.clearCert" />
          <span class="form-tip">{{ t('admin.domain.clearCertTip') }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDomainDialog = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="savingDomain"
          @click="handleSaveDomain"
        >
          {{ t('common.save') }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="showConfigDialog"
      :title="t('admin.domain.editConfig')"
      width="560px"
      destroy-on-close
    >
      <el-form
        ref="configFormRef"
        :model="configForm"
        label-width="150px"
      >
        <el-form-item :label="t('admin.domain.enabled')">
          <el-switch v-model="configForm.enabled" />
        </el-form-item>
        <el-form-item :label="t('admin.domain.maxDomainsPerUser')">
          <el-input-number
            v-model="configForm.maxDomainsPerUser"
            :min="1"
            :max="100"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="t('admin.domain.dnsType')">
          <el-select
            v-model="configForm.dnsType"
            style="width: 100%"
          >
            <el-option
              label="Hosts"
              value="hosts"
            />
            <el-option
              label="Nginx"
              value="nginx"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('admin.domain.allowedSuffixes')">
          <el-input
            v-model="configForm.allowedSuffixes"
            :placeholder="t('admin.domain.allowedSuffixesPlaceholder')"
          />
          <div class="form-tip">
            {{ t('admin.domain.allowedSuffixesTip') }}
          </div>
        </el-form-item>
        <el-form-item :label="t('admin.domain.nginxConfigPath')">
          <el-input
            v-model="configForm.nginxConfigPath"
            placeholder="/etc/nginx/conf.d"
          />
        </el-form-item>
        <el-form-item :label="t('admin.domain.nginxReloadCmd')">
          <el-input
            v-model="configForm.nginxReloadCmd"
            placeholder="systemctl reload nginx"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showConfigDialog = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="savingConfig"
          @click="handleSaveConfig"
        >
          {{ t('common.save') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { nextTick, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Connection, Delete, Edit, Refresh, Search, Warning } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import {
  adminDeleteDomain,
  adminGetDomains,
  adminSyncDomainProxy,
  adminSyncDomainProxies,
  adminUpdateDomain,
  getDomainConfig,
  updateDomainConfig
} from '@/api/features'
import { getAllInstances, getProviderList } from '@/api/admin'
import { formatEndpointHostPort } from '@/utils/endpoint'

const { t, locale } = useI18n()

const activeTab = ref('domains')
const domains = ref([])
const loading = ref(false)
const syncingProxies = ref(false)
const syncingDomainId = ref(null)
const providers = ref([])
const providerConfigs = ref({})
const configLoading = ref(false)
const showConfigDialog = ref(false)
const savingConfig = ref(false)
const configFormRef = ref(null)
const editingProviderId = ref(null)
const showDomainDialog = ref(false)
const savingDomain = ref(false)
const domainFormRef = ref(null)
const instanceOptions = ref([])
const instanceLoading = ref(false)

const filters = reactive({
  keyword: '',
  status: '',
  providerId: null,
  userId: null,
  instanceId: null
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const domainForm = reactive({
  id: null,
  domainName: '',
  instanceId: null,
  internalIP: '',
  internalPort: 80,
  protocol: 'http',
  enableSSL: false,
  sslCertContent: '',
  sslKeyContent: '',
  hasCert: false,
  clearCert: false
})

const domainRules = {
  domainName: [{ required: true, message: () => t('admin.domain.domainRequired'), trigger: 'blur' }],
  instanceId: [{ required: true, message: () => t('admin.domain.instanceRequired'), trigger: 'change' }],
  internalIP: [{ required: true, message: () => t('admin.domain.internalIPRequired'), trigger: 'blur' }],
  internalPort: [{ required: true, message: () => t('admin.domain.portRequired'), trigger: 'change' }]
}

const configForm = reactive({
  enabled: false,
  maxDomainsPerUser: 3,
  dnsType: 'hosts',
  dnsConfigPath: '',
  nginxConfigPath: '/etc/nginx/conf.d',
  nginxReloadCmd: 'systemctl reload nginx',
  allowedSuffixes: ''
})

function getProviderConfig(providerId) {
  return providerConfigs.value[providerId] || null
}

function buildListParams() {
  return {
    page: pagination.page,
    pageSize: pagination.pageSize,
    keyword: filters.keyword || undefined,
    status: filters.status || undefined,
    providerId: filters.providerId || undefined,
    userId: filters.userId || undefined,
    instanceId: filters.instanceId || undefined
  }
}

async function fetchData() {
  loading.value = true
  try {
    const res = await adminGetDomains(buildListParams())
    if (res.code === 200) {
      const data = res.data || {}
      domains.value = Array.isArray(data) ? data : (data.list || [])
      pagination.total = data.total ?? domains.value.length
    }
  } finally {
    loading.value = false
  }
}

async function fetchProviders() {
  configLoading.value = true
  try {
    const res = await getProviderList({ page: 1, pageSize: 999 })
    if (res.code === 200) {
      providers.value = res.data?.list || []
      const configs = {}
      await Promise.all(providers.value.map(async (provider) => {
        try {
          const configRes = await getDomainConfig(provider.id)
          if (configRes.code === 200) {
            configs[provider.id] = configRes.data
          }
        } catch {
          configs[provider.id] = null
        }
      }))
      providerConfigs.value = configs
    }
  } finally {
    configLoading.value = false
  }
}

async function fetchInstances() {
  instanceLoading.value = true
  try {
    const res = await getAllInstances({ page: 1, pageSize: 500 })
    if (res.code === 200) {
      instanceOptions.value = res.data?.list || []
    }
  } catch (error) {
    ElMessage.error(error?.message || t('admin.domain.loadInstancesFailed'))
  } finally {
    instanceLoading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  fetchData()
}

function resetFilters() {
  Object.assign(filters, {
    keyword: '',
    status: '',
    providerId: null,
    userId: null,
    instanceId: null
  })
  handleSearch()
}

function handlePageSizeChange(size) {
  pagination.pageSize = size
  pagination.page = 1
  fetchData()
}

function handleEditConfig(provider) {
  editingProviderId.value = provider.id
  const existing = getProviderConfig(provider.id)
  Object.assign(configForm, {
    enabled: existing?.enabled ?? false,
    maxDomainsPerUser: existing?.maxDomainsPerUser ?? 3,
    dnsType: existing?.dnsType ?? 'hosts',
    dnsConfigPath: existing?.dnsConfigPath ?? '',
    nginxConfigPath: existing?.nginxConfigPath ?? '/etc/nginx/conf.d',
    nginxReloadCmd: existing?.nginxReloadCmd ?? 'systemctl reload nginx',
    allowedSuffixes: existing?.allowedSuffixes ?? ''
  })
  showConfigDialog.value = true
}

async function handleSaveConfig() {
  savingConfig.value = true
  try {
    await updateDomainConfig(editingProviderId.value, { ...configForm })
    ElMessage.success(t('admin.domain.configUpdated'))
    showConfigDialog.value = false
    await fetchProviders()
  } catch (error) {
    ElMessage.error(error?.message || t('admin.domain.saveFailed'))
  } finally {
    savingConfig.value = false
  }
}

async function handleEditDomain(row) {
  if (instanceOptions.value.length === 0) {
    await fetchInstances()
  }
  Object.assign(domainForm, {
    id: row.id,
    domainName: row.domainName,
    instanceId: row.instanceId,
    internalIP: row.internalIP,
    internalPort: row.internalPort,
    protocol: row.protocol || 'http',
    enableSSL: Boolean(row.enableSSL),
    sslCertContent: '',
    sslKeyContent: '',
    hasCert: Boolean(row.hasCert),
    clearCert: false
  })
  showDomainDialog.value = true
  await nextTick()
  domainFormRef.value?.clearValidate()
}

function resetDomainForm() {
  Object.assign(domainForm, {
    id: null,
    domainName: '',
    instanceId: null,
    internalIP: '',
    internalPort: 80,
    protocol: 'http',
    enableSSL: false,
    sslCertContent: '',
    sslKeyContent: '',
    hasCert: false,
    clearCert: false
  })
}

async function handleSaveDomain() {
  try {
    await domainFormRef.value?.validate()
  } catch {
    return
  }
  savingDomain.value = true
  try {
    const payload = {
      domainName: domainForm.domainName,
      instanceId: domainForm.instanceId,
      internalIP: domainForm.internalIP,
      internalPort: domainForm.internalPort,
      protocol: domainForm.protocol,
      enableSSL: domainForm.enableSSL,
      clearCert: domainForm.clearCert
    }
    if (domainForm.sslCertContent || domainForm.sslKeyContent) {
      payload.sslCertContent = domainForm.sslCertContent
      payload.sslKeyContent = domainForm.sslKeyContent
    }
    await adminUpdateDomain(domainForm.id, payload)
    ElMessage.success(t('admin.domain.updateDomainSuccess'))
    showDomainDialog.value = false
    await fetchData()
  } catch (error) {
    if (error !== false) {
      ElMessage.error(error?.message || t('admin.domain.updateDomainFailed'))
    }
  } finally {
    savingDomain.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(t('admin.domain.confirmDeleteDomain', { domain: row.domainName }))
    await adminDeleteDomain(row.id)
    ElMessage.success(t('admin.domain.deleteSuccess'))
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error?.message || t('admin.domain.deleteFailed'))
    }
  }
}

async function handleSyncDomain(row) {
  syncingDomainId.value = row.id
  try {
    await adminSyncDomainProxy(row.id)
    ElMessage.success(t('admin.domain.retrySyncSuccess'))
    await fetchData()
  } catch (error) {
    ElMessage.error(error?.message || t('admin.domain.retrySyncFailed'))
  } finally {
    syncingDomainId.value = null
  }
}

async function handleSyncProxies() {
  syncingProxies.value = true
  try {
    const res = await adminSyncDomainProxies()
    if (res.code === 200) {
      const data = res.data || {}
      ElMessage.success(t('admin.domain.syncProxiesSuccess', {
        success: data.success || 0,
        failed: data.failed || 0,
        skipped: data.skipped || 0,
        removed: data.removed || 0
      }))
      await fetchData()
    }
  } catch (error) {
    ElMessage.error(error?.message || t('admin.domain.syncProxiesFailed'))
  } finally {
    syncingProxies.value = false
  }
}

function getStatusTagType(status) {
  if (status === 'active') return 'success'
  if (status === 'error') return 'danger'
  return 'warning'
}

function getStatusText(status) {
  if (status === 'active') return t('admin.domain.statusActive')
  if (status === 'error') return t('admin.domain.statusError')
  return t('admin.domain.statusPending')
}

function formatDateTime(value) {
  if (!value) return '-'
  return new Date(value).toLocaleString(locale.value)
}

function buildInstanceLabel(instance) {
  const parts = [`#${instance.id}`, instance.name]
  if (instance.providerName) parts.push(instance.providerName)
  if (instance.userName) parts.push(instance.userName)
  return parts.filter(Boolean).join(' / ')
}

onMounted(async () => {
  await Promise.all([fetchData(), fetchProviders()])
})
</script>

<style scoped>
.domain-mgmt-container {
  padding: 20px;
}

.domain-panel {
  background: var(--bg-color-primary);
  border: 1px solid var(--border-color-light);
  border-radius: 8px;
  padding: 16px;
}

.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
  margin-bottom: 16px;
}

.filter-input {
  width: 260px;
}

.filter-select {
  width: 180px;
}

.filter-number {
  width: 140px;
}

.domain-main,
.status-cell {
  display: flex;
  gap: 8px;
  align-items: center;
  min-width: 0;
}

.domain-name {
  overflow: hidden;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sub-text {
  color: var(--text-color-secondary);
  font-size: 12px;
  line-height: 1.4;
}

.target-text {
  margin-left: 8px;
}

.error-icon {
  color: var(--el-color-danger);
}

.pagination {
  justify-content: flex-end;
  margin-top: 16px;
}

.form-tip {
  margin-left: 10px;
  color: var(--text-color-secondary);
  font-size: 12px;
}

@media (max-width: 768px) {
  .domain-mgmt-container {
    padding: 12px;
  }

  .filter-input,
  .filter-select,
  .filter-number {
    width: 100%;
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
