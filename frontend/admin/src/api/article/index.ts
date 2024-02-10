import request from "@/utils/request"
import type {articleDetail,articleViewResponse,articlePost} from "./type"


enum API{
    ARTICLE_INFO_URL = "/article/info",
    ARTICLA_ADD_URL = "/article/add",
    ARTICLE_LIST_URL = '/article/list'
}
export const reqArticleInfo = (id:number)=>request.get<any,articleDetail>(API.ARTICLE_INFO_URL+`/${id}`)
export const reqArticleAdd = (data:articlePost)=>request.post<any,articleDetail>(API.ARTICLA_ADD_URL,data)
export const reqArticleList = ()=>request.get<any,articleViewResponse>(API.ARTICLE_LIST_URL)