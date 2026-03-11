<template>
  <div class="register-container">
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
            <h1>{{ siteConfigs.site_name || t('register.title') }}</h1>
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

    <!-- 注册被禁用的提示 -->
    <div
      v-if="!registrationEnabled"
      class="registration-disabled"
    >
      <el-card>
        <div class="disabled-content">
          <el-icon
            size="60"
            color="#e6a23c"
          >
            <Warning />
          </el-icon>
          <h2>{{ t('register.disabled.title') }}</h2>
          <p>{{ t('register.disabled.message') }}</p>
          <el-button
            type="primary"
            @click="router.push('/login')"
          >
            {{ t('register.disabled.backToLogin') }}
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 正常注册表单 -->
    <div
      v-else-if="!emailVerificationSent"
      class="register-form"
    >
      <div class="register-header">
        <h2>{{ t('register.title') }}</h2>
        <p>{{ t('register.subtitle') }}</p>
      </div>

      <el-form 
        ref="registerFormRef"
        :model="registerForm"
        :rules="registerRules"
        :label-width="locale === 'en-US' ? '140px' : '80px'"
        size="large"
      >
        <el-form-item
          :label="t('register.username')"
          prop="username"
        >
          <el-input 
            v-model="registerForm.username"
            :placeholder="t('register.pleaseEnterUsername')"
          />
        </el-form-item>

        <el-form-item
          :label="t('register.password')"
          prop="password"
        >
          <el-input 
            v-model="registerForm.password"
            type="password"
            :placeholder="t('register.pleaseEnterPassword')"
            show-password
          />
          <div class="password-hint">
            <el-text
              size="small"
              type="info"
            >
              {{ t('register.passwordHint') }}
            </el-text>
          </div>
        </el-form-item>

        <el-form-item
          :label="t('register.confirmPassword')"
          prop="confirmPassword"
        >
          <el-input 
            v-model="registerForm.confirmPassword"
            type="password"
            :placeholder="t('register.pleaseConfirmPassword')"
            show-password
          />
        </el-form-item>

        <el-form-item
          :label="t('register.captcha')"
          prop="captcha"
        >
          <div class="captcha-container">
            <el-input 
              v-model="registerForm.captcha"
              :placeholder="t('register.pleaseEnterCaptcha')"
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

        <el-form-item
          v-if="showInviteCode"
          :label="t('register.inviteCode')"
          prop="inviteCode"
        >
          <el-input 
            v-model="registerForm.inviteCode"
            :placeholder="t('register.pleaseEnterInviteCode')"
          />
        </el-form-item>

        <el-form-item>
          <el-button 
            type="primary" 
            :loading="loading" 
            style="width: 100%;"
            @click="handleRegister"
          >
            {{ t('register.registerButton') }}
          </el-button>
        </el-form-item>

        <div class="form-footer">
          <p>
            {{ t('register.hasAccount') }}<router-link to="/login">
              {{ t('register.loginNow') }}
            </router-link>
          </p>
        </div>
      </el-form>
    </div>

    <!-- 邮箱验证提示 -->
    <div
      v-if="emailVerificationSent"
      class="register-form"
    >
      <div class="register-header">
        <h2>📧 注册成功</h2>
        <p>激活邮件已发送到您的邮箱</p>
      </div>
      <el-result
        icon="success"
        title="请查看邮箱"
        :sub-title="'我们已向 ' + registerForm.email + ' 发送了激活邮件，请点击邮件中的链接完成激活。'"
      >
        <template #extra>
          <el-button
            type="primary"
            @click="router.push('/login')"
          >
            前往登录
          </el-button>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getCaptcha, register } from '@/api/auth'
import { getRegisterConfig } from '@/api/config'
import { useErrorHandler } from '@/composables/useErrorHandler'
import { getPublicSiteConfigs } from '@/api/public'
import { Warning, Operation, HomeFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useLanguageStore } from '@/pinia/modules/language'
import logoUrl from '@/assets/images/logo.png'

const router = useRouter()
const { t, locale } = useI18n()
const { executeAsync, handleSubmit } = useErrorHandler()
const languageStore = useLanguageStore()
const registerFormRef = ref()
const loading = ref(false)
const showInviteCode = ref(false)
const inviteCodeRequired = ref(false)
const captchaImage = ref('')
const registrationEnabled = ref(true)
const emailVerificationRequired = ref(false)
const emailVerificationSent = ref(false)
const siteConfigs = ref({}) // 站点配置

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  captcha: '',
  captchaId: '',
  inviteCode: '',
  registerType: 'normal'
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== registerForm.password) {
    callback(new Error(t('register.passwordMismatch')))
  } else {
    callback()
  }
}

const validateInviteCode = (rule, value, callback) => {
  if (inviteCodeRequired.value && (!value || value.trim() === '')) {
    callback(new Error(t('register.pleaseEnterInviteCode')))
  } else {
    callback()
  }
}

const registerRules = computed(() => ({
  username: [
    { required: true, message: t('register.pleaseEnterUsername'), trigger: 'blur' },
    { min: 3, max: 20, message: t('validation.usernameLength', { min: 3, max: 20 }), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('register.pleaseEnterPassword'), trigger: 'blur' },
    { min: 8, message: t('validation.passwordLength'), trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: t('register.pleaseConfirmPassword'), trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  captcha: [
    { required: true, message: t('register.pleaseEnterCaptcha'), trigger: 'blur' }
  ],
  inviteCode: [
    { validator: validateInviteCode, trigger: 'blur' }
  ]
}))

const refreshCaptcha = async () => {
  await executeAsync(async () => {
    const response = await getCaptcha()
    captchaImage.value = response.data.imageData
    registerForm.captchaId = response.data.captchaId
    registerForm.captcha = ''
  }, {
    showError: false, // 静默处理验证码错误
    showLoading: false
  })
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  // 防止重复提交
  if (loading.value) return

  await registerFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    // 再次检查loading状态，防止表单验证期间的重复点击
    if (loading.value) return

    loading.value = true
    try {
      const result = await handleSubmit(async () => {
        return await register({
          username: registerForm.username,
          password: registerForm.password,
          captcha: registerForm.captcha,
          captchaId: registerForm.captchaId,
          inviteCode: showInviteCode.value ? registerForm.inviteCode : undefined,
          registerType: registerForm.registerType
        })
      }, {
        successMessage: t('register.registerSuccess'),
        showLoading: false // 使用组件自己的loading
      })

      if (result.success && result.data) {
        // 注册成功，直接设置用户登录状态
        const responseData = result.data.data // 正确获取嵌套的data数据

        // 检查是否需要邮箱验证
        if (emailVerificationRequired.value && registerForm.email) {
          emailVerificationSent.value = true
          return
        }

        // 导入用户store
        const { useUserStore } = await import('@/pinia/modules/user')
        const userStore = useUserStore()
        
        // 设置用户登录状态
        userStore.setToken(responseData.token)
        userStore.setUser(responseData.user)
        
        // 跳转到用户仪表盘
        router.push('/user/dashboard')
      } else {
        refreshCaptcha() // 注册失败刷新验证码
      }
    } finally {
      loading.value = false
    }
  })
}

const checkRegistrationEnabled = async () => {
  const result = await executeAsync(async () => {
    const response = await getRegisterConfig()
    const config = response.data
    
    // 新逻辑：如果启用公开注册，或者启用邀请码系统但不强制要求邀请码
    const enablePublicRegistration = config.auth?.enablePublicRegistration ?? false
    const inviteCodeEnabled = config.inviteCode?.enabled ?? false
    
    // 如果启用公开注册，或者启用了邀请码系统，则允许注册
    const canRegister = enablePublicRegistration || inviteCodeEnabled
    
    // 显示邀请码输入框的条件：启用了邀请码系统
    showInviteCode.value = inviteCodeEnabled
    
    // 邀请码必填的条件：启用邀请码系统且未启用公开注册
    inviteCodeRequired.value = inviteCodeEnabled && !enablePublicRegistration

    // 检查是否需要邮箱验证
    emailVerificationRequired.value = config.auth?.enableEmailVerification ?? false

    return canRegister
  }, {
    showError: false, // 不显示错误消息
    showLoading: false
  })
  
  // 如果成功获取配置，使用返回的值；否则默认允许注册
  registrationEnabled.value = result.success ? result.data : true
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

// 切换语言
const switchLanguage = () => {
  const newLang = languageStore.toggleLanguage()
  locale.value = newLang
  ElMessage.success(t('navbar.languageSwitched'))
}

onMounted(async () => {
  await checkRegistrationEnabled()
  await fetchSiteConfigs()
  if (registrationEnabled.value) {
    refreshCaptcha()
  }
})
</script>

<style scoped>
.register-container {
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

.register-form {
  margin: auto;
  margin-top: 40px;
  margin-bottom: 40px;
  width: 500px;
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

.register-form :deep(.el-form-item__label) {
  color: var(--text-color-regular) !important;
  font-weight: 500;
}

.register-form :deep(.el-input__wrapper) {
  border-radius: var(--border-radius-sm);
  background: rgba(255, 255, 255, 0.08) !important;
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.15) inset !important;
}

.register-form :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.25) inset !important;
}

.register-form :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--primary-color-light) inset, 0 0 0 3px rgba(99, 102, 241, 0.2) !important;
}

.register-form :deep(.el-input__inner) {
  color: var(--text-color-primary) !important;
}

.register-form :deep(.el-button--primary) {
  width: 100%;
  height: 45px;
  background: linear-gradient(135deg, #6366F1, #8B5CF6) !important;
  border: none !important;
  border-radius: var(--border-radius-sm) !important;
  font-size: 16px;
  font-weight: 600;
  box-shadow: 0 4px 15px rgba(99, 102, 241, 0.4) !important;
}

.register-form :deep(.el-button--primary:hover) {
  background: linear-gradient(135deg, #4F46E5, #7C3AED) !important;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(99, 102, 241, 0.5) !important;
}

.registration-disabled {
  width: 500px;
  margin: auto;
  margin-top: 60px;
  margin-bottom: 60px;
}

.registration-disabled :deep(.el-card) {
  border-radius: var(--border-radius-xl) !important;
  background: var(--auth-card-bg) !important;
  backdrop-filter: blur(24px);
  border: 1px solid var(--auth-card-border);
}

.disabled-content {
  text-align: center;
  padding: 40px;
}

.disabled-content h2 {
  color: var(--warning-color);
  margin: 20px 0;
  font-size: 24px;
}

.disabled-content p {
  color: var(--text-color-secondary);
  margin-bottom: 30px;
  font-size: 16px;
  line-height: 1.5;
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-header h2 {
  font-size: 24px;
  color: var(--text-color-primary);
  margin-bottom: 10px;
  font-weight: 700;
}

.register-header p {
  font-size: 14px;
  color: var(--text-color-secondary);
}

.form-footer {
  text-align: center;
  margin-top: 20px;
}

.form-footer p {
  color: var(--text-color-secondary);
}

.form-footer a {
  color: var(--primary-color-light);
  text-decoration: none;
  font-weight: 500;
}

.captcha-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.captcha-image {
  width: 38%;
  height: 40px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: var(--border-radius-sm);
  overflow: hidden;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
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

.password-hint {
  margin-top: 5px;
  font-size: 12px;
  line-height: 1.4;
}

@media (max-width: 768px) {
  .register-form {
    width: 90%;
    padding: 24px;
  }
}
</style>