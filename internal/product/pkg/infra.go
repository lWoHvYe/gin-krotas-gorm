package pkg

import (
	"gin-krotas-gorm/internal/pkg/logger"
	"gin-krotas-gorm/internal/product/data/db"
	"gin-krotas-gorm/internal/product/pkg/config"

	"github.com/google/wire"
)

var InfraSet = wire.NewSet(
	config.NewConfig,
	logger.NewLogger,
	db.NewDB,
	config.NewDiscovery,
)
