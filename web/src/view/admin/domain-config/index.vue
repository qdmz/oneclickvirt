<template>
  <div class="domain-config-container">
    <el-card shadow="never">
      <template #header>
        <span>域名系统配置</span>
      </template>

      <el-form :model="config" label-width="180px" v-loading="loading" style="max-width: 600px;">
        <el-form-item label="每用户最大域名数">
          <el-input-number v-model="config.maxDomainsPerUser" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="代理商子用户最大域名数">
          <el-input-number v-model="config.maxDomainsPerAgentUser" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="DNS TTL (秒)">
          <el-input-number v-model="config.defaultTTL" :min="60" :max="86400" :step="60" />
        </el-form-item>
        <el-form-item label="自动SSL">
          <el-switch v-model="config.autoSSL" />
        </el-form-item>
        <el-form-item label="允许的域名后缀">
          <el-input v-model="config.allowedSuffixes" placeholder="如: .com,.net (空=不限)" />
        </el-form-item>
        <el-form-item label="DNS类型">
          <el-radio-group v-model="config.dnsType">
            <el-radio value="dnsmasq">dnsmasq</el-radio>
            <el-radio value="hosts">hosts</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="DNS配置路径">
          <el-input v-model="config.dnsConfigPath" />
        </el-form-item>
        <el-form-item label="Nginx配置路径">
          <el-input v-model="config.nginxConfigPath" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saving">保存配置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getDomainConfig, updateDomainConfig } from '@/api/domain'

const loading = ref(false)
const saving = ref(false)
const config = reactive({
  maxDomainsPerUser: 3,
  maxDomainsPerAgentUser: 5,
  defaultTTL: 300,
  autoSSL: false,
  allowedSuffixes: '',
  dnsType: 'dnsmasq',
  dnsConfigPath: '/etc/dnsmasq.d/oneclickvirt-hosts.conf',
  nginxConfigPath: '/etc/nginx/conf.d/oneclickvirt-domains'
})

async function fetchConfig() {
  loading.value = true
  try {
    const res = await getDomainConfig()
    if (res.code === 0 && res.data) Object.assign(config, res.data)
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    const res = await updateDomainConfig({ ...config })
    if (res.code === 0 && res.data) Object.assign(config, res.data)
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(fetchConfig)
</script>

<style scoped>
.domain-config-container { padding: 20px; }
</style>
