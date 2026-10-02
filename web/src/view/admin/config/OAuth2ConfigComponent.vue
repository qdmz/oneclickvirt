<template>
  <div class="oauth2-config-container">
    <el-form
      ref="formRef"
      v-loading="loading"
      :model="formData"
      :rules="rules"
      label-width="150px"
      class="config-form"
    >
      <el-card
        class="oauth-card"
        shadow="never"
      >
        <template #header>
          <div class="card-header">
            <span>{{ t('admin.config.oauth2LoginConfig') }}</span>
            <el-switch
              v-model="formData.enabled"
              :active-text="t('common.enabled')"
              :inactive-text="t('common.disabled')"
            />
          </div>
        </template>
        <el-divider content-position="left">
          {{ t('admin.config.oauth2BasicConfig') }}
        </el-divider>

        <el-form-item
          :label="t('admin.config.oauth2ClientIdLabel')"
          prop="clientId"
        >
          <el-input
            v-model="formData.clientId"
            :placeholder="t('admin.config.oauth2ClientIdPlaceholder')"
          />
        </el-form-item>

        <el-form-item
          :label="t('admin.config.oauth2ClientSecretLabel')"
          prop="clientSecret"
        >
          <el-input
            v-model="formData.clientSecret"
            type="password"
            :placeholder="t('admin.config.oauth2ClientSecretPlaceholder')"
            show-password
          />
        </el-form-item>

        <el-form-item
          :label="t('admin.config.oauth2RedirectUrlLabel')"
          prop="redirectUrl"
        >
          <el-input
            v-model="formData.redirectUrl"
            :placeholder="t('admin.config.oauth2RedirectUrlPlaceholder')"
          />
        </el-form-item>

        <el-divider content-position="left">
          {{ t('admin.config.oauth2EndpointConfig') }}
        </el-divider>

        <el-form-item
          :label="t('admin.config.oauth2AuthUrlLabel')"
          prop="authUrl"
        >
          <el-input
            v-model="formData.authUrl"
            :placeholder="t('admin.config.oauth2AuthUrlPlaceholder')"
          />
        </el-form-item>

        <el-form-item
          :label="t('admin.config.oauth2TokenUrlLabel')"
          prop="tokenUrl"
        >
          <el-input
            v-model="formData.tokenUrl"
            :placeholder="t('admin.config.oauth2TokenUrlPlaceholder')"
          />
        </el-form-item>

        <el-form-item
          :label="t('admin.config.oauth2UserInfoUrlLabel')"
          prop="userinfoUrl"
        >
          <el-input
            v-model="formData.userinfoUrl"
            :placeholder="t('admin.config.oauth2UserInfoUrlPlaceholder')"
          />
        </el-form-item>

        <el-form-item
          :label="t('admin.config.oauth2ScopesLabel')"
          prop="scopes"
        >
          <el-select
            v-model="formData.scopes"
            multiple
            filterable
            allow-create
            :placeholder="t('admin.config.oauth2ScopesPlaceholder')"
            style="width: 100%"
          >
            <el-option
              label="read"
              value="read"
            />
            <el-option
              label="openid"
              value="openid"
            />
            <el-option
              label="profile"
              value="profile"
            />
            <el-option
              label="email"
              value="email"
            />
          </el-select>
        </el-form-item>

        <el-divider content-position="left">
          {{ t('admin.config.oauth2FieldMapping') }}
        </el-divider>

        <el-form-item
          :label="t('admin.config.oauth2UserIdFieldLabel')"
          prop="userIdField"
        >
          <el-input
            v-model="formData.userIdField"
            :placeholder="t('admin.config.oauth2UserIdFieldPlaceholder')"
          >
            <template #append>
              {{ t('admin.config.oauth2UserIdFieldAppend') }}
            </template>
          </el-input>
        </el-form-item>

        <el-form-item
          :label="t('admin.config.oauth2UsernameFieldLabel')"
          prop="usernameField"
        >
          <el-input
            v-model="formData.usernameField"
            :placeholder="t('admin.config.oauth2UsernameFieldPlaceholder')"
          />
        </el-form-item>

        <el-form-item
          :label="t('admin.config.oauth2EmailFieldLabel')"
          prop="emailField"
        >
          <el-input
            v-model="formData.emailField"
            :placeholder="t('admin.config.oauth2EmailFieldPlaceholder')"
          />
        </el-form-item>

        <el-form-item
          :label="t('admin.config.oauth2AvatarFieldLabel')"
          prop="avatarField"
        >
          <el-input
            v-model="formData.avatarField"
            :placeholder="t('admin.config.oauth2AvatarFieldPlaceholder')"
          />
        </el-form-item>

        <el-form-item
          :label="t('admin.config.oauth2TrustLevelFieldLabel')"
          prop="trustLevelField"
        >
          <el-input
            v-model="formData.trustLevelField"
            :placeholder="t('admin.config.oauth2TrustLevelFieldPlaceholder')"
          />
        </el-form-item>

        <el-divider content-position="left">
          {{ t('admin.config.oauth2RegistrationLimit') }}
        </el-divider>

        <el-form-item
          :label="t('admin.config.oauth2MaxRegistrations')"
          prop="maxRegistrations"
        >
          <el-input-number
            v-model="formData.maxRegistrations"
            :min="0"
            :controls="false"
            :placeholder="t('admin.config.oauth2MaxRegistrationsPlaceholder')"
            style="width: 100%"
          />
          <div class="form-item-tip">
            {{ t('admin.config.oauth2MaxRegistrationsHint') }}
          </div>
        </el-form-item>

        <el-form-item
          :label="t('admin.config.oauth2CurrentRegistrations')"
        >
          <el-input-number
            v-model="formData.currentRegistrations"
            :min="0"
            :controls="false"
            disabled
            style="width: 100%"
          />
          <div class="form-item-tip">
            {{ t('admin.config.oauth2CurrentRegistrationsHint') }}
            <el-button
              type="danger"
              size="small"
              plain
              @click="resetRegistrationCount"
            >
              {{ t('admin.config.oauth2ResetCount') }}
            </el-button>
          </div>
        </el-form-item>

        <el-divider content-position="left">
          {{ t('admin.config.oauth2LevelMapping') }}
        </el-divider>

        <el-form-item :label="t('admin.config.oauth2LevelMapping')">
          <div class="level-mapping-container">
            <div
              v-for="(userLevel, trustLevel) in formData.levelMapping"
              :key="trustLevel"
              class="level-mapping-item"
            >
              <span class="mapping-label">Trust Level {{ trustLevel }}</span>
              <el-icon><Right /></el-icon>
              <el-select
                v-model="formData.levelMapping[trustLevel]"
                :placeholder="t('admin.config.oauth2SelectLevel')"
              >
                <el-option
                  v-for="level in availableLevels"
                  :key="level"
                  :label="t('admin.config.oauth2LevelLabel', { level })"
                  :value="level"
                />
              </el-select>
            </div>
          </div>
          <div class="form-item-tip">
            {{ t('admin.config.oauth2LevelMappingHint') }}
          </div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="saving"
            @click="handleSave"
          >
            {{ t('admin.config.oauth2SaveConfig') }}
          </el-button>
          <el-button @click="loadConfig">
            {{ t('admin.config.oauth2ResetConfig') }}
          </el-button>
        </el-form-item>
      </el-card>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Right } from '@element-plus/icons-vue'
import { getOAuth2Config, updateOAuth2Config, resetOAuth2RegistrationCount } from '@/api/config'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const providerId = ref(null)

const formRef = ref()
const loading = ref(false)
const saving = ref(false)
const availableLevels = [1, 2, 3, 4, 5]

const formData = reactive({
  enabled: false,
  clientId: '',
  clientSecret: '',
  redirectUrl: '',
  authUrl: '',
  tokenUrl: '',
  userinfoUrl: '',
  scopes: ['read', 'openid'],
  userIdField: 'id',
  usernameField: 'username',
  emailField: 'email',
  avatarField: 'avatar_url',
  trustLevelField: 'trust_level',
  maxRegistrations: 0,
  currentRegistrations: 0,
  levelMapping: {
    0: 1,
    1: 1,
    2: 1,
    3: 1,
    4: 1
  }
})

const rules = reactive({
  clientId: [
    { required: true, message: () => t('admin.oauth2.validationClientId'), trigger: 'blur' }
  ],
  clientSecret: [
    { required: true, message: () => t('admin.oauth2.validationClientSecret'), trigger: 'blur' }
  ],
  redirectUrl: [
    { required: true, message: () => t('admin.oauth2.validationRedirectUrl'), trigger: 'blur' }
  ],
  authUrl: [
    { required: true, message: () => t('admin.oauth2.validationAuthUrl'), trigger: 'blur' }
  ],
  tokenUrl: [
    { required: true, message: () => t('admin.oauth2.validationTokenUrl'), trigger: 'blur' }
  ],
  userinfoUrl: [
    { required: true, message: () => t('admin.oauth2.validationUserInfoUrl'), trigger: 'blur' }
  ]
})

const loadConfig = async () => {
  loading.value = true
  try {
    const response = await getOAuth2Config()
    if ((response.code === 200) && response.data) {
      // Response is a providers array; use the first provider as the config
      const providers = Array.isArray(response.data) ? response.data : []
      if (providers.length > 0) {
        const p = providers[0]
        providerId.value = p.id
        Object.assign(formData, p)
      } else {
        providerId.value = null
      }
      
      // 确保levelMapping是对象
      if (!formData.levelMapping || typeof formData.levelMapping !== 'object') {
        formData.levelMapping = {
          0: 1,
          1: 1,
          2: 1,
          3: 1,
          4: 1
        }
      }
    }
  } catch (error) {
    ElMessage.error(t('admin.config.oauth2LoadFailed'))
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    saving.value = true
    try {
      // 转换levelMapping的键为整数
      const levelMappingInt = {}
      Object.keys(formData.levelMapping).forEach(key => {
        levelMappingInt[parseInt(key)] = formData.levelMapping[key]
      })

      const data = {
        ...formData,
        levelMapping: levelMappingInt
      }

      const response = await updateOAuth2Config(providerId.value, data)
      ElMessage.success(t('admin.config.oauth2SaveSuccess'))
      await loadConfig()
    } catch (error) {
      ElMessage.error(error?.message || t('admin.config.oauth2SaveFailed'))
      console.error(error)
    } finally {
      saving.value = false
    }
  })
}

const resetRegistrationCount = async () => {
  try {
    await ElMessageBox.confirm(
      t('admin.config.oauth2ResetCountConfirm'),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )

    const response = await resetOAuth2RegistrationCount(providerId.value)
    ElMessage.success(t('admin.config.oauth2ResetCountSuccess'))
    await loadConfig()
  } catch (error) {
    if (error !== 'cancel' && error?.action !== 'cancel' && error?.action !== 'close') {
      ElMessage.error(error?.message || t('admin.config.oauth2ResetCountFailed'))
      console.error(error)
    }
  }
}

onMounted(() => {
  loadConfig()
})
</script>

<style scoped>
.oauth2-config-container {
  padding: 0;
}

.oauth-card {
  margin-bottom: 20px;
}

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

.form-item-tip {
  font-size: 12px;
  color: var(--text-color-secondary);
  margin-top: 5px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.level-mapping-container {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.level-mapping-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  background: var(--neutral-bg);
  border-radius: 4px;
}

.mapping-label {
  min-width: 120px;
  font-weight: 500;
}

:deep(.el-divider__text) {
  font-weight: 600;
  color: var(--text-color-primary);
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
