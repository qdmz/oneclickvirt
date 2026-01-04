<template>
  <div class="site-config-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>站点配置</span>
          <el-button
            type="primary"
            @click="handleInitialize"
          >
            初始化默认配置
          </el-button>
        </div>
      </template>

      <el-form
        v-loading="loading"
        :model="configForm"
        label-width="120px"
      >
        <el-divider content-position="left">
          基本信息
        </el-divider>
        <el-form-item label="网站名称">
          <el-input
            v-model="configForm.site_name"
            placeholder="请输入网站名称"
          />
        </el-form-item>
        <el-form-item label="网站URL">
          <el-input
            v-model="configForm.site_url"
            placeholder="请输入网站URL"
          />
        </el-form-item>
        <el-form-item label="网站图标">
          <el-input
            v-model="configForm.site_icon_url"
            placeholder="请输入网站图标URL"
          />
        </el-form-item>

        <el-divider content-position="left">
          页面内容
        </el-divider>
        <el-form-item label="页眉内容">
          <el-input
            v-model="configForm.site_header"
            type="textarea"
            :rows="3"
            placeholder="请输入页眉内容"
          />
        </el-form-item>
        <el-form-item label="页脚内容">
          <el-input
            v-model="configForm.site_footer"
            type="textarea"
            :rows="3"
            placeholder="请输入页脚内容"
          />
        </el-form-item>

        <el-divider content-position="left">
          联系方式
        </el-divider>
        <el-form-item label="联系邮箱">
          <el-input
            v-model="configForm.contact_email"
            placeholder="请输入联系邮箱"
          />
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input
            v-model="configForm.contact_phone"
            placeholder="请输入联系电话"
          />
        </el-form-item>

        <el-divider content-position="left">
          公司信息
        </el-divider>
        <el-form-item label="公司名称">
          <el-input
            v-model="configForm.company_name"
            placeholder="请输入公司名称"
          />
        </el-form-item>
        <el-form-item label="ICP备案号">
          <el-input
            v-model="configForm.icp_number"
            placeholder="请输入ICP备案号"
          />
        </el-form-item>

        <el-divider content-position="left">
          高级设置
        </el-divider>
        <el-form-item label="自定义CSS">
          <el-input
            v-model="configForm.custom_css"
            type="textarea"
            :rows="5"
            placeholder="请输入自定义CSS代码"
          />
        </el-form-item>
        <el-form-item label="统计代码">
          <el-input
            v-model="configForm.analytics_code"
            type="textarea"
            :rows="3"
            placeholder="请输入统计代码"
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="saving"
            @click="handleSave"
          >
            保存配置
          </el-button>
          <el-button @click="handleReset">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getSiteConfigs, updateSiteConfigs, initializeSiteConfigs } from '@/api/admin'

const loading = ref(false)
const saving = ref(false)
const configForm = ref({
  site_name: '',
  site_url: '',
  site_icon_url: '',
  site_header: '',
  site_footer: '',
  contact_email: '',
  contact_phone: '',
  company_name: '',
  icp_number: '',
  custom_css: '',
  analytics_code: ''
})

// 加载配置
const loadConfigs = async () => {
  loading.value = true
  try {
    const res = await getSiteConfigs()
    if (res.code === 200) {
      const configs = res.data || []
      configs.forEach(config => {
        if (configForm.value.hasOwnProperty(config.key)) {
          configForm.value[config.key] = config.value
        }
      })
    }
  } catch (error) {
    ElMessage.error('加载配置失败')
  } finally {
    loading.value = false
  }
}

// 保存配置
const handleSave = async () => {
  saving.value = true
  try {
    const configs = Object.keys(configForm.value).map(key => ({
      key,
      value: configForm.value[key],
      type: key.includes('_css') || key.includes('_code') ? 'text' : 'string'
    }))

    const res = await updateSiteConfigs(configs)
    if (res.code === 200) {
      ElMessage.success('保存成功')
    } else {
      ElMessage.error(res.message || '保存失败')
    }
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 初始化默认配置
const handleInitialize = async () => {
  try {
    await ElMessageBox.confirm('确定要初始化默认配置吗？这将覆盖现有配置。', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const res = await initializeSiteConfigs()
    if (res.code === 200) {
      ElMessage.success('初始化成功')
      loadConfigs()
    } else {
      ElMessage.error(res.message || '初始化失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('初始化失败')
    }
  }
}

// 重置
const handleReset = () => {
  loadConfigs()
}

onMounted(() => {
  loadConfigs()
})
</script>

<style scoped>
.site-config-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
