<template>
  <div class="connection-tab">
    <!-- ======================================================
         SSH 模式内容（connectionType === 'ssh'）
         ====================================================== -->
    <el-form
      v-if="showSSHSettings"
      :model="modelValue"
      label-width="120px"
      class="server-form"
    >
      <div
        v-if="isAgentMode"
        class="form-tip"
        style="margin-top: -4px; margin-bottom: 12px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.agentMappedSshOptionalTip') }}
        </el-text>
      </div>

      <el-form-item
        :label="$t('admin.providers.username')"
        prop="username"
      >
        <el-input
          v-model="modelValue.username"
          :placeholder="$t('admin.providers.usernamePlaceholder')"
        />
      </el-form-item>

      <!-- 认证方式选择 -->
      <el-form-item
        :label="$t('admin.providers.authMethod')"
        prop="authMethod"
      >
        <el-radio-group
          v-model="modelValue.authMethod"
          @change="emit('auth-method-change', $event)"
        >
          <el-radio-button label="password">
            {{ $t('admin.providers.usePassword') }}
          </el-radio-button>
          <el-radio-button label="sshKey">
            {{ $t('admin.providers.useSSHKey') }}
          </el-radio-button>
        </el-radio-group>
      </el-form-item>

      <!-- 密码认证 -->
      <el-form-item
        v-if="modelValue.authMethod === 'password'"
        :label="$t('admin.providers.password')"
        prop="password"
      >
        <el-input
          v-model="modelValue.password"
          type="password"
          :placeholder="isEditing ? $t('admin.providers.passwordEditPlaceholder') : $t('admin.providers.passwordPlaceholder')"
          show-password
        />
        <div
          v-if="isEditing"
          class="form-tip"
        >
          <el-text
            size="small"
            type="info"
          >
            {{ $t('admin.providers.passwordKeepTip') }}
          </el-text>
        </div>
      </el-form-item>

      <!-- SSH密钥认证 -->
      <el-form-item
        v-if="modelValue.authMethod === 'sshKey'"
        :label="$t('admin.providers.sshKey')"
        prop="sshKey"
      >
        <el-input
          v-model="modelValue.sshKey"
          type="textarea"
          :rows="4"
          :placeholder="isEditing ? $t('admin.providers.sshKeyEditPlaceholder') : $t('admin.providers.sshKeyPlaceholder')"
        />
        <div
          v-if="isEditing"
          class="form-tip"
        >
          <el-text
            size="small"
            type="info"
          >
            {{ $t('admin.providers.sshKeyEditTip') }}
          </el-text>
        </div>
      </el-form-item>

      <el-divider content-position="left">
        {{ $t('admin.providers.sshTimeoutConfig') }}
      </el-divider>

      <el-form-item
        :label="$t('admin.providers.connectTimeout')"
        prop="sshConnectTimeout"
      >
        <el-input-number
          v-model="modelValue.sshConnectTimeout"
          :min="5"
          :max="300"
          :step="5"
          :controls="false"
          placeholder="30"
        />
        <span style="margin-left: 10px;">{{ $t('admin.providers.seconds') }}</span>
      </el-form-item>
      <div
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.connectTimeoutTip') }}
        </el-text>
      </div>

      <el-form-item
        :label="$t('admin.providers.executeTimeout')"
        prop="sshExecuteTimeout"
      >
        <el-input-number
          v-model="modelValue.sshExecuteTimeout"
          :min="30"
          :max="3600"
          :step="30"
          :controls="false"
          placeholder="300"
        />
        <span style="margin-left: 10px;">{{ $t('admin.providers.seconds') }}</span>
      </el-form-item>
      <div
        class="form-tip"
        style="margin-top: -10px; margin-bottom: 15px; margin-left: 120px;"
      >
        <el-text
          size="small"
          type="info"
        >
          {{ $t('admin.providers.executeTimeoutTip') }}
        </el-text>
      </div>

      <el-form-item :label="$t('admin.providers.connectionTest')">
        <el-button
          type="primary"
          :loading="testingConnection"
          :disabled="!modelValue.host || !modelValue.username || (modelValue.authMethod === 'password' ? !modelValue.password : !modelValue.sshKey)"
          @click="emit('test-connection')"
        >
          <el-icon v-if="!testingConnection">
            <Connection />
          </el-icon>
          {{ testingConnection ? $t('admin.providers.testing') : $t('admin.providers.testSSH') }}
        </el-button>
        <div
          v-if="connectionTestResult"
          class="form-tip"
          style="margin-top: 10px;"
        >
          <el-alert
            :title="connectionTestResult.title"
            :type="connectionTestResult.type"
            :closable="false"
            show-icon
          >
            <template v-if="connectionTestResult.success">
              <div style="margin-top: 8px;">
                <p><strong>{{ $t('admin.providers.testResults') }}:</strong></p>
                <p>{{ $t('admin.providers.minLatency') }}: {{ connectionTestResult.minLatency }}ms</p>
                <p>{{ $t('admin.providers.maxLatency') }}: {{ connectionTestResult.maxLatency }}ms</p>
                <p>{{ $t('admin.providers.avgLatency') }}: {{ connectionTestResult.avgLatency }}ms</p>
                <p style="margin-top: 8px;">
                  <strong>{{ $t('admin.providers.recommendedTimeout') }}: {{ connectionTestResult.recommendedTimeout }}{{ $t('common.seconds') }}</strong>
                </p>
                <el-button
                  type="primary"
                  size="small"
                  style="margin-top: 8px;"
                  @click="emit('apply-timeout')"
                >
                  {{ $t('admin.providers.applyRecommended') }}
                </el-button>
              </div>
            </template>
            <template v-else>
              <p>{{ connectionTestResult.error }}</p>
            </template>
          </el-alert>
        </div>
      </el-form-item>
    </el-form>

    <!-- ======================================================
         本机模式内容（connectionType === 'local'）
         ====================================================== -->
    <div
      v-if="modelValue.connectionType === 'local'"
      class="local-mode-content"
    >
      <el-alert
        type="success"
        :closable="false"
        show-icon
      >
        <template #title>
          {{ $t('admin.providers.localConnection') }}
        </template>
        <div>{{ $t('admin.providers.localConnectionTip') }}</div>
      </el-alert>
      <div class="local-detect-actions">
        <el-button
          type="primary"
          :loading="detectingLocal"
          @click="handleDetectLocalProvider"
        >
          <el-icon v-if="!detectingLocal">
            <CircleCheck />
          </el-icon>
          {{ detectingLocal ? $t('admin.providers.localDetecting') : $t('admin.providers.localDetect') }}
        </el-button>
      </div>
      <el-alert
        v-if="localDetectionResult"
        class="local-detect-result"
        :type="localDetectionResult.available ? 'success' : 'warning'"
        :closable="false"
        show-icon
      >
        <template #title>
          {{ localDetectionResult.available ? $t('admin.providers.localDetectSuccess') : $t('admin.providers.localDetectFailed') }}
        </template>
        <div class="local-detect-summary">
          <span>{{ $t('admin.providers.localDetectKvm') }}: <el-tag
            size="small"
            :type="localDetectionResult.kvmAvailable ? 'success' : 'warning'"
          >{{ formatLocalDetectStatus(localDetectionResult.kvmAvailable) }}</el-tag></span>
          <span>{{ $t('admin.providers.localDetectQemu') }}: <el-tag
            size="small"
            :type="localDetectionResult.qemuAvailable ? 'success' : 'danger'"
          >{{ formatLocalDetectStatus(localDetectionResult.qemuAvailable) }}</el-tag></span>
          <span>{{ $t('admin.providers.localDetectLxc') }}: <el-tag
            size="small"
            :type="localDetectionResult.lxcAvailable ? 'success' : 'warning'"
          >{{ formatLocalDetectStatus(localDetectionResult.lxcAvailable) }}</el-tag></span>
        </div>
        <div
          v-if="localCommandChecks.length"
          class="local-command-list"
        >
          <strong>{{ $t('admin.providers.localDetectCommands') }}</strong>
          <div
            v-for="command in localCommandChecks"
            :key="command.name"
            class="local-command-row"
          >
            <span>{{ command.name }}</span>
            <el-tag
              size="small"
              :type="command.present ? 'success' : 'danger'"
            >
              {{ formatLocalDetectStatus(command.present) }}
            </el-tag>
            <code v-if="command.path">{{ command.path }}</code>
          </div>
        </div>
        <div
          v-if="localDetectionResult.warnings?.length"
          class="local-warning-list"
        >
          <strong>{{ $t('admin.providers.localDetectWarnings') }}</strong>
          <ul>
            <li
              v-for="warning in localDetectionResult.warnings"
              :key="warning"
            >
              {{ warning }}
            </li>
          </ul>
        </div>
      </el-alert>
    </div>

    <!-- ======================================================
         Agent 模式内容（connectionType === 'agent'）
         ====================================================== -->
    <div
      v-if="modelValue.connectionType === 'agent'"
      class="agent-mode-content"
    >
      <!-- Agent 状态（编辑模式） -->
      <el-alert
        v-if="isEditing"
        :type="agentAlertType"
        :closable="false"
        style="margin-bottom: 20px;"
      >
        <template #title>
          <span>
            {{ $t('admin.providers.agentStatus') }}:
            <strong>{{ agentStatusLabel }}</strong>
          </span>
          <span
            v-if="modelValue.agentConnectedAt && effectiveAgentStatus === 'online'"
            style="margin-left: 16px; font-size: 12px; opacity: 0.8;"
          >
            {{ $t('admin.providers.agentOnlineDuration') }}: {{ formatOnlineDuration(modelValue.agentConnectedAt) }}
          </span>
          <span
            v-if="modelValue.agentLastSeen"
            style="margin-left: 16px; font-size: 12px; opacity: 0.8;"
          >
            {{ $t('admin.providers.agentLastSeen') }}: {{ formatDateTime(modelValue.agentLastSeen) }}
          </span>
          <span
            v-if="modelValue.agentControlLastSeen"
            style="margin-left: 16px; font-size: 12px; opacity: 0.8;"
          >
            {{ $t('admin.providers.agentControlLastSeen') }}: {{ formatDateTime(modelValue.agentControlLastSeen) }}
          </span>
          <span
            v-if="modelValue.agentExecLastSeen"
            style="margin-left: 16px; font-size: 12px; opacity: 0.8;"
          >
            {{ $t('admin.providers.agentExecLastSeen') }}: {{ formatDateTime(modelValue.agentExecLastSeen) }}
          </span>
          <span
            v-if="modelValue.agentRemoteIP"
            style="margin-left: 16px; font-size: 12px; opacity: 0.8;"
          >
            {{ $t('admin.providers.agentRemoteIP') }}: {{ modelValue.agentRemoteIP }}
          </span>
        </template>
      </el-alert>

      <!-- 新增模式：只有步骤说明，命令在保存后生成 -->
      <el-alert
        v-if="!isEditing"
        type="info"
        :closable="false"
        style="margin-bottom: 20px;"
      >
        <template #title>
          {{ $t('admin.providers.agentModeNewHint') }}
        </template>
        <div style="margin-top: 8px; line-height: 1.8; font-size: 13px;">
          <p>① {{ $t('admin.providers.agentStep1') }}</p>
          <p>② {{ $t('admin.providers.agentStep2') }}</p>
          <p>③ {{ $t('admin.providers.agentStep3') }}</p>
        </div>
      </el-alert>

      <!-- 编辑模式：Agent密钥生成 + 安装命令 -->
      <template v-if="isEditing">
        <el-divider content-position="left">
          <span style="font-size: 14px; color: #666;">{{ $t('admin.providers.agentInstallSection') }}</span>
        </el-divider>

        <div style="margin-bottom: 16px;">
          <el-button
            type="primary"
            :loading="generatingSecret"
            @click="emit('generate-agent-secret')"
          >
            {{ agentConnectCmd ? $t('admin.providers.regenerateAgentSecret') : $t('admin.providers.generateAgentSecret') }}
          </el-button>
          <div
            class="form-tip"
            style="margin-top: 6px;"
          >
            <el-text
              size="small"
              type="info"
            >
              {{ $t('admin.providers.generateAgentSecretTip') }}
            </el-text>
          </div>
        </div>

        <!-- 安装命令 -->
        <div
          v-if="agentConnectCmd"
          class="install-cmd-box"
        >
          <div class="install-cmd-header">
            <span>{{ $t('admin.providers.agentCmdInstall') }}</span>
            <div style="display:flex;align-items:center;gap:8px;">
              <el-switch
                v-model="useControllerSource"
                size="small"
                :active-text="$t('admin.providers.controllerSource')"
                :inactive-text="$t('admin.providers.githubSource')"
                style="--el-switch-on-color: #13ce66;"
              />
              <el-switch
                v-model="useWSS"
                size="small"
                :disabled="wssUnavailable"
                :active-text="$t('admin.providers.wssSecure')"
                :inactive-text="$t('admin.providers.wsPlain')"
                style="--el-switch-on-color: #13ce66;"
              />
              <el-switch
                v-model="useCDN"
                size="small"
                :disabled="useControllerSource"
                :active-text="$t('admin.providers.cdnAccel')"
                :inactive-text="$t('admin.providers.cdnDirect')"
                style="--el-switch-on-color: #13ce66;"
              />
              <el-button
                size="small"
                @click="copyCmd(installCmdDisplay)"
              >
                {{ $t('common.copy') }}
              </el-button>
            </div>
          </div>
          <!-- wss unavailable warning -->
          <el-alert
            v-if="wssUnavailable"
            :title="$t('admin.providers.wssUnavailable')"
            type="warning"
            :closable="false"
            show-icon
            style="margin: 8px 14px 0;"
          />
          <div class="install-cmd-content">
            {{ installCmdDisplay }}
          </div>
        </div>

        <!-- 卸载命令 -->
        <div
          v-if="agentConnectCmd"
          class="install-cmd-box"
          style="margin-top: 12px;"
        >
          <div class="install-cmd-header">
            <span>{{ $t('admin.providers.agentCmdUninstall') }}</span>
            <el-button
              size="small"
              @click="copyCmd('ocv uninstall')"
            >
              {{ $t('common.copy') }}
            </el-button>
          </div>
          <div class="install-cmd-content">
            ocv uninstall
          </div>
        </div>

        <!-- 升级命令 -->
        <div
          v-if="agentConnectCmd"
          class="install-cmd-box"
          style="margin-top: 12px;"
        >
          <div class="install-cmd-header">
            <span>{{ $t('admin.providers.agentCmdUpgrade') }}</span>
            <el-button
              size="small"
              @click="copyCmd('ocv upgrade')"
            >
              {{ $t('common.copy') }}
            </el-button>
          </div>
          <div class="install-cmd-content">
            ocv upgrade
          </div>
        </div>

        <!-- ocv 快捷命令 -->
        <div
          v-if="agentConnectCmd"
          class="install-cmd-box"
          style="margin-top: 12px;"
        >
          <div class="install-cmd-header">
            <span>{{ $t('admin.providers.agentCmdOcv') }}</span>
            <el-button
              size="small"
              @click="copyCmd('ocv')"
            >
              {{ $t('common.copy') }}
            </el-button>
          </div>
          <div class="install-cmd-content">
            ocv
          </div>
          <div class="install-cmd-tip">
            <el-icon><InfoFilled /></el-icon>
            {{ $t('admin.providers.agentCmdOcvTip') }}
          </div>
        </div>

        <div
          v-if="agentConnectCmd"
          class="form-tip"
          style="margin-top: 10px;"
        >
          <el-text
            size="small"
            type="info"
          >
            {{ $t('admin.providers.agentInstallNote') }}
          </el-text>
        </div>

        <!-- 检测连接 -->
        <div
          v-if="agentConnectCmd"
          style="margin-top: 16px;"
        >
          <el-button
            type="success"
            :loading="checkingAgentStatus"
            @click="emit('check-agent-status')"
          >
            <el-icon><CircleCheck /></el-icon>
            {{ $t('admin.providers.checkAgentConnection') }}
          </el-button>
          <span style="margin-left: 12px; font-size: 13px; color: var(--el-text-color-secondary);">
            {{ $t('admin.providers.checkAgentConnectionTip') }}
          </span>
        </div>

        <!-- Web 终端：Agent online 时显示 -->
        <template v-if="modelValue.agentStatus === 'online'">
          <el-divider content-position="left">
            <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.webTerminal') }}</span>
          </el-divider>
          <el-form
            :model="modelValue"
            label-width="120px"
            class="server-form"
          >
            <el-form-item :label="$t('admin.providers.execCommand')">
              <el-input
                v-model="localCommand"
                :placeholder="$t('admin.providers.execCommandPlaceholder')"
                @keyup.enter="emit('exec-command', localCommand)"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                :loading="execLoading"
                :disabled="!localCommand.trim()"
                @click="emit('exec-command', localCommand)"
              >
                {{ $t('admin.providers.execRun') }}
              </el-button>
              <el-button
                v-if="execResult"
                @click="localCommand = ''; emit('clear-exec-result')"
              >
                {{ $t('common.clear') }}
              </el-button>
            </el-form-item>
            <el-form-item
              v-if="execResult !== null"
              :label="$t('admin.providers.execResult')"
            >
              <div class="exec-output">
                <div v-if="execResult.stdout">
                  {{ execResult.stdout }}
                </div>
                <div
                  v-if="execResult.stderr"
                  style="color: #f48771;"
                >
                  {{ execResult.stderr }}
                </div>
                <div
                  v-if="!execResult.stdout && !execResult.stderr"
                  style="color: #888;"
                >
                  {{ $t('admin.providers.execNoOutput') }}
                </div>
              </div>
            </el-form-item>
          </el-form>
        </template>
      </template>
    </div>

    <!-- ======================================================
         SSH 模式 Web 终端（编辑时显示）
         ====================================================== -->
    <template v-if="isEditing && modelValue.connectionType !== 'agent'">
      <el-divider content-position="left">
        <span style="color: #666; font-size: 14px;">{{ $t('admin.providers.webTerminal') }}</span>
      </el-divider>
      <el-form
        :model="modelValue"
        label-width="120px"
        class="server-form"
      >
        <el-form-item :label="$t('admin.providers.execCommand')">
          <el-input
            v-model="localCommand"
            :placeholder="$t('admin.providers.execCommandPlaceholder')"
            @keyup.enter="emit('exec-command', localCommand)"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="execLoading"
            :disabled="!localCommand.trim()"
            @click="emit('exec-command', localCommand)"
          >
            {{ $t('admin.providers.execRun') }}
          </el-button>
          <el-button
            v-if="execResult"
            @click="localCommand = ''; emit('clear-exec-result')"
          >
            {{ $t('common.clear') }}
          </el-button>
        </el-form-item>
        <el-form-item
          v-if="execResult !== null"
          :label="$t('admin.providers.execResult')"
        >
          <div class="exec-output">
            <div v-if="execResult.stdout">
              {{ execResult.stdout }}
            </div>
            <div
              v-if="execResult.stderr"
              style="color: #f48771;"
            >
              {{ execResult.stderr }}
            </div>
            <div
              v-if="!execResult.stdout && !execResult.stderr"
              style="color: #888;"
            >
              {{ $t('admin.providers.execNoOutput') }}
            </div>
          </div>
        </el-form-item>
      </el-form>
    </template>
  </div>
</template>

<script setup>
import { Connection, WarningFilled, CircleCheck, InfoFilled } from '@element-plus/icons-vue'
import { useConnectionTab } from './composables/useConnectionTab'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true
  },
  isEditing: {
    type: Boolean,
    default: false
  },
  testingConnection: {
    type: Boolean,
    default: false
  },
  connectionTestResult: {
    type: Object,
    default: null
  },
  generatingSecret: {
    type: Boolean,
    default: false
  },
  agentConnectCmd: {
    type: String,
    default: ''
  },
  agentConnectCmdGithub: {
    type: String,
    default: ''
  },
  execLoading: {
    type: Boolean,
    default: false
  },
  execResult: {
    type: Object,
    default: null
  },
  checkingAgentStatus: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'test-connection',
  'apply-timeout',
  'auth-method-change',
  'generate-agent-secret',
  'check-agent-status',
  'exec-command',
  'clear-exec-result'
])

const {
  t,
  localCommand,
  useCDN,
  useWSS,
  useControllerSource,
  wssUnavailable,
  probingWSS,
  detectingLocal,
  localDetectionResult,
  localCommandChecks,
  isAgentMode,
  isLocalMode,
  hasAgentMappedNetworking,
  showSSHSettings,
  effectiveAgentStatus,
  agentAlertType,
  agentStatusLabel,
  installCmdDisplay,
  extractWsOrigin,
  probeWssAvailability,
  formatOnlineDuration,
  formatDateTime,
  formatLocalDetectStatus,
  handleDetectLocalProvider,
  copyCmd,
} = useConnectionTab(props, emit)
</script>

<style scoped>
.connection-tab {
  max-height: 560px;
  overflow-y: auto;
  padding-right: 8px;
}

.server-form {
  padding-right: 10px;
}

.form-tip {
  margin-top: 5px;
}

.local-mode-content {
  padding: 4px 0;
}

.local-detect-actions {
  margin-top: 12px;
}

.local-detect-result {
  margin-top: 12px;
}

.local-detect-summary {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
  margin-top: 8px;
}

.local-command-list,
.local-warning-list {
  margin-top: 10px;
  font-size: 13px;
}

.local-command-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 6px;
}

.local-command-row code {
  color: var(--el-text-color-secondary);
  word-break: break-all;
}

.local-warning-list ul {
  margin: 6px 0 0;
  padding-left: 18px;
}

.agent-mode-content {
  padding: 4px 0;
}

.install-cmd-box {
  border: 1px solid var(--el-color-warning-light-5);
  border-radius: 6px;
  background: var(--el-color-warning-light-9);
  overflow: hidden;
}

.install-cmd-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 14px;
  background: var(--el-color-warning-light-7);
  font-size: 13px;
  font-weight: 500;
}

.install-cmd-content {
  padding: 12px 14px;
  font-family: monospace;
  font-size: 13px;
  word-break: break-all;
  white-space: pre-wrap;
  background: #1e1e1e;
  color: #d4d4d4;
  max-height: 120px;
  overflow-y: auto;
}

.install-cmd-tip {
  padding: 8px 14px;
  font-size: 12px;
  color: var(--el-color-warning);
  display: flex;
  align-items: center;
  gap: 4px;
}

.exec-output {
  background: #1e1e1e;
  color: #d4d4d4;
  font-family: monospace;
  font-size: 12px;
  padding: 12px;
  border-radius: 4px;
  width: 100%;
  max-height: 300px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-all;
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
