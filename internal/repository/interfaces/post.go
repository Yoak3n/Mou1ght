package interfaces

type PostCounter interface {
	BumpView(postType, id string, n int64)
	BumpLike(postType, id string, n int64)
	GetCounterDelta(postType, id string) (viewDelta int64, likeDelta int64)
}

type PostRepository interface {
	UpdatePostStatus(postType string, id string, status int8) error
}
