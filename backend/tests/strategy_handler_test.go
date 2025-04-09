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

// MockStrategyRepository is a mock implementation of StrategyRepository
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

// TestCreateStrategy tests the create strategy endpoint
func TestCreateStrategy(t *testing.T) {
	// Create mock repository
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewStrategyHandler(mockStrategyRepo)

	// Create test strategy
	strategy := models.Strategy{
		Name:        "Test Strategy",
		Description: "A test strategy for automated trading",
		Type:        "momentum",
		Symbol:      "AAPL",
		ProductType: "equity",
		Parameters: map[string]interface{}{
			"lookbackPeriod": 14,
			"threshold":      0.05,
		},
		Active: true,
		Tags:   []string{"test", "momentum", "equity"},
	}

	// Set up expectations
	mockStrategyRepo.On("Create", mock.AnythingOfType("*models.Strategy")).Return("strategy123", nil)

	// Create request with authenticated context
	body, _ := json.Marshal(strategy)
	req, _ := http.NewRequest("POST", "/strategies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.CreateStrategy).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Parse response
	var response models.Strategy
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify strategy data
	assert.Equal(t, "strategy123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "Test Strategy", response.Name)
	assert.Equal(t, "A test strategy for automated trading", response.Description)
	assert.Equal(t, "momentum", response.Type)
	assert.Equal(t, "AAPL", response.Symbol)
	assert.Equal(t, "equity", response.ProductType)
	assert.Equal(t, true, response.Active)
	assert.Equal(t, []string{"test", "momentum", "equity"}, response.Tags)

	// Verify expectations
	mockStrategyRepo.AssertExpectations(t)
}

// TestGetStrategy tests the get strategy endpoint
func TestGetStrategy(t *testing.T) {
	// Create mock repository
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewStrategyHandler(mockStrategyRepo)

	// Create test strategy
	strategy := &models.Strategy{
		ID:          "strategy123",
		UserID:      "user123",
		Name:        "Test Strategy",
		Description: "A test strategy for automated trading",
		Type:        "momentum",
		Symbol:      "AAPL",
		ProductType: "equity",
		Parameters: map[string]interface{}{
			"lookbackPeriod": 14,
			"threshold":      0.05,
		},
		Active:    true,
		Tags:      []string{"test", "momentum", "equity"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(strategy, nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/strategies/strategy123", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/strategies/{id}", handler.GetStrategy).Methods("GET")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Strategy
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify strategy data
	assert.Equal(t, "strategy123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "Test Strategy", response.Name)
	assert.Equal(t, "A test strategy for automated trading", response.Description)
	assert.Equal(t, "momentum", response.Type)
	assert.Equal(t, "AAPL", response.Symbol)
	assert.Equal(t, "equity", response.ProductType)
	assert.Equal(t, true, response.Active)
	assert.Equal(t, []string{"test", "momentum", "equity"}, response.Tags)

	// Verify expectations
	mockStrategyRepo.AssertExpectations(t)
}

// TestGetStrategyUnauthorized tests the get strategy endpoint with unauthorized access
func TestGetStrategyUnauthorized(t *testing.T) {
	// Create mock repository
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewStrategyHandler(mockStrategyRepo)

	// Create test strategy
	strategy := &models.Strategy{
		ID:          "strategy123",
		UserID:      "user456", // Different user ID
		Name:        "Test Strategy",
		Description: "A test strategy for automated trading",
		Type:        "momentum",
		Symbol:      "AAPL",
		ProductType: "equity",
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set up expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(strategy, nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/strategies/strategy123", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/strategies/{id}", handler.GetStrategy).Methods("GET")
	router.ServeHTTP(rr, req)

	// Check response - should be forbidden
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Verify expectations
	mockStrategyRepo.AssertExpectations(t)
}

// TestUpdateStrategy tests the update strategy endpoint
func TestUpdateStrategy(t *testing.T) {
	// Create mock repository
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewStrategyHandler(mockStrategyRepo)

	// Create existing strategy
	existingStrategy := &models.Strategy{
		ID:          "strategy123",
		UserID:      "user123",
		Name:        "Test Strategy",
		Description: "A test strategy for automated trading",
		Type:        "momentum",
		Symbol:      "AAPL",
		ProductType: "equity",
		Parameters: map[string]interface{}{
			"lookbackPeriod": 14,
			"threshold":      0.05,
		},
		Active:    true,
		Tags:      []string{"test", "momentum", "equity"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create strategy update
	strategyUpdate := models.Strategy{
		Name:        "Updated Strategy",
		Description: "An updated test strategy",
		Parameters: map[string]interface{}{
			"lookbackPeriod": 21,
			"threshold":      0.03,
		},
		Tags: []string{"test", "updated", "equity"},
	}

	// Set up expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(existingStrategy, nil)
	mockStrategyRepo.On("Update", mock.AnythingOfType("*models.Strategy")).Return(nil)

	// Create request with authenticated context
	body, _ := json.Marshal(strategyUpdate)
	req, _ := http.NewRequest("PUT", "/strategies/strategy123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/strategies/{id}", handler.UpdateStrategy).Methods("PUT")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Strategy
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify updated strategy data
	assert.Equal(t, "strategy123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "Updated Strategy", response.Name)
	assert.Equal(t, "An updated test strategy", response.Description)
	assert.Equal(t, "momentum", response.Type)
	assert.Equal(t, "AAPL", response.Symbol)
	assert.Equal(t, "equity", response.ProductType)
	assert.Equal(t, true, response.Active)
	assert.Equal(t, []string{"test", "updated", "equity"}, response.Tags)

	// Verify expectations
	mockStrategyRepo.AssertExpectations(t)
}

// TestDeleteStrategy tests the delete strategy endpoint
func TestDeleteStrategy(t *testing.T) {
	// Create mock repository
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewStrategyHandler(mockStrategyRepo)

	// Create existing strategy
	existingStrategy := &models.Strategy{
		ID:          "strategy123",
		UserID:      "user123",
		Name:        "Test Strategy",
		Description: "A test strategy for automated trading",
		Type:        "momentum",
		Symbol:      "AAPL",
		ProductType: "equity",
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set up expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(existingStrategy, nil)
	mockStrategyRepo.On("Delete", "strategy123").Return(nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("DELETE", "/strategies/strategy123", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/strategies/{id}", handler.DeleteStrategy).Methods("DELETE")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify expectations
	mockStrategyRepo.AssertExpectations(t)
}

// TestActivateStrategy tests the activate strategy endpoint
func TestActivateStrategy(t *testing.T) {
	// Create mock repository
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewStrategyHandler(mockStrategyRepo)

	// Create existing strategy
	existingStrategy := &models.Strategy{
		ID:          "strategy123",
		UserID:      "user123",
		Name:        "Test Strategy",
		Description: "A test strategy for automated trading",
		Type:        "momentum",
		Symbol:      "AAPL",
		ProductType: "equity",
		Active:      false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set up expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(existingStrategy, nil)
	mockStrategyRepo.On("Update", mock.AnythingOfType("*models.Strategy")).Return(nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("POST", "/strategies/strategy123/activate", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/strategies/{id}/activate", handler.ActivateStrategy).Methods("POST")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Strategy
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify strategy is activated
	assert.Equal(t, "strategy123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, true, response.Active)

	// Verify expectations
	mockStrategyRepo.AssertExpectations(t)
}

// TestDeactivateStrategy tests the deactivate strategy endpoint
func TestDeactivateStrategy(t *testing.T) {
	// Create mock repository
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewStrategyHandler(mockStrategyRepo)

	// Create existing strategy
	existingStrategy := &models.Strategy{
		ID:          "strategy123",
		UserID:      "user123",
		Name:        "Test Strategy",
		Description: "A test strategy for automated trading",
		Type:        "momentum",
		Symbol:      "AAPL",
		ProductType: "equity",
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set up expectations
	mockStrategyRepo.On("GetByID", "strategy123").Return(existingStrategy, nil)
	mockStrategyRepo.On("Update", mock.AnythingOfType("*models.Strategy")).Return(nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("POST", "/strategies/strategy123/deactivate", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/strategies/{id}/deactivate", handler.DeactivateStrategy).Methods("POST")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Strategy
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify strategy is deactivated
	assert.Equal(t, "strategy123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, false, response.Active)

	// Verify expectations
	mockStrategyRepo.AssertExpectations(t)
}

// TestGetStrategies tests the get strategies endpoint with filtering and pagination
func TestGetStrategies(t *testing.T) {
	// Create mock repository
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewStrategyHandler(mockStrategyRepo)

	// Create test strategies
	strategies := []models.Strategy{
		{
			ID:          "strategy1",
			UserID:      "user123",
			Name:        "Momentum Strategy",
			Description: "A momentum-based trading strategy",
			Type:        "momentum",
			Symbol:      "AAPL",
			ProductType: "equity",
			Active:      true,
			Tags:        []string{"momentum", "equity"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "strategy2",
			UserID:      "user123",
			Name:        "Mean Reversion Strategy",
			Description: "A mean reversion trading strategy",
			Type:        "mean-reversion",
			Symbol:      "MSFT",
			ProductType: "equity",
			Active:      false,
			Tags:        []string{"mean-reversion", "equity"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Set up expectations
	mockStrategyRepo.On("Find", mock.AnythingOfType("models.StrategyFilter"), 1, 20).Return(strategies, 2, nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/strategies?type=momentum&active=true", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.GetStrategies).ServeHTTP(rr, req)

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

	// Verify strategies data
	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))

	// Verify expectations
	mockStrategyRepo.AssertExpectations(t)
}

// TestStrategyNotFound tests the get strategy endpoint with non-existent strategy
func TestStrategyNotFound(t *testing.T) {
	// Create mock repository
	mockStrategyRepo := new(MockStrategyRepository)

	// Create handler
	handler := api.NewStrategyHandler(mockStrategyRepo)

	// Set up expectations
	mockStrategyRepo.On("GetByID", "nonexistent").Return(nil, mongo.ErrNoDocuments)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/strategies/nonexistent", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/strategies/{id}", handler.GetStrategy).Methods("GET")
	router.ServeHTTP(rr, req)

	// Check response - should be not found
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Verify expectations
	mockStrategyRepo.AssertExpectations(t)
}
