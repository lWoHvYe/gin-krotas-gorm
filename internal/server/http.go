package server

import (
	"helloworld-go/internal/conf"
	httpServerHandler "helloworld-go/internal/server/http/handler"
	"helloworld-go/internal/server/http/router"
	"helloworld-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cfg *conf.Bootstrap,
	userHandler *httpServerHandler.UserHandler,
	greeter *service.GreeterService, logger log.Logger) *khttp.Server {

	r := gin.New()
	r.Use(gin.Recovery())
	// 注册router
	router.RegisterUserRoutes(r, userHandler)

	var opts = []khttp.ServerOption{
		//khttp.Handler(r),
		// 手动实现 Option 逻辑，这等同于调用 http.Handler(g)
		func(s *khttp.Server) {
			s.HandlePrefix("/", r) // Kratos 新版底层通过 HandlePrefix 挂载
		},
	}
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
	//v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
