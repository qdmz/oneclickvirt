<template>
  <div class="site-config-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('admin.siteConfig.title') }}</span>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="140px"
        class="config-form"
      >
        <!-- 基础信息配置 -->
        <el-divider content-position="left">
          {{ t('admin.siteConfig.basicInfo') }}
        </el-divider>

        <el-form-item :label="t('admin.siteConfig.siteName')" prop="site_name">
          <el-input
            v-model="form.site_name"
            :placeholder="t('admin.siteConfig.siteNamePlaceholder')"
            maxlength="50"
            show-word-limit
            clearable
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.siteDescription')" prop="site_description">
          <el-input
            v-model="form.site_description"
            type="textarea"
            :rows="2"
            :placeholder="t('admin.siteConfig.siteDescriptionPlaceholder')"
            maxlength="200"
            show-word-limit
            clearable
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.siteKeywords')" prop="site_keywords">
          <el-input
            v-model="form.site_keywords"
            :placeholder="t('admin.siteConfig.siteKeywordsPlaceholder')"
            maxlength="200"
            show-word-limit
            clearable
          />
          <div class="form-hint">{{ t('admin.siteConfig.keywordsHint') }}</div>
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.contactEmail')" prop="contact_email">
          <el-input
            v-model="form.contact_email"
            :placeholder="t('admin.siteConfig.contactEmailPlaceholder')"
            clearable
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.icp')" prop="icp">
          <el-input
            v-model="form.icp"
            :placeholder="t('admin.siteConfig.icpPlaceholder')"
            clearable
          />
        </el-form-item>

        <!-- Logo 配置 -->
        <el-divider content-position="left">
          {{ t('admin.siteConfig.logoConfig') }}
        </el-divider>

        <el-form-item :label="t('admin.siteConfig.siteLogo')">
          <div class="logo-upload-wrapper">
            <el-upload
              class="logo-uploader"
              action="/v1/admin/site-config/logo"
              :show-file-list="false"
              :before-upload="beforeLogoUpload"
              :on-success="handleLogoSuccess"
              :on-error="handleLogoError"
              accept="image/png,image/jpeg,image/svg+xml"
            >
              <img
                v-if="form.logo_url"
                :src="form.logo_url"
                class="logo-preview"
                alt="Logo"
              >
              <div
                v-else
                class="logo-upload-placeholder"
              >
                <el-icon :size="32">
                  <Plus />
                </el-icon>
                <span>{{ t('admin.siteConfig.uploadLogo') }}</span>
              </div>
            </el-upload>
            <div class="logo-actions">
              <el-button
                v-if="form.logo_url"
                type="danger"
                size="small"
                @click="removeLogo"
              >
                {{ t('admin.siteConfig.removeLogo') }}
              </el-button>
            </div>
          </div>
          <div class="form-hint">{{ t('admin.siteConfig.logoHint') }}</div>
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.favicon')">
          <div class="logo-upload-wrapper">
            <el-upload
              class="favicon-uploader"
              action="/v1/admin/site-config/favicon"
              :show-file-list="false"
              :before-upload="beforeFaviconUpload"
              :on-success="handleFaviconSuccess"
              :on-error="handleFaviconError"
              accept="image/x-icon,image/png"
            >
              <img
                v-if="form.favicon_url"
                :src="form.favicon_url"
                class="favicon-preview"
                alt="Favicon"
              >
              <div
                v-else
                class="logo-upload-placeholder"
              >
                <el-icon :size="24">
                  <Plus />
                </el-icon>
                <span>{{ t('admin.siteConfig.uploadFavicon') }}</span>
              </div>
            </el-upload>
            <div class="logo-actions">
              <el-button
                v-if="form.favicon_url"
                type="danger"
                size="small"
                @click="removeFavicon"
              >
                {{ t('admin.siteConfig.removeFavicon') }}
              </el-button>
            </div>
          </div>
          <div class="form-hint">{{ t('admin.siteConfig.faviconHint') }}</div>
        </el-form-item>

        <!-- 主题色配置 -->
        <el-divider content-position="left">
          {{ t('admin.siteConfig.themeConfig') }}
        </el-divider>

        <el-form-item :label="t('admin.siteConfig.primaryColor')" prop="primary_color">
          <div class="color-picker-wrapper">
            <el-color-picker
              v-model="form.primary_color"
              show-alpha
              :predefine="predefineColors"
            />
            <el-input
              v-model="form.primary_color"
              style="width: 200px; margin-left: 12px;"
              placeholder="#16a34a"
            />
          </div>
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.customCss')" prop="custom_css">
          <el-input
            v-model="form.custom_css"
            type="textarea"
            :rows="4"
            :placeholder="t('admin.siteConfig.customCssPlaceholder')"
          />
          <div class="form-hint">{{ t('admin.siteConfig.customCssHint') }}</div>
        </el-form-item>

        <!-- 页眉页脚配置 -->
        <el-divider content-position="left">
          {{ t('admin.siteConfig.headerFooterConfig') }}
        </el-divider>

        <el-form-item :label="t('admin.siteConfig.customHeader')">
          <div class="html-editor-wrapper">
            <el-tabs v-model="headerEditorMode">
              <el-tab-pane
                :label="t('admin.siteConfig.htmlMode')"
                name="html"
              >
                <el-input
                  v-model="form.custom_header"
                  type="textarea"
                  :rows="6"
                  :placeholder="t('admin.siteConfig.customHeaderPlaceholder')"
                />
              </el-tab-pane>
              <el-tab-pane
                :label="t('admin.siteConfig.previewMode')"
                name="preview"
              >
                <div
                  class="html-preview"
                  v-html="form.custom_header"
                />
              </el-tab-pane>
            </el-tabs>
          </div>
          <div class="form-hint">{{ t('admin.siteConfig.customHeaderHint') }}</div>
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.customFooter')">
          <div class="html-editor-wrapper">
            <el-tabs v-model="footerEditorMode">
              <el-tab-pane
                :label="t('admin.siteConfig.htmlMode')"
                name="html"
              >
                <el-input
                  v-model="form.custom_footer"
                  type="textarea"
                  :rows="6"
                  :placeholder="t('admin.siteConfig.customFooterPlaceholder')"
                />
              </el-tab-pane>
              <el-tab-pane
                :label="t('admin.siteConfig.previewMode')"
                name="preview"
              >
                <div
                  class="html-preview"
                  v-html="form.custom_footer"
                />
              </el-tab-pane>
            </el-tabs>
          </div>
          <div class="form-hint">{{ t('admin.siteConfig.customFooterHint') }}</div>
        </el-form-item>

        <!-- 功能开关 -->
        <el-divider content-position="left">
          {{ t('admin.siteConfig.featureSwitches') }}
        </el-divider>

        <el-form-item :label="t('admin.siteConfig.enableRegister')">
          <el-switch
            v-model="form.enable_register"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.enableStore')">
          <el-switch
            v-model="form.enable_store"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.enableRecharge')">
          <el-switch
            v-model="form.enable_recharge"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.enableTicket')">
          <el-switch
            v-model="form.enable_ticket"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.enableAnnouncement')">
          <el-switch
            v-model="form.enable_announcement"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.maintenanceMode')">
          <el-switch
            v-model="form.maintenance_mode"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item
          v-if="form.maintenance_mode"
          :label="t('admin.siteConfig.maintenanceMessage')"
        >
          <el-input
            v-model="form.maintenance_message"
            type="textarea"
            :rows="2"
            :placeholder="t('admin.siteConfig.maintenanceMessagePlaceholder')"
          />
        </el-form-item>
      </el-form>

      <!-- 底部操作按钮 -->
      <div class="form-actions">
        <el-button
          type="primary"
          size="large"
          :loading="submitting"
          @click="handleSubmit"
        >
          {{ t('common.save') }}
        </el-button>
        <el-button
          size="large"
          @click="handleReset"
        >
          {{ t('common.reset') }}
        </el-button>
        <el-button
          type="success"
          size="large"
          @click="handlePreview"
        >
          {{ t('admin.siteConfig.preview') }}
        </el-button>
      </div>
    </el-card>

    <!-- 预览弹窗 -->
    <el-dialog
      v-model="previewVisible"
      :title="t('admin.siteConfig.preview')"
      width="800px"
      destroy-on-close
    >
      <div class="preview-section">
        <h4>{{ t('admin.siteConfig.customHeaderPreview') }}</h4>
        <div
          class="html-preview-box"
          v-html="form.custom_header || t('admin.siteConfig.noContent')"
        />
      </div>
      <div class="preview-section">
        <h4>{{ t('admin.siteConfig.customFooterPreview') }}</h4>
        <div
          class="html-preview-box"
          v-html="form.custom_footer || t('admin.siteConfig.noContent')"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getAdminSiteConfig, updateAdminSiteConfig } from '@/api/site'

const { t } = useI18n()

const formRef = ref(null)
const submitting = ref(false)
const loading = ref(false)
const headerEditorMode = ref('html')
const footerEditorMode = ref('html')
const previewVisible = ref(false)

// 预定义颜色
const predefineColors = [
  '#16a34a',
  '#3b82f6',
  '#f59e0b',
  '#ef4444',
  '#8b5cf6',
  '#ec4899',
  '#06b6d4',
  '#6366f1'
]

// 表单数据
const form = ref({
  site_name: '',
  site_description: '',
  site_keywords: '',
  contact_email: '',
  icp: '',
  logo_url: '',
  favicon_url: '',
  primary_color: '#16a34a',
  custom_css: '',
  custom_header: '',
  custom_footer: '',
  enable_register: true,
  enable_store: true,
  enable_recharge: true,
  enable_ticket: true,
  enable_announcement: true,
  maintenance_mode: false,
  maintenance_message: ''
})

// 表单验证规则
const rules = {
  site_name: [
    { max: 50, message: t('admin.siteConfig.siteNameMaxLength'), trigger: 'blur' }
  ],
  contact_email: [
    { type: 'email', message: t('admin.siteConfig.emailFormatError'), trigger: 'blur' }
  ],
  primary_color: [
    { pattern: /^#[0-9a-fA-F]{6}([0-9a-fA-F]{2})?$/, message: t('admin.siteConfig.colorFormatError'), trigger: 'blur' }
  ]
}

// 加载站点配置
const loadConfig = async () => {
  loading.value = true
  try {
    const res = await getAdminSiteConfig()
    if (res && res.code === 200 && res.data) {
      const data = res.data
      form.value = {
        site_name: data.site_name || '',
        site_description: data.site_description || '',
        site_keywords: data.site_keywords || '',
        contact_email: data.contact_email || '',
        icp: data.icp || '',
        logo_url: data.logo_url || '',
        favicon_url: data.favicon_url || '',
        primary_color: data.primary_color || '#16a34a',
        custom_css: data.custom_css || '',
        custom_header: data.custom_header || '',
        custom_footer: data.custom_footer || '',
        enable_register: data.enable_register !== false,
        enable_store: data.enable_store !== false,
        enable_recharge: data.enable_recharge !== false,
        enable_ticket: data.enable_ticket !== false,
        enable_announcement: data.enable_announcement !== false,
        maintenance_mode: data.maintenance_mode === true,
        maintenance_message: data.maintenance_message || ''
      }
    }
  } catch (error) {
    ElMessage.error(error?.message || t('admin.siteConfig.loadFailed'))
  } finally {
    loading.value = false
  }
}

// Logo 上传前检查
const beforeLogoUpload = (file) => {
  const isImage = file.type === 'image/png' || file.type === 'image/jpeg' || file.type === 'image/svg+xml'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error(t('admin.siteConfig.logoFormatError'))
  }
  if (!isLt2M) {
    ElMessage.error(t('admin.siteConfig.logoSizeError'))
  }
  return isImage && isLt2M
}

// Logo 上传成功
const handleLogoSuccess = (res) => {
  if (res && res.code === 200 && res.data?.url) {
    form.value.logo_url = res.data.url
    ElMessage.success(t('admin.siteConfig.logoUploadSuccess'))
  } else {
    ElMessage.error(res?.message || t('admin.siteConfig.logoUploadFailed'))
  }
}

// Logo 上传失败
const handleLogoError = () => {
  ElMessage.error(t('admin.siteConfig.logoUploadFailed'))
}

// 移除 Logo
const removeLogo = () => {
  form.value.logo_url = ''
  ElMessage.success(t('admin.siteConfig.logoRemoved'))
}

// Favicon 上传前检查
const beforeFaviconUpload = (file) => {
  const isValid = file.type === 'image/x-icon' || file.type === 'image/png' || file.name.endsWith('.ico')
  const isLt1M = file.size / 1024 / 1024 < 1

  if (!isValid) {
    ElMessage.error(t('admin.siteConfig.faviconFormatError'))
  }
  if (!isLt1M) {
    ElMessage.error(t('admin.siteConfig.faviconSizeError'))
  }
  return isValid && isLt1M
}

// Favicon 上传成功
const handleFaviconSuccess = (res) => {
  if (res && res.code === 200 && res.data?.url) {
    form.value.favicon_url = res.data.url
    ElMessage.success(t('admin.siteConfig.faviconUploadSuccess'))
  } else {
    ElMessage.error(res?.message || t('admin.siteConfig.faviconUploadFailed'))
  }
}

// Favicon 上传失败
const handleFaviconError = () => {
  ElMessage.error(t('admin.siteConfig.faviconUploadFailed'))
}

// 移除 Favicon
const removeFavicon = () => {
  form.value.favicon_url = ''
  ElMessage.success(t('admin.siteConfig.faviconRemoved'))
}

// 提交表单
const handleSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const res = await updateAdminSiteConfig(form.value)
    if (res && res.code === 200) {
      ElMessage.success(t('admin.siteConfig.saveSuccess'))
      // 刷新全局站点配置
      const siteStore = (await import('@/pinia/modules/site')).useSiteStore()
      await siteStore.refresh()
    } else {
      ElMessage.error(res?.message || t('admin.siteConfig.saveFailed'))
    }
  } catch (error) {
    ElMessage.error(error?.message || t('admin.siteConfig.saveFailed'))
  } finally {
    submitting.value = false
  }
}

// 重置表单
const handleReset = () => {
  loadConfig()
  ElMessage.info(t('admin.siteConfig.resetSuccess'))
}

// 预览
const handlePreview = () => {
  previewVisible.value = true
}

onMounted(() => {
  loadConfig()
})
</script>

<style scoped lang="scss">
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

.config-form {
  max-width: 800px;
  margin: 0 auto;

  :deep(.el-divider__text) {
    font-size: 15px;
    font-weight: 600;
    color: var(--text-color-primary);
  }
}

.form-hint {
  font-size: 12px;
  color: var(--text-color-tertiary);
  margin-top: 4px;
  line-height: 1.4;
}

.logo-upload-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.logo-uploader,
.favicon-uploader {
  border: 2px dashed var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.3s;
  overflow: hidden;

  &:hover {
    border-color: var(--primary-color);
  }
}

.logo-uploader {
  width: 180px;
  height: 60px;
}

.favicon-uploader {
  width: 64px;
  height: 64px;
}

.logo-preview {
  width: 100%;
  height: 100%;
  object-fit: contain;
  padding: 4px;
}

.favicon-preview {
  width: 100%;
  height: 100%;
  object-fit: contain;
  padding: 4px;
}

.logo-upload-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--text-color-tertiary);
  gap: 4px;

  span {
    font-size: 12px;
  }
}

.logo-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.color-picker-wrapper {
  display: flex;
  align-items: center;
}

.html-editor-wrapper {
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 8px;

  :deep(.el-tabs__header) {
    margin-bottom: 8px;
  }

  :deep(.el-textarea__inner) {
    font-family: 'Courier New', monospace;
  }
}

.html-preview {
  min-height: 120px;
  padding: 12px;
  background-color: var(--bg-color-secondary);
  border-radius: 4px;
  border: 1px solid var(--border-color);
}

.form-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 24px 0 8px;
  border-top: 1px solid var(--border-color);
  margin-top: 24px;
}

.preview-section {
  margin-bottom: 20px;

  h4 {
    margin-bottom: 12px;
    font-size: 14px;
    color: var(--text-color-secondary);
  }
}

.html-preview-box {
  min-height: 80px;
  padding: 16px;
  background-color: var(--bg-color-secondary);
  border-radius: 4px;
  border: 1px solid var(--border-color);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .config-form {
    max-width: 100%;
  }

  .logo-upload-wrapper {
    flex-direction: column;
    align-items: flex-start;
  }

  .form-actions {
    flex-direction: column;
    align-items: center;

    .el-button {
      width: 100%;
      max-width: 200px;
    }
  }

  .color-picker-wrapper {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>
