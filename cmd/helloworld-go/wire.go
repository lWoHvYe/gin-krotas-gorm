//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"helloworld-go/internal/biz"
	"helloworld-go/internal/conf"
	"helloworld-go/internal/data"
	"helloworld-go/internal/server"
	"helloworld-go/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.TransportSet, data.RepoSet, biz.BizSet, service.ServiceSet, newApp))
}
