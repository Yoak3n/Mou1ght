package instance

import (
	"Mou1ght/internal/pkg/database"
	"sync"

	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

var once sync.Once
var db *Database

func NewDatabase() *Database {
	return &Database{DB: database.InitDatabase("", "")}
}

func UseDatabase() *Database {
	once.Do(func() {
		db = NewDatabase()
	})
	return db
}
