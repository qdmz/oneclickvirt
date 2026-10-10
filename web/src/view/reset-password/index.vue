<template>
  <AuthShell>
    <div class="reset-password">
      <h2>{{ t('resetPassword.title') }}</h2>

      <!-- 验证中 -->
      <div
        v-if="state === 'loading'"
        class="reset-state"
      >
        <div class="reset-spinner">
          <el-icon class="is-loading">
            <Loading />
          </el-icon>
        </div>
        <p>{{ t('ui.resetPassword.verifying') }}</p>
      </div>

      <!-- 成功 -->
      <el-result
        v-else-if="state === 'success'"
        icon="success"
        :title="t('ui.resetPassword.successTitle')"
        :sub-title="t('ui.resetPassword.successDesc')"
      >
        <template #extra>
          <el-button
            type="primary"
            class="reset-action"
            @click="router.push('/login')"
          >
            {{ t('forgotPassword.backToLogin') }}
          </el-button>
        </template>
      </el-result>

      <!-- 失败 / 缺少 token -->
      <el-result
        v-else
        icon="error"
        :title="t('ui.resetPassword.invalidTitle')"
        :sub-title="errorMessage"
      >
        <template #extra>
          <div class="reset-actions">
            <el-button
              type="primary"
              class="reset-action"
              @click="router.push('/forgot-password')"
            >
              {{ t('ui.resetPassword.requestAgain') }}
            </el-button>
            <el-button
              v-if="token"
              class="reset-action"
              @click="doReset"
            >
              {{ t('ui.common.retry') }}
            </el-button>
          </div>
        </template>
      </el-result>
    </div>
  </AuthShell>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Loading } from '@element-plus/icons-vue'
import { resetPassword } from '@/api/auth'
import AuthShell from '@/components/AuthShell.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

// 兼容 /#/reset-password?token=xxx，也兼容旧邮件链接 /reset-password?token=xxx（token 在 location.search 中）
const token = computed(() => {
  const q = route.query.token
  if (typeof q === 'string' && q) return q
  const params = new URLSearchParams(window.location.search)
  return params.get('token') || ''
})

const state = ref('loading')
const errorMessage = ref('')

const doReset = async () => {
  if (!token.value) {
    state.value = 'error'
    errorMessage.value = t('ui.resetPassword.missingToken')
    return
  }
  state.value = 'loading'
  try {
    const res = await resetPassword({ token: token.value })
    if (res && res.code === 200) {
      state.value = 'success'
    } else {
      state.value = 'error'
      errorMessage.value = res?.msg || res?.message || t('ui.resetPassword.invalidDesc')
    }
  } catch (error) {
    state.value = 'error'
    errorMessage.value = error?.response?.data?.msg || error?.message || t('ui.resetPassword.invalidDesc')
  }
}

onMounted(doReset)
</script>

<style scoped>
.reset-password :deep(.el-result) {
  padding: 24px 0 0;
}

.reset-password :deep(.el-result__subtitle p) {
  line-height: 1.7;
}

.reset-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 48px 0 24px;
  color: var(--text-color-secondary);
}

.reset-spinner {
  display: grid;
  place-items: center;
  width: 56px;
  height: 56px;
  border-radius: 16px;
  font-size: 26px;
  color: var(--primary-color);
  background: var(--primary-color-bg);
}

.reset-actions {
  display: flex;
  gap: 10px;
  justify-content: center;
  flex-wrap: wrap;
}

.reset-action {
  min-width: 140px;
}
</style>
