<template>
    <div class="layout-wrapper">
        <n-layout has-sider class="menu">
            <n-layout-sider content-style="padding: 12px;fontSize:50px" class="slider" width="300" :native-scrollbar="false">
                <div class="options">
                    <logo />
                    <n-menu :options="menuOptions" inverted @update:value="goRoute" key-field="path" :value="$route.path">
                    </n-menu>

                </div>

            </n-layout-sider>

            <n-layout>
                <Main></Main>
            </n-layout>
        </n-layout>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import {useRoute,useRouter} from 'vue-router';
import { NLayout,  NLayoutSider,  NMenu } from 'naive-ui';
import Logo from '@/components/common/Tiny/Logo.vue';
import menuOptions from './menu/option';
import Main from './main/index.vue'


onMounted(()=>{
    document.title = $route.meta.title as string
})
const $route = useRoute()
const $router = useRouter()
const goRoute = (key: string)=>{
    document.title = $route.meta.title as string
    $router.push(key)
    // $router.push(`/${key}`)
}
</script>

<style scoped lang="scss">
.layout-wrapper {
    color: #fff;
    width: 100%;

    .menu {
        height: 100vh;
       
        .slider {
            font-size: 30px;
            color: #fff;
            background-color: $base-menu-background;
            .options {
                font-size: 50px;
            }
        }
    }

}
</style>