<template>
  <div class="config-container">
    <el-card>
      <template #header>
        <div class="config-header">
          <span>{{ $t('admin.config.title') }}</span>
        </div>
      </template>
      
      <!-- 配置分类标签页 -->
      <el-tabs
        v-model="activeTab"
        type="border-card"
        class="config-tabs"
      >
        <!-- 基础认证配置 -->
        <el-tab-pane
          :label="$t('admin.config.basicAuth')"
          name="auth"
        >
          <el-form
            v-loading="loading"
            :model="config"
            label-width="140px"
            class="config-form"
          >
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.emailLogin')">
                  <el-switch v-model="config.auth.enableEmail" />
                  <div class="form-item-hint">
                    {{ $t('admin.config.emailLoginHint') }}
                  </div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item
                  :label="$t('admin.config.publicRegistration')"
                  :help="$t('admin.config.publicRegistrationHelp')"
                >
                  <el-switch v-model="config.auth.enablePublicRegistration" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.telegramLogin')">
                  <el-switch v-model="config.auth.enableTelegram" />
                  <div class="form-item-hint">
                    {{ $t('admin.config.telegramLoginHint') }}
                  </div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.qqLogin')">
                  <el-switch v-model="config.auth.enableQQ" />
                  <div class="form-item-hint">
                    {{ $t('admin.config.qqLoginHint') }}
                  </div>
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="OAuth2">
                  <el-switch v-model="config.auth.enableOAuth2" />
                  <div class="form-item-hint">
                    {{ $t('admin.config.oauth2Hint') }}
                  </div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.inviteCodeSystem')">
                  <el-switch v-model="config.inviteCode.enabled" />
                  <div class="form-item-hint">
                    {{ $t('admin.config.inviteCodeSystemHint') }}
                  </div>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <!-- 邮箱SMTP配置 -->
        <el-tab-pane
          :label="$t('admin.config.emailConfig')"
          name="email"
        >
          <el-form
            v-loading="loading"
            :model="config"
            label-width="140px"
            class="config-form"
          >
            <el-alert
              :title="$t('admin.config.smtpConfigDesc')"
              type="info"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            >
              {{ $t('admin.config.smtpConfigHint') }}
            </el-alert>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.smtpHost')">
                  <el-input
                    v-model="config.auth.emailSMTPHost"
                    :placeholder="$t('admin.config.smtpHostPlaceholder')"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.smtpPort')">
                  <el-input-number
                    v-model="config.auth.emailSMTPPort"
                    :min="1"
                    :max="65535"
                    :controls="false"
                    :placeholder="$t('admin.config.smtpPortPlaceholder')"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.emailUsername')">
                  <el-input
                    v-model="config.auth.emailUsername"
                    :placeholder="$t('admin.config.emailUsernamePlaceholder')"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.emailPassword')">
                  <el-input
                    v-model="config.auth.emailPassword"
                    type="password"
                    :placeholder="$t('admin.config.emailPasswordPlaceholder')"
                    show-password
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <!-- 第三方登录配置 -->
        <el-tab-pane
          :label="$t('admin.config.thirdPartyLogin')"
          name="oauth"
        >
          <el-form
            v-loading="loading"
            :model="config"
            label-width="140px"
            class="config-form"
          >
            <!-- Telegram配置 -->
            <el-card
              class="oauth-card"
              shadow="never"
            >
              <template #header>
                <div class="oauth-header">
                  <span>{{ $t('admin.config.telegramConfig') }}</span>
                  <el-switch v-model="config.auth.enableTelegram" />
                </div>
              </template>
              <el-form-item label="Bot Token">
                <el-input
                  v-model="config.auth.telegramBotToken"
                  :placeholder="$t('admin.config.telegramBotTokenPlaceholder')"
                  :disabled="!config.auth.enableTelegram"
                />
              </el-form-item>
            </el-card>

            <!-- QQ配置 -->
            <el-card
              class="oauth-card"
              shadow="never"
            >
              <template #header>
                <div class="oauth-header">
                  <span>{{ $t('admin.config.qqConfig') }}</span>
                  <el-switch v-model="config.auth.enableQQ" />
                </div>
              </template>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="App ID">
                    <el-input
                      v-model="config.auth.qqAppID"
                      :placeholder="$t('admin.config.qqAppIdPlaceholder')"
                      :disabled="!config.auth.enableQQ"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="App Key">
                    <el-input
                      v-model="config.auth.qqAppKey"
                      :placeholder="$t('admin.config.qqAppKeyPlaceholder')"
                      :disabled="!config.auth.enableQQ"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
            </el-card>
          </el-form>
        </el-tab-pane>

        <!-- 支付接口配置 -->
        <el-tab-pane
          label="支付接口配置"
          name="payment"
        >
          <el-form
            v-loading="loading"
            :model="config"
            label-width="140px"
            class="config-form"
          >
            <el-alert
              title="支付接口配置说明"
              type="info"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            >
              <div>配置支付宝和微信支付接口,用于订单支付功能</div>
              <div style="margin-top: 8px; color: #E6A23C;">
                <i class="el-icon-warning" />
                请确保填写正确的支付接口信息,否则支付功能将无法正常使用
              </div>
            </el-alert>

            <!-- 支付宝配置 -->
            <el-card
              class="payment-card"
              shadow="never"
            >
              <template #header>
                <div class="payment-header">
                  <span>支付宝配置</span>
                  <el-switch v-model="config.payment.enableAlipay" />
                </div>
              </template>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="应用ID">
                    <el-input
                      v-model="config.payment.alipayAppId"
                      placeholder="请输入支付宝应用ID"
                      :disabled="!config.payment.enableAlipay"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="商户私钥">
                    <el-input
                      v-model="config.payment.alipayPrivateKey"
                      type="textarea"
                      :rows="2"
                      placeholder="请输入商户私钥"
                      :disabled="!config.payment.enableAlipay"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="支付宝公钥">
                    <el-input
                      v-model="config.payment.alipayPublicKey"
                      type="textarea"
                      :rows="2"
                      placeholder="请输入支付宝公钥"
                      :disabled="!config.payment.enableAlipay"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="网关地址">
                    <el-select
                      v-model="config.payment.alipayGateway"
                      placeholder="选择网关地址"
                      style="width: 100%"
                      :disabled="!config.payment.enableAlipay"
                    >
                      <el-option
                        label="正式环境"
                        value="https://openapi.alipay.com/gateway.do"
                      />
                      <el-option
                        label="沙箱环境"
                        value="https://openapi.alipaydev.com/gateway.do"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
            </el-card>

            <!-- 微信支付配置 -->
            <el-card
              class="payment-card"
              shadow="never"
            >
              <template #header>
                <div class="payment-header">
                  <span>微信支付配置</span>
                  <el-switch v-model="config.payment.enableWechat" />
                </div>
              </template>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="商户号">
                    <el-input
                      v-model="config.payment.wechatMchId"
                      placeholder="请输入微信商户号"
                      :disabled="!config.payment.enableWechat"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="应用ID">
                    <el-input
                      v-model="config.payment.wechatAppId"
                      placeholder="请输入微信应用ID"
                      :disabled="!config.payment.enableWechat"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="API密钥">
                    <el-input
                      v-model="config.payment.wechatApiKey"
                      type="password"
                      placeholder="请输入API密钥"
                      show-password
                      :disabled="!config.payment.enableWechat"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="证书序列号">
                    <el-input
                      v-model="config.payment.wechatSerialNo"
                      placeholder="请输入证书序列号"
                      :disabled="!config.payment.enableWechat"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="APIv3密钥">
                    <el-input
                      v-model="config.payment.wechatApiV3Key"
                      type="password"
                      placeholder="请输入APIv3密钥"
                      show-password
                      :disabled="!config.payment.enableWechat"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="支付类型">
                    <el-select
                      v-model="config.payment.wechatType"
                      placeholder="选择支付类型"
                      style="width: 100%"
                      :disabled="!config.payment.enableWechat"
                    >
                      <el-option
                        label="公众号/小程序"
                        value="mp"
                      />
                      <el-option
                        label="APP支付"
                        value="app"
                      />
                      <el-option
                        label="H5支付"
                        value="h5"
                      />
                      <el-option
                        label="Native支付"
                        value="native"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
            </el-card>

            <!-- 余额支付配置 -->
            <el-card
              class="payment-card"
              shadow="never"
            >
              <template #header>
                <div class="payment-header">
                  <span>余额支付配置</span>
                  <el-switch v-model="config.payment.enableBalance" />
                </div>
              </template>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="最低支付金额">
                    <el-input-number
                      v-model="config.payment.balanceMinAmount"
                      :min="0"
                      :precision="2"
                      :disabled="!config.payment.enableBalance"
                      style="width: 100%"
                    />
                    <span style="margin-left: 10px;">元</span>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="每日限额">
                    <el-input-number
                      v-model="config.payment.balanceDailyLimit"
                      :min="0"
                      :precision="2"
                      :disabled="!config.payment.enableBalance"
                      style="width: 100%"
                    />
                    <span style="margin-left: 10px;">元</span>
                  </el-form-item>
                </el-col>
              </el-row>
            </el-card>

            <!-- 易支付配置 -->
            <el-card
              class="payment-card"
              shadow="never"
            >
              <template #header>
                <div class="payment-header">
                  <span>易支付配置</span>
                  <el-switch v-model="config.payment.enableEpay" />
                </div>
              </template>
              <el-form-item label="API地址">
                <el-input
                  v-model="config.payment.epayAPIURL"
                  :disabled="!config.payment.enableEpay"
                  placeholder="请输入易支付API地址"
                />
              </el-form-item>
              <el-form-item label="商户ID">
                <el-input
                  v-model="config.payment.epayPID"
                  :disabled="!config.payment.enableEpay"
                  placeholder="请输入商户ID"
                />
              </el-form-item>
              <el-form-item label="密钥">
                <el-input
                  v-model="config.payment.epayKey"
                  :disabled="!config.payment.enableEpay"
                  placeholder="请输入密钥"
                  type="password"
                />
              </el-form-item>
              <el-form-item label="返回URL">
                <el-input
                  v-model="config.payment.epayReturnURL"
                  :disabled="!config.payment.enableEpay"
                  placeholder="请输入返回URL"
                />
              </el-form-item>
              <el-form-item label="回调URL">
                <el-input
                  v-model="config.payment.epayNotifyURL"
                  :disabled="!config.payment.enableEpay"
                  placeholder="请输入回调URL"
                />
              </el-form-item>
            </el-card>

            <!-- 码支付配置 -->
            <el-card
              class="payment-card"
              shadow="never"
            >
              <template #header>
                <div class="payment-header">
                  <span>码支付配置</span>
                  <el-switch v-model="config.payment.enableMapay" />
                </div>
              </template>
              <el-form-item label="API地址">
                <el-input
                  v-model="config.payment.mapayAPIURL"
                  :disabled="!config.payment.enableMapay"
                  placeholder="请输入码支付API地址"
                />
              </el-form-item>
              <el-form-item label="商户ID">
                <el-input
                  v-model="config.payment.mapayID"
                  :disabled="!config.payment.enableMapay"
                  placeholder="请输入商户ID"
                />
              </el-form-item>
              <el-form-item label="密钥">
                <el-input
                  v-model="config.payment.mapayKey"
                  :disabled="!config.payment.enableMapay"
                  placeholder="请输入密钥"
                  type="password"
                />
              </el-form-item>
              <el-form-item label="返回URL">
                <el-input
                  v-model="config.payment.mapayReturnURL"
                  :disabled="!config.payment.enableMapay"
                  placeholder="请输入返回URL"
                />
              </el-form-item>
              <el-form-item label="回调URL">
                <el-input
                  v-model="config.payment.mapayNotifyURL"
                  :disabled="!config.payment.enableMapay"
                  placeholder="请输入回调URL"
                />
              </el-form-item>
            </el-card>

            <!-- 实名认证配置 -->
            <el-card
              class="payment-card"
              shadow="never"
            >
              <template #header>
                <div class="payment-header">
                  <span>实名认证配置</span>
                  <el-switch v-model="config.payment.enableRealName" />
                </div>
              </template>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="是否强制实名认证">
                    <el-switch v-model="config.payment.requireRealName" :disabled="!config.payment.enableRealName" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="回调URL">
                    <el-input
                      v-model="config.payment.realNameCallbackURL"
                      placeholder="请输入回调URL"
                      :disabled="!config.payment.enableRealName"
                    />
                    <div class="form-item-hint">
                      例如: https://your-domain.com/api/v1/kyc/callback
                    </div>
                  </el-form-item>
                </el-col>
              </el-row>
            </el-card>
          </el-form>
        </el-tab-pane>

        <!-- 用户等级配置 -->
        <el-tab-pane
          :label="$t('admin.config.userLevel')"
          name="quota"
        >
          <el-form
            v-loading="loading"
            :model="config"
            label-width="140px"
            class="config-form"
          >
            <el-alert
              :title="$t('admin.config.userLevelDesc')"
              type="info"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            >
              <div>{{ $t('admin.config.userLevelHint') }}</div>
              <div style="margin-top: 8px; color: #67C23A;">
                <i class="el-icon-check" />
                {{ $t('admin.config.autoSyncHint') }}
              </div>
              <div style="margin-top: 8px; color: #E6A23C;">
                <i class="el-icon-warning" />
                {{ $t('admin.config.resourceLimitWarning') }}
              </div>
            </el-alert>
            
            <el-form-item :label="$t('admin.config.newUserDefaultLevel')">
              <el-select
                v-model="config.quota.defaultLevel"
                :placeholder="$t('admin.config.selectDefaultLevel')"
                style="width: 200px"
              >
                <el-option
                  v-for="level in 5"
                  :key="level"
                  :label="$t('admin.config.levelN', { level })"
                  :value="level"
                />
              </el-select>
            </el-form-item>

            <el-divider content-position="left">
              {{ $t('admin.config.levelLimitsConfig') }}
            </el-divider>
            
            <!-- 等级限制配置 -->
            <el-row :gutter="15">
              <el-col
                v-for="level in 5"
                :key="level"
                :span="24"
                style="margin-bottom: 15px;"
              >
                <el-card 
                  class="level-card"
                  :class="{ 'default-level': config.quota.defaultLevel === level }"
                  shadow="hover"
                >
                  <template #header>
                    <div class="level-header">
                      <span class="level-title">{{ $t('admin.config.levelNLimits', { level }) }}</span>
                      <el-tag
                        v-if="config.quota.defaultLevel === level"
                        type="success"
                        size="small"
                      >
                        {{ $t('admin.config.defaultLevel') }}
                      </el-tag>
                    </div>
                  </template>
                  <el-row :gutter="20">
                    <el-col :span="6">
                      <el-form-item :label="$t('admin.config.maxInstances')">
                        <el-input-number 
                          v-model="config.quota.levelLimits[level]['maxInstances']" 
                          :min="1" 
                          :max="1000"
                          :controls="false"
                          :step="1"
                          style="width: 100%" 
                        />
                      </el-form-item>
                    </el-col>
                    <el-col :span="6">
                      <el-form-item :label="$t('admin.config.maxCPU')">
                        <el-input-number 
                          v-model="config.quota.levelLimits[level]['maxResources']['cpu']" 
                          :min="1" 
                          :max="10240"
                          :controls="false"
                          :step="1"
                          style="width: 100%" 
                        />
                      </el-form-item>
                    </el-col>
                    <el-col :span="6">
                      <el-form-item :label="$t('admin.config.maxMemoryMB')">
                        <el-input-number 
                          v-model="config.quota.levelLimits[level]['maxResources']['memory']" 
                          :min="128" 
                          :max="10485760"
                          :controls="false"
                          :step="128"
                          style="width: 100%" 
                        />
                      </el-form-item>
                    </el-col>
                    <el-col :span="6">
                      <el-form-item :label="$t('admin.config.maxDiskMB')">
                        <el-input-number 
                          v-model="config.quota.levelLimits[level]['maxResources']['disk']" 
                          :min="512" 
                          :max="1024000000"
                          :controls="false"
                          :step="512"
                          style="width: 100%" 
                        />
                      </el-form-item>
                    </el-col>
                  </el-row>
                  <el-row :gutter="20">
                    <el-col :span="6">
                      <el-form-item :label="$t('admin.config.maxBandwidthMbps')">
                        <el-input-number 
                          v-model="config.quota.levelLimits[level]['maxResources']['bandwidth']" 
                          :min="1" 
                          :max="1000000"
                          :controls="false"
                          :step="1"
                          style="width: 100%" 
                        />
                      </el-form-item>
                    </el-col>
                    <el-col :span="6">
                      <el-form-item :label="$t('admin.config.trafficLimitMB')">
                        <el-input-number 
                          v-model="config.quota.levelLimits[level]['maxTraffic']" 
                          :min="1024" 
                          :max="1024000000"
                          :controls="false"
                          :step="1024"
                          style="width: 100%" 
                        />
                      </el-form-item>
                    </el-col>
                  </el-row>
                </el-card>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <!-- 实例类型权限配置 -->
        <el-tab-pane
          :label="$t('admin.config.instancePermissions')"
          name="instancePermissions"
        >
          <el-form
            v-loading="loading"
            :model="instanceTypePermissions"
            label-width="180px"
            class="config-form"
          >
            <el-alert
              :title="$t('admin.config.instancePermissionsDesc')"
              type="info"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            >
              {{ $t('admin.config.instancePermissionsHint') }}
            </el-alert>
            
            <!-- 创建权限 -->
            <el-divider content-position="left">
              <el-icon><Plus /></el-icon> {{ $t('admin.config.createPermissions') }}
            </el-divider>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.containerCreateMinLevel')">
                  <el-select
                    v-model="instanceTypePermissions.minLevelForContainer"
                    :placeholder="$t('admin.config.selectLevel')"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="level in [1, 2, 3, 4, 5]"
                      :key="level"
                      :label="$t('admin.config.levelN', { level })"
                      :value="level"
                    />
                  </el-select>
                  <div class="form-item-hint">
                    {{ $t('admin.config.containerCreateHint') }}
                  </div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.vmCreateMinLevel')">
                  <el-select
                    v-model="instanceTypePermissions.minLevelForVM"
                    :placeholder="$t('admin.config.selectLevel')"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="level in [1, 2, 3, 4, 5]"
                      :key="level"
                      :label="$t('admin.config.levelN', { level })"
                      :value="level"
                    />
                  </el-select>
                  <div class="form-item-hint">
                    {{ $t('admin.config.vmCreateHint') }}
                  </div>
                </el-form-item>
              </el-col>
            </el-row>

            <!-- 删除权限 -->
            <el-divider content-position="left">
              <el-icon><Delete /></el-icon> {{ $t('admin.config.deletePermissions') }}
            </el-divider>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.containerDeleteMinLevel')">
                  <el-select
                    v-model="instanceTypePermissions.minLevelForDeleteContainer"
                    :placeholder="$t('admin.config.selectLevel')"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="level in [1, 2, 3, 4, 5]"
                      :key="level"
                      :label="$t('admin.config.levelN', { level })"
                      :value="level"
                    />
                  </el-select>
                  <div class="form-item-hint">
                    {{ $t('admin.config.containerDeleteHint') }}
                  </div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.vmDeleteMinLevel')">
                  <el-select
                    v-model="instanceTypePermissions.minLevelForDeleteVM"
                    :placeholder="$t('admin.config.selectLevel')"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="level in [1, 2, 3, 4, 5]"
                      :key="level"
                      :label="$t('admin.config.levelN', { level })"
                      :value="level"
                    />
                  </el-select>
                  <div class="form-item-hint">
                    {{ $t('admin.config.vmDeleteHint') }}
                  </div>
                </el-form-item>
              </el-col>
            </el-row>

            <!-- 重置系统权限 -->
            <el-divider content-position="left">
              <el-icon><Refresh /></el-icon> {{ $t('admin.config.resetPermissions') }}
            </el-divider>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.containerResetMinLevel')">
                  <el-select
                    v-model="instanceTypePermissions.minLevelForResetContainer"
                    :placeholder="$t('admin.config.selectLevel')"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="level in [1, 2, 3, 4, 5]"
                      :key="level"
                      :label="$t('admin.config.levelN', { level })"
                      :value="level"
                    />
                  </el-select>
                  <div class="form-item-hint">
                    {{ $t('admin.config.containerResetHint') }}
                  </div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.vmResetMinLevel')">
                  <el-select
                    v-model="instanceTypePermissions.minLevelForResetVM"
                    :placeholder="$t('admin.config.selectLevel')"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="level in [1, 2, 3, 4, 5]"
                      :key="level"
                      :label="$t('admin.config.levelN', { level })"
                      :value="level"
                    />
                  </el-select>
                  <div class="form-item-hint">
                    {{ $t('admin.config.vmResetHint') }}
                  </div>
                </el-form-item>
              </el-col>
            </el-row>

            <el-alert
              :title="$t('admin.config.permissionsSuggestions')"
              type="warning"
              :closable="false"
              show-icon
              style="margin-top: 20px;"
            >
              <ul style="margin: 0; padding-left: 20px;">
                <li>{{ $t('admin.config.containerCreateSuggestion') }}</li>
                <li>{{ $t('admin.config.vmCreateSuggestion') }}</li>
                <li>{{ $t('admin.config.containerDeleteResetSuggestion') }}</li>
                <li>{{ $t('admin.config.vmDeleteResetSuggestion') }}</li>
              </ul>
            </el-alert>
          </el-form>
        </el-tab-pane>

        <!-- 其他配置 -->
        <el-tab-pane
          :label="$t('admin.config.otherConfig')"
          name="other"
        >
          <el-form
            v-loading="loading"
            :model="config"
            label-width="140px"
            class="config-form"
          >
            <el-alert
              :title="$t('admin.config.avatarUploadConfig')"
              type="info"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            >
              {{ $t('admin.config.avatarUploadDesc') }}
            </el-alert>

            <el-divider content-position="left">
              {{ $t('admin.config.avatarUploadSettings') }}
            </el-divider>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.maxAvatarSize')">
                  <el-input-number
                    v-model="config.other.maxAvatarSize"
                    :min="0.5"
                    :max="10"
                    :step="0.5"
                    :precision="1"
                    :controls="false"
                    style="width: 100%"
                  />
                  <div class="form-item-hint">
                    {{ $t('admin.config.maxAvatarSizeHint') }}
                  </div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('admin.config.supportedFormats')">
                  <el-tag
                    type="info"
                    style="margin-right: 8px;"
                  >
                    PNG
                  </el-tag>
                  <el-tag type="info">
                    JPEG
                  </el-tag>
                  <div class="form-item-hint">
                    {{ $t('admin.config.supportedFormatsHint') }}
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
                <strong>{{ $t('admin.config.languageForceNote') || '强制语言设置说明' }}</strong>
              </template>
              {{ $t('admin.config.languageForceDesc') || '当设置了系统默认语言（选择中文或English）后，所有用户将被强制使用该语言，用户的手动语言切换将被覆盖。留空时将根据用户浏览器语言自动选择。' }}
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
        </el-tab-pane>

        <!-- 系统配置 -->
        <el-tab-pane
          label="系统配置"
          name="system"
        >
          <el-form
            v-loading="loading"
            :model="config"
            label-width="140px"
            class="config-form"
          >
            <el-alert
              title="系统配置说明"
              type="info"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            >
              配置系统基本参数，包括服务器端口、环境设置等
            </el-alert>

            <el-divider content-position="left">
              基本设置
            </el-divider>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="服务器端口">
                  <el-input-number
                    v-model="config.system.addr"
                    :min="1"
                    :max="65535"
                    :controls="false"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="运行环境">
                  <el-select
                    v-model="config.system.env"
                    style="width: 100%"
                  >
                    <el-option label="开发环境" value="development" />
                    <el-option label="生产环境" value="production" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="前端URL">
                  <el-input v-model="config.system.frontendURL" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="数据库类型">
                  <el-select
                    v-model="config.system.dbType"
                    style="width: 100%"
                  >
                    <el-option label="MySQL" value="mysql" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-divider content-position="left">
              安全设置
            </el-divider>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="IP限制计数">
                  <el-input-number
                    v-model="config.system.ipLimitCount"
                    :min="0"
                    :controls="false"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="IP限制时间（秒）">
                  <el-input-number
                    v-model="config.system.ipLimitTime"
                    :min="0"
                    :controls="false"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="使用Redis">
                  <el-switch v-model="config.system.useRedis" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="多点登录">
                  <el-switch v-model="config.system.useMultipoint" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <!-- 验证码配置 -->
        <el-tab-pane
          label="验证码配置"
          name="captcha"
        >
          <el-form
            v-loading="loading"
            :model="config"
            label-width="140px"
            class="config-form"
          >
            <el-alert
              title="验证码配置说明"
              type="info"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            >
              配置图形验证码参数，用于登录和注册验证
            </el-alert>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="启用验证码">
                  <el-switch v-model="config.captcha.enabled" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="验证码长度">
                  <el-input-number
                    v-model="config.captcha.length"
                    :min="4"
                    :max="10"
                    :controls="false"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="验证码过期时间（秒）">
                  <el-input-number
                    v-model="config.captcha.expireTime"
                    :min="60"
                    :max="3600"
                    :controls="false"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="验证码宽度">
                  <el-input-number
                    v-model="config.captcha.width"
                    :min="100"
                    :max="300"
                    :controls="false"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="验证码高度">
                  <el-input-number
                    v-model="config.captcha.height"
                    :min="30"
                    :max="100"
                    :controls="false"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <!-- 任务配置 -->
        <el-tab-pane
          label="任务配置"
          name="task"
        >
          <el-form
            v-loading="loading"
            :model="config"
            label-width="140px"
            class="config-form"
          >
            <el-alert
              title="任务配置说明"
              type="info"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            >
              配置任务执行参数，包括删除重试次数和延迟时间
            </el-alert>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="删除重试次数">
                  <el-input-number
                    v-model="config.task.deleteRetryCount"
                    :min="1"
                    :max="10"
                    :controls="false"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="删除重试延迟（秒）">
                  <el-input-number
                    v-model="config.task.deleteRetryDelay"
                    :min="1"
                    :max="60"
                    :controls="false"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <!-- 站点配置 -->
        <el-tab-pane
          label="站点配置"
          name="site"
        >
          <el-form
            v-loading="loading"
            :model="config"
            label-width="140px"
            class="config-form"
          >
            <el-alert
              title="站点配置说明"
              type="info"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            >
              配置站点基本信息，包括站点名称、描述等
            </el-alert>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="站点名称">
                  <el-input v-model="config.site.name" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="站点描述">
                  <el-input v-model="config.site.description" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="站点关键词">
                  <el-input v-model="config.site.keywords" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="系统名称">
                  <el-input v-model="config.systemName" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="系统描述">
                  <el-input v-model="config.systemDescription" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <!-- 底部操作按钮 -->
      <div class="config-actions">
        <el-button
          type="primary"
          size="large"
          :loading="loading"
          @click="saveConfig"
        >
          {{ $t('admin.config.saveCurrentConfig') }}
        </el-button>
        <el-button 
          size="large"
          @click="resetConfig"
        >
          {{ $t('admin.config.resetConfig') }}
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getAdminConfig, updateAdminConfig } from '@/api/config'
import { getInstanceTypePermissions, updateInstanceTypePermissions } from '@/api/admin'
import { useLanguageStore } from '@/pinia/modules/language'

const { t, locale } = useI18n()
const languageStore = useLanguageStore()

// 当前激活的标签页
const activeTab = ref('auth')

const config = ref({
  auth: {
    enableEmail: false,
    enableTelegram: false,
    enableQQ: false,
    enableOAuth2: false,
    enablePublicRegistration: false, // 是否启用公开注册
    emailSMTPHost: '',
    emailSMTPPort: 587,
    emailUsername: '',
    emailPassword: '',
    telegramBotToken: '',
    qqAppID: '',
    qqAppKey: ''
  },
  quota: {
    defaultLevel: 1,
    levelLimits: {
      1: { maxInstances: 1, maxResources: { cpu: 1, memory: 512, disk: 1024, bandwidth: 100 }, maxTraffic: 102400 },    // 磁盘1GB, 流量100MB
      2: { maxInstances: 3, maxResources: { cpu: 2, memory: 1024, disk: 2048, bandwidth: 200 }, maxTraffic: 204800 },   // 磁盘2GB, 流量200MB
      3: { maxInstances: 5, maxResources: { cpu: 4, memory: 2048, disk: 4096, bandwidth: 500 }, maxTraffic: 409600 },   // 磁盘4GB, 流量400MB
      4: { maxInstances: 10, maxResources: { cpu: 8, memory: 4096, disk: 8192, bandwidth: 1000 }, maxTraffic: 819200 },  // 磁盘8GB, 流量800MB
      5: { maxInstances: 20, maxResources: { cpu: 16, memory: 8192, disk: 16384, bandwidth: 2000 }, maxTraffic: 1638400 } // 磁盘16GB, 流量1600MB
    }
  },
  inviteCode: {
    enabled: false,
    required: false
  },
  payment: {
    enableAlipay: false,
    alipayAppId: '',
    alipayPrivateKey: '',
    alipayPublicKey: '',
    alipayGateway: 'https://openapi.alipay.com/gateway.do',
    enableWechat: false,
    wechatMchId: '',
    wechatAppId: '',
    wechatApiKey: '',
    wechatSerialNo: '',
    wechatApiV3Key: '',
    wechatType: 'mp',
    enableBalance: true,
    balanceMinAmount: 0,
    balanceDailyLimit: 10000,
    // 易支付配置
    enableEpay: false,
    epayAPIURL: '',
    epayPID: '',
    epayKey: '',
    epayReturnURL: '',
    epayNotifyURL: '',
    // 码支付配置
    enableMapay: false,
    mapayAPIURL: '',
    mapayID: '',
    mapayKey: '',
    mapayReturnURL: '',
    mapayNotifyURL: '',
    // 实名认证配置
    enableRealName: false,
    requireRealName: false,
    realNameCallbackURL: ''
  },
  other: {
    maxAvatarSize: 2, // MB
    defaultLanguage: '' // 默认语言，空字符串表示使用浏览器语言
  },
  system: {
    addr: 8890,
    dbType: 'mysql',
    env: 'development',
    frontendURL: 'https://heyun.ypvps.com',
    ipLimitCount: 15000,
    ipLimitTime: 3600,
    oauth2StateTokenMinutes: 15,
    ossType: 'local',
    providerInactiveHours: 24,
    useMultipoint: false,
    useRedis: false
  },
  captcha: {
    enabled: false,
    expireTime: 300,
    height: 40,
    length: 4,
    width: 120
  },
  task: {
    deleteRetryCount: 3,
    deleteRetryDelay: 2
  },
  site: {
    name: 'OneClickVirt',
    description: '虚拟化管理平台',
    keywords: '虚拟化,Docker,LXD,Incus,Proxmox'
  },
  systemName: '虚拟化管理平台',
  systemDescription: '支持多种虚拟化技术的管理平台'
})

const instanceTypePermissions = ref({
  minLevelForContainer: 1,
  minLevelForVM: 3,
  minLevelForDeleteContainer: 1,
  minLevelForDeleteVM: 2,
  minLevelForResetContainer: 1,
  minLevelForResetVM: 2
})

const loading = ref(false)

// 记录系统配置的语言，用于判断是否修改
const systemConfigLanguage = ref('')

const loadConfig = async () => {
  loading.value = true
  try {
    const response = await getAdminConfig()
    console.log('加载配置响应:', response)
    console.log('配置数据:', response.data)
    if (response.code === 0 && response.data) {
      // 合并配置，确保所有字段都有默认值
      if (response.data.auth) {
        config.value.auth = {
          ...config.value.auth,
          ...response.data.auth
        }
      }
      
      if (response.data.inviteCode) {
        config.value.inviteCode = {
          ...config.value.inviteCode,
          ...response.data.inviteCode
        }
      }

      // 加载支付配置
      if (response.data.payment) {
        console.log('加载支付配置:', response.data.payment)
        config.value.payment = {
          ...config.value.payment,
          ...response.data.payment
        }
        console.log('合并后的支付配置:', config.value.payment)
      }

      // 加载其他配置
      if (response.data.other) {
        console.log('加载其他配置:', response.data.other)
        config.value.other = {
          ...config.value.other,
          ...response.data.other
        }
        // 记录当前的系统语言配置
        systemConfigLanguage.value = config.value.other.defaultLanguage || ''
        console.log('合并后的其他配置:', config.value.other)
        console.log('当前系统语言配置:', systemConfigLanguage.value)
      }
      
      // 加载系统配置
      if (response.data.system) {
        console.log('加载系统配置:', response.data.system)
        config.value.system = {
          ...config.value.system,
          ...response.data.system
        }
      }

      // 加载验证码配置
      if (response.data.captcha) {
        console.log('加载验证码配置:', response.data.captcha)
        config.value.captcha = {
          ...config.value.captcha,
          ...response.data.captcha
        }
      }

      // 加载任务配置
      if (response.data.task) {
        console.log('加载任务配置:', response.data.task)
        config.value.task = {
          ...config.value.task,
          ...response.data.task
        }
      }

      // 加载站点配置
      if (response.data.site) {
        console.log('加载站点配置:', response.data.site)
        config.value.site = {
          ...config.value.site,
          ...response.data.site
        }
      }

      // 加载系统名称和描述
      if (response.data.systemName) {
        config.value.systemName = response.data.systemName
      }
      if (response.data.systemDescription) {
        config.value.systemDescription = response.data.systemDescription
      }
      
      // 加载等级配置
      if (response.data.quota && response.data.quota.levelLimits) {
        config.value.quota.levelLimits = {}
        for (let level = 1; level <= 5; level++) {
          const levelKey = String(level)
          if (response.data.quota.levelLimits[levelKey]) {
            const limitData = response.data.quota.levelLimits[levelKey]
            config.value.quota.levelLimits[level] = {
              maxInstances: limitData['max-instances'] || (level * 2),
              maxResources: {
                cpu: limitData['max-resources']?.cpu || (level * 2),
                memory: limitData['max-resources']?.memory || (1024 * Math.pow(2, level - 1)),
                disk: limitData['max-resources']?.disk || (10240 * Math.pow(2, level - 1)),
                bandwidth: limitData['max-resources']?.bandwidth || (10 * level)
              },
              maxTraffic: limitData['max-traffic'] || (1024 * level)
            }
          } else {
            // 如果没有数据，使用默认值
            config.value.quota.levelLimits[level] = {
              maxInstances: level * 2,
              maxResources: {
                cpu: level * 2,
                memory: 1024 * Math.pow(2, level - 1),
                disk: 10240 * Math.pow(2, level - 1),
                bandwidth: 10 * level
              },
              maxTraffic: 1024 * level
            }
          }
        }
      }
    }
  } catch (error) {
    console.error('加载配置失败:', error)
    ElMessage.error(t('admin.config.loadConfigFailed'))
  } finally {
    loading.value = false
  }
}

const loadInstanceTypePermissions = async () => {
  try {
    const response = await getInstanceTypePermissions()
    console.log('加载实例类型权限配置响应:', response)
    if (response.code === 0 && response.data) {
      instanceTypePermissions.value = {
        minLevelForContainer: response.data.minLevelForContainer || 1,
        minLevelForVM: response.data.minLevelForVM || 3,
        minLevelForDeleteContainer: response.data.minLevelForDeleteContainer || 1,
        minLevelForDeleteVM: response.data.minLevelForDeleteVM || 2,
        minLevelForResetContainer: response.data.minLevelForResetContainer || 1,
        minLevelForResetVM: response.data.minLevelForResetVM || 2
      }
    }
  } catch (error) {
    console.error('加载实例类型权限配置失败:', error)
    ElMessage.error(t('admin.config.loadPermissionsFailed'))
  }
}

const saveConfig = async () => {
  // 验证配置数据，确保所有资源限制值不为空
  for (let level = 1; level <= 5; level++) {
    const limit = config.value.quota.levelLimits[level]
    if (!limit) {
      ElMessage.error(t('admin.config.levelConfigEmpty', { level }))
      return
    }
    
    // 验证必填字段
    if (!limit.maxInstances || limit.maxInstances <= 0) {
      ElMessage.error(t('admin.config.maxInstancesInvalid', { level }))
      return
    }
    
    if (!limit.maxTraffic || limit.maxTraffic <= 0) {
      ElMessage.error(t('admin.config.trafficLimitInvalid', { level }))
      return
    }
    
    if (!limit.maxResources) {
      ElMessage.error(t('admin.config.resourceConfigEmpty', { level }))
      return
    }
    
    // 验证各项资源限制
    if (!limit.maxResources.cpu || limit.maxResources.cpu <= 0) {
      ElMessage.error(t('admin.config.maxCPUInvalid', { level }))
      return
    }
    
    if (!limit.maxResources.memory || limit.maxResources.memory <= 0) {
      ElMessage.error(t('admin.config.maxMemoryInvalid', { level }))
      return
    }
    
    if (!limit.maxResources.disk || limit.maxResources.disk <= 0) {
      ElMessage.error(t('admin.config.maxDiskInvalid', { level }))
      return
    }
    
    if (!limit.maxResources.bandwidth || limit.maxResources.bandwidth <= 0) {
      ElMessage.error(t('admin.config.maxBandwidthInvalid', { level }))
      return
    }
  }
  
  loading.value = true
  try {
    console.log('开始保存配置...')
    console.log('基础配置:', config.value)
    console.log('实例类型权限配置:', instanceTypePermissions.value)
    console.log('语言配置:', config.value.other.defaultLanguage)
    
    // 记录修改前的语言设置
    const oldLanguage = systemConfigLanguage.value
    const newLanguage = config.value.other.defaultLanguage
    const languageChanged = oldLanguage !== newLanguage
    
    // 转换 levelLimits 为 kebab-case 格式（外层字段），max-resources 内部保持 camelCase
    const configToSave = JSON.parse(JSON.stringify(config.value))
    if (configToSave.quota && configToSave.quota.levelLimits) {
      const convertedLimits = {}
      Object.keys(configToSave.quota.levelLimits).forEach(level => {
        const limit = configToSave.quota.levelLimits[level]
        convertedLimits[level] = {
          'max-instances': limit.maxInstances,
          'max-resources': {
            cpu: limit.maxResources.cpu,
            memory: limit.maxResources.memory,
            disk: limit.maxResources.disk,
            bandwidth: limit.maxResources.bandwidth
          },
          'max-traffic': limit.maxTraffic
        }
      })
      configToSave.quota.levelLimits = convertedLimits
    }
    
    // 转换支付配置为 kebab-case 格式
    if (configToSave.payment) {
      const payment = configToSave.payment
      configToSave.payment = {
        // 支付宝
        'enable-alipay': payment.enableAlipay,
        'alipay-app-id': payment.alipayAppId,
        'alipay-private-key': payment.alipayPrivateKey,
        'alipay-public-key': payment.alipayPublicKey,
        'alipay-gateway': payment.alipayGateway,
        // 微信支付
        'enable-wechat': payment.enableWechat,
        'wechat-app-id': payment.wechatAppId,
        'wechat-mch-id': payment.wechatMchId,
        'wechat-api-key': payment.wechatApiKey,
        'wechat-api-v3-key': payment.wechatApiV3Key,
        'wechat-serial-no': payment.wechatSerialNo,
        'wechat-type': payment.wechatType,
        // 余额支付
        'enable-balance': payment.enableBalance,
        'balance-daily-limit': payment.balanceDailyLimit,
        'balance-min-amount': payment.balanceMinAmount,
        // 易支付
        'enable-epay': payment.enableEpay,
        'epay-api-url': payment.epayAPIURL,
        'epay-pid': payment.epayPID,
        'epay-key': payment.epayKey,
        'epay-return-url': payment.epayReturnURL,
        'epay-notify-url': payment.epayNotifyURL,
        // 码支付
        'enable-mapay': payment.enableMapay,
        'mapay-api-url': payment.mapayAPIURL,
        'mapay-id': payment.mapayID,
        'mapay-key': payment.mapayKey,
        'mapay-return-url': payment.mapayReturnURL,
        'mapay-notify-url': payment.mapayNotifyURL,
        // 实名认证配置
        'enable-real-name': payment.enableRealName,
        'require-real-name': payment.requireRealName,
        'real-name-callback-url': payment.realNameCallbackURL
      }
    }
    
    // 转换系统配置为 kebab-case 格式
    if (configToSave.system) {
      const system = configToSave.system
      configToSave.system = {
        'addr': system.addr,
        'db-type': system.dbType,
        'env': system.env,
        'frontend-url': system.frontendURL,
        'iplimit-count': system.ipLimitCount,
        'iplimit-time': system.ipLimitTime,
        'oauth2-state-token-minutes': system.oauth2StateTokenMinutes,
        'oss-type': system.ossType,
        'provider-inactive-hours': system.providerInactiveHours,
        'use-multipoint': system.useMultipoint,
        'use-redis': system.useRedis
      }
    }
    
    // 转换验证码配置为 kebab-case 格式
    if (configToSave.captcha) {
      const captcha = configToSave.captcha
      configToSave.captcha = {
        'enabled': captcha.enabled,
        'expire-time': captcha.expireTime,
        'height': captcha.height,
        'length': captcha.length,
        'width': captcha.width
      }
    }
    
    // 转换任务配置为 kebab-case 格式
    if (configToSave.task) {
      const task = configToSave.task
      configToSave.task = {
        'delete-retry-count': task.deleteRetryCount,
        'delete-retry-delay': task.deleteRetryDelay
      }
    }
    
    // 转换站点配置为 kebab-case 格式
    if (configToSave.site) {
      const site = configToSave.site
      configToSave.site = {
        'name': site.name,
        'description': site.description,
        'keywords': site.keywords
      }
    }
    
    // 保存基础配置
    const configResult = await updateAdminConfig(configToSave)
    console.log('基础配置保存结果:', configResult)
    
    // 保存实例类型权限配置
    const permissionsResult = await updateInstanceTypePermissions(instanceTypePermissions.value)
    console.log('实例类型权限配置保存结果:', permissionsResult)
    
    ElMessage.success(t('admin.config.saveSuccess'))
    
    // 如果修改了默认语言，强制应用并刷新页面
    if (languageChanged) {
      console.log('[Config] 系统语言已修改，从', oldLanguage, '到', newLanguage)
      
      // 更新 language store 中的系统配置语言并强制应用
      const effectiveLanguage = languageStore.forceApplySystemLanguage(newLanguage)
      console.log('[Config] 强制应用后的有效语言:', effectiveLanguage)
      
      // 更新当前页面的语言
      locale.value = effectiveLanguage
      
      // 显示通知，告知用户页面将刷新
      ElNotification({
        title: t('common.success'),
        message: t('admin.config.languageChangedRefreshing'),
        type: 'success',
        duration: 2000
      })
      
      // 延迟刷新页面，让用户看到通知
      setTimeout(() => {
        window.location.reload()
      }, 2000)
    } else {
      // 保存成功后重新加载配置，确保显示最新数据
      await loadConfig()
      await loadInstanceTypePermissions()
    }
  } catch (error) {
    console.error('保存配置失败:', error)
    ElMessage.error(t('admin.config.saveFailed', { error: error.message || t('common.unknownError') }))
  } finally {
    loading.value = false
  }
}

const resetConfig = async () => {
  await loadConfig()
  await loadInstanceTypePermissions()
  ElMessage.success(t('admin.config.configReset'))
}

onMounted(() => {
  loadConfig()
  loadInstanceTypePermissions()
})
</script>

<style scoped>
.config-header {
  display: flex;
  flex-direction: column;
  gap: 4px;

  > span {
    font-size: 18px;
    font-weight: 600;
    color: #303133;
  }
}

.config-tabs {
  margin-bottom: 20px;
}

.config-tabs :deep(.el-tabs__content) {
  padding: 20px;
}

.payment-card {
  margin-bottom: 20px;
}

.payment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  color: #303133;
}

.config-form {
  max-height: 600px;
  overflow-y: auto;
}

.oauth-card {
  margin-bottom: 16px;
}

.oauth-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.level-card {
  border: 2px solid #f0f0f0;
  transition: all 0.3s ease;
}

.level-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px 0 rgba(64, 158, 255, 0.1);
}

.level-card.default-level {
  border-color: #67c23a;
  background-color: #f0f9ff;
}

.level-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.level-title {
  font-weight: 600;
  color: #303133;
}

.config-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 20px 0;
  border-top: 1px solid #f0f0f0;
  margin-top: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .config-container {
    padding: 10px;
  }
  
  .config-form {
    max-height: none;
  }
  
  .level-card :deep(.el-col) {
    margin-bottom: 10px;
  }
  
  .config-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .config-actions .el-button {
    width: 100%;
    max-width: 200px;
  }
}

/* 标签页样式 */
.config-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
}

.config-tabs :deep(.el-tabs__nav-wrap) {
  padding: 0 10px;
}

.config-tabs :deep(.el-tabs__item) {
  padding: 0 20px;
  font-weight: 500;
}

/* 表单样式 */
.config-form :deep(.el-form-item__label) {
  font-weight: 500;
  color: #606266;
}

.config-form :deep(.el-alert) {
  margin-bottom: 20px;
}

.form-item-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.4;
}
</style>