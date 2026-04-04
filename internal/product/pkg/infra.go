package pkg

import (
	"gin-kratos-gorm/internal/pkg/logger"
	"gin-kratos-gorm/internal/product/data/db"
	"gin-kratos-gorm/internal/product/pkg/config"

	"github.com/google/wire"
)

var InfraSet = wire.NewSet(
	config.NewConfig,
	logger.NewLogger,
	db.NewDB,
	config.NewDiscovery,
)
