package data

import (
	"gin-kratos-gorm/internal/order/biz"
	"gin-kratos-gorm/internal/order/data/persistence"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var RepoSet = wire.NewSet(persistence.NewData, persistence.NewOrderRepo,
	wire.Bind(new(biz.Transaction), new(*persistence.Data)))
