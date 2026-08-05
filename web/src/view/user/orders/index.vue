<template>
  <div class="user-orders">
    <div class="page-header">
      <div>
        <h1>{{ t('user.orders.title') }}</h1>
        <p>{{ t('user.orders.subtitle') }}</p>
      </div>
      <div class="header-right">
        <el-button
          :icon="Goods"
          @click="goProducts"
        >
          {{ t('user.orders.gotoStore') }}
        </el-button>
        <el-button
          :icon="Refresh"
          @click="loadList"
        >
          {{ t('common.refresh') }}
        </el-button>
      </div>
    </div>

    <el-card
      shadow="never"
      class="table-card"
    >
      <el-table
        v-loading="loading"
        :data="list"
        style="width: 100%;"
      >
        <el-table-column
          prop="id"
          label="#"
          width="70"
        />
        <el-table-column
          prop="tradeNo"
          :label="t('user.orders.colTradeNo')"
          min-width="210"
          show-overflow-tooltip
        />
        <el-table-column
          :label="t('user.orders.colAmount')"
          width="110"
        >
          <template #default="{ row }">
            <span class="amt">¥{{ fenToYuan(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('user.orders.colMethod')"
          width="110"
        >
          <template #default="{ row }">
            <el-tag
              size="small"
              effect="plain"
            >
              {{ t('user.orders.method.' + (row.method || 'balance')) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('user.orders.colStatus')"
          width="110"
        >
          <template #default="{ row }">
            <el-tag
              :type="statusTagType(row.status)"
              size="small"
            >
              {{ t('user.orders.status.' + row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('user.orders.colPeriod')"
          width="100"
        >
          <template #default="{ row }">
            {{ row.periodDays }} {{ t('user.orders.days') }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('user.orders.colExpires')"
          width="170"
        >
          <template #default="{ row }">
            <span :class="expiringSoon(row.expiresAt) ? 'exp-warn' : ''">
              {{ formatTime(row.expiresAt) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          prop="remark"
          :label="t('user.orders.colRemark')"
          min-width="200"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ row.remark || '-' }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('user.orders.colTime')"
          width="170"
        >
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('common.actions')"
          width="130"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'paid'"
              link
              type="primary"
              @click="goInstances"
            >
              {{ t('user.orders.createInstance') }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { Refresh, Goods } from '@element-plus/icons-vue'
import { listMyOrders } from '@/api/user'
import { fenToYuan } from '@/utils/epay'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const statusTagType = (s) => ({
  paid: 'success', pending: 'warning', failed: 'danger', cancelled: 'info'
}[s] || 'info')

const formatTime = (val) => {
  if (!val) return '-'
  const d = new Date(val)
  if (isNaN(d.getTime())) return '-'
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

// 7 天内到期高亮
const expiringSoon = (val) => {
  if (!val) return false
  const d = new Date(val).getTime()
  if (isNaN(d)) return false
  const diff = d - Date.now()
  return diff > 0 && diff < 7 * 24 * 3600 * 1000
}

const goProducts = () => router.push('/user/products')
const goInstances = () => router.push('/user/instances')

const loadList = async () => {
  loading.value = true
  try {
    const res = await listMyOrders({ page: page.value, pageSize: pageSize.value })
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

onMounted(loadList)
</script>

<style scoped>
.user-orders {
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
  gap: 10px;
}
.table-card {
  border: 1px solid var(--el-border-color-lighter);
}
.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
.amt {
  font-weight: 600;
  color: var(--el-color-danger);
}
.exp-warn {
  color: var(--el-color-warning);
  font-weight: 600;
}
</style>
