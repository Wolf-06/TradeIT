package models

import (
	"TradeIT/database"
	"time"
)

type User struct {
	Id    int     `gorm:"PrimaryKey"`
	Name  string  `json:"name" validate:"required"`
	Email string  `json:"email" validate:"email,required"`
	Funds float32 `json:"funds"`
}

type Credential struct {
	Id     int    `gorm:"PrimaryKey"`
	Email  string `json:"email" validate:"email,  required"`
	Passwd string `json:"passwd" validate:"required, min=6"`
	Token  string
}

type Order struct {
	Id         int       `gorm:"PrimaryKey"`
	User_id    int       `json:"user_id" validate:"required"`
	Order_type string    `json:"type" validate:"required, oneof= buy sell"`
	Stock      string    `json:"stock" validate:"required"`
	Price      float32   `json:"price" gorm:"type:float" validate:"required gt=0"`
	Quantity   int       `json:"quantity" validate:"required gt=0"`
	Status     string    `json:"status" validate:"required oneof= executed pending cancelled"`
	Created_at time.Time `json:"created_at" validate="required"`
}

func InitDatabase() {
	database.DB.AutoMigrate(User{})
	database.DB.AutoMigrate(Credential{})
	database.DB.AutoMigrate(Order{})
}
