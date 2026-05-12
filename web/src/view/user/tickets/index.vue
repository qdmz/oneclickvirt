<template>
  <div class="tickets-page">
    <div class="page-header">
      <h2>{{ t('user.tickets.title') }}</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        {{ t('user.tickets.createTicket') }}
      </el-button>
    </div>

    <el-card>
      <el-table :data="tickets" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" :label="t('user.tickets.ticketId')" width="80" />
        <el-table-column prop="title" :label="t('user.tickets.title')" min-width="200">
          <template #default="{ row }">
            <el-link type="primary" @click="viewTicket(row.id)">{{ row.title }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="t('user.tickets.type')" width="120">
          <template #default="{ row }">
            <el-tag>{{ t(`user.tickets.type.${row.type}`) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" :label="t('user.tickets.priority')" width="100">
          <template #default="{ row }">
            <el-tag :type="priorityType(row.priority)">{{ t(`user.tickets.priority.${row.priority}`) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="t('user.tickets.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)">{{ t(`user.tickets.status.${row.status}`) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="t('user.tickets.createdAt')" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('user.tickets.actions')" width="100">
          <template #default="{ row }">
            <el-button size="small" @click="viewTicket(row.id)">{{ t('user.tickets.view') }}</el-button>
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

    <el-dialog v-model="showCreateDialog" :title="t('user.tickets.createTicket')" width="600px">
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item :label="t('user.tickets.title')" prop="title">
          <el-input v-model="createForm.title" :placeholder="t('user.tickets.titlePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('user.tickets.type')" prop="type">
          <el-select v-model="createForm.type" style="width: 100%">
            <el-option v-for="t in ticketTypes" :key="t" :label="t(`user.tickets.type.${t}`)" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('user.tickets.priority')" prop="priority">
          <el-select v-model="createForm.priority" style="width: 100%">
            <el-option v-for="p in ticketPriorities" :key="p" :label="t(`user.tickets.priority.${p}`)" :value="p" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('user.tickets.description')" prop="description">
          <el-input v-model="createForm.description" type="textarea" :rows="4" :placeholder="t('user.tickets.descriptionPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleCreate">{{ t('user.tickets.submit') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getTicketList, createTicket } from '@/api/ticket'

const { t } = useI18n()
const router = useRouter()

const tickets = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const showCreateDialog = ref(false)
const submitting = ref(false)
const createFormRef = ref(null)

const ticketTypes = ['question', 'issue', 'feature', 'complaint', 'other']
const ticketPriorities = ['low', 'medium', 'high', 'urgent']

const createForm = reactive({
  title: '',
  type: 'question',
  priority: 'medium',
  description: ''
})

const createRules = computed(() => ({
  title: [{ required: true, message: t('user.tickets.titleRequired'), trigger: 'blur' }],
  type: [{ required: true, message: t('user.tickets.typeRequired'), trigger: 'change' }],
  priority: [{ required: true, message: t('user.tickets.priorityRequired'), trigger: 'change' }],
  description: [{ required: true, message: t('user.tickets.descriptionRequired'), trigger: 'blur' }]
}))

const loadTickets = async () => {
  loading.value = true
  try {
    const res = await getTicketList({ page: page.value, pageSize: pageSize.value })
    tickets.value = res.data || []
    total.value = res.total || 0
  } catch (e) {
    ElMessage.error(t('user.tickets.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleCreate = async () => {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await createTicket(createForm)
      ElMessage.success(t('user.tickets.createSuccess'))
      showCreateDialog.value = false
      Object.assign(createForm, { title: '', type: 'question', priority: 'medium', description: '' })
      loadTickets()
    } catch (e) {
      ElMessage.error(t('user.tickets.createFailed'))
    } finally {
      submitting.value = false
    }
  })
}

const viewTicket = (id) => {
  router.push(`/user/tickets/${id}`)
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
.tickets-page {
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
</style>
