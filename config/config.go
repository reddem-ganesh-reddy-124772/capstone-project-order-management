package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config stores configuration required by application services.
type Config struct {
	AppEnv   string
	LogLevel string

	OrderServicePort   string
	ProductServicePort string
	AuthServicePort    string
	APIGatewayPort     string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	OrderDBName   string
	ProductDBName string

	ConsulAddress    string
	ConsulScheme     string
	ConsulDatacenter string

	OrderServiceURL   string
	ProductServiceURL string
	AuthServiceURL    string

	HTTPClientTimeoutSeconds int
	HTTPClientMaxRetries     int

	CircuitBreakerFailureThreshold    int
	CircuitBreakerResetTimeoutSeconds int

	JWTSecret            string
	JWTExpirationMinutes int

	AuthAdminUsername string
	AuthAdminPassword string
	AuthUserUsername  string
	AuthUserPassword  string
}

// consulConfig stores non-secret settings loaded from Consul KV.
type consulConfig struct {
	OrderServiceURL   string `json:"order_service_url"`
	ProductServiceURL string `json:"product_service_url"`
	AuthServiceURL    string `json:"auth_service_url"`

	HTTPClientTimeoutSeconds int `json:"http_client_timeout_seconds"`
	HTTPClientMaxRetries     int `json:"http_client_max_retries"`

	CircuitBreakerFailureThreshold    int `json:"circuit_breaker_failure_threshold"`
	CircuitBreakerResetTimeoutSeconds int `json:"circuit_breaker_reset_timeout_seconds"`
}

// Load reads configuration from .env/environment variables and then
// applies non-secret overrides from Consul KV when available.
func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		LogLevel: strings.ToUpper(getEnv("LOG_LEVEL", "DEBUG")),

		OrderServicePort:   getEnv("ORDER_SERVICE_PORT", "8081"),
		ProductServicePort: getEnv("PRODUCT_SERVICE_PORT", "8082"),
		AuthServicePort:    getEnv("AUTH_SERVICE_PORT", "8083"),
		APIGatewayPort:     getEnv("API_GATEWAY_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		OrderDBName:   getEnv("ORDER_DB_NAME", "order_db"),
		ProductDBName: getEnv("PRODUCT_DB_NAME", "product_db"),

		ConsulAddress:    getEnv("CONSUL_ADDRESS", "127.0.0.1:8500"),
		ConsulScheme:     getEnv("CONSUL_SCHEME", "http"),
		ConsulDatacenter: getEnv("CONSUL_DATACENTER", "dc1"),

		OrderServiceURL:   getEnv("ORDER_SERVICE_URL", "http://localhost:8081"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "http://localhost:8082"),
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "http://localhost:8083"),

		HTTPClientTimeoutSeconds: getIntEnv("HTTP_CLIENT_TIMEOUT_SECONDS", 5),
		HTTPClientMaxRetries:     getIntEnv("HTTP_CLIENT_MAX_RETRIES", 3),

		CircuitBreakerFailureThreshold:    getIntEnv("CIRCUIT_BREAKER_FAILURE_THRESHOLD", 5),
		CircuitBreakerResetTimeoutSeconds: getIntEnv("CIRCUIT_BREAKER_RESET_TIMEOUT_SECONDS", 30),

		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTExpirationMinutes: getIntEnv("JWT_EXPIRATION_MINUTES", 60),

		AuthAdminUsername: getEnv("AUTH_ADMIN_USERNAME", ""),
		AuthAdminPassword: getEnv("AUTH_ADMIN_PASSWORD", ""),
		AuthUserUsername:  getEnv("AUTH_USER_USERNAME", ""),
		AuthUserPassword:  getEnv("AUTH_USER_PASSWORD", ""),
	}

	applyConsulConfig(&cfg)

	return cfg
}

// applyConsulConfig applies non-secret settings from the Consul KV key
// order-management/config. If Consul is unavailable, existing settings remain.
func applyConsulConfig(cfg *Config) {
	if cfg == nil {
		return
	}

	address := strings.TrimSpace(cfg.ConsulAddress)
	if address == "" {
		return
	}

	scheme := strings.TrimSpace(cfg.ConsulScheme)
	if scheme == "" {
		scheme = "http"
	}

	consulURL := scheme + "://" + address +
		"/v1/kv/order-management/config?raw"

	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	request, err := http.NewRequest(http.MethodGet, consulURL, nil)
	if err != nil {
		return
	}

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return
	}

	var data consulConfig

	// Consul's raw KV response is the stored value. The normal KV response
	// is base64-encoded, so this function requests ?raw above.
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		// If the value is not JSON, try decoding the normal base64 KV format
		// in the helper below only when the response contains valid KV JSON.
		return
	}

	if strings.TrimSpace(data.OrderServiceURL) != "" {
		cfg.OrderServiceURL = strings.TrimRight(data.OrderServiceURL, "/")
	}
	if strings.TrimSpace(data.ProductServiceURL) != "" {
		cfg.ProductServiceURL = strings.TrimRight(data.ProductServiceURL, "/")
	}
	if strings.TrimSpace(data.AuthServiceURL) != "" {
		cfg.AuthServiceURL = strings.TrimRight(data.AuthServiceURL, "/")
	}

	if data.HTTPClientTimeoutSeconds > 0 {
		cfg.HTTPClientTimeoutSeconds = data.HTTPClientTimeoutSeconds
	}
	if data.HTTPClientMaxRetries >= 0 {
		cfg.HTTPClientMaxRetries = data.HTTPClientMaxRetries
	}
	if data.CircuitBreakerFailureThreshold > 0 {
		cfg.CircuitBreakerFailureThreshold = data.CircuitBreakerFailureThreshold
	}
	if data.CircuitBreakerResetTimeoutSeconds > 0 {
		cfg.CircuitBreakerResetTimeoutSeconds = data.CircuitBreakerResetTimeoutSeconds
	}
}

// getEnv returns an environment variable or its fallback value.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists || strings.TrimSpace(value) == "" {
		return fallback
	}

	return strings.TrimSpace(value)
}

// getIntEnv reads an integer environment variable.
func getIntEnv(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

// DatabaseURL builds a PostgreSQL connection URL for a database.
func (c Config) DatabaseURL(databaseName string) string {
	query := url.Values{}
	query.Set("sslmode", c.DBSSLMode)

	connectionURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.DBUser, c.DBPassword),
		Host:     net.JoinHostPort(c.DBHost, c.DBPort),
		Path:     "/" + databaseName,
		RawQuery: query.Encode(),
	}

	return connectionURL.String()
}

// Validate checks the database configuration for services that use PostgreSQL.
// The Auth Service does not need to call this function.
func (c Config) Validate() error {
	if strings.TrimSpace(c.DBHost) == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if strings.TrimSpace(c.DBPort) == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if strings.TrimSpace(c.DBUser) == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if strings.TrimSpace(c.DBPassword) == "" ||
		strings.HasPrefix(c.DBPassword, "replace_with_") {
		return fmt.Errorf("set a real local DB_PASSWORD in .env")
	}

	return nil
}

// Keep the base64 package available for decoding values when needed.
var _ = base64.StdEncoding
