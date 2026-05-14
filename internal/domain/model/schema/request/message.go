package request

type MessagePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type CreateMessageRequest struct {
	Content     string          `json:"content"`
	Position    MessagePosition `json:"position"`
	AuthorIP    string          `json:"author_ip"`
	BoardAnswer string          `json:"board_answer"`
}

type UpdateMessageRequest struct {
	CreateMessageRequest
	ID string `json:"id"`
}

type MessageListRequest struct {
	Sort      string          `json:"sort"`
	DateRange *PostFilterDate `json:"date_range"`
}

type UpdateMessagePositionRequest struct {
	ID       string          `json:"id"`
	Position MessagePosition `json:"position"`
	AuthorIP string          `json:"author_ip"`
}
