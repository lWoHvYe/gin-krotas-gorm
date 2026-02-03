package service

import (
	"context"
	"fmt"
	"helloworld-go/internal/conf"
	"helloworld-go/internal/order/biz/model"
	"helloworld-go/internal/order/data/persistence"
	"strconv"
	"time"

	"github.com/dtm-labs/client/dtmgrpc"

	pb "helloworld-go/api/order/v1"
	pbProduct "helloworld-go/api/product/v1"
)

type OrderService struct {
	pb.UnimplementedOrderServer
	orderRepo *persistence.OrderRepo
	dtmC      *conf.Data_DTM
}

func NewOrderService(or *persistence.OrderRepo, c *conf.Bootstrap) *OrderService {
	return &OrderService{orderRepo: or, dtmC: c.Data.Dtm}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderReply, error) {
	gid := dtmgrpc.MustGenGid(s.dtmC.Address)

	// 1. 构造 SAGA 事务
	saga := dtmgrpc.NewSagaGrpc(s.dtmC.Address, gid)

	// 1. 获取当前用户 (从 JWT 中间件存入的 Context 获取)
	userId := ctx.Value("userID").(uint64)

	var orderSn string

	// B. 预计算价格 (此处简化，实际应从数据库查 SKU 单价)
	// 生成订单号
	orderSn = fmt.Sprintf("%d%06d", time.Now().Unix(), userId%1000)
	req.OrderSn = orderSn

	// 2. 添加子事务：调用商品服务扣库存
	// 注意：这里需要传入商品服务的 gRPC 全路径名
	for _, item := range req.Items {
		saga.Add(pbProduct.Product_ReduceStock_FullMethodName,
			pbProduct.Product_CompensateStock_FullMethodName,
			&pbProduct.ReduceStockReq{OrderSn: orderSn, SkuId: int64(item.SkuId), Num: item.Num})
	}
	// 3. 添加子事务：调用本服务的创建订单接口
	saga.Add(pb.Order_CreateOrderInner_FullMethodName,
		pb.Order_CancelOrder_FullMethodName,
		req)

	// 4. 提交事务
	err := saga.Submit()
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderReply{OrderSn: gid}, nil
}

func (s *OrderService) CreateOrderInner(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderReply, error) {
	// 1. 获取当前用户 (从 JWT 中间件存入的 Context 获取)
	userId := ctx.Value("userID").(uint64)

	// 这个需要计算一下
	var totalPayAmount float64

	// 2. 执行业务逻辑
	// C. 构建订单实体
	orderEntity := &model.Order{
		OrderSn:         req.OrderSn,
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
	err := s.orderRepo.CreateOrder(ctx, orderEntity, orderItems)

	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderReply{
		OrderSn:     req.OrderSn,
		TotalAmount: totalPayAmount,
	}, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderReply, error) {
	err := s.orderRepo.CancelOrder(ctx, req.OrderSn)
	return &pb.CreateOrderReply{}, err
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
