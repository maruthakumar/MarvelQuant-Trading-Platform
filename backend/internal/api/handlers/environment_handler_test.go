package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"trading_platform/backend/internal/auth"
	"trading_platform/backend/internal/models"
	"trading_platform/backend/internal/services/user"
	"trading_platform/backend/internal/services/user/mocks"
)

func TestSwitchEnvironment(t *testing.T) {
	// Create mock service
	mockService := new(mocks.MockEnvironmentService)
	
	// Create handler with mock service
	handler := NewEnvironmentHandler(mockService)
	
	// Test data
	userID := "user123"
	targetEnv := models.EnvironmentSIM
	token := "new-jwt-token"
	
	// Create request body
	reqBody := `{"environment":"SIM"}`
	
	// Create request
	req := httptest.NewRequest("POST", "/api/environment/switch", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	// Create context with user ID
	ctx := auth.SetUserIDInContext(req.Context(), userID)
	req = req.WithContext(ctx)
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Set up expectations
	mockService.On("SwitchEnvironment", mock.Anything, targetEnv).Return(token, nil)
	
	// Call the handler
	handler.SwitchEnvironment(rr, req)
	
	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse response
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify response
	assert.Equal(t, token, response["token"])
	assert.Equal(t, string(targetEnv), response["environment"])
	
	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestGetEnvironmentStatus(t *testing.T) {
	// Create mock service
	mockService := new(mocks.MockEnvironmentService)
	
	// Create handler with mock service
	handler := NewEnvironmentHandler(mockService)
	
	// Test data
	status := map[string]interface{}{
		"environment": string(models.EnvironmentSIM),
		"userType": string(models.UserTypeSIM),
		"isSimulation": true,
	}
	
	// Create request
	req := httptest.NewRequest("GET", "/api/environment/status", nil)
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Set up expectations
	mockService.On("GetEnvironmentStatus", mock.Anything).Return(status)
	
	// Call the handler
	handler.GetEnvironmentStatus(rr, req)
	
	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify response
	assert.Equal(t, status["environment"], response["environment"])
	assert.Equal(t, status["userType"], response["userType"])
	assert.Equal(t, status["isSimulation"], response["isSimulation"])
	
	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestSwitchEnvironmentInvalidRequest(t *testing.T) {
	// Create mock service
	mockService := new(mocks.MockEnvironmentService)
	
	// Create handler with mock service
	handler := NewEnvironmentHandler(mockService)
	
	// Create invalid request body
	reqBody := `{"invalid":"json"`
	
	// Create request
	req := httptest.NewRequest("POST", "/api/environment/switch", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.SwitchEnvironment(rr, req)
	
	// Assertions
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSwitchEnvironmentInvalidEnvironment(t *testing.T) {
	// Create mock service
	mockService := new(mocks.MockEnvironmentService)
	
	// Create handler with mock service
	handler := NewEnvironmentHandler(mockService)
	
	// Create request with invalid environment
	reqBody := `{"environment":"INVALID"}`
	
	// Create request
	req := httptest.NewRequest("POST", "/api/environment/switch", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.SwitchEnvironment(rr, req)
	
	// Assertions
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
