package request

type CreateSharingRequest struct {
    Content    string `json:"content"`
    Author     string `json:"author"`
    Attachment string `json:"attachment"`
    Tags       []Sign `json:"tags"`
}

type UpdateSharingRequest struct {
    CreateSharingRequest
    ID string `json:"id"`
}

