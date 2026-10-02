<template>
  <el-form
    :model="modelValue"
    label-width="120px"
    class="server-form"
  >
    <!-- 并发控制设置 -->
    <el-divider content-position="left">
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.concurrencyControl') }}</span>
    </el-divider>
    
    <el-form-item
      :label="$t('admin.providers.allowConcurrentTasks')"
      prop="allowConcurrentTasks"
    >
      <el-switch
        v-model="modelValue.allowConcurrentTasks"
        :active-text="$t('common.yes')"
        :inactive-text="$t('common.no')"
      />
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.allowConcurrentTasksTip') }}
      </el-text>
    </div>

    <el-form-item
      v-if="modelValue.allowConcurrentTasks"
      :label="$t('admin.providers.maxConcurrentTasks')"
      prop="maxConcurrentTasks"
    >
      <el-input-number
        v-model="modelValue.maxConcurrentTasks"
        :min="1"
        :max="10"
        :step="1"
        :controls="false"
        placeholder="1"
        style="width: 200px"
      />
    </el-form-item>
    <div
      v-if="modelValue.allowConcurrentTasks"
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.maxConcurrentTasksTip') }}
      </el-text>
    </div>

    <!-- 任务轮询设置 -->
    <el-divider content-position="left">
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.taskPollingSettings') }}</span>
    </el-divider>
    
    <el-form-item
      :label="$t('admin.providers.enableTaskPolling')"
      prop="enableTaskPolling"
    >
      <el-switch
        v-model="modelValue.enableTaskPolling"
        :active-text="$t('common.yes')"
        :inactive-text="$t('common.no')"
      />
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.enableTaskPollingTip') }}
      </el-text>
    </div>

    <el-form-item
      v-if="modelValue.enableTaskPolling"
      :label="$t('admin.providers.taskPollInterval')"
      prop="taskPollInterval"
    >
      <el-input-number
        v-model="modelValue.taskPollInterval"
        :min="5"
        :max="300"
        :step="5"
        :controls="false"
        placeholder="60"
        style="width: 200px"
      />
      <span style="margin-left: 10px; color: #666;">{{ $t('common.seconds') }}</span>
    </el-form-item>
    <div
      v-if="modelValue.enableTaskPolling"
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.taskPollIntervalTip') }}
      </el-text>
    </div>

    <!-- 操作执行规则设置 -->
    <el-divider content-position="left">
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.executionRules') }}</span>
    </el-divider>
    
    <el-form-item
      :label="$t('admin.providers.executionRule')"
      prop="executionRule"
    >
      <el-select
        v-model="modelValue.executionRule"
        :placeholder="$t('admin.providers.executionRulePlaceholder')"
        style="width: 200px"
      >
        <el-option
          :label="$t('admin.providers.executionRuleAuto')"
          value="auto"
        >
          <span>{{ $t('admin.providers.executionRuleAuto') }}</span>
          <span style="float: right; color: #8492a6; font-size: 12px;">{{ $t('admin.providers.executionRuleAutoTip') }}</span>
        </el-option>
        <el-option
          :label="$t('admin.providers.executionRuleAPIOnly')"
          value="api_only"
        >
          <span>{{ $t('admin.providers.executionRuleAPIOnly') }}</span>
          <span style="float: right; color: #8492a6; font-size: 12px;">{{ $t('admin.providers.executionRuleAPIOnlyTip') }}</span>
        </el-option>
        <el-option
          :label="$t('admin.providers.executionRuleSSHOnly')"
          value="ssh_only"
        >
          <span>{{ $t('admin.providers.executionRuleSSHOnly') }}</span>
          <span style="float: right; color: #8492a6; font-size: 12px;">{{ $t('admin.providers.executionRuleSSHOnlyTip') }}</span>
        </el-option>
      </el-select>
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.executionRuleTip') }}
      </el-text>
    </div>

    <!-- 申请领取控制 -->
    <el-divider content-position="left">
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.applicationControl') }}</span>
    </el-divider>

    <el-form-item
      :label="$t('admin.providers.redeemCodeOnly')"
      prop="redeemCodeOnly"
    >
      <el-switch
        v-model="modelValue.redeemCodeOnly"
        :active-text="$t('common.yes')"
        :inactive-text="$t('common.no')"
      />
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.redeemCodeOnlyTip') }}
      </el-text>
    </div>

    <!-- 硬件监控 -->
    <el-divider content-position="left">
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.hardwareMonitoring') }}</span>
    </el-divider>

    <el-form-item
      :label="$t('admin.providers.enableResourceMonitoring')"
      prop="enableResourceMonitoring"
    >
      <el-switch
        v-model="modelValue.enableResourceMonitoring"
        :active-text="$t('admin.providers.enabled')"
        :inactive-text="$t('admin.providers.disabled')"
      />
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.enableResourceMonitoringTip') }}
      </el-text>
    </div>

    <!-- 生命周期与流量处置策略 -->
    <el-divider content-position="left">
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.lifecyclePolicy') }}</span>
    </el-divider>

    <el-form-item
      :label="$t('admin.providers.trafficOverLimitAction')"
      prop="trafficOverLimitAction"
    >
      <el-select
        v-model="modelValue.trafficOverLimitAction"
        :placeholder="$t('admin.providers.trafficOverLimitActionPlaceholder')"
        style="width: 260px"
      >
        <el-option
          :label="$t('admin.providers.trafficActionStop')"
          value="stop"
        />
        <el-option
          :label="$t('admin.providers.trafficActionSpeedLimit')"
          value="speed_limit"
        />
        <el-option
          :label="$t('admin.providers.trafficActionFreeze')"
          value="freeze"
        />
        <el-option
          :label="$t('admin.providers.trafficActionMarkOnly')"
          value="mark_only"
        />
      </el-select>
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.trafficOverLimitActionTip') }}
      </el-text>
    </div>

    <el-form-item
      v-if="modelValue.trafficOverLimitAction === 'speed_limit'"
      :label="$t('admin.providers.trafficSpeedLimitKbps')"
      prop="trafficSpeedLimitKbps"
    >
      <el-input-number
        v-model="modelValue.trafficSpeedLimitKbps"
        :min="1"
        :max="1048576"
        :step="128"
        :controls="false"
        style="width: 200px"
      />
      <span style="margin-left: 10px; color: #666;">Kbps</span>
    </el-form-item>

    <el-form-item
      :label="$t('admin.providers.trafficQuotaVisible')"
      prop="trafficQuotaVisible"
    >
      <el-switch
        v-model="modelValue.trafficQuotaVisible"
        :active-text="$t('common.yes')"
        :inactive-text="$t('common.no')"
      />
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.trafficQuotaVisibleTip') }}
      </el-text>
    </div>

    <!-- WebVNC 设置 -->
    <el-divider content-position="left">
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.webVncSettings') }}</span>
    </el-divider>

    <el-form-item
      :label="$t('admin.providers.enableVNC')"
      prop="enableVNC"
    >
      <el-switch
        v-model="modelValue.enableVNC"
        :active-text="$t('common.yes')"
        :inactive-text="$t('common.no')"
      />
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.enableVNCTip') }}
      </el-text>
    </div>

    <template v-if="modelValue.enableVNC">
      <el-form-item
        :label="$t('admin.providers.vncBasePort')"
        prop="vncBasePort"
      >
        <el-input-number
          v-model="modelValue.vncBasePort"
          :min="1"
          :max="65535"
          :step="1"
          :controls="false"
          placeholder="5900"
          style="width: 200px"
        />
      </el-form-item>
      <div
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.vncBasePortTip') }}
        </el-text>
      </div>

      <el-form-item
        :label="$t('admin.providers.vncHost')"
        prop="vncHost"
      >
        <el-input
          v-model="modelValue.vncHost"
          :placeholder="$t('admin.providers.vncHostPlaceholder')"
          clearable
          style="width: 400px"
        />
      </el-form-item>
      <div
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.vncHostTip') }}
        </el-text>
      </div>
    </template>

    <!-- 域名反向代理设置 -->
    <el-divider content-position="left">
      <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.domainProxy') }}</span>
    </el-divider>

    <el-form-item
      :label="$t('admin.providers.enableDomainBinding')"
      prop="enableDomainBinding"
    >
      <el-switch
        v-model="modelValue.enableDomainBinding"
        :active-text="$t('admin.providers.enabled')"
        :inactive-text="$t('admin.providers.disabled')"
      />
    </el-form-item>
    <div
      class="form-tip"
      style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
    >
      <el-text
        size="small"
        type="info"
      >
        {{ $t('admin.providers.enableDomainBindingTip') }}
      </el-text>
    </div>

    <template v-if="modelValue.enableDomainBinding">
      <el-form-item
        :label="$t('admin.providers.proxyEnableHTTP')"
        prop="proxyEnableHttp"
      >
        <el-switch
          v-model="modelValue.proxyEnableHttp"
          :active-text="$t('common.yes')"
          :inactive-text="$t('common.no')"
        />
      </el-form-item>
      <div
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.proxyEnableHTTPTip') }}
        </el-text>
      </div>

      <el-form-item
        v-if="modelValue.proxyEnableHttp"
        :label="$t('admin.providers.proxyHTTPPort')"
        prop="proxyHttpPort"
      >
        <el-input-number
          v-model="modelValue.proxyHttpPort"
          :min="1"
          :max="65535"
          :step="1"
          :controls="false"
          placeholder="80"
          style="width: 200px"
        />
      </el-form-item>
      <div
        v-if="modelValue.proxyEnableHttp"
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.proxyHTTPPortTip') }}
        </el-text>
      </div>

      <el-form-item
        :label="$t('admin.providers.proxyEnableHTTPS')"
        prop="proxyEnableHttps"
      >
        <el-switch
          v-model="modelValue.proxyEnableHttps"
          :active-text="$t('common.yes')"
          :inactive-text="$t('common.no')"
        />
      </el-form-item>
      <div
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.proxyEnableHTTPSTip') }}
        </el-text>
      </div>

      <template v-if="modelValue.proxyEnableHttps">
        <el-form-item
          :label="$t('admin.providers.proxyHTTPSPort')"
          prop="proxyHttpsPort"
        >
          <el-input-number
            v-model="modelValue.proxyHttpsPort"
            :min="1"
            :max="65535"
            :step="1"
            :controls="false"
            placeholder="443"
            style="width: 200px"
          />
        </el-form-item>
        <div
          class="form-tip"
          style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
        >
          <el-text
            size="small"
            type="info"
          >
            {{ $t('admin.providers.proxyHTTPSPortTip') }}
          </el-text>
        </div>

        <el-form-item
          :label="$t('admin.providers.proxyTLSCertPath')"
          prop="proxyTlsCertPath"
        >
          <el-input
            v-model="modelValue.proxyTlsCertPath"
            :placeholder="$t('admin.providers.proxyTLSCertPathPlaceholder')"
            clearable
            style="width: 400px"
          />
        </el-form-item>
        <div
          class="form-tip"
          style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
        >
          <el-text
            size="small"
            type="info"
          >
            {{ $t('admin.providers.proxyTLSCertPathTip') }}
          </el-text>
        </div>

        <el-form-item
          :label="$t('admin.providers.proxyTLSKeyPath')"
          prop="proxyTlsKeyPath"
        >
          <el-input
            v-model="modelValue.proxyTlsKeyPath"
            :placeholder="$t('admin.providers.proxyTLSKeyPathPlaceholder')"
            clearable
            style="width: 400px"
          />
        </el-form-item>
        <div
          class="form-tip"
          style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
        >
          <el-text
            size="small"
            type="info"
          >
            {{ $t('admin.providers.proxyTLSKeyPathTip') }}
          </el-text>
        </div>
      </template>
    </template>
  </el-form>
</template>

<script setup>
defineProps({
  modelValue: {
    type: Object,
    required: true
  }
})
</script>

<style scoped>
.server-form {
  max-height: 500px;
  overflow-y: auto;
  padding-right: 10px;
}

.form-tip {
  margin-top: 5px;
}


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
