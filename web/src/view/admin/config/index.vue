<template>
  <div class="config-container">
    <el-card>
      <template #header>
        <div class="config-header">
          <span>{{ $t('admin.config.title') }}</span>
        </div>
      </template>

      <!-- 配置分类标签页 -->
      <el-tabs
        v-model="activeTab"
        type="border-card"
        class="config-tabs"
      >
        <!-- 基础认证配置 -->
        <el-tab-pane
          :label="$t('admin.config.basicAuth')"
          name="auth"
        >
          <AuthTab
            :config="config"
            :loading="loading"
          />
        </el-tab-pane>

        <!-- 邮箱SMTP配置 -->
        <el-tab-pane
          :label="$t('admin.config.emailConfig')"
          name="email"
        >
          <EmailTab
            :config="config"
            :loading="loading"
          />
        </el-tab-pane>

        <!-- 第三方登录配置 -->
        <el-tab-pane
          :label="$t('admin.config.thirdPartyLogin')"
          name="oauth"
        >
          <OAuthTab
            :config="config"
            :loading="loading"
          />
        </el-tab-pane>

        <!-- 用户等级配置 -->
        <el-tab-pane
          :label="$t('admin.config.userLevel')"
          name="quota"
        >
          <QuotaTab
            :config="config"
            :loading="loading"
          />
        </el-tab-pane>

        <!-- 实例类型权限配置 -->
        <el-tab-pane
          :label="$t('admin.config.instancePermissions')"
          name="instancePermissions"
        >
          <InstancePermissionsTab
            :permissions="instanceTypePermissions"
            :loading="loading"
          />
        </el-tab-pane>

        <!-- 实名认证配置 -->
        <el-tab-pane
          :label="$t('admin.config.kycConfig')"
          name="kyc"
        >
          <KycTab
            :config="config"
            :loading="loading"
          />
        </el-tab-pane>

        <!-- 其他配置 -->
        <el-tab-pane
          :label="$t('admin.config.otherConfig')"
          name="other"
        >
          <OtherTab
            :config="config"
            :loading="loading"
          />
        </el-tab-pane>
      </el-tabs>

      <!-- 底部操作按钮 -->
      <div class="config-actions">
        <el-button
          type="primary"
          size="large"
          :loading="loading"
          @click="saveConfig"
        >
          {{ $t('admin.config.saveCurrentConfig') }}
        </el-button>
        <el-button 
          size="large"
          @click="resetConfig"
        >
          {{ $t('admin.config.resetConfig') }}
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { useConfigManagement } from './composables/useConfigManagement'
import AuthTab from './components/AuthTab.vue'
import EmailTab from './components/EmailTab.vue'
import OAuthTab from './components/OAuthTab.vue'
import QuotaTab from './components/QuotaTab.vue'
import InstancePermissionsTab from './components/InstancePermissionsTab.vue'
import KycTab from './components/KycTab.vue'
import OtherTab from './components/OtherTab.vue'

const {
  activeTab, config, instanceTypePermissions, loading,
  saveConfig, resetConfig, t
} = useConfigManagement()
</script>

<style scoped>
.config-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
  
  > span {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-color-primary);
  }
}

.config-tabs {
  margin-bottom: 20px;
}

.config-tabs :deep(.el-tabs__content) {
  padding: 20px;
}

.config-form {
  max-height: 600px;
  overflow-y: auto;
}

.oauth-card {
  margin-bottom: 16px;
}

.oauth-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.level-card {
  border: 2px solid #f0f0f0;
  transition: all 0.3s ease;
}

.level-card:hover {
  border-color: #16a34a;
  box-shadow: 0 2px 12px 0 rgba(22, 163, 74, 0.15);
}

.level-card.default-level {
  border-color: #16a34a;
  background-color: var(--success-bg);
}

.level-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.level-title {
  font-weight: 600;
  color: var(--text-color-primary);
}

.config-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 20px 0;
  border-top: 1px solid #f0f0f0;
  margin-top: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .config-container {
    padding: 10px;
  }
  
  .config-form {
    max-height: none;
  }
  
  .level-card :deep(.el-col) {
    margin-bottom: 10px;
  }
  
  .config-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .config-actions .el-button {
    width: 100%;
    max-width: 200px;
  }
}

/* 标签页样式 */
.config-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
}

.config-tabs :deep(.el-tabs__nav-wrap) {
  padding: 0 10px;
}

.config-tabs :deep(.el-tabs__item) {
  padding: 0 20px;
  font-weight: 500;
}

/* 表单样式 */
.config-form :deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--text-color-secondary);
}

.config-form :deep(.el-alert) {
  margin-bottom: 20px;
}

.form-item-hint {
  font-size: 12px;
  color: var(--text-color-tertiary);
  margin-top: 4px;
  line-height: 1.4;
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