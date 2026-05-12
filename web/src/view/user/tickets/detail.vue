<template>
  <div class="ticket-detail-page">
    <div class="page-header">
      <el-page-header @back="$router.back()" :content="ticket.title || t('user.tickets.detail')" />
    </div>

    <el-card v-loading="loading" style="margin-top: 20px">
      <template #header>
        <div class="ticket-header">
          <span class="ticket-title">{{ ticket.title }}</span>
          <div class="ticket-meta">
            <el-tag>{{ t(`user.tickets.type.${ticket.type}`) }}</el-tag>
            <el-tag :type="priorityType(ticket.priority)">{{ t(`user.tickets.priority.${ticket.priority}`) }}</el-tag>
            <el-tag :type="statusType(ticket.status)">{{ t(`user.tickets.status.${ticket.status}`) }}</el-tag>
            <span class="ticket-id">#{{ ticket.id }}</span>
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
        <span>{{ t('user.tickets.replies') }}</span>
      </template>

      <div class="replies-list">
        <div v-for="reply in replies" :key="reply.id" class="reply-item">
          <div class="reply-header">
            <el-tag :type="reply.isAdmin ? 'warning' : ''">{{ reply.isAdmin ? t('user.tickets.admin') : t('user.tickets.user') }}</el-tag>
            <span class="reply-time">{{ formatDate(reply.createdAt) }}</span>
          </div>
          <div class="reply-content">{{ reply.content }}</div>
        </div>
        <el-empty v-if="replies.length === 0" :description="t('user.tickets.noReplies')" />
      </div>

      <div class="reply-form" v-if="ticket.status !== 'closed' && ticket.status !== 'resolved'">
        <el-input v-model="replyContent" type="textarea" :rows="3" :placeholder="t('user.tickets.replyPlaceholder')" />
        <el-button type="primary" :loading="replying" @click="handleReply" style="margin-top: 10px">
          {{ t('user.tickets.reply') }}
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
import { getTicketDetail, replyTicket } from '@/api/ticket'

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
    const res = await getTicketDetail(route.params.id)
    ticket.value = res.ticket || {}
    replies.value = res.replies || []
  } catch (e) {
    ElMessage.error(t('user.tickets.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleReply = async () => {
  if (!replyContent.value.trim()) {
    ElMessage.warning(t('user.tickets.replyRequired'))
    return
  }
  replying.value = true
  try {
    await replyTicket(route.params.id, { content: replyContent.value })
    ElMessage.success(t('user.tickets.replySuccess'))
    replyContent.value = ''
    loadDetail()
  } catch (e) {
    ElMessage.error(t('user.tickets.replyFailed'))
  } finally {
    replying.value = false
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
.ticket-detail-page {
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
.ticket-id, .ticket-date {
  color: #909399;
  font-size: 14px;
}
.ticket-description {
  padding: 10px 0;
  white-space: pre-wrap;
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
