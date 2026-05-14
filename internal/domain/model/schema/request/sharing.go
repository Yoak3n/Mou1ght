package request

type CreateSharingRequest struct {
	Content       string   `json:"content"`
	Author        string   `json:"author"`
	Private       bool     `json:"private"`
	AttachmentIDs []string `json:"attachment_ids"`
	Tags          []Sign   `json:"tags"`
}

type UpdateSharingRequest struct {
	CreateSharingRequest
	ID string `json:"id"`
}
