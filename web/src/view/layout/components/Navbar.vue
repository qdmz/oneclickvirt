<template>
  <div class="navbar">
    <div class="navbar-left">
      <!-- 移动端：打开侧边栏抽屉 -->
      <button
        v-if="isMobile"
        class="nav-icon-btn hamburger-btn"
        :title="t('navbar.toggleSidebar')"
        :aria-label="t('navbar.toggleSidebar')"
        @click="toggleSidebar"
      >
        <el-icon><Menu /></el-icon>
      </button>

      <!-- 面包屑 -->
      <nav
        class="breadcrumb"
        aria-label="breadcrumb"
      >
        <router-link
          :to="consoleHome"
          class="crumb crumb-root"
        >
          <el-icon><component :is="isAdminView ? Monitor : HomeFilled" /></el-icon>
          <span class="crumb-text">{{ isAdminView ? t('ui.layout.adminConsole') : t('ui.layout.console') }}</span>
        </router-link>
        <template
          v-for="(c, i) in crumbs"
          :key="c.path"
        >
          <el-icon class="crumb-sep">
            <ArrowRight />
          </el-icon>
          <router-link
            v-if="i < crumbs.length - 1"
            :to="c.path"
            class="crumb"
          >
            {{ c.title }}
          </router-link>
          <span
            v-else
            class="crumb crumb-current"
            aria-current="page"
          >{{ c.title }}</span>
        </template>
      </nav>
    </div>

    <div class="right-menu">
      <button
        class="nav-icon-btn"
        :title="themeStore.isDark ? t('navbar.lightMode') : t('navbar.darkMode')"
        :aria-label="t('navbar.switchTheme')"
        @click="toggleTheme"
      >
        <el-icon><component :is="themeStore.isDark ? Sunny : Moon" /></el-icon>
      </button>

      <button
        class="nav-icon-btn nav-text-btn"
        :title="t('navbar.switchLanguage')"
        @click="switchLanguage"
      >
        <el-icon><Operation /></el-icon>
        <span class="language-text">{{ languageStore.currentLanguage === 'zh-CN' ? 'EN' : '中文' }}</span>
      </button>

      <button
        v-if="fullscreenSupported"
        class="nav-icon-btn hide-mobile"
        :title="isFullscreen ? t('ui.layout.exitFullscreen') : t('ui.layout.fullscreen')"
        :aria-label="t('ui.layout.fullscreen')"
        @click="toggleFullscreen"
      >
        <el-icon><component :is="isFullscreen ? Aim : FullScreen" /></el-icon>
      </button>

      <span class="nav-divider hide-mobile" />

      <el-dropdown
        class="avatar-container"
        trigger="click"
        placement="bottom-end"
        popper-class="oh-user-menu"
      >
        <button class="avatar-wrapper">
          <el-avatar
            :size="32"
            :src="userInfo.headerImg || ''"
            class="avatar"
          >
            {{ avatarInitial }}
          </el-avatar>
          <span class="user-meta">
            <span class="username">{{ displayName }}</span>
            <span class="user-role">{{ isAnyAdmin ? t('ui.layout.adminBadge') : t('ui.layout.userBadge') }}</span>
          </span>
          <el-icon class="caret">
            <ArrowDown />
          </el-icon>
        </button>
        <template #dropdown>
          <div class="user-menu-head">
            <div class="user-menu-label">
              {{ t('ui.layout.signedInAs') }}
            </div>
            <div class="user-menu-name">
              {{ displayName }}
            </div>
            <div
              v-if="userInfo.email"
              class="user-menu-email"
            >
              {{ userInfo.email }}
            </div>
          </div>
          <el-dropdown-menu>
            <el-dropdown-item
              v-if="!isAdminView"
              @click="router.push('/user/profile')"
            >
              <el-icon><User /></el-icon>
              <span>{{ t('ui.layout.profile') }}</span>
            </el-dropdown-item>
            <el-dropdown-item
              v-if="userStore.canSwitchViewMode"
              @click="toggleViewMode"
            >
              <el-icon><Switch /></el-icon>
              <span>{{ t('navbar.switchTo') }}{{ userStore.currentViewMode === 'admin' ? t('navbar.userView') : t('navbar.adminView') }}</span>
            </el-dropdown-item>
            <el-dropdown-item @click="router.push('/home')">
              <el-icon><HomeFilled /></el-icon>
              <span>{{ t('ui.layout.backToSite') }}</span>
            </el-dropdown-item>
            <el-dropdown-item
              divided
              class="logout-item"
              @click="logout"
            >
              <el-icon><SwitchButton /></el-icon>
              <span>{{ t('common.logout') }}</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessageBox, ElMessage } from 'element-plus'
import screenfull from 'screenfull'
import {
  Switch, SwitchButton, User, Menu, Operation, Sunny, Moon,
  ArrowRight, ArrowDown, HomeFilled, Monitor, FullScreen, Aim
} from '@element-plus/icons-vue'
import { useUserStore } from '@/pinia/modules/user'
import { useLanguageStore } from '@/pinia/modules/language'
import { useThemeStore } from '@/pinia/modules/theme'

const emit = defineEmits(['toggle-sidebar'])
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const languageStore = useLanguageStore()
const themeStore = useThemeStore()
const { t, locale } = useI18n()

const isMobile = inject('isMobile', ref(false))
const userInfo = computed(() => userStore.user || {})
const displayName = computed(() => userInfo.value.nickname || userInfo.value.username || '')
const avatarInitial = computed(() => (displayName.value || 'U').charAt(0).toUpperCase())
const isAnyAdmin = computed(() => userStore.isAnyAdmin)
const isAdminView = computed(() => route.path.startsWith('/admin'))
const consoleHome = computed(() => (isAdminView.value ? '/admin/dashboard' : '/user/dashboard'))

// 面包屑：根据当前路由 matched 记录与 meta.title 生成
const crumbs = computed(() => {
  const list = []
  const matched = route.matched.filter(r => r.meta && r.meta.title)
  const last = matched[matched.length - 1]
  if (!last) return list
  const isDashboard = /\/dashboard$/.test(route.path)
  if (isDashboard) {
    list.push({ path: route.path, title: t(last.meta.title) })
    return list
  }
  // 详情页：附带上级列表页（如 /user/instances/:id → 我的实例）
  const segs = route.path.split('/').filter(Boolean)
  if (segs.length >= 3) {
    const parentPath = '/' + segs.slice(0, 2).join('/')
    const parent = router.getRoutes().find(r => r.path === parentPath && r.meta?.title)
    if (parent) list.push({ path: parentPath, title: t(parent.meta.title) })
  }
  list.push({ path: route.path, title: t(last.meta.title) })
  return list
})

const toggleTheme = () => themeStore.toggleTheme()
const toggleSidebar = () => emit('toggle-sidebar')

const switchLanguage = () => {
  const newLang = languageStore.toggleLanguage()
  locale.value = newLang
  ElMessage.success(t('navbar.languageSwitched'))
}

// 全屏
const fullscreenSupported = screenfull.isEnabled
const isFullscreen = ref(false)
const onFullscreenChange = () => { isFullscreen.value = screenfull.isFullscreen }
const toggleFullscreen = () => {
  if (screenfull.isEnabled) screenfull.toggle()
}
onMounted(() => { if (screenfull.isEnabled) screenfull.on('change', onFullscreenChange) })
onBeforeUnmount(() => { if (screenfull.isEnabled) screenfull.off('change', onFullscreenChange) })

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
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 0 20px;
  background: var(--surface-header);
  backdrop-filter: saturate(180%) blur(14px);
  -webkit-backdrop-filter: saturate(180%) blur(14px);
  border-bottom: 1px solid var(--border-color-light);
}

.navbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  flex: 1;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
  font-size: 13.5px;
  white-space: nowrap;
  overflow: hidden;
}

.crumb {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--text-color-secondary);
  text-decoration: none;
  padding: 4px 6px;
  border-radius: 6px;
  transition: color 0.15s, background 0.15s;
  overflow: hidden;
  text-overflow: ellipsis;
}

a.crumb:hover {
  color: var(--text-color-primary);
  background: var(--bg-color-hover);
}

.crumb-root .el-icon {
  font-size: 15px;
  color: var(--primary-color);
}

.crumb-current {
  color: var(--text-color-primary);
  font-weight: 600;
}

.crumb-sep {
  flex-shrink: 0;
  font-size: 12px;
  color: var(--text-color-tertiary);
}

.right-menu {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.nav-icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 34px;
  min-width: 34px;
  padding: 0 8px;
  border: 1px solid transparent;
  border-radius: 9px;
  background: transparent;
  color: var(--text-color-secondary);
  font-size: 16px;
  cursor: pointer;
  transition: color 0.15s, background 0.15s, border-color 0.15s;

  &:hover {
    color: var(--text-color-primary);
    background: var(--bg-color-hover);
    border-color: var(--border-color-light);
  }

  &:focus-visible {
    outline: 2px solid var(--primary-color-focus);
    outline-offset: 1px;
  }
}

.nav-text-btn .language-text {
  font-size: 12.5px;
  font-weight: 600;
}

.nav-divider {
  width: 1px;
  height: 22px;
  margin: 0 6px;
  background: var(--border-color);
}

.avatar-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 8px 4px 4px;
  border: 1px solid transparent;
  border-radius: 12px;
  background: transparent;
  cursor: pointer;
  color: inherit;
  transition: background 0.15s, border-color 0.15s;

  &:hover {
    background: var(--bg-color-hover);
    border-color: var(--border-color-light);
  }

  .avatar {
    flex-shrink: 0;
    background: var(--brand-gradient);
    color: #fff;
    font-weight: 600;
    font-size: 13px;
  }

  .user-meta {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    line-height: 1.2;
    max-width: 140px;
  }

  .username {
    font-size: 13px;
    font-weight: 600;
    color: var(--text-color-primary);
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .user-role {
    font-size: 11px;
    color: var(--text-color-tertiary);
  }

  .caret {
    font-size: 12px;
    color: var(--text-color-tertiary);
  }
}

.user-menu-head {
  padding: 12px 16px 10px;
  border-bottom: 1px solid var(--border-color-light);
  min-width: 200px;
}

.user-menu-label {
  font-size: 11px;
  color: var(--text-color-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.user-menu-name {
  margin-top: 2px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color-primary);
}

.user-menu-email {
  font-size: 12px;
  color: var(--text-color-secondary);
}

:deep(.el-dropdown-menu__item .el-icon) {
  margin-right: 8px;
}

@media (max-width: 768px) {
  .navbar {
    padding: 0 12px;
  }

  .hide-mobile {
    display: none !important;
  }

  .avatar-wrapper .user-meta,
  .avatar-wrapper .caret {
    display: none;
  }

  .crumb-root .crumb-text {
    display: none;
  }
}

@media (max-width: 480px) {
  .nav-text-btn .language-text {
    display: none;
  }

  .breadcrumb .crumb:not(.crumb-current):not(.crumb-root),
  .breadcrumb .crumb-sep {
    display: none;
  }

  .breadcrumb .crumb-sep:has(+ .crumb-current) {
    display: inline-flex;
  }
}
</style>
