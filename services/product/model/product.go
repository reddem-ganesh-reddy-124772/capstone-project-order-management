package model

import "time"

// Product represents a product stored in the Product Service database.
type Product struct {
	// ID uniquely identifies each product.
	ID uint `gorm:"primaryKey"`

	// Name stores the product name and cannot be empty.
	Name string `gorm:"size:150;not null"`

	// Description stores additional information about the product.
	Description string `gorm:"type:text"`

	// Price stores the product price with two decimal places.
	Price float64 `gorm:"type:numeric(12,2);not null"`

	// Stock stores the number of units available for sale.
	Stock int `gorm:"not null;default:0"`

	// CreatedAt stores when the product was created.
	CreatedAt time.Time

	// UpdatedAt stores when the product was last updated.
	UpdatedAt time.Time
}
