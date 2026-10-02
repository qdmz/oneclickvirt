<template>
  <div class="admin-dashboard">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.dashboard.title') }}</span>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row
        :gutter="20"
        class="stats-row"
      >
        <el-col
          :xs="24"
          :sm="12"
          :md="12"
          :lg="6"
          :xl="6"
        >
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon user-icon">
                <i class="fas fa-users" />
              </div>
              <div class="stat-info">
                <div class="stat-number">
                  {{ dashboardData.totalUsers }}
                </div>
                <div class="stat-label">
                  {{ $t('admin.dashboard.totalUsers') }}
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      
        <el-col
          :xs="24"
          :sm="12"
          :md="12"
          :lg="6"
          :xl="6"
        >
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon server-icon">
                <i class="fas fa-server" />
              </div>
              <div class="stat-info">
                <div class="stat-number">
                  {{ dashboardData.totalProviders }}
                </div>
                <div class="stat-label">
                  {{ $t('admin.dashboard.totalProviders') }}
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      
        <el-col
          :xs="24"
          :sm="12"
          :md="12"
          :lg="6"
          :xl="6"
        >
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon vm-icon">
                <i class="fas fa-desktop" />
              </div>
              <div class="stat-info">
                <div class="stat-number">
                  {{ dashboardData.totalVMs }}
                </div>
                <div class="stat-label">
                  {{ $t('admin.dashboard.totalVMs') }}
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      
        <el-col
          :xs="24"
          :sm="12"
          :md="12"
          :lg="6"
          :xl="6"
        >
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon container-icon">
                <i class="fas fa-box" />
              </div>
              <div class="stat-info">
                <div class="stat-number">
                  {{ dashboardData.totalContainers }}
                </div>
                <div class="stat-label">
                  {{ $t('admin.dashboard.totalContainers') }}
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 资源概览 -->
      <el-divider content-position="left">
        <span style="font-size: 16px; font-weight: 600; color: var(--text-color-primary);">{{ $t('admin.dashboard.resourceOverview') }}</span>
      </el-divider>
      <el-row
        :gutter="20"
        class="stats-row"
      >
        <el-col
          :xs="24"
          :sm="12"
          :md="8"
          :lg="8"
          :xl="8"
        >
          <el-card class="resource-card">
            <div class="resource-header">
              <div class="resource-icon cpu-icon">
                <i class="fas fa-microchip" />
              </div>
              <span class="resource-title">{{ $t('admin.dashboard.cpuCores') }}</span>
            </div>
            <div class="resource-numbers">
              <div class="resource-number-item">
                <span class="resource-number-value">{{ resourceUsage.usedCpuCores }}</span>
                <span class="resource-number-label">{{ $t('admin.dashboard.used') }}</span>
              </div>
              <span class="resource-number-separator">/</span>
              <div class="resource-number-item">
                <span class="resource-number-value">{{ resourceUsage.totalCpuCores }}</span>
                <span class="resource-number-label">{{ $t('admin.dashboard.total') }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col
          :xs="24"
          :sm="12"
          :md="8"
          :lg="8"
          :xl="8"
        >
          <el-card class="resource-card">
            <div class="resource-header">
              <div class="resource-icon memory-icon">
                <i class="fas fa-memory" />
              </div>
              <span class="resource-title">{{ $t('admin.dashboard.memory') }}</span>
            </div>
            <div class="resource-numbers">
              <div class="resource-number-item">
                <span class="resource-number-value">{{ formatMB(resourceUsage.usedMemoryMB) }}</span>
                <span class="resource-number-label">{{ $t('admin.dashboard.used') }}</span>
              </div>
              <span class="resource-number-separator">/</span>
              <div class="resource-number-item">
                <span class="resource-number-value">{{ formatMB(resourceUsage.totalMemoryMB) }}</span>
                <span class="resource-number-label">{{ $t('admin.dashboard.total') }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col
          :xs="24"
          :sm="12"
          :md="8"
          :lg="8"
          :xl="8"
        >
          <el-card class="resource-card">
            <div class="resource-header">
              <div class="resource-icon disk-icon">
                <i class="fas fa-hdd" />
              </div>
              <span class="resource-title">{{ $t('admin.dashboard.disk') }}</span>
            </div>
            <div class="resource-numbers">
              <div class="resource-number-item">
                <span class="resource-number-value">{{ formatMB(resourceUsage.usedDiskMB) }}</span>
                <span class="resource-number-label">{{ $t('admin.dashboard.used') }}</span>
              </div>
              <span class="resource-number-separator">/</span>
              <div class="resource-number-item">
                <span class="resource-number-value">{{ formatMB(resourceUsage.totalDiskMB) }}</span>
                <span class="resource-number-label">{{ $t('admin.dashboard.total') }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getAdminDashboard } from '@/api/admin'

const { t } = useI18n()

const dashboardData = reactive({
  totalUsers: 0,
  totalProviders: 0,
  totalVMs: 0,
  totalContainers: 0
})

const resourceUsage = reactive({
  totalCpuCores: 0,
  usedCpuCores: 0,
  totalMemoryMB: 0,
  usedMemoryMB: 0,
  totalDiskMB: 0,
  usedDiskMB: 0
})

const formatMB = (mb) => {
  if (!mb || mb <= 0) return '0 MB'
  const GB = 1024
  const TB = 1024 * 1024
  if (mb >= TB) return (mb / TB).toFixed(2) + ' TB'
  if (mb >= GB) return (mb / GB).toFixed(2) + ' GB'
  return mb + ' MB'
}

const fetchDashboardData = async () => {
  try {
    const response = await getAdminDashboard()
    if (response.code === 200) {
      if (response.data && response.data.statistics) {
        Object.assign(dashboardData, response.data.statistics)
      } else {
        Object.assign(dashboardData, response.data)
      }
      if (response.data && response.data.resourceUsage) {
        Object.assign(resourceUsage, response.data.resourceUsage)
      }
    }
  } catch (error) {
    ElMessage.error(t('admin.dashboard.loadDataFailed'))
    console.error('Dashboard data fetch error:', error)
  }
}

onMounted(async () => {
  await fetchDashboardData()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  > span {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-color-primary);
  }
}

.stats-row {
  margin-bottom: 30px;
}

.stat-card {
  height: 140px;
  border-radius: 12px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stat-content {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 10px;
}

.stat-icon {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  color: white;
  flex-shrink: 0;
  font-size: 32px;
}

.user-icon {
  background: linear-gradient(135deg, #16a34a 0%, #22c55e 100%);
  box-shadow: 0 4px 12px rgba(22, 163, 74, 0.4);
}

.server-icon {
  background: linear-gradient(135deg, #059669 0%, #10b981 100%);
  box-shadow: 0 4px 12px rgba(5, 150, 105, 0.4);
}

.vm-icon {
  background: linear-gradient(135deg, #0f766e 0%, #14b8a6 100%);
  box-shadow: 0 4px 12px rgba(15, 118, 110, 0.4);
}

.container-icon {
  background: linear-gradient(135deg, #65a30d 0%, #84cc16 100%);
  box-shadow: 0 4px 12px rgba(101, 163, 13, 0.4);
}

.resource-card {
  border-radius: 12px;
  transition: all 0.3s ease;
  margin-bottom: 16px;
}

.resource-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 18px rgba(0, 0, 0, 0.1);
}

.resource-header {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.resource-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
  margin-right: 10px;
}

.cpu-icon {
  background: linear-gradient(135deg, #3b82f6 0%, #60a5fa 100%);
}

.memory-icon {
  background: linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%);
}

.disk-icon {
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
}

.resource-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-color-primary);
}

.resource-numbers {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 12px 0;
}

.resource-number-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.resource-number-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-color-primary);
}

.resource-number-label {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.resource-number-separator {
  font-size: 20px;
  color: #c0c4cc;
  font-weight: 300;
}

.stat-info {
  flex: 1;
  min-width: 0;
}

.stat-number {
  font-size: 36px;
  font-weight: 700;
  color: var(--text-color-primary);
  line-height: 1.2;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  font-weight: 500;
}

/* 响应式适配 */
/* 平板端适配 */
@media (max-width: 1024px) {
  .stat-card {
    height: 120px;
    margin-bottom: 16px;
  }
  
  .stat-icon {
    width: 60px;
    height: 60px;
    font-size: 28px;
    margin-right: 16px;
  }
  
  .stat-number {
    font-size: 28px;
  }
  
  .stat-label {
    font-size: 13px;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .stats-row {
    margin-bottom: 20px;
  }
  
  .stat-card {
    height: auto;
    min-height: 100px;
    margin-bottom: 12px;
  }
  
  .stat-card:hover {
    transform: none;
  }
  
  .stat-content {
    padding: 16px;
  }
  
  .stat-icon {
    width: 50px;
    height: 50px;
    font-size: 24px;
    margin-right: 12px;
  }
  
  .stat-number {
    font-size: 24px;
  }
  
  .stat-label {
    font-size: 12px;
  }
}

/* 小屏移动端适配 */
@media (max-width: 480px) {
  .stat-content {
    padding: 12px;
  }
  
  .stat-icon {
    width: 45px;
    height: 45px;
    font-size: 20px;
    margin-right: 10px;
  }
  
  .stat-number {
    font-size: 20px;
  }
  
  .stat-label {
    font-size: 11px;
  }
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