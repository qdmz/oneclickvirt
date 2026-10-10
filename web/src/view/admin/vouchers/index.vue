<template>
  <div class="vouchers-container">
    <!-- 统计卡片 -->
    <el-row
      :gutter="16"
      class="stats-row"
    >
      <el-col
        :xs="12"
        :sm="6"
      >
        <el-card
          shadow="hover"
          class="stat-card"
        >
          <div class="stat-label">
            {{ $t('admin.vouchers.statTotal') }}
          </div>
          <div class="stat-value">
            {{ stats.totalCount }}
          </div>
        </el-card>
      </el-col>
      <el-col
        :xs="12"
        :sm="6"
      >
        <el-card
          shadow="hover"
          class="stat-card"
        >
          <div class="stat-label">
            {{ $t('admin.vouchers.statUnused') }}
          </div>
          <div class="stat-value success">
            {{ stats.unusedCount }}
          </div>
          <div class="stat-sub">
            ¥{{ formatAmount(stats.unusedAmount) }}
          </div>
        </el-card>
      </el-col>
      <el-col
        :xs="12"
        :sm="6"
      >
        <el-card
          shadow="hover"
          class="stat-card"
        >
          <div class="stat-label">
            {{ $t('admin.vouchers.statUsed') }}
          </div>
          <div class="stat-value warning">
            {{ stats.usedCount }}
          </div>
          <div class="stat-sub">
            ¥{{ formatAmount(stats.usedAmount) }}
          </div>
        </el-card>
      </el-col>
      <el-col
        :xs="12"
        :sm="6"
      >
        <el-card
          shadow="hover"
          class="stat-card"
        >
          <div class="stat-label">
            {{ $t('admin.vouchers.statVoid') }}
          </div>
          <div class="stat-value danger">
            {{ stats.voidCount }}
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.vouchers.title') }}</span>
          <div>
            <el-button
              type="danger"
              plain
              :disabled="selectedRows.length === 0"
              @click="handleBatchDelete"
            >
              {{ $t('admin.vouchers.batchDelete') }}
            </el-button>
            <el-button
              type="primary"
              @click="openCreateDialog"
            >
              {{ $t('admin.vouchers.batchCreate') }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-form :inline="true">
          <el-form-item :label="$t('admin.vouchers.code')">
            <el-input
              v-model="filterForm.code"
              :placeholder="$t('admin.vouchers.searchCodePlaceholder')"
              clearable
              style="width: 180px"
              @clear="handleFilterChange"
              @keyup.enter="handleFilterChange"
            />
          </el-form-item>
          <el-form-item :label="$t('common.status')">
            <el-select
              v-model="filterForm.status"
              :placeholder="$t('common.all')"
              clearable
              style="width: 130px"
              @change="handleFilterChange"
            >
              <el-option
                :label="$t('admin.vouchers.statusUnused')"
                :value="0"
              />
              <el-option
                :label="$t('admin.vouchers.statusUsed')"
                :value="1"
              />
              <el-option
                :label="$t('admin.vouchers.statusVoid')"
                :value="2"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('admin.vouchers.batchNo')">
            <el-input
              v-model="filterForm.batchNo"
              :placeholder="$t('admin.vouchers.batchNoPlaceholder')"
              clearable
              style="width: 180px"
              @clear="handleFilterChange"
              @keyup.enter="handleFilterChange"
            />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="handleFilterChange"
            >
              {{ $t('common.search') }}
            </el-button>
            <el-button @click="resetFilter">
              {{ $t('common.reset') }}
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table
        v-loading="loading"
        :data="vouchers"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="46"
          :selectable="row => row.status !== 1"
        />
        <el-table-column
          prop="code"
          :label="$t('admin.vouchers.code')"
          min-width="170"
        >
          <template #default="{ row }">
            <span
              class="code-text"
              @click="copyCode(row.code)"
            >{{ row.code }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('admin.vouchers.amount')"
          width="110"
        >
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.status')"
          width="110"
        >
          <template #default="{ row }">
            <el-tag :type="statusMeta(row).type">
              {{ statusMeta(row).text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="batchNo"
          :label="$t('admin.vouchers.batchNo')"
          min-width="150"
        />
        <el-table-column
          :label="$t('admin.vouchers.usedBy')"
          min-width="120"
        >
          <template #default="{ row }">
            {{ row.usedByUsername || '-' }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('admin.vouchers.usedAt')"
          min-width="160"
        >
          <template #default="{ row }">
            {{ formatDate(row.usedAt) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('admin.vouchers.expireAt')"
          min-width="160"
        >
          <template #default="{ row }">
            {{ row.expireAt ? formatDate(row.expireAt) : $t('admin.vouchers.noExpiry') }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.createdAt')"
          min-width="160"
        >
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="remark"
          :label="$t('admin.vouchers.remark')"
          min-width="120"
          show-overflow-tooltip
        />
        <el-table-column
          :label="$t('common.actions')"
          width="150"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              v-if="row.status === 0"
              link
              type="warning"
              @click="handleVoid(row)"
            >
              {{ $t('admin.vouchers.void') }}
            </el-button>
            <el-button
              v-if="row.status !== 1"
              link
              type="danger"
              @click="handleDelete(row)"
            >
              {{ $t('common.delete') }}
            </el-button>
            <span v-if="row.status === 1">-</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="loadVouchers"
        />
      </div>
    </el-card>

    <!-- 批量生成对话框 -->
    <el-dialog
      v-model="createVisible"
      :title="$t('admin.vouchers.batchCreate')"
      width="520px"
    >
      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-width="110px"
      >
        <el-form-item
          :label="$t('admin.vouchers.amount')"
          prop="amount"
        >
          <el-input-number
            v-model="createForm.amount"
            :min="0.01"
            :precision="2"
            :step="10"
            style="width: 200px"
          />
          <div class="form-tip">
            {{ $t('admin.vouchers.amountTip') }}
          </div>
        </el-form-item>
        <el-form-item
          :label="$t('admin.vouchers.count')"
          prop="count"
        >
          <el-input-number
            v-model="createForm.count"
            :min="1"
            :max="500"
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item :label="$t('admin.vouchers.prefix')">
          <el-input
            v-model="createForm.prefix"
            maxlength="16"
            :placeholder="$t('admin.vouchers.prefixPlaceholder')"
            style="width: 200px"
          />
          <div class="form-tip">
            {{ $t('admin.vouchers.prefixTip') }}
          </div>
        </el-form-item>
        <el-form-item :label="$t('admin.vouchers.expireAt')">
          <el-date-picker
            v-model="createForm.expireAt"
            type="datetime"
            :placeholder="$t('admin.vouchers.expirePlaceholder')"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            style="width: 100%"
          />
          <div class="form-tip">
            {{ $t('admin.vouchers.expireTip') }}
          </div>
        </el-form-item>
        <el-form-item :label="$t('admin.vouchers.remark')">
          <el-input
            v-model="createForm.remark"
            type="textarea"
            :rows="2"
            maxlength="256"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="createLoading"
          @click="submitCreate"
        >
          {{ $t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 生成结果对话框 -->
    <el-dialog
      v-model="resultVisible"
      :title="$t('admin.vouchers.createResult')"
      width="520px"
    >
      <el-alert
        :title="$t('admin.vouchers.createResultTip', { count: createdCodes.length, batchNo: createdBatchNo })"
        type="success"
        :closable="false"
        show-icon
      />
      <el-input
        v-model="createdCodesText"
        type="textarea"
        :rows="10"
        readonly
        class="codes-textarea"
      />
      <template #footer>
        <el-button @click="resultVisible = false">
          {{ $t('common.close') }}
        </el-button>
        <el-button
          type="primary"
          @click="copyCodes"
        >
          {{ $t('admin.vouchers.copyAll') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getVoucherList,
  getVoucherStats,
  createVouchers,
  voidVoucher,
  deleteVoucher,
  batchDeleteVouchers
} from '@/api/voucher'

const { t } = useI18n()

const loading = ref(false)
const vouchers = ref([])
const selectedRows = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const stats = reactive({
  totalCount: 0,
  unusedCount: 0,
  unusedAmount: 0,
  usedCount: 0,
  usedAmount: 0,
  voidCount: 0
})

const filterForm = reactive({
  code: '',
  status: null,
  batchNo: ''
})

const createVisible = ref(false)
const createLoading = ref(false)
const createFormRef = ref(null)
const createForm = reactive({
  amount: 10,
  count: 10,
  prefix: '',
  expireAt: null,
  remark: ''
})
const createRules = {
  amount: [{ required: true, message: t('admin.vouchers.amountRequired'), trigger: 'blur' }],
  count: [{ required: true, message: t('admin.vouchers.countRequired'), trigger: 'blur' }]
}

const resultVisible = ref(false)
const createdCodes = ref([])
const createdBatchNo = ref('')
const createdCodesText = computed(() => createdCodes.value.join('\n'))

const formatAmount = (v) => Number(v || 0).toFixed(2)

const formatDate = (val) => {
  if (!val) return '-'
  const d = new Date(val)
  if (Number.isNaN(d.getTime()) || d.getFullYear() < 1971) return '-'
  return d.toLocaleString('zh-CN', { hour12: false })
}

const statusMeta = (row) => {
  if (row.status === 1) return { text: t('admin.vouchers.statusUsed'), type: 'info' }
  if (row.status === 2) return { text: t('admin.vouchers.statusVoid'), type: 'danger' }
  if (row.expireAt && new Date(row.expireAt).getTime() < Date.now()) {
    return { text: t('admin.vouchers.statusExpired'), type: 'warning' }
  }
  return { text: t('admin.vouchers.statusUnused'), type: 'success' }
}

const loadStats = async () => {
  try {
    const res = await getVoucherStats()
    Object.assign(stats, res.data || {})
  } catch (e) {
    // 统计失败不阻塞列表
  }
}

const loadVouchers = async () => {
  loading.value = true
  try {
    const params = { page: currentPage.value, pageSize: pageSize.value }
    if (filterForm.code) params.code = filterForm.code.trim().toUpperCase()
    if (filterForm.status !== null && filterForm.status !== '') params.status = filterForm.status
    if (filterForm.batchNo) params.batchNo = filterForm.batchNo.trim()
    const res = await getVoucherList(params)
    vouchers.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    ElMessage.error(e.message || t('admin.vouchers.loadFailed'))
  } finally {
    loading.value = false
  }
}

const refresh = () => Promise.all([loadVouchers(), loadStats()])

const handleFilterChange = () => {
  currentPage.value = 1
  loadVouchers()
}

const resetFilter = () => {
  filterForm.code = ''
  filterForm.status = null
  filterForm.batchNo = ''
  handleFilterChange()
}

const handleSizeChange = () => {
  currentPage.value = 1
  loadVouchers()
}

const handleSelectionChange = (rows) => {
  selectedRows.value = rows
}

const openCreateDialog = () => {
  createForm.amount = 10
  createForm.count = 10
  createForm.prefix = ''
  createForm.expireAt = null
  createForm.remark = ''
  createVisible.value = true
}

const submitCreate = async () => {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    createLoading.value = true
    try {
      const payload = {
        amount: Number(createForm.amount),
        count: Number(createForm.count),
        prefix: createForm.prefix ? createForm.prefix.trim().toUpperCase() : '',
        remark: createForm.remark || ''
      }
      if (createForm.expireAt) payload.expireAt = createForm.expireAt
      const res = await createVouchers(payload)
      createdCodes.value = res.data.codes || []
      createdBatchNo.value = res.data.batchNo || ''
      createVisible.value = false
      resultVisible.value = true
      await refresh()
    } catch (e) {
      ElMessage.error(e.message || t('admin.vouchers.createFailed'))
    } finally {
      createLoading.value = false
    }
  })
}

const copyText = async (text, successMsg) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(successMsg)
  } catch (e) {
    const el = document.createElement('textarea')
    el.value = text
    document.body.appendChild(el)
    el.select()
    document.execCommand('copy')
    document.body.removeChild(el)
    ElMessage.success(successMsg)
  }
}

const copyCode = (code) => copyText(code, t('admin.vouchers.copied'))
const copyCodes = () => copyText(createdCodesText.value, t('admin.vouchers.copied'))

const handleVoid = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('admin.vouchers.confirmVoidMsg', { code: row.code }),
      t('admin.vouchers.void'),
      { type: 'warning' }
    )
  } catch (e) {
    return
  }
  try {
    await voidVoucher(row.id)
    ElMessage.success(t('admin.vouchers.voidSuccess'))
    await refresh()
  } catch (e) {
    ElMessage.error(e.message || t('admin.vouchers.voidFailed'))
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('admin.vouchers.confirmDeleteMsg', { code: row.code }),
      t('common.delete'),
      { type: 'warning' }
    )
  } catch (e) {
    return
  }
  try {
    await deleteVoucher(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    await refresh()
  } catch (e) {
    ElMessage.error(e.message || t('common.deleteFailed'))
  }
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) return
  try {
    await ElMessageBox.confirm(
      t('admin.vouchers.confirmBatchDeleteMsg', { count: selectedRows.value.length }),
      t('admin.vouchers.batchDelete'),
      { type: 'warning' }
    )
  } catch (e) {
    return
  }
  try {
    const res = await batchDeleteVouchers({ ids: selectedRows.value.map(r => r.id) })
    ElMessage.success(res.data?.message || t('common.deleteSuccess'))
    selectedRows.value = []
    await refresh()
  } catch (e) {
    ElMessage.error(e.message || t('common.deleteFailed'))
  }
}

onMounted(() => {
  refresh()
})
</script>

<style scoped lang="scss">
.vouchers-container {
  padding: 20px;
}

.stats-row {
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;

  .stat-label {
    font-size: 13px;
    color: var(--el-text-color-secondary);
  }

  .stat-value {
    font-size: 26px;
    font-weight: 600;
    margin-top: 6px;

    &.success { color: var(--el-color-success); }
    &.warning { color: var(--el-color-warning); }
    &.danger { color: var(--el-color-danger); }
  }

  .stat-sub {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-top: 2px;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-bar {
  margin-bottom: 12px;
}

.code-text {
  font-family: 'Courier New', monospace;
  font-weight: 600;
  cursor: pointer;

  &:hover {
    color: var(--el-color-primary);
    text-decoration: underline;
  }
}

.amount-text {
  color: var(--el-color-danger);
  font-weight: 600;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
}

.codes-textarea {
  margin-top: 12px;
  font-family: 'Courier New', monospace;
}
</style>
