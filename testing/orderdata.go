package testing

import (
	"TradeIT/database"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Order struct {
	Id         int       `gorm:"PrimaryKey"`
	User_id    int       `json:"user_id" validate:"required"`
	Order_type string    `json:"type" validate:"required,oneof=buy sell"`
	Stock      string    `json:"stock" validate:"required"`
	Price      float32   `json:"price" validate:"required,gt=0"`
	Quantity   int       `json:"quantity" validate:"required,gt=0"`
	Status     string    `json:"status" validate:"required,oneof=executed pending cancelled"`
	Created_at time.Time `json:"created_at" validate:"required"`
}

func roundToTwoDecimal(num float64) float32 {
	return float32(math.Round(num*100) / 100)
}

func Test() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize database connection
	db := database.SetDB()

	stocks := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "META", "NVDA", "NFLX", "IBM", "INTC"}
	statuses := []string{"executed", "pending", "cancelled"}
	userIDs := []int{6794, 2890}

	for i := 0; i < 20; i++ {
		orderType := "buy"
		if rand.Intn(2) == 1 {
			orderType = "sell"
		}

		stock := stocks[rand.Intn(len(stocks))]
		price := 50 + rand.Float64()*500 // Random price between 50 and 550 as float64
		quantity := 1 + rand.Intn(100)   // Random quantity between 1 and 100
		status := statuses[rand.Intn(len(statuses))]
		userID := userIDs[rand.Intn(len(userIDs))]

		// Create order with random time in the past 30 days
		createdAt := time.Now().Add(-time.Duration(rand.Intn(30*24)) * time.Hour)

		order := Order{
			User_id:    userID,
			Order_type: orderType,
			Stock:      stock,
			Price:      roundToTwoDecimal(price), // Use our rounding function
			Quantity:   quantity,
			Status:     status,
			Created_at: createdAt,
		}

		result := db.Create(&order)
		if result.Error != nil {
			fmt.Printf("Error creating order: %v\n", result.Error)
		} else {
			fmt.Printf("Created order ID: %d, Price: %.2f\n", order.Id, order.Price)
		}
	}
}
