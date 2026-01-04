<template>
  <div class="orders-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>我的订单</span>
          <div>
            <el-select
              v-model="searchForm.status"
              placeholder="订单状态"
              clearable
              @change="handleSearch"
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
          label="订单金额"
          width="120"
        >
          <template #default="{ row }">
            ¥{{ (row.amount / 100).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column
          label="订单状态"
          width="100"
        >
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          label="支付方式"
          width="100"
        >
          <template #default="{ row }">
            {{ getPaymentMethodName(row.paymentMethod) }}
          </template>
        </el-table-column>
        <el-table-column
          label="产品信息"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ getProductInfo(row) }}
          </template>
        </el-table-column>
        <el-table-column
          label="创建时间"
          width="180"
        >
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column
          label="支付时间"
          width="180"
        >
          <template #default="{ row }">
            {{ row.paymentTime ? formatTime(row.paymentTime) : '-' }}
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          width="150"
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
              v-if="row.status === 'pending'"
              type="warning"
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
      width="600px"
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
        <el-descriptions-item label="订单金额">
          ¥{{ (currentOrder.amount / 100).toFixed(2) }}
        </el-descriptions-item>
        <el-descriptions-item label="实付金额">
          ¥{{ (currentOrder.paidAmount / 100).toFixed(2) }}
        </el-descriptions-item>
        <el-descriptions-item label="订单状态">
          <el-tag :type="getStatusTagType(currentOrder.status)">
            {{ getStatusName(currentOrder.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="支付方式">
          {{ getPaymentMethodName(currentOrder.paymentMethod) }}
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
          <pre>{{ JSON.stringify(JSON.parse(currentOrder.productData), null, 2) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserOrders, getUserOrderDetail, cancelOrder } from '@/api/user-payment'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const tableData = ref([])
const detailDialogVisible = ref(false)
const currentOrder = ref(null)

const searchForm = ref({
  status: ''
})

// 加载订单列表
const loadOrders = async () => {
  loading.value = true
  try {
    const res = await getUserOrders({
      page: page.value,
      pageSize: pageSize.value,
      ...searchForm.value
    })
    console.log('订单列表API响应:', res)
    if (res.code === 200 || res.code === 0) {
      tableData.value = res.data?.list || res.data || []
      total.value = res.data?.total || 0
    } else {
      ElMessage.error(res.message || '加载订单列表失败')
    }
  } catch (error) {
    console.error('加载订单列表失败:', error)
    ElMessage.error('加载订单列表失败')
  } finally {
    loading.value = false
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
    return `${product.name} - ${product.level}级`
  } catch {
    return '-'
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
    const res = await getUserOrderDetail(row.id)
    console.log('订单详情API响应:', res)
    if (res.code === 200 || res.code === 0) {
      currentOrder.value = res.data
      detailDialogVisible.value = true
    } else {
      ElMessage.error(res.message || '加载订单详情失败')
    }
  } catch (error) {
    console.error('加载订单详情失败:', error)
    ElMessage.error('加载订单详情失败')
  }
}

// 取消订单
const handleCancel = (row) => {
  ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await cancelOrder(row.id)
      console.log('取消订单API响应:', res)
      if (res.code === 200 || res.code === 0) {
        ElMessage.success('订单已取消')
        loadOrders()
      } else {
        ElMessage.error(res.message || '取消订单失败')
      }
    } catch (error) {
      console.error('取消订单失败:', error)
      ElMessage.error('取消订单失败')
    }
  })
}

// 搜索
const handleSearch = () => {
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

.mt-20 {
  margin-top: 20px;
}
</style>
