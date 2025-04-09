package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trading-platform/backend/internal/broker/common"
)

// TestBrokerManager tests the broker manager
func TestBrokerManager(t *testing.T) {
	// Create a new broker manager
	manager := NewBrokerManager()
	assert.NotNil(t, manager)
	
	// Test register broker
	config := &common.BrokerConfig{
		BrokerType: common.BrokerTypeXTSClient,
		XTSClient: &common.XTSClientConfig{
			APIKey:    "test_api_key",
			SecretKey: "test_secret_key",
			Source:    "WEBAPI",
		},
	}
	
	err := manager.RegisterBroker("client1", config)
	assert.NoError(t, err)
	
	// Test register broker with empty client ID
	err = manager.RegisterBroker("", config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client ID is required")
	
	// Test register broker with nil config
	err = manager.RegisterBroker("client2", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "broker configuration is required")
	
	// Test get client ID for user with empty user ID
	_, err = manager.GetClientIDForUser("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user ID is required")
	
	// Test get client ID for non-existent user
	_, err = manager.GetClientIDForUser("user1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active session found for user ID")
}

// TestBrokerManagerWithMockClient tests the broker manager with a mock client
func TestBrokerManagerWithMockClient(t *testing.T) {
	// Create a new broker manager
	manager := NewBrokerManager()
	
	// Register a mock broker
	config := &common.BrokerConfig{
		BrokerType: common.BrokerTypeXTSClient,
		XTSClient: &common.XTSClientConfig{
			APIKey:    "test_api_key",
			SecretKey: "test_secret_key",
			Source:    "WEBAPI",
		},
	}
	
	err := manager.RegisterBroker("client1", config)
	assert.NoError(t, err)
	
	// Create a mock client
	mockClient := &MockBrokerClient{
		LoginFunc: func(credentials *common.Credentials) (*common.Session, error) {
			return &common.Session{
				Token:     "test_token",
				UserID:    "user1",
				ExpiresAt: 1617345678,
			}, nil
		},
		LogoutFunc: func() error {
			return nil
		},
		PlaceOrderFunc: func(order *common.Order) (*common.OrderResponse, error) {
			return &common.OrderResponse{
				OrderID: "order1",
				Status:  "PLACED",
			}, nil
		},
		ModifyOrderFunc: func(order *common.ModifyOrder) (*common.OrderResponse, error) {
			return &common.OrderResponse{
				OrderID: order.OrderID,
				Status:  "MODIFIED",
			}, nil
		},
		CancelOrderFunc: func(orderID string, clientID string) (*common.OrderResponse, error) {
			return &common.OrderResponse{
				OrderID: orderID,
				Status:  "CANCELLED",
			}, nil
		},
		GetOrderBookFunc: func(clientID string) (*common.OrderBook, error) {
			return &common.OrderBook{
				Orders: []common.OrderDetails{
					{
						OrderID:       "order1",
						OrderSide:     "BUY",
						OrderType:     "LIMIT",
						OrderQuantity: 10,
						LimitPrice:    100.0,
						OrderStatus:   "OPEN",
					},
				},
			}, nil
		},
		GetPositionsFunc: func(clientID string) ([]common.Position, error) {
			return []common.Position{
				{
					ExchangeSegment:      "NSECM",
					ExchangeInstrumentID: "123456",
					ProductType:          "MIS",
					Quantity:             10,
					AveragePrice:         100.0,
					LastPrice:            101.0,
					UnrealizedProfit:     10.0,
				},
			}, nil
		},
		GetHoldingsFunc: func(clientID string) ([]common.Holding, error) {
			return []common.Holding{
				{
					ExchangeSegment:      "NSECM",
					ExchangeInstrumentID: "123456",
					TradingSymbol:        "RELIANCE",
					Quantity:             10,
					AveragePrice:         2000.0,
					LastPrice:            2050.0,
					UnrealizedProfit:     500.0,
				},
			}, nil
		},
		GetQuoteFunc: func(symbols []string) (map[string]common.Quote, error) {
			quotes := make(map[string]common.Quote)
			for _, symbol := range symbols {
				quotes[symbol] = common.Quote{
					ExchangeSegment:      "NSECM",
					ExchangeInstrumentID: "123456",
					TradingSymbol:        symbol,
					LastPrice:            2050.0,
					BidPrice:             2049.0,
					AskPrice:             2051.0,
				}
			}
			return quotes, nil
		},
	}
	
	// Override the client in the manager
	manager.clients["client1"] = mockClient
	
	// Test login
	session, err := manager.Login("client1", nil)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "test_token", session.Token)
	assert.Equal(t, "user1", session.UserID)
	
	// Test place order
	order := &common.Order{
		ExchangeSegment:       "NSECM",
		TradingSymbol:         "RELIANCE",
		OrderSide:             "BUY",
		OrderQuantity:         10,
		ProductType:           "MIS",
		OrderType:             "LIMIT",
		TimeInForce:           "DAY",
		LimitPrice:            2000.0,
		OrderUniqueIdentifier: "test123",
	}
	
	orderResponse, err := manager.PlaceOrder("user1", order)
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "order1", orderResponse.OrderID)
	assert.Equal(t, "PLACED", orderResponse.Status)
	
	// Test modify order
	modifyOrder := &common.ModifyOrder{
		OrderID:       "order1",
		OrderType:     "LIMIT",
		OrderQuantity: 20,
		LimitPrice:    2100.0,
	}
	
	orderResponse, err = manager.ModifyOrder("user1", modifyOrder)
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "order1", orderResponse.OrderID)
	assert.Equal(t, "MODIFIED", orderResponse.Status)
	
	// Test cancel order
	orderResponse, err = manager.CancelOrder("user1", "order1")
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "order1", orderResponse.OrderID)
	assert.Equal(t, "CANCELLED", orderResponse.Status)
	
	// Test get order book
	orderBook, err := manager.GetOrderBook("user1")
	assert.NoError(t, err)
	assert.NotNil(t, orderBook)
	assert.Len(t, orderBook.Orders, 1)
	assert.Equal(t, "order1", orderBook.Orders[0].OrderID)
	
	// Test get positions
	positions, err := manager.GetPositions("user1")
	assert.NoError(t, err)
	assert.NotNil(t, positions)
	assert.Len(t, positions, 1)
	assert.Equal(t, "NSECM", positions[0].ExchangeSegment)
	
	// Test get holdings
	holdings, err := manager.GetHoldings("user1")
	assert.NoError(t, err)
	assert.NotNil(t, holdings)
	assert.Len(t, holdings, 1)
	assert.Equal(t, "RELIANCE", holdings[0].TradingSymbol)
	
	// Test get quote
	quotes, err := manager.GetQuote("user1", []string{"RELIANCE"})
	assert.NoError(t, err)
	assert.NotNil(t, quotes)
	assert.Len(t, quotes, 1)
	assert.Equal(t, 2050.0, quotes["RELIANCE"].LastPrice)
	
	// Test logout
	err = manager.Logout("client1")
	assert.NoError(t, err)
}

// MockBrokerClient is a mock implementation of the BrokerClient interface for testing
type MockBrokerClient struct {
	LoginFunc                 func(credentials *common.Credentials) (*common.Session, error)
	LogoutFunc                func() error
	PlaceOrderFunc            func(order *common.Order) (*common.OrderResponse, error)
	ModifyOrderFunc           func(order *common.ModifyOrder) (*common.OrderResponse, error)
	CancelOrderFunc           func(orderID string, clientID string) (*common.OrderResponse, error)
	GetOrderBookFunc          func(clientID string) (*common.OrderBook, error)
	GetPositionsFunc          func(clientID string) ([]common.Position, error)
	GetHoldingsFunc           func(clientID string) ([]common.Holding, error)
	GetQuoteFunc              func(symbols []string) (map[string]common.Quote, error)
	SubscribeToQuotesFunc     func(symbols []string) (chan common.Quote, error)
	UnsubscribeFromQuotesFunc func(symbols []string) error
}

func (m *MockBrokerClient) Login(credentials *common.Credentials) (*common.Session, error) {
	return m.LoginFunc(credentials)
}

func (m *MockBrokerClient) Logout() error {
	return m.LogoutFunc()
}

func (m *MockBrokerClient) PlaceOrder(order *common.Order) (*common.OrderResponse, error) {
	return m.PlaceOrderFunc(order)
}

func (m *MockBrokerClient) ModifyOrder(order *common.ModifyOrder) (*common.OrderResponse, error) {
	return m.ModifyOrderFunc(order)
}

func (m *MockBrokerClient) CancelOrder(orderID string, clientID string) (*common.OrderResponse, error) {
	return m.CancelOrderFunc(orderID, clientID)
}

func (m *MockBrokerClient) GetOrderBook(clientID string) (*common.OrderBook, error) {
	return m.GetOrderBookFunc(clientID)
}

func (m *MockBrokerClient) GetPositions(clientID string) ([]common.Position, error) {
	return m.GetPositionsFunc(clientID)
}

func (m *MockBrokerClient) GetHoldings(clientID string) ([]common.Holding, error) {
	return m.GetHoldingsFunc(clientID)
}

func (m *MockBrokerClient) GetQuote(symbols []string) (map[string]common.Quote, error) {
	return m.GetQuoteFunc(symbols)
}

func (m *MockBrokerClient) SubscribeToQuotes(symbols []string) (chan common.Quote, error) {
	return m.SubscribeToQuotesFunc(symbols)
}

func (m *MockBrokerClient) UnsubscribeFromQuotes(symbols []string) error {
	return m.UnsubscribeFromQuotesFunc(symbols)
}
