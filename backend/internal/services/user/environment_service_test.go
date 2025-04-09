package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"trading_platform/backend/internal/auth"
	"trading_platform/backend/internal/models"
	"trading_platform/backend/internal/repositories/mocks"
)

// TestSwitchEnvironment tests the environment switching functionality
func TestSwitchEnvironment(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockPrefsRepo := new(mocks.MockUserPreferencesRepository)

	// Create service with mock repositories
	service := NewEnvironmentService(mockUserRepo, mockPrefsRepo)

	// Test data
	userID := "user123"
	username := "simuser"
	role := string(models.UserRoleTrader)
	userType := string(models.UserTypeSIM)
	targetEnv := models.EnvironmentSIM

	simUser := &models.User{
		ID:       userID,
		Username: username,
		Role:     models.UserRoleTrader,
		UserType: models.UserTypeSIM,
	}

	userPrefs := &models.UserPreferences{
		UserID:             userID,
		DefaultEnvironment: models.EnvironmentLive,
	}

	updatedPrefs := &models.UserPreferences{
		UserID:             userID,
		DefaultEnvironment: targetEnv,
	}

	// Create context with user ID
	ctx := context.Background()
	ctx = auth.SetUserIDInContext(ctx, userID)
	ctx = auth.SetRoleInContext(ctx, role)

	// Set up expectations
	mockUserRepo.On("GetByID", mock.Anything, userID).Return(simUser, nil)
	mockPrefsRepo.On("GetByUserID", mock.Anything, userID).Return(userPrefs, nil)
	mockPrefsRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *models.UserPreferences) bool {
		return p.UserID == userID && p.DefaultEnvironment == targetEnv
	})).Return(updatedPrefs, nil)

	// Call the method
	token, err := service.SwitchEnvironment(ctx, targetEnv)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify expectations
	mockUserRepo.AssertExpectations(t)
	mockPrefsRepo.AssertExpectations(t)
}

// TestSwitchEnvironmentNonSimUser tests that non-SIM users cannot switch to SIM environment
func TestSwitchEnvironmentNonSimUser(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockPrefsRepo := new(mocks.MockUserPreferencesRepository)

	// Create service with mock repositories
	service := NewEnvironmentService(mockUserRepo, mockPrefsRepo)

	// Test data
	userID := "user123"
	username := "standarduser"
	role := string(models.UserRoleTrader)
	targetEnv := models.EnvironmentSIM

	standardUser := &models.User{
		ID:       userID,
		Username: username,
		Role:     models.UserRoleTrader,
		UserType: models.UserTypeStandard,
	}

	// Create context with user ID
	ctx := context.Background()
	ctx = auth.SetUserIDInContext(ctx, userID)
	ctx = auth.SetRoleInContext(ctx, role)

	// Set up expectations
	mockUserRepo.On("GetByID", mock.Anything, userID).Return(standardUser, nil)

	// Call the method
	token, err := service.SwitchEnvironment(ctx, targetEnv)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "only SIM users can access")

	// Verify expectations
	mockUserRepo.AssertExpectations(t)
}

// TestGetEnvironmentStatus tests getting the environment status
func TestGetEnvironmentStatus(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockPrefsRepo := new(mocks.MockUserPreferencesRepository)

	// Create service with mock repositories
	service := NewEnvironmentService(mockUserRepo, mockPrefsRepo)

	// Test data
	userID := "user123"
	userType := string(models.UserTypeSIM)
	environment := string(models.EnvironmentSIM)

	// Create context with user ID, user type, and environment
	ctx := context.Background()
	ctx = auth.SetUserIDInContext(ctx, userID)
	ctx = auth.SetUserTypeInContext(ctx, userType)
	ctx = auth.SetEnvironmentInContext(ctx, environment)

	// Call the method
	status := service.GetEnvironmentStatus(ctx)

	// Assertions
	assert.NotNil(t, status)
	assert.Equal(t, environment, status["environment"])
	assert.Equal(t, userType, status["userType"])
	assert.Equal(t, true, status["isSimulation"])
	assert.NotNil(t, status["timestamp"])
}

// TestIsSimulationEnvironment tests checking if current environment is simulation
func TestIsSimulationEnvironment(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockPrefsRepo := new(mocks.MockUserPreferencesRepository)

	// Create service with mock repositories
	service := NewEnvironmentService(mockUserRepo, mockPrefsRepo)

	// Test cases
	testCases := []struct {
		name        string
		environment string
		expected    bool
	}{
		{
			name:        "SIM environment",
			environment: string(models.EnvironmentSIM),
			expected:    true,
		},
		{
			name:        "LIVE environment",
			environment: string(models.EnvironmentLive),
			expected:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context with environment
			ctx := context.Background()
			ctx = auth.SetEnvironmentInContext(ctx, tc.environment)

			// Call the method
			result := service.IsSimulationEnvironment(ctx)

			// Assertions
			assert.Equal(t, tc.expected, result)
		})
	}
}
