<template>
  <div class="site-link-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>站点链接管理</span>
          <el-button type="primary" @click="handleAdd">
            新增链接
          </el-button>
        </div>
      </template>

      <!-- 筛选区域 -->
      <el-form :inline="true" :model="queryParams">
        <el-form-item label="链接类型">
          <el-select v-model="queryParams.linkType" placeholder="全部" clearable>
            <el-option label="虚拟化平台" value="platform" />
            <el-option label="赞助方" value="sponsor" />
            <el-option label="页脚友链" value="footer" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="全部" clearable>
            <el-option label="显示" :value="1" />
            <el-option label="隐藏" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <el-table :data="tableData" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" min-width="120" />
        <el-table-column prop="linkType" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getLinkTypeTagType(row.linkType)">
              {{ getLinkTypeLabel(row.linkType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="url" label="链接地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="iconUrl" label="图标" width="80">
          <template #default="{ row }">
            <el-image
              v-if="row.iconUrl"
              :src="row.iconUrl"
              style="width: 40px; height: 40px"
              fit="cover"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="sortOrder" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '显示' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="queryParams.page"
        v-model:page-size="queryParams.pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadData"
        @current-change="loadData"
        class="pagination"
      />
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? '编辑链接' : '新增链接'"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入名称" maxlength="128" show-word-limit />
        </el-form-item>

        <el-form-item label="链接类型" prop="linkType">
          <el-select v-model="form.linkType" placeholder="请选择类型">
            <el-option label="虚拟化平台" value="platform" />
            <el-option label="赞助方" value="sponsor" />
            <el-option label="页脚友链" value="footer" />
          </el-select>
        </el-form-item>

        <el-form-item label="链接地址" prop="url">
          <el-input v-model="form.url" placeholder="请输入链接地址" maxlength="512" />
        </el-form-item>

        <el-form-item label="图标URL" prop="iconUrl">
          <el-input v-model="form.iconUrl" placeholder="请输入图标URL" maxlength="512" />
        </el-form-item>

        <el-form-item label="排序" prop="sortOrder">
          <el-input-number v-model="form.sortOrder" :min="0" :max="999" />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">显示</el-radio>
            <el-radio :label="0">隐藏</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入描述" maxlength="256" show-word-limit />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAdminSiteLinkList,
  createAdminSiteLink,
  updateAdminSiteLink,
  deleteAdminSiteLink
} from '@/api/admin/ecommerce'

const loading = ref(false)
const submitting = ref(false)
const tableData = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const formRef = ref(null)

const queryParams = reactive({
  linkType: '',
  status: null,
  page: 1,
  pageSize: 20
})

const form = reactive({
  id: null,
  name: '',
  url: '',
  iconUrl: '',
  linkType: 'platform',
  sortOrder: 0,
  status: 1,
  description: ''
})

const rules = {
  name: [
    { required: true, message: '请输入名称', trigger: 'blur' }
  ],
  url: [
    { required: true, message: '请输入链接地址', trigger: 'blur' }
  ],
  linkType: [
    { required: true, message: '请选择链接类型', trigger: 'change' }
  ]
}

const linkTypeLabelMap = {
  platform: '虚拟化平台',
  sponsor: '赞助方',
  footer: '页脚友链'
}
const getLinkTypeLabel = (type) => linkTypeLabelMap[type] || type || '未知'
const getLinkTypeTagType = (type) => {
  if (type === 'platform') return 'primary'
  if (type === 'sponsor') return 'success'
  if (type === 'footer') return 'warning'
  return 'info'
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await getAdminSiteLinkList(queryParams)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error) {
    console.error('加载数据失败', error)
  } finally {
    loading.value = false
  }
}

// 重置查询
const handleReset = () => {
  queryParams.linkType = ''
  queryParams.status = null
  queryParams.page = 1
  queryParams.pageSize = 20
  loadData()
}

// 新增
const handleAdd = () => {
  Object.assign(form, {
    id: null,
    name: '',
    url: '',
    iconUrl: '',
    linkType: 'platform',
    sortOrder: 0,
    status: 1,
    description: ''
  })
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  Object.assign(form, {
    id: row.id,
    name: row.name,
    url: row.url,
    iconUrl: row.iconUrl || '',
    linkType: row.linkType,
    sortOrder: row.sortOrder,
    status: row.status,
    description: row.description || ''
  })
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该链接吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteAdminSiteLink(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败', error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (form.id) {
        await updateAdminSiteLink(form.id, form)
        ElMessage.success('更新成功')
      } else {
        await createAdminSiteLink(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadData()
    } catch (error) {
      console.error('提交失败', error)
    } finally {
      submitting.value = false
    }
  })
}

// 关闭对话框
const handleDialogClose = () => {
  formRef.value?.resetFields()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.site-link-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
