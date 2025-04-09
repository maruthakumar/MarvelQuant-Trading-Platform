package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/trading-platform/backend/internal/broker"
	"github.com/trading-platform/backend/internal/core"
)

// MockBroker is a mock implementation of the Broker interface
type MockBroker struct {
	mock.Mock
}

func (m *MockBroker) Initialize(config broker.BrokerConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockBroker) PlaceOrder(ctx context.Context, request broker.OrderRequest) (*broker.OrderResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*broker.OrderResponse), args.Error(1)
}

func (m *MockBroker) ModifyOrder(ctx context.Context, orderID string, request broker.OrderRequest) (*broker.OrderResponse, error) {
	args := m.Called(ctx, orderID, request)
	return args.Get(0).(*broker.OrderResponse), args.Error(1)
}

func (m *MockBroker) CancelOrder(ctx context.Context, orderID string) (*broker.OrderResponse, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(*broker.OrderResponse), args.Error(1)
}

func (m *MockBroker) GetOrder(ctx context.Context, orderID string) (*broker.Order, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(*broker.Order), args.Error(1)
}

func (m *MockBroker) GetOrders(ctx context.Context) ([]broker.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]broker.Order), args.Error(1)
}

func (m *MockBroker) GetPositions(ctx context.Context) ([]broker.Position, error) {
	args := m.Called(ctx)
	return args.Get(0).([]broker.Position), args.Error(1)
}

func (m *MockBroker) GetQuote(ctx context.Context, symbol, exchange string) (*broker.Quote, error) {
	args := m.Called(ctx, symbol, exchange)
	return args.Get(0).(*broker.Quote), args.Error(1)
}

func (m *MockBroker) SubscribeQuotes(ctx context.Context, symbols []string, exchange string) error {
	args := m.Called(ctx, symbols, exchange)
	return args.Error(0)
}

func (m *MockBroker) UnsubscribeQuotes(ctx context.Context, symbols []string, exchange string) error {
	args := m.Called(ctx, symbols, exchange)
	return args.Error(0)
}

func (m *MockBroker) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockBrokerFactory is a mock implementation of the BrokerFactory
type MockBrokerFactory struct {
	mock.Mock
}

func (m *MockBrokerFactory) RegisterBroker(name string, broker broker.Broker) {
	m.Called(name, broker)
}

func (m *MockBrokerFactory) GetBroker(name string) (broker.Broker, error) {
	args := m.Called(name)
	return args.Get(0).(broker.Broker), args.Error(1)
}

func TestOrderService(t *testing.T) {
	// Create a mock broker factory
	mockFactory := new(MockBrokerFactory)
	
	// Create a mock broker
	mockBroker := new(MockBroker)
	
	// Create an order service with the mock factory
	orderService := core.NewOrderService(mockFactory)
	
	// Test placing an order
	t.Run("PlaceOrder", func(t *testing.T) {
		ctx := context.Background()
		
		// Create an order request
		request := core.OrderRequest{
			UserID:          "user123",
			BrokerName:      "xts",
			Symbol:          "NIFTY",
			Exchange:        "NSE",
			OrderType:       broker.OrderTypeMarket,
			TransactionType: broker.TransactionTypeBuy,
			ProductType:     broker.ProductTypeNRML,
			Quantity:        1,
		}
		
		// Set up the mock broker factory to return our mock broker
		mockFactory.On("GetBroker", "xts").Return(mockBroker, nil)
		
		// Set up the mock broker to return a successful response
		expectedResponse := &broker.OrderResponse{
			Success: true,
			OrderID: "order123",
		}
		mockBroker.On("PlaceOrder", ctx, mock.Anything).Return(expectedResponse, nil)
		
		// Place the order
		response, err := orderService.PlaceOrder(ctx, request)
		
		// Check that there was no error
		require.NoError(t, err)
		
		// Check that the response is not nil
		require.NotNil(t, response)
		
		// Check that the response is as expected
		assert.Equal(t, expectedResponse.Success, response.Success)
		assert.Equal(t, expectedResponse.OrderID, response.OrderID)
		
		// Verify that the mock methods were called as expected
		mockFactory.AssertExpectations(t)
		mockBroker.AssertExpectations(t)
	})
	
	// Test modifying an order
	t.Run("ModifyOrder", func(t *testing.T) {
		ctx := context.Background()
		
		// Create an order request
		request := core.OrderRequest{
			UserID:          "user123",
			BrokerName:      "xts",
			Symbol:          "NIFTY",
			Exchange:        "NSE",
			OrderType:       broker.OrderTypeLimit,
			TransactionType: broker.TransactionTypeBuy,
			ProductType:     broker.ProductTypeNRML,
			Quantity:        2,
			Price:           18000.0,
		}
		
		// Set up the mock broker factory to return our mock broker
		mockFactory.On("GetBroker", "xts").Return(mockBroker, nil)
		
		// Set up the mock broker to return a successful response
		expectedResponse := &broker.OrderResponse{
			Success: true,
			OrderID: "order123",
		}
		mockBroker.On("ModifyOrder", ctx, "order123", mock.Anything).Return(expectedResponse, nil)
		
		// Modify the order
		response, err := orderService.ModifyOrder(ctx, "order123", request)
		
		// Check that there was no error
		require.NoError(t, err)
		
		// Check that the response is not nil
		require.NotNil(t, response)
		
		// Check that the response is as expected
		assert.Equal(t, expectedResponse.Success, response.Success)
		assert.Equal(t, expectedResponse.OrderID, response.OrderID)
		
		// Verify that the mock methods were called as expected
		mockFactory.AssertExpectations(t)
		mockBroker.AssertExpectations(t)
	})
	
	// Test cancelling an order
	t.Run("CancelOrder", func(t *testing.T) {
		ctx := context.Background()
		
		// Set up the mock broker factory to return our mock broker
		mockFactory.On("GetBroker", "xts").Return(mockBroker, nil)
		
		// Set up the mock broker to return a successful response
		expectedResponse := &broker.OrderResponse{
			Success: true,
			OrderID: "order123",
		}
		mockBroker.On("CancelOrder", ctx, "order123").Return(expectedResponse, nil)
		
		// Cancel the order
		response, err := orderService.CancelOrder(ctx, "order123", "xts")
		
		// Check that there was no error
		require.NoError(t, err)
		
		// Check that the response is not nil
		require.NotNil(t, response)
		
		// Check that the response is as expected
		assert.Equal(t, expectedResponse.Success, response.Success)
		assert.Equal(t, expectedResponse.OrderID, response.OrderID)
		
		// Verify that the mock methods were called as expected
		mockFactory.AssertExpectations(t)
		mockBroker.AssertExpectations(t)
	})
	
	// Test getting an order
	t.Run("GetOrder", func(t *testing.T) {
		ctx := context.Background()
		
		// Set up the mock broker factory to return our mock broker
		mockFactory.On("GetBroker", "xts").Return(mockBroker, nil)
		
		// Set up the mock broker to return an order
		expectedOrder := &broker.Order{
			ID:              "order123",
			BrokerOrderID:   "broker-order-123",
			Symbol:          "NIFTY",
			Exchange:        "NSE",
			OrderType:       broker.OrderTypeMarket,
			TransactionType: broker.TransactionTypeBuy,
			ProductType:     broker.ProductTypeNRML,
			Quantity:        1,
			Status:          broker.OrderStatusCompleted,
			OrderTimestamp:  time.Now(),
			UpdatedAt:       time.Now(),
		}
		mockBroker.On("GetOrder", ctx, "order123").Return(expectedOrder, nil)
		
		// Get the order
		order, err := orderService.GetOrder(ctx, "order123", "xts")
		
		// Check that there was no error
		require.NoError(t, err)
		
		// Check that the order is not nil
		require.NotNil(t, order)
		
		// Check that the order is as expected
		assert.Equal(t, expectedOrder.ID, order.ID)
		assert.Equal(t, expectedOrder.Symbol, order.Symbol)
		assert.Equal(t, expectedOrder.OrderType, order.OrderType)
		assert.Equal(t, expectedOrder.Status, order.Status)
		
		// Verify that the mock methods were called as expected
		mockFactory.AssertExpectations(t)
		mockBroker.AssertExpectations(t)
	})
	
	// Test getting all orders
	t.Run("GetOrders", func(t *testing.T) {
		ctx := context.Background()
		
		// Set up the mock broker factory to return our mock broker
		mockFactory.On("GetBroker", "xts").Return(mockBroker, nil)
		
		// Set up the mock broker to return orders
		expectedOrders := []broker.Order{
			{
				ID:              "order123",
				BrokerOrderID:   "broker-order-123",
				Symbol:          "NIFTY",
				Exchange:        "NSE",
				OrderType:       broker.OrderTypeMarket,
				TransactionType: broker.TransactionTypeBuy,
				ProductType:     broker.ProductTypeNRML,
				Quantity:        1,
				Status:          broker.OrderStatusCompleted,
				OrderTimestamp:  time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				ID:              "order456",
				BrokerOrderID:   "broker-order-456",
				Symbol:          "BANKNIFTY",
				Exchange:        "NSE",
				OrderType:       broker.OrderTypeLimit,
				TransactionType: broker.TransactionTypeSell,
				ProductType:     broker.ProductTypeNRML,
				Quantity:        2,
				Price:           40000.0,
				Status:          broker.OrderStatusOpen,
				OrderTimestamp:  time.Now(),
				UpdatedAt:       time.Now(),
			},
		}
		mockBroker.On("GetOrders", ctx).Return(expectedOrders, nil)
		
		// Get the orders
		orders, err := orderService.GetOrders(ctx, "xts")
		
		// Check that there was no error
		require.NoError(t, err)
		
		// Check that the orders are not nil
		require.NotNil(t, orders)
		
		// Check that the orders are as expected
		assert.Equal(t, len(expectedOrders), len(orders))
		assert.Equal(t, expectedOrders[0].ID, orders[0].ID)
		assert.Equal(t, expectedOrders[1].ID, orders[1].ID)
		
		// Verify that the mock methods were called as expected
		mockFactory.AssertExpectations(t)
		mockBroker.AssertExpectations(t)
	})
}
