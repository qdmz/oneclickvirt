<!-- eslint-disable vue/no-v-html -->
<template>
  <div class="home-container">
    <!-- 导航栏 -->
    <header class="home-header">
      <div class="header-content">
        <div class="logo">
          <img
            :src="siteStore.logoSrc"
            :alt="siteStore.displaySiteName"
            class="logo-image"
          >
          <h1>{{ siteStore.displaySiteName }}</h1>
        </div>
        <nav class="nav-menu">
          <!-- 主题切换按钮 -->
          <button
            class="nav-link theme-btn"
            :title="themeStore.isDark ? t('navbar.lightMode') : t('navbar.darkMode')"
            @click="toggleTheme"
          >
            <el-icon><component :is="themeStore.isDark ? Sunny : Moon" /></el-icon>
          </button>
          <!-- 语言切换按钮 -->
          <button
            class="nav-link language-btn"
            @click="switchLanguage"
          >
            <el-icon><Operation /></el-icon>
            {{ languageStore.currentLanguage === 'zh-CN' ? 'English' : '中文' }}
          </button>
          <router-link
            to="/login"
            class="nav-link"
          >
            {{ t('home.nav.login') }}
          </router-link>
          <router-link
            to="/register"
            class="nav-link primary"
          >
            {{ t('home.nav.register') }}
          </router-link>
        </nav>
      </div>
    </header>

    <!-- 主要内容 -->
    <main class="home-main">
      <!-- 英雄区域 -->
      <section
        class="hero-section"
        :style="heroStyle"
      >
        <div class="hero-content">
          <h1 class="hero-title">
            {{ siteStore.homeTitle || t('home.hero.title') }}
          </h1>
          <p class="hero-description">
            {{ siteStore.homeSubtitle || t('home.hero.description') }}
          </p>
          <div class="hero-actions">
            <router-link
              to="/login"
              class="btn btn-primary"
            >
              {{ t('home.hero.loginButton') }}
            </router-link>
            <router-link
              to="/register"
              class="btn btn-secondary"
            >
              {{ t('home.hero.registerButton') }}
            </router-link>
          </div>
        </div>
        <div class="hero-image">
          <div class="feature-preview">
            <div class="preview-card">
              <div class="card-icon">
                <i class="fas fa-server" />
              </div>
              <h3>{{ t('home.features.vm.title') }}</h3>
              <p>{{ t('home.features.vm.description') }}</p>
            </div>
            <div class="preview-card">
              <div class="card-icon">
                <i class="fas fa-box" />
              </div>
              <h3>{{ t('home.features.container.title') }}</h3>
              <p>{{ t('home.features.container.description') }}</p>
            </div>
            <div class="preview-card">
              <div class="card-icon">
                <i class="fas fa-chart-bar" />
              </div>
              <h3>{{ t('home.features.monitoring.title') }}</h3>
              <p>{{ t('home.features.monitoring.description') }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 平台概览 -->
      <section class="overview-section">
        <div class="section-header">
          <h2>{{ t('home.platformOverview.title') }}</h2>
          <p>{{ t('home.platformOverview.description') }}</p>
        </div>
        <div
          class="stats-grid"
          aria-label="platform-stats"
        >
          <div class="platform-item stats-item">
            <div class="platform-icon">
              <i
                class="fas fa-users fa-2x"
                aria-hidden="true"
              />
            </div>
            <h3>{{ t('home.stats.users') }}</h3>
            <p class="stats-value">
              {{ usersCountDisplay }}
            </p>
          </div>

          <div class="platform-item stats-item">
            <div class="platform-icon">
              <i
                class="fas fa-network-wired fa-2x"
                aria-hidden="true"
              />
            </div>
            <h3>{{ t('home.stats.nodes') }}</h3>
            <p class="stats-value">
              {{ nodesCountDisplay }}
            </p>
          </div>

          <div class="platform-item stats-item">
            <div class="platform-icon">
              <i
                class="fas fa-box fa-2x"
                aria-hidden="true"
              />
            </div>
            <h3>{{ t('home.stats.containers') }}</h3>
            <p class="stats-value">
              {{ containersCountDisplay }}
            </p>
          </div>

          <div class="platform-item stats-item">
            <div class="platform-icon">
              <i
                class="fas fa-server fa-2x"
                aria-hidden="true"
              />
            </div>
            <h3>{{ t('home.stats.vms') }}</h3>
            <p class="stats-value">
              {{ vmsCountDisplay }}
            </p>
          </div>
        </div>
      </section>

      <!-- 支持的虚拟化平台 -->
      <section
        v-if="siteStore.showPlatforms"
        class="platforms-section"
      >
        <div class="section-header">
          <h2>{{ t('home.platforms.title') }}</h2>
          <p>{{ t('home.platforms.description') }}</p>
        </div>
        <LogoCarousel
          :items="platformList"
          :speed="35"
          direction="left"
          :gap="24"
        >
          <template #default="{ item }">
            <a
              :href="item.href"
              target="_blank"
              rel="noopener noreferrer"
              class="carousel-platform-card"
              :title="item.name"
            >
              <div class="platform-card-icon">
                <img
                  :src="item.icon"
                  :alt="item.name"
                  width="48"
                  height="48"
                  loading="lazy"
                >
              </div>
              <span class="platform-card-name">{{ item.name }}</span>
              <span class="platform-card-repo">
                <svg
                  width="12"
                  height="12"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                ><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" /></svg>
                {{ item.repo }}
              </span>
            </a>
          </template>
        </LogoCarousel>
      </section>

      <!-- 赞助方 -->
      <section
        v-if="siteStore.showSponsors"
        class="supporters-section"
      >
        <div class="section-header">
          <h2>{{ t('home.supporters.title') }}</h2>
          <p>{{ t('home.supporters.description') }}</p>
        </div>
        <div class="supporters-grid">
          <a
            v-for="item in sponsorList"
            :key="item.name"
            :href="item.href"
            target="_blank"
            rel="noopener noreferrer"
            :class="['supporter-card', item.cardClass]"
            :title="item.name"
            :aria-label="item.name"
          >
            <img
              :src="item.logo"
              :alt="item.name"
              loading="eager"
              decoding="async"
              :class="item.logoClass"
            >
          </a>
        </div>
      </section>

      <!-- 系统公告 -->
      <section
        v-if="announcements.length > 0"
        class="announcements-section"
      >
        <div class="section-header">
          <h2>{{ t('home.announcements.title') }}</h2>
        </div>
        <div class="announcements-list">
          <div
            v-for="announcement in announcements"
            :key="announcement.id"
            class="announcement-item"
          >
            <div class="announcement-header">
              <h3>{{ announcement.title }}</h3>
              <div class="announcement-meta">
                <el-tag
                  :type="announcement.type === 'homepage' ? 'success' : 'warning'"
                  size="small"
                >
                  {{ announcement.type === 'homepage' ? t('home.announcements.typeHomepage') : t('home.announcements.typeTopbar') }}
                </el-tag>
                <span class="announcement-date">{{ formatDate(announcement.createdAt) }}</span>
              </div>
            </div>
            <div
              class="announcement-content"
              v-html="announcement.contentHtml || announcement.content"
            />
          </div>
        </div>
      </section>

      <!-- 推荐产品 -->
      <section
        v-if="siteStore.showRecommended && recommendedProducts.length > 0"
        class="recommended-section"
      >
        <div class="section-header">
          <h2>{{ siteStore.recommendedTitle || t('home.recommendedProducts.title') }}</h2>
          <p>{{ siteStore.recommendedSubtitle || t('home.recommendedProducts.subtitle') }}</p>
        </div>
        <div
          class="recommended-grid"
          :style="{ gridTemplateColumns: 'repeat(' + effectiveCols + ', minmax(0, 1fr))' }"
        >
          <el-card
            v-for="p in recommendedProducts"
            :key="p.id"
            shadow="hover"
            class="recommended-card"
            @click="goToProduct(p.id)"
          >
            <div class="rec-card-head">
              <div class="rec-card-icon">
                <img
                  v-if="p.icon"
                  :src="p.icon"
                  :alt="p.name"
                  class="rec-icon-img"
                >
                <el-icon v-else :size="30" color="#16a34a"><Box /></el-icon>
              </div>
              <div class="rec-card-title">
                <h3 class="rec-name">{{ p.name }}</h3>
                <span class="rec-type-tag">{{ categoryLabel(p) }}</span>
              </div>
            </div>
            <p class="rec-desc">{{ p.description }}</p>

            <div
              v-if="siteStore.recommendedShowSpecs"
              class="rec-specs"
            >
              <div class="rec-spec">
                <span class="rec-spec-label">{{ t('user.store.cpu') }}</span>
                <span class="rec-spec-val">{{ p.cpu }} {{ t('user.store.cores') }}</span>
              </div>
              <div class="rec-spec">
                <span class="rec-spec-label">{{ t('user.store.memory') }}</span>
                <span class="rec-spec-val">{{ formatMB(p.memory) }}</span>
              </div>
              <div class="rec-spec">
                <span class="rec-spec-label">{{ t('user.store.disk') }}</span>
                <span class="rec-spec-val">{{ formatMB(p.disk) }}</span>
              </div>
              <div class="rec-spec">
                <span class="rec-spec-label">{{ t('user.store.bandwidth') }}</span>
                <span class="rec-spec-val">{{ p.bandwidth }} Mbps</span>
              </div>
              <div class="rec-spec">
                <span class="rec-spec-label">{{ t('user.store.traffic') }}</span>
                <span class="rec-spec-val">{{ formatMB(p.traffic) }}</span>
              </div>
            </div>

            <div class="rec-footer">
              <div
                v-if="siteStore.recommendedShowPrice"
                class="rec-price"
              >
                <span class="rec-price-symbol">¥</span>
                <span class="rec-price-num">{{ p.price }}</span>
                <span class="rec-price-unit">/{{ periodLabel(p) }}</span>
              </div>
              <div class="rec-stock">
                {{ p.stock < 0 ? t('user.store.stockUnlimited') : (p.stock + ' ' + t('user.store.stock')) }}
              </div>
            </div>
          </el-card>
        </div>
        <div class="rec-viewall">
          <router-link
            to="/user/store"
            class="btn btn-secondary"
          >{{ t('home.recommendedProducts.viewAll') }}</router-link>
        </div>
      </section>
    </main>

    <!-- 页脚 -->
    <footer class="home-footer">
      <div class="footer-glow-top" />
      <div class="footer-inner">
        <div class="footer-brand">
          <div class="footer-logo">
            <img
              :src="siteStore.logoSrc"
              alt="OneClickVirt Logo"
              class="footer-logo-img"
            >
            <span class="footer-logo-text">{{ siteStore.displaySiteName }}</span>
          </div>
          <p class="footer-tagline">
            {{ t('home.hero.description') }}
          </p>
          <a
            href="https://github.com/oneclickvirt/oneclickvirt"
            target="_blank"
            rel="noopener noreferrer"
            class="footer-github-btn"
          >
            <svg
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="currentColor"
            >
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
            </svg>
            GitHub
          </a>
        </div>

        <div
          v-if="footerLinkList.length > 0"
          class="footer-links-grid"
        >
          <div class="footer-col">
            <h4 class="footer-col-title">
              <span class="footer-col-dot" />
              {{ t('home.footer.friendLinks') }}
            </h4>
            <ul class="footer-link-list">
              <li
                v-for="l in footerLinkList"
                :key="l.id"
              >
                <a
                  :href="l.href"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <span class="link-arrow">›</span>{{ l.name }}
                </a>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <div class="footer-bottom">
        <div class="footer-bottom-inner">
          <span class="footer-copyright">&copy; 2026 {{ siteStore.copyrightText || (siteStore.displaySiteName + '. ' + t('home.footer.allRightsReserved')) }}</span>
          <span class="footer-divider" />
          <a
            href="https://github.com/oneclickvirt"
            target="_blank"
            rel="noopener noreferrer"
            class="footer-bottom-link"
          >
            <svg
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="currentColor"
              style="margin-right:4px;vertical-align:middle"
            >
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
            </svg>
            {{ t('home.footer.openSourceProject') }}
          </a>
          <template v-if="serverVersion">
            <span class="footer-divider" />
            <span
              class="footer-version-tag"
              :title="`${t('home.footer.serverVersion')} ${serverVersion}`"
            >
              <span>{{ t('home.footer.serverVersion') }}</span>
              <span class="footer-version-value">{{ serverVersion }}</span>
            </span>
            <a
              v-if="updateAvailable && latestVersion"
              :href="releaseUrl || 'https://github.com/oneclickvirt/oneclickvirt/releases'"
              target="_blank"
              rel="noopener noreferrer"
              class="footer-bottom-link footer-version-update"
              :title="`${t('home.footer.latestVersion')} ${latestVersion}`"
            >
              <span>{{ t('home.footer.latestVersion') }}</span>
              <span class="footer-version-value">{{ latestVersion }}</span>
            </a>
          </template>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getPublicAnnouncements, getPublicStats, getServerVersion, getRecommendedProducts, getSiteLinks } from '@/api/public'
import { checkSystemInit } from '@/api/init'
import { ElTag, ElMessage } from 'element-plus'
import { Operation, Sunny, Moon, Box } from '@element-plus/icons-vue'
import { useLanguageStore } from '@/pinia/modules/language'
import LogoCarousel from '@/components/LogoCarousel.vue'
import proxmoxPng from '@/assets/images/proxmox.png'
import incusPng from '@/assets/images/incus.png'
import dockerPng from '@/assets/images/docker.png'
import lxdPng from '@/assets/images/lxd.png'
import podmanSvg from '@/assets/images/podman.svg'
import containerdSvg from '@/assets/images/containerd.svg'
import qemuSvg from '@/assets/images/qemu.svg'
import kubevirtPng from '@/assets/images/KubeVirt.png'
import ibmLinuxonePng from '@/assets/images/ibm-linuxone.png'
import dartnodePng from '@/assets/images/dartnode.png'
import { useThemeStore } from '@/pinia/modules/theme'
import { useSiteStore } from '@/pinia/modules/site'

const router = useRouter()
const { t, locale } = useI18n()
const languageStore = useLanguageStore()
const themeStore = useThemeStore()
const siteStore = useSiteStore()
const announcements = ref([])
// 首页推荐产品
const recommendedProducts = ref([])
// 统计数据
const usersCount = ref(null)
const nodesCount = ref(null)
const containersCount = ref(null)
const vmsCount = ref(null)
const serverVersion = ref('')
const latestVersion = ref('')
const releaseUrl = ref('')
const updateAvailable = ref(false)
const platforms = [
  { name: 'Proxmox VE', icon: proxmoxPng, href: 'https://github.com/oneclickvirt/pve', repo: 'oneclickvirt/pve' },
  { name: 'Incus', icon: incusPng, href: 'https://github.com/oneclickvirt/incus', repo: 'oneclickvirt/incus' },
  { name: 'Docker', icon: dockerPng, href: 'https://github.com/oneclickvirt/docker', repo: 'oneclickvirt/docker' },
  { name: 'LXD', icon: lxdPng, href: 'https://github.com/oneclickvirt/lxd', repo: 'oneclickvirt/lxd' },
  { name: 'Podman', icon: podmanSvg, href: 'https://github.com/oneclickvirt/podman', repo: 'oneclickvirt/podman' },
  { name: 'Containerd', icon: containerdSvg, href: 'https://github.com/oneclickvirt/containerd', repo: 'oneclickvirt/containerd' },
  { name: 'QEMU', icon: qemuSvg, href: 'https://github.com/oneclickvirt/qemu', repo: 'oneclickvirt/qemu' },
  { name: 'KubeVirt', icon: kubevirtPng, href: 'https://github.com/oneclickvirt/kubevirt', repo: 'oneclickvirt/kubevirt' }
]

const footerSponsors = [
  {
    name: 'DartNode',
    href: 'https://dartnode.com?aff=bonus',
    logo: dartnodePng,
    cardClass: 'supporter-card-dartnode',
    logoClass: 'supporter-logo-dartnode'
  },
  {
    name: 'zmto',
    href: 'https://console.zmto.com/?affid=1524',
    logo: 'https://console.zmto.com/templates/2019/dist/images/logo_dark.svg'
  },
  {
    name: 'IBM LinuxONE OSS Community Cloud',
    href: 'https://community.ibm.com/zsystems/form/l1cc-oss-vm-request/',
    logo: ibmLinuxonePng,
    cardClass: 'supporter-card-ibm',
    logoClass: 'supporter-logo-ibm'
  },
  {
    name: 'fossvps',
    href: 'https://fossvps.org/',
    logo: 'https://lowendspirit.com/uploads/userpics/793/nHSR7IOVIBO84.png'
  },
  {
    name: 'Linux DO',
    href: 'https://linux.do/',
    logo: 'https://cdn3.ldstatic.com/original/4X/d/1/4/d146c68151340881c884d95e0da4acdf369258c6.png',
    cardClass: 'supporter-card-linuxdo',
    logoClass: 'supporter-logo-linuxdo'
  }
]

const usersCountDisplay = computed(() => (usersCount.value === null ? '-' : usersCount.value))
const nodesCountDisplay = computed(() => (nodesCount.value === null ? '-' : nodesCount.value))
const containersCountDisplay = computed(() => (containersCount.value === null ? '-' : containersCount.value))
const vmsCountDisplay = computed(() => (vmsCount.value === null ? '-' : vmsCount.value))

// 站点链接（来自 site_links 表，后台「站点链接」可配置）——首页平台/赞助方/页脚均从此读取
const siteLinksPlatform = ref([])
const siteLinksSponsor = ref([])
const footerLinks = ref([])

const fetchSiteLinks = async () => {
  try {
    const [p, s, f] = await Promise.all([
      getSiteLinks('platform'),
      getSiteLinks('sponsor'),
      getSiteLinks('footer')
    ])
    if (p && p.code === 200 && Array.isArray(p.data)) siteLinksPlatform.value = p.data
    if (s && s.code === 200 && Array.isArray(s.data)) siteLinksSponsor.value = s.data
    if (f && f.code === 200 && Array.isArray(f.data)) footerLinks.value = f.data
  } catch (error) {
    console.error('获取站点链接失败', error)
  }
}

// 平台轮播：优先使用后台配置的 site_links，无数据时回退到内置默认值
const platformList = computed(() => {
  if (siteLinksPlatform.value.length > 0) {
    return siteLinksPlatform.value.map(l => ({ name: l.name, href: l.url, icon: l.iconUrl, repo: l.description || '' }))
  }
  return platforms
})

// 赞助方：优先后台配置，无数据回退内置默认值
const sponsorList = computed(() => {
  if (siteLinksSponsor.value.length > 0) {
    return siteLinksSponsor.value.map(l => ({ name: l.name, href: l.url, logo: l.iconUrl }))
  }
  return footerSponsors
})

// 页脚友链
const footerLinkList = computed(() => footerLinks.value.map(l => ({ id: l.id, name: l.name, href: l.url })))

// 英雄区背景（站点配置可设置首页背景图）
const heroStyle = computed(() => {
  return siteStore.homeBackground
    ? { backgroundImage: `url(${siteStore.homeBackground})`, backgroundSize: 'cover', backgroundPosition: 'center' }
    : {}
})

const switchLanguage = () => {
  const newLang = languageStore.toggleLanguage()
  locale.value = newLang
  ElMessage.success(t('navbar.languageSwitched'))
}

const toggleTheme = () => {
  themeStore.toggleTheme()
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString(locale.value === 'zh-CN' ? 'zh-CN' : 'en-US')
}

const fetchAnnouncements = async () => {
  try {
    // 获取首页公告
    const response = await getPublicAnnouncements('homepage')
    if (response.code === 200) {
      announcements.value = response.data.slice(0, 3) // 只显示最新3条
    }
  } catch (error) {
    console.error(t('home.errors.fetchAnnouncementsFailed'), error)
  }
}

const fetchPublicStats = async () => {
  try {
    const resp = await getPublicStats()
    if (resp && (resp.code === 200) && resp.data) {
      const d = resp.data
      // 尝试从常见字段拾取数据，做多层回退以兼容不同返回结构
      usersCount.value = d.userStats?.totalUsers ?? d.user_count ?? d.userCount ?? d.userTotal ?? null
      // nodes 可能对应 regionStats 的 count 总和或 provider 总数
      if (Array.isArray(d.regionStats) && d.regionStats.length > 0) {
        let total = 0
        d.regionStats.forEach(r => { total += r.count ?? 0 })
        nodesCount.value = total
      } else {
        nodesCount.value = d.provider_count ?? d.node_count ?? d.nodeCount ?? null
      }

      // 容器/虚拟机：尝试从资源统计中读取
      containersCount.value = d.resourceUsage?.container_count ?? d.resourceUsage?.containerCount ?? d.container_count ?? d.containerCount ?? null
      vmsCount.value = d.resourceUsage?.vm_count ?? d.resourceUsage?.vmCount ?? d.vm_count ?? d.vmCount ?? null
    }
  } catch (error) {
    console.error('获取公开统计数据失败', error)
  }
}

const fetchRecommendedProducts = async () => {
  try {
    const limit = siteStore.recommendedLimit > 0 ? siteStore.recommendedLimit : 8
    const res = await getRecommendedProducts({ limit })
    if (res && res.code === 200 && Array.isArray(res.data)) {
      recommendedProducts.value = res.data.slice(0, limit)
    }
  } catch (error) {
    console.error('获取推荐产品失败', error)
  }
}

const goToProduct = (id) => {
  router.push(`/user/store/${id}`)
}

const goToStore = () => {
  router.push('/user/store')
}

// 推荐产品每行列数（来自站点配置，限制在 2~6 之间）
const effectiveCols = computed(() => {
  const c = Number(siteStore.recommendedCols)
  if (!c || c < 2) return 2
  if (c > 6) return 6
  return c
})

// 产品类型/类别标签
const categoryLabel = (p) => {
  if (!p) return ''
  if (p.category === 'container') return t('user.store.categoryContainer')
  if (p.category === 'vm') return t('user.store.categoryVM')
  return p.type || ''
}

// 将 MB 格式化为 GB / MB
const formatMB = (mb) => {
  const n = Number(mb) || 0
  if (n >= 1024) return (n / 1024) + ' GB'
  return n + ' MB'
}

// 计费周期后缀
const periodLabel = (p) => {
  const map = {
    month: t('user.store.perMonth'),
    day: t('user.store.perDay'),
    year: t('user.store.perYear'),
    hour: t('user.store.perHour')
  }
  return map[p && p.periodType] || t('user.store.perMonth')
}

const checkInitStatus = async () => {
  try {
    const response = await checkSystemInit()
    if (response && (response.code === 200) && response.data && response.data.needInit === true) {
      router.push('/init')
    }
  } catch (error) {
    console.error(t('home.errors.checkInitFailed'), error)
    // 如果是网络错误或服务器错误，可能是数据库未初始化导致的
    if (error.message.includes('Network Error') || 
        error.response?.status >= 500 || 
        error.code === 'ECONNREFUSED') {
      console.warn(t('home.debug.serverConnectionFailed'))
      router.push('/init')
    }
  }
}

onMounted(() => {
  // 首先加载完整站点配置（首页各栏目/英雄区文本/显示开关均由站点配置驱动）
  siteStore.fetchFullSiteConfig()
  // 获取后台可配置的站点链接（平台/赞助方/页脚）
  fetchSiteLinks()
  // 首先检查初始化状态
  checkInitStatus()
  // 然后获取公告
  fetchAnnouncements()
  // 获取公开统计数据（用于未登录首页展示）
  fetchPublicStats()
  // 获取首页推荐产品
  fetchRecommendedProducts()
  // 获取服务器版本信息
  getServerVersion().then(res => {
    if (res && (res.code === 200) && res.data?.server_version) {
      serverVersion.value = res.data.server_version
      latestVersion.value = res.data.latest_version || ''
      releaseUrl.value = res.data.release_url || ''
      updateAvailable.value = Boolean(res.data.update_available)
    }
  }).catch(() => {})
})
</script>

<style src="./home.css" scoped>

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
