<template>
  <div class="forgot-password-container">
    <!-- 顶部栏 -->
    <header class="auth-header">
      <div class="header-content">
        <div class="logo">
          <img
            :src="siteStore.logoSrc"
            alt="OneClickVirt Logo"
            class="logo-image"
          >
          <h1>{{ siteStore.displaySiteName }}</h1>
        </div>
        <nav class="nav-actions">
          <button
            class="nav-link theme-btn"
            :title="themeStore.isDark ? t('navbar.lightMode') : t('navbar.darkMode')"
            @click="toggleTheme"
          >
            <el-icon><component :is="themeStore.isDark ? Sunny : Moon" /></el-icon>
          </button>
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

    <div class="forgot-password-form">
      <div v-if="!emailSent">
        <h2>{{ t('forgotPassword.title') }}</h2>
        <p>{{ t('forgotPassword.subtitle') }}</p>

        <el-form
          ref="forgotFormRef"
          :model="forgotForm"
          :rules="forgotRules"
          label-width="0"
          size="large"
        >
          <el-form-item prop="email">
            <el-input
              v-model="forgotForm.email"
              :placeholder="t('forgotPassword.pleaseEnterEmail')"
              prefix-icon="Message"
            />
          </el-form-item>

          <el-form-item
            v-if="captchaEnabled"
            prop="captcha"
          >
            <div class="captcha-container">
              <el-input
                v-model="forgotForm.captcha"
                :placeholder="t('login.pleaseEnterCaptcha')"
                style="width: 60%"
              />
              <div
                class="captcha-image"
                @click="refreshCaptcha"
              >
                <img
                  v-if="captchaImage"
                  :src="captchaImage"
                  :alt="t('login.captchaAlt')"
                >
                <div
                  v-else
                  class="captcha-loading"
                >
                  {{ t('common.loading') }}
                </div>
              </div>
            </div>
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              style="width: 100%;"
              @click="handleForgotPassword"
            >
              {{ t('forgotPassword.sendResetLink') }}
            </el-button>
          </el-form-item>

          <div class="form-footer">
            <router-link to="/login">
              {{ t('forgotPassword.backToLogin') }}
            </router-link>
          </div>
        </el-form>
      </div>

      <div
        v-else
        class="success-message"
      >
        <el-result
          icon="success"
          :title="t('forgotPassword.emailSent')"
          :sub-title="t('forgotPassword.checkEmail')"
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
    <AppFooter />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { forgotPassword } from '@/api/auth'
import { getCaptcha } from '@/api/auth'
import { getPublicConfig } from '@/api/public'
import { Operation, HomeFilled, Sunny, Moon } from '@element-plus/icons-vue'
import { useLanguageStore } from '@/pinia/modules/language'
import { useThemeStore } from '@/pinia/modules/theme'
import { useSiteStore } from '@/pinia/modules/site'
import AppFooter from '@/view/layout/components/AppFooter.vue'

const router = useRouter()
const { t, locale } = useI18n()
const languageStore = useLanguageStore()
const themeStore = useThemeStore()
const siteStore = useSiteStore()
const forgotFormRef = ref()
const loading = ref(false)
const emailSent = ref(false)
const captchaImage = ref('')
const captchaId = ref('')
const captchaEnabled = ref(false)

const forgotForm = reactive({
  email: '',
  captcha: ''
})

const forgotRules = computed(() => ({
  email: [
    { required: true, message: t('validation.emailRequired'), trigger: 'blur' },
    { type: 'email', message: t('validation.emailFormat'), trigger: 'blur' }
  ],
  ...(captchaEnabled.value ? {
    captcha: [
      { required: true, message: t('validation.captchaRequired'), trigger: 'blur' }
    ]
  } : {})
}))

const handleForgotPassword = async () => {
  if (!forgotFormRef.value) return

  await forgotFormRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const response = await forgotPassword({
        email: forgotForm.email,
        ...(captchaEnabled.value ? {
          captcha: forgotForm.captcha,
          captchaId: captchaId.value
        } : {
          captcha: undefined,
          captchaId: undefined
        })
      })

      if (response.code === 200) {
        emailSent.value = true
      }
    } catch (error) {
      console.error(t('forgotPassword.resetFailed'), error)
      ElMessage.error(t('forgotPassword.resetFailed'))
      if (captchaEnabled.value) {
        refreshCaptcha()
      }
    } finally {
      loading.value = false
    }
  })
}

const refreshCaptcha = async () => {
  if (!captchaEnabled.value) {
    clearCaptchaState()
    return
  }

  try {
    const response = await getCaptcha()
    if (response.code === 200) {
      captchaImage.value = response.data.imageData
      captchaId.value = response.data.captchaId
      forgotForm.captcha = ''
    }
  } catch (error) {
    console.error(t('forgotPassword.captchaFailed'), error)
  }
}

const clearCaptchaState = () => {
  captchaImage.value = ''
  captchaId.value = ''
  forgotForm.captcha = ''
}

const loadCaptchaConfig = async () => {
  try {
    const response = await getPublicConfig()
    captchaEnabled.value = response.data?.captchaEnabled ?? false
  } catch (error) {
    captchaEnabled.value = false
  }
}

const goToLogin = () => {
  router.push('/login')
}

// 切换语言
const switchLanguage = () => {
  const newLang = languageStore.toggleLanguage()
  locale.value = newLang
  ElMessage.success(t('navbar.languageSwitched'))
}

const toggleTheme = () => {
  themeStore.toggleTheme()
}

onMounted(async () => {
  await loadCaptchaConfig()
  if (captchaEnabled.value) {
    refreshCaptcha()
  } else {
    clearCaptchaState()
  }
})
</script>

<style scoped>
.forgot-password-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  min-height: 100dvh;
  background: var(--auth-page-bg);
}

/* 顶部栏样式 */
.auth-header {
  background: var(--auth-header-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  box-shadow: var(--box-shadow-light);
  border-bottom: 1px solid var(--border-color);
  padding-top: env(safe-area-inset-top);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  min-height: 60px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-image {
  width: 42px;
  height: 42px;
  object-fit: contain;
}

.logo h1 {
  font-size: 24px;
  color: var(--primary-color);
  margin: 0;
  font-weight: 700;
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.nav-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 8px;
}

.nav-link.theme-btn {
  padding: 8px 10px;
  min-width: 38px;
  justify-content: center;
}

.nav-link {
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 9px 14px;
  border-radius: 22px;
  border: 1px solid var(--border-color);
  background: transparent;
  color: var(--text-color-primary);
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-link:hover {
  background: var(--primary-color-bg-hover);
  color: var(--accent-text-color);
  transform: translateY(-2px);
}

.nav-link.home-btn {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-light));
  color: white;
  border: none;
  box-shadow: 0 4px 15px var(--primary-color-shadow);
}

.nav-link.home-btn:hover {
  background: linear-gradient(135deg, var(--primary-color-dark), var(--primary-color));
  transform: translateY(-2px);
  box-shadow: 0 6px 20px var(--primary-color-shadow-hover);
}

.forgot-password-form {
  margin: 32px auto;
  width: min(420px, calc(100% - 32px));
  padding: 36px 38px;
  background: var(--card-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 20px;
  box-shadow: var(--box-shadow-heavy);
  border: 1px solid var(--border-color);
}

.forgot-password-form h2 {
  font-size: 26px;
  font-weight: 700;
  color: var(--text-color-primary);
  margin-bottom: 10px;
  text-align: center;
}

.forgot-password-form p {
  font-size: 14px;
  color: var(--text-color-secondary);
  margin-bottom: 30px;
  text-align: center;
}

.form-footer {
  text-align: center;
  margin-top: 20px;
}

.form-footer a {
  color: var(--accent-text-color);
  text-decoration: none;
  font-weight: 500;
}

.form-footer a:hover {
  color: var(--accent-text-color-hover);
}

.success-message {
  text-align: center;
}

.captcha-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.captcha-image {
  width: 38%;
  height: 40px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  overflow: hidden;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.captcha-loading {
  font-size: 12px;
  color: var(--text-color-tertiary);
}

@media (max-width: 768px) {
  .header-content {
    min-height: 56px;
    padding: 8px 16px;
    gap: 10px;
  }

  .logo {
    gap: 8px;
  }

  .logo-image {
    width: 38px;
    height: 38px;
  }

  .logo h1 {
    font-size: 21px;
  }

  .nav-actions {
    gap: 6px;
  }

  .nav-link {
    padding: 8px 10px;
    font-size: 14px;
  }

  .forgot-password-form {
    width: calc(100% - 24px);
    margin: 24px auto;
    padding: 24px;
  }
}

@media (max-width: 480px) {
  .header-content {
    justify-content: center;
    padding: 8px 12px;
  }

  .logo {
    width: 100%;
    justify-content: center;
  }

  .logo-image {
    width: 34px;
    height: 34px;
  }

  .logo h1 {
    font-size: 20px;
  }

  .nav-actions {
    width: 100%;
    justify-content: center;
  }

  .nav-link {
    padding: 7px 9px;
    font-size: 13px;
    border-radius: 18px;
  }

  .forgot-password-form {
    padding: 22px 18px;
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
