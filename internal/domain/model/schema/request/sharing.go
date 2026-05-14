package request

type CreateSharingRequest struct {
	Content       string   `json:"content"`
	Author        string   `json:"author"`
	AttachmentIDs []string `json:"attachment_ids"`
	Attachment    string   `json:"attachment,omitempty"`
	Tags          []Sign   `json:"tags"`
}

type UpdateSharingRequest struct {
	CreateSharingRequest
	ID string `json:"id"`
}
