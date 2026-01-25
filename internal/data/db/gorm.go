// internal/infrastructure/db/gorm.go
package db

import (
	"helloworld-go/internal/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(cfg *conf.Bootstrap) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.Data.Database.Source))
}
