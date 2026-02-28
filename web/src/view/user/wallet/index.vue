<template>
  <div class="wallet-container">
    <!-- 钱包卡片 -->
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="wallet-card">
          <template #header>
            <div class="card-header">
              <span>钱包余额</span>
              <el-icon><Wallet /></el-icon>
            </div>
          </template>
          <div class="balance">
            <span class="amount">¥{{ (walletInfo.balance / 100).toFixed(2) }}</span>
          </div>
          <div class="balance-info">
            <el-descriptions
              :column="1"
              size="small"
            >
              <el-descriptions-item label="累计充值">
                ¥{{ (walletInfo.totalRecharge / 100).toFixed(2) }}
              </el-descriptions-item>
              <el-descriptions-item label="累计消费">
                ¥{{ (walletInfo.totalExpense / 100).toFixed(2) }}
              </el-descriptions-item>
              <el-descriptions-item label="冻结金额">
                ¥{{ (walletInfo.frozen / 100).toFixed(2) }}
              </el-descriptions-item>
            </el-descriptions>
          </div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>充值</span>
              <el-tag type="info">
                快捷充值
              </el-tag>
            </div>
          </template>
          <el-form
            :model="rechargeForm"
            label-width="100px"
          >
            <el-form-item label="充值金额">
              <el-radio-group v-model="rechargeForm.amount">
                <el-radio-button :label="10">
                  ¥10
                </el-radio-button>
                <el-radio-button :label="50">
                  ¥50
                </el-radio-button>
                <el-radio-button :label="100">
                  ¥100
                </el-radio-button>
                <el-radio-button :label="200">
                  ¥200
                </el-radio-button>
                <el-radio-button :label="500">
                  ¥500
                </el-radio-button>
              </el-radio-group>
              <el-input-number
                v-model="rechargeForm.amount"
                :min="1"
                :max="10000"
                :precision="2"
                class="custom-amount"
              />
            </el-form-item>
            <el-form-item label="支付方式">
              <el-radio-group v-model="rechargeForm.paymentMethod">
                <el-radio
                  v-if="paymentConfig.enableAlipay"
                  value="alipay"
                >
                  <el-icon color="#1677ff">
                    <Wallet />
                  </el-icon>
                  支付宝
                </el-radio>
                <el-radio
                  v-if="paymentConfig.enableWechat"
                  value="wechat"
                >
                  <el-icon color="#07c160">
                    <ChatDotRound />
                  </el-icon>
                  微信支付
                </el-radio>
                <el-radio
                  v-if="paymentConfig.enableMapay"
                  value="mapay"
                >
                  <el-icon color="#f7ba2a">
                    <Wallet />
                  </el-icon>
                  码支付
                </el-radio>
                <el-radio
                  v-if="paymentConfig.enableEpay"
                  value="epay"
                >
                  <el-icon color="#409eff">
                    <Wallet />
                  </el-icon>
                  易支付
                </el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                :loading="recharging"
                @click="handleRecharge"
              >
                立即充值
              </el-button>
              <el-button @click="showExchangeDialog = true">
                使用兑换码
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>

    <!-- 交易记录 -->
    <el-card class="mt-20">
      <template #header>
        <div class="card-header">
          <span>交易记录</span>
          <el-tag>{{ transactions.total }} 条记录</el-tag>
        </div>
      </template>
      <el-table
        v-loading="loading"
        :data="transactions.list"
        border
        stripe
      >
        <el-table-column
          prop="id"
          label="ID"
          width="80"
        />
        <el-table-column
          label="交易类型"
          width="120"
        >
          <template #default="{ row }">
            <el-tag :type="getTransactionTypeTag(row.type)">
              {{ getTransactionTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          label="金额"
          width="150"
        >
          <template #default="{ row }">
            <span :class="row.amount >= 0 ? 'income' : 'expense'">
              {{ row.amount >= 0 ? '+' : '' }}¥{{ (row.amount / 100).toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          prop="balance"
          label="余额"
          width="150"
        >
          <template #default="{ row }">
            ¥{{ (row.balance / 100).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="description"
          label="说明"
          show-overflow-tooltip
        />
        <el-table-column
          label="交易时间"
          width="180"
        >
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="transactions.page"
        v-model:page-size="transactions.pageSize"
        :total="transactions.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        class="mt-20"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>

    <!-- 兑换码对话框 -->
    <el-dialog
      v-model="showExchangeDialog"
      title="兑换码充值"
      width="500px"
    >
      <el-form
        :model="exchangeForm"
        label-width="100px"
      >
        <el-form-item label="兑换码">
          <el-input
            v-model="exchangeForm.code"
            placeholder="请输入兑换码"
            clearable
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showExchangeDialog = false">
          取消
        </el-button>
        <el-button
          type="primary"
          :loading="exchanging"
          @click="handleExchange"
        >
          兑换
        </el-button>
      </template>
    </el-dialog>

    <!-- 支付二维码对话框 -->
    <el-dialog
      v-model="showQRDialog"
      title="扫码支付"
      width="400px"
      @close="handleCloseQRDialog"
    >
      <div class="qr-container">
        <el-image
          v-if="qrCodeUrl"
          :src="qrCodeUrl"
          fit="contain"
          class="qr-code"
        />
        <p class="amount-text">
          支付金额: ¥{{ (currentOrder.amount / 100).toFixed(2) }}
        </p>
        <p class="tip-text">
          请使用{{ getPaymentMethodName(rechargeForm.paymentMethod) }}扫描二维码
        </p>
        <el-alert
          title="支付完成后页面将自动跳转"
          type="info"
          :closable="false"
          class="mt-20"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Wallet, ChatDotRound } from '@element-plus/icons-vue'
import {
  getWallet,
  getWalletTransactions,
  createRechargeOrder,
  getAlipayQR,
  getWechatQR,
  exchangeRedemptionCode,
  getRechargeOrderStatus
} from '@/api/user-payment'
import { getPaymentConfig } from '@/api/public'

const loading = ref(false)
const recharging = ref(false)
const exchanging = ref(false)
const showExchangeDialog = ref(false)
const showQRDialog = ref(false)
const qrCodeUrl = ref('')
const pollTimer = ref(null)

const walletInfo = ref({
  balance: 0,
  frozen: 0,
  totalRecharge: 0,
  totalExpense: 0
})

// 支付配置
const paymentConfig = ref({
  enableAlipay: true,
  enableWechat: true,
  enableBalance: true,
  enableEpay: false,
  enableMapay: false
})

const rechargeForm = ref({
  amount: 10,
  paymentMethod: 'alipay'
})

const exchangeForm = ref({
  code: ''
})

const currentOrder = ref({
  orderNo: '',
  amount: 0,
  status: ''
})

const transactions = ref({
  list: [],
  total: 0,
  page: 1,
  pageSize: 20
})

// 加载钱包信息
const loadWallet = async () => {
  try {
    const res = await getWallet()
    if (res.code === 200) {
      walletInfo.value = res.data || {}
    }
  } catch (error) {
    ElMessage.error('加载钱包信息失败')
  }
}

// 加载交易记录
const loadTransactions = async () => {
  loading.value = true
  try {
    const res = await getWalletTransactions({
      page: transactions.value.page,
      pageSize: transactions.value.pageSize
    })
    if (res.code === 200) {
      transactions.value.list = res.data.list || []
      transactions.value.total = res.data.total || 0
    }
  } catch (error) {
    ElMessage.error('加载交易记录失败')
  } finally {
    loading.value = false
  }
}

// 获取交易类型名称
const getTransactionTypeName = (type) => {
  const map = {
    recharge: '充值',
    consume: '消费',
    refund: '退款',
    withdraw: '提现',
    exchange: '兑换',
    system: '系统'
  }
  return map[type] || type
}

// 获取交易类型标签
const getTransactionTypeTag = (type) => {
  const map = {
    recharge: 'success',
    consume: 'warning',
    refund: 'info',
    withdraw: 'danger',
    exchange: 'success',
    system: 'info'
  }
  return map[type] || ''
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

// 获取支付方式名称
const getPaymentMethodName = (method) => {
  const map = {
    alipay: '支付宝',
    wechat: '微信',
    mapay: '码支付',
    epay: '易支付'
  }
  return map[method] || method
}

// 充值
const handleRecharge = async () => {
  if (!rechargeForm.value.amount || rechargeForm.value.amount <= 0) {
    ElMessage.warning('请输入充值金额')
    return
  }

  recharging.value = true
  try {
    const res = await createRechargeOrder({
      amount: Math.round(rechargeForm.value.amount * 100),
      paymentMethod: rechargeForm.value.paymentMethod
    })

    if (res.code === 200) {
      currentOrder.value = res.data
      // 获取支付二维码
      await getQRCode()
      showQRDialog.value = true
      // 开始轮询订单状态
      startPollOrderStatus()
    } else {
      ElMessage.error(res.message || '创建订单失败')
    }
  } catch (error) {
    ElMessage.error('创建订单失败')
  } finally {
    recharging.value = false
  }
}

// 获取支付二维码
const getQRCode = async () => {
  try {
    let res
    if (rechargeForm.value.paymentMethod === 'alipay') {
      res = await getAlipayQR(currentOrder.value.orderNo)
    } else {
      res = await getWechatQR(currentOrder.value.orderNo)
    }

    if (res.code === 200) {
      qrCodeUrl.value = res.data.qrCode
    } else {
      ElMessage.error('获取支付二维码失败')
    }
  } catch (error) {
    ElMessage.error('获取支付二维码失败')
  }
}

// 开始轮询订单状态
const startPollOrderStatus = () => {
  if (pollTimer.value) {
    clearInterval(pollTimer.value)
  }

  pollTimer.value = setInterval(async () => {
    try {
      const res = await getRechargeOrderStatus(currentOrder.value.orderNo)
      if (res.code === 200) {
        const status = res.data.status
        if (status === 'paid') {
          // 支付成功
          clearInterval(pollTimer.value)
          pollTimer.value = null
          ElMessage.success('充值成功')
          showQRDialog.value = false
          qrCodeUrl.value = ''
          loadWallet()
          loadTransactions()
        } else if (status === 'cancelled' || status === 'expired') {
          // 订单已取消或过期
          clearInterval(pollTimer.value)
          pollTimer.value = null
          ElMessage.warning('订单已' + (status === 'cancelled' ? '取消' : '过期'))
          showQRDialog.value = false
          qrCodeUrl.value = ''
        }
      }
    } catch (error) {
      console.error('查询订单状态失败:', error)
    }
  }, 3000)
}

// 关闭二维码对话框
const handleCloseQRDialog = () => {
  if (pollTimer.value) {
    clearInterval(pollTimer.value)
    pollTimer.value = null
  }
  qrCodeUrl.value = ''
  // 刷新钱包和交易记录
  loadWallet()
  loadTransactions()
}

// 兑换码充值
const handleExchange = async () => {
  if (!exchangeForm.value.code) {
    ElMessage.warning('请输入兑换码')
    return
  }

  exchanging.value = true
  try {
    const res = await exchangeRedemptionCode({ code: exchangeForm.value.code })
    console.log('兑换码API响应:', res)
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('兑换成功')
      showExchangeDialog.value = false
      exchangeForm.value.code = ''
      loadWallet()
      loadTransactions()
    } else {
      ElMessage.error(res.message || '兑换失败')
    }
  } catch (error) {
    console.error('兑换码兑换失败:', error)
    ElMessage.error('兑换失败')
  } finally {
    exchanging.value = false
  }
}

// 分页大小变化
const handleSizeChange = (val) => {
  transactions.value.pageSize = val
  transactions.value.page = 1
  loadTransactions()
}

// 页码变化
const handleCurrentChange = (val) => {
  transactions.value.page = val
  loadTransactions()
}

// 加载支付配置
const loadPaymentConfig = async () => {
  try {
    const res = await getPaymentConfig()
    if (res.code === 200 || res.code === 0) {
      paymentConfig.value = { ...paymentConfig.value, ...res.data }
      // 设置默认支付方式
      if (paymentConfig.value.enableAlipay) {
        rechargeForm.value.paymentMethod = 'alipay'
      } else if (paymentConfig.value.enableWechat) {
        rechargeForm.value.paymentMethod = 'wechat'
      } else if (paymentConfig.value.enableEpay) {
        rechargeForm.value.paymentMethod = 'epay'
      } else if (paymentConfig.value.enableMapay) {
        rechargeForm.value.paymentMethod = 'mapay'
      }
    }
  } catch (error) {
    console.error('加载支付配置失败:', error)
  }
}

onMounted(() => {
  loadPaymentConfig()
  loadWallet()
  loadTransactions()
})
</script>

<style scoped>
.wallet-container {
  padding: 20px;
}

.wallet-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.wallet-card :deep(.el-card__header) {
  background: transparent;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
}

.wallet-card :deep(.el-card__body) {
  color: white;
}

.wallet-card .balance {
  font-size: 48px;
  font-weight: bold;
  margin: 20px 0;
}

.wallet-card :deep(.el-descriptions__label) {
  color: rgba(255, 255, 255, 0.8);
}

.wallet-card :deep(.el-descriptions__content) {
  color: white;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.balance-info {
  margin-top: 20px;
}

.mt-20 {
  margin-top: 20px;
}

.custom-amount {
  margin-left: 10px;
}

.income {
  color: #67c23a;
  font-weight: bold;
}

.expense {
  color: #f56c6c;
  font-weight: bold;
}

.qr-container {
  text-align: center;
}

.qr-code {
  width: 280px;
  height: 280px;
  margin: 0 auto;
}

.amount-text {
  font-size: 24px;
  font-weight: bold;
  margin: 20px 0 10px;
}

.tip-text {
  color: #909399;
}
</style>
