package data

import (
	"helloworld-go/internal/product/biz"
	"helloworld-go/internal/product/data/persistence"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var RepoSet = wire.NewSet(persistence.NewData, persistence.NewProductRepo,
	wire.Bind(new(biz.Transaction), new(*persistence.Data)))
