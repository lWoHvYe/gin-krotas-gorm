package pkg

import (
	"helloworld-go/internal/order/data/db"
	"helloworld-go/internal/order/pkg/config"
	"helloworld-go/internal/order/pkg/grpcclient"
	"helloworld-go/internal/pkg/logger"

	"github.com/google/wire"
)

var InfraSet = wire.NewSet(
	config.NewConfig,
	logger.NewLogger,
	db.NewDB,
	config.NewDiscovery,
	grpcclient.NewConnFactory,
)
