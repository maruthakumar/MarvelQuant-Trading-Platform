package xts_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trade-execution-platform/backend/internal/xts/config"
	"github.com/trade-execution-platform/backend/internal/xts/models"
	"github.com/trade-execution-platform/backend/internal/xts/rest"
)

// MockHTTPClient is a mock HTTP client for testing
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestLogin(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockHTTPClient)
	
	// Create a test config
	cfg := config.NewXTSConfig()
	cfg.BaseURL = "https://test-api.xts.com"
	cfg.APIKey = "test-api-key"
	cfg.SecretKey = "test-secret-key"
	
	// Create a test response
	responseBody := `{
		"type": "success",
		"code": "200",
		"description": "Login successful",
		"result": {
			"token": "test-token",
			"userID": "test-user",
			"isInvestorClient": true,
			"exchangeSegments": ["NSECM", "NSEFO"]
		}
	}`
	
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "/interactive/user/session", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		
		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}))
	defer server.Close()
	
	// Update config to use test server
	cfg.BaseURL = server.URL
	
	// Create client
	client, err := rest.NewClient(cfg)
	assert.NoError(t, err)
	
	// Test login
	session, err := client.Login()
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "test-token", session.Token)
	assert.Equal(t, "test-user", session.UserID)
	assert.True(t, session.IsInvestorClient)
}

func TestPlaceOrder(t *testing.T) {
	// Create a test config
	cfg := config.NewXTSConfig()
	cfg.BaseURL = "https://test-api.xts.com"
	cfg.APIKey = "test-api-key"
	cfg.SecretKey = "test-secret-key"
	
	// Create a test response
	responseBody := `{
		"type": "success",
		"code": "200",
		"description": "Order placed successfully",
		"result": {
			"orderID": "123456",
			"exchangeOrderID": "XYZ789",
			"status": "PENDING"
		}
	}`
	
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "/interactive/orders", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		
		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}))
	defer server.Close()
	
	// Update config to use test server
	cfg.BaseURL = server.URL
	
	// Create client
	client, err := rest.NewClient(cfg)
	assert.NoError(t, err)
	
	// Set token
	client.SetToken("test-token", "test-user")
	
	// Create test order
	order := &models.Order{
		ExchangeSegment:      "NSECM",
		ExchangeInstrumentID: "RELIANCE",
		ProductType:          models.ProductNRML,
		OrderType:            models.OrderTypeLimit,
		OrderSide:            models.TransactionTypeBuy,
		TimeInForce:          models.ValidityDay,
		DisclosedQuantity:    0,
		OrderQuantity:        100,
		LimitPrice:          2500.0,
		StopPrice:           0.0,
		OrderUniqueIdentifier: "test-order-123",
	}
	
	// Test place order
	response, err := client.PlaceOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "123456", response.OrderID)
	assert.Equal(t, "XYZ789", response.ExchangeOrderID)
	assert.Equal(t, "PENDING", response.Status)
}

func TestGetOrderBook(t *testing.T) {
	// Create a test config
	cfg := config.NewXTSConfig()
	cfg.BaseURL = "https://test-api.xts.com"
	cfg.APIKey = "test-api-key"
	cfg.SecretKey = "test-secret-key"
	
	// Create a test response
	responseBody := `{
		"type": "success",
		"code": "200",
		"description": "Order book retrieved successfully",
		"result": {
			"orders": [
				{
					"appOrderID": 123456,
					"orderSide": "BUY",
					"orderType": "LIMIT",
					"productType": "NRML",
					"timeInForce": "DAY",
					"orderQuantity": 100,
					"orderPrice": 2500.0,
					"orderStopPrice": 0.0,
					"orderStatus": "PENDING",
					"orderAverageTradedPrice": 0.0,
					"orderDisclosedQuantity": 0,
					"orderTradedQuantity": 0,
					"exchangeOrderID": "XYZ789",
					"exchangeSegment": "NSECM",
					"exchangeInstrumentID": "RELIANCE",
					"orderGeneratedDateTime": "2025-04-02T07:00:00Z",
					"lastUpdateDateTime": "2025-04-02T07:00:00Z",
					"cancelRejectReason": "",
					"orderUniqueIdentifier": "test-order-123"
				}
			]
		}
	}`
	
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "/interactive/orders", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		
		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}))
	defer server.Close()
	
	// Update config to use test server
	cfg.BaseURL = server.URL
	
	// Create client
	client, err := rest.NewClient(cfg)
	assert.NoError(t, err)
	
	// Set token
	client.SetToken("test-token", "test-user")
	
	// Test get order book
	orderBook, err := client.GetOrderBook("")
	assert.NoError(t, err)
	assert.NotNil(t, orderBook)
	assert.Len(t, orderBook.Orders, 1)
	assert.Equal(t, 123456, orderBook.Orders[0].AppOrderID)
	assert.Equal(t, "BUY", orderBook.Orders[0].OrderSide)
	assert.Equal(t, "LIMIT", orderBook.Orders[0].OrderType)
}

func TestRetryMechanism(t *testing.T) {
	// Create a test config
	cfg := config.NewXTSConfig()
	cfg.BaseURL = "https://test-api.xts.com"
	cfg.APIKey = "test-api-key"
	cfg.SecretKey = "test-secret-key"
	
	// Create a counter for retry attempts
	attempts := 0
	
	// Create a test server that fails first, then succeeds
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		
		if attempts <= 2 {
			// Simulate a server error for the first two attempts
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		
		// Succeed on the third attempt
		responseBody := `{
			"type": "success",
			"code": "200",
			"description": "Login successful",
			"result": {
				"token": "test-token",
				"userID": "test-user",
				"isInvestorClient": true,
				"exchangeSegments": ["NSECM", "NSEFO"]
			}
		}`
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}))
	defer server.Close()
	
	// Update config to use test server
	cfg.BaseURL = server.URL
	
	// Create client with retry config
	client, err := rest.NewClient(cfg)
	assert.NoError(t, err)
	
	// Set up retry with context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Test login with retry
	var session *models.Session
	err = recovery.RetryWithBackoff(ctx, recovery.DefaultRetryConfig(), func() error {
		var loginErr error
		session, loginErr = client.Login()
		return loginErr
	})
	
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "test-token", session.Token)
	assert.Equal(t, 3, attempts) // Verify it took 3 attempts
}
