package user

import (
	"context"
	"errors"
	"time"

	"trading_platform/backend/internal/auth"
	"trading_platform/backend/internal/models"
	"trading_platform/backend/internal/repositories"
)

// EnvironmentService handles operations related to environment switching and management
type EnvironmentService struct {
	userRepo       repositories.UserRepository
	preferencesRepo repositories.UserPreferencesRepository
}

// NewEnvironmentService creates a new environment service
func NewEnvironmentService(userRepo repositories.UserRepository, preferencesRepo repositories.UserPreferencesRepository) *EnvironmentService {
	return &EnvironmentService{
		userRepo:       userRepo,
		preferencesRepo: preferencesRepo,
	}
}

// SwitchEnvironment switches the user's environment
func (s *EnvironmentService) SwitchEnvironment(ctx context.Context, environment models.Environment) (string, error) {
	// Get user ID from context
	userID := auth.GetUserIDFromContext(ctx)
	if userID == "" {
		return "", errors.New("user not authenticated")
	}

	// Get user from repository
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	// Check if user is allowed to switch to the requested environment
	if environment == models.EnvironmentSIM && user.UserType != models.UserTypeSIM {
		return "", errors.New("only SIM users can access the simulation environment")
	}

	// Get user role from context
	role := auth.GetRoleFromContext(ctx)
	if role == "" {
		return "", errors.New("user role not found")
	}

	// Generate new token with updated environment
	token, err := auth.GenerateToken(userID, user.Username, string(user.Role), string(user.UserType), string(environment))
	if err != nil {
		return "", err
	}

	// Update user preferences with new default environment
	preferences, err := s.preferencesRepo.GetByUserID(ctx, userID)
	if err == nil {
		preferences.DefaultEnvironment = environment
		preferences.UpdatedAt = time.Now()
		_, err = s.preferencesRepo.Update(ctx, preferences)
		if err != nil {
			// Log error but continue (non-critical)
			// logger.Error("Failed to update user preferences", "error", err)
		}
	}

	return token, nil
}

// GetCurrentEnvironment gets the current environment from context
func (s *EnvironmentService) GetCurrentEnvironment(ctx context.Context) models.Environment {
	environment := auth.GetEnvironmentFromContext(ctx)
	return models.Environment(environment)
}

// IsSimulationEnvironment checks if the current environment is simulation
func (s *EnvironmentService) IsSimulationEnvironment(ctx context.Context) bool {
	return s.GetCurrentEnvironment(ctx) == models.EnvironmentSIM
}

// GetEnvironmentStatus returns the status of the current environment
func (s *EnvironmentService) GetEnvironmentStatus(ctx context.Context) map[string]interface{} {
	environment := s.GetCurrentEnvironment(ctx)
	userType := models.UserType(auth.GetUserTypeFromContext(ctx))
	
	return map[string]interface{}{
		"environment": environment,
		"userType": userType,
		"isSimulation": environment == models.EnvironmentSIM,
		"timestamp": time.Now(),
	}
}
