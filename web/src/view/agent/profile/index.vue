<template>
  <div class="agent-profile">
    <div class="page-header">
      <h1>代理商资料</h1>
    </div>

    <!-- 申请代理商 -->
    <el-card v-if="!profile">
      <template #header><span>申请成为代理商</span></template>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" style="max-width: 500px;">
        <el-form-item label="代理商名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入代理商名称" />
        </el-form-item>
        <el-form-item label="联系邮箱" prop="contactEmail">
          <el-input v-model="form.contactEmail" placeholder="请输入联系邮箱" />
        </el-form-item>
        <el-form-item label="联系人" prop="contactName">
          <el-input v-model="form.contactName" placeholder="请输入联系人姓名" />
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input v-model="form.contactPhone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleApply">提交申请</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 代理商信息 -->
    <template v-else>
      <!-- 审核状态提示 -->
      <el-alert
        v-if="profile.status === 0"
        title="您的代理商申请正在审核中"
        type="warning"
        show-icon
        :closable="false"
        style="margin-bottom: 20px;"
      />

      <el-card>
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
            <el-button v-if="profile.status === 1" type="primary" size="small" @click="editMode = !editMode">
              {{ editMode ? '取消' : '编辑' }}
            </el-button>
          </div>
        </template>

        <el-form v-if="editMode" :model="editForm" label-width="100px" style="max-width: 500px;">
          <el-form-item label="代理商名称">
            <el-input v-model="editForm.name" />
          </el-form-item>
          <el-form-item label="联系邮箱">
            <el-input v-model="editForm.contactEmail" />
          </el-form-item>
          <el-form-item label="联系人">
            <el-input v-model="editForm.contactName" />
          </el-form-item>
          <el-form-item label="联系电话">
            <el-input v-model="editForm.contactPhone" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="submitting" @click="handleUpdate">保存</el-button>
          </el-form-item>
        </el-form>

        <el-descriptions v-else :column="2" border>
          <el-descriptions-item label="代理商名称">{{ profile.name }}</el-descriptions-item>
          <el-descriptions-item label="推广码">
            <el-tag>{{ profile.code }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="{ 0: 'warning', 1: 'success', 2: 'danger' }[profile.status]">
              {{ { 0: '待审核', 1: '正常', 2: '禁用' }[profile.status] }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="佣金比例">{{ profile.commissionRate }}%</el-descriptions-item>
          <el-descriptions-item label="联系人">{{ profile.contactName || '-' }}</el-descriptions-item>
          <el-descriptions-item label="联系邮箱">{{ profile.contactEmail || '-' }}</el-descriptions-item>
          <el-descriptions-item label="联系电话">{{ profile.contactPhone || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 推广链接 -->
      <el-card v-if="profile.status === 1" style="margin-top: 20px;">
        <template #header><span>推广链接</span></template>
        <div class="promo-link">
          <el-input :model-value="promoLink" readonly>
            <template #append>
              <el-button @click="copyLink">复制</el-button>
            </template>
          </el-input>
          <p class="promo-tip">分享此链接给新用户，通过此链接注册的用户将自动成为您的子用户。</p>
        </div>
      </el-card>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getAgentProfile, createAgent, updateAgentProfile } from '@/api/agent'
import { ElMessage } from 'element-plus'

const profile = ref(null)
const editMode = ref(false)
const submitting = ref(false)
const formRef = ref(null)

const form = ref({ name: '', contactEmail: '', contactName: '', contactPhone: '' })
const editForm = ref({ name: '', contactEmail: '', contactName: '', contactPhone: '' })

const rules = {
  name: [{ required: true, message: '请输入代理商名称', trigger: 'blur' }],
  contactEmail: [
    { required: true, message: '请输入联系邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ]
}

const promoLink = computed(() => {
  if (!profile.value) return ''
  const base = window.location.origin + window.location.pathname
  return `${base}#/register?agent=${profile.value.code}`
})

const copyLink = () => {
  navigator.clipboard.writeText(promoLink.value)
  ElMessage.success('推广链接已复制')
}

const loadProfile = async () => {
  try {
    const res = await getAgentProfile()
    if (res.data?.code === 0) {
      profile.value = res.data.data
      editForm.value = {
        name: profile.value.name,
        contactEmail: profile.value.contactEmail,
        contactName: profile.value.contactName,
        contactPhone: profile.value.contactPhone
      }
    }
  } catch { /* 404 = no profile yet */ }
}

const handleApply = async () => {
  try {
    await formRef.value.validate()
  } catch { return }

  submitting.value = true
  try {
    await createAgent(form.value)
    ElMessage.success('申请已提交')
    await loadProfile()
  } catch (error) {
    ElMessage.error(error?.response?.data?.message || '申请失败')
  } finally {
    submitting.value = false
  }
}

const handleUpdate = async () => {
  submitting.value = true
  try {
    await updateAgentProfile(editForm.value)
    ElMessage.success('更新成功')
    editMode.value = false
    await loadProfile()
  } catch { ElMessage.error('更新失败') }
  finally { submitting.value = false }
}

onMounted(loadProfile)
</script>

<style scoped>
.agent-profile { padding: 20px; }
.page-header { margin-bottom: 20px; }
.page-header h1 { font-size: 24px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.promo-link { max-width: 600px; }
.promo-tip { margin-top: 12px; color: #909399; font-size: 13px; }
</style>
