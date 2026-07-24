<template>
  <div class="wallet-container">
    <!-- 页面头部 -->
    <div class="wallet-header">
      <h1>{{ t('user.wallet.title') }}</h1>
    </div>

    <!-- 余额卡片 -->
    <div class="balance-section">
      <el-card class="balance-card">
        <div class="balance-content">
          <div class="balance-info">
            <div class="balance-label">{{ t('user.wallet.availableBalance') }}</div>
            <div class="balance-value">
              <span class="currency">¥</span>
              <span class="amount">{{ formatAmount(balance) }}</span>
            </div>
          </div>
          <el-button type="primary" size="large" @click="rechargeVisible = true">
            <el-icon><Plus /></el-icon>
            {{ t('user.wallet.recharge') }}
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 余额变动记录 -->
    <el-card class="logs-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('user.wallet.balanceLogs') }}</span>
          <div class="filter-group">
            <el-radio-group v-model="logType" size="small" @change="handleLogFilterChange">
              <el-radio-button label="">{{ t('user.wallet.all') }}</el-radio-button>
              <el-radio-button label="income">{{ t('user.wallet.income') }}</el-radio-button>
              <el-radio-button label="expense">{{ t('user.wallet.expense') }}</el-radio-button>
            </el-radio-group>
          </div>
        </div>
      </template>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <el-loading-directive />
        <div class="loading-text">{{ t('common.loading') }}</div>
      </div>

      <template v-else>
        <el-empty v-if="logList.length === 0" :description="t('user.wallet.noLogs')" />

        <el-table v-else :data="logList" stripe style="width: 100%">
          <el-table-column :label="t('user.wallet.logType')" width="120">
            <template #default="{ row }">
              <el-tag :type="row.amount > 0 ? 'success' : 'danger'" size="small">
                {{ row.amount > 0 ? t('user.wallet.income') : t('user.wallet.expense') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('user.wallet.logAmount')" width="150">
            <template #default="{ row }">
              <span :class="row.amount > 0 ? 'amount-income' : 'amount-expense'">
                {{ row.amount > 0 ? '+' : '' }}{{ formatAmount(row.amount) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column :label="t('user.wallet.logDesc')" prop="description" min-width="200" />
          <el-table-column :label="t('user.wallet.logTime')" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div v-if="total > pageSize" class="pagination-wrapper">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :total="total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
          />
        </div>
      </template>
    </el-card>

    <!-- 充值弹窗 -->
    <el-dialog
      v-model="rechargeVisible"
      :title="t('user.wallet.recharge')"
      width="400px"
      :close-on-click-modal="false"
    >
      <div class="recharge-content">
        <div class="recharge-hint">{{ t('user.wallet.selectAmount') }}</div>
        <div class="amount-grid">
          <el-button
            v-for="amt in presetAmounts"
            :key="amt"
            :type="rechargeAmount === amt ? 'primary' : 'default'"
            class="amount-btn"
            @click="rechargeAmount = amt"
          >
            ¥{{ amt }}
          </el-button>
        </div>
        <div class="custom-amount">
          <span>{{ t('user.wallet.customAmount') }}:</span>
          <el-input-number v-model="rechargeAmount" :min="1" :max="10000" :step="10" controls-position="right" />
        </div>
        <div class="recharge-summary">
          <span>{{ t('user.wallet.rechargeAmount') }}:</span>
          <span class="amount-highlight">¥{{ rechargeAmount }}</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="rechargeVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="rechargeLoading" @click="handleRecharge">
          {{ t('user.wallet.confirmRecharge') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getUserBalance, getBalanceLogs, createYiPayOrder } from '@/api/product'

const { t, locale } = useI18n()

const loading = ref(true)
const balance = ref(0)
const logList = ref([])
const logType = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const rechargeVisible = ref(false)
const rechargeAmount = ref(50)
const rechargeLoading = ref(false)

const presetAmounts = [10, 30, 50, 100, 200, 500]

const formatAmount = (amount) => {
  return Number(amount || 0).toFixed(2)
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString(locale.value === 'en-US' ? 'en-US' : 'zh-CN')
}

// 加载余额
const loadBalance = async () => {
  try {
    const res = await getUserBalance()
    if (res.code === 200) {
      balance.value = res.data?.balance || 0
    }
  } catch (error) {
    console.error('加载余额失败:', error)
  }
}

// 加载余额记录
const loadLogs = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      type: logType.value || undefined
    }
    const res = await getBalanceLogs(params)
    if (res.code === 200) {
      logList.value = res.data?.list || res.data?.items || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载余额记录失败:', error)
    ElMessage.error(error?.message || t('user.wallet.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleLogFilterChange = () => {
  currentPage.value = 1
  loadLogs()
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadLogs()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadLogs()
}

// 充值
const handleRecharge = async () => {
  if (!rechargeAmount.value || rechargeAmount.value < 1) {
    ElMessage.warning(t('user.wallet.amountTooSmall'))
    return
  }
  rechargeLoading.value = true
  try {
    const res = await createYiPayOrder({
      amount: rechargeAmount.value,
      type: 'recharge'
    })
    if (res.code === 200 && res.data?.pay_url) {
      window.open(res.data.pay_url, '_blank')
      ElMessage.success(t('user.wallet.openPayPage'))
      rechargeVisible.value = false
    } else {
      throw new Error(res.message || t('user.wallet.createPayFailed'))
    }
  } catch (error) {
    ElMessage.error(error?.message || t('user.wallet.rechargeFailed'))
  } finally {
    rechargeLoading.value = false
  }
}

onMounted(() => {
  loadBalance()
  loadLogs()
})
</script>

<style lang="scss" scoped>
.wallet-container {
  padding: 24px;
}

.wallet-header {
  margin-bottom: 20px;

  h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
    color: var(--text-color-primary);
  }
}

.balance-section {
  margin-bottom: 24px;
}

.balance-card {
  background: linear-gradient(135deg, #16a34a 0%, #15803d 100%);
  border: none;

  :deep(.el-card__body) {
    padding: 32px;
  }
}

.balance-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.balance-info {
  color: #ffffff;
}

.balance-label {
  font-size: 14px;
  margin-bottom: 8px;
  opacity: 0.9;
}

.balance-value {
  display: flex;
  align-items: baseline;
  gap: 4px;

  .currency {
    font-size: 24px;
    font-weight: 600;
  }

  .amount {
    font-size: 40px;
    font-weight: 700;
  }
}

.logs-card {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 200px;
}

.amount-income {
  color: #67c23a;
  font-weight: 600;
}

.amount-expense {
  color: #f56c6c;
  font-weight: 600;
}

.pagination-wrapper {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}

/* 充值弹窗 */
.recharge-content {
  .recharge-hint {
    margin-bottom: 16px;
    color: var(--text-color-secondary);
  }
}

.amount-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.amount-btn {
  width: 100%;
}

.custom-amount {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;

  span {
    color: var(--text-color-secondary);
    white-space: nowrap;
  }
}

.recharge-summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: var(--neutral-bg);
  border-radius: 8px;

  .amount-highlight {
    font-size: 20px;
    font-weight: 700;
    color: #f56c6c;
  }
}

@media (max-width: 768px) {
  .wallet-container {
    padding: 16px;
  }

  .balance-content {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }

  .balance-value .amount {
    font-size: 32px;
  }

  .amount-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
