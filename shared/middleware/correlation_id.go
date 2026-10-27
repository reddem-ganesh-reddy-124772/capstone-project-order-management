package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const correlationIDHeader = "X-Correlation-ID"

type correlationIDContextKey struct{}

// CorrelationID adds a correlation ID to every incoming HTTP request.
// It preserves a valid incoming ID or generates a new one.
func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := strings.TrimSpace(
			c.GetHeader(correlationIDHeader),
		)

		// Avoid accepting excessively long IDs from incoming requests.
		if correlationID == "" || len(correlationID) > 128 {
			correlationID = newCorrelationID()
		}

		// Return the ID to the client in the HTTP response.
		c.Header(correlationIDHeader, correlationID)
c.Request.Header.Set(correlationIDHeader, correlationID)

		// Store the ID in Gin's request context.
		c.Set(correlationIDHeader, correlationID)

		// Also attach the ID to the standard Go context so that
		// downstream HTTP clients can propagate it to other services.
		requestContext := context.WithValue(
			c.Request.Context(),
			correlationIDContextKey{},
			correlationID,
		)
		c.Request = c.Request.WithContext(requestContext)

		c.Next()
	}
}

// GetCorrelationID retrieves the ID from a Gin request.
func GetCorrelationID(c *gin.Context) string {
	value, exists := c.Get(correlationIDHeader)
	if !exists {
		return ""
	}

	correlationID, ok := value.(string)
	if !ok {
		return ""
	}

	return correlationID
}

// FromContext retrieves the ID from a standard Go context.
func FromContext(ctx context.Context) string {
	value, ok := ctx.Value(correlationIDContextKey{}).(string)
	if !ok {
		return ""
	}

	return value
}

// newCorrelationID creates a random identifier without requiring
// an additional UUID library.
func newCorrelationID() string {
	var value [16]byte

	if _, err := rand.Read(value[:]); err != nil {
		// Keep request handling available if the random source fails.
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}

	encoded := hex.EncodeToString(value[:])

	return encoded[:8] + "-" +
		encoded[8:12] + "-" +
		encoded[12:16] + "-" +
		encoded[16:20] + "-" +
		encoded[20:]
}
