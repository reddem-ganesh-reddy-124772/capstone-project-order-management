
package utilities

import (
	"errors"

	"order-management/services/order/model"
)

// ValidateStatusTransition checks whether an order can move
// from its current status to the requested status.
func ValidateStatusTransition(currentStatus string, newStatus string) error {
	// An order cannot be changed to an unknown status.
	if !isValidStatus(newStatus) {
		return errors.New("invalid order status")
	}

	// An order cannot be confirmed or canceled again after cancellation.
	if currentStatus == model.OrderStatusCanceled {
		return errors.New("a canceled order cannot be changed")
	}

	// A newly created order can be confirmed or canceled.
	if currentStatus == model.OrderStatusCreated {
		if newStatus == model.OrderStatusConfirmed ||
			newStatus == model.OrderStatusCanceled {
			return nil
		}
	}

	// A confirmed order can be canceled.
	if currentStatus == model.OrderStatusConfirmed {
		if newStatus == model.OrderStatusCanceled {
			return nil
		}
	}

	// Repeating the same status is not treated as a new transition.
	if currentStatus == newStatus {
		return errors.New("order already has this status")
	}

	// Any transition not explicitly allowed above is rejected.
	return errors.New("order status transition is not allowed")
}

// isValidStatus checks whether the supplied status is one of the
// three statuses defined for this project.
func isValidStatus(status string) bool {
	switch status {
	case model.OrderStatusCreated,
		model.OrderStatusConfirmed,
		model.OrderStatusCanceled:
		return true
	default:
		return false
	}
}
