package engine

import (
	"testing"

	"TradeIT/models"

	"github.com/stretchr/testify/assert"
)

func TestMarketOrderMatching__(t *testing.T) {
	ob := setupOrderbook()
	ob.SetLTP(100)
	//test for two market order matching
	ob.Matcher(models.Metadata{Order: models.Order{Id: 1, Side: "sell", Order_Type: "market", Quantity: 10, Status: "pending"}, Remq: 10})
	ob.Matcher(models.Metadata{Order: models.Order{Id: 2, Side: "buy", Order_Type: "market", Quantity: 5, Status: "pending"}, Remq: 5})
	assert.Equal(t, int64(1), ob.tradeCount, "one trade should have occurred")
	assert.True(t, ob.orderTable[1] != nil, "sell order must be present")
	assert.False(t, ob.orderTable[2] != nil, "buy order must be removed after fill")
	//test for partial order clearance with a market order

	ob.Matcher(models.Metadata{Order: models.Order{Id: 3, Side: "buy", Order_Type: "market", Quantity: 10, Status: "pending"}, Remq: 15})
	assert.Equal(t, int64(2), ob.tradeCount, "trade count is not right")
	assert.False(t, ob.orderTable[1] != nil, "partial order should be removed after the fill")
	assert.False(t, ob.orderTable[3] == nil, "partial order should be present")
	assert.Equal(t, ob.orderTable[3].Metadata.Remq, 10, "Remaining quantity is incorrect for the partial order")
	//
	ob.Matcher(models.Metadata{Order: models.Order{Id: 4, Order_Type: "limit", Side: "sell", Price: 100, Quantity: 10, Status: "pending"}, Remq: 10})
	assert.Equal(t, int64(3), ob.tradeCount, "trade count is not right")
	assert.False(t, ob.orderTable[3] != nil, "partial order should be removed after the fill")
	assert.False(t, ob.orderTable[4] != nil, "partial order should be removed after the fill")
}
