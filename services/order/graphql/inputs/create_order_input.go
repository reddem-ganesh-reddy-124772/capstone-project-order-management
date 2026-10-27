package inputs

import "github.com/graphql-go/graphql"

// CreateOrderInputType defines the input required to create an order.
var CreateOrderInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "CreateOrderInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"customerId": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"items": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(
					graphql.NewList(
						graphql.NewNonNull(OrderItemInputType),
					),
				),
			},
		},
	},
)
