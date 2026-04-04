package data

import (
	"gin-kratos-gorm/internal/biz"
	"gin-kratos-gorm/internal/data/persistence"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var RepoSet = wire.NewSet(persistence.NewData, NewGreeterRepo, persistence.NewUserRepo, NewEnforcer,
	persistence.NewProductRepo,
	persistence.NewOrderRepo,
	wire.Bind(new(biz.Transaction), new(*persistence.Data)))
