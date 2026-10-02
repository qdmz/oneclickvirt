import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getPublicSystemConfig } from '@/api/public'
import { getSiteConfig } from '@/api/site'

export const useSiteStore = defineStore('site', () => {
  // 自定义 Logo URL（从后端读取，空字符串表示使用默认）
  const logoURL = ref('')
  // 自定义网站名称（空表示使用默认 OneClickVirt）
  const siteName = ref('')
  // 站点描述
  const siteDescription = ref('')
  // 站点关键词
  const siteKeywords = ref('')
  // 联系邮箱
  const contactEmail = ref('')
  // ICP备案号
  const icp = ref('')
  // Favicon URL
  const faviconURL = ref('')
  // 主题色
  const primaryColor = ref('#2563eb')
  // 自定义 CSS
  const customCSS = ref('')
  // 自定义页眉 HTML
  const customHeader = ref('')
  // 自定义页脚 HTML
  const customFooter = ref('')
  // 功能开关
  const enableRegister = ref(true)
  const enableStore = ref(true)
  const enableRecharge = ref(true)
  const enableTicket = ref(true)
  const enableAnnouncement = ref(true)
  // 维护模式
  const maintenanceMode = ref(false)
  const maintenanceMessage = ref('')
  // 易支付是否启用（后台“易支付配置”的启用开关）
  const showYiPay = ref(true)
  // 易支付启用的支付方式（后台“启用的支付方式”：alipay / wxpay / qqpay）
  // 前端下单/充值时只展示此处列出的渠道，关闭的渠道不出现供选择
  const enabledPayTypes = ref(['alipay', 'wxpay', 'qqpay'])
  // 是否已经初始化
  const initialized = ref(false)

  // ===== 首页与站点展示配置（来自完整站点配置） =====
  // 首页英雄区
  const homeTitle = ref('')
  const homeSubtitle = ref('')
  const homeBackground = ref('')
  const showHomeStats = ref(true)
  // 首页各显示栏目开关
  const showPlatforms = ref(true)
  const showSponsors = ref(true)
  const showRecommended = ref(true)
  const recommendedLimit = ref(8)
  // 推荐产品详细配置
  const recommendedTitle = ref('')
  const recommendedSubtitle = ref('')
  const recommendedCols = ref(4)
  const recommendedShowPrice = ref(true)
  const recommendedShowSpecs = ref(true)
  // 其它展示项
  const copyrightText = ref('')
  const darkLogoURL = ref('')
  const contactPhone = ref('')
  const contactQQ = ref('')
  const contactTelegram = ref('')
  const showBalance = ref(true)
  const showNav = ref(true)
  const headerEnabled = ref(false)
  const footerEnabled = ref(false)
  const announcementBar = ref('')

  // 默认 Logo 资源路径（用于 img src 属性）
  const defaultLogoSrc = new URL('@/assets/images/logo.png', import.meta.url).href

  // 计算最终使用的 logo src
  const logoSrc = computed(() => {
    return logoURL.value && logoURL.value.trim() !== '' ? logoURL.value.trim() : defaultLogoSrc
  })

  // 计算最终显示的网站名称（空时默认 OneClickVirt）
  const displaySiteName = computed(() => {
    return siteName.value && siteName.value.trim() !== '' ? siteName.value.trim() : 'OneClickVirt'
  })

  // 计算是否显示自定义页眉
  const hasCustomHeader = computed(() => {
    return customHeader.value && customHeader.value.trim().length > 0
  })

  // 计算是否显示自定义页脚
  const hasCustomFooter = computed(() => {
    return customFooter.value && customFooter.value.trim().length > 0
  })

  // 从后端获取站点配置（公开接口）
  async function fetchSiteConfig() {
    if (initialized.value) return
    try {
      const res = await getPublicSystemConfig()
      if (res && (res.code === 200) && res.data) {
        if (res.data.logo_url) {
          logoURL.value = res.data.logo_url
        }
        if (res.data.site_name) {
          siteName.value = res.data.site_name
        }
      }
    } catch (e) {
      // 静默失败，使用默认配置
    } finally {
      initialized.value = true
    }
  }

  // 获取完整的站点配置（包含功能开关等）
  async function fetchFullSiteConfig() {
    try {
      const res = await getSiteConfig()
      if (res && (res.code === 200) && res.data) {
        const data = res.data
        logoURL.value = data.logo_url || logoURL.value
        siteName.value = data.site_name || siteName.value
        siteDescription.value = data.site_description || ''
        siteKeywords.value = data.site_keywords || ''
        contactEmail.value = data.contact_email || ''
        icp.value = data.icp_number || ''
        faviconURL.value = data.favicon_url || ''
        darkLogoURL.value = data.dark_logo_url || ''
        primaryColor.value = data.primary_color || '#2563eb'
        customCSS.value = data.custom_css || ''
        customHeader.value = data.custom_header || ''
        customFooter.value = data.custom_footer || ''
        // 功能开关（字段名须与后端 GetPublicSiteInfo 返回保持一致）
        enableRegister.value = data.enable_registration !== false
        enableStore.value = data.enable_product_store !== false
        enableTicket.value = data.enable_ticket !== false
        enableAnnouncement.value = data.announcement_enabled === true
        showBalance.value = data.show_balance !== false
        showYiPay.value = data.show_yipay !== false
        showNav.value = data.show_nav !== false
        headerEnabled.value = data.header_enabled === true
        footerEnabled.value = data.footer_enabled === true
        maintenanceMode.value = data.maintenance_mode === true
        maintenanceMessage.value = data.maintenance_message || ''
        // 首页配置
        homeTitle.value = data.home_title || ''
        homeSubtitle.value = data.home_subtitle || ''
        homeBackground.value = data.home_background || ''
        showHomeStats.value = data.show_home_stats !== false
        showPlatforms.value = data.show_platforms !== false
        showSponsors.value = data.show_sponsors !== false
        showRecommended.value = data.show_recommended !== false
        recommendedLimit.value = Number(data.recommended_limit) > 0 ? Number(data.recommended_limit) : 8
        recommendedTitle.value = data.recommended_title || ''
        recommendedSubtitle.value = data.recommended_subtitle || ''
        recommendedCols.value = Number(data.recommended_cols) > 0 ? Number(data.recommended_cols) : 4
        recommendedShowPrice.value = data.recommended_show_price !== false
        recommendedShowSpecs.value = data.recommended_show_specs !== false
        copyrightText.value = data.copyright_text || ''
        contactPhone.value = data.contact_phone || ''
        contactQQ.value = data.contact_qq || ''
        contactTelegram.value = data.contact_telegram || ''
        announcementBar.value = data.announcement_bar || ''
        // 启用的支付方式（来自后台“启用的支付方式”）
        if (data.yipay_pay_types) {
          const arr = String(data.yipay_pay_types).split(',').map(s => s.trim()).filter(Boolean)
          if (arr.length > 0) {
            enabledPayTypes.value = arr
          }
        }
        // 应用主题色
        applyPrimaryColor(primaryColor.value)
        // 应用自定义 CSS
        applyCustomCSS(customCSS.value)
        // 应用站点图标
        applyFavicon(faviconURL.value)
        // 站点名称加载晚于路由守卫时，同步刷新一次浏览器标题
        applyDocumentTitle()
      }
    } catch (e) {
      // 静默失败，使用默认配置
    } finally {
      initialized.value = true
    }
  }

  // 应用主题色
  function applyPrimaryColor(color) {
    if (!color || typeof document === 'undefined') return
    const root = document.documentElement
    root.style.setProperty('--primary-color', color)
    // 同时更新 Element Plus 的主题色
    const style = document.getElementById('site-custom-theme-style')
    if (style) {
      style.innerHTML = generateThemeCSS(color)
    } else {
      const newStyle = document.createElement('style')
      newStyle.id = 'site-custom-theme-style'
      newStyle.innerHTML = generateThemeCSS(color)
      document.head.appendChild(newStyle)
    }
  }

  // 生成主题色 CSS
  function generateThemeCSS(color) {
    return `
      :root {
        --el-color-primary: ${color} !important;
      }
      .el-button--primary {
        --el-button-bg-color: ${color} !important;
        --el-button-border-color: ${color} !important;
      }
      .el-switch.is-checked .el-switch__core {
        --el-switch-on-color: ${color} !important;
      }
      .el-radio__input.is-checked .el-radio__inner {
        --el-radio-checked-bg-color: ${color} !important;
        --el-radio-checked-input-border-color: ${color} !important;
      }
      .el-checkbox__input.is-checked .el-checkbox__inner {
        --el-checkbox-checked-bg-color: ${color} !important;
        --el-checkbox-checked-input-border-color: ${color} !important;
      }
      .el-tabs__item.is-active {
        color: ${color} !important;
      }
      .el-tabs__active-bar {
        background-color: ${color} !important;
      }
    `
  }

  // 应用站点图标（favicon）
  function applyFavicon(url) {
    if (!url || typeof document === 'undefined') return
    let link = document.querySelector("link[rel~='icon']")
    if (!link) {
      link = document.createElement('link')
      link.rel = 'icon'
      document.head.appendChild(link)
    }
    link.href = url
  }

  // 用站点名称刷新浏览器标题（保留当前页面名部分）
  function applyDocumentTitle() {
    if (typeof document === 'undefined') return
    const brand = displaySiteName.value
    const cur = document.title || ''
    const idx = cur.lastIndexOf(' - ')
    document.title = idx > 0 ? `${cur.slice(0, idx)} - ${brand}` : brand
  }

  // 应用自定义 CSS
  function applyCustomCSS(css) {
    if (typeof document === 'undefined') return
    let style = document.getElementById('site-custom-css')
    if (!style) {
      style = document.createElement('style')
      style.id = 'site-custom-css'
      document.head.appendChild(style)
    }
    style.innerHTML = css || ''
  }

  // 强制刷新（管理员保存配置后调用）
  function refresh() {
    initialized.value = false
    return fetchFullSiteConfig()
  }

  return {
    logoURL,
    logoSrc,
    siteName,
    displaySiteName,
    siteDescription,
    siteKeywords,
    contactEmail,
    icp,
    faviconURL,
    primaryColor,
    customCSS,
    customHeader,
    customFooter,
    hasCustomHeader,
    hasCustomFooter,
    enableRegister,
    enableStore,
    enableRecharge,
    enableTicket,
    enableAnnouncement,
    maintenanceMode,
    maintenanceMessage,
    showYiPay,
    enabledPayTypes,
    homeTitle,
    homeSubtitle,
    homeBackground,
    showHomeStats,
    showPlatforms,
    showSponsors,
    showRecommended,
    recommendedLimit,
    recommendedTitle,
    recommendedSubtitle,
    recommendedCols,
    recommendedShowPrice,
    recommendedShowSpecs,
    copyrightText,
    darkLogoURL,
    contactPhone,
    contactQQ,
    contactTelegram,
    showBalance,
    showNav,
    headerEnabled,
    footerEnabled,
    announcementBar,
    fetchSiteConfig,
    fetchFullSiteConfig,
    applyPrimaryColor,
    applyCustomCSS,
    refresh
  }
})
