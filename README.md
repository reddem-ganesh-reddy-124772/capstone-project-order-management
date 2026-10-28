# Order Management System

## 1. Project Overview

The **Order Management System** is a Go-based microservices application designed to manage products, inventory reservations, user authentication, and customer orders.

The application follows a microservices architecture, with an API Gateway providing a single entry point for client requests. The individual services communicate through HTTP APIs and use Consul for service discovery.

## 2. Technology Stack

- **Programming Language:** Go (Golang)
- **Web Framework:** Gin
- **Database:** PostgreSQL
- **ORM:** GORM
- **API Gateway:** Go reverse proxy
- **Authentication:** JWT (JSON Web Tokens)
- **Service Discovery:** Consul
- **API Documentation:** Swagger / OpenAPI
- **API Query Language:** GraphQL
- **Logging:** Structured logging
- **Resilience:** Circuit breaker
- **Deployment:** AWS EC2
- **CI/CD:** AWS CodePipeline, AWS CodeBuild, AWS CodeDeploy

## 3. Architecture

The application consists of four main services.

### API Gateway

- Acts as the single entry point for client requests.
- Routes requests to the appropriate backend service.
- Applies authentication and role-based authorization.
- Supports request ID and correlation ID middleware.
- Provides health-check and request-logging middleware.

### Auth Service

- Handles user authentication.
- Generates JWT tokens.
- Supports token-based authentication for protected endpoints.

### Product Service

- Manages product information.
- Supports product-related API operations.
- Handles inventory reservation operations.
- Provides REST and GraphQL endpoints.

### Order Service

- Manages customer orders and order items.
- Supports order creation, retrieval, confirmation, and cancellation.
- Communicates with the Product Service.
- Supports REST and GraphQL endpoints.

### Supporting Components

- **PostgreSQL:** Stores application data.
- **Consul:** Supports service discovery.
- **Shared Middleware:** Provides request IDs, correlation IDs, request logging, recovery, and health checks.
- **Shared HTTP Client:** Contains HTTP client and circuit breaker components.
- **Shared Logger:** Provides centralized logging functionality.

## 4. Project Structure

```text
Order Management/
├── .env
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── buildspec.yml
├── Procfile
├── consul-config.ps1
├── start-all.ps1
├── cmd/
│   ├── api-gateway/
│   │   └── main.go
│   ├── auth-service/
│   │   └── main.go
│   ├── product-service/
│   │   └── main.go
│   └── order-service/
│       └── main.go
├── config/
│   └── config.go
├── database/
│   └── database.go
├── services/
│   ├── auth/
│   ├── product/
│   │   ├── controller/
│   │   ├── docs/
│   │   ├── dto/
│   │   ├── graphql/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── routes/
│   │   ├── service/
│   │   └── utilities/
│   └── order/
│       ├── controller/
│       ├── docs/
│       ├── dto/
│       ├── graphql/
│       ├── model/
│       ├── repository/
│       ├── routes/
│       ├── service/
│       └── utilities/
└── shared/
    ├── discovery/
    ├── httpclient/
    ├── logger/
    └── middleware/
```

## 5. Prerequisites

Install and configure the following before running the application:

- Go
- PostgreSQL
- Consul
- Git (optional, for source control)
- AWS CLI (for AWS deployment and related operations)

Ensure PostgreSQL and Consul are running and accessible to the application.

## 6. Configuration

Application configuration is managed through environment variables and the configuration package.

1. Open the `.env.example` file.
2. Create or update the `.env` file in the project root.
3. Configure the database connection details.
4. Configure the service ports and Consul address.
5. Configure the JWT secret and logging level.
6. Verify the required variable names in `config/config.go` and `.env.example`.

**Important:** Use a strong JWT secret and never commit real credentials or production secrets to source control.

## 7. Running the Application Locally

Open PowerShell in the project root directory.

### Step 1: Download Go Dependencies

```powershell
go mod download
```

### Step 2: Start the Auth Service

```powershell
go run ./cmd/auth-service
```

### Step 3: Start the Product Service

Open another PowerShell window in the project root and run:

```powershell
go run ./cmd/product-service
```

### Step 4: Start the Order Service

Open another PowerShell window and run:

```powershell
go run ./cmd/order-service
```

### Step 5: Start the API Gateway

Open another PowerShell window and run:

```powershell
go run ./cmd/api-gateway
```

Alternatively, if the existing PowerShell startup script is configured for your environment, you can run:

```powershell
.\start-all.ps1
```

Keep the service windows running while testing the application.

## 8. API Endpoints

The API Gateway provides routes for authentication, products, inventory reservations, and orders. The following paths represent the configured route groups.

| Method | Endpoint | Purpose |
|---|---|---|
| POST | `/api/auth/login` | User login |
| GET | `/api/products` | Retrieve products |
| GET | `/api/products/:id` | Retrieve a product by ID |
| POST | `/api/products` | Create a product (admin) |
| PUT/PATCH | `/api/products/:id` | Update a product (admin) |
| DELETE | `/api/products/:id` | Delete a product (admin) |
| GET | `/api/inventory/reservations` | Retrieve inventory reservations |
| POST | `/api/inventory/reservations` | Create an inventory reservation (admin) |
| POST | `/api/inventory/reservations/release` | Release an inventory reservation (admin) |
| POST | `/api/orders` | Create an order |
| GET | `/api/orders` | Retrieve orders |
| GET | `/api/orders/:id` | Retrieve an order by ID |
| PATCH | `/api/orders/:id/confirm` | Confirm an order |
| PATCH | `/api/orders/:id/cancel` | Cancel an order |

Authentication and authorization requirements apply to protected endpoints. Confirm the exact HTTP methods and route definitions in the service route files before using individual endpoints.

## 9. REST API Documentation

The Product Service and Order Service contain Swagger documentation files.

Documentation-related files:

- `services/product/docs/`
- `services/order/docs/`

The directories contain generated or maintained Swagger specification files. Use the Swagger UI route configured in the application, if enabled, to explore and test REST endpoints interactively.

## 10. GraphQL

GraphQL support is implemented for the Product Service and Order Service.

The GraphQL implementation includes schemas, queries, mutations, types, and handlers.

Use the GraphQL endpoint and UI configured in the corresponding service to test supported queries and mutations. Check the service source code for the exact endpoint and UI paths.

## 11. Middleware and Request Tracing

Shared middleware supports request tracing and operational reliability.

- **Request ID:** Identifies an individual request.
- **Correlation ID:** Helps trace a request across multiple services.
- **Request Logging:** Records information about incoming requests.
- **Recovery Middleware:** Helps recover from panics and prevents an unexpected panic from terminating normal request handling.
- **Health Checks:** Provides health-related endpoints.

Request IDs and correlation IDs are useful when debugging requests that pass through the API Gateway, Order Service, and Product Service.

## 12. Service Discovery

Consul is used for service discovery.

The API Gateway uses service discovery to locate the Auth Service, Product Service, and Order Service. Confirm that Consul is running, the services are registered, and the configured service names and ports match the application configuration.

## 13. Circuit Breaker

The shared HTTP client package contains circuit breaker-related components.

A circuit breaker is a resilience mechanism intended to prevent repeated calls to an unhealthy downstream service. It can temporarily stop requests after repeated failures and allow recovery checks later.

The presence of circuit breaker files alone does not confirm that every service-to-service request or API Gateway route is protected. Verify the integration and run the available tests before describing the circuit breaker as active across the application.

## 14. Logging and Error Handling

The application contains a shared logging package and recovery middleware.

Logs can help investigate:

- Service startup failures
- Database connection issues
- Authentication and authorization failures
- Downstream service errors
- Request tracing problems
- Unexpected application errors

Set the logging level through the supported configuration and avoid logging passwords, JWT secrets, or other sensitive information.

## 15. Build and Test

Run the following commands from the project root.

### Download Dependencies

```powershell
go mod download
```

### Build the Application

```powershell
go build ./cmd/...
```

### Run All Go Tests

```powershell
go test ./...
```

### Test the Shared HTTP Client Package

```powershell
go test ./shared/httpclient -v
```

Resolve any build or test failures before deploying the application.

## 16. AWS Deployment

The project includes AWS deployment configuration files for use with AWS CodePipeline, AWS CodeBuild, and AWS CodeDeploy.

The deployment workflow can be organized into the following stages:

1. **Source:** Retrieve the application source from the configured repository.
2. **Build:** Download Go dependencies, run tests, and build the services.
3. **Deploy:** Deploy the application to the configured EC2 environment.
4. **Verify:** Confirm that the services start successfully and can communicate with PostgreSQL and Consul.

Deployment requires the correct AWS permissions, environment configuration, EC2 setup, database connectivity, service ports, and health checks.

The exact build and deployment behavior depends on the current pipeline configuration, build specification, deployment configuration, and startup scripts.

## 17. Troubleshooting

### Service Does Not Start

- Check the terminal output for errors.
- Confirm that the required environment variables are configured.
- Verify that the selected port is available.
- Ensure that dependencies have been downloaded.

### Database Connection Failure

- Confirm that PostgreSQL is running.
- Verify the database host, port, database name, username, and password.
- Check database permissions and network access.

### Service Discovery Failure

- Confirm that Consul is running.
- Verify service registration and health status.
- Check service names and port configuration.

### Unauthorized API Request

- Obtain a valid JWT through the login endpoint.
- Send the token using the `Authorization: Bearer <token>` header.
- Confirm that the token is valid and that the user has the required role.

### Downstream Service Failure

- Confirm that the target service is running.
- Check Consul registration and network connectivity.
- Review the logs using the request ID or correlation ID.

## 18. Security Considerations

- Keep secrets and credentials out of source control.
- Use strong JWT signing secrets.
- Validate authentication and authorization on protected routes.
- Use HTTPS for production traffic.
- Restrict database access and service ports.
- Avoid exposing internal errors or sensitive configuration to clients.
- Keep dependencies updated and run tests after changes.

## 19. Project Status

This repository contains the core Go microservices, shared middleware, service discovery integration, API documentation files, GraphQL components, and AWS deployment configuration.

Before a production release, validate the complete request flow, database behavior, authentication, role-based authorization, API documentation, GraphQL operations, circuit breaker integration, automated tests, and AWS deployment process.

## 20. Conclusion

The Order Management System demonstrates a Go-based microservices architecture with API Gateway routing, authentication, product and order management, inventory reservation, service discovery, request tracing, REST and GraphQL APIs, and AWS deployment automation.
