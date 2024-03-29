import {RouteRecordRaw}from 'vue-router'

const routes:RouteRecordRaw[] =[
    {
        path:'/',
        component:()=>import('@/layout/index.vue'),
        redirect:'/home',
        children:[
            {
                path:'/home',
                component:()=>import('@/views/home/index.vue'),
                name:'Home',
                meta:{
                    title:'首页',
                    hidden:false,
                    icon: 'HomeFiled'
                }
            },{
                path:'/screen',
                component:()=>import('@/views/screen/index.vue'),
                name:'Screen',
                meta:{
                    hidden:false,
                    title:'数据大屏',
                    icon:'LaptopOutline'
                }
            }
            //     path:'/home',    
        ]
    },{
        path:'/acl',
        component:()=>import('@/layout/index.vue'),
        name:'Acl',
        meta:{
            title:'权限管理',
            hidden:false,
            icon:'LockClosed'
        },
        children:[
            {
                path:'/acl/user',
                component:()=>import('@/views/acl/user/index.vue'),
                name:'User',
                meta:{
                    title:'用户管理',
                    hidden:false,
                    icon:'PersonOutline'
                }
            },{
                path:'/acl/role',
                component:()=>import('@/views/acl/role/index.vue'),
                name:'Role',
                meta:{
                    title:'角色管理',
                    hidden:false,
                    icon:'PersonOutline'
                }
            },{
                path:'/acl/permission',
                component:()=>import('@/views/acl/permission/index.vue'),
                name:'Permission',
                meta:{
                    title:'菜单管理',
                    hidden:false,
                    icon:'PersonOutline'
                }
            },
        ]
    },{
        path:'/article',
        component:()=>import('@/layout/index.vue'),
        name:'Product',
        meta:{
            title:'文章管理',
            hidden:false,
            icon:'Goods'
        },
        children:[
            {
                path:'/article/list',
                component:()=>import('@/views/article/list/index.vue'),
                name:'ArticleList',
                meta:{
                    title:'文章列表',
                    icon:'PersonOutline'
                }
            },{
                path:'/article/spu',
                component:()=>import('@/views/article/spu/index.vue'),
                name:'Spu',
                meta:{
                    title:'SPU管理',
                    icon:'PersonOutline'
                }
            },{
                path:'/article/sku',
                component:()=>import('@/views/article/sku/index.vue'),
                name:'Sku',
                meta:{
                    title:'SKU管理',
                    icon:'PersonOutline'
                }
            },{
                path:'/article/attr',
                component:()=>import('@/views/article/attr/index.vue'),
                name:'Attr',
                meta:{
                    title:'属性管理',
                    icon:'PersonOutline'
                }
            },
        ]
    },{
       path:'/login',
        component:()=>import('@/views/login/index.vue'),
        name:'login'
    },{
        path:'/register',
        component:()=>import('@/views/register/index.vue'),
        name:'register'
    },{
        path:'/404',
        component:()=>import('@/views/exceptions/404.vue'),
        name:'404'
    },{
        path:'/:pathMatch(.*)*',
        redirect:'/404',
        name:'Any'
    },{
        path:'/500',
        component:()=>import('@/views/exceptions/500.vue'),
        name:'500'},
] 

export default  routes