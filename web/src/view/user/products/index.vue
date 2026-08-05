<template>
  <div class="user-products">
    <div class="page-header">
      <div>
        <h1>{{ t('user.products.title') }}</h1>
        <p>{{ t('user.products.subtitle') }}</p>
      </div>
      <div class="header-right">
        <div class="balance-chip">
          <span class="chip-label">{{ t('user.products.myBalance') }}</span>
          <span class="chip-value">¥{{ fenToYuan(balance) }}</span>
        </div>
        <el-button
          :icon="Wallet"
          @click="goWallet"
        >
          {{ t('user.products.gotoRecharge') }}
        </el-button>
      </div>
    </div>

    <!-- 类型筛选 -->
    <div class="filter-section">
      <el-radio-group
        v-model="filterType"
        @change="onFilterChange"
      >
        <el-radio-button value="">
          {{ t('user.products.typeAll') }}
        </el-radio-button>
        <el-radio-button value="vm">
          {{ t('user.products.typeVm') }}
        </el-radio-button>
        <el-radio-button value="lxc">
          {{ t('user.products.typeLxc') }}
        </el-radio-button>
        <el-radio-button value="container">
          {{ t('user.products.typeContainer') }}
        </el-radio-button>
      </el-radio-group>
    </div>

    <!-- 商品网格 -->
    <div
      v-loading="loading"
      class="product-grid"
    >
      <el-empty
        v-if="!loading && list.length === 0"
        :description="t('user.products.empty')"
        class="grid-empty"
      />
      <el-card
        v-for="p in list"
        :key="p.id"
        shadow="hover"
        class="product-card"
      >
        <div class="pc-head">
          <div class="pc-name">
            {{ p.name }}
          </div>
          <el-tag
            size="small"
            :type="typeTagType(p.type)"
            effect="light"
          >
            {{ t('user.products.typeLabel.' + p.type) }}
          </el-tag>
        </div>
        <div
          v-if="p.description"
          class="pc-desc"
        >
          {{ p.description }}
        </div>

        <ul class="pc-specs">
          <li>
            <span class="k">{{ t('user.products.cpu') }}</span>
            <span class="v">{{ p.cpu }} {{ t('user.products.cores') }}</span>
          </li>
          <li>
            <span class="k">{{ t('user.products.memory') }}</span>
            <span class="v">{{ formatMB(p.memoryMB) }}</span>
          </li>
          <li>
            <span class="k">{{ t('user.products.disk') }}</span>
            <span class="v">{{ formatMB(p.diskMB) }}</span>
          </li>
          <li>
            <span class="k">{{ t('user.products.bandwidth') }}</span>
            <span class="v">{{ p.bandwidthMbps }} Mbps</span>
          </li>
        </ul>

        <div class="pc-foot">
          <div class="pc-price">
            <span class="cur">¥</span>{{ fenToYuan(p.price) }}
            <span class="period">/ {{ p.periodDays }}{{ t('user.products.days') }}</span>
          </div>
          <el-button
            type="primary"
            :disabled="isSoldOut(p)"
            @click="openBuy(p)"
          >
            {{ isSoldOut(p) ? t('user.products.soldOut') : t('user.products.buy') }}
          </el-button>
        </div>
        <div
          v-if="p.stock >= 0"
          class="pc-stock"
        >
          {{ t('user.products.stock') }}: {{ Math.max(p.stock - (p.sold || 0), 0) }}
        </div>
      </el-card>
    </div>

    <div class="pagination">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[12, 24, 48]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadList"
        @current-change="loadList"
      />
    </div>

    <!-- 购买对话框 -->
    <el-dialog
      v-model="buyVisible"
      :title="t('user.products.buyTitle')"
      width="500px"
    >
      <div
        v-if="current"
        class="buy-body"
      >
        <div class="buy-summary">
          <div class="bs-row">
            <span>{{ t('user.products.product') }}</span>
            <strong>{{ current.name }}</strong>
          </div>
          <div class="bs-row">
            <span>{{ t('user.products.spec') }}</span>
            <strong>
              {{ current.cpu }}C / {{ formatMB(current.memoryMB) }} / {{ formatMB(current.diskMB) }} / {{ current.bandwidthMbps }}Mbps
            </strong>
          </div>
          <div class="bs-row">
            <span>{{ t('user.products.period') }}</span>
            <strong>{{ current.periodDays }} {{ t('user.products.days') }}</strong>
          </div>
          <div class="bs-row total">
            <span>{{ t('user.products.payable') }}</span>
            <strong class="price">¥{{ fenToYuan(current.price) }}</strong>
          </div>
        </div>

        <el-form label-position="top">
          <el-form-item :label="t('user.products.payMethod')">
            <el-radio-group v-model="buyForm.method">
              <el-radio value="balance">
                {{ t('user.products.payBalance') }}
                <span class="rb-sub">({{ t('user.products.myBalance') }} ¥{{ fenToYuan(balance) }})</span>
              </el-radio>
              <el-radio value="epay">
                {{ t('user.products.payEpay') }}
              </el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>

        <el-alert
          v-if="buyForm.method === 'balance' && balance < current.price"
          :title="t('user.products.insufficient')"
          type="warning"
          :closable="false"
          show-icon
        />
        <el-alert
          v-else
          :title="t('user.products.grantTip')"
          type="info"
          :closable="false"
          show-icon
        />
      </div>
      <template #footer>
        <el-button @click="buyVisible = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="buying"
          @click="submitBuy"
        >
          {{ t('user.products.confirmBuy') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Wallet } from '@element-plus/icons-vue'
import { listProducts, buyProduct, getWallet } from '@/api/user'
import { submitToEpay, currentBaseURL, fenToYuan } from '@/utils/epay'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(12)
const filterType = ref('')
const balance = ref(0)

const typeTagType = (ty) => ({ vm: 'primary', lxc: 'success', container: 'warning' }[ty] || 'info')

const formatMB = (mb) => {
  const n = Number(mb || 0)
  if (n >= 1024) {
    const g = n / 1024
    return `${Number.isInteger(g) ? g : g.toFixed(1)} GB`
  }
  return `${n} MB`
}

const isSoldOut = (p) => p.stock >= 0 && (p.sold || 0) >= p.stock

const goWallet = () => router.push('/user/wallet')

const loadBalance = async () => {
  try {
    const res = await getWallet()
    if (res && res.code === 200) balance.value = res.data.balance || 0
  } catch {
    // interceptor handles
  }
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await listProducts({
      page: page.value,
      pageSize: pageSize.value,
      type: filterType.value
    })
    if (res && res.code === 200) {
      list.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch {
    // interceptor handles
  } finally {
    loading.value = false
  }
}

const onFilterChange = () => {
  page.value = 1
  loadList()
}

// 购买
const buyVisible = ref(false)
const buying = ref(false)
const current = ref(null)
const buyForm = reactive({ method: 'balance' })

const openBuy = (p) => {
  current.value = p
  buyForm.method = balance.value >= p.price ? 'balance' : 'epay'
  buyVisible.value = true
}

const submitBuy = async () => {
  if (!current.value) return
  buying.value = true
  try {
    const res = await buyProduct(current.value.id, {
      method: buyForm.method,
      periodDays: current.value.periodDays,
      baseURL: currentBaseURL()
    })
    if (res && res.code === 200) {
      const { epay, epayUrl, status } = res.data || {}
      if (buyForm.method === 'epay' && epay && epayUrl) {
        buyVisible.value = false
        ElMessage.success(t('user.products.redirecting'))
        submitToEpay(epayUrl, epay)
        return
      }
      buyVisible.value = false
      if (status === 'paid') {
        ElMessage.success(t('user.products.buySuccess'))
      } else {
        ElMessage.success(t('user.products.orderCreated'))
      }
      loadBalance()
      loadList()
    }
  } finally {
    buying.value = false
  }
}

onMounted(() => {
  loadBalance()
  loadList()
})
</script>

<style scoped>
.user-products {
  padding: 20px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
  gap: 12px;
  flex-wrap: wrap;
}
.page-header h1 {
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 4px;
}
.page-header p {
  color: var(--el-text-color-secondary);
  margin: 0;
  font-size: 13px;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}
.balance-chip {
  display: flex;
  align-items: baseline;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 999px;
  background: var(--el-color-primary-light-9);
  border: 1px solid var(--el-color-primary-light-7);
}
.chip-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.chip-value {
  font-size: 15px;
  font-weight: 600;
  color: var(--el-color-primary);
}
.filter-section {
  margin-bottom: 16px;
}
.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(270px, 1fr));
  gap: 16px;
  min-height: 120px;
}
.grid-empty {
  grid-column: 1 / -1;
}
.product-card {
  border: 1px solid var(--el-border-color-lighter);
  transition: border-color 0.2s ease;
}
.product-card:hover {
  border-color: var(--el-color-primary-light-5);
}
.pc-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 6px;
}
.pc-name {
  font-size: 16px;
  font-weight: 600;
}
.pc-desc {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
  margin-bottom: 12px;
  min-height: 18px;
}
.pc-specs {
  list-style: none;
  padding: 0;
  margin: 0 0 14px;
  border-top: 1px solid var(--el-border-color-lighter);
  padding-top: 10px;
}
.pc-specs li {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  padding: 4px 0;
}
.pc-specs .k {
  color: var(--el-text-color-secondary);
}
.pc-specs .v {
  font-weight: 500;
}
.pc-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-top: 1px solid var(--el-border-color-lighter);
  padding-top: 12px;
}
.pc-price {
  font-size: 24px;
  font-weight: 600;
  color: var(--el-color-danger);
  line-height: 1;
}
.pc-price .cur {
  font-size: 14px;
  margin-right: 1px;
}
.pc-price .period {
  font-size: 12px;
  font-weight: 400;
  color: var(--el-text-color-secondary);
  margin-left: 2px;
}
.pc-stock {
  margin-top: 8px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
.buy-summary {
  background: var(--el-fill-color-light);
  border-radius: 8px;
  padding: 12px 14px;
  margin-bottom: 16px;
}
.bs-row {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  padding: 5px 0;
}
.bs-row span {
  color: var(--el-text-color-secondary);
}
.bs-row.total {
  border-top: 1px dashed var(--el-border-color);
  margin-top: 6px;
  padding-top: 10px;
}
.bs-row .price {
  font-size: 18px;
  color: var(--el-color-danger);
}
.rb-sub {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-left: 2px;
}
</style>
