import {RouteRecordRaw}from 'vue-router'




const routes:RouteRecordRaw[] =[
    {path:'/',component:()=>import('@/views/HomePage.vue')},
    {path:'/login',component:()=>import('@/views/login/index.vue'),name:'login'},
    {path:'/register',component:()=>import('@/views/register/index.vue'),name:'register'},
    {path:'/404',component:()=>import('@/views/exceptions/404.vue'),name:'404'},
    {path:'/:pathMatch(.*)*',redirect:'/404',name:'Any'},
    {path:'/500',component:()=>import('@/views/exceptions/500.vue'),name:'500'},
] 

export default  routes