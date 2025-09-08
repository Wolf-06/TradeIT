package testing

import (
	"TradeIT/database"
	"TradeIT/engine"
	"TradeIT/models"
	repo "TradeIT/repo"
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

func createTestOrder(userID int, stock, side string, price float64, quantity int) models.MetaOrder {
	return models.MetaOrder{
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
	order := createTestOrder(1, "AAPL", "buy", 180, 10)
	repo.CreateOrder(db, order)
	err := repo.CancelOrder(db, order.Id)
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
	order := createTestOrder(2, "GOOG", "sell", 2800, 5)
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
	order := createTestOrder(3, "TSLA", "buy", 300, 2)
	repo.CreateOrder(db, order)
	// Simulate execution
	order.Status = "executed"
	db.Save(&order)
	err := repo.CancelOrder(db, order.Id)
	if err == nil {
		t.Errorf("Expected error when cancelling executed order, got nil")
	}
}
