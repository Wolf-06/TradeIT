package engine

import (
	"TradeIT/models"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupOrderbook initializes an Orderbook with sample buy and sell orders.
func SetupOrderbook() *Orderbook {
	ob := InitOrderBook_()
	//ob.UnLock()
	// Insert a buy order with ID 1 at price 50 for quantity 10.
	ob.InsertOrder(models.Metadata{
		Order: models.Order{
			Id:       1,
			Side:     "buy",
			Price:    50,
			Quantity: 10,
			Status:   "pending",
		},
		Remq: 10,
	})
	ob.buyCount++

	// Insert a sell order with ID 2 at price 75 for quantity 5.
	ob.InsertOrder(models.Metadata{
		Order: models.Order{
			Id:       2,
			Side:     "sell",
			Price:    75,
			Quantity: 5,
			Status:   "pending",
		},
		Remq: 5,
	})
	ob.sellCount++
	fmt.Println("lock")
	return ob
}

func TestCancelExistingBuyOrder(t *testing.T) {
	ob := SetupOrderbook()
	//ob.Unlock()
	err := ob.CancelOrder(1)
	assert.NoError(t, err, "cancelling existing buy order should not error")
	_, exists := ob.orderTable[1]
	assert.False(t, exists, "orderTable must not contain cancelled buy order")
	_, priceExists := ob.buy_orders[50]
	assert.False(t, priceExists, "buy_orders map should remove price level when empty")
}

func TestCancelExistingSellOrder(t *testing.T) {
	ob := SetupOrderbook()
	//ob.Unlock()
	err := ob.CancelOrder(2)
	assert.NoError(t, err, "cancelling existing sell order should not error")
	_, exists := ob.orderTable[2]
	assert.False(t, exists, "orderTable must not contain cancelled sell order")
	_, priceExists := ob.sell_orders[75]
	assert.False(t, priceExists, "sell_orders map should remove price level when empty")
}

func TestCancelNonExistentOrder(t *testing.T) {
	ob := SetupOrderbook()
	//ob.Unlock()
	err := ob.CancelOrder(999)
	expected := errors.New("order has been proceessed or does not exist")
	assert.EqualError(t, err, expected.Error(), "cancelling non-existent order should return specific error")
}
