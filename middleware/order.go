package middleware

import (
	"TradeIT/models"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Query struct { //used for custom searching the orders
	Parameter string `json:"parameter"`
	Value     any    `json:"value"`
}

var producer = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("redis_addr"),
	Password: os.Getenv("redis_password"),
	DB:       0,
})

func CreateOrder(db *gorm.DB, order models.Order) string {
	if err := db.Create(&order).Error; err != nil {
		fmt.Println("Error in creating the order: ", err)
		return "Failed"
	}
	key := "order"
	meta := models.Metadata{Order: order, Remq: order.Quantity}
	details, err := json.Marshal(meta)
	if err != nil {
		fmt.Println("Error in marshalling for redis queue: ", err)
		return "Failed"
	}
	err = producer.LPush(context.Background(), key, details).Err()
	if err != nil {
		fmt.Println("Error in pushing to redis queue: ", err)
		return "Failed"
	}
	return "Success"
}

func GetAllOrders(db *gorm.DB, user_id float64) []byte {
	var orders []models.Order
	result := db.Where("user_id = ?", user_id).Find(&orders)
	if result.Error != nil {
		fmt.Println("Error in getting the orders from database: ", result.Error)
	}
	order_json, err := json.MarshalIndent(orders, "", "")
	if err != nil {
		fmt.Println("Error in Json Marshalling: ", err)
	}
	return order_json
}

func GetOrders(db *gorm.DB, user_id float64, query Query) []byte {
	var orders []models.Order
	fmt.Println(query)
	cond := "user_id = ? AND " + query.Parameter + " =?"
	fmt.Println(cond)
	result := db.Where(cond, user_id, query.Value.(string)).Find(&orders)
	if result.Error != nil {
		fmt.Println("Error in fetching the data: ", result.Error)
		return nil
	}

	orderJson, err := json.MarshalIndent(orders, "", "")
	if err == nil {
		fmt.Println("Error in Marshalling json")
	}
	return orderJson
}
