package server

import (
	orderAPIV1 "helloworld-go/api/order/v1"
	"helloworld-go/internal/conf"
	"helloworld-go/internal/order/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(cfg *conf.Bootstrap,
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
	orderAPIV1.RegisterOrderServer(srv, order)
	return srv
}
