package middleware

import (
	"net/http"
	"strings"

	"gin_gorm_demo/auth"
	"gin_gorm_demo/config"
	"gin_gorm_demo/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	cfg := config.Load()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, http.StatusUnauthorized, 40101, "missing token")
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(c, http.StatusUnauthorized, 40102, "invalid token format")
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(parts[1], cfg.JWTSecret)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, 40103, "invalid token")
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
