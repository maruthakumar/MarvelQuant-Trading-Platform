package interfaces

import (
	"context"
	"time"
	
	"trading_platform/backend/internal/models"
)

// ExecutionSimulationInterface defines the contract between the execution platform
// and the simulation system. This interface ensures proper isolation while allowing
// controlled communication between the two systems.
type ExecutionSimulationInterface interface {
	// Account Management
	CreateSimulationAccount(ctx context.Context, userID string, account models.SimulationAccount) (*models.SimulationAccount, error)
	GetSimulationAccount(ctx context.Context, accountID string) (*models.SimulationAccount, error)
	GetSimulationAccountsByUser(ctx context.Context, userID string) ([]*models.SimulationAccount, error)
	UpdateSimulationAccount(ctx context.Context, accountID string, updates map[string]interface{}) (*models.SimulationAccount, error)
	DeleteSimulationAccount(ctx context.Context, accountID string) error
	
	// Balance Management
	AddFunds(ctx context.Context, accountID string, amount float64, description string) (*models.SimulationTransaction, error)
	WithdrawFunds(ctx context.Context, accountID string, amount float64, description string) (*models.SimulationTransaction, error)
	GetAccountBalance(ctx context.Context, accountID string) (float64, error)
	GetAccountEquity(ctx context.Context, accountID string) (float64, error)
	GetTransactions(ctx context.Context, accountID string, startDate, endDate time.Time) ([]*models.SimulationTransaction, error)
	
	// Order Management
	CreateOrder(ctx context.Context, accountID string, order models.SimulationOrder) (*models.SimulationOrder, error)
	GetOrder(ctx context.Context, orderID string) (*models.SimulationOrder, error)
	GetOrdersByAccount(ctx context.Context, accountID string) ([]*models.SimulationOrder, error)
	CancelOrder(ctx context.Context, orderID string) (*models.SimulationOrder, error)
	ModifyOrder(ctx context.Context, orderID string, updates models.SimulationOrder) (*models.SimulationOrder, error)
	GetOrderHistory(ctx context.Context, accountID string, startDate, endDate time.Time, symbol string) ([]*models.SimulationOrder, error)
	
	// Position Management
	GetPositions(ctx context.Context, accountID string) ([]*models.SimulationPosition, error)
	GetPosition(ctx context.Context, positionID string) (*models.SimulationPosition, error)
	ClosePosition(ctx context.Context, positionID string, price float64) (*models.SimulationPosition, error)
	GetPositionHistory(ctx context.Context, accountID string, startDate, endDate time.Time) ([]*models.SimulationPosition, error)
	
	// Market Data
	GetCurrentMarketPrice(ctx context.Context, symbol string) (*models.MarketDataSnapshot, error)
	GetHistoricalMarketData(ctx context.Context, symbol string, startDate, endDate time.Time, timeframe string) ([]*models.MarketDataSnapshot, error)
	GetMarketDepth(ctx context.Context, symbol string, levels int) (map[string]interface{}, error)
	
	// Backtesting
	CreateBacktestSession(ctx context.Context, accountID string, session models.BacktestSession) (*models.BacktestSession, error)
	GetBacktestSession(ctx context.Context, sessionID string) (*models.BacktestSession, error)
	GetBacktestSessionsByAccount(ctx context.Context, accountID string) ([]*models.BacktestSession, error)
	RunBacktest(ctx context.Context, sessionID string) error
	StopBacktest(ctx context.Context, sessionID string) error
	GetBacktestResults(ctx context.Context, sessionID string) ([]*models.BacktestResult, error)
	GetBacktestPerformanceMetrics(ctx context.Context, sessionID string) (map[string]interface{}, error)
	
	// System Management
	GetSystemStatus(ctx context.Context) (map[string]interface{}, error)
	SynchronizeMarketData(ctx context.Context, symbols []string) error
	ResetSimulationEnvironment(ctx context.Context, accountID string) error
}

// ExecutionPlatformInterface defines the contract for the execution platform
// to expose functionality to the simulation system. This interface ensures
// that the simulation system can access necessary execution platform features
// while maintaining proper isolation.
type ExecutionPlatformInterface interface {
	// Market Data Access
	GetRealTimeMarketData(ctx context.Context, symbol string) (*models.MarketDataSnapshot, error)
	GetHistoricalMarketData(ctx context.Context, symbol string, startDate, endDate time.Time, timeframe string) ([]*models.MarketDataSnapshot, error)
	SubscribeToMarketData(ctx context.Context, symbol string, callback func(*models.MarketDataSnapshot)) (string, error)
	UnsubscribeFromMarketData(ctx context.Context, subscriptionID string) error
	
	// Reference Data
	GetInstrumentDetails(ctx context.Context, symbol string) (*models.Instrument, error)
	GetExchangeInfo(ctx context.Context, exchange string) (map[string]interface{}, error)
	GetTradingHours(ctx context.Context, exchange string) (map[string]interface{}, error)
	
	// Strategy Templates
	GetStrategyTemplates(ctx context.Context) ([]*models.StrategyTemplate, error)
	GetStrategyTemplate(ctx context.Context, templateID string) (*models.StrategyTemplate, error)
	
	// System Information
	GetSystemVersion(ctx context.Context) (string, error)
	GetSupportedFeatures(ctx context.Context) ([]string, error)
	GetSystemHealth(ctx context.Context) (map[string]interface{}, error)
}

// InterfaceMetadata contains information about the interface version and capabilities
type InterfaceMetadata struct {
	Version            string   `json:"version"`
	SupportedFeatures  []string `json:"supportedFeatures"`
	DeprecatedFeatures []string `json:"deprecatedFeatures"`
	MaxBatchSize       int      `json:"maxBatchSize"`
	RateLimits         map[string]int `json:"rateLimits"`
}

// GetInterfaceMetadata returns metadata about the current interface implementation
func GetInterfaceMetadata() InterfaceMetadata {
	return InterfaceMetadata{
		Version: "1.0.0",
		SupportedFeatures: []string{
			"account_management",
			"balance_management",
			"order_management",
			"position_management",
			"market_data",
			"backtesting",
			"system_management",
		},
		DeprecatedFeatures: []string{},
		MaxBatchSize: 100,
		RateLimits: map[string]int{
			"market_data_requests_per_minute": 300,
			"order_requests_per_minute": 100,
			"account_requests_per_minute": 60,
		},
	}
}
