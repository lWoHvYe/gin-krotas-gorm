// internal/domain/user/repository.go
package user

import (
	"context"
	"helloworld-go/internal/biz"
)

type UserRepository interface {
	biz.Repository[User]
	UpdateRoleByID(ctx context.Context, t User) error
}
