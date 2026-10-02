<template>
  <el-dialog
    :model-value="modelValue"
    :title="t('admin.instances.editAccess')"
    width="520px"
    destroy-on-close
    @update:model-value="emit('update:modelValue', $event)"
  >
    <el-form
      v-if="instance"
      label-position="top"
      @submit.prevent="save"
    >
      <el-form-item :label="t('admin.instances.sshHost')">
        <el-input
          v-model="form.sshHost"
          :placeholder="t('admin.instances.sshHostPlaceholder')"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item :label="t('admin.instances.sshPort')">
        <el-input-number
          v-model="form.sshPort"
          :min="1"
          :max="65535"
          controls-position="right"
          style="width: 100%;"
        />
      </el-form-item>
      <el-form-item :label="t('admin.instances.username')">
        <el-input
          v-model="form.username"
          autocomplete="username"
        />
      </el-form-item>
      <el-form-item :label="t('admin.instances.password')">
        <el-input
          v-model="form.password"
          type="password"
          show-password
          :placeholder="t('admin.instances.keepExistingSecret')"
          autocomplete="new-password"
        />
        <el-checkbox v-model="clearPassword">
          {{ t('admin.instances.clearPassword') }}
        </el-checkbox>
      </el-form-item>
      <el-form-item :label="t('admin.instances.sshKey')">
        <el-input
          v-model="form.sshKey"
          type="textarea"
          :rows="5"
          :placeholder="t('admin.instances.keepExistingSecret')"
          autocomplete="off"
        />
        <div class="secret-row">
          <span v-if="instance.hasSshKey">{{ t('admin.instances.sshKeyConfigured') }}</span>
          <el-checkbox v-model="clearSSHKey">
            {{ t('admin.instances.clearSSHKey') }}
          </el-checkbox>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="emit('update:modelValue', false)">
        {{ t('common.cancel') }}
      </el-button>
      <el-button
        type="primary"
        :loading="saving"
        @click="save"
      >
        {{ t('common.confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { updateInstance } from '@/api/admin'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  instance: { type: Object, default: null }
})
const emit = defineEmits(['update:modelValue', 'updated'])
const { t } = useI18n()
const saving = ref(false)
const clearPassword = ref(false)
const clearSSHKey = ref(false)
const form = reactive({ sshHost: '', sshPort: 22, username: '', password: '', sshKey: '' })

const resetForm = () => {
  const source = props.instance || {}
  form.sshHost = source.sshHost || ''
  form.sshPort = Number(source.sshPort) || 22
  form.username = source.username || ''
  form.password = ''
  form.sshKey = ''
  clearPassword.value = false
  clearSSHKey.value = false
}

watch(() => [props.modelValue, props.instance], ([visible]) => {
  if (visible) resetForm()
}, { immediate: true, deep: false })

const save = async () => {
  if (!props.instance?.id || saving.value) return
  const port = Number(form.sshPort)
  if (!Number.isInteger(port) || port < 1 || port > 65535) {
    ElMessage.error(t('admin.instances.sshPortInvalid'))
    return
  }
  if (form.password && clearPassword.value) {
    ElMessage.error(t('admin.instances.secretConflict'))
    return
  }
  if (form.sshKey && clearSSHKey.value) {
    ElMessage.error(t('admin.instances.secretConflict'))
    return
  }

  const payload = {
    sshHost: form.sshHost.trim(),
    sshPort: port,
    username: form.username.trim()
  }
  if (form.password) payload.password = form.password
  if (clearPassword.value) payload.password = ''
  if (form.sshKey) payload.sshKey = form.sshKey
  if (clearSSHKey.value) payload.sshKey = ''

  saving.value = true
  try {
    await updateInstance(props.instance.id, payload)
    ElMessage.success(t('admin.instances.editAccessSuccess'))
    emit('update:modelValue', false)
    emit('updated')
  } catch (error) {
    ElMessage.error(error?.fullMessage || error?.userMessage || error?.message || t('admin.instances.editAccessFailed'))
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.secret-row {
  display: flex;
  width: 100%;
  justify-content: space-between;
  align-items: center;
  color: var(--el-text-color-secondary);
  font-size: 12px;
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
