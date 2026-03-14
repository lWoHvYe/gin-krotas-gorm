// internal/infrastructure/db/gorm.go
package db

import (
	"gin-krotas-gorm/internal/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(cfg *conf.Bootstrap) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Data.Database.Source))
	if err != nil {
		panic(err)
	}
	// 自动迁移（根据定义同步db表结构，没有表会自动创建）建议还是先生成表，在根据表生成
	//db.AutoMigrate(&model.ProductSpu{}, &model.ProductSku{}, &model.Order{}, &model.OrderItem{})
	return db, err
}
