<template>
  <div class="admin-tickets-page">
    <div class="page-header">
      <h2>{{ t('admin.tickets.title') }}</h2>
    </div>

    <el-card>
      <div class="filter-bar">
        <el-select v-model="statusFilter" :placeholder="t('admin.tickets.filterStatus')" clearable @change="loadTickets">
          <el-option :label="t('admin.tickets.status.all')" value="" />
          <el-option :label="t('admin.tickets.status.open')" value="open" />
          <el-option :label="t('admin.tickets.status.pending')" value="pending" />
          <el-option :label="t('admin.tickets.status.resolved')" value="resolved" />
          <el-option :label="t('admin.tickets.status.closed')" value="closed" />
        </el-select>
      </div>

      <el-table :data="tickets" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" :label="t('admin.tickets.ticketId')" width="80" />
        <el-table-column prop="userId" :label="t('admin.tickets.userId')" width="80" />
        <el-table-column prop="title" :label="t('admin.tickets.title')" min-width="200">
          <template #default="{ row }">
            <el-link type="primary" @click="viewTicket(row.id)">{{ row.title }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="t('admin.tickets.type')" width="120">
          <template #default="{ row }">
            <el-tag>{{ t(`admin.tickets.type.${row.type}`) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" :label="t('admin.tickets.priority')" width="100">
          <template #default="{ row }">
            <el-tag :type="priorityType(row.priority)">{{ t(`admin.tickets.priority.${row.priority}`) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="t('admin.tickets.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)">{{ t(`admin.tickets.status.${row.status}`) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="t('admin.tickets.createdAt')" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.tickets.actions')" width="100">
          <template #default="{ row }">
            <el-button size="small" @click="viewTicket(row.id)">{{ t('admin.tickets.view') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > 0"
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getAdminTicketList } from '@/api/ticket'

const { t } = useI18n()
const router = useRouter()

const tickets = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const statusFilter = ref('')

const loadTickets = async () => {
  loading.value = true
  try {
    const res = await getAdminTicketList({ page: page.value, pageSize: pageSize.value, status: statusFilter.value })
    tickets.value = res.data || []
    total.value = res.total || 0
  } catch (e) {
    ElMessage.error(t('admin.tickets.loadFailed'))
  } finally {
    loading.value = false
  }
}

const viewTicket = (id) => {
  router.push(`/admin/tickets/${id}`)
}

const handlePageChange = (p) => {
  page.value = p
  loadTickets()
}

const priorityType = (p) => {
  const map = { low: '', medium: 'warning', high: 'danger', urgent: 'danger' }
  return map[p] || ''
}

const statusType = (s) => {
  const map = { open: 'success', pending: 'warning', resolved: 'info', closed: 'info' }
  return map[s] || ''
}

const formatDate = (d) => {
  if (!d) return '-'
  return new Date(d).toLocaleString()
}

onMounted(() => {
  loadTickets()
})
</script>

<style scoped>
.admin-tickets-page {
  padding: 20px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.page-header h2 {
  margin: 0;
}
.filter-bar {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}
</style>
