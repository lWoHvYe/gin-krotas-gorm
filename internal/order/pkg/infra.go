package pkg

import (
	"gin-krotas-gorm/internal/order/data/db"
	"gin-krotas-gorm/internal/order/pkg/config"
	"gin-krotas-gorm/internal/order/pkg/grpcclient"
	"gin-krotas-gorm/internal/pkg/logger"

	"github.com/google/wire"
)

var InfraSet = wire.NewSet(
	config.NewConfig,
	logger.NewLogger,
	db.NewDB,
	config.NewDiscovery,
	grpcclient.NewConnFactory,
)
