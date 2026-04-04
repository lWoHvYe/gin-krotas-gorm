package server

import (
	httpServerHandler "gin-kratos-gorm/internal/order/server/http/handler"

	"github.com/google/wire"
)

// ProviderSet is server providers.
var TransportSet = wire.NewSet(NewGRPCServer, NewHTTPServer, httpServerHandler.NewOrderHandler)
