package broker

import (
	"context"
	"testing"
	"time"
	"os"
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trading-platform/backend/internal/broker"
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

func TestBrokerFactory(t *testing.T) {
	// Create a broker factory
	brokerFactory := broker.NewBrokerFactory()

	// Test registering and getting brokers
	t.Run("RegisterAndGetBroker", func(t *testing.T) {
		// Register a mock broker
		mockBroker := &MockBroker{}
		brokerFactory.RegisterBroker("xts", mockBroker)

		// Get the registered broker
		retrievedBroker, err := brokerFactory.GetBroker("xts")
		require.NoError(t, err)
		assert.Equal(t, mockBroker, retrievedBroker)

		// Try to get a non-existent broker
		_, err = brokerFactory.GetBroker("nonexistent")
		assert.Error(t, err)
	})
}

func TestXTSBrokerIntegration(t *testing.T) {
	// Skip this test in short mode or if no credentials are available
	if testing.Short() {
		t.Skip("Skipping broker integration test in short mode")
	}

	// Load configuration
	config, err := loadConfig()
	if err != nil || config.Broker.XTS.APIKey == "test-api-key" {
		t.Skip("Skipping broker integration test due to missing configuration")
	}

	// Create an XTS broker
	xtsBroker := broker.NewXTSBroker()

	// Initialize the broker with configuration
	brokerConfig := broker.BrokerConfig{
		Endpoint:  config.Broker.XTS.Endpoint,
		APIKey:    config.Broker.XTS.APIKey,
		APISecret: config.Broker.XTS.APISecret,
	}
	err = xtsBroker.Initialize(brokerConfig)
	require.NoError(t, err)

	// Test the broker functionality
	ctx := context.Background()

	// Test getting a quote
	t.Run("GetQuote", func(t *testing.T) {
		quote, err := xtsBroker.GetQuote(ctx, "NIFTY", "NSE")
		require.NoError(t, err)
		require.NotNil(t, quote)
		assert.Equal(t, "NIFTY", quote.Symbol)
		assert.Equal(t, "NSE", quote.Exchange)
		assert.Greater(t, quote.LastPrice, 0.0)
	})

	// Test subscribing to quotes
	t.Run("SubscribeQuotes", func(t *testing.T) {
		err := xtsBroker.SubscribeQuotes(ctx, []string{"NIFTY", "BANKNIFTY"}, "NSE")
		require.NoError(t, err)

		// Wait a bit to ensure subscription is processed
		time.Sleep(2 * time.Second)

		// Unsubscribe
		err = xtsBroker.UnsubscribeQuotes(ctx, []string{"NIFTY", "BANKNIFTY"}, "NSE")
		require.NoError(t, err)
	})

	// Test placing, modifying, and canceling an order
	// Note: This is commented out to avoid actual order placement in tests
	// Uncomment and modify for actual broker testing with real credentials
	/*
	t.Run("OrderLifecycle", func(t *testing.T) {
		// Place an order
		orderRequest := broker.OrderRequest{
			Symbol:          "NIFTY",
			Exchange:        "NSE",
			OrderType:       broker.OrderTypeLimit,
			TransactionType: broker.TransactionTypeBuy,
			ProductType:     broker.ProductTypeNRML,
			Quantity:        1,
			Price:           17000.0, // Set a price that won't execute
		}

		placeResponse, err := xtsBroker.PlaceOrder(ctx, orderRequest)
		require.NoError(t, err)
		require.NotNil(t, placeResponse)
		assert.True(t, placeResponse.Success)
		assert.NotEmpty(t, placeResponse.OrderID)

		// Get the order
		orderID := placeResponse.OrderID
		order, err := xtsBroker.GetOrder(ctx, orderID)
		require.NoError(t, err)
		require.NotNil(t, order)
		assert.Equal(t, orderID, order.ID)

		// Modify the order
		modifyRequest := orderRequest
		modifyRequest.Price = 17100.0
		modifyResponse, err := xtsBroker.ModifyOrder(ctx, orderID, modifyRequest)
		require.NoError(t, err)
		require.NotNil(t, modifyResponse)
		assert.True(t, modifyResponse.Success)

		// Cancel the order
		cancelResponse, err := xtsBroker.CancelOrder(ctx, orderID)
		require.NoError(t, err)
		require.NotNil(t, cancelResponse)
		assert.True(t, cancelResponse.Success)
	})
	*/

	// Test getting positions
	t.Run("GetPositions", func(t *testing.T) {
		positions, err := xtsBroker.GetPositions(ctx)
		require.NoError(t, err)
		// We can't assert much about the positions, as they depend on the account state
		// Just verify that we got a response
		assert.NotNil(t, positions)
	})

	// Clean up
	err = xtsBroker.Close()
	require.NoError(t, err)
}

func TestMockBrokerImplementation(t *testing.T) {
	// Create a mock broker
	mockBroker := &MockBroker{}

	// Initialize the broker
	err := mockBroker.Initialize(broker.BrokerConfig{})
	require.NoError(t, err)

	// Test the broker functionality
	ctx := context.Background()

	// Test placing an order
	t.Run("PlaceOrder", func(t *testing.T) {
		// Set up the mock to return a specific response
		mockBroker.PlaceOrderFunc = func(ctx context.Context, request broker.OrderRequest) (*broker.OrderResponse, error) {
			return &broker.OrderResponse{
				Success: true,
				OrderID: "test-order-123",
			}, nil
		}

		// Place an order
		orderRequest := broker.OrderRequest{
			Symbol:          "NIFTY",
			Exchange:        "NSE",
			OrderType:       broker.OrderTypeMarket,
			TransactionType: broker.TransactionTypeBuy,
			ProductType:     broker.ProductTypeNRML,
			Quantity:        1,
		}

		response, err := mockBroker.PlaceOrder(ctx, orderRequest)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.True(t, response.Success)
		assert.Equal(t, "test-order-123", response.OrderID)
	})

	// Test getting an order
	t.Run("GetOrder", func(t *testing.T) {
		// Set up the mock to return a specific order
		mockBroker.GetOrderFunc = func(ctx context.Context, orderID string) (*broker.Order, error) {
			return &broker.Order{
				ID:              orderID,
				Symbol:          "NIFTY",
				Exchange:        "NSE",
				OrderType:       broker.OrderTypeMarket,
				TransactionType: broker.TransactionTypeBuy,
				ProductType:     broker.ProductTypeNRML,
				Quantity:        1,
				Status:          broker.OrderStatusCompleted,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}, nil
		}

		// Get an order
		order, err := mockBroker.GetOrder(ctx, "test-order-123")
		require.NoError(t, err)
		require.NotNil(t, order)
		assert.Equal(t, "test-order-123", order.ID)
		assert.Equal(t, "NIFTY", order.Symbol)
		assert.Equal(t, broker.OrderStatusCompleted, order.Status)
	})

	// Test getting a quote
	t.Run("GetQuote", func(t *testing.T) {
		// Set up the mock to return a specific quote
		mockBroker.GetQuoteFunc = func(ctx context.Context, symbol, exchange string) (*broker.Quote, error) {
			return &broker.Quote{
				Symbol:    symbol,
				Exchange:  exchange,
				LastPrice: 18500.0,
				BidPrice:  18499.0,
				AskPrice:  18501.0,
				Volume:    1000000,
				Timestamp: time.Now(),
			}, nil
		}

		// Get a quote
		quote, err := mockBroker.GetQuote(ctx, "NIFTY", "NSE")
		require.NoError(t, err)
		require.NotNil(t, quote)
		assert.Equal(t, "NIFTY", quote.Symbol)
		assert.Equal(t, "NSE", quote.Exchange)
		assert.Equal(t, 18500.0, quote.LastPrice)
	})

	// Clean up
	err = mockBroker.Close()
	require.NoError(t, err)
}

// MockBroker is a mock implementation of the Broker interface for testing
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
