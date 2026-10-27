package schema

import (
	"order-management/services/product/graphql/mutations"
	"order-management/services/product/graphql/queries"
	"order-management/services/product/service"

	"github.com/graphql-go/graphql"
)

// NewProductSchema combines product queries and mutations into one GraphQL schema.
func NewProductSchema(productService service.ProductService) (graphql.Schema, error) {
	// Create the root query object for product read operations.
	query := queries.NewProductQueries(productService)

	// Create the root mutation object for product write operations.
	mutation := mutations.NewProductMutations(productService)

	// Build the schema used to execute Product GraphQL requests.
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
}
