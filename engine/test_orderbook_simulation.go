package engine

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"TradeIT/models"

	"github.com/stretchr/testify/assert"
)

// setupOrderbook initializes an empty Orderbook with no orders.
func setupOrderbook() *Orderbook {
	ob := InitOrderBook_()
	return ob
}

// TestNormalMatch tests basic buy–sell matching and full fills.
func TestNormalMatch(t *testing.T) {
	ob := setupOrderbook()
	ob.Lock()
	// Insert a sell order at price 100, quantity 10
	ob.InsertOrder(models.Metadata{Order: models.Order{Id: 1, Side: "sell", Price: 100, Quantity: 10, Status: "pending"}, Remq: 10})
	fmt.Println("Inserted 1")
	ob.sellCount++
	// Insert a buy order at price 110, quantity 10
	ob.InsertOrder(models.Metadata{Order: models.Order{Id: 2, Side: "buy", Price: 110, Quantity: 10, Status: "pending"}, Remq: 10})
	fmt.Println("Inserted 2")
	ob.buyCount++
	ob.Unlock()
	ob.Matcher(models.Metadata{Order: models.Order{Id: 2, Side: "buy", Price: 110, Quantity: 10, Status: "pending"}, Remq: 10})
	assert.Equal(t, int64(1), ob.tradeCount, "one trade should have occurred")
	assert.False(t, ob.orderTable[1] != nil, "sell order must be removed after fill")
	assert.False(t, ob.orderTable[2] != nil, "buy order must be removed after fill")
}

// TestPartialFill tests partial fills and remaining quantities.
func TestPartialFill(t *testing.T) {
	ob := setupOrderbook()
	ob.Lock()
	// Insert a sell order with quantity 5
	ob.InsertOrder(models.Metadata{Order: models.Order{Id: 3, Side: "sell", Price: 200, Quantity: 5, Status: "pending"}, Remq: 5})
	ob.sellCount++
	// Insert a buy order with quantity 10
	ob.InsertOrder(models.Metadata{Order: models.Order{Id: 4, Side: "buy", Price: 200, Quantity: 10, Status: "pending"}, Remq: 10})
	ob.buyCount++
	ob.Unlock()
	// Matcher should partially fill the buy order
	ob.Matcher(models.Metadata{Order: models.Order{Id: 4, Side: "buy", Price: 200, Quantity: 10, Status: "pending"}, Remq: 10})
	assert.Equal(t, int64(1), ob.tradeCount, "one partial trade should occur")
	node, exists := ob.orderTable[4]
	assert.True(t, exists, "buy order should remain after partial fill")
	assert.Equal(t, 5, node.Metadata.Remq, "remaining quantity should be 5")
}

// TestCancelBeforeMatch tests cancellation of orders before matching.
func TestCancelBeforeMatch(t *testing.T) {
	ob := setupOrderbook()
	ob.Lock()
	ob.InsertOrder(models.Metadata{Order: models.Order{Id: 5, Side: "buy", Price: 150, Quantity: 5, Status: "pending"}, Remq: 5})
	ob.buyCount++
	ob.Unlock()
	err := ob.CancelOrder(5)
	assert.NoError(t, err, "cancellation before match should not error")
	_, exists := ob.orderTable[5]
	assert.False(t, exists, "orderTable should not contain cancelled order")
	_, levelExists := ob.buy_orders[150]
	assert.False(t, levelExists, "price level must be removed when empty")
}

// TestCancelAfterPartial tests cancellation after a partial fill.
func TestCancelAfterPartial(t *testing.T) {
	ob := setupOrderbook()
	ob.Lock()
	ob.InsertOrder(models.Metadata{Order: models.Order{Id: 6, Side: "sell", Price: 250, Quantity: 20, Status: "pending"}, Remq: 20})
	ob.sellCount++
	ob.Unlock()
	// Trigger partial match of 10 units
	ob.Matcher(models.Metadata{Order: models.Order{Id: 7, Side: "buy", Price: 250, Quantity: 10, Status: "pending"}, Remq: 10})
	// Now cancel remaining sell order
	err := ob.CancelOrder(6)
	assert.NoError(t, err, "cancellation after partial fill should succeed")
	_, exists := ob.orderTable[6]
	assert.False(t, exists, "partially filled order should be removed on cancel")
}

// TestConcurrentCancellations tests thread safety with concurrent cancellations.
func TestConcurrentCancellations(t *testing.T) {
	ob := setupOrderbook()
	ob.Lock()
	// Insert many orders
	for i := uint64(100); i < 110; i++ {
		ob.InsertOrder(models.Metadata{Order: models.Order{Id: i, Side: "buy", Price: 300, Quantity: 1, Status: "pending"}, Remq: 1})
		ob.buyCount++
	}
	ob.Unlock()
	var wg sync.WaitGroup
	for i := uint64(100); i < 110; i++ {
		wg.Add(1)
		go func(id uint64) {
			defer wg.Done()
			err := ob.CancelOrder(id)
			assert.NoError(t, err, "concurrent cancellation should not error")
		}(i)
	}
	wg.Wait()

	for i := uint64(100); i < 110; i++ {
		_, exists := ob.orderTable[i]
		assert.False(t, exists, "all orders must be removed concurrently")
	}
}

// TestCancelNonExistent tests error on cancelling non-existent order.
func TestCancelNonExistent(t *testing.T) {
	ob := setupOrderbook()
	err := ob.CancelOrder(9999)
	assert.EqualError(t, err, "order has been proceessed or does not exist", "should return specific error for invalid ID")
}

// TestInvalidOrderType tests unsupported order type handling.
func TestInvalidOrderType(t *testing.T) {
	ob := setupOrderbook()
	// Manually inject invalid Metadata to simulate bad data
	node := &Node{Metadata: models.Metadata{Order: models.Order{Id: 10, Side: "hold", Price: 100, Quantity: 1, Status: "pending"}, Remq: 1}}
	ob.orderTable[10] = node

	err := ob.CancelOrder(10)
	assert.EqualError(t, err, "unsupported order type", "should error on invalid order side")
}

// TestStressScenario simulates rapid inserts, matches, and cancels under load.
func TestStressScenario(t *testing.T) {
	ob := setupOrderbook()
	for i := 1; i <= 10; i++ {
		side := "buy"
		if i%2 == 0 {
			side = "sell"
		}
		ob.Lock()
		ob.InsertOrder(models.Metadata{Order: models.Order{Id: uint64(1000 + i), Side: side, Price: float64(100 + i%5), Quantity: 1, Status: "pending"}, Remq: 1})
		ob.Unlock()
		if i%10 == 0 {
			ob.Matcher(models.Metadata{Order: models.Order{Id: uint64(2000 + i), Side: side, Price: float64(100 + i%5), Quantity: 1, Status: "pending"}, Remq: 1})
		}
		if i%15 == 0 {
			_ = ob.CancelOrder(uint64(1000 + i))
		}
	}
	// Allow matcher operations to complete
	time.Sleep(100 * time.Millisecond)
	// Verify that no panics occurred and internal counts are non-negative
	ob.DisplayResult()
	assert.True(t, ob.buyCount >= 0, "buyCount must be non-negative")
	assert.True(t, ob.sellCount >= 0, "sellCount must be non-negative")
}
