// tests/user_service_test.go
package tests

import (
	"context"
	pb "helloworld-go/api/user/v1"
	"helloworld-go/internal/biz"
	"helloworld-go/internal/biz/user"
	"helloworld-go/internal/service"
	"testing"
)

type MockUserRepo struct {
	biz.RepositoryImpl[user.User]
}

func (r *MockUserRepo) UpdateRoleByID(ctx context.Context, t user.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *MockUserRepo) Save(ctx context.Context, u *user.User) error {
	u.ID = 1
	return nil
}

func (r *MockUserRepo) FindByID(ctx context.Context, id any) (*user.User, error) {
	return &user.User{ID: id, Name: "mock"}, nil
}

func TestRegister(t *testing.T) {
	repo := &MockUserRepo{}
	svc := service.NewUserService(repo)

	u, err := svc.Register(context.Background(), &pb.RegisterRequest{Name: "Alice"})
	if err != nil || u.Id != 1 {
		t.Fatal("Register failed")
	}
}
