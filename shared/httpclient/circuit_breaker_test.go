package httpclient

import (
	"testing"
	"time"

	"order-management/config"
)

// TestCircuitBreakerAllowsRequest checks that requests are allowed initially.
func TestCircuitBreakerAllowsRequest(t *testing.T) {
	cfg := config.Config{
		CircuitBreakerFailureThreshold:    3,
		CircuitBreakerResetTimeoutSeconds: 30,
	}

	breaker := NewCircuitBreaker(cfg)

	// A new circuit breaker should allow the first request.
	if err := breaker.Allow(); err != nil {
		t.Fatalf("expected request to be allowed, got error: %v", err)
	}
}

// TestCircuitBreakerOpensAfterFailures checks that too many failures block requests.
func TestCircuitBreakerOpensAfterFailures(t *testing.T) {
	cfg := config.Config{
		CircuitBreakerFailureThreshold:    2,
		CircuitBreakerResetTimeoutSeconds: 30,
	}

	breaker := NewCircuitBreaker(cfg)

	// Record two failures to reach the configured limit.
	breaker.RecordFailure()
	breaker.RecordFailure()

	// The circuit should now reject requests.
	if err := breaker.Allow(); err == nil {
		t.Fatal("expected circuit breaker to block the request")
	}
}

// TestCircuitBreakerResetsAfterSuccess checks that success clears failures.
func TestCircuitBreakerResetsAfterSuccess(t *testing.T) {
	cfg := config.Config{
		CircuitBreakerFailureThreshold:    2,
		CircuitBreakerResetTimeoutSeconds: 30,
	}

	breaker := NewCircuitBreaker(cfg)

	// Record one failure, then a successful service call.
	breaker.RecordFailure()
	breaker.RecordSuccess()

	// One new failure should not open the circuit.
	breaker.RecordFailure()

	if err := breaker.Allow(); err != nil {
		t.Fatalf("expected request to be allowed after success, got: %v", err)
	}
}

// TestCircuitBreakerAllowsRequestAfterTimeout checks that the circuit resets.
func TestCircuitBreakerAllowsRequestAfterTimeout(t *testing.T) {
	cfg := config.Config{
		CircuitBreakerFailureThreshold:    1,
		CircuitBreakerResetTimeoutSeconds: 1,
	}

	breaker := NewCircuitBreaker(cfg)

	// One failure opens the circuit because the limit is one.
	breaker.RecordFailure()

	// Move the opening time back so the timeout has already passed.
	breaker.mu.Lock()
	breaker.openedAt = time.Now().Add(-2 * time.Second)
	breaker.mu.Unlock()

	// The request should be allowed after the reset timeout.
	if err := breaker.Allow(); err != nil {
		t.Fatalf("expected request to be allowed after timeout, got: %v", err)
	}
}
