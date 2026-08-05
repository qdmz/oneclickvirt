<template>
  <div class="admin-tickets">
    <div class="page-header">
      <div>
        <h1>{{ t('admin.tickets.title') }}</h1>
        <p>{{ t('admin.tickets.subtitle') }}</p>
      </div>
      <el-button :icon="Refresh" @click="refreshAll">
        {{ t('admin.tickets.refresh') }}
      </el-button>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-row">
      <el-card
        shadow="never"
        class="stat-card"
      >
        <div class="stat-value">
          {{ statsTotal }}
        </div>
        <div class="stat-label">
          {{ t('admin.tickets.statsTotal') }}
        </div>
      </el-card>
      <el-card
        shadow="never"
        class="stat-card stat-open"
      >
        <div class="stat-value">
          {{ stats.open }}
        </div>
        <div class="stat-label">
          {{ t('admin.tickets.statsOpen') }}
        </div>
      </el-card>
      <el-card
        shadow="never"
        class="stat-card stat-answered"
      >
        <div class="stat-value">
          {{ stats.answered }}
        </div>
        <div class="stat-label">
          {{ t('admin.tickets.statsAnswered') }}
        </div>
      </el-card>
      <el-card
        shadow="never"
        class="stat-card stat-closed"
      >
        <div class="stat-value">
          {{ stats.closed }}
        </div>
        <div class="stat-label">
          {{ t('admin.tickets.statsClosed') }}
        </div>
      </el-card>
    </div>

    <!-- 筛选 -->
    <div class="filter-section">
      <el-form :inline="true">
        <el-form-item>
          <el-select
            v-model="filter.status"
            :placeholder="t('admin.tickets.filterStatus')"
            clearable
            style="width: 150px;"
            @change="loadList"
          >
            <el-option
              :label="t('admin.tickets.statusOpen')"
              value="open"
            />
            <el-option
              :label="t('admin.tickets.statusAnswered')"
              value="answered"
            />
            <el-option
              :label="t('admin.tickets.statusClosed')"
              value="closed"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="filter.priority"
            :placeholder="t('admin.tickets.filterPriority')"
            clearable
            style="width: 150px;"
            @change="loadList"
          >
            <el-option
              :label="t('admin.tickets.priority.low')"
              value="low"
            />
            <el-option
              :label="t('admin.tickets.priority.normal')"
              value="normal"
            />
            <el-option
              :label="t('admin.tickets.priority.high')"
              value="high"
            />
            <el-option
              :label="t('admin.tickets.priority.urgent')"
              value="urgent"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="filter.keyword"
            :placeholder="t('admin.tickets.searchPlaceholder')"
            clearable
            style="width: 240px;"
            @keyup.enter="loadList"
            @clear="loadList"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            @click="loadList"
          >
            {{ t('common.search') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 列表 -->
    <el-card
      shadow="never"
      class="table-card"
    >
      <el-table
        v-loading="loading"
        :data="list"
        style="width: 100%;"
      >
        <el-table-column
          prop="id"
          :label="t('admin.tickets.colId')"
          width="70"
        />
        <el-table-column
          prop="title"
          :label="t('admin.tickets.colTitle')"
          min-width="220"
          show-overflow-tooltip
        />
        <el-table-column
          prop="username"
          :label="t('admin.tickets.colUser')"
          width="140"
        />
        <el-table-column
          :label="t('admin.tickets.colPriority')"
          width="110"
        >
          <template #default="{ row }">
            <el-tag
              :type="priorityTagType(row.priority)"
              size="small"
              effect="light"
            >
              {{ t('admin.tickets.priority.' + row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.tickets.colStatus')"
          width="120"
        >
          <template #default="{ row }">
            <el-tag
              :type="statusTagType(row.status)"
              size="small"
            >
              {{ t('admin.tickets.status' + capitalize(row.status)) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="replyCount"
          :label="t('admin.tickets.colReplies')"
          width="90"
        />
        <el-table-column
          :label="t('admin.tickets.colLastReply')"
          width="180"
        >
          <template #default="{ row }">
            {{ formatTime(row.lastReplyAt) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('common.operation')"
          width="140"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              @click="openDetail(row.id)"
            >
              {{ t('admin.tickets.view') }}
            </el-button>
            <el-button
              type="danger"
              link
              @click="doDelete(row)"
            >
              {{ t('admin.tickets.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @current-change="loadList"
          @size-change="loadList"
        />
      </div>
    </el-card>

    <!-- 详情抽屉 -->
    <el-drawer
      v-model="detailVisible"
      :title="detail ? ('#' + detail.id + ' ' + detail.title) : t('admin.tickets.detail')"
      size="600px"
    >
      <div
        v-if="detail"
        class="ticket-detail"
      >
        <div class="detail-meta">
          <el-tag
            :type="statusTagType(detail.status)"
            size="small"
          >
            {{ t('admin.tickets.status' + capitalize(detail.status)) }}
          </el-tag>
          <el-tag
            :type="priorityTagType(detail.priority)"
            size="small"
            effect="light"
          >
            {{ t('admin.tickets.priority.' + detail.priority) }}
          </el-tag>
          <span class="meta-user">{{ t('admin.tickets.detailUser') }}: {{ detail.username || '-' }}</span>
          <span class="meta-time">{{ t('admin.tickets.detailCreatedAt') }} {{ formatTime(detail.createdAt) }}</span>
        </div>

        <div class="thread">
          <div
            v-for="r in detail.replies"
            :key="r.id"
            class="msg"
            :class="{ 'msg-admin': r.isAdmin, 'msg-user': !r.isAdmin }"
          >
            <div class="msg-header">
              <span class="msg-author">
                <el-tag
                  v-if="r.isAdmin"
                  type="warning"
                  size="small"
                  effect="dark"
                >{{ t('admin.tickets.staff') }}</el-tag>
                {{ r.username || (r.isAdmin ? t('admin.tickets.staff') : t('admin.tickets.user')) }}
              </span>
              <span class="msg-time">{{ formatTime(r.createdAt) }}</span>
            </div>
            <div class="msg-body">
              {{ r.content }}
            </div>
          </div>
        </div>

        <div
          v-if="detail.status !== 'closed'"
          class="reply-box"
        >
          <el-input
            v-model="replyContent"
            type="textarea"
            :rows="3"
            :placeholder="t('admin.tickets.replyPlaceholder')"
            maxlength="5000"
          />
          <div class="reply-actions">
            <el-select
              v-model="statusSelect"
              size="default"
              style="width: 140px;"
              @change="onStatusChange"
            >
              <el-option
                :label="t('admin.tickets.statusOpen')"
                value="open"
              />
              <el-option
                :label="t('admin.tickets.statusAnswered')"
                value="answered"
              />
              <el-option
                :label="t('admin.tickets.statusClosed')"
                value="closed"
              />
            </el-select>
            <el-popconfirm
              :title="t('admin.tickets.confirmClose')"
              @confirm="doClose"
            >
              <template #reference>
                <el-button plain>
                  {{ t('admin.tickets.close') }}
                </el-button>
              </template>
            </el-popconfirm>
            <el-button
              type="primary"
              :loading="replying"
              @click="doReply"
            >
              {{ t('admin.tickets.reply') }}
            </el-button>
          </div>
        </div>
        <el-alert
          v-else
          :title="t('admin.tickets.closedTip')"
          type="info"
          :closable="false"
          show-icon
        />
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import {
  getAdminTickets, getAdminTicketStats, getAdminTicketDetail,
  replyAdminTicket, updateAdminTicketStatus, closeAdminTicket, deleteAdminTicket
} from '@/api/admin'

const { t } = useI18n()

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const filter = reactive({ status: '', priority: '', keyword: '' })

const stats = reactive({ open: 0, answered: 0, closed: 0 })
const statsTotal = computed(() => stats.open + stats.answered + stats.closed)

const capitalize = (s) => (s ? s.charAt(0).toUpperCase() + s.slice(1) : '')
const statusTagType = (s) => ({ open: 'danger', answered: 'success', closed: 'info' }[s] || 'info')
const priorityTagType = (p) => ({ low: 'info', normal: '', high: 'warning', urgent: 'danger' }[p] || '')

const formatTime = (val) => {
  if (!val) return '-'
  const d = new Date(val)
  if (isNaN(d.getTime())) return '-'
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const loadStats = async () => {
  try {
    const res = await getAdminTicketStats()
    if (res && res.code === 200 && res.data) {
      stats.open = res.data.open || 0
      stats.answered = res.data.answered || 0
      stats.closed = res.data.closed || 0
    }
  } catch (e) {
    // ignore
  }
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await getAdminTickets({
      page: page.value,
      pageSize: pageSize.value,
      status: filter.status,
      priority: filter.priority,
      keyword: filter.keyword
    })
    if (res && res.code === 200) {
      list.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (e) {
    // ignore
  } finally {
    loading.value = false
  }
}

const refreshAll = () => {
  loadList()
  loadStats()
}

// 详情 + 回复
const detailVisible = ref(false)
const detail = ref(null)
const replyContent = ref('')
const statusSelect = ref('open')
const replying = ref(false)

const openDetail = async (id) => {
  try {
    const res = await getAdminTicketDetail(id)
    if (res && res.code === 200) {
      detail.value = res.data
      replyContent.value = ''
      statusSelect.value = res.data.status
      detailVisible.value = true
    }
  } catch (e) {
    // ignore
  }
}

const doReply = async () => {
  if (!replyContent.value.trim()) {
    ElMessage.warning(t('admin.tickets.replyRequired'))
    return
  }
  replying.value = true
  try {
    const res = await replyAdminTicket(detail.value.id, replyContent.value)
    if (res && res.code === 200) {
      ElMessage.success(t('admin.tickets.replySuccess'))
      await openDetail(detail.value.id)
      refreshAll()
    }
  } finally {
    replying.value = false
  }
}

const onStatusChange = async (val) => {
  if (!detail.value || val === detail.value.status) return
  const res = await updateAdminTicketStatus(detail.value.id, val)
  if (res && res.code === 200) {
    ElMessage.success(t('admin.tickets.statusUpdateSuccess'))
    detail.value.status = val
    refreshAll()
  } else {
    statusSelect.value = detail.value.status
  }
}

const doClose = async () => {
  const res = await closeAdminTicket(detail.value.id)
  if (res && res.code === 200) {
    ElMessage.success(t('admin.tickets.closeSuccess'))
    detail.value.status = 'closed'
    refreshAll()
  }
}

const doDelete = async (row) => {
  try {
    await ElMessageBox.confirm(t('admin.tickets.confirmDelete'), t('admin.tickets.delete'), {
      type: 'warning'
    })
  } catch (e) {
    return
  }
  const res = await deleteAdminTicket(row.id)
  if (res && res.code === 200) {
    ElMessage.success(t('admin.tickets.deleteSuccess'))
    refreshAll()
  }
}

onMounted(refreshAll)
</script>

<style scoped>
.admin-tickets {
  padding: 20px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}
.page-header h1 {
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 4px;
}
.page-header p {
  color: var(--el-text-color-secondary);
  margin: 0;
  font-size: 13px;
}
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}
.stat-card {
  border: 1px solid var(--el-border-color-lighter);
  text-align: center;
}
.stat-value {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
}
.stat-label {
  margin-top: 4px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
}
.stat-open .stat-value {
  color: var(--el-color-danger);
}
.stat-answered .stat-value {
  color: var(--el-color-success);
}
.stat-closed .stat-value {
  color: var(--el-text-color-secondary);
}
.filter-section {
  margin-bottom: 16px;
}
.table-card {
  border: 1px solid var(--el-border-color-lighter);
}
.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
.ticket-detail {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.detail-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  margin-bottom: 12px;
}
.meta-user {
  font-size: 13px;
}
.meta-time {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  margin-left: auto;
}
.thread {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-right: 4px;
}
.msg {
  padding: 10px 12px;
  border-radius: 10px;
  max-width: 92%;
}
.msg-user {
  background: var(--el-color-primary-light-9);
  align-self: flex-end;
}
.msg-admin {
  background: var(--el-fill-color-light);
  align-self: flex-start;
}
.msg-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}
.msg-author {
  font-size: 12px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}
.msg-time {
  font-size: 11px;
  color: var(--el-text-color-secondary);
}
.msg-body {
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}
.reply-box {
  border-top: 1px solid var(--el-border-color-lighter);
  padding-top: 12px;
  margin-top: 12px;
}
.reply-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
}
</style>
