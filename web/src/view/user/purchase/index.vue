<template>
  <div class="purchase-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>产品购买</span>
          <el-tag type="info">
            选择适合您的套餐
          </el-tag>
        </div>
      </template>

      <el-row
        v-loading="loading"
        :gutter="20"
      >
        <el-col
          v-for="product in filteredProducts"
          :key="product.id"
          :span="8"
          class="mb-20"
        >
          <el-card
            :body-style="{ padding: '20px' }"
            class="product-card"
          >
            <template #header>
              <div class="product-header">
                <span class="product-name">{{ product.name }}</span>
                <el-tag
                  v-if="product.period === 0"
                  type="success"
                >
                  永久
                </el-tag>
                <el-tag
                  v-else
                  type="info"
                >
                  {{ product.period }}天
                </el-tag>
              </div>
            </template>
            <div class="product-body">
              <div class="price">
                <span class="amount">¥{{ (product.price / 100).toFixed(2) }}</span>
                <span class="unit">{{ product.period === 0 ? '/永久' : `/${product.period}天` }}</span>
              </div>
              <div class="level-badge">
                Lv.{{ product.level }}
              </div>
              <div class="features">
                <el-descriptions
                  :column="1"
                  size="small"
                >
                  <el-descriptions-item label="CPU">
                    {{ product.cpu }} 核
                  </el-descriptions-item>
                  <el-descriptions-item label="内存">
                    {{ formatMemory(product.memory) }}
                  </el-descriptions-item>
                  <el-descriptions-item label="磁盘">
                    {{ formatDisk(product.disk) }}
                  </el-descriptions-item>
                  <el-descriptions-item label="带宽">
                    {{ product.bandwidth }} Mbps
                  </el-descriptions-item>
                  <el-descriptions-item label="流量">
                    {{ formatTraffic(product.traffic) }}
                  </el-descriptions-item>
                  <el-descriptions-item label="最大实例">
                    {{ product.maxInstances }} 个
                  </el-descriptions-item>
                </el-descriptions>
              </div>
              <el-button
                type="primary"
                :disabled="!product.isEnabled"
                class="purchase-btn"
                @click="handlePurchase(product)"
              >
                {{ product.isEnabled ? '立即购买' : '已下架' }}
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>

    <!-- 支付方式选择对话框 -->
    <el-dialog
      v-model="paymentDialogVisible"
      title="选择支付方式"
      width="400px"
    >
      <div class="payment-info">
        <p class="product-title">
          {{ selectedProduct?.name }}
        </p>
        <p class="order-amount">
          订单金额: ¥{{ (selectedProduct?.price / 100).toFixed(2) }}
        </p>
      </div>
      <el-radio-group
        v-model="paymentMethod"
        class="payment-methods"
      >
        <el-radio
          v-if="paymentConfig.enableAlipay"
          value="alipay"
          class="payment-option"
        >
          <el-icon color="#1677ff">
            <Wallet />
          </el-icon>
          <span>支付宝支付</span>
        </el-radio>
        <el-radio
          v-if="paymentConfig.enableWechat"
          value="wechat"
          class="payment-option"
        >
          <el-icon color="#07c160">
            <ChatDotRound />
          </el-icon>
          <span>微信支付</span>
        </el-radio>
        <el-radio
          v-if="paymentConfig.enableMapay"
          value="mapay"
          class="payment-option"
        >
          <el-icon color="#f7ba2a">
            <Wallet />
          </el-icon>
          <span>码支付</span>
        </el-radio>
        <el-radio
          v-if="paymentConfig.enableEpay"
          value="epay"
          class="payment-option"
        >
          <el-icon color="#409eff">
            <Wallet />
          </el-icon>
          <span>易支付</span>
        </el-radio>
        <el-radio
          v-if="paymentConfig.enableBalance"
          value="balance"
          class="payment-option"
        >
          <el-icon color="#67c23a">
            <Coin />
          </el-icon>
          <span>余额支付</span>
        </el-radio>
      </el-radio-group>
      <template #footer>
        <el-button @click="paymentDialogVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          :loading="purchasing"
          @click="handleConfirmPurchase"
        >
          确认购买
        </el-button>
      </template>
    </el-dialog>

    <!-- 支付二维码对话框 -->
    <el-dialog
      v-model="qrDialogVisible"
      title="扫码支付"
      width="400px"
      @close="handleCloseQRDialog"
    >
      <div class="qr-container">
        <div
          v-if="qrCodeUrl && (paymentMethod === 'epay' || paymentMethod === 'mapay')"
          class="payment-link-container"
        >
          <el-alert
            type="info"
            :closable="false"
            class="mb-20"
          >
            <template #title>
              请点击下方按钮跳转支付页面
            </template>
          </el-alert>
          <el-button
            type="primary"
            size="large"
            @click="openPaymentUrl"
          >
            立即支付
          </el-button>
        </div>
        <el-image
          v-else-if="qrCodeUrl"
          :src="qrCodeUrl"
          fit="contain"
          class="qr-code"
        />
        <p class="amount-text">
          支付金额: ¥{{ (currentOrder?.amount / 100).toFixed(2) }}
        </p>
        <p class="tip-text">
          请使用{{ getPaymentMethodName(paymentMethod) }}{{ paymentMethod === 'epay' || paymentMethod === 'mapay' ? '完成支付' : '扫描二维码' }}
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
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Wallet, ChatDotRound, Coin } from '@element-plus/icons-vue'
import { getUserProducts, purchaseProduct, getAlipayQR, getWechatQR, getPurchaseOrderStatus, getPurchaseEpayQR, getPurchaseMapayQR } from '@/api/user-payment'
import { getUserProfile } from '@/api/user'
import { getPaymentConfig } from '@/api/public'

const router = useRouter()

const loading = ref(false)
const purchasing = ref(false)
// const products = ref([])
const allProducts = ref([])
const userLevel = ref(0)
const paymentDialogVisible = ref(false)
const qrDialogVisible = ref(false)
const selectedProduct = ref(null)
const paymentMethod = ref('alipay')
const qrCodeUrl = ref('')
const currentOrder = ref(null)
const pollTimer = ref(null)

// 支付配置
const paymentConfig = ref({
  enableAlipay: true,
  enableWechat: true,
  enableBalance: true,
  enableEpay: false,
  enableMapay: false
})

// 过滤后的产品列表（只显示等于和高于用户等级的产品）
const filteredProducts = computed(() => {
  return allProducts.value.filter(product => product.level >= userLevel.value)
})

// 加载支付配置
const loadPaymentConfig = async () => {
  try {
    const res = await getPaymentConfig()
    if (res.code === 200 || res.code === 0) {
      paymentConfig.value = { ...paymentConfig.value, ...res.data }
      // 设置默认支付方式
      if (paymentConfig.value.enableAlipay) {
        paymentMethod.value = 'alipay'
      } else if (paymentConfig.value.enableWechat) {
        paymentMethod.value = 'wechat'
      } else if (paymentConfig.value.enableEpay) {
        paymentMethod.value = 'epay'
      } else if (paymentConfig.value.enableMapay) {
        paymentMethod.value = 'mapay'
      } else if (paymentConfig.value.enableBalance) {
        paymentMethod.value = 'balance'
      }
    }
  } catch (error) {
    console.error('加载支付配置失败:', error)
  }
}

// 加载用户信息
const loadUserInfo = async () => {
  try {
    const res = await getUserProfile()
    console.log('用户信息API响应:', res)
    if (res.code === 200 || res.code === 0) {
      // 检查返回数据结构
      console.log('用户信息数据结构:', res.data)
      
      // 尝试多种可能的数据路径，优先检查小写user
      if (res.data.user && res.data.user.level) {
        userLevel.value = res.data.user.level
        console.log('从 res.data.user.level 获取用户等级:', userLevel.value)
      } else if (res.data.User && res.data.User.level) {
        userLevel.value = res.data.User.level
        console.log('从 res.data.User.level 获取用户等级:', userLevel.value)
      } else if (res.data.level) {
        userLevel.value = res.data.level
        console.log('从 res.data.level 获取用户等级:', userLevel.value)
      } else {
        console.error('无法从API响应中获取用户等级，响应数据:', res.data)
        userLevel.value = 0
      }
    } else {
      console.error('获取用户信息失败，响应码:', res.code, '消息:', res.message)
      userLevel.value = 0
    }
  } catch (error) {
    console.error('获取用户信息异常:', error)
    // 获取用户信息失败不影响产品列表加载，使用默认等级0
    userLevel.value = 0
  }
  console.log('最终用户等级:', userLevel.value)
}

// 加载产品列表
const loadProducts = async () => {
  loading.value = true
  try {
    // 先获取用户信息
    await loadUserInfo()
    
    // 再获取产品列表
    const res = await getUserProducts()
    console.log('产品列表API响应:', res)
    if (res.code === 200 || res.code === 0) {
      // 将后端返回的整数isEnabled转换为布尔值
      allProducts.value = (res.data || []).map(product => ({
        ...product,
        isEnabled: product.isEnabled === 1
      }))
      
      // 显示所有产品的信息
      console.log('所有产品列表:', allProducts.value)
      console.log('所有产品数量:', allProducts.value.length)
      
      // 显示每个产品的等级
      allProducts.value.forEach(product => {
        console.log(`产品 ${product.name} (ID: ${product.id}) 等级: ${product.level}`)
      })
    } else {
      console.error('加载产品列表失败，响应码:', res.code, '消息:', res.message)
      ElMessage.error(res.message || '加载产品列表失败')
    }
  } catch (error) {
    console.error('加载产品列表异常:', error)
    console.error('错误详情:', error.message)
    console.error('错误栈:', error.stack)
    ElMessage.error('加载产品列表失败')
  } finally {
    loading.value = false
    // 显示过滤后的产品信息
    console.log('过滤条件 - 用户等级:', userLevel.value)
    console.log('过滤前产品数量:', allProducts.value.length)
    console.log('过滤后产品数量:', filteredProducts.value.length)
    console.log('过滤后的产品列表:', filteredProducts.value)
    
    // 显示过滤后的产品信息
    filteredProducts.value.forEach(product => {
      console.log(`过滤后产品 ${product.name} (ID: ${product.id}) 等级: ${product.level}`)
    })
  }
}

// 格式化内存
const formatMemory = (memory) => {
  if (memory < 1024) {
    return memory + 'MB'
  } else {
    return (memory / 1024).toFixed(2) + 'GB'
  }
}

// 格式化磁盘
const formatDisk = (disk) => {
  if (disk < 1024) {
    return disk + 'MB'
  } else {
    return (disk / 1024).toFixed(2) + 'GB'
  }
}

// 格式化流量
const formatTraffic = (traffic) => {
  if (traffic < 1024) {
    return traffic + 'MB'
  } else if (traffic < 1024 * 1024) {
    return (traffic / 1024).toFixed(2) + 'GB'
  } else {
    return (traffic / 1024 / 1024).toFixed(2) + 'TB'
  }
}

// 获取支付方式名称
const getPaymentMethodName = (method) => {
  const map = {
    alipay: '支付宝',
    wechat: '微信',
    mapay: '码支付',
    epay: '易支付',
    balance: '余额'
  }
  return map[method] || method
}

// 购买产品
const handlePurchase = (product) => {
  selectedProduct.value = product
  paymentMethod.value = 'alipay'
  paymentDialogVisible.value = true
}

// 确认购买
const handleConfirmPurchase = async () => {
  if (!selectedProduct.value) return

  purchasing.value = true
  try {
    const res = await purchaseProduct(selectedProduct.value.id, {
      paymentMethod: paymentMethod.value
    })

    if (res.code === 200) {
      currentOrder.value = res.data
      paymentDialogVisible.value = false

      // 如果是余额支付,直接完成
      if (paymentMethod.value === 'balance') {
        ElMessage.success('购买成功')
        loadProducts()
        // 跳转到实例创建页面
        router.push('/user/apply')
      } else {
        // 获取支付二维码
        await getQRCode()
        qrDialogVisible.value = true
        // 开始轮询订单状态
        startPollOrderStatus()
      }
    } else {
      ElMessage.error(res.message || '创建订单失败')
    }
  } catch (error) {
    ElMessage.error('创建订单失败')
  } finally {
    purchasing.value = false
  }
}

// 获取支付二维码
const getQRCode = async () => {
  try {
    let res
    if (paymentMethod.value === 'alipay') {
      res = await getAlipayQR(currentOrder.value.orderNo)
    } else if (paymentMethod.value === 'wechat') {
      res = await getWechatQR(currentOrder.value.orderNo)
    } else if (paymentMethod.value === 'epay') {
      res = await getPurchaseEpayQR(currentOrder.value.orderNo)
    } else if (paymentMethod.value === 'mapay') {
      res = await getPurchaseMapayQR(currentOrder.value.orderNo)
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
      const res = await getPurchaseOrderStatus(currentOrder.value.orderNo)
      if (res.code === 200) {
        const status = res.data.status
        if (status === 'paid') {
          // 支付成功
          clearInterval(pollTimer.value)
          pollTimer.value = null
          ElMessage.success('购买成功')
          qrDialogVisible.value = false
          qrCodeUrl.value = ''
          loadProducts()
          // 跳转到实例创建页面
          router.push('/user/apply')
        } else if (status === 'cancelled' || status === 'expired') {
          // 订单已取消或过期
          clearInterval(pollTimer.value)
          pollTimer.value = null
          ElMessage.warning('订单已' + (status === 'cancelled' ? '取消' : '过期'))
          qrDialogVisible.value = false
          qrCodeUrl.value = ''
        }
      }
    } catch (error) {
      console.error('查询订单状态失败:', error)
    }
  }, 3000)
}

// 打开支付链接
const openPaymentUrl = () => {
  if (qrCodeUrl.value) {
    window.open(qrCodeUrl.value, '_blank')
  }
}

// 关闭二维码对话框
const handleCloseQRDialog = () => {
  if (pollTimer.value) {
    clearInterval(pollTimer.value)
    pollTimer.value = null
  }
  qrCodeUrl.value = ''
}

onMounted(() => {
  loadPaymentConfig()
  loadProducts()
})
</script>

<style scoped>
.purchase-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.mb-20 {
  margin-bottom: 20px;
}

.product-card {
  transition: all 0.3s;
  border: 2px solid transparent;
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  border-color: #409eff;
}

.product-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.product-name {
  font-size: 18px;
  font-weight: bold;
}

.product-body {
  text-align: center;
}

.price {
  margin: 20px 0;
}

.price .amount {
  font-size: 36px;
  font-weight: bold;
  color: #409eff;
}

.price .unit {
  font-size: 14px;
  color: #909399;
}

.level-badge {
  display: inline-block;
  padding: 5px 15px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 20px;
  margin-bottom: 15px;
  font-weight: bold;
}

.features {
  margin: 20px 0;
  text-align: left;
}

.purchase-btn {
  width: 100%;
  margin-top: 10px;
}

.payment-info {
  text-align: center;
  margin-bottom: 20px;
}

.product-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
}

.order-amount {
  font-size: 24px;
  color: #409eff;
  font-weight: bold;
}

.payment-methods {
  display: block;
  text-align: center;
}

.payment-option {
  display: block;
  margin: 15px 0;
  padding: 15px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  transition: all 0.3s;
}

.payment-option:hover {
  border-color: #409eff;
  background-color: #f0f9ff;
}

.payment-option :deep(.el-radio__label) {
  display: flex;
  align-items: center;
  font-size: 16px;
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
