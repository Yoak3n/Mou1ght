<script setup lang="ts">
import {ref } from 'vue';
import { NConfigProvider, } from 'naive-ui'
import type { GlobalTheme } from 'naive-ui'
import NaiveProvider from './components/common/NaiveProvider/index.vue'

let theme = ref(<GlobalTheme | null>null)
</script>

<template>
  <NConfigProvider
    class="h-full"
    :theme="theme"
  >
    <NaiveProvider>
      <RouterView />
    </NaiveProvider>
  </NConfigProvider>
</template>

<style scoped></style>
