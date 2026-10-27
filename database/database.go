package database

import (
	"context"
	"fmt"
	"time"

	"order-management/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Connect opens and verifies a PostgreSQL connection for the specified database.
// Examples: order_db, product_db, and auth_db.
func Connect(cfg config.Config, databaseName string) (*gorm.DB, error) {
	// Validate the database configuration before connecting.
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	if databaseName == "" {
		return nil, fmt.Errorf("database name is required")
	}

	// Open the database using GORM and the PostgreSQL driver.
	db, err := gorm.Open(
		postgres.Open(cfg.DatabaseURL(databaseName)),
		&gorm.Config{
			Logger: gormlogger.Default.LogMode(getGORMLogLevel(cfg.LogLevel)),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open database %q: %w", databaseName, err)
	}

	// Get the underlying connection pool to configure its limits.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to access database connection pool: %w", err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// Verify that PostgreSQL is reachable before returning the connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("failed to connect to database %q: %w", databaseName, err)
	}

	return db, nil
}

// getGORMLogLevel converts the application's log level into a GORM log level.
func getGORMLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "DEBUG":
		return gormlogger.Info
	case "INFO":
		return gormlogger.Warn
	case "WARN":
		return gormlogger.Error
	default:
		return gormlogger.Silent
	}
}
