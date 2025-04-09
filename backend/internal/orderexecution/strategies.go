package orderexecution

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// OrderSplitter is responsible for splitting large orders into smaller chunks
type OrderSplitter struct {
	maxOrderSize      int
	timeBetweenChunks time.Duration
}

// NewOrderSplitter creates a new order splitter
func NewOrderSplitter(maxOrderSize int, timeBetweenChunks time.Duration) *OrderSplitter {
	return &OrderSplitter{
		maxOrderSize:      maxOrderSize,
		timeBetweenChunks: timeBetweenChunks,
	}
}

// SplitOrder splits a large order into smaller chunks
func (s *OrderSplitter) SplitOrder(request *OrderRequest) []*OrderRequest {
	if request.Quantity <= s.maxOrderSize {
		return []*OrderRequest{request}
	}

	// Calculate number of chunks needed
	numChunks := (request.Quantity + s.maxOrderSize - 1) / s.maxOrderSize
	chunks := make([]*OrderRequest, numChunks)

	remaining := request.Quantity
	for i := 0; i < numChunks; i++ {
		chunkSize := s.maxOrderSize
		if remaining < chunkSize {
			chunkSize = remaining
		}

		// Create a copy of the request with adjusted quantity
		chunks[i] = &OrderRequest{
			Symbol:          request.Symbol,
			Quantity:        chunkSize,
			Price:           request.Price,
			OrderType:       request.OrderType,
			TransactionType: request.TransactionType,
			Validity:        request.Validity,
			TriggerPrice:    request.TriggerPrice,
			Exchange:        request.Exchange,
			Product:         request.Product,
			StrategyID:      request.StrategyID,
			Tags:            append([]string{}, request.Tags...),
		}

		// Add a tag to indicate this is part of a split order
		chunks[i].Tags = append(chunks[i].Tags, fmt.Sprintf("split_order_part_%d_of_%d", i+1, numChunks))

		remaining -= chunkSize
	}

	return chunks
}

// ExecuteSplitOrder executes a split order with time delays between chunks
func (s *OrderSplitter) ExecuteSplitOrder(ctx context.Context, engine *OrderExecutionEngine, request *OrderRequest) ([]*OrderResponse, error) {
	chunks := s.SplitOrder(request)
	responses := make([]*OrderResponse, len(chunks))

	for i, chunk := range chunks {
		// Execute this chunk
		response, err := engine.ExecuteOrder(ctx, chunk)
		if err != nil {
			return responses[:i], fmt.Errorf("error executing chunk %d: %w", i+1, err)
		}

		responses[i] = response

		// Wait between chunks, but not after the last one
		if i < len(chunks)-1 && s.timeBetweenChunks > 0 {
			select {
			case <-ctx.Done():
				return responses[:i+1], ctx.Err()
			case <-time.After(s.timeBetweenChunks):
				// Continue to next chunk
			}
		}
	}

	return responses, nil
}

// OrderExecutionStrategy defines the interface for order execution strategies
type OrderExecutionStrategy interface {
	Execute(ctx context.Context, engine *OrderExecutionEngine, request *OrderRequest) (*OrderResponse, error)
	Name() string
}

// ImmediateOrCancelStrategy executes an order with IOC validity
type ImmediateOrCancelStrategy struct{}

// Execute implements the OrderExecutionStrategy interface
func (s *ImmediateOrCancelStrategy) Execute(ctx context.Context, engine *OrderExecutionEngine, request *OrderRequest) (*OrderResponse, error) {
	// Set the validity to IOC
	request.Validity = "IOC"
	return engine.ExecuteOrder(ctx, request)
}

// Name returns the strategy name
func (s *ImmediateOrCancelStrategy) Name() string {
	return "ImmediateOrCancel"
}

// TwapStrategy implements Time-Weighted Average Price strategy
type TwapStrategy struct {
	duration      time.Duration
	numSlices     int
	orderSplitter *OrderSplitter
}

// NewTwapStrategy creates a new TWAP strategy
func NewTwapStrategy(duration time.Duration, numSlices int) *TwapStrategy {
	return &TwapStrategy{
		duration:      duration,
		numSlices:     numSlices,
		orderSplitter: NewOrderSplitter(0, 0), // Will be configured during execution
	}
}

// Execute implements the OrderExecutionStrategy interface
func (s *TwapStrategy) Execute(ctx context.Context, engine *OrderExecutionEngine, request *OrderRequest) (*OrderResponse, error) {
	// Calculate slice size and time between slices
	sliceSize := (request.Quantity + s.numSlices - 1) / s.numSlices
	timeBetweenSlices := s.duration / time.Duration(s.numSlices)

	// Configure the order splitter
	s.orderSplitter.maxOrderSize = sliceSize
	s.orderSplitter.timeBetweenChunks = timeBetweenSlices

	// Execute the split order
	responses, err := s.orderSplitter.ExecuteSplitOrder(ctx, engine, request)
	if err != nil {
		return nil, err
	}

	// Combine the responses into a single response
	combinedResponse := &OrderResponse{
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

	combinedResponse.Order.FilledQuantity = totalFilled
	if totalFilled > 0 {
		combinedResponse.Order.AveragePrice = totalValue / float64(totalFilled)
	}

	return combinedResponse, nil
}

// Name returns the strategy name
func (s *TwapStrategy) Name() string {
	return "TWAP"
}

// VwapStrategy implements Volume-Weighted Average Price strategy
type VwapStrategy struct {
	duration       time.Duration
	volumeProfile  []float64 // Percentage of volume to execute in each time slice
	orderSplitter  *OrderSplitter
	marketDataFunc func(symbol string) ([]float64, error) // Function to get volume profile
}

// NewVwapStrategy creates a new VWAP strategy
func NewVwapStrategy(duration time.Duration, marketDataFunc func(symbol string) ([]float64, error)) *VwapStrategy {
	return &VwapStrategy{
		duration:       duration,
		volumeProfile:  nil, // Will be populated during execution
		orderSplitter:  NewOrderSplitter(0, 0), // Will be configured during execution
		marketDataFunc: marketDataFunc,
	}
}

// Execute implements the OrderExecutionStrategy interface
func (s *VwapStrategy) Execute(ctx context.Context, engine *OrderExecutionEngine, request *OrderRequest) (*OrderResponse, error) {
	// Get volume profile for the symbol
	volumeProfile, err := s.marketDataFunc(request.Symbol)
	if err != nil {
		// Fall back to equal distribution if volume profile is not available
		numSlices := 5 // Default to 5 slices
		volumeProfile = make([]float64, numSlices)
		for i := range volumeProfile {
			volumeProfile[i] = 1.0 / float64(numSlices)
		}
	}
	s.volumeProfile = volumeProfile

	// Calculate time between slices
	numSlices := len(s.volumeProfile)
	timeBetweenSlices := s.duration / time.Duration(numSlices)

	// Create order slices based on volume profile
	slices := make([]*OrderRequest, numSlices)
	for i := 0; i < numSlices; i++ {
		sliceQuantity := int(float64(request.Quantity) * s.volumeProfile[i])
		if sliceQuantity <= 0 {
			sliceQuantity = 1 // Ensure at least 1 quantity per slice
		}

		// Create a copy of the request with adjusted quantity
		slices[i] = &OrderRequest{
			Symbol:          request.Symbol,
			Quantity:        sliceQuantity,
			Price:           request.Price,
			OrderType:       request.OrderType,
			TransactionType: request.TransactionType,
			Validity:        request.Validity,
			TriggerPrice:    request.TriggerPrice,
			Exchange:        request.Exchange,
			Product:         request.Product,
			StrategyID:      request.StrategyID,
			Tags:            append([]string{}, request.Tags...),
		}

		// Add a tag to indicate this is part of a VWAP order
		slices[i].Tags = append(slices[i].Tags, fmt.Sprintf("vwap_slice_%d_of_%d", i+1, numSlices))
	}

	// Execute the slices with time delays
	responses := make([]*OrderResponse, numSlices)
	for i, slice := range slices {
		// Execute this slice
		response, err := engine.ExecuteOrder(ctx, slice)
		if err != nil {
			return nil, fmt.Errorf("error executing VWAP slice %d: %w", i+1, err)
		}

		responses[i] = response

		// Wait between slices, but not after the last one
		if i < numSlices-1 && timeBetweenSlices > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(timeBetweenSlices):
				// Continue to next slice
			}
		}
	}

	// Combine the responses into a single response
	combinedResponse := &OrderResponse{
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

	combinedResponse.Order.FilledQuantity = totalFilled
	if totalFilled > 0 {
		combinedResponse.Order.AveragePrice = totalValue / float64(totalFilled)
	}

	return combinedResponse, nil
}

// Name returns the strategy name
func (s *VwapStrategy) Name() string {
	return "VWAP"
}

// StrategyManager manages order execution strategies
type StrategyManager struct {
	strategies map[string]OrderExecutionStrategy
	mutex      sync.RWMutex
}

// NewStrategyManager creates a new strategy manager
func NewStrategyManager() *StrategyManager {
	return &StrategyManager{
		strategies: make(map[string]OrderExecutionStrategy),
	}
}

// RegisterStrategy registers a strategy with the manager
func (m *StrategyManager) RegisterStrategy(strategy OrderExecutionStrategy) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.strategies[strategy.Name()] = strategy
}

// GetStrategy returns a strategy by name
func (m *StrategyManager) GetStrategy(name string) (OrderExecutionStrategy, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	strategy, exists := m.strategies[name]
	if !exists {
		return nil, fmt.Errorf("strategy not found: %s", name)
	}
	return strategy, nil
}

// ExecuteWithStrategy executes an order using the specified strategy
func (m *StrategyManager) ExecuteWithStrategy(ctx context.Context, engine *OrderExecutionEngine, strategyName string, request *OrderRequest) (*OrderResponse, error) {
	strategy, err := m.GetStrategy(strategyName)
	if err != nil {
		return nil, err
	}
	return strategy.Execute(ctx, engine, request)
}

// ListStrategies returns a list of available strategy names
func (m *StrategyManager) ListStrategies() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	names := make([]string, 0, len(m.strategies))
	for name := range m.strategies {
		names = append(names, name)
	}
	return names
}
