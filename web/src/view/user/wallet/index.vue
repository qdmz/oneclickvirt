<template>
  <div class="user-wallet">
    <div class="page-header">
      <div>
        <h1>{{ t('user.wallet.title') }}</h1>
        <p>{{ t('user.wallet.subtitle') }}</p>
      </div>
      <el-button
        :icon="Refresh"
        @click="reloadAll"
      >
        {{ t('common.refresh') }}
      </el-button>
    </div>

    <!-- 余额卡片 -->
    <el-card
      shadow="never"
      class="balance-card"
    >
      <div class="balance-inner">
        <div class="balance-left">
          <div class="balance-label">
            {{ t('user.wallet.currentBalance') }}
          </div>
          <div class="balance-value">
            <span class="cur">¥</span>{{ fenToYuan(balance) }}
          </div>
        </div>
        <div class="balance-right">
          <el-button
            type="primary"
            :icon="Wallet"
            size="large"
            @click="openRecharge"
          >
            {{ t('user.wallet.recharge') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 流水 -->
    <el-card
      shadow="never"
      class="table-card"
    >
      <template #header>
        <span class="card-title">{{ t('user.wallet.records') }}</span>
      </template>
      <el-table
        v-loading="loading"
        :data="list"
        style="width: 100%;"
      >
        <el-table-column
          prop="id"
          label="#"
          width="70"
        />
        <el-table-column
          prop="tradeNo"
          :label="t('user.wallet.colTradeNo')"
          min-width="220"
          show-overflow-tooltip
        />
        <el-table-column
          :label="t('user.wallet.colAmount')"
          width="140"
        >
          <template #default="{ row }">
            <span :class="row.amount >= 0 ? 'amt-in' : 'amt-out'">
              {{ row.amount >= 0 ? '+' : '-' }}¥{{ fenToYuan(Math.abs(row.amount)) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('user.wallet.colMethod')"
          width="120"
        >
          <template #default="{ row }">
            <el-tag
              size="small"
              effect="plain"
            >
              {{ t('user.wallet.method.' + (row.method || 'epay')) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('user.wallet.colStatus')"
          width="110"
        >
          <template #default="{ row }">
            <el-tag
              :type="statusTagType(row.status)"
              size="small"
            >
              {{ t('user.wallet.status.' + row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="remark"
          :label="t('user.wallet.colRemark')"
          min-width="180"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ row.remark || '-' }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('user.wallet.colTime')"
          width="170"
        >
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('common.actions')"
          width="110"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending' && row.method === 'epay'"
              link
              type="primary"
              @click="repay(row)"
            >
              {{ t('user.wallet.continuePay') }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>

    <!-- 充值对话框 -->
    <el-dialog
      v-model="rechargeVisible"
      :title="t('user.wallet.rechargeTitle')"
      width="440px"
    >
      <el-form
        ref="rechargeFormRef"
        :model="rechargeForm"
        :rules="rechargeRules"
        label-position="top"
      >
        <el-form-item
          :label="t('user.wallet.amountLabel')"
          prop="amount"
        >
          <el-input-number
            v-model="rechargeForm.amount"
            :min="1"
            :max="100000"
            :precision="2"
            :step="10"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item>
          <div class="quick-amounts">
            <el-button
              v-for="a in quickAmounts"
              :key="a"
              size="small"
              :type="rechargeForm.amount === a ? 'primary' : 'default'"
              @click="rechargeForm.amount = a"
            >
              ¥{{ a }}
            </el-button>
          </div>
        </el-form-item>
        <el-alert
          :title="t('user.wallet.rechargeTip')"
          type="info"
          :closable="false"
          show-icon
        />
      </el-form>
      <template #footer>
        <el-button @click="rechargeVisible = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="submitRecharge"
        >
          {{ t('user.wallet.goPay') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh, Wallet } from '@element-plus/icons-vue'
import { getWallet, createRecharge, listRechargeOrders } from '@/api/user'
import { submitToEpay, currentBaseURL, fenToYuan } from '@/utils/epay'

const { t } = useI18n()

const balance = ref(0)
const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const quickAmounts = [10, 30, 50, 100, 200, 500]

const statusTagType = (s) => ({ paid: 'success', pending: 'warning', failed: 'danger' }[s] || 'info')

const formatTime = (val) => {
  if (!val) return '-'
  const d = new Date(val)
  if (isNaN(d.getTime())) return '-'
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const loadBalance = async () => {
  try {
    const res = await getWallet()
    if (res && res.code === 200) {
      balance.value = res.data.balance || 0
    }
  } catch {
    // interceptor handles
  }
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await listRechargeOrders({ page: page.value, pageSize: pageSize.value })
    if (res && res.code === 200) {
      list.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch {
    // interceptor handles
  } finally {
    loading.value = false
  }
}

const reloadAll = () => {
  loadBalance()
  loadList()
}

// 充值
const rechargeVisible = ref(false)
const rechargeFormRef = ref()
const submitting = ref(false)
const rechargeForm = reactive({ amount: 50 })
const rechargeRules = {
  amount: [{ required: true, message: t('user.wallet.amountRequired'), trigger: 'blur' }]
}

const openRecharge = () => {
  rechargeForm.amount = 50
  rechargeVisible.value = true
}

const submitRecharge = async () => {
  if (!rechargeFormRef.value) return
  await rechargeFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const res = await createRecharge({
        amount: rechargeForm.amount,
        baseURL: currentBaseURL()
      })
      if (res && res.code === 200) {
        rechargeVisible.value = false
        const { epay, epayUrl } = res.data || {}
        if (epay && epayUrl) {
          ElMessage.success(t('user.wallet.redirecting'))
          submitToEpay(epayUrl, epay)
        } else {
          ElMessage.warning(t('user.wallet.noGateway'))
        }
        loadList()
      }
    } finally {
      submitting.value = false
    }
  })
}

// 待支付订单继续支付：重新发起同额充值
const repay = async (row) => {
  try {
    const res = await createRecharge({
      amount: Number(row.amount) / 100,
      baseURL: currentBaseURL()
    })
    if (res && res.code === 200) {
      const { epay, epayUrl } = res.data || {}
      if (epay && epayUrl) {
        submitToEpay(epayUrl, epay)
      } else {
        ElMessage.warning(t('user.wallet.noGateway'))
      }
    }
  } catch {
    // interceptor handles
  }
}

onMounted(() => {
  reloadAll()
  // 支付返回页带 ?pay=success 时提示并刷新
  const params = new URLSearchParams(window.location.search)
  if (params.get('pay') === 'success') {
    ElMessage.success(t('user.wallet.paySuccess'))
  }
})
</script>

<style scoped>
.user-wallet {
  padding: 20px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}
.page-header h1 {
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 4px;
}
.page-header p {
  color: var(--el-text-color-secondary);
  margin: 0;
  font-size: 13px;
}
.balance-card {
  border: 1px solid var(--el-border-color-lighter);
  margin-bottom: 16px;
}
.balance-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 4px;
}
.balance-label {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin-bottom: 6px;
}
.balance-value {
  font-size: 34px;
  font-weight: 600;
  line-height: 1;
  color: var(--el-color-primary);
}
.balance-value .cur {
  font-size: 20px;
  margin-right: 2px;
}
.table-card {
  border: 1px solid var(--el-border-color-lighter);
}
.card-title {
  font-weight: 600;
  font-size: 15px;
}
.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
.amt-in {
  color: var(--el-color-success);
  font-weight: 600;
}
.amt-out {
  color: var(--el-color-danger);
  font-weight: 600;
}
.quick-amounts {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
</style>
