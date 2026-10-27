package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"order-management/config"
	"order-management/shared/middleware"
)

// Client is used for HTTP communication between services.
type Client struct {
	httpClient *http.Client
	maxRetries int
}

// New creates an HTTP client using the application configuration.
func New(cfg config.Config) *Client {
	timeout := time.Duration(cfg.HTTPClientTimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	retries := cfg.HTTPClientMaxRetries
	if retries < 0 {
		retries = 0
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		maxRetries: retries,
	}
}

// Do sends an HTTP request with timeout, retries, and tracking IDs.
func (c *Client) Do(
	ctx context.Context,
	method string,
	url string,
	body []byte,
	headers map[string]string,
) (*http.Response, error) {
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(
			ctx,
			method,
			url,
			bytes.NewReader(body),
		)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		// Preserve caller-provided headers.
		for name, value := range headers {
			req.Header.Set(name, value)
		}

		// Propagate the IDs from the incoming request context.
		requestID := middleware.RequestIDFromContext(ctx)
		if requestID != "" {
			req.Header.Set("X-Request-ID", requestID)
		}

		correlationID := middleware.FromContext(ctx)
		if correlationID != "" {
			req.Header.Set("X-Correlation-ID", correlationID)
		}

		resp, err := c.httpClient.Do(req)

		// Retry network errors.
		if err != nil {
			if attempt == c.maxRetries {
				return nil, fmt.Errorf("send request: %w", err)
			}

			time.Sleep(time.Duration(attempt+1) * 200 * time.Millisecond)
			continue
		}

		// Retry server errors (HTTP 5xx).
		if resp.StatusCode >= 500 && attempt < c.maxRetries {
			resp.Body.Close()
			time.Sleep(time.Duration(attempt+1) * 200 * time.Millisecond)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed")
}
