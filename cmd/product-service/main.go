package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order-management/config"
	"order-management/database"
	"order-management/services/product/controller"
	_ "order-management/services/product/docs"
	productgraphql "order-management/services/product/graphql"
	"order-management/services/product/model"
	"order-management/services/product/repository"
	"order-management/services/product/routes"
	"order-management/services/product/service"
	"order-management/shared/logger"
	"order-management/shared/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Product Service API
// @version 1.0
// @description REST and GraphQL APIs for managing products.
// @host localhost:8082
// @BasePath /

func main() {
	// Load and validate the application configuration.
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize the shared application logger.
	logger.Init(cfg.LogLevel)
	log := logger.Get()

	// Connect to the separate Product database.
	db, err := database.Connect(cfg, cfg.ProductDBName)
	if err != nil {
		log.Error("Failed to connect to Product database", "error", err)
		os.Exit(1)
	}

	// Create or update the Product and InventoryReservation tables.
	if err := db.AutoMigrate(
		&model.Product{},
		&model.InventoryReservation{},
	); err != nil {
		log.Error("Failed to migrate Product database", "error", err)
		os.Exit(1)
	}

	// Create the Product repository, service, and controller.
	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productController := controller.NewProductController(productService)

	// Create the Inventory Reservation repository.
	inventoryRepository := repository.NewInventoryReservationRepository(db)

	// Create the Inventory Reservation service.
	inventoryService := service.NewInventoryReservationService(
		inventoryRepository,
	)

	// Create the Inventory Reservation controller.
	inventoryController := controller.NewInventoryReservationController(
		inventoryService,
	)

	// Use Gin's release mode outside development.
	if cfg.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create the HTTP router.
	router := gin.New()

	router.Use(middleware.RequestID())
router.Use(middleware.CorrelationID())
router.Use(middleware.RequestLogging())
router.Use(middleware.Recovery())

	// Register the health endpoint.
	router.GET("/health", middleware.HealthCheck())

	// Register the existing Product REST CRUD endpoints.
	routes.RegisterProductRoutes(router, productController)

	// Register the Inventory Reservation endpoints.
	routes.RegisterInventoryReservationRoutes(
		router.Group("/api"),
		inventoryController,
	)

	// Register the generated Swagger documentation UI.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create the Product GraphQL HTTP handler.
	graphQLHandler, err := productgraphql.NewGraphQLHandler(productService)
	if err != nil {
		log.Error("Failed to initialize Product GraphQL", "error", err)
		os.Exit(1)
	}

	// GET supports browser-based GraphQL queries.
	router.GET("/graphql", gin.WrapH(graphQLHandler))

	// POST supports GraphQL queries and mutations.
	router.POST("/graphql", gin.WrapH(graphQLHandler))

	// Configure the Product Service HTTP server.
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.ProductServicePort),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Start the HTTP server in the background.
	serverErrors := make(chan error, 1)
	go func() {
		log.Info("Product Service starting", "port", cfg.ProductServicePort)
		serverErrors <- server.ListenAndServe()
	}()

	// Wait for a server error or a shutdown signal.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Product Service stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	case sig := <-stop:
		log.Info("Shutdown signal received", "signal", sig.String())
	}

	// Allow active requests to finish during graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Graceful shutdown failed", "error", err)

		if closeErr := server.Close(); closeErr != nil {
			log.Error("Failed to close HTTP server", "error", closeErr)
		}
	}

	log.Info("Product Service stopped")
}
