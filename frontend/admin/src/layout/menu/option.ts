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
        label: '商品管理',
        key: 'product',
        path:'/product',
        icon: renderIcon(icon.CarOutline),
        children:[
            {
                label:'品牌管理',key:'product/trademark',path:'/product/trademark'
            },{
                label:'SPU管理',key:'product/spu',path:'/product/spu'
            },{
                label:'SKU管理',key:'product/sku',path:'/product/sku'
            },{
                label:'属性管理',key:'product/attr',path:'/product/attr'
            }
        ]
    }
]
export default menuOptions