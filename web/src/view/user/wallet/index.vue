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
          <div class="balance-actions">
            <el-button size="large" @click="openVoucherDialog">
              <el-icon><Ticket /></el-icon>
              {{ t('user.wallet.redeemVoucher') }}
            </el-button>
            <el-button type="primary" size="large" @click="rechargeVisible = true">
              <el-icon><Plus /></el-icon>
              {{ t('user.wallet.recharge') }}
            </el-button>
          </div>
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
          <el-table-column :label="t('user.wallet.logDesc')" prop="remark" min-width="200" />
          <el-table-column :label="t('user.wallet.logTime')" width="180">
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
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
        <div class="pay-type-select">
          <span>{{ t('user.wallet.payType') }}:</span>
          <el-radio-group v-model="selectedPayType">
            <el-radio v-if="enabledPayTypes.includes('alipay')" label="alipay">支付宝</el-radio>
            <el-radio v-if="enabledPayTypes.includes('wxpay')" label="wxpay">微信支付</el-radio>
            <el-radio v-if="enabledPayTypes.includes('qqpay')" label="qqpay">QQ钱包</el-radio>
          </el-radio-group>
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

    <!-- 代金券兑换弹窗 -->
    <el-dialog
      v-model="voucherVisible"
      :title="t('user.wallet.redeemVoucher')"
      width="420px"
      :close-on-click-modal="false"
    >
      <div class="voucher-content">
        <div class="voucher-hint">{{ t('user.wallet.voucherHint') }}</div>
        <el-input
          v-model="voucherCode"
          :placeholder="t('user.wallet.voucherPlaceholder')"
          size="large"
          clearable
          maxlength="64"
          @keyup.enter="handleRedeemVoucher"
        />
      </div>
      <template #footer>
        <el-button @click="voucherVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="voucherLoading" @click="handleRedeemVoucher">
          {{ t('user.wallet.confirmRedeem') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus, Ticket } from '@element-plus/icons-vue'
import { getUserBalance, getBalanceLogs, createYiPayOrder, redeemVoucher } from '@/api/product'
import { useSiteStore } from '@/pinia/modules/site'

const { t, locale } = useI18n()
const siteStore = useSiteStore()

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
const selectedPayType = ref('alipay')
const voucherVisible = ref(false)
const voucherLoading = ref(false)
const voucherCode = ref('')
// 启用的支付方式来自全局站点配置（后台“启用的支付方式”），与下单页保持一致
const enabledPayTypes = computed(() => siteStore.enabledPayTypes)

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

// 代金券兑换
const openVoucherDialog = () => {
  voucherCode.value = ''
  voucherVisible.value = true
}

const handleRedeemVoucher = async () => {
  const code = (voucherCode.value || '').trim().toUpperCase()
  if (!code) {
    ElMessage.warning(t('user.wallet.voucherPlaceholder'))
    return
  }
  voucherLoading.value = true
  try {
    const res = await redeemVoucher(code)
    if (res?.code === 200) {
      ElMessage.success(
        t('user.wallet.redeemSuccess', { amount: formatAmount(res.data?.amount) })
      )
      voucherVisible.value = false
      // 兑换成功后刷新余额与流水
      await Promise.all([loadBalance(), loadLogs()])
    } else {
      throw new Error(res?.message || t('user.wallet.redeemFailed'))
    }
  } catch (error) {
    ElMessage.error(error?.message || t('user.wallet.redeemFailed'))
  } finally {
    voucherLoading.value = false
  }
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
      payType: selectedPayType.value
    })
    if (res.code === 200 && res.data?.payURL) {
      window.open(res.data.payURL, '_blank')
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

// 若后台关闭了当前选中的支付方式，自动回退到第一个启用的渠道
watch(() => siteStore.enabledPayTypes, (list) => {
  if (list.length > 0 && !list.includes(selectedPayType.value)) {
    selectedPayType.value = list[0]
  }
}, { immediate: true })

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
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-color-dark) 100%);
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
.balance-actions {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.recharge-content {
  .recharge-hint {
    margin-bottom: 16px;
    color: var(--text-color-secondary);
  }
}

.voucher-content {
  .voucher-hint {
    margin-bottom: 16px;
    color: var(--text-color-secondary);
    line-height: 1.6;
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

.pay-type-select {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;

  > span {
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
