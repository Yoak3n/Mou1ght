<template>
    <div class="homepage-wrapper">
        <n-layout has-sider class="homepage">
            <n-layout-sider 
                content-style="padding: 24px;" 
                class="slider" 
                width="300"
                :native-scrollbar="false">
                <div>
                    <logo />
                    
                    <n-menu
                    :options="menuOptions"
                    inverted
                    >
                        
                    </n-menu>
        
                </div>

            </n-layout-sider>
            <n-layout>
                <n-layout-header>面包屑</n-layout-header>
                <n-layout-content content-style="padding: 24px;">
                    789
                </n-layout-content>
                <n-layout-footer>
                    12345
                </n-layout-footer>
            </n-layout>
        </n-layout>
    </div>
</template>

<script setup lang="ts">
import {h,Component} from 'vue'
import { NLayout, NLayoutContent, NLayoutFooter, NLayoutSider, NLayoutHeader, NMenu,MenuOption,NIcon} from 'naive-ui';
import Logo from '@/components/common/Tiny/Logo.vue'
import {
  BookOutline as BookIcon,
  PersonOutline as PersonIcon,
  WineOutline as WineIcon
} from '@vicons/ionicons5'


function renderIcon (icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

let menuOptions:MenuOption[] = [
    {
        label:'主页',
        key: 'hear-the-wind-sing',
        icon: renderIcon(BookIcon)
    }, {
        label: '数据大图',
        key: 'a-wild-sheep-chase',
        icon: renderIcon(PersonIcon),
    }, {
        label: '权限管理',
        key: 'auth-admin',
        icon: renderIcon(WineIcon),
    },
]

</script>

<style scoped lang="scss">
.homepage-wrapper {
    color: #fff;
    width: 100%;
    .homepage{
        height: 100vh;
        .slider{
            color: #fff;
            background-color: $base-menu-background;
        }
    }

}
</style>