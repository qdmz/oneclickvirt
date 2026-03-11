<template>
  <div class="agent-dashboard">
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="8" animated />
    </div>

    <div v-else>
      <div class="dashboard-header">
        <h1>代理商仪表盘</h1>
        <p>管理您的子用户和佣金收益</p>
      </div>

      <!-- 统计卡片 -->
      <el-row :gutter="20" class="stats-row">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-icon" style="background: #ecf5ff;">
              <el-icon :size="28" color="#409eff"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.subUserCount }}</div>
              <div class="stat-label">子用户总数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-icon" style="background: #f0f9eb;">
              <el-icon :size="28" color="#67c23a"><UserFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.activeUserCount }}</div>
              <div class="stat-label">活跃用户</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-icon" style="background: #fdf6ec;">
              <el-icon :size="28" color="#e6a23c"><Coin /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">¥{{ (stats.monthCommission / 100).toFixed(2) }}</div>
              <div class="stat-label">本月佣金</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-icon" style="background: #fef0f0;">
              <el-icon :size="28" color="#f56c6c"><Wallet /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">¥{{ (stats.totalCommission / 100).toFixed(2) }}</div>
              <div class="stat-label">累计佣金</div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 最近佣金记录 -->
      <el-card class="content-card">
        <template #header>
          <div class="card-header">
            <span>最近佣金记录</span>
            <el-button size="small" @click="$router.push('/agent/commissions')">
              查看全部
            </el-button>
          </div>
        </template>
        <el-table :data="recentCommissions" stripe style="width: 100%">
          <el-table-column prop="description" label="描述" min-width="200" />
          <el-table-column prop="amount" label="金额" width="120" align="right">
            <template #default="{ row }">
              <span :style="{ color: row.amount >= 0 ? '#67c23a' : '#f56c6c' }">
                ¥{{ (row.amount / 100).toFixed(2) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="statusTagType(row.status)" size="small">
                {{ statusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { User, UserFilled, Coin, Wallet } from '@element-plus/icons-vue'
import { getAgentStatistics, getCommissions } from '@/api/agent'
import { ElMessage } from 'element-plus'

const loading = ref(true)
const stats = ref({
  subUserCount: 0,
  activeUserCount: 0,
  monthCommission: 0,
  totalCommission: 0,
  totalWithdrawn: 0,
  balance: 0
})
const recentCommissions = ref([])

const statusTagType = (status) => {
  const map = { 0: 'warning', 1: 'success', 2: 'info' }
  return map[status] || 'info'
}

const statusText = (status) => {
  const map = { 0: '待结算', 1: '已结算', 2: '已取消' }
  return map[status] || '未知'
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const loadData = async () => {
  try {
    const [statsRes, commRes] = await Promise.all([
      getAgentStatistics(),
      getCommissions({ page: 1, pageSize: 5 })
    ])
    if (statsRes.data?.code === 0) {
      stats.value = statsRes.data.data
    }
    if (commRes.data?.code === 0) {
      recentCommissions.value = commRes.data.data?.list || []
    }
  } catch (error) {
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.agent-dashboard {
  padding: 20px;
}
.dashboard-header {
  margin-bottom: 24px;
}
.dashboard-header h1 {
  font-size: 24px;
  margin-bottom: 8px;
}
.dashboard-header p {
  color: #909399;
}
.stats-row {
  margin-bottom: 24px;
}
.stat-card {
  display: flex;
  align-items: center;
}
.stat-card :deep(.el-card__body) {
  display: flex;
  align-items: center;
  width: 100%;
  gap: 16px;
}
.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.stat-info {
  flex: 1;
}
.stat-value {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 4px;
}
.stat-label {
  font-size: 14px;
  color: #909399;
}
.content-card {
  margin-bottom: 24px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
