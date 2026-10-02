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
              <el-tag type="info" size="small">{{ getCategoryText(ticket.category) }}</el-tag>
            </div>
            <span class="ticket-time">{{ formatDate(ticket.updated_at || ticket.created_at) }}</span>
          </div>
          <div class="ticket-subject">{{ ticket.title || ticket.subject }}</div>
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
        <el-form-item :label="t('user.tickets.subject')" prop="title">
          <el-input v-model="createForm.title" :placeholder="t('user.tickets.inputSubject')" />
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
            <span class="meta-value">{{ currentTicket.title || currentTicket.subject }}</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">{{ t('user.tickets.category') }}:</span>
            <el-tag size="small">{{ getCategoryText(currentTicket.category) }}</el-tag>
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
        <div v-if="currentTicket.status !== 3" class="reply-area">
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
  title: '',
  content: ''
})

const createRules = {
  category: [{ required: true, message: t('user.tickets.categoryRequired'), trigger: 'change' }],
  title: [{ required: true, message: t('user.tickets.subjectRequired'), trigger: 'blur' }],
  content: [{ required: true, message: t('user.tickets.contentRequired'), trigger: 'blur' }]
}

// 状态映射 (后端返回数字: 0=待处理 1=处理中 2=已解决 3=已关闭)
const statusMap = {
  0: { text: t('user.tickets.openStatus'), type: 'primary' },
  1: { text: t('user.tickets.repliedStatus'), type: 'warning' },
  2: { text: t('user.tickets.solvedStatus') || '已解决', type: 'success' },
  3: { text: t('user.tickets.closedStatus'), type: 'info' }
}

// 筛选映射 (前端字符串 -> 后端数字)
const filterStatusMap = {
  '': undefined,
  'open': 0,
  'replied': 1,
  'closed': 3
}

const getStatusType = (status) => statusMap[status]?.type || 'info'
const getStatusText = (status) => statusMap[status]?.text || status

// 分类翻译
const categoryMap = {
  general: 'catGeneral',
  technical: 'catTechnical',
  billing: 'catBilling',
  other: 'catOther'
}
const getCategoryText = (category) => {
  const key = categoryMap[category]
  return key ? t(`user.tickets.${key}`) : category
}

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
      status: filterStatusMap[filterStatus.value]
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
      createForm.value = { category: 'general', title: '', content: '' }
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
      // 兼容 replies 和 messages 字段，并映射 is_user 字段
      const replies = data.messages || data.replies || []
      messages.value = replies.map(r => ({
        ...r,
        is_user: r.userType === 'user' || r.is_user === true
      }))
      currentTicket.value = { ...currentTicket.value, ...data.ticket, ...data }
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


/* ── 现代 UI 增强 ── */

/* 卡片 */
:deep(.el-card),
:deep(.el-card.is-always-shadow),
:deep(.el-card.is-hover-shadow) {
  border-radius: 16px !important;
  border: 1px solid var(--border-color) !important;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05) !important;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1),
              box-shadow 0.3s cubic-bezier(0.4, 0, 0.2, 1),
              border-color 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.el-card:hover) {
  transform: translateY(-3px) !important;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.09) !important;
}

/* 统计卡片 */
.stat-card :deep(.el-card__body) {
  padding: 20px !important;
}

.stat-content {
  display: flex !important;
  align-items: center !important;
  gap: 14px !important;
}

.stat-icon {
  width: 48px !important;
  height: 48px !important;
  border-radius: 14px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  font-size: 22px !important;
  flex-shrink: 0 !important;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.08) !important;
}

.stat-icon.user-icon { background: linear-gradient(135deg, #6366f1, #818cf8); color: #fff; }
.stat-icon.server-icon { background: linear-gradient(135deg, #10b981, #34d399); color: #fff; }
.stat-icon.vm-icon { background: linear-gradient(135deg, #f59e0b, #fbbf24); color: #fff; }
.stat-icon.container-icon { background: linear-gradient(135deg, #ec4899, #f472b6); color: #fff; }
.stat-icon.cpu-icon { background: linear-gradient(135deg, #6366f1, #818cf8); color: #fff; }
.stat-icon.memory-icon { background: linear-gradient(135deg, #10b981, #34d399); color: #fff; }
.stat-icon.disk-icon { background: linear-gradient(135deg, #f59e0b, #fbbf24); color: #fff; }

.stat-number {
  font-size: 28px !important;
  font-weight: 800 !important;
  color: var(--text-color-primary) !important;
  line-height: 1.1 !important;
  font-variant-numeric: tabular-nums !important;
}

.stat-label {
  font-size: 13px !important;
  color: var(--text-color-secondary) !important;
  font-weight: 600 !important;
  margin-top: 2px !important;
}

/* 资源卡片 */
.resource-card :deep(.el-card__body) {
  padding: 16px 20px !important;
}

.resource-header {
  display: flex !important;
  align-items: center !important;
  gap: 10px !important;
  margin-bottom: 12px !important;
}

.resource-icon {
  width: 36px !important;
  height: 36px !important;
  border-radius: 10px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  font-size: 16px !important;
  flex-shrink: 0 !important;
}

.resource-title {
  font-weight: 700 !important;
  font-size: 14px !important;
  color: var(--text-color-primary) !important;
}

.resource-number-value {
  font-size: 22px !important;
  font-weight: 800 !important;
  color: var(--primary-color) !important;
  font-variant-numeric: tabular-nums !important;
}

.resource-number-label {
  font-size: 11px !important;
  color: var(--text-color-secondary) !important;
  font-weight: 600 !important;
}

.resource-number-separator {
  color: var(--text-color-tertiary) !important;
  font-weight: 700 !important;
  font-size: 16px !important;
  margin: 0 8px !important;
}

/* 按钮 */
:deep(.el-button--primary) {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-light)) !important;
  border: none !important;
  border-radius: 12px !important;
  font-weight: 600 !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1) !important;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.2) !important;
}

:deep(.el-button--primary:hover) {
  transform: translateY(-2px) !important;
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.3) !important;
}

:deep(.el-button--success) {
  border-radius: 12px !important;
  font-weight: 600 !important;
}

:deep(.el-button--warning) {
  border-radius: 12px !important;
  font-weight: 600 !important;
}

:deep(.el-button--danger) {
  border-radius: 12px !important;
  font-weight: 600 !important;
}

:deep(.el-button--default) {
  border-radius: 12px !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.el-button--default:hover) {
  transform: translateY(-1px) !important;
}

/* 表格 */
:deep(.el-table) {
  --el-table-border-color: var(--border-color) !important;
}

:deep(.el-table th) {
  background: rgba(99, 102, 241, 0.04) !important;
  color: var(--text-color-primary) !important;
  font-weight: 700 !important;
  font-size: 13px !important;
  letter-spacing: 0.3px !important;
}

:deep(.el-table td) {
  font-size: 13px !important;
  color: var(--text-color-primary) !important;
}

:deep(.el-table__row) {
  transition: all 0.2s ease !important;
}

:deep(.el-table__row:hover) {
  background: rgba(99, 102, 241, 0.04) !important;
  transform: scale(1.005) !important;
}

/* 徽章 */
:deep(.el-tag) {
  border-radius: 999px !important;
  font-weight: 600 !important;
  font-size: 12px !important;
  padding: 2px 12px !important;
  border: none !important;
}

:deep(.el-tag--success) {
  background: rgba(16, 185, 129, 0.12) !important;
  color: #10b981 !important;
}

:deep(.el-tag--warning) {
  background: rgba(245, 158, 11, 0.12) !important;
  color: #f59e0b !important;
}

:deep(.el-tag--danger) {
  background: rgba(239, 68, 68, 0.12) !important;
  color: #ef4444 !important;
}

:deep(.el-tag--info) {
  background: rgba(107, 114, 128, 0.1) !important;
  color: #6b7280 !important;
}

/* 输入框 */
:deep(.el-input__wrapper) {
  border-radius: 12px !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2) !important;
}

:deep(.el-select .el-input__wrapper),
:deep(.el-cascader .el-input__wrapper) {
  border-radius: 12px !important;
}

/* 分页 */
:deep(.el-pagination) {
  margin-top: 24px !important;
}

:deep(.el-pager li),
:deep(.el-pager li.btn-quicknext),
:deep(.el-pager li.btn-quickprev),
:deep(.el-pagination .btn-prev),
:deep(.el-pagination .btn-next) {
  border-radius: 10px !important;
  font-weight: 600 !important;
  transition: all 0.2s ease !important;
}

:deep(.el-pager li.is-active) {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-light)) !important;
  border-color: transparent !important;
}

/* 加载态 */
.el-loading-mask {
  backdrop-filter: blur(4px) !important;
  -webkit-backdrop-filter: blur(4px) !important;
  background: rgba(255, 255, 255, 0.7) !important;
}

.el-dark .el-loading-mask {
  background: rgba(15, 23, 42, 0.7) !important;
}

/* 标题 */
h1, h2, h3, h4 {
  color: var(--text-color-primary) !important;
}

/* 进场动画 */
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(18px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 响应式断点 */
@media (max-width: 768px) {
  .stat-content {
    flex-direction: column !important;
    align-items: flex-start !important;
  }

  .stat-number {
    font-size: 22px !important;
  }

  .stat-icon {
    width: 40px !important;
    height: 40px !important;
    font-size: 18px !important;
  }

  .resource-header {
    flex-direction: column !important;
    align-items: flex-start !important;
  }

  .resource-numbers {
    display: flex !important;
    align-items: center !important;
    gap: 4px !important;
    flex-wrap: wrap !important;
  }
}

@media (max-width: 480px) {
  .stat-number {
    font-size: 20px !important;
  }

  .stat-label {
    font-size: 12px !important;
  }

  .resource-number-value {
    font-size: 18px !important;
  }
}
</style>
