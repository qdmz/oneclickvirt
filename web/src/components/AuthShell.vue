<template>
  <div class="auth-shell">
    <!-- 左侧品牌面板（≤960px 隐藏） -->
    <aside class="auth-brand">
      <div class="auth-brand__grid" />
      <div class="auth-brand__glow auth-brand__glow--a" />
      <div class="auth-brand__glow auth-brand__glow--b" />

      <router-link
        to="/"
        class="auth-brand__logo"
      >
        <img
          :src="siteStore.logoSrc"
          alt="logo"
        >
        <span>{{ siteStore.displaySiteName }}</span>
      </router-link>

      <div class="auth-brand__body">
        <span class="auth-brand__tag">{{ t('ui.brand.tagline') }}</span>
        <h2>{{ t('ui.brand.headline') }}</h2>
        <p>{{ t('ui.brand.description') }}</p>
        <ul>
          <li
            v-for="i in 4"
            :key="i"
          >
            <el-icon><CircleCheckFilled /></el-icon>
            <span>{{ t(`ui.brand.point${i}`) }}</span>
          </li>
        </ul>
      </div>

      <div class="auth-brand__foot">
        <el-icon><Lock /></el-icon>
        <span>{{ t('ui.brand.secure') }}</span>
      </div>
    </aside>

    <!-- 右侧表单区 -->
    <main class="auth-main">
      <header class="auth-topbar">
        <router-link
          to="/"
          class="auth-topbar__logo"
        >
          <img
            :src="siteStore.logoSrc"
            alt="logo"
          >
          <span>{{ siteStore.displaySiteName }}</span>
        </router-link>
        <div class="auth-topbar__actions">
          <button
            class="auth-icon-btn"
            :title="themeStore.isDark ? t('navbar.lightMode') : t('navbar.darkMode')"
            :aria-label="t('navbar.switchTheme')"
            @click="themeStore.toggleTheme()"
          >
            <el-icon><component :is="themeStore.isDark ? Sunny : Moon" /></el-icon>
          </button>
          <button
            class="auth-icon-btn auth-icon-btn--text"
            :title="t('navbar.switchLanguage')"
            @click="switchLanguage"
          >
            <el-icon><Operation /></el-icon>
            <span>{{ languageStore.currentLanguage === 'zh-CN' ? 'EN' : '中文' }}</span>
          </button>
          <router-link
            to="/"
            class="auth-icon-btn auth-icon-btn--text"
          >
            <el-icon><HomeFilled /></el-icon>
            <span class="hide-xs">{{ t('common.backToHome') }}</span>
          </router-link>
        </div>
      </header>

      <section class="auth-content">
        <div class="auth-card">
          <slot />
        </div>
      </section>

      <AppFooter class="auth-footer" />
    </main>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { Sunny, Moon, Operation, HomeFilled, CircleCheckFilled, Lock } from '@element-plus/icons-vue'
import { useThemeStore } from '@/pinia/modules/theme'
import { useLanguageStore } from '@/pinia/modules/language'
import { useSiteStore } from '@/pinia/modules/site'
import AppFooter from '@/view/layout/components/AppFooter.vue'

const { t, locale } = useI18n()
const themeStore = useThemeStore()
const languageStore = useLanguageStore()
const siteStore = useSiteStore()

const switchLanguage = () => {
  locale.value = languageStore.toggleLanguage()
}
</script>

<style scoped>
.auth-shell {
  display: grid;
  grid-template-columns: minmax(380px, 46%) 1fr;
  min-height: 100vh;
  min-height: 100dvh;
  background: var(--bg-color-page);
}

/* ---------- 品牌面板 ---------- */
.auth-brand {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 40px 48px;
  overflow: hidden;
  color: #e2e8f0;
  background: radial-gradient(120% 90% at 0% 0%, #312e81 0%, #1e1b4b 38%, #0b1020 100%);
}

.auth-brand__grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(148, 163, 184, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148, 163, 184, 0.08) 1px, transparent 1px);
  background-size: 44px 44px;
  mask-image: radial-gradient(ellipse at 30% 40%, #000 30%, transparent 75%);
  -webkit-mask-image: radial-gradient(ellipse at 30% 40%, #000 30%, transparent 75%);
  pointer-events: none;
}

.auth-brand__glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.55;
  pointer-events: none;
}

.auth-brand__glow--a {
  width: 380px;
  height: 380px;
  right: -120px;
  top: -80px;
  background: #6366f1;
}

.auth-brand__glow--b {
  width: 320px;
  height: 320px;
  left: -100px;
  bottom: -120px;
  background: #0ea5e9;
  opacity: 0.4;
}

.auth-brand__logo,
.auth-brand__body,
.auth-brand__foot {
  position: relative;
  z-index: 1;
}

.auth-brand__logo {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  color: #fff;
  text-decoration: none;
  font-size: 18px;
  font-weight: 650;
  letter-spacing: -0.01em;
}

.auth-brand__logo img {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  object-fit: contain;
  background: rgba(255, 255, 255, 0.08);
  padding: 4px;
}

.auth-brand__tag {
  display: inline-block;
  padding: 4px 12px;
  border: 1px solid rgba(165, 180, 252, 0.35);
  border-radius: 999px;
  font-size: 12px;
  color: #c7d2fe;
  background: rgba(99, 102, 241, 0.14);
}

.auth-brand__body h2 {
  margin: 20px 0 14px;
  font-size: clamp(28px, 3vw, 38px);
  line-height: 1.2;
  font-weight: 700;
  letter-spacing: -0.02em;
  background: linear-gradient(120deg, #fff 0%, #c7d2fe 55%, #7dd3fc 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.auth-brand__body p {
  max-width: 460px;
  margin: 0 0 28px;
  font-size: 15px;
  line-height: 1.7;
  color: #a5b4fc;
  opacity: 0.9;
}

.auth-brand__body ul {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 14px;
}

.auth-brand__body li {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: #e2e8f0;
}

.auth-brand__body li .el-icon {
  color: #38bdf8;
  font-size: 18px;
}

.auth-brand__foot {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #94a3b8;
}

/* ---------- 表单区 ---------- */
.auth-main {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.auth-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 20px 32px;
}

.auth-topbar__logo {
  display: none;
  align-items: center;
  gap: 10px;
  color: var(--text-color-primary);
  text-decoration: none;
  font-weight: 650;
  font-size: 16px;
  min-width: 0;
}

.auth-topbar__logo span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.auth-topbar__logo img {
  width: 30px;
  height: 30px;
  object-fit: contain;
}

.auth-topbar__actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

.auth-icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 36px;
  min-width: 36px;
  padding: 0 10px;
  border: 1px solid var(--border-color);
  border-radius: 10px;
  background: var(--surface-elevated, var(--bg-color-secondary));
  color: var(--text-color-secondary);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  text-decoration: none;
  transition: color 0.15s, border-color 0.15s, background 0.15s;
}

.auth-icon-btn:hover {
  color: var(--primary-color);
  border-color: var(--border-color-hover);
  background: var(--primary-color-bg);
}

.auth-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 32px 40px;
}

.auth-card {
  width: 100%;
  max-width: 420px;
}

.auth-footer {
  background: transparent !important;
  border-top: none !important;
}

/* 统一子页面表单标题样式 */
.auth-card :deep(h2) {
  margin: 0 0 8px;
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--text-color-primary);
  background: none;
  -webkit-text-fill-color: currentColor;
  text-align: left;
}

.auth-card :deep(h2 + p),
.auth-card :deep([class$='-header'] p) {
  text-align: left;
  color: var(--text-color-secondary);
  font-size: 14px;
}

.auth-card :deep([class$='-header']) {
  text-align: left;
}

.auth-card :deep(.el-input__wrapper) {
  border-radius: 10px;
}

.auth-card :deep(.el-button--primary) {
  height: 44px;
  border-radius: 10px;
  font-weight: 600;
  border: none;
  background: var(--brand-gradient);
  box-shadow: var(--brand-glow);
}

.auth-card :deep(.el-button--primary:hover) {
  filter: brightness(1.06);
}

@media (max-width: 960px) {
  .auth-shell {
    grid-template-columns: 1fr;
  }

  .auth-brand {
    display: none;
  }

  .auth-topbar__logo {
    display: inline-flex;
  }

  .auth-main {
    background:
      radial-gradient(60% 40% at 100% 0%, rgba(99, 102, 241, 0.10), transparent 70%),
      radial-gradient(50% 35% at 0% 100%, rgba(14, 165, 233, 0.08), transparent 70%);
  }

  .auth-card {
    padding: 28px 24px;
    border: 1px solid var(--border-color);
    border-radius: 16px;
    background: var(--surface-elevated, var(--bg-color-secondary));
    box-shadow: var(--shadow-lg);
  }
}

@media (max-width: 768px) {
  .auth-topbar {
    padding: 14px 16px;
  }

  .auth-content {
    padding: 12px 16px 28px;
    align-items: flex-start;
  }
}

@media (max-width: 480px) {
  .auth-card {
    padding: 22px 18px;
  }

  .auth-card :deep(h2) {
    font-size: 22px;
  }

  .hide-xs {
    display: none;
  }

  .auth-topbar__logo span {
    max-width: 120px;
  }
}
</style>
