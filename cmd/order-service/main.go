package main

import (
	"net/http"
	"strconv"

	"order-management/config"
	"order-management/database"
	"order-management/services/order/controller"
	_ "order-management/services/order/docs"
	graphqlhandler "order-management/services/order/graphql/handler"
	"order-management/services/order/model"
	"order-management/services/order/repository"
	"order-management/services/order/routes"
	"order-management/services/order/service"
	"order-management/shared/discovery"
	"order-management/shared/logger"
	"order-management/shared/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Order Service API
// @version 1.0
// @description REST and GraphQL APIs for managing customer orders.
// @host localhost:8081
// @BasePath /

func main() {
	cfg := config.Load()

	// Initialize the shared structured logger.
	logger.Init(cfg.LogLevel)
	log := logger.Get()

	// Connect to the Order database.
	db, err := database.Connect(cfg, cfg.OrderDBName)
	if err != nil {
		log.Error("Failed to connect to Order database", "error", err)
		return
	}

	// Create the required database tables.
	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Error("Failed to migrate Order database", "error", err)
		return
	}

	// Initialize the repository, Product client, and service.
	orderRepository := repository.NewOrderRepository(db)
	productClient := service.NewHTTPProductClient(cfg)
	orderService := service.NewOrderService(orderRepository, productClient)

	// Initialize the REST controller.
	orderController := controller.NewOrderController(orderService)

	// Initialize the GraphQL handler.
	orderGraphQLHandler, err := graphqlhandler.NewOrderGraphQLHandler(orderService)
	if err != nil {
		log.Error("Failed to create Order GraphQL handler", "error", err)
		return
	}

	router := gin.New()

	// Tracking middleware for direct and Gateway requests.
	router.Use(middleware.RequestID())
	router.Use(middleware.CorrelationID())
	router.Use(middleware.RequestLogging())
	router.Use(middleware.Recovery())

	// Health endpoint used by Consul.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "order-service",
			"status":  "UP",
		})
	})

	// Register Order REST APIs.
	routes.RegisterOrderRoutes(router, orderController)

	// Register GraphQL.
	router.GET("/graphql", gin.WrapH(orderGraphQLHandler))
	router.POST("/graphql", gin.WrapH(orderGraphQLHandler))

	// Register Swagger UI.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Convert the configured port to an integer for Consul registration.
	port, err := strconv.Atoi(cfg.OrderServicePort)
	if err != nil {
		log.Error("Invalid Order Service port", "port", cfg.OrderServicePort, "error", err)
		return
	}

	// Register the Order Service with Consul.
	consulClient := discovery.New(cfg)
	err = consulClient.RegisterService(
		"order-service",
		"order-service-1",
		"127.0.0.1",
		port,
		"http://127.0.0.1:"+cfg.OrderServicePort+"/health",
	)
	if err != nil {
		// Keep the service available even if Consul registration fails.
		log.Error("Failed to register Order Service with Consul", "error", err)
	} else {
		log.Info("Order Service registered with Consul", "service", "order-service")
	}

	address := ":" + cfg.OrderServicePort
	log.Info("Order Service starting", "address", address)
	log.Info("GraphQL UI", "url", "http://localhost"+address+"/graphql")
	log.Info("Swagger UI", "url", "http://localhost"+address+"/swagger/index.html")

	if err := router.Run(address); err != nil {
		log.Error("Failed to start Order Service", "error", err)
	}
}
