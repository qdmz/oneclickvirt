<template>
  <footer class="app-footer">
    <div class="app-footer-inner">
      <span class="footer-copyright">&copy; 2026 {{ siteStore.displaySiteName }}. {{ t('home.footer.allRightsReserved') }}</span>
      <span class="footer-divider" />
      <a
        href="https://github.com/oneclickvirt"
        target="_blank"
        rel="noopener noreferrer"
        class="footer-link"
      >
        <svg
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="currentColor"
          class="footer-github-icon"
        >
          <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
        </svg>
        {{ t('home.footer.openSourceProject') }}
      </a>
      <span
        v-if="serverVersion"
        class="footer-divider"
      />
      <span
        v-if="serverVersion"
        class="footer-version"
        :title="`${t('home.footer.serverVersion')} ${serverVersion}`"
      >
        <span>{{ t('home.footer.serverVersion') }}</span>
        <span class="footer-version-value">{{ serverVersion }}</span>
      </span>
      <el-button
        v-if="isSuperAdmin && serverVersion"
        link
        size="small"
        class="footer-update-button"
        :title="t('home.footer.manageUpdates')"
        :aria-label="t('home.footer.manageUpdates')"
        @click="updateDialogVisible = true"
      >
        <el-icon><Tools /></el-icon>
        <span>{{ t('home.footer.manageUpdates') }}</span>
      </el-button>
      <a
        v-if="updateAvailable && latestVersion"
        :href="releaseUrl || 'https://github.com/oneclickvirt/oneclickvirt/releases'"
        target="_blank"
        rel="noopener noreferrer"
        class="footer-update-link"
        :title="`${t('home.footer.latestVersion')} ${latestVersion}`"
      >
        <span>{{ t('home.footer.latestVersion') }}</span>
        <span class="footer-version-value">{{ latestVersion }}</span>
      </a>
      <span
        v-if="versionFetchFailed"
        class="footer-divider"
      />
      <span
        v-if="versionFetchFailed"
        class="footer-version-error"
      >
        {{ t('home.footer.versionFetchFailed') }}
      </span>
    </div>
    <SystemUpdateDialog v-model="updateDialogVisible" />
  </footer>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useSiteStore } from '@/pinia/modules/site'
import { useUserStore } from '@/pinia/modules/user'
import { Tools } from '@element-plus/icons-vue'
import { getServerVersion } from '@/api/public'
import SystemUpdateDialog from './SystemUpdateDialog.vue'

const { t } = useI18n()
const siteStore = useSiteStore()
const userStore = useUserStore()
const isSuperAdmin = computed(() => userStore.userType === 'admin' && Boolean(userStore.token))
const updateDialogVisible = ref(false)
const serverVersion = ref('')
const latestVersion = ref('')
const releaseUrl = ref('')
const updateAvailable = ref(false)
const versionFetchFailed = ref(false)

onMounted(async () => {
  try {
    const res = await getServerVersion()
    if (res && (res.code === 200) && res.data?.server_version) {
      serverVersion.value = res.data.server_version
      latestVersion.value = res.data.latest_version || ''
      releaseUrl.value = res.data.release_url || ''
      updateAvailable.value = Boolean(res.data.update_available)
      versionFetchFailed.value = res.data.version_check_status === 'failed'
    } else {
      versionFetchFailed.value = true
    }
  } catch {
    versionFetchFailed.value = true
  }
})
</script>

<style lang="scss" scoped>
.app-footer {
  width: 100%;
  background-color: var(--footer-bg, var(--bg-color-secondary));
  border-top: 1px solid var(--border-color);
  padding: 9px 0;
  margin-top: auto;
  flex-shrink: 0;
}

.app-footer-inner {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  flex-wrap: wrap;
  min-width: 0;
  padding: 0 var(--spacing-lg);
}

@media (max-width: 768px) {
  .app-footer {
    padding: 8px 0 calc(8px + env(safe-area-inset-bottom));
  }

  .app-footer-inner {
    gap: 8px;
    padding: 0 12px;
  }

  .footer-copyright,
  .footer-link {
    font-size: 12px;
  }
}

.footer-copyright {
  font-size: 13px;
  color: var(--text-color-secondary);
}

.footer-divider {
  display: inline-block;
  width: 1px;
  height: 14px;
  background-color: var(--border-color);
  vertical-align: middle;
}

.footer-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  max-width: 100%;
  font-size: 13px;
  color: var(--footer-link-color, var(--primary-color));
  text-decoration: none;
  transition: var(--transition-all);

  &:hover {
    color: var(--footer-link-hover-color, var(--primary-color-dark));
  }
}

.footer-github-icon {
  flex-shrink: 0;
}

.footer-version {
  display: inline-flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  max-width: 100%;
  min-width: 0;
  font-size: 12px;
  color: var(--text-color-tertiary);
  font-family: monospace;
  line-height: 1.5;
  overflow-wrap: anywhere;
  white-space: normal;
}

.footer-update-link {
  display: inline-flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  max-width: 100%;
  min-width: 0;
  font-size: 12px;
  color: var(--el-color-success);
  text-decoration: none;
  font-family: monospace;
  line-height: 1.5;
  overflow-wrap: anywhere;
  white-space: normal;

  &:hover {
    text-decoration: underline;
  }
}

.footer-version-value {
  min-width: 0;
  word-break: break-all;
  overflow-wrap: anywhere;
}

.footer-update-button {
  min-height: 24px;
  padding: 0 4px;
  color: var(--footer-link-color, var(--primary-color));
  font-size: 12px;
}

.footer-update-button :deep(.el-icon) {
  margin-right: 3px;
}

.footer-version-error {
  font-size: 12px;
  color: var(--el-color-warning);
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
