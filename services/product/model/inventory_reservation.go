
package model

import "time"

// InventoryReservation records stock reserved for an order.
// Product Service owns this record because it manages product stock.
type InventoryReservation struct {
	// ID is the database-generated primary key.
	ID uint `gorm:"primaryKey" json:"id"`

	// ReservationID uniquely identifies this reservation.
	// Retried requests must reuse the same ID.
	ReservationID string `gorm:"size:100;not null;uniqueIndex" json:"reservation_id"`

	// ProductID identifies the product whose stock is reserved.
	ProductID uint `gorm:"not null;index" json:"product_id"`

	// Quantity is the number of units reserved.
	Quantity int `gorm:"not null" json:"quantity"`

	// Status records whether stock is currently reserved or released.
	// Allowed values are reserved and released.
	Status string `gorm:"size:20;not null;default:reserved" json:"status"`

	// CreatedAt records when this reservation was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt records when this reservation was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}
