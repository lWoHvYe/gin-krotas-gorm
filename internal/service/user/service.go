// internal/application/user/service.go
package user

import (
	"context"
	"helloworld-go/internal/biz/user"
)

type Service struct {
	repo user.UserRepository
}

func NewService(repo user.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, name string) (*user.User, error) {
	u := &user.User{Name: name}
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*user.User, error) {
	return s.repo.FindByID(ctx, id)
}
