package orderexecution

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
	"github.com/google/uuid"
)

// OrderPlacementService handles the placement of orders
type OrderPlacementService struct {
	engine           *OrderExecutionEngine
	strategyManager  *StrategyManager
	orderSplitter    *OrderSplitter
	rateLimiter      *RateLimiter
	orderQueue       chan *OrderRequest
	workers          int
	isRunning        bool
	stopChan         chan struct{}
	wg               sync.WaitGroup
	mutex            sync.RWMutex
	orderCallbacks   map[string]OrderCallback
	callbacksMutex   sync.RWMutex
}

// OrderCallback is a function that gets called when an order is processed
type OrderCallback func(response *OrderResponse, err error)

// RateLimiter controls the rate of order placement
type RateLimiter struct {
	maxOrdersPerSecond int
	tokens             int
	lastRefill         time.Time
	mutex              sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxOrdersPerSecond int) *RateLimiter {
	return &RateLimiter{
		maxOrdersPerSecond: maxOrdersPerSecond,
		tokens:             maxOrdersPerSecond,
		lastRefill:         time.Now(),
	}
}

// TakeToken attempts to take a token from the rate limiter
func (r *RateLimiter) TakeToken() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(r.lastRefill)
	
	// Refill tokens based on elapsed time
	tokensToAdd := int(elapsed.Seconds() * float64(r.maxOrdersPerSecond))
	if tokensToAdd > 0 {
		r.tokens = min(r.tokens+tokensToAdd, r.maxOrdersPerSecond)
		r.lastRefill = now
	}
	
	// Check if we have tokens available
	if r.tokens > 0 {
		r.tokens--
		return true
	}
	
	return false
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// NewOrderPlacementService creates a new order placement service
func NewOrderPlacementService(engine *OrderExecutionEngine, strategyManager *StrategyManager, queueSize int, workers int, maxOrdersPerSecond int) *OrderPlacementService {
	return &OrderPlacementService{
		engine:          engine,
		strategyManager: strategyManager,
		orderSplitter:   NewOrderSplitter(1000, 100*time.Millisecond), // Default to 1000 max order size with 100ms between chunks
		rateLimiter:     NewRateLimiter(maxOrdersPerSecond),
		orderQueue:      make(chan *OrderRequest, queueSize),
		workers:         workers,
		stopChan:        make(chan struct{}),
		orderCallbacks:  make(map[string]OrderCallback),
	}
}

// Start starts the order placement service
func (s *OrderPlacementService) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if s.isRunning {
		return
	}
	
	s.isRunning = true
	s.stopChan = make(chan struct{})
	
	// Start worker goroutines
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
	
	log.Printf("Order placement service started with %d workers", s.workers)
}

// Stop stops the order placement service
func (s *OrderPlacementService) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if !s.isRunning {
		return
	}
	
	close(s.stopChan)
	s.wg.Wait()
	s.isRunning = false
	
	log.Println("Order placement service stopped")
}

// worker processes orders from the queue
func (s *OrderPlacementService) worker(id int) {
	defer s.wg.Done()
	
	log.Printf("Order placement worker %d started", id)
	
	for {
		select {
		case <-s.stopChan:
			log.Printf("Order placement worker %d stopping", id)
			return
		case request := <-s.orderQueue:
			// Wait for rate limiter
			for !s.rateLimiter.TakeToken() {
				select {
				case <-s.stopChan:
					return
				case <-time.After(10 * time.Millisecond):
					// Try again after a short delay
				}
			}
			
			// Process the order
			s.processOrder(request)
		}
	}
}

// processOrder processes a single order
func (s *OrderPlacementService) processOrder(request *OrderRequest) {
	// Generate a unique ID for this request if not already set
	requestID := fmt.Sprintf("%s-%s", request.Symbol, uuid.New().String())
	
	// Check if we have a callback registered for this request
	s.callbacksMutex.RLock()
	callback, hasCallback := s.orderCallbacks[requestID]
	s.callbacksMutex.RUnlock()
	
	// Execute the order
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	var response *OrderResponse
	var err error
	
	// Check if we need to use a specific strategy
	if request.Tags != nil {
		for _, tag := range request.Tags {
			if tag == "TWAP" || tag == "VWAP" || tag == "ImmediateOrCancel" {
				// Execute with strategy
				response, err = s.strategyManager.ExecuteWithStrategy(ctx, s.engine, tag, request)
				goto done
			}
		}
	}
	
	// Check if we need to split the order
	if request.Quantity > s.orderSplitter.maxOrderSize {
		// Split and execute the order
		responses, err := s.orderSplitter.ExecuteSplitOrder(ctx, s.engine, request)
		if err != nil {
			if hasCallback {
				callback(nil, err)
			}
			log.Printf("Error executing split order: %v", err)
			return
		}
		
		// Combine the responses
		response = &OrderResponse{
			Status: true,
			Order: &Order{
				ID:              responses[0].Order.ID, // Use the ID of the first order
				Symbol:          request.Symbol,
				Quantity:        request.Quantity,
				Price:           request.Price,
				OrderType:       request.OrderType,
				TransactionType: request.TransactionType,
				Status:          "EXECUTED", // Assume success for now
				FilledQuantity:  0,
				AveragePrice:    0,
				PlacedAt:        time.Now(),
				UpdatedAt:       time.Now(),
				Validity:        request.Validity,
				TriggerPrice:    request.TriggerPrice,
				Exchange:        request.Exchange,
				Product:         request.Product,
				StrategyID:      request.StrategyID,
				Tags:            append([]string{}, request.Tags...),
			},
		}
		
		// Calculate filled quantity and average price
		totalFilled := 0
		totalValue := 0.0
		for _, resp := range responses {
			if resp != nil && resp.Order != nil {
				totalFilled += resp.Order.FilledQuantity
				totalValue += float64(resp.Order.FilledQuantity) * resp.Order.AveragePrice
			}
		}
		
		response.Order.FilledQuantity = totalFilled
		if totalFilled > 0 {
			response.Order.AveragePrice = totalValue / float64(totalFilled)
		}
	} else {
		// Execute the order directly
		response, err = s.engine.ExecuteOrder(ctx, request)
	}
	
done:
	// Call the callback if registered
	if hasCallback {
		callback(response, err)
		
		// Remove the callback
		s.callbacksMutex.Lock()
		delete(s.orderCallbacks, requestID)
		s.callbacksMutex.Unlock()
	}
	
	// Log the result
	if err != nil {
		log.Printf("Error executing order: %v", err)
	} else {
		log.Printf("Order executed: %s, Status: %s", response.Order.ID, response.Order.Status)
	}
}

// PlaceOrder places an order asynchronously
func (s *OrderPlacementService) PlaceOrder(request *OrderRequest, callback OrderCallback) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	if !s.isRunning {
		return errors.New("order placement service is not running")
	}
	
	// Generate a unique ID for this request
	requestID := fmt.Sprintf("%s-%s", request.Symbol, uuid.New().String())
	
	// Register the callback if provided
	if callback != nil {
		s.callbacksMutex.Lock()
		s.orderCallbacks[requestID] = callback
		s.callbacksMutex.Unlock()
	}
	
	// Add the request to the queue
	select {
	case s.orderQueue <- request:
		return nil
	default:
		// Queue is full
		if callback != nil {
			// Remove the callback
			s.callbacksMutex.Lock()
			delete(s.orderCallbacks, requestID)
			s.callbacksMutex.Unlock()
			
			// Call the callback with an error
			callback(nil, errors.New("order queue is full"))
		}
		return errors.New("order queue is full")
	}
}

// PlaceOrderSync places an order synchronously
func (s *OrderPlacementService) PlaceOrderSync(ctx context.Context, request *OrderRequest) (*OrderResponse, error) {
	s.mutex.RLock()
	if !s.isRunning {
		s.mutex.RUnlock()
		return nil, errors.New("order placement service is not running")
	}
	s.mutex.RUnlock()
	
	// Wait for rate limiter
	for !s.rateLimiter.TakeToken() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(10 * time.Millisecond):
			// Try again after a short delay
		}
	}
	
	// Check if we need to use a specific strategy
	if request.Tags != nil {
		for _, tag := range request.Tags {
			if tag == "TWAP" || tag == "VWAP" || tag == "ImmediateOrCancel" {
				// Execute with strategy
				return s.strategyManager.ExecuteWithStrategy(ctx, s.engine, tag, request)
			}
		}
	}
	
	// Check if we need to split the order
	if request.Quantity > s.orderSplitter.maxOrderSize {
		// Split and execute the order
		return s.executeSplitOrderSync(ctx, request)
	}
	
	// Execute the order directly
	return s.engine.ExecuteOrder(ctx, request)
}

// executeSplitOrderSync executes a split order synchronously
func (s *OrderPlacementService) executeSplitOrderSync(ctx context.Context, request *OrderRequest) (*OrderResponse, error) {
	responses, err := s.orderSplitter.ExecuteSplitOrder(ctx, s.engine, request)
	if err != nil {
		return nil, err
	}
	
	// Combine the responses
	response := &OrderResponse{
		Status: true,
		Order: &Order{
			ID:              responses[0].Order.ID, // Use the ID of the first order
			Symbol:          request.Symbol,
			Quantity:        request.Quantity,
			Price:           request.Price,
			OrderType:       request.OrderType,
			TransactionType: request.TransactionType,
			Status:          "EXECUTED", // Assume success for now
			FilledQuantity:  0,
			AveragePrice:    0,
			PlacedAt:        time.Now(),
			UpdatedAt:       time.Now(),
			Validity:        request.Validity,
			TriggerPrice:    request.TriggerPrice,
			Exchange:        request.Exchange,
			Product:         request.Product,
			StrategyID:      request.StrategyID,
			Tags:            append([]string{}, request.Tags...),
		},
	}
	
	// Calculate filled quantity and average price
	totalFilled := 0
	totalValue := 0.0
	for _, resp := range responses {
		if resp != nil && resp.Order != nil {
			totalFilled += resp.Order.FilledQuantity
			totalValue += float64(resp.Order.FilledQuantity) * resp.Order.AveragePrice
		}
	}
	
	response.Order.FilledQuantity = totalFilled
	if totalFilled > 0 {
		response.Order.AveragePrice = totalValue / float64(totalFilled)
	}
	
	return response, nil
}

// PlaceOrderBatch places multiple orders asynchronously
func (s *OrderPlacementService) PlaceOrderBatch(requests []*OrderRequest, callback func([]*OrderResponse, error)) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	if !s.isRunning {
		return errors.New("order placement service is not running")
	}
	
	// Create a wait group to track completion of all orders
	var wg sync.WaitGroup
	wg.Add(len(requests))
	
	// Create a slice to hold responses
	responses := make([]*OrderResponse, len(requests))
	errors := make([]error, len(requests))
	
	// Place each order
	for i, request := range requests {
		go func(idx int, req *OrderRequest) {
			defer wg.Done()
			
			// Place the order
			err := s.PlaceOrder(req, func(response *OrderResponse, err error) {
				responses[idx] = response
				errors[idx] = err
			})
			
			if err != nil {
				errors[idx] = err
			}
		}(i, request)
	}
	
	// Wait for all orders to complete in a separate goroutine
	go func() {
		wg.Wait()
		
		// Check if any errors occurred
		var batchErr error
		for _, err := range errors {
			if err != nil {
				batchErr = errors.New("one or more orders failed")
				break
			}
		}
		
		// Call the callback
		if callback != nil {
			callback(responses, batchErr)
		}
	}()
	
	return nil
}

// PlaceOrderBatchSync places multiple orders synchronously
func (s *OrderPlacementService) PlaceOrderBatchSync(ctx context.Context, requests []*OrderRequest) ([]*OrderResponse, error) {
	s.mutex.RLock()
	if !s.isRunning {
		s.mutex.RUnlock()
		return nil, errors.New("order placement service is not running")
	}
	s.mutex.RUnlock()
	
	// Create a slice to hold responses
	responses := make([]*OrderResponse, len(requests))
	
	// Create a channel to limit concurrent order executions
	semaphore := make(chan struct{}, s.workers)
	
	// Create a wait group to track completion of all orders
	var wg sync.WaitGroup
	wg.Add(len(requests))
	
	// Create a mutex to protect access to the error variable
	var errMutex sync.Mutex
	var batchErr error
	
	// Place each order
	for i, request := range requests {
		semaphore <- struct{}{} // Acquire semaphore
		
		go func(idx int, req *OrderRequest) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore
			
			// Place the order
			response, err := s.PlaceOrderSync(ctx, req)
			
			if err != nil {
				errMutex.Lock()
				batchErr = errors.New("one or more orders failed")
				errMutex.Unlock()
				responses[idx] = &OrderResponse{
					Status: false,
					Error:  err.Error(),
				}
			} else {
				responses[idx] = response
			}
		}(i, request)
	}
	
	// Wait for all orders to complete
	wg.Wait()
	
	return responses, batchErr
}

// SetOrderSplitterConfig configures the order splitter
func (s *OrderPlacementService) SetOrderSplitterConfig(maxOrderSize int, timeBetweenChunks time.Duration) {
	s.orderSplitter.maxOrderSize = maxOrderSize
	s.orderSplitter.timeBetweenChunks = timeBetweenChunks
}

// SetRateLimiterConfig configures the rate limiter
func (s *OrderPlacementService) SetRateLimiterConfig(maxOrdersPerSecond int) {
	s.rateLimiter.mutex.Lock()
	defer s.rateLimiter.mutex.Unlock()
	
	s.rateLimiter.maxOrdersPerSecond = maxOrdersPerSecond
	s.rateLimiter.tokens = maxOrdersPerSecond
	s.rateLimiter.lastRefill = time.Now()
}

// GetQueueStats returns statistics about the order queue
func (s *OrderPlacementService) GetQueueStats() (int, int) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	return len(s.orderQueue), cap(s.orderQueue)
}
