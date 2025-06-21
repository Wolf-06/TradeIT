package engine

import (
	"fmt"
	"testing"

	"TradeIT/models"

	"github.com/stretchr/testify/assert"
)

func TestNormalQuantityModification(t *testing.T) {
	ob := setupOrderbook()
	ob.Lock()
	ob.InsertOrder(models.Metadata{Order: models.Order{Id: 1, Side: "sell", Price: 100, Quantity: 5, Status: "pending"}, Remq: 5})
	fmt.Println("Inserted 1")
	// Insert a buy order at price 110, quantity 10
	//ob.InsertOrder(models.Metadata{Order: models.Order{Id: 2, Side: "buy", Price: 110, Quantity: 10, Status: "pending"}, Remq: 10})
	ob.Unlock()
	// Matching engine should fully match both orders
	ob.Matcher(models.Metadata{Order: models.Order{Id: 2, Side: "buy", Price: 110, Quantity: 10, Status: "pending"}, Remq: 10})
	err := ob.ModifyQuantity(2, -5)
	assert.EqualError(t, err, "invalid quantity (negative/zero quantity)")
	err = ob.ModifyQuantity(2, -5)
	assert.EqualError(t, err, "invalid quantity (negative/zero quantity)")
	err = ob.ModifyQuantity(2, 1)
	assert.EqualError(t, err, "order has been partially filled and quantity can't be decreased ")
	err = ob.ModifyQuantity(2, 15)
	assert.NoError(t, err, "Updation is not working properly")
	assert.False(t, ob.orderTable[2] == nil, "order has been removed")
	assert.True(t, ob.orderTable[2].Metadata.Quantity == 15, "Quantity didn't modify properly")
	assert.True(t, ob.orderTable[2].Metadata.Remq == 10, "Remq didn't modify properly")
	ob.Matcher(models.Metadata{Order: models.Order{Id: 1, Side: "sell", Price: 110, Quantity: 10, Status: "pending"}, Remq: 10})
	assert.Equal(t, int64(2), ob.tradeCount, "one trade should have occurred")
	assert.False(t, ob.orderTable[1] != nil, "sell order must be removed after fill")
	assert.False(t, ob.orderTable[2] != nil, "buy order must be removed after fill")
	ob.DisplayResult()
}

func TestNormalPriceModification(t *testing.T) {
	ob := setupOrderbook()
	ob.Lock()
	ob.InsertOrder(models.Metadata{Order: models.Order{Id: 1, Side: "sell", Price: 100, Quantity: 5, Status: "pending"}, Remq: 5})
	ob.Unlock()
	ob.Matcher(models.Metadata{Order: models.Order{Id: 2, Side: "buy", Price: 90, Quantity: 5, Status: "pending"}, Remq: 5})
	err := ob.ModifyPrice(2, -90)
	assert.EqualError(t, err, "invalid price (negative or zero price)")
	err = ob.ModifyPrice(2, 0)
	assert.EqualError(t, err, "invalid price (negative or zero price)")
	err = ob.ModifyPrice(2, 95)
	assert.NoError(t, err, "price updation logic has issues")
	err = ob.ModifyPrice(2, 96)
	assert.NoError(t, err, "price updation logic has issues")
	err = ob.ModifyPrice(1, 93)
	assert.NoError(t, err, "price updation logic has issues")
	assert.Equal(t, int64(1), ob.tradeCount, "one trade should have occurred")
	assert.False(t, ob.orderTable[1] != nil, "sell order must be removed after fill")
	assert.False(t, ob.orderTable[2] != nil, "buy order must be removed after fill")
	err = ob.ModifyPrice(2, 0)
	assert.EqualError(t, err, "order has been processed or order doesn't exists")
}

func TestBothPriceAndQuantityModification(t *testing.T) {
	ob := setupOrderbook()
	ob.Lock()
	ob.InsertOrder(models.Metadata{Order: models.Order{Id: 1, Side: "sell", Price: 100, Quantity: 15, Status: "pending"}, Remq: 10})
	ob.Unlock()
	ob.Matcher(models.Metadata{Order: models.Order{Id: 2, Side: "buy", Price: 90, Quantity: 10, Status: "pending"}, Remq: 10})
	err := ob.ModifyPrice(2, 95)
	assert.NoError(t, err, "price updation logic has issues")
	err = ob.ModifyQuantity(1, 10)
	assert.NoError(t, err, "quantity updation logic has issues")
	err = ob.ModifyPrice(1, 93)
	assert.NoError(t, err, "price updation logic has issues")
	assert.Equal(t, int64(1), ob.tradeCount, "one trade should have occurred")
	assert.False(t, ob.orderTable[1] != nil, "sell order must be removed after fill")
	assert.False(t, ob.orderTable[2] == nil, "buy order must be removed after fill")
	fmt.Println(len(ob.sell_orders))
	ob.Matcher(models.Metadata{Order: models.Order{Id: 3, Side: "buy", Price: 110, Quantity: 10, Status: "pending"}, Remq: 10})
	ob.Matcher(models.Metadata{Order: models.Order{Id: 4, Side: "sell", Price: 80, Quantity: 100, Status: "pending"}, Remq: 100})
	ob.DisplayResult()

}
