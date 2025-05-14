package services

import (
	"TradeIT/database"
	"TradeIT/middleware"
	"TradeIT/models"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func InitOrderService() *OrderService {
	return &OrderService{
		db: database.SetDB(),
	}
}

func (os *OrderService) CreateOrderService(c *gin.Context) string {
	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		log.Fatalln("Json Binding error: ", err)
		return "Error"
	}
	return middleware.CreateOrder(os.db, order)

}

func (os *OrderService) GetAllOrderService(c *gin.Context) []byte {
	user_id, _ := c.Get("userid")
	json := middleware.GetAllOrders(os.db, user_id.(float64))
	return json
}

func (os *OrderService) GetOrderByParameterService(c *gin.Context) []byte {
	var Query middleware.Query
	user_id, _ := c.Get("userid")
	if err := c.BindJSON(&Query); err != nil {
		fmt.Println("Error in binding JSON: ", err)
	}
	json := middleware.GetOrders(os.db, user_id.(float64), Query)
	return json
}
