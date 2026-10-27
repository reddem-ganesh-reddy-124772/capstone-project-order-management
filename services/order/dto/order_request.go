
package dto

// OrderItemRequest contains the product and quantity
// requested for one item in an order.
type OrderItemRequest struct {
	// ProductID identifies the product in the Product Service.
	ProductID uint `json:"product_id" binding:"required"`

	// Quantity is the number of units requested.
	// It must be greater than zero.
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

// CreateOrderRequest contains the information needed
// to create an order with one or more products.
type CreateOrderRequest struct {
	// CustomerID identifies the customer placing the order.
	CustomerID uint `json:"customer_id" binding:"required"`

	// Items contains all products requested in this order.
	// At least one item must be provided.
	Items []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}
