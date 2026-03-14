package service

import (
	"context"
	"fmt"
	"gin-krotas-gorm/internal/biz"
	"gin-krotas-gorm/internal/biz/model"
	"gin-krotas-gorm/internal/data/persistence"
	"strconv"
	"time"

	pb "gin-krotas-gorm/api/order/v1"
)

type OrderService struct {
	pb.UnimplementedOrderServer
	orderRepo   *persistence.OrderRepo
	productRepo *persistence.ProductRepo // 需要调用商品 Repo 扣库存
	tx          biz.Transaction
}

func NewOrderService(or *persistence.OrderRepo, pr *persistence.ProductRepo, tx biz.Transaction) *OrderService {
	return &OrderService{orderRepo: or, productRepo: pr, tx: tx}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderReply, error) {
	// 1. 获取当前用户 (从 JWT 中间件存入的 Context 获取)
	userId := ctx.Value("userID").(uint64)

	var orderSn string
	var totalPayAmount float64

	// 2. 执行业务逻辑 (使用事务封装)
	err := s.tx.InTx(ctx, func(ctx context.Context) error {
		// A. 扣减库存 (内部调用 ProductRepo 的 DB(ctx) 确保在事务内)
		for _, item := range req.Items {
			err := s.productRepo.ReduceStock(ctx, int64(item.SkuId), item.Num)
			if err != nil {
				return err // 失败自动回滚
			}
		}

		// B. 预计算价格 (此处简化，实际应从数据库查 SKU 单价)
		orderSn = fmt.Sprintf("%d%06d", time.Now().Unix(), userId%1000)

		// C. 构建订单实体
		orderEntity := &model.Order{
			OrderSn:         orderSn,
			UserID:          int64(userId),
			TotalAmount:     319.0,
			PayAmount:       99.0, // 示例数值
			ReceiverName:    "张三",
			ReceiverAddress: strconv.Itoa(int(req.AddressId)),
			Remark:          req.Remark,
			PayTime:         time.Now(),
		}

		orderItems := make([]*model.OrderItem, 0)
		for _, item := range req.Items {
			oi := &model.OrderItem{
				Price: 99.0,
				SpuID: 1,
				SkuID: int64(item.SkuId),
				Num:   item.Num,
			}
			orderItems = append(orderItems, oi)
		}

		// D. 写入订单表
		return s.orderRepo.CreateOrder(ctx, orderEntity, orderItems)
	})

	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderReply{
		OrderSn:     orderSn,
		TotalAmount: totalPayAmount,
	}, nil
}

// 伪代码示例
func (s *OrderService) sendDelayCancelMessage(orderSn string) {
	// 使用 Redis ZSet 或 RocketMQ 延迟消息
	// 30 分钟后检查该订单状态，若仍为“待支付”，则执行：
	// 1. 修改订单状态为“已取消”
	// 2. 调用 productRepo.AddStock 返还库存
}

func (s *OrderService) GetOrderDetail(ctx context.Context, req *pb.GetOrderDetailReq) (*pb.GetOrderDetailReply, error) {
	order, err := s.orderRepo.GetOrderWithItems(ctx, req.OrderSn)
	if err != nil {
		return nil, err
	}
	reply := &pb.GetOrderDetailReply{
		OrderSn:         order.OrderSn,
		TotalAmount:     order.TotalAmount,
		PayAmount:       order.PayAmount,
		ReceiverName:    order.ReceiverName,
		ReceiverPhone:   order.ReceiverPhone,
		ReceiverAddress: order.ReceiverAddress,
	}

	for _, item := range order.OrderItems {
		reply.Items = append(reply.Items, &pb.OrderItemInfo{
			SkuId:   uint32(item.SkuID),
			SkuName: item.SkuName,
			SkuPic:  item.SkuPic,
			Price:   item.Price,
			Num:     item.Num,
		})
	}

	return reply, nil
}
