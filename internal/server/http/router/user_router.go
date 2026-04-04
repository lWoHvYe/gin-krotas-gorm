package router

import (
	"gin-kratos-gorm/internal/server/http/handler"
	"gin-kratos-gorm/internal/server/middleware"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userHandler *handler.UserHandler
	enforcer    *casbin.Enforcer
	// ... 其他 API
}

func NewUserRouter(e *casbin.Enforcer, h *handler.UserHandler) *UserRouter {
	return &UserRouter{enforcer: e, userHandler: h}
}

func (r *UserRouter) InitUserRoutes(private *gin.RouterGroup, public *gin.RouterGroup) {
	// 不需要登录的公开接口
	publicGroup := public.Group("base")
	{
		publicGroup.POST("login", r.userHandler.Login)
	}
	// 基础鉴权中间件 (JWT)
	auth := private.Group("/users")
	{
		// 需要登录的私有接口
		auth.GET("getUserInfo", r.userHandler.GetUserInfo)

		// 权限鉴权中间件 (RBAC)
		// 如果还要加 RBAC，就在后面叠加中间件
		rbac := auth.Use(middleware.CasbinMiddleware(r.enforcer))
		rbac.POST("", r.userHandler.Create)
		rbac.GET("/:id", r.userHandler.GetByID)
		rbac.PUT("/:id/role", r.userHandler.UpdateRoleById)
	}
}
