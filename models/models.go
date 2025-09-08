package models

import (
	"TradeIT/database"
	"time"
)

type User struct {
	Id    uint64  `json:"userid" gorm:"PrimaryKey"`
	Name  string  `json:"name" validate:"required"`
	Email string  `json:"email" validate:"email,required"`
	Funds float32 `json:"funds"`
}

type Credential struct {
	Id     uint64 `json:"userid" gorm:"PrimaryKey"`
	Email  string `json:"email" validate:"email,  required"`
	Passwd string `json:"passwd" validate:"required, min=6"`
	Token  string
}

type MetaOrder struct {
	Id         uint64    `json:"orderid" gorm:"PrimaryKey"`
	User_id    int       `json:"user_id" validate:"required"`
	Order_Type string    `json:"orderType" validate:"required"`
	Side       string    `json:"type" validate:"required, oneof= buy sell"`
	Stock      string    `json:"stock" validate:"required"`
	Price      float64   `json:"price" gorm:"type:decimal" validate:"required gt=0"`
	AvgPrice   float64   `json:"avgPrice" gorm:"type:decimal"`
	Quantity   int       `json:"quantity" validate:"required gt=0"`
	Status     string    `json:"status" validate:"required oneof= executed pending cancelled"`
	Created_at time.Time `json:"created_at" validate:"required"`
}

type Order struct {
	MetaOrder
	Remq int `json:"rem_quantity" validate:"required gt=0"`
}

type TradeDetails struct {
	TradeId     uint64    `json:"tradeid" gorm:"PrimaryKey"`
	Stock       string    `json:"stock"`
	Buyer       int       `json:"buyer"`
	Seller      int       `json:"seller"`
	BidOrderID  uint64    `json:"bidorder_id"`
	AskOrderID  uint64    `json:"askorder_id"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price" gorm:"type:decimal"`
	Executed_at time.Time `json:"execution_time"`
}

func InitTables() {
	database.DB.AutoMigrate(User{})
	database.DB.AutoMigrate(Credential{})
	database.DB.AutoMigrate(Order{})
	database.DB.AutoMigrate(TradeDetails{})
}
