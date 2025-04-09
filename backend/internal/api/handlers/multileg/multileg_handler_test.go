package multileg

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trading-platform/backend/internal/models"
)

// MockMultilegService is a mock implementation of the MultilegService interface
type MockMultilegService struct {
	mock.Mock
}

func (m *MockMultilegService) CreateMultilegStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error) {
	args := m.Called(strategy)
	return args.Get(0).(*models.MultilegStrategy), args.Error(1)
}

func (m *MockMultilegService) GetMultilegStrategyByID(id string) (*models.MultilegStrategy, error) {
	args := m.Called(id)
	return args.Get(0).(*models.MultilegStrategy), args.Error(1)
}

func (m *MockMultilegService) GetMultilegStrategiesByUser(userID string) ([]models.MultilegStrategy, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.MultilegStrategy), args.Error(1)
}

func (m *MockMultilegService) GetMultilegStrategiesByPortfolio(portfolioID string) ([]models.MultilegStrategy, error) {
	args := m.Called(portfolioID)
	return args.Get(0).([]models.MultilegStrategy), args.Error(1)
}

func (m *MockMultilegService) UpdateMultilegStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error) {
	args := m.Called(strategy)
	return args.Get(0).(*models.MultilegStrategy), args.Error(1)
}

func (m *MockMultilegService) DeleteMultilegStrategy(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMultilegService) AddLeg(strategyID string, leg *models.Leg) (*models.Leg, error) {
	args := m.Called(strategyID, leg)
	return args.Get(0).(*models.Leg), args.Error(1)
}

func (m *MockMultilegService) UpdateLeg(strategyID string, leg *models.Leg) (*models.Leg, error) {
	args := m.Called(strategyID, leg)
	return args.Get(0).(*models.Leg), args.Error(1)
}

func (m *MockMultilegService) RemoveLeg(strategyID string, legID string) error {
	args := m.Called(strategyID, legID)
	return args.Error(0)
}

func (m *MockMultilegService) GetLegsByStrategy(strategyID string) ([]models.Leg, error) {
	args := m.Called(strategyID)
	return args.Get(0).([]models.Leg), args.Error(1)
}

func (m *MockMultilegService) ExecuteMultilegStrategy(strategyID string) error {
	args := m.Called(strategyID)
	return args.Error(0)
}

func (m *MockMultilegService) PauseMultilegStrategy(strategyID string) error {
	args := m.Called(strategyID)
	return args.Error(0)
}

func (m *MockMultilegService) ResumeMultilegStrategy(strategyID string) error {
	args := m.Called(strategyID)
	return args.Error(0)
}

func (m *MockMultilegService) CancelMultilegStrategy(strategyID string) error {
	args := m.Called(strategyID)
	return args.Error(0)
}

func (m *MockMultilegService) GetMultilegStrategyStatus(strategyID string) (string, error) {
	args := m.Called(strategyID)
	return args.String(0), args.Error(1)
}

func (m *MockMultilegService) GetMultilegStrategyPerformance(strategyID string) (*models.StrategyPerformance, error) {
	args := m.Called(strategyID)
	return args.Get(0).(*models.StrategyPerformance), args.Error(1)
}

// TestCreateMultilegStrategyHandler tests the CreateMultilegStrategy handler
func TestCreateMultilegStrategyHandler(t *testing.T) {
	// Create a mock service
	mockService := new(MockMultilegService)
	
	// Create the handler
	handler := NewMultilegHandler(mockService)
	
	// Create a sample strategy
	strategy := &models.MultilegStrategy{
		Name:        "Test Multileg Strategy",
		Description: "A test multileg strategy",
		UserID:      "user123",
		PortfolioID: "portfolio123",
	}
	
	// Set up the mock expectations
	mockService.On("CreateMultilegStrategy", mock.AnythingOfType("*models.MultilegStrategy")).Return(strategy, nil)
	
	// Create a request
	requestBody := `{
		"name": "Test Multileg Strategy",
		"description": "A test multileg strategy",
		"userId": "user123",
		"portfolioId": "portfolio123"
	}`
	req, err := http.NewRequest("POST", "/api/multileg", strings.NewReader(requestBody))
	assert.NoError(t, err)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.CreateMultilegStrategy(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)
	
	// Verify that the mock was called
	mockService.AssertExpectations(t)
}

// TestGetMultilegStrategyHandler tests the GetMultilegStrategy handler
func TestGetMultilegStrategyHandler(t *testing.T) {
	// Create a mock service
	mockService := new(MockMultilegService)
	
	// Create the handler
	handler := NewMultilegHandler(mockService)
	
	// Create a sample strategy
	strategy := &models.MultilegStrategy{
		ID:          "strategy123",
		Name:        "Test Multileg Strategy",
		Description: "A test multileg strategy",
		UserID:      "user123",
		PortfolioID: "portfolio123",
	}
	
	// Set up the mock expectations
	mockService.On("GetMultilegStrategyByID", "strategy123").Return(strategy, nil)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/multileg/strategy123", nil)
	assert.NoError(t, err)
	
	// Add URL parameters
	vars := map[string]string{
		"strategyId": "strategy123",
	}
	req = mux.SetURLVars(req, vars)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.GetMultilegStrategy(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify that the mock was called
	mockService.AssertExpectations(t)
}

// TestGetUserMultilegStrategiesHandler tests the GetUserMultilegStrategies handler
func TestGetUserMultilegStrategiesHandler(t *testing.T) {
	// Create a mock service
	mockService := new(MockMultilegService)
	
	// Create the handler
	handler := NewMultilegHandler(mockService)
	
	// Create sample strategies
	strategies := []models.MultilegStrategy{
		{
			ID:          "strategy1",
			Name:        "Strategy 1",
			Description: "First strategy",
			UserID:      "user123",
			PortfolioID: "portfolio123",
		},
		{
			ID:          "strategy2",
			Name:        "Strategy 2",
			Description: "Second strategy",
			UserID:      "user123",
			PortfolioID: "portfolio456",
		},
	}
	
	// Set up the mock expectations
	mockService.On("GetMultilegStrategiesByUser", "user123").Return(strategies, nil)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/users/user123/multileg", nil)
	assert.NoError(t, err)
	
	// Add URL parameters
	vars := map[string]string{
		"userId": "user123",
	}
	req = mux.SetURLVars(req, vars)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.GetUserMultilegStrategies(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify that the mock was called
	mockService.AssertExpectations(t)
}

// TestUpdateMultilegStrategyHandler tests the UpdateMultilegStrategy handler
func TestUpdateMultilegStrategyHandler(t *testing.T) {
	// Create a mock service
	mockService := new(MockMultilegService)
	
	// Create the handler
	handler := NewMultilegHandler(mockService)
	
	// Create a sample strategy
	strategy := &models.MultilegStrategy{
		ID:          "strategy123",
		Name:        "Updated Multileg Strategy",
		Description: "An updated test multileg strategy",
		UserID:      "user123",
		PortfolioID: "portfolio123",
	}
	
	// Set up the mock expectations
	mockService.On("UpdateMultilegStrategy", mock.AnythingOfType("*models.MultilegStrategy")).Return(strategy, nil)
	
	// Create a request
	requestBody := `{
		"id": "strategy123",
		"name": "Updated Multileg Strategy",
		"description": "An updated test multileg strategy",
		"userId": "user123",
		"portfolioId": "portfolio123"
	}`
	req, err := http.NewRequest("PUT", "/api/multileg/strategy123", strings.NewReader(requestBody))
	assert.NoError(t, err)
	
	// Add URL parameters
	vars := map[string]string{
		"strategyId": "strategy123",
	}
	req = mux.SetURLVars(req, vars)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.UpdateMultilegStrategy(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify that the mock was called
	mockService.AssertExpectations(t)
}

// TestDeleteMultilegStrategyHandler tests the DeleteMultilegStrategy handler
func TestDeleteMultilegStrategyHandler(t *testing.T) {
	// Create a mock service
	mockService := new(MockMultilegService)
	
	// Create the handler
	handler := NewMultilegHandler(mockService)
	
	// Set up the mock expectations
	mockService.On("DeleteMultilegStrategy", "strategy123").Return(nil)
	
	// Create a request
	req, err := http.NewRequest("DELETE", "/api/multileg/strategy123", nil)
	assert.NoError(t, err)
	
	// Add URL parameters
	vars := map[string]string{
		"strategyId": "strategy123",
	}
	req = mux.SetURLVars(req, vars)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.DeleteMultilegStrategy(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify that the mock was called
	mockService.AssertExpectations(t)
}

// TestAddLegHandler tests the AddLeg handler
func TestAddLegHandler(t *testing.T) {
	// Create a mock service
	mockService := new(MockMultilegService)
	
	// Create the handler
	handler := NewMultilegHandler(mockService)
	
	// Create a sample leg
	leg := &models.Leg{
		Symbol:        "AAPL",
		Type:          models.LegTypeBuy,
		Quantity:      10,
		ExecutionType: models.ExecutionTypeMarket,
		Sequence:      1,
	}
	
	// Set up the mock expectations
	mockService.On("AddLeg", "strategy123", mock.AnythingOfType("*models.Leg")).Return(leg, nil)
	
	// Create a request
	requestBody := `{
		"symbol": "AAPL",
		"type": "BUY",
		"quantity": 10,
		"executionType": "MARKET",
		"sequence": 1
	}`
	req, err := http.NewRequest("POST", "/api/multileg/strategy123/legs", strings.NewReader(requestBody))
	assert.NoError(t, err)
	
	// Add URL parameters
	vars := map[string]string{
		"strategyId": "strategy123",
	}
	req = mux.SetURLVars(req, vars)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.AddLeg(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)
	
	// Verify that the mock was called
	mockService.AssertExpectations(t)
}

// TestExecuteMultilegStrategyHandler tests the ExecuteMultilegStrategy handler
func TestExecuteMultilegStrategyHandler(t *testing.T) {
	// Create a mock service
	mockService := new(MockMultilegService)
	
	// Create the handler
	handler := NewMultilegHandler(mockService)
	
	// Set up the mock expectations
	mockService.On("ExecuteMultilegStrategy", "strategy123").Return(nil)
	
	// Create a request
	req, err := http.NewRequest("POST", "/api/multileg/strategy123/execute", nil)
	assert.NoError(t, err)
	
	// Add URL parameters
	vars := map[string]string{
		"strategyId": "strategy123",
	}
	req = mux.SetURLVars(req, vars)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.ExecuteMultilegStrategy(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify that the mock was called
	mockService.AssertExpectations(t)
}

// TestGetMultilegStrategyStatusHandler tests the GetMultilegStrategyStatus handler
func TestGetMultilegStrategyStatusHandler(t *testing.T) {
	// Create a mock service
	mockService := new(MockMultilegService)
	
	// Create the handler
	handler := NewMultilegHandler(mockService)
	
	// Set up the mock expectations
	mockService.On("GetMultilegStrategyStatus", "strategy123").Return("ACTIVE", nil)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/multileg/strategy123/status", nil)
	assert.NoError(t, err)
	
	// Add URL parameters
	vars := map[string]string{
		"strategyId": "strategy123",
	}
	req = mux.SetURLVars(req, vars)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.GetMultilegStrategyStatus(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify that the mock was called
	mockService.AssertExpectations(t)
}

// TestGetMultilegStrategyPerformanceHandler tests the GetMultilegStrategyPerformance handler
func TestGetMultilegStrategyPerformanceHandler(t *testing.T) {
	// Create a mock service
	mockService := new(MockMultilegService)
	
	// Create the handler
	handler := NewMultilegHandler(mockService)
	
	// Create a sample performance
	performance := &models.StrategyPerformance{
		StrategyID:    "strategy123",
		TotalPnL:      100.0,
		WinCount:      1,
		LossCount:     0,
		TotalTrades:   1,
		WinRate:       100.0,
		OrderCount:    2,
		PositionCount: 1,
	}
	
	// Set up the mock expectations
	mockService.On("GetMultilegStrategyPerformance", "strategy123").Return(performance, nil)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/multileg/strategy123/performance", nil)
	assert.NoError(t, err)
	
	// Add URL parameters
	vars := map[string]string{
		"strategyId": "strategy123",
	}
	req = mux.SetURLVars(req, vars)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.GetMultilegStrategyPerformance(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify that the mock was called
	mockService.AssertExpectations(t)
}
