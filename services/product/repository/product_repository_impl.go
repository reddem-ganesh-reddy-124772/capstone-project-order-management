package repository

import (
	"context"
	"errors"
	"fmt"

	"order-management/services/product/model"

	"gorm.io/gorm"
)

// ProductRepositoryImpl implements product database operations.
type ProductRepositoryImpl struct {
	db *gorm.DB
}

// NewProductRepository creates a product repository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

// Create saves a new product.
func (r *ProductRepositoryImpl) Create(ctx context.Context, product *model.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("create product: %w", err)
	}

	return nil
}

// GetByID finds a product by its ID.
func (r *ProductRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product

	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, fmt.Errorf("get product: %w", err)
	}

	return &product, nil
}

// GetAll returns products with optional filters, sorting, and pagination.
func (r *ProductRepositoryImpl) GetAll(
	ctx context.Context,
	query ProductQuery,
) ([]model.Product, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Product{})

	// Filter products by name when a name is provided.
	if query.Name != "" {
		db = db.Where("name ILIKE ?", "%"+query.Name+"%")
	}

	// Apply the minimum and maximum price filters when provided.
	if query.MinPrice != nil {
		db = db.Where("price >= ?", *query.MinPrice)
	}

	if query.MaxPrice != nil {
		db = db.Where("price <= ?", *query.MaxPrice)
	}

	// Filter by stock availability when requested.
	if query.InStock != nil {
		if *query.InStock {
			db = db.Where("stock > 0")
		} else {
			db = db.Where("stock = 0")
		}
	}

	// Count matching products before applying pagination.
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	// Use default pagination values if none were supplied.
	page := query.Page
	if page < 1 {
		page = 1
	}

	limit := query.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Allow only known columns to prevent unsafe sort expressions.
	sortColumns := map[string]string{
		"id":         "id",
		"name":       "name",
		"price":      "price",
		"stock":      "stock",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	sortColumn, ok := sortColumns[query.SortBy]
	if !ok {
		sortColumn = "id"
	}

	// Use ascending order unless descending order was requested.
	sortOrder := "ASC"
	if query.SortOrder == "desc" {
		sortOrder = "DESC"
	}

	// Add ID as a secondary sort to keep results stable.
	orderBy := sortColumn + " " + sortOrder
	if sortColumn != "id" {
		orderBy += ", id ASC"
	}

	// Fetch only the requested page from PostgreSQL.
	var products []model.Product
	offset := (page - 1) * limit

	err := db.
		Order(orderBy).
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, 0, fmt.Errorf("get products: %w", err)
	}

	return products, total, nil
}

// Update saves changes to an existing product.
func (r *ProductRepositoryImpl) Update(ctx context.Context, product *model.Product) error {
	result := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", product.ID).
		Select("Name", "Description", "Price", "Stock", "UpdatedAt").
		Updates(product)

	if result.Error != nil {
		return fmt.Errorf("update product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete removes a product by its ID.
func (r *ProductRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Product{}, id)

	if result.Error != nil {
		return fmt.Errorf("delete product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
