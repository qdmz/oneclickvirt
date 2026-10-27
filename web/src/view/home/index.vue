<!-- eslint-disable vue/no-v-html -->
<template>
  <div class="home">
    <!-- 顶部导航（玻璃态，吸顶） -->
    <header class="nav">
      <div class="nav-inner">
        <router-link
          to="/"
          class="brand"
        >
          <img
            :src="siteStore.logoSrc"
            :alt="siteStore.displaySiteName"
            class="brand-logo"
          >
          <span class="brand-name">{{ siteStore.displaySiteName }}</span>
        </router-link>
        <nav class="nav-actions">
          <button
            class="icon-btn"
            :title="themeStore.isDark ? t('navbar.lightMode') : t('navbar.darkMode')"
            :aria-label="themeStore.isDark ? t('navbar.lightMode') : t('navbar.darkMode')"
            @click="toggleTheme"
          >
            <el-icon><component :is="themeStore.isDark ? Sunny : Moon" /></el-icon>
          </button>
          <button
            class="icon-btn lang-btn"
            @click="switchLanguage"
          >
            <el-icon><Operation /></el-icon>
            <span>{{ languageStore.currentLanguage === 'zh-CN' ? 'EN' : '中文' }}</span>
          </button>
          <template v-if="userStore.isLoggedIn">
            <router-link
              :to="consolePath"
              class="btn btn-primary btn-sm"
            >
              {{ t('ui.home.ctaConsole') }}
            </router-link>
          </template>
          <template v-else>
            <router-link
              to="/login"
              class="btn btn-ghost btn-sm hide-xs"
            >
              {{ t('home.nav.login') }}
            </router-link>
            <router-link
              to="/register"
              class="btn btn-primary btn-sm"
            >
              {{ t('home.nav.register') }}
            </router-link>
          </template>
        </nav>
      </div>
    </header>

    <main>
      <!-- Hero -->
      <section
        class="hero"
        :class="{ 'has-bg': !!siteStore.homeBackground }"
        :style="heroStyle"
      >
        <div
          class="hero-glow"
          aria-hidden="true"
        />
        <div
          class="hero-grid"
          aria-hidden="true"
        />
        <div class="container hero-inner">
          <span class="badge"><span class="badge-dot" />{{ t('ui.home.badge') }}</span>
          <h1 class="hero-title">
            <span class="oh-gradient-text">{{ siteStore.homeTitle || t('home.hero.title') }}</span>
          </h1>
          <p class="hero-sub">
            {{ siteStore.homeSubtitle || t('home.hero.description') }}
          </p>
          <div class="hero-cta">
            <router-link
              :to="userStore.isLoggedIn ? consolePath : '/register'"
              class="btn btn-primary btn-lg"
            >
              {{ userStore.isLoggedIn ? t('ui.home.ctaConsole') : t('home.hero.registerButton') }}
              <el-icon class="btn-arrow">
                <Right />
              </el-icon>
            </router-link>
            <router-link
              to="/user/store"
              class="btn btn-ghost btn-lg"
            >
              {{ t('ui.home.ctaStore') }}
            </router-link>
          </div>

          <div
            class="stats"
            aria-label="platform-stats"
          >
            <div
              v-for="s in statItems"
              :key="s.key"
              class="stat"
            >
              <i
                :class="['fas', s.icon, 'stat-icon']"
                aria-hidden="true"
              />
              <div class="stat-value">
                {{ s.value }}
              </div>
              <div class="stat-label">
                {{ s.label }}
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 栏目加载骨架屏：推荐/平台/赞助/公告数据返回前占位，避免空白等待（P2-8） -->
      <section
        v-if="homeLoading"
        class="section"
        aria-hidden="true"
      >
        <div class="container">
          <el-skeleton
            :rows="2"
            animated
          />
          <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 16px; margin-top: 16px;">
            <el-skeleton-item
              v-for="n in 3"
              :key="n"
              variant="rect"
              style="height: 140px; border-radius: 12px;"
            />
          </div>
        </div>
      </section>

      <!-- 六大特性 -->
      <section class="section">
        <div class="container">
          <div class="section-head">
            <h2>{{ t('ui.home.featuresTitle') }}</h2>
            <p>{{ t('ui.home.featuresSubtitle') }}</p>
          </div>
          <div class="features">
            <div
              v-for="f in featureItems"
              :key="f.n"
              :class="['feature', f.tone]"
            >
              <div class="feature-icon">
                <i
                  :class="['fas', f.icon]"
                  aria-hidden="true"
                />
              </div>
              <h3>{{ t(`ui.home.feature${f.n}Title`) }}</h3>
              <p>{{ t(`ui.home.feature${f.n}Desc`) }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 推荐产品（价格卡） -->
      <section
        v-if="siteStore.showRecommended && recommendedProducts.length > 0"
        class="section section-alt"
      >
        <div class="container">
          <div class="section-head">
            <h2>{{ siteStore.recommendedTitle || t('home.recommendedProducts.title') }}</h2>
            <p>{{ siteStore.recommendedSubtitle || t('home.recommendedProducts.subtitle') }}</p>
          </div>
          <div
            class="pricing"
            :style="{ '--cols': effectiveCols }"
          >
            <article
              v-for="(p, i) in recommendedProducts"
              :key="p.id"
              :class="['plan', { featured: i === featuredIndex }]"
              @click="goToProduct(p.id)"
            >
              <header class="plan-head">
                <div class="plan-icon">
                  <img
                    v-if="p.icon"
                    :src="p.icon"
                    :alt="p.name"
                  >
                  <el-icon
                    v-else
                    :size="22"
                  >
                    <Box />
                  </el-icon>
                </div>
                <div class="plan-title">
                  <h3>{{ p.name }}</h3>
                  <span class="chip">{{ categoryLabel(p) }}</span>
                </div>
              </header>
              <div
                v-if="siteStore.recommendedShowPrice"
                class="plan-price"
              >
                <span class="cur">¥</span><span class="num">{{ p.price }}</span><span class="per">/{{ periodLabel(p) }}</span>
              </div>
              <p
                v-if="p.description"
                class="plan-desc"
              >
                {{ p.description }}
              </p>
              <ul
                v-if="siteStore.recommendedShowSpecs"
                class="plan-specs"
              >
                <li><span>{{ t('user.store.cpu') }}</span><b>{{ p.cpu }} {{ t('user.store.cores') }}</b></li>
                <li><span>{{ t('user.store.memory') }}</span><b>{{ formatMB(p.memory) }}</b></li>
                <li><span>{{ t('user.store.disk') }}</span><b>{{ formatMB(p.disk) }}</b></li>
                <li><span>{{ t('user.store.bandwidth') }}</span><b>{{ p.bandwidth }} Mbps</b></li>
                <li><span>{{ t('user.store.traffic') }}</span><b>{{ formatMB(p.traffic) }}</b></li>
              </ul>
              <footer class="plan-foot">
                <span class="plan-stock">{{ p.stock < 0 ? t('user.store.stockUnlimited') : (p.stock + ' ' + t('user.store.stock')) }}</span>
                <span :class="['btn', 'btn-sm', i === featuredIndex ? 'btn-primary' : 'btn-ghost']">{{ t('ui.home.viewDetail') }}</span>
              </footer>
            </article>
          </div>
          <div class="center-row">
            <router-link
              to="/user/store"
              class="btn btn-ghost"
            >
              {{ t('ui.home.viewAll') }} <el-icon class="btn-arrow">
                <Right />
              </el-icon>
            </router-link>
          </div>
        </div>
      </section>

      <!-- 支持的虚拟化平台 -->
      <section
        v-if="siteStore.showPlatforms"
        class="section"
      >
        <div class="container">
          <div class="section-head">
            <h2>{{ t('home.platforms.title') }}</h2>
            <p>{{ t('home.platforms.description') }}</p>
          </div>
        </div>
        <div class="container">
        <LogoCarousel
          class="carousel"
          :items="platformList"
          :speed="35"
          direction="left"
          :gap="16"
        >
          <template #default="{ item }">
            <a
              :href="item.href"
              target="_blank"
              rel="noopener noreferrer"
              class="platform"
              :title="item.name"
            >
              <img
                v-if="item.icon"
                :src="item.icon"
                :alt="item.name"
                width="36"
                height="36"
                loading="eager"
                decoding="async"
                @error="e => (e.target.style.visibility = 'hidden')"
              >
              <i
                v-else
                class="fas fa-server platform-fallback-icon"
              />
              <span class="platform-text">
                <b>{{ item.name }}</b>
                <small v-if="item.repo"><i class="fab fa-github" /> {{ item.repo }}</small>
              </span>
            </a>
          </template>
        </LogoCarousel>
        </div>
      </section>

      <!-- 赞助方 -->
      <section
        v-if="siteStore.showSponsors"
        class="section section-tight"
      >
        <div class="container">
          <div class="section-head">
            <h2>{{ t('home.supporters.title') }}</h2>
            <p>{{ t('home.supporters.description') }}</p>
          </div>
          <div class="sponsors">
            <a
              v-for="item in sponsorList"
              :key="item.name"
              :href="item.href"
              target="_blank"
              rel="noopener noreferrer"
              :class="['sponsor', item.cardClass]"
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
        </div>
      </section>

      <!-- 系统公告 -->
      <section
        v-if="announcements.length > 0"
        class="section section-alt"
      >
        <div class="container">
          <div class="section-head">
            <h2>{{ t('home.announcements.title') }}</h2>
          </div>
          <div class="news">
            <article
              v-for="announcement in announcements"
              :key="announcement.id"
              class="news-item"
            >
              <div class="news-head">
                <h3>{{ announcement.title }}</h3>
                <div class="news-meta">
                  <el-tag
                    :type="announcement.type === 'homepage' ? 'success' : 'warning'"
                    size="small"
                    effect="light"
                    round
                  >
                    {{ announcement.type === 'homepage' ? t('home.announcements.typeHomepage') : t('home.announcements.typeTopbar') }}
                  </el-tag>
                  <time>{{ formatDate(announcement.createdAt) }}</time>
                </div>
              </div>
              <div
                class="news-body"
                v-html="announcement.contentHtml || announcement.content"
              />
            </article>
          </div>
        </div>
      </section>

      <!-- CTA -->
      <section class="section">
        <div class="container">
          <div class="cta">
            <div
              class="cta-glow"
              aria-hidden="true"
            />
            <div class="cta-text">
              <h2>{{ t('ui.home.ctaTitle') }}</h2>
              <p>{{ t('ui.home.ctaSubtitle') }}</p>
            </div>
            <div class="cta-actions">
              <router-link
                :to="userStore.isLoggedIn ? consolePath : '/register'"
                class="btn btn-light btn-lg"
              >
                {{ userStore.isLoggedIn ? t('ui.home.ctaConsole') : t('home.hero.registerButton') }}
              </router-link>
              <router-link
                v-if="!userStore.isLoggedIn"
                to="/login"
                class="btn btn-outline-light btn-lg"
              >
                {{ t('home.hero.loginButton') }}
              </router-link>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- 页脚 -->
    <footer class="footer">
      <div class="container footer-top">
        <div class="footer-brand">
          <div class="brand">
            <img
              :src="siteStore.logoSrc"
              alt="logo"
              class="brand-logo"
            >
            <span class="brand-name">{{ siteStore.displaySiteName }}</span>
          </div>
          <p>{{ t('home.hero.description') }}</p>
          <a
            href="https://github.com/oneclickvirt/oneclickvirt"
            target="_blank"
            rel="noopener noreferrer"
            class="btn btn-ghost btn-sm"
          ><i class="fab fa-github" /> GitHub</a>
        </div>
        <div
          v-if="footerLinkList.length > 0"
          class="footer-links"
        >
          <h4>{{ t('home.footer.friendLinks') }}</h4>
          <ul>
            <li
              v-for="l in footerLinkList"
              :key="l.id"
            >
              <a
                :href="l.href"
                target="_blank"
                rel="noopener noreferrer"
              >{{ l.name }}</a>
            </li>
          </ul>
        </div>
      </div>
      <div class="container footer-bottom">
        <span>&copy; 2026 {{ siteStore.copyrightText || (siteStore.displaySiteName + '. ' + t('home.footer.allRightsReserved')) }}</span>
        <a
          href="https://github.com/oneclickvirt"
          target="_blank"
          rel="noopener noreferrer"
        ><i class="fab fa-github" /> {{ t('home.footer.openSourceProject') }}</a>
        <template v-if="serverVersion">
          <span
            class="ver"
            :title="`${t('home.footer.serverVersion')} ${serverVersion}`"
          >{{ t('home.footer.serverVersion') }} <b>{{ serverVersion }}</b></span>
          <a
            v-if="updateAvailable && latestVersion"
            :href="releaseUrl || 'https://github.com/oneclickvirt/oneclickvirt/releases'"
            target="_blank"
            rel="noopener noreferrer"
            class="ver ver-update"
            :title="`${t('home.footer.latestVersion')} ${latestVersion}`"
          >{{ t('home.footer.latestVersion') }} <b>{{ latestVersion }}</b></a>
        </template>
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
import { Operation, Sunny, Moon, Box, Right } from '@element-plus/icons-vue'
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
import { useUserStore } from '@/pinia/modules/user'

const router = useRouter()
const { t, locale } = useI18n()
const languageStore = useLanguageStore()
const themeStore = useThemeStore()
const siteStore = useSiteStore()
const userStore = useUserStore()
// 已登录用户的控制台入口
const consolePath = computed(() => (userStore.isAnyAdmin ? '/admin/dashboard' : '/user/dashboard'))
const announcements = ref([])
// 首页栏目（推荐/平台/赞助/公告）加载中标记，驱动栏目区骨架屏（P2-8）
const homeLoading = ref(true)
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

// Hero 实时统计
const statItems = computed(() => [
  { key: 'users', icon: 'fa-users', label: t('home.stats.users'), value: usersCountDisplay.value },
  { key: 'nodes', icon: 'fa-network-wired', label: t('home.stats.nodes'), value: nodesCountDisplay.value },
  { key: 'containers', icon: 'fa-box', label: t('home.stats.containers'), value: containersCountDisplay.value },
  { key: 'vms', icon: 'fa-server', label: t('home.stats.vms'), value: vmsCountDisplay.value }
])

// 六大特性（文案见 ui.home.featureN*）
const featureItems = [
  { n: 1, icon: 'fa-layer-group', tone: 'oh-tone-indigo' },
  { n: 2, icon: 'fa-bolt', tone: 'oh-tone-sky' },
  { n: 3, icon: 'fa-sliders-h', tone: 'oh-tone-violet' },
  { n: 4, icon: 'fa-chart-line', tone: 'oh-tone-emerald' },
  { n: 5, icon: 'fa-store', tone: 'oh-tone-amber' },
  { n: 6, icon: 'fa-shield-alt', tone: 'oh-tone-rose' }
]

// 价格卡高亮：产品数 ≥3 时高亮第二个
const featuredIndex = computed(() => (recommendedProducts.value.length >= 3 ? 1 : -1))

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

const norm = (v) => String(v || '').toLowerCase().replace(/[^a-z0-9]/g, '')
const hostOf = (u) => { try { return new URL(u).host.replace(/^www\./, '') } catch { return '' } }
const repoFromUrl = (u) => {
  const m = String(u || '').match(/github\.com\/([^/]+\/[^/?#]+)/i)
  return m ? m[1] : ''
}
const findBuiltinPlatform = (l) => {
  const n = norm(l.name)
  return platforms.find(p => norm(p.name) === n || (n && norm(p.name).startsWith(n)) || (l.url && p.href === l.url.replace(/\/$/, '')))
}

// 平台轮播：优先使用后台配置的 site_links，无数据时回退到内置默认值
const platformList = computed(() => {
  if (siteLinksPlatform.value.length > 0) {
    return siteLinksPlatform.value.map(l => {
      // 后台未配置图标时，按名称或链接回退到内置图标，避免显示破图
      const builtin = findBuiltinPlatform(l)
      return {
        name: l.name,
        href: l.url || builtin?.href || '#',
        icon: l.iconUrl || builtin?.icon || '',
        repo: l.description || builtin?.repo || repoFromUrl(l.url)
      }
    })
  }
  return platforms
})

// 赞助方：优先后台配置，无数据回退内置默认值
const sponsorList = computed(() => {
  if (siteLinksSponsor.value.length > 0) {
    return siteLinksSponsor.value.map(l => {
      const builtin = footerSponsors.find(b => norm(b.name) === norm(l.name) || (l.url && b.href && hostOf(b.href) === hostOf(l.url)))
      return {
        name: l.name,
        href: l.url || builtin?.href || '#',
        logo: l.iconUrl || builtin?.logo || '',
        cardClass: builtin?.cardClass,
        logoClass: builtin?.logoClass
      }
    })
  }
  return footerSponsors
})

// 页脚友链
const footerLinkList = computed(() => footerLinks.value.map(l => ({ id: l.id, name: l.name, href: l.url })))

// 英雄区背景（站点配置可设置首页背景图）
const heroStyle = computed(() => {
  return siteStore.homeBackground
    ? { '--hero-bg': `url(${siteStore.homeBackground})` }
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

onMounted(async () => {
  // 首先加载完整站点配置（首页各栏目/英雄区文本/显示开关均由站点配置驱动）
  const configPromise = siteStore.fetchFullSiteConfig()
  // 获取后台可配置的站点链接（平台/赞助方/页脚）
  const linksPromise = fetchSiteLinks()
  // 首先检查初始化状态
  checkInitStatus()
  // 然后获取公告
  const announcementsPromise = fetchAnnouncements()
  // 获取公开统计数据（用于未登录首页展示）
  const statsPromise = fetchPublicStats()
  // 获取首页推荐产品
  const recommendedPromise = fetchRecommendedProducts()
  // 栏目数据全部落定（成功或失败）后撤掉骨架屏，避免空白等待（P2-8）
  await Promise.allSettled([configPromise, linksPromise, announcementsPromise, statsPromise, recommendedPromise])
  homeLoading.value = false
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

<style src="./home.css" scoped></style>
