
package service

import "context"

// ProductClient defines the Product Service operations
// required by the Order Service.
// It does not connect to the product database directly.
type ProductClient interface {
	// GetProduct retrieves the latest details of a product.
	GetProduct(ctx context.Context, productID uint) (*ProductInfo, error)

	// ReserveStock reserves the requested quantity using a unique
	// reservation ID so retrying the same request cannot reserve twice.
	ReserveStock(
		ctx context.Context,
		productID uint,
		reservationID string,
		quantity int,
	) (*ProductReservation, error)

	// ReleaseStock releases a reservation.
	// Repeating the same release must not restore stock twice.
	ReleaseStock(ctx context.Context, reservationID string) error
}

// ProductInfo contains the product details needed when creating an order.
// These details are received from Product Service through HTTP.
type ProductInfo struct {
	// ID identifies the product.
	ID uint `json:"id"`

	// Name is saved as a snapshot in the order item.
	Name string `json:"name"`

	// Price is the current product price.
	Price float64 `json:"price"`

	// Stock is the available stock reported by Product Service.
	Stock int `json:"stock"`
}

// ProductReservation describes the result of a stock reservation.
type ProductReservation struct {
	// ReservationID uniquely identifies this reservation.
	ReservationID string `json:"reservation_id"`

	// ProductID identifies the reserved product.
	ProductID uint `json:"product_id"`

	// Quantity is the number of units reserved.
	Quantity int `json:"quantity"`

	// Status indicates whether the reservation is reserved or released.
	Status string `json:"status"`
}
