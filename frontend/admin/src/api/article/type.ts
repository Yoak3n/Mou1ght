export interface articlePost{
    title:string
    content:string
    author_id: number
    description:string
    category:string
}

export interface articleDetail{
    id:number
    title:string
    content:string
    author_id: number
    author_name: string
    description:string
    category:string
}

export interface articleViewResponse {
    code : number
    data : articles
}
export interface articleView{
    id:number
    title:string
    author_id: number
    author_name: string
    description:string
    category:string
}

interface articles {
    articles : articleView[]
}

