
package repository

import (
	"errors"

	"order-management/services/order/model"

	"gorm.io/gorm"
)

// orderRepository is the concrete implementation of OrderRepository.
// It uses GORM to communicate with the Order Service database.
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a repository using the supplied database connection.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// Create saves an order and its associated items in one database operation.
// GORM creates the associated items when saving the order.
func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(order).Error
	})
}

// GetByID retrieves an order and loads all its associated items.
func (r *orderRepository) GetByID(id uint) (*model.Order, error) {
	var order model.Order

	// Preload loads the order items along with the order.
	err := r.db.Preload("Items").First(&order, id).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetAll retrieves all orders and their items.
// Orders are returned newest first.
func (r *orderRepository) GetAll() ([]model.Order, error) {
	var orders []model.Order

	err := r.db.
		Preload("Items").
		Order("created_at DESC").
		Order("id DESC").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// Update saves changes to an existing order's main fields.
// Order items are managed separately.
func (r *orderRepository) Update(order *model.Order) error {
	result := r.db.Model(&model.Order{}).
		Where("id = ?", order.ID).
		Select("CustomerID", "Status", "TotalAmount").
		Updates(order)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// UpdateStatus updates only the status of an existing order.
// The service layer must validate status transitions before calling this method.
func (r *orderRepository) UpdateStatus(id uint, status string) error {
	result := r.db.Model(&model.Order{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		// Confirm that the order does not exist.
		var order model.Order
		err := r.db.Select("id").First(&order, id).Error
		if err != nil {
			return err
		}

		// No row changed, so report that the requested update did not occur.
		return errors.New("order status was not updated")
	}

	return nil
}
