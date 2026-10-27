package input

import "github.com/graphql-go/graphql"

// CreateProductInput defines the fields required to create a product.
var CreateProductInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateProductInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"description": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"price": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"stock": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
})

// UpdateProductInput defines the editable fields for an existing product.
var UpdateProductInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateProductInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"description": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"price": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"stock": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
})

// ProductQueryInput defines optional product filters, sorting, and pagination.
var ProductQueryInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ProductQueryInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"page": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"limit": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"minPrice": &graphql.InputObjectFieldConfig{
			Type: graphql.Float,
		},
		"maxPrice": &graphql.InputObjectFieldConfig{
			Type: graphql.Float,
		},
		"inStock": &graphql.InputObjectFieldConfig{
			Type: graphql.Boolean,
		},
		"sortBy": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"sortOrder": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})
