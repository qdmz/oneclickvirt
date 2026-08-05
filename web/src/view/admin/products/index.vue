<template>
  <div class="admin-products">
    <div class="page-header">
      <div>
        <h1>{{ t('admin.products.title') }}</h1>
        <p>{{ t('admin.products.subtitle') }}</p>
      </div>
      <div class="header-right">
        <el-button
          :icon="Refresh"
          @click="loadList"
        >
          {{ t('common.refresh') }}
        </el-button>
        <el-button
          type="primary"
          :icon="Plus"
          @click="openCreate"
        >
          {{ t('admin.products.create') }}
        </el-button>
      </div>
    </div>

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
          prop="name"
          :label="t('admin.products.colName')"
          min-width="160"
          show-overflow-tooltip
        />
        <el-table-column
          :label="t('admin.products.colType')"
          width="110"
        >
          <template #default="{ row }">
            <el-tag
              size="small"
              :type="typeTagType(row.type)"
              effect="light"
            >
              {{ t('admin.products.typeLabel.' + row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.products.colSpec')"
          min-width="220"
        >
          <template #default="{ row }">
            {{ row.cpu }}C / {{ formatMB(row.memoryMB) }} / {{ formatMB(row.diskMB) }} / {{ row.bandwidthMbps }}Mbps
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.products.colPrice')"
          width="120"
        >
          <template #default="{ row }">
            <span class="price">¥{{ fenToYuan(row.price) }}</span>
            <span class="period">/{{ row.periodDays }}d</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.products.colStock')"
          width="110"
        >
          <template #default="{ row }">
            <span v-if="row.stock < 0">{{ t('admin.products.unlimited') }}</span>
            <span v-else>{{ row.sold || 0 }} / {{ row.stock }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('admin.products.colStatus')"
          width="110"
        >
          <template #default="{ row }">
            <el-switch
              :model-value="row.status === 1"
              @change="(v) => toggleStatus(row, v)"
            />
          </template>
        </el-table-column>
        <el-table-column
          prop="sort"
          :label="t('admin.products.colSort')"
          width="80"
        />
        <el-table-column
          :label="t('common.actions')"
          width="150"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              @click="openEdit(row)"
            >
              {{ t('common.edit') }}
            </el-button>
            <el-popconfirm
              :title="t('admin.products.deleteConfirm')"
              @confirm="doDelete(row)"
            >
              <template #reference>
                <el-button
                  link
                  type="danger"
                >
                  {{ t('common.delete') }}
                </el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>

    <!-- 新建/编辑 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? t('admin.products.editTitle') : t('admin.products.createTitle')"
      width="640px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="110px"
      >
        <el-form-item
          :label="t('admin.products.fName')"
          prop="name"
        >
          <el-input
            v-model="form.name"
            :placeholder="t('admin.products.fNamePh')"
          />
        </el-form-item>
        <el-form-item :label="t('admin.products.fDescription')">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="2"
            :placeholder="t('admin.products.fDescriptionPh')"
          />
        </el-form-item>
        <el-form-item
          :label="t('admin.products.fType')"
          prop="type"
        >
          <el-radio-group v-model="form.type">
            <el-radio-button value="vm">
              {{ t('admin.products.typeLabel.vm') }}
            </el-radio-button>
            <el-radio-button value="lxc">
              {{ t('admin.products.typeLabel.lxc') }}
            </el-radio-button>
            <el-radio-button value="container">
              {{ t('admin.products.typeLabel.container') }}
            </el-radio-button>
          </el-radio-group>
        </el-form-item>

        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item
              :label="t('admin.products.fCpu')"
              prop="cpu"
            >
              <el-input-number
                v-model="form.cpu"
                :min="1"
                :max="128"
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="t('admin.products.fMemory')"
              prop="memoryMB"
            >
              <el-input-number
                v-model="form.memoryMB"
                :min="128"
                :step="512"
                :max="524288"
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item
              :label="t('admin.products.fDisk')"
              prop="diskMB"
            >
              <el-input-number
                v-model="form.diskMB"
                :min="1024"
                :step="10240"
                :max="20971520"
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="t('admin.products.fBandwidth')"
              prop="bandwidthMbps"
            >
              <el-input-number
                v-model="form.bandwidthMbps"
                :min="1"
                :max="100000"
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item
              :label="t('admin.products.fPrice')"
              prop="priceYuan"
            >
              <el-input-number
                v-model="form.priceYuan"
                :min="0"
                :precision="2"
                :step="1"
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="t('admin.products.fPeriod')"
              prop="periodDays"
            >
              <el-input-number
                v-model="form.periodDays"
                :min="1"
                :max="3650"
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item :label="t('admin.products.fStock')">
              <el-input-number
                v-model="form.stock"
                :min="-1"
                style="width: 100%;"
              />
              <div class="form-tip">
                {{ t('admin.products.fStockTip') }}
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('admin.products.fSort')">
              <el-input-number
                v-model="form.sort"
                :min="0"
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item :label="t('admin.products.fNodeId')">
              <el-input-number
                v-model="form.nodeId"
                :min="0"
                style="width: 100%;"
              />
              <div class="form-tip">
                {{ t('admin.products.fNodeIdTip') }}
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('admin.products.fStatus')">
              <el-switch
                v-model="form.statusOn"
                :active-text="t('admin.products.onShelf')"
                :inactive-text="t('admin.products.offShelf')"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="submitForm"
        >
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import {
  getAdminProducts, createAdminProduct, updateAdminProduct, deleteAdminProduct
} from '@/api/admin'
import { fenToYuan } from '@/utils/epay'

const { t } = useI18n()

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const typeTagType = (ty) => ({ vm: 'primary', lxc: 'success', container: 'warning' }[ty] || 'info')

const formatMB = (mb) => {
  const n = Number(mb || 0)
  if (n >= 1024) {
    const g = n / 1024
    return `${Number.isInteger(g) ? g : g.toFixed(1)}GB`
  }
  return `${n}MB`
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await getAdminProducts({ page: page.value, pageSize: pageSize.value })
    if (res && res.code === 200) {
      list.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch {
    // interceptor handles
  } finally {
    loading.value = false
  }
}

// 表单
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref(0)
const formRef = ref()
const submitting = ref(false)

const defaultForm = () => ({
  name: '',
  description: '',
  type: 'vm',
  cpu: 1,
  memoryMB: 1024,
  diskMB: 10240,
  bandwidthMbps: 100,
  priceYuan: 9.9,
  periodDays: 30,
  stock: -1,
  sort: 0,
  nodeId: 0,
  statusOn: true
})

const form = reactive(defaultForm())

const rules = {
  name: [{ required: true, message: t('admin.products.nameRequired'), trigger: 'blur' }],
  cpu: [{ required: true, message: t('admin.products.required'), trigger: 'blur' }],
  memoryMB: [{ required: true, message: t('admin.products.required'), trigger: 'blur' }],
  diskMB: [{ required: true, message: t('admin.products.required'), trigger: 'blur' }],
  periodDays: [{ required: true, message: t('admin.products.required'), trigger: 'blur' }]
}

const resetForm = (data) => {
  const d = data || defaultForm()
  Object.keys(defaultForm()).forEach((k) => { form[k] = d[k] })
}

const openCreate = () => {
  isEdit.value = false
  editingId.value = 0
  resetForm(defaultForm())
  dialogVisible.value = true
}

const openEdit = (row) => {
  isEdit.value = true
  editingId.value = row.id
  resetForm({
    name: row.name,
    description: row.description || '',
    type: row.type || 'vm',
    cpu: row.cpu,
    memoryMB: row.memoryMB,
    diskMB: row.diskMB,
    bandwidthMbps: row.bandwidthMbps,
    priceYuan: Number(fenToYuan(row.price)),
    periodDays: row.periodDays,
    stock: row.stock,
    sort: row.sort || 0,
    nodeId: row.nodeId || 0,
    statusOn: row.status === 1
  })
  dialogVisible.value = true
}

const buildPayload = () => ({
  name: form.name,
  description: form.description,
  type: form.type,
  cpu: form.cpu,
  memoryMB: form.memoryMB,
  diskMB: form.diskMB,
  bandwidthMbps: form.bandwidthMbps,
  price: Math.round(form.priceYuan * 100),
  periodDays: form.periodDays,
  stock: form.stock,
  sort: form.sort,
  nodeId: form.nodeId,
  status: form.statusOn ? 1 : 0
})

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const payload = buildPayload()
      const res = isEdit.value
        ? await updateAdminProduct(editingId.value, payload)
        : await createAdminProduct(payload)
      if (res && res.code === 200) {
        ElMessage.success(isEdit.value ? t('admin.products.updateSuccess') : t('admin.products.createSuccess'))
        dialogVisible.value = false
        loadList()
      }
    } finally {
      submitting.value = false
    }
  })
}

const toggleStatus = async (row, on) => {
  try {
    const res = await updateAdminProduct(row.id, { status: on ? 1 : 0 })
    if (res && res.code === 200) {
      row.status = on ? 1 : 0
      ElMessage.success(t('admin.products.statusUpdated'))
    }
  } catch {
    // interceptor handles
  }
}

const doDelete = async (row) => {
  try {
    const res = await deleteAdminProduct(row.id)
    if (res && res.code === 200) {
      ElMessage.success(t('admin.products.deleteSuccess'))
      loadList()
    }
  } catch {
    // interceptor handles
  }
}

onMounted(loadList)
</script>

<style scoped>
.admin-products {
  padding: 20px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
  gap: 12px;
  flex-wrap: wrap;
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
.header-right {
  display: flex;
  gap: 10px;
}
.table-card {
  border: 1px solid var(--el-border-color-lighter);
}
.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
.price {
  font-weight: 600;
  color: var(--el-color-danger);
}
.period {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-left: 2px;
}
.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.4;
  margin-top: 2px;
}
</style>
