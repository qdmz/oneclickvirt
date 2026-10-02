<template>
  <div class="navbar">
    <!-- 移动端汉堡菜单按钮 -->
    <div class="hamburger-container">
      <el-button
        class="hamburger-btn"
        :icon="Menu"
        :title="t('navbar.toggleSidebar')"
        :aria-label="t('navbar.toggleSidebar')"
        circle
        @click="toggleSidebar"
      />
    </div>
    
    <div class="right-menu">
      <!-- 主题切换按鈕 -->
      <div class="theme-switcher">
        <el-button
          :title="themeStore.isDark ? t('navbar.lightMode') : t('navbar.darkMode')"
          @click="toggleTheme"
        >
          <el-icon><component :is="themeStore.isDark ? Sunny : Moon" /></el-icon>
        </el-button>
      </div>
      <!-- 语言切换按钮 -->
      <div class="language-switcher">
        <el-button
          :title="t('navbar.switchLanguage')"
          @click="switchLanguage"
        >
          <el-icon><Operation /></el-icon>
          <span class="language-text">{{ languageStore.currentLanguage === 'zh-CN' ? 'English' : '中文' }}</span>
        </el-button>
      </div>

      <el-dropdown
        class="avatar-container"
        trigger="click"
      >
        <div class="avatar-wrapper">
          <el-avatar
            :size="40"
            :src="userInfo.headerImg || ''"
          >
            <el-icon><User /></el-icon>
          </el-avatar>
          <span class="username">{{ userInfo.nickname || userInfo.username }}</span>
          <el-icon class="el-icon-caret-bottom">
            <CaretBottom />
          </el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <!-- 管理员视图切换按钮 -->
            <el-dropdown-item
              v-if="userStore.canSwitchViewMode"
              @click="toggleViewMode"
            >
              <el-icon style="margin-right: 8px;">
                <Switch />
              </el-icon>
              <span>{{ t('navbar.switchTo') }}{{ userStore.currentViewMode === 'admin' ? t('navbar.userView') : t('navbar.adminView') }}</span>
            </el-dropdown-item>
            <el-dropdown-item
              divided
              @click="logout"
            >
              <el-icon style="margin-right: 8px;">
                <SwitchButton />
              </el-icon>
              <span>{{ t('common.logout') }}</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Switch, SwitchButton, User, CaretBottom, Menu, Operation, Sunny, Moon } from '@element-plus/icons-vue'
import { useUserStore } from '@/pinia/modules/user'
import { useLanguageStore } from '@/pinia/modules/language'
import { useThemeStore } from '@/pinia/modules/theme'

const emit = defineEmits(['toggle-sidebar'])
const router = useRouter()
const userStore = useUserStore()
const languageStore = useLanguageStore()
const { t, locale } = useI18n()

const themeStore = useThemeStore()
const userInfo = computed(() => userStore.user || {})

const toggleTheme = () => {
  themeStore.toggleTheme()
}

const toggleSidebar = () => {
  emit('toggle-sidebar')
}

const switchLanguage = () => {
  const newLang = languageStore.toggleLanguage()
  locale.value = newLang
  ElMessage.success(t('navbar.languageSwitched'))
}

const toggleViewMode = () => {
  if (!userStore.canSwitchViewMode) {
    ElMessage.warning(t('navbar.onlyAdminCanSwitch'))
    return
  }
  
  const newMode = userStore.currentViewMode === 'admin' ? 'user' : 'admin'
  const success = userStore.switchViewMode(newMode)
  
  if (success) {
    const viewName = newMode === 'admin' ? t('navbar.adminView') : t('navbar.userView')
    ElMessage.success(`${t('navbar.switchedTo')}${viewName}`)
    
    const targetPath = newMode === 'admin' ? '/admin/dashboard' : '/user/dashboard'
    router.push(targetPath)
  }
}

const logout = async () => {
  try {
    await ElMessageBox.confirm(t('navbar.confirmLogout'), t('navbar.tip'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    
    userStore.logout()
    router.push('/home')
  } catch (error) {
    void error
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: var(--navbar-height);
  overflow: hidden;
  position: relative;
  background: var(--auth-header-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  box-shadow: 0 2px 20px rgba(22, 163, 74, 0.08);
  border-bottom: 1px solid rgba(22, 163, 74, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;

  .hamburger-container {
    display: none;
    
    .hamburger-btn {
      color: var(--text-color-primary);
      background: transparent;
      border: none;
      
      &:hover {
        background: var(--bg-color-hover);
      }
    }
  }

  .right-menu {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-left: auto;

    &:focus {
      outline: none;
    }

    .theme-switcher {
      display: flex;
      align-items: center;

      .el-button {
        color: var(--text-color-primary);
        background: transparent;
        border: 1px solid var(--border-color);
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 8px 12px;
        border-radius: 25px;
        min-width: 40px;
        font-size: 16px;

        &:hover {
          background: rgba(22, 163, 74, 0.08);
          border-color: rgba(22, 163, 74, 0.4);
          color: var(--primary-color);
        }
      }
    }

    .language-switcher {
      display: flex;
      align-items: center;

      .el-button {
        color: var(--text-color-primary);
        background: transparent;
        border: 1px solid var(--border-color);
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 8px 14px;
        border-radius: 25px;
        font-size: 16px;
        
        &:hover {
          background: rgba(22, 163, 74, 0.08);
          border-color: rgba(22, 163, 74, 0.4);
          color: var(--primary-color);
        }

        .language-text {
          font-size: 16px;
          font-weight: 500;
        }
      }
    }

    .right-menu-item {
      display: inline-block;
      padding: 0 8px;
      height: 100%;
      font-size: 18px;
      color: var(--text-color-secondary);
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background .3s;

        &:hover {
          background: rgba(0, 0, 0, .025)
        }
      }
    }

    .avatar-container {
      .avatar-wrapper {
        position: relative;
        display: flex;
        align-items: center;
        cursor: pointer;
        padding: 6px 10px;
        border-radius: 24px;
        transition: all 0.25s ease;

        &:hover {
          background: rgba(22, 163, 74, 0.08);
        }

        .username {
          margin-left: 10px;
          margin-right: 5px;
          font-size: 15px;
          font-weight: 500;
        }

        .el-icon-caret-bottom {
          cursor: pointer;
          font-size: 12px;
          margin-left: 4px;
        }
      }
    }
  }
}

/* 平板和移动端适配 */
@media (max-width: 1024px), (max-width: 1100px) and (orientation: portrait) {
  .navbar {
    .hamburger-container {
      display: block;
    }
    
    .right-menu {
      .avatar-container .avatar-wrapper .username {
        display: none;
      }
    }
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .navbar {
    padding: 0 12px;
    height: var(--navbar-height);
    
    .right-menu {
      gap: 6px;

      .theme-switcher .el-button,
      .language-switcher .el-button {
        min-width: 36px;
        padding: 7px 9px;
        font-size: 14px;
      }

      .language-switcher .el-button .language-text {
        font-size: 14px;
      }

      .avatar-container {
        .avatar-wrapper {
          padding: 4px 6px;

          .el-avatar {
            width: 32px !important;
            height: 32px !important;
          }
          
          .el-icon-caret-bottom {
            display: none;
          }
        }
      }
    }
  }
}

@media (max-width: 480px) {
  .navbar {
    padding: 0 8px;

    .right-menu {
      gap: 4px;

      .language-switcher .el-button .language-text {
        display: none;
      }

      .avatar-container .avatar-wrapper {
        padding: 4px;
      }
    }
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
