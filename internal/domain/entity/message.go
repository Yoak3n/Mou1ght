package entity

type MessageEntity struct {
	ID       string          `json:"id"`
	Content  string          `json:"content"`
	Position MessagePosition `json:"position"`
	State    PostState       `json:"state"`
	Time     PostTimeInfo    `json:"time"`
	AuthorIP string          `json:"-"`
}

type MessagePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}
