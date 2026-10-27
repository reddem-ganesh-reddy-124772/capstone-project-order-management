package queries

import (
	"order-management/services/order/graphql/resolvers"
	"order-management/services/order/graphql/types"

	"github.com/graphql-go/graphql"
)

// NewOrderQueries defines the read operations for Order Service.
func NewOrderQueries(orderResolver *resolvers.OrderResolver) graphql.Fields {
	return graphql.Fields{
		// Retrieve all orders.
		"orders": &graphql.Field{
			Type: graphql.NewNonNull(
				graphql.NewList(
					graphql.NewNonNull(types.OrderType),
				),
			),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return orderResolver.GetOrders()
			},
		},

		// Retrieve one order using its ID.
		"order": &graphql.Field{
			Type: types.OrderType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return orderResolver.GetOrder(p.Args["id"])
			},
		},
	}
}
