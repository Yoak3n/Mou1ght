export interface loginForm {
    username:string,
    password:string
}

interface dataType{
    token?:string
    message?:string
} 
//登录接口返回数据
export interface loginResponseData{
    code :number,
    data :dataType
}

export interface userInfo {
    name:number,
    email:string,
    avatar:string,
    nick_name:string,
    // password:string,
    desc:string,
    role:string[],
    // buttons:string[],
    // routes:string[],
    // token:string

}
interface user{
    user:userInfo
}

// 定义服务器返回用户信息相关的数据类型
export interface userResponseData {
    code:number
    data:user
}