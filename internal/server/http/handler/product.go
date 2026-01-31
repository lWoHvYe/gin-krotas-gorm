package handler

import (
	pb "helloworld-go/api/product/v1"
	"helloworld-go/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (r *ProductHandler) GetProductDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req = &pb.GetProductDetailReq{Id: uint32(id)}
	// 调用 Service 层（Kratos 风格的 Service）
	reply, err := r.svc.GetProductDetail(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, reply)
}
