package data

import (
	"gin-krotas-gorm/internal/biz"
	"gin-krotas-gorm/internal/data/persistence"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var RepoSet = wire.NewSet(persistence.NewData, NewGreeterRepo, persistence.NewUserRepo, NewEnforcer,
	persistence.NewProductRepo,
	persistence.NewOrderRepo,
	wire.Bind(new(biz.Transaction), new(*persistence.Data)))
