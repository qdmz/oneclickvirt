<template>
  <div class="login-container">
    <!-- 顶部栏 -->
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
            <h1>{{ siteConfigs.site_name || t('login.title') }}</h1>
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

    <div class="login-form">
      <div class="login-header">
        <h2>{{ t('login.title') }}</h2>
        <p>{{ t('login.subtitle') }}</p>
      </div>

      <!-- 邮箱未验证提示 -->
      <el-alert
        v-if="emailNotVerified"
        :title="emailNotVerifiedMsg"
        type="warning"
        show-icon
        :closable="false"
        style="margin-bottom: 20px;"
      >
        <template #default>
          <el-button
            type="text"
            size="small"
            :loading="resendLoading"
            @click="handleResendVerification"
          >
            重新发送激活邮件
          </el-button>
        </template>
      </el-alert>

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

        <el-form-item prop="captcha">
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/pinia/modules/user'
import { getCaptcha } from '@/api/auth'
import { useErrorHandler } from '@/composables/useErrorHandler'
import { getPublicConfig, getPublicSiteConfigs } from '@/api/public'
import { getEnabledOAuth2Providers } from '@/api/oauth2'
import { resendVerification } from '@/api/auth'
import { Connection, Operation, HomeFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useLanguageStore } from '@/pinia/modules/language'
import logoUrl from '@/assets/images/logo.png'

const router = useRouter()
const userStore = useUserStore()
const { t, locale } = useI18n()
const { executeAsync, handleSubmit } = useErrorHandler()
const languageStore = useLanguageStore()

const loginFormRef = ref()
const loading = ref(false)
const captchaImage = ref('')
const captchaId = ref('')
const oauth2Enabled = ref(false)
const oauth2Providers = ref([])
const oauth2Loading = ref(false) // OAuth2登录防重复点击
const siteConfigs = ref({}) // 站点配置
const emailNotVerified = ref(false)
const emailNotVerifiedMsg = ref('')
const resendLoading = ref(false)

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
  captcha: [
    { required: true, message: t('validation.captchaRequired'), trigger: 'blur' }
  ]
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
        emailNotVerified.value = false
        return await userStore.userLogin({
          ...loginForm,
          captchaId: captchaId.value
        })
      }, {
        successMessage: t('login.loginSuccess'),
        showLoading: false // 使用组件自己的loading
      })

      if (result.success) {
        // 根据用户类型跳转到相应的仪表盘
        if (userStore.userType === 'admin') {
          router.push('/admin/dashboard')
        } else if (userStore.userType === 'agent') {
          router.push('/agent/dashboard')
        } else {
          router.push('/user/dashboard')
        }
      } else {
        // 检查是否是邮箱未验证错误 (code 4009)
        const errData = result.error?.response?.data
        if (errData && errData.code === 4009) {
          emailNotVerified.value = true
          emailNotVerifiedMsg.value = errData.message || '邮箱未验证，请先验证邮箱'
        }
        refreshCaptcha() // 登录失败刷新验证码
      }
    } finally {
      loading.value = false
    }
  })
}

const refreshCaptcha = async () => {
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
  ElMessage.success(t('common.languageChanged'))
}

// 获取站点配置
const fetchSiteConfigs = async () => {
  try {
    const resp = await getPublicSiteConfigs()
    if (resp && (resp.code === 0 || resp.code === 200) && resp.data) {
      const configs = resp.data
      // 检查数据格式，如果是对象直接使用，如果是数组则遍历
      if (Array.isArray(configs)) {
        configs.forEach(config => {
          siteConfigs.value[config.key] = config.value
        })
      } else {
        // 直接将对象赋值给siteConfigs
        siteConfigs.value = { ...siteConfigs.value, ...configs }
      }
    }
  } catch (error) {
    console.error('获取站点配置失败', error)
  }
}

// 重新发送激活邮件
const handleResendVerification = async () => {
  if (!loginForm.username) {
    ElMessage.warning(t('login.pleaseEnterUsername'))
    return
  }
  // 尝试用用户名查找邮箱，简单处理：提示用户输入邮箱
  // 这里用一个简单的方式：用用户名作为查询
  resendLoading.value = true
  try {
    const resp = await resendVerification({ email: loginForm.username })
    if (resp.code === 0 || resp.code === 200) {
      ElMessage.success('激活邮件已重新发送，请查收')
      emailNotVerified.value = false
    } else {
      ElMessage.error(resp.message || '发送失败')
    }
  } catch (error) {
    ElMessage.error('发送激活邮件失败')
  } finally {
    resendLoading.value = false
  }
}

onMounted(() => {
  refreshCaptcha()
  checkOAuth2Config()
  fetchSiteConfigs()
})
</script>

<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: var(--auth-page-bg);
  position: relative;
}

/* 顶部栏样式 */
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
  width: 40px;
  height: 40px;
  object-fit: contain;
}

.site-name-link {
  text-decoration: none;
}

.logo h1 {
  font-size: 22px;
  color: #fff;
  margin: 0;
  font-weight: 700;
  background: linear-gradient(135deg, #fff, rgba(255,255,255,0.8));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  transition: all 0.3s ease;
}

.logo h1:hover {
  transform: scale(1.05);
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
  padding: 10px 20px;
  border-radius: var(--border-radius-sm);
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.nav-link:hover {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  transform: translateY(-2px);
}

.nav-link.home-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.nav-link.home-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.login-form {
  margin: auto;
  margin-top: 60px;
  margin-bottom: 60px;
  width: 420px;
  padding: 40px;
  background: var(--auth-card-bg);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-radius: var(--border-radius-xl);
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
  border: 1px solid var(--auth-card-border);
  animation: fadeIn 0.6s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
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
  border-radius: var(--border-radius-sm);
  background: rgba(255, 255, 255, 0.08) !important;
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.15) inset !important;
}

.login-form :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.25) inset !important;
}

.login-form :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--primary-color-light) inset, 0 0 0 3px rgba(99, 102, 241, 0.2) !important;
}

.login-form :deep(.el-input__inner) {
  color: var(--text-color-primary) !important;
}

.login-form :deep(.el-input__prefix .el-icon) {
  color: rgba(255, 255, 255, 0.5) !important;
}

.login-form :deep(.el-checkbox__label) {
  color: var(--text-color-secondary) !important;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
  font-size: 24px;
  color: var(--text-color-primary);
  margin-bottom: 10px;
  font-weight: 700;
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
  color: var(--primary-color-light);
  text-decoration: none;
}

.form-actions {
  margin-bottom: 20px;
  width: 100%;
}

.form-actions :deep(.el-button) {
  width: 100% !important;
  height: 45px;
  background: linear-gradient(135deg, #6366F1, #8B5CF6) !important;
  border: none !important;
  border-radius: var(--border-radius-sm) !important;
  font-size: 16px;
  font-weight: 600;
  box-shadow: 0 4px 15px rgba(99, 102, 241, 0.4) !important;
  transition: all 0.3s ease !important;
}

.form-actions :deep(.el-button:hover) {
  background: linear-gradient(135deg, #4F46E5, #7C3AED) !important;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(99, 102, 241, 0.5) !important;
}

.form-footer {
  text-align: center;
  margin-bottom: 20px;
  width: 100%;
}

.form-footer p {
  color: var(--text-color-secondary);
}

.form-footer a {
  color: var(--primary-color-light);
  text-decoration: none;
  font-weight: 500;
}

.admin-login {
  text-align: center;
  font-size: 14px;
}

.admin-link {
  color: var(--text-color-secondary);
  text-decoration: none;
  margin: 0 5px;
}

.admin-link:hover {
  color: var(--primary-color-light);
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
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: var(--border-radius-sm);
  overflow: hidden;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: rgba(255, 255, 255, 0.05);
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.captcha-loading {
  font-size: 12px;
  color: var(--text-color-secondary);
}

.oauth2-login {
  margin: 20px 0 0 0;
  width: 100%;
  padding: 0;
}

.oauth2-login :deep(.el-divider) {
  margin: 20px 0;
}

.oauth2-login :deep(.el-divider__text) {
  background: transparent !important;
  color: var(--text-color-secondary) !important;
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
  border: 1px solid rgba(255, 255, 255, 0.15) !important;
  background: rgba(255, 255, 255, 0.08) !important;
  color: var(--text-color-primary) !important;
  margin: 0 !important;
  padding: 0 20px !important;
  box-sizing: border-box;
  border-radius: var(--border-radius-sm) !important;
  backdrop-filter: blur(10px);
}

.oauth2-button:hover {
  border-color: var(--primary-color-light) !important;
  background: rgba(99, 102, 241, 0.15) !important;
  color: var(--primary-color-light) !important;
}

.oauth2-providers :deep(.el-button) {
  width: 100% !important;
  margin: 0 !important;
}

/* Alert override for auth pages */
.login-form :deep(.el-alert) {
  --el-alert-bg-color: rgba(245, 158, 11, 0.1);
}

@media (max-width: 768px) {
  .login-form {
    width: 90%;
    padding: 24px;
  }
}
</style>