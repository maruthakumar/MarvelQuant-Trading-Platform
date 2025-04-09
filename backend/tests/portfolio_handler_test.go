package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"

	"trading_platform/backend/internal/api"
	"trading_platform/backend/internal/auth"
	"trading_platform/backend/internal/models"
)

// MockPortfolioRepository is a mock implementation of PortfolioRepository
type MockPortfolioRepository struct {
	mock.Mock
}

func (m *MockPortfolioRepository) Create(portfolio *models.Portfolio) (string, error) {
	args := m.Called(portfolio)
	return args.String(0), args.Error(1)
}

func (m *MockPortfolioRepository) GetByID(id string) (*models.Portfolio, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Portfolio), args.Error(1)
}

func (m *MockPortfolioRepository) Update(portfolio *models.Portfolio) error {
	args := m.Called(portfolio)
	return args.Error(0)
}

func (m *MockPortfolioRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPortfolioRepository) Find(filter models.PortfolioFilter, page, limit int) ([]models.Portfolio, int, error) {
	args := m.Called(filter, page, limit)
	return args.Get(0).([]models.Portfolio), args.Int(1), args.Error(2)
}

// MockStrategyRepository is a mock implementation of StrategyRepository for portfolio tests
type MockStrategyRepository struct {
	mock.Mock
}

func (m *MockStrategyRepository) Create(strategy *models.Strategy) (string, error) {
	args := m.Called(strategy)
	return args.String(0), args.Error(1)
}

func (m *MockStrategyRepository) GetByID(id string) (*models.Strategy, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Strategy), args.Error(1)
}

func (m *MockStrategyRepository) Update(strategy *models.Strategy) error {
	args := m.Called(strategy)
	return args.Error(0)
}

func (m *MockStrategyRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStrategyRepository) Find(filter models.StrategyFilter, page, limit int) ([]models.Strategy, int, error) {
	args := m.Called(filter, page, limit)
	return args.Get(0).([]models.Strategy), args.Int(1), args.Error(2)
}

// TestCreatePortfolio tests the create portfolio endpoint
func TestCreatePortfolio(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Create test portfolio
	portfolio := models.Portfolio{
		Name:        "Test Portfolio",
		Description: "A test portfolio for automated trading",
		StrategyID:  "strategy123",
		Status:      models.PortfolioStatusDraft,
		Capital:     10000.0,
		Currency:    "USD",
		RiskLevel:   "medium",
	}

	// Create test strategy
	strategy := &models.Strategy{
		ID:     "strategy123",
		UserID: "user123",
		Name:   "Test Strategy",
		Active: true,
	}

	// Set up expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(strategy, nil)
	mockPortfolioRepo.On("Create", mock.AnythingOfType("*models.Portfolio")).Return("portfolio123", nil)

	// Create request with authenticated context
	body, _ := json.Marshal(portfolio)
	req, _ := http.NewRequest("POST", "/portfolios", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.CreatePortfolio).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Parse response
	var response models.Portfolio
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify portfolio data
	assert.Equal(t, "portfolio123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "Test Portfolio", response.Name)
	assert.Equal(t, "A test portfolio for automated trading", response.Description)
	assert.Equal(t, "strategy123", response.StrategyID)
	assert.Equal(t, models.PortfolioStatusDraft, response.Status)
	assert.Equal(t, 10000.0, response.Capital)
	assert.Equal(t, "USD", response.Currency)
	assert.Equal(t, "medium", response.RiskLevel)

	// Verify expectations
	mockStrategyRepo.AssertExpectations(t)
	mockPortfolioRepo.AssertExpectations(t)
}

// TestGetPortfolio tests the get portfolio endpoint
func TestGetPortfolio(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Create test portfolio
	portfolio := &models.Portfolio{
		ID:          "portfolio123",
		UserID:      "user123",
		Name:        "Test Portfolio",
		Description: "A test portfolio for automated trading",
		StrategyID:  "strategy123",
		Status:      models.PortfolioStatusActive,
		Capital:     10000.0,
		Currency:    "USD",
		RiskLevel:   "medium",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Legs:        []models.Leg{},
	}

	// Set up expectations
	mockPortfolioRepo.On("GetByID", "portfolio123").Return(portfolio, nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/portfolios/portfolio123", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/portfolios/{id}", handler.GetPortfolio).Methods("GET")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Portfolio
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify portfolio data
	assert.Equal(t, "portfolio123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "Test Portfolio", response.Name)
	assert.Equal(t, "A test portfolio for automated trading", response.Description)
	assert.Equal(t, "strategy123", response.StrategyID)
	assert.Equal(t, models.PortfolioStatusActive, response.Status)
	assert.Equal(t, 10000.0, response.Capital)
	assert.Equal(t, "USD", response.Currency)
	assert.Equal(t, "medium", response.RiskLevel)

	// Verify expectations
	mockPortfolioRepo.AssertExpectations(t)
}

// TestGetPortfolioUnauthorized tests the get portfolio endpoint with unauthorized access
func TestGetPortfolioUnauthorized(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Create test portfolio
	portfolio := &models.Portfolio{
		ID:          "portfolio123",
		UserID:      "user456", // Different user ID
		Name:        "Test Portfolio",
		Description: "A test portfolio for automated trading",
		StrategyID:  "strategy123",
		Status:      models.PortfolioStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set up expectations
	mockPortfolioRepo.On("GetByID", "portfolio123").Return(portfolio, nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/portfolios/portfolio123", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/portfolios/{id}", handler.GetPortfolio).Methods("GET")
	router.ServeHTTP(rr, req)

	// Check response - should be forbidden
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Verify expectations
	mockPortfolioRepo.AssertExpectations(t)
}

// TestUpdatePortfolio tests the update portfolio endpoint
func TestUpdatePortfolio(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Create existing portfolio
	existingPortfolio := &models.Portfolio{
		ID:          "portfolio123",
		UserID:      "user123",
		Name:        "Test Portfolio",
		Description: "A test portfolio for automated trading",
		StrategyID:  "strategy123",
		Status:      models.PortfolioStatusDraft,
		Capital:     10000.0,
		Currency:    "USD",
		RiskLevel:   "medium",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Legs:        []models.Leg{},
	}

	// Create portfolio update
	portfolioUpdate := models.Portfolio{
		Name:        "Updated Portfolio",
		Description: "An updated test portfolio",
		Capital:     15000.0,
		RiskLevel:   "high",
	}

	// Set up expectations
	mockPortfolioRepo.On("GetByID", "portfolio123").Return(existingPortfolio, nil)
	mockPortfolioRepo.On("Update", mock.AnythingOfType("*models.Portfolio")).Return(nil)

	// Create request with authenticated context
	body, _ := json.Marshal(portfolioUpdate)
	req, _ := http.NewRequest("PUT", "/portfolios/portfolio123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/portfolios/{id}", handler.UpdatePortfolio).Methods("PUT")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Portfolio
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify updated portfolio data
	assert.Equal(t, "portfolio123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "Updated Portfolio", response.Name)
	assert.Equal(t, "An updated test portfolio", response.Description)
	assert.Equal(t, "strategy123", response.StrategyID)
	assert.Equal(t, models.PortfolioStatusDraft, response.Status)
	assert.Equal(t, 15000.0, response.Capital)
	assert.Equal(t, "USD", response.Currency)
	assert.Equal(t, "high", response.RiskLevel)

	// Verify expectations
	mockPortfolioRepo.AssertExpectations(t)
}

// TestDeletePortfolio tests the delete portfolio endpoint
func TestDeletePortfolio(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Create existing portfolio
	existingPortfolio := &models.Portfolio{
		ID:          "portfolio123",
		UserID:      "user123",
		Name:        "Test Portfolio",
		Description: "A test portfolio for automated trading",
		StrategyID:  "strategy123",
		Status:      models.PortfolioStatusDraft,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set up expectations
	mockPortfolioRepo.On("GetByID", "portfolio123").Return(existingPortfolio, nil)
	mockPortfolioRepo.On("Delete", "portfolio123").Return(nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("DELETE", "/portfolios/portfolio123", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/portfolios/{id}", handler.DeletePortfolio).Methods("DELETE")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify expectations
	mockPortfolioRepo.AssertExpectations(t)
}

// TestActivatePortfolio tests the activate portfolio endpoint
func TestActivatePortfolio(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Create existing portfolio
	existingPortfolio := &models.Portfolio{
		ID:          "portfolio123",
		UserID:      "user123",
		Name:        "Test Portfolio",
		Description: "A test portfolio for automated trading",
		StrategyID:  "strategy123",
		Status:      models.PortfolioStatusDraft,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set up expectations
	mockPortfolioRepo.On("GetByID", "portfolio123").Return(existingPortfolio, nil)
	mockPortfolioRepo.On("Update", mock.AnythingOfType("*models.Portfolio")).Return(nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("POST", "/portfolios/portfolio123/activate", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/portfolios/{id}/activate", handler.ActivatePortfolio).Methods("POST")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Portfolio
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify portfolio is activated
	assert.Equal(t, "portfolio123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, models.PortfolioStatusActive, response.Status)

	// Verify expectations
	mockPortfolioRepo.AssertExpectations(t)
}

// TestDeactivatePortfolio tests the deactivate portfolio endpoint
func TestDeactivatePortfolio(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Create existing portfolio
	existingPortfolio := &models.Portfolio{
		ID:          "portfolio123",
		UserID:      "user123",
		Name:        "Test Portfolio",
		Description: "A test portfolio for automated trading",
		StrategyID:  "strategy123",
		Status:      models.PortfolioStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set up expectations
	mockPortfolioRepo.On("GetByID", "portfolio123").Return(existingPortfolio, nil)
	mockPortfolioRepo.On("Update", mock.AnythingOfType("*models.Portfolio")).Return(nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("POST", "/portfolios/portfolio123/deactivate", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/portfolios/{id}/deactivate", handler.DeactivatePortfolio).Methods("POST")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Portfolio
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify portfolio is deactivated
	assert.Equal(t, "portfolio123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, models.PortfolioStatusInactive, response.Status)

	// Verify expectations
	mockPortfolioRepo.AssertExpectations(t)
}

// TestAddLegToPortfolio tests the add leg to portfolio endpoint
func TestAddLegToPortfolio(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Create existing portfolio
	existingPortfolio := &models.Portfolio{
		ID:          "portfolio123",
		UserID:      "user123",
		Name:        "Test Portfolio",
		Description: "A test portfolio for automated trading",
		StrategyID:  "strategy123",
		Status:      models.PortfolioStatusDraft,
		Capital:     10000.0,
		Currency:    "USD",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Legs:        []models.Leg{},
	}

	// Create leg to add
	leg := models.Leg{
		Symbol:      "AAPL",
		Direction:   models.DirectionLong,
		Allocation:  0.25,
		EntryPrice:  150.0,
		StopLoss:    140.0,
		TakeProfit:  170.0,
		Description: "Apple Inc. long position",
	}

	// Set up expectations
	mockPortfolioRepo.On("GetByID", "portfolio123").Return(existingPortfolio, nil)
	mockPortfolioRepo.On("Update", mock.AnythingOfType("*models.Portfolio")).Return(nil)

	// Create request with authenticated context
	body, _ := json.Marshal(leg)
	req, _ := http.NewRequest("POST", "/portfolios/portfolio123/legs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/portfolios/{id}/legs", handler.AddLegToPortfolio).Methods("POST")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Portfolio
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify leg was added
	assert.Equal(t, 1, len(response.Legs))
	assert.Equal(t, "AAPL", response.Legs[0].Symbol)
	assert.Equal(t, models.DirectionLong, response.Legs[0].Direction)
	assert.Equal(t, 0.25, response.Legs[0].Allocation)
	assert.Equal(t, 150.0, response.Legs[0].EntryPrice)
	assert.Equal(t, 140.0, response.Legs[0].StopLoss)
	assert.Equal(t, 170.0, response.Legs[0].TakeProfit)
	assert.Equal(t, "Apple Inc. long position", response.Legs[0].Description)
	assert.Equal(t, 1, response.Legs[0].ID) // First leg should have ID 1

	// Verify expectations
	mockPortfolioRepo.AssertExpectations(t)
}

// TestUpdateLegInPortfolio tests the update leg in portfolio endpoint
func TestUpdateLegInPortfolio(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Create existing leg
	existingLeg := models.Leg{
		ID:          1,
		PortfolioID: "portfolio123",
		Symbol:      "AAPL",
		Direction:   models.DirectionLong,
		Allocation:  0.25,
		EntryPrice:  150.0,
		StopLoss:    140.0,
		TakeProfit:  170.0,
		Description: "Apple Inc. long position",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create existing portfolio with leg
	existingPortfolio := &models.Portfolio{
		ID:          "portfolio123",
		UserID:      "user123",
		Name:        "Test Portfolio",
		Description: "A test portfolio for automated trading",
		StrategyID:  "strategy123",
		Status:      models.PortfolioStatusDraft,
		Capital:     10000.0,
		Currency:    "USD",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Legs:        []models.Leg{existingLeg},
	}

	// Create leg update
	legUpdate := models.Leg{
		Symbol:      "AAPL",
		Direction:   models.DirectionLong,
		Allocation:  0.30,
		EntryPrice:  155.0,
		StopLoss:    145.0,
		TakeProfit:  175.0,
		Description: "Updated Apple Inc. long position",
	}

	// Set up expectations
	mockPortfolioRepo.On("GetByID", "portfolio123").Return(existingPortfolio, nil)
	mockPortfolioRepo.On("Update", mock.AnythingOfType("*models.Portfolio")).Return(nil)

	// Create request with authenticated context
	body, _ := json.Marshal(legUpdate)
	req, _ := http.NewRequest("PUT", "/portfolios/portfolio123/legs/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/portfolios/{id}/legs/{legId}", handler.UpdateLegInPortfolio).Methods("PUT")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Portfolio
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify leg was updated
	assert.Equal(t, 1, len(response.Legs))
	assert.Equal(t, 1, response.Legs[0].ID)
	assert.Equal(t, "AAPL", response.Legs[0].Symbol)
	assert.Equal(t, models.DirectionLong, response.Legs[0].Direction)
	assert.Equal(t, 0.30, response.Legs[0].Allocation)
	assert.Equal(t, 155.0, response.Legs[0].EntryPrice)
	assert.Equal(t, 145.0, response.Legs[0].StopLoss)
	assert.Equal(t, 175.0, response.Legs[0].TakeProfit)
	assert.Equal(t, "Updated Apple Inc. long position", response.Legs[0].Description)

	// Verify expectations
	mockPortfolioRepo.AssertExpectations(t)
}

// TestRemoveLegFromPortfolio tests the remove leg from portfolio endpoint
func TestRemoveLegFromPortfolio(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Create existing legs
	leg1 := models.Leg{
		ID:          1,
		PortfolioID: "portfolio123",
		Symbol:      "AAPL",
		Direction:   models.DirectionLong,
		Allocation:  0.25,
		EntryPrice:  150.0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	leg2 := models.Leg{
		ID:          2,
		PortfolioID: "portfolio123",
		Symbol:      "MSFT",
		Direction:   models.DirectionLong,
		Allocation:  0.25,
		EntryPrice:  250.0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create existing portfolio with legs
	existingPortfolio := &models.Portfolio{
		ID:          "portfolio123",
		UserID:      "user123",
		Name:        "Test Portfolio",
		Description: "A test portfolio for automated trading",
		StrategyID:  "strategy123",
		Status:      models.PortfolioStatusDraft,
		Capital:     10000.0,
		Currency:    "USD",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Legs:        []models.Leg{leg1, leg2},
	}

	// Set up expectations
	mockPortfolioRepo.On("GetByID", "portfolio123").Return(existingPortfolio, nil)
	mockPortfolioRepo.On("Update", mock.AnythingOfType("*models.Portfolio")).Return(nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("DELETE", "/portfolios/portfolio123/legs/1", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/portfolios/{id}/legs/{legId}", handler.RemoveLegFromPortfolio).Methods("DELETE")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Portfolio
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify leg was removed
	assert.Equal(t, 1, len(response.Legs))
	assert.Equal(t, 2, response.Legs[0].ID) // Only leg 2 should remain
	assert.Equal(t, "MSFT", response.Legs[0].Symbol)

	// Verify expectations
	mockPortfolioRepo.AssertExpectations(t)
}

// TestGetPortfolios tests the get portfolios endpoint with filtering and pagination
func TestGetPortfolios(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Create test portfolios
	portfolios := []models.Portfolio{
		{
			ID:          "portfolio1",
			UserID:      "user123",
			Name:        "Growth Portfolio",
			Description: "A growth-focused portfolio",
			StrategyID:  "strategy1",
			Status:      models.PortfolioStatusActive,
			Capital:     10000.0,
			Currency:    "USD",
			RiskLevel:   "high",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "portfolio2",
			UserID:      "user123",
			Name:        "Income Portfolio",
			Description: "An income-focused portfolio",
			StrategyID:  "strategy2",
			Status:      models.PortfolioStatusInactive,
			Capital:     20000.0,
			Currency:    "USD",
			RiskLevel:   "low",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Set up expectations
	mockPortfolioRepo.On("Find", mock.AnythingOfType("models.PortfolioFilter"), 1, 20).Return(portfolios, 2, nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/portfolios?status=active&riskLevel=high", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.GetPortfolios).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify pagination data
	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(20), response["limit"])
	assert.Equal(t, float64(2), response["total"])
	assert.Equal(t, float64(1), response["totalPages"])

	// Verify portfolios data
	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))

	// Verify expectations
	mockPortfolioRepo.AssertExpectations(t)
}

// TestPortfolioNotFound tests the get portfolio endpoint with non-existent portfolio
func TestPortfolioNotFound(t *testing.T) {
	// Create mock repositories
	mockPortfolioRepo := new(MockPortfolioRepository)
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewPortfolioHandler(mockPortfolioRepo, mockStrategyRepo)

	// Set up expectations
	mockPortfolioRepo.On("GetByID", "nonexistent").Return(nil, mongo.ErrNoDocuments)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/portfolios/nonexistent", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/portfolios/{id}", handler.GetPortfolio).Methods("GET")
	router.ServeHTTP(rr, req)

	// Check response - should be not found
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Verify expectations
	mockPortfolioRepo.AssertExpectations(t)
}
