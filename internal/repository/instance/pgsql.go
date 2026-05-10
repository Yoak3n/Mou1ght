package instance

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/database"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Database struct {
	DB      *gorm.DB
	counter *counterBuffer
}

var once sync.Once
var db *Database

func NewDatabase() *Database {
	d := &Database{DB: database.InitDatabase()}
	d.counter = newCounterBuffer(d, 5*time.Second, 10000)
	d.counter.start()
	return d
}

func UseDatabase() *Database {
	once.Do(func() {
		db = NewDatabase()
	})
	return db
}

func (d *Database) applyCounterDelta(postType, id string, viewDelta, likeDelta int64) {
	if viewDelta == 0 && likeDelta == 0 {
		return
	}

	var model any
	switch postType {
	case "article":
		model = &table.ArticleTable{}
	case "sharing":
		model = &table.SharingTable{}
	case "message":
		model = &table.MessageTable{}
	default:
		return
	}

	updates := make(map[string]any, 2)
	if viewDelta != 0 {
		updates["view"] = gorm.Expr("view + ?", viewDelta)
	}
	if likeDelta != 0 {
		updates["like"] = gorm.Expr("like + ?", likeDelta)
	}
	_ = d.DB.Model(model).Where("id = ? AND status = 1", id).Updates(updates).Error
}

func (d *Database) GetCounter() *counterBuffer {
	return d.counter
}

type counterKey struct {
	postType string
	id       string
}

type counterDelta struct {
	view int64
	like int64
}

type counterBuffer struct {
	db            *Database
	flushInterval time.Duration
	maxKeys       int
	mu            sync.Mutex
	deltas        map[counterKey]counterDelta
}

func newCounterBuffer(db *Database, flushInterval time.Duration, maxKeys int) *counterBuffer {
	return &counterBuffer{
		db:            db,
		flushInterval: flushInterval,
		maxKeys:       maxKeys,
		deltas:        make(map[counterKey]counterDelta),
	}
}

func (c *counterBuffer) start() {
	go func() {
		t := time.NewTicker(c.flushInterval)
		defer t.Stop()
		for range t.C {
			c.flush()
		}
	}()
}

func (c *counterBuffer) bumpView(postType, id string, n int64) {
	if n <= 0 {
		return
	}
	c.mu.Lock()
	k := counterKey{postType: postType, id: id}
	d := c.deltas[k]
	d.view += n
	c.deltas[k] = d
	needFlush := c.maxKeys > 0 && len(c.deltas) >= c.maxKeys
	c.mu.Unlock()
	if needFlush {
		c.flush()
	}
}

func (c *counterBuffer) bumpLike(postType, id string, n int64) {
	if n <= 0 {
		return
	}
	c.mu.Lock()
	k := counterKey{postType: postType, id: id}
	d := c.deltas[k]
	d.like += n
	c.deltas[k] = d
	needFlush := c.maxKeys > 0 && len(c.deltas) >= c.maxKeys
	c.mu.Unlock()
	if needFlush {
		c.flush()
	}
}

func (c *counterBuffer) get(postType, id string) (viewDelta int64, likeDelta int64) {
	c.mu.Lock()
	d := c.deltas[counterKey{postType: postType, id: id}]
	c.mu.Unlock()
	return d.view, d.like
}

func (c *counterBuffer) flush() {
	c.mu.Lock()
	if len(c.deltas) == 0 {
		c.mu.Unlock()
		return
	}
	snapshot := c.deltas
	c.deltas = make(map[counterKey]counterDelta)
	c.mu.Unlock()

	for k, d := range snapshot {
		if d.view == 0 && d.like == 0 {
			continue
		}
		c.db.applyCounterDelta(k.postType, k.id, d.view, d.like)
	}
}
