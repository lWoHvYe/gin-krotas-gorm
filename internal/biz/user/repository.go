// internal/domain/user/repository.go
package user

import "context"

// 泛型 Repository 接口
type Repository[T any] interface {
	Save(ctx context.Context, t *T) error
	FindByID(ctx context.Context, id int64) (*T, error)
}

type UserRepository interface {
	Repository[User]
}
