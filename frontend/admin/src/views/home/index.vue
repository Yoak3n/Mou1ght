<template>
    <div class="homepage">
        <n-layout-header :bordered="true">
            <n-breadcrumb v-for="(item, index) in path">
                <n-breadcrumb-item :key="index">
                    <n-icon :component="MdCash" />{{ item }}</n-breadcrumb-item>
            </n-breadcrumb>
        </n-layout-header>
        <n-layout-content content-style="padding: 24px;height: 100vh;background-color:#fff">
            <div class="user-info">
                <n-card title="个人信息" hoverable embedded>
                    <template #header-extra>
                        <n-button text  @click="updateUserInformation">修改信息</n-button>
                    </template>
                    <n-space align="center">
                        <n-avatar lazy
                            :src="userInfo.avatar != '' ? userInfo.avatar : 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg'"
                            size="large" />
                        {{ userInfo.nick_name }}
                    </n-space>
                    <template #footer>
                        <n-ellipsis style="max-width: 240px">
                            住在我心里孤独的 孤独的海怪 痛苦之王 开始厌倦 深海的光 停滞的海浪
                        </n-ellipsis>
                    </template>
                </n-card>
            </div>
        </n-layout-content>
        <n-layout-footer>
        </n-layout-footer>
    </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
    NBreadcrumb, NBreadcrumbItem,
    NIcon,NButton,
    NLayoutHeader, NLayoutContent, NLayoutFooter,
    NCard,
    NAvatar, NEllipsis,
    NSpace
} from 'naive-ui'
import { MdCash } from '@vicons/ionicons4'
import { useRoute, useRouter } from 'vue-router';

import { reqUserInfo } from '@/api/user';
import type { userInfo } from '@/api/user/type'
const $route = useRoute();
const $router = useRouter();
let path = ref<Array<string>>([])
let userInfo = ref<userInfo>({ name: 0, nick_name: '', email: '', avatar: '', role: [], desc: '' })
onMounted(() => {
    reqUserInfo().then((v) => {
        userInfo.value = v.data.user
    }).catch((err) => {
        window.$message.error({ content: err, duration: 2500 })
        $router.push('/login')
    })
})


$route.path.split('/').forEach((item) => {
    if (item != '') {
        path.value.push(item)
    }

})


const updateUserInformation = ()=>{
    console.log(userInfo.value);
}
</script>

<style scoped lang="scss">
.homepage {
    margin-top: 1rem;
}
</style>