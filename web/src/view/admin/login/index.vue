<template>
  <div class="admin-login-container">
    <!-- 顶部�?-->
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

    <div class="login-form">
      <div class="login-header">
        <h2>{{ t('adminLogin.title') }}</h2>
        <p>{{ t('adminLogin.subtitle') }}</p>
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
            :placeholder="t('login.pleaseEnterAdminUsername')"
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

        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            style="width: 100%;"
            @click="handleLogin"
          >
            {{ t('common.login') }}
          </el-button>
        </el-form-item>

        <div class="form-footer">
          <router-link
            to="/login"
            class="back-link"
          >
            {{ t('login.backToUserLogin') }}
          </router-link>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/pinia/modules/user'
import { ElMessage } from 'element-plus'
import { useErrorHandler } from '@/composables/useErrorHandler'

import { getCaptcha } from '@/api/auth'
import { getPublicSiteConfigs } from '@/api/public'
import { Operation, HomeFilled } from '@element-plus/icons-vue'
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
const siteConfigs = ref({}) // 站点配置

const loginForm = reactive({
  username: '',
  password: '',
  captcha: '',
  userType: 'admin',
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
    
    // 再次检查loading状态，防止表单验证期间的重复点�?
    if (loading.value) return
    
    loading.value = true
    
    try {
      const result = await handleSubmit(async () => {
        return await userStore.adminLogin({
          ...loginForm,
          captchaId: captchaId.value
        })
      }, {
        successMessage: t('login.loginSuccess'),
        showLoading: false // 使用组件自己的loading
      })

      if (result.success) {
        router.push('/admin/dashboard')
      } else {
        refreshCaptcha() // 登录失败刷新验证�?
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
    showError: false, // 静默处理验证码错�?
    showLoading: false
  })
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
  await fetchSiteConfigs()
  refreshCaptcha()
})
</script>

<style scoped>
.admin-login-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: var(--auth-page-bg);
}

/* 顶部栏样�?*/
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

.site-name-link {
  text-decoration: none;
}

.logo-image {
  width: 64px;
  height: 64px;
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

.admin-login-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  background-size: cover;
  opacity: 0.1;
  z-index: -1;
}

.login-form {
  width: 400px;
  padding: 40px;
  background: var(--auth-card-bg);
  backdrop-filter: blur(24px);
  border-radius: var(--border-radius-xl);
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
  margin: auto;
  margin-top: 60px;
  margin-bottom: 60px;
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
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
  font-size: 24px;
  color: #303133;
  margin-bottom: 10px;
}

.login-header p {
  font-size: 14px;
  color: #909399;
}

.form-footer {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: #909399;
  width: 100%;
}

.login-form :deep(.el-button) {
  width: 100% !important;
  height: 45px;
}

.back-link {
  color: #909399;
  text-decoration: none;
  margin: 0 5px;
}

.back-link:hover {
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
  border: 1px solid #dcdfe6;
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
  color: #909399;
}

@media (max-width: 768px) {
  .login-form {
    width: 90%;
    padding: 20px;
  }
}
</style>