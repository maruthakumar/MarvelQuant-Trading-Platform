package orderexecution

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ConsoleLogger implements the Logger interface for console logging
type ConsoleLogger struct{}

// Debug logs a debug message
func (l *ConsoleLogger) Debug(msg string, fields ...interface{}) {
	l.log("DEBUG", msg, fields...)
}

// Info logs an info message
func (l *ConsoleLogger) Info(msg string, fields ...interface{}) {
	l.log("INFO", msg, fields...)
}

// Warn logs a warning message
func (l *ConsoleLogger) Warn(msg string, fields ...interface{}) {
	l.log("WARN", msg, fields...)
}

// Error logs an error message
func (l *ConsoleLogger) Error(msg string, fields ...interface{}) {
	l.log("ERROR", msg, fields...)
}

// Fatal logs a fatal message
func (l *ConsoleLogger) Fatal(msg string, fields ...interface{}) {
	l.log("FATAL", msg, fields...)
}

// log logs a message with the given level
func (l *ConsoleLogger) log(level string, msg string, fields ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] [%s] %s", timestamp, level, msg)
	
	if len(fields) > 0 {
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				fmt.Printf(" %v=%v", fields[i], fields[i+1])
			}
		}
	}
	fmt.Println()
}

// NewConsoleLogger creates a new console logger
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

// ErrorHandlerMiddleware is middleware for handling errors in the execution engine
type ErrorHandlerMiddleware struct {
	handler ErrorHandler
	next    ExecutionMiddleware
}

// ExecutionMiddleware is middleware for the execution engine
type ExecutionMiddleware interface {
	Execute(ctx context.Context, order Order) (OrderResponse, error)
}

// NewErrorHandlerMiddleware creates a new error handler middleware
func NewErrorHandlerMiddleware(handler ErrorHandler, next ExecutionMiddleware) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		handler: handler,
		next:    next,
	}
}

// Execute executes an order with error handling
func (m *ErrorHandlerMiddleware) Execute(ctx context.Context, order Order) (OrderResponse, error) {
	var response OrderResponse
	var err error
	
	// Define the operation to retry
	operation := func() error {
		var opErr error
		response, opErr = m.next.Execute(ctx, order)
		return opErr
	}
	
	// Retry the operation
	err = RetryOperation(ctx, m.handler, operation)
	
	return response, err
}

// DeadLetterQueue handles orders that failed execution
type DeadLetterQueue struct {
	orders map[string]Order
	mutex  sync.RWMutex
	logger Logger
}

// NewDeadLetterQueue creates a new dead letter queue
func NewDeadLetterQueue(logger Logger) *DeadLetterQueue {
	return &DeadLetterQueue{
		orders: make(map[string]Order),
		logger: logger,
	}
}

// AddOrder adds an order to the dead letter queue
func (q *DeadLetterQueue) AddOrder(order Order, err error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	
	// Update order status and error message
	order.Status = OrderStatusFailed
	order.ErrorMessage = err.Error()
	order.UpdatedAt = time.Now()
	
	// Add to queue
	q.orders[order.ID] = order
	
	// Log the event
	q.logger.Error("Order added to dead letter queue",
		"orderId", order.ID,
		"portfolioId", order.PortfolioID,
		"strategyId", order.StrategyID,
		"symbol", order.Symbol,
		"error", err.Error(),
	)
}

// GetOrders returns all orders in the dead letter queue
func (q *DeadLetterQueue) GetOrders() []Order {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	
	orders := make([]Order, 0, len(q.orders))
	for _, order := range q.orders {
		orders = append(orders, order)
	}
	
	return orders
}

// GetOrder returns an order from the dead letter queue
func (q *DeadLetterQueue) GetOrder(orderID string) (Order, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	
	order, ok := q.orders[orderID]
	if !ok {
		return Order{}, fmt.Errorf("order not found in dead letter queue: %s", orderID)
	}
	
	return order, nil
}

// RemoveOrder removes an order from the dead letter queue
func (q *DeadLetterQueue) RemoveOrder(orderID string) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	
	if _, ok := q.orders[orderID]; !ok {
		return fmt.Errorf("order not found in dead letter queue: %s", orderID)
	}
	
	delete(q.orders, orderID)
	
	// Log the event
	q.logger.Info("Order removed from dead letter queue",
		"orderId", orderID,
	)
	
	return nil
}

// RetryOrder retries an order from the dead letter queue
func (q *DeadLetterQueue) RetryOrder(ctx context.Context, orderID string, engine *OrderExecutionEngine) error {
	// Get the order
	order, err := q.GetOrder(orderID)
	if err != nil {
		return err
	}
	
	// Remove from queue
	if err := q.RemoveOrder(orderID); err != nil {
		return err
	}
	
	// Reset order status
	order.Status = OrderStatusPending
	order.ErrorMessage = ""
	order.UpdatedAt = time.Now()
	
	// Submit for execution
	_, err = engine.ExecuteOrder(ctx, order)
	if err != nil {
		// If execution fails again, add back to queue
		q.AddOrder(order, err)
		return err
	}
	
	return nil
}

// RetryAllOrders retries all orders in the dead letter queue
func (q *DeadLetterQueue) RetryAllOrders(ctx context.Context, engine *OrderExecutionEngine) (int, int, error) {
	orders := q.GetOrders()
	
	successCount := 0
	failureCount := 0
	
	for _, order := range orders {
		err := q.RetryOrder(ctx, order.ID, engine)
		if err != nil {
			failureCount++
		} else {
			successCount++
		}
	}
	
	return successCount, failureCount, nil
}

// EnhancedOrderExecutionEngine extends OrderExecutionEngine with error handling
type EnhancedOrderExecutionEngine struct {
	*OrderExecutionEngine
	errorHandler ErrorHandler
	deadLetterQueue *DeadLetterQueue
	circuitBreakers map[string]*CircuitBreaker
	mutex sync.RWMutex
}

// NewEnhancedOrderExecutionEngine creates a new enhanced order execution engine
func NewEnhancedOrderExecutionEngine(
	engine *OrderExecutionEngine,
	errorHandler ErrorHandler,
	deadLetterQueue *DeadLetterQueue,
) *EnhancedOrderExecutionEngine {
	return &EnhancedOrderExecutionEngine{
		OrderExecutionEngine: engine,
		errorHandler:         errorHandler,
		deadLetterQueue:      deadLetterQueue,
		circuitBreakers:      make(map[string]*CircuitBreaker),
	}
}

// RegisterCircuitBreaker registers a circuit breaker
func (e *EnhancedOrderExecutionEngine) RegisterCircuitBreaker(name string, cb *CircuitBreaker) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.circuitBreakers[name] = cb
}

// GetCircuitBreaker gets a circuit breaker by name
func (e *EnhancedOrderExecutionEngine) GetCircuitBreaker(name string) (*CircuitBreaker, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	cb, ok := e.circuitBreakers[name]
	if !ok {
		return nil, fmt.Errorf("circuit breaker not found: %s", name)
	}
	
	return cb, nil
}

// ExecuteOrderWithErrorHandling executes an order with error handling
func (e *EnhancedOrderExecutionEngine) ExecuteOrderWithErrorHandling(ctx context.Context, order Order) (OrderResponse, error) {
	// Get the appropriate circuit breaker
	cbName := "default"
	if order.Exchange != "" {
		cbName = order.Exchange
	}
	
	cb, err := e.GetCircuitBreaker(cbName)
	if err != nil {
		// If circuit breaker doesn't exist, create a default one
		cb = NewCircuitBreaker(cbName, 5, 30*time.Second, 1)
		e.RegisterCircuitBreaker(cbName, cb)
	}
	
	// Execute with circuit breaker
	var response OrderResponse
	err = cb.Execute(func() error {
		var execErr error
		response, execErr = e.OrderExecutionEngine.ExecuteOrder(ctx, order)
		
		// Handle specific error types
		if execErr != nil {
			var executionError *ExecutionError
			if !errors.As(execErr, &executionError) {
				// Convert to ExecutionError
				executionError = NewExecutionError(
					ErrorTypeExecution,
					ErrorSeverityError,
					ErrCodeExecutionFailed,
					execErr.Error(),
					execErr,
					"EnhancedOrderExecutionEngine",
				).WithOrderID(order.ID).WithPortfolioID(order.PortfolioID).WithStrategyID(order.StrategyID)
			}
			
			// Handle the error
			shouldRetry, handledErr := e.errorHandler.HandleError(ctx, executionError)
			if !shouldRetry {
				// Add to dead letter queue if not retrying
				e.deadLetterQueue.AddOrder(order, handledErr)
			}
			
			return handledErr
		}
		
		return nil
	})
	
	return response, err
}

// RecoverFailedOrders attempts to recover failed orders from the dead letter queue
func (e *EnhancedOrderExecutionEngine) RecoverFailedOrders(ctx context.Context) (int, int, error) {
	return e.deadLetterQueue.RetryAllOrders(ctx, e.OrderExecutionEngine)
}

// GetFailedOrders returns all failed orders from the dead letter queue
func (e *EnhancedOrderExecutionEngine) GetFailedOrders() []Order {
	return e.deadLetterQueue.GetOrders()
}
