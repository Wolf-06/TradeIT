package main

import (
	cntrl "TradeIT/controller"
	database "TradeIT/database"
	"TradeIT/services"
	"fmt"

	gin "github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("starting")
	router := gin.Default()
	db := database.InitDb()
	if db == nil {
		fmt.Print("Database error: \n")
	}
	userService := &services.UserService{}
	userService.InitService(db)
	userController := &cntrl.UserController{}
	userController.InitUserControllerRoutes(router, *userService)
	router.Run(":8000")

}
