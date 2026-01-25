//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package bootstrap

import (
	"helloworld-go/internal/biz"
	"helloworld-go/internal/data"
	"helloworld-go/internal/pkg"
	"helloworld-go/internal/server"
	"helloworld-go/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func InitApp(string, log.Logger) (*kratos.App, func(), error) {
	wire.Build(
		biz.BizSet,
		pkg.InfraSet,
		data.RepoSet,
		service.ServiceSet,
		server.TransportSet,
		NewApp,
	)
	return nil, nil, nil
}
