<template>
  <div class="orders-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('admin.orders.title') }}</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="toolbar">
        <el-input
          v-model="searchKeyword"
          :placeholder="t('admin.orders.searchPlaceholder')"
          style="width: 250px;"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select
          v-model="searchStatus"
          :placeholder="t('admin.orders.selectStatus')"
          style="width: 150px; margin-left: 10px;"
          clearable
        >
          <el-option :label="t('admin.orders.all')" value="" />
          <el-option :label="t('admin.orders.pending')" :value="0" />
          <el-option :label="t('admin.orders.paid')" :value="1" />
          <el-option :label="t('admin.orders.paymentFailed')" :value="2" />
          <el-option :label="t('admin.orders.refunded')" :value="3" />
        </el-select>
        <el-button type="primary" style="margin-left: 10px;" @click="handleSearch">
          {{ t('common.search') }}
        </el-button>
        <el-button style="margin-left: 10px;" @click="resetFilters">
          {{ t('common.reset') }}
        </el-button>
      </div>

      <!-- 订单表格 -->
      <el-table v-loading="loading" :data="orderList" stripe style="width: 100%">
        <el-table-column :label="t('admin.orders.orderNo')" prop="orderNo" min-width="160" />
        <el-table-column :label="t('admin.orders.user')" min-width="120">
          <template #default="{ row }">
            {{ row.username || row.userId }}
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.orders.product')" prop="productName" min-width="150" />
        <el-table-column :label="t('admin.orders.totalPrice')" width="120">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.totalAmount }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.orders.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.paymentStatus)" size="small">
              {{ getStatusText(row.paymentStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.orders.createTime')" width="160">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.orders.actions')" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="showDetail(row)">
              {{ t('common.detail') }}
            </el-button>
            <el-button
              v-if="row.paymentStatus === 1 && row.provisionStatus !== 2"
              type="success"
              size="small"
              @click="handleManualOpen(row)"
            >
              {{ t('admin.orders.manualOpen') }}
            </el-button>
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
    </el-card>

    <!-- 订单详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      :title="t('admin.orders.orderDetail')"
      width="600px"
    >
      <div v-if="currentOrder" class="order-detail-content">
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.orderNo') }}</span>
          <span class="detail-value">{{ currentOrder.orderNo }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.user') }}</span>
          <span class="detail-value">{{ currentOrder.username || currentOrder.userId }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.product') }}</span>
          <span class="detail-value">{{ currentOrder.productName }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.status') }}</span>
          <el-tag :type="getStatusType(currentOrder.paymentStatus)">
            {{ getStatusText(currentOrder.paymentStatus) }}
          </el-tag>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.totalPrice') }}</span>
          <span class="detail-value price-highlight">¥{{ currentOrder.totalAmount }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.period') }}</span>
          <span class="detail-value">{{ currentOrder.quantity }}{{ t('admin.orders.months') }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.createTime') }}</span>
          <span class="detail-value">{{ formatDate(currentOrder.createdAt) }}</span>
        </div>
        <div v-if="currentOrder.paidAt" class="detail-row">
          <span class="detail-label">{{ t('admin.orders.payTime') }}</span>
          <span class="detail-value">{{ formatDate(currentOrder.paidAt) }}</span>
        </div>
        <div v-if="currentOrder.expireAt" class="detail-row">
          <span class="detail-label">{{ t('admin.orders.expireTime') }}</span>
          <span class="detail-value">{{ formatDate(currentOrder.expireAt) }}</span>
        </div>
        <div v-if="currentOrder.instanceId" class="detail-row">
          <span class="detail-label">{{ t('admin.orders.instanceId') }}</span>
          <span class="detail-value">{{ currentOrder.instanceId }}</span>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { getAdminOrderList, getAdminOrderDetail, manualProvision } from '@/api/admin'

const { t, locale } = useI18n()

const loading = ref(true)
const orderList = ref([])
const searchKeyword = ref('')
const searchStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const detailVisible = ref(false)
const currentOrder = ref(null)

const statusMap = {
  0: { text: t('admin.orders.pendingPay'), type: 'warning' },
  1: { text: t('admin.orders.paidStatus'), type: 'success' },
  2: { text: t('admin.orders.paymentFailedStatus'), type: 'danger' },
  3: { text: t('admin.orders.refundedStatus'), type: 'info' }
}

const getStatusType = (status) => statusMap[status]?.type || 'info'
const getStatusText = (status) => statusMap[status]?.text || status

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
      keyword: searchKeyword.value || undefined,
      paymentStatus: searchStatus.value === '' ? undefined : searchStatus.value
    }
    const res = await getAdminOrderList(params)
    if (res.code === 200) {
      orderList.value = res.data?.list || res.data?.items || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载订单列表失败:', error)
    ElMessage.error(error?.message || t('admin.orders.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadOrders()
}

const resetFilters = () => {
  searchKeyword.value = ''
  searchStatus.value = ''
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

const showDetail = async (row) => {
  currentOrder.value = row
  detailVisible.value = true
  try {
    const res = await getAdminOrderDetail(row.id)
    if (res.code === 200 && res.data) {
      currentOrder.value = { ...row, ...res.data }
    }
  } catch (error) {
    console.error('Failed to load order detail:', error)
  }
}

// 手动开通
const handleManualOpen = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('admin.orders.confirmManualOpen', { orderNo: row.orderNo }),
      t('common.tip'),
      { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' }
    )
    const res = await manualProvision(row.id)
    if (res.code === 200) {
      ElMessage.success(t('admin.orders.manualOpenSuccess'))
      loadOrders()
    } else {
      ElMessage.error(res.message || t('admin.orders.manualOpenFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error?.message || t('admin.orders.manualOpenFailed'))
    }
  }
}

onMounted(() => {
  loadOrders()
})
</script>

<style lang="scss" scoped>
.orders-container {
  padding: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.toolbar {
  margin-bottom: 20px;
}

.price-text {
  color: #f56c6c;
  font-weight: 600;
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
