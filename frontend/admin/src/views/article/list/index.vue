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
      <n-data-table :columns="columns" :data="articleList" :pagination="pagination" :bordered="false" />
    </n-layout-footer>
  </div>
</template>
<script setup lang="ts">
import { reactive, h, onMounted } from 'vue'
import {
  NLayoutHeader,
  NLayoutContent,
  NLayoutFooter,
  NBreadcrumb,
  NBreadcrumbItem,
  NIcon,
  NDataTable,
  DataTableColumns,
  NButton
} from 'naive-ui'
import { MdCash } from '@vicons/ionicons4';

import { articleView } from '@/api/article/type'
import { reqArticleList } from '@/api/article'

type Article = {
  no: number
  title: string
  category: string
  author: string
}
const createColumns = ({
  play
}: {
  play: (row: Article) => void
}): DataTableColumns<Article> => {
  return [
    {
      title: 'No.',
      key: 'no'
    },
    {
      title: 'Title',
      key: 'title'
    }, {
      title: 'Category',
      key: 'category'
    }, {
      title: '作者',
      key: 'author'
    },
    {
      title: 'Action',
      key: 'action',
      render(row) {
        return h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: 'small',
            onClick: () => play(row)
          },
          { default: () => '查看' }
        )
      }
    }
  ]
}



const columns = createColumns({
  play(row: Article) {
    window.$message.info(`Play ${row.title}`)
  }
})



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
const pagination = paginationReactive

let articleList = reactive<Article[]>([])

onMounted(() => {
  reqArticleList().then((v) => {
    const as: articleView[] = v.data.articles
    if (Array.isArray(as)){
      as.forEach((a) => {
      articleList.push({
        no: a.id,
        title: a.title,
        category: a.category,
        author: a.author_name
      })
    })
    }
  })
})


</script>

<style scoped>
.article-list-wrapper {
  height: 100%;
  background-color: aqua
}</style>