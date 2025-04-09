package position

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/repositories"
)

// MockPositionRepository is a mock implementation of the PositionRepository interface
type MockPositionRepository struct {
	mock.Mock
}

func (m *MockPositionRepository) Create(position *models.Position) (*models.Position, error) {
	args := m.Called(position)
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionRepository) GetByID(id string) (*models.Position, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionRepository) GetAll(filter models.PositionFilter, offset, limit int) ([]models.Position, int, error) {
	args := m.Called(filter, offset, limit)
	return args.Get(0).([]models.Position), args.Int(1), args.Error(2)
}

func (m *MockPositionRepository) Update(position *models.Position) (*models.Position, error) {
	args := m.Called(position)
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockOrderRepository is a mock implementation of the OrderRepository interface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *models.Order) (*models.Order, error) {
	args := m.Called(order)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByID(id string) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetAll(filter models.OrderFilter, offset, limit int) ([]models.Order, int, error) {
	args := m.Called(filter, offset, limit)
	return args.Get(0).([]models.Order), args.Int(1), args.Error(2)
}

func (m *MockOrderRepository) Update(order *models.Order) (*models.Order, error) {
	args := m.Called(order)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreatePositionFromOrder(t *testing.T) {
	// Create mock repositories
	mockPositionRepo := new(MockPositionRepository)
	mockOrderRepo := new(MockOrderRepository)
	
	// Create a sample executed order
	order := &models.Order{
		ID:             "order123",
		UserID:         "user123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		OrderType:      models.OrderTypeLimit,
		Direction:      models.OrderDirectionBuy,
		Quantity:       10,
		Price:          500.50,
		ExecutionPrice: 500.75,
		FilledQuantity: 10,
		Status:         models.OrderStatusExecuted,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
		OptionType:     models.OptionTypeCall,
		StrikePrice:    18000,
		Expiry:         time.Now().AddDate(0, 1, 0),
	}
	
	// Set up the mock repository expectations
	mockPositionRepo.On("GetAll", mock.AnythingOfType("models.PositionFilter"), 0, 1).Return([]models.Position{}, 0, nil)
	mockPositionRepo.On("Create", mock.AnythingOfType("*models.Position")).Return(func(position *models.Position) *models.Position {
		position.ID = "position123"
		return position
	}, nil)
	
	// Create the service with the mock repositories
	service := NewPositionService(mockPositionRepo, mockOrderRepo)
	
	// Call the service method
	createdPosition, err := service.CreatePositionFromOrder(order)
	
	// Check the result
	assert.NoError(t, err)
	assert.NotNil(t, createdPosition)
	assert.Equal(t, "position123", createdPosition.ID)
	assert.Equal(t, order.UserID, createdPosition.UserID)
	assert.Equal(t, order.ID, createdPosition.OrderID)
	assert.Equal(t, order.Symbol, createdPosition.Symbol)
	assert.Equal(t, models.PositionDirectionLong, createdPosition.Direction)
	assert.Equal(t, order.ExecutionPrice, createdPosition.EntryPrice)
	assert.Equal(t, order.FilledQuantity, createdPosition.Quantity)
	assert.Equal(t, models.PositionStatusOpen, createdPosition.Status)
	
	// Verify that the mock repositories were called
	mockPositionRepo.AssertExpectations(t)
}

func TestGetPositionByID(t *testing.T) {
	// Create mock repositories
	mockPositionRepo := new(MockPositionRepository)
	mockOrderRepo := new(MockOrderRepository)
	
	// Create a sample position
	position := &models.Position{
		ID:             "position123",
		UserID:         "user123",
		OrderID:        "order123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		Direction:      models.PositionDirectionLong,
		EntryPrice:     500.75,
		Quantity:       10,
		Status:         models.PositionStatusOpen,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
		OptionType:     models.OptionTypeCall,
		StrikePrice:    18000,
		Expiry:         time.Now().AddDate(0, 1, 0),
		UnrealizedPnL:  100.0,
		Greeks: models.Greeks{
			Delta: 0.6,
			Gamma: 0.05,
			Theta: -0.1,
			Vega:  0.2,
		},
	}
	
	// Set up the mock repository expectations
	mockPositionRepo.On("GetByID", "position123").Return(position, nil)
	mockPositionRepo.On("GetByID", "nonexistent").Return(nil, assert.AnError)
	
	// Create the service with the mock repositories
	service := NewPositionService(mockPositionRepo, mockOrderRepo)
	
	// Test successful retrieval
	retrievedPosition, err := service.GetPositionByID("position123")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedPosition)
	assert.Equal(t, position.ID, retrievedPosition.ID)
	
	// Test error case
	retrievedPosition, err = service.GetPositionByID("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, retrievedPosition)
	
	// Test empty ID
	retrievedPosition, err = service.GetPositionByID("")
	assert.Error(t, err)
	assert.Nil(t, retrievedPosition)
	
	// Verify that the mock repositories were called
	mockPositionRepo.AssertExpectations(t)
}

func TestGetPositions(t *testing.T) {
	// Create mock repositories
	mockPositionRepo := new(MockPositionRepository)
	mockOrderRepo := new(MockOrderRepository)
	
	// Create sample positions
	positions := []models.Position{
		{
			ID:             "position123",
			UserID:         "user123",
			OrderID:        "order123",
			Symbol:         "NIFTY",
			Exchange:       "NSE",
			Direction:      models.PositionDirectionLong,
			EntryPrice:     500.75,
			Quantity:       10,
			Status:         models.PositionStatusOpen,
			ProductType:    models.ProductTypeMIS,
			InstrumentType: models.InstrumentTypeOption,
		},
		{
			ID:             "position456",
			UserID:         "user123",
			OrderID:        "order456",
			Symbol:         "BANKNIFTY",
			Exchange:       "NSE",
			Direction:      models.PositionDirectionShort,
			EntryPrice:     1200.50,
			Quantity:       5,
			Status:         models.PositionStatusOpen,
			ProductType:    models.ProductTypeMIS,
			InstrumentType: models.InstrumentTypeFuture,
		},
	}
	
	// Set up the mock repository expectations
	mockPositionRepo.On("GetAll", mock.AnythingOfType("models.PositionFilter"), 0, 50).Return(positions, 2, nil)
	
	// Create the service with the mock repositories
	service := NewPositionService(mockPositionRepo, mockOrderRepo)
	
	// Test successful retrieval with default pagination
	filter := models.PositionFilter{UserID: "user123"}
	retrievedPositions, total, err := service.GetPositions(filter, 1, 50)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(retrievedPositions))
	assert.Equal(t, 2, total)
	
	// Verify that the mock repositories were called
	mockPositionRepo.AssertExpectations(t)
}

func TestUpdatePosition(t *testing.T) {
	// Create mock repositories
	mockPositionRepo := new(MockPositionRepository)
	mockOrderRepo := new(MockOrderRepository)
	
	// Create sample positions
	existingPosition := &models.Position{
		ID:             "position123",
		UserID:         "user123",
		OrderID:        "order123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		Direction:      models.PositionDirectionLong,
		EntryPrice:     500.75,
		Quantity:       10,
		Status:         models.PositionStatusOpen,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
		CreatedAt:      time.Now().Add(-time.Hour),
	}
	
	updatedPosition := &models.Position{
		ID:             "position123",
		UserID:         "user123",
		OrderID:        "order123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		Direction:      models.PositionDirectionLong,
		EntryPrice:     500.75,
		Quantity:       10,
		Status:         models.PositionStatusOpen,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
		Tags:           []string{"tag1", "tag2"}, // Added tags
	}
	
	closedPosition := &models.Position{
		ID:             "position456",
		UserID:         "user123",
		OrderID:        "order456",
		Symbol:         "BANKNIFTY",
		Exchange:       "NSE",
		Direction:      models.PositionDirectionShort,
		EntryPrice:     1200.50,
		Quantity:       5,
		Status:         models.PositionStatusClosed,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeFuture,
	}
	
	// Set up the mock repository expectations
	mockPositionRepo.On("GetByID", "position123").Return(existingPosition, nil)
	mockPositionRepo.On("GetByID", "position456").Return(closedPosition, nil)
	mockPositionRepo.On("GetByID", "nonexistent").Return(nil, assert.AnError)
	mockPositionRepo.On("Update", mock.AnythingOfType("*models.Position")).Return(updatedPosition, nil)
	
	// Create the service with the mock repositories
	service := NewPositionService(mockPositionRepo, mockOrderRepo)
	
	// Test successful update
	result, err := service.UpdatePosition(updatedPosition)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedPosition.ID, result.ID)
	assert.Equal(t, updatedPosition.Tags, result.Tags)
	
	// Test update of closed position (should fail)
	_, err = service.UpdatePosition(closedPosition)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "closed positions cannot be updated")
	
	// Test update of non-existent position
	nonexistentPosition := &models.Position{ID: "nonexistent"}
	_, err = service.UpdatePosition(nonexistentPosition)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "position not found")
	
	// Verify that the mock repositories were called
	mockPositionRepo.AssertExpectations(t)
}

func TestClosePosition(t *testing.T) {
	// Create mock repositories
	mockPositionRepo := new(MockPositionRepository)
	mockOrderRepo := new(MockOrderRepository)
	
	// Create sample positions
	openPosition := &models.Position{
		ID:             "position123",
		UserID:         "user123",
		OrderID:        "order123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		Direction:      models.PositionDirectionLong,
		EntryPrice:     500.75,
		Quantity:       10,
		ExitQuantity:   0,
		Status:         models.PositionStatusOpen,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
	}
	
	closedPosition := &models.Position{
		ID:             "position456",
		UserID:         "user123",
		OrderID:        "order456",
		Symbol:         "BANKNIFTY",
		Exchange:       "NSE",
		Direction:      models.PositionDirectionShort,
		EntryPrice:     1200.50,
		Quantity:       5,
		ExitQuantity:   5,
		Status:         models.PositionStatusClosed,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeFuture,
	}
	
	// Set up the mock repository expectations
	mockPositionRepo.On("GetByID", "position123").Return(openPosition, nil)
	mockPositionRepo.On("GetByID", "position456").Return(closedPosition, nil)
	mockPositionRepo.On("GetByID", "nonexistent").Return(nil, assert.AnError)
	mockPositionRepo.On("Update", mock.AnythingOfType("*models.Position")).Return(func(position *models.Position) *models.Position {
		return position
	}, nil)
	
	// Create the service with the mock repositories
	service := NewPositionService(mockPositionRepo, mockOrderRepo)
	
	// Test successful full close
	result, err := service.ClosePosition("position123", 550.0, 10)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, models.PositionStatusClosed, result.Status)
	assert.Equal(t, 10, result.ExitQuantity)
	assert.Equal(t, 550.0, result.ExitPrice)
	assert.Equal(t, (550.0-500.75)*10.0, result.RealizedPnL)
	
	// Test successful partial close
	openPosition.ExitQuantity = 0 // Reset for the next test
	openPosition.Status = models.PositionStatusOpen
	result, err = service.ClosePosition("position123", 550.0, 5)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, models.PositionStatusPartial, result.Status)
	assert.Equal(t, 5, result.ExitQuantity)
	assert.Equal(t, 550.0, result.ExitPrice)
	assert.Equal(t, (550.0-500.75)*5.0, result.RealizedPnL)
	
	// Test close of already closed position
	_, err = service.ClosePosition("position456", 1100.0, 5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "position is already closed")
	
	// Test close with invalid exit price
	_, err = service.ClosePosition("position123", 0, 5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exit price must be greater than zero")
	
	// Test close with invalid exit quantity
	_, err = service.ClosePosition("position123", 550.0, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exit quantity must be greater than zero")
	
	// Test close with excessive exit quantity
	_, err = service.ClosePosition("position123", 550.0, 15)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exit quantity cannot exceed position quantity")
	
	// Test close of non-existent position
	_, err = service.ClosePosition("nonexistent", 550.0, 5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "position not found")
	
	// Verify that the mock repositories were called
	mockPositionRepo.AssertExpectations(t)
}

func TestCalculatePnL(t *testing.T) {
	// Create mock repositories
	mockPositionRepo := new(MockPositionRepository)
	mockOrderRepo := new(MockOrderRepository)
	
	// Create the service with the mock repositories
	service := NewPositionService(mockPositionRepo, mockOrderRepo)
	
	// Test long position
	longPosition := &models.Position{
		Direction:  models.PositionDirectionLong,
		EntryPrice: 100.0,
		Quantity:   10,
		Status:     models.PositionStatusOpen,
	}
	pnl, err := service.CalculatePnL(longPosition)
	assert.NoError(t, err)
	assert.Greater(t, pnl, 0.0) // Should be positive with our mock implementation
	
	// Test short position
	shortPosition := &models.Position{
		Direction:  models.PositionDirectionShort,
		EntryPrice: 100.0,
		Quantity:   10,
		Status:     models.PositionStatusOpen,
	}
	pnl, err = service.CalculatePnL(shortPosition)
	assert.NoError(t, err)
	assert.Less(t, pnl, 0.0) // Should be negative with our mock implementation
	
	// Test closed position
	closedPosition := &models.Position{
		Direction:   models.PositionDirectionLong,
		EntryPrice:  100.0,
		ExitPrice:   110.0,
		Quantity:    10,
		ExitQuantity: 10,
		Status:      models.PositionStatusClosed,
		RealizedPnL: 100.0,
	}
	pnl, err = service.CalculatePnL(closedPosition)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, pnl) // Should return the realized P&L
	
	// Test nil position
	pnl, err = service.CalculatePnL(nil)
	assert.Error(t, err)
	assert.Equal(t, 0.0, pnl)
}

func TestCalculateGreeks(t *testing.T) {
	// Create mock repositories
	mockPositionRepo := new(MockPositionRepository)
	mockOrderRepo := new(MockOrderRepository)
	
	// Create the service with the mock repositories
	service := NewPositionService(mockPositionRepo, mockOrderRepo)
	
	// Test option position
	optionPosition := &models.Position{
		InstrumentType: models.InstrumentTypeOption,
		OptionType:     models.OptionTypeCall,
		Direction:      models.PositionDirectionLong,
		Quantity:       10,
	}
	greeks, err := service.CalculateGreeks(optionPosition)
	assert.NoError(t, err)
	assert.NotNil(t, greeks)
	assert.Equal(t, 0.6*10.0, greeks.Delta)
	assert.Equal(t, 0.05*10.0, greeks.Gamma)
	assert.Equal(t, -0.1*10.0, greeks.Theta)
	assert.Equal(t, 0.2*10.0, greeks.Vega)
	
	// Test future position
	futurePosition := &models.Position{
		InstrumentType: models.InstrumentTypeFuture,
		Direction:      models.PositionDirectionLong,
		Quantity:       10,
	}
	greeks, err = service.CalculateGreeks(futurePosition)
	assert.NoError(t, err)
	assert.NotNil(t, greeks)
	assert.Equal(t, 1.0*10.0, greeks.Delta)
	assert.Equal(t, 0.0, greeks.Gamma)
	assert.Equal(t, 0.0, greeks.Theta)
	assert.Equal(t, 0.0, greeks.Vega)
	
	// Test nil position
	greeks, err = service.CalculateGreeks(nil)
	assert.Error(t, err)
	assert.Nil(t, greeks)
}

func TestCalculateExposure(t *testing.T) {
	// Create mock repositories
	mockPositionRepo := new(MockPositionRepository)
	mockOrderRepo := new(MockOrderRepository)
	
	// Create the service with the mock repositories
	service := NewPositionService(mockPositionRepo, mockOrderRepo)
	
	// Test with multiple positions
	positions := []models.Position{
		{
			Direction:  models.PositionDirectionLong,
			EntryPrice: 100.0,
			Quantity:   10,
			Status:     models.PositionStatusOpen,
		},
		{
			Direction:  models.PositionDirectionShort,
			EntryPrice: 200.0,
			Quantity:   5,
			Status:     models.PositionStatusOpen,
		},
		{
			Direction:  models.PositionDirectionLong,
			EntryPrice: 300.0,
			Quantity:   3,
			Status:     models.PositionStatusClosed, // Should be ignored
		},
	}
	
	exposure, err := service.CalculateExposure(positions)
	assert.NoError(t, err)
	// Expected: (100.0 * 10) + (200.0 * 5) = 1000 + 1000 = 2000
	assert.Equal(t, 2000.0, exposure)
	
	// Test with empty positions
	exposure, err = service.CalculateExposure([]models.Position{})
	assert.NoError(t, err)
	assert.Equal(t, 0.0, exposure)
}

func TestAggregatePositions(t *testing.T) {
	// Create mock repositories
	mockPositionRepo := new(MockPositionRepository)
	mockOrderRepo := new(MockOrderRepository)
	
	// Create the service with the mock repositories
	service := NewPositionService(mockPositionRepo, mockOrderRepo)
	
	// Test with multiple positions
	positions := []models.Position{
		{
			Symbol:        "NIFTY",
			Direction:     models.PositionDirectionLong,
			EntryPrice:    100.0,
			Quantity:      10,
			Status:        models.PositionStatusOpen,
			UnrealizedPnL: 50.0,
			Greeks: models.Greeks{
				Delta: 10.0,
				Gamma: 1.0,
				Theta: -2.0,
				Vega:  3.0,
			},
		},
		{
			Symbol:        "NIFTY",
			Direction:     models.PositionDirectionShort,
			EntryPrice:    200.0,
			Quantity:      5,
			Status:        models.PositionStatusOpen,
			UnrealizedPnL: -20.0,
			Greeks: models.Greeks{
				Delta: -5.0,
				Gamma: 0.5,
				Theta: -1.0,
				Vega:  1.5,
			},
		},
		{
			Symbol:     "BANKNIFTY",
			Direction:  models.PositionDirectionLong,
			EntryPrice: 300.0,
			Quantity:   3,
			Status:     models.PositionStatusOpen,
			UnrealizedPnL: 30.0,
			Greeks: models.Greeks{
				Delta: 3.0,
				Gamma: 0.3,
				Theta: -0.6,
				Vega:  0.9,
			},
		},
	}
	
	// Test aggregation by symbol
	aggregated, err := service.AggregatePositions(positions, "symbol")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(aggregated))
	
	// Check NIFTY aggregation
	nifty := aggregated["NIFTY"]
	assert.Equal(t, "NIFTY", nifty.Key)
	assert.Equal(t, "symbol", nifty.GroupBy)
	assert.Equal(t, 15, nifty.TotalQuantity) // 10 + 5
	assert.Equal(t, 5, nifty.NetQuantity)    // 10 - 5
	assert.Equal(t, 2000.0, nifty.TotalValue) // (100 * 10) + (200 * 5)
	assert.Equal(t, 0.0, nifty.NetValue)      // (100 * 10) - (200 * 5)
	assert.Equal(t, 30.0, nifty.PnL)          // 50 + (-20)
	assert.Equal(t, 5.0, nifty.Greeks.Delta)  // 10 + (-5)
	assert.Equal(t, 1.5, nifty.Greeks.Gamma)  // 1.0 + 0.5
	assert.Equal(t, -3.0, nifty.Greeks.Theta) // -2.0 + (-1.0)
	assert.Equal(t, 4.5, nifty.Greeks.Vega)   // 3.0 + 1.5
	assert.Equal(t, 2, nifty.PositionCount)
	
	// Test invalid groupBy
	_, err = service.AggregatePositions(positions, "invalid")
	assert.Error(t, err)
	
	// Test with empty positions
	aggregated, err = service.AggregatePositions([]models.Position{}, "symbol")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(aggregated))
}
