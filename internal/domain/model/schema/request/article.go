package request

type CreateArticleRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Author     string `json:"author"`
	Categories []Sign `json:"categories"`
	Tags       []Sign `json:"tags"`
}

type UpdateArticleRequest struct {
	CreateArticleRequest
	ID string `json:"id"`
}
