<template>
  <div class="product-detail-container">
    <!-- 返回按钮 -->
    <el-button class="back-btn" @click="goBack">
      <el-icon><ArrowLeft /></el-icon>
      {{ t('common.back') }}
    </el-button>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-loading-directive />
      <div class="loading-text">{{ t('common.loading') }}</div>
    </div>

    <template v-else-if="product.id">
      <div class="detail-layout">
        <!-- 左侧：产品信息 -->
        <div class="detail-left">
          <el-card>
            <div class="product-header">
              <div class="product-icon-large">
                <el-icon :size="64" color="#16a34a">
                  <component :is="getProductIcon(product.type)" />
                </el-icon>
              </div>
              <div class="product-title">
                <h1>{{ product.name }}</h1>
                <p>{{ product.description }}</p>
                <div class="product-tags">
                  <el-tag v-if="product.is_new" type="success" size="small">{{ t('user.store.newProduct') }}</el-tag>
                  <el-tag v-if="product.is_hot" type="danger" size="small">{{ t('user.store.hotProduct') }}</el-tag>
                </div>
              </div>
            </div>

            <!-- 资源配置详情 -->
            <div class="specs-detail">
              <h3>{{ t('user.store.resourceConfig') }}</h3>
              <div class="specs-grid">
                <div class="spec-card">
                  <el-icon><Cpu /></el-icon>
                  <div class="spec-label">{{ t('user.store.cpu') }}</div>
                  <div class="spec-value">{{ product.cpu }} {{ t('user.store.cores') }}</div>
                </div>
                <div class="spec-card">
                  <el-icon><Memo /></el-icon>
                  <div class="spec-label">{{ t('user.store.memory') }}</div>
                  <div class="spec-value">{{ formatMemory(product.memory) }}</div>
                </div>
                <div class="spec-card">
                  <el-icon><Coin /></el-icon>
                  <div class="spec-label">{{ t('user.store.disk') }}</div>
                  <div class="spec-value">{{ formatDisk(product.disk) }}</div>
                </div>
                <div class="spec-card">
                  <el-icon><TopRight /></el-icon>
                  <div class="spec-label">{{ t('user.store.bandwidth') }}</div>
                  <div class="spec-value">{{ formatBandwidth(product.bandwidth) }}</div>
                </div>
                <div class="spec-card">
                  <el-icon><DataLine /></el-icon>
                  <div class="spec-label">{{ t('user.store.traffic') }}</div>
                  <div class="spec-value">
                    {{ product.traffic > 0 ? formatTraffic(product.traffic) : t('user.store.unlimitedTraffic') }}
                  </div>
                </div>
                <div class="spec-card">
                  <el-icon><OfficeBuilding /></el-icon>
                  <div class="spec-label">{{ t('user.store.node') }}</div>
                  <div class="spec-value">{{ product.node_name || t('user.store.autoAssign') }}</div>
                </div>
              </div>
            </div>
          </el-card>
        </div>

        <!-- 右侧：购买配置 -->
        <div class="detail-right">
          <el-card>
            <template #header>
              <span>{{ t('user.store.purchaseConfig') }}</span>
            </template>

            <!-- 操作系统选择 -->
            <div class="config-section">
              <label class="config-label">{{ t('user.store.selectOS') }}</label>
              <el-select
                v-model="selectedImage"
                :placeholder="t('user.store.pleaseSelectOS')"
                class="config-select"
              >
                <el-option-group
                  v-for="group in imageGroups"
                  :key="group.label"
                  :label="group.label"
                >
                  <el-option
                    v-for="image in group.options"
                    :key="image.id"
                    :label="image.name"
                    :value="image.id"
                  />
                </el-option-group>
              </el-select>
            </div>

            <!-- 购买周期 -->
            <div class="config-section">
              <label class="config-label">{{ t('user.store.selectPeriod') }}</label>
              <el-radio-group v-model="selectedPeriod" class="period-group">
                <el-radio-button
                  v-for="period in periodOptions"
                  :key="period.value"
                  :label="period.value"
                >
                  {{ period.label }}
                </el-radio-button>
              </el-radio-group>
            </div>

            <!-- 数量 -->
            <div class="config-section">
              <label class="config-label">{{ t('user.store.quantity') }}</label>
              <el-input-number v-model="quantity" :min="1" :max="10" />
            </div>

            <el-divider />

            <!-- 价格汇总 -->
            <div class="price-summary">
              <div class="price-row">
                <span>{{ t('user.store.unitPrice') }}</span>
                <span>¥{{ product.price }}</span>
              </div>
              <div class="price-row">
                <span>{{ t('user.store.period') }}</span>
                <span>{{ getPeriodLabel() }}</span>
              </div>
              <div class="price-row">
                <span>{{ t('user.store.quantity') }}</span>
                <span>x{{ quantity }}</span>
              </div>
              <div class="price-row total">
                <span>{{ t('user.store.totalPrice') }}</span>
                <span class="total-price">¥{{ totalPrice }}</span>
              </div>
            </div>

            <!-- 支付方式 -->
            <div class="payment-methods">
              <label class="config-label">{{ t('user.store.paymentMethod') }}</label>
              <el-radio-group v-model="paymentMethod">
                <el-radio label="balance">
                  {{ t('user.store.payWithBalance') }}
                  <span class="balance-hint">(¥{{ userBalance }})</span>
                </el-radio>
                <el-radio label="yipay">{{ t('user.store.payWithYiPay') }}</el-radio>
              </el-radio-group>
            </div>

            <!-- 购买按钮 -->
            <div class="purchase-actions">
              <el-button
                type="primary"
                size="large"
                :loading="submitting"
                :disabled="!canSubmit"
                class="buy-btn"
                @click="handlePurchase"
              >
                {{ t('user.store.confirmPurchase') }}
              </el-button>
            </div>
          </el-card>
        </div>
      </div>
    </template>

    <el-empty v-else :description="t('user.store.productNotFound')" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft, Cpu, Memo, Coin, TopRight, DataLine,
  OfficeBuilding, Monitor, Box, Grid
} from '@element-plus/icons-vue'
import { getProductDetail } from '@/api/product'
import { createOrder, payWithBalance, createYiPayOrder, getUserBalance } from '@/api/product'
import { getSystemImages } from '@/api/user'
import { formatMemorySize, formatDiskSize, formatBandwidthSpeed } from '@/utils/unit-formatter'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const loading = ref(true)
const submitting = ref(false)
const product = ref({})
const imageList = ref([])
const selectedImage = ref('')
const selectedPeriod = ref(1)
const quantity = ref(1)
const paymentMethod = ref('balance')
const userBalance = ref(0)

// 周期选项
const periodOptions = [
  { value: 1, label: t('user.store.oneMonth') },
  { value: 3, label: t('user.store.threeMonths') },
  { value: 6, label: t('user.store.sixMonths') },
  { value: 12, label: t('user.store.oneYear') }
]

// 计算总价
const totalPrice = computed(() => {
  const price = product.value.price || 0
  return (price * selectedPeriod.value * quantity.value).toFixed(2)
})

// 是否可以提交
const canSubmit = computed(() => {
  return selectedImage.value && selectedPeriod.value > 0 && quantity.value > 0
})

// 操作系统分组
const imageGroups = computed(() => {
  const groups = {}
  imageList.value.forEach(img => {
    const category = img.category || t('user.store.other')
    if (!groups[category]) {
      groups[category] = []
    }
    groups[category].push(img)
  })
  return Object.keys(groups).map(label => ({
    label,
    options: groups[label]
  }))
})

// 获取产品图标
const getProductIcon = (type) => {
  const iconMap = { vm: Monitor, container: Box, gpu: Grid }
  return iconMap[type] || Monitor
}

const formatMemory = (memory) => formatMemorySize(memory)
const formatDisk = (disk) => formatDiskSize(disk)
const formatBandwidth = (bandwidth) => formatBandwidthSpeed(bandwidth)
const formatTraffic = (traffic) => formatDiskSize(traffic)

const getPeriodLabel = () => {
  const opt = periodOptions.find(p => p.value === selectedPeriod.value)
  return opt ? opt.label : ''
}

// 加载产品详情
const loadProduct = async () => {
  const id = route.params.id
  if (!id) return
  loading.value = true
  try {
    const res = await getProductDetail(id)
    if (res.code === 200) {
      product.value = res.data || {}
    }
  } catch (error) {
    console.error('加载产品详情失败:', error)
    ElMessage.error(error?.message || t('user.store.loadDetailFailed'))
  } finally {
    loading.value = false
  }
}

// 加载系统镜像
const loadImages = async () => {
  try {
    const res = await getSystemImages()
    if (res.code === 200) {
      imageList.value = res.data || []
    }
  } catch (error) {
    console.error('加载系统镜像失败:', error)
  }
}

// 加载用户余额
const loadBalance = async () => {
  try {
    const res = await getUserBalance()
    if (res.code === 200) {
      userBalance.value = res.data?.balance || 0
    }
  } catch (error) {
    console.error('加载余额失败:', error)
  }
}

// 购买
const handlePurchase = async () => {
  if (!canSubmit.value) {
    ElMessage.warning(t('user.store.pleaseCompleteConfig'))
    return
  }

  // 余额检查
  if (paymentMethod.value === 'balance' && parseFloat(userBalance.value) < parseFloat(totalPrice.value)) {
    ElMessageBox.confirm(
      t('user.store.insufficientBalance'),
      t('common.tip'),
      {
        confirmButtonText: t('user.store.goRecharge'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    ).then(() => {
      router.push('/user/wallet')
    }).catch(() => {})
    return
  }

  submitting.value = true
  try {
    // 创建订单
    const orderRes = await createOrder({
      product_id: product.value.id,
      image_id: selectedImage.value,
      period: selectedPeriod.value,
      quantity: quantity.value
    })

    if (orderRes.code !== 200) {
      throw new Error(orderRes.message || t('user.store.createOrderFailed'))
    }

    const orderId = orderRes.data?.id

    if (paymentMethod.value === 'balance') {
      // 余额支付
      const payRes = await payWithBalance(orderId)
      if (payRes.code === 200) {
        ElMessage.success(t('user.store.paySuccess'))
        // 跳转到实例详情
        const instanceId = payRes.data?.instance_id
        if (instanceId) {
          router.push(`/user/instances/${instanceId}`)
        } else {
          router.push('/user/orders')
        }
      } else {
        throw new Error(payRes.message || t('user.store.payFailed'))
      }
    } else {
      // 易支付
      const payRes = await createYiPayOrder({
        order_id: orderId,
        amount: totalPrice.value
      })
      if (payRes.code === 200 && payRes.data?.pay_url) {
        window.open(payRes.data.pay_url, '_blank')
        ElMessage.success(t('user.store.openPayPage'))
        router.push('/user/orders')
      } else {
        throw new Error(payRes.message || t('user.store.createPayFailed'))
      }
    }
  } catch (error) {
    console.error('购买失败:', error)
    ElMessage.error(error?.message || t('user.store.purchaseFailed'))
  } finally {
    submitting.value = false
  }
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  loadProduct()
  loadImages()
  loadBalance()
})
</script>

<style lang="scss" scoped>
.product-detail-container {
  padding: 24px;
}

.back-btn {
  margin-bottom: 16px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.detail-layout {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 24px;
}

.product-header {
  display: flex;
  gap: 20px;
  margin-bottom: 24px;
}

.product-icon-large {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100px;
  height: 100px;
  background: linear-gradient(135deg, #dcfce7 0%, #bbf7d0 100%);
  border-radius: 16px;
  flex-shrink: 0;
}

.product-title {
  h1 {
    margin: 0 0 8px 0;
    font-size: 24px;
    font-weight: 600;
    color: var(--text-color-primary);
  }

  p {
    margin: 0 0 12px 0;
    color: var(--text-color-secondary);
    line-height: 1.5;
  }
}

.product-tags {
  display: flex;
  gap: 8px;
}

.specs-detail {
  h3 {
    margin: 0 0 16px 0;
    font-size: 16px;
    font-weight: 600;
    color: var(--text-color-primary);
  }
}

.specs-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.spec-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
  background: var(--neutral-bg);
  border-radius: 8px;
  gap: 6px;

  .el-icon {
    font-size: 24px;
    color: #16a34a;
  }

  .spec-label {
    font-size: 12px;
    color: var(--text-color-secondary);
  }

  .spec-value {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-color-primary);
  }
}

.config-section {
  margin-bottom: 20px;
}

.config-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color-primary);
}

.config-select {
  width: 100%;
}

.period-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.price-summary {
  margin-bottom: 20px;
}

.price-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  font-size: 14px;
  color: var(--text-color-secondary);

  &.total {
    border-top: 1px solid var(--border-color);
    margin-top: 8px;
    padding-top: 12px;
    font-size: 16px;
    font-weight: 600;
    color: var(--text-color-primary);

    .total-price {
      color: #f56c6c;
      font-size: 24px;
    }
  }
}

.payment-methods {
  margin-bottom: 20px;

  .balance-hint {
    color: #909399;
    font-size: 12px;
    margin-left: 4px;
  }
}

.purchase-actions {
  .buy-btn {
    width: 100%;
    font-size: 16px;
    font-weight: 600;
  }
}

@media (max-width: 992px) {
  .detail-layout {
    grid-template-columns: 1fr;
  }

  .specs-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .product-detail-container {
    padding: 16px;
  }

  .product-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .specs-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
