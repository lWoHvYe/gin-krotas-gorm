package persistence

import (
	"context"
	"helloworld-go/internal/biz/model"
	"helloworld-go/internal/biz/product"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db  *gorm.DB
	log *log.Helper
}

func NewProductRepo(db *gorm.DB, logger log.Logger) *ProductRepo {
	return &ProductRepo{db: db, log: log.NewHelper(logger)}
}

func (p *ProductRepo) GetProductWithSkus(ctx context.Context, productId int64) (*model.ProductSpu, error) {
	q := product.Use(p.db)
	return q.ProductSpu.
		WithContext(ctx).
		Preload(q.ProductSpu.Skus). // 单纯一对多，用Preload
		Where(q.ProductSpu.ID.Eq(productId)).
		First()
}

func (p *ProductRepo) GetProductByPage(ctx context.Context, Page int, PageSize int) ([]*model.ProductSpu, int64, error) {
	q := product.Use(p.db)
	offset := (Page - 1) * PageSize
	return q.ProductSpu.
		WithContext(ctx).
		Preload(q.ProductSpu.Skus).
		Where(q.ProductSpu.IsOnSale.Is(true)).
		FindByPage(offset, PageSize)
}

func (p *ProductRepo) GetProductHasStock(ctx context.Context, productId int64) (*model.ProductSpu, error) {
	q := product.Use(p.db)
	return q.ProductSpu.
		WithContext(ctx).
		Preload(q.ProductSpu.Skus.Where(product.ProductSku.Stock.Gt(0))).
		Where(q.ProductSpu.ID.Eq(productId)).
		First()
}

func (p *ProductRepo) GetBestProduct(ctx context.Context) ([]product.SpuWithMinPrice, error) {
	q := product.Use(p.db)
	var res []product.SpuWithMinPrice
	err := q.ProductSpu.WithContext(ctx).
		Select(
			q.ProductSpu.ID,
			q.ProductSpu.Name,
			q.ProductSku.Price.Min().As("min_price"),
		).
		Join(q.ProductSku,
			q.ProductSku.SpuID.EqCol(q.ProductSpu.ID)).
		Group(q.ProductSpu.ID).
		Scan(&res)
	return res, err
}

// 不含sku
func (p *ProductRepo) GetSpu(ctx context.Context, id int64) (*model.ProductSpu, error) {
	q := product.Use(p.db)
	return q.ProductSpu.WithContext(ctx).Where(q.ProductSpu.ID.Eq(id)).First()
}

func (p *ProductRepo) ListSkusBySpu(ctx context.Context, spuId int64) ([]*model.ProductSku, error) {
	q := product.Use(p.db)
	return q.ProductSku.WithContext(ctx).Where(q.ProductSku.SpuID.Eq(spuId)).Find()
}
