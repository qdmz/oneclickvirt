<template>
  <div class="orders-container">
    <!-- 页面头部 -->
    <div class="orders-header">
      <h1>{{ t('user.orders.title') }}</h1>
    </div>

    <!-- 状态筛选 -->
    <div class="filter-bar">
      <el-radio-group v-model="filterStatus" @change="handleFilterChange">
        <el-radio-button label="">{{ t('user.orders.all') }}</el-radio-button>
        <el-radio-button label="pending">{{ t('user.orders.pending') }}</el-radio-button>
        <el-radio-button label="paid">{{ t('user.orders.paid') }}</el-radio-button>
        <el-radio-button label="active">{{ t('user.orders.active') }}</el-radio-button>
        <el-radio-button label="expired">{{ t('user.orders.expired') }}</el-radio-button>
      </el-radio-group>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-loading-directive />
      <div class="loading-text">{{ t('common.loading') }}</div>
    </div>

    <!-- 订单列表 -->
    <template v-else>
      <el-empty v-if="orderList.length === 0" :description="t('user.orders.noOrders')" />

      <div v-else class="order-list">
        <el-card
          v-for="order in orderList"
          :key="order.id"
          class="order-card"
          shadow="hover"
        >
          <div class="order-header">
            <div class="order-info">
              <span class="order-id">{{ t('user.orders.orderNo') }}: {{ order.order_no }}</span>
              <span class="order-date">{{ formatDate(order.created_at) }}</span>
            </div>
            <el-tag :type="getStatusType(order.status)">
              {{ getStatusText(order.status) }}
            </el-tag>
          </div>

          <el-divider />

          <div class="order-body">
            <div class="product-info">
              <el-icon :size="32" color="#16a34a"><Box /></el-icon>
              <div class="product-detail">
                <div class="product-name">{{ order.product_name }}</div>
                <div class="product-specs">
                  {{ order.cpu }}核 / {{ formatMemory(order.memory) }} / {{ formatDisk(order.disk) }}
                </div>
              </div>
            </div>
            <div class="order-price">
              <div class="price">¥{{ order.total_price }}</div>
              <div class="period">{{ order.period }}{{ t('user.orders.months') }}</div>
            </div>
          </div>

          <div v-if="order.instance_id" class="order-instance">
            <el-link type="primary" @click="goToInstance(order.instance_id)">
              {{ t('user.orders.viewInstance') }}
            </el-link>
          </div>

          <div class="order-actions">
            <el-button
              v-if="order.status === 'pending'"
              type="primary"
              size="small"
              @click="handlePay(order)"
            >
              {{ t('user.orders.payNow') }}
            </el-button>
            <el-button
              v-if="order.status === 'pending'"
              size="small"
              @click="handleCancel(order)"
            >
              {{ t('user.orders.cancel') }}
            </el-button>
            <el-button
              v-if="['active', 'paid'].includes(order.status)"
              type="primary"
              size="small"
              @click="handleRenew(order)"
            >
              {{ t('user.orders.renew') }}
            </el-button>
            <el-button
              size="small"
              @click="showDetail(order)"
            >
              {{ t('user.orders.detail') }}
            </el-button>
          </div>
        </el-card>
      </div>

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

    <!-- 订单详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      :title="t('user.orders.orderDetail')"
      width="600px"
    >
      <div v-if="currentOrder" class="order-detail-content">
        <div class="detail-row">
          <span class="detail-label">{{ t('user.orders.orderNo') }}</span>
          <span class="detail-value">{{ currentOrder.order_no }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('user.orders.productName') }}</span>
          <span class="detail-value">{{ currentOrder.product_name }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('user.orders.status') }}</span>
          <el-tag :type="getStatusType(currentOrder.status)">
            {{ getStatusText(currentOrder.status) }}
          </el-tag>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('user.orders.totalPrice') }}</span>
          <span class="detail-value price-highlight">¥{{ currentOrder.total_price }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('user.orders.period') }}</span>
          <span class="detail-value">{{ currentOrder.period }}{{ t('user.orders.months') }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('user.orders.createTime') }}</span>
          <span class="detail-value">{{ formatDate(currentOrder.created_at) }}</span>
        </div>
        <div v-if="currentOrder.paid_at" class="detail-row">
          <span class="detail-label">{{ t('user.orders.payTime') }}</span>
          <span class="detail-value">{{ formatDate(currentOrder.paid_at) }}</span>
        </div>
        <div v-if="currentOrder.expire_at" class="detail-row">
          <span class="detail-label">{{ t('user.orders.expireTime') }}</span>
          <span class="detail-value">{{ formatDate(currentOrder.expire_at) }}</span>
        </div>
      </div>
    </el-dialog>

    <!-- 续费弹窗 -->
    <el-dialog
      v-model="renewVisible"
      :title="t('user.orders.renewOrder')"
      width="400px"
    >
      <div v-if="currentOrder" class="renew-content">
        <div class="renew-product">{{ currentOrder.product_name }}</div>
        <div class="renew-price">{{ t('user.orders.unitPrice') }}: ¥{{ currentOrder.product_price }}/{{ t('user.orders.month') }}</div>
        <el-form :model="renewForm" label-width="80px">
          <el-form-item :label="t('user.orders.renewPeriod')">
            <el-radio-group v-model="renewForm.period">
              <el-radio-button :label="1">1{{ t('user.orders.month') }}</el-radio-button>
              <el-radio-button :label="3">3{{ t('user.orders.month') }}</el-radio-button>
              <el-radio-button :label="6">6{{ t('user.orders.month') }}</el-radio-button>
              <el-radio-button :label="12">12{{ t('user.orders.month') }}</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item :label="t('user.orders.renewPrice')">
            <span class="price-highlight">¥{{ renewPrice }}</span>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="renewVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="renewLoading" @click="confirmRenew">
          {{ t('user.orders.confirmRenew') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Box } from '@element-plus/icons-vue'
import { getOrderList, cancelOrder, renewOrder, payWithBalance, createYiPayOrder } from '@/api/product'
import { formatMemorySize, formatDiskSize } from '@/utils/unit-formatter'

const router = useRouter()
const { t, locale } = useI18n()

const loading = ref(true)
const orderList = ref([])
const filterStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const detailVisible = ref(false)
const currentOrder = ref(null)
const renewVisible = ref(false)
const renewLoading = ref(false)
const renewForm = ref({ period: 1 })

// 续费价格
const renewPrice = computed(() => {
  if (!currentOrder.value) return '0.00'
  return (currentOrder.value.product_price * renewForm.value.period).toFixed(2)
})

// 状态映射
const statusMap = {
  pending: { text: t('user.orders.pendingPay'), type: 'warning' },
  paid: { text: t('user.orders.paidStatus'), type: 'success' },
  active: { text: t('user.orders.activeStatus'), type: 'primary' },
  expired: { text: t('user.orders.expiredStatus'), type: 'info' },
  cancelled: { text: t('user.orders.cancelledStatus'), type: 'info' }
}

const getStatusType = (status) => statusMap[status]?.type || 'info'
const getStatusText = (status) => statusMap[status]?.text || status

const formatMemory = (memory) => formatMemorySize(memory)
const formatDisk = (disk) => formatDiskSize(disk)

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString(locale.value === 'en-US' ? 'en-US' : 'zh-CN')
}

// 加载订单列表
const loadOrders = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      status: filterStatus.value || undefined
    }
    const res = await getOrderList(params)
    if (res.code === 200) {
      orderList.value = res.data?.list || res.data?.items || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载订单列表失败:', error)
    ElMessage.error(error?.message || t('user.orders.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleFilterChange = () => {
  currentPage.value = 1
  loadOrders()
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadOrders()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadOrders()
}

// 支付
const handlePay = async (order) => {
  try {
    await ElMessageBox.confirm(
      t('user.orders.confirmPay', { amount: order.total_price }),
      t('common.tip'),
      { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel') }
    )
    const res = await payWithBalance(order.id)
    if (res.code === 200) {
      ElMessage.success(t('user.orders.paySuccess'))
      loadOrders()
    } else {
      // 余额不足，尝试易支付
      const payRes = await createYiPayOrder({ order_id: order.id, amount: order.total_price })
      if (payRes.code === 200 && payRes.data?.pay_url) {
        window.open(payRes.data.pay_url, '_blank')
      }
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error?.message || t('user.orders.payFailed'))
    }
  }
}

// 取消订单
const handleCancel = async (order) => {
  try {
    await ElMessageBox.confirm(
      t('user.orders.confirmCancel'),
      t('common.tip'),
      { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' }
    )
    const res = await cancelOrder(order.id)
    if (res.code === 200) {
      ElMessage.success(t('user.orders.cancelSuccess'))
      loadOrders()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error?.message || t('user.orders.cancelFailed'))
    }
  }
}

// 续费
const handleRenew = (order) => {
  currentOrder.value = order
  renewForm.value.period = 1
  renewVisible.value = true
}

const confirmRenew = async () => {
  renewLoading.value = true
  try {
    const res = await renewOrder(currentOrder.value.id, { period: renewForm.value.period })
    if (res.code === 200) {
      ElMessage.success(t('user.orders.renewSuccess'))
      renewVisible.value = false
      loadOrders()
    }
  } catch (error) {
    ElMessage.error(error?.message || t('user.orders.renewFailed'))
  } finally {
    renewLoading.value = false
  }
}

// 查看详情
const showDetail = (order) => {
  currentOrder.value = order
  detailVisible.value = true
}

// 跳转实例
const goToInstance = (id) => {
  router.push(`/user/instances/${id}`)
}

onMounted(() => {
  loadOrders()
})
</script>

<style lang="scss" scoped>
.orders-container {
  padding: 24px;
}

.orders-header {
  margin-bottom: 20px;

  h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
    color: var(--text-color-primary);
  }
}

.filter-bar {
  margin-bottom: 20px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
}

.order-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.order-card {
  :deep(.el-card__body) {
    padding: 16px 20px;
  }
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-info {
  display: flex;
  gap: 16px;
  align-items: center;

  .order-id {
    font-weight: 600;
    color: var(--text-color-primary);
  }

  .order-date {
    font-size: 13px;
    color: var(--text-color-secondary);
  }
}

.order-body {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.product-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.product-detail {
  .product-name {
    font-weight: 600;
    color: var(--text-color-primary);
    margin-bottom: 4px;
  }

  .product-specs {
    font-size: 13px;
    color: var(--text-color-secondary);
  }
}

.order-price {
  text-align: right;

  .price {
    font-size: 20px;
    font-weight: 700;
    color: #f56c6c;
  }

  .period {
    font-size: 12px;
    color: var(--text-color-secondary);
  }
}

.order-instance {
  margin-bottom: 12px;
}

.order-actions {
  display: flex;
  gap: 8px;
}

.pagination-wrapper {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}

.order-detail-content {
  .detail-row {
    display: flex;
    justify-content: space-between;
    padding: 10px 0;
    border-bottom: 1px solid var(--border-color);

    &:last-child {
      border-bottom: none;
    }
  }

  .detail-label {
    color: var(--text-color-secondary);
  }

  .detail-value {
    font-weight: 500;
    color: var(--text-color-primary);
  }

  .price-highlight {
    color: #f56c6c;
    font-size: 18px;
    font-weight: 700;
  }
}

.renew-content {
  .renew-product {
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 8px;
  }

  .renew-price {
    color: var(--text-color-secondary);
    margin-bottom: 16px;
  }
}

.price-highlight {
  color: #f56c6c;
  font-size: 18px;
  font-weight: 700;
}

@media (max-width: 768px) {
  .orders-container {
    padding: 16px;
  }

  .order-body {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .order-price {
    text-align: left;
  }
}
</style>
