<template>
  <div class="announcements-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.announcements.title') }}</span>
          <div class="header-actions">
            <el-button 
              v-if="selectedRows.length > 0" 
              type="danger" 
              :disabled="selectedRows.length === 0"
              @click="handleBatchDelete"
            >
              {{ $t('admin.announcements.batchDelete') }} ({{ selectedRows.length }})
            </el-button>
            <el-button 
              v-if="selectedRows.length > 0" 
              type="warning" 
              :disabled="selectedRows.length === 0"
              :loading="batchUpdating"
              @click="handleBatchToggleStatus"
            >
              {{ $t('admin.announcements.batchToggleStatus') }} ({{ selectedRows.length }})
            </el-button>
            <el-button
              type="primary"
              @click="addAnnouncement"
            >
              {{ $t('admin.announcements.addAnnouncement') }}
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 筛选条件 -->
      <div class="filter-container">
        <el-row
          :gutter="20"
          style="margin-bottom: 20px;"
        >
          <el-col :span="6">
            <el-select
              v-model="filterType"
              :placeholder="$t('admin.announcements.selectType')"
              clearable
              @change="loadAnnouncements"
            >
              <el-option
                :label="$t('admin.announcements.all')"
                value=""
              />
              <el-option
                :label="$t('admin.announcements.homepageAnnouncement')"
                value="homepage"
              />
              <el-option
                :label="$t('admin.announcements.topbarAnnouncement')"
                value="topbar"
              />
            </el-select>
          </el-col>
          <el-col :span="6">
            <el-select
              v-model="filterStatus"
              :placeholder="$t('admin.announcements.selectStatus')"
              clearable
              @change="loadAnnouncements"
            >
              <el-option
                :label="$t('admin.announcements.all')"
                :value="null"
              />
              <el-option
                :label="$t('common.enabled')"
                :value="1"
              />
              <el-option
                :label="$t('common.disabled')"
                :value="0"
              />
            </el-select>
          </el-col>
          <el-col :span="6">
            <el-input 
              v-model="filterTitle" 
              :placeholder="$t('admin.announcements.searchTitle')" 
              clearable 
              @clear="loadAnnouncements"
              @keyup.enter="loadAnnouncements"
            >
              <template #append>
                <el-button
                  icon="Search"
                  @click="loadAnnouncements"
                />
              </template>
            </el-input>
          </el-col>
          <el-col :span="6">
            <el-button @click="resetFilters">
              {{ $t('admin.announcements.resetFilters') }}
            </el-button>
          </el-col>
        </el-row>
      </div>
      
      <el-table 
        v-loading="loading" 
        :data="announcements" 
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column
          prop="title"
          :label="$t('admin.announcements.announcementTitle')"
          width="200"
        />
        <el-table-column
          prop="type"
          :label="$t('common.name')"
          width="120"
        >
          <template #default="scope">
            <el-tag :type="scope.row.type === 'homepage' ? 'success' : 'warning'">
              {{ scope.row.type === 'homepage' ? $t('admin.announcements.homepage') : $t('admin.announcements.topbar') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="priority"
          :label="$t('admin.announcements.priority')"
          min-width="110"
        />
        <el-table-column
          prop="isSticky"
          :label="$t('admin.announcements.isSticky')"
          min-width="90"
        >
          <template #default="scope">
            <el-tag
              :type="scope.row.isSticky ? 'danger' : 'info'"
              size="small"
            >
              {{ scope.row.isSticky ? $t('common.yes') : $t('common.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="status"
          :label="$t('common.status')"
          min-width="90"
        >
          <template #default="scope">
            <el-tag
              :type="scope.row.status === 1 ? 'success' : 'danger'"
              size="small"
            >
              {{ scope.row.status === 1 ? $t('common.enabled') : $t('common.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="content"
          :label="$t('admin.announcements.content')"
          show-overflow-tooltip
        />
        <el-table-column
          prop="createdAt"
          :label="$t('common.createTime')"
          width="160"
        >
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.actions')"
          width="260"
          fixed="right"
        >
          <template #default="scope">
            <div class="row-actions">
              <el-button
                size="small"
                @click="editAnnouncement(scope.row)"
              >
                {{ $t('common.edit') }}
              </el-button>
              <el-button
                size="small"
                :type="scope.row.status === 1 ? 'warning' : 'success'"
                @click="toggleAnnouncementStatus(scope.row)"
              >
                {{ scope.row.status === 1 ? $t('common.disable') : $t('common.enable') }}
              </el-button>
              <el-button
                size="small"
                type="danger"
                @click="deleteAnnouncementHandler(scope.row.id)"
              >
                {{ $t('common.delete') }}
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑公告对话框 -->
    <el-dialog 
      v-model="showAddDialog" 
      :title="isEditing ? $t('admin.announcements.editAnnouncement') : $t('admin.announcements.addAnnouncement')" 
      width="1200px"
      :append-to-body="true"
      class="announcement-dialog"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        label-width="100px"
        :rules="rules"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              :label="$t('admin.announcements.title')"
              prop="title"
            >
              <el-input
                v-model="form.title"
                :placeholder="$t('admin.announcements.titlePlaceholder')"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              :label="$t('admin.announcements.type')"
              prop="type"
            >
              <el-select
                v-model="form.type"
                :placeholder="$t('admin.announcements.typePlaceholder')"
                style="width: 100%"
              >
                <el-option
                  :label="$t('admin.announcements.typeHomepage')"
                  value="homepage"
                />
                <el-option
                  :label="$t('admin.announcements.typeTopbar')"
                  value="topbar"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="$t('admin.announcements.priority')">
              <el-input-number
                v-model="form.priority"
                :min="0"
                :max="100"
                :controls="false"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('admin.announcements.isSticky')">
              <el-switch v-model="form.isSticky" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item
              v-if="isEditing"
              :label="$t('common.status')"
            >
              <el-select
                v-model="form.status"
                style="width: 100%"
              >
                <el-option
                  :label="$t('common.enabled')"
                  :value="1"
                />
                <el-option
                  :label="$t('common.disabled')"
                  :value="0"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item
          :label="$t('admin.announcements.content')"
          prop="content"
        >
          <div class="quill-editor-wrapper">
            <QuillEditor
              v-model:content="form.content"
              content-type="html"
              theme="snow"
              :options="editorOptions"
              @update:content="handleContentChange"
            />
          </div>
        </el-form-item>
        
        <el-row
          v-if="isEditing"
          :gutter="20"
        >
          <el-col :span="12">
            <el-form-item :label="$t('admin.announcements.startTime')">
              <el-date-picker
                v-model="startTime"
                type="datetime"
                :placeholder="$t('admin.announcements.selectStartTime')"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('admin.announcements.endTime')">
              <el-date-picker
                v-model="endTime"
                type="datetime"
                :placeholder="$t('admin.announcements.selectEndTime')"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      
      <template #footer>
        <el-button @click="handleDialogClose">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="saveAnnouncement"
        >
          {{ isEditing ? $t('common.update') : $t('common.save') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { useAnnouncementManagement } from './composables/useAnnouncementManagement'
import { QuillEditor } from '@vueup/vue-quill'
import '@vueup/vue-quill/dist/vue-quill.snow.css'

const {
  announcements,
  showAddDialog,
  loading,
  submitting,
  isEditing,
  formRef,
  selectedRows,
  batchUpdating,
  filterType,
  filterStatus,
  filterTitle,
  startTime,
  endTime,
  form,
  rules,
  editorOptions,
  formatDate,
  handleContentChange,
  loadAnnouncements,
  addAnnouncement,
  editAnnouncement,
  deleteAnnouncementHandler,
  saveAnnouncement,
  handleDialogClose,
  handleSelectionChange,
  handleBatchDelete,
  handleBatchToggleStatus,
  toggleAnnouncementStatus,
  resetFilters
} = useAnnouncementManagement()
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

.filter-container {
  margin-bottom: 20px;
}

.row-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: nowrap;
  white-space: nowrap;
}

.row-actions :deep(.el-button) {
  margin-left: 0;
}

/* 富文本编辑器容器固定大小 */
.quill-editor-wrapper {
  width: 100%;
  height: 400px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

:deep(.quill-editor-wrapper .ql-container) {
  height: calc(100% - 42px);
  font-size: 14px;
}

:deep(.quill-editor-wrapper .ql-editor) {
  min-height: 100%;
  max-height: 100%;
  overflow-y: auto;
}

:deep(.quill-editor-wrapper .ql-toolbar) {
  border-top-left-radius: 4px;
  border-top-right-radius: 4px;
}

:deep(.quill-editor-wrapper .ql-container) {
  border-bottom-left-radius: 4px;
  border-bottom-right-radius: 4px;
}

/* 确保对话框宽度固定 */
:deep(.announcement-dialog) {
  width: 1200px !important;
  max-width: 90vw;
}

:deep(.announcement-dialog .el-dialog) {
  width: 1200px !important;
  max-width: 90vw;
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
