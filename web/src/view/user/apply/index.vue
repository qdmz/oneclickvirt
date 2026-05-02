<template>
  <div class="apply-page">
    <div class="page-header">
      <el-page-header @back="$router.back()" :content="t('user.apply.title')"></el-page-header>
    </div>

    <el-card class="config-card">
      <template #header>
        <span>{{ t('user.apply.configTitle') }}</span>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="150px">
        <!-- 实例名称 -->
        <el-form-item :label="t('user.apply.name')" prop="name">
          <el-input v-model="form.name" :placeholder="t('user.apply.namePlaceholder')" />
        </el-form-item>

        <!-- 虚拟机类型 -->
        <el-form-item :label="t('user.apply.type')" prop="instanceType">
          <el-radio-group v-model="form.instanceType">
            <el-radio value="vm">
              <el-icon><Monitor /></el-icon>
              {{ t('user.apply.vm') }}
            </el-radio>
            <el-radio value="container">
              <el-icon><Box /></el-icon>
              {{ t('user.apply.container') }}
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 计算 - CPU和内存 -->
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="t('user.apply.cpu')" prop="cpu">
              <el-input-number
                v-model="form.cpu"
                :min="1"
                :max="64"
                :step="1"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('user.apply.memory')" prop="memory">
              <el-input-number
                v-model="form.memory"
                :min="1"
                :max="128"
                :step="1"
                :unit="t('user.apply.gb')"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 存储 -->
        <el-form-item :label="t('user.apply.disk')" prop="disk">
          <el-input-number
            v-model="form.disk"
            :min="10"
            :max="2048"
            :step="10"
          />
          <span class="mb-2">GB</span>
        </el-form-item>

        <!-- 镜像选择 -->
        <el-form-item :label="t('user.apply.image')" prop="imageId">
          <el-select v-model="form.imageId" :placeholder="t('user.apply.imagePlaceholder')" style="width: 100%">
            <el-option
              v-for="img in images"
              :key="img.id"
              :label="img.name"
              :value="img.id"
            >
              <div style="display: flex; justify-content: space-between;">
                <span>{{ img.name }}</span>
                <el-tag size="small">{{ img.os }}</el-tag>
              </div>
            </el-option>
          </el-select>
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

const { t } = useI18n()
const router = useRouter()

// 表单数据和验证规则
const form = ref({
  name: '',
  instanceType: 'vm',
  cpu: 1,
  memory: 1,
  disk: 50,
  imageId: ''
})

const rules = {
  name: [
    { required: true, message: t('user.apply.nameRequired'), trigger: 'blur' }
  ],
  instanceType: [
    { required: true, message: t('user.apply.typeRequired'), trigger: 'change' }
  ],
  cpu: [
    { required: true, message: t('user.apply.cpuRequired'), trigger: 'blur' },
    { type: 'number', min: 1, max: 64, message: t('user.apply.cpuRange'), trigger: 'blur' }
  ],
  memory: [
    { required: true, message: t('user.apply.memoryRequired'), trigger: 'blur' },
    { type: 'number', min: 1, max: 128, message: t('user.apply.memoryRange'), trigger: 'blur' }
  ],
  disk: [
    { required: true, message: t('user.apply.diskRequired'), trigger: 'blur' },
    { type: 'number', min: 10, max: 2048, message: t('user.apply.diskRange'), trigger: 'blur' }
  ],
  imageId: [
    { required: true, message: t('user.apply.imageRequired'), trigger: 'change' }
  ]
}

// 加载系统镜像
const images = ref([])

const loadImages = async () => {
  try {
    const response = await fetch('/v1/user/images', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    const data = await response.json()
    images.value = data.images || data.imagesList || []
  } catch (error) {
    console.error('Failed to load images:', error)
  }
}

// 提交表单
const submitForm = async () => {
  try {
    if (!formRef.value.validate) throw new Error('表单未初始化')
    await formRef.value.validate()

    submitting.value = true

    const response = await fetch('/v1/user/instances', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(form.value)
    })

    const data = await response.json()

    if (response.ok) {
      ElMessage.success(t('user.apply.success'))
      router.push(`/user/instances`)
    } else {
      ElMessage.error(data.message || t('user.apply.failed'))
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
  loadImages()
})
</script>

<style scoped>
.apply-page {
  padding: 20px;
}

.config-card {
  margin-top: 20px;
}

.mb-2 {
  margin-bottom: 10px;
}

:deep(.el-input-number) {
  width: 100%;
}
</style>
