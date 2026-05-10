package instance

import (
	"Mou1ght/internal/repository/interfaces"
	"errors"

	"gorm.io/gorm"
)

type PostRepository struct {
	counter interfaces.PostCounter
	db      *gorm.DB
}

func NewPostRepository(counter interfaces.PostCounter, db *gorm.DB) *PostRepository {
	return &PostRepository{
		counter: counter,
		db:      db,
	}
}

func (p *PostRepository) UpdatePostStatus(postType string, id string, status int8) error {
	tableName := ""
	switch postType {
	case "article":
		tableName = "article_tables"
	case "message":
		tableName = "message_tables"
	case "sharing":
		tableName = "sharing_tables"
	default:
		return errors.New("invalid post type")
	}
	if status < 0 || status > 2 {
		return errors.New("status must be 0, 1, or 2")
	}

	result := p.db.Table(tableName).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("post not found")
	}
	return nil
}

type PostCounter struct {
	counter *counterBuffer
}

func NewPostCounter(counter *counterBuffer) *PostCounter {
	return &PostCounter{
		counter: counter,
	}
}

func (p *PostCounter) BumpView(postType, id string, n int64) {
	if p.counter == nil {
		return
	}
	p.counter.bumpView(postType, id, n)
}

func (p *PostCounter) BumpLike(postType, id string, n int64) {
	if p.counter == nil {
		return
	}
	p.counter.bumpLike(postType, id, n)
}

func (p *PostCounter) GetCounterDelta(postType, id string) (viewDelta int64, likeDelta int64) {
	if p.counter == nil {
		return 0, 0
	}
	return p.counter.get(postType, id)
}
