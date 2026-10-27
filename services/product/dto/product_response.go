package dto

import "time"

// ProductResponse contains the product information returned to clients.
type ProductResponse struct {
	// ID is the unique product identifier.
	ID uint `json:"id"`

	// Name is the product name.
	Name string `json:"name"`

	// Description explains the product.
	Description string `json:"description"`

	// Price is the product's selling price.
	Price float64 `json:"price"`

	// Stock is the available quantity.
	Stock int `json:"stock"`

	// CreatedAt records when the product was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt records when the product was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductListResponse contains a page of products and pagination details.
type ProductListResponse struct {
	// Items contains the products on the current page.
	Items []ProductResponse `json:"items"`

	// Page is the current page number.
	Page int `json:"page"`

	// Limit is the maximum number of products per page.
	Limit int `json:"limit"`

	// TotalRecords is the number of products matching the filters.
	TotalRecords int64 `json:"total_records"`

	// TotalPages is the number of pages available for the matching products.
	TotalPages int64 `json:"total_pages"`
}
