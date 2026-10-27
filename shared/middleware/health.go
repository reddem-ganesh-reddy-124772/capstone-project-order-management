package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse contains the status returned by the health endpoint.
type HealthResponse struct {
	// Status tells us whether the service is running.
	Status string `json:"status"`
}

// HealthCheck returns a handler that reports whether the service is running.
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Return HTTP 200 to indicate that the service is responding.
		c.JSON(http.StatusOK, HealthResponse{
			Status: "UP",
		})
	}
}
