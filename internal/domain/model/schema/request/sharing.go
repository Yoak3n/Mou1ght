package request

type CreateSharingRequest struct {
	Content string `json:"content"`
	Author  string `json:"author"`
	// TODO 需要看附件模块怎么实现
	Attachment string `json:"attachment"`
	Tags       []Sign `json:"tags"`
}

type UpdateSharingRequest struct {
	CreateSharingRequest
	ID string `json:"id"`
}
