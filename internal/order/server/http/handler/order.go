package handler

import (
	pb "gin-krotas-gorm/api/order/v1"
	"gin-krotas-gorm/internal/order/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req pb.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o, err := h.svc.CreateOrder(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderSn := c.Param("orderSn")
	var req = &pb.GetOrderDetailReq{
		OrderSn: orderSn,
	}
	o, err := h.svc.GetOrderDetail(c, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}
