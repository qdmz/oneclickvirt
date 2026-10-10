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
</style>
