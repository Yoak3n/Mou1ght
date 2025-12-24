package request

type CategoryRequest struct {
	Label  string `json:"label"`
	Parent string `json:"parent"`
}
