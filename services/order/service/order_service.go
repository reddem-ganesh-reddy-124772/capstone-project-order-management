package service

import (
	"context"

	"order-management/services/order/dto"
)

// OrderService defines the business operations available for orders.
type OrderService interface {
	// CreateOrder validates the request and creates an order.
	// The context carries the correlation ID and request cancellation.
	// Product details and prices must come from Product Service.
	CreateOrder(ctx context.Context, request dto.CreateOrderRequest) (*dto.OrderResponse, error)

	// GetOrderByID retrieves an order and its items using its ID.
	GetOrderByID(id uint) (*dto.OrderResponse, error)

	// GetAllOrders retrieves the available orders.
	GetAllOrders() ([]dto.OrderResponse, error)

	// ConfirmOrder confirms an order only after required stock
	// has been reserved successfully.
	ConfirmOrder(id uint) (*dto.OrderResponse, error)

	// CancelOrder cancels an eligible order and releases its
	// reserved stock safely through Product Service.
	// The context carries the correlation ID and request cancellation.
	CancelOrder(ctx context.Context, id uint) (*dto.OrderResponse, error)
}