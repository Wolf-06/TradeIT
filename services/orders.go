package services

import (
	"TradeIT/database"

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

func (os *OrderService) GetAllOrders() {}
