package main

import (
	cntrl "TradeIT/controller"
	database "TradeIT/database"
	"TradeIT/models"
	"fmt"

	gin "github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("starting")

	router := gin.Default()
	database.InitDb()
	models.InitDatabase()

	userController := cntrl.InitUserController()
	userController.InitUserControllerRoutes(router)

	orderController := cntrl.InitOrderController()
	orderController.InitOrderRoutes(router)
	router.Run(":8000")
}
