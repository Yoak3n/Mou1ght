import request from "@/utils/request"
import type{loginForm,loginResponseData,userResponseData} from './type'
enum API{
    REGISTER_URL = "/user/register",
    LOGIN_URL = "/user/login",
    USERINFO_URL = "/user/info"
}
// 暴露请求函数
// 登录接口方法
export const reqLogin = (data:loginForm)=>request.post<any,loginResponseData>(API.LOGIN_URL+`/${data.username}/${data.password}`);
export const reqRegister = (data:loginForm)=>request.post<any,loginResponseData>(API.REGISTER_URL+`/${data.username}/${data.password}`,)
export const reqUserInfo = ()=>request.get<any,userResponseData>(API.USERINFO_URL)