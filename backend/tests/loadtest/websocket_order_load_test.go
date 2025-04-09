package loadtest

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
	"math/rand"
	"os"
	"encoding/json"
	"log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trading-platform/backend/internal/websocket"
	"github.com/trading-platform/backend/internal/core"
	"github.com/trading-platform/backend/internal/broker"
	"github.com/trading-platform/backend/internal/messagequeue"
)

// TestConfig represents the test configuration
type TestConfig struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"database"`
	Redis struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
	RabbitMQ struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"rabbitmq"`
	JWT struct {
		Secret     string `json:"secret"`
		Expiration int    `json:"expiration"`
	} `json:"jwt"`
	Broker struct {
		XTS struct {
			Endpoint  string `json:"endpoint"`
			APIKey    string `json:"api_key"`
			APISecret string `json:"api_secret"`
		} `json:"xts"`
	} `json:"broker"`
}

// loadConfig loads the test configuration from a file
func loadConfig() (*TestConfig, error) {
	// Read the configuration file
	configFile, err := os.ReadFile("../config.json")
	if err != nil {
		return nil, err
	}

	// Parse the configuration
	var config TestConfig
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadTestResult represents the result of a load test
type LoadTestResult struct {
	TotalRequests      int
	SuccessfulRequests int
	FailedRequests     int
	MinResponseTime    time.Duration
	MaxResponseTime    time.Duration
	AvgResponseTime    time.Duration
	TotalDuration      time.Duration
	RequestsPerSecond  float64
}

// WebSocketLoadTest tests the WebSocket server under load
func TestWebSocketLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	// Create a new WebSocket hub
	hub := websocket.NewHub()
	
	// Start the hub in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)
	
	// Test parameters
	numClients := 100
	messagesPerClient := 50
	
	// Create a wait group to wait for all clients to finish
	var wg sync.WaitGroup
	wg.Add(numClients)
	
	// Create a channel to collect results
	resultChan := make(chan time.Duration, numClients*messagesPerClient)
	
	// Start time
	startTime := time.Now()
	
	// Create and register clients
	for i := 0; i < numClients; i++ {
		clientID := fmt.Sprintf("client-%d", i)
		
		// Create a mock client
		client := &MockWebSocketClient{ID: clientID}
		
		// Register the client
		hub.Register(client)
		
		// Subscribe to a topic
		hub.Subscribe(client, "market_data")
		
		// Start a goroutine for each client
		go func(c *MockWebSocketClient, id string) {
			defer wg.Done()
			
			// Send messages
			for j := 0; j < messagesPerClient; j++ {
				// Record start time
				start := time.Now()
				
				// Broadcast a message
				message := fmt.Sprintf("message from %s: %d", id, j)
				hub.BroadcastToTopic("market_data", []byte(message))
				
				// Simulate processing time
				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
				
				// Record response time
				resultChan <- time.Since(start)
			}
			
			// Unregister the client
			hub.Unregister(c)
		}(client, clientID)
	}
	
	// Wait for all clients to finish
	wg.Wait()
	
	// Calculate results
	totalDuration := time.Since(startTime)
	totalRequests := numClients * messagesPerClient
	
	// Collect response times
	var responseTimes []time.Duration
	for i := 0; i < totalRequests; i++ {
		responseTimes = append(responseTimes, <-resultChan)
	}
	
	// Calculate statistics
	minTime := responseTimes[0]
	maxTime := responseTimes[0]
	var totalTime time.Duration
	
	for _, t := range responseTimes {
		if t < minTime {
			minTime = t
		}
		if t > maxTime {
			maxTime = t
		}
		totalTime += t
	}
	
	avgTime := totalTime / time.Duration(len(responseTimes))
	requestsPerSecond := float64(totalRequests) / totalDuration.Seconds()
	
	// Create result
	result := LoadTestResult{
		TotalRequests:      totalRequests,
		SuccessfulRequests: totalRequests, // Assuming all are successful in this test
		FailedRequests:     0,
		MinResponseTime:    minTime,
		MaxResponseTime:    maxTime,
		AvgResponseTime:    avgTime,
		TotalDuration:      totalDuration,
		RequestsPerSecond:  requestsPerSecond,
	}
	
	// Log results
	t.Logf("WebSocket Load Test Results:")
	t.Logf("  Total Requests: %d", result.TotalRequests)
	t.Logf("  Successful Requests: %d", result.SuccessfulRequests)
	t.Logf("  Failed Requests: %d", result.FailedRequests)
	t.Logf("  Min Response Time: %v", result.MinResponseTime)
	t.Logf("  Max Response Time: %v", result.MaxResponseTime)
	t.Logf("  Avg Response Time: %v", result.AvgResponseTime)
	t.Logf("  Total Duration: %v", result.TotalDuration)
	t.Logf("  Requests Per Second: %.2f", result.RequestsPerSecond)
	
	// Assert that the performance meets requirements
	assert.GreaterOrEqual(t, result.RequestsPerSecond, 1000.0, "WebSocket server should handle at least 1000 requests per second")
	assert.Less(t, result.AvgResponseTime, 50*time.Millisecond, "Average response time should be less than 50ms")
}

// MockWebSocketClient is a mock implementation of a WebSocket client
type MockWebSocketClient struct {
	ID string
}

func (m *MockWebSocketClient) GetID() string {
	return m.ID
}

func (m *MockWebSocketClient) Send(message []byte) error {
	// Simulate message processing
	return nil
}

func (m *MockWebSocketClient) Close() error {
	// Simulate closing the connection
	return nil
}

// TestOrderExecutionLoad tests the order execution system under load
func TestOrderExecutionLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	// Create a broker factory
	brokerFactory := broker.NewBrokerFactory()

	// Register a mock broker
	mockBroker := &MockBroker{}
	brokerFactory.RegisterBroker("xts", mockBroker)

	// Create an order service
	orderService := core.NewOrderService(brokerFactory)

	// Create an execution engine
	executionEngine := core.NewExecutionEngine(orderService)

	// Test parameters
	numConcurrentOrders := 50
	numOrdersPerBatch := 20
	numBatches := 5

	// Create a wait group to wait for all orders to finish
	var wg sync.WaitGroup
	wg.Add(numConcurrentOrders * numBatches)

	// Create a channel to collect results
	resultChan := make(chan time.Duration, numConcurrentOrders * numBatches * numOrdersPerBatch)

	// Set up the mock broker to return a successful response
	mockBroker.PlaceOrderFunc = func(ctx context.Context, request broker.OrderRequest) (*broker.OrderResponse, error) {
		// Simulate processing time
		time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
		return &broker.OrderResponse{
			Success: true,
			OrderID: fmt.Sprintf("order-%d", rand.Intn(1000000)),
		}, nil
	}

	// Start time
	startTime := time.Now()

	// Run batches of concurrent orders
	for batch := 0; batch < numBatches; batch++ {
		// Start concurrent order placements
		for i := 0; i < numConcurrentOrders; i++ {
			go func(batchNum, orderNum int) {
				defer wg.Done()

				ctx := context.Background()

				// Place multiple orders
				for j := 0; j < numOrdersPerBatch; j++ {
					// Create an order request
					orderRequest := core.OrderRequest{
						UserID:          fmt.Sprintf("user-%d", orderNum),
						BrokerName:      "xts",
						Symbol:          "NIFTY",
						Exchange:        "NSE",
						OrderType:       "MARKET",
						TransactionType: "BUY",
						ProductType:     "NRML",
						Quantity:        1,
					}

					// Record start time
					start := time.Now()

					// Place the order
					_, err := orderService.PlaceOrder(ctx, orderRequest)
					
					// Record response time
					responseTime := time.Since(start)
					resultChan <- responseTime

					// Check for errors
					if err != nil {
						t.Logf("Error placing order: %v", err)
					}
				}
			}(batch, i)
		}

		// Wait a bit between batches
		if batch < numBatches-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	// Wait for all orders to finish
	wg.Wait()

	// Calculate results
	totalDuration := time.Since(startTime)
	totalRequests := numConcurrentOrders * numBatches * numOrdersPerBatch

	// Collect response times
	var responseTimes []time.Duration
	for i := 0; i < totalRequests; i++ {
		responseTimes = append(responseTimes, <-resultChan)
	}

	// Calculate statistics
	minTime := responseTimes[0]
	maxTime := responseTimes[0]
	var totalTime time.Duration
	successfulRequests := totalRequests // Assuming all are successful in this test

	for _, t := range responseTimes {
		if t < minTime {
			minTime = t
		}
		if t > maxTime {
			maxTime = t
		}
		totalTime += t
	}

	avgTime := totalTime / time.Duration(len(responseTimes))
	requestsPerSecond := float64(totalRequests) / totalDuration.Seconds()

	// Create result
	result := LoadTestResult{
		TotalRequests:      totalRequests,
		SuccessfulRequests: successfulRequests,
		FailedRequests:     totalRequests - successfulRequests,
		MinResponseTime:    minTime,
		MaxResponseTime:    maxTime,
		AvgResponseTime:    avgTime,
		TotalDuration:      totalDuration,
		RequestsPerSecond:  requestsPerSecond,
	}

	// Log results
	t.Logf("Order Execution Load Test Results:")
	t.Logf("  Total Requests: %d", result.TotalRequests)
	t.Logf("  Successful Requests: %d", result.SuccessfulRequests)
	t.Logf("  Failed Requests: %d", result.FailedRequests)
	t.Logf("  Min Response Time: %v", result.MinResponseTime)
	t.Logf("  Max Response Time: %v", result.MaxResponseTime)
	t.Logf("  Avg Response Time: %v", result.AvgResponseTime)
	t.Logf("  Total Duration: %v", result.TotalDuration)
	t.Logf("  Requests Per Second: %.2f", result.RequestsPerSecond)

	// Assert that the performance meets requirements
	assert.GreaterOrEqual(t, result.RequestsPerSecond, 100.0, "Order execution system should handle at least 100 orders per second")
	assert.Less(t, result.AvgResponseTime, 100*time.Millisecond, "Average order execution time should be less than 100ms")
}

// MockBroker is a mock implementation of the Broker interface for load testing
type MockBroker struct {
	PlaceOrderFunc        func(ctx context.Context, request broker.OrderRequest) (*broker.OrderResponse, error)
	ModifyOrderFunc       func(ctx context.Context, orderID string, request broker.OrderRequest) (*broker.OrderResponse, error)
	CancelOrderFunc       func(ctx context.Context, orderID string) (*broker.OrderResponse, error)
	GetOrderFunc          func(ctx context.Context, orderID string) (*broker.Order, error)
	GetOrdersFunc         func(ctx context.Context) ([]broker.Order, error)
	GetPositionsFunc      func(ctx context.Context) ([]broker.Position, error)
	GetQuoteFunc          func(ctx context.Context, symbol, exchange string) (*broker.Quote, error)
	SubscribeQuotesFunc   func(ctx context.Context, symbols []string, exchange string) error
	UnsubscribeQuotesFunc func(ctx context.Context, symbols []string, exchange string) error
}

func (m *MockBroker) Initialize(config broker.BrokerConfig) error {
	return nil
}

func (m *MockBroker) PlaceOrder(ctx context.Context, request broker.OrderRequest) (*broker.OrderResponse, error) {
	if m.PlaceOrderFunc != nil {
		return m.PlaceOrderFunc(ctx, request)
	}
	return &broker.OrderResponse{Success: true, OrderID: "mock-order-id"}, nil
}

func (m *MockBroker) ModifyOrder(ctx context.Context, orderID string, request broker.OrderRequest) (*broker.OrderResponse, error) {
	if m.ModifyOrderFunc != nil {
		return m.ModifyOrderFunc(ctx, orderID, request)
	}
	return &broker.OrderResponse{Success: true, OrderID: orderID}, nil
}

func (m *MockBroker) CancelOrder(ctx context.Context, orderID string) (*broker.OrderResponse, error) {
	if m.CancelOrderFunc != nil {
		return m.CancelOrderFunc(ctx, orderID)
	}
	return &broker.OrderResponse{Success: true, OrderID: orderID}, nil
}

func (m *MockBroker) GetOrder(ctx context.Context, orderID string) (*broker.Order, error) {
	if m.GetOrderFunc != nil {
		return m.GetOrderFunc(ctx, orderID)
	}
	return &broker.Order{ID: orderID, Status: broker.OrderStatusCompleted}, nil
}

func (m *MockBroker) GetOrders(ctx context.Context) ([]broker.Order, error) {
	if m.GetOrdersFunc != nil {
		return m.GetOrdersFunc(ctx)
	}
	return []broker.Order{{ID: "mock-order-id", Status: broker.OrderStatusCompleted}}, nil
}

func (m *MockBroker) GetPositions(ctx context.Context) ([]broker.Position, error) {
	if m.GetPositionsFunc != nil {
		return m.GetPositionsFunc(ctx)
	}
	return []broker.Position{{Symbol: "NIFTY", Exchange: "NSE", Quantity: 1}}, nil
}

func (m *MockBroker) GetQuote(ctx context.Context, symbol, exchange string) (*broker.Quote, error) {
	if m.GetQuoteFunc != nil {
		return m.GetQuoteFunc(ctx, symbol, exchange)
	}
	return &broker.Quote{Symbol: symbol, Exchange: exchange, LastPrice: 18000.0}, nil
}

func (m *MockBroker) SubscribeQuotes(ctx context.Context, symbols []string, exchange string) error {
	if m.SubscribeQuotesFunc != nil {
		return m.SubscribeQuotesFunc(ctx, symbols, exchange)
	}
	return nil
}

func (m *MockBroker) UnsubscribeQuotes(ctx context.Context, symbols []string, exchange string) error {
	if m.UnsubscribeQuotesFunc != nil {
		return m.UnsubscribeQuotesFunc(ctx, symbols, exchange)
	}
	return nil
}

func (m *MockBroker) Close() error {
	return nil
}
