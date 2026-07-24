<template>
  <div class="tickets-container">
    <!-- 页面头部 -->
    <div class="tickets-header">
      <h1>{{ t('user.tickets.title') }}</h1>
      <el-button type="primary" @click="createVisible = true">
        <el-icon><Plus /></el-icon>
        {{ t('user.tickets.createTicket') }}
      </el-button>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-radio-group v-model="filterStatus" @change="handleFilterChange">
        <el-radio-button label="">{{ t('user.tickets.all') }}</el-radio-button>
        <el-radio-button label="open">{{ t('user.tickets.open') }}</el-radio-button>
        <el-radio-button label="replied">{{ t('user.tickets.replied') }}</el-radio-button>
        <el-radio-button label="closed">{{ t('user.tickets.closed') }}</el-radio-button>
      </el-radio-group>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-loading-directive />
      <div class="loading-text">{{ t('common.loading') }}</div>
    </div>

    <!-- 工单列表 -->
    <template v-else>
      <el-empty v-if="ticketList.length === 0" :description="t('user.tickets.noTickets')" />

      <div v-else class="ticket-list">
        <el-card
          v-for="ticket in ticketList"
          :key="ticket.id"
          class="ticket-card"
          shadow="hover"
          @click="openTicketDetail(ticket)"
        >
          <div class="ticket-header">
            <div class="ticket-info">
              <span class="ticket-id">#{{ ticket.id }}</span>
              <el-tag :type="getStatusType(ticket.status)" size="small">
                {{ getStatusText(ticket.status) }}
              </el-tag>
              <el-tag type="info" size="small">{{ ticket.category }}</el-tag>
            </div>
            <span class="ticket-time">{{ formatDate(ticket.updated_at || ticket.created_at) }}</span>
          </div>
          <div class="ticket-subject">{{ ticket.subject }}</div>
          <div class="ticket-preview">{{ ticket.last_message || ticket.content }}</div>
        </el-card>
      </div>

      <!-- 分页 -->
      <div v-if="total > pageSize" class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </template>

    <!-- 创建工单弹窗 -->
    <el-dialog
      v-model="createVisible"
      :title="t('user.tickets.createTicket')"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-width="80px"
      >
        <el-form-item :label="t('user.tickets.category')" prop="category">
          <el-select v-model="createForm.category" :placeholder="t('user.tickets.selectCategory')">
            <el-option :label="t('user.tickets.catGeneral')" value="general" />
            <el-option :label="t('user.tickets.catTechnical')" value="technical" />
            <el-option :label="t('user.tickets.catBilling')" value="billing" />
            <el-option :label="t('user.tickets.catOther')" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('user.tickets.subject')" prop="subject">
          <el-input v-model="createForm.subject" :placeholder="t('user.tickets.inputSubject')" />
        </el-form-item>
        <el-form-item :label="t('user.tickets.content')" prop="content">
          <el-input
            v-model="createForm.content"
            type="textarea"
            :rows="5"
            :placeholder="t('user.tickets.inputContent')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleCreate">
          {{ t('common.submit') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 工单详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      :title="t('user.tickets.ticketDetail')"
      width="700px"
      class="ticket-detail-dialog"
    >
      <div v-if="currentTicket" class="ticket-detail">
        <!-- 工单信息 -->
        <div class="ticket-meta">
          <div class="meta-row">
            <span class="meta-label">{{ t('user.tickets.subject') }}:</span>
            <span class="meta-value">{{ currentTicket.subject }}</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">{{ t('user.tickets.category') }}:</span>
            <el-tag size="small">{{ currentTicket.category }}</el-tag>
            <el-tag :type="getStatusType(currentTicket.status)" size="small" style="margin-left: 8px;">
              {{ getStatusText(currentTicket.status) }}
            </el-tag>
          </div>
        </div>

        <el-divider />

        <!-- 对话列表 -->
        <el-scrollbar class="message-list" max-height="400px">
          <div
            v-for="(msg, index) in messages"
            :key="index"
            class="message-item"
            :class="{ 'message-self': msg.is_user }"
          >
            <div class="message-avatar">
              <el-avatar :size="36" :icon="msg.is_user ? User : Service" />
            </div>
            <div class="message-content">
              <div class="message-header">
                <span class="message-author">{{ msg.is_user ? t('user.tickets.me') : t('user.tickets.staff') }}</span>
                <span class="message-time">{{ formatDate(msg.created_at) }}</span>
              </div>
              <div class="message-body">{{ msg.content }}</div>
            </div>
          </div>
        </el-scrollbar>

        <!-- 回复区域 -->
        <div v-if="currentTicket.status !== 'closed'" class="reply-area">
          <el-divider />
          <el-input
            v-model="replyContent"
            type="textarea"
            :rows="3"
            :placeholder="t('user.tickets.replyPlaceholder')"
          />
          <div class="reply-actions">
            <el-button @click="detailVisible = false">{{ t('common.close') }}</el-button>
            <el-button type="primary" :loading="replyLoading" @click="handleReply">
              {{ t('user.tickets.sendReply') }}
            </el-button>
          </div>
        </div>

        <div v-else class="closed-hint">
          <el-alert :title="t('user.tickets.closedHint')" type="info" :closable="false" />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus, User, Service } from '@element-plus/icons-vue'
import { getTicketList, getTicketDetail, createTicket, replyTicket, closeTicket } from '@/api/ticket'

const { t, locale } = useI18n()

const loading = ref(true)
const ticketList = ref([])
const filterStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const createVisible = ref(false)
const detailVisible = ref(false)
const currentTicket = ref(null)
const messages = ref([])
const replyContent = ref('')
const replyLoading = ref(false)
const submitting = ref(false)
const createFormRef = ref(null)

const createForm = ref({
  category: 'general',
  subject: '',
  content: ''
})

const createRules = {
  category: [{ required: true, message: t('user.tickets.categoryRequired'), trigger: 'change' }],
  subject: [{ required: true, message: t('user.tickets.subjectRequired'), trigger: 'blur' }],
  content: [{ required: true, message: t('user.tickets.contentRequired'), trigger: 'blur' }]
}

// 状态映射
const statusMap = {
  open: { text: t('user.tickets.openStatus'), type: 'primary' },
  replied: { text: t('user.tickets.repliedStatus'), type: 'warning' },
  closed: { text: t('user.tickets.closedStatus'), type: 'info' }
}

const getStatusType = (status) => statusMap[status]?.type || 'info'
const getStatusText = (status) => statusMap[status]?.text || status

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString(locale.value === 'en-US' ? 'en-US' : 'zh-CN')
}

// 加载工单列表
const loadTickets = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      status: filterStatus.value || undefined
    }
    const res = await getTicketList(params)
    if (res.code === 200) {
      ticketList.value = res.data?.list || res.data?.items || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载工单列表失败:', error)
    ElMessage.error(error?.message || t('user.tickets.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleFilterChange = () => {
  currentPage.value = 1
  loadTickets()
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadTickets()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadTickets()
}

// 创建工单
const handleCreate = async () => {
  const valid = await createFormRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const res = await createTicket(createForm.value)
    if (res.code === 200) {
      ElMessage.success(t('user.tickets.createSuccess'))
      createVisible.value = false
      createForm.value = { category: 'general', subject: '', content: '' }
      loadTickets()
    }
  } catch (error) {
    ElMessage.error(error?.message || t('user.tickets.createFailed'))
  } finally {
    submitting.value = false
  }
}

// 打开工单详情
const openTicketDetail = async (ticket) => {
  currentTicket.value = ticket
  detailVisible.value = true
  replyContent.value = ''
  await loadTicketDetail(ticket.id)
}

// 加载工单详情
const loadTicketDetail = async (id) => {
  try {
    const res = await getTicketDetail(id)
    if (res.code === 200) {
      const data = res.data || {}
      messages.value = data.messages || []
      currentTicket.value = { ...currentTicket.value, ...data }
    }
  } catch (error) {
    console.error('加载工单详情失败:', error)
  }
}

// 回复工单
const handleReply = async () => {
  if (!replyContent.value.trim()) {
    ElMessage.warning(t('user.tickets.inputReply'))
    return
  }
  replyLoading.value = true
  try {
    const res = await replyTicket(currentTicket.value.id, { content: replyContent.value })
    if (res.code === 200) {
      ElMessage.success(t('user.tickets.replySuccess'))
      replyContent.value = ''
      await loadTicketDetail(currentTicket.value.id)
      loadTickets()
    }
  } catch (error) {
    ElMessage.error(error?.message || t('user.tickets.replyFailed'))
  } finally {
    replyLoading.value = false
  }
}

onMounted(() => {
  loadTickets()
})
</script>

<style lang="scss" scoped>
.tickets-container {
  padding: 24px;
}

.tickets-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
    color: var(--text-color-primary);
  }
}

.filter-bar {
  margin-bottom: 20px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
}

.ticket-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ticket-card {
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    transform: translateX(4px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  :deep(.el-card__body) {
    padding: 16px 20px;
  }
}

.ticket-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.ticket-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.ticket-id {
  font-weight: 600;
  color: var(--text-color-primary);
}

.ticket-time {
  font-size: 12px;
  color: var(--text-color-secondary);
}

.ticket-subject {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-color-primary);
  margin-bottom: 6px;
}

.ticket-preview {
  font-size: 13px;
  color: var(--text-color-secondary);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.pagination-wrapper {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}

/* 工单详情对话样式 */
.ticket-detail {
  .ticket-meta {
    margin-bottom: 8px;
  }

  .meta-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
  }

  .meta-label {
    color: var(--text-color-secondary);
    font-size: 14px;
  }

  .meta-value {
    font-weight: 500;
    color: var(--text-color-primary);
  }
}

.message-list {
  padding: 8px 0;
}

.message-item {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;

  &.message-self {
    flex-direction: row-reverse;

    .message-content {
      align-items: flex-end;
    }

    .message-body {
      background-color: #dcfce7;
      color: #166534;
    }
  }
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-width: 70%;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.message-author {
  font-weight: 500;
  color: var(--text-color-primary);
}

.message-time {
  color: var(--text-color-secondary);
}

.message-body {
  padding: 10px 14px;
  background-color: var(--neutral-bg);
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.5;
  color: var(--text-color-primary);
  word-break: break-word;
}

.reply-area {
  margin-top: 8px;
}

.reply-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 12px;
}

.closed-hint {
  margin-top: 16px;
}

@media (max-width: 768px) {
  .tickets-container {
    padding: 16px;
  }

  .tickets-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}
</style>
