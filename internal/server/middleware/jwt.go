package middleware

import (
	"errors"
	"helloworld-go/internal/pkg/utils" // 替换为你的路径
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWT Auth middleware
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check if the header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Use 'Bearer TOKEN'"})
			return
		}

		token := parts[1]
		// Parse and validate the token
		j := utils.NewJWT()
		// 解析 token
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "授权已过期"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
			return
		}

		// 将解析出的信息存入上下文，方便后续 Service 层或 RBAC 中间件使用
		c.Set("claims", claims)
		c.Set("roleID", claims.AuthorityId) // 供 Casbin 使用
		c.Set("userID", claims.UID)

		c.Next()
	}
}
