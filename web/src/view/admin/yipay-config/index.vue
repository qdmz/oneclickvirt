<template>
  <div class="yipay-config-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('admin.yipayConfig.title') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="form" label-width="160px" class="config-form" v-loading="loading">
        <el-divider content-position="left">{{ t('admin.yipayConfig.basicConfig') }}</el-divider>

        <el-form-item :label="t('admin.yipayConfig.enabled')">
          <el-switch v-model="form.enabled" />
        </el-form-item>

        <el-form-item :label="t('admin.yipayConfig.name')" prop="name">
          <el-input v-model="form.name" style="width: 300px;" />
        </el-form-item>

        <el-form-item :label="t('admin.yipayConfig.apiUrl')" prop="apiUrl">
          <el-input v-model="form.apiUrl" :placeholder="t('admin.yipayConfig.apiUrlPlaceholder')" style="width: 400px;" />
        </el-form-item>

        <el-form-item :label="t('admin.yipayConfig.pid')" prop="pid">
          <el-input v-model="form.pid" :placeholder="t('admin.yipayConfig.pidPlaceholder')" style="width: 300px;" />
        </el-form-item>

        <el-form-item :label="t('admin.yipayConfig.key')" prop="key">
          <el-input v-model="form.key" type="password" show-password :placeholder="t('admin.yipayConfig.keyPlaceholder')" style="width: 400px;" />
        </el-form-item>

        <el-form-item label="启用的支付方式" prop="enabledPayTypes">
          <el-checkbox-group v-model="enabledPayTypesArray">
            <el-checkbox label="支付宝" value="alipay" />
            <el-checkbox label="微信支付" value="wxpay" />
            <el-checkbox label="QQ钱包" value="qqpay" />
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="默认支付方式" prop="payType">
          <el-select v-model="form.payType" style="width: 300px;">
            <el-option v-if="enabledPayTypesArray.includes('alipay')" label="支付宝" value="alipay" />
            <el-option v-if="enabledPayTypesArray.includes('wxpay')" label="微信支付" value="wxpay" />
            <el-option v-if="enabledPayTypesArray.includes('qqpay')" label="QQ支付" value="qqpay" />
          </el-select>
        </el-form-item>

        <el-divider content-position="left">{{ t('admin.yipayConfig.callbackConfig') }}</el-divider>

        <el-form-item :label="t('admin.yipayConfig.notifyUrl')">
          <el-input v-model="form.notifyUrl" readonly style="width: 500px;" />
          <div class="form-hint">{{ t('admin.yipayConfig.notifyUrlHint') }}</div>
        </el-form-item>

        <el-form-item :label="t('admin.yipayConfig.returnUrl')">
          <el-input v-model="form.returnUrl" readonly style="width: 500px;" />
          <div class="form-hint">{{ t('admin.yipayConfig.returnUrlHint') }}</div>
        </el-form-item>

        <el-divider content-position="left">{{ t('admin.yipayConfig.feeConfig') }}</el-divider>

        <el-form-item :label="t('admin.yipayConfig.feePercent')">
          <el-input-number v-model="form.feePercent" :min="0" :max="100" :precision="2" style="width: 200px;" />
          <span style="margin-left: 8px;">%</span>
        </el-form-item>

        <el-form-item :label="t('admin.yipayConfig.minAmount')">
          <el-input-number v-model="form.minAmount" :min="0" :precision="2" style="width: 200px;" />
        </el-form-item>

        <el-form-item :label="t('admin.yipayConfig.maxAmount')">
          <el-input-number v-model="form.maxAmount" :min="0" :precision="2" style="width: 200px;" />
        </el-form-item>

        <el-alert type="info" :closable="false" style="margin: 20px 0;">
          {{ t('admin.yipayConfig.hint') }}
        </el-alert>

        <div class="form-actions">
          <el-button type="primary" size="large" :loading="submitting" @click="handleSubmit">
            {{ t('common.save') }}
          </el-button>
          <el-button size="large" :loading="testing" @click="handleTest">
            {{ testing ? t('admin.yipayConfig.testing') : t('admin.yipayConfig.test') }}
          </el-button>
        </div>

        <el-alert
          v-if="testResultVisible"
          :type="testResultType"
          :closable="true"
          @close="testResultVisible = false"
          style="margin-top: 16px;"
        >
          <template #title>{{ t('admin.yipayConfig.testResultTitle') }}</template>
          <div style="margin-top: 6px; white-space: pre-wrap; word-break: break-all;">
            {{ testResultMessage }}
          </div>
        </el-alert>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getYiPayConfig, updateYiPayConfig, testYiPayConfig } from '@/api/admin'

const { t } = useI18n()
const formRef = ref(null)
const loading = ref(false)
const submitting = ref(false)
const testing = ref(false)
const testResultVisible = ref(false)
const testResultType = ref('info')
const testResultMessage = ref('')

const form = ref({
  enabled: false,
  name: '易支付',
  apiUrl: '',
  pid: '',
  key: '',
  payType: 'alipay',
  enabledPayTypes: 'alipay,wxpay,qqpay',
  notifyUrl: '',
  returnUrl: '',
  feePercent: 0,
  minAmount: 1,
  maxAmount: 10000
})

const enabledPayTypesArray = computed({
  get() {
    if (!form.value.enabledPayTypes) return []
    return form.value.enabledPayTypes.split(',').map(t => t.trim()).filter(Boolean)
  },
  set(val) {
    form.value.enabledPayTypes = (val && val.length > 0) ? val.join(',') : ''
  }
})

const loadConfig = async () => {
  loading.value = true
  try {
    const res = await getYiPayConfig()
    if (res && res.code === 200 && res.data) {
      Object.assign(form.value, res.data)
    }
    // 处理 enabledPayTypes 字段：如果为空或不存在，默认启用全部支付方式
    if (!form.value.enabledPayTypes) {
      form.value.enabledPayTypes = 'alipay,wxpay,qqpay'
    }
    // 确保默认支付方式在已启用的支付方式中，否则回退到第一个启用的支付方式
    const enabledArr = form.value.enabledPayTypes.split(',').map(t => t.trim()).filter(Boolean)
    if (enabledArr.length > 0 && !enabledArr.includes(form.value.payType)) {
      form.value.payType = enabledArr[0]
    }
    // Auto-fill callback URLs based on current domain
    const origin = window.location.origin
    if (!form.value.notifyUrl) {
      form.value.notifyUrl = `${origin}/api/v1/public/payments/yipay/notify`
    }
    if (!form.value.returnUrl) {
      form.value.returnUrl = `${origin}/api/v1/public/payments/yipay/return`
    }
  } catch (error) {
    ElMessage.error(error?.message || t('admin.yipayConfig.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    const res = await updateYiPayConfig(form.value)
    if (res && res.code === 200) {
      ElMessage.success(t('admin.yipayConfig.saveSuccess'))
    } else {
      ElMessage.error(res?.message || t('admin.yipayConfig.saveFailed'))
    }
  } catch (error) {
    ElMessage.error(error?.message || t('admin.yipayConfig.saveFailed'))
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadConfig()
})

// 易支付连通性/密钥自检：用当前配置的 key 向网关发起订单查询
const handleTest = async () => {
  testing.value = true
  testResultVisible.value = false
  try {
    const res = await testYiPayConfig()
    if (res && res.code === 200 && res.data) {
      const query = res.data.query || {}
      const raw = query.raw || ''
      // 网关返回含 -3 / 商户密钥错误 => 密钥不匹配
      const isKeyMismatch = raw.includes('-3') || raw.includes('商户密钥错误') || raw.includes('merchant key')
      if (raw.includes('"code":1') || raw.includes('"code": 1') || (!isKeyMismatch && raw.includes('"code"'))) {
        testResultType.value = 'success'
        testResultMessage.value = t('admin.yipayConfig.testSuccessHint') + '\n' +
          t('admin.yipayConfig.testRaw') + ':\n' + raw
      } else if (isKeyMismatch) {
        testResultType.value = 'error'
        testResultMessage.value = t('admin.yipayConfig.testKeyMismatchHint') + '\n' +
          t('admin.yipayConfig.testRaw') + ':\n' + raw
      } else {
        testResultType.value = 'warning'
        testResultMessage.value = t('admin.yipayConfig.testRaw') + ':\n' + raw
      }
      testResultVisible.value = true
    } else {
      ElMessage.error(res?.message || t('admin.yipayConfig.testNetworkErrorHint'))
    }
  } catch (error) {
    ElMessage.error(error?.message || t('admin.yipayConfig.testNetworkErrorHint'))
  } finally {
    testing.value = false
  }
}
</script>

<style scoped lang="scss">
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  > span {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-color-primary);
  }
}

.config-form {
  max-width: 800px;
  margin: 0 auto;

  :deep(.el-divider__text) {
    font-size: 15px;
    font-weight: 600;
    color: var(--text-color-primary);
  }
}

.form-hint {
  font-size: 12px;
  color: var(--text-color-tertiary);
  margin-top: 4px;
  line-height: 1.4;
}

.form-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 24px 0 8px;
  border-top: 1px solid var(--border-color);
  margin-top: 24px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .config-form {
    max-width: 100%;
  }

  .form-actions {
    flex-direction: column;
    align-items: center;

    .el-button {
      width: 100%;
      max-width: 200px;
    }
  }
}
</style>
