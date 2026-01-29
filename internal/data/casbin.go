package data

import (
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewEnforcer 提供 Casbin 权限执行器
func NewEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	// 使用 GORM 适配器，将权限策略持久化到数据库
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 加载模型配置文件
	enforcer, err := casbin.NewEnforcer("configs/rbac_model.conf", adapter)
	if err != nil {
		return nil, err
	}

	// 启用自动加载策略（可选，用于多实例同步）
	err = enforcer.LoadPolicy()
	return enforcer, err
}
