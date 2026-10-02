<template>
  <div class="profile-container">
    <!-- 加载状态 -->
    <div
      v-if="loading"
      class="loading-container"
    >
      <el-loading-directive />
      <div class="loading-text">
        {{ t('user.profile.loadingProfile') }}
      </div>
    </div>
    
    <!-- 主要内容 -->
    <div v-else>
      <el-card class="profile-card">
        <template #header>
          <div class="card-header">
            <span>{{ t('user.profile.title') }}</span>
          </div>
        </template>

        <!-- 用户头像和基本信息 -->
        <div class="profile-header">
          <div class="avatar-section">
            <el-avatar
              :size="100"
              :src="userStore.getUserAvatar()"
            />
          </div>
          <div class="user-info">
            <h2>{{ userStore.getUserDisplayName() }}</h2>
            <p class="username">
              @{{ userStore.user?.username }}
            </p>
            <el-tag :type="getUserTypeTagType()">
              >
              {{ getUserTypeText() }}
            </el-tag>
          </div>
        </div>

        <!-- 标签页内容 -->
        <div class="profile-content">
          <el-divider />
        
          <el-tabs
            v-model="activeTab"
            type="card"
            class="profile-tabs"
          >
            <!-- 基本信息标签页 -->
            <el-tab-pane
              :label="t('user.profile.basicInfo')"
              name="basic"
            >
              <el-form
                ref="profileFormRef"
                :model="profileForm"
                :rules="profileRules"
                label-width="100px"
                size="large"
              >
                <el-form-item :label="t('user.profile.username')">
                  <el-input
                    v-model="profileForm.username"
                    disabled
                  />
                  <div class="form-tip">
                    {{ t('user.profile.usernameCannotChange') }}
                  </div>
                </el-form-item>

                <el-form-item
                  :label="t('user.profile.nickname')"
                  prop="nickname"
                >
                  <el-input
                    v-model="profileForm.nickname"
                    :placeholder="t('user.profile.pleaseEnterNickname')"
                    clearable
                  />
                </el-form-item>

                <el-form-item
                  :label="t('user.profile.email')"
                  prop="email"
                >
                  <el-input
                    v-model="profileForm.email"
                    :placeholder="t('user.profile.pleaseEnterEmail')"
                    clearable
                  />
                </el-form-item>

                <el-form-item
                  :label="t('user.profile.phone')"
                  prop="phone"
                >
                  <el-input
                    v-model="profileForm.phone"
                    :placeholder="t('user.profile.pleaseEnterPhone')"
                    clearable
                  />
                </el-form-item>

                <el-form-item>
                  <el-button
                    type="primary"
                    :loading="updating"
                    @click="updateProfile"
                  >
                    {{ t('user.profile.saveChanges') }}
                  </el-button>
                  <el-button @click="resetForm">
                    {{ t('common.reset') }}
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <!-- 密码管理标签页 -->
            <el-tab-pane
              :label="t('user.profile.passwordManagement')"
              name="password"
            >
              <div class="password-section">
                <!-- 修改密码 -->
                <div class="password-change-section" style="margin-bottom: 30px;">
                  <h3>修改密码</h3>
                  <el-form label-width="100px" style="max-width: 500px;">
                    <el-form-item label="旧密码">
                      <el-input
                        v-model="changePasswordForm.oldPassword"
                        type="password"
                        show-password
                        placeholder="请输入当前密码"
                      />
                    </el-form-item>
                    <el-form-item label="新密码">
                      <el-input
                        v-model="changePasswordForm.newPassword"
                        type="password"
                        show-password
                        placeholder="请输入新密码（至少6位）"
                      />
                    </el-form-item>
                    <el-form-item label="确认密码">
                      <el-input
                        v-model="changePasswordForm.confirmPassword"
                        type="password"
                        show-password
                        placeholder="请再次输入新密码"
                      />
                    </el-form-item>
                    <el-form-item>
                      <el-button
                        type="primary"
                        :loading="changingPassword"
                        @click="handleChangePassword"
                      >
                        确认修改
                      </el-button>
                    </el-form-item>
                  </el-form>
                </div>
                <!-- 自动重置密码 -->
                <div class="password-reset-section">
                  <h3>{{ t('user.profile.autoResetPassword') }}</h3>
                  <div class="reset-intro">
                    <el-alert
                      :title="t('user.profile.passwordAutoReset')"
                      type="warning"
                      :closable="false"
                      show-icon
                    >
                      <template #default>
                        <p>{{ t('user.profile.autoResetDescription1') }}</p>
                        <p>{{ t('user.profile.autoResetDescription2') }}</p>
                        <p><strong>{{ t('user.profile.autoResetDescription3') }}</strong></p>
                      </template>
                    </el-alert>
                  
                    <!-- 显示生成的新密码 -->
                    <div
                      v-if="generatedPassword"
                      class="generated-password"
                    >
                      <el-result
                        icon="success"
                        :title="t('user.profile.passwordResetSuccess')"
                        :sub-title="t('user.profile.newPasswordGenerated')"
                      >
                        <template #extra>
                          <div style="margin: 20px 0;">
                            <el-text
                              type="info"
                              style="display: block; margin-bottom: 10px;"
                            >
                              {{ t('user.profile.newPassword') }}：
                            </el-text>
                            <el-input
                              v-model="generatedPassword"
                              readonly
                              style="width: 350px; font-family: monospace; font-size: 16px;"
                            >
                              <template #append>
                                <el-button @click="copyPassword">
                                  {{ t('common.copy') }}
                                </el-button>
                              </template>
                            </el-input>
                          </div>
                          <div style="margin: 20px 0;">
                            <el-text
                              size="small"
                              type="warning"
                            >
                              {{ t('user.profile.passwordSentToChannel') }}
                            </el-text>
                          </div>
                          <div style="margin-top: 20px;">
                            <el-button @click="closePasswordDialog">
                              {{ t('common.close') }}
                            </el-button>
                          </div>
                        </template>
                      </el-result>
                    </div>
                  
                    <!-- 重置密码按钮 -->
                    <div
                      v-else
                      style="margin-top: 20px;"
                    >
                      <el-button
                        type="danger"
                        :loading="resetPasswordLoading"
                        @click="confirmPasswordReset"
                      >
                        {{ t('user.profile.resetPassword') }}
                      </el-button>
                    </div>
                  </div>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </el-card>
    </div> <!-- 结束主要内容区域 -->
  </div>
</template>

<script setup>
import { h, ref, reactive, onMounted, onActivated, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { copyToClipboard as copyToClipboardUtil } from '@/utils/clipboard'
import { useUserStore } from '@/pinia/modules/user'
import { updateProfile as updateProfileApi, resetPassword, changePassword } from '@/api/user'

const { t } = useI18n()
const userStore = useUserStore()
const router = useRouter()

// 当前活动标签页
const activeTab = ref('basic')

// 表单引用
const profileFormRef = ref()

// 加载状态
const loading = ref(true)
const updating = ref(false)
const resetPasswordLoading = ref(false)

// 密码重置相关
const generatedPassword = ref('')

// 修改密码表单
const changePasswordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const changingPassword = ref(false)

// 个人信息表单
const profileForm = reactive({
  username: '',
  nickname: '',
  email: '',
  phone: ''
})

const profileRules = reactive({
  nickname: [
    { required: true, message: () => t('user.profile.pleaseEnterNickname'), trigger: 'blur' },
    { min: 2, max: 20, message: () => t('user.profile.nicknameLengthRange'), trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: () => t('user.profile.invalidEmailFormat'), trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: () => t('user.profile.invalidPhoneFormat'), trigger: 'blur' }
  ]
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== changePasswordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const getUserTypeTagType = () => {
  switch (userStore.userType) {
    case 'admin':
      return 'danger'
    default:
      return 'primary'
  }
}

const getUserTypeText = () => {
  switch (userStore.userType) {
    case 'admin':
      return t('common.admin')
    case 'user':
      return t('common.normalUser')
    default:
      return t('common.unknown')
  }
}

const initForm = () => {
  if (userStore.user) {
    profileForm.username = userStore.user.username
    profileForm.nickname = userStore.user.nickname || ''
    profileForm.email = userStore.user.email || ''
    profileForm.phone = userStore.user.phone || ''
  }
}

const updateProfile = async () => {
  if (!profileFormRef.value) return
  
  await profileFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    updating.value = true
    try {
      const response = await updateProfileApi(profileForm)
      if (response.code === 200) {
        ElMessage.success(t('user.profile.updateSuccess'))
        await userStore.fetchUserInfo()
      } else {
        ElMessage.error(response.msg || t('user.profile.updateFailed'))
      }
    } catch (error) {
      ElMessage.error(t('user.profile.updateFailedRetry'))
    } finally {
      updating.value = false
    }
  })
}

const resetForm = () => {
  initForm()
}

const handleChangePassword = async () => {
  if (!changePasswordForm.oldPassword || !changePasswordForm.newPassword) {
    ElMessage.warning('请填写旧密码和新密码')
    return
  }
  if (changePasswordForm.newPassword !== changePasswordForm.confirmPassword) {
    ElMessage.warning('两次输入的密码不一致')
    return
  }
  if (changePasswordForm.newPassword.length < 6) {
    ElMessage.warning('新密码长度不能少于6位')
    return
  }

  changingPassword.value = true
  try {
    const response = await changePassword({
      oldPassword: changePasswordForm.oldPassword,
      newPassword: changePasswordForm.newPassword
    })
    if (response.code === 200) {
      ElMessage.success('密码修改成功')
      changePasswordForm.oldPassword = ''
      changePasswordForm.newPassword = ''
      changePasswordForm.confirmPassword = ''
    } else {
      ElMessage.error(response.msg || '密码修改失败')
    }
  } catch (error) {
    ElMessage.error(error?.message || '密码修改失败')
  } finally {
    changingPassword.value = false
  }
}

// 确认密码重置
const confirmPasswordReset = async () => {
  try {
    await ElMessageBox.confirm(
      t('user.profile.passwordResetConfirm'),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      }
    )
    
    await resetUserPassword()
  } catch {
    // 用户取消操作
  }
}

// 重置密码
const resetUserPassword = async () => {
  resetPasswordLoading.value = true
  let response
  try {
    response = await resetPassword()
  } catch {
    ElMessage.error(t('user.profile.passwordResetFailedRetry'))
    return
  } finally {
    resetPasswordLoading.value = false
  }
  if (response.code !== 200) {
    ElMessage.error(response.msg || t('user.profile.passwordResetFailed'))
    return
  }

  const loginPath = userStore.isAnyAdmin ? '/admin/login' : '/login'
  // The server has revoked this JWT. Do not leave the UI using it or call
  // logout with it. Keep the returned password in a global dialog so a
  // background 401/navigation cannot erase it before the user saves it.
  userStore.clearUserData()
  generatedPassword.value = typeof response.data?.newPassword === 'string' ? response.data.newPassword : ''
  const message = [h('p', t('user.profile.passwordResetRequiresLogin'))]
  if (generatedPassword.value) {
    message.push(h('p', t('user.profile.newPassword')))
    message.push(h('input', {
      value: generatedPassword.value, readonly: true,
      'aria-label': t('user.profile.newPassword'),
      style: 'width: 100%; box-sizing: border-box; font-family: monospace;',
      onFocus: event => event.target.select()
    }))
    message.push(h('button', { type: 'button', onClick: copyPassword }, t('common.copy')))
  } else {
    message.push(h('p', response.msg || t('user.profile.passwordResetSuccessDefault')))
  }
  await ElMessageBox.alert(h('div', message), t('user.profile.passwordResetSuccess'), {
    confirmButtonText: t('user.profile.passwordSavedLogin'),
    showClose: false, closeOnClickModal: false, closeOnPressEscape: false
  }).catch(() => {})
  generatedPassword.value = ''
  await router.replace(loginPath)
}

// 复制密码到剪贴板
const copyPassword = async () => {
  await copyToClipboardUtil(generatedPassword.value, t('user.profile.passwordCopied'))
}

// 关闭密码对话框
const closePasswordDialog = () => {
  generatedPassword.value = ''
}

onMounted(() => {
  // 强制页面刷新监听器
  window.addEventListener('force-page-refresh', handleForceRefresh)
  
  loading.value = true
  try {
    initForm()
  } finally {
    loading.value = false
  }
})

// 使用 onActivated 确保每次页面激活时都重新加载数据
onActivated(() => {
  loading.value = true
  try {
    initForm()
  } finally {
    loading.value = false
  }
})

// 处理强制刷新事件
const handleForceRefresh = (event) => {
  if (event.detail && event.detail.path === '/user/profile') {
    loading.value = true
    try {
      initForm()
    } finally {
      loading.value = false
    }
  }
}

onUnmounted(() => {
  // 清理事件监听器
  window.removeEventListener('force-page-refresh', handleForceRefresh)
})
</script>

<style scoped>
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  color: #666;
}

.loading-text {
  margin-top: 16px;
  font-size: 14px;
}

.profile-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.profile-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.profile-content {
  padding: 20px;
}

.profile-header {
  display: flex;
  align-items: center;
  margin-bottom: 30px;
}

.avatar-section {
  text-align: center;
  margin-right: 30px;
}

.avatar-section .el-button {
  margin-top: 10px;
}

.user-info h2 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 24px;
}

.username {
  margin: 0 0 10px 0;
  color: #666;
  font-size: 14px;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.password-section h3 {
  margin: 0 0 20px 0;
  color: #333;
  font-size: 16px;
  font-weight: 600;
}

.avatar-uploader {
  text-align: center;
}

.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: 0.2s;
}

.avatar-uploader .el-upload:hover {
  border-color: #16a34a;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  line-height: 178px;
  text-align: center;
}

.avatar-preview {
  width: 178px;
  height: 178px;
  display: block;
}

.password-hint {
  margin-top: 5px;
  font-size: 12px;
  line-height: 1.4;
}

.password-reset-section {
  margin-bottom: 30px;
}

.reset-intro {
  margin-top: 15px;
}

.generated-password {
  margin-top: 20px;
  text-align: center;
}

.profile-tabs {
  margin-top: 20px;
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
