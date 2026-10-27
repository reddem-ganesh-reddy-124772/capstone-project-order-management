package utilities

import (
	"strings"

	"order-management/services/product/model"
	"order-management/services/product/repository"
)

// ValidateProduct checks the product data before creating or updating it.
func ValidateProduct(product *model.Product) error {
	if product == nil {
		return ErrInvalidProduct
	}

	product.Name = strings.TrimSpace(product.Name)

	if product.Name == "" || len(product.Name) > 150 {
		return ErrInvalidProduct
	}

	if product.Price <= 0 || product.Stock < 0 {
		return ErrInvalidProduct
	}

	return nil
}

// ValidateProductID checks whether a product ID is valid.
func ValidateProductID(id uint) error {
	if id == 0 {
		return ErrInvalidProductID
	}

	return nil
}

// ValidatePage checks the requested page number.
func ValidatePage(page int) error {
	if page < 1 {
		return ErrInvalidPage
	}

	return nil
}

// ValidatePageSize checks the number of products requested per page.
func ValidatePageSize(limit int) error {
	if limit < 1 || limit > 100 {
		return ErrInvalidPageSize
	}

	return nil
}

// ValidatePriceRange checks the optional minimum and maximum prices.
func ValidatePriceRange(minPrice, maxPrice *float64) error {
	if minPrice != nil && *minPrice < 0 {
		return ErrInvalidPriceRange
	}

	if maxPrice != nil && *maxPrice < 0 {
		return ErrInvalidPriceRange
	}

	if minPrice != nil && maxPrice != nil && *minPrice > *maxPrice {
		return ErrInvalidPriceRange
	}

	return nil
}

// ValidateSorting checks the requested sort field and direction.
func ValidateSorting(query *repository.ProductQuery) error {
	if query == nil {
		return ErrInvalidProduct
	}

	query.SortBy = strings.ToLower(strings.TrimSpace(query.SortBy))
	query.SortOrder = strings.ToLower(strings.TrimSpace(query.SortOrder))

	// Use defaults when the client does not specify sorting.
	if query.SortBy == "" {
		query.SortBy = "id"
	}
	if query.SortOrder == "" {
		query.SortOrder = "asc"
	}

	// Only allow known database fields for sorting.
	allowedFields := map[string]bool{
		"id":         true,
		"name":       true,
		"price":      true,
		"stock":      true,
		"created_at": true,
		"updated_at": true,
	}

	if !allowedFields[query.SortBy] {
		return ErrInvalidSortField
	}

	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		return ErrInvalidSortOrder
	}

	return nil
}

// ValidateProductQuery checks pagination, price filters, and sorting.
func ValidateProductQuery(query *repository.ProductQuery) error {
	if query == nil {
		return ErrInvalidProduct
	}

	// Apply default pagination values when omitted.
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	if err := ValidatePage(query.Page); err != nil {
		return err
	}

	if err := ValidatePageSize(query.Limit); err != nil {
		return err
	}

	if err := ValidatePriceRange(query.MinPrice, query.MaxPrice); err != nil {
		return err
	}

	return ValidateSorting(query)
}
