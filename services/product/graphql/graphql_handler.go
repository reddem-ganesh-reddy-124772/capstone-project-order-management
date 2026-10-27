package graphqlapi

import (
	"net/http"

	productschema "order-management/services/product/graphql/schema"
	"order-management/services/product/service"

	graphqlhandler "github.com/graphql-go/handler"
)

// NewGraphQLHandler creates an interactive GraphiQL UI and
// a handler for executing Product GraphQL queries and mutations.
func NewGraphQLHandler(productService service.ProductService) (http.Handler, error) {
	// Build the GraphQL schema using the Product Service.
	productSchema, err := productschema.NewProductSchema(productService)
	if err != nil {
		return nil, err
	}

	// Enable the browser-based GraphiQL interface.
	// GraphiQL supports writing and executing GraphQL operations.
	graphQLHandler := graphqlhandler.New(&graphqlhandler.Config{
		Schema:   &productSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	return graphQLHandler, nil
}
