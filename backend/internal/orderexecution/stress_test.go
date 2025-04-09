package orderexecution

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

// BenchmarkConfig represents the configuration for benchmarking
type BenchmarkConfig struct {
	// Number of orders to execute
	NumOrders int
	
	// Number of concurrent clients
	NumClients int
	
	// Order distribution
	OrderDistribution string // "uniform", "burst", "random"
	
	// Order types to test
	OrderTypes []OrderType
	
	// Product types to test
	ProductTypes []ProductType
	
	// Symbols to use
	Symbols []string
	
	// Price range
	MinPrice float64
	MaxPrice float64
	
	// Quantity range
	MinQuantity int
	MaxQuantity int
	
	// Duration of the benchmark
	Duration time.Duration
	
	// Whether to enable broker latency simulation
	EnableLatencySimulation bool
	
	// Broker latency range (ms)
	MinLatencyMs int
	MaxLatencyMs int
	
	// Whether to enable broker error simulation
	EnableErrorSimulation bool
	
	// Broker error rate (percentage)
	ErrorRate int
	
	// Performance configuration
	PerformanceConfig *PerformanceConfig
}

// DefaultBenchmarkConfig returns the default benchmark configuration
func DefaultBenchmarkConfig() *BenchmarkConfig {
	return &BenchmarkConfig{
		NumOrders:              1000,
		NumClients:             10,
		OrderDistribution:      "uniform",
		OrderTypes:             []OrderType{OrderTypeLimit, OrderTypeMarket},
		ProductTypes:           []ProductType{ProductTypeDelivery, ProductTypeIntraday},
		Symbols:                []string{"RELIANCE", "TCS", "HDFCBANK", "INFY", "ICICIBANK"},
		MinPrice:               100.0,
		MaxPrice:               5000.0,
		MinQuantity:            10,
		MaxQuantity:            1000,
		Duration:               1 * time.Minute,
		EnableLatencySimulation: true,
		MinLatencyMs:           5,
		MaxLatencyMs:           50,
		EnableErrorSimulation:   true,
		ErrorRate:              5,
		PerformanceConfig:      DefaultPerformanceConfig(),
	}
}

// HighLoadBenchmarkConfig returns a benchmark configuration for high load testing
func HighLoadBenchmarkConfig() *BenchmarkConfig {
	return &BenchmarkConfig{
		NumOrders:              10000,
		NumClients:             50,
		OrderDistribution:      "burst",
		OrderTypes:             []OrderType{OrderTypeLimit, OrderTypeMarket, OrderTypeStopLoss},
		ProductTypes:           []ProductType{ProductTypeDelivery, ProductTypeIntraday},
		Symbols:                []string{"RELIANCE", "TCS", "HDFCBANK", "INFY", "ICICIBANK", "SBIN", "TATASTEEL", "WIPRO", "AXISBANK", "KOTAKBANK"},
		MinPrice:               100.0,
		MaxPrice:               5000.0,
		MinQuantity:            10,
		MaxQuantity:            1000,
		Duration:               5 * time.Minute,
		EnableLatencySimulation: true,
		MinLatencyMs:           10,
		MaxLatencyMs:           100,
		EnableErrorSimulation:   true,
		ErrorRate:              10,
		PerformanceConfig:      HighFrequencyPerformanceConfig(),
	}
}

// BenchmarkResult represents the result of a benchmark
type BenchmarkResult struct {
	// Configuration used
	Config *BenchmarkConfig
	
	// Total number of orders executed
	TotalOrders int
	
	// Number of successful orders
	SuccessfulOrders int
	
	// Number of failed orders
	FailedOrders int
	
	// Orders per second
	OrdersPerSecond float64
	
	// Average latency
	AverageLatency time.Duration
	
	// Latency percentiles
	LatencyP50 time.Duration
	LatencyP90 time.Duration
	LatencyP99 time.Duration
	
	// CPU usage
	CPUUsage float64
	
	// Memory usage
	MemoryUsage uint64
	
	// Error distribution
	ErrorDistribution map[string]int
	
	// Duration of the benchmark
	Duration time.Duration
}

// BenchmarkRunner runs benchmarks for the order execution engine
type BenchmarkRunner struct {
	config *BenchmarkConfig
	logger Logger
}

// NewBenchmarkRunner creates a new benchmark runner
func NewBenchmarkRunner(config *BenchmarkConfig, logger Logger) *BenchmarkRunner {
	return &BenchmarkRunner{
		config: config,
		logger: logger,
	}
}

// RunBenchmark runs a benchmark
func (r *BenchmarkRunner) RunBenchmark(engine *HighPerformanceOrderExecutionEngine) (*BenchmarkResult, error) {
	r.logger.Info("Starting benchmark",
		"numOrders", r.config.NumOrders,
		"numClients", r.config.NumClients,
		"distribution", r.config.OrderDistribution,
		"duration", r.config.Duration,
	)
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), r.config.Duration)
	defer cancel()
	
	// Create result
	result := &BenchmarkResult{
		Config:            r.config,
		ErrorDistribution: make(map[string]int),
	}
	
	// Create order generator
	generator := NewOrderGenerator(r.config)
	
	// Create latency tracker
	tracker := NewLatencyTracker()
	
	// Create client workers
	var wg sync.WaitGroup
	orderCh := make(chan Order, r.config.NumOrders)
	resultCh := make(chan OrderResult, r.config.NumOrders)
	
	// Start client workers
	for i := 0; i < r.config.NumClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			r.runClient(ctx, engine, orderCh, resultCh, clientID, tracker)
		}(i)
	}
	
	// Start order generator
	go generator.GenerateOrders(ctx, orderCh)
	
	// Start result collector
	go r.collectResults(ctx, resultCh, result)
	
	// Wait for benchmark to complete
	startTime := time.Now()
	
	// Wait for context to be done (timeout or cancellation)
	<-ctx.Done()
	
	// Close order channel to signal clients to stop
	close(orderCh)
	
	// Wait for all clients to finish
	wg.Wait()
	
	// Close result channel
	close(resultCh)
	
	// Calculate benchmark duration
	result.Duration = time.Since(startTime)
	
	// Calculate orders per second
	result.OrdersPerSecond = float64(result.TotalOrders) / result.Duration.Seconds()
	
	// Calculate latency percentiles
	latencies := tracker.GetLatencies()
	if len(latencies) > 0 {
		result.AverageLatency = tracker.GetAverageLatency()
		result.LatencyP50 = tracker.GetPercentileLatency(50)
		result.LatencyP90 = tracker.GetPercentileLatency(90)
		result.LatencyP99 = tracker.GetPercentileLatency(99)
	}
	
	// Get resource usage
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	result.MemoryUsage = memStats.Alloc
	
	// Log benchmark results
	r.logger.Info("Benchmark completed",
		"totalOrders", result.TotalOrders,
		"successfulOrders", result.SuccessfulOrders,
		"failedOrders", result.FailedOrders,
		"ordersPerSecond", result.OrdersPerSecond,
		"averageLatency", result.AverageLatency,
		"latencyP50", result.LatencyP50,
		"latencyP90", result.LatencyP90,
		"latencyP99", result.LatencyP99,
		"memoryUsage", result.MemoryUsage,
		"duration", result.Duration,
	)
	
	return result, nil
}

// runClient runs a client worker
func (r *BenchmarkRunner) runClient(
	ctx context.Context,
	engine *HighPerformanceOrderExecutionEngine,
	orderCh <-chan Order,
	resultCh chan<- OrderResult,
	clientID int,
	tracker *LatencyTracker,
) {
	r.logger.Info("Starting client worker", "clientID", clientID)
	
	for {
		select {
		case <-ctx.Done():
			// Context cancelled or timed out
			return
		case order, ok := <-orderCh:
			if !ok {
				// Channel closed
				return
			}
			
			// Execute order and track latency
			startTime := time.Now()
			response, err := engine.ExecuteOrder(ctx, order)
			latency := time.Since(startTime)
			
			// Track latency
			tracker.AddLatency(latency)
			
			// Send result
			result := OrderResult{
				Order:    order,
				Response: response,
				Error:    err,
				Latency:  latency,
			}
			resultCh <- result
		}
	}
}

// collectResults collects benchmark results
func (r *BenchmarkRunner) collectResults(
	ctx context.Context,
	resultCh <-chan OrderResult,
	result *BenchmarkResult,
) {
	for {
		select {
		case <-ctx.Done():
			// Context cancelled or timed out
			return
		case res, ok := <-resultCh:
			if !ok {
				// Channel closed
				return
			}
			
			// Update result
			result.TotalOrders++
			
			if res.Error != nil {
				result.FailedOrders++
				
				// Update error distribution
				errorType := "unknown"
				if execErr, ok := res.Error.(*ExecutionError); ok {
					errorType = string(execErr.Type)
				}
				result.ErrorDistribution[errorType]++
			} else {
				result.SuccessfulOrders++
			}
		}
	}
}

// OrderResult represents the result of an order execution
type OrderResult struct {
	Order    Order
	Response OrderResponse
	Error    error
	Latency  time.Duration
}

// OrderGenerator generates orders for benchmarking
type OrderGenerator struct {
	config *BenchmarkConfig
	rand   *rand.Rand
}

// NewOrderGenerator creates a new order generator
func NewOrderGenerator(config *BenchmarkConfig) *OrderGenerator {
	return &OrderGenerator{
		config: config,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateOrders generates orders and sends them to the order channel
func (g *OrderGenerator) GenerateOrders(ctx context.Context, orderCh chan<- Order) {
	switch g.config.OrderDistribution {
	case "uniform":
		g.generateUniformOrders(ctx, orderCh)
	case "burst":
		g.generateBurstOrders(ctx, orderCh)
	case "random":
		g.generateRandomOrders(ctx, orderCh)
	default:
		g.generateUniformOrders(ctx, orderCh)
	}
}

// generateUniformOrders generates orders at a uniform rate
func (g *OrderGenerator) generateUniformOrders(ctx context.Context, orderCh chan<- Order) {
	// Calculate interval between orders
	interval := g.config.Duration / time.Duration(g.config.NumOrders)
	
	// Create ticker
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	// Generate orders
	for i := 0; i < g.config.NumOrders; i++ {
		select {
		case <-ctx.Done():
			// Context cancelled or timed out
			return
		case <-ticker.C:
			// Generate order
			order := g.generateOrder(i)
			
			// Send order
			select {
			case <-ctx.Done():
				return
			case orderCh <- order:
				// Order sent
			}
		}
	}
}

// generateBurstOrders generates orders in bursts
func (g *OrderGenerator) generateBurstOrders(ctx context.Context, orderCh chan<- Order) {
	// Calculate number of bursts
	numBursts := 10
	ordersPerBurst := g.config.NumOrders / numBursts
	burstInterval := g.config.Duration / time.Duration(numBursts)
	
	// Create ticker
	ticker := time.NewTicker(burstInterval)
	defer ticker.Stop()
	
	// Generate orders
	orderCount := 0
	for i := 0; i < numBursts; i++ {
		select {
		case <-ctx.Done():
			// Context cancelled or timed out
			return
		case <-ticker.C:
			// Generate burst of orders
			for j := 0; j < ordersPerBurst; j++ {
				// Generate order
				order := g.generateOrder(orderCount)
				orderCount++
				
				// Send order
				select {
				case <-ctx.Done():
					return
				case orderCh <- order:
					// Order sent
				}
			}
		}
	}
}

// generateRandomOrders generates orders at random intervals
func (g *OrderGenerator) generateRandomOrders(ctx context.Context, orderCh chan<- Order) {
	// Calculate average interval between orders
	avgInterval := g.config.Duration / time.Duration(g.config.NumOrders)
	
	// Generate orders
	for i := 0; i < g.config.NumOrders; i++ {
		// Calculate random interval
		interval := time.Duration(float64(avgInterval) * (0.5 + g.rand.Float64()))
		
		// Wait for interval
		select {
		case <-ctx.Done():
			// Context cancelled or timed out
			return
		case <-time.After(interval):
			// Generate order
			order := g.generateOrder(i)
			
			// Send order
			select {
			case <-ctx.Done():
				return
			case orderCh <- order:
				// Order sent
			}
		}
	}
}

// generateOrder generates a random order
func (g *OrderGenerator) generateOrder(index int) Order {
	// Select random values
	orderType := g.config.OrderTypes[g.rand.Intn(len(g.config.OrderTypes))]
	productType := g.config.ProductTypes[g.rand.Intn(len(g.config.ProductTypes))]
	symbol := g.config.Symbols[g.rand.Intn(len(g.config.Symbols))]
	side := OrderSideBuy
	if g.rand.Intn(2) == 1 {
		side = OrderSideSell
	}
	
	// Generate random price and quantity
	price := g.config.MinPrice + g.rand.Float64()*(g.config.MaxPrice-g.config.MinPrice)
	quantity := g.config.MinQuantity + g.rand.Intn(g.config.MaxQuantity-g.config.MinQuantity+1)
	
	// Create order
	order := Order{
		ID:          fmt.Sprintf("BENCH-%06d", index+1),
		PortfolioID: fmt.Sprintf("PORTFOLIO-%03d", g.rand.Intn(10)+1),
		StrategyID:  fmt.Sprintf("STRATEGY-%03d", g.rand.Intn(5)+1),
		Symbol:      symbol,
		Exchange:    "NSE",
		OrderType:   orderType,
		ProductType: productType,
		Side:        side,
		Quantity:    quantity,
		Price:       price,
		Status:      OrderStatusNew,
	}
	
	return order
}

// LatencyTracker tracks order execution latencies
type LatencyTracker struct {
	latencies []time.Duration
	mutex     sync.RWMutex
}

// NewLatencyTracker creates a new latency tracker
func NewLatencyTracker() *LatencyTracker {
	return &LatencyTracker{
		latencies: make([]time.Duration, 0, 1000),
	}
}

// AddLatency adds a latency measurement
func (t *LatencyTracker) AddLatency(latency time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	
	t.latencies = append(t.latencies, latency)
}

// GetLatencies gets all latency measurements
func (t *LatencyTracker) GetLatencies() []time.Duration {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	
	// Create a copy of the latencies
	latencies := make([]time.Duration, len(t.latencies))
	copy(latencies, t.latencies)
	
	return latencies
}

// GetAverageLatency gets the average latency
func (t *LatencyTracker) GetAverageLatency() time.Duration {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	
	if len(t.latencies) == 0 {
		return 0
	}
	
	var sum time.Duration
	for _, latency := range t.latencies {
		sum += latency
	}
	
	return sum / time.Duration(len(t.latencies))
}

// GetPercentileLatency gets the latency at the specified percentile
func (t *LatencyTracker) GetPercentileLatency(percentile int) time.Duration {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	
	if len(t.latencies) == 0 {
		return 0
	}
	
	// Sort latencies
	sorted := make([]time.Duration, len(t.latencies))
	copy(sorted, t.latencies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	
	// Calculate index
	index := len(sorted) * percentile / 100
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	
	return sorted[index]
}

// StressTest represents a stress test for the order execution engine
type StressTest struct {
	Name        string
	Description string
	Config      *BenchmarkConfig
	Duration    time.Duration
	Validator   func(result *BenchmarkResult) (bool, string)
}

// RunStressTests runs a series of stress tests
func RunStressTests(t *testing.T, engine *HighPerformanceOrderExecutionEngine, logger Logger) {
	tests := []StressTest{
		{
			Name:        "BasicLoadTest",
			Description: "Basic load test with moderate order volume",
			Config:      DefaultBenchmarkConfig(),
			Duration:    30 * time.Second,
			Validator: func(result *BenchmarkResult) (bool, string) {
				if result.OrdersPerSecond < 10 {
					return false, fmt.Sprintf("Orders per second too low: %.2f", result.OrdersPerSecond)
				}
				if float64(result.FailedOrders)/float64(result.TotalOrders) > 0.1 {
					return false, fmt.Sprintf("Failure rate too high: %.2f%%", float64(result.FailedOrders)/float64(result.TotalOrders)*100)
				}
				return true, ""
			},
		},
		{
			Name:        "HighVolumeTest",
			Description: "High volume test with many concurrent clients",
			Config: &BenchmarkConfig{
				NumOrders:              5000,
				NumClients:             50,
				OrderDistribution:      "uniform",
				OrderTypes:             []OrderType{OrderTypeLimit, OrderTypeMarket},
				ProductTypes:           []ProductType{ProductTypeDelivery, ProductTypeIntraday},
				Symbols:                []string{"RELIANCE", "TCS", "HDFCBANK", "INFY", "ICICIBANK"},
				MinPrice:               100.0,
				MaxPrice:               5000.0,
				MinQuantity:            10,
				MaxQuantity:            1000,
				Duration:               30 * time.Second,
				EnableLatencySimulation: true,
				MinLatencyMs:           5,
				MaxLatencyMs:           50,
				EnableErrorSimulation:   true,
				ErrorRate:              5,
				PerformanceConfig:      HighFrequencyPerformanceConfig(),
			},
			Duration: 30 * time.S
(Content truncated due to size limit. Use line ranges to read in chunks)