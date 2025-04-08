package main

import (
	cntrl "TradeIT/controller"
	database "TradeIT/database"
	"TradeIT/models"
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
	models.InitDatabase(db)
	userService := &services.UserService{}
	userService.SetDB(db)
	userController := &cntrl.UserController{}
	userController.InitUserControllerRoutes(router, *userService)
	router.Run(":8000")

}
