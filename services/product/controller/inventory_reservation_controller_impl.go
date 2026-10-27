
package controller

import (
	"net/http"

	"order-management/services/product/model"

	"github.com/gin-gonic/gin"
)

// inventoryReservationController handles inventory reservation HTTP requests.
type inventoryReservationController struct {
	service InventoryReservationService
}

// NewInventoryReservationController creates the controller.
func NewInventoryReservationController(
	service InventoryReservationService,
) *inventoryReservationController {
	return &inventoryReservationController{
		service: service,
	}
}

// ReserveStockRequest contains the information needed to reserve stock.
type ReserveStockRequest struct {
	ReservationID string `json:"reservation_id" binding:"required"`
	ProductID     uint   `json:"product_id" binding:"required"`
	Quantity      int    `json:"quantity" binding:"required,gt=0"`
}

// ReleaseStockRequest contains the reservation ID to release.
type ReleaseStockRequest struct {
	ReservationID string `json:"reservation_id" binding:"required"`
}

// ReserveStock godoc
// @Summary Reserve product stock
// @Description Reserves a specified quantity of an existing product.
// @Tags Inventory Reservations
// @Accept json
// @Produce json
// @Param request body ReserveStockRequest true "Stock reservation details"
// @Success 200 {object} map[string]interface{} "Stock reserved successfully"
// @Failure 400 {object} map[string]string "Invalid request or reservation error"
// @Router /api/inventory/reservations [post]
func (c *inventoryReservationController) ReserveStock(ctx *gin.Context) {
	var request ReserveStockRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid reservation request",
		})
		return
	}

	reservation := &model.InventoryReservation{
		ReservationID: request.ReservationID,
		ProductID:     request.ProductID,
		Quantity:      request.Quantity,
	}

	if err := c.service.ReserveStock(
		ctx.Request.Context(),
		reservation,
	); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":        "stock reserved successfully",
		"reservation_id": request.ReservationID,
		"product_id":     request.ProductID,
		"quantity":       request.Quantity,
	})
}

// ReleaseStock godoc
// @Summary Release reserved stock
// @Description Releases stock associated with an existing reservation.
// @Tags Inventory Reservations
// @Accept json
// @Produce json
// @Param request body ReleaseStockRequest true "Reservation to release"
// @Success 200 {object} map[string]string "Stock released successfully"
// @Failure 400 {object} map[string]string "Invalid request or release error"
// @Router /api/inventory/reservations/release [post]
func (c *inventoryReservationController) ReleaseStock(ctx *gin.Context) {
	var request ReleaseStockRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid release request",
		})
		return
	}

	if err := c.service.ReleaseStock(
		ctx.Request.Context(),
		request.ReservationID,
	); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":        "stock released successfully",
		"reservation_id": request.ReservationID,
	})
}

// GetReservation godoc
// @Summary Get a stock reservation
// @Description Retrieves a reservation using its reservation ID.
// @Tags Inventory Reservations
// @Produce json
// @Param reservation_id path string true "Reservation ID"
// @Success 200 {object} model.InventoryReservation
// @Failure 400 {object} map[string]string "Reservation ID is required"
// @Failure 404 {object} map[string]string "Reservation not found"
// @Router /api/inventory/reservations/{reservation_id} [get]
func (c *inventoryReservationController) GetReservation(ctx *gin.Context) {
	reservationID := ctx.Param("reservation_id")

	if reservationID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "reservation ID is required",
		})
		return
	}

	reservation, err := c.service.GetReservation(
		ctx.Request.Context(),
		reservationID,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

// GetAllReservations godoc
// @Summary Get all inventory reservations
// @Description Retrieves all inventory reservation records.
// @Tags Inventory Reservations
// @Produce json
// @Success 200 {array} model.InventoryReservation
// @Failure 500 {object} map[string]string "Failed to retrieve reservations"
// @Router /api/inventory/reservations [get]
func (c *inventoryReservationController) GetAllReservations(ctx *gin.Context) {
	reservations, err := c.service.GetAllReservations(
		ctx.Request.Context(),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve inventory reservations",
		})
		return
	}

	if reservations == nil {
		reservations = []model.InventoryReservation{}
	}

	ctx.JSON(http.StatusOK, reservations)
}
