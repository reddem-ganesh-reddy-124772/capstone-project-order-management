package routes

import (
	"order-management/services/order/controller"

	"github.com/gin-gonic/gin"
)

// RegisterOrderRoutes registers all Order Service endpoints.
func RegisterOrderRoutes(
	router *gin.Engine,
	orderController *controller.OrderController,
) {
	orderRoutes := router.Group("/api/orders")

	orderRoutes.POST("", orderController.CreateOrder)
	orderRoutes.GET("", orderController.GetAllOrders)
	orderRoutes.GET("/:id", orderController.GetOrderByID)
	orderRoutes.PATCH("/:id/confirm", orderController.ConfirmOrder)
	orderRoutes.PATCH("/:id/cancel", orderController.CancelOrder)
}
