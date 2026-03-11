<template>
  <div class="domain-container">
    <!-- 配额信息 -->
    <el-card class="quota-card" shadow="never">
      <el-row :gutter="20">
        <el-col :span="8">
          <div class="quota-item">
            <span class="quota-label">已绑定域名</span>
            <span class="quota-value">{{ quota.used }}</span>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="quota-item">
            <span class="quota-label">最大配额</span>
            <span class="quota-value">{{ quota.max }}</span>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="quota-item">
            <span class="quota-label">剩余可绑</span>
            <span class="quota-value" :class="{ 'quota-warn': quota.remain <= 0 }">{{ quota.remain }}</span>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 操作栏 -->
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>域名管理</span>
          <el-button type="primary" @click="showAddDialog">
            添加域名
          </el-button>
        </div>
      </template>

      <el-table :data="domainList" v-loading="loading" border>
        <el-table-column prop="domain" label="域名" min-width="180" />
        <el-table-column label="内部地址" min-width="160">
          <template #default="{ row }">
            {{ row.internalIp }}:{{ row.internalPort }}
          </template>
        </el-table-column>
        <el-table-column prop="externalPort" label="外部端口" width="100">
          <template #default="{ row }">
            {{ row.externalPort || '自动' }}
          </template>
        </el-table-column>
        <el-table-column prop="protocol" label="协议" width="80" />
        <el-table-column label="SSL" width="70" align="center">
          <template #default="{ row }">
            <el-tag :type="row.ssl ? 'success' : 'info'" size="small">
              {{ row.ssl ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">
              {{ statusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="showEditDialog(row)">
              编辑
            </el-button>
            <el-popconfirm title="确定删除此域名？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button type="danger" link>
                  删除
                </el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchDomains"
          @current-change="fetchDomains"
        />
      </div>
    </el-card>

    <!-- 添加/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑域名' : '添加域名'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="域名" prop="domain">
          <el-input v-model="form.domain" placeholder="example.com" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="实例ID" prop="instanceId">
          <el-input-number v-model="form.instanceId" :min="1" />
        </el-form-item>
        <el-form-item label="内部IP" prop="internalIp">
          <el-input v-model="form.internalIp" placeholder="10.0.0.x" />
        </el-form-item>
        <el-form-item label="内部端口" prop="internalPort">
          <el-input-number v-model="form.internalPort" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="外部端口">
          <el-input-number v-model="form.externalPort" :min="0" :max="65535" placeholder="0=自动分配" />
          <span class="form-tip">0表示自动分配</span>
        </el-form-item>
        <el-form-item label="协议" prop="protocol">
          <el-select v-model="form.protocol">
            <el-option label="HTTP" value="http" />
            <el-option label="HTTPS" value="https" />
            <el-option label="TCP" value="tcp" />
            <el-option label="UDP" value="udp" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getDomains, createDomain, updateDomain, deleteDomain, getDomainQuota } from '@/api/domain'

const loading = ref(false)
const domainList = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const quota = reactive({ used: 0, max: 3, remain: 3 })

const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const form = reactive({
  id: null,
  domain: '',
  instanceId: null,
  internalIp: '',
  internalPort: 80,
  externalPort: 0,
  protocol: 'http'
})

const rules = {
  domain: [{ required: true, message: '请输入域名', trigger: 'blur' }],
  instanceId: [{ required: true, message: '请选择实例', trigger: 'change' }],
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
    const { data } = await getDomains({ page: page.value, pageSize: pageSize.value })
    domainList.value = data?.data?.list || []
    total.value = data?.data?.total || 0
  } finally {
    loading.value = false
  }
}

async function fetchQuota() {
  try {
    const { data } = await getDomainQuota()
    Object.assign(quota, data?.data || {})
  } catch {}
}

function showAddDialog() {
  isEdit.value = false
  Object.assign(form, { id: null, domain: '', instanceId: null, internalIp: '', internalPort: 80, externalPort: 0, protocol: 'http' })
  dialogVisible.value = true
}

function showEditDialog(row) {
  isEdit.value = true
  Object.assign(form, { id: row.id, domain: row.domain, instanceId: row.instanceId, internalIp: row.internalIp, internalPort: row.internalPort, externalPort: row.externalPort, protocol: row.protocol })
  dialogVisible.value = true
}

async function handleSubmit() {
  await formRef.value?.validate()
  submitting.value = true
  try {
    const payload = { instanceId: form.instanceId, domain: form.domain, internalIp: form.internalIp, internalPort: form.internalPort, externalPort: form.externalPort, protocol: form.protocol }
    if (isEdit.value) {
      await updateDomain(form.id, payload)
      ElMessage.success('更新成功')
    } else {
      await createDomain(payload)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchDomains()
    fetchQuota()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id) {
  try {
    await deleteDomain(id)
    ElMessage.success('删除成功')
    fetchDomains()
    fetchQuota()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || '删除失败')
  }
}

onMounted(() => {
  fetchDomains()
  fetchQuota()
})
</script>

<style scoped>
.domain-container { padding: 20px; }
.quota-card { margin-bottom: 20px; }
.quota-item { text-align: center; padding: 10px; }
.quota-label { display: block; font-size: 14px; color: #909399; margin-bottom: 8px; }
.quota-value { font-size: 28px; font-weight: 600; }
.quota-warn { color: #f56c6c; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.pagination-wrap { display: flex; justify-content: flex-end; margin-top: 20px; }
.form-tip { margin-left: 8px; color: #909399; font-size: 12px; }
</style>
