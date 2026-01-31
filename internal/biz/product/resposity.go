package product

import "context"

type ProductRepository interface {
	GetSpu(ctx context.Context, id uint) (ProductSpu, error)
	ListSkusBySpu(ctx context.Context, spuId uint) ([]ProductSku, error)
	GetProductWithSkus(ctx context.Context, productId uint) (ProductSpu, error)
}
