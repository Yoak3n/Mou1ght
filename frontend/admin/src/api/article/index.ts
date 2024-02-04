import request from "@/utils/request"
import type {articleView} from "./type"


enum API{
    // REGISTER_URL = "/user/register",
    // LOGIN_URL = "/user/login",
    ARTICLELIST_URL = '/article/list'
}



export const getArticleList = ()=>request.get<any,articleView[]>(API.ARTICLELIST_URL)