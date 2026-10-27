
package routes

import "github.com/gin-gonic/gin"

// InventoryReservationHandler defines the controller methods
// required by the inventory reservation routes.
type InventoryReservationHandler interface {
	ReserveStock(ctx *gin.Context)
	ReleaseStock(ctx *gin.Context)
	GetReservation(ctx *gin.Context)
	GetAllReservations(ctx *gin.Context)
}

// RegisterInventoryReservationRoutes registers the inventory endpoints.
func RegisterInventoryReservationRoutes(
	router *gin.RouterGroup,
	inventoryController InventoryReservationHandler,
) {
	// Create a route group for inventory reservation operations.
	inventoryRoutes := router.Group("/inventory/reservations")

	// Reserve stock for an order.
	inventoryRoutes.POST("", inventoryController.ReserveStock)

	// Release stock associated with a reservation.
	inventoryRoutes.POST("/release", inventoryController.ReleaseStock)

	// Retrieve all inventory reservations.
	inventoryRoutes.GET("", inventoryController.GetAllReservations)

	// Retrieve one reservation using its ID.
	inventoryRoutes.GET("/:reservation_id", inventoryController.GetReservation)
}
