package orderexecution

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

// ExecutionAlgorithm defines the interface for execution algorithms
type ExecutionAlgorithm interface {
	Execute(ctx context.Context, engine *OrderExecutionEngine, request *OrderRequest) (*OrderResponse, error)
	Name() string
	Description() string
}

// POVAlgorithm implements Percentage of Volume algorithm
type POVAlgorithm struct {
	targetPercentage float64
	maxDuration      time.Duration
	volumeCheckInterval time.Duration
	marketDataFunc   func(symbol string) (int, error) // Function to get current market volume
}

// NewPOVAlgorithm creates a new POV algorithm
func NewPOVAlgorithm(targetPercentage float64, maxDuration time.Duration, marketDataFunc func(symbol string) (int, error)) *POVAlgorithm {
	return &POVAlgorithm{
		targetPercentage:    targetPercentage,
		maxDuration:         maxDuration,
		volumeCheckInterval: 30 * time.Second,
		marketDataFunc:      marketDataFunc,
	}
}

// Execute implements the ExecutionAlgorithm interface
func (a *POVAlgorithm) Execute(ctx context.Context, engine *OrderExecutionEngine, request *OrderRequest) (*OrderResponse, error) {
	// Create a new context with timeout
	ctx, cancel := context.WithTimeout(ctx, a.maxDuration)
	defer cancel()

	// Initialize variables
	totalQuantity := request.Quantity
	remainingQuantity := totalQuantity
	executedQuantity := 0
	totalValue := 0.0
	
	// Create a slice to hold responses
	var responses []*OrderResponse
	
	// Start execution loop
	startTime := time.Now()
	ticker := time.NewTicker(a.volumeCheckInterval)
	defer ticker.Stop()
	
	for remainingQuantity > 0 {
		select {
		case <-ctx.Done():
			// Context timeout or cancellation
			if executedQuantity == 0 {
				return nil, ctx.Err()
			}
			// Return partial execution
			break
		case <-ticker.C:
			// Get current market volume
			marketVolume, err := a.marketDataFunc(request.Symbol)
			if err != nil {
				log.Printf("Error getting market volume: %v", err)
				continue
			}
			
			// Calculate quantity to execute based on POV
			targetVolume := int(float64(marketVolume) * a.targetPercentage / 100.0)
			executeQuantity := min(targetVolume, remainingQuantity)
			
			if executeQuantity <= 0 {
				continue
			}
			
			// Create a new order request for this slice
			sliceRequest := &OrderRequest{
				Symbol:          request.Symbol,
				Quantity:        executeQuantity,
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
			
			// Add a tag to indicate this is part of a POV order
			sliceRequest.Tags = append(sliceRequest.Tags, fmt.Sprintf("pov_order_%.2f_percent", a.targetPercentage))
			
			// Execute the slice
			response, err := engine.ExecuteOrder(ctx, sliceRequest)
			if err != nil {
				log.Printf("Error executing POV slice: %v", err)
				continue
			}
			
			responses = append(responses, response)
			
			// Update remaining quantity
			if response.Status && response.Order != nil {
				executedQuantity += response.Order.FilledQuantity
				remainingQuantity -= response.Order.FilledQuantity
				totalValue += float64(response.Order.FilledQuantity) * response.Order.AveragePrice
			}
		}
		
		// Check if we've reached the max duration
		if time.Since(startTime) >= a.maxDuration {
			break
		}
	}
	
	// Combine the responses into a single response
	combinedResponse := &OrderResponse{
		Status: true,
		Order: &Order{
			ID:              responses[0].Order.ID, // Use the ID of the first order
			Symbol:          request.Symbol,
			Quantity:        totalQuantity,
			Price:           request.Price,
			OrderType:       request.OrderType,
			TransactionType: request.TransactionType,
			Status:          "EXECUTED", // Assume success for now
			FilledQuantity:  executedQuantity,
			AveragePrice:    0,
			PlacedAt:        startTime,
			UpdatedAt:       time.Now(),
			Validity:        request.Validity,
			TriggerPrice:    request.TriggerPrice,
			Exchange:        request.Exchange,
			Product:         request.Product,
			StrategyID:      request.StrategyID,
			Tags:            append([]string{}, request.Tags...),
		},
	}
	
	// Calculate average price
	if executedQuantity > 0 {
		combinedResponse.Order.AveragePrice = totalValue / float64(executedQuantity)
	}
	
	// Set appropriate status
	if executedQuantity == 0 {
		combinedResponse.Order.Status = "REJECTED"
		combinedResponse.Order.Message = "Failed to execute any quantity"
	} else if executedQuantity < totalQuantity {
		combinedResponse.Order.Status = "PARTIALLY_EXECUTED"
	}
	
	return combinedResponse, nil
}

// Name returns the algorithm name
func (a *POVAlgorithm) Name() string {
	return "POV"
}

// Description returns the algorithm description
func (a *POVAlgorithm) Description() string {
	return fmt.Sprintf("Percentage of Volume (%.2f%%) algorithm with max duration %v", a.targetPercentage, a.maxDuration)
}

// ImpactMinimizationAlgorithm implements an algorithm to minimize market impact
type ImpactMinimizationAlgorithm struct {
	maxDuration      time.Duration
	slices           int
	randomizeFactor  float64
	orderSplitter    *OrderSplitter
	marketDataFunc   func(symbol string) (float64, error) // Function to get current market volatility
}

// NewImpactMinimizationAlgorithm creates a new impact minimization algorithm
func NewImpactMinimizationAlgorithm(maxDuration time.Duration, slices int, randomizeFactor float64, marketDataFunc func(symbol string) (float64, error)) *ImpactMinimizationAlgorithm {
	return &ImpactMinimizationAlgorithm{
		maxDuration:     maxDuration,
		slices:          slices,
		randomizeFactor: randomizeFactor,
		orderSplitter:   NewOrderSplitter(0, 0), // Will be configured during execution
		marketDataFunc:  marketDataFunc,
	}
}

// Execute implements the ExecutionAlgorithm interface
func (a *ImpactMinimizationAlgorithm) Execute(ctx context.Context, engine *OrderExecutionEngine, request *OrderRequest) (*OrderResponse, error) {
	// Get market volatility
	volatility, err := a.marketDataFunc(request.Symbol)
	if err != nil {
		// Default to medium volatility if data is not available
		volatility = 0.5
	}
	
	// Adjust number of slices based on volatility and order size
	adjustedSlices := a.slices
	if volatility > 0.7 {
		// High volatility - increase slices
		adjustedSlices = int(float64(a.slices) * 1.5)
	} else if volatility < 0.3 {
		// Low volatility - decrease slices
		adjustedSlices = int(float64(a.slices) * 0.7)
	}
	
	// Ensure minimum number of slices
	if adjustedSlices < 2 {
		adjustedSlices = 2
	}
	
	// Calculate slice size and time between slices
	sliceSize := (request.Quantity + adjustedSlices - 1) / adjustedSlices
	timeBetweenSlices := a.maxDuration / time.Duration(adjustedSlices)
	
	// Configure the order splitter
	a.orderSplitter.maxOrderSize = sliceSize
	a.orderSplitter.timeBetweenChunks = timeBetweenSlices
	
	// Execute the split order
	responses, err := a.orderSplitter.ExecuteSplitOrder(ctx, engine, request)
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
	
	// Set appropriate status
	if totalFilled == 0 {
		combinedResponse.Order.Status = "REJECTED"
		combinedResponse.Order.Message = "Failed to execute any quantity"
	} else if totalFilled < request.Quantity {
		combinedResponse.Order.Status = "PARTIALLY_EXECUTED"
	}
	
	return combinedResponse, nil
}

// Name returns the algorithm name
func (a *ImpactMinimizationAlgorithm) Name() string {
	return "ImpactMinimization"
}

// Description returns the algorithm description
func (a *ImpactMinimizationAlgorithm) Description() string {
	return fmt.Sprintf("Impact Minimization algorithm with %d slices over %v", a.slices, a.maxDuration)
}

// AdaptiveAlgorithm implements an adaptive execution algorithm
type AdaptiveAlgorithm struct {
	maxDuration      time.Duration
	initialSlices    int
	adaptiveInterval time.Duration
	marketDataFunc   func(symbol string) (float64, float64, error) // Function to get price and volume
	orderSplitter    *OrderSplitter
}

// NewAdaptiveAlgorithm creates a new adaptive algorithm
func NewAdaptiveAlgorithm(maxDuration time.Duration, initialSlices int, adaptiveInterval time.Duration, marketDataFunc func(symbol string) (float64, float64, error)) *AdaptiveAlgorithm {
	return &AdaptiveAlgorithm{
		maxDuration:      maxDuration,
		initialSlices:    initialSlices,
		adaptiveInterval: adaptiveInterval,
		marketDataFunc:   marketDataFunc,
		orderSplitter:    NewOrderSplitter(0, 0), // Will be configured during execution
	}
}

// Execute implements the ExecutionAlgorithm interface
func (a *AdaptiveAlgorithm) Execute(ctx context.Context, engine *OrderExecutionEngine, request *OrderRequest) (*OrderResponse, error) {
	// Create a new context with timeout
	ctx, cancel := context.WithTimeout(ctx, a.maxDuration)
	defer cancel()
	
	// Initialize variables
	totalQuantity := request.Quantity
	remainingQuantity := totalQuantity
	executedQuantity := 0
	totalValue := 0.0
	
	// Get initial price and volume
	initialPrice, initialVolume, err := a.marketDataFunc(request.Symbol)
	if err != nil {
		return nil, fmt.Errorf("error getting initial market data: %w", err)
	}
	
	// Calculate initial slice size
	sliceSize := (request.Quantity + a.initialSlices - 1) / a.initialSlices
	
	// Create a slice to hold responses
	var responses []*OrderResponse
	
	// Start execution loop
	startTime := time.Now()
	ticker := time.NewTicker(a.adaptiveInterval)
	defer ticker.Stop()
	
	for remainingQuantity > 0 {
		select {
		case <-ctx.Done():
			// Context timeout or cancellation
			if executedQuantity == 0 {
				return nil, ctx.Err()
			}
			// Return partial execution
			break
		case <-ticker.C:
			// Get current price and volume
			currentPrice, currentVolume, err := a.marketDataFunc(request.Symbol)
			if err != nil {
				log.Printf("Error getting market data: %v", err)
				continue
			}
			
			// Calculate price movement
			priceChange := (currentPrice - initialPrice) / initialPrice
			
			// Calculate volume change
			volumeChange := (currentVolume - initialVolume) / initialVolume
			
			// Adjust slice size based on market conditions
			adjustedSliceSize := sliceSize
			
			// If price is moving in our favor, increase slice size
			if (request.TransactionType == Buy && priceChange < 0) || (request.TransactionType == Sell && priceChange > 0) {
				adjustedSliceSize = int(float64(sliceSize) * (1.0 + math.Abs(priceChange)))
			} else {
				// If price is moving against us, decrease slice size
				adjustedSliceSize = int(float64(sliceSize) * (1.0 - math.Abs(priceChange)*0.5))
			}
			
			// Adjust based on volume
			if volumeChange > 0.2 {
				// Volume increasing, we can be more aggressive
				adjustedSliceSize = int(float64(adjustedSliceSize) * (1.0 + volumeChange*0.5))
			} else if volumeChange < -0.2 {
				// Volume decreasing, be more cautious
				adjustedSliceSize = int(float64(adjustedSliceSize) * (1.0 + volumeChange*0.5))
			}
			
			// Ensure slice size is at least 1 and not more than remaining quantity
			adjustedSliceSize = max(1, min(adjustedSliceSize, remainingQuantity))
			
			// Create a new order request for this slice
			sliceRequest := &OrderRequest{
				Symbol:          request.Symbol,
				Quantity:        adjustedSliceSize,
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
			
			// Add a tag to indicate this is part of an adaptive order
			sliceRequest.Tags = append(sliceRequest.Tags, "adaptive_order")
			
			// Execute the slice
			response, err := engine.ExecuteOrder(ctx, sliceRequest)
			if err != nil {
				log.Printf("Error executing adaptive slice: %v", err)
				continue
			}
			
			responses = append(responses, response)
			
			// Update remaining quantity
			if response.Status && response.Order != nil {
				executedQuantity += response.Order.FilledQuantity
				remainingQuantity -= response.Order.FilledQuantity
				totalValue += float64(response.Order.FilledQuantity) * response.Order.AveragePrice
			}
			
			// Update initial values for next comparison
			initialPrice = currentPrice
			initialVolume = currentVolume
		}
		
		// Check if we've reached the max duration
		if time.Since(startTime) >= a.maxDuration {
			break
		}
	}
	
	// Combine the responses into a single response
	combinedResponse := &OrderResponse{
		Status: true,
		Order: &Order{
			ID:              responses[0].Order.ID, // Use the ID of the first order
			Symbol:          request.Symbol,
			Quantity:        totalQuantity,
			Price:           request.Price,
			OrderType:       request.OrderType,
			TransactionType: request.TransactionType,
			Status:          "EXECUTED", // Assume success for now
			FilledQuantity:  executedQuantity,
			AveragePrice:    0,
			PlacedAt:        startTime,
			UpdatedAt:       time.Now(),
			Validity:        request.Validity,
			TriggerPrice:    request.TriggerPrice,
			Exchange:        request.Exchange,
			Product:         request.Product,
			StrategyID:      request.StrategyID,
			Tags:            append([]string{}, request.Tags...),
		},
	}
	
	// Calculate average price
	if executedQuantity > 0 {
		combinedResponse.Order.AveragePrice = totalValue / float64(executedQuantity)
	}
	
	// Set appropriate status
	if executedQuantity == 0 {
		combinedResponse.Order.Status = "REJECTED"
		combinedResponse.Order.Message = "Failed to execute any quantity"
	} else if executedQuantity < totalQuantity {
		combinedResponse.Order.Status = "PARTIALLY_EXECUTED"
	}
	
	return combinedResponse, nil
}

// Name returns the algorithm name
func (a *AdaptiveAlgorithm) Name() string {
	return "Adaptive"
}

// Description returns the algorithm description
func (a *AdaptiveAlgorithm) Description() string {
	return fmt.Sprintf("Adaptive algorithm with initial %d slices over %v", a.initialSlices, a.maxDuration)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// AlgorithmManager manages execution algorithms
type AlgorithmManager struct {
	algorithms map[string]ExecutionAlgorithm
	mutex      sync.RWMutex
}

// NewAlgorithmManager creates a new algorithm manager
func NewAlgorithmManager() *AlgorithmManager {
	return &AlgorithmManager{
		algorithms: make(map[string]ExecutionAlgorithm),
	}
}

// RegisterAlgorithm registers an algorithm with the manager
func (m *AlgorithmManager) RegisterAlgorithm(algorithm ExecutionAlgorithm) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.algorithms[algorithm.Name()] = algorithm
}

// GetAlgorithm returns an algorithm by name
func (m *AlgorithmManager) GetAlgorithm(name string) (ExecutionAlgorithm, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	algorithm, exists := m.algorithms[name]
	if !exists {
		return nil, fmt.Errorf("algorithm not found: %s", name)
	}
	return algorithm, nil
}

// ExecuteWithAlgorithm executes an order using the specified algorithm
func (m *AlgorithmManager) ExecuteWithAlgorithm(ctx context.Context, engine *OrderExecutionEngine, algorithmName string, request *OrderRequest) (*OrderResponse, error) {
	algorithm, err := m.GetAlgorithm(algorithmName)
	if err != nil {
		return nil, err
	}
	return algorithm.Execute(ctx, engine, request)
}

// ListAlgorithms returns a list of available algorithm names
func (m *AlgorithmManager) ListAlgorithms() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	names := make([]string, 0, len(m.algorithms))
	for name := range m.algorithms {
		names = append(names, name)
	}
	return names
}

// GetAlgorithmDescription returns the description of an algorithm
func (m *AlgorithmManager) GetAlgorithmDescription(name string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	algorithm, exists := m.algorithms[name]
	if !exists {
		return "", fmt.Errorf("algorithm not found: %s", name)
	}
	return algorithm.Description(), nil
}
