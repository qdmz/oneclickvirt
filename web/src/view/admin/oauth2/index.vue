<template>
  <div class="oauth2-providers-container">
    <!-- OAuth2 功能未启用提示 -->
    <el-alert
      v-if="!oauth2Enabled"
      :title="$t('admin.oauth2.notEnabled')"
      type="warning"
      :closable="false"
      show-icon
      style="margin-bottom: 20px;"
    >
      <template #default>
        <div>
          {{ $t('admin.oauth2.notEnabledHint') }}
          <br>
          {{ $t('admin.oauth2.enableHint') }}
          <el-link
            type="primary"
            :underline="false"
            @click="goToConfig"
          >
            <strong>{{ $t('admin.oauth2.systemConfig') }}</strong>
          </el-link>
          {{ $t('admin.oauth2.enableHint2') }}
        </div>
      </template>
    </el-alert>

    <el-card
      shadow="never"
      class="providers-card"
    >
      <template #header>
        <div class="card-header">
          <span>{{ $t('admin.oauth2.title') }}</span>
          <el-button
            type="primary"
            @click="handleAdd"
          >
            {{ $t('admin.oauth2.addProvider') }}
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="providers"
        class="providers-table"
        :cell-style="{ padding: '12px 0' }"
        :header-cell-style="{ background: '#f5f7fa', padding: '14px 0', fontWeight: '600' }"
      >
        <el-table-column
          prop="id"
          label="ID"
          width="80"
          align="center"
        />
        <el-table-column
          prop="displayName"
          :label="$t('admin.oauth2.displayName')"
          min-width="140"
        />
        <el-table-column
          prop="name"
          :label="$t('admin.oauth2.identifierName')"
          min-width="160"
        />
        <el-table-column
          :label="$t('common.status')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <el-tag
              :type="row.enabled ? 'success' : 'info'"
              size="default"
            >
              {{ row.enabled ? $t('common.enabled') : $t('common.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('admin.oauth2.registrationStats')"
          min-width="190"
          align="center"
        >
          <template #default="{ row }">
            <span v-if="row.maxRegistrations > 0">
              {{ row.currentRegistrations }} / {{ row.maxRegistrations }}
            </span>
            <span v-else>
              {{ row.totalUsers }} ({{ $t('admin.oauth2.unlimited') }})
            </span>
          </template>
        </el-table-column>
        <el-table-column
          prop="clientId"
          :label="$t('admin.oauth2.clientIdLabel')"
          min-width="220"
          show-overflow-tooltip
        />
        <el-table-column
          prop="redirectUrl"
          :label="$t('admin.oauth2.callbackUrl')"
          min-width="200"
          show-overflow-tooltip
        />
        <el-table-column
          :label="$t('common.actions')"
          width="340"
          fixed="right"
          align="center"
        >
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button
                size="small"
                @click="handleEdit(row)"
              >
                {{ $t('common.edit') }}
              </el-button>
              <el-button
                size="small"
                type="warning"
                @click="handleResetCount(row)"
              >
                {{ $t('admin.oauth2.resetCount') }}
              </el-button>
              <el-button
                size="small"
                type="danger"
                @click="handleDelete(row)"
              >
                {{ $t('common.delete') }}
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="900px"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="dialog-header">
          <span>{{ dialogTitle }}</span>
          <div
            v-if="!isEdit"
            class="preset-selector"
          >
            <el-select
              v-model="selectedPreset"
              :placeholder="$t('admin.oauth2.selectPreset')"
              size="small"
              style="width: 200px"
              clearable
              @change="handlePresetChange"
            >
              <el-option
                label="Linux.do"
                value="linuxdo"
              />
              <el-option
                label="IDCFlare"
                value="idcflare"
              />
              <el-option
                label="GitHub"
                value="github"
              />
              <el-option
                label="GitLab"
                value="gitlab"
              />
              <el-option
                label="Gitea"
                value="gitea"
              />
              <el-option
                label="Google"
                value="google"
              />
              <el-option
                label="Microsoft"
                value="microsoft"
              />
              <el-option
                label="Discord"
                value="discord"
              />
              <el-option
                :label="$t('admin.oauth2.genericOAuth2')"
                value="generic"
              />
              <el-option
                :label="$t('admin.oauth2.customOAuth2')"
                value="custom"
              />
            </el-select>
          </div>
        </div>
      </template>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
        class="oauth2-form"
      >
        <el-tabs
          v-model="activeTab"
          class="oauth2-tabs"
        >
          <el-tab-pane
            :label="$t('admin.oauth2.basicConfig')"
            name="basic"
          >
            <div class="form-section">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item
                    :label="$t('admin.oauth2.displayName')"
                    prop="displayName"
                  >
                    <el-input
                      v-model="formData.displayName"
                      :placeholder="$t('admin.oauth2.displayNamePlaceholder')"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item
                    :label="$t('admin.oauth2.identifierName')"
                    prop="name"
                  >
                    <el-input
                      v-model="formData.name"
                      :placeholder="$t('admin.oauth2.identifierNamePlaceholder')"
                      :disabled="isEdit"
                    />
                  </el-form-item>
                </el-col>
              </el-row>

              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item :label="$t('admin.oauth2.enableStatus')">
                    <el-switch
                      v-model="formData.enabled"
                      :active-text="$t('common.enable')"
                      :inactive-text="$t('common.disable')"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item
                    :label="$t('admin.oauth2.displayOrder')"
                    prop="sort"
                  >
                    <el-input-number
                      v-model="formData.sort"
                      :min="0"
                      :max="999"
                      :controls="false"
                      style="width: 100%"
                    />
                    <span class="form-tip">{{ $t('admin.oauth2.displayOrderHint') }}</span>
                  </el-form-item>
                </el-col>
              </el-row>

              <el-divider content-position="left">
                {{ $t('admin.oauth2.oauth2Credentials') }}
              </el-divider>

              <el-form-item
                :label="$t('admin.oauth2.clientIdLabel')"
                prop="clientId"
              >
                <el-input
                  v-model="formData.clientId"
                  :placeholder="$t('admin.oauth2.clientIdPlaceholder')"
                />
              </el-form-item>

              <el-form-item
                :label="$t('admin.oauth2.clientSecretLabel')"
                prop="clientSecret"
              >
                <el-input
                  v-model="formData.clientSecret"
                  type="password"
                  :placeholder="isEdit ? $t('admin.oauth2.secretPlaceholderEdit') : $t('admin.oauth2.clientSecretPlaceholder')"
                  show-password
                />
              </el-form-item>
            </div>
          </el-tab-pane>

          <el-tab-pane
            :label="$t('admin.oauth2.oauth2Endpoints')"
            name="endpoints"
          >
            <el-form-item
              :label="$t('admin.oauth2.callbackUrl')"
              prop="redirectUrl"
            >
              <el-input
                v-model="formData.redirectUrl"
                placeholder="http://localhost:8888/api/v1/auth/oauth2/callback"
              />
            </el-form-item>

            <el-form-item
              :label="$t('admin.oauth2.authUrl')"
              prop="authUrl"
            >
              <el-input
                v-model="formData.authUrl"
                placeholder="https://provider.com/oauth2/authorize"
              />
            </el-form-item>

            <el-form-item
              :label="$t('admin.oauth2.tokenUrl')"
              prop="tokenUrl"
            >
              <el-input
                v-model="formData.tokenUrl"
                placeholder="https://provider.com/oauth2/token"
              />
            </el-form-item>

            <el-form-item
              :label="$t('admin.oauth2.userInfoUrl')"
              prop="userInfoUrl"
            >
              <el-input
                v-model="formData.userInfoUrl"
                placeholder="https://provider.com/api/user"
              />
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane
            :label="$t('admin.oauth2.fieldMapping')"
            name="fields"
          >
            <el-alert
              type="info"
              :closable="false"
              style="margin-bottom: 20px"
            >
              <p>{{ $t('admin.oauth2.fieldMappingDesc') }}</p>
              <p>• {{ $t('admin.oauth2.requiredFields') }}</p>
              <p>• {{ $t('admin.oauth2.optionalFields') }}</p>
              <p>• {{ $t('admin.oauth2.nestedFieldsSupport') }}</p>
              <p>• {{ $t('admin.oauth2.defaultValuesInfo') }}</p>
            </el-alert>

            <el-form-item
              :label="$t('admin.oauth2.userIdField')"
              prop="userIdField"
            >
              <el-input
                v-model="formData.userIdField"
                :placeholder="$t('admin.oauth2.userIdFieldPlaceholder')"
              />
            </el-form-item>

            <el-form-item
              :label="$t('admin.oauth2.usernameField')"
              prop="usernameField"
            >
              <el-input
                v-model="formData.usernameField"
                :placeholder="$t('admin.oauth2.usernameFieldPlaceholder')"
              />
            </el-form-item>

            <el-form-item :label="$t('admin.oauth2.emailField')">
              <el-input
                v-model="formData.emailField"
                :placeholder="$t('admin.oauth2.emailFieldPlaceholder')"
              />
            </el-form-item>

            <el-form-item :label="$t('admin.oauth2.avatarField')">
              <el-input
                v-model="formData.avatarField"
                :placeholder="$t('admin.oauth2.avatarFieldPlaceholder')"
              />
            </el-form-item>

            <el-form-item :label="$t('admin.oauth2.nicknameField')">
              <el-input
                v-model="formData.nicknameField"
                :placeholder="$t('admin.oauth2.nicknameFieldPlaceholder')"
              />
            </el-form-item>

            <el-form-item :label="$t('admin.oauth2.trustLevelField')">
              <el-input
                v-model="formData.trustLevelField"
                :placeholder="$t('admin.oauth2.trustLevelFieldPlaceholder')"
              />
              <span class="form-tip">{{ $t('admin.oauth2.trustLevelFieldHint') }}</span>
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane
            :label="$t('admin.oauth2.levelAndLimits')"
            name="level"
          >
            <el-form-item
              :label="$t('admin.oauth2.defaultUserLevel')"
              prop="defaultLevel"
            >
              <el-input-number
                v-model="formData.defaultLevel"
                :min="1"
                :max="10"
                :controls="false"
              />
              <span class="form-tip">{{ $t('admin.oauth2.defaultUserLevelHint') }}</span>
            </el-form-item>

            <el-form-item :label="$t('admin.oauth2.levelMappingConfig')">
              <div class="level-mapping">
                <div
                  v-for="(level, key) in formData.levelMapping"
                  :key="key"
                  class="mapping-item"
                >
                  <span>{{ $t('admin.oauth2.externalLevel') }} {{ key }} →</span>
                  <el-input-number
                    v-model="formData.levelMapping[key]"
                    :min="1"
                    :max="10"
                    :controls="false"
                    size="small"
                  />
                  <el-button
                    size="small"
                    type="danger"
                    text
                    @click="removeLevelMapping(key)"
                  >
                    {{ $t('common.delete') }}
                  </el-button>
                </div>
                <el-button
                  size="small"
                  @click="addLevelMapping"
                >
                  <el-icon><Plus /></el-icon>
                  {{ $t('admin.oauth2.addMapping') }}
                </el-button>
              </div>
              <span class="form-tip">{{ $t('admin.oauth2.levelMappingHint') }}</span>
            </el-form-item>

            <el-form-item :label="$t('admin.oauth2.registrationLimit')">
              <el-input-number
                v-model="formData.maxRegistrations"
                :min="0"
                :max="999999"
                :controls="false"
              />
              <span class="form-tip">{{ $t('admin.oauth2.registrationLimitHint') }}</span>
            </el-form-item>

            <el-form-item
              v-if="isEdit"
              :label="$t('admin.oauth2.currentRegistrations')"
            >
              <el-input-number
                v-model="formData.currentRegistrations"
                :controls="false"
                disabled
              />
            </el-form-item>
          </el-tab-pane>
        </el-tabs>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="handleSubmit"
        >
          {{ $t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 添加等级映射对话框 -->
    <el-dialog
      v-model="mappingDialogVisible"
      :title="$t('admin.oauth2.addLevelMapping')"
      width="400px"
    >
      <el-form label-width="120px">
        <el-form-item :label="$t('admin.oauth2.externalLevelValue')">
          <el-input
            v-model="newMapping.externalLevel"
            :placeholder="$t('admin.oauth2.externalLevelPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="$t('admin.oauth2.systemUserLevel')">
          <el-input-number
            v-model="newMapping.systemLevel"
            :min="1"
            :max="10"
            :controls="false"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="mappingDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          @click="confirmAddMapping"
        >
          {{ $t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { Plus, Connection, Setting } from '@element-plus/icons-vue'
import useOAuth2 from './useOAuth2'

const {
  loading, providers,
  dialogVisible, dialogTitle, isEdit, submitting, activeTab, formRef, oauth2Enabled,
  mappingDialogVisible, selectedPreset, newMapping,
  formData, formRules,
  goToConfig, loadProviders,
  handlePresetChange, handleAdd, handleEdit, handleSubmit,
  handleDelete, handleResetCount,
  addLevelMapping, confirmAddMapping, removeLevelMapping
} = useOAuth2()
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

.providers-table {
  width: 100%;
  
  .action-buttons {
    display: flex;
    gap: 10px;
    justify-content: center;
    flex-wrap: wrap;
    padding: 4px 0;
    
    .el-button {
      margin: 0 !important;
    }
  }
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  
  .preset-selector {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.oauth2-form {
  .oauth2-tabs {
    :deep(.el-tabs__content) {
      padding-top: 20px;
    }
  }

  .form-section {
    padding: 10px 0;
  }

  :deep(.el-form-item) {
    margin-bottom: 24px;
  }

  :deep(.el-divider) {
    margin: 30px 0 24px 0;
  }

  :deep(.el-input-number) {
    width: 100%;
  }
  
  :deep(.el-col) {
    .el-form-item {
      margin-right: 0;
    }
  }
}

.form-tip {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
}

.level-mapping {
  .mapping-item {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 10px;

    span {
      min-width: 120px;
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
