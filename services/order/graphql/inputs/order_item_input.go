package inputs

import "github.com/graphql-go/graphql"

// OrderItemInputType defines the input fields for an order item.
var OrderItemInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "OrderItemInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"productId": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"quantity": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
	},
)
