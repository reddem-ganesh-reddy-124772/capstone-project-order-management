package main

import (
	"log"
	"net/http"

	"order-management/config"
	"order-management/services/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load application settings from .env.
	cfg := config.Load()

	// Use production logging behavior when configured for production.
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create the Gin router with logging and panic recovery.
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Health endpoint for checking whether the service is running.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"service": "auth-service",
		})
	})

	// Register authentication routes.
	auth.RegisterRoutes(router, cfg)

	// Start the Auth Service using AUTH_SERVICE_PORT from .env.
	address := ":" + cfg.AuthServicePort
	log.Printf("Auth Service starting on port %s", cfg.AuthServicePort)

	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start Auth Service: %v", err)
	}
}
