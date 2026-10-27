package handler

import (
	"net/http"

	"order-management/services/order/graphql/schema"
	"order-management/services/order/service"

	graphqlhandler "github.com/graphql-go/handler"
)

// NewOrderGraphQLHandler creates the Order Service GraphQL handler
// and enables the built-in GraphiQL browser interface.
func NewOrderGraphQLHandler(
	orderService service.OrderService,
) (http.Handler, error) {
	graphqlSchema, err := schema.NewOrderSchema(orderService)
	if err != nil {
		return nil, err
	}

	handler := graphqlhandler.New(&graphqlhandler.Config{
		Schema:   &graphqlSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	return handler, nil
}
