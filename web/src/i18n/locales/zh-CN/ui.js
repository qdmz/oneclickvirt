// OneHost 设计系统新增文案（布局、认证品牌面板、首页、仪表盘、404 等）
export default {
  brand: {
    tagline: '新一代云主机与虚拟化管理平台',
    headline: '一站式管理你的云端算力',
    description: '统一纳管 LXD / Incus / Docker / Proxmox VE 等多种虚拟化节点，开通、计费、监控、运维全流程自动化。',
    point1: '多虚拟化平台统一纳管',
    point2: '秒级开通 · 自动化运维',
    point3: '流量 / 资源实时监控',
    point4: '内置商城、工单与钱包',
    secure: '数据传输全程加密'
  },
  layout: {
    console: '控制台',
    adminConsole: '管理后台',
    collapse: '收起侧边栏',
    expand: '展开侧边栏',
    fullscreen: '全屏',
    exitFullscreen: '退出全屏',
    backToSite: '访问首页',
    profile: '个人中心',
    signedInAs: '当前登录',
    adminBadge: '管理员',
    userBadge: '用户'
  },
  home: {
    badge: '企业级云主机管理平台',
    ctaConsole: '进入控制台',
    ctaStore: '浏览产品',
    featuresTitle: '为规模化运营而生',
    featuresSubtitle: '从节点接入到用户交付，覆盖云主机业务的每一个环节',
    feature1Title: '多平台统一纳管',
    feature1Desc: '一套面板对接 LXD、Incus、Docker、Proxmox VE、Podman 等，屏蔽底层差异。',
    feature2Title: '自动化开通交付',
    feature2Desc: '下单即开通，镜像、端口映射、IPv6、快照全自动处理。',
    feature3Title: '精细化资源计量',
    feature3Desc: '按用户等级配额限制 CPU、内存、磁盘与流量，超限自动处置。',
    feature4Title: '实时监控与告警',
    feature4Desc: '节点与实例流量、资源使用趋势一目了然，异常及时发现。',
    feature5Title: '商业化闭环',
    feature5Desc: '内置产品商城、订单、钱包、代金券与在线支付，开箱即可运营。',
    feature6Title: '安全与合规',
    feature6Desc: '实名认证、封禁规则、API Token 与 OAuth2 登录，权限分级清晰。',
    ctaTitle: '准备好开始了吗？',
    ctaSubtitle: '注册账号，几分钟内即可拥有你的第一台云主机。',
    from: '起',
    viewDetail: '查看详情',
    viewAll: '查看全部产品'
  },
  dashboard: {
    goodMorning: '早上好',
    goodAfternoon: '下午好',
    goodEvening: '晚上好',
    adminSubtitle: '这是平台当前的运行概况',
    instanceDistribution: '实例分布',
    resourceAllocation: '资源分配率',
    serverStatus: '面板服务器状态',
    uptime: '运行时长',
    running: '运行中',
    activeUsers: '活跃用户',
    activeProviders: '在线节点',
    totalInstances: '实例总数',
    quickEntry: '快捷入口',
    allocated: '已分配',
    available: '可用',
    noData: '暂无数据',
    quotaOverview: '配额使用概览',
    recentAnnouncements: '最新公告',
    noAnnouncements: '暂无公告'
  },
  notFound: {
    hint: '页面可能已被移动、删除，或您输入的地址有误。'
  },
  resetPassword: {
    verifying: '正在验证重置链接…',
    successTitle: '密码已重置',
    successDesc: '系统已为您生成新密码，并发送到您的邮箱 / 已绑定的通知渠道，请查收后尽快登录并修改密码。',
    invalidTitle: '链接无效或已过期',
    invalidDesc: '请重新申请找回密码，重置链接有效期为 24 小时。',
    missingToken: '缺少重置令牌，请通过邮件中的链接访问本页面。',
    requestAgain: '重新申请'
  },
  common: {
    viewAll: '查看全部',
    retry: '重试'
  }
}
