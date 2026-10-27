package types

import (
	"order-management/services/product/dto"

	"github.com/graphql-go/graphql"
)

// ProductListType describes a paginated list of products returned by GraphQL.
var ProductListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProductList",
	Fields: graphql.Fields{
		"items": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(
				graphql.NewNonNull(ProductType),
			)),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				result, ok := params.Source.(*dto.ProductListResponse)
				if !ok || result == nil {
					return nil, nil
				}

				return result.Items, nil
			},
		},
		"page": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"limit": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"totalRecords": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				result, ok := params.Source.(*dto.ProductListResponse)
				if !ok || result == nil {
					return nil, nil
				}

				return int(result.TotalRecords), nil
			},
		},
		"totalPages": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				result, ok := params.Source.(*dto.ProductListResponse)
				if !ok || result == nil {
					return nil, nil
				}

				return int(result.TotalPages), nil
			},
		},
	},
})
