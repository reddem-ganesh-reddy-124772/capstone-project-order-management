
package controller

import (
	"context"

	"order-management/services/product/model"
)

// InventoryReservationService defines the operations
// required by the inventory reservation controller.
type InventoryReservationService interface {
	ReserveStock(
		ctx context.Context,
		reservation *model.InventoryReservation,
	) error

	ReleaseStock(
		ctx context.Context,
		reservationID string,
	) error

	GetReservation(
		ctx context.Context,
		reservationID string,
	) (*model.InventoryReservation, error)

	// GetAllReservations retrieves all inventory reservations.
	GetAllReservations(
		ctx context.Context,
	) ([]model.InventoryReservation, error)
}
