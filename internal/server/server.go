package server

import (
	httpServerHandler "gin-krotas-gorm/internal/server/http/handler"

	"github.com/google/wire"
)

// ProviderSet is server providers.
var TransportSet = wire.NewSet(NewGRPCServer, NewHTTPServer,
	httpServerHandler.NewUserHandler, httpServerHandler.NewProductHandler, httpServerHandler.NewOrderHandler)
