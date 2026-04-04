package server

import (
	"gin-kratos-gorm/internal/conf"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cfg *conf.Bootstrap,
	logger log.Logger) *khttp.Server {

	r := gin.New()
	r.Use(gin.Recovery())

	// 注册Router

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
