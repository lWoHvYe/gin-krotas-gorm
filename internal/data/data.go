package data

import (
	"helloworld-go/internal/biz"
	"helloworld-go/internal/data/persistence"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var RepoSet = wire.NewSet(persistence.NewData, NewGreeterRepo, persistence.NewUserRepo, NewEnforcer,
	persistence.NewProductRepo,
	persistence.NewOrderRepo,
	wire.Bind(new(biz.Transaction), new(*persistence.Data)))
