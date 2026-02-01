package middleware

import (
	"helloworld-go/internal/pkg/utils"
	"net/http"
	"strconv"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

func CasbinMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 注意：sub 必须是字符串类型，对应 Casbin 策略中的主体
		// 1. 获取当前用户标识（假设你之前在 JWT 中间件里设置了 "userID"）
		val, _ := c.Get("claims")
		claims := val.(*utils.CustomClaims)
		sub := "u:" + strconv.FormatUint(claims.UID, 10) // 或者使用用户唯一 ID

		// 2. 获取请求的资源 (对象) 和 操作 (动作)
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 3. 校验权限
		ok, err := e.Enforce(sub, obj, act)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "权限校验出错"})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"msg": "没有访问权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
