package server

import (
	"helloworld-go/internal/conf"
	httpServerHandler "helloworld-go/internal/server/http/handler"
	"helloworld-go/internal/server/http/router"
	"helloworld-go/internal/server/middleware"
	"helloworld-go/internal/service"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cfg *conf.Bootstrap,
	userHandler *httpServerHandler.UserHandler,
	enforcer *casbin.Enforcer,
	greeter *service.GreeterService, logger log.Logger) *khttp.Server {

	r := gin.New()
	r.Use(gin.Recovery())

	public := r.Group("/api")
	private := r.Group("/api")
	private.Use(middleware.JWTAuth())

	// 注册Router
	userRouter := router.NewUserRouter(enforcer, userHandler)
	userRouter.RegisterUserRoutes(private, public)

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
