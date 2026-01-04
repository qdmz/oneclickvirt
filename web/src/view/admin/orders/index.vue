<template>
  <div class="orders-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>订单管理</span>
          <div class="header-actions">
            <el-select
              v-model="searchForm.status"
              placeholder="订单状态"
              clearable
              style="width: 120px; margin-right: 10px;"
            >
              <el-option
                label="待支付"
                value="pending"
              />
              <el-option
                label="已支付"
                value="paid"
              />
              <el-option
                label="已取消"
                value="cancelled"
              />
              <el-option
                label="已退款"
                value="refunded"
              />
              <el-option
                label="已过期"
                value="expired"
              />
            </el-select>
            <el-input
              v-model="searchForm.orderNo"
              placeholder="订单号"
              clearable
              style="width: 200px; margin-right: 10px;"
            />
            <el-input
              v-model="searchForm.username"
              placeholder="用户名"
              clearable
              style="width: 150px;"
            />
            <el-button
              type="primary"
              @click="handleSearch"
            >
              搜索
            </el-button>
            <el-button @click="handleReset">
              重置
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
      >
        <el-table-column
          prop="id"
          label="ID"
          width="80"
        />
        <el-table-column
          prop="orderNo"
          label="订单号"
          width="200"
        />
        <el-table-column
          prop="username"
          label="用户"
          width="120"
        />
        <el-table-column
          label="订单金额"
          width="100"
        >
          <template #default="{ row }">
            ¥{{ (row.amount / 100).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column
          label="实付金额"
          width="100"
        >
          <template #default="{ row }">
            ¥{{ (row.paidAmount / 100).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column
          label="订单状态"
          width="90"
        >
          <template #default="{ row }">
            <el-tag
              :type="getStatusTagType(row.status)"
              size="small"
            >
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          label="支付方式"
          width="90"
        >
          <template #default="{ row }">
            {{ getPaymentMethodName(row.paymentMethod) }}
          </template>
        </el-table-column>
        <el-table-column
          label="产品信息"
          width="150"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ getProductInfo(row) }}
          </template>
        </el-table-column>
        <el-table-column
          label="创建时间"
          width="160"
        >
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column
          label="支付时间"
          width="160"
        >
          <template #default="{ row }">
            {{ row.paymentTime ? formatTime(row.paymentTime) : '-' }}
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          width="180"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="handleView(row)"
            >
              详情
            </el-button>
            <el-button
              v-if="row.status === 'paid'"
              type="warning"
              size="small"
              @click="handleRefund(row)"
            >
              退款
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="danger"
              size="small"
              @click="handleCancel(row)"
            >
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        class="mt-20"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>

    <!-- 订单详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="订单详情"
      width="700px"
    >
      <el-descriptions
        v-if="currentOrder"
        :column="2"
        border
      >
        <el-descriptions-item
          label="订单号"
          :span="2"
        >
          {{ currentOrder.orderNo }}
        </el-descriptions-item>
        <el-descriptions-item label="用户">
          {{ currentOrder.username || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="用户ID">
          {{ currentOrder.userId || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="订单金额">
          ¥{{ (currentOrder.amount / 100).toFixed(2) }}
        </el-descriptions-item>
        <el-descriptions-item label="实付金额">
          ¥{{ (currentOrder.paidAmount / 100).toFixed(2) }}
        </el-descriptions-item>
        <el-descriptions-item
          label="订单状态"
          :span="2"
        >
          <el-tag :type="getStatusTagType(currentOrder.status)">
            {{ getStatusName(currentOrder.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="支付方式">
          {{ getPaymentMethodName(currentOrder.paymentMethod) }}
        </el-descriptions-item>
        <el-descriptions-item label="支付交易号">
          {{ currentOrder.transactionId || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ formatTime(currentOrder.createdAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="支付时间">
          {{ currentOrder.paymentTime ? formatTime(currentOrder.paymentTime) : '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="过期时间">
          {{ formatTime(currentOrder.expireAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">
          {{ currentOrder.clientIp || '-' }}
        </el-descriptions-item>
        <el-descriptions-item
          label="备注"
          :span="2"
        >
          {{ currentOrder.remark || '-' }}
        </el-descriptions-item>
        <el-descriptions-item
          v-if="currentOrder.productData"
          label="产品信息"
          :span="2"
        >
          <pre style="white-space: pre-wrap; word-break: break-all;">{{ formatProductData(currentOrder.productData) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 退款对话框 -->
    <el-dialog
      v-model="refundDialogVisible"
      title="订单退款"
      width="500px"
    >
      <el-form
        :model="refundForm"
        label-width="80px"
      >
        <el-form-item label="退款金额">
          <el-input-number
            v-model="refundForm.amount"
            :min="1"
            :max="currentOrder?.amount/100 || 0"
            :precision="2"
          />
          <span style="margin-left: 10px; color: #909399;">订单金额: ¥{{ currentOrder ? (currentOrder.amount / 100).toFixed(2) : 0 }}</span>
        </el-form-item>
        <el-form-item label="退款原因">
          <el-input
            v-model="refundForm.reason"
            type="textarea"
            :rows="3"
            placeholder="请输入退款原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="refundDialogVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          @click="handleConfirmRefund"
        >
          确认退款
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getOrders, getOrder, deleteOrder, cancelOrder, refundOrder } from '@/api/admin'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const tableData = ref([])
const detailDialogVisible = ref(false)
const refundDialogVisible = ref(false)
const currentOrder = ref(null)

const searchForm = ref({
  status: '',
  orderNo: '',
  username: ''
})

const refundForm = ref({
  amount: 0,
  reason: ''
})

// 加载订单列表
const loadOrders = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value
    }
    if (searchForm.value.status) params.status = searchForm.value.status
    if (searchForm.value.orderNo) params.orderNo = searchForm.value.orderNo
    if (searchForm.value.username) params.username = searchForm.value.username

    const res = await getOrders(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    } else {
      ElMessage.error(res.message || '加载订单列表失败')
    }
  } catch (error) {
    console.error('加载订单列表失败:', error)
    ElMessage.error('加载订单列表失败,请检查后端服务')
  } finally {
    loading.value = false
  }
}

// 取消订单
const handleCancelOrder = async (id) => {
  try {
    const res = await cancelOrder(id)
    return res
  } catch (error) {
    console.error('取消订单失败:', error)
    return { code: 500, message: '取消订单失败' }
  }
}

// 退款订单
const handleRefundOrder = async (id, amount, reason) => {
  try {
    const res = await refundOrder(id, { amount: Math.round(amount * 100), reason })
    return res
  } catch (error) {
    console.error('退款失败:', error)
    return { code: 500, message: '退款失败' }
  }
}

// 获取状态名称
const getStatusName = (status) => {
  const map = {
    pending: '待支付',
    paid: '已支付',
    cancelled: '已取消',
    refunded: '已退款',
    expired: '已过期'
  }
  return map[status] || status
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const map = {
    pending: 'warning',
    paid: 'success',
    cancelled: 'info',
    refunded: 'danger',
    expired: 'info'
  }
  return map[status] || ''
}

// 获取支付方式名称
const getPaymentMethodName = (method) => {
  const map = {
    alipay: '支付宝',
    wechat: '微信支付',
    balance: '余额',
    exchange: '兑换码'
  }
  return map[method] || method
}

// 获取产品信息
const getProductInfo = (row) => {
  if (!row.productData) return '-'
  try {
    const product = JSON.parse(row.productData)
    return product.name || product.level || '-'
  } catch {
    return '-'
  }
}

// 格式化产品数据
const formatProductData = (data) => {
  try {
    const parsed = typeof data === 'string' ? JSON.parse(data) : data
    return JSON.stringify(parsed, null, 2)
  } catch {
    return data
  }
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

// 查看详情
const handleView = async (row) => {
  try {
    const res = await getOrderDetail(row.id)
    if (res.code === 200) {
      currentOrder.value = res.data
    } else {
      currentOrder.value = row
    }
  } catch (error) {
    console.error('获取订单详情失败:', error)
    currentOrder.value = row
  }
  detailDialogVisible.value = true
}

// 取消订单
const handleCancel = (row) => {
  ElMessageBox.confirm(`确定要取消订单 ${row.orderNo} 吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const result = await handleCancelOrder(row.id)
    if (result.code === 200) {
      ElMessage.success('订单已取消')
      loadOrders()
    } else {
      ElMessage.error(result.message || '取消订单失败')
    }
  })
}

// 退款
const handleRefund = (row) => {
  currentOrder.value = row
  refundForm.value.amount = row.paidAmount / 100
  refundForm.value.reason = ''
  refundDialogVisible.value = true
}

// 确认退款
const handleConfirmRefund = () => {
  if (!refundForm.value.reason) {
    ElMessage.warning('请输入退款原因')
    return
  }

  ElMessageBox.confirm(
    `确认退款 ¥${refundForm.value.amount.toFixed(2)} 给订单 ${currentOrder.value.orderNo} 吗?`,
    '退款确认',
    {
      confirmButtonText: '确认退款',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    const result = await handleRefundOrder(
      currentOrder.value.id,
      refundForm.value.amount,
      refundForm.value.reason
    )
    if (result.code === 200) {
      ElMessage.success('退款成功')
      refundDialogVisible.value = false
      loadOrders()
    } else {
      ElMessage.error(result.message || '退款失败')
    }
  })
}

// 搜索
const handleSearch = () => {
  page.value = 1
  loadOrders()
}

// 重置
const handleReset = () => {
  searchForm.value = {
    status: '',
    orderNo: '',
    username: ''
  }
  page.value = 1
  loadOrders()
}

// 分页大小变化
const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  loadOrders()
}

// 页码变化
const handleCurrentChange = (val) => {
  page.value = val
  loadOrders()
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.orders-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
}

.mt-20 {
  margin-top: 20px;
}

pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #606266;
}
</style>
