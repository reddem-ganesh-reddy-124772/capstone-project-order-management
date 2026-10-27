package auth

import (
	"order-management/config"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers authentication API endpoints.
func RegisterRoutes(router *gin.Engine, cfg config.Config) {
	controller := NewController(cfg)

	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/login", controller.Login)
	}
}
