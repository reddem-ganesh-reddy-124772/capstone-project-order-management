
package service

import (
	"context"

	"order-management/services/product/model"
)

// InventoryReservationService defines the business operations
// for managing inventory reservations.
type InventoryReservationService interface {
	// ReserveStock checks the request and reserves product stock.
	ReserveStock(
		ctx context.Context,
		reservation *model.InventoryReservation,
	) error

	// ReleaseStock restores stock for an existing reservation.
	// Repeating the request must not restore stock twice.
	ReleaseStock(
		ctx context.Context,
		reservationID string,
	) error

	// GetReservation retrieves one reservation by its ID.
	GetReservation(
		ctx context.Context,
		reservationID string,
	) (*model.InventoryReservation, error)

	// GetAllReservations retrieves all inventory reservations.
	GetAllReservations(
		ctx context.Context,
	) ([]model.InventoryReservation, error)
}
