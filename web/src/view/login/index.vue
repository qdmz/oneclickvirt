<template>
  <div class="login-container">
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

    <div class="login-form">
      <div class="login-header">
        <h2>{{ t('login.title') }}</h2>
        <p>{{ t('login.subtitle') }}</p>
      </div>

      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        label-width="0"
        size="large"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            :placeholder="t('login.pleaseEnterUsername')"
            prefix-icon="User"
            clearable
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            :placeholder="t('login.pleaseEnterPassword')"
            prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item
          v-if="captchaEnabled"
          prop="captcha"
        >
          <div class="captcha-container">
            <el-input
              v-model="loginForm.captcha"
              :placeholder="t('login.pleaseEnterCaptcha')"
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

        <div class="form-options">
          <el-checkbox v-model="loginForm.rememberMe">
            {{ t('login.rememberMe') }}
          </el-checkbox>
          <router-link
            to="/forgot-password"
            class="forgot-link"
          >
            {{ t('login.forgotPassword') }}
          </router-link>
        </div>

        <div class="form-actions">
          <el-button
            type="primary"
            :loading="loading"
            style="width: 100%;"
            @click="handleLogin"
          >
            {{ t('common.login') }}
          </el-button>
        </div>

        <div class="form-footer">
          <p>
            {{ t('login.noAccount') }} <router-link to="/register">
              {{ t('login.registerNow') }}
            </router-link>
          </p>
        </div>

        <div class="admin-login">
          <router-link
            to="/admin/login"
            class="admin-link"
          >
            {{ t('login.adminLogin') }}
          </router-link>
        </div>
      </el-form>

      <!-- OAuth2登录 -->
      <div
        v-if="oauth2Enabled && oauth2Providers.length > 0"
        class="oauth2-login"
      >
        <el-divider>{{ t('login.thirdPartyLogin') }}</el-divider>
        <div class="oauth2-providers">
          <el-button
            v-for="provider in oauth2Providers"
            :key="provider.id"
            class="oauth2-button"
            :loading="oauth2Loading"
            :disabled="oauth2Loading"
            @click="handleOAuth2Login(provider)"
          >
            <el-icon><Connection /></el-icon>
            {{ provider.displayName }}
          </el-button>
        </div>
      </div>
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/pinia/modules/user'
import { getCaptcha } from '@/api/auth'
import { useErrorHandler } from '@/composables/useErrorHandler'
import { getPublicConfig } from '@/api/public'
import { getEnabledOAuth2Providers } from '@/api/oauth2'
import { Connection, Operation, HomeFilled, Sunny, Moon } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useLanguageStore } from '@/pinia/modules/language'
import { useThemeStore } from '@/pinia/modules/theme'
import { useSiteStore } from '@/pinia/modules/site'
import AppFooter from '@/view/layout/components/AppFooter.vue'

const router = useRouter()
const userStore = useUserStore()
const { t, locale } = useI18n()
const { executeAsync, handleSubmit } = useErrorHandler()
const languageStore = useLanguageStore()
const themeStore = useThemeStore()
const siteStore = useSiteStore()

const loginFormRef = ref()
const loading = ref(false)
const captchaImage = ref('')
const captchaId = ref('')
const oauth2Enabled = ref(false)
const oauth2Providers = ref([])
const oauth2Loading = ref(false) // OAuth2登录防重复点击
const captchaEnabled = ref(false)

const loginForm = reactive({
  username: '',
  password: '',
  captcha: '',
  rememberMe: false,
  userType: 'user',
  loginType: 'password'
})

const loginRules = computed(() => ({
  username: [
    { required: true, message: t('validation.usernameRequired'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('validation.passwordRequired'), trigger: 'blur' }
  ],
  ...(captchaEnabled.value ? {
    captcha: [
      { required: true, message: t('validation.captchaRequired'), trigger: 'blur' }
    ]
  } : {})
}))

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  // 防止重复提交
  if (loading.value) return

  await loginFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    // 再次检查loading状态，防止表单验证期间的重复点击
    if (loading.value) return
    
    loading.value = true
    
    try {
      const result = await handleSubmit(async () => {
        return await userStore.userLogin({
          ...loginForm,
          ...(captchaEnabled.value ? { captchaId: captchaId.value } : { captcha: undefined, captchaId: undefined })
        })
      }, {
        successMessage: t('login.loginSuccess'),
        showLoading: false // 使用组件自己的loading
      })

      if (result.success) {
        // 二次验证：确认 userStore 实际拥有 token 和用户数据
        // 因为 executeAsync 会捕获抛出的错误并返回 { success: false }，
        // 但如果 userLogin 返回 { success: false } 而未抛出异常，
        // executeAsync 会将其视为成功，因此需要额外检查
        if (!userStore.token || !userStore.user) {
          // 登录实际未成功，尽管没有抛出异常
          if (captchaEnabled.value) {
            refreshCaptcha()
          }
          return
        }

        // 根据用户类型和视图模式跳转
        const userType = userStore.userType
        const viewMode = userStore.viewMode || userType
        
        // 管理员（包括普通管理员）可以访问管理员界面
        if ((userType === 'admin' || userType === 'normal_admin') && (viewMode === 'admin' || viewMode === 'normal_admin')) {
          router.push('/admin/dashboard')
        } else {
          // 普通用户或管理员的用户视图
          router.push('/user/dashboard')
        }
      } else {
        if (captchaEnabled.value) {
          refreshCaptcha() // 登录失败刷新验证码
        }
      }
    } finally {
      loading.value = false
    }
  })
}

const refreshCaptcha = async () => {
  if (!captchaEnabled.value) {
    captchaImage.value = ''
    captchaId.value = ''
    loginForm.captcha = ''
    return
  }

  await executeAsync(async () => {
    const response = await getCaptcha()
    captchaImage.value = response.data.imageData
    captchaId.value = response.data.captchaId
    loginForm.captcha = ''
  }, {
    showError: false, // 静默处理验证码错误
    showLoading: false
  })
}

// OAuth2登录
const handleOAuth2Login = (provider) => {
  // 防止重复点击
  if (oauth2Loading.value) return
  
  oauth2Loading.value = true
  
  // 跳转到后端的OAuth2登录接口，使用provider_id参数
  window.location.href = `/api/v1/auth/oauth2/login?provider_id=${provider.id}`
  
  // 页面跳转后loading状态会自动重置，这里不需要手动重置
}

// 检查OAuth2配置并加载提供商列表
const checkOAuth2Config = async () => {
  try {
    // 获取OAuth2全局开关状态
    const configResponse = await getPublicConfig()
    oauth2Enabled.value = configResponse.data?.oauth2Enabled || false
    captchaEnabled.value = configResponse.data?.captchaEnabled || false
    
    // 如果启用了OAuth2，加载提供商列表
    if (oauth2Enabled.value) {
      const providersResponse = await getEnabledOAuth2Providers()
      oauth2Providers.value = providersResponse.data || []
    }
  } catch (error) {
    console.error(t('login.getOAuth2ConfigFailed'), error)
  }
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
  await checkOAuth2Config()
  if (captchaEnabled.value) {
    refreshCaptcha()
  }
})
</script>

<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  min-height: 100dvh;
  background: var(--auth-page-bg);
  position: relative;
  overflow: hidden;
}

/* 装饰性背景光晕 */
.login-container::before {
  content: '';
  position: absolute;
  top: -20%;
  left: -10%;
  width: 500px;
  height: 500px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.12) 0%, transparent 70%);
  pointer-events: none;
  animation: floatGlow 8s ease-in-out infinite;
}

.login-container::after {
  content: '';
  position: absolute;
  bottom: -20%;
  right: -10%;
  width: 600px;
  height: 600px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(168, 85, 247, 0.1) 0%, transparent 70%);
  pointer-events: none;
  animation: floatGlow 10s ease-in-out infinite reverse;
}

@keyframes floatGlow {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(30px, -30px) scale(1.05); }
}

/* 顶部栏样式 */
.auth-header {
  background: var(--auth-header-bg);
  backdrop-filter: blur(24px) saturate(1.2);
  -webkit-backdrop-filter: blur(24px) saturate(1.2);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
  border-bottom: 1px solid var(--border-color);
  padding-top: env(safe-area-inset-top);
  position: relative;
  z-index: 10;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  min-height: 64px;
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
  filter: drop-shadow(0 4px 12px rgba(99, 102, 241, 0.2));
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.logo:hover .logo-image {
  transform: scale(1.08) rotate(-3deg);
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
  padding: 9px 16px;
  border-radius: 999px;
  border: 1px solid var(--border-color);
  background: transparent;
  color: var(--text-color-primary);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.nav-link:hover {
  background: var(--primary-color-bg-hover);
  color: var(--accent-text-color);
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.12);
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
  box-shadow: 0 8px 24px var(--primary-color-shadow-hover);
}

.login-form {
  margin: 40px auto;
  width: min(440px, calc(100% - 32px));
  padding: 40px 36px;
  background: var(--card-bg);
  backdrop-filter: blur(28px) saturate(1.3);
  -webkit-backdrop-filter: blur(28px) saturate(1.3);
  border-radius: 24px;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.12), 0 2px 8px rgba(0, 0, 0, 0.06);
  border: 1px solid var(--border-color);
  position: relative;
  z-index: 10;
  animation: fadeInUp 0.6s cubic-bezier(0.4, 0, 0.2, 1) both;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(24px); }
  to { opacity: 1; transform: translateY(0); }
}

.login-form :deep(.el-form) {
  width: 100%;
}

.login-form :deep(.el-form-item) {
  width: 100%;
  margin-bottom: 20px;
}

.login-form :deep(.el-form-item__content) {
  width: 100%;
  line-height: normal;
}

.login-form :deep(.el-input) {
  width: 100%;
}

.login-form :deep(.el-input__wrapper) {
  width: 100%;
  box-sizing: border-box;
  border-radius: 14px !important;
  padding: 4px 14px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.login-form :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--primary-color) inset, 0 0 0 4px rgba(99, 102, 241, 0.12) !important;
}

.login-form :deep(.el-button--primary) {
  height: 48px;
  font-size: 16px;
  font-weight: 700;
  border-radius: 14px;
  letter-spacing: 0.5px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.login-form :deep(.el-button--primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px var(--primary-color-shadow), 0 0 20px rgba(99, 102, 241, 0.2);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h2 {
  font-size: 28px;
  font-weight: 800;
  margin-bottom: 10px;
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.5px;
}

.login-header p {
  font-size: 14px;
  color: var(--text-color-secondary);
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  width: 100%;
}

.forgot-link {
  color: var(--accent-text-color);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;
}

.forgot-link:hover {
  color: var(--accent-text-color-hover);
}

.form-actions {
  margin-bottom: 20px;
  width: 100%;
}

.form-actions .el-button {
  width: 100% !important;
  height: 45px;
}

.form-footer {
  text-align: center;
  margin-bottom: 20px;
  width: 100%;
  color: var(--text-color-secondary);
}

.form-footer a {
  color: var(--accent-text-color);
  text-decoration: none;
  font-weight: 500;
}

.form-footer a:hover {
  color: var(--accent-text-color-hover);
}

.admin-login {
  text-align: center;
  font-size: 14px;
  color: var(--text-color-secondary);
}

.admin-link {
  color: var(--text-color-tertiary);
  text-decoration: none;
  margin: 0 5px;
}

.admin-link:hover {
  color: var(--accent-text-color);
}

.captcha-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  width: 100%;
}

.captcha-container .el-input {
  flex: 1;
}

.captcha-image {
  width: 120px;
  height: 40px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  overflow: hidden;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
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

.oauth2-login {
  margin: 20px 0 0 0;
  width: 100%;
  padding: 0;
}

.oauth2-login :deep(.el-divider) {
  margin: 20px 0;
}

.oauth2-providers {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
  padding: 0;
  margin: 0;
}

.oauth2-button {
  width: 100% !important;
  height: 45px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border-color);
  background: var(--surface-input);
  color: var(--text-color-primary);
  margin: 0 !important;
  padding: 0 20px !important;
  box-sizing: border-box;
}

.oauth2-button:hover {
  border-color: var(--border-color-hover);
  color: var(--accent-text-color);
}

.oauth2-providers :deep(.el-button) {
  width: 100% !important;
  margin: 0 !important;
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

  .login-form {
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

  .login-form {
    padding: 22px 18px;
  }
}
</style>
