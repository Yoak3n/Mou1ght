<script setup lang="ts">
import {ref } from 'vue';
import { NConfigProvider, } from 'naive-ui'
import type { GlobalTheme } from 'naive-ui'


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
