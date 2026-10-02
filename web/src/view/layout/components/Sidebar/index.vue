<template>
  <div
    class="sidebar-container"
    :class="{ 
      'is-collapse': isCollapse && !isMobile,
      'mobile': isMobile
    }"
  >
    <div class="sidebar-logo">
      <img
        v-show="!isCollapse || isMobile"
        :src="siteStore.logoSrc"
        alt="Logo"
        class="sidebar-logo-img"
      >
      <h1 v-show="!isCollapse || isMobile">
        {{ siteStore.displaySiteName }}
      </h1>
      <el-button 
        v-if="!isMobile"
        class="collapse-btn" 
        :icon="isCollapse ? Expand : Fold" 
        size="small" 
        circle 
        @click="toggleCollapse" 
      />
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
        active-text-color="#16a34a"
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
const activeMenu = computed(() => {
  return route.path
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
        path: '/user/instances',
        name: 'UserInstances',
        meta: {
          title: 'sidebar.myInstances',
          icon: 'Box'
        }
      },
      {
        path: '/user/apply',
        name: 'UserApply',
        meta: {
          title: 'sidebar.apply',
          icon: 'Plus'
        }
      },
      {
        path: '/user/tasks',
        name: 'UserTasks',
        meta: {
          title: 'sidebar.taskList',
          icon: 'List'
        }
      },
      {
        path: '/user/profile',
        name: 'UserProfile',
        meta: {
          title: 'sidebar.personalCenter',
          icon: 'User'
        }
      },
      {
        path: '/user/domain',
        name: 'UserDomain',
        meta: {
          title: 'sidebar.domainBinding',
          icon: 'Link'
        }
      },
      {
        path: '/user/kyc',
        name: 'UserKYC',
        meta: {
          title: 'sidebar.kycVerification',
          icon: 'Postcard'
        }
      },
      {
        path: '/user/checkin',
        name: 'UserCheckin',
        meta: {
          title: 'sidebar.checkinRenewal',
          icon: 'Calendar'
        }
      },
      {
        path: '/user/api-tokens',
        name: 'UserApiTokens',
        meta: {
          title: 'sidebar.apiTokenManagement',
          icon: 'Key'
        }
      },
      {
        path: '/user/store',
        name: 'UserStore',
        meta: {
          title: 'sidebar.store',
          icon: 'Shop'
        }
      },
      {
        path: '/user/orders',
        name: 'UserOrders',
        meta: {
          title: 'sidebar.orders',
          icon: 'Document'
        }
      },
      {
        path: '/user/tickets',
        name: 'UserTickets',
        meta: {
          title: 'sidebar.tickets',
          icon: 'Service'
        }
      },
      {
        path: '/user/wallet',
        name: 'UserWallet',
        meta: {
          title: 'sidebar.wallet',
          icon: 'Wallet'
        }
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

const defaultOpeneds = computed(() => {
  return []
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
.sidebar-container {
  transition: width 0.28s;
  width: var(--sidebar-width);
  background-color: var(--bg-color-sidebar-light);

  .sidebar-logo {
    height: var(--navbar-height);
    line-height: var(--navbar-height);
    background: #16a34a; /* 绿色背景 */
    text-align: center;
    overflow: hidden;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-start;
    padding: 0 var(--spacing-md);
    position: relative;

    h1 {
      color: #ffffff; /* 白色文字 */
      font-weight: var(--font-weight-semibold);
      font-size: var(--font-size-md);
      font-family: Avenir, Helvetica Neue, Arial, Helvetica, sans-serif;
      margin: 0;
      transition: opacity 0.28s;
    }
    
    .sidebar-logo-img {
      height: 28px;
      width: 28px;
      border-radius: 4px;
      object-fit: contain;
      flex-shrink: 0;
      margin-right: 8px;
    }
    
    span {
      font-size: var(--font-size-xs);
      color: #dcfce7; /* 浅绿色文字 */
    }

    .collapse-btn {
      position: absolute;
      top: 50%;
      right: 10px;
      transform: translateY(-50%);
      color: #dcfce7; /* 浅绿色 */
      background: transparent;
      border: none;
      transition: all 0.28s;
      
      &:hover {
        color: #ffffff; /* 悬停时白色 */
      }
    }
  }

  .scrollbar-wrapper {
    overflow-x: hidden !important;
  }

  .el-scrollbar__bar.is-vertical {
    right: 0px;
  }

  .el-scrollbar {
    height: calc(100% - var(--navbar-height));
  }

  .is-horizontal {
    display: none;
  }

  a {
    display: inline-block;
    width: 100%;
    overflow: hidden;
  }

  .svg-icon {
    margin-right: 16px;
  }

  .sub-el-icon {
    margin-right: 12px;
    margin-left: -2px;
  }

  .el-menu {
    border: none;
    height: 100%;
    background-color: var(--bg-color-sidebar-light) !important;
  }

  /* 菜单项悬停效果 */
  :deep(.el-menu-item) {
    height: 48px;
    line-height: 48px;
    background-color: transparent !important;
    
    &:hover {
      background-color: var(--bg-color-hover) !important;
      color: #16a34a !important;
    }
    
    &.is-active {
      background-color: var(--bg-color-active) !important;
      color: #16a34a !important;
      border-right: 3px solid #16a34a;
    }
  }

  :deep(.el-sub-menu__title) {
    height: 48px;
    line-height: 48px;
    background-color: transparent !important;
    color: var(--text-color-sidebar-primary) !important;
    font-weight: var(--font-weight-semibold);
    
    &:hover {
      background-color: var(--bg-color-hover) !important;
      color: #16a34a !important;
    }
  }

  :deep(.el-sub-menu.is-opened > .el-sub-menu__title),
  :deep(.el-sub-menu.is-active > .el-sub-menu__title) {
    color: #16a34a !important;
  }

  :deep(.el-sub-menu .el-menu) {
    background-color: var(--bg-color-sidebar-light) !important;
  }

  :deep(.el-sub-menu .el-menu-item) {
    height: 42px;
    line-height: 42px;
    padding-left: 44px !important;
    font-size: 13px;
  }

  // 收缩状态样式
  &.is-collapse {
    width: var(--sidebar-width-collapsed);

    :deep(.el-menu) {
      width: var(--sidebar-width-collapsed) !important;
    }

    :deep(.el-menu-item),
    :deep(.el-sub-menu__title) {
      justify-content: center;
      padding-left: 0 !important;
      padding-right: 0 !important;
    }

    :deep(.el-menu-item .el-menu-tooltip__trigger),
    :deep(.el-sub-menu__title .el-menu-tooltip__trigger) {
      justify-content: center;
    }

    :deep(.menu-title),
    :deep(.el-sub-menu__icon-arrow),
    :deep(.el-menu-item > span:not(.el-icon)),
    :deep(.el-sub-menu__title > span:not(.el-icon)) {
      display: none !important;
      width: 0 !important;
      min-width: 0 !important;
      overflow: hidden !important;
    }

    :deep(.menu-item) {
      justify-content: center;
      width: var(--sidebar-width-collapsed);
    }

    :deep(.menu-icon),
    :deep(.el-icon) {
      margin-right: 0 !important;
    }

    .sidebar-logo {
      .collapse-btn {
        right: 50%;
        transform: translate(50%, -50%);
      }
    }
  }
  
  // 移动端样式
  &.mobile {
    width: var(--sidebar-width);
    
    .sidebar-logo {
      .collapse-btn {
        display: none;
      }
    }
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .sidebar-container {
    .sidebar-logo {
      h1 {
        font-size: var(--font-size-base);
      }
    }
    
    :deep(.el-menu-item),
    :deep(.el-sub-menu__title) {
      height: 48px;
      line-height: 48px;
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
