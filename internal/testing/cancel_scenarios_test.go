package testing

import (
	"TradeIT/engine"
	"TradeIT/internal/database"
	"TradeIT/internal/models"
	repo "TradeIT/services/user-order-service/repo"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

var testOrderEngine *engine.OrderEngine

func setupTestDBAndEngine() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalln("Error in loading the env file: ", err)
	}
	database.InitDb()
	testOrderEngine = engine.InitOrderEngine()
}

func createTestOrder(orderId uint64, userID int, stock, side string, price float64, quantity int) models.MetaOrder {
	return models.MetaOrder{
<<<<<<< HEAD:internal/testing/cancel_scenarios_test.go
		Id:         orderId,
=======
>>>>>>> parent of 6daa9c4 (Completed the Order Creation and Matching Pipeline.):testing/cancel_scenarios_test.go
		User_id:    userID,
		Order_Type: "limit",
		Side:       side,
		Stock:      stock,
		Price:      price,
		Quantity:   quantity,
		Status:     "pending",
		Created_at: time.Now(),
	}
}

func TestCancelExistingOrder(t *testing.T) {
	setupTestDBAndEngine()
	db := database.SetDB()
	order := createTestOrder(1001, 1, "AAPL", "buy", 180, 10)
	repo.CreateOrder(db, order)
	err := repo.CancelOrder(db, order.Id)
	order = createTestOrder(1002, 1, "XXXX", "buy", 0, 0)
	repo.CreateOrder(db, order)
	if err != nil {
		t.Errorf("Expected no error when cancelling existing order, got: %v", err)
	}
}

func TestCancelNonExistentOrder(t *testing.T) {
	setupTestDBAndEngine()
	db := database.SetDB()
	nonExistentOrderID := uint64(999999)
	err := repo.CancelOrder(db, nonExistentOrderID)
	if err == nil {
		t.Errorf("Expected error when cancelling non-existent order, got nil")
	}
}

func TestCancelAlreadyCancelledOrder(t *testing.T) {
	setupTestDBAndEngine()
	db := database.SetDB()
	order := createTestOrder(1003, 2, "GOOG", "sell", 2800, 5)
	repo.CreateOrder(db, order)
	repo.CancelOrder(db, order.Id)
	err := repo.CancelOrder(db, order.Id)
	if err == nil {
		t.Errorf("Expected error when cancelling already cancelled order, got nil")
	}
}

func TestCancelExecutedOrder(t *testing.T) {
	setupTestDBAndEngine()
	db := database.SetDB()
	order := createTestOrder(1004, 3, "TSLA", "sell", 300, 2)
	repo.CreateOrder(db, order)
	order = createTestOrder(1005, 4, "TSLA", "buy", 300, 10)
	repo.CreateOrder(db, order)
	err := repo.CancelOrder(db, order.Id)
	if err == nil {
		t.Errorf("Expected error when cancelling executed order, got nil")
	}
}
