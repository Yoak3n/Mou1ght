<script setup lang="ts">
import { NForm, NInput, NFormItem, NButton } from 'naive-ui'
import {ref} from 'vue'

interface registerForm {
    username:string |null
    password:string |null
    repeatpassword:string |null
}
const rules=  {
      model: {
        username: {
          required: true,
          message: '请输入用户名',
          trigger: 'blur'
        },
        password: {
          required: true,
          message: '请输入正确的密码',
          trigger: [ 'input','blur']
        }
      }
    }
const model = ref<registerForm|null>(null)




</script>

<template>
    <n-form ref="formRef" :model?="model" :rules="rules">
    <n-form-item path="age" label="用户名">
      <n-input v-model:value="model?.username" @keydown.enter.prevent />
    </n-form-item>
    <n-form-item path="password" label="密码">
      <n-input
        v-model:value="model?.password"
        type="password"
        @input=""
        @keydown.enter.prevent
      />
    </n-form-item>
    <n-form-item
      ref="rPasswordFormItemRef"
      first
      path="reenteredPassword"
      label="重复密码"
    >
      <n-input
        v-model:value="model?.repeatpassword"
        :disabled="!model?.password"
        type="password"
        @keydown.enter.prevent
      />
    </n-form-item>
    <n-row :gutter="[0, 24]">
      <n-col :span="24">
        <div style="display: flex; justify-content: flex-end">
          <n-button
            :disabled="model?.username === null"
            round
            type="primary"
            @click=""
          >
            验证
          </n-button>
        </div>
      </n-col>
    </n-row>
  </n-form>
</template>

<style scoped>

</style>