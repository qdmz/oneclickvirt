<template>
  <div class="sub-users-page">
    <div class="page-header">
      <h1>子用户管理</h1>
    </div>

    <el-card>
      <!-- 搜索和操作栏 -->
      <div class="toolbar">
        <el-input
          v-model="keyword"
          placeholder="搜索用户名/邮箱/昵称"
          clearable
          style="width: 250px"
          @clear="loadData"
          @keyup.enter="loadData"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <div class="toolbar-actions">
          <el-button type="primary" @click="loadData">搜索</el-button>
          <el-button
            v-if="selectedUsers.length > 0"
            type="success"
            @click="handleBatchStatus(1)"
          >
            批量启用
          </el-button>
          <el-button
            v-if="selectedUsers.length > 0"
            type="warning"
            @click="handleBatchStatus(0)"
          >
            批量禁用
          </el-button>
          <el-button
            v-if="selectedUsers.length > 0"
            type="danger"
            @click="handleBatchDelete"
          >
            批量删除
          </el-button>
        </div>
      </div>

      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="subUsers"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="nickname" label="昵称" width="150">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="200">
          <template #default="{ row }">{{ row.email || '-' }}</template>
        </el-table-column>
        <el-table-column prop="user_status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.user_status === 1 ? 'success' : 'danger'" size="small">
              {{ row.user_status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="关联时间" width="180">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDelete(row)">
              移除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { getSubUsers, deleteSubUser, batchUpdateSubUserStatus, batchDeleteSubUsers } from '@/api/agent'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const subUsers = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')
const selectedUsers = ref([])

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getSubUsers({
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value
    })
    if (res.data?.code === 0) {
      subUsers.value = res.data.data?.list || []
      total.value = res.data.data?.total || 0
    }
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (selection) => {
  selectedUsers.value = selection
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定移除该子用户吗？', '确认', { type: 'warning' })
    await deleteSubUser(row.userId)
    ElMessage.success('已移除')
    loadData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

const handleBatchStatus = async (status) => {
  try {
    await ElMessageBox.confirm(
      `确定${status === 1 ? '启用' : '禁用'}选中的 ${selectedUsers.value.length} 个用户吗？`,
      '确认',
      { type: 'warning' }
    )
    await batchUpdateSubUserStatus(selectedUsers.value.map(u => u.userId), status)
    ElMessage.success('操作成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定批量移除选中的 ${selectedUsers.value.length} 个子用户吗？`,
      '确认',
      { type: 'warning' }
    )
    await batchDeleteSubUsers(selectedUsers.value.map(u => u.userId))
    ElMessage.success('批量移除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

onMounted(loadData)
</script>

<style scoped>
.sub-users-page { padding: 20px; }
.page-header { margin-bottom: 20px; }
.page-header h1 { font-size: 24px; }
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.toolbar-actions { display: flex; gap: 8px; }
.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
