package persistence

import (
	pb "helloworld-go/api/product/v1"
	"helloworld-go/internal/order/biz/model"
	"helloworld-go/internal/order/pkg/grpcclient"

	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type OrderRepo struct {
	data          *Data
	productClient pb.ProductClient // 注入远程客户端
}

func NewOrderRepo(data *Data, factory *grpcclient.ConnFactory, logger log.Logger) *OrderRepo {
	conn := factory.Get(pb.Product_ServiceDesc.ServiceName)
	return &OrderRepo{
		data:          data,
		productClient: pb.NewProductClient(conn),
	}
}

// ReduceStock 远程调用商品服务扣库存
func (r *OrderRepo) ReduceStock(ctx context.Context, skuId int64, num int32) error {
	_, err := r.productClient.ReduceStock(ctx, &pb.ReduceStockReq{
		SkuId: skuId,
		Num:   num,
	})
	return err
}

func (r *OrderRepo) CreateOrder(ctx context.Context, order *model.Order, items []*model.OrderItem) error {
	// 1. 获取当前上下文的 Query (如果是事务中，这里的 q 会绑定事务 DB)
	q := r.data.Q(ctx)
	// 1. 先创建订单主表
	if err := q.Order.WithContext(ctx).Create(order); err != nil {
		return err
	}

	// 2. 为每个 Item 绑定刚刚生成的 OrderId (假设 order.ID 已被 GORM 回填)
	for i := range items {
		items[i].OrderID = order.ID
		items[i].OrderSn = order.OrderSn
	}

	// 3. 批量创建订单项
	return q.OrderItem.WithContext(ctx).Create(items...)
}

func (r *OrderRepo) GetOrderWithItems(ctx context.Context, orderSn string) (*model.Order, error) {
	q := r.data.Q(ctx)
	return q.Order.
		WithContext(ctx).Preload(q.Order.OrderItems).
		Where(q.Order.OrderSn.Eq(orderSn)).First()
}

func (r *OrderRepo) CancelOrder(ctx context.Context, orderSn string) error {
	q := r.data.Q(ctx)
	_, err := q.Order.WithContext(ctx).Where(q.Order.OrderSn.Eq(orderSn)).Delete()
	if err != nil {
		return err
	}
	_, err = q.OrderItem.WithContext(ctx).Where(q.OrderItem.OrderSn.Eq(orderSn)).Delete()
	if err != nil {
		return err
	}
	return nil
}
