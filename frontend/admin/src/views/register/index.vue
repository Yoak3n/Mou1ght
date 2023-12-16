<script setup lang="ts">
import { ref } from 'vue'
import { NForm, NInput, NFormItem, NButton, FormInst, NGrid, NGi, FormItemRule, FormRules } from 'naive-ui'

import useUserStore from '@/store/modules/user';

const userStore = useUserStore()
const model = ref({ username: '', password: '', repeat_password: '', check: false })
let isLoading = ref(false)

let formRef = ref<FormInst | null>(null)
let rules: FormRules = {
  username: {
    required: true,
    message: '请输入用户名',
    trigger: ['input'],
  },
  password: {
    required: true,
    validator(_: FormItemRule, value: string) {
      if (!value) {
        return new Error("请输入密码")
      } else if (value.length < 6) {
        return new Error("请输入至少6位密码")
      } else {
        return true
      }
    },
    trigger: ['input', 'blur'],
  },
  repeat_password: {
    required: true,
    validator(_: FormItemRule, value: string) {
      if (!value) {
        return new Error("请再次输入密码")
      } else if (value !== model.value.password) {
        return new Error("两次密码输入不一致")
      } else {
        return true
      }
    },
    trigger: ['input', 'blur'],
  }
}


const registerHandle = async () => {
  isLoading.value = true
  if (model.value.password === model.value.repeat_password) {
    try {
      await userStore.userRegister(model.value)
    } catch (error) {
      window.$notification.error({
        content: 'Error',
        meta: error,
        duration: 2500,
        closable: false,
      })
    } finally {
      model.value.username = ''
      model.value.password = ''
      isLoading.value = false
    }
  }

}

</script>

<template>
  <div class="register-wrapper">
    <n-grid cols="24" item-responsive>
      <n-gi span="0 600:0 1200:8"></n-gi>
      <n-gi span="24 600:18 1200:8">
        <div class="register-card">
          <h1>Hello!</h1>
          <h2>欢迎来到Mou1ght~</h2>
          <n-form ref="formRef" :model="model" :rules="rules" :show-require-mark="false">
            <n-form-item label="用户名" path="username">
              <n-input v-model:value="model.username" placeholder="请输入您的用户名" />
            </n-form-item>
            <n-form-item label="密码" path="password" show-require-mark>
              <n-input v-model:value="model.password" type="password" placeholder="请输入您的密码"
                show-password-on="mousedown" />
            </n-form-item>
            <n-form-item label="确认密码" path="repeat_password" show-require-mark>
              <n-input v-model:value="model.repeat_password" type="password" placeholder="请再次输入密码"
                show-password-on="mousedown"></n-input>
            </n-form-item>
            <!-- <n-form-item size="small" inline :show-label="false">
              <n-checkbox v-model:checked="model.check" size="large" style="margin: 10px 0;" path="check"><n-button
                  text>请阅读并勾选用户协议 </n-button></n-checkbox>
            </n-form-item> -->
            <!-- 添加窗口 -->
            <div style="display: flex; justify-content: flex-end">
              <n-button class="register-button" createLocalStorage type="primary" :loading="isLoading"
                @click="registerHandle">
                注册
              </n-button>
            </div>
          </n-form>
        </div>
      </n-gi>
      <n-gi span="0 600:6 1200:8">

      </n-gi>

    </n-grid>
  </div>
</template>

<style scoped>
.register-wrapper {
  background: url('@/assets/images/login_background.png')no-repeat;
  background-size: cover;
  height: 100vh;

  .register-card {
    font-size: large;
    color: #fff;
    position: relative;
    width: 80%;
    top: 30vh;
    background: url('@/assets/images/login_form.png')no-repeat;
    background-size: cover;
    padding: 40px;
    margin: 0 auto;

    h1 {
      color: white;
      font-size: 40px;
    }

    h2 {
      font-size: 20px;
      color: #fff;
      margin: 20px 0;
    }

    .register-button {
      width: 100%;
    }
  }
}
</style>