import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

const app = createApp(App)


// 注册路由插件
import router from './router/index'
app.use(router)



app.mount('#app')

