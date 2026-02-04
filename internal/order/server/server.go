package server

import (
	httpServerHandler "helloworld-go/internal/order/server/http/handler"

	"github.com/google/wire"
)

// ProviderSet is server providers.
var TransportSet = wire.NewSet(NewGRPCServer, NewHTTPServer, httpServerHandler.NewOrderHandler)
