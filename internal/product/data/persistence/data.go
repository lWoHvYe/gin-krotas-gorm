package persistence

import (
	"context"
	"helloworld-go/internal/conf"
	"helloworld-go/internal/product/biz/query"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// Data .
type Data struct {
	// TODO wrapped database client
	db    *gorm.DB
	Query *query.Query
}

// NewData .
func NewData(c *conf.Bootstrap, logger log.Logger, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db, Query: query.Use(db)}, cleanup, nil
}

// InTx 事务封装
func (d *Data) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, "contextTxKey", tx) // 将事务对象注入 Context
		return fn(ctx)
	})
}

// DB 获取连接（判断 Context 是否有事务）
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("contextTxKey").(*gorm.DB)
	if ok {
		return tx
	}
	return d.db.WithContext(ctx)
}

// 获取当前 Context 下的 Query 对象（自动判断是否在事务中）
func (d *Data) Q(ctx context.Context) *query.Query {
	tx, ok := ctx.Value("contextTxKey").(*gorm.DB)
	if ok {
		// 如果在事务中，基于事务 DB 重新生成 Query 映射
		return query.Use(tx)
	}
	return d.Query
}
