<template>
  <div class="admin-wallet">
    <div class="page-header">
      <div>
        <h1>{{ t('admin.wallet.title') }}</h1>
        <p>{{ t('admin.wallet.subtitle') }}</p>
      </div>
      <el-button
        type="primary"
        :icon="Wallet"
        @click="openAdjust()"
      >
        {{ t('admin.wallet.adjust') }}
      </el-button>
    </div>

    <el-card
      shadow="never"
      class="table-card"
    >
      <el-tabs
        v-model="activeTab"
        @tab-change="onTabChange"
      >
        <!-- 余额流水 -->
        <el-tab-pane
          :label="t('admin.wallet.tabRecharge')"
          name="recharge"
        >
          <div class="filter-section">
            <el-form :inline="true">
              <el-form-item :label="t('admin.wallet.filterUser')">
                <el-select
                  v-model="rechargeFilter.userId"
                  filterable
                  clearable
                  remote
                  :remote-method="searchUsers"
                  :loading="userLoading"
                  :placeholder="t('admin.wallet.filterUserPh')"
                  style="width: 220px;"
                  @change="loadRecharge"
                >
                  <el-option
                    v-for="u in userOptions"
                    :key="u.id"
                    :label="`#${u.id} ${u.username}`"
                    :value="u.id"
                  />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button
                  type="primary"
                  @click="loadRecharge"
                >
                  {{ t('common.search') }}
                </el-button>
              </el-form-item>
            </el-form>
          </div>

          <el-table
            v-loading="rechargeLoading"
            :data="rechargeList"
            style="width: 100%;"
          >
            <el-table-column
              prop="id"
              label="#"
              width="70"
            />
            <el-table-column
              prop="userId"
              :label="t('admin.wallet.colUser')"
              width="100"
            >
              <template #default="{ row }">
                #{{ row.userId }}
              </template>
            </el-table-column>
            <el-table-column
              prop="tradeNo"
              :label="t('admin.wallet.colTradeNo')"
              min-width="210"
              show-overflow-tooltip
            />
            <el-table-column
              :label="t('admin.wallet.colAmount')"
              width="130"
            >
              <template #default="{ row }">
                <span :class="row.amount >= 0 ? 'amt-in' : 'amt-out'">
                  {{ row.amount >= 0 ? '+' : '-' }}¥{{ fenToYuan(Math.abs(row.amount)) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.wallet.colMethod')"
              width="110"
            >
              <template #default="{ row }">
                <el-tag
                  size="small"
                  effect="plain"
                >
                  {{ t('admin.wallet.method.' + (row.method || 'epay')) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.wallet.colStatus')"
              width="110"
            >
              <template #default="{ row }">
                <el-tag
                  :type="rechargeStatusType(row.status)"
                  size="small"
                >
                  {{ t('admin.wallet.status.' + row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              prop="remark"
              :label="t('admin.wallet.colRemark')"
              min-width="180"
              show-overflow-tooltip
            >
              <template #default="{ row }">
                {{ row.remark || '-' }}
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.wallet.colTime')"
              width="170"
            >
              <template #default="{ row }">
                {{ formatTime(row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column
              :label="t('common.actions')"
              width="110"
              fixed="right"
            >
              <template #default="{ row }">
                <el-button
                  link
                  type="primary"
                  @click="openAdjust(row.userId)"
                >
                  {{ t('admin.wallet.adjustShort') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination">
            <el-pagination
              v-model:current-page="rechargePage"
              v-model:page-size="rechargePageSize"
              :total="rechargeTotal"
              :page-sizes="[20, 50, 100]"
              layout="total, sizes, prev, pager, next"
              @size-change="loadRecharge"
              @current-change="loadRecharge"
            />
          </div>
        </el-tab-pane>

        <!-- 商品订单 -->
        <el-tab-pane
          :label="t('admin.wallet.tabOrders')"
          name="orders"
        >
          <div class="filter-section">
            <el-form :inline="true">
              <el-form-item :label="t('admin.wallet.filterUser')">
                <el-select
                  v-model="orderFilter.userId"
                  filterable
                  clearable
                  remote
                  :remote-method="searchUsers"
                  :loading="userLoading"
                  :placeholder="t('admin.wallet.filterUserPh')"
                  style="width: 220px;"
                  @change="loadOrders"
                >
                  <el-option
                    v-for="u in userOptions"
                    :key="u.id"
                    :label="`#${u.id} ${u.username}`"
                    :value="u.id"
                  />
                </el-select>
              </el-form-item>
              <el-form-item :label="t('admin.wallet.filterStatus')">
                <el-select
                  v-model="orderFilter.status"
                  clearable
                  style="width: 150px;"
                  @change="loadOrders"
                >
                  <el-option
                    :label="t('admin.wallet.orderStatus.pending')"
                    value="pending"
                  />
                  <el-option
                    :label="t('admin.wallet.orderStatus.paid')"
                    value="paid"
                  />
                  <el-option
                    :label="t('admin.wallet.orderStatus.failed')"
                    value="failed"
                  />
                  <el-option
                    :label="t('admin.wallet.orderStatus.cancelled')"
                    value="cancelled"
                  />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button
                  type="primary"
                  @click="loadOrders"
                >
                  {{ t('common.search') }}
                </el-button>
              </el-form-item>
            </el-form>
          </div>

          <el-table
            v-loading="orderLoading"
            :data="orderList"
            style="width: 100%;"
          >
            <el-table-column
              prop="id"
              label="#"
              width="70"
            />
            <el-table-column
              :label="t('admin.wallet.colUser')"
              width="100"
            >
              <template #default="{ row }">
                #{{ row.userId }}
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.wallet.colProduct')"
              width="100"
            >
              <template #default="{ row }">
                #{{ row.productId }}
              </template>
            </el-table-column>
            <el-table-column
              prop="tradeNo"
              :label="t('admin.wallet.colTradeNo')"
              min-width="200"
              show-overflow-tooltip
            />
            <el-table-column
              :label="t('admin.wallet.colAmount')"
              width="110"
            >
              <template #default="{ row }">
                <span class="amt-out">¥{{ fenToYuan(row.amount) }}</span>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.wallet.colMethod')"
              width="110"
            >
              <template #default="{ row }">
                <el-tag
                  size="small"
                  effect="plain"
                >
                  {{ t('admin.wallet.method.' + (row.method || 'balance')) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.wallet.colStatus')"
              width="110"
            >
              <template #default="{ row }">
                <el-tag
                  :type="orderStatusType(row.status)"
                  size="small"
                >
                  {{ t('admin.wallet.orderStatus.' + row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.wallet.colExpires')"
              width="170"
            >
              <template #default="{ row }">
                {{ formatTime(row.expiresAt) }}
              </template>
            </el-table-column>
            <el-table-column
              prop="remark"
              :label="t('admin.wallet.colRemark')"
              min-width="180"
              show-overflow-tooltip
            >
              <template #default="{ row }">
                {{ row.remark || '-' }}
              </template>
            </el-table-column>
            <el-table-column
              :label="t('admin.wallet.colTime')"
              width="170"
            >
              <template #default="{ row }">
                {{ formatTime(row.createdAt) }}
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination">
            <el-pagination
              v-model:current-page="orderPage"
              v-model:page-size="orderPageSize"
              :total="orderTotal"
              :page-sizes="[20, 50, 100]"
              layout="total, sizes, prev, pager, next"
              @size-change="loadOrders"
              @current-change="loadOrders"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 调整余额 -->
    <el-dialog
      v-model="adjustVisible"
      :title="t('admin.wallet.adjustTitle')"
      width="460px"
    >
      <el-form
        ref="adjustFormRef"
        :model="adjustForm"
        :rules="adjustRules"
        label-position="top"
      >
        <el-form-item
          :label="t('admin.wallet.fUser')"
          prop="userId"
        >
          <el-select
            v-model="adjustForm.userId"
            filterable
            remote
            :remote-method="searchUsers"
            :loading="userLoading"
            :placeholder="t('admin.wallet.fUserPh')"
            style="width: 100%;"
          >
            <el-option
              v-for="u in userOptions"
              :key="u.id"
              :label="`#${u.id} ${u.username}`"
              :value="u.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="t('admin.wallet.fAmount')"
          prop="amount"
        >
          <el-input-number
            v-model="adjustForm.amount"
            :precision="2"
            :step="10"
            style="width: 100%;"
          />
          <div class="form-tip">
            {{ t('admin.wallet.fAmountTip') }}
          </div>
        </el-form-item>
        <el-form-item :label="t('admin.wallet.fRemark')">
          <el-input
            v-model="adjustForm.remark"
            :placeholder="t('admin.wallet.fRemarkPh')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustVisible = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="adjusting"
          @click="submitAdjust"
        >
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Wallet } from '@element-plus/icons-vue'
import {
  adjustUserBalance, getAdminRechargeOrders, getAdminOrders, getUserList
} from '@/api/admin'
import { fenToYuan } from '@/utils/epay'

const { t } = useI18n()

const activeTab = ref('recharge')

const formatTime = (val) => {
  if (!val) return '-'
  const d = new Date(val)
  if (isNaN(d.getTime())) return '-'
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const rechargeStatusType = (s) => ({ paid: 'success', pending: 'warning', failed: 'danger' }[s] || 'info')
const orderStatusType = (s) => ({
  paid: 'success', pending: 'warning', failed: 'danger', cancelled: 'info'
}[s] || 'info')

// 用户下拉
const userOptions = ref([])
const userLoading = ref(false)
const searchUsers = async (keyword) => {
  userLoading.value = true
  try {
    const res = await getUserList({ page: 1, pageSize: 20, keyword: keyword || '' })
    if (res && res.code === 200) {
      userOptions.value = res.data.list || res.data.users || []
    }
  } catch {
    // interceptor handles
  } finally {
    userLoading.value = false
  }
}

// 充值流水
const rechargeLoading = ref(false)
const rechargeList = ref([])
const rechargeTotal = ref(0)
const rechargePage = ref(1)
const rechargePageSize = ref(20)
const rechargeFilter = reactive({ userId: '' })

const loadRecharge = async () => {
  rechargeLoading.value = true
  try {
    const res = await getAdminRechargeOrders({
      page: rechargePage.value,
      pageSize: rechargePageSize.value,
      userId: rechargeFilter.userId || 0
    })
    if (res && res.code === 200) {
      rechargeList.value = res.data.list || []
      rechargeTotal.value = res.data.total || 0
    }
  } catch {
    // interceptor handles
  } finally {
    rechargeLoading.value = false
  }
}

// 商品订单
const orderLoading = ref(false)
const orderList = ref([])
const orderTotal = ref(0)
const orderPage = ref(1)
const orderPageSize = ref(20)
const orderFilter = reactive({ userId: '', status: '' })

const loadOrders = async () => {
  orderLoading.value = true
  try {
    const res = await getAdminOrders({
      page: orderPage.value,
      pageSize: orderPageSize.value,
      userId: orderFilter.userId || 0,
      status: orderFilter.status || ''
    })
    if (res && res.code === 200) {
      orderList.value = res.data.list || []
      orderTotal.value = res.data.total || 0
    }
  } catch {
    // interceptor handles
  } finally {
    orderLoading.value = false
  }
}

const onTabChange = (name) => {
  if (name === 'orders' && orderList.value.length === 0) loadOrders()
  if (name === 'recharge' && rechargeList.value.length === 0) loadRecharge()
}

// 调整余额
const adjustVisible = ref(false)
const adjustFormRef = ref()
const adjusting = ref(false)
const adjustForm = reactive({ userId: '', amount: 10, remark: '' })
const adjustRules = {
  userId: [{ required: true, message: t('admin.wallet.userRequired'), trigger: 'change' }],
  amount: [{ required: true, message: t('admin.wallet.amountRequired'), trigger: 'blur' }]
}

const openAdjust = (userId) => {
  adjustForm.userId = userId || ''
  adjustForm.amount = 10
  adjustForm.remark = ''
  if (userOptions.value.length === 0) searchUsers('')
  adjustVisible.value = true
}

const submitAdjust = async () => {
  if (!adjustFormRef.value) return
  await adjustFormRef.value.validate(async (valid) => {
    if (!valid) return
    adjusting.value = true
    try {
      const res = await adjustUserBalance(adjustForm.userId, {
        amount: adjustForm.amount,
        remark: adjustForm.remark
      })
      if (res && res.code === 200) {
        ElMessage.success(
          t('admin.wallet.adjustSuccess', { balance: fenToYuan(res.data.balance) })
        )
        adjustVisible.value = false
        loadRecharge()
      }
    } finally {
      adjusting.value = false
    }
  })
}

onMounted(() => {
  loadRecharge()
  searchUsers('')
})
</script>

<style scoped>
.admin-wallet {
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
.table-card {
  border: 1px solid var(--el-border-color-lighter);
}
.filter-section {
  margin-bottom: 8px;
}
.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
.amt-in {
  color: var(--el-color-success);
  font-weight: 600;
}
.amt-out {
  color: var(--el-color-danger);
  font-weight: 600;
}
.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.4;
  margin-top: 2px;
}
</style>
