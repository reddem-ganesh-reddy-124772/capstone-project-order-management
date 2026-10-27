package types

import "github.com/graphql-go/graphql"

// OrderItemType defines the GraphQL output fields for an order item.
var OrderItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "OrderItem",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"productId": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"productName": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"quantity": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"unitPrice": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"subtotal": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Float),
		},
	},
})
