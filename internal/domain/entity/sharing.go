package entity

type SharingEntity struct {
	ID          string             `json:"id"`
	Content     string             `json:"content"`
	Author      UserEntity         `json:"author"`
	Tags        []PostSign         `json:"tags"`
	State       PostState          `json:"state"`
	Time        PostTimeInfo       `json:"time"`
	Attachments []AttachmentEntity `json:"attachments"`
}
