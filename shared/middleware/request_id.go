
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

const requestIDHeader = "X-Request-ID"

type requestIDContextKey struct{}

// RequestID assigns an ID to each incoming HTTP request.
// It preserves an incoming ID or generates one if missing.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(
			c.GetHeader(requestIDHeader),
		)

		if requestID == "" || len(requestID) > 128 {
			requestID = newRequestID()
		}

		c.Header(requestIDHeader, requestID)
		c.Set(requestIDHeader, requestID)

		ctx := context.WithValue(
			c.Request.Context(),
			requestIDContextKey{},
			requestID,
		)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetRequestID retrieves the ID from a Gin request.
func GetRequestID(c *gin.Context) string {
	value, exists := c.Get(requestIDHeader)
	if !exists {
		return ""
	}

	id, ok := value.(string)
	if !ok {
		return ""
	}

	return id
}

// RequestIDFromContext retrieves the ID from a Go context.
func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDContextKey{}).(string)
	return id
}

func newRequestID() string {
	var value [16]byte

	if _, err := rand.Read(value[:]); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}

	encoded := hex.EncodeToString(value[:])

	return encoded[:8] + "-" +
		encoded[8:12] + "-" +
		encoded[12:16] + "-" +
		encoded[16:20] + "-" +
		encoded[20:]
}
