package entity

type ArticleEntity struct {
	ID         string            `json:"id"`
	Title      string            `json:"title"`
	Content    string            `json:"content"`
	Categories []PostSign        `json:"categories"`
	Tags       []PostSign        `json:"tags"`
	Author     UserEntity        `json:"author"`
	State      PostState         `json:"state"`
	Time       PostTimeInfo      `json:"time"`
	// Attachments 仅在文章详情中填充（列表接口不加载）
	Attachments []AttachmentEntity `json:"attachments,omitempty"`
}
