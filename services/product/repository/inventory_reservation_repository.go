
package repository

import "order-management/services/product/model"

// InventoryReservationRepository defines the database operations
// required to manage inventory reservations.
type InventoryReservationRepository interface {
	// ReserveStock deducts stock and saves its reservation safely.
	ReserveStock(reservation *model.InventoryReservation) error

	// ReleaseStock restores stock only once for a reservation.
	ReleaseStock(reservationID string) error

	// GetReservation retrieves one reservation using its unique ID.
	GetReservation(reservationID string) (*model.InventoryReservation, error)

	// GetAllReservations retrieves all inventory reservations.
	GetAllReservations() ([]model.InventoryReservation, error)
}
