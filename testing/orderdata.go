package testing

import (
	"TradeIT/database"
	"TradeIT/engine"
	"TradeIT/middleware"
	"TradeIT/models"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func roundToTwoDecimal(num float64) float64 {
	return float64((math.Round(num*100) / 100))
}

func OrderTest() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize database connection
	db := database.SetDB()

	stocks := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "META", "NVDA", "NFLX", "IBM", "INTC"}
	statuses := []string{"executed", "pending", "cancelled"}
	userIDs := []int{6794, 2890}
	var i uint64
	for i = 0; i < 20; i++ {
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

		order := models.Order{
			Id:         100 + i,
			User_id:    userID,
			Side:       orderType,
			Stock:      stock,
			Price:      roundToTwoDecimal(price), // Use our rounding function
			Quantity:   quantity,
			Status:     status,
			Created_at: createdAt,
		}

		result := middleware.CreateOrder(db, order)
		if result == "Failed" {
			fmt.Printf("Error creating order")
		} else {
			fmt.Printf("Created order ID: %d, Price: %.2f\n", order.Id, order.Price)
		}
	}
}

func Test() {
	OrderTest()
	engine.EngineTest()
}
