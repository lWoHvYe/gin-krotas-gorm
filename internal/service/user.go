package service

import (
	"context"
	pb "helloworld-go/api/user/v1"
	"helloworld-go/internal/biz/model"
	"helloworld-go/internal/data/persistence"
)

type UserService struct {
	pb.UnimplementedUserServer
	repo *persistence.UserRepo
}

func NewUserService(repo *persistence.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserReply, error) {
	u := &model.User{Name: req.GetName()}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: uint64(u.ID), Name: u.Name, RoleName: "admin"}, nil
}
func (s *UserService) GetByID(ctx context.Context, req *pb.SingleIDRequest) (*pb.UserReply, error) {
	u, err := s.repo.FindByID(ctx, int32(req.GetId()))
	if err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: uint64(u.ID), Name: u.Name}, nil
}
func (s *UserService) UpdateRoleByID(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.SimpleResponse, error) {
	u := &model.User{ID: int32(req.Id), RoleName: req.RoleName}
	err := s.repo.UpdateRoleByID(ctx, u)
	return &pb.SimpleResponse{Msg: req.RoleName}, err
}
