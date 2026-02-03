package server

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
var TransportSet = wire.NewSet(NewGRPCServer, NewHTTPServer)
