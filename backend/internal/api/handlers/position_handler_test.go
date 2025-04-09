package handlers

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
	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/services/position"
)

// MockPositionService is a mock implementation of the PositionService interface
type MockPositionService struct {
	mock.Mock
}

func (m *MockPositionService) CreatePositionFromOrder(order *models.Order) (*models.Position, error) {
	args := m.Called(order)
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionService) GetPositionByID(id string) (*models.Position, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionService) GetPositions(filter models.PositionFilter, page, limit int) ([]models.Position, int, error) {
	args := m.Called(filter, page, limit)
	return args.Get(0).([]models.Position), args.Int(1), args.Error(2)
}

func (m *MockPositionService) UpdatePosition(position *models.Position) (*models.Position, error) {
	args := m.Called(position)
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionService) ClosePosition(id string, exitPrice float64, exitQuantity int) (*models.Position, error) {
	args := m.Called(id, exitPrice, exitQuantity)
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionService) CalculatePnL(position *models.Position) (float64, error) {
	args := m.Called(position)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPositionService) CalculateGreeks(position *models.Position) (*models.Greeks, error) {
	args := m.Called(position)
	return args.Get(0).(*models.Greeks), args.Error(1)
}

func (m *MockPositionService) CalculateExposure(positions []models.Position) (float64, error) {
	args := m.Called(positions)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPositionService) AggregatePositions(positions []models.Position, groupBy string) (map[string]models.AggregatedPosition, error) {
	args := m.Called(positions, groupBy)
	return args.Get(0).(map[string]models.AggregatedPosition), args.Error(1)
}

func TestCreatePositionFromOrder(t *testing.T) {
	// Create a mock position service
	mockService := new(MockPositionService)
	
	// Create a sample order
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
	}
	
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
	}
	
	// Set up the mock service expectations
	mockService.On("CreatePositionFromOrder", mock.AnythingOfType("*models.Order")).Return(position, nil)
	
	// Create the handler with the mock service
	handler := NewPositionHandler(mockService)
	
	// Create a request body
	orderJSON, _ := json.Marshal(order)
	req, err := http.NewRequest("POST", "/api/positions/create-from-order", bytes.NewBuffer(orderJSON))
	if err != nil {
		t.Fatal(err)
	}
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.CreatePositionFromOrder(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)
	
	// Parse the response
	var response models.Position
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, position.ID, response.ID)
	assert.Equal(t, position.UserID, response.UserID)
	assert.Equal(t, position.Symbol, response.Symbol)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetPosition(t *testing.T) {
	// Create a mock position service
	mockService := new(MockPositionService)
	
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
	}
	
	// Set up the mock service expectations
	mockService.On("GetPositionByID", "position123").Return(position, nil)
	
	// Create the handler with the mock service
	handler := NewPositionHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/positions/position123", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/positions/{id}", handler.GetPosition)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response models.Position
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, position.ID, response.ID)
	assert.Equal(t, position.UserID, response.UserID)
	assert.Equal(t, position.Symbol, response.Symbol)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetPositions(t *testing.T) {
	// Create a mock position service
	mockService := new(MockPositionService)
	
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
	
	// Set up the mock service expectations
	mockService.On("GetPositions", mock.AnythingOfType("models.PositionFilter"), 1, 50).Return(positions, 2, nil)
	
	// Create the handler with the mock service
	handler := NewPositionHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/positions?userId=user123", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.GetPositions(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, float64(2), response["total"])
	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(50), response["limit"])
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestUpdatePosition(t *testing.T) {
	// Create a mock position service
	mockService := new(MockPositionService)
	
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
		Tags:           []string{"tag1", "tag2"},
	}
	
	// Set up the mock service expectations
	mockService.On("UpdatePosition", mock.AnythingOfType("*models.Position")).Return(position, nil)
	
	// Create the handler with the mock service
	handler := NewPositionHandler(mockService)
	
	// Create a request body
	positionJSON, _ := json.Marshal(position)
	req, err := http.NewRequest("PUT", "/api/positions/position123", bytes.NewBuffer(positionJSON))
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/positions/{id}", handler.UpdatePosition)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response models.Position
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, position.ID, response.ID)
	assert.Equal(t, position.Tags, response.Tags)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestClosePosition(t *testing.T) {
	// Create a mock position service
	mockService := new(MockPositionService)
	
	// Create a sample position
	position := &models.Position{
		ID:             "position123",
		UserID:         "user123",
		OrderID:        "order123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		Direction:      models.PositionDirectionLong,
		EntryPrice:     500.75,
		ExitPrice:      550.0,
		Quantity:       10,
		ExitQuantity:   10,
		Status:         models.PositionStatusClosed,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
		RealizedPnL:    492.5, // (550 - 500.75) * 10
	}
	
	// Set up the mock service expectations
	mockService.On("ClosePosition", "position123", 550.0, 10).Return(position, nil)
	
	// Create the handler with the mock service
	handler := NewPositionHandler(mockService)
	
	// Create a request body
	closeParams := map[string]interface{}{
		"exitPrice":    550.0,
		"exitQuantity": 10,
	}
	paramsJSON, _ := json.Marshal(closeParams)
	req, err := http.NewRequest("POST", "/api/positions/position123/close", bytes.NewBuffer(paramsJSON))
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/positions/{id}/close", handler.ClosePosition)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response models.Position
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, position.ID, response.ID)
	assert.Equal(t, position.Status, response.Status)
	assert.Equal(t, position.ExitPrice, response.ExitPrice)
	assert.Equal(t, position.ExitQuantity, response.ExitQuantity)
	assert.Equal(t, position.RealizedPnL, response.RealizedPnL)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestCalculatePnL(t *testing.T) {
	// Create a mock position service
	mockService := new(MockPositionService)
	
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
	}
	
	// Set up the mock service expectations
	mockService.On("GetPositionByID", "position123").Return(position, nil)
	mockService.On("CalculatePnL", position).Return(100.0, nil)
	
	// Create the handler with the mock service
	handler := NewPositionHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/positions/position123/pnl", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/positions/{id}/pnl", handler.CalculatePnL)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, "position123", response["positionId"])
	assert.Equal(t, float64(100.0), response["pnl"])
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestCalculateGreeks(t *testing.T) {
	// Create a mock position service
	mockService := new(MockPositionService)
	
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
	}
	
	// Create sample Greeks
	greeks := &models.Greeks{
		Delta: 0.6,
		Gamma: 0.05,
		Theta: -0.1,
		Vega:  0.2,
	}
	
	// Set up the mock service expectations
	mockService.On("GetPositionByID", "position123").Return(position, nil)
	mockService.On("CalculateGreeks", position).Return(greeks, nil)
	
	// Create the handler with the mock service
	handler := NewPositionHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/positions/position123/greeks", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/positions/{id}/greeks", handler.CalculateGreeks)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, "position123", response["positionId"])
	greeksMap := response["greeks"].(map[string]interface{})
	assert.Equal(t, 0.6, greeksMap["delta"])
	assert.Equal(t, 0.05, greeksMap["gamma"])
	assert.Equal(t, -0.1, greeksMap["theta"])
	assert.Equal(t, 0.2, greeksMap["vega"])
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestCalculateExposure(t *testing.T) {
	// Create a mock position service
	mockService := new(MockPositionService)
	
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
	
	// Set up the mock service expectations
	mockService.On("GetPositions", mock.AnythingOfType("models.PositionFilter"), 1, 1000).Return(positions, 2, nil)
	mockService.On("CalculateExposure", positions).Return(10000.0, nil)
	
	// Create the handler with the mock service
	handler := NewPositionHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/users/user123/exposure", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/exposure", handler.CalculateExposure)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, "user123", response["userId"])
	assert.Equal(t, float64(10000.0), response["exposure"])
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestAggregatePositions(t *testing.T) {
	// Create a mock position service
	mockService := new(MockPositionService)
	
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
	
	// Create sample aggregated positions
	aggregated := map[string]models.AggregatedPosition{
		"NIFTY": {
			Key:           "NIFTY",
			GroupBy:       "symbol",
			TotalQuantity: 10,
			NetQuantity:   10,
			TotalValue:    5007.5,
			NetValue:      5007.5,
			PnL:           100.0,
			Greeks: models.Greeks{
				Delta: 6.0,
				Gamma: 0.5,
				Theta: -1.0,
				Vega:  2.0,
			},
			PositionCount: 1,
		},
		"BANKNIFTY": {
			Key:           "BANKNIFTY",
			GroupBy:       "symbol",
			TotalQuantity: 5,
			NetQuantity:   -5,
			TotalValue:    6002.5,
			NetValue:      -6002.5,
			PnL:           50.0,
			Greeks: models.Greeks{
				Delta: -5.0,
				Gamma: 0.25,
				Theta: -0.5,
				Vega:  1.0,
			},
			PositionCount: 1,
		},
	}
	
	// Set up the mock service expectations
	mockService.On("GetPositions", mock.AnythingOfType("models.PositionFilter"), 1, 1000).Return(positions, 2, nil)
	mockService.On("AggregatePositions", positions, "symbol").Return(aggregated, nil)
	
	// Create the handler with the mock service
	handler := NewPositionHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/positions/aggregate?userId=user123&groupBy=symbol", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.AggregatePositions(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response []models.AggregatedPosition
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, 2, len(response))
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetPositionsByUser(t *testing.T) {
	// Create a mock position service
	mockService := new(MockPositionService)
	
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
	
	// Set up the mock service expectations
	mockService.On("GetPositions", mock.MatchedBy(func(filter models.PositionFilter) bool {
		return filter.UserID == "user123"
	}), 1, 50).Return(positions, 2, nil)
	
	// Create the handler with the mock service
	handler := NewPositionHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/users/user123/positions", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/positions", handler.GetPositionsByUser)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, float64(2), response["total"])
	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(50), response["limit"])
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}
