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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trading-platform/backend/internal/messagequeue"
)

// TestConfig is imported from websocket_order_load_test.go

// loadConfig is imported from websocket_order_load_test.go

// LoadTestResult is imported from websocket_order_load_test.go

// TestMessageQueueLoad tests the message queue system under load
func TestMessageQueueLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	// Load configuration
	config, err := loadConfig()
	require.NoError(t, err)

	// Create Redis client (mock for testing)
	redisClient := &MockRedisClient{}

	// Create RabbitMQ client (mock for testing)
	rabbitMQClient := &MockRabbitMQClient{}

	// Create message queue service
	mqService := messagequeue.NewService(redisClient, rabbitMQClient)

	// Initialize the service
	err = mqService.Initialize()
	require.NoError(t, err)

	// Test parameters
	numPublishers := 20
	numSubscribers := 5
	messagesPerPublisher := 1000
	topicCount := 10

	// Create a wait group to wait for all publishers and subscribers to finish
	var wg sync.WaitGroup
	wg.Add(numPublishers + numSubscribers)

	// Create a channel to collect results
	resultChan := make(chan time.Duration, numPublishers*messagesPerPublisher)

	// Create a channel to signal subscribers to stop
	done := make(chan struct{})

	// Set up the mock Redis client
	redisClient.PublishFunc = func(channel string, message []byte) error {
		// Simulate processing time (0-5ms)
		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
		return nil
	}

	redisClient.SubscribeFunc = func(channel string, callback func([]byte)) error {
		// Start a goroutine to simulate receiving messages
		go func() {
			ticker := time.NewTicker(1 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					// Simulate receiving a message
					callback([]byte(fmt.Sprintf("test message on %s", channel)))
				}
			}
		}()
		return nil
	}

	// Set up the mock RabbitMQ client
	rabbitMQClient.PublishFunc = func(exchange string, routingKey string, message []byte) error {
		// Simulate processing time (1-10ms)
		time.Sleep(time.Duration(rand.Intn(10)+1) * time.Millisecond)
		return nil
	}

	rabbitMQClient.ConsumeFunc = func(queueName string, callback func([]byte)) error {
		// Start a goroutine to simulate receiving messages
		go func() {
			ticker := time.NewTicker(2 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					// Simulate receiving a message
					callback([]byte(fmt.Sprintf("test message on %s", queueName)))
				}
			}
		}()
		return nil
	}

	// Start time
	startTime := time.Now()

	// Start subscribers
	for i := 0; i < numSubscribers; i++ {
		go func(subscriberID int) {
			defer wg.Done()

			// Subscribe to market data for different symbols
			for j := 0; j < topicCount; j++ {
				symbol := fmt.Sprintf("SYMBOL%d", j)
				err := mqService.SubscribeToMarketData(symbol, func(data []byte) {
					// Just receive the data
				})
				require.NoError(t, err)
			}

			// Consume orders
			err := mqService.ConsumeOrders(func(data []byte) {
				// Just receive the data
			})
			require.NoError(t, err)

			// Keep the subscriber running until signaled to stop
			<-done
		}(i)
	}

	// Start publishers
	for i := 0; i < numPublishers; i++ {
		go func(publisherID int) {
			defer wg.Done()

			ctx := context.Background()

			// Publish messages
			for j := 0; j < messagesPerPublisher; j++ {
				// Determine if this is a market data or order message
				isMarketData := rand.Intn(2) == 0

				// Record start time
				start := time.Now()

				if isMarketData {
					// Publish market data
					symbol := fmt.Sprintf("SYMBOL%d", rand.Intn(topicCount))
					message := fmt.Sprintf("market data %d from publisher %d", j, publisherID)
					err := mqService.PublishMarketData(symbol, []byte(message))
					require.NoError(t, err)
				} else {
					// Publish order
					orderType := "new"
					if rand.Intn(3) == 0 {
						orderType = "cancel"
					}
					message := fmt.Sprintf("order %d from publisher %d", j, publisherID)
					err := mqService.PublishOrder(orderType, []byte(message))
					require.NoError(t, err)
				}

				// Record response time
				resultChan <- time.Since(start)

				// Small delay between messages
				time.Sleep(time.Duration(rand.Intn(2)) * time.Millisecond)
			}
		}(i)
	}

	// Wait for all publishers to finish
	publishersDone := make(chan struct{})
	go func() {
		for i := 0; i < numPublishers; i++ {
			<-publishersDone
		}
		close(done) // Signal subscribers to stop
	}()

	// Wait for all publishers and subscribers to finish
	wg.Wait()

	// Calculate results
	totalDuration := time.Since(startTime)
	totalRequests := numPublishers * messagesPerPublisher

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
	t.Logf("Message Queue Load Test Results:")
	t.Logf("  Total Requests: %d", result.TotalRequests)
	t.Logf("  Successful Requests: %d", result.SuccessfulRequests)
	t.Logf("  Failed Requests: %d", result.FailedRequests)
	t.Logf("  Min Response Time: %v", result.MinResponseTime)
	t.Logf("  Max Response Time: %v", result.MaxResponseTime)
	t.Logf("  Avg Response Time: %v", result.AvgResponseTime)
	t.Logf("  Total Duration: %v", result.TotalDuration)
	t.Logf("  Requests Per Second: %.2f", result.RequestsPerSecond)

	// Assert that the performance meets requirements
	assert.GreaterOrEqual(t, result.RequestsPerSecond, 5000.0, "Message queue system should handle at least 5000 messages per second")
	assert.Less(t, result.AvgResponseTime, 10*time.Millisecond, "Average message queue response time should be less than 10ms")

	// Clean up
	err = mqService.Close()
	require.NoError(t, err)
}

// MockRedisClient is a mock implementation of the Redis client for load testing
type MockRedisClient struct {
	ConnectFunc    func() error
	CloseFunc      func() error
	PublishFunc    func(channel string, message []byte) error
	SubscribeFunc  func(channel string, callback func([]byte)) error
	UnsubscribeFunc func(channel string) error
}

func (m *MockRedisClient) Connect() error {
	if m.ConnectFunc != nil {
		return m.ConnectFunc()
	}
	return nil
}

func (m *MockRedisClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *MockRedisClient) Publish(channel string, message []byte) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(channel, message)
	}
	return nil
}

func (m *MockRedisClient) Subscribe(channel string, callback func([]byte)) error {
	if m.SubscribeFunc != nil {
		return m.SubscribeFunc(channel, callback)
	}
	return nil
}

func (m *MockRedisClient) Unsubscribe(channel string) error {
	if m.UnsubscribeFunc != nil {
		return m.UnsubscribeFunc(channel)
	}
	return nil
}

// MockRabbitMQClient is a mock implementation of the RabbitMQ client for load testing
type MockRabbitMQClient struct {
	ConnectFunc         func() error
	CloseFunc           func() error
	DeclareExchangeFunc func(name string, exchangeType string) error
	DeclareQueueFunc    func(name string) error
	BindQueueFunc       func(queueName string, exchangeName string, routingKey string) error
	PublishFunc         func(exchange string, routingKey string, message []byte) error
	ConsumeFunc         func(queueName string, callback func([]byte)) error
}

func (m *MockRabbitMQClient) Connect() error {
	if m.ConnectFunc != nil {
		return m.ConnectFunc()
	}
	return nil
}

func (m *MockRabbitMQClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *MockRabbitMQClient) DeclareExchange(name string, exchangeType string) error {
	if m.DeclareExchangeFunc != nil {
		return m.DeclareExchangeFunc(name, exchangeType)
	}
	return nil
}

func (m *MockRabbitMQClient) DeclareQueue(name string) error {
	if m.DeclareQueueFunc != nil {
		return m.DeclareQueueFunc(name)
	}
	return nil
}

func (m *MockRabbitMQClient) BindQueue(queueName string, exchangeName string, routingKey string) error {
	if m.BindQueueFunc != nil {
		return m.BindQueueFunc(queueName, exchangeName, routingKey)
	}
	return nil
}

func (m *MockRabbitMQClient) Publish(exchange string, routingKey string, message []byte) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(exchange, routingKey, message)
	}
	return nil
}

func (m *MockRabbitMQClient) Consume(queueName string, callback func([]byte)) error {
	if m.ConsumeFunc != nil {
		return m.ConsumeFunc(queueName, callback)
	}
	return nil
}
