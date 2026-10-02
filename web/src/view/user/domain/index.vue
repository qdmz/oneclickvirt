<template>
  <div class="domain-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('user.domain.title') }}</span>
          <el-button
            type="primary"
            @click="handleCreate"
          >
            {{ t('user.domain.addDomain') }}
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="domains"
        stripe
      >
        <el-table-column
          prop="domainName"
          :label="t('user.domain.domainName')"
        />
        <el-table-column
          prop="instanceId"
          :label="t('user.domain.instanceId')"
          min-width="110"
        />
        <el-table-column
          prop="protocol"
          :label="t('user.domain.protocol')"
          min-width="110"
        />
        <el-table-column
          prop="internalIP"
          :label="t('user.domain.internalIp')"
          min-width="110"
        />
        <el-table-column
          prop="internalPort"
          :label="t('user.domain.internalPort')"
          min-width="150"
        />
        <el-table-column
          prop="enableSSL"
          :label="t('user.domain.enableSsl')"
          min-width="120"
        >
          <template #default="{ row }">
            <el-tag
              :type="row.enableSSL ? (row.hasCert ? 'success' : 'warning') : 'info'"
              size="small"
            >
              {{ row.enableSSL ? (row.hasCert ? 'SSL' : t('user.domain.noCert')) : 'HTTP' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="status"
          :label="t('user.domain.status')"
          width="100"
        >
          <template #default="{ row }">
            <el-tag
              :type="row.status === 'active' ? 'success' : row.status === 'error' ? 'danger' : 'warning'"
              size="small"
            >
              {{ row.status === 'active' ? t('user.domain.statusActive') : row.status === 'error' ? t('user.domain.statusError') : t('user.domain.statusPending') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('user.domain.actions')"
          width="150"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              @click="handleEdit(row)"
            >
              <el-icon><Edit /></el-icon>
            </el-button>
            <el-button
              link
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
        :description="t('user.domain.noDomains')"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="isEdit ? t('user.domain.edit') : t('user.domain.addDomain')"
      width="560px"
      destroy-on-close
    >
      <!-- DNS 绑定说明 -->
      <el-alert
        v-if="!isEdit"
        type="info"
        :closable="false"
        style="margin-bottom: 16px;"
      >
        <template #title>
          <span>{{ t('user.domain.dnsGuideTitle') }}</span>
        </template>
        <div style="margin-top: 6px; line-height: 1.7;">
          <div>{{ t('user.domain.dnsStep1') }}</div>
          <div>{{ t('user.domain.dnsStep2') }}</div>
          <div>{{ t('user.domain.dnsStep3') }}</div>
        </div>
      </el-alert>

      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="130px"
      >
        <el-form-item
          :label="t('user.domain.domainName')"
          prop="domainName"
        >
          <el-input
            v-model="form.domainName"
            placeholder="app.example.com"
            :disabled="isEdit"
          />
        </el-form-item>
        <el-form-item
          v-if="!isEdit"
          :label="t('user.domain.instanceId')"
          prop="instanceId"
        >
          <el-select
            v-model="form.instanceId"
            :placeholder="t('user.domain.selectInstance')"
            filterable
            style="width: 100%"
            @change="onInstanceChange"
          >
            <el-option
              v-for="inst in userInstances"
              :key="inst.id"
              :value="inst.id"
              :label="`#${inst.id} - ${inst.name || ''}`"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="!isEdit && selectedInstancePublicIP"
          :label="t('user.domain.nodeIpLabel')"
        >
          <el-input
            :model-value="selectedInstancePublicIP"
            readonly
          />
          <div class="form-tip">
            {{ t('user.domain.nodeIpTip') }}
          </div>
        </el-form-item>
        <el-form-item
          :label="t('user.domain.internalIp')"
          prop="internalIP"
        >
          <el-input
            v-model="form.internalIP"
            placeholder="172.17.0.2"
          />
          <div class="form-tip">
            {{ t('user.domain.internalIpTip') }}
          </div>
        </el-form-item>
        <el-form-item
          :label="t('user.domain.internalPort')"
          prop="internalPort"
        >
          <el-input-number
            v-model="form.internalPort"
            :min="1"
            :max="65535"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item
          :label="t('user.domain.protocol')"
          prop="protocol"
        >
          <el-select
            v-model="form.protocol"
            style="width: 100%"
          >
            <el-option
              label="HTTP"
              value="http"
            />
            <el-option
              label="HTTPS"
              value="https"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('user.domain.enableSsl')">
          <el-switch v-model="form.enableSSL" />
          <span
            class="form-tip"
            style="margin-left: 8px;"
          >{{ t('user.domain.enableSslTip') }}</span>
        </el-form-item>
        <template v-if="form.enableSSL">
          <el-form-item :label="t('user.domain.sslCert')">
            <el-input
              v-model="form.sslCertContent"
              type="textarea"
              :rows="4"
              :placeholder="t('user.domain.sslCertPlaceholder')"
            />
            <div class="form-tip">
              {{ t('user.domain.sslCertTip') }}
            </div>
          </el-form-item>
          <el-form-item :label="t('user.domain.sslKey')">
            <el-input
              v-model="form.sslKeyContent"
              type="textarea"
              :rows="4"
              :placeholder="t('user.domain.sslKeyPlaceholder')"
            />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="handleSubmit"
        >
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Delete } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { getUserDomains, createUserDomain, updateUserDomain, deleteUserDomain } from '@/api/features'
import { getUserInstances } from '@/api/user'

const { t } = useI18n()

const domains = ref([])
const loading = ref(false)
const submitting = ref(false)
const showDialog = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const formRef = ref(null)
const userInstances = ref([])

const form = reactive({
  domainName: '',
  instanceId: null,
  protocol: 'http',
  internalIP: '',
  internalPort: 80,
  enableSSL: false,
  sslCertContent: '',
  sslKeyContent: ''
})

const formRules = {
  domainName: [{ required: true, message: () => t('user.domain.domainRequired'), trigger: 'blur' }],
  instanceId: [{ required: true, message: () => t('user.domain.instanceIdRequired'), trigger: 'change' }],
  internalIP: [{ required: true, message: () => t('user.domain.internalIPRequired'), trigger: 'blur' }],
  internalPort: [{ required: true, message: () => t('user.domain.portRequired'), trigger: 'blur' }]
}

function toArray(payload) {
  if (Array.isArray(payload)) return payload
  if (Array.isArray(payload?.list)) return payload.list
  if (Array.isArray(payload?.data)) return payload.data
  return []
}

const selectedInstancePublicIP = computed(() => {
  if (!form.instanceId) return ''
  const inst = userInstances.value.find(i => i.id === form.instanceId)
  return inst?.publicIP || inst?.publicIp || ''
})

function onInstanceChange(instanceId) {
  const inst = userInstances.value.find(i => i.id === instanceId)
  if (inst) {
    if (inst.privateIP || inst.privateIp) {
      form.internalIP = inst.privateIP || inst.privateIp
    }
  }
}

async function fetchInstances() {
  try {
    const res = await getUserInstances({ page: 1, pageSize: 999 })
    userInstances.value = res.data?.list || res.data?.data || res.data || []
  } catch (_) {
    // ignore
  }
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getUserDomains()
    if (res.code === 200) {
      domains.value = toArray(res.data)
    }
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  isEdit.value = false
  editId.value = null
  Object.assign(form, { domainName: '', instanceId: null, protocol: 'http', internalIP: '', internalPort: 80, enableSSL: false, sslCertContent: '', sslKeyContent: '' })
  showDialog.value = true
}

function handleEdit(row) {
  isEdit.value = true
  editId.value = row.id
  Object.assign(form, {
    domainName: row.domainName,
    instanceId: row.instanceId,
    protocol: row.protocol || 'http',
    internalIP: row.internalIP || '',
    internalPort: row.internalPort,
    enableSSL: row.enableSSL || false,
    sslCertContent: '',
    sslKeyContent: ''
  })
  showDialog.value = true
}

async function handleSubmit() {
  try {
    await formRef.value.validate()
  } catch (_) {
    return
  }
  submitting.value = true
  try {
    if (isEdit.value) {
      await updateUserDomain(editId.value, {
        internalIP: form.internalIP,
        internalPort: form.internalPort,
        protocol: form.protocol,
        enableSSL: form.enableSSL,
        sslCertContent: form.sslCertContent,
        sslKeyContent: form.sslKeyContent
      })
      ElMessage.success(t('user.domain.updateSuccess'))
    } else {
      await createUserDomain({
        domainName: form.domainName,
        instanceId: form.instanceId,
        protocol: form.protocol,
        internalIP: form.internalIP,
        internalPort: form.internalPort,
        enableSSL: form.enableSSL,
        sslCertContent: form.sslCertContent,
        sslKeyContent: form.sslKeyContent
      })
      ElMessage.success(t('user.domain.createSuccess'))
    }
    showDialog.value = false
    fetchData()
  } catch (error) {
    ElMessage.error(error?.message || (isEdit.value ? t('user.domain.updateFailed') : t('user.domain.createFailed')))
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(t('user.domain.confirmDelete'))
  } catch (_) {
    return
  }
  try {
    await deleteUserDomain(row.id)
    ElMessage.success(t('user.domain.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error(error?.message || t('user.domain.deleteFailed'))
  }
}

onMounted(() => {
  fetchData()
  fetchInstances()
})
</script>

<style scoped>
.domain-container { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.form-tip { font-size: 12px; color: #909399; line-height: 1.4; margin-top: 2px; }


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
