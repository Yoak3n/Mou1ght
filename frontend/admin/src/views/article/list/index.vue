<template>
  <div class="article-list-wrapper">
    <n-layout-header :bordered="true" class="header">
      <n-space justify="space-between" >
        <BreadCrumb :title="path"/>
        <n-button>创建文章</n-button>
      </n-space>


   

    </n-layout-header>
    <n-layout-content content-style="padding: 24px;">


    </n-layout-content>
    <n-layout-footer>
      <n-data-table :columns="columns" :data="articleList" :pagination="pagination" :bordered="false" />
    </n-layout-footer>
  </div>
</template>
<script setup lang="ts">
import { reactive, h, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import {
  NLayoutHeader,
  NLayoutContent,
  NSpace,
  NLayoutFooter,
  NDataTable,
  DataTableColumns,
  NButton
} from 'naive-ui'

import { articleView } from '@/api/article/type'
import { reqArticleList } from '@/api/article'
import BreadCrumb  from '@/components/common/BreadCrumb/index.vue'
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


const $route = useRoute()
let path = reactive<string[]>([])
$route.matched.forEach((item) => {
  if(item.meta.title  != undefined){
    path.push(item.meta.title as string )
  }
})
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

<style scoped lang="scss">
.article-list-wrapper {
  margin-top: 1rem;
  height: 100%;
  background-color: aqua;
  .header{
    padding: 0 1rem
  }
}</style> 