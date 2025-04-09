package user

import (
	"errors"
	"time"

	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/repositories"
)

// UserService defines the interface for user-related operations
type UserService interface {
	GetUserSettings(userID string) (*models.UserSettings, error)
	UpdateUserSettings(userID string, settings *models.UserSettings) (*models.UserSettings, error)
	GetUserPreferences(userID string) (*models.UserPreferences, error)
	UpdateUserPreferences(userID string, preferences *models.UserPreferences) (*models.UserPreferences, error)
	GetUserTheme(userID string) (*models.UserTheme, error)
	UpdateUserTheme(userID string, theme *models.UserTheme) (*models.UserTheme, error)
	GetUserLayout(userID string, layoutName string) (*models.UserLayout, error)
	GetAllUserLayouts(userID string) ([]models.UserLayout, error)
	SaveUserLayout(userID string, layout *models.UserLayout) (*models.UserLayout, error)
	DeleteUserLayout(userID string, layoutName string) error
	GetUserApiKeys(userID string) ([]models.UserApiKey, error)
	AddUserApiKey(userID string, apiKey *models.UserApiKey) (*models.UserApiKey, error)
	UpdateUserApiKey(userID string, keyID string, apiKey *models.UserApiKey) (*models.UserApiKey, error)
	DeleteUserApiKey(userID string, keyID string) error
	GetUserNotificationSettings(userID string) (*models.UserNotificationSettings, error)
	UpdateUserNotificationSettings(userID string, settings *models.UserNotificationSettings) (*models.UserNotificationSettings, error)
}

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// GetUserSettings retrieves user settings
func (s *UserServiceImpl) GetUserSettings(userID string) (*models.UserSettings, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get the user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get the user settings
	settings, err := s.userRepo.GetUserSettings(userID)
	if err != nil {
		// If settings don't exist, create default settings
		defaultSettings := &models.UserSettings{
			UserID:    userID,
			Language:  "en",
			TimeZone:  "UTC",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return s.userRepo.CreateUserSettings(defaultSettings)
	}

	return settings, nil
}

// UpdateUserSettings updates user settings
func (s *UserServiceImpl) UpdateUserSettings(userID string, settings *models.UserSettings) (*models.UserSettings, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if settings == nil {
		return nil, errors.New("settings cannot be nil")
	}

	// Validate settings
	if err := settings.Validate(); err != nil {
		return nil, err
	}

	// Ensure the user ID matches
	if settings.UserID != userID {
		return nil, errors.New("user ID mismatch")
	}

	// Check if settings exist
	existingSettings, err := s.userRepo.GetUserSettings(userID)
	if err != nil {
		// If settings don't exist, create them
		settings.CreatedAt = time.Now()
		return s.userRepo.CreateUserSettings(settings)
	}

	// Preserve creation time
	settings.CreatedAt = existingSettings.CreatedAt
	settings.UpdatedAt = time.Now()

	// Update settings
	return s.userRepo.UpdateUserSettings(settings)
}

// GetUserPreferences retrieves user preferences
func (s *UserServiceImpl) GetUserPreferences(userID string) (*models.UserPreferences, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get the user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get the user preferences
	preferences, err := s.userRepo.GetUserPreferences(userID)
	if err != nil {
		// If preferences don't exist, create default preferences
		defaultPreferences := &models.UserPreferences{
			UserID:                userID,
			DefaultOrderQuantity:  1,
			DefaultProductType:    models.ProductTypeMIS,
			DefaultExchange:       "NSE",
			ShowConfirmationDialog: true,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}
		return s.userRepo.CreateUserPreferences(defaultPreferences)
	}

	return preferences, nil
}

// UpdateUserPreferences updates user preferences
func (s *UserServiceImpl) UpdateUserPreferences(userID string, preferences *models.UserPreferences) (*models.UserPreferences, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if preferences == nil {
		return nil, errors.New("preferences cannot be nil")
	}

	// Validate preferences
	if err := preferences.Validate(); err != nil {
		return nil, err
	}

	// Ensure the user ID matches
	if preferences.UserID != userID {
		return nil, errors.New("user ID mismatch")
	}

	// Check if preferences exist
	existingPreferences, err := s.userRepo.GetUserPreferences(userID)
	if err != nil {
		// If preferences don't exist, create them
		preferences.CreatedAt = time.Now()
		return s.userRepo.CreateUserPreferences(preferences)
	}

	// Preserve creation time
	preferences.CreatedAt = existingPreferences.CreatedAt
	preferences.UpdatedAt = time.Now()

	// Update preferences
	return s.userRepo.UpdateUserPreferences(preferences)
}

// GetUserTheme retrieves user theme settings
func (s *UserServiceImpl) GetUserTheme(userID string) (*models.UserTheme, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get the user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get the user theme
	theme, err := s.userRepo.GetUserTheme(userID)
	if err != nil {
		// If theme doesn't exist, create default theme
		defaultTheme := &models.UserTheme{
			UserID:     userID,
			ThemeMode:  "light",
			PrimaryColor: "#1976d2",
			SecondaryColor: "#dc004e",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		return s.userRepo.CreateUserTheme(defaultTheme)
	}

	return theme, nil
}

// UpdateUserTheme updates user theme settings
func (s *UserServiceImpl) UpdateUserTheme(userID string, theme *models.UserTheme) (*models.UserTheme, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if theme == nil {
		return nil, errors.New("theme cannot be nil")
	}

	// Validate theme
	if err := theme.Validate(); err != nil {
		return nil, err
	}

	// Ensure the user ID matches
	if theme.UserID != userID {
		return nil, errors.New("user ID mismatch")
	}

	// Check if theme exists
	existingTheme, err := s.userRepo.GetUserTheme(userID)
	if err != nil {
		// If theme doesn't exist, create it
		theme.CreatedAt = time.Now()
		return s.userRepo.CreateUserTheme(theme)
	}

	// Preserve creation time
	theme.CreatedAt = existingTheme.CreatedAt
	theme.UpdatedAt = time.Now()

	// Update theme
	return s.userRepo.UpdateUserTheme(theme)
}

// GetUserLayout retrieves a specific user layout
func (s *UserServiceImpl) GetUserLayout(userID string, layoutName string) (*models.UserLayout, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if layoutName == "" {
		return nil, errors.New("layout name is required")
	}

	// Get the user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get the user layout
	return s.userRepo.GetUserLayout(userID, layoutName)
}

// GetAllUserLayouts retrieves all user layouts
func (s *UserServiceImpl) GetAllUserLayouts(userID string) ([]models.UserLayout, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get the user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get all user layouts
	return s.userRepo.GetAllUserLayouts(userID)
}

// SaveUserLayout saves a user layout
func (s *UserServiceImpl) SaveUserLayout(userID string, layout *models.UserLayout) (*models.UserLayout, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if layout == nil {
		return nil, errors.New("layout cannot be nil")
	}
	if layout.Name == "" {
		return nil, errors.New("layout name is required")
	}

	// Validate layout
	if err := layout.Validate(); err != nil {
		return nil, err
	}

	// Ensure the user ID matches
	if layout.UserID != userID {
		return nil, errors.New("user ID mismatch")
	}

	// Check if layout exists
	existingLayout, err := s.userRepo.GetUserLayout(userID, layout.Name)
	if err == nil {
		// If layout exists, update it
		layout.CreatedAt = existingLayout.CreatedAt
		layout.UpdatedAt = time.Now()
		return s.userRepo.UpdateUserLayout(layout)
	}

	// If layout doesn't exist, create it
	layout.CreatedAt = time.Now()
	layout.UpdatedAt = time.Now()
	return s.userRepo.CreateUserLayout(layout)
}

// DeleteUserLayout deletes a user layout
func (s *UserServiceImpl) DeleteUserLayout(userID string, layoutName string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}
	if layoutName == "" {
		return errors.New("layout name is required")
	}

	// Check if layout exists
	_, err := s.userRepo.GetUserLayout(userID, layoutName)
	if err != nil {
		return errors.New("layout not found")
	}

	// Delete the layout
	return s.userRepo.DeleteUserLayout(userID, layoutName)
}

// GetUserApiKeys retrieves user API keys
func (s *UserServiceImpl) GetUserApiKeys(userID string) ([]models.UserApiKey, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get the user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get the user API keys
	return s.userRepo.GetUserApiKeys(userID)
}

// AddUserApiKey adds a new user API key
func (s *UserServiceImpl) AddUserApiKey(userID string, apiKey *models.UserApiKey) (*models.UserApiKey, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if apiKey == nil {
		return nil, errors.New("API key cannot be nil")
	}

	// Validate API key
	if err := apiKey.Validate(); err != nil {
		return nil, err
	}

	// Ensure the user ID matches
	if apiKey.UserID != userID {
		return nil, errors.New("user ID mismatch")
	}

	// Generate a unique ID if not provided
	if apiKey.ID == "" {
		apiKey.ID = generateUniqueID()
	}

	// Set timestamps
	apiKey.CreatedAt = time.Now()
	apiKey.UpdatedAt = time.Now()

	// Add the API key
	return s.userRepo.CreateUserApiKey(apiKey)
}

// UpdateUserApiKey updates a user API key
func (s *UserServiceImpl) UpdateUserApiKey(userID string, keyID string, apiKey *models.UserApiKey) (*models.UserApiKey, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if keyID == "" {
		return nil, errors.New("key ID is required")
	}
	if apiKey == nil {
		return nil, errors.New("API key cannot be nil")
	}

	// Validate API key
	if err := apiKey.Validate(); err != nil {
		return nil, err
	}

	// Ensure the user ID and key ID match
	if apiKey.UserID != userID {
		return nil, errors.New("user ID mismatch")
	}
	if apiKey.ID != keyID {
		return nil, errors.New("key ID mismatch")
	}

	// Check if API key exists
	existingApiKey, err := s.userRepo.GetUserApiKey(userID, keyID)
	if err != nil {
		return nil, errors.New("API key not found")
	}

	// Preserve creation time
	apiKey.CreatedAt = existingApiKey.CreatedAt
	apiKey.UpdatedAt = time.Now()

	// Update the API key
	return s.userRepo.UpdateUserApiKey(apiKey)
}

// DeleteUserApiKey deletes a user API key
func (s *UserServiceImpl) DeleteUserApiKey(userID string, keyID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}
	if keyID == "" {
		return errors.New("key ID is required")
	}

	// Check if API key exists
	_, err := s.userRepo.GetUserApiKey(userID, keyID)
	if err != nil {
		return errors.New("API key not found")
	}

	// Delete the API key
	return s.userRepo.DeleteUserApiKey(userID, keyID)
}

// GetUserNotificationSettings retrieves user notification settings
func (s *UserServiceImpl) GetUserNotificationSettings(userID string) (*models.UserNotificationSettings, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get the user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get the user notification settings
	settings, err := s.userRepo.GetUserNotificationSettings(userID)
	if err != nil {
		// If settings don't exist, create default settings
		defaultSettings := &models.UserNotificationSettings{
			UserID:                 userID,
			EnableEmailNotifications: true,
			EnablePushNotifications:  true,
			OrderExecutionAlerts:     true,
			PriceAlerts:              true,
			MarginCallAlerts:         true,
			NewsAlerts:               false,
			CreatedAt:              time.Now(),
			UpdatedAt:              time.Now(),
		}
		return s.userRepo.CreateUserNotificationSettings(defaultSettings)
	}

	return settings, nil
}

// UpdateUserNotificationSettings updates user notification settings
func (s *UserServiceImpl) UpdateUserNotificationSettings(userID string, settings *models.UserNotificationSettings) (*models.UserNotificationSettings, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if settings == nil {
		return nil, errors.New("settings cannot be nil")
	}

	// Validate settings
	if err := settings.Validate(); err != nil {
		return nil, err
	}

	// Ensure the user ID matches
	if settings.UserID != userID {
		return nil, errors.New("user ID mismatch")
	}

	// Check if settings exist
	existingSettings, err := s.userRepo.GetUserNotificationSettings(userID)
	if err != nil {
		// If settings don't exist, create them
		settings.CreatedAt = time.Now()
		return s.userRepo.CreateUserNotificationSettings(settings)
	}

	// Preserve creation time
	settings.CreatedAt = existingSettings.CreatedAt
	settings.UpdatedAt = time.Now()

	// Update settings
	return s.userRepo.UpdateUserNotificationSettings(settings)
}

// Helper function to generate a unique ID
func generateUniqueID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// Helper function to generate a random string
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond)
	}
	return string(result)
}
