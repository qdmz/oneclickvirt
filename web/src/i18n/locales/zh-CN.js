// 中文（简体）语言包 - 模块化结构
// 公共模块
import common from './zh-CN/common.js'
import navbar from './zh-CN/navbar.js'
import validation from './zh-CN/validation.js'
import message from './zh-CN/message.js'
import errors from './zh-CN/errors.js'

// 认证模块
import login from './zh-CN/auth/login.js'
import adminLogin from './zh-CN/auth/adminLogin.js'
import register from './zh-CN/auth/register.js'
import forgotPassword from './zh-CN/auth/forgotPassword.js'
import resetPassword from './zh-CN/auth/resetPassword.js'
import oauth2Callback from './zh-CN/auth/oauth2Callback.js'
import init from './zh-CN/auth/init.js'

// 公共页面模块
import home from './zh-CN/public/home.js'
import sidebar from './zh-CN/public/sidebar.js'
import notFound from './zh-CN/public/notFound.js'

// 用户模块
import userDashboard from './zh-CN/user/dashboard.js'
import userProfile from './zh-CN/user/profile.js'
import userInstances from './zh-CN/user/instances.js'
import userInstanceDetail from './zh-CN/user/instanceDetail.js'
import userTasks from './zh-CN/user/tasks.js'
import userTrafficOverview from './zh-CN/user/trafficOverview.js'
import userTraffic from './zh-CN/user/traffic.js'
import userResources from './zh-CN/user/resources.js'
import userApply from './zh-CN/user/apply.js'
import userDomain from './zh-CN/user/domain.js'
import userKyc from './zh-CN/user/kyc.js'
import userCheckin from './zh-CN/user/checkin.js'
import userApiTokens from './zh-CN/user/apiTokens.js'
import userStore from './zh-CN/user/store.js'
import userOrders from './zh-CN/user/orders.js'
import userTickets from './zh-CN/user/tickets.js'
import userWallet from './zh-CN/user/wallet.js'

// 管理员模块
import adminDashboard from './zh-CN/admin/dashboard.js'
import adminUsers from './zh-CN/admin/users.js'
import adminProviders from './zh-CN/admin/providers.js'
import adminConfig from './zh-CN/admin/config.js'
import adminAnnouncements from './zh-CN/admin/announcements.js'
import adminInviteCodes from './zh-CN/admin/inviteCodes.js'
import adminRedemptionCodes from './zh-CN/admin/redemptionCodes.js'
import adminVouchers from './zh-CN/admin/vouchers.js'
import adminSystemImages from './zh-CN/admin/systemImages.js'
import adminInstances from './zh-CN/admin/instances.js'
import adminTasks from './zh-CN/admin/tasks.js'
import adminTraffic from './zh-CN/admin/traffic.js'
import adminPortMapping from './zh-CN/admin/portMapping.js'
import adminOauth2 from './zh-CN/admin/oauth2.js'
import adminPerformance from './zh-CN/admin/performance.js'
import adminLogs from './zh-CN/admin/logs.js'
import adminBlockRules from './zh-CN/admin/blockRules.js'
import adminDomain from './zh-CN/admin/domain.js'
import adminKyc from './zh-CN/admin/kyc.js'
import adminGroup from './zh-CN/admin/group.js'
import adminApiTokens from './zh-CN/admin/apiTokens.js'
import adminSnapshots from './zh-CN/admin/snapshots.js'
import adminProducts from './zh-CN/admin/products.js'
import adminOrders from './zh-CN/admin/orders.js'
import adminTickets from './zh-CN/admin/tickets.js'
import adminSiteConfig from './zh-CN/admin/siteConfig.js'
import adminYipayConfig from './zh-CN/admin/yipayConfig.js'
import adminUpstream from './zh-CN/admin/upstream.js'

export default {
  common,
  navbar,
  login,
  adminLogin,
  register,
  forgotPassword,
  resetPassword,
  oauth2Callback,
  init,
  home,
  sidebar,
  user: {
    dashboard: userDashboard,
    profile: userProfile,
    instances: userInstances,
    instanceDetail: userInstanceDetail,
    tasks: userTasks,
    trafficOverview: userTrafficOverview,
    traffic: userTraffic,
    resources: userResources,
    apply: userApply,
    domain: userDomain,
    kyc: userKyc,
    checkin: userCheckin,
    apiTokens: userApiTokens,
    store: userStore,
    orders: userOrders,
    tickets: userTickets,
    wallet: userWallet
  },
  admin: {
    dashboard: adminDashboard,
    users: adminUsers,
    providers: adminProviders,
    config: adminConfig,
    announcements: adminAnnouncements,
    inviteCodes: adminInviteCodes,
    redemptionCodes: adminRedemptionCodes,
    vouchers: adminVouchers,
    systemImages: adminSystemImages,
    instances: adminInstances,
    tasks: adminTasks,
    traffic: adminTraffic,
    portMapping: adminPortMapping,
    oauth2: adminOauth2,
    performance: adminPerformance,
    logs: adminLogs,
    blockRules: adminBlockRules,
    domain: adminDomain,
    kyc: adminKyc,
    group: adminGroup,
    apiTokens: adminApiTokens,
    snapshots: adminSnapshots,
    products: adminProducts,
    orders: adminOrders,
    tickets: adminTickets,
    siteConfig: adminSiteConfig,
    yipayConfig: adminYipayConfig,
    upstream: adminUpstream
  },
  validation,
  message,
  errors,
  notFound
}
