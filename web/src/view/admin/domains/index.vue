<template>
  <div class="admin-domain-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>域名管理</span>
          <div class="header-actions">
            <el-button type="primary" @click="handleAddDomain">
              添加域名
            </el-button>
            <el-button type="success" @click="handleSyncDNS" :loading="syncing">
              同步DNS
            </el-button>
          </div>
        </div>
      </template>

      <!-- 筛选 -->
      <el-form :inline="true" class="filter-form">
        <el-form-item label="域名">
          <el-input v-model="filters.domain" placeholder="搜索域名" clearable @clear="fetchDomains" />
        </el-form-item>
        <el-form-item label="用户ID">
          <el-input-number v-model="filters.userId" :min="0" clearable @change="fetchDomains" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部" clearable @change="fetchDomains">
            <el-option label="待配置" :value="0" />
            <el-option label="正常" :value="1" />
            <el-option label="配置失败" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchDomains">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="domainList" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="userId" label="用户ID" width="80" />
        <el-table-column prop="domain" label="域名" min-width="180" />
        <el-table-column label="内部地址" min-width="140">
          <template #default="{ row }">
            {{ row.internalIp }}:{{ row.internalPort }}
          </template>
        </el-table-column>
        <el-table-column prop="externalPort" label="外部端口" width="90" />
        <el-table-column prop="protocol" label="协议" width="80" />
        <el-table-column label="SSL" width="60" align="center">
          <template #default="{ row }">
            <el-tag :type="row.ssl ? 'success' : 'info'" size="small">{{ row.ssl ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEditDomain(row)">
              编辑
            </el-button>
            <el-popconfirm title="确定删除此域名？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10, 20, 50]" layout="total, sizes, prev, pager, next" @size-change="fetchDomains" @current-change="fetchDomains" />
      </div>
    </el-card>

    <!-- 添加/编辑域名对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="domainForm" :rules="domainRules" ref="domainFormRef" label-width="100px">
        <el-form-item label="用户ID" prop="userId">
          <el-input-number v-model="domainForm.userId" :min="1" />
        </el-form-item>
        <el-form-item label="实例ID" prop="instanceId">
          <el-input-number v-model="domainForm.instanceId" :min="1" />
        </el-form-item>
        <el-form-item label="域名" prop="domain">
          <el-input v-model="domainForm.domain" placeholder="请输入域名" />
        </el-form-item>
        <el-form-item label="内部IP" prop="internalIp">
          <el-input v-model="domainForm.internalIp" placeholder="请输入内部IP" />
        </el-form-item>
        <el-form-item label="内部端口" prop="internalPort">
          <el-input-number v-model="domainForm.internalPort" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="协议" prop="protocol">
          <el-select v-model="domainForm.protocol" placeholder="请选择协议">
            <el-option label="http" value="http" />
            <el-option label="https" value="https" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="domainForm.status" placeholder="请选择状态">
            <el-option label="待配置" value="0" />
            <el-option label="正常" value="1" />
            <el-option label="配置失败" value="2" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { adminGetDomains, adminDeleteDomain, adminCreateDomain, adminUpdateDomain, syncDNS } from '@/api/domain'

const loading = ref(false)
const syncing = ref(false)
const domainList = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filters = reactive({ domain: '', userId: undefined, status: undefined })

// 对话框相关
const dialogVisible = ref(false)
const dialogTitle = ref('添加域名')
const domainFormRef = ref(null)
const editingDomain = ref(null)

// 表单数据
const domainForm = reactive({
  userId: '',
  instanceId: '',
  domain: '',
  internalIp: '',
  internalPort: '',
  protocol: 'http',
  status: 0
})

// 表单验证规则
const domainRules = {
  userId: [{ required: true, message: '请输入用户ID', trigger: 'blur' }],
  instanceId: [{ required: true, message: '请输入实例ID', trigger: 'blur' }],
  domain: [{ required: true, message: '请输入域名', trigger: 'blur' }],
  internalIp: [{ required: true, message: '请输入内部IP', trigger: 'blur' }],
  internalPort: [{ required: true, message: '请输入内部端口', trigger: 'blur' }],
  protocol: [{ required: true, message: '请选择协议', trigger: 'change' }]
}

const statusMap = { 0: ['待配置', 'warning'], 1: ['正常', 'success'], 2: ['配置失败', 'danger'] }
const statusText = (s) => statusMap[s]?.[0] || '未知'
const statusType = (s) => statusMap[s]?.[1] || 'info'

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleString('zh-CN')
}

async function fetchDomains() {
  loading.value = true
  try {
    const params = { page: page.value, pageSize: pageSize.value }
    if (filters.domain) params.domain = filters.domain
    if (filters.userId) params.userId = filters.userId
    if (filters.status !== undefined && filters.status !== null) params.status = filters.status
    const res = await adminGetDomains(params)
    if (res.code === 0) {
      domainList.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } finally {
    loading.value = false
  }
}

async function handleDelete(id) {
  try {
    await adminDeleteDomain(id)
    ElMessage.success('删除成功')
    fetchDomains()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || '删除失败')
  }
}

async function handleSyncDNS() {
  syncing.value = true
  try {
    await syncDNS()
    ElMessage.success('DNS同步完成')
  } catch (e) {
    ElMessage.error('DNS同步失败')
  } finally {
    syncing.value = false
  }
}

// 添加域名
function handleAddDomain() {
  editingDomain.value = null
  dialogTitle.value = '添加域名'
  Object.assign(domainForm, {
    userId: '',
    instanceId: '',
    domain: '',
    internalIp: '',
    internalPort: '',
    protocol: 'http',
    status: 0
  })
  dialogVisible.value = true
}

// 编辑域名
function handleEditDomain(row) {
  editingDomain.value = row
  dialogTitle.value = '编辑域名'
  Object.assign(domainForm, {
    userId: row.userId,
    instanceId: row.instanceId,
    domain: row.domain,
    internalIp: row.internalIp,
    internalPort: row.internalPort,
    protocol: row.protocol,
    status: row.status
  })
  dialogVisible.value = true
}

// 提交表单
async function handleSubmit() {
  if (!domainFormRef.value) return
  
  try {
    await domainFormRef.value.validate()
    
    if (editingDomain.value) {
      // 编辑
      await adminUpdateDomain(editingDomain.value.id, domainForm)
      ElMessage.success('更新成功')
    } else {
      // 添加
      await adminCreateDomain(domainForm)
      ElMessage.success('添加成功')
    }
    
    dialogVisible.value = false
    fetchDomains()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || '操作失败')
  }
}

onMounted(fetchDomains)
</script>

<style scoped>
.admin-domain-container { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; gap: 10px; }
.filter-form { margin-bottom: 16px; }
.pagination-wrap { display: flex; justify-content: flex-end; margin-top: 20px; }
.dialog-footer { display: flex; justify-content: flex-end; gap: 10px; }
</style>
