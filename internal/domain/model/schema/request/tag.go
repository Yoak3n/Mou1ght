package request

type CreateTagRequest struct {
	Label string `json:"label"`
}

type UpdateTagRequest struct {
	Label string `json:"label"`
}
