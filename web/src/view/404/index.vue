<template>
  <div class="nf-page">
    <div
      class="nf-grid"
      aria-hidden="true"
    />
    <div
      class="nf-glow nf-glow--a"
      aria-hidden="true"
    />
    <div
      class="nf-glow nf-glow--b"
      aria-hidden="true"
    />

    <header class="nf-topbar">
      <router-link
        to="/"
        class="nf-logo"
      >
        <img
          :src="siteStore.logoSrc"
          alt="logo"
        >
        <span>{{ siteStore.displaySiteName }}</span>
      </router-link>
    </header>

    <main class="nf-main">
      <div class="nf-card">
        <span class="nf-badge">
          <span class="nf-badge__dot" />
          HTTP 404
        </span>

        <h1
          class="nf-code"
          aria-hidden="true"
        >
          404
        </h1>

        <h2 class="nf-title">
          {{ t('notFound.title') }}
        </h2>
        <p class="nf-message">
          {{ t('notFound.message') }}
        </p>
        <p class="nf-hint">
          {{ t('ui.notFound.hint') }}
        </p>

        <div class="nf-actions">
          <router-link
            to="/"
            class="nf-btn nf-btn--primary"
          >
            <el-icon><HomeFilled /></el-icon>
            <span>{{ t('notFound.goHome') }}</span>
          </router-link>
          <router-link
            :to="consolePath"
            class="nf-btn nf-btn--ghost"
          >
            <el-icon><Monitor /></el-icon>
            <span>{{ t('ui.layout.console') }}</span>
          </router-link>
        </div>

        <button
          type="button"
          class="nf-back"
          @click="goBack"
        >
          <el-icon><Back /></el-icon>
          <span>{{ t('notFound.goBack') }}</span>
        </button>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { HomeFilled, Monitor, Back } from '@element-plus/icons-vue'
import { useUserStore } from '@/pinia/modules/user'
import { useSiteStore } from '@/pinia/modules/site'

const router = useRouter()
const { t } = useI18n()
const userStore = useUserStore()
const siteStore = useSiteStore()

// 控制台入口：已登录按身份/视图模式进入对应仪表盘，未登录进入登录页
const consolePath = computed(() => {
  if (!userStore.token) return '/login'
  const userType = userStore.userType
  const viewMode = userStore.viewMode || userType
  const isAdmin = userType === 'admin' || userType === 'normal_admin'
  const inAdminView = viewMode === 'admin' || viewMode === 'normal_admin'
  return isAdmin && inAdminView ? '/admin/dashboard' : '/user/dashboard'
})

const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push('/')
  }
}
</script>

<style scoped>
.nf-page {
  position: relative;
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-color-page);
  color: var(--text-color-primary);
  font-family: var(--font-sans);
}

/* 背景：细网格 + 品牌光晕 */
.nf-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(var(--border-color-light) 1px, transparent 1px),
    linear-gradient(90deg, var(--border-color-light) 1px, transparent 1px);
  background-size: 48px 48px;
  mask-image: radial-gradient(ellipse 70% 60% at 50% 45%, #000 20%, transparent 75%);
  -webkit-mask-image: radial-gradient(ellipse 70% 60% at 50% 45%, #000 20%, transparent 75%);
  opacity: 0.6;
  pointer-events: none;
}

.nf-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  pointer-events: none;
}

.nf-glow--a {
  width: 520px;
  height: 520px;
  top: -180px;
  left: 50%;
  transform: translateX(-70%);
  background: var(--brand-gradient-soft);
  opacity: 0.9;
}

.nf-glow--b {
  width: 360px;
  height: 360px;
  bottom: -140px;
  right: -80px;
  background: color-mix(in srgb, var(--brand-accent) 18%, transparent);
}

/* 顶栏 */
.nf-topbar {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  padding: 20px 32px;
}

.nf-logo {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: var(--text-color-primary);
  text-decoration: none;
  font-weight: var(--font-weight-semibold);
  font-size: 16px;
}

.nf-logo img {
  width: 28px;
  height: 28px;
  object-fit: contain;
  border-radius: 6px;
}

/* 主体 */
.nf-main {
  position: relative;
  z-index: 1;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 20px calc(48px + env(safe-area-inset-bottom));
}

.nf-card {
  width: 100%;
  max-width: 560px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.nf-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 5px 12px;
  border-radius: var(--border-radius-full);
  border: 1px solid var(--border-color);
  background: var(--surface-elevated, var(--card-bg-solid));
  color: var(--text-color-secondary);
  font-family: var(--font-mono);
  font-size: 12px;
  letter-spacing: 0.04em;
  box-shadow: var(--shadow-xs);
}

.nf-badge__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--warning-color);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--warning-color) 22%, transparent);
}

.nf-code {
  margin: 20px 0 8px;
  font-size: clamp(112px, 22vw, 188px);
  line-height: 0.95;
  font-weight: var(--font-weight-black, 900);
  letter-spacing: -0.06em;
  background: var(--brand-gradient-text);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  color: transparent;
  user-select: none;
  filter: drop-shadow(0 18px 40px color-mix(in srgb, var(--primary-color) 22%, transparent));
}

.nf-title {
  margin: 0 0 10px;
  font-size: clamp(22px, 3.2vw, 28px);
  font-weight: var(--font-weight-bold);
  letter-spacing: -0.01em;
  color: var(--text-color-primary);
}

.nf-message {
  margin: 0;
  font-size: 16px;
  line-height: var(--line-height-relaxed, 1.65);
  color: var(--text-color-secondary);
}

.nf-hint {
  margin: 6px 0 0;
  font-size: 14px;
  line-height: var(--line-height-relaxed, 1.65);
  color: var(--text-color-tertiary);
}

.nf-actions {
  margin-top: 32px;
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
}

.nf-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 44px;
  padding: 0 22px;
  border-radius: var(--border-radius-full);
  font-size: 15px;
  font-weight: var(--font-weight-semibold);
  text-decoration: none;
  transition: transform var(--duration-fast, 0.15s) var(--ease-out, ease-out),
    box-shadow var(--duration-fast, 0.15s) var(--ease-out, ease-out),
    background-color var(--duration-fast, 0.15s) var(--ease-out, ease-out),
    border-color var(--duration-fast, 0.15s) var(--ease-out, ease-out);
}

.nf-btn:focus-visible,
.nf-back:focus-visible,
.nf-logo:focus-visible {
  outline: 2px solid var(--primary-color);
  outline-offset: 3px;
}

.nf-btn--primary {
  color: #fff;
  background: var(--brand-gradient);
  box-shadow: var(--brand-glow);
}

.nf-btn--primary:hover {
  transform: translateY(-1px);
  box-shadow: var(--brand-glow), 0 0 0 4px color-mix(in srgb, var(--primary-color) 14%, transparent);
}

.nf-btn--ghost {
  color: var(--text-color-primary);
  background: var(--card-bg-solid, var(--bg-color-primary));
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-xs);
}

.nf-btn--ghost:hover {
  transform: translateY(-1px);
  border-color: var(--border-color-hover);
  background: var(--bg-color-hover);
}

.nf-back {
  margin-top: 18px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  border: none;
  border-radius: var(--border-radius-sm);
  background: transparent;
  color: var(--text-color-tertiary);
  font-size: 14px;
  cursor: pointer;
  transition: color var(--duration-fast, 0.15s) var(--ease-out, ease-out);
}

.nf-back:hover {
  color: var(--accent-text-color, var(--primary-color));
}

@media (prefers-reduced-motion: reduce) {
  .nf-btn,
  .nf-back {
    transition: none;
  }

  .nf-btn:hover {
    transform: none;
  }
}

@media (max-width: 640px) {
  .nf-topbar {
    padding: 16px 18px;
  }

  .nf-glow--a {
    width: 360px;
    height: 360px;
  }

  .nf-message {
    font-size: 15px;
  }

  .nf-actions {
    flex-direction: column;
    width: 100%;
    max-width: 320px;
  }

  .nf-btn {
    width: 100%;
  }
}
</style>
