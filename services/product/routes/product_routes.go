package routes

import (
	"order-management/services/product/controller"

	"github.com/gin-gonic/gin"
)

// RegisterProductRoutes connects product URLs to their controller methods.
func RegisterProductRoutes(router *gin.Engine, productController *controller.ProductController) {
	// Keep all product endpoints under the /api/products path.
	products := router.Group("/api/products")

	// Register the CRUD endpoints.
	products.POST("", productController.Create)
	products.GET("", productController.GetAll)
	products.GET("/:id", productController.GetByID)
	products.PUT("/:id", productController.Update)
	products.DELETE("/:id", productController.Delete)
}
