package repositories

import (
	"errors"
	"time"

	"github.com/trading-platform/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// User operations
	Create(user *models.User) (*models.User, error)
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id string) error
	
	// User settings operations
	GetUserSettings(userID string) (*models.UserSettings, error)
	CreateUserSettings(settings *models.UserSettings) (*models.UserSettings, error)
	UpdateUserSettings(settings *models.UserSettings) (*models.UserSettings, error)
	
	// User preferences operations
	GetUserPreferences(userID string) (*models.UserPreferences, error)
	CreateUserPreferences(preferences *models.UserPreferences) (*models.UserPreferences, error)
	UpdateUserPreferences(preferences *models.UserPreferences) (*models.UserPreferences, error)
	
	// User theme operations
	GetUserTheme(userID string) (*models.UserTheme, error)
	CreateUserTheme(theme *models.UserTheme) (*models.UserTheme, error)
	UpdateUserTheme(theme *models.UserTheme) (*models.UserTheme, error)
	
	// User layout operations
	GetUserLayout(userID string, layoutName string) (*models.UserLayout, error)
	GetAllUserLayouts(userID string) ([]models.UserLayout, error)
	CreateUserLayout(layout *models.UserLayout) (*models.UserLayout, error)
	UpdateUserLayout(layout *models.UserLayout) (*models.UserLayout, error)
	DeleteUserLayout(userID string, layoutName string) error
	
	// User API key operations
	GetUserApiKeys(userID string) ([]models.UserApiKey, error)
	GetUserApiKey(userID string, keyID string) (*models.UserApiKey, error)
	CreateUserApiKey(apiKey *models.UserApiKey) (*models.UserApiKey, error)
	UpdateUserApiKey(apiKey *models.UserApiKey) (*models.UserApiKey, error)
	DeleteUserApiKey(userID string, keyID string) error
	
	// User notification settings operations
	GetUserNotificationSettings(userID string) (*models.UserNotificationSettings, error)
	CreateUserNotificationSettings(settings *models.UserNotificationSettings) (*models.UserNotificationSettings, error)
	UpdateUserNotificationSettings(settings *models.UserNotificationSettings) (*models.UserNotificationSettings, error)
}

// MongoUserRepository implements UserRepository using MongoDB
type MongoUserRepository struct {
	db *mongo.Database
}

// NewMongoUserRepository creates a new MongoUserRepository
func NewMongoUserRepository(db *mongo.Database) UserRepository {
	return &MongoUserRepository{
		db: db,
	}
}

// Create adds a new user to the database
func (r *MongoUserRepository) Create(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user with the same email already exists
	existingUser, _ := r.GetByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Generate a new ID if not provided
	if user.ID == "" {
		user.ID = primitive.NewObjectID().Hex()
	}

	// Insert the user
	_, err := r.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *MongoUserRepository) GetByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	filter := bson.M{"_id": id}

	err := r.db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *MongoUserRepository) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	filter := bson.M{"email": email}

	err := r.db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user
func (r *MongoUserRepository) Update(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the user
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	_, err := r.db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete removes a user from the database
func (r *MongoUserRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := r.db.Collection("users").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetUserSettings retrieves user settings
func (r *MongoUserRepository) GetUserSettings(userID string) (*models.UserSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var settings models.UserSettings
	filter := bson.M{"userId": userID}

	err := r.db.Collection("user_settings").FindOne(ctx, filter).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user settings not found")
		}
		return nil, err
	}

	return &settings, nil
}

// CreateUserSettings creates new user settings
func (r *MongoUserRepository) CreateUserSettings(settings *models.UserSettings) (*models.UserSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the settings
	_, err := r.db.Collection("user_settings").InsertOne(ctx, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateUserSettings updates existing user settings
func (r *MongoUserRepository) UpdateUserSettings(settings *models.UserSettings) (*models.UserSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the settings
	filter := bson.M{"userId": settings.UserID}
	update := bson.M{"$set": settings}

	_, err := r.db.Collection("user_settings").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// GetUserPreferences retrieves user preferences
func (r *MongoUserRepository) GetUserPreferences(userID string) (*models.UserPreferences, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var preferences models.UserPreferences
	filter := bson.M{"userId": userID}

	err := r.db.Collection("user_preferences").FindOne(ctx, filter).Decode(&preferences)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user preferences not found")
		}
		return nil, err
	}

	return &preferences, nil
}

// CreateUserPreferences creates new user preferences
func (r *MongoUserRepository) CreateUserPreferences(preferences *models.UserPreferences) (*models.UserPreferences, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the preferences
	_, err := r.db.Collection("user_preferences").InsertOne(ctx, preferences)
	if err != nil {
		return nil, err
	}

	return preferences, nil
}

// UpdateUserPreferences updates existing user preferences
func (r *MongoUserRepository) UpdateUserPreferences(preferences *models.UserPreferences) (*models.UserPreferences, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the preferences
	filter := bson.M{"userId": preferences.UserID}
	update := bson.M{"$set": preferences}

	_, err := r.db.Collection("user_preferences").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return preferences, nil
}

// GetUserTheme retrieves user theme settings
func (r *MongoUserRepository) GetUserTheme(userID string) (*models.UserTheme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var theme models.UserTheme
	filter := bson.M{"userId": userID}

	err := r.db.Collection("user_themes").FindOne(ctx, filter).Decode(&theme)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user theme not found")
		}
		return nil, err
	}

	return &theme, nil
}

// CreateUserTheme creates new user theme settings
func (r *MongoUserRepository) CreateUserTheme(theme *models.UserTheme) (*models.UserTheme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the theme
	_, err := r.db.Collection("user_themes").InsertOne(ctx, theme)
	if err != nil {
		return nil, err
	}

	return theme, nil
}

// UpdateUserTheme updates existing user theme settings
func (r *MongoUserRepository) UpdateUserTheme(theme *models.UserTheme) (*models.UserTheme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the theme
	filter := bson.M{"userId": theme.UserID}
	update := bson.M{"$set": theme}

	_, err := r.db.Collection("user_themes").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return theme, nil
}

// GetUserLayout retrieves a specific user layout
func (r *MongoUserRepository) GetUserLayout(userID string, layoutName string) (*models.UserLayout, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var layout models.UserLayout
	filter := bson.M{"userId": userID, "name": layoutName}

	err := r.db.Collection("user_layouts").FindOne(ctx, filter).Decode(&layout)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user layout not found")
		}
		return nil, err
	}

	return &layout, nil
}

// GetAllUserLayouts retrieves all user layouts
func (r *MongoUserRepository) GetAllUserLayouts(userID string) ([]models.UserLayout, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var layouts []models.UserLayout
	filter := bson.M{"userId": userID}

	cursor, err := r.db.Collection("user_layouts").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &layouts); err != nil {
		return nil, err
	}

	return layouts, nil
}

// CreateUserLayout creates a new user layout
func (r *MongoUserRepository) CreateUserLayout(layout *models.UserLayout) (*models.UserLayout, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the layout
	_, err := r.db.Collection("user_layouts").InsertOne(ctx, layout)
	if err != nil {
		return nil, err
	}

	return layout, nil
}

// UpdateUserLayout updates an existing user layout
func (r *MongoUserRepository) UpdateUserLayout(layout *models.UserLayout) (*models.UserLayout, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the layout
	filter := bson.M{"userId": layout.UserID, "name": layout.Name}
	update := bson.M{"$set": layout}

	_, err := r.db.Collection("user_layouts").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return layout, nil
}

// DeleteUserLayout deletes a user layout
func (r *MongoUserRepository) DeleteUserLayout(userID string, layoutName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userId": userID, "name": layoutName}
	result, err := r.db.Collection("user_layouts").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user layout not found")
	}

	return nil
}

// GetUserApiKeys retrieves all user API keys
func (r *MongoUserRepository) GetUserApiKeys(userID string) ([]models.UserApiKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var apiKeys []models.UserApiKey
	filter := bson.M{"userId": userID}

	cursor, err := r.db.Collection("user_api_keys").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &apiKeys); err != nil {
		return nil, err
	}

	return apiKeys, nil
}

// GetUserApiKey retrieves a specific user API key
func (r *MongoUserRepository) GetUserApiKey(userID string, keyID string) (*models.UserApiKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var apiKey models.UserApiKey
	filter := bson.M{"userId": userID, "_id": keyID}

	err := r.db.Collection("user_api_keys").FindOne(ctx, filter).Decode(&apiKey)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user API key not found")
		}
		return nil, err
	}

	return &apiKey, nil
}

// CreateUserApiKey creates a new user API key
func (r *MongoUserRepository) CreateUserApiKey(apiKey *models.UserApiKey) (*models.UserApiKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the API key
	_, err := r.db.Collection("user_api_keys").InsertOne(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

// UpdateUserApiKey updates an existing user API key
func (r *MongoUserRepository) UpdateUserApiKey(apiKey *models.UserApiKey) (*models.UserApiKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the API key
	filter := bson.M{"userId": apiKey.UserID, "_id": apiKey.ID}
	update := bson.M{"$set": apiKey}

	_, err := r.db.Collection("user_api_keys").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

// DeleteUserApiKey deletes a user API key
func (r *MongoUserRepository) DeleteUserApiKey(userID string, keyID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userId": userID, "_id": keyID}
	result, err := r.db.Collection("user_api_keys").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user API key not found")
	}

	return nil
}

// GetUserNotificationSettings retrieves user notification settings
func (r *MongoUserRepository) GetUserNotificationSettings(userID string) (*models.UserNotificationSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var settings models.UserNotificationSettings
	filter := bson.M{"userId": userID}

	err := r.db.Collection("user_notification_settings").FindOne(ctx, filter).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user notification settings not found")
		}
		return nil, err
	}

	return &settings, nil
}

// CreateUserNotificationSettings creates new user notification settings
func (r *MongoUserRepository) CreateUserNotificationSettings(settings *models.UserNotificationSettings) (*models.UserNotificationSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the settings
	_, err := r.db.Collection("user_notification_settings").InsertOne(ctx, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateUserNotificationSettings updates existing user notification settings
func (r *MongoUserRepository) UpdateUserNotificationSettings(settings *models.UserNotificationSettings) (*models.UserNotificationSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the settings
	filter := bson.M{"userId": settings.UserID}
	update := bson.M{"$set": settings}

	_, err := r.db.Collection("user_notification_settings").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
