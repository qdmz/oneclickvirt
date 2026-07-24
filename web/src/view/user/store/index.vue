<template>
  <div class="store-container">
    <!-- 页面头部 -->
    <div class="store-header">
      <h1>{{ t('user.store.title') }}</h1>
      <p>{{ t('user.store.subtitle') }}</p>
    </div>

    <!-- 产品分类筛选 -->
    <div class="category-filter">
      <el-radio-group v-model="selectedCategory" @change="handleCategoryChange">
        <el-radio-button label="">{{ t('user.store.allCategories') }}</el-radio-button>
        <el-radio-button label="vm">{{ t('user.store.categoryVM') }}</el-radio-button>
        <el-radio-button label="container">{{ t('user.store.categoryContainer') }}</el-radio-button>
        <el-radio-button label="gpu">{{ t('user.store.categoryGPU') }}</el-radio-button>
      </el-radio-group>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-loading-directive />
      <div class="loading-text">{{ t('common.loading') }}</div>
    </div>

    <!-- 产品列表 -->
    <div v-else-if="productList.length > 0" class="product-grid">
      <el-card
        v-for="product in productList"
        :key="product.id"
        class="product-card"
        shadow="hover"
        @click="goToDetail(product.id)"
      >
        <!-- 产品状态标签 -->
        <div class="product-badge">
          <el-tag v-if="product.is_new" type="success" size="small">{{ t('user.store.newProduct') }}</el-tag>
          <el-tag v-if="product.is_hot" type="danger" size="small">{{ t('user.store.hotProduct') }}</el-tag>
          <el-tag v-if="!product.is_active" type="info" size="small">{{ t('user.store.offShelf') }}</el-tag>
        </div>

        <!-- 产品图标/图片 -->
        <div class="product-icon">
          <el-icon :size="48" color="#16a34a">
            <component :is="getProductIcon(product.type)" />
          </el-icon>
        </div>

        <!-- 产品信息 -->
        <div class="product-info">
          <h3 class="product-name">{{ product.name }}</h3>
          <p class="product-desc">{{ product.description }}</p>

          <!-- 资源配置 -->
          <div class="product-specs">
            <div class="spec-item">
              <el-icon><Cpu /></el-icon>
              <span>{{ product.cpu }} {{ t('user.store.cores') }}</span>
            </div>
            <div class="spec-item">
              <el-icon><Memo /></el-icon>
              <span>{{ formatMemory(product.memory) }}</span>
            </div>
            <div class="spec-item">
              <el-icon><Coin /></el-icon>
              <span>{{ formatDisk(product.disk) }}</span>
            </div>
            <div class="spec-item">
              <el-icon><TopRight /></el-icon>
              <span>{{ formatBandwidth(product.bandwidth) }}</span>
            </div>
            <div v-if="product.traffic > 0" class="spec-item">
              <el-icon><DataLine /></el-icon>
              <span>{{ formatTraffic(product.traffic) }}</span>
            </div>
            <div v-else class="spec-item">
              <el-icon><DataLine /></el-icon>
              <span>{{ t('user.store.unlimitedTraffic') }}</span>
            </div>
          </div>
        </div>

        <!-- 价格与操作 -->
        <div class="product-footer">
          <div class="product-price">
            <span class="price-symbol">¥</span>
            <span class="price-value">{{ product.price }}</span>
            <span class="price-unit">/{{ t('user.store.perMonth') }}</span>
          </div>
          <el-button type="primary" size="default" @click.stop="goToDetail(product.id)">
            {{ t('user.store.buyNow') }}
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 空状态 -->
    <el-empty v-else :description="t('user.store.noProducts')" />

    <!-- 分页 -->
    <div v-if="total > pageSize" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[12, 24, 48]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Cpu, Memo, Coin, TopRight, DataLine, Monitor, Box, Grid } from '@element-plus/icons-vue'
import { getProductList } from '@/api/product'
import { formatMemorySize, formatDiskSize, formatBandwidthSpeed } from '@/utils/unit-formatter'

const router = useRouter()
const { t } = useI18n()

const loading = ref(true)
const productList = ref([])
const selectedCategory = ref('')
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

// 获取产品图标
const getProductIcon = (type) => {
  const iconMap = {
    vm: Monitor,
    container: Box,
    gpu: Grid
  }
  return iconMap[type] || Monitor
}

// 格式化内存
const formatMemory = (memory) => formatMemorySize(memory)

// 格式化磁盘
const formatDisk = (disk) => formatDiskSize(disk)

// 格式化带宽
const formatBandwidth = (bandwidth) => formatBandwidthSpeed(bandwidth)

// 格式化流量
const formatTraffic = (traffic) => formatDiskSize(traffic)

// 加载产品列表
const loadProducts = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      category: selectedCategory.value || undefined
    }
    const res = await getProductList(params)
    if (res.code === 200) {
      productList.value = res.data?.list || res.data?.items || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载产品列表失败:', error)
    ElMessage.error(error?.message || t('user.store.loadFailed'))
  } finally {
    loading.value = false
  }
}

// 分类切换
const handleCategoryChange = () => {
  currentPage.value = 1
  loadProducts()
}

// 分页切换
const handlePageChange = (page) => {
  currentPage.value = page
  loadProducts()
}

// 每页条数切换
const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadProducts()
}

// 跳转到详情页
const goToDetail = (id) => {
  router.push(`/user/store/${id}`)
}

onMounted(() => {
  loadProducts()
})
</script>

<style lang="scss" scoped>
.store-container {
  padding: 24px;
}

.store-header {
  margin-bottom: 24px;

  h1 {
    margin: 0 0 8px 0;
    color: var(--text-color-primary);
    font-size: 28px;
    font-weight: 600;
  }

  p {
    margin: 0;
    color: var(--text-color-secondary);
    font-size: 16px;
  }
}

.category-filter {
  margin-bottom: 24px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  color: #666;
}

.loading-text {
  margin-top: 16px;
  font-size: 14px;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.product-card {
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  }

  :deep(.el-card__body) {
    padding: 20px;
    display: flex;
    flex-direction: column;
    height: 100%;
  }
}

.product-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  display: flex;
  gap: 6px;
  z-index: 1;
}

.product-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #dcfce7 0%, #bbf7d0 100%);
  border-radius: 12px;
  margin-bottom: 16px;
}

.product-info {
  flex: 1;
}

.product-name {
  margin: 0 0 8px 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color-primary);
}

.product-desc {
  margin: 0 0 16px 0;
  font-size: 13px;
  color: var(--text-color-secondary);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-specs {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
  margin-bottom: 16px;
}

.spec-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-color-secondary);

  .el-icon {
    color: #16a34a;
    font-size: 14px;
  }
}

.product-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
}

.product-price {
  display: flex;
  align-items: baseline;
  gap: 2px;

  .price-symbol {
    font-size: 14px;
    color: #f56c6c;
    font-weight: 600;
  }

  .price-value {
    font-size: 24px;
    color: #f56c6c;
    font-weight: 700;
  }

  .price-unit {
    font-size: 12px;
    color: var(--text-color-secondary);
  }
}

.pagination-wrapper {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .store-container {
    padding: 16px;
  }

  .store-header h1 {
    font-size: 22px;
  }

  .product-grid {
    grid-template-columns: 1fr;
  }
}
</style>
