<template>
  <div class="apply-page">
    <div class="page-header">
      <el-page-header @back="$router.back()" :content="t('user.apply.title')"></el-page-header>
    </div>

    <el-card class="config-card">
      <template #header>
        <span>{{ t('user.apply.configTitle') }}</span>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <!-- 实例名称 -->
        <el-form-item :label="t('user.apply.name')" prop="name">
          <el-input v-model="form.name" :placeholder="t('user.apply.namePlaceholder')" />
        </el-form-item>

        <!-- 节点选择 -->
        <el-form-item label="节点" prop="providerId">
          <el-select v-model="form.providerId" :placeholder="t('user.apply.providerPlaceholder')" style="width: 100%" @change="onProviderChange">
            <el-option
              v-for="p in providers"
              :key="p.id"
              :label="`${p.name} (${p.type})`"
              :value="p.id"
            >
              <div style="display: flex; justify-content: space-between;">
                <span>{{ p.name }}</span>
                <el-tag size="small" :type="p.status === 'active' ? 'success' : 'danger'">{{ p.type }}</el-tag>
              </div>
            </el-option>
          </el-select>
        </el-form-item>

        <!-- 镜像选择 -->
        <el-form-item label="镜像" prop="imageId">
          <el-select v-model="form.imageId" :placeholder="t('user.apply.imagePlaceholder')" style="width: 100%">
            <el-option
              v-for="img in images"
              :key="img.id"
              :label="img.name"
              :value="img.id"
            >
              <div style="display: flex; justify-content: space-between;">
                <span>{{ img.name }}</span>
                <el-tag size="small">{{ img.providerType }}/{{ img.instanceType }}</el-tag>
              </div>
            </el-option>
          </el-select>
        </el-form-item>

        <!-- 计算资源 -->
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="CPU" prop="cpuId">
              <el-select v-model="form.cpuId" placeholder="请选择CPU规格" style="width: 100%">
                <el-option
                  v-for="spec in cpuSpecs"
                  :key="spec.id"
                  :label="`${spec.name} (${spec.cores}核)`"
                  :value="spec.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="内存" prop="memoryId">
              <el-select v-model="form.memoryId" placeholder="请选择内存规格" style="width: 100%">
                <el-option
                  v-for="spec in memorySpecs"
                  :key="spec.id"
                  :label="`${spec.name} (${formatMemory(spec.sizeMB)})`"
                  :value="spec.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="磁盘" prop="diskId">
              <el-select v-model="form.diskId" placeholder="请选择磁盘规格" style="width: 100%">
                <el-option
                  v-for="spec in diskSpecs"
                  :key="spec.id"
                  :label="`${spec.name} (${formatDisk(spec.sizeMB)})`"
                  :value="spec.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="带宽" prop="bandwidthId">
              <el-select v-model="form.bandwidthId" placeholder="请选择带宽规格" style="width: 100%">
                <el-option
                  v-for="spec in bandwidthSpecs"
                  :key="spec.id"
                  :label="`${spec.name} (${spec.speedMbps}Mbps)`"
                  :value="spec.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 描述 -->
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入实例描述（可选）"
          />
        </el-form-item>

        <!-- 操作 -->
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            {{ t('user.apply.submit') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getAvailableProviders, getSystemImages, createInstance } from '@/api/user'
import { getSpecs } from '@/api/public'

const { t } = useI18n()
const router = useRouter()

// 表单数据
const form = ref({
  name: '',
  providerId: null,
  imageId: null,
  cpuId: '',
  memoryId: '',
  diskId: '',
  bandwidthId: '',
  description: ''
})

const rules = {
  name: [
    { required: true, message: t('user.apply.nameRequired'), trigger: 'blur' }
  ],
  providerId: [
    { required: true, message: '请选择节点', trigger: 'change' }
  ],
  imageId: [
    { required: true, message: '请选择镜像', trigger: 'change' }
  ],
  cpuId: [
    { required: true, message: '请选择CPU规格', trigger: 'change' }
  ],
  memoryId: [
    { required: true, message: '请选择内存规格', trigger: 'change' }
  ],
  diskId: [
    { required: true, message: '请选择磁盘规格', trigger: 'change' }
  ],
  bandwidthId: [
    { required: true, message: '请选择带宽规格', trigger: 'change' }
  ]
}

// 数据
const providers = ref([])
const images = ref([])
const cpuSpecs = ref([])
const memorySpecs = ref([])
const diskSpecs = ref([])
const bandwidthSpecs = ref([])

// 格式化内存
const formatMemory = (mb) => {
  if (mb >= 1024) {
    return `${(mb / 1024).toFixed(0)}GB`
  }
  return `${mb}MB`
}

// 格式化磁盘
const formatDisk = (mb) => {
  if (mb >= 1024) {
    return `${(mb / 1024).toFixed(0)}GB`
  }
  return `${mb}MB`
}

// 加载节点列表
const loadProviders = async () => {
  try {
    const res = await getAvailableProviders()
    providers.value = res.data || []
  } catch (error) {
    console.error('Failed to load providers:', error)
  }
}

// 加载镜像列表
const loadImages = async () => {
  try {
    const res = await getSystemImages()
    images.value = res.data.images || res.data.imagesList || []
  } catch (error) {
    console.error('Failed to load images:', error)
  }
}

// 加载规格配置
const loadSpecs = async () => {
  try {
    const res = await getSpecs()
    cpuSpecs.value = res.data.cpu || []
    memorySpecs.value = res.data.memory || []
    diskSpecs.value = res.data.disk || []
    bandwidthSpecs.value = res.data.bandwidth || []
  } catch (error) {
    console.error('Failed to load specs:', error)
  }
}

// 节点变化时过滤镜像
const onProviderChange = () => {
  // 可以在此处根据节点类型过滤镜像
}

// 提交表单
const submitForm = async () => {
  try {
    if (!formRef.value.validate) throw new Error('Form not initialized')
    await formRef.value.validate()

    submitting.value = true

    const res = await createInstance({
      providerId: form.value.providerId,
      imageId: form.value.imageId,
      cpuId: form.value.cpuId,
      memoryId: form.value.memoryId,
      diskId: form.value.diskId,
      bandwidthId: form.value.bandwidthId,
      description: form.value.description
    })

    if (res.code === 200 || res.code === 0) {
      ElMessage.success(t('user.apply.success'))
      router.push('/user/instances')
    } else {
      ElMessage.error(res.message || res.msg || t('user.apply.failed'))
    }
  } catch (error) {
    if (!error.message.includes('validate')) {
      ElMessage.error(t('user.apply.error'))
    }
  } finally {
    submitting.value = false
  }
}

const formRef = ref(null)
const submitting = ref(false)

onMounted(() => {
  loadProviders()
  loadImages()
  loadSpecs()
})
</script>

<style scoped>
.apply-page {
  padding: 20px;
}

.config-card {
  margin-top: 20px;
}

:deep(.el-input-number) {
  width: 100%;
}
</style>
