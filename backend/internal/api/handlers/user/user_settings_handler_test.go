package user

import (
	"net/http/httptest"
	"testing"
	"encoding/json"
	"net/http"
	"bytes"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trading-platform/backend/internal/models"
)

// MockUserService is a mock implementation of the UserService interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserSettings(userID string) (*models.UserSettings, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.UserSettings), args.Error(1)
}

func (m *MockUserService) UpdateUserSettings(userID string, settings *models.UserSettings) (*models.UserSettings, error) {
	args := m.Called(userID, settings)
	return args.Get(0).(*models.UserSettings), args.Error(1)
}

func (m *MockUserService) GetUserPreferences(userID string) (*models.UserPreferences, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.UserPreferences), args.Error(1)
}

func (m *MockUserService) UpdateUserPreferences(userID string, preferences *models.UserPreferences) (*models.UserPreferences, error) {
	args := m.Called(userID, preferences)
	return args.Get(0).(*models.UserPreferences), args.Error(1)
}

func (m *MockUserService) GetUserTheme(userID string) (*models.UserTheme, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.UserTheme), args.Error(1)
}

func (m *MockUserService) UpdateUserTheme(userID string, theme *models.UserTheme) (*models.UserTheme, error) {
	args := m.Called(userID, theme)
	return args.Get(0).(*models.UserTheme), args.Error(1)
}

func (m *MockUserService) GetUserLayout(userID string, layoutName string) (*models.UserLayout, error) {
	args := m.Called(userID, layoutName)
	return args.Get(0).(*models.UserLayout), args.Error(1)
}

func (m *MockUserService) GetAllUserLayouts(userID string) ([]models.UserLayout, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.UserLayout), args.Error(1)
}

func (m *MockUserService) SaveUserLayout(userID string, layout *models.UserLayout) (*models.UserLayout, error) {
	args := m.Called(userID, layout)
	return args.Get(0).(*models.UserLayout), args.Error(1)
}

func (m *MockUserService) DeleteUserLayout(userID string, layoutName string) error {
	args := m.Called(userID, layoutName)
	return args.Error(0)
}

func (m *MockUserService) GetUserApiKeys(userID string) ([]models.UserApiKey, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.UserApiKey), args.Error(1)
}

func (m *MockUserService) AddUserApiKey(userID string, apiKey *models.UserApiKey) (*models.UserApiKey, error) {
	args := m.Called(userID, apiKey)
	return args.Get(0).(*models.UserApiKey), args.Error(1)
}

func (m *MockUserService) UpdateUserApiKey(userID string, keyID string, apiKey *models.UserApiKey) (*models.UserApiKey, error) {
	args := m.Called(userID, keyID, apiKey)
	return args.Get(0).(*models.UserApiKey), args.Error(1)
}

func (m *MockUserService) DeleteUserApiKey(userID string, keyID string) error {
	args := m.Called(userID, keyID)
	return args.Error(0)
}

func (m *MockUserService) GetUserNotificationSettings(userID string) (*models.UserNotificationSettings, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.UserNotificationSettings), args.Error(1)
}

func (m *MockUserService) UpdateUserNotificationSettings(userID string, settings *models.UserNotificationSettings) (*models.UserNotificationSettings, error) {
	args := m.Called(userID, settings)
	return args.Get(0).(*models.UserNotificationSettings), args.Error(1)
}

func TestGetUserSettings(t *testing.T) {
	// Create a mock user service
	mockService := new(MockUserService)
	
	// Create sample user settings
	settings := &models.UserSettings{
		UserID:    "user123",
		Language:  "en",
		TimeZone:  "UTC",
		DateFormat: "MM/DD/YYYY",
		TimeFormat: "HH:mm:ss",
		Currency:  "USD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Set up the mock service expectations
	mockService.On("GetUserSettings", "user123").Return(settings, nil)
	
	// Create the handler with the mock service
	handler := NewUserSettingsHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/users/user123/settings", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/settings", handler.GetUserSettings)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response models.UserSettings
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, settings.UserID, response.UserID)
	assert.Equal(t, settings.Language, response.Language)
	assert.Equal(t, settings.TimeZone, response.TimeZone)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestUpdateUserSettings(t *testing.T) {
	// Create a mock user service
	mockService := new(MockUserService)
	
	// Create sample user settings
	settings := &models.UserSettings{
		UserID:    "user123",
		Language:  "fr",
		TimeZone:  "Europe/Paris",
		DateFormat: "DD/MM/YYYY",
		TimeFormat: "HH:mm",
		Currency:  "EUR",
		UpdatedAt: time.Now(),
	}
	
	// Set up the mock service expectations
	mockService.On("UpdateUserSettings", "user123", mock.AnythingOfType("*models.UserSettings")).Return(settings, nil)
	
	// Create the handler with the mock service
	handler := NewUserSettingsHandler(mockService)
	
	// Create a request body
	settingsJSON, _ := json.Marshal(settings)
	req, err := http.NewRequest("PUT", "/api/users/user123/settings", bytes.NewBuffer(settingsJSON))
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/settings", handler.UpdateUserSettings)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response models.UserSettings
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, settings.UserID, response.UserID)
	assert.Equal(t, settings.Language, response.Language)
	assert.Equal(t, settings.TimeZone, response.TimeZone)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetUserPreferences(t *testing.T) {
	// Create a mock user service
	mockService := new(MockUserService)
	
	// Create sample user preferences
	preferences := &models.UserPreferences{
		UserID:                "user123",
		DefaultOrderQuantity:  10,
		DefaultProductType:    models.ProductTypeMIS,
		DefaultExchange:       "NSE",
		ShowConfirmationDialog: true,
		DefaultInstrumentType: models.InstrumentTypeOption,
		DefaultSymbols:        []string{"NIFTY", "BANKNIFTY"},
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
	
	// Set up the mock service expectations
	mockService.On("GetUserPreferences", "user123").Return(preferences, nil)
	
	// Create the handler with the mock service
	handler := NewUserSettingsHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/users/user123/preferences", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/preferences", handler.GetUserPreferences)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response models.UserPreferences
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, preferences.UserID, response.UserID)
	assert.Equal(t, preferences.DefaultOrderQuantity, response.DefaultOrderQuantity)
	assert.Equal(t, preferences.DefaultExchange, response.DefaultExchange)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetUserTheme(t *testing.T) {
	// Create a mock user service
	mockService := new(MockUserService)
	
	// Create sample user theme
	theme := &models.UserTheme{
		UserID:         "user123",
		ThemeMode:      "dark",
		PrimaryColor:   "#1976d2",
		SecondaryColor: "#dc004e",
		ChartColors:    []string{"#ff0000", "#00ff00", "#0000ff"},
		FontSize:       "medium",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	
	// Set up the mock service expectations
	mockService.On("GetUserTheme", "user123").Return(theme, nil)
	
	// Create the handler with the mock service
	handler := NewUserSettingsHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/users/user123/theme", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/theme", handler.GetUserTheme)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response models.UserTheme
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, theme.UserID, response.UserID)
	assert.Equal(t, theme.ThemeMode, response.ThemeMode)
	assert.Equal(t, theme.PrimaryColor, response.PrimaryColor)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetUserLayout(t *testing.T) {
	// Create a mock user service
	mockService := new(MockUserService)
	
	// Create sample user layout
	layout := &models.UserLayout{
		UserID:    "user123",
		Name:      "default",
		Type:      "grid",
		IsDefault: true,
		Layout: map[string]interface{}{
			"columns": 12,
			"rows":    6,
			"widgets": []interface{}{
				map[string]interface{}{
					"id":   "widget1",
					"type": "chart",
					"x":    0,
					"y":    0,
					"w":    6,
					"h":    3,
				},
				map[string]interface{}{
					"id":   "widget2",
					"type": "orderbook",
					"x":    6,
					"y":    0,
					"w":    6,
					"h":    3,
				},
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Set up the mock service expectations
	mockService.On("GetUserLayout", "user123", "default").Return(layout, nil)
	
	// Create the handler with the mock service
	handler := NewUserSettingsHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/users/user123/layouts/default", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/layouts/{layoutName}", handler.GetUserLayout)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response models.UserLayout
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, layout.UserID, response.UserID)
	assert.Equal(t, layout.Name, response.Name)
	assert.Equal(t, layout.Type, response.Type)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetAllUserLayouts(t *testing.T) {
	// Create a mock user service
	mockService := new(MockUserService)
	
	// Create sample user layouts
	layouts := []models.UserLayout{
		{
			UserID:    "user123",
			Name:      "default",
			Type:      "grid",
			IsDefault: true,
			Layout: map[string]interface{}{
				"columns": 12,
				"rows":    6,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID:    "user123",
			Name:      "custom",
			Type:      "flex",
			IsDefault: false,
			Layout: map[string]interface{}{
				"direction": "row",
				"wrap":      "wrap",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	
	// Set up the mock service expectations
	mockService.On("GetAllUserLayouts", "user123").Return(layouts, nil)
	
	// Create the handler with the mock service
	handler := NewUserSettingsHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/users/user123/layouts", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/layouts", handler.GetAllUserLayouts)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response []models.UserLayout
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, 2, len(response))
	assert.Equal(t, layouts[0].UserID, response[0].UserID)
	assert.Equal(t, layouts[0].Name, response[0].Name)
	assert.Equal(t, layouts[1].UserID, response[1].UserID)
	assert.Equal(t, layouts[1].Name, response[1].Name)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetUserApiKeys(t *testing.T) {
	// Create a mock user service
	mockService := new(MockUserService)
	
	// Create sample user API keys
	apiKeys := []models.UserApiKey{
		{
			ID:          "key1",
			UserID:      "user123",
			Name:        "Zerodha",
			ApiKey:      "abc123",
			ApiSecret:   "xyz789",
			Broker:      "zerodha",
			IsActive:    true,
			Permissions: []string{"trade", "data"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "key2",
			UserID:      "user123",
			Name:        "ICICI Direct",
			ApiKey:      "def456",
			ApiSecret:   "uvw321",
			Broker:      "icici",
			IsActive:    true,
			Permissions: []string{"data"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	
	// Set up the mock service expectations
	mockService.On("GetUserApiKeys", "user123").Return(apiKeys, nil)
	
	// Create the handler with the mock service
	handler := NewUserSettingsHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/users/user123/apikeys", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/apikeys", handler.GetUserApiKeys)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response []models.UserApiKey
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, 2, len(response))
	assert.Equal(t, apiKeys[0].ID, response[0].ID)
	assert.Equal(t, apiKeys[0].Name, response[0].Name)
	assert.Equal(t, "********", response[0].ApiSecret) // Should be masked
	assert.Equal(t, apiKeys[1].ID, response[1].ID)
	assert.Equal(t, apiKeys[1].Name, response[1].Name)
	assert.Equal(t, "********", response[1].ApiSecret) // Should be masked
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestAddUserApiKey(t *testing.T) {
	// Create a mock user service
	mockService := new(MockUserService)
	
	// Create sample user API key
	apiKey := &models.UserApiKey{
		UserID:      "user123",
		Name:        "Zerodha",
		ApiKey:      "abc123",
		ApiSecret:   "xyz789",
		Broker:      "zerodha",
		IsActive:    true,
		Permissions: []string{"trade", "data"},
	}
	
	// Create expected result with ID
	resultApiKey := &models.UserApiKey{
		ID:          "key1",
		UserID:      "user123",
		Name:        "Zerodha",
		ApiKey:      "abc123",
		ApiSecret:   "xyz789",
		Broker:      "zerodha",
		IsActive:    true,
		Permissions: []string{"trade", "data"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Set up the mock service expectations
	mockService.On("AddUserApiKey", "user123", mock.AnythingOfType("*models.UserApiKey")).Return(resultApiKey, nil)
	
	// Create the handler with the mock service
	handler := NewUserSettingsHandler(mockService)
	
	// Create a request body
	apiKeyJSON, _ := json.Marshal(apiKey)
	req, err := http.NewRequest("POST", "/api/users/user123/apikeys", bytes.NewBuffer(apiKeyJSON))
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/apikeys", handler.AddUserApiKey)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)
	
	// Parse the response
	var response models.UserApiKey
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, resultApiKey.ID, response.ID)
	assert.Equal(t, resultApiKey.Name, response.Name)
	assert.Equal(t, "********", response.ApiSecret) // Should be masked
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetUserNotificationSettings(t *testing.T) {
	// Create a mock user service
	mockService := new(MockUserService)
	
	// Create sample user notification settings
	settings := &models.UserNotificationSettings{
		UserID:                 "user123",
		EnableEmailNotifications: true,
		EnablePushNotifications:  true,
		OrderExecutionAlerts:     true,
		PriceAlerts:              true,
		MarginCallAlerts:         true,
		NewsAlerts:               false,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
	
	// Set up the mock service expectations
	mockService.On("GetUserNotificationSettings", "user123").Return(settings, nil)
	
	// Create the handler with the mock service
	handler := NewUserSettingsHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/users/user123/notifications", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/notifications", handler.GetUserNotificationSettings)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response models.UserNotificationSettings
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, settings.UserID, response.UserID)
	assert.Equal(t, settings.EnableEmailNotifications, response.EnableEmailNotifications)
	assert.Equal(t, settings.EnablePushNotifications, response.EnablePushNotifications)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}
