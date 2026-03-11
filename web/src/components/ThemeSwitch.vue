<template>
  <button
    class="theme-switch"
    :title="isDark ? t('theme.lightMode') : t('theme.darkMode')"
    @click="toggleTheme"
  >
    <transition name="theme-icon" mode="out-in">
      <el-icon v-if="isDark" key="sun" :size="18">
        <Sunny />
      </el-icon>
      <el-icon v-else key="moon" :size="18">
        <Moon />
      </el-icon>
    </transition>
  </button>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Sunny, Moon } from '@element-plus/icons-vue'

const { t } = useI18n()
const isDark = ref(true)

const toggleTheme = () => {
  isDark.value = !isDark.value
  applyTheme()
}

const applyTheme = () => {
  document.documentElement.setAttribute('data-theme', isDark.value ? 'dark' : 'light')
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

onMounted(() => {
  const saved = localStorage.getItem('theme')
  if (saved) {
    isDark.value = saved === 'dark'
  }
  applyTheme()
})
</script>

<style scoped>
.theme-switch {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: var(--border-radius-sm);
  border: 1px solid var(--border-color);
  background: var(--bg-color-secondary);
  color: var(--text-color-secondary);
  cursor: pointer;
  transition: var(--transition-normal);
  padding: 0;
  flex-shrink: 0;
}

.theme-switch:hover {
  color: var(--primary-color);
  border-color: var(--primary-color);
  background: var(--primary-color-bg);
  transform: rotate(15deg);
}

.theme-icon-enter-active,
.theme-icon-leave-active {
  transition: all 0.2s ease;
}

.theme-icon-enter-from {
  opacity: 0;
  transform: rotate(-90deg) scale(0.5);
}

.theme-icon-leave-to {
  opacity: 0;
  transform: rotate(90deg) scale(0.5);
}
</style>
