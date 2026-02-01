package service

import (
	"context"

	pb "helloworld-go/api/order/v1"
)

type OrderService struct {
	pb.UnimplementedOrderServer
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderReply, error) {
	return &pb.CreateOrderReply{}, nil
}
func (s *OrderService) GetOrderDetail(ctx context.Context, req *pb.GetOrderDetailReq) (*pb.GetOrderDetailReply, error) {
	return &pb.GetOrderDetailReply{}, nil
}
