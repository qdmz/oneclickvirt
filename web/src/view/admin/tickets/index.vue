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
      color: var(--accent-text-color-hover);
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
</style>
