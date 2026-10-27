
package repository

import "order-management/services/order/model"

// OrderRepository defines the database operations
// required by the Order Service.
type OrderRepository interface {
	// Create saves a new order and its items.
	Create(order *model.Order) error

	// GetByID retrieves one order, including its items.
	GetByID(id uint) (*model.Order, error)

	// GetAll retrieves orders with their associated items.
	GetAll() ([]model.Order, error)

	// Update saves changes to an existing order.
	Update(order *model.Order) error

	// UpdateStatus changes the status of an existing order.
	UpdateStatus(id uint, status string) error
}
