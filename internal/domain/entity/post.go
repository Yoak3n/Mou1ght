package entity

type PostState struct {
	Like   int64  `json:"like"`
	View   int64  `json:"view"`
	Length int64  `json:"length"`
	Status string `json:"status"`
}

type PostTimeInfo struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PostSign struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

func StatusIntToString(status int8) string {
	switch status {
	case 0:
		return "draft"
	case 1:
		return "published"
	case 2:
		return "archived"
	case 3:
		return "pending"
	default:
		return "unknown"
	}
}
