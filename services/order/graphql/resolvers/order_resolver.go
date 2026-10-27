
package resolvers

import (
	"context"
	"fmt"
	"strconv"

	"order-management/services/order/dto"
	"order-management/services/order/service"
)

// OrderResolver connects GraphQL operations to the Order Service.
type OrderResolver struct {
	orderService service.OrderService
}

// NewOrderResolver creates an OrderResolver.
func NewOrderResolver(orderService service.OrderService) *OrderResolver {
	return &OrderResolver{
		orderService: orderService,
	}
}

// GetOrders retrieves all orders.
func (r *OrderResolver) GetOrders() (interface{}, error) {
	return r.orderService.GetAllOrders()
}

// GetOrder retrieves one order by ID.
func (r *OrderResolver) GetOrder(idValue interface{}) (interface{}, error) {
	id, err := parseID(idValue)
	if err != nil {
		return nil, err
	}

	return r.orderService.GetOrderByID(id)
}

// CreateOrder creates an order using the existing Order Service.
func (r *OrderResolver) CreateOrder(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	customerID, err := parseID(input["customerId"])
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %w", err)
	}

	rawItems, ok := input["items"].([]interface{})
	if !ok || len(rawItems) == 0 {
		return nil, fmt.Errorf("at least one order item is required")
	}

	request := dto.CreateOrderRequest{
		CustomerID: customerID,
		Items:      make([]dto.OrderItemRequest, 0, len(rawItems)),
	}

	for _, rawItem := range rawItems {
		item, ok := rawItem.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid order item input")
		}

		productID, err := parseID(item["productId"])
		if err != nil {
			return nil, fmt.Errorf("invalid product ID: %w", err)
		}

		quantity, err := parseID(item["quantity"])
		if err != nil {
			return nil, fmt.Errorf("invalid quantity: %w", err)
		}

		request.Items = append(request.Items, dto.OrderItemRequest{
			ProductID: productID,
			Quantity:  int(quantity),
		})
	}

	return r.orderService.CreateOrder(ctx, request)
}

// ConfirmOrder confirms an existing order.
func (r *OrderResolver) ConfirmOrder(idValue interface{}) (interface{}, error) {
	id, err := parseID(idValue)
	if err != nil {
		return nil, err
	}

	return r.orderService.ConfirmOrder(id)
}

// CancelOrder cancels an existing order.
func (r *OrderResolver) CancelOrder(ctx context.Context, idValue interface{}) (interface{}, error) {
	id, err := parseID(idValue)
	if err != nil {
		return nil, err
	}

	return r.orderService.CancelOrder(ctx, id)
}

// parseID converts a GraphQL ID or integer argument to uint.
func parseID(value interface{}) (uint, error) {
	var parsed uint64

	switch v := value.(type) {
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("ID must be a positive integer")
		}
		parsed = id

	case int:
		if v <= 0 {
			return 0, fmt.Errorf("ID must be a positive integer")
		}
		parsed = uint64(v)

	case int32:
		if v <= 0 {
			return 0, fmt.Errorf("ID must be a positive integer")
		}
		parsed = uint64(v)

	case int64:
		if v <= 0 {
			return 0, fmt.Errorf("ID must be a positive integer")
		}
		parsed = uint64(v)

	case float64:
		if v <= 0 || v != float64(uint64(v)) {
			return 0, fmt.Errorf("ID must be a positive integer")
		}
		parsed = uint64(v)

	default:
		return 0, fmt.Errorf("invalid ID value")
	}

	if parsed == 0 || parsed > uint64(^uint32(0)) {
		return 0, fmt.Errorf("ID is outside the supported range")
	}

	return uint(parsed), nil
}
