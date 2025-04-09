package integration

import (
	"context"
	"testing"
	"time"
	"os"
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trading-platform/backend/internal/auth"
	"github.com/trading-platform/backend/internal/database"
	"github.com/trading-platform/backend/internal/messagequeue"
	"github.com/trading-platform/backend/internal/core"
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

// setupTestEnvironment sets up the test environment
func setupTestEnvironment(t *testing.T) (*database.Service, *messagequeue.Service, *auth.AuthService, func()) {
	// Load the configuration
	config, err := loadConfig()
	require.NoError(t, err)

	// Create a database service
	dbService := database.NewService(nil) // In a real test, this would connect to a test database

	// Create a message queue service
	mqService := messagequeue.NewService(nil, nil) // In a real test, this would connect to test Redis and RabbitMQ

	// Create an auth service
	authService := auth.NewAuthService(config.JWT.Secret)

	// Return the services and a cleanup function
	cleanup := func() {
		// Clean up resources
		dbService.Close()
		mqService.Close()
	}

	return dbService, mqService, authService, cleanup
}

func TestAuthenticationFlow(t *testing.T) {
	// Set up the test environment
	dbService, mqService, authService, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test the complete authentication flow
	t.Run("RegisterLoginVerify", func(t *testing.T) {
		ctx := context.Background()

		// 1. Register a new user
		registerRequest := auth.RegisterRequest{
			Username: "integrationtest",
			Email:    "integration@test.com",
			Password: "password123",
		}

		registerResponse, err := authService.Register(ctx, registerRequest)
		require.NoError(t, err)
		require.NotNil(t, registerResponse)
		assert.NotEmpty(t, registerResponse.Token)
		assert.Equal(t, registerRequest.Username, registerResponse.User.Username)
		assert.Equal(t, registerRequest.Email, registerResponse.User.Email)

		// 2. Login with the registered user
		loginRequest := auth.LoginRequest{
			Username: "integrationtest",
			Password: "password123",
		}

		loginResponse, err := authService.Login(ctx, loginRequest)
		require.NoError(t, err)
		require.NotNil(t, loginResponse)
		assert.NotEmpty(t, loginResponse.Token)
		assert.Equal(t, registerRequest.Username, loginResponse.User.Username)

		// 3. Verify the token
		claims, err := authService.VerifyToken(loginResponse.Token)
		require.NoError(t, err)
		require.NotNil(t, claims)
		assert.Equal(t, loginResponse.User.ID, claims.UserID)
		assert.Equal(t, loginResponse.User.Username, claims.Username)
	})
}

func TestOrderExecutionFlow(t *testing.T) {
	// Set up the test environment
	dbService, mqService, authService, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a broker factory
	brokerFactory := broker.NewBrokerFactory()

	// Register a mock broker
	mockBroker := &MockBroker{}
	brokerFactory.RegisterBroker("xts", mockBroker)

	// Create an order service
	orderService := core.NewOrderService(brokerFactory)

	// Create an execution engine
	executionEngine := core.NewExecutionEngine(orderService)

	// Test the complete order execution flow
	t.Run("PlaceOrderExecuteMonitor", func(t *testing.T) {
		ctx := context.Background()

		// 1. Place an order
		orderRequest := core.OrderRequest{
			UserID:          "user123",
			BrokerName:      "xts",
			Symbol:          "NIFTY",
			Exchange:        "NSE",
			OrderType:       "MARKET",
			TransactionType: "BUY",
			ProductType:     "NRML",
			Quantity:        1,
		}

		// Set up the mock broker to return a successful response
		mockBroker.PlaceOrderFunc = func(ctx context.Context, request broker.OrderRequest) (*broker.OrderResponse, error) {
			return &broker.OrderResponse{
				Success: true,
				OrderID: "order123",
			}, nil
		}

		orderResponse, err := orderService.PlaceOrder(ctx, orderRequest)
		require.NoError(t, err)
		require.NotNil(t, orderResponse)
		assert.True(t, orderResponse.Success)
		assert.Equal(t, "order123", orderResponse.OrderID)

		// 2. Execute a strategy
		executionRequest := core.ExecutionRequest{
			UserID:      "user123",
			PortfolioID: "portfolio123",
			StrategyID:  "strategy123",
		}

		executionResponse, err := executionEngine.ExecuteStrategy(ctx, executionRequest)
		require.NoError(t, err)
		require.NotNil(t, executionResponse)
		assert.True(t, executionResponse.Success)
		assert.NotEmpty(t, executionResponse.OrderIDs)

		// 3. Monitor execution
		err = executionEngine.MonitorExecution(ctx, "portfolio123")
		require.NoError(t, err)
	})
}

// MockBroker is a mock implementation of the Broker interface for integration testing
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
