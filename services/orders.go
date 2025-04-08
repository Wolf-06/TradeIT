package services

import (
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

//func (o *OrderService) PlaceOrder(c *gin.Context) gin.HandlerFunc {
//	var order models.Order
//	if err := C.BindJson(&order); err != nil {
//		log.Fatal("Error in Binding json to order ", err)
//	}
//	err := o.db.Create(&models.Order{})
//	if err != nil {
//		log.Fatal(err)
//	}
//}
