package service

import (
	"context"
	"errors"
	"fmt"

	"order-management/services/product/dto"
	"order-management/services/product/model"
	"order-management/services/product/repository"
	"order-management/services/product/utilities"

	"gorm.io/gorm"
)

// ProductServiceImpl contains the business logic for products.
// It uses the repository to access product data.
type ProductServiceImpl struct {
	repo repository.ProductRepository
}

// NewProductService creates the Product Service with a repository.
func NewProductService(repo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{repo: repo}
}

// Create validates the request, saves the product, and returns its response.
func (s *ProductServiceImpl) Create(
	ctx context.Context,
	request dto.CreateProductRequest,
) (*dto.ProductResponse, error) {
	// Convert the request DTO into a database model.
	product := &model.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	// Validate product fields using the shared utility.
	if err := utilities.ValidateProduct(product); err != nil {
		return nil, err
	}

	// Save the product through the repository.
	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	return toProductResponse(product), nil
}

// GetByID finds a product and returns its response DTO.
func (s *ProductServiceImpl) GetByID(
	ctx context.Context,
	id uint,
) (*dto.ProductResponse, error) {
	// Check the ID before accessing the database.
	if err := utilities.ValidateProductID(id); err != nil {
		return nil, err
	}

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utilities.ErrProductNotFound
		}
		return nil, fmt.Errorf("get product: %w", err)
	}

	return toProductResponse(product), nil
}

// GetAll validates the query and returns a page of product responses.
func (s *ProductServiceImpl) GetAll(
	ctx context.Context,
	request dto.ProductQueryRequest,
) (*dto.ProductListResponse, error) {
	// Convert the request DTO into repository query options.
	query := repository.ProductQuery{
		Page:      request.Page,
		Limit:     request.Limit,
		Name:      request.Name,
		MinPrice:  request.MinPrice,
		MaxPrice:  request.MaxPrice,
		InStock:   request.InStock,
		SortBy:    request.SortBy,
		SortOrder: request.SortOrder,
	}

	// Validate filters, sorting, and pagination in the utilities package.
	if err := utilities.ValidateProductQuery(&query); err != nil {
		return nil, err
	}

	// Retrieve the requested page and total matching record count.
	products, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	// Convert database models into response DTOs.
	items := make([]dto.ProductResponse, 0, len(products))
	for i := range products {
		items = append(items, *toProductResponse(&products[i]))
	}

	// Calculate the number of pages needed for all matching products.
	totalPages := total / int64(query.Limit)
	if total%int64(query.Limit) != 0 {
		totalPages++
	}

	// Return the products together with their pagination details.
	return &dto.ProductListResponse{
		Items:        items,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalRecords: total,
		TotalPages:   totalPages,
	}, nil
}

// Update validates the request and saves changes to an existing product.
func (s *ProductServiceImpl) Update(
	ctx context.Context,
	id uint,
	request dto.UpdateProductRequest,
) (*dto.ProductResponse, error) {
	// Check the ID before querying the database.
	if err := utilities.ValidateProductID(id); err != nil {
		return nil, err
	}

	// Find the existing product so its creation details are preserved.
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utilities.ErrProductNotFound
		}
		return nil, fmt.Errorf("find product before update: %w", err)
	}

	// Replace the editable fields with the supplied request values.
	product.Name = request.Name
	product.Description = request.Description
	product.Price = request.Price
	product.Stock = request.Stock

	// Validate the updated product using the shared utility.
	if err := utilities.ValidateProduct(product); err != nil {
		return nil, err
	}

	// Save the changes through the repository.
	if err := s.repo.Update(ctx, product); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utilities.ErrProductNotFound
		}
		return nil, fmt.Errorf("update product: %w", err)
	}

	return toProductResponse(product), nil
}

// Delete removes an existing product.
func (s *ProductServiceImpl) Delete(ctx context.Context, id uint) error {
	// Validate the ID before accessing the database.
	if err := utilities.ValidateProductID(id); err != nil {
		return err
	}

	// Check that the product exists before deleting it.
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utilities.ErrProductNotFound
		}
		return fmt.Errorf("find product before delete: %w", err)
	}

	// Delete the product through the repository.
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utilities.ErrProductNotFound
		}
		return fmt.Errorf("delete product: %w", err)
	}

	return nil
}

// toProductResponse converts a database model into a response DTO.
// This keeps database model details separate from API responses.
func toProductResponse(product *model.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// This check ensures the implementation matches the ProductService interface.
var _ ProductService = (*ProductServiceImpl)(nil)
