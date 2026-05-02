<template>
  <el-form ref="formRef" :model="form" :rules="rules" label-width="150px">
    <el-form-item label="实例名称" prop="name">
      <el-input v-model="form.name" placeholder="请输入实例名称" />
    </el-form-item>

    <el-form-item label="实例类型" prop="instance_type">
      <el-radio-group v-model="form.instance_type">
        <el-radio label="vm">虚拟机 (VM)</el-radio>
        <el-radio label="container">容器</el-radio>
      </el-radio-group>
    </el-form-item>

    <el-form-item label="所属用户" prop="user_id">
      <el-select v-model="form.user_id" placeholder="请选择用户" clearable filterable>
        <el-option
          v-for="user in users"
          :key="user.id"
          :label="`${user.username}`"
          :value="user.id"
        />
      </el-select>
    </el-form-item>

    <el-form-item label="设置密码" prop="password">
      <el-input
        v-model="form.password"
        type="password"
        placeholder="请输入密码"
        show-password
      >
        <template #append>
          <el-button @click="generatePassword">
            随机生成
          </el-button>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item>
      <el-button type="primary" @click="submitForm">创建实例</el-button>
      <el-button @click="resetForm">重置</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'

const emit = defineEmits(['created', 'cancel'])
const formRef = ref()

const form = ref({
  name: '',
  instance_type: 'vm',
  user_id: null,
  password: ''
})

const rules = {
  name: [{ required: true, message: '请输入实例名称', trigger: 'blur' }],
  instance_type: [{ required: true, message: '请选择实例类型', trigger: 'change' }],
  user_id: [{ required: true, message: '请选择用户', trigger: 'change' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const users = ref([])

const loadUsers = async () => {
  try {
    const response = await fetch('/v1/users')
    if (response.ok) {
      users.value = await response.json()
    }
  } catch (error) {
    console.error('加载用户列表失败:', error)
  }
}

const generatePassword = () => {
  const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*'
  let password = ''
  for (let i = 0; i < 12; i++) {
    password += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  form.value.password = password
}

const submitForm = async () => {
  try {
    if (!formRef.value) throw new Error('表单引用无效')
    await formRef.value.validate()

    const response = await fetch('/v1/admin/instances', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: form.value.name,
        instance_type: form.value.instance_type,
        user_id: form.value.user_id,
        password: form.value.password
      })
    })

    if (response.ok) {
      ElMessage.success('实例创建成功')
      emit('created')
    } else {
      ElMessage.error('创建失败')
    }
  } catch (error) {
    if (!error.message.includes('validate')) {
      ElMessage.error('创建失败')
    }
  }
}

const resetForm = () => {
  form.value = {
    name: '',
    instance_type: 'vm',
    user_id: null,
    password: ''
  }
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

onMounted(() => {
  loadUsers()
})
</script>
