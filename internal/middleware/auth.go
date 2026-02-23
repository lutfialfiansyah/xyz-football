package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"xyz-football-api/internal/modules/admin"
	"xyz-football-api/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header is required", "Header otorisasi wajib diisi")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}", "Format header otorisasi harus Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &admin.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				response.Error(c, http.StatusUnauthorized, "Token is expired", "Token sudah kedaluwarsa")
				c.Abort()
				return
			}
			response.Error(c, http.StatusUnauthorized, "Invalid token", "Token tidak valid")
			c.Abort()
			return
		}

		if !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Invalid token", "Token tidak valid")
			c.Abort()
			return
		}

		// Set Context values for next handlers
		c.Set("id", claims.ID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware Example usage: router.Use(middleware.RoleMiddleware("admin"))
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusForbidden, "Role not found in context", "Role tidak ditemukan")
			c.Abort()
			return
		}

		if roleStr, ok := role.(string); ok && roleStr != requiredRole {
			response.Error(c, http.StatusForbidden, "You don't have permission to access this resource", "Anda tidak memiliki izin untuk mengakses resource ini")
			c.Abort()
			return
		}

		c.Next()
	}
}
