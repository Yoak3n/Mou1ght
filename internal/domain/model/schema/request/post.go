package request

// Sign is a struct for category or tag
type Sign struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type PostListRequest struct {
	Filter PostFilter     `json:"filter"`
	Data   PostFilterData `json:"data"`
}

type PostFilter struct {
	Typ       string          `json:"type"`
	DateRange *PostFilterDate `json:"date_range,omitempty"`
	Sort      string          `json:"sort"`
}

type PostFilterData struct {
	Keyword []string `json:"keyword"`
}

type PostFilterDate struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
