package dto

// ProductQueryRequest contains optional filters, sorting, and pagination
// values that a client can provide when requesting a product list.
type ProductQueryRequest struct {
	// Page is the page number to return. The default is 1.
	Page int `form:"page" json:"page"`

	// Limit is the maximum number of products per page. The default is 10.
	Limit int `form:"limit" json:"limit"`

	// Name filters products whose names contain this text.
	Name string `form:"name" json:"name"`

	// MinPrice filters products at or above this price.
	MinPrice *float64 `form:"min_price" json:"min_price"`

	// MaxPrice filters products at or below this price.
	MaxPrice *float64 `form:"max_price" json:"max_price"`

	// InStock filters products by stock availability when provided.
	InStock *bool `form:"in_stock" json:"in_stock"`

	// SortBy identifies the field used to sort the results.
	SortBy string `form:"sort_by" json:"sort_by"`

	// SortOrder specifies ascending or descending order.
	SortOrder string `form:"sort_order" json:"sort_order"`
}
