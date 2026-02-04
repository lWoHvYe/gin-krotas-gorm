package pkg

import (
	"helloworld-go/internal/pkg/logger"
	"helloworld-go/internal/product/data/db"
	"helloworld-go/internal/product/pkg/config"

	"github.com/google/wire"
)

var InfraSet = wire.NewSet(
	config.NewConfig,
	logger.NewLogger,
	db.NewDB,
	config.NewDiscovery,
)
