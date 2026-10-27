package schema

import (
	"order-management/services/order/graphql/mutations"
	"order-management/services/order/graphql/queries"
	"order-management/services/order/graphql/resolvers"
	"order-management/services/order/service"

	"github.com/graphql-go/graphql"
)

// NewOrderSchema assembles the Order Service GraphQL schema.
func NewOrderSchema(orderService service.OrderService) (graphql.Schema, error) {
	// Initialize the resolver with the existing Order Service.
	orderResolver := resolvers.NewOrderResolver(orderService)

	// Define the root query type.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: queries.NewOrderQueries(orderResolver),
	})

	// Define the root mutation type.
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: mutations.NewOrderMutations(orderResolver),
	})

	// Assemble the schema.
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}
