package mutations

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

// NewProductMutations creates the write operations supported by Product GraphQL.
func NewProductMutations(productService service.ProductService) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "ProductMutation",
		Fields: graphql.Fields{
			// Create a product using the supplied input fields.
			"createProduct": &graphql.Field{
				Type: types.ProductType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(input.CreateProductInput),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					values, ok := params.Args["input"].(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("product input is required")
					}

					request := dto.CreateProductRequest{
						Name:        stringValue(values, "name"),
						Description: stringValue(values, "description"),
						Price:       floatValue(values, "price"),
						Stock:       intValue(values, "stock"),
					}

					return productService.Create(graphQLContext(params), request)
				},
			},

			// Update a product identified by its ID.
			"updateProduct": &graphql.Field{
				Type: types.ProductType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(input.UpdateProductInput),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, err := parseProductID(params.Args["id"])
					if err != nil {
						return nil, err
					}

					values, ok := params.Args["input"].(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("product input is required")
					}

					request := dto.UpdateProductRequest{
						Name:        stringValue(values, "name"),
						Description: stringValue(values, "description"),
						Price:       floatValue(values, "price"),
						Stock:       intValue(values, "stock"),
					}

					return productService.Update(graphQLContext(params), id, request)
				},
			},

			// Delete a product identified by its ID.
			"deleteProduct": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, err := parseProductID(params.Args["id"])
					if err != nil {
						return false, err
					}

					if err := productService.Delete(graphQLContext(params), id); err != nil {
						return false, err
					}

					return true, nil
				},
			},
		},
	})
}

// parseProductID converts a GraphQL ID to a positive numeric product ID.
func parseProductID(value interface{}) (uint, error) {
	id, err := strconv.ParseUint(fmt.Sprint(value), 10, 64)
	if err != nil || id == 0 {
		return 0, fmt.Errorf("product ID must be a positive number")
	}

	return uint(id), nil
}

// graphQLContext returns the request context for service calls.
func graphQLContext(params graphql.ResolveParams) context.Context {
	if params.Context != nil {
		return params.Context
	}

	return context.Background()
}

// stringValue safely reads a string field from GraphQL input.
func stringValue(values map[string]interface{}, key string) string {
	value, _ := values[key].(string)
	return value
}

// floatValue safely reads a numeric field from GraphQL input.
func floatValue(values map[string]interface{}, key string) float64 {
	switch value := values[key].(type) {
	case float64:
		return value
	case int:
		return float64(value)
	default:
		return 0
	}
}

// intValue safely reads an integer field from GraphQL input.
func intValue(values map[string]interface{}, key string) int {
	value, _ := values[key].(int)
	return value
}

// Ensure the expected request DTO fields remain aligned with GraphQL inputs.
var _ = dto.CreateProductRequest{}
