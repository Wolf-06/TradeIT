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
	router.POST("/login", u.LoginUser())
	router.PUT("/update/email", u.updateEmail())
	router.PUT("/update/password", u.updatePasswd())
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
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"token": u.userService.LoginUserService(c),
		})
	}
}

func (u *UserController) updateEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"token": u.userService.UpdateUserEmailService(c),
		})
	}
}

func (u *UserController) updatePasswd() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"token": u.userService.UpdateUserPasswdService(c),
		})
	}
}
