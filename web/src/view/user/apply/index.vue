<template>
  <div class="user-apply">
    <div class="page-header">
      <h1>{{ t('user.apply.title') }}</h1>
      <p>{{ t('user.apply.subtitle') }}</p>
    </div>

    <!-- 用户等级和限制信息 -->
    <el-card class="user-limits-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('user.apply.userQuotaInfo') }}</span>
        </div>
      </template>
      <div class="limits-grid">
        <div class="limit-item">
          <span class="label">{{ t('user.apply.maxInstances') }}</span>
          <span class="value">
            {{ userLimits.usedInstances }} / {{ userLimits.maxInstances }}
            <span
              v-if="userLimits.containerCount !== undefined || userLimits.vmCount !== undefined"
              style="color: #909399; font-size: 12px; margin-left: 8px;"
            >
              ({{ t('user.dashboard.containerCount') }}: {{ userLimits.containerCount || 0 }} / {{ t('user.dashboard.vmCount') }}: {{ userLimits.vmCount || 0 }})
            </span>
          </span>
        </div>
        <div class="limit-item">
          <span class="label">{{ t('user.apply.cpuCoreLimit') }}</span>
          <span class="value">{{ userLimits.usedCpu }} / {{ userLimits.maxCpu }}{{ t('user.apply.cores') }}</span>
        </div>
        <div class="limit-item">
          <span class="label">{{ t('user.apply.memoryLimit') }}</span>
          <span class="value">{{ formatResourceUsage(userLimits.usedMemory, userLimits.maxMemory, 'memory') }}</span>
        </div>
        <div class="limit-item">
          <span class="label">{{ t('user.apply.diskLimit') }}</span>
          <span class="value">{{ formatResourceUsage(userLimits.usedDisk, userLimits.maxDisk, 'disk') }}</span>
        </div>
        <div class="limit-item">
          <span class="label">{{ t('user.apply.trafficLimit') }}</span>
          <span class="value">{{ formatResourceUsage(userLimits.usedTraffic, userLimits.maxTraffic, 'disk') }}</span>
        </div>
      </div>
    </el-card>

    <!-- 常驻兑换码区域：兑换码不再依赖节点选择 -->
    <el-card class="global-redeem-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('user.apply.redeemCodeTitle') }}</span>
        </div>
      </template>
      <el-alert
        v-if="redeemCodeOnlyProviders.length > 0"
        class="redeem-alert"
        type="info"
        :closable="false"
        show-icon
      >
        <template #title>
          {{ t('user.apply.redeemCodeOnlyNotice') }}
        </template>
        <div class="redeem-only-provider-list">
          <el-tag
            v-for="provider in redeemCodeOnlyProviders"
            :key="provider.id"
            size="small"
            type="info"
          >
            {{ provider.name }}
          </el-tag>
        </div>
      </el-alert>
      <el-form
        class="global-redeem-form"
        label-width="120px"
      >
        <el-form-item :label="t('user.apply.redeemCodeTitle')">
          <el-input
            v-model="redeemCodeInput"
            :placeholder="t('user.apply.redeemCodePlaceholder')"
            class="redeem-input"
            clearable
            @keyup.enter="submitRedemption"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="redeemSubmitting"
            size="large"
            @click="submitRedemption"
          >
            {{ t('user.apply.redeemCodeSubmit') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 服务器选择 -->
    <el-card class="providers-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('user.apply.selectProvider') }}</span>
          <el-button
            size="small"
            :loading="refreshing"
            @click="refreshData"
          >
            {{ t('user.apply.refresh') }}
          </el-button>
        </div>
      </template>

      <el-collapse
        v-if="providerGroups.length > 0"
        v-model="expandedProviderGroups"
        class="provider-groups-collapse"
      >
        <el-collapse-item
          v-for="group in providerGroups"
          :key="group.tabName"
          :name="group.tabName"
        >
          <template #title>
            <div class="group-title">
              <span>{{ group.name }}</span>
              <el-tag
                size="small"
                type="info"
              >
                {{ group.providers.length }}
              </el-tag>
            </div>
          </template>
          <div
            v-if="group.description"
            class="group-description"
            v-html="group.description"
          />
          <div class="providers-grid">
            <div
              v-for="provider in group.providers"
              :key="provider.id"
              class="provider-card"
              :class="{
                'selected': selectedProvider?.id === provider.id,
                'active': provider.status === 'active',
                'offline': provider.status === 'offline' || provider.status === 'inactive',
                'partial': provider.status === 'partial'
              }"
              @click="selectProvider(provider)"
            >
              <div class="provider-header">
                <h3>{{ provider.name }}</h3>
                <el-tag
                  :type="getProviderStatusType(provider.status)"
                  size="small"
                >
                  {{ getProviderStatusText(provider.status) }}
                </el-tag>
              </div>
              <div class="provider-info">
                <div class="info-item">
                  <span class="location-info">
                    <span
                      v-if="provider.countryCode"
                      class="flag-icon"
                    >{{ getFlagEmoji(provider.countryCode) }}</span>
                    {{ t('user.apply.location') }}: {{ formatProviderLocation(provider) }}
                  </span>
                </div>
                <div class="info-item provider-meta-line">
                  <el-tag
                    size="small"
                    type="info"
                  >
                    {{ t('user.apply.virtualizationType') }}: {{ formatProviderType(provider.type) }}
                  </el-tag>
                  <el-tag
                    size="small"
                    type="info"
                  >
                    {{ t('user.apply.networkMode') }}: {{ formatNetworkType(provider.networkType) }}
                  </el-tag>
                  <el-tag
                    v-if="provider.redeemCodeOnly"
                    size="small"
                    type="warning"
                  >
                    {{ t('user.apply.redeemCodeOnlyTag') }}
                  </el-tag>
                </div>
                <div class="info-item">
                  <span>CPU: {{ provider.cpu }}{{ t('user.apply.cores') }}</span>
                </div>
                <div class="info-item">
                  <span>{{ t('user.apply.memoryLimit') }}: {{ formatMemorySize(provider.memory || 0) }}</span>
                </div>
                <div class="info-item">
                  <span>{{ t('user.apply.diskLimit') }}: {{ formatDiskSize(provider.disk || 0) }}</span>
                </div>
                <div
                  v-if="provider.containerEnabled && provider.vmEnabled"
                  class="info-item"
                >
                  <span>
                    {{ t('user.apply.availableInstances') }}:
                    {{ t('user.apply.container') }}{{ provider.availableContainerSlots === -1 ? t('user.apply.unlimited') : provider.availableContainerSlots }} /
                    {{ t('user.apply.vm') }}{{ provider.availableVMSlots === -1 ? t('user.apply.unlimited') : provider.availableVMSlots }}
                  </span>
                </div>
                <div
                  v-else-if="provider.containerEnabled"
                  class="info-item"
                >
                  <span>{{ t('user.apply.availableInstances') }}: {{ provider.availableContainerSlots === -1 ? t('user.apply.unlimited') : provider.availableContainerSlots }}</span>
                </div>
                <div
                  v-else-if="provider.vmEnabled"
                  class="info-item"
                >
                  <span>{{ t('user.apply.availableInstances') }}: {{ provider.availableVMSlots === -1 ? t('user.apply.unlimited') : provider.availableVMSlots }}</span>
                </div>
                <div class="info-item">
                  <el-link
                    type="primary"
                    :underline="false"
                    @click.stop="viewHardwareReport(provider.id)"
                  >
                    {{ t('user.apply.viewHardwareReport') }}
                  </el-link>
                </div>
              </div>
            </div>
          </div>
        </el-collapse-item>
      </el-collapse>
    </el-card>

    <!-- 配置表单 -->
    <el-card
      v-if="selectedProvider"
      class="config-card"
    >
      <template #header>
        <div class="card-header">
          <span>{{ t('user.apply.configInstance') }} - {{ selectedProvider.name }}</span>
        </div>
      </template>
      <el-form 
        v-if="!selectedProvider.redeemCodeOnly"
        ref="formRef"
        :model="configForm"
        :rules="configRules"
        label-width="120px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              :label="t('user.apply.instanceType')"
              prop="type"
            >
              <el-select
                v-model="configForm.type"
                :placeholder="t('user.apply.selectInstanceType')"
                @change="onInstanceTypeChange"
              >
                <el-option 
                  :label="t('user.apply.container')" 
                  value="container" 
                  :disabled="!canCreateInstanceType('container')"
                />
                <el-option 
                  :label="t('user.apply.vm')" 
                  value="vm" 
                  :disabled="!canCreateInstanceType('vm')"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="t('user.apply.systemImage')"
              prop="imageId"
            >
              <el-select
                v-model="configForm.imageId"
                :placeholder="t('user.apply.selectSystemImage')"
              >
                <el-option 
                  v-for="image in availableImages" 
                  :key="image.id" 
                  :label="image.name" 
                  :value="image.id"
                >
                  <span style="display: inline-flex; align-items: center; gap: 6px;">
                    <OsIcon
                      :name="image.name"
                      :size="18"
                    />
                    {{ image.name }}
                  </span>
                  <span style="float: right; color: #8492a6; font-size: 12px; margin-left: 10px">
                    {{ formatImageRequirements(image) }}
                  </span>
                </el-option>
              </el-select>
              <div
                v-if="selectedImageInfo"
                class="form-hint"
                style="margin-top: 5px; font-size: 12px; color: #909399;"
              >
                {{ t('user.apply.imageRequirements', { 
                  memory: selectedImageInfo.minMemoryMB, 
                  disk: Math.round(selectedImageInfo.minDiskMB / 1024 * 10) / 10 
                }) }}
              </div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              :label="t('user.apply.cpuSpec')"
              prop="cpuId"
            >
              <el-select
                v-model="configForm.cpuId"
                :placeholder="t('user.apply.selectCpuSpec')"
              >
                <el-option 
                  v-for="cpu in availableCpuSpecs" 
                  :key="cpu.id" 
                  :label="cpu.name" 
                  :value="cpu.id"
                  :disabled="!canSelectSpec('cpu', cpu)"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="t('user.apply.memorySpec')"
              prop="memoryId"
            >
              <el-select
                v-model="configForm.memoryId"
                :placeholder="t('user.apply.selectMemorySpec')"
              >
                <el-option 
                  v-for="memory in availableMemorySpecs" 
                  :key="memory.id" 
                  :label="memory.name" 
                  :value="memory.id"
                  :disabled="!canSelectSpec('memory', memory)"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              :label="t('user.apply.diskSpec')"
              prop="diskId"
            >
              <el-select
                v-model="configForm.diskId"
                :placeholder="t('user.apply.selectDiskSpec')"
              >
                <el-option 
                  v-for="disk in availableDiskSpecs" 
                  :key="disk.id" 
                  :label="disk.name" 
                  :value="disk.id"
                  :disabled="!canSelectSpec('disk', disk)"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="t('user.apply.bandwidthSpec')"
              prop="bandwidthId"
            >
              <el-select
                v-model="configForm.bandwidthId"
                :placeholder="t('user.apply.selectBandwidthSpec')"
              >
                <el-option 
                  v-for="bandwidth in availableBandwidthSpecs" 
                  :key="bandwidth.id" 
                  :label="bandwidth.name" 
                  :value="bandwidth.id"
                  :disabled="!canSelectSpec('bandwidth', bandwidth)"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- GPU passthrough, using backend-cached device data. -->
        <el-form-item
          v-if="canConfigureGpuPassthrough"
          :label="t('user.apply.gpuPassthrough')"
        >
          <div class="gpu-config-wrap">
            <el-checkbox
              v-model="configForm.gpuEnabled"
              @change="onGpuEnabledChange"
            >
              {{ t('user.apply.gpuEnabled') }}
            </el-checkbox>
            <div
              v-if="configForm.gpuEnabled"
              style="margin-top: 8px;"
            >
              <!-- Available GPUs (cached, checkbox selection) -->
              <div
                v-if="detectedGpus.length > 0"
                class="gpu-options-wrap"
              >
                <el-checkbox-group
                  v-model="selectedGpuIndices"
                  class="gpu-options"
                >
                  <el-checkbox
                    v-for="(gpu, idx) in detectedGpus"
                    :key="idx"
                    :value="idx"
                    :label="idx"
                    class="gpu-option-item"
                  >
                    GPU {{ idx }} — {{ gpu.label || gpu.name || '' }}
                  </el-checkbox>
                </el-checkbox-group>
                <div class="gpu-batch-actions">
                  <el-button
                    size="small"
                    text
                    @click="selectAllGpus"
                  >
                    {{ t('user.apply.gpuSelectAll') }}
                  </el-button>
                  <el-button
                    size="small"
                    text
                    @click="deselectAllGpus"
                  >
                    {{ t('user.apply.gpuDeselectAll') }}
                  </el-button>
                </div>
                <div style="font-size: 11px; color: #909399; margin-top: 4px;">
                  {{ t('user.apply.gpuDeviceIdsHint') }}
                </div>
              </div>
              <div
                v-else-if="gpuInfoMsg"
                style="font-size: 12px; color: #909399;"
              >
                {{ gpuInfoMsg }}
              </div>
              <!-- Fallback text input if no cached GPUs -->
              <div v-else>
                <span style="font-size: 12px; color: #909399; margin-right: 8px;">{{ t('user.apply.gpuDeviceIds') }}:</span>
                <el-input
                  v-model="configForm.gpuDeviceIds"
                  :placeholder="t('user.apply.gpuDeviceIdsPlaceholder')"
                  style="max-width: 340px;"
                  size="small"
                />
                <div style="font-size: 11px; color: #c0c4cc; margin-top: 4px;">
                  {{ t('user.apply.gpuDeviceIdsHint') }}
                </div>
              </div>
            </div>
          </div>
        </el-form-item>

        <el-form-item :label="t('user.apply.remarks')">
          <el-input 
            v-model="configForm.description"
            type="textarea"
            :rows="3"
            :placeholder="t('user.apply.remarksPlaceholder')"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>

        <el-form-item>
          <el-button 
            type="primary" 
            :loading="submitting"
            size="large"
            @click="submitApplication"
          >
            {{ t('user.apply.submitApplication') }}
          </el-button>
          <el-button
            size="large"
            @click="resetForm"
          >
            {{ t('user.apply.resetConfig') }}
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 仅兑换码节点不再重复展示兑换入口 -->
      <div
        v-else
        class="redeem-card disabled-redeem-card"
      >
        <el-alert
          type="warning"
          :closable="false"
          show-icon
          :title="t('user.apply.redeemCodeOnlySelectedTitle', { provider: selectedProvider.name })"
          :description="t('user.apply.redeemCodeOnlySelectedDescription')"
        />
        <el-form
          label-width="120px"
          class="disabled-redeem-form"
        >
          <el-form-item :label="t('user.apply.redeemCodeTitle')">
            <el-input
              :model-value="''"
              :placeholder="t('user.apply.useGlobalRedeemArea')"
              class="redeem-input"
              disabled
            />
          </el-form-item>
          <el-form-item>
            <el-button
              size="large"
              disabled
            >
              {{ t('user.apply.redeemCodeSubmit') }}
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>

    <!-- 空状态 -->
    <el-empty 
      v-if="providers.length === 0 && !loading"
      :description="t('user.apply.noProvidersDescription')"
    >
      <template #description>
        <p>{{ t('user.apply.noProvidersMessage') }}</p>
        <p style="font-size: 12px; color: #909399; margin-top: 8px;">
          {{ t('user.apply.noProvidersHint') }}
        </p>
      </template>
      <el-button
        type="primary"
        @click="() => loadProviders(true)"
      >
        {{ t('user.apply.refresh') }}
      </el-button>
    </el-empty>

    <!-- 加载状态 -->
    <div
      v-if="loading"
      class="loading-container"
    >
      <el-skeleton
        :rows="5"
        animated
      />
    </div>
  </div>
</template>

<script setup>
import OsIcon from '@/components/OsIcon.vue'
import useApplyPage from './useApplyPage'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

const {
  loading, refreshing, providers, selectedProvider, providerCapabilities,
  instanceTypePermissions, redeemCodeOnlyProviders,
  getProviderStatusType, getProviderStatusText, formatProviderLocation,
  canCreateInstanceType,
  submitting, redeemCodeInput, redeemSubmitting,
  availableImages, instanceConfig, userLimits, configForm, formRef,
  configRules, availableCpuSpecs, selectedImageInfo,
  availableMemorySpecs, availableDiskSpecs, availableBandwidthSpecs,
  canSelectSpec, formatCpuSpecName, formatImageRequirements,
  onInstanceTypeChange, submitApplication, submitRedemption, resetForm,
  gpuLoading, detectedGpus, selectedGpuIndices, gpuInfoMsg,
  selectAllGpus, deselectAllGpus,
  expandedProviderGroups, canConfigureGpuPassthrough,
  onGpuEnabledChange, providerGroups,
  viewHardwareReport, refreshData, selectProvider,
  formatProviderType, formatNetworkType,
  formatMemorySize, formatDiskSize, formatResourceUsage, getFlagEmoji
} = useApplyPage()
</script>

<style scoped>
.user-apply {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
  color: var(--text-color-primary);
}

.page-header p {
  margin: 0;
  color: var(--text-color-secondary);
}

.user-limits-card,
.global-redeem-card,
.providers-card,
.config-card {
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.redeem-card {
  padding: 8px 0;
}

.global-redeem-form {
  margin-top: 16px;
}

.redeem-input {
  max-width: 420px;
}

.redeem-alert {
  margin-bottom: 8px;
}

.redeem-only-provider-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
}

.disabled-redeem-card {
  padding: 12px 0;
  opacity: 0.72;
}

.disabled-redeem-form {
  margin-top: 16px;
  pointer-events: none;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  color: var(--text-color-primary);
}

.provider-groups-collapse {
  border-top: none;
}

.group-title {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.limits-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.limit-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: var(--neutral-bg);
  border-radius: 8px;
}

.limit-item .label {
  color: var(--text-color-secondary);
  font-weight: 500;
}

.limit-item .value {
  color: var(--text-color-primary);
  font-weight: 600;
}

.providers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.provider-meta-line {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.provider-card {
  border: 2px solid var(--border-color);
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  background-color: var(--card-bg-solid);
}

.provider-card:hover {
  border-color: #3b82f6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.15);
  transform: translateY(-2px);
}

.provider-card.selected {
  border-color: #3b82f6;
  background-color: var(--info-bg);
  box-shadow: 0 4px 16px rgba(59, 130, 246, 0.2);
}

/* Active状态 - 绿色 */
.provider-card.active {
  border-color: #10b981;
  background-color: var(--success-bg);
}

.provider-card.active:hover {
  border-color: #059669;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.2);
}

.provider-card.active.selected {
  border-color: #059669;
  background-color: var(--success-bg);
  box-shadow: 0 4px 16px rgba(16, 185, 129, 0.25);
}

/* Offline状态 - 红色 */
.provider-card.offline {
  border-color: #ef4444;
  background-color: var(--error-bg);
  cursor: not-allowed;
  opacity: 0.7;
  position: relative;
}

.provider-card.offline::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(239, 68, 68, 0.1);
  border-radius: 10px;
  pointer-events: none;
}

.provider-card.offline:hover {
  border-color: #dc2626;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.2);
  transform: none;
}

.provider-card.offline * {
  color: #9ca3af !important;
}

/* Partial状态 - 黄色 */
.provider-card.partial {
  border-color: #f59e0b;
  background-color: var(--warning-bg);
}

.provider-card.partial:hover {
  border-color: #d97706;
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.2);
}

.provider-card.partial.selected {
  border-color: #d97706;
  background-color: var(--warning-bg);
  box-shadow: 0 4px 16px rgba(245, 158, 11, 0.25);
}

.provider-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.provider-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color-primary);
}

.provider-info {
  margin-bottom: 12px;
}

.location-info {
  display: flex;
  align-items: center;
  gap: 6px;
}

.country-flag {
  width: 16px;
  height: 12px;
  border-radius: 2px;
  flex-shrink: 0;
}

.info-item {
  margin-bottom: 4px;
  font-size: 14px;
  color: var(--text-color-secondary);
}

.loading-container {
  padding: 24px;
}

.group-description {
  padding: 12px 16px;
  margin-bottom: 16px;
  background: var(--neutral-bg, #f5f7fa);
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-color-secondary);
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
