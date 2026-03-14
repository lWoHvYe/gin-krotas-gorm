//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package bootstrap

import (
	"gin-krotas-gorm/internal/biz"
	"gin-krotas-gorm/internal/data"
	"gin-krotas-gorm/internal/pkg"
	"gin-krotas-gorm/internal/server"
	"gin-krotas-gorm/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func InitApp(string, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		biz.BizSet,
		pkg.InfraSet,
		data.RepoSet,
		service.ServiceSet,
		server.TransportSet,
		NewApp,
	))
}
