package persistence

import (
	"context"
	"errors"
	"fmt"
	"helloworld-go/internal/product/biz/model"
	"helloworld-go/internal/product/biz/query"

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

func (p *ProductRepo) GetDB(ctx context.Context) *gorm.DB {
	return p.data.DB(ctx)
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

// 使用 GORM 执行带幂等校验的库存扣减
func (p *ProductRepo) ReduceStock(ctx context.Context, orderSn string, skuId int64, num int32) error {
	q := p.data.Q(ctx)

	return q.Transaction(func(tx *query.Query) error {
		// 1. 插入幂等记录，如果已存在则报错回滚（利用数据库唯一索引）
		if err := tx.StockLog.WithContext(ctx).Create(&model.StockLog{ID: orderSn, SkuID: skuId, Num: num}); err != nil {
			return errors.New("duplicate request")
		}

		// 2. 扣减库存 (乐观锁/条件更新)
		res, err := tx.ProductSku.WithContext(ctx).
			Where(q.ProductSku.ID.Eq(skuId), q.ProductSku.Stock.Gte(int64(num))).
			Update(q.ProductSku.Stock, gorm.Expr("stock - ?", num))

		if err != nil {
			return err
		}

		if res.RowsAffected == 0 {
			return errors.New("insufficient stock")
		}
		return nil
	})
}

func (p *ProductRepo) CompensateStock(ctx context.Context, orderSn string, skuId int64, num int32) error {
	q := p.data.Q(ctx)

	return q.Transaction(func(tx *query.Query) error {
		if info, err := tx.StockLog.WithContext(ctx).Where(q.StockLog.ID.Eq(orderSn)).Delete(); err != nil {
			return err
		} else if info.RowsAffected == 0 {
			return nil // duplicate call
		}

		_, err := tx.ProductSku.WithContext(ctx).
			Where(q.ProductSku.ID.Eq(skuId)). // 这里需考虑并发更新问题
			Update(q.ProductSku.Stock, gorm.Expr("stock + ?", num))
		if err != nil {
			return err
		}
		return nil
	})
}

func (p *ProductRepo) ReduceStockInTx(ctx context.Context, gtx *gorm.DB, orderSn string, skuId int64, num int32) error {
	q := query.Use(gtx)
	return q.Transaction(func(tx *query.Query) error {
		// 1. 插入幂等记录，如果已存在则报错回滚（利用数据库唯一索引）
		stockLogId := fmt.Sprintf("%s_%06d", orderSn, skuId) // 避免与其他业务的 stock log 冲突
		if err := tx.StockLog.WithContext(ctx).Create(&model.StockLog{ID: stockLogId, SkuID: skuId, Num: num}); err != nil {
			return errors.New("duplicate request")
		}

		_, err := tx.ProductSku.WithContext(ctx).
			Where(q.ProductSku.ID.Eq(skuId)). // 这里需考虑并发更新问题
			Update(q.ProductSku.Stock, gorm.Expr("stock - ?", num))
		if err != nil {
			return err
		}
		return nil
	})
}
