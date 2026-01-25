package biz

import "github.com/google/wire"

// BizSet is biz providers.
var BizSet = wire.NewSet(NewGreeterUsecase)
