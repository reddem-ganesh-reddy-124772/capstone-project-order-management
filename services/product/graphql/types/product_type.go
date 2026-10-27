package types

import (
	"fmt"
	"strconv"

	"order-management/services/product/dto"

	"github.com/graphql-go/graphql"
)

// ProductType defines the product fields available through GraphQL.
var ProductType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveProductValue(params, func(product *dto.ProductResponse) interface{} {
					return strconv.FormatUint(uint64(product.ID), 10)
				})
			},
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveProductValue(params, func(product *dto.ProductResponse) interface{} {
					return product.Name
				})
			},
		},
		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveProductValue(params, func(product *dto.ProductResponse) interface{} {
					return product.Description
				})
			},
		},
		"price": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Float),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveProductValue(params, func(product *dto.ProductResponse) interface{} {
					return product.Price
				})
			},
		},
		"stock": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveProductValue(params, func(product *dto.ProductResponse) interface{} {
					return product.Stock
				})
			},
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveProductValue(params, func(product *dto.ProductResponse) interface{} {
					return product.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
				})
			},
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveProductValue(params, func(product *dto.ProductResponse) interface{} {
					return product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
				})
			},
		},
	},
})

// resolveProductValue accepts either a ProductResponse pointer or value.
func resolveProductValue(
	params graphql.ResolveParams,
	getValue func(*dto.ProductResponse) interface{},
) (interface{}, error) {
	switch product := params.Source.(type) {
	case *dto.ProductResponse:
		if product != nil {
			return getValue(product), nil
		}
	case dto.ProductResponse:
		return getValue(&product), nil
	}

	return nil, fmt.Errorf(
		"invalid product data received by GraphQL: expected dto.ProductResponse, received %T",
		params.Source,
	)
}
