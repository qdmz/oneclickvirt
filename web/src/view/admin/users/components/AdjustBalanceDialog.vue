<template>
  <el-dialog
    :model-value="visible"
    :title="$t('admin.users.adjustBalance')"
    width="520px"
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-form label-width="120px">
      <el-form-item :label="$t('admin.users.username')">
        <el-input
          :model-value="form.username"
          disabled
        />
      </el-form-item>
      <el-form-item :label="$t('admin.users.currentBalance')">
        <span class="current-balance">¥{{ Number(form.currentBalance || 0).toFixed(2) }}</span>
      </el-form-item>
      <el-form-item :label="$t('admin.users.adjustMode')">
        <el-radio-group v-model="form.mode">
          <el-radio-button value="add">
            {{ $t('admin.users.modeAdd') }}
          </el-radio-button>
          <el-radio-button value="set">
            {{ $t('admin.users.modeSet') }}
          </el-radio-button>
        </el-radio-group>
        <div class="form-tip">
          {{ form.mode === 'add' ? $t('admin.users.modeAddTip') : $t('admin.users.modeSetTip') }}
        </div>
      </el-form-item>
      <el-form-item :label="$t('admin.users.adjustAmount')">
        <el-input-number
          v-model="form.amount"
          :precision="2"
          :step="10"
          :min="form.mode === 'set' ? 0 : undefined"
          style="width: 220px"
        />
        <span class="unit">{{ $t('admin.users.yuan') }}</span>
      </el-form-item>
      <el-form-item :label="$t('admin.users.afterBalance')">
        <span class="after-balance">¥{{ previewBalance }}</span>
      </el-form-item>
      <el-form-item :label="$t('admin.vouchers.remark')">
        <el-input
          v-model="form.remark"
          type="textarea"
          :rows="2"
          maxlength="256"
          :placeholder="$t('admin.users.adjustRemarkPlaceholder')"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="$emit('update:visible', false)">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="loading"
          @click="$emit('confirm')"
        >
          {{ $t('common.confirm') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  visible: { type: Boolean, default: false },
  form: { type: Object, required: true },
  loading: { type: Boolean, default: false }
})

defineEmits(['update:visible', 'confirm'])

const previewBalance = computed(() => {
  const current = Number(props.form.currentBalance || 0)
  const amount = Number(props.form.amount || 0)
  const result = props.form.mode === 'set' ? amount : current + amount
  return result.toFixed(2)
})
</script>

<style scoped>
.current-balance {
  font-size: 16px;
  font-weight: 600;
}

.after-balance {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-color-danger);
}

.unit {
  margin-left: 8px;
  color: var(--el-text-color-secondary);
}

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
  width: 100%;
}

.dialog-footer {
  text-align: right;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
