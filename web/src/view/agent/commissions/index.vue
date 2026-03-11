<template>
  <div class="commissions-page">
    <div class="page-header">
      <h1>佣金记录</h1>
    </div>

    <el-card>
      <div class="toolbar">
        <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 150px" @change="loadData">
          <el-option label="全部" :value="null" />
          <el-option label="待结算" :value="0" />
          <el-option label="已结算" :value="1" />
          <el-option label="已取消" :value="2" />
        </el-select>
        <el-button type="primary" @click="$router.push('/agent/wallet')">钱包管理</el-button>
      </div>

      <el-table v-loading="loading" :data="commissions" stripe>
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="rate" label="佣金比例" width="100" align="center">
          <template #default="{ row }">
            <span v-if="row.rate">{{ row.rate }}%</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="130" align="right">
          <template #default="{ row }">
            <span :style="{ color: row.amount >= 0 ? '#67c23a' : '#f56c6c', fontWeight: '600' }">
              {{ row.amount >= 0 ? '+' : '' }}¥{{ (row.amount / 100).toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="时间" width="180">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column prop="settledAt" label="结算时间" width="180">
          <template #default="{ row }">{{ row.settledAt ? formatDate(row.settledAt) : '-' }}</template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCommissions } from '@/api/agent'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const commissions = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const statusFilter = ref(null)

const statusTagType = (s) => ({ 0: 'warning', 1: 'success', 2: 'info' }[s] || 'info')
const statusText = (s) => ({ 0: '待结算', 1: '已结算', 2: '已取消' }[s] || '未知')
const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try {
    const res = await getCommissions({
      page: page.value,
      pageSize: pageSize.value,
      ...(statusFilter.value !== null ? { status: statusFilter.value } : {})
    })
    if (res.data?.code === 0) {
      commissions.value = res.data.data?.list || []
      total.value = res.data.data?.total || 0
    }
  } catch { ElMessage.error('加载失败') }
  finally { loading.value = false }
}

onMounted(loadData)
</script>

<style scoped>
.commissions-page { padding: 20px; }
.page-header { margin-bottom: 20px; }
.page-header h1 { font-size: 24px; }
.toolbar { display: flex; justify-content: space-between; margin-bottom: 16px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
