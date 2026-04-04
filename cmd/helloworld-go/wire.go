//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"gin-kratos-gorm/internal/biz"
	"gin-kratos-gorm/internal/data"
	"gin-kratos-gorm/internal/pkg"
	"gin-kratos-gorm/internal/server"
	"gin-kratos-gorm/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(string, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(pkg.InfraSet, server.TransportSet, data.RepoSet, biz.BizSet, service.ServiceSet, newApp))
}
