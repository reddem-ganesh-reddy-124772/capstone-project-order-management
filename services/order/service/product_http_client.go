package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"order-management/config"
	"order-management/shared/discovery"
	"order-management/shared/httpclient"
	"order-management/shared/middleware"
)

// httpProductClient communicates with Product Service through REST APIs.
type httpProductClient struct {
	baseURL string
	client  *httpclient.Client
	breaker *httpclient.CircuitBreaker
}

// NewHTTPProductClient creates the Product Service client using Consul discovery.
// If Consul discovery fails, it uses the URL from configuration.
func NewHTTPProductClient(cfg config.Config) ProductClient {
	productURL := strings.TrimRight(cfg.ProductServiceURL, "/")

	consulClient := discovery.New(cfg)

	discoveredURL, err := consulClient.DiscoverService("product-service")
	if err == nil && strings.TrimSpace(discoveredURL) != "" {
		productURL = strings.TrimRight(discoveredURL, "/")
	}

	return &httpProductClient{
		baseURL: productURL,
		client:  httpclient.New(cfg),
		breaker: httpclient.NewCircuitBreaker(cfg),
	}
}

// GetProduct retrieves product details from Product Service.
func (c *httpProductClient) GetProduct(
	ctx context.Context,
	productID uint,
) (*ProductInfo, error) {
	if productID == 0 {
		return nil, errors.New("product ID must be greater than zero")
	}

	url := c.baseURL + "/api/products/" +
		strconv.FormatUint(uint64(productID), 10)

	body, status, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf(
			"Product Service returned HTTP %d: %s",
			status,
			string(body),
		)
	}

	var product ProductInfo
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, fmt.Errorf("decode product response: %w", err)
	}

	if product.ID == 0 {
		return nil, errors.New("Product Service returned an invalid product")
	}

	return &product, nil
}

// ReserveStock reserves stock for an order item.
func (c *httpProductClient) ReserveStock(
	ctx context.Context,
	productID uint,
	reservationID string,
	quantity int,
) (*ProductReservation, error) {
	if productID == 0 {
		return nil, errors.New("product ID must be greater than zero")
	}
	if strings.TrimSpace(reservationID) == "" {
		return nil, errors.New("reservation ID is required")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	requestBody := map[string]interface{}{
		"reservation_id": reservationID,
		"product_id":     productID,
		"quantity":       quantity,
	}

	body, status, err := c.doRequest(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/inventory/reservations",
		requestBody,
	)
	if err != nil {
		return nil, err
	}

	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"stock reservation failed with HTTP %d: %s",
			status,
			string(body),
		)
	}

	reservation := &ProductReservation{
		ReservationID: reservationID,
		ProductID:     productID,
		Quantity:      quantity,
		Status:        "reserved",
	}

	var response ProductReservation
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("decode reservation response: %w", err)
	}

	if response.ReservationID != "" {
		reservation.ReservationID = response.ReservationID
	}
	if response.ProductID != 0 {
		reservation.ProductID = response.ProductID
	}
	if response.Quantity != 0 {
		reservation.Quantity = response.Quantity
	}
	if response.Status != "" {
		reservation.Status = response.Status
	}

	return reservation, nil
}

// ReleaseStock releases an existing reservation.
func (c *httpProductClient) ReleaseStock(
	ctx context.Context,
	reservationID string,
) error {
	if strings.TrimSpace(reservationID) == "" {
		return errors.New("reservation ID is required")
	}

	body, status, err := c.doRequest(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/inventory/reservations/release",
		map[string]string{"reservation_id": reservationID},
	)
	if err != nil {
		return err
	}

	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		return fmt.Errorf(
			"stock release failed with HTTP %d: %s",
			status,
			string(body),
		)
	}

	return nil
}

// doRequest uses the shared HTTP client and circuit breaker.
func (c *httpProductClient) doRequest(
	ctx context.Context,
	method string,
	url string,
	requestData interface{},
) ([]byte, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	// Reject the call when the circuit is open.
	if err := c.breaker.Allow(); err != nil {
		return nil, 0, err
	}

	var requestBody []byte
	var err error

	if requestData != nil {
		requestBody, err = json.Marshal(requestData)
		if err != nil {
			return nil, 0, fmt.Errorf("encode HTTP request: %w", err)
		}
	}

	headers := map[string]string{
		"Accept": "application/json",
	}
	if requestBody != nil {
		headers["Content-Type"] = "application/json"
	}

	// Forward the existing correlation ID to Product Service.
	// Do not generate a new correlation ID in this HTTP client.
	if correlationID := middleware.FromContext(ctx); correlationID != "" {
		headers["X-Correlation-ID"] = correlationID
	}

	// The shared client handles timeout, retry, and HTTP requests.
	response, err := c.client.Do(
		ctx,
		method,
		url,
		requestBody,
		headers,
	)
	if err != nil {
		c.breaker.RecordFailure()
		return nil, 0, fmt.Errorf("call Product Service: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		c.breaker.RecordFailure()
		return nil, 0, fmt.Errorf("read Product Service response: %w", err)
	}

	// Count temporary/server failures toward opening the circuit.
	if response.StatusCode == http.StatusRequestTimeout ||
		response.StatusCode == http.StatusTooManyRequests ||
		response.StatusCode >= http.StatusInternalServerError {
		c.breaker.RecordFailure()
	} else {
		c.breaker.RecordSuccess()
	}

	return body, response.StatusCode, nil
}
