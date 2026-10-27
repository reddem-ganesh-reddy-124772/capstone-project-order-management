package main

/*
# Test Notes — Order and Inventory Reservation

## Test 1: Create an Order and Reserve Stock

**Step 1: Check initial product stock**

- Open Product Service Swagger: `http://localhost:8082/swagger/index.html`
- Use `GET /api/products/{id}`.
- Note the product ID and initial stock quantity.

**Step 2: Create an order**

- Open Order Service Swagger: `http://localhost:8081/swagger/index.html`
- Use `POST /api/orders`.
- Pass a valid product ID and quantity in the request body, along with any other required fields.

**Expected result:**
- Order is created successfully.
- Product stock decreases by the ordered quantity.
- An inventory reservation is created.

## Test 2: Verify the Reservation

In Product Service Swagger, use:

`GET /api/inventory/reservations`

**Expected result:**
- The new reservation appears.
- Its reservation ID, product ID, quantity, and status are correct.
- Status is `Reserved` (or the equivalent status defined in the application).

## Test 3: Cancel the Order and Restore Stock

In Order Service Swagger, use:

`PATCH /api/orders/{id}/cancel`

Replace `{id}` with the created order's ID.

**Expected result:**
- Order is cancelled.
- Reservation status changes to `Released`.
- Product stock returns to its original quantity.

## Test 4: Cancel the Same Order Again

Send the same cancellation request again.

**Expected result:**
- The API returns an error because the order is already cancelled.
- Product stock does not increase again.
- Reservation remains released.

## Final Verification

The test passes when the order creation reserves stock, cancellation releases it, stock returns to its original quantity, and repeated cancellation does not restore stock twice.

One important point: these tests verify both service-to-service communication and duplicate stock-restoration protection.
*/