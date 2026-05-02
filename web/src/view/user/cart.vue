<!-- 购物车页面 -->
<template>
  <div class="cart-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>购物车</span>
          <div class="header-actions">
            <el-button
              v-if="cartItems.length > 0"
              type="danger"
              @click="handleClearAll"
            >
              清空购物车
            </el-button>
          </div>
        </div>
      </template>

      <div v-loading="loading" class="cart-content">
        <!-- 购物车为空 -->
        <div v-if="cartItems.length === 0" class="empty-cart">
          <el-empty description="您的购物车是空的">
            <el-button type="primary" @click="goToProducts">
              去选购商品
            </el-button>
          </el-empty>
        </div>

        <!-- 购物车列表 -->
        <div v-else>
          <el-table :data="cartItems" border stripe>
            <el-table-column
              prop="productId"
              label="商品ID"
              width="100"
            />

            <el-table-column
              label="商品信息"
            >
              <template #default="{ row }">
                <div>
                  <div class="product-name">{{ row.name }}</div>
                  <el-tag size="small" type="info">{{ row.level }}级</el-tag>
                </div>
              </template>
            </el-table-column>

            <el-table-column
              label="价格"
              width="120"
            >
              <template #default="{ row }">
                ¥{{ (row.price / 100).toFixed(2) }}
              </template>
            </el-table-column>

            <el-table-column
              label="数量"
              width="180"
            >
              <template #default="{ row }">
                <el-input-number
                  v-model="row.quantity"
                  :min="1"
                  :max="999"
                  size="small"
                  @change="handleQuantityChange(row)"
                />
              </template>
            </el-table-column>

            <el-table-column
              label="小计"
              width="120"
            >
              <template #default="{ row }">
                ¥{{ ((row.price * row.quantity) / 100).toFixed(2) }}
              </template>
            </el-table-column>

            <el-table-column
              label="操作"
              width="150"
              fixed="right"
            >
              <template #default="{ row }">
                <el-button
                  type="danger"
                  size="small"
                  @click="handleRemove(row)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <!-- 购物车底部 -->
          <div class="cart-footer">
            <div class="total-info">
              <span>商品数量：</span>
              <span class="count">{{ totalQuantity }}</span>
            </div>
            <div class="total-price">
              <span>总计：</span>
              <span class="price">¥{{ (cartTotal / 100).toFixed(2) }}</span>
            </div>
            <el-button
              type="primary"
              size="large"
              @click="handleCheckout"
              :disabled="cartItems.length === 0"
            >
              立即结算
            </el-button>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getCartItems, removeFromCart, clearCart } from '@/api/user-cart'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const cartItems = ref([])
const cartTotal = ref(0)

// 加载购物车数据
const loadCart = async () => {
  try {
    loading.value = true
    const res = await getCartItems()
    cartItems.value = res.data || []
  } catch (error) {
    console.error('加载购物车失败:', error)
    ElMessage.error('加载购物车失败')
  } finally {
    loading.value = false
  }
}

// 计算商品总数
const totalQuantity = computed(() => {
  return cartItems.value.reduce((sum, item) => sum + item.quantity, 0)
})

// 数量变更
const handleQuantityChange = async (item) => {
  try {
    await updateCartItem(item.itemId, item.quantity)
    ElMessage.success('数量已更新')
    loadCart()
  } catch (error) {
    console.error('更新数量失败:', error)
    ElMessage.error('更新失败，请重试')
    loadCart() // 重新加载恢复原数值
  }
}

// 删除商品
const handleRemove = async (item) => {
  try {
    await removeFromCart(item.itemId)
    ElMessage.success('已删除商品')
    loadCart()
  } catch (error) {
    console.error('删除失败:', error)
    ElMessage.error('删除失败，请重试')
  }
}

// 清空购物车
const handleClearAll = async () => {
  try {
    await clearCart()
    ElMessage.success('购物车已清空')
    loadCart()
  } catch (error) {
    console.error('清空失败:', error)
    ElMessage.error('清空失败，请重试')
  }
}

// 结算
const handleCheckout = () => {
  ElMessageBox.confirm(
    '确定要结算吗？',
    '结算确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    router.push({
      path: '/checkout',
      query: { cart: 1 }
    })
  }).catch(() => {})
}

// 去选购商品
const goToProducts = () => {
  router.push({ path: '/product-server' })
}

onMounted(() => {
  loadCart()
})
</script>

<style scoped>
.cart-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.cart-content {
  min-height: 400px;
}

.empty-cart {
  padding: 100px 0;
}

.cart-footer {
  margin-top: 20px;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.total-info {
  font-size: 14px;
  color: #606266;
}

.total-info .count {
  font-weight: bold;
  margin-left: 5px;
}

.total-price {
  font-size: 16px;
  font-weight: bold;
}

.total-price .price {
  color: #f56c6c;
  font-size: 20px;
  margin-left: 10px;
}
</style>
