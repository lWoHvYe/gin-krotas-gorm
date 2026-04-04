package pkg

import (
	"gin-kratos-gorm/internal/data/db"
	"gin-kratos-gorm/internal/pkg/config"
	"gin-kratos-gorm/internal/pkg/logger"

	"github.com/google/wire"
)

var InfraSet = wire.NewSet(
	config.NewConfig,
	logger.NewLogger,
	db.NewDB,
)
