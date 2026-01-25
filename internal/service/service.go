package service

import (
	userService "helloworld-go/internal/service/user"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ServiceSet = wire.NewSet(NewGreeterService, userService.NewUserService)
