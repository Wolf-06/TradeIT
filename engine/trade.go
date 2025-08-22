package engine

import (
	"TradeIT/database"
	"TradeIT/models"
	"log"

	"github.com/svicknesh/tsid"
)

var db = database.SetDB()
var ts = tsid.NewDefault()

func registerTrades(trade *models.TradeDetails) {
	trade.TradeId = ts.Generate().Number
	if err := db.Create(&trade).Error; err != nil {
		log.Println("Error in creating the entry in trade database: ", err)
	}
}

func updateOrderStatus(order models.Metadata) {
	if err := db.Save(&order).Error; err != nil {
		log.Println("Error while updating the executed order in databse: ", err)
	}
}
