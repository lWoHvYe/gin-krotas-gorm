package router

import (
	"gin-krotas-gorm/internal/order/server/http/handler"

	"github.com/gin-gonic/gin"
)

type OrderRouter struct {
	orderHandler *handler.OrderHandler
}

func NewOrderRouter(orderHandler *handler.OrderHandler) *OrderRouter {
	return &OrderRouter{orderHandler: orderHandler}
}

func (router *OrderRouter) RegisterRouter(private *gin.RouterGroup, public *gin.RouterGroup) {
	orderGroup := private.Group("order")
	{
		orderGroup.POST("/create", router.orderHandler.CreateOrder)
		orderGroup.GET("/:orderSn", router.orderHandler.GetOrder)
	}
}
