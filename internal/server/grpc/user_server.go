// internal/transport/grpc/user_server.go
package grpc

import (
	"context"
	pb "helloworld-go/api/user-service/v1"
	"helloworld-go/internal/service/user"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	svc *user.Service
}

func NewServer(svc *user.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserReply, error) {
	u, err := s.svc.Register(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: u.ID, Name: u.Name}, nil
}

func (s *Server) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.UserReply, error) {
	u, err := s.svc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: u.ID, Name: u.Name}, nil
}
