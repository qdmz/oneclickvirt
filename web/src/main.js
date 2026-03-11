import '@fortawesome/fontawesome-free/css/all.css'
import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import router from './router'
import { createPinia } from 'pinia'
import App from './App.vue'
import { initUserStatusMonitor } from '@/utils/userStatusMonitor'
import i18n from './i18n'
import { getPublicSystemConfig } from '@/api/public'
import { useLanguageStore } from '@/pinia/modules/language'
import '@/assets/styles/variables.css'
import './style/main.scss'
import './style/dialog-override.css'

// Apply saved theme before mounting
const savedTheme = localStorage.getItem('theme') || 'dark'
document.documentElement.setAttribute('data-theme', savedTheme)

const app = createApp(App)
app.config.productionTip = false

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

const pinia = createPinia()
app.use(ElementPlus).use(pinia).use(i18n).use(router)

// 初始化语言设置
const initLanguage = async () => {
  const languageStore = useLanguageStore()
  
  try {
    const response = await getPublicSystemConfig()
    if (response && response.data) {
      const defaultLang = response.data.default_language
      if (defaultLang !== undefined) {
        languageStore.setSystemConfigLanguage(defaultLang)
      }
    }
  } catch (error) {
    console.warn('获取系统语言配置失败，将使用浏览器语言检测:', error)
  }
  
  const effectiveLanguage = languageStore.initLanguage()
  i18n.global.locale.value = effectiveLanguage
}

initLanguage().then(() => {
  initUserStatusMonitor()
  app.mount('#app')
})

export default app
