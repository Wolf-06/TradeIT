package controller

import (
	"TradeIT/middleware"
	"TradeIT/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService services.OrderService
}

func InitOrderController() *OrderController {
	return &OrderController{
		orderService: *services.InitOrderService(),
	}
}

func (o *OrderController) InitOrderRoutes(router *gin.Engine) {
	protected := router.Group("/order")
	protected.Use(middleware.VerifyToken())
	{
		protected.GET("/", o.GetAllOrders())
		protected.POST("/sort", o.GetOrders())
	}
}

func (o *OrderController) GetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		json, err := c.Writer.Write(o.orderService.GetAllOrderService(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "error in json",
			})
		} else {
			c.JSON(200, json)
		}
	}
}

func (o *OrderController) GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		json, err := c.Writer.Write(o.orderService.GetOrderByParameterService(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "error in json",
			})
		} else {
			c.JSON(200, json)
		}
	}
}
