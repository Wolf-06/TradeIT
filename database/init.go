package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	cred := map[string]string{
		"host":   os.Getenv("db_host"),
		"user":   os.Getenv("db_user"),
		"passwd": os.Getenv("db_passwd"),
		"dbName": os.Getenv("db_name"),
		"port":   os.Getenv("db_port"),
	}
	fmt.Println(cred)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", cred["host"], cred["user"], cred["passwd"], cred["dbName"], cred["port"])
	fmt.Println(dsn)
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
