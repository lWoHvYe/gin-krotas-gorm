package product

import (
	"gorm.io/gorm"
)

// ProductSpu 标准产品单位
type ProductSpu struct {
	gorm.Model
	Name        string `gorm:"size:255;not null;comment:商品名称"`
	Description string `gorm:"type:text;comment:商品详情介绍"`
	CategoryId  uint   `gorm:"index;comment:分类ID"`
	BrandId     uint   `gorm:"index;comment:品牌ID"`
	MainImage   string `gorm:"size:512;comment:主图URL"`
	IsOnSale    bool   `gorm:"default:false;comment:是否上架"`
	// 一对多关系：一个 SPU 对应多个 SKU
	Skus []ProductSku `gorm:"foreignKey:SpuId"`
}

// ProductSku 库存保存单位
type ProductSku struct {
	gorm.Model
	SpuId     uint    `gorm:"index;comment:所属SPU_ID"`
	SkuCode   string  `gorm:"size:64;uniqueIndex;comment:SKU唯一编码"`
	AttrValue string  `gorm:"type:json;comment:属性组合JSON"` // GORM 支持 JSON 类型
	Price     float64 `gorm:"type:decimal(10,2);comment:售价"`
	Stock     int     `gorm:"default:0;comment:库存数量"`
}
