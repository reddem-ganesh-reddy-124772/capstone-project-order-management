
package repository

import (
	"errors"

	"order-management/services/product/model"

	"gorm.io/gorm"
)

// inventoryReservationRepository performs reservation operations
// using the Product Service database.
type inventoryReservationRepository struct {
	db *gorm.DB
}

// NewInventoryReservationRepository creates the repository.
func NewInventoryReservationRepository(
	db *gorm.DB,
) InventoryReservationRepository {
	return &inventoryReservationRepository{db: db}
}

// ReserveStock deducts available stock and saves the reservation
// in one transaction.
func (r *inventoryReservationRepository) ReserveStock(
	reservation *model.InventoryReservation,
) error {
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

	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing model.InventoryReservation

		err := tx.Where(
			"reservation_id = ?",
			reservation.ReservationID,
		).First(&existing).Error

		if err == nil {
			if existing.ProductID != reservation.ProductID ||
				existing.Quantity != reservation.Quantity {
				return errors.New("reservation ID belongs to a different request")
			}

			if existing.Status == model.InventoryReservationStatusReserved {
				return nil
			}

			return errors.New("reservation was already released")
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		result := tx.Model(&model.Product{}).
			Where(
				"id = ? AND stock >= ?",
				reservation.ProductID,
				reservation.Quantity,
			).
			UpdateColumn(
				"stock",
				gorm.Expr("stock - ?", reservation.Quantity),
			)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("product not found or insufficient stock")
		}

		reservation.Status = model.InventoryReservationStatusReserved

		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		return nil
	})
}

// ReleaseStock restores stock only when the reservation is still active.
func (r *inventoryReservationRepository) ReleaseStock(
	reservationID string,
) error {
	if reservationID == "" {
		return errors.New("reservation ID is required")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		var reservation model.InventoryReservation

		err := tx.Where(
			"reservation_id = ?",
			reservationID,
		).First(&reservation).Error

		if err != nil {
			return err
		}

		if reservation.Status == model.InventoryReservationStatusReleased {
			return nil
		}

		if reservation.Status != model.InventoryReservationStatusReserved {
			return errors.New("reservation has an invalid status")
		}

		result := tx.Model(&model.InventoryReservation{}).
			Where(
				"id = ? AND status = ?",
				reservation.ID,
				model.InventoryReservationStatusReserved,
			).
			Update("status", model.InventoryReservationStatusReleased)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			var current model.InventoryReservation

			err := tx.Where(
				"reservation_id = ?",
				reservationID,
			).First(&current).Error

			if err != nil {
				return err
			}

			if current.Status == model.InventoryReservationStatusReleased {
				return nil
			}

			return errors.New("reservation could not be released")
		}

		result = tx.Model(&model.Product{}).
			Where("id = ?", reservation.ProductID).
			UpdateColumn(
				"stock",
				gorm.Expr("stock + ?", reservation.Quantity),
			)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("reserved product was not found")
		}

		return nil
	})
}

// GetReservation retrieves one reservation using its unique ID.
func (r *inventoryReservationRepository) GetReservation(
	reservationID string,
) (*model.InventoryReservation, error) {
	var reservation model.InventoryReservation

	err := r.db.Where(
		"reservation_id = ?",
		reservationID,
	).First(&reservation).Error

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

// GetAllReservations retrieves all inventory reservations.
func (r *inventoryReservationRepository) GetAllReservations() ([]model.InventoryReservation, error) {
	var reservations []model.InventoryReservation

	err := r.db.Order("id DESC").Find(&reservations).Error
	if err != nil {
		return nil, err
	}

	return reservations, nil
}
