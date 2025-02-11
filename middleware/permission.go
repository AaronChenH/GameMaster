package middleware

import (
	"strings"

	"game-admin/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var permissionMap = map[string][]string{
	"admin":   {"player_info", "give_item", "execute_script", "manage_users"},
	"service": {"player_info", "give_item"},
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "未授权访问"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.SecretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

func CheckPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		permissions := permissionMap[role]

		hasPermission := false
		for _, p := range permissions {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(403, gin.H{"error": "无权限访问"})
			c.Abort()
			return
		}

		c.Next()
	}
}
