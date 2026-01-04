<template>
  <div class="redemption-codes-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>兑换码管理</span>
          <div>
            <el-button
              type="primary"
              @click="handleAdd"
            >
              添加兑换码
            </el-button>
            <el-button
              type="success"
              @click="handleBatchGenerate"
            >
              批量生成
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form
        :inline="true"
        :model="searchForm"
        class="search-form"
      >
        <el-form-item label="兑换码">
          <el-input
            v-model="searchForm.code"
            placeholder="请输入兑换码"
            clearable
          />
        </el-form-item>
        <el-form-item label="类型">
          <el-select
            v-model="searchForm.type"
            placeholder="请选择类型"
            clearable
          >
            <el-option
              label="余额"
              value="balance"
            />
            <el-option
              label="等级"
              value="level"
            />
            <el-option
              label="产品"
              value="product"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.enabled"
            placeholder="请选择状态"
            clearable
          >
            <el-option
              label="启用"
              :value="true"
            />
            <el-option
              label="禁用"
              :value="false"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            @click="handleSearch"
          >
            搜索
          </el-button>
          <el-button @click="handleReset">
            重置
          </el-button>
        </el-form-item>
      </el-form>

      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
      >
        <el-table-column
          prop="id"
          label="ID"
          width="80"
        />
        <el-table-column
          prop="code"
          label="兑换码"
          width="200"
        />
        <el-table-column
          label="类型"
          width="100"
        >
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)">
              {{ getTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          label="金额/等级"
          width="120"
        >
          <template #default="{ row }">
            {{ row.type === 'balance' ? `¥${(row.amount / 100).toFixed(2)}` : row.amount }}
          </template>
        </el-table-column>
        <el-table-column
          prop="maxUses"
          label="最大使用次数"
          width="120"
        />
        <el-table-column
          prop="usedCount"
          label="已使用"
          width="100"
        />
        <el-table-column
          label="状态"
          width="80"
        >
          <template #default="{ row }">
            <el-tag :type="row.isEnabled ? 'success' : 'danger'">
              {{ row.isEnabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          label="过期时间"
          width="180"
        >
          <template #default="{ row }">
            {{ row.expireAt ? formatTime(row.expireAt) : '永久' }}
          </template>
        </el-table-column>
        <el-table-column
          prop="remark"
          label="备注"
          show-overflow-tooltip
        />
        <el-table-column
          label="操作"
          width="200"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="handleViewUsages(row)"
            >
              使用记录
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="handleToggle(row)"
            >
              {{ row.isEnabled ? '禁用' : '启用' }}
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加兑换码对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="添加兑换码"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item
          label="兑换码"
          prop="code"
        >
          <el-input
            v-model="form.code"
            placeholder="留空自动生成"
          />
          <el-button
            type="primary"
            size="small"
            @click="generateCode"
          >
            生成
          </el-button>
        </el-form-item>
        <el-form-item
          label="类型"
          prop="type"
        >
          <el-select
            v-model="form.type"
            placeholder="请选择类型"
            @change="handleTypeChange"
          >
            <el-option
              label="余额"
              value="balance"
            />
            <el-option
              label="等级"
              value="level"
            />
            <el-option
              label="产品"
              value="product"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="form.type === 'balance'"
          label="金额(元)"
          prop="amount"
        >
          <el-input-number
            v-model="form.amount"
            :min="0"
            :precision="2"
          />
          <span class="ml-2">元</span>
        </el-form-item>
        <el-form-item
          v-if="form.type === 'level'"
          label="等级"
          prop="amount"
        >
          <el-input-number
            v-model="form.amount"
            :min="1"
            :max="5"
          />
        </el-form-item>
        <el-form-item
          v-if="form.type === 'product'"
          label="产品ID"
          prop="productId"
        >
          <el-input-number
            v-model="form.productId"
            :min="1"
          />
        </el-form-item>
        <el-form-item
          label="最大使用次数"
          prop="maxUses"
        >
          <el-input-number
            v-model="form.maxUses"
            :min="1"
          />
        </el-form-item>
        <el-form-item
          label="过期时间"
          prop="expireAt"
        >
          <el-date-picker
            v-model="form.expireAt"
            type="datetime"
            placeholder="选择过期时间"
            clearable
          />
          <span class="ml-2">留空表示永久</span>
        </el-form-item>
        <el-form-item
          label="备注"
          prop="remark"
        >
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="3"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          :loading="saving"
          @click="handleSubmit"
        >
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 批量生成对话框 -->
    <el-dialog
      v-model="batchDialogVisible"
      title="批量生成兑换码"
      width="500px"
    >
      <el-form
        :model="batchForm"
        label-width="120px"
      >
        <el-form-item
          label="生成数量"
          prop="count"
        >
          <el-input-number
            v-model="batchForm.count"
            :min="1"
            :max="1000"
          />
          <span class="ml-2">个</span>
        </el-form-item>
        <el-form-item
          label="类型"
          prop="type"
        >
          <el-select
            v-model="batchForm.type"
            placeholder="请选择类型"
          >
            <el-option
              label="余额"
              value="balance"
            />
            <el-option
              label="等级"
              value="level"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="batchForm.type === 'balance'"
          label="金额(元)"
          prop="amount"
        >
          <el-input-number
            v-model="batchForm.amount"
            :min="0"
            :precision="2"
          />
          <span class="ml-2">元</span>
        </el-form-item>
        <el-form-item
          v-if="batchForm.type === 'level'"
          label="等级"
          prop="amount"
        >
          <el-input-number
            v-model="batchForm.amount"
            :min="1"
            :max="5"
          />
        </el-form-item>
        <el-form-item
          label="最大使用次数"
          prop="maxUses"
        >
          <el-input-number
            v-model="batchForm.maxUses"
            :min="1"
            :value="1"
          />
        </el-form-item>
        <el-form-item
          label="过期时间"
          prop="expireAt"
        >
          <el-date-picker
            v-model="batchForm.expireAt"
            type="datetime"
            placeholder="选择过期时间"
            clearable
          />
        </el-form-item>
        <el-form-item
          label="备注"
          prop="remark"
        >
          <el-input
            v-model="batchForm.remark"
            type="textarea"
            :rows="2"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="batchDialogVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          :loading="saving"
          @click="handleBatchSubmit"
        >
          生成
        </el-button>
      </template>
    </el-dialog>

    <!-- 使用记录对话框 -->
    <el-dialog
      v-model="usagesDialogVisible"
      title="使用记录"
      width="800px"
    >
      <el-table
        v-loading="usagesLoading"
        :data="usagesData"
        border
        stripe
      >
        <el-table-column
          prop="id"
          label="ID"
          width="80"
        />
        <el-table-column
          prop="user.username"
          label="用户名"
          width="150"
        />
        <el-table-column
          prop="user.nickname"
          label="昵称"
          width="150"
        />
        <el-table-column
          label="奖励"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ JSON.stringify(JSON.parse(row.reward), null, 2) }}
          </template>
        </el-table-column>
        <el-table-column
          label="使用时间"
          width="180"
        >
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getRedemptionCodes,
  createRedemptionCode,
  generateRedemptionCodes,
  deleteRedemptionCode,
  getRedemptionCodeUsages,
  toggleRedemptionCode
} from '@/api/admin'

const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const batchDialogVisible = ref(false)
const usagesDialogVisible = ref(false)
const usagesLoading = ref(false)
const tableData = ref([])
const usagesData = ref([])
const formRef = ref()

const searchForm = ref({
  code: '',
  type: '',
  enabled: null
})

const form = ref({
  code: '',
  type: 'balance',
  amount: 0,
  productId: null,
  maxUses: 1,
  expireAt: null,
  remark: ''
})

const batchForm = ref({
  count: 10,
  type: 'balance',
  amount: 0,
  maxUses: 1,
  expireAt: null,
  remark: ''
})

const rules = {
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  amount: [{ required: true, message: '请输入金额或等级', trigger: 'blur' }],
  maxUses: [{ required: true, message: '请输入最大使用次数', trigger: 'blur' }]
}

// 加载兑换码列表
const loadRedemptionCodes = async () => {
  loading.value = true
  try {
    const res = await getRedemptionCodes(searchForm.value)
    if (res.code === 200) {
      tableData.value = res.data || []
    }
  } catch (error) {
    ElMessage.error('加载兑换码列表失败')
  } finally {
    loading.value = false
  }
}

// 获取类型名称
const getTypeName = (type) => {
  const map = {
    balance: '余额',
    level: '等级',
    product: '产品'
  }
  return map[type] || type
}

// 获取类型标签类型
const getTypeTagType = (type) => {
  const map = {
    balance: 'success',
    level: 'warning',
    product: 'info'
  }
  return map[type] || ''
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

// 生成兑换码
const generateCode = () => {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'
  let code = ''
  for (let i = 0; i < 16; i++) {
    code += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  form.value.code = code.substring(0, 4) + '-' +
                  code.substring(4, 8) + '-' +
                  code.substring(8, 12) + '-' +
                  code.substring(12, 16)
}

// 类型变化
const handleTypeChange = () => {
  form.value.amount = 0
  form.value.productId = null
}

// 添加兑换码
const handleAdd = () => {
  form.value = {
    code: '',
    type: 'balance',
    amount: 0,
    productId: null,
    maxUses: 1,
    expireAt: null,
    remark: ''
  }
  dialogVisible.value = true
}

// 批量生成
const handleBatchGenerate = () => {
  batchForm.value = {
    count: 10,
    type: 'balance',
    amount: 0,
    maxUses: 1,
    expireAt: null,
    remark: ''
  }
  batchDialogVisible.value = true
}

// 提交添加
const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        const submitData = {
          ...form.value,
          amount: form.value.type === 'balance' ? Math.round(form.value.amount * 100) : form.value.amount
        }

        const res = await createRedemptionCode(submitData)
        if (res.code === 200) {
          ElMessage.success('添加成功')
          dialogVisible.value = false
          loadRedemptionCodes()
        } else {
          ElMessage.error(res.message || '添加失败')
        }
      } catch (error) {
        ElMessage.error('添加失败')
      } finally {
        saving.value = false
      }
    }
  })
}

// 批量生成提交
const handleBatchSubmit = async () => {
  saving.value = true
  try {
    const submitData = {
      ...batchForm.value,
      amount: batchForm.value.type === 'balance' ? Math.round(batchForm.value.amount * 100) : batchForm.value.amount
    }

    const res = await generateRedemptionCodes(submitData)
    if (res.code === 200) {
      ElMessage.success(`成功生成${batchForm.value.count}个兑换码`)
      batchDialogVisible.value = false
      loadRedemptionCodes()
    } else {
      ElMessage.error(res.message || '生成失败')
    }
  } catch (error) {
    ElMessage.error('生成失败')
  } finally {
    saving.value = false
  }
}

// 查看使用记录
const handleViewUsages = async (row) => {
  usagesDialogVisible.value = true
  usagesLoading.value = true
  try {
    const res = await getRedemptionCodeUsages(row.id)
    if (res.code === 200) {
      usagesData.value = res.data || []
    }
  } catch (error) {
    ElMessage.error('加载使用记录失败')
  } finally {
    usagesLoading.value = false
  }
}

// 启用/禁用
const handleToggle = async (row) => {
  try {
    const res = await toggleRedemptionCode(row.id)
    if (res.code === 200) {
      ElMessage.success('操作成功')
      loadRedemptionCodes()
    } else {
      ElMessage.error(res.message || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该兑换码吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteRedemptionCode(row.id)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        loadRedemptionCodes()
      } else {
        ElMessage.error(res.message || '删除失败')
      }
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

// 搜索
const handleSearch = () => {
  loadRedemptionCodes()
}

// 重置
const handleReset = () => {
  searchForm.value = {
    code: '',
    type: '',
    enabled: null
  }
  loadRedemptionCodes()
}

// 对话框关闭
const handleDialogClose = () => {
  formRef.value?.resetFields()
}

onMounted(() => {
  loadRedemptionCodes()
})
</script>

<style scoped>
.redemption-codes-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.ml-2 {
  margin-left: 8px;
}
</style>
