package entity

import "time"

type PostState struct {
	Like   int64 `json:"like"`
	View   int64 `json:"view"`
	Length int64 `json:"length"`
	Status int8  `json:"status"`
}

type PostTimeInfo struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
