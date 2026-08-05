<template>
  <div class="user-tickets">
    <div class="page-header">
      <div>
        <h1>{{ t('user.tickets.title') }}</h1>
        <p>{{ t('user.tickets.subtitle') }}</p>
      </div>
      <el-button
        type="primary"
        :icon="Plus"
        @click="openCreate"
      >
        {{ t('user.tickets.create') }}
      </el-button>
    </div>

    <!-- 筛选器 -->
    <div class="filter-section">
      <el-form :inline="true">
        <el-form-item>
          <el-select
            v-model="filter.status"
            :placeholder="t('user.tickets.filterStatus')"
            clearable
            style="width: 160px;"
            @change="loadList"
          >
            <el-option
              :label="t('user.tickets.statusOpen')"
              value="open"
            />
            <el-option
              :label="t('user.tickets.statusAnswered')"
              value="answered"
            />
            <el-option
              :label="t('user.tickets.statusClosed')"
              value="closed"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="filter.keyword"
            :placeholder="t('user.tickets.searchPlaceholder')"
            clearable
            style="width: 220px;"
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

    <!-- 工单列表 -->
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
          label="#"
          width="70"
        />
        <el-table-column
          prop="title"
          :label="t('user.tickets.colTitle')"
          min-width="220"
          show-overflow-tooltip
        />
        <el-table-column
          :label="t('user.tickets.colPriority')"
          width="110"
        >
          <template #default="{ row }">
            <el-tag
              :type="priorityTagType(row.priority)"
              size="small"
              effect="light"
            >
              {{ t('user.tickets.priority.' + row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('user.tickets.colStatus')"
          width="120"
        >
          <template #default="{ row }">
            <el-tag
              :type="statusTagType(row.status)"
              size="small"
            >
              {{ t('user.tickets.status' + capitalize(row.status)) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="replyCount"
          :label="t('user.tickets.colReplies')"
          width="90"
        />
        <el-table-column
          :label="t('user.tickets.colLastReply')"
          width="180"
        >
          <template #default="{ row }">
            {{ formatTime(row.lastReplyAt) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="t('common.operation')"
          width="120"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              @click="openDetail(row.id)"
            >
              {{ t('user.tickets.view') }}
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

    <!-- 创建工单对话框 -->
    <el-dialog
      v-model="createVisible"
      :title="t('user.tickets.create')"
      width="560px"
    >
      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-width="90px"
      >
        <el-form-item
          :label="t('user.tickets.colTitle')"
          prop="title"
        >
          <el-input
            v-model="createForm.title"
            :placeholder="t('user.tickets.titlePlaceholder')"
            maxlength="255"
          />
        </el-form-item>
        <el-form-item :label="t('user.tickets.colPriority')">
          <el-select
            v-model="createForm.priority"
            style="width: 100%;"
          >
            <el-option
              :label="t('user.tickets.priority.low')"
              value="low"
            />
            <el-option
              :label="t('user.tickets.priority.normal')"
              value="normal"
            />
            <el-option
              :label="t('user.tickets.priority.high')"
              value="high"
            />
            <el-option
              :label="t('user.tickets.priority.urgent')"
              value="urgent"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="t('user.tickets.content')"
          prop="content"
        >
          <el-input
            v-model="createForm.content"
            type="textarea"
            :rows="5"
            :placeholder="t('user.tickets.contentPlaceholder')"
            maxlength="5000"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="submitCreate"
        >
          {{ t('common.submit') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 工单详情抽屉 -->
    <el-drawer
      v-model="detailVisible"
      :title="detail ? ('#' + detail.id + ' ' + detail.title) : t('user.tickets.detail')"
      size="560px"
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
            {{ t('user.tickets.status' + capitalize(detail.status)) }}
          </el-tag>
          <el-tag
            :type="priorityTagType(detail.priority)"
            size="small"
            effect="light"
          >
            {{ t('user.tickets.priority.' + detail.priority) }}
          </el-tag>
          <span class="meta-time">{{ t('user.tickets.createdAt') }} {{ formatTime(detail.createdAt) }}</span>
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
                >{{ t('user.tickets.staff') }}</el-tag>
                {{ r.username || (r.isAdmin ? t('user.tickets.staff') : t('user.tickets.me')) }}
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
            :placeholder="t('user.tickets.replyPlaceholder')"
            maxlength="5000"
          />
          <div class="reply-actions">
            <el-popconfirm
              :title="t('user.tickets.confirmClose')"
              @confirm="doClose"
            >
              <template #reference>
                <el-button plain>
                  {{ t('user.tickets.close') }}
                </el-button>
              </template>
            </el-popconfirm>
            <el-button
              type="primary"
              :loading="replying"
              @click="doReply"
            >
              {{ t('user.tickets.reply') }}
            </el-button>
          </div>
        </div>
        <el-alert
          v-else
          :title="t('user.tickets.closedTip')"
          type="info"
          :closable="false"
          show-icon
        />
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getMyTickets, createMyTicket, getMyTicketDetail, replyMyTicket, closeMyTicket
} from '@/api/user'

const { t } = useI18n()

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const filter = reactive({ status: '', keyword: '' })

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

const loadList = async () => {
  loading.value = true
  try {
    const res = await getMyTickets({
      page: page.value,
      pageSize: pageSize.value,
      status: filter.status,
      keyword: filter.keyword
    })
    if (res && res.code === 200) {
      list.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (e) {
    // ignore, interceptor handles messaging
  } finally {
    loading.value = false
  }
}

// 创建
const createVisible = ref(false)
const createFormRef = ref()
const submitting = ref(false)
const createForm = reactive({ title: '', content: '', priority: 'normal' })
const createRules = {
  title: [{ required: true, message: t('user.tickets.titleRequired'), trigger: 'blur' }],
  content: [{ required: true, message: t('user.tickets.contentRequired'), trigger: 'blur' }]
}
const openCreate = () => {
  createForm.title = ''
  createForm.content = ''
  createForm.priority = 'normal'
  createVisible.value = true
}
const submitCreate = async () => {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const res = await createMyTicket({ ...createForm })
      if (res && res.code === 200) {
        ElMessage.success(t('user.tickets.createSuccess'))
        createVisible.value = false
        page.value = 1
        loadList()
      }
    } finally {
      submitting.value = false
    }
  })
}

// 详情 + 回复
const detailVisible = ref(false)
const detail = ref(null)
const replyContent = ref('')
const replying = ref(false)

const openDetail = async (id) => {
  try {
    const res = await getMyTicketDetail(id)
    if (res && res.code === 200) {
      detail.value = res.data
      replyContent.value = ''
      detailVisible.value = true
    }
  } catch (e) {
    // ignore
  }
}
const doReply = async () => {
  if (!replyContent.value.trim()) {
    ElMessage.warning(t('user.tickets.replyRequired'))
    return
  }
  replying.value = true
  try {
    const res = await replyMyTicket(detail.value.id, replyContent.value)
    if (res && res.code === 200) {
      ElMessage.success(t('user.tickets.replySuccess'))
      await openDetail(detail.value.id)
      loadList()
    }
  } finally {
    replying.value = false
  }
}
const doClose = async () => {
  const res = await closeMyTicket(detail.value.id)
  if (res && res.code === 200) {
    ElMessage.success(t('user.tickets.closeSuccess'))
    await openDetail(detail.value.id)
    loadList()
  }
}

onMounted(loadList)
</script>

<style scoped>
.user-tickets {
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
  padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  margin-bottom: 12px;
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
  margin-top: 10px;
}
</style>
