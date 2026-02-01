// internal/infrastructure/persistence/user.go
package persistence

import (
	"context"
	"helloworld-go/internal/biz/model"
	biz "helloworld-go/internal/biz/user"
	"strconv"

	"github.com/casbin/casbin/v3"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
	e  *casbin.Enforcer
}

func NewUserRepo(db *gorm.DB, e *casbin.Enforcer) *UserRepo {
	return &UserRepo{db: db, e: e}
}

func (r *UserRepo) Create(ctx context.Context, u *model.User) error {
	q := biz.Use(r.db)
	return q.User.WithContext(ctx).Create(u)
}

func (r *UserRepo) FindByID(ctx context.Context, id int32) (*model.User, error) {
	q := biz.Use(r.db)
	return q.User.
		WithContext(ctx).
		Preload(q.User.UserRoles.Where(q.UserRole.Status.Eq("1"))).
		Preload(q.User.UserRoles.Role).
		Where(q.User.ID.Eq(id)).
		First()
}

func (r *UserRepo) UpdateRoleByID(ctx context.Context, userId int32, roleIds []uint64) error {
	q := biz.Use(r.db)
	// 开启事务
	return q.Transaction(func(tx *biz.Query) error {
		// 更新 g p表（先从g表移除该用户）
		gV0 := "u:" + strconv.FormatUint(uint64(userId), 10)
		if _, err := r.e.DeleteUser(gV0); err != nil {
			return err
		}
		// 在 tx 上做所有 DB 操作
		if _, err := tx.UserRole.WithContext(ctx).Where(tx.UserRole.UserID.Eq(userId)).Delete(); err != nil {
			return err
		}
		userRoles := make([]*model.UserRole, 0)
		for _, roleId := range roleIds {
			userRoles = append(userRoles, &model.UserRole{UserID: userId, RoleID: int32(roleId), Status: "1"})
			// 再根据新的权限添加
			if _, err := r.e.AddGroupingPolicy(gV0, "r:"+strconv.FormatUint(roleId, 10)); err != nil {
				return err
			}
		}
		return tx.UserRole.WithContext(ctx).Create(userRoles...)
	})
}

func (r *UserRepo) UpdateByID(ctx context.Context, u *model.User) error {
	q := biz.Use(r.db)
	_, err := q.User.WithContext(ctx).Updates(u)
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id int32) error {
	q := biz.Use(r.db)
	return q.Transaction(func(tx *biz.Query) error {
		if _, err := tx.User.WithContext(ctx).Where(tx.User.ID.Eq(id)).Delete(); err != nil {
			return err
		}
		if _, err := tx.UserRole.WithContext(ctx).Where(tx.UserRole.UserID.Eq(id)).Delete(); err != nil {
			return err
		}
		return nil
	})
}

func (r *UserRepo) FindByName(ctx context.Context, name string) (*model.User, error) {
	q := biz.Use(r.db)
	u := q.User
	ur := q.UserRole
	return q.User.
		WithContext(ctx).
		Preload(u.UserRoles.On(ur.Status.Eq("1")), u.UserRoles.Role).
		Where(q.User.Name.Eq(name)).
		First()
}
