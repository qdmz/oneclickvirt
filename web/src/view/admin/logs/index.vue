<template>
  <div class="log-viewer">
    <!-- 页头 -->
    <el-card
      class="header-card"
      shadow="never"
    >
      <div class="header-content">
        <div class="title-section">
          <h2>
            <el-icon><Document /></el-icon>
            {{ $t('admin.logs.title') }}
          </h2>
          <p class="subtitle">
            {{ $t('admin.logs.subtitle') }}
          </p>
        </div>
      </div>
    </el-card>

    <!-- 工具栏 -->
    <el-card
      class="toolbar-card"
      shadow="never"
    >
      <el-row
        :gutter="16"
        align="middle"
      >
        <!-- 日期选择 -->
        <el-col
          :xs="24"
          :sm="8"
          :md="6"
        >
          <el-select
            v-model="selectedDate"
            :placeholder="$t('admin.logs.selectDate')"
            style="width: 100%"
            @change="onDateChange"
          >
            <el-option-group
              v-if="rootFiles.length"
              :label="$t('admin.logs.rootFiles')"
            >
              <el-option
                v-for="f in rootFiles"
                :key="'root::' + f"
                :label="f"
                :value="'root::' + f"
              />
            </el-option-group>
            <el-option-group
              v-if="dates.length"
              :label="$t('admin.logs.dateFolders')"
            >
              <el-option
                v-for="d in dates"
                :key="d.date"
                :label="d.date"
                :value="d.date"
              />
            </el-option-group>
          </el-select>
        </el-col>

        <!-- 日志级别选择（仅当选择了日期目录时显示） -->
        <el-col
          v-if="isDateSelected"
          :xs="24"
          :sm="6"
          :md="4"
        >
          <el-select
            v-model="selectedType"
            :placeholder="$t('admin.logs.selectLevel')"
            style="width: 100%"
            @change="loadLog"
          >
            <el-option
              v-for="t in currentTypes"
              :key="t"
              :label="t + '.log'"
              :value="t"
            />
          </el-select>
        </el-col>

        <!-- 行数 -->
        <el-col
          :xs="12"
          :sm="5"
          :md="4"
        >
          <el-input-number
            v-model="tailLines"
            :min="50"
            :max="5000"
            :step="100"
            controls-position="right"
            style="width: 100%"
          />
        </el-col>

        <!-- 搜索 -->
        <el-col
          :xs="24"
          :sm="6"
          :md="4"
        >
          <el-input
            v-model="searchKeyword"
            :placeholder="$t('admin.logs.search')"
            clearable
            :prefix-icon="Search"
          />
        </el-col>

        <!-- 操作按钮 -->
        <el-col
          :xs="12"
          :sm="5"
          :md="6"
        >
          <div class="btn-group">
            <el-button
              type="primary"
              :icon="Refresh"
              :loading="loading"
              @click="refreshAll"
            >
              {{ $t('admin.logs.refresh') }}
            </el-button>
            <el-button
              :icon="CopyDocument"
              @click="copyContent"
            >
              {{ $t('admin.logs.copy') }}
            </el-button>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 日志内容 -->
    <el-card
      class="content-card"
      shadow="never"
    >
      <template #header>
        <div class="content-header">
          <span class="file-label">
            <el-icon><Folder /></el-icon>
            {{ currentFileLabel }}
          </span>
          <span
            v-if="lineCount > 0"
            class="line-count"
          >
            <template v-if="searchKeyword && filteredLineCount !== lineCount">
              {{ filteredLineCount }} / {{ lineCount }} {{ $t('admin.logs.lines') }}
            </template>
            <template v-else>
              {{ lineCount }} {{ $t('admin.logs.lines') }}
            </template>
          </span>
        </div>
      </template>

      <div
        ref="logContainerRef"
        class="log-container"
      >
        <div
          v-if="loading"
          class="log-placeholder"
        >
          <el-icon class="is-loading">
            <Loading />
          </el-icon>
          <span>{{ $t('admin.logs.loading') }}</span>
        </div>
        <div
          v-else-if="!hasSelection"
          class="log-placeholder"
        >
          <el-icon><InfoFilled /></el-icon>
          <span>{{ $t('admin.logs.pleaseSelect') }}</span>
        </div>
        <div
          v-else-if="!logContent"
          class="log-placeholder"
        >
          <el-icon><DocumentDelete /></el-icon>
          <span>{{ $t('admin.logs.noContent') }}</span>
        </div>
        <pre
          v-else
          ref="logPreRef"
          class="log-pre"
          v-html="displayContent"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { copyToClipboard } from '@/utils/clipboard'
import {
  Document, Refresh, CopyDocument, Folder, Loading, InfoFilled,
  DocumentDelete, Search
} from '@element-plus/icons-vue'
import { getLogDates, getLogContent } from '@/api/admin'

const { t } = useI18n()

// 状态
const dates = ref([])
const rootFiles = ref([])
const selectedDate = ref('')
const selectedType = ref('')
const tailLines = ref(200)
const logContent = ref('')
const lineCount = ref(0)
const loading = ref(false)
const logContainerRef = ref(null)
const logPreRef = ref(null)
const searchKeyword = ref('')

// 搜索过滤后的显示内容
const filteredLines = computed(() => {
  if (!logContent.value) return []
  const lines = logContent.value.split('\n')
  if (!searchKeyword.value) return lines
  const kw = searchKeyword.value.toLowerCase()
  return lines.filter(line => line.toLowerCase().includes(kw))
})

const filteredLineCount = computed(() => filteredLines.value.length)

const displayContent = computed(() => {
  const lines = filteredLines.value
  if (!lines.length) return ''
  const text = lines.join('\n')
  if (!searchKeyword.value) {
    // escape HTML
    return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
  }
  // escape HTML first, then highlight
  const escaped = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
  const kwEscaped = searchKeyword.value.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
  const regex = new RegExp(kwEscaped.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi')
  return escaped.replace(regex, match => `<mark class="log-highlight">${match}</mark>`)
})

// 是否选择的是日期目录（非根目录文件）
const isDateSelected = computed(() => {
  return selectedDate.value && !selectedDate.value.startsWith('root::')
})

// 当前日期目录下的类型列表
const currentTypes = computed(() => {
  if (!isDateSelected.value) return []
  const found = dates.value.find(d => d.date === selectedDate.value)
  return found ? found.types : []
})

// 是否已做出完整选择
const hasSelection = computed(() => {
  if (!selectedDate.value) return false
  if (selectedDate.value.startsWith('root::')) return true
  return !!selectedType.value
})

// 当前文件标签
const currentFileLabel = computed(() => {
  if (!selectedDate.value) return '—'
  if (selectedDate.value.startsWith('root::')) {
    return selectedDate.value.replace('root::', '')
  }
  if (selectedType.value) {
    return `${selectedDate.value} / ${selectedType.value}.log`
  }
  return selectedDate.value
})

// 加载日期列表
const loadDates = async () => {
  try {
    const res = await getLogDates()
    dates.value = res.data?.dates || []
    rootFiles.value = res.data?.root_files || []
  } catch (e) {
    console.error('获取日志日期失败', e)
    ElMessage.error(t('admin.logs.loadDatesFailed'))
  }
}

// 同时刷新日期列表和日志内容
const refreshAll = async () => {
  await loadDates()
  loadLog()
}

// 当日期变化时，重置类型并尝试加载
const onDateChange = () => {
  selectedType.value = ''
  logContent.value = ''
  lineCount.value = 0

  if (selectedDate.value.startsWith('root::')) {
    // 根目录文件，直接加载
    loadLog()
  } else if (currentTypes.value.length === 1) {
    // 只有一种类型时自动选择
    selectedType.value = currentTypes.value[0]
    loadLog()
  }
}

// 加载日志内容
const loadLog = async () => {
  if (!hasSelection.value) return

  let params
  if (selectedDate.value.startsWith('root::')) {
    params = { file: selectedDate.value.replace('root::', ''), tail: tailLines.value }
  } else {
    params = { date: selectedDate.value, file: selectedType.value, tail: tailLines.value }
  }

  loading.value = true
  try {
    const res = await getLogContent(params)
    logContent.value = res.data?.content || ''
    lineCount.value = res.data?.lines || 0
    // 滚动到底部
    nextTick(() => {
      if (logContainerRef.value) {
        logContainerRef.value.scrollTop = logContainerRef.value.scrollHeight
      }
    })
  } catch (e) {
    logContent.value = ''
    lineCount.value = 0
    if (e?.response?.status === 404 || e?.message?.includes('not found') || e?.message?.includes('不存在')) {
      ElMessage.warning(t('admin.logs.fileNotFound'))
    } else {
      ElMessage.error(t('admin.logs.loadFailed'))
    }
  } finally {
    loading.value = false
  }
}

// 复制日志内容
const copyContent = async () => {
  if (!logContent.value) return
  await copyToClipboard(logContent.value, t('admin.logs.copySuccess'))
}

onMounted(async () => {
  await loadDates()
  // 默认选择最新日期 + 第一种类型
  if (dates.value.length > 0) {
    selectedDate.value = dates.value[0].date
    if (dates.value[0].types.length > 0) {
      selectedType.value = dates.value[0].types[0]
      loadLog()
    }
  } else if (rootFiles.value.length > 0) {
    selectedDate.value = 'root::' + rootFiles.value[0]
    loadLog()
  }
})

</script>

<style lang="scss" scoped>
.log-viewer {
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 16px;

  .header-card {
    .header-content {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .title-section {
        h2 {
          display: flex;
          align-items: center;
          gap: 8px;
          margin: 0 0 4px 0;
          font-size: 20px;
          font-weight: 600;
          color: var(--el-text-color-primary);
        }
        .subtitle {
          margin: 0;
          font-size: 13px;
          color: var(--el-text-color-secondary);
        }
      }
    }
  }

  .toolbar-card {
    :deep(.el-card__body) {
      padding: 16px 20px;
    }

    .btn-group {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }

  }

  .content-card {
    flex: 1;

    .content-header {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .file-label {
        display: flex;
        align-items: center;
        gap: 6px;
        font-family: monospace;
        font-size: 13px;
        color: var(--el-text-color-regular);
      }

      .line-count {
        font-size: 12px;
        color: var(--el-text-color-secondary);
        background: var(--el-fill-color-light);
        padding: 2px 8px;
        border-radius: 10px;
      }
    }

    .log-container {
      background: #0d1117;
      border-radius: 6px;
      min-height: 500px;
      max-height: calc(100vh - 380px);
      overflow-y: auto;
      display: flex;
      flex-direction: column;

      .log-placeholder {
        flex: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 12px;
        color: #6e7681;
        min-height: 300px;
        font-size: 14px;

        .el-icon {
          font-size: 32px;
        }
      }

      .log-pre {
        margin: 0;
        padding: 16px;
        font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', Consolas, monospace;
        font-size: 12px;
        line-height: 1.6;
        color: #c9d1d9;
        white-space: pre-wrap;
        word-break: break-all;
        flex: 1;

        :deep(.log-highlight) {
          background: #e2b714;
          color: #0d1117;
          border-radius: 2px;
          padding: 0 1px;
        }
      }
    }
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
