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


/* ── 现代 UI 增强 ── */

/* 卡片 */
:deep(.el-card),
:deep(.el-card.is-always-shadow),
:deep(.el-card.is-hover-shadow) {
  border-radius: 16px !important;
  border: 1px solid var(--border-color) !important;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05) !important;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1),
              box-shadow 0.3s cubic-bezier(0.4, 0, 0.2, 1),
              border-color 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.el-card:hover) {
  transform: translateY(-3px) !important;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.09) !important;
}

/* 统计卡片 */
.stat-card :deep(.el-card__body) {
  padding: 20px !important;
}

.stat-content {
  display: flex !important;
  align-items: center !important;
  gap: 14px !important;
}

.stat-icon {
  width: 48px !important;
  height: 48px !important;
  border-radius: 14px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  font-size: 22px !important;
  flex-shrink: 0 !important;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.08) !important;
}

.stat-icon.user-icon { background: linear-gradient(135deg, #6366f1, #818cf8); color: #fff; }
.stat-icon.server-icon { background: linear-gradient(135deg, #10b981, #34d399); color: #fff; }
.stat-icon.vm-icon { background: linear-gradient(135deg, #f59e0b, #fbbf24); color: #fff; }
.stat-icon.container-icon { background: linear-gradient(135deg, #ec4899, #f472b6); color: #fff; }
.stat-icon.cpu-icon { background: linear-gradient(135deg, #6366f1, #818cf8); color: #fff; }
.stat-icon.memory-icon { background: linear-gradient(135deg, #10b981, #34d399); color: #fff; }
.stat-icon.disk-icon { background: linear-gradient(135deg, #f59e0b, #fbbf24); color: #fff; }

.stat-number {
  font-size: 28px !important;
  font-weight: 800 !important;
  color: var(--text-color-primary) !important;
  line-height: 1.1 !important;
  font-variant-numeric: tabular-nums !important;
}

.stat-label {
  font-size: 13px !important;
  color: var(--text-color-secondary) !important;
  font-weight: 600 !important;
  margin-top: 2px !important;
}

/* 资源卡片 */
.resource-card :deep(.el-card__body) {
  padding: 16px 20px !important;
}

.resource-header {
  display: flex !important;
  align-items: center !important;
  gap: 10px !important;
  margin-bottom: 12px !important;
}

.resource-icon {
  width: 36px !important;
  height: 36px !important;
  border-radius: 10px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  font-size: 16px !important;
  flex-shrink: 0 !important;
}

.resource-title {
  font-weight: 700 !important;
  font-size: 14px !important;
  color: var(--text-color-primary) !important;
}

.resource-number-value {
  font-size: 22px !important;
  font-weight: 800 !important;
  color: var(--primary-color) !important;
  font-variant-numeric: tabular-nums !important;
}

.resource-number-label {
  font-size: 11px !important;
  color: var(--text-color-secondary) !important;
  font-weight: 600 !important;
}

.resource-number-separator {
  color: var(--text-color-tertiary) !important;
  font-weight: 700 !important;
  font-size: 16px !important;
  margin: 0 8px !important;
}

/* 按钮 */
:deep(.el-button--primary) {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-light)) !important;
  border: none !important;
  border-radius: 12px !important;
  font-weight: 600 !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1) !important;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.2) !important;
}

:deep(.el-button--primary:hover) {
  transform: translateY(-2px) !important;
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.3) !important;
}

:deep(.el-button--success) {
  border-radius: 12px !important;
  font-weight: 600 !important;
}

:deep(.el-button--warning) {
  border-radius: 12px !important;
  font-weight: 600 !important;
}

:deep(.el-button--danger) {
  border-radius: 12px !important;
  font-weight: 600 !important;
}

:deep(.el-button--default) {
  border-radius: 12px !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.el-button--default:hover) {
  transform: translateY(-1px) !important;
}

/* 表格 */
:deep(.el-table) {
  --el-table-border-color: var(--border-color) !important;
}

:deep(.el-table th) {
  background: rgba(99, 102, 241, 0.04) !important;
  color: var(--text-color-primary) !important;
  font-weight: 700 !important;
  font-size: 13px !important;
  letter-spacing: 0.3px !important;
}

:deep(.el-table td) {
  font-size: 13px !important;
  color: var(--text-color-primary) !important;
}

:deep(.el-table__row) {
  transition: all 0.2s ease !important;
}

:deep(.el-table__row:hover) {
  background: rgba(99, 102, 241, 0.04) !important;
  transform: scale(1.005) !important;
}

/* 徽章 */
:deep(.el-tag) {
  border-radius: 999px !important;
  font-weight: 600 !important;
  font-size: 12px !important;
  padding: 2px 12px !important;
  border: none !important;
}

:deep(.el-tag--success) {
  background: rgba(16, 185, 129, 0.12) !important;
  color: #10b981 !important;
}

:deep(.el-tag--warning) {
  background: rgba(245, 158, 11, 0.12) !important;
  color: #f59e0b !important;
}

:deep(.el-tag--danger) {
  background: rgba(239, 68, 68, 0.12) !important;
  color: #ef4444 !important;
}

:deep(.el-tag--info) {
  background: rgba(107, 114, 128, 0.1) !important;
  color: #6b7280 !important;
}

/* 输入框 */
:deep(.el-input__wrapper) {
  border-radius: 12px !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2) !important;
}

:deep(.el-select .el-input__wrapper),
:deep(.el-cascader .el-input__wrapper) {
  border-radius: 12px !important;
}

/* 分页 */
:deep(.el-pagination) {
  margin-top: 24px !important;
}

:deep(.el-pager li),
:deep(.el-pager li.btn-quicknext),
:deep(.el-pager li.btn-quickprev),
:deep(.el-pagination .btn-prev),
:deep(.el-pagination .btn-next) {
  border-radius: 10px !important;
  font-weight: 600 !important;
  transition: all 0.2s ease !important;
}

:deep(.el-pager li.is-active) {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-light)) !important;
  border-color: transparent !important;
}

/* 加载态 */
.el-loading-mask {
  backdrop-filter: blur(4px) !important;
  -webkit-backdrop-filter: blur(4px) !important;
  background: rgba(255, 255, 255, 0.7) !important;
}

.el-dark .el-loading-mask {
  background: rgba(15, 23, 42, 0.7) !important;
}

/* 标题 */
h1, h2, h3, h4 {
  color: var(--text-color-primary) !important;
}

/* 进场动画 */
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(18px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 响应式断点 */
@media (max-width: 768px) {
  .stat-content {
    flex-direction: column !important;
    align-items: flex-start !important;
  }

  .stat-number {
    font-size: 22px !important;
  }

  .stat-icon {
    width: 40px !important;
    height: 40px !important;
    font-size: 18px !important;
  }

  .resource-header {
    flex-direction: column !important;
    align-items: flex-start !important;
  }

  .resource-numbers {
    display: flex !important;
    align-items: center !important;
    gap: 4px !important;
    flex-wrap: wrap !important;
  }
}

@media (max-width: 480px) {
  .stat-number {
    font-size: 20px !important;
  }

  .stat-label {
    font-size: 12px !important;
  }

  .resource-number-value {
    font-size: 18px !important;
  }
}
</style>
