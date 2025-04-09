package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"trading_platform/backend/internal/models"
	"trading_platform/backend/internal/repositories/mocks"
)

// TestCreateSIMUser tests the creation of a SIM user
func TestCreateSIMUser(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockPrefsRepo := new(mocks.MockUserPreferencesRepository)

	// Create service with mock repositories
	service := NewUserService(mockUserRepo, mockPrefsRepo)

	// Test data
	now := time.Now()
	simUser := &models.User{
		Username:     "simuser",
		Email:        "sim@example.com",
		PasswordHash: "hashedpassword",
		FirstName:    "Sim",
		LastName:     "User",
		Role:         models.UserRoleTrader,
		UserType:     models.UserTypeSIM,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Set up expectations
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Username == simUser.Username && u.UserType == models.UserTypeSIM
	})).Return(simUser, nil)

	mockPrefsRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *models.UserPreferences) bool {
		return p.UserID == simUser.ID && p.DefaultEnvironment == models.EnvironmentSIM
	})).Return(&models.UserPreferences{
		UserID:             simUser.ID,
		DefaultEnvironment: models.EnvironmentSIM,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil)

	// Call the method
	createdUser, err := service.CreateUser(context.Background(), simUser)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, models.UserTypeSIM, createdUser.UserType)
	
	// Verify expectations
	mockUserRepo.AssertExpectations(t)
	mockPrefsRepo.AssertExpectations(t)
}

// TestSwitchUserType tests switching a user's type
func TestSwitchUserType(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockPrefsRepo := new(mocks.MockUserPreferencesRepository)

	// Create service with mock repositories
	service := NewUserService(mockUserRepo, mockPrefsRepo)

	// Test data
	userID := "user123"
	standardUser := &models.User{
		ID:        userID,
		Username:  "testuser",
		UserType:  models.UserTypeStandard,
		UpdatedAt: time.Now(),
	}
	
	updatedUser := &models.User{
		ID:        userID,
		Username:  "testuser",
		UserType:  models.UserTypeSIM,
		UpdatedAt: time.Now(),
	}

	// Set up expectations
	mockUserRepo.On("GetByID", mock.Anything, userID).Return(standardUser, nil)
	mockUserRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.ID == userID && u.UserType == models.UserTypeSIM
	})).Return(updatedUser, nil)

	// Call the method
	result, err := service.SwitchUserType(context.Background(), userID, models.UserTypeSIM)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, models.UserTypeSIM, result.UserType)
	
	// Verify expectations
	mockUserRepo.AssertExpectations(t)
}

// TestValidateUserTypeAccess tests validation of user type access
func TestValidateUserTypeAccess(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockPrefsRepo := new(mocks.MockUserPreferencesRepository)

	// Create service with mock repositories
	service := NewUserService(mockUserRepo, mockPrefsRepo)

	// Test cases
	testCases := []struct {
		name        string
		userType    models.UserType
		targetType  models.UserType
		shouldAllow bool
	}{
		{
			name:        "Admin can access SIM",
			userType:    models.UserTypeAdmin,
			targetType:  models.UserTypeSIM,
			shouldAllow: true,
		},
		{
			name:        "Admin can access Standard",
			userType:    models.UserTypeAdmin,
			targetType:  models.UserTypeStandard,
			shouldAllow: true,
		},
		{
			name:        "Standard cannot access SIM",
			userType:    models.UserTypeStandard,
			targetType:  models.UserTypeSIM,
			shouldAllow: false,
		},
		{
			name:        "SIM can access SIM",
			userType:    models.UserTypeSIM,
			targetType:  models.UserTypeSIM,
			shouldAllow: true,
		},
		{
			name:        "SIM cannot access Admin",
			userType:    models.UserTypeSIM,
			targetType:  models.UserTypeAdmin,
			shouldAllow: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the method
			err := service.ValidateUserTypeAccess(tc.userType, tc.targetType)
			
			// Assertions
			if tc.shouldAllow {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
