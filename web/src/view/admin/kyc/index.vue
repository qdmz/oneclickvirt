<template>
  <div class="admin-kyc-container">
    <!-- Stats -->
    <el-row :gutter="16" style="margin-bottom: 20px">
      <el-col :span="4" v-for="stat in statCards" :key="stat.key">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value" :style="{ color: stat.color }">{{ stats[stat.key] }}</div>
          <div class="stat-label">{{ stat.label }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Filters -->
    <el-card shadow="hover" style="margin-bottom: 16px">
      <el-form inline>
        <el-form-item :label="$t('admin.kyc.status')">
          <el-select v-model="filters.status" clearable @change="fetchRecords" style="width: 150px">
            <el-option :label="$t('admin.kyc.all')" value="" />
            <el-option :label="$t('admin.kyc.pending')" :value="0" />
            <el-option :label="$t('admin.kyc.certified')" :value="1" />
            <el-option :label="$t('admin.kyc.rejected')" :value="2" />
            <el-option :label="$t('admin.kyc.expired')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('admin.kyc.username')">
          <el-input v-model="filters.username" clearable @clear="fetchRecords" @keyup.enter="fetchRecords" style="width: 200px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchRecords">{{ $t('common.search') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Table -->
    <el-card shadow="hover">
      <el-table :data="records" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="userId" label="User ID" width="80" />
        <el-table-column prop="realName" :label="$t('admin.kyc.realName')" width="120" />
        <el-table-column prop="statusText" :label="$t('admin.kyc.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ row.statusText }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="certifiedAt" :label="$t('admin.kyc.certifiedAt')" width="180">
          <template #default="{ row }">
            {{ row.certifiedAt ? new Date(row.certifiedAt).toLocaleString() : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="$t('admin.kyc.createdAt')" width="180">
          <template #default="{ row }">
            {{ new Date(row.createdAt).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" :label="$t('admin.kyc.remark')" show-overflow-tooltip />
        <el-table-column :label="$t('common.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 0 || row.status === 2">
              <el-button size="small" type="success" @click="handleApprove(row)">{{ $t('admin.kyc.approve') }}</el-button>
              <el-button size="small" type="danger" @click="handleReject(row)">{{ $t('admin.kyc.reject') }}</el-button>
            </template>
            <span v-else style="color: #999; font-size: 12px">{{ $t('admin.kyc.noAction') }}</span>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end"
        @size-change="fetchRecords"
        @current-change="fetchRecords"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getKYCRecords, updateKYCStatus, getKYCStats } from '@/api/kyc'

const { t } = useI18n()

const records = ref([])
const loading = ref(false)
const stats = ref({ total: 0, pending: 0, certified: 0, rejected: 0, expired: 0 })

const filters = reactive({ status: '', username: '' })
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const statCards = computed(() => [
  { key: 'total', label: t('admin.kyc.total'), color: '#409eff' },
  { key: 'pending', label: t('admin.kyc.pending'), color: '#e6a23c' },
  { key: 'certified', label: t('admin.kyc.certified'), color: '#67c23a' },
  { key: 'rejected', label: t('admin.kyc.rejected'), color: '#f56c6c' },
  { key: 'expired', label: t('admin.kyc.expired'), color: '#909399' },
])

const statusTagType = (status) => {
  const map = { 0: 'warning', 1: 'success', 2: 'danger', 3: 'info' }
  return map[status] || 'info'
}

const fetchRecords = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize
    }
    if (filters.status !== '') params.status = filters.status
    if (filters.username) params.username = filters.username

    const res = await getKYCRecords(params)
    if (res.data && res.data.code === 0) {
      records.value = res.data.data.list || []
      pagination.total = res.data.data.total || 0
    }
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const res = await getKYCStats()
    if (res.data && res.data.code === 0) {
      stats.value = res.data.data
    }
  } catch {
    // ignore
  }
}

const handleApprove = async (row) => {
  try {
    await ElMessageBox.confirm(t('admin.kyc.approveConfirm'), t('common.confirm'), { type: 'warning' })
    await updateKYCStatus(row.id, { status: 1, remark: 'approved by admin' })
    ElMessage.success(t('admin.kyc.approveSuccess'))
    fetchRecords()
    fetchStats()
  } catch {
    // cancelled
  }
}

const handleReject = async (row) => {
  try {
    const { value } = await ElMessageBox.prompt(t('admin.kyc.rejectReason'), t('admin.kyc.reject'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      inputPlaceholder: t('admin.kyc.rejectReasonPlaceholder'),
      type: 'warning'
    })
    await updateKYCStatus(row.id, { status: 2, remark: value || 'rejected by admin' })
    ElMessage.success(t('admin.kyc.rejectSuccess'))
    fetchRecords()
    fetchStats()
  } catch {
    // cancelled
  }
}

onMounted(() => {
  fetchRecords()
  fetchStats()
})
</script>

<style scoped>
.admin-kyc-container { padding: 0; }
.stat-card { text-align: center; }
.stat-value { font-size: 28px; font-weight: 700; }
.stat-label { font-size: 12px; color: #909399; margin-top: 4px; }
</style>
