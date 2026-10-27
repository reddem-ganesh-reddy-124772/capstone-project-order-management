
package model

import "time"

// OrderItem represents one product included in an order.
// Each item has its own stock reservation ID.
type OrderItem struct {
	// ID uniquely identifies this order item.
	ID uint `gorm:"primaryKey" json:"id"`

	// OrderID identifies the order to which this item belongs.
	OrderID uint `gorm:"not null;index" json:"order_id"`

	// ProductID identifies the product in Product Service.
	ProductID uint `gorm:"not null;index" json:"product_id"`

	// ProductName stores the product name captured when ordering.
	ProductName string `gorm:"size:150;not null" json:"product_name"`

	// Quantity is the number of units requested.
	Quantity int `gorm:"not null" json:"quantity"`

	// UnitPrice stores the product price captured when ordering.
	UnitPrice float64 `gorm:"type:numeric(12,2);not null" json:"unit_price"`

	// ReservationID uniquely identifies this item's stock reservation
	// in Product Service. Retries must reuse the same ID.
	ReservationID string `gorm:"size:100;uniqueIndex" json:"reservation_id"`

	// StockReserved indicates that the reservation was confirmed locally.
	StockReserved bool `gorm:"not null;default:false" json:"stock_reserved"`

	// StockReleased indicates that the reservation was successfully released.
	// Product Service must also make release operations idempotent.
	StockReleased bool `gorm:"not null;default:false" json:"stock_released"`

	// CreatedAt records when this order item was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt records when this order item was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}
