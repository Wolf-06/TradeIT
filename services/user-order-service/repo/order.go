package repo

import (
	"TradeIT/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var orderPool = InitOrderPool()

// redis connector for placing order
var ord_producer = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("redis_addr"),
	Password: os.Getenv("redis_password"),
	DB:       0,
})

// redis connectors for placing order modification request to order engine
var mod_producer = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("redis_addr"),
	Password: os.Getenv("redis_password"),
	DB:       1,
})

var mod_client = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("redis_addr"),
	Password: os.Getenv("redis_password"),
	DB:       1,
})

type Query struct {
	Parameter string `json:"parameter"`
	Value     any    `json:"value"`
}

func CreateOrder(db *gorm.DB, metaOrder models.MetaOrder) bool {
	Order := orderPool.AcquireOrder()
	Order.MetaOrder = metaOrder
	Order.Remq = Order.Quantity
	if err := db.Create(&Order).Error; err != nil {
		fmt.Println("Error in creating the order: ", err)
		return false
	}
	key := "IncomingOrder"
	details, err := json.Marshal(Order)
	if err != nil {
		fmt.Println("Error in marshalling for redis queue: ", err)
		return false
	}

	err = ord_producer.RPush(context.Background(), key, details).Err()
	if err != nil {
		fmt.Println("Error in pushing to redis queue: ", err)
		return false
	}
	orderPool.ReleaseOrder(Order)
	return true
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

func CancelOrder(db *gorm.DB, order_id uint64) error {
	key := "cancel"
	order := orderPool.AcquireOrder()

	if result := db.Where("id=?", order_id).Find(&order); result.Error != nil {
		log.Println("error in fetching data from order database for cancelling order id: ", order_id, "\nError: ", result.Error)
		return errors.New("internal error")
	}
	value := string(order_id) + ":" + order.Stock
	if err := mod_producer.RPush(context.Background(), key, value).Err(); err != nil {
		log.Println("Error in publishing a request to cancel order with id: ", order_id, "\nError: ", err)
		return errors.New("internal error")
	}

	key = key + ":" + string(order_id)
	status, err := mod_client.BLPop(context.Background(), 10, key).Result()
	if err != nil {
		log.Println("Error in getting the cancellation status from the order with ")
	}
	if (status[1] == "nil") || status[1] == "order not reached the orderbook" {
		return nil
	}
	return errors.New(status[1])
}

func ModifyPrice(db *gorm.DB, order_id uint64, newPrice float64) error {
	key := "PriceModifyQueue"
	order := orderPool.AcquireOrder()

	if result := db.Where("id=?", order_id).Find(&order); result.Error != nil {
		log.Println("error in fetching data from order database for cancelling order id: ", order_id, "\nError: ", result.Error)
		return errors.New("internal error")
	}
	id_str := strconv.FormatUint(order_id, 10)
	price_str := fmt.Sprintf("%.2f", newPrice)
	value := id_str + ":" + order.Stock + ":" + price_str
	fmt.Println(value)
	if err := mod_producer.RPush(context.Background(), key, value).Err(); err != nil {
		log.Println("Error in publishing a request to modify price of order with id: ", order_id, "\nError: ", err)
		return errors.New("internal error")
	}

	key = key + ":" + id_str
	status, err := mod_client.BLPop(context.Background(), 0, key).Result()
	if err != nil {
		log.Println("Error in getting the cancellation status from the order")
	}
	if (status[1] == "nil") || status[1] == "order not reached the orderbook" {
		return nil
	}
	return errors.New(status[1])
}
