package router

import (
	"helloworld-go/internal/server/http/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, h *handler.UserHandler) {
	grp := r.Group("/users")
	grp.POST("", h.Create)
	grp.GET("/:id", h.GetByID)
	grp.PUT("/:id/role", h.UpdateRoleById)
}
