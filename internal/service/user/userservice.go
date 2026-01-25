package service

import (
	"context"
	"helloworld-go/internal/biz/user"

	pb "helloworld-go/api/user/v1"
)

type UserServiceService struct {
	pb.UnimplementedUserServiceServer
	repo user.UserRepository
}

func NewUserServiceService(repo user.UserRepository) *UserServiceService {
	return &UserServiceService{repo: repo}
}

func (s *UserServiceService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserReply, error) {
	u := &user.User{Name: req.GetName()}
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: u.ID, Name: u.Name, RoleName: "admin"}, nil
}
func (s *UserServiceService) GetByID(ctx context.Context, req *pb.SingleIDRequest) (*pb.UserReply, error) {
	u, err := s.repo.FindByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: u.ID, Name: u.Name}, nil
}
func (s *UserServiceService) UpdateRoleByID(ctx context.Context, req *pb.SingleIDRequest) (*pb.SimpleResponse, error) {
	return &pb.SimpleResponse{}, nil
}
