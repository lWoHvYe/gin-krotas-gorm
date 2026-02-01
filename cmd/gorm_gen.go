package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func main() {
	db, _ := gorm.Open(mysql.Open("root:root@tcp(10.211.55.29:3306)/unicorn?parseTime=True&loc=Local"))

	g := gen.NewGenerator(gen.Config{
		OutPath:           "internal/biz/user",
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldWithTypeTag:  true,
		FieldWithIndexTag: true,
	})

	g.UseDB(db)

	user := g.GenerateModel("users")
	role := g.GenerateModel("roles")
	userRole := g.GenerateModel("user_roles")

	user = g.GenerateModel("users",
		gen.FieldRelate(field.HasMany, "UserRoles", userRole, nil),
	)

	role = g.GenerateModel("roles",
		gen.FieldRelate(field.HasMany, "UserRoles", userRole, nil),
	)

	userRole = g.GenerateModel("user_roles",
		gen.FieldRelate(field.BelongsTo, "User", user, nil),
		gen.FieldRelate(field.BelongsTo, "Role", role, nil),
	)

	g.ApplyBasic(user, role, userRole)

	/*sku := g.GenerateModel("product_skus")

	product_tags := make(field.GormTag)
	product_tags.Append("foreignKey", "SpuId")

	spu := g.GenerateModel("product_spus",
		gen.FieldRelate(
			field.HasMany,
			"Skus",
			sku, &field.RelateConfig{GORMTag: product_tags},
		),
	)

	// 生成之前定义的商品和订单模型 这种直接根据结构体生成似乎也可以
	//g.ApplyBasic(data.ProductSpu{}, data.ProductSku{}, data.Order{}, data.OrderItem{})

	order_item := g.GenerateModel("order_items")

	order_tags := make(field.GormTag)
	order_tags.Append("foreignKey", "OrderId")

	order := g.GenerateModel("orders",
		gen.FieldRelate(
			field.HasMany,
			"OrderItems",
			order_item, &field.RelateConfig{GORMTag: order_tags},
		),
	)

	g.ApplyBasic(spu, sku, order, order_item)*/

	g.Execute()
}
