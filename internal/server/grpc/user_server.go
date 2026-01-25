// internal/transport/grpc/user_server.go
package grpc

import (
	"context"
	pb "helloworld-go/api/user/v1"
	service "helloworld-go/internal/service/user"

	grpc "google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	svc         *service.Service
	ServiceDesc grpc.ServiceDesc
}

func NewServer(svc *service.Service) *Server {
	return &Server{svc: svc, ServiceDesc: pb.UserService_ServiceDesc}
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserReply, error) {
	u, err := s.svc.Register(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: u.ID, Name: u.Name}, nil
}

func (s *Server) GetByID(ctx context.Context, req *pb.SingleIDRequest) (*pb.UserReply, error) {
	u, err := s.svc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: u.ID, Name: u.Name}, nil
}
