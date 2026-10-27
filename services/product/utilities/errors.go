package utilities

import "errors"

// Product-related errors shared by the Product Service.
var (
	// ErrProductNotFound means the requested product does not exist.
	ErrProductNotFound = errors.New("product not found")

	// ErrInvalidProduct means the supplied product data is invalid.
	ErrInvalidProduct = errors.New("invalid product data")

	// ErrInvalidProductID means the product ID is invalid.
	ErrInvalidProductID = errors.New("product ID must be greater than zero")

	// ErrInvalidPage means the requested page number is invalid.
	ErrInvalidPage = errors.New("page must be greater than zero")

	// ErrInvalidPageSize means the requested page size is outside the allowed range.
	ErrInvalidPageSize = errors.New("page size must be between 1 and 100")

	// ErrInvalidPriceRange means the supplied price range is invalid.
	ErrInvalidPriceRange = errors.New("minimum price cannot exceed maximum price")

	// ErrInvalidSortField means the requested sort field is not supported.
	ErrInvalidSortField = errors.New("invalid sort field")

	// ErrInvalidSortOrder means the requested sort direction is not supported.
	ErrInvalidSortOrder = errors.New("sort order must be asc or desc")
)
