<template>
  <div class="admin-group-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <div>
            <div class="title">
              {{ t('admin.group.title') }}
            </div>
            <div class="subtitle">
              使用分组控制普通用户申请领取页的节点展示；描述支持 Markdown，并会安全渲染部分 HTML。
            </div>
          </div>
          <el-button
            type="primary"
            @click="openCreateDialog"
          >
            新增分组
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="groups"
        stripe
      >
        <el-table-column
          prop="groupName"
          label="分组名称"
          min-width="150"
        />
        <el-table-column
          label="分组描述"
          min-width="260"
        >
          <template #default="{ row }">
            <div
              v-if="row.groupDescriptionHtml"
              class="description-preview"
              v-html="row.groupDescriptionHtml"
            />
            <el-text
              v-else
              type="info"
            >
              未填写
            </el-text>
          </template>
        </el-table-column>
        <el-table-column
          label="节点"
          min-width="260"
        >
          <template #default="{ row }">
            <el-space wrap>
              <el-tag
                v-for="provider in row.providers"
                :key="provider.id"
                size="small"
              >
                {{ provider.name }}
              </el-tag>
              <el-text
                v-if="!row.providers?.length"
                type="info"
              >
                暂无节点
              </el-text>
            </el-space>
          </template>
        </el-table-column>
        <el-table-column
          prop="providerCount"
          label="节点数量"
          width="100"
        />
        <el-table-column
          label="操作"
          width="180"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              text
              type="primary"
              @click="openEditDialog(row)"
            >
              编辑
            </el-button>
            <el-button
              text
              type="danger"
              @click="deleteGroup(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty
        v-if="!loading && groups.length === 0"
        description="暂无分组，点击右上角新增分组后再勾选节点"
      />
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="editingGroup?.id ? '编辑分组' : '新增分组'"
      width="860px"
      destroy-on-close
      @closed="resetDialog"
    >
      <el-form
        :model="form"
        label-width="110px"
      >
        <el-form-item
          label="分组名称"
          required
        >
          <el-input
            v-model="form.groupName"
            maxlength="64"
            show-word-limit
            placeholder="例如：香港高性能 / 免费体验 / 美国节点"
          />
        </el-form-item>
        <el-form-item label="分组描述">
          <el-input
            v-model="form.groupDescription"
            type="textarea"
            :rows="6"
            maxlength="20000"
            show-word-limit
            placeholder="支持 Markdown，例如：**注意事项**、列表、链接；也可使用部分安全 HTML 标签"
          />
          <div class="form-item-hint">
            建议用 Markdown 编写说明，会像 GitHub 一样渲染常用标题、列表、粗体、链接与部分安全 HTML。
          </div>
        </el-form-item>
        <el-form-item
          v-if="descriptionPreview"
          label="预览"
        >
          <div
            class="description-preview full"
            v-html="descriptionPreview"
          />
        </el-form-item>
        <el-form-item label="包含节点">
          <el-table
            ref="providerTableRef"
            :data="providers"
            height="320"
            row-key="id"
            @selection-change="onProviderSelectionChange"
          >
            <el-table-column
              type="selection"
              width="48"
              :selectable="canSelectProvider"
            />
            <el-table-column
              prop="name"
              label="节点名称"
              min-width="180"
            />
            <el-table-column
              prop="type"
              label="类型"
              width="100"
            />
            <el-table-column
              prop="status"
              label="状态"
              width="100"
            />
            <el-table-column
              label="当前分组"
              min-width="140"
            >
              <template #default="{ row }">
                <el-tag
                  v-if="row.groupId && row.groupId !== editingGroup?.id"
                  type="warning"
                  size="small"
                >
                  {{ row.groupName || '其他分组' }}
                </el-tag>
                <el-tag
                  v-else-if="row.groupId === editingGroup?.id"
                  type="success"
                  size="small"
                >
                  本分组
                </el-tag>
                <el-text
                  v-else
                  type="info"
                >
                  未分组
                </el-text>
              </template>
            </el-table-column>
          </el-table>
          <div class="form-item-hint">
            一个节点同一时间只属于一个分组；保存后会自动从原分组移出。
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          :loading="saving"
          @click="saveGroup"
        >
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import service from '@/utils/request'
import { renderMarkdown } from '@/utils/markdown'
import { useUserStore } from '@/pinia/modules/user'

const { t } = useI18n()
const userStore = useUserStore()
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const editingGroup = ref(null)
const providerTableRef = ref(null)
const groups = ref([])
const providers = ref([])
const selectedProviderIds = ref([])

const form = reactive({
  groupName: '',
  groupDescription: '',
  sortOrder: 0
})

const descriptionPreview = computed(() => renderMarkdown(form.groupDescription))

const fetchGroups = async () => {
  loading.value = true
  try {
    const res = await service({ url: '/v1/admin/groups', method: 'get' })
    if (res.code === 200 && res.data) {
      groups.value = res.data.groups || []
      providers.value = res.data.providers || []
    }
  } catch (e) {
    console.error('Failed to fetch groups:', e)
    ElMessage.error('获取分组失败')
  } finally {
    loading.value = false
  }
}

const resetDialog = () => {
  editingGroup.value = null
  selectedProviderIds.value = []
  form.groupName = ''
  form.groupDescription = ''
  form.sortOrder = 0
}

const openCreateDialog = async () => {
  resetDialog()
  dialogVisible.value = true
  await nextTick()
  providerTableRef.value?.clearSelection()
}

const openEditDialog = async (row) => {
  editingGroup.value = row
  form.groupName = row.groupName || ''
  form.groupDescription = row.groupDescription || ''
  form.sortOrder = row.sortOrder || 0
  selectedProviderIds.value = [...(row.providerIds || [])]
  dialogVisible.value = true
  await nextTick()
  providerTableRef.value?.clearSelection()
  providers.value.forEach(provider => {
    if (selectedProviderIds.value.includes(provider.id)) providerTableRef.value?.toggleRowSelection(provider, true)
  })
}

const canSelectProvider = (provider) => {
  if (editingGroup.value?.id) return provider.ownerAdminId === editingGroup.value.ownerAdminId
  if (userStore.userType === 'admin') return !provider.ownerAdminId
  return true
}

const onProviderSelectionChange = (selection) => {
  selectedProviderIds.value = selection.map(item => item.id)
}

const saveGroup = async () => {
  if (!form.groupName.trim()) {
    ElMessage.warning('请填写分组名称')
    return
  }
  saving.value = true
  try {
    const data = {
      groupName: form.groupName.trim(),
      groupDescription: form.groupDescription,
      sortOrder: form.sortOrder,
      providerIds: selectedProviderIds.value
    }
    const id = editingGroup.value?.id
    const res = await service({
      url: id ? `/v1/admin/groups/${id}` : '/v1/admin/groups',
      method: id ? 'put' : 'post',
      data
    })
    if (res.code === 200) {
      ElMessage.success('保存成功')
      dialogVisible.value = false
      await fetchGroups()
    } else {
      ElMessage.error(res.msg || res.message || '保存失败')
    }
  } catch (e) {
    ElMessage.error(e?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const deleteGroup = async (row) => {
  await ElMessageBox.confirm(`确认删除分组「${row.groupName}」？节点会回到未分组状态。`, '删除分组', { type: 'warning' })
  try {
    const res = await service({ url: `/v1/admin/groups/${row.id}`, method: 'delete' })
    if (res.code === 200) {
      ElMessage.success('删除成功')
      await fetchGroups()
    } else {
      ElMessage.error(res.msg || res.message || '删除失败')
    }
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e?.message || '删除失败')
  }
}

onMounted(fetchGroups)
</script>

<style scoped>
.admin-group-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}
.title {
  font-weight: 600;
  font-size: 16px;
}
.subtitle,
.form-item-hint {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  margin-top: 4px;
}
.description-preview {
  max-height: 150px;
  overflow: auto;
  padding: 8px 10px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 6px;
  background: var(--el-fill-color-lighter);
  line-height: 1.6;
}
.description-preview.full {
  max-height: 260px;
  width: 100%;
}
.description-preview :deep(p) {
  margin: 0 0 6px;
}
.description-preview :deep(ul),
.description-preview :deep(ol) {
  padding-left: 20px;
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
