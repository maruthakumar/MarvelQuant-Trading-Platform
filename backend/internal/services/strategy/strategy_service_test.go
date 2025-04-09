package strategy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trading-platform/backend/internal/models"
)

// MockStrategyRepository is a mock implementation of the StrategyRepository interface
type MockStrategyRepository struct {
	mock.Mock
}

func (m *MockStrategyRepository) Create(strategy *models.Strategy) (*models.Strategy, error) {
	args := m.Called(strategy)
	return args.Get(0).(*models.Strategy), args.Error(1)
}

func (m *MockStrategyRepository) GetByID(id string) (*models.Strategy, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Strategy), args.Error(1)
}

func (m *MockStrategyRepository) GetByUser(userID string) ([]models.Strategy, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Strategy), args.Error(1)
}

func (m *MockStrategyRepository) GetByTag(tag string) ([]models.Strategy, error) {
	args := m.Called(tag)
	return args.Get(0).([]models.Strategy), args.Error(1)
}

func (m *MockStrategyRepository) Update(strategy *models.Strategy) (*models.Strategy, error) {
	args := m.Called(strategy)
	return args.Get(0).(*models.Strategy), args.Error(1)
}

func (m *MockStrategyRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStrategyRepository) SaveSchedule(schedule *models.StrategySchedule) error {
	args := m.Called(schedule)
	return args.Error(0)
}

func (m *MockStrategyRepository) GetSchedule(strategyID string) (*models.StrategySchedule, error) {
	args := m.Called(strategyID)
	return args.Get(0).(*models.StrategySchedule), args.Error(1)
}

func (m *MockStrategyRepository) DeleteSchedule(strategyID string) error {
	args := m.Called(strategyID)
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

// TestCreateStrategy tests the CreateStrategy method
func TestCreateStrategy(t *testing.T) {
	// Create mock repositories
	mockStrategyRepo := new(MockStrategyRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	
	// Create the service
	service := NewStrategyService(mockStrategyRepo, mockOrderRepo, mockPositionRepo)
	
	// Create a sample strategy
	strategy := &models.Strategy{
		Name:        "Test Strategy",
		Description: "A test strategy",
		UserID:      "user123",
		Type:        models.StrategyTypeManual,
		Instruments: []string{"AAPL"},
		RiskParameters: models.RiskParameters{
			MaxPositionSize: 100,
			MaxLoss:         1000,
		},
	}
	
	// Set up the mock expectations
	mockStrategyRepo.On("Create", mock.AnythingOfType("*models.Strategy")).Return(strategy, nil)
	
	// Call the service method
	result, err := service.CreateStrategy(strategy)
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, strategy.Name, result.Name)
	assert.Equal(t, strategy.UserID, result.UserID)
	assert.Equal(t, models.StrategyStatusDraft, result.Status)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
	
	// Verify that the mock was called
	mockStrategyRepo.AssertExpectations(t)
}

// TestGetStrategyByID tests the GetStrategyByID method
func TestGetStrategyByID(t *testing.T) {
	// Create mock repositories
	mockStrategyRepo := new(MockStrategyRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	
	// Create the service
	service := NewStrategyService(mockStrategyRepo, mockOrderRepo, mockPositionRepo)
	
	// Create a sample strategy
	strategy := &models.Strategy{
		ID:          "strategy123",
		Name:        "Test Strategy",
		Description: "A test strategy",
		UserID:      "user123",
		Type:        models.StrategyTypeManual,
		Status:      models.StrategyStatusDraft,
		Instruments: []string{"AAPL"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Set up the mock expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(strategy, nil)
	
	// Call the service method
	result, err := service.GetStrategyByID("strategy123")
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, strategy.ID, result.ID)
	assert.Equal(t, strategy.Name, result.Name)
	
	// Verify that the mock was called
	mockStrategyRepo.AssertExpectations(t)
}

// TestGetStrategiesByUser tests the GetStrategiesByUser method
func TestGetStrategiesByUser(t *testing.T) {
	// Create mock repositories
	mockStrategyRepo := new(MockStrategyRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	
	// Create the service
	service := NewStrategyService(mockStrategyRepo, mockOrderRepo, mockPositionRepo)
	
	// Create sample strategies
	strategies := []models.Strategy{
		{
			ID:          "strategy1",
			Name:        "Strategy 1",
			Description: "First strategy",
			UserID:      "user123",
			Type:        models.StrategyTypeManual,
			Status:      models.StrategyStatusDraft,
			Instruments: []string{"AAPL"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "strategy2",
			Name:        "Strategy 2",
			Description: "Second strategy",
			UserID:      "user123",
			Type:        models.StrategyTypeAutomated,
			Status:      models.StrategyStatusActive,
			Instruments: []string{"MSFT"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	
	// Set up the mock expectations
	mockStrategyRepo.On("GetByUser", "user123").Return(strategies, nil)
	
	// Call the service method
	result, err := service.GetStrategiesByUser("user123")
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, strategies[0].ID, result[0].ID)
	assert.Equal(t, strategies[1].ID, result[1].ID)
	
	// Verify that the mock was called
	mockStrategyRepo.AssertExpectations(t)
}

// TestUpdateStrategy tests the UpdateStrategy method
func TestUpdateStrategy(t *testing.T) {
	// Create mock repositories
	mockStrategyRepo := new(MockStrategyRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	
	// Create the service
	service := NewStrategyService(mockStrategyRepo, mockOrderRepo, mockPositionRepo)
	
	// Create a sample strategy
	createdAt := time.Now().Add(-24 * time.Hour)
	strategy := &models.Strategy{
		ID:          "strategy123",
		Name:        "Test Strategy",
		Description: "A test strategy",
		UserID:      "user123",
		Type:        models.StrategyTypeManual,
		Status:      models.StrategyStatusDraft,
		Instruments: []string{"AAPL"},
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
	
	// Create an updated strategy
	updatedStrategy := &models.Strategy{
		ID:          "strategy123",
		Name:        "Updated Strategy",
		Description: "An updated test strategy",
		UserID:      "user123",
		Type:        models.StrategyTypeManual,
		Status:      models.StrategyStatusDraft,
		Instruments: []string{"AAPL", "MSFT"},
		RiskParameters: models.RiskParameters{
			MaxPositionSize: 200,
			MaxLoss:         2000,
		},
	}
	
	// Set up the mock expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(strategy, nil)
	mockStrategyRepo.On("Update", mock.AnythingOfType("*models.Strategy")).Return(updatedStrategy, nil)
	
	// Call the service method
	result, err := service.UpdateStrategy(updatedStrategy)
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, updatedStrategy.Name, result.Name)
	assert.Equal(t, updatedStrategy.Description, result.Description)
	assert.Equal(t, 2, len(result.Instruments))
	assert.Equal(t, float64(200), result.RiskParameters.MaxPositionSize)
	assert.Equal(t, createdAt, result.CreatedAt)
	assert.NotEqual(t, createdAt, result.UpdatedAt)
	
	// Verify that the mock was called
	mockStrategyRepo.AssertExpectations(t)
}

// TestDeleteStrategy tests the DeleteStrategy method
func TestDeleteStrategy(t *testing.T) {
	// Create mock repositories
	mockStrategyRepo := new(MockStrategyRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	
	// Create the service
	service := NewStrategyService(mockStrategyRepo, mockOrderRepo, mockPositionRepo)
	
	// Create a sample strategy
	strategy := &models.Strategy{
		ID:          "strategy123",
		Name:        "Test Strategy",
		Description: "A test strategy",
		UserID:      "user123",
		Type:        models.StrategyTypeManual,
		Status:      models.StrategyStatusDraft,
		Instruments: []string{"AAPL"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Set up the mock expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(strategy, nil)
	mockStrategyRepo.On("Delete", "strategy123").Return(nil)
	
	// Call the service method
	err := service.DeleteStrategy("strategy123")
	
	// Assert the result
	assert.NoError(t, err)
	
	// Verify that the mock was called
	mockStrategyRepo.AssertExpectations(t)
}

// TestExecuteStrategy tests the ExecuteStrategy method
func TestExecuteStrategy(t *testing.T) {
	// Create mock repositories
	mockStrategyRepo := new(MockStrategyRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	
	// Create the service
	service := NewStrategyService(mockStrategyRepo, mockOrderRepo, mockPositionRepo)
	
	// Create a sample strategy
	strategy := &models.Strategy{
		ID:          "strategy123",
		Name:        "Test Strategy",
		Description: "A test strategy",
		UserID:      "user123",
		Type:        models.StrategyTypeManual,
		Status:      models.StrategyStatusDraft,
		Instruments: []string{"AAPL"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Set up the mock expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(strategy, nil)
	mockStrategyRepo.On("Update", mock.AnythingOfType("*models.Strategy")).Return(strategy, nil)
	
	// Call the service method
	err := service.ExecuteStrategy("strategy123")
	
	// Assert the result
	assert.NoError(t, err)
	
	// Verify that the mock was called
	mockStrategyRepo.AssertExpectations(t)
}

// TestGetStrategyPerformance tests the GetStrategyPerformance method
func TestGetStrategyPerformance(t *testing.T) {
	// Create mock repositories
	mockStrategyRepo := new(MockStrategyRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	
	// Create the service
	service := NewStrategyService(mockStrategyRepo, mockOrderRepo, mockPositionRepo)
	
	// Create a sample strategy
	strategy := &models.Strategy{
		ID:          "strategy123",
		Name:        "Test Strategy",
		Description: "A test strategy",
		UserID:      "user123",
		Type:        models.StrategyTypeManual,
		Status:      models.StrategyStatusActive,
		Instruments: []string{"AAPL"},
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
	mockStrategyRepo.On("GetByID", "strategy123").Return(strategy, nil)
	mockOrderRepo.On("GetByStrategy", "strategy123").Return(orders, nil)
	mockPositionRepo.On("GetByStrategy", "strategy123").Return(positions, nil)
	
	// Call the service method
	performance, err := service.GetStrategyPerformance("strategy123")
	
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
	
	// Verify that the mock was called
	mockStrategyRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
	mockPositionRepo.AssertExpectations(t)
}

// TestScheduleStrategy tests the ScheduleStrategy method
func TestScheduleStrategy(t *testing.T) {
	// Create mock repositories
	mockStrategyRepo := new(MockStrategyRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	
	// Create the service
	service := NewStrategyService(mockStrategyRepo, mockOrderRepo, mockPositionRepo)
	
	// Create a sample strategy
	strategy := &models.Strategy{
		ID:          "strategy123",
		Name:        "Test Strategy",
		Description: "A test strategy",
		UserID:      "user123",
		Type:        models.StrategyTypeAutomated,
		Status:      models.StrategyStatusDraft,
		Instruments: []string{"AAPL"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Create a sample schedule
	schedule := &models.StrategySchedule{
		StrategyID: "strategy123",
		Frequency:  models.ScheduleFrequencyDaily,
		StartTime:  time.Now().Add(1 * time.Hour),
		Enabled:    true,
	}
	
	// Set up the mock expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(strategy, nil)
	mockStrategyRepo.On("SaveSchedule", mock.AnythingOfType("*models.StrategySchedule")).Return(nil)
	
	// Call the service method
	err := service.ScheduleStrategy("strategy123", schedule)
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, "strategy123", schedule.StrategyID)
	assert.NotZero(t, schedule.CreatedAt)
	assert.NotZero(t, schedule.UpdatedAt)
	
	// Verify that the mock was called
	mockStrategyRepo.AssertExpectations(t)
}

// TestAddStrategyTag tests the AddStrategyTag method
func TestAddStrategyTag(t *testing.T) {
	// Create mock repositories
	mockStrategyRepo := new(MockStrategyRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	
	// Create the service
	service := NewStrategyService(mockStrategyRepo, mockOrderRepo, mockPositionRepo)
	
	// Create a sample strategy
	strategy := &models.Strategy{
		ID:          "strategy123",
		Name:        "Test Strategy",
		Description: "A test strategy",
		UserID:      "user123",
		Type:        models.StrategyTypeManual,
		Status:      models.StrategyStatusDraft,
		Instruments: []string{"AAPL"},
		Tags:        []string{"tag1"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Set up the mock expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(strategy, nil)
	mockStrategyRepo.On("Update", mock.AnythingOfType("*models.Strategy")).Return(strategy, nil)
	
	// Call the service method
	err := service.AddStrategyTag("strategy123", "tag2")
	
	// Assert the result
	assert.NoError(t, err)
	
	// Verify that the mock was called
	mockStrategyRepo.AssertExpectations(t)
}

// TestGetStrategiesByTag tests the GetStrategiesByTag method
func TestGetStrategiesByTag(t *testing.T) {
	// Create mock repositories
	mockStrategyRepo := new(MockStrategyRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockPositionRepo := new(MockPositionRepository)
	
	// Create the service
	service := NewStrategyService(mockStrategyRepo, mockOrderRepo, mockPositionRepo)
	
	// Create sample strategies
	strategies := []models.Strategy{
		{
			ID:          "strategy1",
			Name:        "Strategy 1",
			Description: "First strategy",
			UserID:      "user123",
			Type:        models.StrategyTypeManual,
			Status:      models.StrategyStatusDraft,
			Instruments: []string{"AAPL"},
			Tags:        []string{"tag1", "tag2"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "strategy2",
			Name:        "Strategy 2",
			Description: "Second strategy",
			UserID:      "user456",
			Type:        models.StrategyTypeAutomated,
			Status:      models.StrategyStatusActive,
			Instruments: []string{"MSFT"},
			Tags:        []string{"tag1"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	
	// Set up the mock expectations
	mockStrategyRepo.On("GetByTag", "tag1").Return(strategies, nil)
	
	// Call the service method
	result, err := service.GetStrategiesByTag("tag1")
	
	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, strategies[0].ID, result[0].ID)
	assert.Equal(t, strategies[1].ID, result[1].ID)
	
	// Verify that the mock was called
	mockStrategyRepo.AssertExpectations(t)
}
