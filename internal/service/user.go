package service

import (
	"context"
	"errors"
	pb "gin-kratos-gorm/api/user/v1"
	"gin-kratos-gorm/internal/biz/model"
	"gin-kratos-gorm/internal/data/persistence"
	"strconv"
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
	return &pb.UserReply{Id: uint64(u.ID), Name: u.Name, RoleId: []uint64{1, 2}}, nil
}
func (s *UserService) GetByID(ctx context.Context, req *pb.SingleIDRequest) (*pb.UserReply, error) {
	u, err := s.repo.FindByID(ctx, int32(req.GetId()))
	if err != nil {
		return nil, err
	}
	roleIds := make([]uint64, 0)
	for _, userRole := range u.UserRoles {
		roleIds = append(roleIds, uint64(userRole.RoleID))
	}
	return &pb.UserReply{Id: uint64(u.ID), Name: u.Name, RoleId: roleIds}, nil
}
func (s *UserService) UpdateRoleByID(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.SimpleResponse, error) {
	err := s.repo.UpdateRoleByID(ctx, int32(req.Id), req.RoleId)
	return &pb.SimpleResponse{Msg: strconv.FormatUint(req.RoleId[0], 10)}, err
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.UserReply, error) {
	if req.Password == "" {
		return nil, errors.New("password is empty")
	}
	u, err := s.repo.FindByName(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	roleIds := make([]uint64, 0)
	for _, userRole := range u.UserRoles {
		roleIds = append(roleIds, uint64(userRole.RoleID))
	}
	return &pb.UserReply{Id: uint64(u.ID), Name: u.Name, RoleId: roleIds}, nil
}
