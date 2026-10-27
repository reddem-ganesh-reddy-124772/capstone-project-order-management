package httpclient

import (
	"fmt"
	"sync"
	"time"

	"order-management/config"
)

// CircuitBreaker prevents repeated calls to a service that keeps failing.
type CircuitBreaker struct {
	// mu protects the circuit breaker fields when multiple requests use them.
	mu sync.Mutex

	// failureCount stores the number of consecutive failures.
	failureCount int

	// maxFailures is the number of failures allowed before blocking calls.
	maxFailures int

	// resetTimeout is how long the circuit stays open.
	resetTimeout time.Duration

	// openedAt stores the time when the circuit was opened.
	openedAt time.Time

	// isOpen tells us whether requests are currently blocked.
	isOpen bool
}

// NewCircuitBreaker creates a circuit breaker using application settings.
func NewCircuitBreaker(cfg config.Config) *CircuitBreaker {
	// Use five failures if the configured value is zero or negative.
	maxFailures := cfg.CircuitBreakerFailureThreshold
	if maxFailures <= 0 {
		maxFailures = 5
	}

	// Wait 30 seconds if the configured timeout is invalid.
	resetSeconds := cfg.CircuitBreakerResetTimeoutSeconds
	if resetSeconds <= 0 {
		resetSeconds = 30
	}

	// Convert the timeout from seconds into Go's time duration.
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: time.Duration(resetSeconds) * time.Second,
	}
}

// Allow checks whether a request can be sent to the service.
func (b *CircuitBreaker) Allow() error {
	// Lock the fields so another request cannot change them at the same time.
	b.mu.Lock()

	// Release the lock when this function finishes.
	defer b.mu.Unlock()

	// Allow the request when the circuit is closed.
	if !b.isOpen {
		return nil
	}

	// Check whether the waiting period has finished.
	if time.Since(b.openedAt) >= b.resetTimeout {
		// Close the circuit so another request can be attempted.
		b.isOpen = false

		// Start counting failures again.
		b.failureCount = 0

		return nil
	}

	// Reject the request while the circuit is still open.
	return fmt.Errorf("circuit breaker is open; try again later")
}

// RecordSuccess clears the failure count after a successful service call.
func (b *CircuitBreaker) RecordSuccess() {
	// Lock the fields before changing their values.
	b.mu.Lock()
	defer b.mu.Unlock()

	// A successful call means the service responded correctly.
	b.failureCount = 0

	// Allow future requests.
	b.isOpen = false
}

// RecordFailure counts a failed service call.
func (b *CircuitBreaker) RecordFailure() {
	// Lock the fields before updating the failure count.
	b.mu.Lock()
	defer b.mu.Unlock()

	// Increase the number of recorded failures by one.
	b.failureCount++

	// Open the circuit when the failure limit is reached.
	if b.failureCount >= b.maxFailures {
		b.isOpen = true

		// Store the current time to calculate when requests can resume.
		b.openedAt = time.Now()
	}
}
