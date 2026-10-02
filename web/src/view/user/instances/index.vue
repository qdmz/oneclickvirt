<template>
  <div class="user-instances">
    <!-- 加载状态 -->
    <div
      v-if="loading"
      class="loading-container"
    >
      <el-loading-directive />
      <div class="loading-text">
        {{ t('user.instances.loadingInstances') }}
      </div>
    </div>
    
    <!-- 主要内容 -->
    <div v-else>
      <div class="page-header">
        <h1>{{ t('user.instances.title') }}</h1>
        <p>{{ t('user.instances.subtitle') }}</p>
      </div>

      <!-- 筛选和搜索 -->
      <div class="filter-section">
        <el-form
          :inline="true"
          :model="filterForm"
        >
          <el-form-item>
            <el-input
              v-model="filterForm.name"
              :placeholder="t('user.instances.searchByName')"
              clearable
              style="width: 200px;"
            />
          </el-form-item>
          <el-form-item>
            <el-input
              v-model="filterForm.providerName"
              :placeholder="t('user.instances.searchByProvider')"
              clearable
              style="width: 200px;"
            />
          </el-form-item>
          <el-form-item>
            <el-select
              v-model="filterForm.type"
              :placeholder="t('user.instances.selectType')"
              clearable
              style="width: 150px;"
            >
              <el-option
                :label="t('user.instances.allTypes')"
                value=""
              />
              <el-option
                :label="t('user.instances.vm')"
                value="vm"
              />
              <el-option
                :label="t('user.instances.container')"
                value="container"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-select
              v-model="filterForm.status"
              :placeholder="t('user.instances.selectStatus')"
              clearable
              style="width: 150px;"
            >
              <el-option
                :label="t('user.instances.allStatuses')"
                value=""
              />
              <el-option
                :label="t('user.instances.statusRunning')"
                value="running"
              />
              <el-option
                :label="t('user.instances.statusStopped')"
                value="stopped"
              />
              <el-option
                :label="t('user.instances.statusPaused')"
                value="paused"
              />
              <el-option
                :label="t('user.instances.statusCreating')"
                value="creating"
              />
              <el-option
                :label="t('user.instances.statusStarting')"
                value="starting"
              />
              <el-option
                :label="t('user.instances.statusStopping')"
                value="stopping"
              />
              <el-option
                :label="t('user.instances.statusRestarting')"
                value="restarting"
              />
              <el-option
                :label="t('user.instances.statusRebuilding')"
                value="rebuilding"
              />
              <el-option
                :label="t('user.instances.statusResetting')"
                value="resetting"
              />
              <el-option
                :label="t('user.instances.statusError')"
                value="error"
              />
              <el-option
                :label="t('user.instances.statusFailed')"
                value="failed"
              />
              <el-option
                :label="t('user.instances.statusDeleting')"
                value="deleting"
              />
              <el-option
                :label="t('user.instances.statusDeleted')"
                value="deleted"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="handleSearch"
            >
              {{ t('user.instances.search') }}
            </el-button>
            <el-button @click="resetFilter">
              {{ t('user.instances.reset') }}
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 实例列表 -->
      <div class="instances-grid">
        <div 
          v-for="instance in instances" 
          :key="instance.id" 
          class="instance-card"
          :class="{ 'instance-card-disabled': !canOpenInstanceDetail(instance) }"
          @click="viewInstanceDetail(instance)"
        >
          <div class="instance-header">
            <div class="instance-info">
              <h3>{{ instance.name }}</h3>
              <div class="instance-type">
                <el-tag :type="instance.instance_type === 'vm' ? 'primary' : 'success'">
                  {{ instance.instance_type === 'vm' ? t('user.instances.vm') : t('user.instances.container') }}
                </el-tag>
                <el-tag 
                  v-if="instance.providerType"
                  :type="getProviderTypeColor(instance.providerType)"
                  style="margin-left: 8px;"
                >
                  {{ getProviderTypeName(instance.providerType) }}
                </el-tag>
              </div>
            </div>
            <div class="instance-status">
              <el-tag 
                :type="getStatusType(instance.status)"
                effect="dark"
              >
                {{ getStatusText(instance.status) }}
              </el-tag>
            </div>
          </div>

          <div class="instance-details">
            <div class="detail-item">
              <span class="label">{{ t('user.instances.configuration') }}:</span>
              <span class="value">{{ instance.cpu }}{{ t('user.instances.cores') }} / {{ formatMemorySize(instance.memory) }} / {{ formatDiskSize(instance.disk) }}</span>
            </div>
            <div class="detail-item">
              <span class="label">{{ t('user.instances.bandwidth') }}:</span>
              <span class="value">{{ instance.bandwidth || 100 }}Mbps</span>
            </div>
            <div class="detail-item">
              <span class="label">{{ t('user.instances.system') }}:</span>
              <span class="value"><OsIcon
                :name="instance.osType || instance.image"
                :size="18"
                style="margin-right: 4px;"
              />{{ instance.osType }}</span>
            </div>
            <div class="detail-item">
              <span class="label">{{ t('user.instances.createdAt') }}:</span>
              <span class="value">{{ formatDate(instance.createdAt) }}</span>
            </div>
            <!-- 端口映射信息 -->
            <div
              v-if="instance.portMappings && instance.portMappings.length > 0"
              class="detail-item port-info"
            >
              <span class="label">{{ t('user.instances.portMapping') }}:</span>
              <div class="port-mappings">
                <div class="public-ip">
                  <el-tag
                    type="info"
                    size="small"
                  >
                    {{ t('user.instances.publicIP') }}: {{ instance.publicIP || t('user.instances.unassigned') }}
                  </el-tag>
                </div>
                <!-- 普通用户不显示端口区间 -->
                <div class="port-list">
                  <el-tag 
                    v-for="port in instance.portMappings.slice(0, 3)" 
                    :key="port.id"
                    size="small"
                    effect="plain"
                    class="port-tag"
                    :type="port.isSSH ? 'warning' : 'primary'"
                  >
                    <span v-if="port.isSSH">SSH: {{ port.hostPort }}</span>
                    <span v-else>{{ port.hostPort }}:{{ port.guestPort }}/{{ port.protocol }}</span>
                  </el-tag>
                  <el-tag 
                    v-if="instance.portMappings.length > 3"
                    size="small"
                    type="info"
                    effect="plain"
                  >
                    {{ t('user.instances.morePortsCount', { count: instance.portMappings.length - 3 }) }}
                  </el-tag>
                </div>
              </div>
            </div>
          </div>

          <!-- 实例操作按钮 -->
          <div
            class="instance-actions"
            @click.stop
          >
            <el-button
              v-if="instance.trafficQuotaVisible !== false"
              size="small"
              type="primary"
              :disabled="!canOpenInstanceDetail(instance)"
              @click="showTrafficDetail(instance)"
            >
              <el-icon><TrendCharts /></el-icon>
              {{ t('user.instances.trafficDetail') }}
            </el-button>
            <el-button
              size="small"
              :disabled="!canOpenInstanceDetail(instance)"
              @click="viewInstanceDetail(instance)"
            >
              <el-icon><View /></el-icon>
              {{ t('user.instances.viewDetail') }}
            </el-button>
            <el-button
              size="small"
              type="success"
              :disabled="!canOpenInstanceDetail(instance) || instance.trafficOperationLocked"
              :title="instance.trafficOperationLockMessage || ''"
              @click="createShareLink(instance)"
            >
              <el-icon><Link /></el-icon>
              {{ t('user.instances.share') }}
            </el-button>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <el-empty
        v-if="instances.length === 0 && !loading"
        :description="t('user.instances.noInstances')"
      >
        <el-button
          type="primary"
          @click="$router.push('/user/apply')"
        >
          {{ t('user.instances.applyNow') }}
        </el-button>
      </el-empty>

      <!-- 分页 -->
      <div
        v-if="total > 0"
        class="pagination"
      >
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="() => loadInstances()"
          @current-change="() => loadInstances()"
        />
      </div>

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
    </div> <!-- 结束主要内容区域 -->

    <!-- 流量详情对话框 -->
    <InstanceTrafficDetail
      v-model="showTrafficDialog"
      :instance-id="selectedInstanceForTraffic?.id"
      :instance-name="selectedInstanceForTraffic?.name"
    />
  </div>
</template>

<script setup>
import { TrendCharts, View, Link } from '@element-plus/icons-vue'
import { formatDiskSize, formatMemorySize } from '@/utils/unit-formatter'
import InstanceTrafficDetail from '@/components/InstanceTrafficDetail.vue'
import OsIcon from '@/components/OsIcon.vue'
import { useUserInstances } from './composables/useUserInstances.js'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

const {
  loading,
  instances,
  total,
  showTrafficDialog,
  selectedInstanceForTraffic,
  filterForm,
  pagination,
  handleSearch,
  loadInstances,
  resetFilter,
  getStatusType,
  getStatusText,
  getProviderTypeName,
  getProviderTypeColor,
  formatDate,
  canOpenInstanceDetail,
  viewInstanceDetail,
  showTrafficDetail,
  createShareLink
} = useUserInstances()
</script>

<style scoped>
.user-instances {
  padding: 24px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  color: #666;
}

.loading-text {
  margin-top: 16px;
  font-size: 14px;
}

.page-header {
  margin-bottom: 24px;
  animation: fadeInUp 0.5s cubic-bezier(0.4, 0, 0.2, 1) both;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.page-header h1 {
  margin: 0 0 8px 0;
  font-size: 26px;
  font-weight: 700;
  color: var(--text-color-primary);
  letter-spacing: -0.5px;
}

.page-header p {
  margin: 0;
  color: var(--text-color-secondary);
  font-size: 14px;
}

.filter-section {
  background: var(--card-bg-solid);
  padding: 18px 20px;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  margin-bottom: 24px;
  border: 1px solid var(--border-color);
  animation: fadeInUp 0.55s cubic-bezier(0.4, 0, 0.2, 1) 0.05s both;
}

.filter-section :deep(.el-input__wrapper) {
  border-radius: 12px !important;
}

.instances-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 400px), 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.instance-card {
  background: var(--card-bg-solid);
  border: 1px solid var(--border-color);
  border-radius: 18px;
  padding: 22px;
  min-width: 0;
  max-width: 100%;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.35s cubic-bezier(0.4, 0, 0.2, 1),
              box-shadow 0.35s cubic-bezier(0.4, 0, 0.2, 1),
              border-color 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  position: relative;
}

.instance-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #10b981, #06b6d4, #6366f1);
  border-radius: 18px 18px 0 0;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.instance-card:hover {
  border-color: #10b981;
  box-shadow: 0 12px 32px rgba(16, 185, 129, 0.12), 0 0 0 1px rgba(16, 185, 129, 0.06);
  transform: translateY(-4px);
}

.instance-card:hover::before {
  opacity: 1;
}

.instance-card-disabled {
  cursor: not-allowed;
  opacity: 0.78;
}

.instance-card-disabled:hover {
  border-color: var(--border-color);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  transform: none;
}

.instance-card-disabled:hover::before {
  opacity: 0;
}

.instance-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.instance-info {
  min-width: 0;
}

.instance-info h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  font-weight: 700;
  color: var(--text-color-primary);
  overflow-wrap: anywhere;
}

.instance-type {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.instance-type .el-tag {
  margin-left: 0 !important;
  border-radius: 999px !important;
}

.instance-status {
  flex-shrink: 0;
}

.instance-status :deep(.el-tag) {
  border-radius: 999px !important;
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

.instance-status :deep(.el-tag)::before {
  content: '';
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 6px currentColor;
}

.instance-details {
  margin-bottom: 16px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
  font-size: 14px;
}

.detail-item.port-info {
  flex-direction: column;
  align-items: flex-start;
}

.detail-item .label {
  color: var(--text-color-secondary);
  font-weight: 600;
  min-width: 80px;
  font-size: 13px;
}

.detail-item .value {
  color: var(--text-color-primary);
  text-align: right;
  flex: 1;
  min-width: 0;
  overflow-wrap: anywhere;
  font-variant-numeric: tabular-nums;
}

.port-mappings {
  margin-top: 8px;
  width: 100%;
}

.public-ip, .port-range, .ipv6-info {
  margin-bottom: 8px;
}

.port-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 4px;
}

.port-tag {
  margin: 2px;
  font-size: 12px;
  border-radius: 8px !important;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

.loading-container {
  padding: 24px;
}

.instance-actions {
  border-top: 1px solid var(--el-border-color-lighter);
  padding-top: 14px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: flex-end;
}

.instance-actions .el-button {
  font-size: 12px;
  margin-left: 0;
  border-radius: 10px !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.instance-actions .el-button:hover {
  transform: translateY(-1px);
}

.instance-card {
  animation: fadeInUp 0.5s cubic-bezier(0.4, 0, 0.2, 1) both;
}

.instance-card:nth-child(1) { animation-delay: 0.06s; }
.instance-card:nth-child(2) { animation-delay: 0.12s; }
.instance-card:nth-child(3) { animation-delay: 0.18s; }
.instance-card:nth-child(4) { animation-delay: 0.24s; }
.instance-card:nth-child(5) { animation-delay: 0.30s; }
.instance-card:nth-child(6) { animation-delay: 0.36s; }
.instance-card:nth-child(7) { animation-delay: 0.42s; }
.instance-card:nth-child(8) { animation-delay: 0.48s; }

@media (max-width: 1024px) {
  .instances-grid {
    grid-template-columns: repeat(auto-fill, minmax(min(100%, 320px), 1fr));
    gap: 16px;
  }
}

@media (max-width: 768px) {
  .user-instances {
    padding: 0;
  }

  .page-header h1 {
    font-size: 22px;
  }

  .filter-section {
    padding: 14px;
    border-radius: 14px;
    margin-bottom: 16px;
  }

  .filter-section :deep(.el-form--inline) {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }

  .filter-section :deep(.el-form--inline .el-form-item) {
    margin-right: 0;
    margin-bottom: 0;
  }

  .filter-section :deep(.el-input),
  .filter-section :deep(.el-select),
  .filter-section :deep(.el-button) {
    width: 100% !important;
  }

  .instances-grid {
    gap: 12px;
  }

  .instance-card {
    padding: 16px;
    border-radius: 14px;
  }

  .instance-header {
    gap: 12px;
  }

  .instance-info h3 {
    font-size: 16px;
  }

  .detail-item {
    gap: 12px;
    font-size: 13px;
  }

  .detail-item .label {
    min-width: 72px;
    font-size: 12px;
  }

  .instance-actions {
    justify-content: flex-start;
  }

  .instance-actions .el-button {
    font-size: 11px;
    padding: 6px 10px;
  }
}

@media (max-width: 480px) {
  .page-header h1 {
    font-size: 20px;
  }

  .instance-card {
    padding: 14px;
  }

  .instance-info h3 {
    font-size: 15px;
  }

  .detail-item {
    font-size: 12px;
    margin-bottom: 6px;
  }

  .detail-item .label {
    min-width: 64px;
  }

  .port-tag {
    font-size: 10px;
  }

  .instance-actions .el-button {
    font-size: 10px;
    padding: 5px 8px;
  }
}
</style>
