package main

import (
	"errors"
	"log"
"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"order-management/config"
	"order-management/shared/discovery"
	"order-management/shared/logger"
	"order-management/shared/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sony/gobreaker"
)

const (
	serviceTimeout      = 5 * time.Second
	circuitOpenTimeout  = 15 * time.Second
	circuitFailureLimit = 3
)

var (
	breakers   = make(map[string]*gobreaker.CircuitBreaker)
	breakersMu sync.RWMutex
)

func main() {
	cfg := config.Load()

	// Discover service URLs through Consul.
	// If discovery fails, retain the URLs configured in .env.
	consulClient := discovery.New(cfg)

	if serviceURL, err := consulClient.DiscoverService("auth-service"); err == nil &&
		strings.TrimSpace(serviceURL) != "" {
		cfg.AuthServiceURL = serviceURL
		log.Println("Auth Service discovered through Consul")
	} else {
		log.Println("Using configured Auth Service URL")
	}

	if serviceURL, err := consulClient.DiscoverService("product-service"); err == nil &&
		strings.TrimSpace(serviceURL) != "" {
		cfg.ProductServiceURL = serviceURL
		log.Println("Product Service discovered through Consul")
	} else {
		log.Println("Using configured Product Service URL")
	}

	if serviceURL, err := consulClient.DiscoverService("order-service"); err == nil &&
		strings.TrimSpace(serviceURL) != "" {
		cfg.OrderServiceURL = serviceURL
		log.Println("Order Service discovered through Consul")
	} else {
		log.Println("Using configured Order Service URL")
	}

	if cfg.JWTSecret == "" ||
		cfg.JWTSecret == "replace_with_a_long_random_secret" {
		log.Fatal("JWT_SECRET must be configured in .env")
	}

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.Init(cfg.LogLevel)
	appLogger := logger.Get()

	router := gin.New()

	// Tracking middleware applies to all Gateway routes.
	router.Use(middleware.RequestID())
	router.Use(middleware.CorrelationID())
	router.Use(middleware.RequestLogging())
	router.Use(middleware.Recovery())

	router.GET("/health", middleware.HealthCheck())

	// Create independent circuit breakers for each downstream service.
	initCircuitBreaker("auth-service")
	initCircuitBreaker("product-service")
	initCircuitBreaker("order-service")

	// Login is public.
	authProxy := newProxy("auth-service", cfg.AuthServiceURL)
	router.Any("/api/auth/login", authProxy)

	// Both Admin and User can view products.
	productProxy := newProxy("product-service", cfg.ProductServiceURL)
	router.GET("/api/products", AuthMiddleware(cfg), productProxy)
	router.GET("/api/products/:id", AuthMiddleware(cfg), productProxy)

	// Only Admin can create, update, or delete products.
	router.POST("/api/products", AuthMiddleware(cfg), RequireRole("admin"), productProxy)
	router.PUT("/api/products/:id", AuthMiddleware(cfg), RequireRole("admin"), productProxy)
	router.PATCH("/api/products/:id", AuthMiddleware(cfg), RequireRole("admin"), productProxy)
	router.DELETE("/api/products/:id", AuthMiddleware(cfg), RequireRole("admin"), productProxy)

	// Both roles can view inventory reservations.
	inventoryProxy := newProxy("product-service", cfg.ProductServiceURL)
	router.GET("/api/inventory/reservations", AuthMiddleware(cfg), inventoryProxy)
	router.GET("/api/inventory/reservations/:reservation_id", AuthMiddleware(cfg), inventoryProxy)

	// Only Admin can directly reserve or release inventory.
	router.POST("/api/inventory/reservations", AuthMiddleware(cfg), RequireRole("admin"), inventoryProxy)
	router.POST("/api/inventory/reservations/release", AuthMiddleware(cfg), RequireRole("admin"), inventoryProxy)

	// Both roles can create and view orders.
	orderProxy := newProxy("order-service", cfg.OrderServiceURL)
	router.POST("/api/orders", AuthMiddleware(cfg), orderProxy)
	router.GET("/api/orders", AuthMiddleware(cfg), orderProxy)
	router.GET("/api/orders/:id", AuthMiddleware(cfg), orderProxy)

	// Both Admin and User can confirm or cancel orders.
	router.PATCH("/api/orders/:id/confirm", AuthMiddleware(cfg), orderProxy)
	router.PATCH("/api/orders/:id/cancel", AuthMiddleware(cfg), orderProxy)

	server := &http.Server{
		Addr:              ":" + cfg.APIGatewayPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	appLogger.Info("API Gateway starting", "port", cfg.APIGatewayPort)

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		appLogger.Error("API Gateway failed", "error", err)
	}
}

// initCircuitBreaker creates a circuit breaker for a downstream service.
func initCircuitBreaker(serviceName string) {
	breakersMu.Lock()
	defer breakersMu.Unlock()

	breakers[serviceName] = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        serviceName,
		MaxRequests: 1,
		Interval:    30 * time.Second,
		Timeout:     circuitOpenTimeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= circuitFailureLimit
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf(
				"Circuit breaker %s changed state: %s -> %s",
				name,
				from,
				to,
			)
		},
	})
}

func getCircuitBreaker(serviceName string) *gobreaker.CircuitBreaker {
	breakersMu.RLock()
	defer breakersMu.RUnlock()
	return breakers[serviceName]
}

// newProxy forwards requests through a service-specific circuit breaker.
func newProxy(serviceName, rawURL string) gin.HandlerFunc {
	target, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || target.Scheme == "" || target.Host == "" {
		log.Fatalf("invalid service URL: %q", rawURL)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// The timeout limits each outbound request.
	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: serviceTimeout}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   serviceTimeout,
		ResponseHeaderTimeout: serviceTimeout,
		ExpectContinueTimeout: 1 * time.Second,
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Gateway proxy error for %s: %v", serviceName, err)
		http.Error(w, "downstream service unavailable", http.StatusBadGateway)
	}

	return func(c *gin.Context) {
		breaker := getCircuitBreaker(serviceName)
		if breaker == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "circuit breaker is not initialized",
			})
			return
		}

		// Run the proxy through the circuit breaker.
		// The response is written inside the proxy, so the breaker wraps
		// the request lifecycle rather than returning an HTTP response.
		_, err := breaker.Execute(func() (interface{}, error) {
			// Use a recorder so we can decide whether the breaker considers
			// the downstream response successful before writing to the client.
			recorder := newResponseRecorder()
			proxy.ServeHTTP(recorder, c.Request)

			if recorder.statusCode >= http.StatusInternalServerError {
				return nil, errors.New("downstream service returned server error")
			}

			c.Status(recorder.statusCode)
			for key, values := range recorder.header {
				for _, value := range values {
					c.Writer.Header().Add(key, value)
				}
			}
			_, _ = c.Writer.Write(recorder.body)
			return nil, nil
		})

		if err != nil {
			if errors.Is(err, gobreaker.ErrOpenState) ||
				errors.Is(err, gobreaker.ErrTooManyRequests) {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error": "downstream service temporarily unavailable",
				})
				return
			}

			// A downstream 5xx response also counts as a breaker failure.
			// The recorder's response is not sent when Execute returns an error.
			if !c.Writer.Written() {
				c.JSON(http.StatusBadGateway, gin.H{
					"error": "downstream service failed",
				})
			}
		}
	}
}

type responseRecorder struct {
	header     http.Header
	body       []byte
	statusCode int
}

func newResponseRecorder() *responseRecorder {
	return &responseRecorder{
		header:     make(http.Header),
		statusCode: http.StatusOK,
		body:       make([]byte, 0),
	}
}

func (r *responseRecorder) Header() http.Header {
	return r.header
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func (r *responseRecorder) Write(body []byte) (int, error) {
	r.body = append(r.body, body...)
	return len(body), nil
}

// AuthMiddleware validates the JWT issued by the Auth Service.
func AuthMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "provide Authorization: Bearer <token>",
			})
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(
			parts[1],
			claims,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("unexpected JWT signing method")
				}
				return []byte(cfg.JWTSecret), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithIssuer("order-management-auth"),
			jwt.WithExpirationRequired(),
		)

		if err != nil || token == nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		username, usernameOK := claims["sub"].(string)
		role, roleOK := claims["role"].(string)

		if !usernameOK || username == "" ||
			!roleOK || (role != "admin" && role != "user") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			return
		}

		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}

func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you do not have permission to access this resource",
			})
			return
		}
		c.Next()
	}
}
