package dto

// CreateProductRequest contains the fields required to create a product.
type CreateProductRequest struct {
	// Name is the name of the product.
	Name string `json:"name" binding:"required,max=150"`

	// Description provides additional information about the product.
	Description string `json:"description"`

	// Price is the product's selling price and must be greater than zero.
	Price float64 `json:"price" binding:"required,gt=0"`

	// Stock is the available quantity and cannot be negative.
	Stock int `json:"stock" binding:"gte=0"`
}

// UpdateProductRequest contains the fields required to update a product.
type UpdateProductRequest struct {
	// Name is the updated product name.
	Name string `json:"name" binding:"required,max=150"`

	// Description is the updated product description.
	Description string `json:"description"`

	// Price is the updated selling price and must be greater than zero.
	Price float64 `json:"price" binding:"required,gt=0"`

	// Stock is the updated available quantity and cannot be negative.
	Stock int `json:"stock" binding:"gte=0"`
}
