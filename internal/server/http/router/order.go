package router

import (
	"gin-krotas-gorm/internal/server/http/handler"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

type OrderRouter struct {
	orderHandler *handler.OrderHandler
	enforcer     *casbin.Enforcer
}

func NewOrderRouter(e *casbin.Enforcer, orderHandler *handler.OrderHandler) *OrderRouter {
	return &OrderRouter{enforcer: e, orderHandler: orderHandler}
}

func (router *OrderRouter) RegisterRouter(private *gin.RouterGroup, public *gin.RouterGroup) {
	orderGroup := private.Group("order")
	{
		orderGroup.POST("/create", router.orderHandler.CreateOrder)
		orderGroup.GET("/:orderSn", router.orderHandler.GetOrder)
	}
}
