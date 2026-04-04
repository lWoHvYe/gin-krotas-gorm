package server

import (
	"gin-kratos-gorm/internal/conf"
	httpServerHandler "gin-kratos-gorm/internal/server/http/handler"
	"gin-kratos-gorm/internal/server/http/router"
	"gin-kratos-gorm/internal/server/middleware"
	"gin-kratos-gorm/internal/service"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cfg *conf.Bootstrap,
	userHandler *httpServerHandler.UserHandler,
	productHandler *httpServerHandler.ProductHandler,
	orderHandler *httpServerHandler.OrderHandler,
	enforcer *casbin.Enforcer,
	greeter *service.GreeterService, logger log.Logger) *khttp.Server {

	r := gin.New()
	r.Use(gin.Recovery())

	public := r.Group("/api")
	private := r.Group("/api")
	private.Use(middleware.JWTAuth())

	// 注册Router
	userRouter := router.NewUserRouter(enforcer, userHandler)
	userRouter.InitUserRoutes(private, public)

	productRouter := router.NewProductRouter(enforcer, productHandler)
	productRouter.RegisterProductRouter(private, public)

	orderRouter := router.NewOrderRouter(enforcer, orderHandler)
	orderRouter.RegisterRouter(private, public)

	var opts []khttp.ServerOption
	if cfg.Server.Http.Network != "" {
		opts = append(opts, khttp.Network(cfg.Server.Http.Network))
	}
	if cfg.Server.Http.Addr != "" {
		opts = append(opts, khttp.Address(cfg.Server.Http.Addr))
	}
	if cfg.Server.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(cfg.Server.Http.Timeout.AsDuration()))
	}

	srv := khttp.NewServer(opts...)
	srv.HandlePrefix("/", r) // 这里的 r 是 *gin.Engine
	//v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
