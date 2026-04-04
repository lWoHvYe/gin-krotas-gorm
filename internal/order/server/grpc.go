package server

import (
	"context"
	orderAPIV1 "gin-kratos-gorm/api/order/v1"
	"gin-kratos-gorm/internal/conf"
	"gin-kratos-gorm/internal/order/service"
	"net/url"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(cfg *conf.Bootstrap,
	order *service.OrderService,
	logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			// 1. 注入 JWT 中间件
			//jwt.Server(
			//	func(token *jwtv5.Token) (interface{}, error) {
			//		return []byte("lWoHvYe"), nil // auth.Key 是你的加密私钥
			//	},
			//	// 可选：指定排除不需要验证的路由（例如登录、注册）
			//	jwt.WithSigningMethod(jwtv5.SigningMethodHS256),
			//),
			selector.Server(
				jwt.Server(
					func(token *jwtv5.Token) (interface{}, error) {
						return []byte("lWoHvYe"), nil // auth.Key 是你的加密私钥
					},
					// 可选：指定排除不需要验证的路由（例如登录、注册）
					jwt.WithSigningMethod(jwtv5.SigningMethodHS256),
				),
			).Match(NewWhiteListMatcher()).Build(),
		),
	}
	if cfg.Server.Grpc.Network != "" {
		opts = append(opts, grpc.Network(cfg.Server.Grpc.Network))
	}
	if cfg.Server.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(cfg.Server.Grpc.Addr))
	}
	if cfg.Server.Grpc.Endpoint != "" {
		u, _ := url.Parse(cfg.Server.Grpc.Endpoint)
		opts = append(opts, grpc.Endpoint(u))
	}
	if cfg.Server.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(cfg.Server.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	orderAPIV1.RegisterOrderServer(srv, order)
	return srv
}

// NewWhiteListMatcher 处理白名单（跳过验证）
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/api.auth.v1.Auth/Login"] = struct{}{} // 这里的字符串是 Proto 定义的全名
	whiteList["/api.auth.v1.Auth/Register"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false // 返回 false 表示不执行该中间件
		}
		return true
	}
}
