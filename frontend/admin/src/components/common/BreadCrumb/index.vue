<template>
    <n-breadcrumb>
        <n-breadcrumb-item v-for="(item, index) in title">
            <n-icon :key="index" :component="MdCash" />{{ item }}</n-breadcrumb-item>
    </n-breadcrumb>
</template>

<script setup lang="ts">
import { toRef } from 'vue';
import { NBreadcrumb, NBreadcrumbItem, NIcon } from 'naive-ui'
import { MdCash } from '@vicons/ionicons4'
const props = defineProps({
    title: {
        type: Array<String>,
        default: () => []
    }
    // icon: {
    //     type: Array<String>,
    //     default: () => []
    // }
})
const title = toRef(props.title)
</script>