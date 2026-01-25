// internal/transport/grpc/user.go
package grpc

import (
	"context"
	pb "helloworld-go/api/user/v1"
	service "helloworld-go/internal/service/user"
)

type UserServer struct {
	pb.UnimplementedUserServer
	svc *service.UserService
}

func NewUserServer(svc *service.UserService) *UserServer {
	return &UserServer{svc: svc}
}

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserReply, error) {
	u, err := s.svc.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: u.Id, Name: u.Name}, nil
}

func (s *UserServer) GetByID(ctx context.Context, req *pb.SingleIDRequest) (*pb.UserReply, error) {
	u, err := s.svc.GetByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: u.Id, Name: u.Name}, nil
}
