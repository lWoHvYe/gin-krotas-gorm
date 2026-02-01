package persistence

import (
	"context"
	"helloworld-go/internal/biz/model"

	"github.com/go-kratos/kratos/v2/log"
)

type OrderRepo struct {
	data *Data
	log  *log.Helper
}

func NewOrderRepo(data *Data, logger log.Logger) *OrderRepo {
	return &OrderRepo{data: data, log: log.NewHelper(logger)}
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
