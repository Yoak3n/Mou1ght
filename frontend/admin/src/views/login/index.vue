<script setup lang="ts">
import useUserStore from '@/store/modules/user';
import { NForm, NInput, NFormItem, NButton, FormInst, NGrid, NGi,NCheckbox } from 'naive-ui'
import { ref } from 'vue';
import type { loginForm } from "@/api/user/type";


import { getTimeDutation } from '@/utils/time';
import { useRouter } from 'vue-router';

const model = ref<loginForm>({ username: '', password: '' })
let check = ref(false)
let isLoading=  ref(false)

const $router = useRouter()
const formRef = ref<FormInst | null>(null)


const userStore = useUserStore()
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
          trigger: ['input', 'blur']
        }
      }
    }



const showAgreement = ()=>{
  window.$dialog.success({
          content: '用户协议内容',
          title:'用户协议',
          positiveText: '确认',
        })
}

const loginHandle =  async()=>{
  isLoading.value = true
  isLoading.value = true
    try{
        await userStore.userLogin(model.value)
        $router.push('/')
        window.$notification.success({
            content: `Hi,${getTimeDutation()}好！`,
            meta:'登录成功',
            duration:2500,
            closable: false,
        })
       
       // 登录成功的提示信息
    }catch(error){
      window.$notification.error({
            content: 'Error',
            meta:error,
            duration:2500,
            closable: false,
        })
    }finally{
        model.value.username=''
        model.value.password = ''
        isLoading.value = false
    }
  
}

</script>

<template>
  <div class="login-wrapper">
    <n-grid  cols="24" item-responsive>
      <n-gi span="0 600:0 1200:8"></n-gi>
      <n-gi span="0 600:6 1200:8"></n-gi>
      <n-gi span="24 600:18 1200:8">
        <div class="login-card">
          <h1>Hello!</h1>
          <h2>欢迎登录~</h2>
          <n-form ref="formRef" :model="model" :rules="rules">
            <n-form-item label="用户名" > 
              <n-input v-model:value="model.username"  placeholder="请输入您的用户名" />
            </n-form-item>
            <n-form-item label="密码">
              <n-input v-model:value="model.password" type="password" placeholder="请输入您的密码" @input=""
                @keydown.enter="loginHandle" />
            </n-form-item>
              <n-checkbox v-model:checked="check" size="large" style="margin: 10px 0;"><n-button text @click="showAgreement">请阅读用户协议并勾选 </n-button></n-checkbox>
              <!-- 添加窗口 -->
   
            <div style="display: flex; justify-content: flex-end">
                <n-button 
                  class="login-button" 
                  createLocalStorage
                  type="primary" @click="loginHandle"
                  :loading="isLoading"
                  >
                  
                  登录
                </n-button>
              </div>

                
          </n-form>
        </div>
      </n-gi>

    </n-grid>

  </div>
</template>

<style scoped lang="scss">
.login-wrapper {
  .login-card {
    font-size: large;
    color: #fff;
    position: relative;
    width: 80%;
    top: 30vh;
    background: url('@/assets/images/login_form.png')no-repeat;
    background-size: cover;
    padding: 40px;
    margin: 0  auto;
    h1 {
      color: white;
      font-size: 40px;
    }

    h2 {
      font-size: 20px;
      color: #fff;
      margin: 20px 0;
    }

    .login-button {
      
      width: 100%;
    }
  }
}
</style>