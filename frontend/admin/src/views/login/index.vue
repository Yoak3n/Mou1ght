<script setup lang="ts">
import useUserStore from '@/store/modules/user';
import { NForm, NInput, NFormItem, NButton, FormInst, NGrid, NGi,NCheckbox,FormItemRule,FormRules } from 'naive-ui'
import { computed, ref } from 'vue';



import { getTimeDutation } from '@/utils/time';
import { useRouter } from 'vue-router';
const model = ref({ username: '', password: '',check:false })
let isLoading=  ref(false)

const $router = useRouter()
let formRef = ref<FormInst | null>(null)

const userStore = useUserStore()
// 自定义表单验证 
let rules:FormRules=  {
      
        username: {
          required: true,
          message: '请输入用户名',
          trigger: [ 'input'],
        },
        password: {
          required: true,
          validator(_:FormItemRule,value:string){
            if(!value){
              return new Error("请输入密码")
            }else if(value.length<6){
              return new Error("请输入至少6位密码")
            }else{
              return true
            }
          },
          trigger: ['input', 'blur'],
        }
    }
// let vaild = ref(false)
// watch(
//   ()=>formRef.value?.validate,
//   (currentValue,preValue)=>{
//     if (currentValue != preValue){
//       if (!currentValue){
//         window.$message.success('Vaild')
//         return vaild.value = true
//       }else{
//         window.$message.error('Invalid')
//         return vaild.value = false
//       }
//     }
//   }
// )


const showAgreement = ()=>{
  window.$dialog.success({
          content: '用户协议内容',
          title:'用户协议',
          positiveText: '确认',
        })
}

const createStatus =(value:boolean)=>{
  if (!value){
    return "error"
  }else{
    return undefined
  }
}

const showAgreementFeedback = computed(()=>{
  return createFeedback(model.value.check)
})

const createFeedback = (value :boolean):string=>{
  if (!value){
    return "请阅读并勾选用户协议"
  }else{
    return ""
  }
}

const checkAgreementStatus = computed(()=>{
  return createStatus(model.value.check)
})


let showFeedback = ref(false)
const loginHandle =  async()=>{
  showFeedback.value = true
  isLoading.value = true
  if (createStatus(model.value.check) != "error"){
    // 屎代码，之后去看看别人的方案
    // 或者放弃用表单验证是否勾选，直接弹出消息
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
  }else{
    isLoading.value = false
    return
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
          <n-form ref="formRef" :model="model" :rules="rules" :show-require-mark="false">
            <n-form-item label="用户名" path="username" > 
              <n-input v-model:value="model.username"  placeholder="请输入您的用户名" />
            </n-form-item>
            <n-form-item label="密码" path="password">
              <n-input v-model:value="model.password" type="password" placeholder="请输入您的密码" 
                @keydown.enter="loginHandle" />
            </n-form-item >
            <n-form-item size="small" inline :show-label="false" :feedback="showAgreementFeedback" :show-feedback="showFeedback" :validation-status="checkAgreementStatus">
              <n-checkbox v-model:checked="model.check" size="large" style="margin: 10px 0;" path="check"><n-button text @click="showAgreement">请阅读并勾选用户协议 </n-button></n-checkbox>
            </n-form-item >
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
  background: url('@/assets/images/login_background.png')no-repeat;
  background-size: cover;
  height: 100vh;
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