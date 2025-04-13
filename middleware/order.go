package middleware

import (
	"TradeIT/models"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type Query struct {
	Parameter string `json:"parameter"`
	Value     any    `json:"value"`
}

func CreateOrder(db *gorm.DB, order models.Order) string {
	if err := db.Create(&order).Error; err != nil {
		fmt.Println("Error in creating the order: ", err)
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
