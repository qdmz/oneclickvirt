<template>
  <div class="reset-password-container">
    <header class="auth-header">
      <div class="header-content">
        <div class="logo">
          <img
            :src="siteConfigs.site_icon_url || logoUrl"
            :alt="siteConfigs.site_name || 'OneClickVirt Logo'"
            class="logo-image"
          >
          <a
            :href="siteConfigs.site_url || '#'"
            target="_self"
            class="site-name-link"
          >
            <h1>{{ siteConfigs.site_name || 'OneClickVirt' }}</h1>
          </a>
        </div>
        <nav class="nav-actions">
          <button
            class="nav-link language-btn"
            @click="switchLanguage"
          >
            <el-icon><Operation /></el-icon>
            {{ languageStore.currentLanguage === 'zh-CN' ? 'English' : '中文' }}
          </button>
          <router-link
            to="/"
            class="nav-link home-btn"
          >
            <el-icon><HomeFilled /></el-icon>
            {{ t('common.backToHome') }}
          </router-link>
        </nav>
      </div>
    </header>

    <div class="reset-password-form">
      <div v-if="!resetSuccess">
        <h2>{{ t('resetPassword.title') }}</h2>
        <p>{{ t('resetPassword.subtitle') }}</p>

        <el-form
          ref="resetFormRef"
          :model="resetForm"
          :rules="resetRules"
          label-width="0"
          size="large"
        >
          <el-form-item prop="newPassword">
            <el-input
              v-model="resetForm.newPassword"
              type="password"
              :placeholder="t('resetPassword.pleaseEnterNewPassword')"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>

          <el-form-item prop="confirmPassword">
            <el-input
              v-model="resetForm.confirmPassword"
              type="password"
              :placeholder="t('resetPassword.pleaseConfirmNewPassword')"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              style="width: 100%;"
              @click="handleResetPassword"
            >
              {{ t('resetPassword.resetButton') }}
            </el-button>
          </el-form-item>

          <div class="form-footer">
            <router-link to="/login">
              {{ t('forgotPassword.backToLogin') }}
            </router-link>
          </div>
        </el-form>
      </div>

      <div v-else class="success-message">
        <el-result
          icon="success"
          :title="t('resetPassword.resetSuccess')"
          sub-title=""
        >
          <template #extra>
            <el-button
              type="primary"
              @click="goToLogin"
            >
              {{ t('forgotPassword.backToLogin') }}
            </el-button>
          </template>
        </el-result>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { resetPassword } from '@/api/auth'
import { Operation, HomeFilled } from '@element-plus/icons-vue'
import { useLanguageStore } from '@/pinia/modules/language'
import { getPublicSiteConfigs } from '@/api/public'
import logoUrl from '@/assets/images/logo.png'

const router = useRouter()
const route = useRoute()
const { t, locale } = useI18n()
const languageStore = useLanguageStore()
const resetFormRef = ref()
const loading = ref(false)
const resetSuccess = ref(false)
const siteConfigs = ref({})

const resetForm = reactive({
  newPassword: '',
  confirmPassword: ''
})

const resetRules = computed(() => ({
  newPassword: [
    { required: true, message: t('validation.passwordRequired'), trigger: 'blur' },
    { min: 6, message: t('validation.passwordMinLength'), trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: t('validation.passwordRequired'), trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== resetForm.newPassword) {
          callback(new Error(t('validation.passwordMismatch')))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}))

const handleResetPassword = async () => {
  if (!resetFormRef.value) return

  await resetFormRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const token = route.query.token || route.hash.split('token=')[1]?.split('&')[0]
      const response = await resetPassword({
        token,
        newPassword: resetForm.newPassword
      })

      if (response.code === 0 || response.code === 200) {
        resetSuccess.value = true
      } else {
        ElMessage.error(response.message || t('resetPassword.resetFailed'))
      }
    } catch (error) {
      console.error(t('resetPassword.resetFailed'), error)
      ElMessage.error(t('resetPassword.resetFailed'))
    } finally {
      loading.value = false
    }
  })
}

const goToLogin = () => {
  router.push('/login')
}

const switchLanguage = () => {
  const newLang = languageStore.toggleLanguage()
  locale.value = newLang
  ElMessage.success(t('navbar.languageSwitched'))
}

const fetchSiteConfigs = async () => {
  try {
    const resp = await getPublicSiteConfigs()
    if (resp && (resp.code === 0 || resp.code === 200) && resp.data) {
      const configs = resp.data
      if (Array.isArray(configs)) {
        configs.forEach(config => {
          siteConfigs.value[config.key] = config.value
        })
      } else {
        siteConfigs.value = { ...siteConfigs.value, ...configs }
      }
    }
  } catch (error) {
    console.error('获取站点配置失败', error)
  }
}

onMounted(() => {
  fetchSiteConfigs()
})
</script>

<style scoped>
.reset-password-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: var(--auth-page-bg);
}

.auth-header {
  background: var(--auth-header-bg);
  backdrop-filter: blur(20px);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
  border-bottom: 1px solid var(--auth-card-border);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 70px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-image {
  width: 48px;
  height: 48px;
  object-fit: contain;
}

.logo h1 {
  font-size: 28px;
  color: #fff;
  margin: 0;
  font-weight: 700;
  background: linear-gradient(135deg, #fff, rgba(255,255,255,0.8));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.nav-link {
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px 24px;
  border-radius: 25px;
  border: 1px solid #e5e7eb;
  background: transparent;
  color: #374151;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-link:hover {
  background: rgba(22, 163, 74, 0.1);
  color: #fff;
  transform: translateY(-2px);
}

.nav-link.home-btn {
  background: linear-gradient(135deg, #fff, rgba(255,255,255,0.8));
  color: white;
  border: none;
  box-shadow: 0 4px 15px rgba(22, 163, 74, 0.3);
}

.nav-link.home-btn:hover {
  background: linear-gradient(135deg, #15803d, #16a34a);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(22, 163, 74, 0.4);
}

.reset-password-form {
  margin: auto;
  margin-top: 60px;
  margin-bottom: 60px;
  width: 400px;
  padding: 40px;
  background: var(--auth-card-bg);
  backdrop-filter: blur(24px);
  border-radius: var(--border-radius-xl);
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
}

.reset-password-form h2 {
  font-size: 24px;
  color: #303133;
  margin-bottom: 10px;
  text-align: center;
}

.reset-password-form p {
  font-size: 14px;
  color: #909399;
  margin-bottom: 30px;
  text-align: center;
}

.form-footer {
  text-align: center;
  margin-top: 20px;
}

.form-footer a {
  color: var(--primary-color-light);
  text-decoration: none;
}

.success-message {
  text-align: center;
}

@media (max-width: 768px) {
  .reset-password-form {
    width: 90%;
    padding: 20px;
  }
}
</style>
