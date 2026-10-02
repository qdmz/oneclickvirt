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
        <el-radio-button label="unpaid">{{ t('user.orders.pending') }}</el-radio-button>
        <el-radio-button label="paid">{{ t('user.orders.paid') }}</el-radio-button>
        <el-radio-button label="provisioned">{{ t('user.orders.active') }}</el-radio-button>
        <el-radio-button label="provisionFailed">{{ t('user.orders.provisionFailed') }}</el-radio-button>
      </el-radio-group>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="4" animated />
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
              <span class="order-id">{{ t('user.orders.orderNo') }}: {{ order.orderNo }}</span>
              <span class="order-date">{{ formatDate(order.createdAt) }}</span>
              <el-tag v-if="order.isRenewal" size="small" type="warning" effect="plain">
                {{ t('user.orders.renewalOrder') }}
              </el-tag>
            </div>
            <div class="order-tags">
              <el-tag :type="payStatusMeta(order.paymentStatus).type" size="small">
                {{ payStatusMeta(order.paymentStatus).text }}
              </el-tag>
              <el-tag
                v-if="order.paymentStatus === 1"
                :type="provisionStatusMeta(order).type"
                size="small"
                effect="plain"
              >
                {{ provisionStatusMeta(order).text }}
              </el-tag>
            </div>
          </div>

          <el-divider />

          <div class="order-body">
            <div class="product-info">
              <el-icon :size="32" color="#16a34a"><Box /></el-icon>
              <div class="product-detail">
                <div class="product-name">{{ order.productName || '-' }}</div>
                <div class="product-specs">{{ specsText(order) }}</div>
              </div>
            </div>
            <div class="order-price">
              <div class="price">¥{{ formatAmount(order.totalAmount) }}</div>
              <div class="period">{{ periodText(order) }}</div>
            </div>
          </div>

          <div v-if="order.instanceId" class="order-instance">
            <el-link type="primary" @click="goToInstance(order.instanceId)">
              {{ t('user.orders.viewInstance') }}
            </el-link>
          </div>

          <div class="order-actions">
            <el-button
              v-if="order.paymentStatus === 0"
              type="primary"
              size="small"
              @click="handlePay(order)"
            >
              {{ t('user.orders.payNow') }}
            </el-button>
            <el-button
              v-if="order.paymentStatus === 0"
              size="small"
              @click="handleCancel(order)"
            >
              {{ t('user.orders.cancel') }}
            </el-button>
            <el-button
              v-if="order.paymentStatus === 1 && order.provisionStatus === 2"
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
      width="760px"
      class="order-detail-dialog"
    >
      <div v-loading="detailLoading" class="order-detail-content">
        <template v-if="currentOrder">
          <el-descriptions :title="t('user.orders.baseInfo')" :column="2" border size="small">
            <el-descriptions-item :label="t('user.orders.orderNo')">
              {{ currentOrder.orderNo || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.orderType')">
              {{ currentOrder.isRenewal ? t('user.orders.renewalOrder') : t('user.orders.newOrder') }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.paymentStatus')">
              <el-tag :type="payStatusMeta(currentOrder.paymentStatus).type" size="small">
                {{ payStatusMeta(currentOrder.paymentStatus).text }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.provisionStatus')">
              <el-tag :type="provisionStatusMeta(currentOrder).type" size="small" effect="plain">
                {{ provisionStatusMeta(currentOrder).text }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.payMethod')">
              {{ payMethodText(currentOrder.paymentMethod) }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.tradeNo')">
              {{ currentOrder.tradeNo || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.createTime')">
              {{ formatDate(currentOrder.createdAt) }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.payTime')">
              {{ formatDate(currentOrder.paidAt) }}
            </el-descriptions-item>
          </el-descriptions>

          <el-descriptions
            :title="t('user.orders.configInfo')"
            :column="2"
            border
            size="small"
            class="detail-section"
          >
            <el-descriptions-item :label="t('user.orders.productName')">
              {{ currentOrder.productName || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.productType')">
              {{ productTypeText(currentOrder.productType) }}
            </el-descriptions-item>
            <el-descriptions-item label="CPU">
              {{ currentOrder.cpu ? currentOrder.cpu + ' ' + t('user.orders.core') : '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.memory')">
              {{ currentOrder.memory ? formatMemory(currentOrder.memory) : '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.disk')">
              {{ currentOrder.disk ? formatDisk(currentOrder.disk) : '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.bandwidth')">
              {{ currentOrder.bandwidth ? currentOrder.bandwidth + ' Mbps' : '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.traffic')">
              {{ trafficText(currentOrder.traffic) }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.image')">
              {{ currentOrder.imageName || '-' }}
            </el-descriptions-item>
          </el-descriptions>

          <el-descriptions
            :title="t('user.orders.feeInfo')"
            :column="2"
            border
            size="small"
            class="detail-section"
          >
            <el-descriptions-item :label="t('user.orders.unitPrice')">
              ¥{{ formatAmount(currentOrder.price) }} / {{ periodUnitText(currentOrder.periodType, currentOrder.periodValue) }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.quantity')">
              {{ currentOrder.quantity || 1 }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.period')">
              {{ periodText(currentOrder) }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.totalPrice')">
              <span class="price-highlight">¥{{ formatAmount(currentOrder.totalAmount) }}</span>
            </el-descriptions-item>
          </el-descriptions>

          <el-descriptions
            :title="t('user.orders.provisionInfo')"
            :column="2"
            border
            size="small"
            class="detail-section"
          >
            <el-descriptions-item :label="t('user.orders.instance')">
              <el-link
                v-if="currentOrder.instanceId"
                type="primary"
                @click="goToInstance(currentOrder.instanceId)"
              >
                #{{ currentOrder.instanceId }} {{ t('user.orders.viewInstance') }}
              </el-link>
              <span v-else>-</span>
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.provisionedAt')">
              {{ formatDate(currentOrder.provisionedAt) }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.expireTime')">
              {{ formatDate(currentOrder.expireAt) }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('user.orders.remark')">
              {{ currentOrder.remark || '-' }}
            </el-descriptions-item>
          </el-descriptions>
        </template>
      </div>
    </el-dialog>

    <!-- 续费弹窗 -->
    <el-dialog
      v-model="renewVisible"
      :title="t('user.orders.renewOrder')"
      width="420px"
    >
      <div v-if="currentOrder" class="renew-content">
        <div class="renew-product">{{ currentOrder.productName }}</div>
        <div class="renew-price">
          {{ t('user.orders.unitPrice') }}: ¥{{ formatAmount(currentOrder.price) }}
          / {{ periodUnitText(currentOrder.periodType, currentOrder.periodValue) }}
        </div>
        <el-form :model="renewForm" label-width="90px">
          <el-form-item :label="t('user.orders.renewPeriod')">
            <el-input-number v-model="renewForm.quantity" :min="1" :max="36" />
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
import {
  getOrderList,
  getOrderDetail,
  cancelOrder,
  renewOrder,
  payWithBalance,
  createYiPayOrder
} from '@/api/product'
import { formatMemorySize, formatDiskSize } from '@/utils/unit-formatter'
import { useSiteStore } from '@/pinia/modules/site'

const router = useRouter()
const { t, locale } = useI18n()
const siteStore = useSiteStore()

const loading = ref(true)
const orderList = ref([])
const filterStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const detailVisible = ref(false)
const detailLoading = ref(false)
const currentOrder = ref(null)
const renewVisible = ref(false)
const renewLoading = ref(false)
const renewForm = ref({ quantity: 1 })

// 续费价格
const renewPrice = computed(() => {
  if (!currentOrder.value) return '0.00'
  return (Number(currentOrder.value.price || 0) * Number(renewForm.value.quantity || 1)).toFixed(2)
})

const formatAmount = (val) => Number(val || 0).toFixed(2)

// 支付状态: 0=未支付 1=已支付 2=支付失败 3=已退款
const payStatusMeta = (status) => {
  switch (Number(status)) {
    case 1: return { text: t('user.orders.paidStatus'), type: 'success' }
    case 2: return { text: t('user.orders.payFailedStatus'), type: 'danger' }
    case 3: return { text: t('user.orders.refundedStatus'), type: 'info' }
    default: return { text: t('user.orders.pendingPay'), type: 'warning' }
  }
}

// 开通状态: 0=待开通 1=开通中 2=已开通 3=开通失败（已开通但已过期时展示为已过期）
const provisionStatusMeta = (order) => {
  const status = Number(order?.provisionStatus)
  if (status === 2) {
    if (order?.expireAt && new Date(order.expireAt).getTime() < Date.now()) {
      return { text: t('user.orders.expiredStatus'), type: 'info' }
    }
    return { text: t('user.orders.activeStatus'), type: 'success' }
  }
  if (status === 1) return { text: t('user.orders.provisioning'), type: 'primary' }
  if (status === 3) return { text: t('user.orders.provisionFailed'), type: 'danger' }
  return { text: t('user.orders.provisionPending'), type: 'warning' }
}

const periodUnitMap = () => ({
  hour: t('user.orders.unitHour'),
  day: t('user.orders.unitDay'),
  month: t('user.orders.unitMonth'),
  year: t('user.orders.unitYear')
})

// 单个周期的描述，例如 "1 个月" / "3 天"
const periodUnitText = (periodType, periodValue) => {
  const unit = periodUnitMap()[periodType] || t('user.orders.unitMonth')
  const value = Number(periodValue || 1)
  return value > 1 ? `${value} ${unit}` : unit
}

// 订单总周期：periodValue * quantity
const periodText = (order) => {
  if (!order) return '-'
  const unit = periodUnitMap()[order.periodType] || t('user.orders.unitMonth')
  const totalUnits = Number(order.periodValue || 1) * Number(order.quantity || 1)
  return `${totalUnits} ${unit}`
}

const specsText = (order) => {
  const parts = []
  if (order.cpu) parts.push(`${order.cpu} ${t('user.orders.core')}`)
  if (order.memory) parts.push(formatMemorySize(order.memory))
  if (order.disk) parts.push(formatDiskSize(order.disk))
  if (order.bandwidth) parts.push(`${order.bandwidth} Mbps`)
  return parts.length ? parts.join(' / ') : '-'
}

const trafficText = (traffic) => {
  const value = Number(traffic || 0)
  if (!value || value < 0) return t('user.orders.unlimited')
  return `${value} GB`
}

const payMethodText = (method) => {
  const map = {
    balance: t('user.orders.payMethodBalance'),
    yipay: t('user.orders.payMethodYiPay'),
    alipay: t('user.orders.payMethodAlipay'),
    wxpay: t('user.orders.payMethodWxPay'),
    admin: t('user.orders.payMethodAdmin')
  }
  return map[method] || method || '-'
}

const productTypeText = (type) => {
  const map = {
    vm: t('user.orders.productTypeVm'),
    container: t('user.orders.productTypeContainer'),
    lxc: t('user.orders.productTypeContainer')
  }
  return map[type] || type || '-'
}

const formatMemory = (memory) => formatMemorySize(memory)
const formatDisk = (disk) => formatDiskSize(disk)

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  if (Number.isNaN(date.getTime()) || date.getFullYear() < 1971) return '-'
  return date.toLocaleString(locale.value === 'en-US' ? 'en-US' : 'zh-CN')
}

// 筛选条件 → 后端查询参数（后端只认 paymentStatus / provisionStatus 两个整型字段）
const buildFilterParams = () => {
  switch (filterStatus.value) {
    case 'unpaid': return { paymentStatus: 0 }
    case 'paid': return { paymentStatus: 1 }
    case 'provisioned': return { provisionStatus: 2 }
    case 'provisionFailed': return { provisionStatus: 3 }
    default: return {}
  }
}

// 加载订单列表
const loadOrders = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      ...buildFilterParams()
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

// 判断是否为「余额不足」错误：请求拦截器对非 200 的业务码会 reject，
// 因此余额支付失败只能在 catch 中识别，不能靠 res.code 分支。
const isInsufficientBalance = (error) => {
  const msg = String(error?.message || error?.msg || '')
  return msg.includes('余额不足') || msg.toLowerCase().includes('insufficient balance')
}

// 支付
const handlePay = async (order) => {
  try {
    await ElMessageBox.confirm(
      t('user.orders.confirmPay', { amount: formatAmount(order.totalAmount) }),
      t('common.tip'),
      { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel') }
    )
  } catch {
    return // 用户取消
  }

  // 优先使用余额支付
  try {
    const res = await payWithBalance(order.id)
    if (res?.code === 200) {
      ElMessage.success(t('user.orders.paySuccess'))
      loadOrders()
      return
    }
    throw new Error(res?.message || t('user.orders.payFailed'))
  } catch (error) {
    if (!isInsufficientBalance(error)) {
      ElMessage.error(error?.message || t('user.orders.payFailed'))
      return
    }
  }

  // 余额不足：/payments/yipay 目前只支持余额充值（仅接受 amount + payType），
  // 因此按订单金额发起充值，用户支付完成后回到订单页用余额完成支付。
  try {
    const payRes = await createYiPayOrder({
      amount: Number(order.totalAmount),
      payType: siteStore.enabledPayTypes[0] || 'alipay'
    })
    const payUrl = payRes?.data?.payURL || payRes?.data?.pay_url
    if (payRes?.code === 200 && payUrl) {
      window.open(payUrl, '_blank')
      ElMessage.success(t('user.orders.balanceNotEnough'))
    } else {
      throw new Error(payRes?.message || t('user.orders.payFailed'))
    }
  } catch (error) {
    ElMessage.error(error?.message || t('user.orders.payFailed'))
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
  renewForm.value.quantity = 1
  renewVisible.value = true
}

const confirmRenew = async () => {
  renewLoading.value = true
  try {
    const res = await renewOrder(currentOrder.value.id, { quantity: renewForm.value.quantity })
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

// 查看详情：先用列表数据占位，再请求详情接口补全字段（列表接口返回的是完整订单对象，
// 但详情接口才是权威来源，避免列表缓存导致状态过期）
const showDetail = async (order) => {
  currentOrder.value = order
  detailVisible.value = true
  detailLoading.value = true
  try {
    const res = await getOrderDetail(order.id)
    if (res?.code === 200 && res.data) {
      currentOrder.value = res.data
    }
  } catch (error) {
    console.error('加载订单详情失败:', error)
    ElMessage.error(error?.message || t('user.orders.detailLoadFailed'))
  } finally {
    detailLoading.value = false
  }
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
  min-height: 300px;
  padding: 12px;
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
  gap: 12px;
  flex-wrap: wrap;
}

.order-tags {
  display: flex;
  gap: 8px;
  align-items: center;
}

.order-info {
  display: flex;
  gap: 16px;
  align-items: center;
  flex-wrap: wrap;

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
  min-height: 120px;

  .detail-section {
    margin-top: 20px;
  }

  :deep(.el-descriptions__title) {
    font-size: 14px;
    font-weight: 600;
  }

  :deep(.el-descriptions__label) {
    width: 110px;
    color: var(--text-color-secondary);
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
  font-size: 16px;
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
