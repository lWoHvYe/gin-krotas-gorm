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
		OutPath:           "internal/biz/order",
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldWithTypeTag:  true,
		FieldWithIndexTag: true,
	})

	g.UseDB(db)

	/*	sku := g.GenerateModel("product_skus")

		tags := make(field.GormTag)
		tags.Append("foreignKey", "SpuId")

		spu := g.GenerateModel("product_spus",
			gen.FieldRelate(
				field.HasMany,
				"Skus",
				sku, &field.RelateConfig{GORMTag: tags},
			),
		)

		g.ApplyBasic(spu, sku)*/

	order_item := g.GenerateModel("order_items")

	tags := make(field.GormTag)
	tags.Append("foreignKey", "OrderId")

	order := g.GenerateModel("orders",
		gen.FieldRelate(
			field.HasMany,
			"OrderItems",
			order_item, &field.RelateConfig{GORMTag: tags},
		),
	)

	g.ApplyBasic(order, order_item)

	g.Execute()
}
