package orderexecution

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// Error types
	ErrorTypeValidation    ErrorType = "VALIDATION"
	ErrorTypeConnection    ErrorType = "CONNECTION"
	ErrorTypeAuthentication ErrorType = "AUTHENTICATION"
	ErrorTypeExecution     ErrorType = "EXECUTION"
	ErrorTypeTimeout       ErrorType = "TIMEOUT"
	ErrorTypeSystem        ErrorType = "SYSTEM"
	ErrorTypeUnknown       ErrorType = "UNKNOWN"
)

// ErrorSeverity represents the severity of an error
type ErrorSeverity string

const (
	// Error severities
	ErrorSeverityInfo    ErrorSeverity = "INFO"
	ErrorSeverityWarning ErrorSeverity = "WARNING"
	ErrorSeverityError   ErrorSeverity = "ERROR"
	ErrorSeverityCritical ErrorSeverity = "CRITICAL"
)

// ExecutionError represents an error in the execution engine
type ExecutionError struct {
	Type        ErrorType     `json:"type"`
	Severity    ErrorSeverity `json:"severity"`
	Code        string        `json:"code"`
	Message     string        `json:"message"`
	Details     string        `json:"details,omitempty"`
	OrderID     string        `json:"orderId,omitempty"`
	PortfolioID string        `json:"portfolioId,omitempty"`
	StrategyID  string        `json:"strategyId,omitempty"`
	Timestamp   time.Time     `json:"timestamp"`
	RetryCount  int           `json:"retryCount"`
	RetryDelay  time.Duration `json:"retryDelay"`
	Retryable   bool          `json:"retryable"`
	Source      string        `json:"source"`
	OriginalErr error         `json:"-"`
}

// Error implements the error interface
func (e *ExecutionError) Error() string {
	return fmt.Sprintf("[%s][%s] %s: %s", e.Type, e.Severity, e.Code, e.Message)
}

// Unwrap returns the original error
func (e *ExecutionError) Unwrap() error {
	return e.OriginalErr
}

// WithOrderID adds order ID to the error
func (e *ExecutionError) WithOrderID(orderID string) *ExecutionError {
	e.OrderID = orderID
	return e
}

// WithPortfolioID adds portfolio ID to the error
func (e *ExecutionError) WithPortfolioID(portfolioID string) *ExecutionError {
	e.PortfolioID = portfolioID
	return e
}

// WithStrategyID adds strategy ID to the error
func (e *ExecutionError) WithStrategyID(strategyID string) *ExecutionError {
	e.StrategyID = strategyID
	return e
}

// WithDetails adds details to the error
func (e *ExecutionError) WithDetails(details string) *ExecutionError {
	e.Details = details
	return e
}

// WithRetry configures retry parameters
func (e *ExecutionError) WithRetry(retryable bool, retryCount int, retryDelay time.Duration) *ExecutionError {
	e.Retryable = retryable
	e.RetryCount = retryCount
	e.RetryDelay = retryDelay
	return e
}

// NewExecutionError creates a new execution error
func NewExecutionError(
	errType ErrorType,
	severity ErrorSeverity,
	code string,
	message string,
	originalErr error,
	source string,
) *ExecutionError {
	return &ExecutionError{
		Type:        errType,
		Severity:    severity,
		Code:        code,
		Message:     message,
		Timestamp:   time.Now(),
		RetryCount:  0,
		RetryDelay:  0,
		Retryable:   false,
		Source:      source,
		OriginalErr: originalErr,
	}
}

// Common error codes
const (
	ErrCodeInvalidOrder          = "ERR_INVALID_ORDER"
	ErrCodeInvalidParameter      = "ERR_INVALID_PARAMETER"
	ErrCodeInsufficientMargin    = "ERR_INSUFFICIENT_MARGIN"
	ErrCodePositionLimitExceeded = "ERR_POSITION_LIMIT_EXCEEDED"
	ErrCodeRateLimitExceeded     = "ERR_RATE_LIMIT_EXCEEDED"
	ErrCodeOrderNotFound         = "ERR_ORDER_NOT_FOUND"
	ErrCodeConnectionFailed      = "ERR_CONNECTION_FAILED"
	ErrCodeAuthenticationFailed  = "ERR_AUTHENTICATION_FAILED"
	ErrCodeExecutionFailed       = "ERR_EXECUTION_FAILED"
	ErrCodeTimeout               = "ERR_TIMEOUT"
	ErrCodeInternalError         = "ERR_INTERNAL_ERROR"
)

// ErrorHandler handles errors in the execution engine
type ErrorHandler interface {
	// HandleError handles an error and returns whether it was handled successfully
	HandleError(ctx context.Context, err *ExecutionError) (bool, error)
	
	// LogError logs an error
	LogError(err *ExecutionError)
	
	// ShouldRetry determines if an operation should be retried
	ShouldRetry(err *ExecutionError) bool
	
	// GetRetryDelay returns the delay before the next retry
	GetRetryDelay(err *ExecutionError) time.Duration
}

// DefaultErrorHandler implements the ErrorHandler interface
type DefaultErrorHandler struct {
	maxRetries     int
	baseRetryDelay time.Duration
	logger         Logger
}

// Logger interface for logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

// NewDefaultErrorHandler creates a new default error handler
func NewDefaultErrorHandler(maxRetries int, baseRetryDelay time.Duration, logger Logger) *DefaultErrorHandler {
	return &DefaultErrorHandler{
		maxRetries:     maxRetries,
		baseRetryDelay: baseRetryDelay,
		logger:         logger,
	}
}

// HandleError handles an error
func (h *DefaultErrorHandler) HandleError(ctx context.Context, err *ExecutionError) (bool, error) {
	// Log the error
	h.LogError(err)
	
	// Check if we should retry
	if !h.ShouldRetry(err) {
		return false, err
	}
	
	// Increment retry count
	err.RetryCount++
	
	// Calculate retry delay
	err.RetryDelay = h.GetRetryDelay(err)
	
	// Return true to indicate retry
	return true, nil
}

// LogError logs an error
func (h *DefaultErrorHandler) LogError(err *ExecutionError) {
	// Log based on severity
	switch err.Severity {
	case ErrorSeverityInfo:
		h.logger.Info(err.Error(),
			"type", err.Type,
			"code", err.Code,
			"orderId", err.OrderID,
			"portfolioId", err.PortfolioID,
			"strategyId", err.StrategyID,
			"details", err.Details,
			"retryCount", err.RetryCount,
			"source", err.Source,
		)
	case ErrorSeverityWarning:
		h.logger.Warn(err.Error(),
			"type", err.Type,
			"code", err.Code,
			"orderId", err.OrderID,
			"portfolioId", err.PortfolioID,
			"strategyId", err.StrategyID,
			"details", err.Details,
			"retryCount", err.RetryCount,
			"source", err.Source,
		)
	case ErrorSeverityError:
		h.logger.Error(err.Error(),
			"type", err.Type,
			"code", err.Code,
			"orderId", err.OrderID,
			"portfolioId", err.PortfolioID,
			"strategyId", err.StrategyID,
			"details", err.Details,
			"retryCount", err.RetryCount,
			"source", err.Source,
		)
	case ErrorSeverityCritical:
		h.logger.Fatal(err.Error(),
			"type", err.Type,
			"code", err.Code,
			"orderId", err.OrderID,
			"portfolioId", err.PortfolioID,
			"strategyId", err.StrategyID,
			"details", err.Details,
			"retryCount", err.RetryCount,
			"source", err.Source,
		)
	}
}

// ShouldRetry determines if an operation should be retried
func (h *DefaultErrorHandler) ShouldRetry(err *ExecutionError) bool {
	// Don't retry if not retryable
	if !err.Retryable {
		return false
	}
	
	// Don't retry if max retries exceeded
	if err.RetryCount >= h.maxRetries {
		return false
	}
	
	// Determine if error type is retryable
	switch err.Type {
	case ErrorTypeConnection, ErrorTypeTimeout:
		// Connection and timeout errors are always retryable
		return true
	case ErrorTypeExecution:
		// Execution errors may be retryable depending on the specific error
		return err.Code != ErrCodeInsufficientMargin && 
			   err.Code != ErrCodePositionLimitExceeded &&
			   err.Code != ErrCodeRateLimitExceeded
	case ErrorTypeAuthentication:
		// Authentication errors may be retryable (e.g., token expired)
		return err.Code == ErrCodeAuthenticationFailed
	default:
		// Other error types are generally not retryable
		return false
	}
}

// GetRetryDelay returns the delay before the next retry
func (h *DefaultErrorHandler) GetRetryDelay(err *ExecutionError) time.Duration {
	// Exponential backoff with jitter
	delay := h.baseRetryDelay * time.Duration(1<<uint(err.RetryCount))
	
	// Add jitter (±20%)
	jitter := float64(delay) * 0.2
	delay = time.Duration(float64(delay) - jitter + 2*jitter*float64(time.Now().UnixNano()%100)/100)
	
	return delay
}

// RetryOperation retries an operation with exponential backoff
func RetryOperation(ctx context.Context, handler ErrorHandler, operation func() error) error {
	var lastErr error
	
	for {
		// Execute the operation
		err := operation()
		if err == nil {
			return nil
		}
		
		// Convert to ExecutionError if needed
		var execErr *ExecutionError
		if !errors.As(err, &execErr) {
			// If it's not already an ExecutionError, create a generic one
			execErr = NewExecutionError(
				ErrorTypeUnknown,
				ErrorSeverityError,
				ErrCodeInternalError,
				err.Error(),
				err,
				"RetryOperation",
			)
		}
		
		// Handle the error
		shouldRetry, handledErr := handler.HandleError(ctx, execErr)
		if !shouldRetry {
			return handledErr
		}
		
		// Store the last error
		lastErr = handledErr
		
		// Wait for retry delay
		select {
		case <-time.After(execErr.RetryDelay):
			// Continue with retry
		case <-ctx.Done():
			// Context cancelled
			return ctx.Err()
		}
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name           string
	maxFailures    int
	resetTimeout   time.Duration
	halfOpenMaxOps int
	
	failures       int
	lastFailure    time.Time
	state          CircuitBreakerState
	mutex          sync.RWMutex
}

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState string

const (
	CircuitBreakerClosed    CircuitBreakerState = "CLOSED"
	CircuitBreakerOpen      CircuitBreakerState = "OPEN"
	CircuitBreakerHalfOpen  CircuitBreakerState = "HALF_OPEN"
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, maxFailures int, resetTimeout time.Duration, halfOpenMaxOps int) *CircuitBreaker {
	return &CircuitBreaker{
		name:           name,
		maxFailures:    maxFailures,
		resetTimeout:   resetTimeout,
		halfOpenMaxOps: halfOpenMaxOps,
		state:          CircuitBreakerClosed,
	}
}

// Execute executes an operation with circuit breaker protection
func (cb *CircuitBreaker) Execute(operation func() error) error {
	// Check if circuit is open
	if !cb.AllowRequest() {
		return NewExecutionError(
			ErrorTypeSystem,
			ErrorSeverityError,
			"ERR_CIRCUIT_OPEN",
			fmt.Sprintf("Circuit breaker '%s' is open", cb.name),
			nil,
			"CircuitBreaker",
		)
	}
	
	// Execute the operation
	err := operation()
	
	// Record the result
	if err != nil {
		cb.RecordFailure()
		return err
	}
	
	cb.RecordSuccess()
	return nil
}

// AllowRequest checks if a request should be allowed
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	
	switch cb.state {
	case CircuitBreakerClosed:
		return true
	case CircuitBreakerOpen:
		// Check if reset timeout has elapsed
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			// Transition to half-open state
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			cb.state = CircuitBreakerHalfOpen
			cb.failures = 0
			cb.mutex.Unlock()
			cb.mutex.RLock()
			return true
		}
		return false
	case CircuitBreakerHalfOpen:
		// Allow limited requests in half-open state
		return cb.failures < cb.halfOpenMaxOps
	default:
		return false
	}
}

// RecordSuccess records a successful operation
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	if cb.state == CircuitBreakerHalfOpen {
		// Transition to closed state after successful operation in half-open state
		cb.state = CircuitBreakerClosed
		cb.failures = 0
	}
}

// RecordFailure records a failed operation
func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	cb.failures++
	cb.lastFailure = time.Now()
	
	if cb.state == CircuitBreakerClosed && cb.failures >= cb.maxFailures {
		// Transition to open state
		cb.state = CircuitBreakerOpen
	} else if cb.state == CircuitBreakerHalfOpen {
		// Transition back to open state on failure in half-open state
		cb.state = CircuitBreakerOpen
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.state = CircuitBreakerClosed
	cb.failures = 0
}
