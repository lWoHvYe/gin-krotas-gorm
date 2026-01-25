package pkg

import (
	"helloworld-go/internal/data/db"
	"helloworld-go/internal/pkg/config"
	"helloworld-go/internal/pkg/logger"

	"github.com/google/wire"
)

var InfraSet = wire.NewSet(
	config.NewConfig,
	logger.NewLogger,
	db.NewDB,
)
