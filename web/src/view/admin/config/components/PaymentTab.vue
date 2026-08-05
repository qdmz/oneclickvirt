<template>
  <el-form
    v-loading="loading"
    :model="config"
    label-width="160px"
    class="config-form"
  >
    <el-alert
      :title="$t('admin.config.payment.desc')"
      type="info"
      :closable="false"
      show-icon
      style="margin-bottom: 20px;"
    >
      <div>{{ $t('admin.config.payment.hint') }}</div>
    </el-alert>

    <el-form-item :label="$t('admin.config.payment.enabled')">
      <el-switch v-model="config.payment.epayEnabled" />
      <div class="form-tip">
        {{ $t('admin.config.payment.enabledTip') }}
      </div>
    </el-form-item>

    <el-divider content-position="left">
      {{ $t('admin.config.payment.merchantSection') }}
    </el-divider>

    <el-form-item :label="$t('admin.config.payment.pid')">
      <el-input
        v-model="config.payment.epayPid"
        :placeholder="$t('admin.config.payment.pidPh')"
        style="max-width: 460px;"
      />
    </el-form-item>

    <el-form-item :label="$t('admin.config.payment.key')">
      <el-input
        v-model="config.payment.epayKey"
        type="password"
        show-password
        :placeholder="$t('admin.config.payment.keyPh')"
        style="max-width: 460px;"
      />
      <div class="form-tip">
        {{ $t('admin.config.payment.keyTip') }}
      </div>
    </el-form-item>

    <el-form-item :label="$t('admin.config.payment.apiUrl')">
      <el-input
        v-model="config.payment.epayApiUrl"
        :placeholder="$t('admin.config.payment.apiUrlPh')"
        style="max-width: 460px;"
      />
      <div class="form-tip">
        {{ $t('admin.config.payment.apiUrlTip') }}
      </div>
    </el-form-item>

    <el-divider content-position="left">
      {{ $t('admin.config.payment.callbackSection') }}
    </el-divider>

    <el-form-item :label="$t('admin.config.payment.notifyUrl')">
      <el-input
        v-model="config.payment.epayNotifyUrl"
        :placeholder="defaultNotifyUrl"
        style="max-width: 460px;"
      />
      <div class="form-tip">
        {{ $t('admin.config.payment.notifyUrlTip', { url: defaultNotifyUrl }) }}
      </div>
    </el-form-item>

    <el-form-item :label="$t('admin.config.payment.returnUrl')">
      <el-input
        v-model="config.payment.epayReturnUrl"
        :placeholder="defaultReturnUrl"
        style="max-width: 460px;"
      />
      <div class="form-tip">
        {{ $t('admin.config.payment.returnUrlTip', { url: defaultReturnUrl }) }}
      </div>
    </el-form-item>

    <el-alert
      v-if="config.payment.epayEnabled && !isConfigComplete"
      :title="$t('admin.config.payment.incomplete')"
      type="warning"
      :closable="false"
      show-icon
      style="margin-top: 12px;"
    />
  </el-form>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  config: { type: Object, required: true },
  loading: { type: Boolean, default: false }
})

const origin = `${window.location.protocol}//${window.location.host}`
const defaultNotifyUrl = `${origin}/api/v1/payment/epay/notify`
const defaultReturnUrl = `${origin}/api/v1/payment/epay/return`

const isConfigComplete = computed(() => {
  const p = props.config.payment || {}
  return !!(p.epayPid && p.epayKey && p.epayApiUrl)
})
</script>

<style scoped>
.config-form {
  max-width: 900px;
}
.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
  margin-top: 4px;
  word-break: break-all;
}
</style>
