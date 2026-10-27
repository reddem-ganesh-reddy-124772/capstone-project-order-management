package queries

import (
	"context"
	"fmt"
	"strconv"

	"order-management/services/product/dto"
	"order-management/services/product/graphql/input"
	"order-management/services/product/graphql/types"
	"order-management/services/product/service"

	"github.com/graphql-go/graphql"
)

// NewProductQueries creates the read operations supported by Product GraphQL.
func NewProductQueries(productService service.ProductService) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "ProductQuery",
		Fields: graphql.Fields{
			// Get one product using its ID.
			"product": &graphql.Field{
				Type: types.ProductType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					idText := fmt.Sprint(params.Args["id"])
					idValue, err := strconv.ParseUint(idText, 10, 64)
					if err != nil || idValue == 0 {
						return nil, fmt.Errorf("product ID must be a positive number")
					}

					return productService.GetByID(
						contextFromGraphQL(params),
						uint(idValue),
					)
				},
			},

			// Get products using optional filters, sorting, and pagination.
			"products": &graphql.Field{
				Type: types.ProductListType,
				Args: graphql.FieldConfigArgument{
					"filter": &graphql.ArgumentConfig{
						Type: input.ProductQueryInput,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					request := dto.ProductQueryRequest{}

					filter, exists := params.Args["filter"].(map[string]interface{})
					if exists {
						if value, ok := filter["page"].(int); ok {
							request.Page = value
						}
						if value, ok := filter["limit"].(int); ok {
							request.Limit = value
						}
						if value, ok := filter["name"].(string); ok {
							request.Name = value
						}
						if value, ok := filter["minPrice"].(float64); ok {
							request.MinPrice = &value
						}
						if value, ok := filter["maxPrice"].(float64); ok {
							request.MaxPrice = &value
						}
						if value, ok := filter["inStock"].(bool); ok {
							request.InStock = &value
						}
						if value, ok := filter["sortBy"].(string); ok {
							request.SortBy = value
						}
						if value, ok := filter["sortOrder"].(string); ok {
							request.SortOrder = value
						}
					}

					return productService.GetAll(
						contextFromGraphQL(params),
						request,
					)
				},
			},
		},
	})
}

// contextFromGraphQL returns the request context for service calls.
func contextFromGraphQL(params graphql.ResolveParams) context.Context {
	if params.Context != nil {
		return params.Context
	}

	return context.Background()
}
