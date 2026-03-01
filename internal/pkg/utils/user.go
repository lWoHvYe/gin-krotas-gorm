package utils

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
)

func GetLoginUserId(ctx context.Context) uint64 {
	// Kratos grpc auth middleware
	if claims, ok := jwt.FromContext(ctx); ok {
		// 强制断言为你定义的类型
		c := claims.(CustomClaims)
		return c.UID
	}

	// or http handler
	// 1. 获取当前用户 (从 JWT 中间件存入的 Context 获取)
	return ctx.Value("userID").(uint64)
}
