package types

import "github.com/graphql-go/graphql"

// OrderType defines the GraphQL output fields for an order.
var OrderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"customerId": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"status": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"totalAmount": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"items": &graphql.Field{
			Type: graphql.NewNonNull(
				graphql.NewList(
					graphql.NewNonNull(OrderItemType),
				),
			),
		},
		"createdAt": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"updatedAt": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})
