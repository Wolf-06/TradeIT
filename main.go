package main

import (
	cntrl "TradeIT/controller"
	database "TradeIT/database"
	"TradeIT/models"
	"TradeIT/testing"
	"fmt"

	gin "github.com/gin-gonic/gin"
)

func main_() {
	fmt.Println("starting")

	router := gin.Default()
	database.InitDb()
	models.InitDatabase()

	userController := cntrl.InitUserController()
	userController.InitUserControllerRoutes(router)

	orderController := cntrl.InitOrderController()
	orderController.InitOrderRoutes(router)
	go testing.Test()
	router.Run(":8000")
}
