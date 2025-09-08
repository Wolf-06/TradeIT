package main

import (
	cntrl "TradeIT/controller"
	database "TradeIT/database"
	"TradeIT/models"

	test "TradeIT/testing"
	"fmt"
	"log"

	gin "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("starting")

	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error in loading the env file: ", err)
	}

	router := gin.Default()
	database.InitDb()
	models.InitTables()

	userController := cntrl.InitUserController()
	userController.InitUserControllerRoutes(router)

	orderController := cntrl.InitOrderController()
	orderController.InitOrderRoutes(router)
	go test.Test()
	router.Run(":8000")

}
