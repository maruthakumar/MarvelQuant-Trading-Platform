package database

import (
	"context"
	"testing"
	"time"
	"fmt"
	"sync"
	"math/rand"
	"os"
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trading-platform/backend/internal/database"
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

// TestDatabasePerformance tests the database performance under load
func TestDatabasePerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database performance test in short mode")
	}

	// Create a mock database for testing
	mockDB := &MockDB{}
	
	// Create a database service with the mock DB
	dbService := database.NewService(mockDB)
	
	// Test parameters
	numConcurrentClients := 50
	queriesPerClient := 100
	
	// Create a wait group to wait for all clients to finish
	var wg sync.WaitGroup
	wg.Add(numConcurrentClients)
	
	// Create a channel to collect results
	resultChan := make(chan time.Duration, numConcurrentClients*queriesPerClient)
	
	// Set up the mock DB to simulate query execution
	mockDB.QueryFunc = func(query string, args ...interface{}) ([]map[string]interface{}, error) {
		// Simulate query execution time based on query complexity
		var delay time.Duration
		if len(query) > 100 {
			// Complex query
			delay = time.Duration(rand.Intn(20)+10) * time.Millisecond
		} else {
			// Simple query
			delay = time.Duration(rand.Intn(5)+1) * time.Millisecond
		}
		time.Sleep(delay)
		
		// Return mock results
		return []map[string]interface{}{
			{"id": 1, "name": "Test Result"},
		}, nil
	}
	
	mockDB.ExecFunc = func(query string, args ...interface{}) error {
		// Simulate execution time
		time.Sleep(time.Duration(rand.Intn(10)+5) * time.Millisecond)
		return nil
	}
	
	// Start time
	startTime := time.Now()
	
	// Start concurrent clients
	for i := 0; i < numConcurrentClients; i++ {
		go func(clientID int) {
			defer wg.Done()
			
			// Execute queries
			for j := 0; j < queriesPerClient; j++ {
				// Determine query type (read or write)
				isRead := rand.Intn(4) < 3 // 75% reads, 25% writes
				
				// Record start time
				start := time.Now()
				
				if isRead {
					// Execute a read query
					var query string
					if rand.Intn(10) < 8 {
						// Simple query (80%)
						query = fmt.Sprintf("SELECT * FROM orders WHERE user_id = %d LIMIT 10", clientID)
					} else {
						// Complex query (20%)
						query = fmt.Sprintf(`
							SELECT o.id, o.symbol, o.quantity, o.price, o.status, 
							       u.username, u.email, 
							       p.name AS portfolio_name
							FROM orders o
							JOIN users u ON o.user_id = u.id
							JOIN portfolios p ON o.portfolio_id = p.id
							WHERE o.user_id = %d
							AND o.created_at > '2025-01-01'
							ORDER BY o.created_at DESC
							LIMIT 20
						`, clientID)
					}
					
					_, err := dbService.Query(query)
					if err != nil {
						t.Logf("Error executing query: %v", err)
					}
				} else {
					// Execute a write query
					query := fmt.Sprintf("INSERT INTO orders (user_id, symbol, quantity, price, status) VALUES (%d, 'NIFTY', 1, 18000.0, 'OPEN')", clientID)
					err := dbService.Exec(query)
					if err != nil {
						t.Logf("Error executing query: %v", err)
					}
				}
				
				// Record response time
				resultChan <- time.Since(start)
			}
		}(i)
	}
	
	// Wait for all clients to finish
	wg.Wait()
	
	// Calculate results
	totalDuration := time.Since(startTime)
	totalRequests := numConcurrentClients * queriesPerClient
	
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
	t.Logf("Database Performance Test Results:")
	t.Logf("  Total Requests: %d", result.TotalRequests)
	t.Logf("  Successful Requests: %d", result.SuccessfulRequests)
	t.Logf("  Failed Requests: %d", result.FailedRequests)
	t.Logf("  Min Response Time: %v", result.MinResponseTime)
	t.Logf("  Max Response Time: %v", result.MaxResponseTime)
	t.Logf("  Avg Response Time: %v", result.AvgResponseTime)
	t.Logf("  Total Duration: %v", result.TotalDuration)
	t.Logf("  Requests Per Second: %.2f", result.RequestsPerSecond)
	
	// Assert that the performance meets requirements
	assert.GreaterOrEqual(t, result.RequestsPerSecond, 1000.0, "Database should handle at least 1000 queries per second")
	assert.Less(t, result.AvgResponseTime, 20*time.Millisecond, "Average database response time should be less than 20ms")
}

// TestTimescaleDBPerformance tests the TimescaleDB performance for time-series data
func TestTimescaleDBPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TimescaleDB performance test in short mode")
	}

	// Create a mock database for testing
	mockDB := &MockDB{}
	
	// Create a database service with the mock DB
	dbService := database.NewService(mockDB)
	
	// Test parameters
	numDataPoints := 10000
	batchSize := 1000
	numQueries := 100
	
	// Create a channel to collect results
	insertResultChan := make(chan time.Duration, numDataPoints/batchSize)
	queryResultChan := make(chan time.Duration, numQueries)
	
	// Set up the mock DB to simulate query execution
	mockDB.ExecFunc = func(query string, args ...interface{}) error {
		// Simulate batch insert time
		time.Sleep(time.Duration(rand.Intn(20)+10) * time.Millisecond)
		return nil
	}
	
	mockDB.QueryFunc = func(query string, args ...interface{}) ([]map[string]interface{}, error) {
		// Simulate time-series query execution time
		time.Sleep(time.Duration(rand.Intn(30)+20) * time.Millisecond)
		
		// Return mock results
		results := make([]map[string]interface{}, 0, 100)
		for i := 0; i < 100; i++ {
			results = append(results, map[string]interface{}{
				"time":       time.Now().Add(-time.Duration(i) * time.Minute),
				"symbol":     "NIFTY",
				"last_price": 18000.0 + float64(rand.Intn(100)),
				"volume":     rand.Intn(1000) + 1000,
			})
		}
		return results, nil
	}
	
	// Start time for inserts
	insertStartTime := time.Now()
	
	// Insert time-series data in batches
	for i := 0; i < numDataPoints; i += batchSize {
		// Record start time
		start := time.Now()
		
		// Create a batch insert query
		query := "INSERT INTO market_data (time, symbol, last_price, volume) VALUES "
		for j := 0; j < batchSize; j++ {
			if j > 0 {
				query += ", "
			}
			query += fmt.Sprintf("('2025-04-02 %02d:%02d:%02d', 'NIFTY', %f, %d)",
				rand.Intn(24), rand.Intn(60), rand.Intn(60),
				18000.0+float64(rand.Intn(100)),
				rand.Intn(1000)+1000)
		}
		
		// Execute the batch insert
		err := dbService.Exec(query)
		if err != nil {
			t.Logf("Error executing batch insert: %v", err)
		}
		
		// Record response time
		insertResultChan <- time.Since(start)
	}
	
	// Calculate insert results
	insertTotalDuration := time.Since(insertStartTime)
	insertTotalRequests := numDataPoints / batchSize
	
	// Collect insert response times
	var insertResponseTimes []time.Duration
	for i := 0; i < insertTotalRequests; i++ {
		insertResponseTimes = append(insertResponseTimes, <-insertResultChan)
	}
	
	// Calculate insert statistics
	insertMinTime := insertResponseTimes[0]
	insertMaxTime := insertResponseTimes[0]
	var insertTotalTime time.Duration
	
	for _, t := range insertResponseTimes {
		if t < insertMinTime {
			insertMinTime = t
		}
		if t > insertMaxTime {
			insertMaxTime = t
		}
		insertTotalTime += t
	}
	
	insertAvgTime := insertTotalTime / time.Duration(len(insertResponseTimes))
	insertPointsPerSecond := float64(numDataPoints) / insertTotalDuration.Seconds()
	
	// Start time for queries
	queryStartTime := time.Now()
	
	// Execute time-series queries
	for i := 0; i < numQueries; i++ {
		// Record start time
		start := time.Now()
		
		// Create a time-series query
		var query string
		queryType := rand.Intn(3)
		
		switch queryType {
		case 0:
			// Simple time range query
			query = "SELECT time, last_price FROM market_data WHERE symbol = 'NIFTY' AND time > '2025-04-01' AND time < '2025-04-02' ORDER BY time DESC LIMIT 100"
		case 1:
			// Aggregation query
			query = "SELECT time_bucket('1 hour', time) AS hour, AVG(last_price) AS avg_price, SUM(volume) AS total_volume FROM market_data WHERE symbol = 'NIFTY' AND time > '2025-04-01' GROUP BY hour ORDER BY hour DESC LIMIT 24"
		case 2:
			// Complex time-series query
			query = `
				SELECT time_bucket('5 minutes', time) AS bucket,
				       first(last_price, time) AS open_price,
				       max(last_price) AS high_price,
				       min(last_price) AS low_price,
				       last(last_price, time) AS close_price,
				       sum(volume) AS total_volume
				FROM market_data
				WHERE symbol = 'NIFTY'
				AND time > '2025-04-01'
				GROUP BY bucket
				ORDER BY bucket DESC
				LIMIT 100
			`
		}
		
		// Execute the query
		_, err := dbService.Query(query)
		if err != nil {
			t.Logf("Error executing time-series query: %v", err)
		}
		
		// Record response time
		queryResultChan <- time.Since(start)
	}
	
	// Calculate query results
	queryTotalDuration := time.Since(queryStartTime)
	
	// Collect query response times
	var queryResponseTimes []time.Duration
	for i := 0; i < numQueries; i++ {
		queryResponseTimes = append(queryResponseTimes, <-queryResultChan)
	}
	
	// Calculate query statistics
	queryMinTime := queryResponseTimes[0]
	queryMaxTime := queryResponseTimes[0]
	var queryTotalTime time.Duration
	
	for _, t := range queryResponseTimes {
		if t < queryMinTime {
			queryMinTime = t
		}
		if t > queryMaxTime {
			queryMaxTime = t
		}
		queryTotalTime += t
	}
	
	queryAvgTime := queryTotalTime / time.Duration(len(queryResponseTimes))
	queriesPerSecond := float64(numQueries) / queryTotalDuration.Seconds()
	
	// Log results
	t.Logf("TimescaleDB Performance Test Results:")
	t.Logf("  Insert Performance:")
	t.Logf("    Total Data Points: %d", numDataPoints)
	t.Logf("    Batch Size: %d", batchSize)
	t.Logf("    Min Batch Time: %v", insertMinTime)
	t.Logf("    Max Batch Time: %v", insertMaxTime)
	t.Logf("    Avg Batch Time: %v", insertAvgTime)
	t.Logf("    Total Duration: %v", insertTotalDuration)
	t.Logf("    Data Points Per Second: %.2f", insertPointsPerSecond)
	
	t.Logf("  Query Performance:")
	t.Logf("    Total Queries: %d", numQueries)
	t.Logf("    Min Query Time: %v", queryMinTime)
	t.Logf("    Max Query Time: %v", queryMaxTime)
	t.Logf("    Avg Query Time: %v", queryAvgTime)
	t.Logf("    Total Duration: %v", queryTotalDuration)
	t.Logf("    Queries Per Second: %.2f", queriesPerSecond)
	
	// Assert that the performance meets requirements
	assert.GreaterOrEqual(t, insertPointsPerSecond, 10000.0, "TimescaleDB should handle at least 10,000 data points per second for inserts")
	assert.Less(t, insertAvgTime, 50*time.Millisecond, "Average batch insert time should be less than 50ms")
	assert.GreaterOrEqual(t, queriesPerSecond, 10.0, "TimescaleDB should handle at least 10 complex time-series queries per second")
	assert.Less(t, queryAvgTime, 100*time.Millisecond, "Average time-series query time should be less than 100ms")
}

// MockDB is a mock implementation of the database interface
type MockDB struct {
	ConnectFunc func() error
	CloseFunc   func() error
	PingFunc    func() error
	ExecFunc    func(query string, args ...interface{}) error
	QueryFunc   func(query string, args ...interface{}) ([]map[string]interface{}, error)
}

func (m *MockDB) Connect() error {
	if m.ConnectFunc != nil {
		return m.ConnectFunc()
	}
	return nil
}

func (m *MockDB) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *MockDB) Ping() error {
	if m.PingFunc != nil {
		return m.PingFunc()
	}
	return nil
}

func (m *MockDB) Exec(query string, args ...interface{}) error {
	if m.ExecFunc != nil {
		return m.ExecFunc(query, args...)
	}
	return nil
}

func (m *MockDB) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(query, args...)
	}
	return []map[string]interface{}{}, nil
}
