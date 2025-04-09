package orderexecution

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

// SmartRoutingStrategy defines the strategy for smart routing
type SmartRoutingStrategy string

const (
	// BestPrice routes orders to the broker with the best price
	BestPrice SmartRoutingStrategy = "BEST_PRICE"
	
	// LowestLatency routes orders to the broker with the lowest latency
	LowestLatency SmartRoutingStrategy = "LOWEST_LATENCY"
	
	// HighestFillRate routes orders to the broker with the highest fill rate
	HighestFillRate SmartRoutingStrategy = "HIGHEST_FILL_RATE"
	
	// LowestCost routes orders to the broker with the lowest transaction cost
	LowestCost SmartRoutingStrategy = "LOWEST_COST"
	
	// RoundRobin distributes orders evenly among brokers
	RoundRobin SmartRoutingStrategy = "ROUND_ROBIN"
	
	// VolumeWeighted routes orders based on volume available at each broker
	VolumeWeighted SmartRoutingStrategy = "VOLUME_WEIGHTED"
)

// BrokerMetrics stores performance metrics for a broker
type BrokerMetrics struct {
	AverageLatency    time.Duration
	FillRate          float64
	TransactionCost   float64
	AvailableVolume   map[string]int // Symbol -> Volume
	LastOrderTime     time.Time
	OrderCount        int
	SuccessfulOrders  int
	FailedOrders      int
	TotalOrdersPlaced int
}

// DefaultSmartRouter implements the SmartRouter interface
type DefaultSmartRouter struct {
	brokers           map[string]BrokerAdapter
	brokerMetrics     map[string]*BrokerMetrics
	defaultStrategy   SmartRoutingStrategy
	symbolStrategies  map[string]SmartRoutingStrategy
	roundRobinCounter int
	mutex             sync.RWMutex
}

// NewDefaultSmartRouter creates a new default smart router
func NewDefaultSmartRouter(defaultStrategy SmartRoutingStrategy) *DefaultSmartRouter {
	return &DefaultSmartRouter{
		brokers:          make(map[string]BrokerAdapter),
		brokerMetrics:    make(map[string]*BrokerMetrics),
		defaultStrategy:  defaultStrategy,
		symbolStrategies: make(map[string]SmartRoutingStrategy),
	}
}

// RegisterBroker registers a broker with the smart router
func (r *DefaultSmartRouter) RegisterBroker(name string, broker BrokerAdapter) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.brokers[name] = broker
	r.brokerMetrics[name] = &BrokerMetrics{
		AverageLatency:  0,
		FillRate:        1.0, // Start with optimistic fill rate
		TransactionCost: 0.0,
		AvailableVolume: make(map[string]int),
		OrderCount:      0,
	}
}

// SetStrategyForSymbol sets a specific routing strategy for a symbol
func (r *DefaultSmartRouter) SetStrategyForSymbol(symbol string, strategy SmartRoutingStrategy) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.symbolStrategies[symbol] = strategy
}

// UpdateBrokerMetrics updates the metrics for a broker
func (r *DefaultSmartRouter) UpdateBrokerMetrics(brokerName string, latency time.Duration, success bool, cost float64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	metrics, exists := r.brokerMetrics[brokerName]
	if !exists {
		return
	}
	
	// Update latency (exponential moving average)
	if metrics.OrderCount == 0 {
		metrics.AverageLatency = latency
	} else {
		alpha := 0.2 // Weight for new value
		metrics.AverageLatency = time.Duration(float64(metrics.AverageLatency)*(1-alpha) + float64(latency)*alpha)
	}
	
	// Update fill rate
	metrics.TotalOrdersPlaced++
	if success {
		metrics.SuccessfulOrders++
	} else {
		metrics.FailedOrders++
	}
	metrics.FillRate = float64(metrics.SuccessfulOrders) / float64(metrics.TotalOrdersPlaced)
	
	// Update transaction cost (exponential moving average)
	if metrics.OrderCount == 0 {
		metrics.TransactionCost = cost
	} else {
		alpha := 0.2 // Weight for new value
		metrics.TransactionCost = metrics.TransactionCost*(1-alpha) + cost*alpha
	}
	
	metrics.LastOrderTime = time.Now()
	metrics.OrderCount++
}

// UpdateAvailableVolume updates the available volume for a symbol at a broker
func (r *DefaultSmartRouter) UpdateAvailableVolume(brokerName string, symbol string, volume int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	metrics, exists := r.brokerMetrics[brokerName]
	if !exists {
		return
	}
	
	metrics.AvailableVolume[symbol] = volume
}

// RouteOrder routes an order to the appropriate broker based on the strategy
func (r *DefaultSmartRouter) RouteOrder(ctx context.Context, request *OrderRequest) (BrokerAdapter, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if len(r.brokers) == 0 {
		return nil, errors.New("no brokers registered")
	}
	
	// Determine the strategy to use
	strategy := r.defaultStrategy
	if symbolStrategy, exists := r.symbolStrategies[request.Symbol]; exists {
		strategy = symbolStrategy
	}
	
	// Route based on the strategy
	switch strategy {
	case BestPrice:
		return r.routeByBestPrice(ctx, request)
	case LowestLatency:
		return r.routeByLowestLatency(request)
	case HighestFillRate:
		return r.routeByHighestFillRate(request)
	case LowestCost:
		return r.routeByLowestCost(request)
	case RoundRobin:
		return r.routeByRoundRobin()
	case VolumeWeighted:
		return r.routeByVolumeWeighted(request)
	default:
		// Default to round-robin if strategy is unknown
		return r.routeByRoundRobin()
	}
}

// routeByBestPrice routes to the broker with the best price
func (r *DefaultSmartRouter) routeByBestPrice(ctx context.Context, request *OrderRequest) (BrokerAdapter, error) {
	// In a real implementation, this would query each broker for current prices
	// For now, we'll just return the first broker as a placeholder
	for _, broker := range r.brokers {
		return broker, nil
	}
	
	return nil, errors.New("no suitable broker found")
}

// routeByLowestLatency routes to the broker with the lowest latency
func (r *DefaultSmartRouter) routeByLowestLatency(request *OrderRequest) (BrokerAdapter, error) {
	var lowestLatencyBroker BrokerAdapter
	var lowestLatency time.Duration = math.MaxInt64
	
	for name, metrics := range r.brokerMetrics {
		if metrics.AverageLatency < lowestLatency {
			lowestLatency = metrics.AverageLatency
			lowestLatencyBroker = r.brokers[name]
		}
	}
	
	if lowestLatencyBroker == nil {
		return nil, errors.New("no suitable broker found")
	}
	
	return lowestLatencyBroker, nil
}

// routeByHighestFillRate routes to the broker with the highest fill rate
func (r *DefaultSmartRouter) routeByHighestFillRate(request *OrderRequest) (BrokerAdapter, error) {
	var highestFillRateBroker BrokerAdapter
	var highestFillRate float64 = -1
	
	for name, metrics := range r.brokerMetrics {
		if metrics.FillRate > highestFillRate {
			highestFillRate = metrics.FillRate
			highestFillRateBroker = r.brokers[name]
		}
	}
	
	if highestFillRateBroker == nil {
		return nil, errors.New("no suitable broker found")
	}
	
	return highestFillRateBroker, nil
}

// routeByLowestCost routes to the broker with the lowest transaction cost
func (r *DefaultSmartRouter) routeByLowestCost(request *OrderRequest) (BrokerAdapter, error) {
	var lowestCostBroker BrokerAdapter
	var lowestCost float64 = math.MaxFloat64
	
	for name, metrics := range r.brokerMetrics {
		if metrics.TransactionCost < lowestCost {
			lowestCost = metrics.TransactionCost
			lowestCostBroker = r.brokers[name]
		}
	}
	
	if lowestCostBroker == nil {
		return nil, errors.New("no suitable broker found")
	}
	
	return lowestCostBroker, nil
}

// routeByRoundRobin distributes orders evenly among brokers
func (r *DefaultSmartRouter) routeByRoundRobin() (BrokerAdapter, error) {
	brokerNames := make([]string, 0, len(r.brokers))
	for name := range r.brokers {
		brokerNames = append(brokerNames, name)
	}
	
	sort.Strings(brokerNames) // Sort for deterministic order
	
	if len(brokerNames) == 0 {
		return nil, errors.New("no brokers available")
	}
	
	// Select broker using round-robin
	selectedBroker := brokerNames[r.roundRobinCounter%len(brokerNames)]
	r.roundRobinCounter++
	
	return r.brokers[selectedBroker], nil
}

// routeByVolumeWeighted routes based on volume available at each broker
func (r *DefaultSmartRouter) routeByVolumeWeighted(request *OrderRequest) (BrokerAdapter, error) {
	totalVolume := 0
	brokerVolumes := make(map[string]int)
	
	// Calculate total available volume across all brokers
	for name, metrics := range r.brokerMetrics {
		volume, exists := metrics.AvailableVolume[request.Symbol]
		if exists && volume > 0 {
			totalVolume += volume
			brokerVolumes[name] = volume
		}
	}
	
	if totalVolume == 0 {
		// Fall back to round-robin if no volume information is available
		return r.routeByRoundRobin()
	}
	
	// Select a broker based on its proportion of the total volume
	randomPoint := int(float64(totalVolume) * 0.5) // Use 0.5 as a simple random point
	cumulativeVolume := 0
	
	for name, volume := range brokerVolumes {
		cumulativeVolume += volume
		if randomPoint <= cumulativeVolume {
			return r.brokers[name], nil
		}
	}
	
	// Fallback to the broker with the highest volume
	var highestVolumeBroker string
	var highestVolume int = -1
	
	for name, volume := range brokerVolumes {
		if volume > highestVolume {
			highestVolume = volume
			highestVolumeBroker = name
		}
	}
	
	if highestVolumeBroker == "" {
		return nil, errors.New("no suitable broker found")
	}
	
	return r.brokers[highestVolumeBroker], nil
}

// GetBrokerMetrics returns the metrics for all brokers
func (r *DefaultSmartRouter) GetBrokerMetrics() map[string]*BrokerMetrics {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	// Create a copy to avoid concurrent access issues
	metricsCopy := make(map[string]*BrokerMetrics)
	for name, metrics := range r.brokerMetrics {
		metricsCopy[name] = &BrokerMetrics{
			AverageLatency:    metrics.AverageLatency,
			FillRate:          metrics.FillRate,
			TransactionCost:   metrics.TransactionCost,
			LastOrderTime:     metrics.LastOrderTime,
			OrderCount:        metrics.OrderCount,
			SuccessfulOrders:  metrics.SuccessfulOrders,
			FailedOrders:      metrics.FailedOrders,
			TotalOrdersPlaced: metrics.TotalOrdersPlaced,
		}
		
		// Copy available volume map
		metricsCopy[name].AvailableVolume = make(map[string]int)
		for symbol, volume := range metrics.AvailableVolume {
			metricsCopy[name].AvailableVolume[symbol] = volume
		}
	}
	
	return metricsCopy
}

// GetRoutingDecision returns the routing decision for a given order request
func (r *DefaultSmartRouter) GetRoutingDecision(request *OrderRequest) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	// Determine the strategy to use
	strategy := r.defaultStrategy
	if symbolStrategy, exists := r.symbolStrategies[request.Symbol]; exists {
		strategy = symbolStrategy
	}
	
	var brokerName string
	var err error
	
	// Get the broker based on the strategy
	switch strategy {
	case BestPrice:
		broker, err := r.routeByBestPrice(context.Background(), request)
		if err != nil {
			return "", err
		}
		// Find the broker name
		for name, b := range r.brokers {
			if b == broker {
				brokerName = name
				break
			}
		}
	case LowestLatency:
		broker, err := r.routeByLowestLatency(request)
		if err != nil {
			return "", err
		}
		// Find the broker name
		for name, b := range r.brokers {
			if b == broker {
				brokerName = name
				break
			}
		}
	case HighestFillRate:
		broker, err := r.routeByHighestFillRate(request)
		if err != nil {
			return "", err
		}
		// Find the broker name
		for name, b := range r.brokers {
			if b == broker {
				brokerName = name
				break
			}
		}
	case LowestCost:
		broker, err := r.routeByLowestCost(request)
		if err != nil {
			return "", err
		}
		// Find the broker name
		for name, b := range r.brokers {
			if b == broker {
				brokerName = name
				break
			}
		}
	case RoundRobin:
		broker, err := r.routeByRoundRobin()
		if err != nil {
			return "", err
		}
		// Find the broker name
		for name, b := range r.brokers {
			if b == broker {
				brokerName = name
				break
			}
		}
	case VolumeWeighted:
		broker, err := r.routeByVolumeWeighted(request)
		if err != nil {
			return "", err
		}
		// Find the broker name
		for name, b := range r.brokers {
			if b == broker {
				brokerName = name
				break
			}
		}
	default:
		return "", fmt.Errorf("unknown routing strategy: %s", strategy)
	}
	
	if brokerName == "" {
		return "", errors.New("could not determine broker name")
	}
	
	return brokerName, nil
}
