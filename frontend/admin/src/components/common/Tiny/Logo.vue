<template>
    <div class="logo-wrapper">
        <n-space align="center">
            <img :src="setting.icon" alt="logo" class="logo">
            <p>{{ setting.title }}</p>
        </n-space>
        
        
    </div>
</template>

<script setup lang="ts">
import { NSpace } from 'naive-ui';
import setting from '@/setting'
</script>

<style scoped lang="scss">
.logo-wrapper{
    width: 100%;
    height: $base-menu-logo-height;
    img{
        width: 40px;
    }
    p{
        font-size: $base-logo-title-fontSize;
    }
}

</style>