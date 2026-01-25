package service

import (
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ServiceSet = wire.NewSet(NewGreeterService, NewUserService)
