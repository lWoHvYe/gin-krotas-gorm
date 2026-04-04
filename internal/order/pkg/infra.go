package pkg

import (
	"gin-kratos-gorm/internal/order/data/db"
	"gin-kratos-gorm/internal/order/pkg/config"
	"gin-kratos-gorm/internal/order/pkg/grpcclient"
	"gin-kratos-gorm/internal/pkg/logger"

	"github.com/google/wire"
)

var InfraSet = wire.NewSet(
	config.NewConfig,
	logger.NewLogger,
	db.NewDB,
	config.NewDiscovery,
	grpcclient.NewConnFactory,
)
