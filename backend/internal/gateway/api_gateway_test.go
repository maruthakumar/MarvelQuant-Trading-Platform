package gateway

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	
	"trading_platform/backend/internal/interfaces"
	"trading_platform/backend/internal/models"
)

// MockExecutionPlatform is a mock implementation of the ExecutionPlatformInterface
type MockExecutionPlatform struct {
	mock.Mock
}

// GetRealTimeMarketData mocks the GetRealTimeMarketData method
func (m *MockExecutionPlatform) GetRealTimeMarketData(ctx context.Context, symbol string) (*models.MarketDataSnapshot, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MarketDataSnapshot), args.Error(1)
}

// GetHistoricalMarketData mocks the GetHistoricalMarketData method
func (m *MockExecutionPlatform) GetHistoricalMarketData(ctx context.Context, symbol string, startDate, endDate time.Time, timeframe string) ([]*models.MarketDataSnapshot, error) {
	args := m.Called(ctx, symbol, startDate, endDate, timeframe)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MarketDataSnapshot), args.Error(1)
}

// SubscribeToMarketData mocks the SubscribeToMarketData method
func (m *MockExecutionPlatform) SubscribeToMarketData(ctx context.Context, symbol string, callback func(*models.MarketDataSnapshot)) (string, error) {
	args := m.Called(ctx, symbol, callback)
	return args.String(0), args.Error(1)
}

// UnsubscribeFromMarketData mocks the UnsubscribeFromMarketData method
func (m *MockExecutionPlatform) UnsubscribeFromMarketData(ctx context.Context, subscriptionID string) error {
	args := m.Called(ctx, subscriptionID)
	return args.Error(0)
}

// GetInstrumentDetails mocks the GetInstrumentDetails method
func (m *MockExecutionPlatform) GetInstrumentDetails(ctx context.Context, symbol string) (*models.Instrument, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Instrument), args.Error(1)
}

// GetExchangeInfo mocks the GetExchangeInfo method
func (m *MockExecutionPlatform) GetExchangeInfo(ctx context.Context, exchange string) (map[string]interface{}, error) {
	args := m.Called(ctx, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// GetTradingHours mocks the GetTradingHours method
func (m *MockExecutionPlatform) GetTradingHours(ctx context.Context, exchange string) (map[string]interface{}, error) {
	args := m.Called(ctx, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// GetStrategyTemplates mocks the GetStrategyTemplates method
func (m *MockExecutionPlatform) GetStrategyTemplates(ctx context.Context) ([]*models.StrategyTemplate, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.StrategyTemplate), args.Error(1)
}

// GetStrategyTemplate mocks the GetStrategyTemplate method
func (m *MockExecutionPlatform) GetStrategyTemplate(ctx context.Context, templateID string) (*models.StrategyTemplate, error) {
	args := m.Called(ctx, templateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StrategyTemplate), args.Error(1)
}

// GetSystemVersion mocks the GetSystemVersion method
func (m *MockExecutionPlatform) GetSystemVersion(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

// GetSupportedFeatures mocks the GetSupportedFeatures method
func (m *MockExecutionPlatform) GetSupportedFeatures(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// GetSystemHealth mocks the GetSystemHealth method
func (m *MockExecutionPlatform) GetSystemHealth(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// SynchronizeMarketData mocks the SynchronizeMarketData method
func (m *MockExecutionPlatform) SynchronizeMarketData(ctx context.Context, symbols []string) error {
	args := m.Called(ctx, symbols)
	return args.Error(0)
}

// MockSimulationService is a mock implementation of the SimulationAccountService
type MockSimulationService struct {
	mock.Mock
}

// CreateSimulationAccount mocks the CreateSimulationAccount method
func (m *MockSimulationService) CreateSimulationAccount(userID string, account models.SimulationAccount) (*models.SimulationAccount, error) {
	args := m.Called(userID, account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SimulationAccount), args.Error(1)
}

// GetSimulationAccount mocks the GetSimulationAccount method
func (m *MockSimulationService) GetSimulationAccount(accountID string) (*models.SimulationAccount, error) {
	args := m.Called(accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SimulationAccount), args.Error(1)
}

// GetSimulationAccountsByUser mocks the GetSimulationAccountsByUser method
func (m *MockSimulationService) GetSimulationAccountsByUser(userID string) ([]*models.SimulationAccount, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SimulationAccount), args.Error(1)
}

// UpdateSimulationAccount mocks the UpdateSimulationAccount method
func (m *MockSimulationService) UpdateSimulationAccount(accountID string, updates map[string]interface{}) (*models.SimulationAccount, error) {
	args := m.Called(accountID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SimulationAccount), args.Error(1)
}

// DeleteSimulationAccount mocks the DeleteSimulationAccount method
func (m *MockSimulationService) DeleteSimulationAccount(accountID string) error {
	args := m.Called(accountID)
	return args.Error(0)
}

// AddFunds mocks the AddFunds method
func (m *MockSimulationService) AddFunds(accountID string, amount float64, description string) (*models.SimulationTransaction, error) {
	args := m.Called(accountID, amount, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SimulationTransaction), args.Error(1)
}

// WithdrawFunds mocks the WithdrawFunds method
func (m *MockSimulationService) WithdrawFunds(accountID string, amount float64, description string) (*models.SimulationTransaction, error) {
	args := m.Called(accountID, amount, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SimulationTransaction), args.Error(1)
}

// GetTransactions mocks the GetTransactions method
func (m *MockSimulationService) GetTransactions(accountID string, startDate, endDate time.Time) ([]*models.SimulationTransaction, error) {
	args := m.Called(accountID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SimulationTransaction), args.Error(1)
}

// ResetAccount mocks the ResetAccount method
func (m *MockSimulationService) ResetAccount(accountID string) error {
	args := m.Called(accountID)
	return args.Error(0)
}

// MockVirtualBalanceService is a mock implementation of the VirtualBalanceService
type MockVirtualBalanceService struct {
	mock.Mock
}

// GetAccountBalance mocks the GetAccountBalance method
func (m *MockVirtualBalanceService) GetAccountBalance(accountID string) (float64, error) {
	args := m.Called(accountID)
	return args.Get(0).(float64), args.Error(1)
}

// GetAccountEquity mocks the GetAccountEquity method
func (m *MockVirtualBalanceService) GetAccountEquity(accountID string) (float64, error) {
	args := m.Called(accountID)
	return args.Get(0).(float64), args.Error(1)
}

// MockSimulationOrderService is a mock implementation of the SimulationOrderService
type MockSimulationOrderService struct {
	mock.Mock
}

// CreateOrder mocks the CreateOrder method
func (m *MockSimulationOrderService) CreateOrder(accountID string, order models.SimulationOrder) (*models.SimulationOrder, error) {
	args := m.Called(accountID, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SimulationOrder), args.Error(1)
}

// GetOrder mocks the GetOrder method
func (m *MockSimulationOrderService) GetOrder(orderID string) (*models.SimulationOrder, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SimulationOrder), args.Error(1)
}

// GetOrdersByAccount mocks the GetOrdersByAccount method
func (m *MockSimulationOrderService) GetOrdersByAccount(accountID string) ([]*models.SimulationOrder, error) {
	args := m.Called(accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SimulationOrder), args.Error(1)
}

// CancelOrder mocks the CancelOrder method
func (m *MockSimulationOrderService) CancelOrder(orderID string) (*models.SimulationOrder, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SimulationOrder), args.Error(1)
}

// ModifyOrder mocks the ModifyOrder method
func (m *MockSimulationOrderService) ModifyOrder(orderID string, updates models.SimulationOrder) (*models.SimulationOrder, error) {
	args := m.Called(orderID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SimulationOrder), args.Error(1)
}

// GetOrderHistory mocks the GetOrderHistory method
func (m *MockSimulationOrderService) GetOrderHistory(accountID string, startDate, endDate time.Time, symbol string) ([]*models.SimulationOrder, error) {
	args := m.Called(accountID, startDate, endDate, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SimulationOrder), args.Error(1)
}

// GetPositions mocks the GetPositions method
func (m *MockSimulationOrderService) GetPositions(accountID string) ([]*models.SimulationPosition, error) {
	args := m.Called(accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SimulationPosition), args.Error(1)
}

// GetPosition mocks the GetPosition method
func (m *MockSimulationOrderService) GetPosition(positionID string) (*models.SimulationPosition, error) {
	args := m.Called(positionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SimulationPosition), args.Error(1)
}

// ClosePosition mocks the ClosePosition method
func (m *MockSimulationOrderService) ClosePosition(positionID string, price float64) (*models.SimulationPosition, error) {
	args := m.Called(positionID, price)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SimulationPosition), args.Error(1)
}

// GetPositionHistory mocks the GetPositionHistory method
func (m *MockSimulationOrderService) GetPositionHistory(accountID string, startDate, endDate time.Time) ([]*models.SimulationPosition, error) {
	args := m.Called(accountID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SimulationPosition), args.Error(1)
}

// MockMarketSimulationService is a mock implementation of the MarketSimulationService
type MockMarketSimulationService struct {
	mock.Mock
}

// GetCurrentMarketPrice mocks the GetCurrentMarketPrice method
func (m *MockMarketSimulationService) GetCurrentMarketPrice(symbol string) (*models.MarketDataSnapshot, error) {
	args := m.Called(symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MarketDataSnapshot), args.Error(1)
}

// GetHistoricalMarketData mocks the GetHistoricalMarketData method
func (m *MockMarketSimulationService) GetHistoricalMarketData(symbol string, startDate, endDate time.Time, timeframe string) ([]*models.MarketDataSnapshot, error) {
	args := m.Called(symbol, startDate, endDate, timeframe)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MarketDataSnapshot), args.Error(1)
}

// GetMarketDepth mocks the GetMarketDepth method
func (m *MockMarketSimulationService) GetMarketDepth(symbol string, levels int) (map[string]interface{}, error) {
	args := m.Called(symbol, levels)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// UpdateMarketData mocks the UpdateMarketData method
func (m *MockMarketSimulationService) UpdateMarketData(symbol string, data *models.MarketDataSnapshot) error {
	args := m.Called(symbol, data)
	return args.Error(0)
}

// MockBacktestService is a mock implementation of the BacktestService
type MockBacktestService struct {
	mock.Mock
}

// CreateBacktestSession mocks the CreateBacktestSession method
func (m *MockBacktestService) CreateBacktestSession(accountID string, session models.BacktestSession) (*models.BacktestSession, error) {
	args := m.Called(accountID, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BacktestSession), args.Error(1)
}

// GetBacktestSession mocks the GetBacktestSession method
func (m *MockBacktestService) GetBacktestSession(sessionID string) (*models.BacktestSession, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BacktestSession), args.Error(1)
}

// GetBacktestSessionsByAccount mocks the GetBacktestSessionsByAccount method
func (m *MockBacktestService) GetBacktestSessionsByAccount(accountID string) ([]*models.BacktestSession, error) {
	args := m.Called(accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.BacktestSession), args.Error(1)
}

// RunBacktest mocks the RunBacktest method
func (m *MockBacktestService) RunBacktest(sessionID string) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

// StopBacktest mocks the StopBacktest method
func (m *MockBacktestService) StopBacktest(sessionID string) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

// GetBacktestResults mocks the GetBacktestResults method
func (m *MockBacktestService) GetBacktestResults(sessionID string) ([]*models.BacktestResult, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.BacktestResult), args.Error(1)
}

// GetBacktestPerformanceMetrics mocks the GetBacktestPerformanceMetrics method
func (m *MockBacktestService) GetBacktestPerformanceMetrics(sessionID string) (map[string]interface{}, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// TestAPIGateway tests the API Gateway implementation
func TestAPIGateway(t *testing.T) {
	// Create mock services
	mockExecutionPlatform := new(MockExecutionPlatform)
	mockSimulationService := new(MockSimulationService)
	mockVirtualBalanceService := new(MockVirtualBalanceService)
	mockSimulationOrderService := new(MockSimulationOrderService)
	mockMarketSimulationService := new(MockMarketSimulationService)
	mockBacktestService := new(MockBacktestService)
	
	// Create API Gateway with mock execution platform
	gateway := NewAPIGateway(mockExecutionPlatform)
	
	// Replace services with mocks
	gateway.simulationService = mockSimulationService
	gateway.virtualBalanceService = mockVirtualBalanceService
	gateway.simulationOrderService = mockSimulationOrderService
	gateway.marketSimulationService = mockMarketSimulationService
	gateway.backtestService = mockBacktestService
	
	// Add permissions for test user
	gateway.accessControlList["user123"] = []string{
		"simulation:account:create",
		"simulation:account:read",
		"simulation:account:update",
		"simulation:account:delete",
		"simulation:balance:add",
		"simulation:balance:withdraw",
		"simulation:balance:read",
		"simulation:transaction:read",
		"simulation:order:create",
		"simulation:order:read",
		"simulation:order:update",
		"simulation:order:cancel",
		"simulation:position:read",
		"simulation:position:close",
		"simulation:market:read",
		"backtest:session:create",
		"backtest:session:read",
		"backtest:session:run",
		"backtest:session:stop",
		"backtest:result:read",
		"system:status:read",
		"system:sync:execute",
	}
	
	// Create context with user ID and user type
	ctx := context.WithValue(context.Background(), "userID", "user123")
	ctx = context.WithValue(ctx, "userType", "SIM")
	
	t.Run("CreateSimulationAccount", func(t *testing.T) {
		// Setup mock
		account := models.SimulationAccount{
			Name:           "Test Account",
			InitialBalance: 100000.0,
			Currency:       "USD",
			SimulationType: "PAPER",
		}
		
		expectedAccount := &models.SimulationAccount{
			ID:             "sim123",
			Name:           "Test Account",
			InitialBalance: 100000.0,
			CurrentBalance: 100000.0,
			Currency:       "USD",
			SimulationType: "PAPER",
			IsActive:       true,
		}
		
		mockSimulationService.On("CreateSimulationAccount", "user123", account).Return(expectedAccount, nil)
		
		// Call method
		result, err := gateway.CreateSimulationAccount(ctx, "user123", account)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedAccount, result)
		mockSimulationService.AssertExpectations(t)
	})
	
	t.Run("GetSimulationAccount", func(t *testing.T) {
		// Setup mock
		expectedAccount := &models.SimulationAccount{
			ID:             "sim123",
			Name:           "Test Account",
			InitialBalance: 100000.0,
			CurrentBalance: 100000.0,
			Currency:       "USD",
			SimulationType: "PAPER",
			IsActive:       true,
		}
		
		mockSimulationService.On("GetSimulationAccount", "sim123").Return(expectedAccount, nil)
		
		// Call method
		result, err := gateway.GetSimulationAccount(ctx, "sim123")
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedAccount, result)
		mockSimulationService.AssertExpectations(t)
	})
	
	t.Run("GetAccountBalance", func(t *testing.T) {
		// Setup mock
		mockVirtualBalanceService.On("GetAccountBalance", "sim123").Return(100000.0, nil)
		
		// Call method
		result, err := gateway.GetAccountBalance(ctx, "sim123")
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 100000.0, result)
		mockVirtualBalanceService.AssertExpectations(t)
	})
	
	t.Run("CreateOrder", func(t *testing.T) {
		// Setup mock
		order := models.SimulationOrder{
			Order: models.Order{
				Symbol:    "AAPL",
				Quantity:  100,
				Side:      "BUY",
				OrderType: "MARKET",
			},
		}
		
		expectedOrder := &models.SimulationOrder{
			Order: models.Order{
				ID:        "order123",
				Symbol:    "AAPL",
				Quantity:  100,
				Side:      "BUY",
				OrderType: "MARKET",
				Status:    "PENDING",
			},
		}
		
		mockSimulationOrderService.On("CreateOrder", "sim123", order).Return(expectedOrder, nil)
		
		// Setup market data synchronization
		symbols := []string{"AAPL", "MSFT", "GOOGL"}
		mockExecutionPlatform.On("SynchronizeMarketData", ctx, symbols).Return(nil)
		
		// Call method
		result, err := gateway.CreateOrder(ctx, "sim123", order)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, result)
		mockSimulationOrderService.AssertExpectations(t)
	})
	
	t.Run("GetCurrentMarketPrice", func(t *testing.T) {
		// Setup mock
		expectedData := &models.MarketDataSnapshot{
			Symbol:    "AAPL",
			Price:     150.25,
			Timestamp: time.Now(),
		}
		
		mockMarketSimulationService.On("GetCurrentMarketPrice", "AAPL").Return(expectedData, nil)
		
		// Setup market data synchronization
		symbols := []string{"AAPL", "MSFT", "GOOGL"}
		mockExecutionPlatform.On("SynchronizeMarketData", ctx, symbols).Return(nil)
		
		// Call method
		result, err := gateway.GetCurrentMarketPrice(ctx, "AAPL")
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedData, result)
		mockMarketSimulationService.AssertExpectations(t)
	})
	
	t.Run("CreateBacktestSession", func(t *testing.T) {
		// Setup mock
		session := models.BacktestSession{
			Name:           "Test Backtest",
			StartDate:      time.Now().Add(-30 * 24 * time.Hour),
			EndDate:        time.Now(),
			Symbols:        []string{"AAPL", "MSFT"},
			InitialBalance: 100000.0,
		}
		
		expectedSession := &models.BacktestSession{
			ID:             "session123",
			Name:           "Test Backtest",
			StartDate:      time.Now().Add(-30 * 24 * time.Hour),
			EndDate:        time.Now(),
			Symbols:        []string{"AAPL", "MSFT"},
			InitialBalance: 100000.0,
			Status:         "PENDING",
		}
		
		mockBacktestService.On("CreateBacktestSession", "sim123", session).Return(expectedSession, nil)
		
		// Call method
		result, err := gateway.CreateBacktestSession(ctx, "sim123", session)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedSession, result)
		mockBacktestService.AssertExpectations(t)
	})
	
	t.Run("SynchronizeMarketData", func(t *testing.T) {
		// Setup mocks
		symbols := []string{"AAPL", "MSFT"}
		
		marketData1 := &models.MarketDataSnapshot{
			Symbol:    "AAPL",
			Price:     150.25,
			Timestamp: time.Now(),
		}
		
		marketData2 := &models.MarketDataSnapshot{
			Symbol:    "MSFT",
			Price:     290.75,
			Timestamp: time.Now(),
		}
		
		mockExecutionPlatform.On("GetRealTimeMarketData", ctx, "AAPL").Return(marketData1, nil)
		mockExecutionPlatform.On("GetRealTimeMarketData", ctx, "MSFT").Return(marketData2, nil)
		
		mockMarketSimulationService.On("UpdateMarketData", "AAPL", marketData1).Return(nil)
		mockMarketSimulationService.On("UpdateMarketData", "MSFT", marketData2).Return(nil)
		
		// Call method
		err := gateway.SynchronizeMarketData(ctx, symbols)
		
		// Assert
		assert.NoError(t, err)
		mockExecutionPlatform.AssertExpectations(t)
		mockMarketSimulationService.AssertExpectations(t)
	})
	
	t.Run("Permission Denied", func(t *testing.T) {
		// Create context with unauthorized user
		unauthorizedCtx := context.WithValue(context.Background(), "userID", "unauthorized")
		unauthorizedCtx = context.WithValue(unauthorizedCtx, "userType", "SIM")
		
		// Call method
		_, err := gateway.GetSimulationAccount(unauthorizedCtx, "sim123")
		
		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authorization")
	})
	
	t.Run("Rate Limit Exceeded", func(t *testing.T) {
		// Set up a rate limit that will be exceeded
		gateway.rateLimits["account_management"] = RateLimit{
			MaxRequests:     1,
			TimeWindow:      time.Minute,
			CurrentRequests: make(map[string][]time.Time),
		}
		
		// First request should succeed
		mockSimulationService.On("GetSimulationAccount", "sim123").Return(&models.SimulationAccount{}, nil)
		_, err := gateway.GetSimulationAccount(ctx, "sim123")
		assert.NoError(t, err)
		
		// Second request should fail due to rate limit
		_, err = gateway.GetSimulationAccount(ctx, "sim123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate limit")
	})
}

// TestInterfaceMetadata tests the interface metadata functionality
func TestInterfaceMetadata(t *testing.T) {
	metadata := interfaces.GetInterfaceMetadata()
	
	assert.Equal(t, "1.0.0", metadata.Version)
	assert.Contains(t, metadata.SupportedFeatures, "account_management")
	assert.Contains(t, metadata.SupportedFeatures, "order_management")
	assert.Contains(t, metadata.SupportedFeatures, "market_data")
	assert.Contains(t, metadata.SupportedFeatures, "backtesting")
	assert.Equal(t, 100, metadata.MaxBatchSize)
	assert.Contains(t, metadata.RateLimits, "market_data_requests_per_minute")
	assert.Contains(t, metadata.RateLimits, "order_requests_per_minute")
}
