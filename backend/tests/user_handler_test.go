package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"trading_platform/backend/internal/api"
	"trading_platform/backend/internal/database"
	"trading_platform/backend/internal/models"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) Find(filter models.UserFilter, page, limit int) ([]models.User, int, error) {
	args := m.Called(filter, page, limit)
	return args.Get(0).([]models.User), args.Int(1), args.Error(2)
}

// MockUserPreferenceRepository is a mock implementation of UserPreferenceRepository
type MockUserPreferenceRepository struct {
	mock.Mock
}

func (m *MockUserPreferenceRepository) Create(preferences *models.UserPreferences) (string, error) {
	args := m.Called(preferences)
	return args.String(0), args.Error(1)
}

func (m *MockUserPreferenceRepository) GetByID(id string) (*models.UserPreferences, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPreferences), args.Error(1)
}

func (m *MockUserPreferenceRepository) GetByUserID(userID string) (*models.UserPreferences, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPreferences), args.Error(1)
}

func (m *MockUserPreferenceRepository) Update(preferences *models.UserPreferences) error {
	args := m.Called(preferences)
	return args.Error(0)
}

// MockAPIKeyRepository is a mock implementation of APIKeyRepository
type MockAPIKeyRepository struct {
	mock.Mock
}

func (m *MockAPIKeyRepository) Create(apiKey *models.APIKey) (string, error) {
	args := m.Called(apiKey)
	return args.String(0), args.Error(1)
}

func (m *MockAPIKeyRepository) GetByID(id string) (*models.APIKey, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.APIKey), args.Error(1)
}

func (m *MockAPIKeyRepository) GetByUserID(userID string) ([]*models.APIKey, error) {
	args := m.Called(userID)
	return args.Get(0).([]*models.APIKey), args.Error(1)
}

func (m *MockAPIKeyRepository) Update(apiKey *models.APIKey) error {
	args := m.Called(apiKey)
	return args.Error(0)
}

func (m *MockAPIKeyRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestUserRegistration tests the user registration endpoint
func TestUserRegistration(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockPreferenceRepo := new(MockUserPreferenceRepository)
	mockAPIKeyRepo := new(MockAPIKeyRepository)

	// Create handler
	handler := api.NewUserHandler(mockUserRepo, mockPreferenceRepo, mockAPIKeyRepo)

	// Create test user
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Set up expectations
	mockUserRepo.On("GetByUsername", "testuser").Return(nil, database.ErrNotFound)
	mockUserRepo.On("GetByEmail", "test@example.com").Return(nil, database.ErrNotFound)
	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return("user123", nil)
	mockPreferenceRepo.On("Create", mock.AnythingOfType("*models.UserPreferences")).Return("pref123", nil)

	// Create request
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.Register).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Verify expectations
	mockUserRepo.AssertExpectations(t)
	mockPreferenceRepo.AssertExpectations(t)
}

// TestUserLogin tests the user login endpoint
func TestUserLogin(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockPreferenceRepo := new(MockUserPreferenceRepository)
	mockAPIKeyRepo := new(MockAPIKeyRepository)

	// Create handler
	handler := api.NewUserHandler(mockUserRepo, mockPreferenceRepo, mockAPIKeyRepo)

	// Create test user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		ID:       "user123",
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     models.RoleUser,
		Active:   true,
	}

	// Set up expectations
	mockUserRepo.On("GetByUsername", "testuser").Return(user, nil)
	mockUserRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

	// Create login request
	loginRequest := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	body, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.Login).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify token exists
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "refreshToken")
	assert.Contains(t, response, "user")

	// Verify expectations
	mockUserRepo.AssertExpectations(t)
}

// TestGetProfile tests the get profile endpoint
func TestGetProfile(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockPreferenceRepo := new(MockUserPreferenceRepository)
	mockAPIKeyRepo := new(MockAPIKeyRepository)

	// Create handler
	handler := api.NewUserHandler(mockUserRepo, mockPreferenceRepo, mockAPIKeyRepo)

	// Create test user
	user := &models.User{
		ID:        "user123",
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      models.RoleUser,
		Active:    true,
	}

	// Set up expectations
	mockUserRepo.On("GetByID", "user123").Return(user, nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/users/profile", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.GetProfile).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.User
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify user data
	assert.Equal(t, "user123", response.ID)
	assert.Equal(t, "testuser", response.Username)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "Test", response.FirstName)
	assert.Equal(t, "User", response.LastName)
	assert.Equal(t, "", response.Password) // Password should be removed

	// Verify expectations
	mockUserRepo.AssertExpectations(t)
}

// TestUpdateProfile tests the update profile endpoint
func TestUpdateProfile(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockPreferenceRepo := new(MockUserPreferenceRepository)
	mockAPIKeyRepo := new(MockAPIKeyRepository)

	// Create handler
	handler := api.NewUserHandler(mockUserRepo, mockPreferenceRepo, mockAPIKeyRepo)

	// Create existing user
	existingUser := &models.User{
		ID:        "user123",
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      models.RoleUser,
		Active:    true,
	}

	// Create profile update
	profileUpdate := map[string]string{
		"firstName": "Updated",
		"lastName":  "Name",
		"email":     "updated@example.com",
		"phone":     "1234567890",
	}

	// Set up expectations
	mockUserRepo.On("GetByID", "user123").Return(existingUser, nil)
	mockUserRepo.On("GetByEmail", "updated@example.com").Return(nil, database.ErrNotFound)
	mockUserRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

	// Create request with authenticated context
	body, _ := json.Marshal(profileUpdate)
	req, _ := http.NewRequest("PUT", "/users/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.UpdateProfile).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.User
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify updated user data
	assert.Equal(t, "user123", response.ID)
	assert.Equal(t, "testuser", response.Username)
	assert.Equal(t, "updated@example.com", response.Email)
	assert.Equal(t, "Updated", response.FirstName)
	assert.Equal(t, "Name", response.LastName)
	assert.Equal(t, "1234567890", response.Phone)
	assert.Equal(t, "", response.Password) // Password should be removed

	// Verify expectations
	mockUserRepo.AssertExpectations(t)
}

// TestChangePassword tests the change password endpoint
func TestChangePassword(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockPreferenceRepo := new(MockUserPreferenceRepository)
	mockAPIKeyRepo := new(MockAPIKeyRepository)

	// Create handler
	handler := api.NewUserHandler(mockUserRepo, mockPreferenceRepo, mockAPIKeyRepo)

	// Create existing user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	existingUser := &models.User{
		ID:        "user123",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		FirstName: "Test",
		LastName:  "User",
		Role:      models.RoleUser,
		Active:    true,
	}

	// Create password change request
	passwordChange := map[string]string{
		"currentPassword": "oldpassword",
		"newPassword":     "newpassword123",
	}

	// Set up expectations
	mockUserRepo.On("GetByID", "user123").Return(existingUser, nil)
	mockUserRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

	// Create request with authenticated context
	body, _ := json.Marshal(passwordChange)
	req, _ := http.NewRequest("PUT", "/users/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.ChangePassword).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify expectations
	mockUserRepo.AssertExpectations(t)
}

// TestGetPreferences tests the get preferences endpoint
func TestGetPreferences(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockPreferenceRepo := new(MockUserPreferenceRepository)
	mockAPIKeyRepo := new(MockAPIKeyRepository)

	// Create handler
	handler := api.NewUserHandler(mockUserRepo, mockPreferenceRepo, mockAPIKeyRepo)

	// Create test preferences
	preferences := &models.UserPreferences{
		ID:     "pref123",
		UserID: "user123",
		Theme:  "dark",
		Language: "en",
		NotificationsEnabled: true,
	}

	// Set up expectations
	mockPreferenceRepo.On("GetByUserID", "user123").Return(preferences, nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/users/preferences", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.GetPreferences).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.UserPreferences
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify preferences data
	assert.Equal(t, "pref123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "dark", response.Theme)
	assert.Equal(t, "en", response.Language)
	assert.Equal(t, true, response.NotificationsEnabled)

	// Verify expectations
	mockPreferenceRepo.AssertExpectations(t)
}

// TestUpdatePreferences tests the update preferences endpoint
func TestUpdatePreferences(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockPreferenceRepo := new(MockUserPreferenceRepository)
	mockAPIKeyRepo := new(MockAPIKeyRepository)

	// Create handler
	handler := api.NewUserHandler(mockUserRepo, mockPreferenceRepo, mockAPIKeyRepo)

	// Create existing preferences
	existingPreferences := &models.UserPreferences{
		ID:     "pref123",
		UserID: "user123",
		Theme:  "light",
		Language: "en",
		NotificationsEnabled: false,
	}

	// Create preferences update
	preferencesUpdate := models.UserPreferences{
		Theme:  "dark",
		Language: "fr",
		NotificationsEnabled: true,
	}

	// Set up expectations
	mockPreferenceRepo.On("GetByUserID", "user123").Return(existingPreferences, nil)
	mockPreferenceRepo.On("Update", mock.AnythingOfType("*models.UserPreferences")).Return(nil)

	// Create request with authenticated context
	body, _ := json.Marshal(preferencesUpdate)
	req, _ := http.NewRequest("PUT", "/users/preferences", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.UpdatePreferences).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.UserPreferences
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify updated preferences data
	assert.Equal(t, "pref123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "dark", response.Theme)
	assert.Equal(t, "fr", response.Language)
	assert.Equal(t, true, response.NotificationsEnabled)

	// Verify expectations
	mockPreferenceRepo.AssertExpectations(t)
}
