<template>
    <div class="homepage">
        <n-layout-header :bordered="true">
            
        </n-layout-header>
        <n-layout-content content-style="padding: 24px;">
            <n-breadcrumb>
                <n-breadcrumb-item>
                    <n-icon :component="MdCash" /> 北京总行</n-breadcrumb-item>
                <n-breadcrumb-item>
                    <n-icon :component="MdCash" /> 天津分行</n-breadcrumb-item>
                <n-breadcrumb-item>
                    <n-icon :component="MdCash" /> 平山道支行</n-breadcrumb-item>
            </n-breadcrumb>
        </n-layout-content>
        <n-layout-footer>
            12345
        </n-layout-footer>
    </div>
</template>
<script setup lang="ts">
import {NBreadcrumb,NBreadcrumbItem,NIcon, NLayoutHeader, NLayoutContent, NLayoutFooter } from 'naive-ui'
import { MdCash } from '@vicons/ionicons4'

</script>

<style scoped></style>