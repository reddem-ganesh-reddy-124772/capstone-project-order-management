
package mutations

import (
	"fmt"

	"order-management/services/order/graphql/inputs"
	"order-management/services/order/graphql/resolvers"
	"order-management/services/order/graphql/types"

	"github.com/graphql-go/graphql"
)

func NewOrderMutations(orderResolver *resolvers.OrderResolver) graphql.Fields {
	return graphql.Fields{
		"createOrder": &graphql.Field{
			Type: types.OrderType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(inputs.CreateOrderInputType),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				input, ok := p.Args["input"].(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("invalid create-order input")
				}

				return orderResolver.CreateOrder(p.Context, input)
			},
		},

		"confirmOrder": &graphql.Field{
			Type: types.OrderType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return orderResolver.ConfirmOrder(p.Args["id"])
			},
		},

		"cancelOrder": &graphql.Field{
			Type: types.OrderType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return orderResolver.CancelOrder(p.Context, p.Args["id"])
			},
		},
	}
}
