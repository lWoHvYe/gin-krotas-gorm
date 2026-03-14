package server

import (
	"gin-krotas-gorm/internal/conf"
	httpServerHandler "gin-krotas-gorm/internal/order/server/http/handler"
	"gin-krotas-gorm/internal/order/server/http/router"
	"gin-krotas-gorm/internal/order/server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cfg *conf.Bootstrap,
	orderHandler *httpServerHandler.OrderHandler,
	logger log.Logger) *khttp.Server {

	r := gin.New()
	r.Use(gin.Recovery())

	public := r.Group("/api")
	private := r.Group("/api")
	private.Use(middleware.JWTAuth())

	// 注册Router
	orderRouter := router.NewOrderRouter(orderHandler)
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
