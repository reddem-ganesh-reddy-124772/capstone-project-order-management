package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"order-management/services/order/dto"
	"order-management/services/order/model"
	"order-management/services/order/repository"
	"order-management/services/order/utilities"
)

type orderService struct {
	orderRepository repository.OrderRepository
	productClient   ProductClient
}

func NewOrderService(
	orderRepository repository.OrderRepository,
	productClient ProductClient,
) OrderService {
	return &orderService{
		orderRepository: orderRepository,
		productClient:   productClient,
	}
}

// CreateOrder retrieves product details, reserves stock, calculates
// prices on the server, and saves the order.
func (s *orderService) CreateOrder(
	parentCtx context.Context,
	request dto.CreateOrderRequest,
) (*dto.OrderResponse, error) {
	if request.CustomerID == 0 {
		return nil, errors.New("customer ID is required")
	}

	if len(request.Items) == 0 {
		return nil, errors.New("at least one order item is required")
	}

	// Preserve the incoming correlation ID while applying a timeout.
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	order := &model.Order{
		CustomerID: request.CustomerID,
		Status:     model.OrderStatusCreated,
	}

	// Track reservations so they can be released if a later step fails.
	reservedIDs := make([]string, 0, len(request.Items))

	releaseReservations := func() error {
		var releaseErrors []error

		for i := len(reservedIDs) - 1; i >= 0; i-- {
			err := s.productClient.ReleaseStock(ctx, reservedIDs[i])
			if err != nil {
				releaseErrors = append(
					releaseErrors,
					fmt.Errorf(
						"release reservation %s: %w",
						reservedIDs[i],
						err,
					),
				)
			}
		}

		return errors.Join(releaseErrors...)
	}

	for _, itemRequest := range request.Items {
		if itemRequest.ProductID == 0 || itemRequest.Quantity <= 0 {
			return nil, errors.New("product ID and positive quantity are required")
		}

		product, err := s.productClient.GetProduct(ctx, itemRequest.ProductID)
		if err != nil {
			releaseErr := releaseReservations()
			return nil, errors.Join(
				fmt.Errorf("get product %d: %w", itemRequest.ProductID, err),
				releaseErr,
			)
		}

		if product == nil || product.ID == 0 {
			releaseErr := releaseReservations()
			return nil, errors.Join(
				fmt.Errorf("product %d was not found", itemRequest.ProductID),
				releaseErr,
			)
		}

		// A reservation ID is generated for each order item.
		reservationID := fmt.Sprintf(
			"ORD-%d-%d-%d",
			time.Now().UnixNano(),
			product.ID,
			len(order.Items)+1,
		)

		_, err = s.productClient.ReserveStock(
			ctx,
			product.ID,
			reservationID,
			itemRequest.Quantity,
		)
		if err != nil {
			releaseErr := releaseReservations()
			return nil, errors.Join(
				fmt.Errorf("reserve stock for product %d: %w", product.ID, err),
				releaseErr,
			)
		}

		reservedIDs = append(reservedIDs, reservationID)

		item := model.OrderItem{
			ProductID:     product.ID,
			ProductName:   product.Name,
			Quantity:      itemRequest.Quantity,
			UnitPrice:     product.Price,
			ReservationID: reservationID,
			StockReserved: true,
			StockReleased: false,
		}

		order.Items = append(order.Items, item)
		order.TotalAmount += product.Price * float64(itemRequest.Quantity)
	}

	// Save the order and its items after stock has been reserved.
	if err := s.orderRepository.Create(order); err != nil {
		releaseErr := releaseReservations()
		return nil, errors.Join(
			fmt.Errorf("save order: %w", err),
			releaseErr,
		)
	}

	response := convertOrderToResponse(order)
	return &response, nil
}

// GetOrderByID retrieves one order.
func (s *orderService) GetOrderByID(
	id uint,
) (*dto.OrderResponse, error) {
	if id == 0 {
		return nil, errors.New("order ID must be greater than zero")
	}

	order, err := s.orderRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := convertOrderToResponse(order)
	return &response, nil
}

// GetAllOrders retrieves all orders.
func (s *orderService) GetAllOrders() ([]dto.OrderResponse, error) {
	orders, err := s.orderRepository.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.OrderResponse, 0, len(orders))

	for i := range orders {
		response := convertOrderToResponse(&orders[i])
		responses = append(responses, response)
	}

	return responses, nil
}

// ConfirmOrder changes a created order to confirmed.
func (s *orderService) ConfirmOrder(
	id uint,
) (*dto.OrderResponse, error) {
	order, err := s.orderRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	err = utilities.ValidateStatusTransition(
		order.Status,
		model.OrderStatusConfirmed,
	)
	if err != nil {
		return nil, err
	}

	err = s.orderRepository.UpdateStatus(id, model.OrderStatusConfirmed)
	if err != nil {
		return nil, err
	}

	order.Status = model.OrderStatusConfirmed

	response := convertOrderToResponse(order)
	return &response, nil
}

// CancelOrder releases all reserved stock before changing the order
// status to canceled.
func (s *orderService) CancelOrder(
	parentCtx context.Context,
	id uint,
) (*dto.OrderResponse, error) {
	order, err := s.orderRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	err = utilities.ValidateStatusTransition(
		order.Status,
		model.OrderStatusCanceled,
	)
	if err != nil {
		return nil, err
	}

	// Preserve the incoming correlation ID while applying a timeout.
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	// Release each reservation before marking the order as canceled.
	for i := range order.Items {
		item := &order.Items[i]

		if !item.StockReserved || item.StockReleased {
			continue
		}

		if item.ReservationID == "" {
			return nil, errors.New("order item has no reservation ID")
		}

		err := s.productClient.ReleaseStock(ctx, item.ReservationID)
		if err != nil {
			return nil, fmt.Errorf(
				"release stock for reservation %s: %w",
				item.ReservationID,
				err,
			)
		}

		item.StockReleased = true
	}

	// Do not change the order status if releasing stock fails.
	err = s.orderRepository.UpdateStatus(id, model.OrderStatusCanceled)
	if err != nil {
		return nil, err
	}

	order.Status = model.OrderStatusCanceled

	response := convertOrderToResponse(order)
	return &response, nil
}

// convertOrderToResponse converts database models into response DTOs.
func convertOrderToResponse(order *model.Order) dto.OrderResponse {
	response := dto.OrderResponse{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		Items:       make([]dto.OrderItemResponse, 0, len(order.Items)),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}

	for _, item := range order.Items {
		response.Items = append(response.Items, dto.OrderItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    item.UnitPrice * float64(item.Quantity),
		})
	}

	return response
}