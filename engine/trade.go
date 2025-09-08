package engine

import (
	"TradeIT/database"
	"TradeIT/models"
	"fmt"
	"log"

	"github.com/svicknesh/tsid"
	"gorm.io/gorm"
)

var trade_db *gorm.DB
var ts *tsid.TSID

func InitTrade() {
	trade_db = database.SetDB()
	ts = tsid.NewDefault()
}

func TestDB() bool {
	return (trade_db == nil)
}

func checkTrade(orderid uint64) bool {
	temp := tradePool.acquireTrade()
	if result := trade_db.Where("bidorderid= ? or askorderid = ?", orderid, orderid).Find(&temp); result.Error != nil {
		log.Fatalf("error in getting the trade details")
		tradePool.releaseTrade(temp)
		return false
	}
	tradePool.releaseTrade(temp)
	return true
}

// registers the trade in trade database
func registerTrades(trade *models.TradeDetails) {
	trade.TradeId = ts.Generate().Number
	fmt.Println(trade)
	fmt.Println("creating trade for ", trade.TradeId)
	if err := trade_db.Create(&trade).Error; err != nil {
		log.Println("Error in creating the entry in trade database: ", err)
	}
}

// updates the order details in main order database
func updateOrderStatusExecuted(order models.Order) {
	order.Status = "executed"
	if err := trade_db.Save(&order).Error; err != nil {
		log.Println("Error while updating the executed order in databse: ", err)
	}
}

func updateOrderStatusUpdated(order models.Order) {
	order.Status = "filled"
	order.Quantity = order.Quantity - order.Remq
	if err := trade_db.Save(&order).Error; err != nil {
		log.Println("error while updating the partially filled order: ")
	}
}

func updateOrderStatusCancelled(order models.Order) {
	order.Status = "Cancelled"
	if err := trade_db.Save(&order).Error; err != nil {
		log.Println("error while updating the partially filled order: ")
	}
}
