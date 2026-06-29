<template>
  <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
    <el-form-item label="实例名称" prop="name">
      <el-input v-model="form.name" placeholder="请输入实例名称" />
    </el-form-item>

    <el-form-item label="实例类型" prop="instance_type">
      <el-radio-group v-model="form.instance_type">
        <el-radio label="vm">虚拟机 (VM)</el-radio>
        <el-radio label="container">容器</el-radio>
      </el-radio-group>
    </el-form-item>

    <el-form-item label="所属用户" prop="userId">
      <el-select v-model="form.userId" placeholder="请选择用户" clearable filterable style="width: 100%">
        <el-option
          v-for="user in users"
          :key="user.id"
          :label="`${user.username} (${user.email || '无邮箱'})`"
          :value="user.id"
        />
      </el-select>
    </el-form-item>

    <el-form-item label="节点" prop="provider">
      <el-select v-model="form.provider" placeholder="请选择节点" clearable filterable style="width: 100%">
        <el-option
          v-for="p in providers"
          :key="p.id"
          :label="`${p.name} (${p.type})`"
          :value="p.name"
        />
      </el-select>
    </el-form-item>

    <el-form-item label="镜像" prop="image">
      <el-select v-model="form.image" placeholder="请选择镜像" clearable filterable style="width: 100%">
        <el-option
          v-for="img in images"
          :key="img.id"
          :label="`${img.name} (${img.providerType}/${img.instanceType}/${img.architecture})`"
          :value="img.name"
        />
      </el-select>
    </el-form-item>

    <el-form-item label="CPU核心数" prop="cpu">
      <el-input-number v-model="form.cpu" :min="1" :max="64" :step="1" style="width: 100%" />
    </el-form-item>

    <el-form-item label="内存(GB)" prop="memory">
      <el-input-number v-model="form.memory" :min="1" :max="128" :step="1" style="width: 100%" />
    </el-form-item>

    <el-form-item label="磁盘(GB)" prop="disk">
      <el-input-number v-model="form.disk" :min="10" :max="2048" :step="10" style="width: 100%" />
    </el-form-item>

    <el-form-item>
      <el-button type="primary" @click="submitForm" :loading="submitting">创建实例</el-button>
      <el-button @click="resetForm">重置</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const emit = defineEmits(['created', 'cancel'])
const formRef = ref()
const submitting = ref(false)

const form = ref({
  name: '',
  instance_type: 'vm',
  userId: null,
  provider: '',
  image: '',
  cpu: 2,
  memory: 4,
  disk: 40
})

const rules = {
  name: [{ required: true, message: '请输入实例名称', trigger: 'blur' }],
  instance_type: [{ required: true, message: '请选择实例类型', trigger: 'change' }],
  userId: [{ required: true, message: '请选择用户', trigger: 'change' }],
  provider: [{ required: true, message: '请选择节点', trigger: 'change' }],
  image: [{ required: true, message: '请选择镜像', trigger: 'change' }],
  cpu: [{ required: true, message: '请输入CPU核心数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
  disk: [{ required: true, message: '请输入磁盘大小', trigger: 'blur' }]
}

const users = ref([])
const providers = ref([])
const images = ref([])

const loadUsers = async () => {
  try {
    const res = await request.get('/v1/admin/users')
    if (res.code === 200) {
      users.value = res.data || []
    }
  } catch (error) {
    console.error('加载用户列表失败:', error)
  }
}

const loadProviders = async () => {
  try {
    const res = await request.get('/v1/admin/providers')
    if (res.code === 200) {
      providers.value = res.data || []
    }
  } catch (error) {
    console.error('加载节点列表失败:', error)
  }
}

const loadImages = async () => {
  try {
    const res = await request.get('/v1/admin/system-images')
    if (res.code === 200) {
      images.value = res.data || []
    }
  } catch (error) {
    console.error('加载镜像列表失败:', error)
  }
}

const submitForm = async () => {
  try {
    if (!formRef.value) throw new Error('表单引用无效')
    await formRef.value.validate()

    submitting.value = true

    const payload = {
      name: form.value.name,
      instance_type: form.value.instance_type,
      userId: form.value.userId,
      provider: form.value.provider,
      image: form.value.image,
      cpu: form.value.cpu,
      memory: form.value.memory,
      disk: form.value.disk
    }

    const res = await request.post('/v1/admin/instances', payload)

    if (res.code === 200) {
      ElMessage.success('创建实例成功')
      emit('created')
    } else {
      ElMessage.error(res.msg || '创建实例失败')
    }
  } catch (error) {
    if (error.message && error.message !== '表单引用无效') {
      ElMessage.error('创建实例失败: ' + error.message)
    }
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  form.value = {
    name: '',
    instance_type: 'vm',
    userId: null,
    provider: '',
    image: '',
    cpu: 2,
    memory: 4,
    disk: 40
  }
  formRef.value?.resetFields()
}

onMounted(() => {
  loadUsers()
  loadProviders()
  loadImages()
})
</script>
