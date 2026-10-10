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

        <el-form-item :label="t('admin.siteConfig.icp')" prop="icp_number">
          <el-input
            v-model="form.icp_number"
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
              :show-file-list="false"
              :before-upload="beforeLogoUpload"
              :http-request="customLogoUpload"
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
              :show-file-list="false"
              :before-upload="beforeFaviconUpload"
              :http-request="customFaviconUpload"
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
              placeholder="#4f46e5"
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
            v-model="form.enable_registration"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.enableStore')">
          <el-switch
            v-model="form.enable_product_store"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.enableRecharge')">
          <el-switch
            v-model="form.show_yipay"
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
            v-model="form.announcement_enabled"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <!-- 首页设置 -->
        <el-divider content-position="left">
          {{ t('admin.siteConfig.homepageSettings') }}
        </el-divider>

        <el-form-item :label="t('admin.siteConfig.homeTitle')">
          <el-input v-model="form.home_title" :placeholder="t('admin.siteConfig.homeTitlePlaceholder')" maxlength="256" />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.homeSubtitle')">
          <el-input v-model="form.home_subtitle" :placeholder="t('admin.siteConfig.homeSubtitlePlaceholder')" maxlength="512" />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.homeBackground')">
          <el-input v-model="form.home_background" placeholder="https://..." maxlength="512" />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.showHomeStats')">
          <el-switch
            v-model="form.show_home_stats"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-divider content-position="left">
          {{ t('admin.siteConfig.homepageSections') }}
        </el-divider>

        <el-form-item :label="t('admin.siteConfig.showPlatforms')">
          <el-switch
            v-model="form.show_platforms"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.showSponsors')">
          <el-switch
            v-model="form.show_sponsors"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.showRecommended')">
          <el-switch
            v-model="form.show_recommended"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.recommendedLimit')">
          <el-input-number v-model="form.recommended_limit" :min="1" :max="20" />
        </el-form-item>

        <el-divider content-position="left">
          {{ t('admin.siteConfig.recommendedSettings') }}
        </el-divider>

        <el-form-item :label="t('admin.siteConfig.recommendedTitle')">
          <el-input v-model="form.recommended_title" :placeholder="t('admin.siteConfig.recommendedTitlePlaceholder')" maxlength="256" />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.recommendedSubtitle')">
          <el-input v-model="form.recommended_subtitle" :placeholder="t('admin.siteConfig.recommendedSubtitlePlaceholder')" maxlength="512" />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.recommendedCols')">
          <el-select v-model="form.recommended_cols" style="width: 160px;">
            <el-option label="2" :value="2" />
            <el-option label="3" :value="3" />
            <el-option label="4" :value="4" />
            <el-option label="6" :value="6" />
          </el-select>
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.recommendedShowPrice')">
          <el-switch
            v-model="form.recommended_show_price"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.recommendedShowSpecs')">
          <el-switch
            v-model="form.recommended_show_specs"
            :active-text="t('admin.siteConfig.on')"
            :inactive-text="t('admin.siteConfig.off')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.siteConfig.copyrightText')">
          <el-input v-model="form.copyright_text" :placeholder="t('admin.siteConfig.copyrightTextPlaceholder')" maxlength="512" />
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
import { uploadSiteImage } from '@/api/admin'

const { t } = useI18n()

const formRef = ref(null)
const submitting = ref(false)
const loading = ref(false)
const headerEditorMode = ref('html')
const footerEditorMode = ref('html')
const previewVisible = ref(false)

// 预定义颜色
const predefineColors = [
  '#4f46e5',
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
  icp_number: '',
  logo_url: '',
  favicon_url: '',
  primary_color: '#4f46e5',
  custom_css: '',
  custom_header: '',
  custom_footer: '',
  enable_registration: true,
  enable_product_store: true,
  show_yipay: false,
  enable_ticket: true,
  announcement_enabled: false,
  home_title: '',
  home_subtitle: '',
  home_background: '',
  show_home_stats: true,
  show_platforms: true,
  show_sponsors: true,
  show_recommended: true,
  recommended_limit: 8,
  recommended_title: '',
  recommended_subtitle: '',
  recommended_cols: 4,
  recommended_show_price: true,
  recommended_show_specs: true,
  copyright_text: ''
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
    { pattern: /^(#[0-9a-fA-F]{6}([0-9a-fA-F]{2})?|rgba?\([^)]+\))$/, message: t('admin.siteConfig.colorFormatError'), trigger: 'blur' }
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
        icp_number: data.icp_number || '',
        logo_url: data.logo_url || '',
        favicon_url: data.favicon_url || '',
        primary_color: data.primary_color || '#4f46e5',
        custom_css: data.custom_css || '',
        custom_header: data.custom_header || '',
        custom_footer: data.custom_footer || '',
        enable_registration: data.enable_registration !== false,
        enable_product_store: data.enable_product_store !== false,
        show_yipay: data.show_yipay === true,
        enable_ticket: data.enable_ticket !== false,
        announcement_enabled: data.announcement_enabled === true,
        home_title: data.home_title || '',
        home_subtitle: data.home_subtitle || '',
        home_background: data.home_background || '',
        show_home_stats: data.show_home_stats !== false,
        show_platforms: data.show_platforms !== false,
        show_sponsors: data.show_sponsors !== false,
        show_recommended: data.show_recommended !== false,
        recommended_limit: Number(data.recommended_limit) > 0 ? Number(data.recommended_limit) : 8,
        recommended_title: data.recommended_title || '',
        recommended_subtitle: data.recommended_subtitle || '',
        recommended_cols: Number(data.recommended_cols) > 0 ? Number(data.recommended_cols) : 4,
        recommended_show_price: data.recommended_show_price !== false,
        recommended_show_specs: data.recommended_show_specs !== false,
        copyright_text: data.copyright_text || ''
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

// 自定义 Logo 上传
const customLogoUpload = async (options) => {
  const { file } = options
  const formData = new FormData()
  formData.append('image', file)
  formData.append('type', 'logo')
  try {
    const res = await uploadSiteImage(formData)
    if (res && res.code === 200 && res.data?.url) {
      form.value.logo_url = res.data.url
      ElMessage.success(t('admin.siteConfig.logoUploadSuccess'))
    } else {
      ElMessage.error(res?.message || t('admin.siteConfig.logoUploadFailed'))
    }
  } catch (error) {
    ElMessage.error(error?.message || t('admin.siteConfig.logoUploadFailed'))
  }
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

// 自定义 Favicon 上传
const customFaviconUpload = async (options) => {
  const { file } = options
  const formData = new FormData()
  formData.append('image', file)
  formData.append('type', 'favicon')
  try {
    const res = await uploadSiteImage(formData)
    if (res && res.code === 200 && res.data?.url) {
      form.value.favicon_url = res.data.url
      ElMessage.success(t('admin.siteConfig.faviconUploadSuccess'))
    } else {
      ElMessage.error(res?.message || t('admin.siteConfig.faviconUploadFailed'))
    }
  } catch (error) {
    ElMessage.error(error?.message || t('admin.siteConfig.faviconUploadFailed'))
  }
}

// 移除 Favicon
const removeFavicon = () => {
  form.value.favicon_url = ''
  ElMessage.success(t('admin.siteConfig.faviconRemoved'))
}

// 提交表单
const handleSubmit = async () => {
  try {
    const valid = await formRef.value?.validate()
    if (!valid) return
  } catch {
    ElMessage.warning(t('admin.siteConfig.validationFailed'))
    return
  }

  submitting.value = true
  try {
    const res = await updateAdminSiteConfig(form.value)
    if (res && res.code === 200) {
      ElMessage.success(t('admin.siteConfig.saveSuccess'))
      // 重新拉取配置回填表单，避免"看似没保存"
      await loadConfig()
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
