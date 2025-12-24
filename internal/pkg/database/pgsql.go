package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initPostgres(dsn string) *gorm.DB {
	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return postgresDB
}
