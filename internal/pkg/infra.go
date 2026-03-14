package pkg

import (
	"gin-krotas-gorm/internal/data/db"
	"gin-krotas-gorm/internal/pkg/config"
	"gin-krotas-gorm/internal/pkg/logger"

	"github.com/google/wire"
)

var InfraSet = wire.NewSet(
	config.NewConfig,
	logger.NewLogger,
	db.NewDB,
)
