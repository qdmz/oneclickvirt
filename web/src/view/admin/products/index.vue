<template>
  <div class="products-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>产品管理</span>
          <el-button
            type="primary"
            @click="handleAdd"
          >
            添加产品
          </el-button>
        </div>
      </template>

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
          prop="name"
          label="产品名称"
        />
        <el-table-column
          prop="level"
          label="等级"
          width="80"
        />
        <el-table-column
          label="价格"
          width="120"
        >
          <template #default="{ row }">
            ¥{{ (row.price / 100).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="period"
          label="有效期"
          width="100"
        >
          <template #default="{ row }">
            {{ row.period === 0 ? '永久' : row.period + '天' }}
          </template>
        </el-table-column>
        <el-table-column
          label="配置"
          width="300"
        >
          <template #default="{ row }">
            {{ row.cpu }}核 / {{ row.memory }}MB / {{ row.disk }}GB / {{ row.bandwidth }}Mbps
          </template>
        </el-table-column>
        <el-table-column
          label="流量"
          width="100"
        >
          <template #default="{ row }">
            {{ formatTraffic(row.traffic) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="maxInstances"
          label="最大实例"
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
          prop="sortOrder"
          label="排序"
          width="80"
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
              @click="handleEdit(row)"
            >
              编辑
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

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="800px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item
          label="产品名称"
          prop="name"
        >
          <el-input
            v-model="form.name"
            placeholder="请输入产品名称"
          />
        </el-form-item>
        <el-form-item
          label="产品描述"
          prop="description"
        >
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入产品描述"
          />
        </el-form-item>
        <el-form-item
          label="等级"
          prop="level"
        >
          <el-input-number
            v-model="form.level"
            :min="1"
            :max="5"
          />
        </el-form-item>
        <el-form-item
          label="价格(元)"
          prop="price"
        >
          <el-input-number
            v-model="form.price"
            :min="0"
            :precision="2"
          />
          <span class="ml-2">元</span>
        </el-form-item>
        <el-form-item
          label="有效期(天)"
          prop="period"
        >
          <el-input-number
            v-model="form.period"
            :min="0"
          />
          <span class="ml-2">天, 0表示永久</span>
        </el-form-item>
        <el-divider content-position="left">
          资源配置
        </el-divider>
        <div class="auto-fill-tip">
          <el-alert
            type="info"
            :closable="false"
            show-icon
          >
            资源配置会根据产品等级自动从用户等级限制中填充，您也可以手动修改
          </el-alert>
        </div>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              label="CPU核心数"
              prop="cpu"
            >
              <el-input-number
                v-model="form.cpu"
                :min="1"
              />
              <span class="ml-2">核</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              label="内存"
              prop="memory"
            >
              <el-input-number
                v-model="form.memory"
                :min="128"
                :step="128"
              />
              <span class="ml-2">MB</span>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              label="磁盘空间"
              prop="disk"
            >
              <el-input-number
                v-model="form.disk"
                :min="1"
              />
              <span class="ml-2">MB</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              label="带宽"
              prop="bandwidth"
            >
              <el-input-number
                v-model="form.bandwidth"
                :min="10"
              />
              <span class="ml-2">Mbps</span>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              label="流量配额"
              prop="traffic"
            >
              <el-input-number
                v-model="form.traffic"
                :min="0"
                :step="1024"
              />
              <span class="ml-2">MB</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              label="最大实例数"
              prop="maxInstances"
            >
              <el-input-number
                v-model="form.maxInstances"
                :min="1"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              label="排序"
              prop="sortOrder"
            >
              <el-input-number
                v-model="form.sortOrder"
                :min="0"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              label="状态"
              prop="isEnabled"
            >
              <el-switch v-model="form.isEnabled" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              label="是否允许重复购买"
              prop="allowRepeat"
            >
              <el-switch v-model="form.allowRepeat" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="特性(JSON)">
          <el-input
            v-model="form.features"
            type="textarea"
            :rows="3"
            placeholder="[&quot;特性1&quot;, &quot;特性2&quot;]"
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getProducts,
  createProduct,
  updateProduct,
  deleteProduct,
  toggleProduct
} from '@/api/admin'

const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('添加产品')
const tableData = ref([])
const formRef = ref()

const form = ref({
  name: '',
  description: '',
  level: 1,
  price: 0,
  period: 30,
  cpu: 1,
  memory: 512,
  disk: 10240,
  bandwidth: 100,
  traffic: 102400,
  maxInstances: 1,
  isEnabled: true,
  sortOrder: 0,
  allowRepeat: true,
  features: '[]'
})

const rules = {
  name: [{ required: true, message: '请输入产品名称', trigger: 'blur' }],
  level: [{ required: true, message: '请输入等级', trigger: 'blur' }],
  price: [{ required: true, message: '请输入价格', trigger: 'blur' }],
  cpu: [{ required: true, message: '请输入CPU核心数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存', trigger: 'blur' }],
  disk: [{ required: true, message: '请输入磁盘空间', trigger: 'blur' }],
  bandwidth: [{ required: true, message: '请输入带宽', trigger: 'blur' }],
  traffic: [{ required: true, message: '请输入流量配额', trigger: 'blur' }],
  maxInstances: [{ required: true, message: '请输入最大实例数', trigger: 'blur' }]
}

// 加载产品列表
const loadProducts = async () => {
  loading.value = true
  try {
    const res = await getProducts()
    console.log('产品列表API响应:', res)
    if (res.code === 200 || res.code === 0) {
      tableData.value = res.data || []
      ElMessage.success(`加载成功，共 ${tableData.value.length} 个产品`)
      
      // 检查是否每个等级都有产品
      const levels = new Set(tableData.value.map(p => p.level))
      console.log('当前产品等级分布:', Array.from(levels))
      
      // 如果有等级没有产品，提示管理员
      for (let i = 1; i <= 5; i++) {
        if (!levels.has(i)) {
          console.warn(`等级 ${i} 没有对应的产品，请创建`)
        }
      }
    } else {
      console.error('加载产品列表失败，响应码:', res.code, '消息:', res.message)
      ElMessage.error(res.message || '加载产品列表失败')
    }
  } catch (error) {
    console.error('加载产品列表失败:', error)
    console.error('错误详情:', error.message)
    console.error('错误栈:', error.stack)
    ElMessage.error('加载产品列表失败')
  } finally {
    loading.value = false
  }
}

// 格式化流量
const formatTraffic = (traffic) => {
  if (traffic < 1024) {
    return traffic + 'MB'
  } else if (traffic < 1024 * 1024) {
    return (traffic / 1024).toFixed(2) + 'GB'
  } else {
    return (traffic / 1024 / 1024).toFixed(2) + 'TB'
  }
}

// 添加产品
const handleAdd = () => {
  dialogTitle.value = '添加产品'
  form.value = {
    name: '',
    description: '',
    level: 1,
    price: 0,
    period: 30,
    cpu: 1,
    memory: 512,
    disk: 10240,
    bandwidth: 100,
    traffic: 102400,
    maxInstances: 1,
    isEnabled: true,
    sortOrder: 0,
    features: '[]'
  }
  dialogVisible.value = true
}

// 编辑产品
const handleEdit = (row) => {
  dialogTitle.value = '编辑产品'
  form.value = {
    ...row,
    price: row.price / 100,
    isEnabled: row.isEnabled === 1 || row.isEnabled === true, // 将整数或布尔值转换为布尔值
    allowRepeat: row.allowRepeat === 1 || row.allowRepeat === true // 将整数或布尔值转换为布尔值
  }
  dialogVisible.value = true
}

// 启用/禁用产品
const handleToggle = async (row) => {
  try {
    console.log('开始切换产品状态，产品ID:', row.id, '当前状态:', row.isEnabled)
    const res = await toggleProduct(row.id)
    console.log('切换产品状态API响应:', res)
    if (res.code === 200) {
      ElMessage.success('操作成功')
      loadProducts()
    } else {
      console.error('切换产品状态失败，响应码:', res.code, '消息:', res.message)
      ElMessage.error(res.message || '操作失败')
    }
  } catch (error) {
    console.error('切换产品状态失败:', error)
    console.error('错误详情:', error.message)
    console.error('错误栈:', error.stack)
    ElMessage.error('操作失败')
  }
}

// 删除产品
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该产品吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      console.log('开始删除产品，产品ID:', row.id)
      const res = await deleteProduct(row.id)
      console.log('删除产品API响应:', res)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        loadProducts()
      } else {
        console.error('删除产品失败，响应码:', res.code, '消息:', res.message)
        ElMessage.error(res.message || '删除失败')
      }
    } catch (error) {
      console.error('删除产品失败:', error)
      console.error('错误详情:', error.message)
      console.error('错误栈:', error.stack)
      ElMessage.error('删除失败')
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        const submitData = {
          ...form.value,
          price: Math.round(form.value.price * 100), // 转换为分
          isEnabled: form.value.isEnabled ? 1 : 0, // 将布尔值转换为整数(1:启用, 0:禁用)
          allowRepeat: form.value.allowRepeat ? 1 : 0 // 将布尔值转换为整数(1:允许, 0:不允许)
        }

        let res
        console.log(`开始${dialogTitle.value}操作，数据:`, submitData)
        if (dialogTitle.value === '添加产品') {
          res = await createProduct(submitData)
        } else {
          res = await updateProduct(form.value.id, submitData)
        }
        console.log(`${dialogTitle.value}产品API响应:`, res)

        if (res.code === 200) {
          ElMessage.success(dialogTitle.value === '添加产品' ? '添加成功' : '更新成功')
          dialogVisible.value = false
          loadProducts()
        } else {
          console.error(`${dialogTitle.value}产品失败，响应码:`, res.code, '消息:', res.message)
          ElMessage.error(res.message || '操作失败')
        }
      } catch (error) {
        console.error(`${dialogTitle.value}产品失败:`, error)
        console.error('错误详情:', error.message)
        console.error('错误栈:', error.stack)
        ElMessage.error('操作失败')
      } finally {
        saving.value = false
      }
    }
  })
}

// 对话框关闭
const handleDialogClose = () => {
  formRef.value?.resetFields()
}

onMounted(() => {
  loadProducts()
})
</script>

<style scoped>
.products-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.ml-2 {
  margin-left: 8px;
}

.auto-fill-tip {
  margin-bottom: 20px;
}
</style>
