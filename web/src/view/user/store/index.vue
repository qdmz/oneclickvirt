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
          <el-tag v-if="product.status !== 1" type="info" size="small">{{ t('user.store.offShelf') }}</el-tag>
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
            <div class="spec-item">
              <el-icon><Box /></el-icon>
              <span>{{ product.stock < 0 ? (t('user.store.stock') + '：' + t('user.store.stockUnlimited')) : (t('user.store.stock') + '：' + product.stock) }}</span>
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
