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
