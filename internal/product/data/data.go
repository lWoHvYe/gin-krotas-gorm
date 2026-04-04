package data

import (
	"gin-kratos-gorm/internal/product/biz"
	"gin-kratos-gorm/internal/product/data/persistence"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var RepoSet = wire.NewSet(persistence.NewData, persistence.NewProductRepo,
	wire.Bind(new(biz.Transaction), new(*persistence.Data)))
