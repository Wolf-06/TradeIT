package main

import (
	"TradeIT/engine"
	"TradeIT/models"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func roundToTwoDecimal(num float64) float64 {
	return float64((math.Round(num*10) / 10))
}

func generateOrders(n uint64) []models.Metadata {
	orders := make([]models.Metadata, 0, n)
	var i uint64
	var limitCount, marketCount int
	var sellOrders, buyOrders int
	for i = 0; i < n; i++ {
		orderSide := "buy"
		orderType := "limit"
		if rand.Float64() < 0.4 {
			orderSide = "sell"
			sellOrders++
		} else {
			buyOrders++
		}
		if rand.Float32() > 0.6 {
			orderType = "market"
			marketCount++
		} else {
			limitCount++
		}
		price := 90 + rand.Float64()*10 // Prices between 90 and 110
		quantity := rand.Intn(100) + 1  // Quantity between 1 and 100
		order := models.Metadata{
			Order: models.Order{
				Id:         i,
				User_id:    rand.Intn(1000),
				Order_Type: orderType,
				Side:       orderSide,
				Stock:      "TEST",
				Price:      roundToTwoDecimal(price),
				Quantity:   quantity,
				Status:     "pending",
				Created_at: time.Now(),
			},
			Remq: quantity,
		}
		orders = append(orders, order)
	}
	fmt.Println("Limit orders: ", limitCount, " | Market Orders ", marketCount)
	fmt.Println("Buy orders: ", buyOrders, " | Sell orders: ", sellOrders)
	return orders
}

func main_() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	ob := engine.InitOrderBook_()
	orders := generateOrders(1000000) //generate orders
	fmt.Println("Starting the Matching")
	start := time.Now()
	for _, order := range orders {
		ob.Matcher(order)
	}
	elapsed := time.Since(start)
	fmt.Printf("Processed %d orders in %s\n", len(orders), elapsed)
	ob.DisplayResult()
}
