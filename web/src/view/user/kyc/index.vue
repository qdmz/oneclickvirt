<template>
  <div class="kyc-container">
    <el-card>
      <template #header>
        <span>{{ t('user.kyc.title') }}</span>
      </template>

      <!-- 已通过 -->
      <el-result
        v-if="kycStatus === 'approved'"
        icon="success"
        :title="t('user.kyc.statusApproved')"
        :sub-title="t('user.kyc.alreadyVerified')"
      />

      <!-- 审核中(手动) -->
      <el-result
        v-else-if="kycStatus === 'pending' && kycRecord?.method === 'manual'"
        icon="warning"
        :title="t('user.kyc.statusPending')"
        :sub-title="t('user.kyc.pendingReview')"
      />

      <!-- 审核中(支付宝) — 可查询结果 -->
      <div v-else-if="kycStatus === 'pending' && kycRecord?.method === 'alipay'">
        <el-result
          icon="warning"
          :title="t('user.kyc.statusPending')"
          :sub-title="t('user.kyc.alipayPendingTip')"
        />
        <div style="text-align: center; margin-top: 16px;">
          <el-button
            type="primary"
            :loading="queryLoading"
            @click="handleQueryAlipay"
          >
            {{ t('user.kyc.queryAlipayResult') }}
          </el-button>
        </div>
      </div>

      <!-- 已拒绝 -->
      <div v-else-if="kycStatus === 'rejected'">
        <el-alert
          type="error"
          :title="t('user.kyc.statusRejected')"
          :description="kycRecord?.rejectReason"
          show-icon
          :closable="false"
          style="margin-bottom: 20px;"
        />
        <kyc-form-component
          :kyc-method="kycMethodConfig"
          @submitted="fetchKYC"
        />
      </div>

      <!-- 未认证 -->
      <div v-else>
        <el-alert
          type="info"
          :title="t('user.kyc.statusNone')"
          show-icon
          :closable="false"
          style="margin-bottom: 20px;"
        />
        <kyc-form-component
          :kyc-method="kycMethodConfig"
          @submitted="fetchKYC"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getUserKYC, queryAlipayKYCResult } from '@/api/features'
import { useFeatureStore } from '@/pinia/modules/feature'

const { t } = useI18n()
const featureStore = useFeatureStore()

const kycStatus = ref('none')
const kycRecord = ref(null)
const queryLoading = ref(false)

const kycMethodConfig = computed(() => featureStore.kycMethod || 'manual')

async function fetchKYC() {
  try {
    const res = await getUserKYC()
    if (res.code === 200) {
      if (res.data?.status) {
        kycStatus.value = res.data.status
        kycRecord.value = res.data
      } else {
        kycStatus.value = 'none'
      }
    }
  } catch {
    kycStatus.value = 'none'
  }
}

async function handleQueryAlipay() {
  queryLoading.value = true
  try {
    const res = await queryAlipayKYCResult()
    if (res.code === 200) {
      if (res.data?.passed) {
        ElMessage.success(t('user.kyc.alipayVerifySuccess'))
        await fetchKYC()
      } else {
        ElMessage.warning(t('user.kyc.alipayNotPassed'))
      }
    }
  } catch (e) {
    ElMessage.error(e?.response?.data?.msg || t('user.kyc.queryFailed'))
  } finally {
    queryLoading.value = false
  }
}

onMounted(() => fetchKYC())
</script>

<script>
import { defineComponent, ref as componentRef, h as vueH, resolveComponent as resolveVueComponent } from 'vue'
import { ElMessage as ElementMessage } from 'element-plus'
import { useI18n as useVueI18n } from 'vue-i18n'
import { submitUserKYC as submitUserKYCRequest, submitAlipayKYC as submitAlipayKYCRequest } from '@/api/features'

// 内联KYC表单组件
const KycFormComponent = defineComponent({
  name: 'KycFormComponent',
  props: {
    kycMethod: { type: String, default: 'manual' }
  },
  emits: ['submitted'],
  setup(props, { emit }) {
    const { t } = useVueI18n()
    const submitting = componentRef(false)
    const selectedMethod = componentRef(props.kycMethod === 'both' ? 'manual' : props.kycMethod)
    const form = componentRef({ realName: '', idNumber: '' })

    async function handleSubmit() {
      if (!form.value.realName || !form.value.idNumber) {
        ElementMessage.warning(t('user.kyc.fillAllFields'))
        return
      }
      submitting.value = true
      try {
        if (selectedMethod.value === 'alipay') {
          const res = await submitAlipayKYCRequest(form.value)
          if (res.code === 200) {
            ElementMessage.success(t('user.kyc.alipayRedirectTip'))
            if (res.data?.certifyUrl) {
              window.open(res.data.certifyUrl, '_blank')
            }
            emit('submitted')
          }
        } else {
          const res = await submitUserKYCRequest(form.value)
          if (res.code === 200) {
            ElementMessage.success(t('user.kyc.submitSuccess'))
            emit('submitted')
          }
        }
      } finally {
        submitting.value = false
      }
    }

    return () => {
      const ElForm = resolveVueComponent('el-form')
      const ElFormItem = resolveVueComponent('el-form-item')
      const ElInput = resolveVueComponent('el-input')
      const ElButton = resolveVueComponent('el-button')
      const ElRadioGroup = resolveVueComponent('el-radio-group')
      const ElRadio = resolveVueComponent('el-radio')
      const ElAlert = resolveVueComponent('el-alert')

      const children = []

      // Method selection if both methods enabled
      if (props.kycMethod === 'both') {
        children.push(
          vueH(ElFormItem, { label: t('user.kyc.verifyMethod') }, () => [
            vueH(ElRadioGroup, {
              modelValue: selectedMethod.value,
              'onUpdate:modelValue': v => { selectedMethod.value = v }
            }, () => [
              vueH(ElRadio, { value: 'manual' }, () => t('user.kyc.methodManual')),
              vueH(ElRadio, { value: 'alipay' }, () => t('user.kyc.methodAlipay'))
            ])
          ])
        )
      }

      if (selectedMethod.value === 'alipay') {
        children.push(
          vueH(ElAlert, {
            type: 'info',
            title: t('user.kyc.alipayTip'),
            showIcon: true,
            closable: false,
            style: 'margin-bottom: 16px;'
          })
        )
      }

      children.push(
        vueH(ElFormItem, { label: t('user.kyc.realName'), prop: 'realName' }, () => [
          vueH(ElInput, {
            modelValue: form.value.realName,
            'onUpdate:modelValue': v => { form.value.realName = v },
            placeholder: t('user.kyc.realNamePlaceholder')
          })
        ]),
        vueH(ElFormItem, { label: t('user.kyc.idNumber'), prop: 'idNumber' }, () => [
          vueH(ElInput, {
            modelValue: form.value.idNumber,
            'onUpdate:modelValue': v => { form.value.idNumber = v },
            placeholder: t('user.kyc.idNumberPlaceholder')
          })
        ]),
        vueH(ElFormItem, {}, () => [
          vueH(ElButton, {
            type: 'primary',
            loading: submitting.value,
            onClick: handleSubmit
          }, () => selectedMethod.value === 'alipay'
            ? t('user.kyc.startAlipayVerify')
            : t('user.kyc.submit'))
        ])
      )

      return vueH(ElForm, { model: form.value, labelWidth: '120px' }, () => children)
    }
  }
})

export default { components: { KycFormComponent } }
</script>

<style scoped>
.kyc-container { padding: 20px; }


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
