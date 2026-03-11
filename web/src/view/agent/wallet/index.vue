<template>
  <div class="agent-wallet">
    <div class="page-header">
      <h1>钱包管理</h1>
    </div>

    <!-- 余额卡片 -->
    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="balance-info">
            <div class="balance-label">可用余额</div>
            <div class="balance-value">¥{{ (wallet.balance / 100).toFixed(2) }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="balance-info">
            <div class="balance-label">累计佣金</div>
            <div class="balance-value">¥{{ (wallet.totalCommission / 100).toFixed(2) }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="balance-info">
            <div class="balance-label">累计提现</div>
            <div class="balance-value">¥{{ (wallet.totalWithdrawn / 100).toFixed(2) }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 提现 -->
    <el-card style="margin-bottom: 20px;">
      <template #header><span>提现操作</span></template>
      <div style="display: flex; gap: 12px; align-items: center; max-width: 400px;">
        <el-input-number
          v-model="withdrawAmount"
          :min="1"
          :step="100"
          :precision="0"
          placeholder="提现金额（分）"
          controls-position="right"
          style="flex: 1;"
        />
        <el-button type="primary" :loading="withdrawing" @click="handleWithdraw">提现</el-button>
      </div>
      <p style="margin-top: 8px; color: #909399; font-size: 13px;">
        单位：分（1元 = 100分）。当前可用：¥{{ (wallet.balance / 100).toFixed(2) }}
      </p>
    </el-card>

    <!-- 交易记录 -->
    <el-card>
      <template #header><span>交易记录</span></template>
      <el-table v-loading="loading" :data="transactions" stripe>
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="amount" label="金额" width="130" align="right">
          <template #default="{ row }">
            <span :style="{ color: row.amount >= 0 ? '#67c23a' : '#f56c6c', fontWeight: '600' }">
              {{ row.amount >= 0 ? '+' : '' }}¥{{ (row.amount / 100).toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="{ 0: 'warning', 1: 'success', 2: 'info' }[row.status]" size="small">
              {{ { 0: '待结算', 1: '已结算', 2: '已取消' }[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="时间" width="180">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadTransactions"
          @current-change="loadTransactions"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAgentWallet, getWalletTransactions, withdraw } from '@/api/agent'
import { ElMessage, ElMessageBox } from 'element-plus'

const wallet = ref({ balance: 0, totalCommission: 0, totalWithdrawn: 0 })
const transactions = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const withdrawAmount = ref(0)
const withdrawing = ref(false)

const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'

const loadWallet = async () => {
  try {
    const res = await getAgentWallet()
    if (res.data?.code === 0) wallet.value = res.data.data
  } catch { /* */ }
}

const loadTransactions = async () => {
  loading.value = true
  try {
    const res = await getWalletTransactions({ page: page.value, pageSize: pageSize.value })
    if (res.data?.code === 0) {
      transactions.value = res.data.data?.list || []
      total.value = res.data.data?.total || 0
    }
  } catch { ElMessage.error('加载失败') }
  finally { loading.value = false }
}

const handleWithdraw = async () => {
  if (!withdrawAmount.value || withdrawAmount.value <= 0) {
    ElMessage.warning('请输入提现金额')
    return
  }
  try {
    await ElMessageBox.confirm(
      `确定提现 ¥${(withdrawAmount.value / 100).toFixed(2)} 吗？`,
      '确认提现',
      { type: 'warning' }
    )
    withdrawing.value = true
    await withdraw(withdrawAmount.value)
    ElMessage.success('提现申请已提交')
    withdrawAmount.value = 0
    await loadWallet()
    await loadTransactions()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(error?.response?.data?.message || '提现失败')
  } finally {
    withdrawing.value = false
  }
}

onMounted(() => {
  loadWallet()
  loadTransactions()
})
</script>

<style scoped>
.agent-wallet { padding: 20px; }
.page-header { margin-bottom: 20px; }
.page-header h1 { font-size: 24px; }
.balance-info { text-align: center; }
.balance-label { font-size: 14px; color: #909399; margin-bottom: 8px; }
.balance-value { font-size: 28px; font-weight: 600; color: #16a34a; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
