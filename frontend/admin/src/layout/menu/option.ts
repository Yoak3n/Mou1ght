import { h, Component} from 'vue'
import { MenuOption,NIcon } from "naive-ui"
import * as icon from '@vicons/ionicons5'

function renderIcon(icon: Component) {
    return () => h(NIcon, null, { default: () => h(icon) })
}
let  menuOptions: MenuOption[] = [
    {
        label: '主页',
        key: 'home',
        path:'/home',
        icon: renderIcon(icon.BookOutline)
    }, {
        label: '数据大图',
        key: 'screen',
        path:'/screen',
        icon: renderIcon(icon.PersonAddOutline),
    }, {
        label: '权限管理',
        key: 'acl',
        path:'/acl',
        icon: renderIcon(icon.LockClosed),
        children:[
            {
                label:'用户管理',key:'user',path:'/acl/user'
            },{
                label:'角色管理',key:'role',path:'/acl/role'
            },{
                label:'菜单管理',key:'permission',path:'/acl/permission'
            },
        ]
    },{
        label: '文章管理',
        key: 'article',
        path:'/article',
        icon: renderIcon(icon.CarOutline),
        children:[
            {
                label:'文章列表',key:'article/list',path:'/article/list'
            },{
                label:'SPU管理',key:'article/spu',path:'/article/spu'
            },{
                label:'SKU管理',key:'article/sku',path:'/article/sku'
            },{
                label:'属性管理',key:'article/attr',path:'/article/attr'
            }
        ]
    }
]
export default menuOptions