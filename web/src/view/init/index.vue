<template>
  <div class="init-page">
    <div class="init-container">
      <div class="init-bg-pattern" />
      <div class="init-card">
        <!-- Header -->
        <div class="init-header">
          <div class="init-logo">
            <svg
              viewBox="0 0 24 24"
              width="40"
              height="40"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M12 2L2 7l10 5 10-5-10-5z" />
              <path d="M2 17l10 5 10-5" />
              <path d="M2 12l10 5 10-5" />
            </svg>
          </div>
          <h1>{{ t('init.title') }}</h1>
          <p>{{ t('init.subtitle') }}</p>
        </div>

        <!-- Progress Panel (shown after submit) -->
        <InitProgressPanel
          v-if="showProgress"
          :status="progressStatus"
          :steps="progressSteps"
          :percent="progressPercent"
          :bar-status="progressBarStatus"
          :tag-type="progressTagType"
          :status-text="progressStatusText"
          @retry="retryInit"
          @go-home="goHome"
        />

        <!-- Form Panel (hidden after submit) -->
        <template v-if="!showProgress">
          <!-- Steps indicator -->
          <el-steps
            :active="stepIndex"
            finish-status="success"
            align-center
            class="init-steps"
          >
            <el-step :title="t('init.database.tabLabel')" />
            <el-step :title="t('init.admin.tabLabel')" />
            <el-step :title="t('init.user.tabLabel')" />
          </el-steps>

          <!-- Step 1: Database -->
          <div
            v-show="activeTab === 'database'"
            class="step-content"
          >
            <el-form
              ref="databaseFormRef"
              :model="databaseForm"
              :rules="databaseRules"
              label-width="140px"
              label-position="top"
              size="large"
            >
              <el-form-item
                :label="t('init.database.type')"
                prop="type"
              >
                <el-radio-group
                  v-model="databaseForm.type"
                  @change="onDatabaseTypeChange"
                >
                  <el-radio-button label="mysql">
                    MySQL
                  </el-radio-button>
                  <el-radio-button label="mariadb">
                    MariaDB
                  </el-radio-button>
                </el-radio-group>
                <div class="field-hint">
                  <el-text
                    v-if="dbRecommendation"
                    size="small"
                    type="success"
                  >
                    {{ t('init.database.autoDetectHint') }}
                  </el-text>
                  <el-text
                    v-else
                    size="small"
                    type="info"
                  >
                    {{ t('init.database.autoSelectHint') }}
                  </el-text>
                </div>
              </el-form-item>

              <el-row :gutter="16">
                <el-col :span="16">
                  <el-form-item
                    :label="t('init.database.host')"
                    prop="host"
                  >
                    <el-input
                      v-model="databaseForm.host"
                      placeholder="127.0.0.1"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item
                    :label="t('init.database.port')"
                    prop="port"
                  >
                    <el-input
                      v-model="databaseForm.port"
                      placeholder="3306"
                    />
                  </el-form-item>
                </el-col>
              </el-row>

              <el-form-item
                :label="t('init.database.dbName')"
                prop="database"
              >
                <el-input
                  v-model="databaseForm.database"
                  placeholder="oneclickvirt"
                />
              </el-form-item>

              <el-row :gutter="16">
                <el-col :span="12">
                  <el-form-item
                    :label="t('init.database.username')"
                    prop="username"
                  >
                    <el-input
                      v-model="databaseForm.username"
                      placeholder="root"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item
                    :label="t('init.database.password')"
                    prop="password"
                  >
                    <el-input
                      v-model="databaseForm.password"
                      type="password"
                      :placeholder="t('init.database.passwordPlaceholder')"
                      show-password
                    />
                  </el-form-item>
                </el-col>
              </el-row>

              <el-form-item>
                <el-button
                  type="info"
                  plain
                  :loading="testingConnection"
                  @click="testDatabaseConnection"
                >
                  {{ t('init.database.testConnection') }}
                </el-button>
                <span
                  v-if="connectionTestResult"
                  :class="connectionTestResult.success ? 'test-success' : 'test-error'"
                >
                  {{ connectionTestResult.message }}
                </span>
              </el-form-item>
            </el-form>
          </div>

          <!-- Step 2: Admin -->
          <div
            v-show="activeTab === 'admin'"
            class="step-content"
          >
            <el-form
              ref="adminFormRef"
              :model="initForm.admin"
              :rules="adminRules"
              label-position="top"
              size="large"
            >
              <el-form-item
                :label="t('init.admin.username')"
                prop="username"
              >
                <el-input
                  v-model="initForm.admin.username"
                  :placeholder="t('init.admin.usernamePlaceholder')"
                  clearable
                  prefix-icon="User"
                />
              </el-form-item>
              <el-form-item
                :label="t('init.admin.password')"
                prop="password"
              >
                <el-input
                  v-model="initForm.admin.password"
                  type="password"
                  :placeholder="t('init.admin.passwordPlaceholder')"
                  show-password
                  clearable
                  prefix-icon="Lock"
                />
                <div class="field-hint">
                  <el-text
                    size="small"
                    type="info"
                  >
                    {{ t('init.admin.passwordHint') }}
                  </el-text>
                </div>
              </el-form-item>
              <el-form-item
                :label="t('init.admin.confirmPassword')"
                prop="confirmPassword"
              >
                <el-input
                  v-model="initForm.admin.confirmPassword"
                  type="password"
                  :placeholder="t('init.admin.confirmPasswordPlaceholder')"
                  show-password
                  clearable
                  prefix-icon="Lock"
                />
              </el-form-item>
              <el-form-item
                :label="t('init.admin.email')"
                prop="email"
              >
                <el-input
                  v-model="initForm.admin.email"
                  :placeholder="t('init.admin.emailPlaceholder')"
                  clearable
                  prefix-icon="Message"
                />
              </el-form-item>
            </el-form>
          </div>

          <!-- Step 3: User -->
          <div
            v-show="activeTab === 'user'"
            class="step-content"
          >
            <el-form
              ref="userFormRef"
              :model="initForm.user"
              :rules="userRules"
              label-position="top"
              size="large"
            >
              <el-form-item :label="t('init.user.enableStatus')">
                <div class="enable-toggle">
                  <el-switch
                    v-model="initForm.user.enabled"
                    :active-text="t('common.enabled')"
                    :inactive-text="t('common.disabled')"
                  />
                  <el-text
                    size="small"
                    type="warning"
                    class="enable-hint"
                  >
                    {{ t('init.user.enableHint') }}
                  </el-text>
                </div>
              </el-form-item>
              <template v-if="initForm.user.enabled">
                <el-form-item
                  :label="t('init.user.username')"
                  prop="username"
                >
                  <el-input
                    v-model="initForm.user.username"
                    :placeholder="t('init.user.usernamePlaceholder')"
                    clearable
                    prefix-icon="User"
                  />
                </el-form-item>
                <el-form-item
                  :label="t('init.user.password')"
                  prop="password"
                >
                  <el-input
                    v-model="initForm.user.password"
                    type="password"
                    :placeholder="t('init.user.passwordPlaceholder')"
                    show-password
                    clearable
                    prefix-icon="Lock"
                  />
                  <div class="field-hint">
                    <el-text
                      size="small"
                      type="info"
                    >
                      {{ t('init.user.passwordHint') }}
                    </el-text>
                  </div>
                </el-form-item>
                <el-form-item
                  :label="t('init.user.confirmPassword')"
                  prop="confirmPassword"
                >
                  <el-input
                    v-model="initForm.user.confirmPassword"
                    type="password"
                    :placeholder="t('init.user.confirmPasswordPlaceholder')"
                    show-password
                    clearable
                    prefix-icon="Lock"
                  />
                </el-form-item>
                <el-form-item
                  :label="t('init.user.email')"
                  prop="email"
                >
                  <el-input
                    v-model="initForm.user.email"
                    :placeholder="t('init.user.emailPlaceholder')"
                    clearable
                    prefix-icon="Message"
                  />
                </el-form-item>
              </template>
            </el-form>
          </div>

          <!-- Navigation buttons -->
          <div class="init-actions">
            <el-button
              v-if="activeTab !== 'database'"
              size="large"
              @click="prevStep"
            >
              {{ t('common.back') }}
            </el-button>
            <el-button
              type="info"
              plain
              size="large"
              @click="fillDefaultData"
            >
              {{ t('init.fillDefaults') }}
            </el-button>
            <div style="flex: 1" />
            <el-button
              v-if="activeTab !== 'user'"
              type="primary"
              size="large"
              @click="nextStep"
            >
              {{ t('init.nextStep') }}
            </el-button>
            <el-button
              v-else
              type="primary"
              :loading="loading"
              :disabled="loading || !isFormValid"
              size="large"
              @click="handleInit"
            >
              {{ t('init.initSystem') }}
            </el-button>
          </div>
        </template>
      </div>
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import InitProgressPanel from './components/InitProgressPanel.vue'
import useInit from './useInit'
import { useI18n } from 'vue-i18n'
import AppFooter from '@/view/layout/components/AppFooter.vue'
const { t } = useI18n()

const {
  adminFormRef, userFormRef, databaseFormRef,
  loading, testingConnection, connectionTestResult,
  activeTab, dbRecommendation,
  showProgress, progressStatus, progressSteps,
  progressPercent, progressBarStatus, progressTagType, progressStatusText,
  stepIndex,
  databaseForm, initForm,
  adminRules, userRules, databaseRules,
  isFormValid,
  nextStep, prevStep,
  onDatabaseTypeChange, fillDefaultData,
  testDatabaseConnection, handleInit,
  retryInit, goHome
} = useInit()
</script>

<style scoped>
.init-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--auth-page-bg);
}

.init-container {
  flex: 1;
  min-height: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  position: relative;
  overflow-x: hidden;
}

.init-bg-pattern {
  position: absolute;
  inset: 0;
  background-image:
    radial-gradient(circle at 20% 30%, var(--accent-soft-bg) 0%, transparent 50%),
    radial-gradient(circle at 80% 70%, color-mix(in srgb, var(--info-color) 7%, transparent) 0%, transparent 50%);
  z-index: 0;
}

.init-card {
  background: var(--card-bg);
  backdrop-filter: blur(20px);
  padding: 36px 34px 30px;
  border-radius: 20px;
  box-shadow: var(--box-shadow-heavy);
  width: 100%;
  max-width: 620px;
  border: 1px solid var(--border-color);
  position: relative;
  z-index: 1;
}

.init-header {
  text-align: center;
  margin-bottom: 24px;
}

.init-logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-light));
  color: white;
  margin-bottom: 16px;
  box-shadow: 0 4px 12px var(--primary-color-shadow);
}

.init-header h1 {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-color-primary);
  margin: 0 0 8px;
}

.init-header p {
  font-size: 15px;
  color: var(--text-color-secondary);
  margin: 0;
}

.init-steps {
  margin-bottom: 24px;
}

.step-content {
  min-height: 240px;
}

.init-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 18px;
  padding-top: 18px;
  border-top: 1px solid var(--border-color);
}

.field-hint {
  margin-top: 6px;
}

.enable-toggle {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.enable-hint {
  line-height: 1.4;
}

.test-success {
  color: var(--success-color);
  margin-left: 12px;
  font-size: 14px;
}

.test-error {
  color: var(--error-color);
  margin-left: 12px;
  font-size: 14px;
}

:deep(.el-steps) {
  padding: 0 20px;
}

:deep(.el-step__title.is-finish) {
  color: var(--primary-color);
}

:deep(.el-step__head.is-finish) {
  color: var(--primary-color);
  border-color: var(--primary-color);
}

:deep(.el-step__head.is-process) {
  color: var(--primary-color);
  border-color: var(--primary-color);
}

:deep(.el-step__title.is-process) {
  color: var(--primary-color);
  font-weight: 600;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--text-color-primary);
}

:deep(.el-input__wrapper) {
  border-radius: 10px;
  background: var(--el-fill-color-blank);
  box-shadow: 0 0 0 1px var(--el-border-color) inset;
  transition: all 0.2s;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px var(--border-color-hover) inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px var(--primary-color-focus) inset;
}

:deep(.el-button--primary) {
  background: var(--primary-color);
  border-color: var(--primary-color);
  border-radius: 10px;
}

:deep(.el-button--primary:hover) {
  background: var(--primary-color-dark);
  border-color: var(--primary-color-dark);
}

:deep(.el-radio-button__inner) {
  border-radius: 8px;
}

@media (max-width: 640px) {
  .init-container {
    padding: 12px;
    align-items: flex-start;
  }

  .init-card {
    padding: 24px 18px 20px;
    border-radius: 16px;
  }

  .init-header h1 {
    font-size: 24px;
  }

  .init-actions {
    flex-wrap: wrap;
  }

  :deep(.el-steps) {
    padding: 0;
  }

  :deep(.el-form-item__label) {
    line-height: 1.3;
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
