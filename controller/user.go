package controller

import (
	"TradeIT/middleware"
	"TradeIT/services"

	gin "github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func InitUserController() *UserController {
	return &UserController{
		userService: *services.InitUserService(),
	}
}

func (u *UserController) InitUserControllerRoutes(router *gin.Engine) {
	router.POST("/register", u.RegisterUser())
	router.POST("/login", u.LoginUser())

	protected := router.Group("/user")
	protected.Use(middleware.VerifyToken())
	{
		protected.PUT("/update/email", u.updateEmail())
		protected.PUT("/update/password", u.updatePasswd())
		protected.GET("/fund", u.getFunds())
	}
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

func (u *UserController) getFunds() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"funds": u.userService.GetFundsService(c)})
	}
}
