package services

import (
	"errors"
	"time"
	"github.com/google/uuid"
	"trading_platform/backend/internal/models"
)

// SimulationAccountService handles operations related to simulation accounts
type SimulationAccountService struct {
	// Dependencies would be injected here in a real implementation
	// For example: database connection, market data service, etc.
}

// NewSimulationAccountService creates a new instance of SimulationAccountService
func NewSimulationAccountService() *SimulationAccountService {
	return &SimulationAccountService{}
}

// CreateSimulationAccount creates a new simulation account
func (s *SimulationAccountService) CreateSimulationAccount(userID string, accountData models.SimulationAccount) (*models.SimulationAccount, error) {
	// Validate input
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	
	if accountData.Name == "" {
		return nil, errors.New("account name is required")
	}
	
	if accountData.InitialBalance <= 0 {
		return nil, errors.New("initial balance must be greater than zero")
	}
	
	if accountData.SimulationType != "PAPER" && accountData.SimulationType != "BACKTEST" {
		return nil, errors.New("simulation type must be either PAPER or BACKTEST")
	}
	
	// Create new account
	account := models.SimulationAccount{
		ID:              uuid.New().String(),
		UserID:          userID,
		Name:            accountData.Name,
		Description:     accountData.Description,
		InitialBalance:  accountData.InitialBalance,
		CurrentBalance:  accountData.InitialBalance, // Initially set to initial balance
		Currency:        accountData.Currency,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsActive:        true,
		SimulationType:  accountData.SimulationType,
		RiskSettings:    accountData.RiskSettings,
		MarketSettings:  accountData.MarketSettings,
	}
	
	// Set default risk settings if not provided
	if account.RiskSettings == nil {
		account.RiskSettings = &models.RiskSettings{
			MaxPositionSize:       account.InitialBalance * 0.1, // 10% of initial balance
			MaxDrawdown:           account.InitialBalance * 0.2, // 20% of initial balance
			MaxDailyLoss:          account.InitialBalance * 0.05, // 5% of initial balance
			MaxOpenPositions:      10,
			MaxLeverage:           1.0, // No leverage by default
			StopLossRequired:      true,
			TakeProfitRecommended: true,
		}
	}
	
	// Set default market settings if not provided
	if account.MarketSettings == nil {
		account.MarketSettings = &models.MarketSettings{
			SlippageModel:       "PERCENTAGE",
			SlippageValue:       0.001, // 0.1% slippage
			LatencyModel:        "FIXED",
			LatencyValue:        100, // 100ms latency
			PriceFeedSource:     "REAL_TIME",
			CommissionModel:     "PERCENTAGE",
			CommissionValue:     0.001, // 0.1% commission
			SpreadModel:         "REALISTIC",
			SpreadValue:         0.0,
			AllowShortSelling:   true,
			AllowFractionalLots: true,
		}
	}
	
	// In a real implementation, we would save the account to the database here
	
	// Create initial deposit transaction
	depositTransaction := models.SimulationTransaction{
		ID:                 uuid.New().String(),
		SimulationAccountID: account.ID,
		Type:               "DEPOSIT",
		Amount:             account.InitialBalance,
		Balance:            account.InitialBalance,
		Description:        "Initial deposit",
		ReferenceID:        "",
		ReferenceType:      "",
		Timestamp:          time.Now(),
	}
	
	// In a real implementation, we would save the transaction to the database here
	
	return &account, nil
}

// GetSimulationAccount retrieves a simulation account by ID
func (s *SimulationAccountService) GetSimulationAccount(accountID string) (*models.SimulationAccount, error) {
	if accountID == "" {
		return nil, errors.New("account ID is required")
	}
	
	// In a real implementation, we would retrieve the account from the database here
	
	// For now, return a mock account
	return &models.SimulationAccount{
		ID:              accountID,
		UserID:          "user123",
		Name:            "Test Simulation Account",
		Description:     "Test account for simulation",
		InitialBalance:  100000.0,
		CurrentBalance:  105000.0,
		Currency:        "USD",
		CreatedAt:       time.Now().Add(-24 * time.Hour),
		UpdatedAt:       time.Now(),
		IsActive:        true,
		SimulationType:  "PAPER",
		RiskSettings:    &models.RiskSettings{},
		MarketSettings:  &models.MarketSettings{},
	}, nil
}

// UpdateSimulationAccount updates a simulation account
func (s *SimulationAccountService) UpdateSimulationAccount(accountID string, accountData models.SimulationAccount) (*models.SimulationAccount, error) {
	if accountID == "" {
		return nil, errors.New("account ID is required")
	}
	
	// In a real implementation, we would retrieve the account from the database,
	// update it with the new data, and save it back to the database
	
	// For now, return a mock updated account
	return &models.SimulationAccount{
		ID:              accountID,
		UserID:          "user123",
		Name:            accountData.Name,
		Description:     accountData.Description,
		InitialBalance:  accountData.InitialBalance,
		CurrentBalance:  accountData.CurrentBalance,
		Currency:        accountData.Currency,
		CreatedAt:       time.Now().Add(-24 * time.Hour),
		UpdatedAt:       time.Now(),
		IsActive:        accountData.IsActive,
		SimulationType:  accountData.SimulationType,
		RiskSettings:    accountData.RiskSettings,
		MarketSettings:  accountData.MarketSettings,
	}, nil
}

// DeleteSimulationAccount deletes a simulation account
func (s *SimulationAccountService) DeleteSimulationAccount(accountID string) error {
	if accountID == "" {
		return errors.New("account ID is required")
	}
	
	// In a real implementation, we would delete the account from the database
	
	return nil
}

// GetSimulationAccountsByUser retrieves all simulation accounts for a user
func (s *SimulationAccountService) GetSimulationAccountsByUser(userID string) ([]models.SimulationAccount, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	
	// In a real implementation, we would retrieve the accounts from the database
	
	// For now, return mock accounts
	return []models.SimulationAccount{
		{
			ID:              "sim1",
			UserID:          userID,
			Name:            "Paper Trading Account",
			Description:     "Account for paper trading",
			InitialBalance:  100000.0,
			CurrentBalance:  105000.0,
			Currency:        "USD",
			CreatedAt:       time.Now().Add(-24 * time.Hour),
			UpdatedAt:       time.Now(),
			IsActive:        true,
			SimulationType:  "PAPER",
			RiskSettings:    &models.RiskSettings{},
			MarketSettings:  &models.MarketSettings{},
		},
		{
			ID:              "sim2",
			UserID:          userID,
			Name:            "Backtesting Account",
			Description:     "Account for backtesting strategies",
			InitialBalance:  50000.0,
			CurrentBalance:  48000.0,
			Currency:        "USD",
			CreatedAt:       time.Now().Add(-48 * time.Hour),
			UpdatedAt:       time.Now(),
			IsActive:        true,
			SimulationType:  "BACKTEST",
			RiskSettings:    &models.RiskSettings{},
			MarketSettings:  &models.MarketSettings{},
		},
	}, nil
}

// AddFunds adds funds to a simulation account
func (s *SimulationAccountService) AddFunds(accountID string, amount float64, description string) (*models.SimulationTransaction, error) {
	if accountID == "" {
		return nil, errors.New("account ID is required")
	}
	
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	
	// In a real implementation, we would retrieve the account from the database,
	// update its balance, and create a transaction
	
	// For now, create a mock transaction
	transaction := models.SimulationTransaction{
		ID:                 uuid.New().String(),
		SimulationAccountID: accountID,
		Type:               "DEPOSIT",
		Amount:             amount,
		Balance:            105000.0 + amount, // Mock current balance + deposit
		Description:        description,
		ReferenceID:        "",
		ReferenceType:      "",
		Timestamp:          time.Now(),
	}
	
	return &transaction, nil
}

// WithdrawFunds withdraws funds from a simulation account
func (s *SimulationAccountService) WithdrawFunds(accountID string, amount float64, description string) (*models.SimulationTransaction, error) {
	if accountID == "" {
		return nil, errors.New("account ID is required")
	}
	
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	
	// In a real implementation, we would retrieve the account from the database,
	// check if it has sufficient balance, update its balance, and create a transaction
	
	// For now, create a mock transaction
	transaction := models.SimulationTransaction{
		ID:                 uuid.New().String(),
		SimulationAccountID: accountID,
		Type:               "WITHDRAWAL",
		Amount:             -amount, // Negative amount for withdrawal
		Balance:            105000.0 - amount, // Mock current balance - withdrawal
		Description:        description,
		ReferenceID:        "",
		ReferenceType:      "",
		Timestamp:          time.Now(),
	}
	
	return &transaction, nil
}

// GetTransactions retrieves transactions for a simulation account
func (s *SimulationAccountService) GetTransactions(accountID string, startDate, endDate time.Time) ([]models.SimulationTransaction, error) {
	if accountID == "" {
		return nil, errors.New("account ID is required")
	}
	
	// In a real implementation, we would retrieve the transactions from the database
	
	// For now, return mock transactions
	return []models.SimulationTransaction{
		{
			ID:                 "trans1",
			SimulationAccountID: accountID,
			Type:               "DEPOSIT",
			Amount:             100000.0,
			Balance:            100000.0,
			Description:        "Initial deposit",
			ReferenceID:        "",
			ReferenceType:      "",
			Timestamp:          time.Now().Add(-24 * time.Hour),
		},
		{
			ID:                 "trans2",
			SimulationAccountID: accountID,
			Type:               "P&L",
			Amount:             5000.0,
			Balance:            105000.0,
			Description:        "Realized profit from AAPL trade",
			ReferenceID:        "order123",
			ReferenceType:      "ORDER",
			Timestamp:          time.Now().Add(-12 * time.Hour),
		},
	}, nil
}

// ResetAccount resets a simulation account to its initial state
func (s *SimulationAccountService) ResetAccount(accountID string) error {
	if accountID == "" {
		return errors.New("account ID is required")
	}
	
	// In a real implementation, we would:
	// 1. Retrieve the account from the database
	// 2. Close all open positions
	// 3. Cancel all open orders
	// 4. Reset the balance to the initial balance
	// 5. Create a reset transaction
	// 6. Update the account in the database
	
	return nil
}

// GetAccountPerformance retrieves performance metrics for a simulation account
func (s *SimulationAccountService) GetAccountPerformance(accountID string, startDate, endDate time.Time) (map[string]interface{}, error) {
	if accountID == "" {
		return nil, errors.New("account ID is required")
	}
	
	// In a real implementation, we would calculate performance metrics based on
	// transactions, positions, and orders
	
	// For now, return mock performance metrics
	return map[string]interface{}{
		"totalPnL":           5000.0,
		"percentageReturn":   5.0,
		"winRate":            0.65,
		"averageWin":         1200.0,
		"averageLoss":        -500.0,
		"profitFactor":       2.4,
		"maxDrawdown":        -2000.0,
		"maxDrawdownPercent": -2.0,
		"sharpeRatio":        1.8,
		"totalTrades":        20,
		"winningTrades":      13,
		"losingTrades":       7,
	}, nil
}
