package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/DavoodHakimi/warehouse-app/internal/auth"
	"github.com/DavoodHakimi/warehouse-app/internal/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RBAC struct {
	db *gorm.DB
}

func NewRBAC(db *gorm.DB) *RBAC {
	return &RBAC{db: db}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			slog.Warn("auth - missing authorization header", "ip", c.ClientIP(), "path", c.FullPath())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			slog.Warn("auth - invalid token", "ip", c.ClientIP(), "path", c.FullPath(), "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		toInt := func(val interface{}) int {
			switch v := val.(type) {
			case float64:
				return int(v)
			case int:
				return v
			case string:
				i, _ := strconv.Atoi(v)
				return i
			default:
				return 0
			}
		}

		c.Set("role", claims["role"])
		c.Set("user_id", toInt(claims["user_id"]))
		c.Set("username", claims["username"])
		c.Set("company_id", toInt(claims["company_id"]))

		slog.Info("auth - authenticated", "user_id", toInt(claims["user_id"]), "path", c.FullPath())
		c.Next()
	}
}

func (r *RBAC) RBACMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")

		if !database.HasAccess(r.db, userID, permission) {
			slog.Warn("rbac - permission denied", "user_id", userID, "permission", permission, "path", c.FullPath())
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}
		slog.Info("rbac - allowed", "user_id", userID, "permission", permission, "path", c.FullPath())
		c.Next()
	}
}
