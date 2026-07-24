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
          <el-option :label="t('admin.orders.pending')" value="pending" />
          <el-option :label="t('admin.orders.paid')" value="paid" />
          <el-option :label="t('admin.orders.active')" value="active" />
          <el-option :label="t('admin.orders.expired')" value="expired" />
          <el-option :label="t('admin.orders.cancelled')" value="cancelled" />
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
        <el-table-column :label="t('admin.orders.orderNo')" prop="order_no" min-width="160" />
        <el-table-column :label="t('admin.orders.user')" min-width="120">
          <template #default="{ row }">
            {{ row.user_name || row.user_id }}
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.orders.product')" prop="product_name" min-width="150" />
        <el-table-column :label="t('admin.orders.totalPrice')" width="120">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.total_price }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.orders.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.orders.createTime')" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.orders.actions')" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="showDetail(row)">
              {{ t('common.detail') }}
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
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
          <span class="detail-value">{{ currentOrder.order_no }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.user') }}</span>
          <span class="detail-value">{{ currentOrder.user_name || currentOrder.user_id }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.product') }}</span>
          <span class="detail-value">{{ currentOrder.product_name }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.status') }}</span>
          <el-tag :type="getStatusType(currentOrder.status)">
            {{ getStatusText(currentOrder.status) }}
          </el-tag>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.totalPrice') }}</span>
          <span class="detail-value price-highlight">¥{{ currentOrder.total_price }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.period') }}</span>
          <span class="detail-value">{{ currentOrder.period }}{{ t('admin.orders.months') }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('admin.orders.createTime') }}</span>
          <span class="detail-value">{{ formatDate(currentOrder.created_at) }}</span>
        </div>
        <div v-if="currentOrder.paid_at" class="detail-row">
          <span class="detail-label">{{ t('admin.orders.payTime') }}</span>
          <span class="detail-value">{{ formatDate(currentOrder.paid_at) }}</span>
        </div>
        <div v-if="currentOrder.expire_at" class="detail-row">
          <span class="detail-label">{{ t('admin.orders.expireTime') }}</span>
          <span class="detail-value">{{ formatDate(currentOrder.expire_at) }}</span>
        </div>
        <div v-if="currentOrder.instance_id" class="detail-row">
          <span class="detail-label">{{ t('admin.orders.instanceId') }}</span>
          <span class="detail-value">{{ currentOrder.instance_id }}</span>
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
  pending: { text: t('admin.orders.pendingPay'), type: 'warning' },
  paid: { text: t('admin.orders.paidStatus'), type: 'success' },
  active: { text: t('admin.orders.activeStatus'), type: 'primary' },
  expired: { text: t('admin.orders.expiredStatus'), type: 'info' },
  cancelled: { text: t('admin.orders.cancelledStatus'), type: 'info' }
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
    const { getOrderList } = await import('@/api/product')
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: searchKeyword.value || undefined,
      status: searchStatus.value || undefined
    }
    const res = await getOrderList(params)
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

const showDetail = (row) => {
  currentOrder.value = row
  detailVisible.value = true
}

// 手动开通
const handleManualOpen = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('admin.orders.confirmManualOpen', { orderNo: row.order_no }),
      t('common.tip'),
      { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' }
    )
    // await manualOpenOrder(row.id)
    ElMessage.success(t('admin.orders.manualOpenSuccess'))
    loadOrders()
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
</style>
