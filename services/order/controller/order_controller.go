
package controller

import (
	"net/http"
	"strconv"

	"order-management/services/order/dto"
	"order-management/services/order/service"

	"github.com/gin-gonic/gin"
)

// OrderController handles HTTP requests for orders.
type OrderController struct {
	orderService service.OrderService
}

// NewOrderController creates an Order Controller.
func NewOrderController(orderService service.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// CreateOrder creates an order.
//
// @Summary Create an order
// @Description Creates an order and reserves stock for its products.
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderRequest true "Order details"
// @Success 201 {object} dto.OrderResponse
// @Failure 400 {object} map[string]string
// @Router /api/orders [post]
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var request dto.CreateOrderRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order request",
		})
		return
	}

	order, err := c.orderService.CreateOrder(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, order)
}

// GetAllOrders retrieves all orders.
//
// @Summary Get all orders
// @Description Returns all orders with their items.
// @Tags Orders
// @Produce json
// @Success 200 {array} dto.OrderResponse
// @Failure 500 {object} map[string]string
// @Router /api/orders [get]
func (c *OrderController) GetAllOrders(ctx *gin.Context) {
	orders, err := c.orderService.GetAllOrders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve orders",
		})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

// GetOrderByID retrieves an order by ID.
//
// @Summary Get an order
// @Description Returns an order by its ID.
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/orders/{id} [get]
func (c *OrderController) GetOrderByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order ID",
		})
		return
	}

	order, err := c.orderService.GetOrderByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "order not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// ConfirmOrder confirms an order.
//
// @Summary Confirm an order
// @Description Confirms an order that is eligible for confirmation.
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]string
// @Router /api/orders/{id}/confirm [patch]
func (c *OrderController) ConfirmOrder(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order ID",
		})
		return
	}

	order, err := c.orderService.ConfirmOrder(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// CancelOrder cancels an order.
//
// @Summary Cancel an order
// @Description Cancels an order and requests release of reserved stock.
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]string
// @Router /api/orders/{id}/cancel [patch]
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order ID",
		})
		return
	}

	order, err := c.orderService.CancelOrder(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, order)
}
