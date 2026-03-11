<template>
  <div class="admin-domain-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>域名管理</span>
          <el-button type="success" @click="handleSyncDNS" :loading="syncing">
            同步DNS
          </el-button>
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
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { adminGetDomains, adminDeleteDomain, syncDNS } from '@/api/domain'

const loading = ref(false)
const syncing = ref(false)
const domainList = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filters = reactive({ domain: '', userId: undefined, status: undefined })

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
    const { data } = await adminGetDomains(params)
    domainList.value = data?.data?.list || []
    total.value = data?.data?.total || 0
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

onMounted(fetchDomains)
</script>

<style scoped>
.admin-domain-container { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.filter-form { margin-bottom: 16px; }
.pagination-wrap { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
