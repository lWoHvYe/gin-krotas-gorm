//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

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

// wireApp init kratos application.
func wireApp(string, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(pkg.InfraSet, server.TransportSet, data.RepoSet, biz.BizSet, service.ServiceSet, newApp))
}
