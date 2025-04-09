package multileg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trading-platform/backend/internal/models"
)

// MockMultilegRepository is a mock implementation of the MultilegRepository interface
type MockMultilegRepository struct {
	mock.Mock
}

func (m *MockMultilegRepository) CreateStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error) {
	args := m.Called(strategy)
	return args.Get(0).(*models.MultilegStrategy), args.Error(1)
}

func (m *MockMultilegRepository) GetStrategyByID(id string) (*models.MultilegStrategy, error) {
	args := m.Called(id)
	return args.Get(0).(*models.MultilegStrategy), args.Error(1)
}

func (m *MockMultilegRepository) GetStrategiesByUser(userID string) ([]models.MultilegStrategy, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.MultilegStrategy), args.Error(1)
}

func (m *MockMultilegRepository) GetStrategiesByPortfolio(portfolioID string) ([]models.MultilegStrategy, error) {
	args := m.Called(portfolioID)
	return args.Get(0).([]models.MultilegStrategy), args.Error(1)
}

func (m *MockMultilegRepository) UpdateStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error) {
	args := m.Called(strategy)
	return args.Get(0).(*models.MultilegStrategy), args.Error(1)
}

func (m *MockMultilegRepository) DeleteStrategy(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockOrderRepository is a mock implementation of the OrderRepository interface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) GetByStrategy(strategyID string) ([]models.Order, error) {
	args := m.Called(strategyID)
	return args.Get(0).([]models.Order), args.Error(1)
}

// MockPositionRepository is a mock implementation of the PositionRepository interface
type MockPositionRepository struct {
	mock.Mock
}

func (m *MockPositionRepository) GetByStrategy(strategyID string) ([]models.Position, error) {
	args := m.Called(strategyID)
	return args.Get(0).([]models.Position), args.Error(1)
}

// MockPortfolioRepository is a mock implementation of the PortfolioRepository interface
type MockPortfolioRepository struct {
	mock.Mock
}

func (m *MockPortfolioRepository) GetByID(id string) (*models.Portfolio, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Portfolio), args.Error(1)
}

// TestCreateMultilegStrategy tests the CreateMultilegStrategy method
func TestCreateMultilegStrategy(t *testing.T) {
	// Create mock repositories
	mockMultilegRepo := new(MockMultilegRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	mockPortfolioRepo := new(MockPortfolioRepository)
	
	// Create the service
	service := NewMultilegService(mockMultilegRepo, mockOrderRepo, mockPositionRepo, mockPortfolioRepo)
	
	// Create a sample portfolio
	portfolio := &models.Portfolio{
		ID:     "portfolio123",
		UserID: "user123",
		Name:   "Test Portfolio",
	}
	
	// Create a sample strategy
	strategy := &models.MultilegStrategy{
		Name:        "Test Multileg Strategy",
		Description: "A test multileg strategy",
		UserID:      "user123",
		PortfolioID: "portfolio123",
		Legs: []models.Leg{
			{
				Symbol:        "AAPL",
				Type:          models.LegTypeBuy,
				Quantity:      10,
				ExecutionType: models.ExecutionTypeMarket,
				Sequence:      1,
			},
		},
		ExecutionParams: models.ExecutionParams{
			Sequential:      true,
			SimultaneousLegs: false,
			MaxSlippage:     0.5,
		},
		RiskParams: models.RiskParameters{
			MaxLoss:      1000,
			MaxDailyLoss: 2000,
		},
	}
	
	// Set up the mock expectations
	mockPortfolioRepo.On("GetByID", "portfolio123").Return(portfolio, nil)
	mockMultilegRepo.On("CreateStrategy", mock.AnythingOfType("*models.MultilegStrategy")).Return(strategy, nil)
	
	// Call the service method
	result, err := service.CreateMultilegStrategy(strategy)
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, strategy.Name, result.Name)
	assert.Equal(t, strategy.UserID, result.UserID)
	assert.Equal(t, "DRAFT", result.Status)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
	
	// Verify that the mocks were called
	mockPortfolioRepo.AssertExpectations(t)
	mockMultilegRepo.AssertExpectations(t)
}

// TestGetMultilegStrategyByID tests the GetMultilegStrategyByID method
func TestGetMultilegStrategyByID(t *testing.T) {
	// Create mock repositories
	mockMultilegRepo := new(MockMultilegRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	mockPortfolioRepo := new(MockPortfolioRepository)
	
	// Create the service
	service := NewMultilegService(mockMultilegRepo, mockOrderRepo, mockPositionRepo, mockPortfolioRepo)
	
	// Create a sample strategy
	strategy := &models.MultilegStrategy{
		ID:          "strategy123",
		Name:        "Test Multileg Strategy",
		Description: "A test multileg strategy",
		UserID:      "user123",
		PortfolioID: "portfolio123",
		Status:      "DRAFT",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Set up the mock expectations
	mockMultilegRepo.On("GetStrategyByID", "strategy123").Return(strategy, nil)
	
	// Call the service method
	result, err := service.GetMultilegStrategyByID("strategy123")
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, strategy.ID, result.ID)
	assert.Equal(t, strategy.Name, result.Name)
	
	// Verify that the mock was called
	mockMultilegRepo.AssertExpectations(t)
}

// TestGetMultilegStrategiesByUser tests the GetMultilegStrategiesByUser method
func TestGetMultilegStrategiesByUser(t *testing.T) {
	// Create mock repositories
	mockMultilegRepo := new(MockMultilegRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	mockPortfolioRepo := new(MockPortfolioRepository)
	
	// Create the service
	service := NewMultilegService(mockMultilegRepo, mockOrderRepo, mockPositionRepo, mockPortfolioRepo)
	
	// Create sample strategies
	strategies := []models.MultilegStrategy{
		{
			ID:          "strategy1",
			Name:        "Strategy 1",
			Description: "First strategy",
			UserID:      "user123",
			PortfolioID: "portfolio123",
			Status:      "DRAFT",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "strategy2",
			Name:        "Strategy 2",
			Description: "Second strategy",
			UserID:      "user123",
			PortfolioID: "portfolio456",
			Status:      "ACTIVE",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	
	// Set up the mock expectations
	mockMultilegRepo.On("GetStrategiesByUser", "user123").Return(strategies, nil)
	
	// Call the service method
	result, err := service.GetMultilegStrategiesByUser("user123")
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, strategies[0].ID, result[0].ID)
	assert.Equal(t, strategies[1].ID, result[1].ID)
	
	// Verify that the mock was called
	mockMultilegRepo.AssertExpectations(t)
}

// TestUpdateMultilegStrategy tests the UpdateMultilegStrategy method
func TestUpdateMultilegStrategy(t *testing.T) {
	// Create mock repositories
	mockMultilegRepo := new(MockMultilegRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	mockPortfolioRepo := new(MockPortfolioRepository)
	
	// Create the service
	service := NewMultilegService(mockMultilegRepo, mockOrderRepo, mockPositionRepo, mockPortfolioRepo)
	
	// Create a sample portfolio
	portfolio := &models.Portfolio{
		ID:     "portfolio123",
		UserID: "user123",
		Name:   "Test Portfolio",
	}
	
	// Create a sample strategy
	createdAt := time.Now().Add(-24 * time.Hour)
	strategy := &models.MultilegStrategy{
		ID:          "strategy123",
		Name:        "Test Multileg Strategy",
		Description: "A test multileg strategy",
		UserID:      "user123",
		PortfolioID: "portfolio123",
		Status:      "DRAFT",
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
	
	// Create an updated strategy
	updatedStrategy := &models.MultilegStrategy{
		ID:          "strategy123",
		Name:        "Updated Multileg Strategy",
		Description: "An updated test multileg strategy",
		UserID:      "user123",
		PortfolioID: "portfolio123",
		Legs: []models.Leg{
			{
				Symbol:        "AAPL",
				Type:          models.LegTypeBuy,
				Quantity:      20,
				ExecutionType: models.ExecutionTypeMarket,
				Sequence:      1,
			},
		},
		ExecutionParams: models.ExecutionParams{
			Sequential:      false,
			SimultaneousLegs: true,
			MaxSlippage:     1.0,
		},
		RiskParams: models.RiskParameters{
			MaxLoss:      2000,
			MaxDailyLoss: 4000,
		},
	}
	
	// Set up the mock expectations
	mockMultilegRepo.On("GetStrategyByID", "strategy123").Return(strategy, nil)
	mockPortfolioRepo.On("GetByID", "portfolio123").Return(portfolio, nil)
	mockMultilegRepo.On("UpdateStrategy", mock.AnythingOfType("*models.MultilegStrategy")).Return(updatedStrategy, nil)
	
	// Call the service method
	result, err := service.UpdateMultilegStrategy(updatedStrategy)
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, updatedStrategy.Name, result.Name)
	assert.Equal(t, updatedStrategy.Description, result.Description)
	assert.Equal(t, 1, len(result.Legs))
	assert.Equal(t, float64(2000), result.RiskParams.MaxLoss)
	
	// Verify that the mocks were called
	mockMultilegRepo.AssertExpectations(t)
	mockPortfolioRepo.AssertExpectations(t)
}

// TestDeleteMultilegStrategy tests the DeleteMultilegStrategy method
func TestDeleteMultilegStrategy(t *testing.T) {
	// Create mock repositories
	mockMultilegRepo := new(MockMultilegRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	mockPortfolioRepo := new(MockPortfolioRepository)
	
	// Create the service
	service := NewMultilegService(mockMultilegRepo, mockOrderRepo, mockPositionRepo, mockPortfolioRepo)
	
	// Create a sample strategy
	strategy := &models.MultilegStrategy{
		ID:          "strategy123",
		Name:        "Test Multileg Strategy",
		Description: "A test multileg strategy",
		UserID:      "user123",
		PortfolioID: "portfolio123",
		Status:      "DRAFT",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Set up the mock expectations
	mockMultilegRepo.On("GetStrategyByID", "strategy123").Return(strategy, nil)
	mockMultilegRepo.On("DeleteStrategy", "strategy123").Return(nil)
	
	// Call the service method
	err := service.DeleteMultilegStrategy("strategy123")
	
	// Assert the result
	assert.NoError(t, err)
	
	// Verify that the mocks were called
	mockMultilegRepo.AssertExpectations(t)
}

// TestAddLeg tests the AddLeg method
func TestAddLeg(t *testing.T) {
	// Create mock repositories
	mockMultilegRepo := new(MockMultilegRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	mockPortfolioRepo := new(MockPortfolioRepository)
	
	// Create the service
	service := NewMultilegService(mockMultilegRepo, mockOrderRepo, mockPositionRepo, mockPortfolioRepo)
	
	// Create a sample strategy
	strategy := &models.MultilegStrategy{
		ID:          "strategy123",
		Name:        "Test Multileg Strategy",
		Description: "A test multileg strategy",
		UserID:      "user123",
		PortfolioID: "portfolio123",
		Status:      "DRAFT",
		Legs:        []models.Leg{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Create a sample leg
	leg := &models.Leg{
		Symbol:        "AAPL",
		Type:          models.LegTypeBuy,
		Quantity:      10,
		ExecutionType: models.ExecutionTypeMarket,
		Sequence:      1,
	}
	
	// Set up the mock expectations
	mockMultilegRepo.On("GetStrategyByID", "strategy123").Return(strategy, nil)
	mockMultilegRepo.On("UpdateStrategy", mock.AnythingOfType("*models.MultilegStrategy")).Return(strategy, nil)
	
	// Call the service method
	result, err := service.AddLeg("strategy123", leg)
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, leg.Symbol, result.Symbol)
	assert.Equal(t, leg.Type, result.Type)
	assert.Equal(t, leg.Quantity, result.Quantity)
	assert.Equal(t, models.LegStatusPending, result.Status)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
	
	// Verify that the mocks were called
	mockMultilegRepo.AssertExpectations(t)
}

// TestExecuteMultilegStrategy tests the ExecuteMultilegStrategy method
func TestExecuteMultilegStrategy(t *testing.T) {
	// Create mock repositories
	mockMultilegRepo := new(MockMultilegRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	mockPortfolioRepo := new(MockPortfolioRepository)
	
	// Create the service
	service := NewMultilegService(mockMultilegRepo, mockOrderRepo, mockPositionRepo, mockPortfolioRepo)
	
	// Create a sample strategy
	strategy := &models.MultilegStrategy{
		ID:          "strategy123",
		Name:        "Test Multileg Strategy",
		Description: "A test multileg strategy",
		UserID:      "user123",
		PortfolioID: "portfolio123",
		Status:      "DRAFT",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Set up the mock expectations
	mockMultilegRepo.On("GetStrategyByID", "strategy123").Return(strategy, nil)
	mockMultilegRepo.On("UpdateStrategy", mock.AnythingOfType("*models.MultilegStrategy")).Return(strategy, nil)
	
	// Call the service method
	err := service.ExecuteMultilegStrategy("strategy123")
	
	// Assert the result
	assert.NoError(t, err)
	
	// Verify that the mocks were called
	mockMultilegRepo.AssertExpectations(t)
}

// TestGetMultilegStrategyPerformance tests the GetMultilegStrategyPerformance method
func TestGetMultilegStrategyPerformance(t *testing.T) {
	// Create mock repositories
	mockMultilegRepo := new(MockMultilegRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	mockPortfolioRepo := new(MockPortfolioRepository)
	
	// Create the service
	service := NewMultilegService(mockMultilegRepo, mockOrderRepo, mockPositionRepo, mockPortfolioRepo)
	
	// Create a sample strategy
	strategy := &models.MultilegStrategy{
		ID:          "strategy123",
		Name:        "Test Multileg Strategy",
		Description: "A test multileg strategy",
		UserID:      "user123",
		PortfolioID: "portfolio123",
		Status:      "ACTIVE",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}
	
	// Create sample orders
	orders := []models.Order{
		{
			ID:         "order1",
			UserID:     "user123",
			StrategyID: "strategy123",
			Symbol:     "AAPL",
			Side:       models.OrderSideBuy,
			Quantity:   10,
			Price:      150.0,
			Status:     models.OrderStatusFilled,
			CreatedAt:  time.Now().Add(-12 * time.Hour),
			UpdatedAt:  time.Now().Add(-11 * time.Hour),
		},
		{
			ID:         "order2",
			UserID:     "user123",
			StrategyID: "strategy123",
			Symbol:     "AAPL",
			Side:       models.OrderSideSell,
			Quantity:   10,
			Price:      160.0,
			Status:     models.OrderStatusFilled,
			CreatedAt:  time.Now().Add(-6 * time.Hour),
			UpdatedAt:  time.Now().Add(-5 * time.Hour),
		},
	}
	
	// Create sample positions
	positions := []models.Position{
		{
			ID:          "position1",
			UserID:      "user123",
			StrategyID:  "strategy123",
			Symbol:      "AAPL",
			Quantity:    0,
			EntryPrice:  150.0,
			ExitPrice:   160.0,
			RealizedPnL: 100.0,
			Status:      models.PositionStatusClosed,
			CreatedAt:   time.Now().Add(-12 * time.Hour),
			UpdatedAt:   time.Now().Add(-5 * time.Hour),
		},
	}
	
	// Set up the mock expectations
	mockMultilegRepo.On("GetStrategyByID", "strategy123").Return(strategy, nil)
	mockOrderRepo.On("GetByStrategy", "strategy123").Return(orders, nil)
	mockPositionRepo.On("GetByStrategy", "strategy123").Return(positions, nil)
	
	// Call the service method
	performance, err := service.GetMultilegStrategyPerformance("strategy123")
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, "strategy123", performance.StrategyID)
	assert.Equal(t, float64(100.0), performance.TotalPnL)
	assert.Equal(t, 1, performance.WinCount)
	assert.Equal(t, 0, performance.LossCount)
	assert.Equal(t, 1, performance.TotalTrades)
	assert.Equal(t, float64(100.0), performance.WinRate)
	assert.Equal(t, 2, performance.OrderCount)
	assert.Equal(t, 1, performance.PositionCount)
	
	// Verify that the mocks were called
	mockMultilegRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
	mockPositionRepo.AssertExpectations(t)
}
