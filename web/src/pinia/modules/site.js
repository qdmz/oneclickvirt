import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getPublicSystemConfig } from '@/api/public'

export const useSiteStore = defineStore('site', () => {
  // 自定义 Logo URL（从后端读取，空字符串表示使用默认）
  const logoURL = ref('')
  // 自定义网站名称（空表示使用默认 OneClickVirt）
  const siteName = ref('')
  // 自定义站点标语（首页副标题）
  const slogan = ref('')
  // 自定义页脚文字
  const footerText = ref('')
  // 备案号文字
  const icpText = ref('')
  // 备案号跳转链接
  const icpLink = ref('')
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

  // 计算最终显示的标语（空时使用默认）
  const displaySlogan = computed(() => {
    return slogan.value && slogan.value.trim() !== ''
      ? slogan.value.trim()
      : '一站式虚拟化资源管理与销售平台'
  })

  // 从后端获取站点配置
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
        if (res.data.slogan !== undefined && res.data.slogan !== null) {
          slogan.value = res.data.slogan
        }
        if (res.data.footer_text !== undefined && res.data.footer_text !== null) {
          footerText.value = res.data.footer_text
        }
        if (res.data.icp_text !== undefined && res.data.icp_text !== null) {
          icpText.value = res.data.icp_text
        }
        if (res.data.icp_link !== undefined && res.data.icp_link !== null) {
          icpLink.value = res.data.icp_link
        }
      }
    } catch (e) {
      // 静默失败，使用默认配置
    } finally {
      initialized.value = true
    }
  }

  // 强制刷新（管理员保存配置后调用）
  function refresh() {
    initialized.value = false
    return fetchSiteConfig()
  }

  return {
    logoURL, logoSrc, siteName, displaySiteName,
    slogan, displaySlogan, footerText, icpText, icpLink,
    fetchSiteConfig, refresh
  }
})
