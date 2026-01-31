// internal/infrastructure/db/gorm.go
package db

import (
	"helloworld-go/internal/biz/product"
	"helloworld-go/internal/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(cfg *conf.Bootstrap) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Data.Database.Source))
	if err != nil {
		panic(err)
	}
	// 自动迁移（根据定义同步db表结构，没有表会自动创建）
	db.AutoMigrate(&product.ProductSpu{}, &product.ProductSku{})
	return db, err
}
