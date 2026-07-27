package middle

import (
	"strings"

	"github.com/gin-gonic/gin"
	"trashcan/ginServer/api/common"
	"trashcan/utils"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.FailWithAuthority(c)
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.FailWithAuthority(c)
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			common.FailWithAuthority(c)
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(token)
		if err != nil {
			common.FailWithAuthority(c)
			c.Abort()
			return
		}

		// 将用户信息存入context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("isAdmin", claims.IsAdmin)

		c.Next()
	}
}
