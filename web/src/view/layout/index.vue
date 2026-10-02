<template>
  <div
    class="app-wrapper"
    :class="{ 'mobile': isMobile, 'has-topbar-announcement': hasTopbarAnnouncement }"
  >
    <!-- 顶部栏公告 -->
    <TopbarAnnouncement @visible-change="hasTopbarAnnouncement = $event" />
    
    <!-- 移动端遮罩层 -->
    <div
      v-if="isMobile && sidebar.opened"
      class="drawer-bg"
      @click="closeSidebar"
    />
    
    <!-- 侧边栏 -->
    <component
      :is="Sidebar"
      :key="userStore.userType"
      class="sidebar-container"
      :class="{ 
        'is-collapse': isCollapse && !isMobile,
        'mobile': isMobile,
        'hidden': isMobile && !sidebar.opened
      }"
    />
    
    <!-- 主容器 -->
    <div
      class="main-container"
      :class="{ 
        'main-container-collapsed': isCollapse && !isMobile,
        'mobile': isMobile
      }"
    >
      <div
        class="fixed-header"
        :class="{ 
          'fixed-header-collapsed': isCollapse && !isMobile,
          'mobile': isMobile
        }"
      >
        <navbar @toggle-sidebar="toggleSidebar" />
      </div>
      <app-main />
      <app-footer />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick, provide } from 'vue'
import { Navbar, Sidebar, AppMain, AppFooter } from './components'
import { useUserStore } from '@/pinia/modules/user'
import TopbarAnnouncement from '@/components/TopbarAnnouncement.vue'
import { shouldUseSidebarDrawer } from '@/utils/layout'

const userStore = useUserStore()
const SIDEBAR_COLLAPSE_STORAGE_KEY = 'sidebarCollapsed'
const isMobile = ref(false)
const sidebar = ref({
  opened: true
})
const isCollapse = ref(true)
const hasTopbarAnnouncement = ref(false)
let deviceModeInitialized = false

const readStoredCollapse = () => {
  const stored = localStorage.getItem(SIDEBAR_COLLAPSE_STORAGE_KEY)
  if (stored === null) return true
  return stored === 'true'
}

const saveStoredCollapse = (collapsed) => {
  localStorage.setItem(SIDEBAR_COLLAPSE_STORAGE_KEY, String(collapsed))
}

// 检测设备类型
const checkDevice = () => {
  const useDrawer = shouldUseSidebarDrawer(window.innerWidth, window.innerHeight)
  if (deviceModeInitialized && useDrawer === isMobile.value) return

  isMobile.value = useDrawer
  
  // 手机和竖屏平板使用完整宽度抽屉，避免固定窄栏截断菜单名称。
  if (isMobile.value) {
    sidebar.value.opened = false
    isCollapse.value = false
  } else {
    sidebar.value.opened = true
    isCollapse.value = readStoredCollapse()
  }
  deviceModeInitialized = true
}

// 切换侧边栏
const toggleSidebar = () => {
  if (isMobile.value) {
    sidebar.value.opened = !sidebar.value.opened
  } else {
    isCollapse.value = !isCollapse.value
    saveStoredCollapse(isCollapse.value)
    if (toggleSidebarCollapse) {
      toggleSidebarCollapse(isCollapse.value)
    }
  }
}

// 关闭侧边栏（移动端）
const closeSidebar = () => {
  sidebar.value.opened = false
}

// 提供给子组件的方法
const toggleSidebarCollapse = (collapsed) => {
  if (!isMobile.value) {
    isCollapse.value = collapsed
    saveStoredCollapse(collapsed)
  }
}

// 提供收缩状态和移动端状态给子组件
provide('toggleSidebarCollapse', toggleSidebarCollapse)
provide('sidebarCollapsed', computed(() => isCollapse.value))
provide('isMobile', computed(() => isMobile.value))
provide('sidebarOpened', computed(() => sidebar.value.opened))
provide('closeSidebar', closeSidebar)

onMounted(() => {
  checkDevice()
  window.addEventListener('resize', checkDevice)
  
  nextTick(() => {
    const sidebarEl = document.querySelector('.sidebar-container')
    if (!sidebarEl || sidebarEl.children.length === 0) {
      userStore.$patch({ userType: userStore.userType })
    }
  })
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', checkDevice)
})
</script>

<style lang="scss" scoped>
.app-wrapper {
  position: relative;
  min-height: 100%;
  min-height: 100dvh;
  width: 100%;
  background-color: var(--bg-color-primary);
  --topbar-announcement-height: 48px;

  &.mobile {
    overflow-x: hidden;
    --sidebar-width: min(280px, calc(100vw - 40px));
  }

  &.has-topbar-announcement {
    .fixed-header {
      top: var(--topbar-announcement-height);
    }

    .sidebar-container {
      top: var(--topbar-announcement-height);
      height: calc(100% - var(--topbar-announcement-height));
    }
  }
}

.drawer-bg {
  background: rgba(0, 0, 0, 0.3);
  width: 100%;
  top: 0;
  height: 100%;
  position: fixed;
  z-index: var(--z-drawer-bg);
}

.fixed-header {
  position: fixed;
  top: 0;
  right: 0;
  z-index: var(--z-navbar);
  width: calc(100% - var(--sidebar-width));
  transition: width 0.28s;
  background-color: var(--bg-color-secondary);
  box-shadow: var(--box-shadow-light);
  border-bottom: 1px solid var(--border-color);
  
  &.fixed-header-collapsed {
    width: calc(100% - var(--sidebar-width-collapsed));
  }
  
  &.mobile {
    width: 100%;
  }
}

.sidebar-container {
  transition: transform 0.28s, width 0.28s;
  width: var(--sidebar-width);
  background-color: var(--bg-color-sidebar);
  height: 100%;
  position: fixed;
  font-size: 0px;
  top: 0;
  bottom: 0;
  left: 0;
  z-index: var(--z-sidebar);
  overflow: hidden;
  box-shadow: 2px 0 6px rgba(0, 0, 0, 0.1);
  
  &.is-collapse {
    width: var(--sidebar-width-collapsed);
  }
  
  &.mobile {
    width: var(--sidebar-width);
    transform: translateX(0);
    
    &.hidden {
      transform: translateX(-100%);
    }
  }
}

.main-container {
  min-height: 100vh;
  min-height: 100dvh;
  transition: margin-left 0.28s;
  margin-left: var(--sidebar-width);
  position: relative;
  padding-top: var(--navbar-height);
  padding-bottom: env(safe-area-inset-bottom);
  display: flex;
  flex-direction: column;
  
  &.main-container-collapsed {
    margin-left: var(--sidebar-width-collapsed);
  }
  
  &.mobile {
    margin-left: 0;
    width: 100%;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .app-wrapper {
    --topbar-announcement-height: 48px;
  }

  .sidebar-container {
    width: var(--sidebar-width);
  }
  
  .main-container {
    margin-left: 0;
  }
  
  .fixed-header {
    width: 100%;
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
