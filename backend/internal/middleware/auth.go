package middleware

import (
	"net/http"
	"strings"

	"github.com/DavoodHakimi/warehouse-app/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("role", claims["role"])
		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Set("company_id", claims["company_id"])

		c.Next()
	}
}
