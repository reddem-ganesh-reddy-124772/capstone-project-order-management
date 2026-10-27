
package service

import (
	"context"
	"errors"

	"order-management/services/product/model"
	"order-management/services/product/repository"
)

// inventoryReservationService implements the inventory reservation rules.
type inventoryReservationService struct {
	repository repository.InventoryReservationRepository
}

// NewInventoryReservationService creates the inventory reservation service.
func NewInventoryReservationService(
	repo repository.InventoryReservationRepository,
) InventoryReservationService {
	return &inventoryReservationService{
		repository: repo,
	}
}

// ReserveStock validates the request before reserving stock.
func (s *inventoryReservationService) ReserveStock(
	ctx context.Context,
	reservation *model.InventoryReservation,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if reservation == nil {
		return errors.New("reservation is required")
	}

	if reservation.ReservationID == "" {
		return errors.New("reservation ID is required")
	}

	if reservation.ProductID == 0 {
		return errors.New("product ID is required")
	}

	if reservation.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	return s.repository.ReserveStock(reservation)
}

// ReleaseStock releases a reservation safely.
func (s *inventoryReservationService) ReleaseStock(
	ctx context.Context,
	reservationID string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if reservationID == "" {
		return errors.New("reservation ID is required")
	}

	return s.repository.ReleaseStock(reservationID)
}

// GetReservation retrieves a reservation by its ID.
func (s *inventoryReservationService) GetReservation(
	ctx context.Context,
	reservationID string,
) (*model.InventoryReservation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if reservationID == "" {
		return nil, errors.New("reservation ID is required")
	}

	return s.repository.GetReservation(reservationID)
}

// GetAllReservations retrieves all inventory reservations.
func (s *inventoryReservationService) GetAllReservations(
	ctx context.Context,
) ([]model.InventoryReservation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return s.repository.GetAllReservations()
}
