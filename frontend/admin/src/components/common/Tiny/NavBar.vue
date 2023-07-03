<template>
    <div class="nav-bar-wrapper">
        <div class="logo-wrapper">
            <img src="" alt="">
        </div>
        <div class="categories-wrapper" >

            <ul class="categories" >
                <li class="list" v-for="item in option" @click="jumpCategory(item.url)">
                    <button class="list-button" >{{ item.label }}</button>
                </li>
            </ul>
            <div class="search-wrapper">
                <n-icon :component="Search" />
                <input placeholder="搜索" size="large">            
            </div>
        </div>
        
        
        
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { NIcon } from 'naive-ui';
import { Search } from '@vicons/ionicons5';

let $router = useRouter()
// let $route = useRoute()
let option = ref([{label:'',icon:'',url:''}])
option.value=[
    {label:'主页',icon:'',url:'/'},
    {label:'文章',icon:'',url:'/article'},
    {label:'登录',icon:'',url:'/login'},
    {label:'注册',icon:'',url:'/register'}
]
    
const jumpCategory= (url:string)=>{
    $router.push(url)
}
</script>

<style scoped lang="scss">
.nav-bar-wrapper{
    background-color: aqua;
    display: flex;
    justify-content: space-between;
    box-shadow: 0 1px 2px #666;
    .logo-wrapper{
        width: 20px;
    }
    .categories-wrapper{
        height: 50px;
        display: flex;
        background-color: antiquewhite;
        .categories{
            display: flex;
            list-style: none;
            -ms-flex-align: end;
            flex-direction: row;
            .list{
                height: 100%;
                
                line-height: 50px;
                background-color: #fff;
                text-align: center;
                .list-button{
                    height: 100%;
                    padding: 0 1rem;
                    border: none;
                    background-color: transparent;
                    &:hover{
                        background-color: antiquewhite;
                    }
                }
            }

        }
    }
    .search-wrapper{
        margin: 0 10px;
        line-height: 50px;
        height: 50px;
    }
}


</style>