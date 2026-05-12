<template>
  <div class="admin-ticket-detail-page">
    <div class="page-header">
      <el-page-header @back="$router.back()" :content="ticket.title || t('admin.tickets.detail')" />
    </div>

    <el-card v-loading="loading" style="margin-top: 20px">
      <template #header>
        <div class="ticket-header">
          <span class="ticket-title">{{ ticket.title }}</span>
          <div class="ticket-meta">
            <el-tag>{{ t(`admin.tickets.type.${ticket.type}`) }}</el-tag>
            <el-tag :type="priorityType(ticket.priority)">{{ t(`admin.tickets.priority.${ticket.priority}`) }}</el-tag>
            <el-tag :type="statusType(ticket.status)">{{ t(`admin.tickets.status.${ticket.status}`) }}</el-tag>
            <span class="ticket-id">#{{ ticket.id }}</span>
            <span class="ticket-user">UID: {{ ticket.userId }}</span>
            <span class="ticket-date">{{ formatDate(ticket.createdAt) }}</span>
          </div>
        </div>
      </template>

      <div class="ticket-description">
        <p>{{ ticket.description }}</p>
      </div>
    </el-card>

    <el-card style="margin-top: 20px">
      <template #header>
        <div class="reply-header">
          <span>{{ t('admin.tickets.replies') }}</span>
          <el-dropdown @command="handleStatusChange" v-if="ticket.status !== 'closed' && ticket.status !== 'resolved'">
            <el-button size="small">{{ t('admin.tickets.updateStatus') }}<el-icon><ArrowDown /></el-icon></el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="pending">{{ t('admin.tickets.status.pending') }}</el-dropdown-item>
                <el-dropdown-item command="resolved">{{ t('admin.tickets.status.resolved') }}</el-dropdown-item>
                <el-dropdown-item command="closed">{{ t('admin.tickets.status.closed') }}</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </template>

      <div class="replies-list">
        <div v-for="reply in replies" :key="reply.id" class="reply-item">
          <div class="reply-header">
            <el-tag :type="reply.isAdmin ? 'warning' : 'success'">{{ reply.isAdmin ? t('admin.tickets.admin') : t('admin.tickets.user') }}</el-tag>
            <span class="reply-time">{{ formatDate(reply.createdAt) }}</span>
          </div>
          <div class="reply-content">{{ reply.content }}</div>
        </div>
        <el-empty v-if="replies.length === 0" :description="t('admin.tickets.noReplies')" />
      </div>

      <div class="reply-form" v-if="ticket.status !== 'closed' && ticket.status !== 'resolved'">
        <el-input v-model="replyContent" type="textarea" :rows="3" :placeholder="t('admin.tickets.replyPlaceholder')" />
        <el-button type="primary" :loading="replying" @click="handleReply" style="margin-top: 10px">
          {{ t('admin.tickets.reply') }}
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { getAdminTicketDetail, adminReplyTicket, updateTicket, closeTicket } from '@/api/ticket'

const { t } = useI18n()
const route = useRoute()

const ticket = ref({})
const replies = ref([])
const loading = ref(false)
const replyContent = ref('')
const replying = ref(false)

const loadDetail = async () => {
  loading.value = true
  try {
    const res = await getAdminTicketDetail(route.params.id)
    ticket.value = res.ticket || {}
    replies.value = res.replies || []
  } catch (e) {
    ElMessage.error(t('admin.tickets.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleReply = async () => {
  if (!replyContent.value.trim()) {
    ElMessage.warning(t('admin.tickets.replyRequired'))
    return
  }
  replying.value = true
  try {
    await adminReplyTicket(route.params.id, { content: replyContent.value })
    ElMessage.success(t('admin.tickets.replySuccess'))
    replyContent.value = ''
    loadDetail()
  } catch (e) {
    ElMessage.error(t('admin.tickets.replyFailed'))
  } finally {
    replying.value = false
  }
}

const handleStatusChange = async (status) => {
  try {
    if (status === 'closed') {
      await closeTicket(route.params.id, { resolutionNotes: '' })
    } else {
      await updateTicket(route.params.id, { status })
    }
    ElMessage.success(t('admin.tickets.statusUpdated'))
    loadDetail()
  } catch (e) {
    ElMessage.error(t('admin.tickets.updateFailed'))
  }
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
  loadDetail()
})
</script>

<style scoped>
.admin-ticket-detail-page {
  padding: 20px;
}
.ticket-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}
.ticket-title {
  font-size: 18px;
  font-weight: bold;
}
.ticket-meta {
  display: flex;
  gap: 8px;
  align-items: center;
}
.ticket-id, .ticket-user, .ticket-date {
  color: #909399;
  font-size: 14px;
}
.ticket-description {
  padding: 10px 0;
  white-space: pre-wrap;
}
.reply-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.reply-item {
  padding: 15px 0;
  border-bottom: 1px solid #ebeef5;
}
.reply-item:last-child {
  border-bottom: none;
}
.reply-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}
.reply-time {
  color: #909399;
  font-size: 14px;
}
.reply-content {
  white-space: pre-wrap;
  line-height: 1.6;
}
.reply-form {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}
</style>
