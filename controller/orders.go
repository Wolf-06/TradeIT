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

func (oc *OrderController) InitOrderRoutes(router *gin.Engine) {
	protected := router.Group("/order")
	protected.Use(middleware.VerifyToken())
	{
		protected.GET("/", oc.GetAllOrders())
		protected.POST("/sort", oc.GetOrders())
		protected.POST("/create", oc.PlaceOrder())
	}
}

func (oc *OrderController) GetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		json, err := c.Writer.Write(oc.orderService.GetAllOrderService(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "error in json",
			})
		} else {
			c.JSON(200, json)
		}
	}
}

func (oc *OrderController) GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		json, err := c.Writer.Write(oc.orderService.GetOrderByParameterService(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "error in json",
			})
		} else {
			c.JSON(200, json)
		}
	}
}

func (oc *OrderController) PlaceOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status := oc.orderService.CreateOrderService(ctx)
		if !status {
			ctx.JSON(500, "Failed due to internal Server Error")
		} else {
			ctx.JSON(201, "Success")
		}
	}
}
