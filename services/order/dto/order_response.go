
package dto

import "time"

// OrderItemResponse contains the details of one product in an order.
type OrderItemResponse struct {
	// ID uniquely identifies this order item.
	ID uint `json:"id"`

	// ProductID identifies the product in Product Service.
	ProductID uint `json:"product_id"`

	// ProductName stores the product name captured for this order.
	ProductName string `json:"product_name"`

	// Quantity is the number of units ordered.
	Quantity int `json:"quantity"`

	// UnitPrice is the product price captured when the order was placed.
	UnitPrice float64 `json:"unit_price"`

	// Subtotal is the unit price multiplied by the quantity.
	Subtotal float64 `json:"subtotal"`
}

// OrderResponse contains the order details returned to the client.
type OrderResponse struct {
	// ID uniquely identifies the order.
	ID uint `json:"id"`

	// CustomerID identifies the customer who placed the order.
	CustomerID uint `json:"customer_id"`

	// Status is created, confirmed, or canceled.
	Status string `json:"status"`

	// TotalAmount is the total cost of all order items.
	TotalAmount float64 `json:"total_amount"`

	// Items contains the products included in the order.
	Items []OrderItemResponse `json:"items"`

	// CreatedAt records when the order was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt records when the order was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}
