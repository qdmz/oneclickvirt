<template>
  <div
    class="sidebar-container oh-sidebar"
    :class="{ 
      'is-collapse': isCollapse && !isMobile,
      'mobile': isMobile
    }"
  >
    <div class="sidebar-logo">
      <router-link
        :to="homePath"
        class="sidebar-brand"
        :title="siteStore.displaySiteName"
      >
        <span class="sidebar-logo-mark">
          <img
            :src="siteStore.logoSrc"
            alt="Logo"
            class="sidebar-logo-img"
          >
        </span>
        <span
          v-show="!isCollapse || isMobile"
          class="sidebar-brand-text"
        >
          <h1>{{ siteStore.displaySiteName }}</h1>
          <small>{{ isAdminView ? t('ui.layout.adminConsole') : t('ui.layout.console') }}</small>
        </span>
      </router-link>
    </div>
    <el-scrollbar wrap-class="scrollbar-wrapper">
      <el-menu
        :key="menuRenderKey"
        :default-active="activeMenu"
        :collapse="isCollapse && !isMobile"
        :unique-opened="false"
        :default-openeds="defaultOpeneds"
        :collapse-transition="false"
        mode="vertical"
        @select="handleMenuSelect"
      >
        <!-- 首页链接 - 仅在未登录时显示 -->
        <el-menu-item
          v-if="!userStore.isLoggedIn"
          index="/home"
        >
          <el-icon><HomeFilled /></el-icon>
          <template #title>
            {{ t('navbar.home') }}
          </template>
        </el-menu-item>
        
        <!-- 动态生成的菜单项 -->
        <sidebar-item
          v-for="route in userRoutes"
          :key="route.path"
          :item="route"
          :base-path="route.path"
          :is-collapse="isCollapse && !isMobile"
        />
      </el-menu>
    </el-scrollbar>
    <div
      v-if="!isMobile"
      class="sidebar-footer"
    >
      <button
        type="button"
        class="collapse-btn"
        :title="isCollapse ? t('ui.layout.expand') : t('ui.layout.collapse')"
        :aria-label="isCollapse ? t('ui.layout.expand') : t('ui.layout.collapse')"
        @click="toggleCollapse"
      >
        <el-icon><component :is="isCollapse ? Expand : Fold" /></el-icon>
        <span v-show="!isCollapse">{{ t('ui.layout.collapse') }}</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, watch, nextTick, inject, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/pinia/modules/user'
import { HomeFilled, Expand, Fold } from '@element-plus/icons-vue'
import SidebarItem from './SidebarItem.vue'
import { useSiteStore } from '@/pinia/modules/site'
import { useFeatureStore } from '@/pinia/modules/feature'

const route = useRoute()
const { t, locale } = useI18n()
const userStore = useUserStore()
const siteStore = useSiteStore()
const featureStore = useFeatureStore()

// 从父组件注入的状态和方法
const toggleSidebarCollapse = inject('toggleSidebarCollapse', null)
const sidebarCollapsed = inject('sidebarCollapsed', computed(() => false))
const isMobile = inject('isMobile', ref(false))
const closeSidebar = inject('closeSidebar', null)
const isCollapse = computed(() => sidebarCollapsed.value)

const toggleCollapse = () => {
  if (toggleSidebarCollapse) {
    toggleSidebarCollapse(!isCollapse.value)
  }
}

// 移动端点击菜单后关闭侧边栏
const handleMenuSelect = () => {
  if (isMobile.value && closeSidebar) {
    closeSidebar()
  }
}

// 获取当前活动菜单
// 详情页（如 /user/instances/123）也高亮其所属菜单
const activeMenu = computed(() => {
  const flat = []
  const walk = (items) => items.forEach(i => { if (i.children?.length) walk(i.children); else flat.push(i.path) })
  walk(userRoutes.value || [])
  if (flat.includes(route.path)) return route.path
  const parent = flat
    .filter(p => route.path.startsWith(p + '/'))
    .sort((a, b) => b.length - a.length)[0]
  return parent || route.path
})

// 根据用户类型获取对应的路由
const userRoutes = computed(() => {
  // 使用 viewMode 来决定显示哪个视图的菜单
  // 管理员(含normal_admin)可以切换视图，普通用户只能看到用户视图
  const viewMode = userStore.currentViewMode || userStore.userType
  // normal_admin 也使用 admin 路由集，但会过滤掉超管专属项
  const effectiveMode = (viewMode === 'normal_admin' || viewMode === 'admin') ? 'admin' : viewMode
  
  // 强制依赖 locale，确保语言切换时重新计算
  locale.value
  
  // 用户特定路由
  const userTypeRoutes = {
    // 普通用户路由
    user: [
      {
        path: '/user/dashboard',
        name: 'UserDashboard',
        meta: {
          title: 'sidebar.dashboard',
          icon: 'Odometer'
        }
      },
      {
        path: '/user/_cloud',
        name: 'UserMenuCloud',
        alwaysShow: true,
        meta: { title: 'sidebar.groupCloud', icon: 'Cloudy' },
        children: [
          { path: '/user/instances', name: 'UserInstances', meta: { title: 'sidebar.myInstances', icon: 'Box' } },
          { path: '/user/apply', name: 'UserApply', meta: { title: 'sidebar.apply', icon: 'Plus' } },
          { path: '/user/tasks', name: 'UserTasks', meta: { title: 'sidebar.taskList', icon: 'List' } }
        ]
      },
      {
        path: '/user/_billing',
        name: 'UserMenuBilling',
        alwaysShow: true,
        meta: { title: 'sidebar.groupBilling', icon: 'ShoppingBag' },
        children: [
          { path: '/user/store', name: 'UserStore', meta: { title: 'sidebar.store', icon: 'Shop' } },
          { path: '/user/orders', name: 'UserOrders', meta: { title: 'sidebar.orders', icon: 'Document' } },
          { path: '/user/wallet', name: 'UserWallet', meta: { title: 'sidebar.wallet', icon: 'Wallet' } }
        ]
      },
      {
        path: '/user/_service',
        name: 'UserMenuService',
        alwaysShow: true,
        meta: { title: 'sidebar.groupService', icon: 'Service' },
        children: [
          { path: '/user/tickets', name: 'UserTickets', meta: { title: 'sidebar.tickets', icon: 'ChatDotSquare' } },
          { path: '/user/checkin', name: 'UserCheckin', meta: { title: 'sidebar.checkinRenewal', icon: 'Calendar' } },
          { path: '/user/domain', name: 'UserDomain', meta: { title: 'sidebar.domainBinding', icon: 'Link' } }
        ]
      },
      {
        path: '/user/_account',
        name: 'UserMenuAccount',
        alwaysShow: true,
        meta: { title: 'sidebar.groupAccount', icon: 'User' },
        children: [
          { path: '/user/profile', name: 'UserProfile', meta: { title: 'sidebar.personalCenter', icon: 'User' } },
          { path: '/user/kyc', name: 'UserKYC', meta: { title: 'sidebar.kycVerification', icon: 'Postcard' } },
          { path: '/user/api-tokens', name: 'UserApiTokens', meta: { title: 'sidebar.apiTokenManagement', icon: 'Key' } }
        ]
      }
    ],
    // 管理员路由
    admin: [
      {
        path: '/admin/_overview',
        name: 'AdminMenuOverview',
        alwaysShow: true,
        meta: {
          title: 'sidebar.groupOverview',
          icon: 'Odometer'
        },
        children: [
          {
            path: '/admin/dashboard',
            name: 'AdminDashboard',
            meta: {
              title: 'sidebar.dashboard',
              icon: 'Odometer'
            }
          },
          {
            path: '/admin/performance',
            name: 'AdminPerformance',
            meta: {
              title: 'sidebar.performanceMonitoring',
              icon: 'Histogram'
            }
          },
          {
            path: '/admin/logs',
            name: 'AdminLogs',
            meta: {
              title: 'sidebar.logViewer',
              icon: 'Document'
            }
          }
        ]
      },
      {
        path: '/admin/_users',
        name: 'AdminMenuUsers',
        alwaysShow: true,
        meta: {
          title: 'sidebar.groupUsers',
          icon: 'User'
        },
        children: [
          {
            path: '/admin/users',
            name: 'AdminUsers',
            meta: {
              title: 'sidebar.userManagement',
              icon: 'User'
            }
          },
          {
            path: '/admin/invite-codes',
            name: 'AdminInviteCodes',
            meta: {
              title: 'sidebar.inviteCodeManagement',
              icon: 'Ticket'
            }
          },
          {
            path: '/admin/redemption-codes',
            name: 'AdminRedemptionCodes',
            meta: {
              title: 'sidebar.redemptionCodeManagement',
              icon: 'Discount'
            }
          },
          {
            path: '/admin/vouchers',
            name: 'AdminVouchers',
            meta: {
              title: 'sidebar.voucherManagement',
              icon: 'Wallet'
            }
          }
        ]
      },
      {
        path: '/admin/_resources',
        name: 'AdminMenuResources',
        alwaysShow: true,
        meta: {
          title: 'sidebar.groupResources',
          icon: 'Box'
        },
        children: [
          {
            path: '/admin/instances',
            name: 'AdminInstances',
            meta: {
              title: 'sidebar.instanceManagement',
              icon: 'Box'
            }
          },
          {
            path: '/admin/providers',
            name: 'AdminProviders',
            meta: {
              title: 'sidebar.providerManagement',
              icon: 'Monitor'
            }
          },
          {
            path: '/admin/group',
            name: 'AdminGroup',
            meta: {
              title: 'sidebar.groupManagement',
              icon: 'Collection'
            }
          },
          {
            path: '/admin/port-mappings',
            name: 'AdminPortMappings',
            meta: {
              title: 'sidebar.portManagement',
              icon: 'Connection'
            }
          }
        ]
      },
      {
        path: '/admin/_images',
        name: 'AdminMenuImages',
        alwaysShow: true,
        meta: {
          title: 'sidebar.groupImages',
          icon: 'Folder'
        },
        children: [
          {
            path: '/admin/system-images',
            name: 'AdminSystemImages',
            meta: {
              title: 'sidebar.systemImages',
              icon: 'Folder'
            }
          },
          {
            path: '/admin/snapshots',
            name: 'AdminSnapshots',
            meta: {
              title: 'sidebar.snapshotManagement',
              icon: 'Camera'
            }
          }
        ]
      },
      {
        path: '/admin/_tasks',
        name: 'AdminMenuTasks',
        alwaysShow: true,
        meta: {
          title: 'sidebar.groupTasks',
          icon: 'List'
        },
        children: [
          {
            path: '/admin/traffic',
            name: 'AdminTraffic',
            meta: {
              title: 'sidebar.trafficManagement',
              icon: 'TrendCharts'
            }
          },
          {
            path: '/admin/tasks',
            name: 'AdminTasks',
            meta: {
              title: 'sidebar.taskManagement',
              icon: 'List'
            }
          }
        ]
      },
      {
        path: '/admin/_store',
        name: 'AdminMenuStore',
        alwaysShow: true,
        meta: {
          title: 'sidebar.groupStore',
          icon: 'ShoppingBag'
        },
        children: [
          {
            path: '/admin/products',
            name: 'AdminProducts',
            meta: {
              title: 'sidebar.productManagement',
              icon: 'Goods'
            }
          },
          {
            path: '/admin/orders',
            name: 'AdminOrders',
            meta: {
              title: 'sidebar.orderManagement',
              icon: 'DocumentChecked'
            }
          }
        ]
      },
      {
        path: '/admin/_tickets',
        name: 'AdminMenuTickets',
        alwaysShow: true,
        meta: {
          title: 'sidebar.groupTickets',
          icon: 'ChatDotSquare'
        },
        children: [
          {
            path: '/admin/tickets',
            name: 'AdminTickets',
            meta: {
              title: 'sidebar.ticketManagement',
              icon: 'Service'
            }
          }
        ]
      },
      {
        path: '/admin/_system',
        name: 'AdminMenuSystem',
        alwaysShow: true,
        meta: {
          title: 'sidebar.groupSystem',
          icon: 'Setting'
        },
        children: [
          {
            path: '/admin/config',
            name: 'AdminConfig',
            meta: {
              title: 'sidebar.systemConfiguration',
              icon: 'Setting'
            }
          },
          {
            path: '/admin/site-config',
            name: 'AdminSiteConfig',
            meta: {
              title: 'sidebar.siteConfig',
              icon: 'Monitor'
            }
          },
          {
            path: '/admin/site-link',
            name: 'AdminSiteLink',
            meta: {
              title: 'sidebar.siteLink',
              icon: 'Link'
            }
          },
          {
            path: '/admin/yipay-config',
            name: 'AdminYipayConfig',
            meta: {
              title: 'sidebar.yipayConfig',
              icon: 'Wallet'
            }
          },
          {
            path: '/admin/upstream',
            name: 'AdminUpstream',
            meta: {
              title: 'sidebar.upstreamManagement',
              icon: 'Connection'
            }
          },
          {
            path: '/admin/block-rules',
            name: 'AdminBlockRules',
            meta: {
              title: 'sidebar.blockRulesManagement',
              icon: 'Lock'
            }
          },
          {
            path: '/admin/api-tokens',
            name: 'AdminApiTokens',
            meta: {
              title: 'sidebar.adminApiTokenManagement',
              icon: 'Key'
            }
          },
          {
            path: '/admin/oauth2-providers',
            name: 'AdminOAuth2Providers',
            meta: {
              title: 'sidebar.oauth2Management',
              icon: 'Connection'
            }
          },
          {
            path: '/admin/announcements',
            name: 'AdminAnnouncements',
            meta: {
              title: 'sidebar.announcementManagement',
              icon: 'Bell'
            }
          }
        ]
      }
    ]
  }
  
  // 根据视图模式返回对应路由
  const routes = userTypeRoutes[effectiveMode] || []
  
  // 超级管理员专属路由名称集（normal_admin 不可见）
  const superAdminOnlyRoutes = new Set([
    'AdminUsers', 'AdminConfig', 'AdminPerformance', 'AdminLogs', 'AdminOAuth2Providers',
    'AdminInviteCodes', 'AdminAnnouncements', 'AdminSystemImages', 'AdminApiTokens', 'AdminKYC'
  ])
  
  // 判断是否为普通管理员
  const isNormalAdmin = userStore.userType === 'normal_admin'
  
  const shouldShowRoute = (route) => {
    if (['UserKYC', 'AdminKYC'].includes(route.name) && !featureStore.kycEnabled) return false
    if (['UserDomain', 'AdminDomain'].includes(route.name) && !featureStore.domainEnabled) return false
    if (['UserCheckin'].includes(route.name) && !featureStore.checkinEnabled) return false
    if (isNormalAdmin && superAdminOnlyRoutes.has(route.name)) return false
    return true
  }

  const filterRoutes = (menuRoutes) => {
    return menuRoutes.reduce((result, menuRoute) => {
      if (menuRoute.children?.length) {
        const children = filterRoutes(menuRoute.children)
        if (children.length > 0) {
          result.push({ ...menuRoute, children })
        }
        return result
      }

      if (shouldShowRoute(menuRoute)) {
        result.push(menuRoute)
      }
      return result
    }, [])
  }

  const filteredRoutes = filterRoutes(routes)
  
  return filteredRoutes
})

const isAdminView = computed(() => {
  const viewMode = userStore.currentViewMode || userStore.userType
  return viewMode === 'admin' || viewMode === 'normal_admin'
})

const homePath = computed(() => {
  if (!userStore.isLoggedIn) return '/home'
  return isAdminView.value ? '/admin/dashboard' : '/user/dashboard'
})

// 默认展开：用户视图展开全部分组（菜单项少，一屏可见）；管理员视图仅展开当前页面所在分组
const defaultOpeneds = computed(() => {
  const groups = userRoutes.value.filter(r => r.children?.length)
  if (!isAdminView.value) return groups.map(g => g.path)
  const current = groups.find(g => g.children.some(c => route.path === c.path || route.path.startsWith(c.path + '/')))
  return current ? [current.path] : []
})

const menuRenderKey = computed(() => {
  const viewMode = userStore.currentViewMode || userStore.userType || 'guest'
  const collapseState = isCollapse.value && !isMobile.value ? 'collapsed' : 'expanded'
  return `${viewMode}-${locale.value}-${collapseState}-${defaultOpeneds.value.join('|')}`
})

// 生命周期钩子，检查DOM渲染
onMounted(() => {
  // 确保组件在DOM中
  nextTick(() => {
    document.querySelector('.sidebar-container')
  })
})

// 监听用户类型变化
watch([
  () => userStore.userType,
  () => userStore.currentViewMode,
  () => featureStore.kycEnabled,
  () => featureStore.domainEnabled,
  () => featureStore.checkinEnabled
], () => {
  nextTick(() => {
    userRoutes.value
  })
}, { immediate: true })
</script>

<style lang="scss" scoped>
.oh-sidebar {
  display: flex;
  flex-direction: column;
  width: var(--sidebar-width);
  background-color: var(--bg-color-sidebar);
  border-right: 1px solid var(--border-color);
  box-shadow: none;
  transition: width 0.26s var(--ease-in-out), transform 0.26s var(--ease-in-out);

  .sidebar-logo {
    height: var(--navbar-height);
    flex-shrink: 0;
    display: flex;
    align-items: center;
    padding: 0 14px;
    border-bottom: 1px solid var(--border-color-light);
    background: transparent;
  }

  .sidebar-brand {
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 0;
    width: 100%;
    color: inherit;
    text-decoration: none;
  }

  .sidebar-logo-mark {
    width: 36px;
    height: 36px;
    flex-shrink: 0;
    border-radius: 10px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background: var(--brand-gradient-soft);
    border: 1px solid color-mix(in srgb, var(--primary-color) 18%, transparent);
    box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.4);
  }

  .sidebar-logo-img {
    width: 24px;
    height: 24px;
    object-fit: contain;
    border-radius: 6px;
  }

  .sidebar-brand-text {
    display: flex;
    flex-direction: column;
    min-width: 0;
    line-height: 1.2;

    h1 {
      margin: 0;
      font-size: 15px;
      font-weight: 700;
      letter-spacing: -0.01em;
      color: var(--text-color-primary);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    small {
      margin-top: 2px;
      font-size: 11px;
      font-weight: 500;
      color: var(--text-color-tertiary);
      letter-spacing: 0.04em;
      text-transform: uppercase;
    }
  }

  .el-scrollbar {
    flex: 1;
    height: auto;
    min-height: 0;
  }

  :deep(.scrollbar-wrapper) {
    overflow-x: hidden !important;
  }

  :deep(.el-menu) {
    border: none;
    background-color: transparent !important;
    padding: 10px 10px 16px;
    --el-menu-item-height: 40px;
    --el-menu-sub-item-height: 38px;
  }

  :deep(.el-menu .el-menu) {
    padding: 2px 0 4px;
  }

  :deep(a) {
    display: block;
    width: 100%;
    text-decoration: none;
  }

  /* 一级与二级菜单项 */
  :deep(.el-menu-item),
  :deep(.el-sub-menu__title) {
    height: 40px;
    line-height: 40px;
    margin: 2px 0;
    padding: 0 12px !important;
    border-radius: 8px;
    color: var(--text-color-sidebar-inactive) !important;
    background-color: transparent !important;
    font-size: 13.5px;
    font-weight: 500;
    transition: background-color 0.15s ease, color 0.15s ease;
  }

  :deep(.el-sub-menu .el-menu-item) {
    height: 38px;
    line-height: 38px;
    padding-left: 38px !important;
    font-size: 13px;
    position: relative;
  }

  :deep(.el-sub-menu .el-menu-item)::before {
    content: '';
    position: absolute;
    left: 21px;
    top: 0;
    bottom: 0;
    width: 1px;
    background: var(--border-color);
  }

  :deep(.el-menu-item:hover),
  :deep(.el-sub-menu__title:hover) {
    background-color: var(--bg-color-hover) !important;
    color: var(--text-color-primary) !important;
  }

  :deep(.el-menu-item.is-active) {
    background-color: var(--bg-color-active) !important;
    color: var(--accent-text-color) !important;
    font-weight: 600;
  }

  :deep(.el-sub-menu .el-menu-item.is-active)::before {
    width: 2px;
    left: 20.5px;
    top: 9px;
    bottom: 9px;
    border-radius: 2px;
    background: var(--brand-gradient);
  }

  :deep(> .el-scrollbar .el-menu > div > a > .el-menu-item.is-active),
  :deep(.el-menu > div > a > .el-menu-item.is-active) {
    box-shadow: inset 3px 0 0 var(--primary-color);
  }

  :deep(.el-sub-menu.is-active > .el-sub-menu__title) {
    color: var(--text-color-primary) !important;
    font-weight: 600;
  }

  :deep(.el-sub-menu.is-active > .el-sub-menu__title .menu-icon) {
    color: var(--primary-color);
  }

  :deep(.el-sub-menu__icon-arrow) {
    right: 10px;
    font-size: 12px;
    color: var(--text-color-tertiary);
  }

  :deep(.menu-icon) {
    font-size: 17px;
    margin-right: 10px;
  }

  :deep(.el-sub-menu .el-menu-item .menu-icon) {
    font-size: 15px;
    opacity: 0.9;
  }

  .sidebar-footer {
    flex-shrink: 0;
    padding: 10px;
    border-top: 1px solid var(--border-color-light);
  }

  .collapse-btn {
    width: 100%;
    height: 36px;
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 0 12px;
    border: none;
    border-radius: 8px;
    background: transparent;
    color: var(--text-color-tertiary);
    font-size: 13px;
    cursor: pointer;
    font-family: inherit;
    transition: background-color 0.15s ease, color 0.15s ease;

    .el-icon {
      font-size: 16px;
    }

    &:hover {
      background: var(--bg-color-hover);
      color: var(--text-color-primary);
    }
  }

  /* 收缩状态 */
  &.is-collapse {
    width: var(--sidebar-width-collapsed);

    .sidebar-logo {
      padding: 0;
      justify-content: center;
    }

    .sidebar-brand {
      justify-content: center;
    }

    :deep(.el-menu) {
      width: var(--sidebar-width-collapsed) !important;
      padding: 10px 8px;
    }

    :deep(.el-menu-item),
    :deep(.el-sub-menu__title) {
      justify-content: center;
      padding: 0 !important;
    }

    :deep(.el-menu-item .el-menu-tooltip__trigger),
    :deep(.el-sub-menu__title .el-menu-tooltip__trigger) {
      justify-content: center;
      padding: 0 !important;
    }

    :deep(.menu-title),
    :deep(.el-sub-menu__icon-arrow),
    :deep(.el-menu-item > span:not(.el-icon)),
    :deep(.el-sub-menu__title > span:not(.el-icon)) {
      display: none !important;
    }

    :deep(.menu-item) {
      justify-content: center;
    }

    :deep(.menu-icon),
    :deep(.el-icon) {
      margin-right: 0 !important;
    }

    .collapse-btn {
      justify-content: center;
      padding: 0;
    }
  }

  &.mobile {
    width: var(--sidebar-width);
    box-shadow: var(--shadow-xl);
  }
}
</style>
