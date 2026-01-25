// tests/user_service_test.go
package tests

import (
	"context"
	"helloworld-go/internal/biz/user"
	service "helloworld-go/internal/service/user"
	"testing"
)

type MockUserRepo struct{}

func (r *MockUserRepo) Save(ctx context.Context, u *user.User) error {
	u.ID = 1
	return nil
}

func (r *MockUserRepo) FindByID(ctx context.Context, id int64) (*user.User, error) {
	return &user.User{ID: id, Name: "mock"}, nil
}

func TestRegister(t *testing.T) {
	repo := &MockUserRepo{}
	svc := service.NewService(repo)

	u, err := svc.Register(context.Background(), "Alice")
	if err != nil || u.ID != 1 {
		t.Fatal("Register failed")
	}
}
