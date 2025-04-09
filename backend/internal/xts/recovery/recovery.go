package recovery

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/trade-execution-platform/backend/internal/xts/errors"
)

// RetryConfig defines configuration for retry mechanisms
type RetryConfig struct {
	MaxRetries      int           // Maximum number of retry attempts
	InitialDelay    time.Duration // Initial delay before first retry
	MaxDelay        time.Duration // Maximum delay between retries
	BackoffFactor   float64       // Factor by which to increase delay after each retry
	RetryableErrors []error       // List of errors that should trigger a retry
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:    3,
		InitialDelay:  500 * time.Millisecond,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
		RetryableErrors: []error{
			errors.ErrConnectionFailed,
			errors.ErrRequestTimeout,
			errors.ErrNetworkFailure,
			errors.ErrAPIUnavailable,
			errors.ErrWebSocketConnFailed,
		},
	}
}

// IsRetryable checks if an error is retryable based on the configuration
func (c *RetryConfig) IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check if error is in the list of retryable errors
	for _, retryableErr := range c.RetryableErrors {
		if errors.Is(err, retryableErr) {
			return true
		}
	}

	// Check if error is a network error
	if errors.IsNetworkError(err) {
		return true
	}

	return false
}

// RetryWithBackoff executes a function with exponential backoff retry
func RetryWithBackoff(ctx context.Context, config *RetryConfig, operation func() error) error {
	var err error
	delay := config.InitialDelay

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Execute the operation
		err = operation()
		
		// If no error or error is not retryable, return immediately
		if err == nil || !config.IsRetryable(err) {
			return err
		}

		// If this was the last attempt, return the error
		if attempt == config.MaxRetries {
			return err
		}

		// Log retry attempt
		log.Printf("Operation failed with error: %v. Retrying in %v (attempt %d/%d)...", 
			err, delay, attempt+1, config.MaxRetries)

		// Wait before retrying
		select {
		case <-ctx.Done():
			// Context was canceled, return context error
			return ctx.Err()
		case <-time.After(delay):
			// Continue with retry
		}

		// Increase delay for next attempt
		delay = time.Duration(float64(delay) * config.BackoffFactor)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return err
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mutex           sync.RWMutex
	failureCount    int
	lastFailureTime time.Time
	state           State
	failureThreshold int
	resetTimeout    time.Duration
}

// State represents the state of the circuit breaker
type State int

const (
	StateClosed State = iota    // Circuit is closed, requests are allowed
	StateOpen                   // Circuit is open, requests are not allowed
	StateHalfOpen               // Circuit is half-open, limited requests are allowed
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:           StateClosed,
		failureThreshold: failureThreshold,
		resetTimeout:    resetTimeout,
	}
}

// Execute executes a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(operation func() error) error {
	if !cb.AllowRequest() {
		return errors.ErrAPIUnavailable
	}

	err := operation()

	cb.RecordResult(err)
	return err
}

// AllowRequest checks if a request should be allowed
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		// Check if reset timeout has elapsed
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			// Transition to half-open state
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			if cb.state == StateOpen {
				cb.state = StateHalfOpen
			}
			cb.mutex.Unlock()
			cb.mutex.RLock()
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return true
	}
}

// RecordResult records the result of an operation
func (cb *CircuitBreaker) RecordResult(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err != nil && !errors.Is(err, errors.ErrAPIRateLimited) {
		// Record failure
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		// Check if threshold is reached
		if cb.state == StateClosed && cb.failureCount >= cb.failureThreshold {
			cb.state = StateOpen
		} else if cb.state == StateHalfOpen {
			cb.state = StateOpen
		}
	} else {
		// Record success
		if cb.state == StateHalfOpen {
			// Reset on success in half-open state
			cb.failureCount = 0
			cb.state = StateClosed
		} else if cb.state == StateClosed {
			// Reset failure count on success
			cb.failureCount = 0
		}
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.failureCount = 0
	cb.state = StateClosed
}

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	mutex       sync.Mutex
	rate        int           // Tokens per second
	burst       int           // Maximum burst size
	tokens      float64       // Current token count
	lastRefill  time.Time     // Last time tokens were refilled
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastRefill: time.Now(),
	}
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens += elapsed * float64(rl.rate)
	if rl.tokens > float64(rl.burst) {
		rl.tokens = float64(rl.burst)
	}
	rl.lastRefill = now

	// Check if we have enough tokens
	if rl.tokens < 1 {
		return false
	}

	// Consume a token
	rl.tokens--
	return true
}

// Wait waits until a request is allowed
func (rl *RateLimiter) Wait(ctx context.Context) error {
	for {
		if rl.Allow() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second / time.Duration(rl.rate)):
			// Try again
		}
	}
}
