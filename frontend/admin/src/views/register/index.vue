<script setup lang="ts">
import { NForm, NInput, NFormItem, NButton, FormInst, NGrid, NGi, NCheckbox, FormItemRule, FormRules } from 'naive-ui'



</script>

<template>
  <div class="register-wrapper">
    <n-grid cols="24" item-responsive>
      <n-gi span="0 600:0 1200:8"></n-gi>
      <n-gi span="0 600:6 1200:8"></n-gi>
      <n-gi span="24 600:18 1200:8">
        <div class="login-card">
          <h1>Hello!</h1>
          <h2>欢迎登录~</h2>
          <n-form ref="formRef" :model="model" :rules="rules" :show-require-mark="false">
            <n-form-item label="用户名" path="username">
              <n-input v-model:value="model.username" placeholder="请输入您的用户名" />
            </n-form-item>
            <n-form-item label="密码" path="password">
              <n-input v-model:value="model.password" type="password" placeholder="请输入您的密码"
                @keydown.enter="loginHandle" />
            </n-form-item>
            <n-form-item size="small" inline :show-label="false" :feedback="showAgreementFeedback"
              :show-feedback="showFeedback" :validation-status="checkAgreementStatus">
              <n-checkbox v-model:checked="model.check" size="large" style="margin: 10px 0;" path="check"><n-button text
                  @click="showAgreement">请阅读并勾选用户协议 </n-button></n-checkbox>
            </n-form-item>
            <!-- 添加窗口 -->
            <div style="display: flex; justify-content: flex-end">
              <n-button class="login-button" createLocalStorage type="primary" @click="loginHandle" :loading="isLoading">
                登录
              </n-button>
            </div>


          </n-form>
        </div>
      </n-gi>

    </n-grid>
  </div>
</template>

<style scoped></style>