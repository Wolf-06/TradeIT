package main

import (
	"TradeIT/internal/database"
	"TradeIT/internal/models"
	"TradeIT/services/user-order-service/controller"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("starting User-Order Service on :8000")

	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error in loading the env file: ", err)
	}

	router := gin.Default()
	database.InitDb()
	models.InitTables()

	userController := controller.InitUserController()
	userController.InitUserControllerRoutes(router)

	orderController := controller.InitOrderController()
	orderController.InitOrderRoutes(router)

	router.Run(":8000")
}
