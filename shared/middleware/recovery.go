package middleware

import (
	"net/http"

	"order-management/shared/logger"

	"github.com/gin-gonic/gin"
)

// Recovery catches unexpected panics and keeps the service running.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log := logger.Get()

		// Log the panic along with the request's correlation ID.
		log.Error(
			"Unexpected panic recovered",
			"correlation_id", GetCorrelationID(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"error", recovered,
		)

		// Stop processing the request and return a safe error response.
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":          "internal server error",
			"correlation_id": GetCorrelationID(c),
		})
	})
}
