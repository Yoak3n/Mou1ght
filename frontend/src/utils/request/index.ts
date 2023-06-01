// 进行axios二次封装：使用请求与响应拦截器
import  axios   from "axios"
// 第一步：利用axios的create方法，去创建axios实例（其他的配置：基础路径、超时限制）
let request = axios.create({
    baseURL:import.meta.env.VITE_APP_BASE_API,
    timeout:5000,
});

import {useUserStore} from "@/store/modules/user";

let userStore = useUserStore()


// 第二步：为axios实例添加请求与响应拦截器
request.interceptors.request.use((config)=>{
    // config配置对象，headers属性请求头，经常给服务器携带公共参数

    // 添加token认证请求头
    config.headers.set("Authorization","Bearer "+userStore.token)
    return config
})

request.interceptors.response.use((response)=>{
    // 简化输出
    return response.data
},(error)=>{
    // 失败的回调，处理http网络错误
    let message = ''
    let statusCode = error.response.status;
    switch(statusCode){
        case 401:
            message = "TOKEN过期"
            break;
        case 403:
            message = "无授权"
            break;
        case 404:
            message = "请求地址错误"
            break;
        case 500:
            message = "服务器出现问题"
            break;
        default:
            message = "网络出现问题"
            break;
    }
    window.$message.error(message)

    return Promise.reject(error)
})

export default request;