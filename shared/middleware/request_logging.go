package middleware

import (
	"time"

	"order-management/shared/logger"

	"github.com/gin-gonic/gin"
)

// RequestLogging records useful details about every HTTP request.
// Register after RequestID() and CorrelationID().
func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		log := logger.Get()

		log.Info(
			"HTTP request completed",
			"request_id", GetRequestID(c),
			"correlation_id", GetCorrelationID(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
