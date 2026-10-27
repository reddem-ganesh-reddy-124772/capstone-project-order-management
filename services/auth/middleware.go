package auth

import (
	"net/http"
	"strings"

	"order-management/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens on protected routes.
func AuthMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization token is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "use Authorization: Bearer <token>",
			})
			return
		}

		if cfg.JWTSecret == "" ||
			cfg.JWTSecret == "replace_with_a_long_random_secret" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "JWT secret is not configured",
			})
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(
			parts[1],
			claims,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.JWTSecret), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithIssuer("order-management-auth"),
			jwt.WithExpirationRequired(),
		)

		if err != nil || token == nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		username, usernameOK := claims["sub"].(string)
		role, roleOK := claims["role"].(string)

		if !usernameOK || username == "" ||
			!roleOK || (role != "admin" && role != "user") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			return
		}

		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}

// RequireRole restricts an endpoint to the specified role.
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you do not have permission to access this resource",
			})
			return
		}

		c.Next()
	}
}
