// internal/transport/grpc/user_server.go
package grpc

import (
	"context"
	pb "helloworld-go/api/user/v1"
	service "helloworld-go/internal/service/user"

	grpc "google.golang.org/grpc"
)

type UserServer struct {
	pb.UnimplementedUserServer
	svc         *service.UserService
	ServiceDesc grpc.ServiceDesc
}

func NewUserServer(svc *service.UserService) *UserServer {
	return &UserServer{svc: svc, ServiceDesc: pb.User_ServiceDesc}
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
