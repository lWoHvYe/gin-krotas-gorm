package server

import (
	grpcTransport "helloworld-go/internal/server/grpc"
	httpServerHandler "helloworld-go/internal/server/http/handler"

	"github.com/google/wire"
)

// ProviderSet is server providers.
var TransportSet = wire.NewSet(NewGRPCServer, NewHTTPServer,
	grpcTransport.NewUserServer, httpServerHandler.NewUserHandler)
