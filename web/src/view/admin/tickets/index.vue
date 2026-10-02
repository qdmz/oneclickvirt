<template>
  <div class="tickets-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('admin.tickets.title') }}</span>
          <div class="header-stats">
            <el-tag type="primary">{{ t('admin.tickets.total') }}: {{ stats.total }}</el-tag>
            <el-tag type="warning">{{ t('admin.tickets.pending') }}: {{ stats.pending }}</el-tag>
            <el-tag type="success">{{ t('admin.tickets.resolved') }}: {{ stats.solved }}</el-tag>
          </div>
        </div>
      </template>

      <!-- 筛选栏 -->
      <div class="toolbar">
        <el-input
          v-model="searchKeyword"
          :placeholder="t('admin.tickets.searchPlaceholder')"
          style="width: 250px;"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select
          v-model="searchStatus"
          :placeholder="t('admin.tickets.selectStatus')"
          style="width: 150px; margin-left: 10px;"
          clearable
        >
          <el-option :label="t('admin.tickets.all')" value="" />
          <el-option :label="t('admin.tickets.open')" :value="0" />
          <el-option :label="t('admin.tickets.processing')" :value="1" />
          <el-option :label="t('admin.tickets.resolvedStatus')" :value="2" />
          <el-option :label="t('admin.tickets.closed')" :value="3" />
        </el-select>
        <el-button type="primary" style="margin-left: 10px;" @click="handleSearch">
          {{ t('common.search') }}
        </el-button>
        <el-button style="margin-left: 10px;" @click="resetFilters">
          {{ t('common.reset') }}
        </el-button>
      </div>

      <!-- 工单表格 -->
      <el-table v-loading="loading" :data="ticketList" stripe style="width: 100%">
        <el-table-column :label="t('admin.tickets.id')" prop="id" width="80" />
        <el-table-column :label="t('admin.tickets.user')" min-width="120">
          <template #default="{ row }">
            {{ row.username || row.userId }}
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.tickets.subject')" prop="title" min-width="180" />
        <el-table-column :label="t('admin.tickets.category')" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getCategoryText(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.tickets.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.tickets.createTime')" width="160">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.tickets.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="openDetail(row)">
              {{ t('common.reply') }}
            </el-button>
            <el-button
              v-if="row.status !== 3"
              type="warning"
              size="small"
              @click="handleClose(row)"
            >
              {{ t('common.close') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

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
    </el-card>

    <!-- 工单详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      :title="t('admin.tickets.ticketDetail')"
      width="700px"
      class="ticket-detail-dialog"
    >
      <div v-if="currentTicket" class="ticket-detail">
        <div class="ticket-meta">
          <div class="meta-row">
            <span class="meta-label">{{ t('admin.tickets.subject') }}:</span>
            <span class="meta-value">{{ currentTicket.title }}</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">{{ t('admin.tickets.user') }}:</span>
            <span class="meta-value">{{ currentTicket.username || currentTicket.userId }}</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">{{ t('admin.tickets.category') }}:</span>
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
            :class="{ 'message-self': msg.userType === 'admin' }"
          >
            <div class="message-avatar">
              <el-avatar :size="36" :icon="msg.userType === 'user' ? User : Service" />
            </div>
            <div class="message-content">
              <div class="message-header">
                <span class="message-author">{{ msg.userType === 'user' ? t('admin.tickets.customer') : t('admin.tickets.staff') }}</span>
                <span class="message-time">{{ formatDate(msg.createdAt) }}</span>
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
            :placeholder="t('admin.tickets.replyPlaceholder')"
          />
          <div class="reply-actions">
            <el-button type="warning" @click="handleClose(currentTicket)">
              {{ t('common.close') }}
            </el-button>
            <el-button type="primary" :loading="replyLoading" @click="handleReply">
              {{ t('admin.tickets.sendReply') }}
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, User, Service } from '@element-plus/icons-vue'
import { getAdminTicketList, getAdminTicketDetail, adminReplyTicket, updateAdminTicketStatus, getAdminTicketStats } from '@/api/admin'

const { t, locale } = useI18n()

const loading = ref(true)
const ticketList = ref([])
const searchKeyword = ref('')
const searchStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const detailVisible = ref(false)
const currentTicket = ref(null)
const messages = ref([])
const replyContent = ref('')
const replyLoading = ref(false)

const stats = ref({ total: 0, pending: 0, solved: 0 })

const statusMap = {
  0: { text: t('admin.tickets.openStatus'), type: 'primary' },
  1: { text: t('admin.tickets.processingStatus'), type: 'warning' },
  2: { text: t('admin.tickets.resolvedStatus'), type: 'success' },
  3: { text: t('admin.tickets.closedStatus'), type: 'info' }
}

const getStatusType = (status) => statusMap[status]?.type || 'info'
const getStatusText = (status) => statusMap[status]?.text || status

// 分类翻译映射
const categoryMap = {
  general: 'catGeneral',
  technical: 'catTechnical',
  billing: 'catBilling',
  other: 'catOther'
}
const getCategoryText = (category) => {
  const key = categoryMap[category]
  return key ? t(`admin.tickets.${key}`) : category
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
      keyword: searchKeyword.value || undefined,
      status: searchStatus.value === '' ? undefined : searchStatus.value
    }
    const res = await getAdminTicketList(params)
    if (res.code === 200) {
      ticketList.value = res.data?.list || res.data?.items || []
      total.value = res.data?.total || 0
    }
    // Load stats
    const statsRes = await getAdminTicketStats()
    if (statsRes.code === 200 && statsRes.data) {
      stats.value = statsRes.data
    }
  } catch (error) {
    console.error('加载工单列表失败:', error)
    ElMessage.error(error?.message || t('admin.tickets.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadTickets()
}

const resetFilters = () => {
  searchKeyword.value = ''
  searchStatus.value = ''
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

// 打开详情
const openDetail = async (row) => {
  currentTicket.value = row
  detailVisible.value = true
  replyContent.value = ''
  await loadTicketDetail(row.id)
}

const loadTicketDetail = async (id) => {
  try {
    const res = await getAdminTicketDetail(id)
    if (res.code === 200) {
      const data = res.data || {}
      messages.value = data.replies || []
      if (data.ticket) {
        currentTicket.value = { ...currentTicket.value, ...data.ticket }
      }
      if (data.user && data.user.username) {
        currentTicket.value = { ...currentTicket.value, username: data.user.username }
      }
    }
  } catch (error) {
    console.error('加载工单详情失败:', error)
  }
}

// 回复
const handleReply = async () => {
  if (!replyContent.value.trim()) {
    ElMessage.warning(t('admin.tickets.inputReply'))
    return
  }
  replyLoading.value = true
  try {
    const res = await adminReplyTicket(currentTicket.value.id, { content: replyContent.value })
    if (res.code === 200) {
      ElMessage.success(t('admin.tickets.replySuccess'))
      replyContent.value = ''
      await loadTicketDetail(currentTicket.value.id)
      loadTickets()
    }
  } catch (error) {
    ElMessage.error(error?.message || t('admin.tickets.replyFailed'))
  } finally {
    replyLoading.value = false
  }
}

// 关闭工单
const handleClose = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('admin.tickets.confirmClose'),
      t('common.tip'),
      { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' }
    )
    const res = await updateAdminTicketStatus(row.id, { status: 3 })  // 3 = closed
    if (res.code === 200) {
      ElMessage.success(t('admin.tickets.closeSuccess'))
      if (detailVisible.value) detailVisible.value = false
      loadTickets()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error?.message || t('admin.tickets.closeFailed'))
    }
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

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.header-stats {
  display: flex;
  gap: 8px;
}

.toolbar {
  margin-bottom: 20px;
}

.pagination-wrapper {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}

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
