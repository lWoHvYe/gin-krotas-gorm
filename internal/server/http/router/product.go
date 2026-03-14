package router

import (
	"gin-krotas-gorm/internal/server/http/handler"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

type ProductRouter struct {
	productHandler *handler.ProductHandler
	enforcer       *casbin.Enforcer
}

func NewProductRouter(e *casbin.Enforcer, productHandler *handler.ProductHandler) *ProductRouter {
	return &ProductRouter{enforcer: e, productHandler: productHandler}
}

func (r *ProductRouter) RegisterProductRouter(private *gin.RouterGroup, public *gin.RouterGroup) {
	productGroup := public.Group("product")
	{
		// 商品详情
		productGroup.GET("/:id", r.productHandler.GetProductDetail)
	}
}
