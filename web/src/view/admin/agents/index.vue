<template>
  <div class="admin-agents">
    <div class="page-header">
      <h1>代理商管理</h1>
    </div>

    <el-card>
      <!-- 工具栏 -->
      <div class="toolbar">
        <el-input
          v-model="keyword"
          placeholder="搜索名称/推广码/邮箱"
          clearable
          style="width: 250px"
          @clear="loadData"
          @keyup.enter="loadData"
        />
        <div class="toolbar-actions">
          <el-select v-model="statusFilter" placeholder="状态" clearable style="width: 120px" @change="loadData">
            <el-option label="全部" :value="null" />
            <el-option label="待审核" :value="0" />
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
          <el-button type="primary" @click="loadData">搜索</el-button>
        </div>
      </div>

      <el-table v-loading="loading" :data="agents" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="code" label="推广码" width="140">
          <template #default="{ row }">
            <el-tag size="small">{{ row.code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="关联用户" width="140">
          <template #default="{ row }">{{ row.user?.username || '-' }}</template>
        </el-table-column>
        <el-table-column prop="contactEmail" label="联系邮箱" min-width="180">
          <template #default="{ row }">{{ row.contactEmail || '-' }}</template>
        </el-table-column>
        <el-table-column prop="commissionRate" label="佣金比例" width="100" align="center">
          <template #default="{ row }">{{ row.commissionRate }}%</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="{ 0: 'warning', 1: 'success', 2: 'danger' }[row.status]" size="small">
              {{ { 0: '待审核', 1: '正常', 2: '禁用' }[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 0" size="small" type="success" @click="handleApprove(row)">
              通过
            </el-button>
            <el-button size="small" @click="handleAdjustCommission(row)">
              佣金
            </el-button>
            <el-button size="small" @click="handleViewDetail(row)">
              详情
            </el-button>
            <el-button v-if="row.status === 1" size="small" type="danger" @click="handleDisable(row)">
              禁用
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="代理商详情" width="800px">
      <template v-if="detailData">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ detailData.agent?.id }}</el-descriptions-item>
          <el-descriptions-item label="名称">{{ detailData.agent?.name }}</el-descriptions-item>
          <el-descriptions-item label="推广码">{{ detailData.agent?.code }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="{ 0: 'warning', 1: 'success', 2: 'danger' }[detailData.agent?.status]">
              {{ { 0: '待审核', 1: '正常', 2: '禁用' }[detailData.agent?.status] }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="佣金比例">{{ detailData.agent?.commissionRate }}%</el-descriptions-item>
          <el-descriptions-item label="子用户数">{{ detailData.statistics?.subUserCount }}</el-descriptions-item>
          <el-descriptions-item label="累计佣金">¥{{ (detailData.statistics?.totalCommission / 100).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="余额">¥{{ (detailData.statistics?.balance / 100).toFixed(2) }}</el-descriptions-item>
        </el-descriptions>

        <h4 style="margin: 20px 0 10px;">子用户列表</h4>
        <el-table :data="subUsers" stripe size="small">
          <el-table-column prop="username" label="用户名" width="150" />
          <el-table-column prop="nickname" label="昵称" width="150">
            <template #default="{ row }">{{ row.nickname || '-' }}</template>
          </el-table-column>
          <el-table-column prop="email" label="邮箱" min-width="200" />
          <el-table-column prop="user_status" label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.user_status === 1 ? 'success' : 'danger'" size="small">
                {{ row.user_status === 1 ? '正常' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </template>
    </el-dialog>

    <!-- 调整佣金对话框 -->
    <el-dialog v-model="commissionVisible" title="调整佣金比例" width="400px">
      <el-form label-width="80px">
        <el-form-item label="当前比例">
          {{ currentAgent?.commissionRate }}%
        </el-form-item>
        <el-form-item label="新比例">
          <el-input-number v-model="newRate" :min="0" :max="100" :precision="1" :step="1" /> %
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="commissionVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="doAdjustCommission">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  getAgentList, approveAgent, updateAgentStatus,
  adjustCommission, getAgentDetail, getAgentSubUsers
} from '@/api/agent'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const agents = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')
const statusFilter = ref(null)

const detailVisible = ref(false)
const detailData = ref(null)
const subUsers = ref([])

const commissionVisible = ref(false)
const currentAgent = ref(null)
const newRate = ref(0)
const saving = ref(false)

const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try {
    const res = await getAgentList({
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value,
      ...(statusFilter.value !== null ? { status: statusFilter.value } : {})
    })
    if (res.code === 0 || res.code === 200) {
      agents.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch { ElMessage.error('加载失败') }
  finally { loading.value = false }
}

const handleApprove = async (row) => {
  try {
    await ElMessageBox.confirm('确定通过该代理商申请？', '确认', { type: 'info' })
    await approveAgent(row.id)
    ElMessage.success('审核通过')
    loadData()
  } catch (e) { if (e !== 'cancel') ElMessage.error('操作失败') }
}

const handleDisable = async (row) => {
  try {
    await ElMessageBox.confirm('确定禁用该代理商？', '确认', { type: 'warning' })
    await updateAgentStatus(row.id, 2)
    ElMessage.success('已禁用')
    loadData()
  } catch (e) { if (e !== 'cancel') ElMessage.error('操作失败') }
}

const handleAdjustCommission = (row) => {
  currentAgent.value = row
  newRate.value = row.commissionRate
  commissionVisible.value = true
}

const doAdjustCommission = async () => {
  saving.value = true
  try {
    await adjustCommission(currentAgent.value.id, newRate.value)
    ElMessage.success('佣金比例已调整')
    commissionVisible.value = false
    loadData()
  } catch { ElMessage.error('操作失败') }
  finally { saving.value = false }
}

const handleViewDetail = async (row) => {
  detailVisible.value = true
  try {
    const [detailRes, subRes] = await Promise.all([
      getAgentDetail(row.id),
      getAgentSubUsers(row.id, { page: 1, pageSize: 50 })
    ])
    if (detailRes.code === 0 || detailRes.code === 200) detailData.value = detailRes.data
    if (subRes.code === 0 || subRes.code === 200) subUsers.value = subRes.data?.list || []
  } catch { ElMessage.error('加载详情失败') }
}

onMounted(loadData)
</script>

<style scoped>
.admin-agents { padding: 20px; }
.page-header { margin-bottom: 20px; }
.page-header h1 { font-size: 24px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.toolbar-actions { display: flex; gap: 8px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
