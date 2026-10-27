package repository

import (
	"context"

	"order-management/services/product/model"
)

// ProductRepository defines the database operations for products.
type ProductRepository interface {
	// Create saves a new product.
	Create(ctx context.Context, product *model.Product) error

	// GetByID finds a product by its ID.
	GetByID(ctx context.Context, id uint) (*model.Product, error)

	// GetAll returns products using optional filters, sorting, and pagination.
	GetAll(ctx context.Context, query ProductQuery) ([]model.Product, int64, error)

	// Update saves changes to an existing product.
	Update(ctx context.Context, product *model.Product) error

	// Delete removes a product by its ID.
	Delete(ctx context.Context, id uint) error
}

// ProductQuery contains optional product search and listing options.
type ProductQuery struct {
	Page      int
	Limit     int
	Name      string
	MinPrice  *float64
	MaxPrice  *float64
	InStock   *bool
	SortBy    string
	SortOrder string
}
