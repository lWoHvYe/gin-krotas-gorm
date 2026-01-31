// internal/infrastructure/persistence/user.go
package persistence

import (
	"context"
	"helloworld-go/internal/biz/model"
	biz "helloworld-go/internal/biz/user"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *model.User) error {
	q := biz.Use(r.db)
	return q.User.WithContext(ctx).Create(u)
}

func (r *UserRepo) FindByID(ctx context.Context, id int32) (*model.User, error) {
	q := biz.Use(r.db)
	return q.User.WithContext(ctx).Where(q.User.ID.Eq(id)).First()
}

func (r *UserRepo) UpdateRoleByID(ctx context.Context, u *model.User) error {
	q := biz.Use(r.db)
	_, err := q.User.WithContext(ctx).Where(q.User.ID.Eq(u.ID)).Update(q.User.RoleName, u.RoleName)
	return err
}

func (r *UserRepo) UpdateByID(ctx context.Context, u *model.User) error {
	q := biz.Use(r.db)
	_, err := q.User.WithContext(ctx).Updates(u)
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id int32) error {
	q := biz.Use(r.db)
	_, err := q.User.WithContext(ctx).Where(q.User.ID.Eq(id)).Delete()
	return err
}
