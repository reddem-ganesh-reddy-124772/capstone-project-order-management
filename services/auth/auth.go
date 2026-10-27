package auth

import (
	"errors"
	"fmt"
	"time"

	"order-management/config"

	"github.com/golang-jwt/jwt/v5"
)

// User represents an authenticated user.
// Passwords are never included in API responses.
type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

// LoginRequest contains the credentials sent by the client.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse contains the JWT returned after successful login.
type LoginResponse struct {
	Token string `json:"token"`
	Type  string `json:"token_type"`
	Role  string `json:"role"`
}

// Authenticate validates credentials loaded from configuration.
func Authenticate(cfg config.Config, username, password string) (User, error) {
	if username == "" || password == "" {
		return User{}, errors.New("username and password are required")
	}

	if cfg.AuthAdminUsername != "" &&
		cfg.AuthAdminPassword != "" &&
		username == cfg.AuthAdminUsername &&
		password == cfg.AuthAdminPassword {
		return User{
			Username: username,
			Role:     "admin",
		}, nil
	}

	if cfg.AuthUserUsername != "" &&
		cfg.AuthUserPassword != "" &&
		username == cfg.AuthUserUsername &&
		password == cfg.AuthUserPassword {
		return User{
			Username: username,
			Role:     "user",
		}, nil
	}

	return User{}, errors.New("invalid username or password")
}

// GenerateToken creates a signed JWT for an authenticated user.
func GenerateToken(cfg config.Config, user User) (string, error) {
	if cfg.JWTSecret == "" ||
		cfg.JWTSecret == "replace_with_a_long_random_secret" {
		return "", errors.New("configure a strong JWT_SECRET in .env")
	}

	if user.Username == "" || user.Role == "" {
		return "", errors.New("valid user details are required")
	}

	expirationMinutes := cfg.JWTExpirationMinutes
	if expirationMinutes <= 0 {
		return "", errors.New("JWT_EXPIRATION_MINUTES must be greater than zero")
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   user.Username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expirationMinutes) * time.Minute)),
		Issuer:    "order-management-auth",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  claims.Subject,
		"role": user.Role,
		"iat":  claims.IssuedAt.Unix(),
		"exp":  claims.ExpiresAt.Unix(),
		"iss":  claims.Issuer,
	})

	signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("generate JWT: %w", err)
	}

	return signedToken, nil
}
