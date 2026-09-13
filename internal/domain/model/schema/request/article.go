package request

type CreateArticleRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Author     string `json:"author"`
	Categories []Sign `json:"categories"`
	Tags       []Sign `json:"tags"`
	// AttachmentIDs 文章关联的附件（可为音频等）；nil 表示不修改，空切片表示清空
	AttachmentIDs []string `json:"attachment_ids"`
}

type UpdateArticleRequest struct {
	CreateArticleRequest
	ID string `json:"id"`
}
