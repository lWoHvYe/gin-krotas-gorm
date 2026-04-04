package server

import (
	HelloWorldV1 "gin-kratos-gorm/api/helloworld/v1"
	orderAPIV1 "gin-kratos-gorm/api/order/v1"
	productV1 "gin-kratos-gorm/api/product/v1"
	UserAPIV1 "gin-kratos-gorm/api/user/v1"
	"gin-kratos-gorm/internal/conf"
	"gin-kratos-gorm/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(cfg *conf.Bootstrap,
	greeter *service.GreeterService,
	user *service.UserService,
	product *service.ProductService,
	order *service.OrderService,
	logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if cfg.Server.Grpc.Network != "" {
		opts = append(opts, grpc.Network(cfg.Server.Grpc.Network))
	}
	if cfg.Server.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(cfg.Server.Grpc.Addr))
	}
	if cfg.Server.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(cfg.Server.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	HelloWorldV1.RegisterGreeterServer(srv, greeter)
	UserAPIV1.RegisterUserServer(srv, user)
	productV1.RegisterProductServer(srv, product)
	orderAPIV1.RegisterOrderServer(srv, order)
	return srv
}
