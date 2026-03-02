package middleWare

import (
	"LifeNavigator/backend/pkg/jwt"
	"LifeNavigator/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Code(c, jwt.StatusMissedToken)
			c.Abort()
			return
		}

		claims, err := jwt.VerifyAccessToken(authHeader)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user", jwt.AuthUser{
			UserID:     claims.UserID,
			Department: claims.Department,
			Role:       claims.Role,
		})
		c.Next()
	}
}
