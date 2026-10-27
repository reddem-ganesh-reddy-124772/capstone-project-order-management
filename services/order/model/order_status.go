
package model

// These constants define the three allowed order statuses.
const (
	// OrderStatusCreated means the order has been created.
	OrderStatusCreated = "created"

	// OrderStatusConfirmed means all required stock has been reserved.
	OrderStatusConfirmed = "confirmed"

	// OrderStatusCanceled means cancellation is complete
	// and all reserved stock has been released.
	OrderStatusCanceled = "canceled"
)
