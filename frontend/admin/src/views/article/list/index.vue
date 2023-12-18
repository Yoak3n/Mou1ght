<template>
  <div class="article-list-wrapper">
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
      <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" />
    </n-layout-footer>
  </div>
</template>
<script setup lang="ts">
import {reactive,h} from 'vue'
import { 
  NLayoutHeader, 
  NLayoutContent,
   NLayoutFooter, 
   NBreadcrumb, 
   NBreadcrumbItem, 
   NIcon,
   NDataTable,
   DataTableColumns,
   NButton } from 'naive-ui'
import { MdCash } from '@vicons/ionicons4';

type Song = {
  no: number
  title: string
  author: string
}
const createColumns = ({
  play
}: {
  play: (row: Song) => void
}): DataTableColumns<Song> => {
  return [
    {
      title: 'No',
      key: 'no'
    },
    {
      title: 'Title',
      key: 'title'
    },
    {
      title: '作者',
      key: 'author'
    },
    {
      title: 'Action',
      key: 'actions',
      render (row) {
        return h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: 'small',
            onClick: () => play(row)
          },
          { default: () => 'Play' }
        )
      }
    }
  ]
}

const columns = createColumns({play (row: Song) {
          window.$message.info(`Play ${row.title}`)
        }})


const paginationReactive = reactive({
  page: 2,
  pageSize: 5,
  showSizePicker: true,
  pageSizes: [3, 5, 7],
  onChange: (page: number) => {
  paginationReactive.page = page
  },
  onUpdatePageSize: (pageSize: number) => {
        paginationReactive.pageSize = pageSize
        paginationReactive.page = 1
      }
})
const pagination= paginationReactive 
const data: Song[] = [
  { no: 3, title: 'Wonderwall', author: '4:18' },
  { no: 4, title: "Don't Look Back in Anger", author: '4:48' },
  { no: 12, title: 'Champagne Supernova', author: '7:27' }
]

</script>

<style scoped>
.article-list-wrapper {
  height: 100%;
  background-color: aqua
}
</style>