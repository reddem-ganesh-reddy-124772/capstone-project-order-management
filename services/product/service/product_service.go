package service

import (
	"context"

	"order-management/services/product/dto"
)

// ProductService defines the business operations for products.
// REST controllers and GraphQL resolvers use this interface.
type ProductService interface {
	// Create validates and saves a product, then returns its response.
	Create(ctx context.Context, request dto.CreateProductRequest) (*dto.ProductResponse, error)

	// GetByID returns a product response for the given ID.
	GetByID(ctx context.Context, id uint) (*dto.ProductResponse, error)

	// GetAll returns products and pagination details.
	GetAll(ctx context.Context, request dto.ProductQueryRequest) (*dto.ProductListResponse, error)

	// Update validates and saves product changes.
	Update(ctx context.Context, id uint, request dto.UpdateProductRequest) (*dto.ProductResponse, error)

	// Delete removes a product using its ID.
	Delete(ctx context.Context, id uint) error
}
