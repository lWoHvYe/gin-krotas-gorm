package bootstrap

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewApp(
	logger log.Logger,
	grpcSrv *grpc.Server,
	httpSrv *http.Server,
) *kratos.App {
	return kratos.New(
		kratos.Name("gin-grpc.service"),
		kratos.Logger(logger),
		kratos.Server(
			grpcSrv,
			httpSrv,
		),
	)
}
