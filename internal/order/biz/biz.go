package biz

import (
	"context"

	"github.com/google/wire"
)

// BizSet is biz providers.
var BizSet = wire.NewSet()

// Transaction 事务管理器接口 (由 data 层实现)
type Transaction interface {
	InTx(context.Context, func(ctx context.Context) error) error
}
