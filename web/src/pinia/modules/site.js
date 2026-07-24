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
  // 是否已经初始化
  const initialized = ref(false)

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
        icp.value = data.icp || ''
        faviconURL.value = data.favicon_url || ''
        primaryColor.value = data.primary_color || '#2563eb'
        customCSS.value = data.custom_css || ''
        customHeader.value = data.custom_header || ''
        customFooter.value = data.custom_footer || ''
        enableRegister.value = data.enable_register !== false
        enableStore.value = data.enable_store !== false
        enableRecharge.value = data.enable_recharge !== false
        enableTicket.value = data.enable_ticket !== false
        enableAnnouncement.value = data.enable_announcement !== false
        maintenanceMode.value = data.maintenance_mode === true
        maintenanceMessage.value = data.maintenance_message || ''
        // 应用主题色
        applyPrimaryColor(primaryColor.value)
        // 应用自定义 CSS
        applyCustomCSS(customCSS.value)
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
    fetchSiteConfig,
    fetchFullSiteConfig,
    applyPrimaryColor,
    applyCustomCSS,
    refresh
  }
})
