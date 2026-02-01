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

	// 1. 定义基本模型（不要在这里加关系，防止循环引用导致空指针）
	u := g.GenerateModel("users")
	r := g.GenerateModel("roles")
	ur := g.GenerateModel("user_roles")

	f_userId_tags := make(field.GormTag)
	f_userId_tags.Append("foreignKey", "UserID")
	f_roleId_tags := make(field.GormTag)
	f_roleId_tags.Append("foreignKey", "RoleID")
	// 2. 重新定义带关系的 UserRole (它需要引用 User 和 Role)
	ur = g.GenerateModel("user_roles",
		gen.FieldRelate(field.BelongsTo, "User", u, &field.RelateConfig{GORMTag: f_userId_tags}),
		gen.FieldRelate(field.BelongsTo, "Role", r, &field.RelateConfig{GORMTag: f_roleId_tags}),
	)

	// 3. 重新定义带关系的 User 和 Role (它们需要引用已经带了关系的 UserRole)
	u = g.GenerateModel("users",
		gen.FieldRelate(field.HasMany, "UserRoles", ur, &field.RelateConfig{GORMTag: f_userId_tags}),
	)
	r = g.GenerateModel("roles",
		gen.FieldRelate(field.HasMany, "UserRoles", ur, &field.RelateConfig{GORMTag: f_roleId_tags}),
	)

	// 4. 最后统一应用
	g.ApplyBasic(u, r, ur)

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
