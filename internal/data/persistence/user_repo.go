// internal/infrastructure/persistence/user_repo.go
package persistence

import (
	"context"
	"helloworld-go/internal/biz"
	"helloworld-go/internal/biz/user"

	"gorm.io/gorm"
)

type UserRepo struct {
	biz.RepositoryImpl[user.User]
}

func NewUserRepo(db *gorm.DB) user.UserRepository {
	return &UserRepo{biz.RepositoryImpl[user.User]{DB: db}}
}

func (r *UserRepo) UpdateRoleByID(ctx context.Context, t user.User) error {
	_, err := gorm.G[user.User](r.DB).Where("id = ?", t.ID).Update(ctx, "role_name", t.RoleName)
	return err
}
