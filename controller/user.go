package controller

import (
	"TradeIT/services"

	gin "github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func (u *UserController) InitUserControllerRoutes(router *gin.Engine, initialisedUserService services.UserService) {
	router.POST("/register", u.RegisterUser())
	router.Post("/login", u.LoginUser())
	u.userService = initialisedUserService
}

func (u *UserController) RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"user id": u.userService.RegisterUserService(c),
		})
	}
}

func (u *UserController) LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "success",
		})
	}
}
