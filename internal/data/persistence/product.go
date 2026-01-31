package persistence

import (
	"context"
	"helloworld-go/internal/biz/product"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db  *gorm.DB
	log *log.Helper
}

func NewProductRepo(db *gorm.DB, logger log.Logger) product.ProductRepository {
	return &ProductRepo{db: db, log: log.NewHelper(logger)}
}

func (p *ProductRepo) GetProductWithSkus(ctx context.Context, productId uint) (product.ProductSpu, error) {
	return gorm.G[product.ProductSpu](p.db).Where("id = ?", productId).Preload("Skus", nil).First(ctx)
}

// 不含sku
func (p *ProductRepo) GetSpu(ctx context.Context, id uint) (product.ProductSpu, error) {
	return gorm.G[product.ProductSpu](p.db).Where("id = ?", id).First(ctx)
}

func (p *ProductRepo) ListSkusBySpu(ctx context.Context, spuId uint) ([]product.ProductSku, error) {
	return gorm.G[product.ProductSku](p.db).Where("spu_id = ?", spuId).Find(ctx)
}
