
package model

import "time"

// Order represents one order placed by a customer.
// This model is stored in the Order Service database.
type Order struct {
	// ID is the unique ID generated for each order.
	ID uint `gorm:"primaryKey" json:"id"`

	// CustomerID identifies the customer who placed the order.
	// Customer information is managed separately.
	CustomerID uint `gorm:"not null;index" json:"customer_id"`

	// Status stores the current state of the order.
	// Examples: pending, confirmed, cancelled.
	Status string `gorm:"size:30;not null;default:pending;index" json:"status"`

	// TotalAmount stores the total price of all items in the order.
	// The Order Service will calculate this amount from product prices.
	TotalAmount float64 `gorm:"type:numeric(12,2);not null;default:0" json:"total_amount"`

	// Items contains the individual products included in this order.
	// GORM connects these records using Order.ID and OrderItem.OrderID.
	// An order can contain multiple items.
	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items"`

	// CreatedAt records when the order was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt records when the order was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}
