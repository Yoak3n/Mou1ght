// 对外暴露全局组件插件
import SvgIcon from './common/SvgIcon/index.vue';



const allGlobalComponent:any ={SvgIcon}
export default {
    // app.use自动执行install方法
    install(app:any){
        Object.keys(allGlobalComponent).forEach(key=>{
            app.component(key,allGlobalComponent[key])
        })
    }
}