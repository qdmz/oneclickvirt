<template>
  <div class="kyc-container">
    <el-card class="status-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>{{ $t('user.kyc.title') }}</span>
        </div>
      </template>

      <!-- Not verified -->
      <div v-if="!kycData || kycData.status === -1">
        <el-alert type="info" :closable="false" show-icon :title="$t('user.kyc.notVerified')" style="margin-bottom: 20px" />
        <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" label-position="top">
          <el-form-item :label="$t('user.kyc.realName')" prop="realName">
            <el-input v-model="form.realName" :placeholder="$t('user.kyc.realNamePlaceholder')" maxlength="32" />
          </el-form-item>
          <el-form-item :label="$t('user.kyc.idCardNumber')" prop="idCardNumber">
            <el-input v-model="form.idCardNumber" :placeholder="$t('user.kyc.idCardPlaceholder')" maxlength="18" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="submitting" @click="handleSubmit">
              {{ $t('user.kyc.submit') }}
            </el-button>
          </el-form-item>
        </el-form>
        <el-alert type="warning" :closable="false" show-icon>
          <template #title>{{ $t('user.kyc.notice') }}</template>
        </el-alert>
      </div>

      <!-- Pending -->
      <div v-else-if="kycData.status === 0">
        <el-result icon="info" :title="$t('user.kyc.pending')">
          <template #extra>
            <el-button type="primary" :loading="querying" @click="handleQuery">
              {{ $t('user.kyc.queryStatus') }}
            </el-button>
          </template>
        </el-result>
      </div>

      <!-- Certified -->
      <div v-else-if="kycData.status === 1">
        <el-result icon="success" :title="$t('user.kyc.certified')">
          <template #extra>
            <el-descriptions :column="1" border size="small" style="max-width: 400px; margin: 0 auto">
              <el-descriptions-item :label="$t('user.kyc.realName')">{{ kycData.realName }}</el-descriptions-item>
              <el-descriptions-item :label="$t('user.kyc.certifiedAt')">{{ formatDate(kycData.certifiedAt) }}</el-descriptions-item>
              <el-descriptions-item :label="$t('user.kyc.idType')">{{ kycData.idType }}</el-descriptions-item>
            </el-descriptions>
          </template>
        </el-result>
      </div>

      <!-- Rejected -->
      <div v-else-if="kycData.status === 2">
        <el-result icon="error" :title="$t('user.kyc.rejected')" :sub-title="kycData.remark">
          <template #extra>
            <el-button type="primary" @click="resetForm">
              {{ $t('user.kyc.retry') }}
            </el-button>
          </template>
        </el-result>
      </div>

      <!-- Expired -->
      <div v-else-if="kycData.status === 3">
        <el-result icon="warning" :title="$t('user.kyc.expired')">
          <template #extra>
            <el-button type="primary" @click="resetForm">
              {{ $t('user.kyc.retry') }}
            </el-button>
          </template>
        </el-result>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { submitKYC, getKYCStatus, queryKYC } from '@/api/kyc'

const { t } = useI18n()

const formRef = ref(null)
const kycData = ref(null)
const submitting = ref(false)
const querying = ref(false)

const form = ref({
  realName: '',
  idCardNumber: ''
})

const rules = {
  realName: [{ required: true, message: () => t('user.kyc.realNameRequired'), trigger: 'blur' }],
  idCardNumber: [
    { required: true, message: () => t('user.kyc.idCardRequired'), trigger: 'blur' },
    { pattern: /^\d{17}[\dXx]$/, message: () => t('user.kyc.idCardInvalid'), trigger: 'blur' }
  ]
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

const fetchStatus = async () => {
  try {
    const res = await getKYCStatus()
    if (res.data && res.data.code === 0) {
      kycData.value = res.data.data
    } else if (res.data && res.data.data && res.data.data.status === -1) {
      kycData.value = null
    }
  } catch (e) {
    // ignore
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  submitting.value = true
  try {
    const res = await submitKYC(form.value)
    if (res.data.code === 0) {
      ElMessage.success(t('user.kyc.submitSuccess'))
      kycData.value = res.data.data.record
      // Open alipay certify page
      const certifyURL = res.data.data.certifyURL
      if (certifyURL) {
        window.open(certifyURL, '_blank')
      }
    } else {
      ElMessage.error(res.data.message || t('user.kyc.submitFailed'))
    }
  } catch (e) {
    ElMessage.error(e.message || t('user.kyc.submitFailed'))
  } finally {
    submitting.value = false
  }
}

const handleQuery = async () => {
  querying.value = true
  try {
    const res = await queryKYC()
    if (res.data.code === 0) {
      kycData.value = res.data.data
      ElMessage.success(t('user.kyc.querySuccess'))
    } else {
      ElMessage.error(res.data.message || t('user.kyc.queryFailed'))
    }
  } catch (e) {
    ElMessage.error(e.message || t('user.kyc.queryFailed'))
  } finally {
    querying.value = false
  }
}

const resetForm = () => {
  kycData.value = null
  form.value = { realName: '', idCardNumber: '' }
}

onMounted(() => {
  fetchStatus()
  // Check callback param
  const url = new URL(window.location.href)
  if (url.searchParams.get('callback') === 'success') {
    setTimeout(handleQuery, 1500)
  }
})
</script>

<style scoped>
.kyc-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}
.card-header {
  font-size: 16px;
  font-weight: 600;
}
</style>
