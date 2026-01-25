package service

import (
	"context"
	"helloworld-go/internal/biz/user"

	pb "helloworld-go/api/user/v1"
)

type UserService struct {
	pb.UnimplementedUserServer
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserReply, error) {
	u := &user.User{Name: req.GetName()}
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: u.ID, Name: u.Name, RoleName: "admin"}, nil
}
func (s *UserService) GetByID(ctx context.Context, req *pb.SingleIDRequest) (*pb.UserReply, error) {
	u, err := s.repo.FindByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: u.ID, Name: u.Name}, nil
}
func (s *UserService) UpdateRoleByID(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.SimpleResponse, error) {
	u := user.User{ID: req.Id, RoleName: req.RoleName}
	err := s.repo.UpdateRoleByID(ctx, u)
	return &pb.SimpleResponse{}, err
}
