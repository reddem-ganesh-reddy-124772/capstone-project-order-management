package main

/*
http://localhost:8082/swagger/index.html**?utm_source=chatgpt.com

Okay. We'll test all Product REST API test cases using Swagger UI first. We won't move to GraphQL until the REST tests are complete.

Open Swagger UI:

[http://localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html?utm_source=chatgpt.com)

Make sure the Product Service is running on port `8082` and the `products` table exists. In Swagger, expand each endpoint and click Try it out, enter the parameters or request body, then click Execute.

# Product Service — REST API test cases

## 1. Create product — `POST /api/products`

Use this endpoint to create products for the remaining tests.

POST

`/api/products`

Test 1: Create a valid product

Request body:

```
{
  "name": "Laptop",
  "description": "Business laptop",
  "price": 55000,
  "stock": 10
}
```

Expected: 201 Created. The response should include the generated `id` and product details.

Test 2: Create another product

```
{
  "name": "Mouse",
  "description": "Wireless mouse",
  "price": 800,
  "stock": 25
}
```

Expected: 201 Created. Save the returned ID.

Test 3: Create a product with zero stock

```
{
  "name": "Keyboard",
  "description": "USB keyboard",
  "price": 1200,
  "stock": 0
}
```

Expected: 201 Created, because zero stock is allowed.

Test 4: Empty product name

```
{
  "name": "",
  "description": "Invalid product",
  "price": 100,
  "stock": 5
}
```

Expected: 400 Bad Request.

Test 5: Negative price

```
{
  "name": "Invalid Item",
  "description": "Negative price",
  "price": -100,
  "stock": 5
}
```

Expected: 400 Bad Request.

Test 6: Negative stock

```
{
  "name": "Invalid Stock Item",
  "description": "Negative stock",
  "price": 100,
  "stock": -5
}
```

Expected: 400 Bad Request.

Test 7: Missing required fields

```
{
  "description": "Missing name and price",
  "stock": 5
}
```

Expected: 400 Bad Request.

Test 8: Invalid JSON

```
{
  "name": "Laptop",
  "price":
}
```

Expected: 400 Bad Request.

## 2. List products — `GET /api/products`

GET

`/api/products`

Test 9: Get all products

No parameters.

Expected: 200 OK, with the list and pagination metadata.

Test 10: Default pagination

Leave `page` and `limit` empty.

Expected: Page 1 and the default page size of 10.

Test 11: Custom pagination

```
page = 1
limit = 2
```

Expected: At most 2 products.

Test 12: Second page

```
page = 2
limit = 2
```

Expected: The next page of results, if available.

Test 13: Filter by product name

```
name = Laptop
```

Expected: Products whose names match `Laptop`; the filter uses case-insensitive partial matching.

Test 14: Minimum price

```
min_price = 1000
```

Expected: Only products priced at or above 1000.

Test 15: Maximum price

```
max_price = 2000
```

Expected: Only products priced at or below 2000.

Test 16: Price range

```
min_price = 500
max_price = 2000
```

Expected: Products within that price range.

Test 17: In-stock products

```
in_stock = true
```

Expected: Products with stock greater than zero.

Test 18: Out-of-stock products

```
in_stock = false
```

Expected: Products with stock equal to zero.

Test 19: Sort by price ascending

```
sort_by = price
sort_order = asc
```

Expected: Lowest-priced products first.

Test 20: Sort by price descending

```
sort_by = price
sort_order = desc
```

Expected: Highest-priced products first.

Test 21: Sort by name

```
sort_by = name
sort_order = asc
```

Expected: Product names in ascending alphabetical order.

Test 22: Combine filters, sorting, and pagination

```
name = a
min_price = 500
max_price = 60000
in_stock = true
page = 1
limit = 5
sort_by = price
sort_order = desc
```

Expected: Matching products, ordered by price descending, limited to five results on page 1.

Test 23: No matching products

```
name = product-that-does-not-exist
```

Expected: 200 OK with an empty `items` list and zero matching records.

Test 24: Invalid page

```
page = 0
```

Expected: 400 Bad Request.

Test 25: Invalid page size

```
limit = 101
```

Expected: 400 Bad Request.

Test 26: Minimum price greater than maximum

```
min_price = 5000
max_price = 1000
```

Expected: 400 Bad Request.

Test 27: Invalid sort field

```
sort_by = invalid_column
```

Expected: 400 Bad Request.

Test 28: Invalid sort direction

```
sort_order = sideways
```

Expected: 400 Bad Request.

## 3. Get product by ID — `GET /api/products/{id}`

Use the actual product IDs returned by your successful create requests.

GET

`/api/products/{id}`

Test 29: Existing product

```
id = the ID returned when creating Laptop
```

Expected: 200 OK with the product details.

Test 30: Nonexistent product

```
id = 999999
```

Expected: 404 Not Found, assuming that ID does not exist.

Test 31: Invalid ID

```
id = abc
```

Expected: 400 Bad Request.

Test 32: Zero ID

```
id = 0
```

Expected: 400 Bad Request.

## 4. Update product — `PUT /api/products/{id}`

Use the ID of the Laptop created in Test 1.

PUT

`/api/products/{id}`

Test 33: Update a valid product

```
{
  "name": "Laptop Pro",
  "description": "Updated business laptop",
  "price": 65000,
  "stock": 15
}
```

Expected: 200 OK, with the updated details.

Test 34: Verify the update

Call `GET /api/products/{id}` again.

Expected: The updated name, price, description, and stock.

Test 35: Update with an invalid ID

```
id = abc
```

Expected: 400 Bad Request.

Test 36: Update a nonexistent product

```
id = 999999
```

Use a valid request body.

Expected: 404 Not Found, assuming that ID does not exist.

Test 37: Update with an empty name

Use the existing Laptop ID and this body:

```
{
  "name": "",
  "description": "Invalid update",
  "price": 65000,
  "stock": 15
}
```

Expected: 400 Bad Request.

Test 38: Update with a negative price

```
{
  "name": "Laptop Pro",
  "description": "Invalid price",
  "price": -10,
  "stock": 15
}
```

Expected: 400 Bad Request.

Test 39: Update with negative stock

```
{
  "name": "Laptop Pro",
  "description": "Invalid stock",
  "price": 65000,
  "stock": -1
}
```

Expected: 400 Bad Request.

Test 40: Update with invalid JSON

Submit a malformed JSON body.

Expected: 400 Bad Request.

## 5. Delete product — `DELETE /api/products/{id}`

Use a product ID that you no longer need, such as the Keyboard ID from Test 3.

DELETE

`/api/products/{id}`

Test 41: Delete an existing product

Enter the Keyboard ID.

Expected: 200 OK, with a success message.

Test 42: Verify deletion

Call `GET /api/products/{id}` using the deleted ID.

Expected: 404 Not Found.

Test 43: Delete the same product again

Call DELETE again using the same ID.

Expected: 404 Not Found.

Test 44: Delete a nonexistent product

Use an ID that does not exist, such as `999999`.

Expected: 404 Not Found.

Test 45: Delete using an invalid ID

Enter `abc`.

Expected: 400 Bad Request.

## 6. Health check — `GET /health`

GET

`/health`

Test 46: Check service health

Open `http://localhost:8082/health`.

Expected: 200 OK, with a response similar to:

```
{
  "status": "UP"
}
```

This is a health endpoint, not a Product CRUD endpoint. It may not appear in Swagger unless we document it separately.

## Test tracking checklist

Progress

### 0 / 6 groups

Create product Tests 1–8

List, filter, sort and paginate Tests 9–28

Get product by ID Tests 29–32

Update product Tests 33–40

Delete product Tests 41–45

Health check Tests 46

Important: These are the expected results based on your current implementation. If a test produces a different status code, mark it as failed and share the endpoint, request parameters/body, actual status code, and response. In particular, run the invalid-input tests rather than assuming they pass.

Start with Tests 1–8 in Swagger. Once those pass, continue with Tests 9–28, and so on. We will finish and fix the REST API tests before touching GraphQL.
*/