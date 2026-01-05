package engine

import (
	"TradeIT/internal/database"
	"TradeIT/internal/models"
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
	var status = true
	result := trade_db.Where("bid_order_id= ? or ask_order_id = ?", orderid, orderid).Find(&temp)
	if result.Error != nil {
		log.Fatalf("error occured in databse query: ", result.Error)
		tradePool.releaseTrade(temp)
		return false
	}
	if result.RowsAffected == 0 {
		status = false
		fmt.Println("No trades Found")
	}
	tradePool.releaseTrade(temp)
	return status
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
