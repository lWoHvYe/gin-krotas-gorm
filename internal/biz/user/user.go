// internal/domain/user/user.go
package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	RoleName string
}
