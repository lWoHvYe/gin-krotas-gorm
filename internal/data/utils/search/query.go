package search

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	FromQueryTag = "search"
	Mysql        = "mysql"
	Postgres     = "postgres"
)

// ResolveSearchQuery 解析搜索结构体
func ResolveSearchQuery(driver string, q interface{}, condition Condition) {
	qType := reflect.TypeOf(q)
	qValue := reflect.ValueOf(q)

	// 安全检查：如果不是结构体，直接返回，防止 NumField 崩溃
	if qType.Kind() == reflect.Ptr {
		qType = qType.Elem()
		qValue = qValue.Elem()
	}
	if qType.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < qType.NumField(); i++ {
		field := qType.Field(i)
		value := qValue.Field(i)

		tag, ok := field.Tag.Lookup(FromQueryTag)
		if !ok {
			// 如果字段是结构体（如内嵌的 PageInfo），尝试递归解析
			if value.Kind() == reflect.Struct {
				ResolveSearchQuery(driver, value.Interface(), condition)
			}
			continue
		}

		if tag == "-" {
			continue
		}

		// 跳过零值处理。对于 time.Time，IsZero() 是有效的判断
		if value.IsZero() {
			continue
		}

		t := makeTag(tag)
		// 统一处理 SQL 构建逻辑
		buildCondition(driver, t, condition, value)
	}
}

// buildCondition 统一构建查询语句，兼容 Postgres 和其他驱动
func buildCondition(driver string, t *resolveSearchTag, condition Condition, value reflect.Value) {
	quote := "`"
	if driver == Postgres {
		quote = "" // Postgres 通常不需要反引号，或使用双引号
	}

	column := fmt.Sprintf("%s%s%s.%s%s%s", quote, t.Table, quote, quote, t.Column, quote)
	val := value.Interface()

	switch t.Type {
	case "left":
		joinOn := fmt.Sprintf("left join %s%s%s on %s%s%s.%s = %s%s%s.%s",
			quote, t.Join, quote, quote, t.Join, quote, t.On[0], quote, t.Table, quote, t.On[1])
		join := condition.SetJoinOn(t.Type, joinOn)
		ResolveSearchQuery(driver, val, join)
	case "eq":
		condition.SetWhere(column+" = ?", []interface{}{val})
	case "neq":
		condition.SetWhere(column+" <> ?", []interface{}{val})
	case "contains":
		condition.SetWhere(column+" LIKE ?", []interface{}{"%" + fmt.Sprint(val) + "%"})
	case "gt":
		condition.SetWhere(column+" > ?", []interface{}{val})
	case "gte":
		condition.SetWhere(column+" >= ?", []interface{}{val})
	case "lt":
		condition.SetWhere(column+" < ?", []interface{}{val})
	case "lte":
		condition.SetWhere(column+" <= ?", []interface{}{val})
	case "between":
		// 支持 CreatedAtRange []time.Time 这种切片
		if value.Kind() == reflect.Slice && value.Len() == 2 {
			condition.SetWhere(column+" BETWEEN ? AND ?", []interface{}{value.Index(0).Interface(), value.Index(1).Interface()})
		}
	case "in":
		condition.SetWhere(column+" IN (?)", []interface{}{val})
	case "isnull":
		condition.SetWhere(column+" IS NULL", nil)
	case "order":
		orderType := strings.ToLower(fmt.Sprint(val))
		if orderType == "desc" || orderType == "descending" {
			condition.SetOrder(column + " DESC")
		} else {
			condition.SetOrder(column + " ASC")
		}
	}
}
