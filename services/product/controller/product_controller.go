package controller

import (
	"errors"
	"net/http"
	"strconv"

	"order-management/services/product/dto"
	"order-management/services/product/service"
	"order-management/services/product/utilities"

	"github.com/gin-gonic/gin"
)

// ProductController handles REST requests for products.
type ProductController struct {
	// service contains the product business logic.
	service service.ProductService
}

// NewProductController creates a controller with the Product Service.
func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{
		service: productService,
	}
}

// Create creates a new product.
//
// @Summary Create a product
// @Description Creates a product after validating its details.
// @Tags Products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product details"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products [post]
func (c *ProductController) Create(ctx *gin.Context) {
	var request dto.CreateProductRequest

	// Read the JSON request body.
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Ask the service to validate and create the product.
	product, err := c.service.Create(ctx.Request.Context(), request)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

// GetByID returns one product by its ID.
//
// @Summary Get a product by ID
// @Description Returns the details of a product.
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products/{id} [get]
func (c *ProductController) GetByID(ctx *gin.Context) {
	// Convert the URL ID into an unsigned integer.
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := c.service.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// GetAll returns products with optional filtering, sorting, and pagination.
//
// @Summary List products
// @Description Lists products with optional filters, sorting, and pagination. Parameters can be combined.
// @Tags Products
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Records per page (default 10, maximum 100)"
// @Param name query string false "Filter by product name"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param in_stock query bool false "Filter by stock availability"
// @Param sort_by query string false "Sort field: id, name, price, stock, created_at, updated_at"
// @Param sort_order query string false "Sort direction: asc or desc"
// @Success 200 {object} dto.ProductListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products [get]
func (c *ProductController) GetAll(ctx *gin.Context) {
	var request dto.ProductQueryRequest

	// Read query parameters such as page, limit, name, and sort_by.
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	products, err := c.service.GetAll(ctx.Request.Context(), request)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// Update updates an existing product.
//
// @Summary Update a product
// @Description Updates the details of an existing product.
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body dto.UpdateProductRequest true "Updated product details"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products/{id} [put]
func (c *ProductController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var request dto.UpdateProductRequest

	// Read the updated product data from the JSON body.
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	product, err := c.service.Update(ctx.Request.Context(), uint(id), request)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// Delete deletes a product by its ID.
//
// @Summary Delete a product
// @Description Deletes an existing product.
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products/{id} [delete]
func (c *ProductController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

// handleError converts service errors into suitable HTTP responses.
func (c *ProductController) handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, utilities.ErrProductNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, utilities.ErrInvalidProduct),
		errors.Is(err, utilities.ErrInvalidProductID),
		errors.Is(err, utilities.ErrInvalidPage),
		errors.Is(err, utilities.ErrInvalidPageSize),
		errors.Is(err, utilities.ErrInvalidPriceRange),
		errors.Is(err, utilities.ErrInvalidSortField),
		errors.Is(err, utilities.ErrInvalidSortOrder):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
