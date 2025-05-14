package middleware

import (
	"TradeIT/models"
	"log"

	"gorm.io/gorm"
)

func FundsInfo(db *gorm.DB, userid float64) float32 {
	var detail models.User
	if err := db.Where("id = ?", userid).Find(&detail).Error; err != nil {
		log.Fatalln("Error in getting the details: ", err)
		return -1
	}
	return detail.Funds
}
