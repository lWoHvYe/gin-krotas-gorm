// internal/infrastructure/persistence/user_repo.go
package persistence

import (
	"context"
	"helloworld-go/internal/biz/user"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) user.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Save(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepo) FindByID(ctx context.Context, id int64) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) UpdateRoleByID(ctx context.Context, t user.User) error {
	return r.db.WithContext(ctx).Where("id = ?", t.ID).Update("role_name", t.RoleName).Error
}
