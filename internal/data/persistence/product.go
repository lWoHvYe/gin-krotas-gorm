package persistence

import (
	"context"
	"errors"
	"helloworld-go/internal/biz/model"
	"helloworld-go/internal/biz/query"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ProductRepo struct {
	data *Data
	log  *log.Helper
}

func NewProductRepo(data *Data, logger log.Logger) *ProductRepo {
	return &ProductRepo{data: data, log: log.NewHelper(logger)}
}

func (p *ProductRepo) GetProductWithSkus(ctx context.Context, productId int64) (*model.ProductSpu, error) {
	q := p.data.Q(ctx)
	return q.ProductSpu.
		WithContext(ctx).
		Preload(q.ProductSpu.Skus). // 单纯一对多，用Preload
		Where(q.ProductSpu.ID.Eq(productId)).
		First()
}

func (p *ProductRepo) GetProductByPage(ctx context.Context, Page int, PageSize int) ([]*model.ProductSpu, int64, error) {
	q := p.data.Q(ctx)
	offset := (Page - 1) * PageSize
	return q.ProductSpu.
		WithContext(ctx).
		Preload(q.ProductSpu.Skus).
		Where(q.ProductSpu.IsOnSale.Is(true)).
		FindByPage(offset, PageSize)
}

func (p *ProductRepo) GetProductHasStock(ctx context.Context, productId int64) (*model.ProductSpu, error) {
	q := p.data.Q(ctx)
	return q.ProductSpu.
		WithContext(ctx).
		Preload(q.ProductSpu.Skus.Where(query.ProductSku.Stock.Gt(0))).
		Where(q.ProductSpu.ID.Eq(productId)).
		First()
}

func (p *ProductRepo) GetBestProduct(ctx context.Context) ([]query.SpuWithMinPrice, error) {
	q := p.data.Q(ctx)
	var res []query.SpuWithMinPrice
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
	q := p.data.Q(ctx)
	return q.ProductSpu.WithContext(ctx).Where(q.ProductSpu.ID.Eq(id)).First()
}

func (p *ProductRepo) ListSkusBySpu(ctx context.Context, spuId int64) ([]*model.ProductSku, error) {
	q := p.data.Q(ctx)
	return q.ProductSku.WithContext(ctx).Where(q.ProductSku.SpuID.Eq(spuId)).Find()
}

func (p *ProductRepo) ReduceStock(ctx context.Context, skuId int64, num int32) error {
	q := p.data.Q(ctx)
	// 2. 操作商品表 (扣库存)
	info, err := q.ProductSku.WithContext(ctx).
		Where(q.ProductSku.ID.Eq(skuId), q.ProductSku.Stock.Gte(int64(num))).
		Update(q.ProductSku.Stock, gorm.Expr("stock - ?", num))
	if info.RowsAffected == 0 {
		return errors.New("库存不足")
	}
	return err
}
