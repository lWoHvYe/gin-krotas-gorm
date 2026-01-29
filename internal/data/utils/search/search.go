package search

import (
	"gorm.io/gorm"
)

// MakeCondition 生成 GORM Scope
func MakeCondition(q interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		condition := &GormCondition{
			GormPublic: GormPublic{
				Where: make(map[string][]interface{}),
				Or:    make(map[string][]interface{}),
				Order: make([]string, 0),
			},
			Join: make([]*GormJoin, 0),
		}

		// 解析查询对象
		ResolveSearchQuery(Mysql, q, condition)

		// 应用 Join 逻辑
		for _, join := range condition.Join {
			if join == nil {
				continue
			}
			db = db.Joins(join.JoinOn)
			applyPublic(db, &join.GormPublic)
		}

		// 应用 Where/Order 逻辑
		applyPublic(db, &condition.GormPublic)

		return db
	}
}

// 辅助函数：应用 Where, Or, Order 到 DB 实例
func applyPublic(db *gorm.DB, p *GormPublic) {
	for k, v := range p.Where {
		db.Where(k, v...)
	}
	for k, v := range p.Or {
		db.Or(k, v...)
	}
	for _, o := range p.Order {
		db.Order(o)
	}
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(offset).Limit(pageSize)
	}
}
