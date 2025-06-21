package models

import (
	"TradeIT/database"
	"time"
)

type User struct {
	Id    uint64  `gorm:"PrimaryKey"`
	Name  string  `json:"name" validate:"required"`
	Email string  `json:"email" validate:"email,required"`
	Funds float32 `json:"funds"`
}

type Credential struct {
	Id     uint64 `gorm:"PrimaryKey"`
	Email  string `json:"email" validate:"email,  required"`
	Passwd string `json:"passwd" validate:"required, min=6"`
	Token  string
}

type Order struct {
	Id         uint64    `gorm:"PrimaryKey"`
	User_id    int       `json:"user_id" validate:"required"`
	Order_Type string    `json:"orderType validatae:"required"`
	Side       string    `json:"type" validate:"required, oneof= buy sell"`
	Stock      string    `json:"stock" validate:"required"`
	Price      float64   `json:"price" gorm:"type:decimal" validate:"required gt=0"`
	AvgPrice   float64   `json:"avgPrice" gorm:"decimal"`
	Quantity   int       `json:"quantity" validate:"required gt=0"`
	Status     string    `json:"status" validate:"required oneof= executed pending cancelled"`
	Created_at time.Time `json:"created_at" validate:"required"`
}

type Metadata struct {
	Order
	Remq int `json:"rem_quantity" validate:"required gt=0"`
}

func InitDatabase() {
	database.DB.AutoMigrate(User{})
	database.DB.AutoMigrate(Credential{})
	database.DB.AutoMigrate(Order{})
}
