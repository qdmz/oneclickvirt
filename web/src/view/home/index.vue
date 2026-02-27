<template>
  <div class="home-container">
    <!-- 导航栏 -->
    <header class="home-header">
      <div class="header-content">
        <div class="logo">
          <img
            :src="siteConfigs.site_icon_url || logoUrl"
            :alt="siteConfigs.site_name || 'OneClickVirt Logo'"
            class="logo-image"
          >
          <a
            :href="siteConfigs.site_url || '#'"
            target="_self"
            class="site-name-link"
          >
            <h1>{{ siteConfigs.site_name || t('home.title') }}</h1>
          </a>
        </div>
        <nav class="nav-menu">
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
      <section class="hero-section">
        <div class="hero-content">
          <div
            v-if="siteConfigs.site_header"
            class="custom-header-content"
            v-html="siteConfigs.site_header"
          />
          <div v-else>
            <h1 class="hero-title">
              {{ t('home.hero.title') }}
            </h1>
            <p class="hero-description">
              {{ t('home.hero.description') }}
            </p>
          </div>
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

      <!-- 支持的虚拟化平台 -->
      <section class="platforms-section">
        <div class="section-header">
          <h2>{{ t('home.platforms.title') }}</h2>
          <p>{{ t('home.platforms.description') }}</p>
        </div>
        <div class="platforms-grid">
          <div class="platform-item">
            <div class="platform-icon pve-icon">
              <img
                src="@/assets/images/proxmox.png"
                alt="Proxmox VE"
                width="60"
                height="60"
              >
            </div>
            <h3>Proxmox VE</h3>
          </div>
          
          <div class="platform-item">
            <div class="platform-icon incus-icon">
              <img
                src="@/assets/images/incus.png"
                alt="Incus"
                width="60"
                height="60"
              >
            </div>
            <h3>Incus</h3>
          </div>
          
          <div class="platform-item">
            <div class="platform-icon docker-icon">
              <img
                src="@/assets/images/docker.png"
                alt="Docker"
                width="60"
                height="60"
              >
            </div>
            <h3>Docker</h3>
          </div>
          
          <div class="platform-item">
            <div class="platform-icon lxd-icon">
              <img
                src="@/assets/images/lxd.png"
                alt="LXD"
                width="60"
                height="60"
              >
            </div>
            <h3>LXD</h3>
          </div>
        </div>
        <!-- 统计概况：与平台卡片相同的框架风格，显示用户/节点/容器/虚拟机数量 -->
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

      <!-- 产品列表 -->
      <section
        v-if="products.length > 0"
        class="products-section"
      >
        <div class="section-header">
          <h2>{{ siteConfigs.site_name ? siteConfigs.site_name + ' 产品' : t('home.products.title') }}</h2>
          <p>{{ t('home.products.description') }}</p>
        </div>
        <div class="products-grid">
          <div
            v-for="product in products"
            :key="product.id"
            class="product-item"
          >
            <div class="product-card">
              <div class="product-header">
                <span class="product-name">{{ product.name }}</span>
                <el-tag
                  v-if="product.period === 0"
                  type="success"
                >
                  永久
                </el-tag>
                <el-tag
                  v-else
                  type="info"
                >
                  {{ product.period }}天
                </el-tag>
              </div>
              <div class="product-body">
                <div class="price">
                  <span class="amount">¥{{ (product.price / 100).toFixed(2) }}</span>
                  <span class="unit">{{ product.period === 0 ? '/永久' : `/${product.period}天` }}</span>
                </div>
                <div class="level-badge">
                  Lv.{{ product.level }}
                </div>
                <div class="features">
                  <div class="feature-item">
                    <i class="fas fa-microchip" />
                    <span>{{ product.cpu }} 核</span>
                  </div>
                  <div class="feature-item">
                    <i class="fas fa-memory" />
                    <span>{{ product.memory }} MB</span>
                  </div>
                  <div class="feature-item">
                    <i class="fas fa-hdd" />
                    <span>{{ product.disk }} MB</span>
                  </div>
                  <div class="feature-item">
                    <i class="fas fa-wifi" />
                    <span>{{ product.bandwidth }} Mbps</span>
                  </div>
                  <div class="feature-item">
                    <i class="fas fa-network-wired" />
                    <span>{{ product.traffic }} MB</span>
                  </div>
                  <div class="feature-item">
                    <i class="fas fa-server" />
                    <span>{{ product.maxInstances }} 个</span>
                  </div>
                </div>
                <router-link
                  :to="'/login'"
                  type="primary"
                  class="purchase-btn"
                >
                  立即购买
                </router-link>
              </div>
            </div>
          </div>
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
    </main>
    
    <!-- 页脚 -->
    <footer class="home-footer">
      <div class="footer-content">
        <div class="footer-section">
          <h3>{{ siteConfigs.site_name || 'OneClickVirt' }}</h3>
          <div class="company-info">
            <p v-if="siteConfigs.company_name">
              公司名称: {{ siteConfigs.company_name }}
            </p>
            <p v-if="siteConfigs.company_address">
              公司地址: {{ siteConfigs.company_address }}
            </p>
            <p v-if="siteConfigs.contact_email">
              联系邮箱: <a :href="'mailto:' + siteConfigs.contact_email">{{ siteConfigs.contact_email }}</a>
            </p>
            <!-- 直接显示联系电话，不使用v-if条件判断，确保始终显示 -->
            <p>
              联系电话: {{ siteConfigs.contact_phone || '888-888-8888' }}
            </p>
          </div>
        </div>
        <div class="footer-section">
          <h4>联系我们</h4>
          <ul>
            <li v-if="siteConfigs.contact_email">
              <a :href="'mailto:' + siteConfigs.contact_email">{{ siteConfigs.contact_email }}</a>
            </li>
            <!-- 直接显示联系电话，不使用v-if条件判断，确保始终显示 -->
            <li>
              {{ siteConfigs.contact_phone || '888-888-8888' }}
            </li>
            <li v-if="siteConfigs.company_address">
              {{ siteConfigs.company_address }}
            </li>
          </ul>
        </div>
        <div class="footer-section">
          <h4>法律信息</h4>
          <ul>
            <li v-if="siteConfigs.icp_number">
              <a
                :href="'https://beian.miit.gov.cn/'"
                target="_blank"
                rel="noopener noreferrer"
              >
                ICP备案号: {{ siteConfigs.icp_number }}
              </a>
            </li>
            <li v-if="siteConfigs.company_name">
              公司名称: {{ siteConfigs.company_name }}
            </li>
          </ul>
        </div>
        <div class="footer-section">
          <h4>关于我们</h4>
          <ul>
            <li v-if="siteConfigs.site_url">
              <a
                :href="siteConfigs.site_url"
                target="_self"
              >
                {{ siteConfigs.site_name || '官方网站' }}
              </a>
            </li>
            <li v-if="siteConfigs.site_description">
              {{ siteConfigs.site_description }}
            </li>
          </ul>
        </div>
      </div>
      <div class="footer-bottom">
        <div
          v-if="siteConfigs.site_footer"
          v-html="siteConfigs.site_footer"
        />
        <div v-else>
          <p>
            &copy; 2026 {{ siteConfigs.site_name || 'OneClickVirt' }}. {{ t('home.footer.allRightsReserved') }}
            <span v-if="siteConfigs.icp_number">
              | <a
                :href="'https://beian.miit.gov.cn/'"
                target="_blank"
                rel="noopener noreferrer"
              >
                ICP备案号: {{ siteConfigs.icp_number }}
              </a>
            </span>
            |
            <a
              href="https://github.com/oneclickvirt"
              target="_blank"
              rel="noopener noreferrer"
            >
              {{ t('home.footer.openSourceProject') }}
            </a>
          </p>
        </div>
        <!-- 统计代码 -->
        <div
          v-if="siteConfigs.analytics_code"
          class="analytics-code"
          v-html="siteConfigs.analytics_code"
        />
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getPublicAnnouncements, getPublicStats, getPublicSiteConfigs, getPublicProducts } from '@/api/public'
import { checkSystemInit } from '@/api/init'
import { ElTag, ElMessage } from 'element-plus'
import { Operation } from '@element-plus/icons-vue'
import { useLanguageStore } from '@/pinia/modules/language'
import logoUrl from '@/assets/images/logo.png'

const router = useRouter()
const { t, locale } = useI18n()
const languageStore = useLanguageStore()
const announcements = ref([])
const siteConfigs = ref({})
const products = ref([])
// 统计数据
const usersCount = ref(null)
const nodesCount = ref(null)
const containersCount = ref(null)
const vmsCount = ref(null)
// 调试用变量
const debugApiResponse = ref('')
const debugSiteConfigs = ref('')

const usersCountDisplay = computed(() => (usersCount.value === null ? '-' : usersCount.value))
const nodesCountDisplay = computed(() => (nodesCount.value === null ? '-' : nodesCount.value))
const containersCountDisplay = computed(() => (containersCount.value === null ? '-' : containersCount.value))
const vmsCountDisplay = computed(() => (vmsCount.value === null ? '-' : vmsCount.value))

// 动态样式注入相关
const styleElement = ref(null)

// 注入自定义CSS
const injectCustomCSS = () => {
  // 清理现有样式
  if (styleElement.value) {
    styleElement.value.remove()
    styleElement.value = null
  }
  
  // 检查配置是否有效
  if (siteConfigs.value && siteConfigs.value.custom_css && 
      typeof siteConfigs.value.custom_css === 'string' && 
      siteConfigs.value.custom_css.trim()) {
    
    // 创建新的样式标签
    const style = document.createElement('style')
    style.textContent = siteConfigs.value.custom_css
    document.head.appendChild(style)
    styleElement.value = style
  }
}

// 监听站点配置变化
watch(
  () => siteConfigs.value,
  (newConfigs) => {
    injectCustomCSS()
  },
  { deep: true }
)

// 组件卸载时清理
onUnmounted(() => {
  if (styleElement.value) {
    styleElement.value.remove()
    styleElement.value = null
  }
})

// 获取站点配置
const fetchSiteConfigs = async () => {
  try {
    const resp = await getPublicSiteConfigs()
    console.log('站点配置API响应完整内容:', JSON.stringify(resp, null, 2))
    debugApiResponse.value = JSON.stringify(resp, null, 2)
    if (resp && (resp.code === 0 || resp.code === 200) && resp.data) {
      const configs = resp.data
      // 检查数据格式，如果是对象直接使用，如果是数组则遍历
      if (Array.isArray(configs)) {
        configs.forEach(config => {
          siteConfigs.value[config.key] = config.value
        })
      } else {
        // 直接将对象赋值给siteConfigs
        siteConfigs.value = { ...siteConfigs.value, ...configs }
      }
      // 添加默认的联系电话，确保前端能正常显示
      if (!siteConfigs.value.contact_phone) {
        siteConfigs.value.contact_phone = '888-888-8888'
      }
      debugSiteConfigs.value = JSON.stringify(siteConfigs.value, null, 2)
    }
    console.log('获取到的站点配置:', JSON.stringify(siteConfigs.value, null, 2))
    console.log('联系电话(contact_phone):', siteConfigs.value.contact_phone)
    console.log('联系电话(contactPhone):', siteConfigs.value.contactPhone)
    console.log('联系电话字段是否存在:', 'contact_phone' in siteConfigs.value)
    console.log('联系电话字段值:', siteConfigs.value.contact_phone)
  } catch (error) {
    console.error('获取站点配置失败', error)
    console.error('错误详情:', JSON.stringify(error, null, 2))
    debugApiResponse.value = JSON.stringify(error, null, 2)
    // 即使API调用失败，也设置默认的联系电话
    siteConfigs.value.contact_phone = '888-888-8888'
  }
  
  // 注入自定义CSS
  injectCustomCSS()
}

// 获取产品列表
const fetchProducts = async () => {
  try {
    const resp = await getPublicProducts()
    if (resp && (resp.code === 0 || resp.code === 200) && resp.data) {
      // 将后端返回的整数isEnabled转换为布尔值
      products.value = resp.data.map(product => ({
        ...product,
        isEnabled: product.isEnabled === 1
      }))
    }
  } catch (error) {
    console.error('获取产品列表失败', error)
  }
}

const switchLanguage = () => {
  const newLang = languageStore.toggleLanguage()
  locale.value = newLang
  ElMessage.success(t('navbar.languageSwitched'))
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString(locale.value === 'zh-CN' ? 'zh-CN' : 'en-US')
}

const fetchAnnouncements = async () => {
  try {
    // 获取首页公告
    const response = await getPublicAnnouncements('homepage')
    if (response.code === 0 || response.code === 200) {
      announcements.value = response.data.slice(0, 3) // 只显示最新3条
    }
  } catch (error) {
    console.error(t('home.errors.fetchAnnouncementsFailed'), error)
  }
}

const fetchPublicStats = async () => {
  try {
    const resp = await getPublicStats()
    if (resp && (resp.code === 0 || resp.code === 200) && resp.data) {
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

const checkInitStatus = async () => {
  try {
    const response = await checkSystemInit()
    console.log(t('home.debug.checkingInit'), response)
    if (response && response.code === 0 && response.data && response.data.needInit === true) {
      console.log(t('home.debug.needInitRedirect'))
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
  console.log('VITE_BASE_API:', import.meta.env.VITE_BASE_API)
  console.log('VITE_BASE_PATH:', import.meta.env.VITE_BASE_PATH)
  console.log('VITE_SERVER_PORT:', import.meta.env.VITE_SERVER_PORT)
  console.log('All env vars:', import.meta.env)
  
  // 首先检查初始化状态
  checkInitStatus()
  // 获取站点配置
  fetchSiteConfigs()
  // 获取产品列表
  fetchProducts()
  // 然后获取公告
  fetchAnnouncements()
  // 获取公开统计数据（用于未登录首页展示）
  fetchPublicStats()
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f0fdf4 0%, #ecfdf5 100%);
}

/* 头部样式 */
.home-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  box-shadow: 0 2px 20px rgba(22, 163, 74, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
  border-bottom: 1px solid rgba(22, 163, 74, 0.1);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 70px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-image {
  width: 64px;
  height: 64px;
  object-fit: contain;
}

.site-name-link {
  text-decoration: none;
}

.logo h1 {
  font-size: 28px;
  color: #16a34a;
  margin: 0;
  font-weight: 700;
  background: linear-gradient(135deg, #16a34a, #22c55e);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  transition: all 0.3s ease;
}

.logo h1:hover {
  transform: scale(1.05);
}

.nav-menu {
  display: flex;
  align-items: center;
}

.nav-link {
  text-decoration: none;
  color: #374151;
  padding: 12px 24px;
  border-radius: 25px;
  transition: all 0.3s ease;
  font-weight: 500;
  margin-left: 12px;
  position: relative;
  overflow: hidden;
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.nav-link.language-btn {
  border: 1px solid #e5e7eb;
}

.nav-link:hover {
  background: rgba(22, 163, 74, 0.1);
  color: #16a34a;
  transform: translateY(-2px);
}

.nav-link.primary {
  background: linear-gradient(135deg, #16a34a, #22c55e);
  color: white;
  box-shadow: 0 4px 15px rgba(22, 163, 74, 0.3);
}

.nav-link.primary:hover {
  background: linear-gradient(135deg, #15803d, #16a34a);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(22, 163, 74, 0.4);
}

/* 主要内容 */
.home-main {
  padding: 60px 0;
}

/* 英雄区域 */
.hero-section {
  display: flex;
  justify-content: center;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  padding: 60px 24px;
  gap: 60px;
  flex-wrap: wrap;
}

.hero-content {
  flex: 1;
  min-width: 400px;
}

.hero-title {
  font-size: 52px;
  color: #1f2937;
  margin-bottom: 24px;
  line-height: 1.2;
  font-weight: 800;
  background: linear-gradient(135deg, #1f2937, #374151);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-description {
  font-size: 20px;
  color: #6b7280;
  margin-bottom: 40px;
  line-height: 1.6;
  font-weight: 400;
}

/* 自定义页眉内容样式 */
.custom-header-content {
  margin-bottom: 30px;
}

.custom-header-content h1 {
  font-size: 52px;
  color: #1f2937;
  margin-bottom: 24px;
  line-height: 1.2;
  font-weight: 800;
  background: linear-gradient(135deg, #1f2937, #374151);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.custom-header-content p {
  font-size: 20px;
  color: #6b7280;
  margin-bottom: 40px;
  line-height: 1.6;
  font-weight: 400;
}

.hero-actions {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.btn {
  display: inline-block;
  padding: 16px 32px;
  border-radius: 30px;
  text-decoration: none;
  font-weight: 600;
  font-size: 16px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  border: none;
  cursor: pointer;
}

.btn-primary {
  background: linear-gradient(135deg, #16a34a, #22c55e);
  color: white;
  box-shadow: 0 4px 15px rgba(22, 163, 74, 0.3);
}

.btn-primary:hover {
  background: linear-gradient(135deg, #15803d, #16a34a);
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(22, 163, 74, 0.4);
}

.btn-secondary {
  background: transparent;
  color: #16a34a;
  border: 2px solid #16a34a;
  box-shadow: 0 4px 15px rgba(22, 163, 74, 0.1);
}

.btn-secondary:hover {
  background: #16a34a;
  color: white;
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(22, 163, 74, 0.3);
}

.hero-image {
  flex: 1;
  min-width: 400px;
}

.feature-preview {
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
}

.preview-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  padding: 24px;
  border-radius: 20px;
  box-shadow: 0 8px 25px rgba(22, 163, 74, 0.1);
  text-align: center;
  transition: all 0.3s ease;
  border: 1px solid rgba(22, 163, 74, 0.1);
}

.preview-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 15px 35px rgba(22, 163, 74, 0.2);
  border-color: rgba(22, 163, 74, 0.3);
}

.card-icon {
  font-size: 42px;
  margin-bottom: 16px;
}

.preview-card h3 {
  font-size: 18px;
  color: #1f2937;
  margin-bottom: 8px;
  font-weight: 600;
}

.preview-card p {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.5;
}

/* 支持的虚拟化平台 */
.platforms-section {
  max-width: 1200px;
  margin: 100px auto;
  padding: 60px 24px;
  text-align: center;
}

.section-header {
  margin-bottom: 60px;
}

.section-header h2 {
  font-size: 42px;
  color: #1f2937;
  margin: 0 0 16px 0;
  font-weight: 700;
  background: linear-gradient(135deg, #1f2937, #374151);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.section-header p {
  font-size: 18px;
  color: #6b7280;
  margin: 0;
  font-weight: 400;
}

.platforms-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 32px;
  margin-top: 60px;
}

/* 统计概况网格：复用 platform-item 的视觉样式，使其与图标卡片一致 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 24px;
  margin-top: 36px;
}

.stats-item .platform-icon {
  height: 56px;
}

.stats-value {
  font-size: 28px;
  color: #16a34a;
  font-weight: 700;
  margin-top: 12px;
}

.platform-item {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  padding: 40px 24px;
  border-radius: 24px;
  box-shadow: 0 8px 25px rgba(22, 163, 74, 0.08);
  transition: all 0.3s ease;
  border: 1px solid rgba(22, 163, 74, 0.1);
  text-align: center;
}

.platform-item:hover {
  transform: translateY(-10px) scale(1.03);
  box-shadow: 0 20px 40px rgba(22, 163, 74, 0.15);
  border-color: rgba(22, 163, 74, 0.3);
}

.platform-icon {
  margin-bottom: 24px;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 80px;
}

.platform-item h3 {
  font-size: 20px;
  color: #1f2937;
  margin-bottom: 12px;
  font-weight: 600;
}

.platform-item p {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.5;
}

/* 产品列表 */
.products-section {
  max-width: 1200px;
  margin: 100px auto;
  padding: 60px 24px;
  text-align: center;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 32px;
  margin-top: 60px;
}

.product-item {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  padding: 24px;
  border-radius: 24px;
  box-shadow: 0 8px 25px rgba(22, 163, 74, 0.08);
  transition: all 0.3s ease;
  border: 1px solid rgba(22, 163, 74, 0.1);
}

.product-item:hover {
  transform: translateY(-10px) scale(1.03);
  box-shadow: 0 20px 40px rgba(22, 163, 74, 0.15);
  border-color: rgba(22, 163, 74, 0.3);
}

.product-card {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.product-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.product-name {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

.product-body {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.price {
  text-align: center;
  margin-bottom: 20px;
}

.amount {
  font-size: 32px;
  font-weight: 700;
  color: #16a34a;
}

.unit {
  font-size: 16px;
  color: #6b7280;
  margin-left: 4px;
}

.level-badge {
  display: inline-block;
  background: linear-gradient(135deg, #16a34a, #22c55e);
  color: white;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 20px;
  align-self: center;
}

.features {
  flex: 1;
  margin-bottom: 24px;
  text-align: left;
}

.feature-item {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  font-size: 16px;
  color: #6b7280;
}

.feature-item i {
  margin-right: 12px;
  color: #16a34a;
  width: 20px;
  text-align: center;
}

.purchase-btn {
  display: block;
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #16a34a, #22c55e);
  color: white;
  text-align: center;
  text-decoration: none;
  border-radius: 12px;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(22, 163, 74, 0.3);
  border: none;
  cursor: pointer;
  font-size: 16px;
}

.purchase-btn:hover {
  background: linear-gradient(135deg, #15803d, #16a34a);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(22, 163, 74, 0.4);
}

/* 系统公告 */
.announcements-section {
  max-width: 1200px;
  margin: 100px auto;
  padding: 60px 24px;
}

.announcements-list {
  display: grid;
  gap: 20px;
  margin-top: 40px;
}

.announcement-item {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 4px 15px rgba(22, 163, 74, 0.05);
  border: 1px solid rgba(22, 163, 74, 0.1);
  transition: all 0.3s ease;
}

.announcement-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(22, 163, 74, 0.1);
  border-color: rgba(22, 163, 74, 0.2);
}

.announcement-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 8px;
}

.announcement-header h3 {
  font-size: 18px;
  color: #1f2937;
  font-weight: 600;
  margin: 0;
  flex: 1;
  min-width: 200px;
}

.announcement-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.announcement-date {
  font-size: 14px;
  color: #6b7280;
  font-weight: 400;
}

.announcement-content {
  font-size: 16px;
  color: #6b7280;
  line-height: 1.6;
  margin: 0;
}

/* 富文本内容样式 */
.announcement-content :deep(p) {
  margin: 8px 0;
}

.announcement-content :deep(ul),
.announcement-content :deep(ol) {
  padding-left: 20px;
  margin: 8px 0;
}

.announcement-content :deep(blockquote) {
  border-left: 4px solid #16a34a;
  padding-left: 16px;
  margin: 16px 0;
  font-style: italic;
  background: rgba(22, 163, 74, 0.05);
  padding: 12px 16px;
  border-radius: 4px;
}

.announcement-content :deep(strong) {
  color: #1f2937;
  font-weight: 600;
}

.announcement-content :deep(code) {
  background: rgba(22, 163, 74, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
}

/* 页脚 */
.home-footer {
  background: linear-gradient(135deg, #1f2937, #374151);
  color: white;
  padding: 60px 24px 24px;
  font-size: 14px;
  margin-top: 100px;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 40px;
}

.footer-section {
  flex: 1;
  min-width: 200px;
}

.footer-section h3,
.footer-section h4 {
  color: white;
  margin-bottom: 20px;
  font-size: 18px;
  font-weight: 600;
}

.company-info {
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.6;
}

.company-info p {
  margin-bottom: 8px;
}

.company-info a {
  color: #22c55e;
  text-decoration: none;
  transition: all 0.3s ease;
}

.company-info a:hover {
  color: #34d399;
}

.footer-section ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-section ul li {
  margin-bottom: 8px;
}

.footer-section ul li a,
.social-link {
  color: rgba(255, 255, 255, 0.7);
  text-decoration: none;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  font-weight: 400;
}

.footer-section ul li a:hover,
.social-link:hover {
  color: #22c55e;
  transform: translateX(5px);
}

.social-link svg {
  margin-right: 8px;
}

.footer-bottom {
  text-align: center;
  margin-top: 40px;
  padding-top: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.footer-bottom p {
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
}

.footer-bottom a {
  color: #22c55e;
  text-decoration: none;
  transition: all 0.3s ease;
}

.footer-bottom a:hover {
  color: #34d399;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .hero-section {
    flex-direction: column;
    text-align: center;
    gap: 40px;
    padding: 40px 20px;
  }

  .hero-content {
    min-width: unset;
  }

  .hero-title {
    font-size: 36px;
  }

  .hero-description {
    font-size: 18px;
  }

  .hero-actions {
    justify-content: center;
  }

  .hero-image {
    min-width: unset;
    width: 100%;
  }

  .platforms-grid {
    grid-template-columns: 1fr;
    gap: 24px;
  }

  .platform-item {
    padding: 32px 20px;
  }

  .footer-content {
    flex-direction: column;
    text-align: center;
  }

  .footer-section {
    margin-bottom: 32px;
  }

  .footer-section ul li a,
  .social-link {
    justify-content: center;
  }

  .footer-section ul li a:hover {
    transform: none;
  }

  .header-content {
    padding: 0 20px;
  }

  .section-header h2 {
    font-size: 32px;
  }

  .section-header p {
    font-size: 16px;
  }
}

@media (max-width: 480px) {
  .hero-title {
    font-size: 28px;
  }

  .hero-description {
    font-size: 16px;
  }

  .btn {
    padding: 14px 28px;
    font-size: 15px;
  }

  .platforms-section,
  .announcements-section {
    padding: 40px 20px;
  }

  .section-header h2 {
    font-size: 28px;
  }
}
</style>