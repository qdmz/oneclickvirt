<template>
  <el-form
    v-loading="loading"
    :model="config"
    label-width="140px"
    class="config-form"
  >
    <!-- Logo 设置 -->
    <el-divider content-position="left">
      {{ $t('admin.config.logoSettings') }}
    </el-divider>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-form-item :label="$t('admin.config.customLogoURL')">
          <el-input
            v-model="config.other.logoURL"
            :placeholder="$t('admin.config.customLogoURLPlaceholder')"
            clearable
          />
          <div class="form-item-hint">
            {{ $t('admin.config.customLogoURLHint') }}
          </div>
        </el-form-item>
      </el-col>
      <el-col
        v-if="config.other.logoURL"
        :span="8"
      >
        <el-form-item :label="$t('admin.config.logoPreview')">
          <img
            :src="config.other.logoURL"
            alt="Logo Preview"
            style="height:40px;max-width:160px;object-fit:contain;border:1px solid var(--border-color);border-radius:4px;padding:4px;background:var(--bg-color-secondary);"
            @error="$event.target.style.display='none'"
            @load="$event.target.style.display=''"
          >
        </el-form-item>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-form-item :label="$t('admin.config.customSiteName')">
          <el-input
            v-model="config.other.siteName"
            :placeholder="$t('admin.config.customSiteNamePlaceholder')"
            clearable
          />
          <div class="form-item-hint">
            {{ $t('admin.config.customSiteNameHint') }}
          </div>
        </el-form-item>
      </el-col>
    </el-row>

    <!-- 站点内容设置：标语 / 页脚 / 备案 -->
    <el-divider content-position="left">
      站点内容设置
    </el-divider>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-form-item label="站点标语">
          <el-input
            v-model="config.other.slogan"
            placeholder="首页副标题，如：一站式虚拟化资源管理与销售平台"
            clearable
          />
          <div class="form-item-hint">
            显示在首页主标题下方，留空使用默认标语
          </div>
        </el-form-item>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-form-item label="页脚文字">
          <el-input
            v-model="config.other.footerText"
            placeholder="自定义页脚文字，如公司名称或版权信息"
            clearable
          />
          <div class="form-item-hint">
            显示在页脚区域，留空则不显示
          </div>
        </el-form-item>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="10">
        <el-form-item label="备案号">
          <el-input
            v-model="config.other.icpText"
            placeholder="如：京ICP备xxxxxxxx号"
            clearable
          />
          <div class="form-item-hint">
            显示在页脚，留空则不显示
          </div>
        </el-form-item>
      </el-col>
      <el-col :span="10">
        <el-form-item label="备案链接">
          <el-input
            v-model="config.other.icpLink"
            placeholder="如：https://beian.miit.gov.cn/"
            clearable
          />
          <div class="form-item-hint">
            备案号点击跳转地址，留空则不可点击
          </div>
        </el-form-item>
      </el-col>
    </el-row>

    <el-divider content-position="left">
      {{ $t('admin.config.languageSettings') }}
    </el-divider>

    <el-alert
      type="warning"
      :closable="false"
      show-icon
      style="margin-bottom: 20px;"
    >
      <template #title>
        <strong>{{ $t('admin.config.languageForceNote') }}</strong>
      </template>
      {{ $t('admin.config.languageForceDesc') }}
    </el-alert>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="$t('admin.config.defaultLanguage')">
          <el-select
            v-model="config.other.defaultLanguage"
            :placeholder="$t('admin.config.selectDefaultLanguage')"
            style="width: 100%"
            clearable
          >
            <el-option
              value=""
              :label="$t('admin.config.browserLanguage')"
            />
            <el-option
              value="zh-CN"
              label="中文"
            />
            <el-option
              value="en-US"
              label="English"
            />
          </el-select>
          <div class="form-item-hint">
            {{ $t('admin.config.defaultLanguageHint') }}
          </div>
        </el-form-item>
      </el-col>
    </el-row>
  </el-form>
</template>

<script setup>
defineProps({
  config: { type: Object, required: true },
  loading: { type: Boolean, default: false }
})
</script>

<style scoped>
.config-form {
  max-height: 600px;
  overflow-y: auto;
}
.form-item-hint {
  font-size: 12px;
  color: var(--text-color-tertiary);
  margin-top: 4px;
  line-height: 1.4;
}
</style>
