package middleware

import (
	"demos/pkg/response"
	"demos/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, response.Error(response.StatusUnauthorized, "Authorization header is required"))
			c.Abort()
			return
		}

		// 从Bearer格式中提取Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, response.Error(response.StatusUnauthorized, "Invalid Authorization header format"))
			c.Abort()
			return
		}

		// 解析并验证Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(401, response.Error(response.StatusUnauthorized, "Invalid token"))
			c.Abort()
			return
		}

		// 将用户ID和角色信息存入上下文
		// 这样在后续的Handler中可以直接从上下文中获取
		c.Set("user_id", claims.UserId)
		c.Set("role", claims.Role)
		c.Next()
	}
}
