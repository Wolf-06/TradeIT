package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	dsn := "host=localhost user=postgres password=dev123 dbname=tradeit port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error connecting to database: %e", err)
		return
	}
	DB = db
}

func SetDB() *gorm.DB {
	return DB
}
