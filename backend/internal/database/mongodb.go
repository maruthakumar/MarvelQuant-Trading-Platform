package database

import (
	"context"
	"time"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	
	"trading_platform/backend/internal/models"
	"trading_platform/backend/internal/config"
)

// MongoDB represents the MongoDB client and database connection
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Collections
const (
	OrderCollection     = "orders"
	PositionCollection  = "positions"
	UserCollection      = "users"
	StrategyCollection  = "strategies"
	PortfolioCollection = "portfolios"
	PreferenceCollection = "preferences"
	APIKeyCollection    = "apikeys"
)

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(cfg *config.Config) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	
	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	db := client.Database(cfg.MongoDB.Database)
	
	// Create indexes
	createIndexes(ctx, db)
	
	return &MongoDB{
		Client:   client,
		Database: db,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	return m.Client.Disconnect(ctx)
}

// createIndexes creates indexes for collections
func createIndexes(ctx context.Context, db *mongo.Database) {
	// User collection indexes
	userIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	db.Collection(UserCollection).Indexes().CreateMany(ctx, userIndexes)
	
	// Order collection indexes
	orderIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "symbol", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "createdAt", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	}
	db.Collection(OrderCollection).Indexes().CreateMany(ctx, orderIndexes)
	
	// Position collection indexes
	positionIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "symbol", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	}
	db.Collection(PositionCollection).Indexes().CreateMany(ctx, positionIndexes)
	
	// Strategy collection indexes
	strategyIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "active", Value: 1},
			},
		},
	}
	db.Collection(StrategyCollection).Indexes().CreateMany(ctx, strategyIndexes)
	
	// Portfolio collection indexes
	portfolioIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "strategyId", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	}
	db.Collection(PortfolioCollection).Indexes().CreateMany(ctx, portfolioIndexes)
	
	// Preference collection indexes
	preferenceIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	db.Collection(PreferenceCollection).Indexes().CreateMany(ctx, preferenceIndexes)
	
	// API Key collection indexes
	apiKeyIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "brokerName", Value: 1},
			},
		},
	}
	db.Collection(APIKeyCollection).Indexes().CreateMany(ctx, apiKeyIndexes)
}

// OrderRepository provides methods for working with orders
type OrderRepository struct {
	db *MongoDB
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db *MongoDB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order
func (r *OrderRepository) Create(order *models.Order) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	
	result, err := r.db.Database.Collection(OrderCollection).InsertOne(ctx, order)
	if err != nil {
		return "", err
	}
	
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// GetByID retrieves an order by ID
func (r *OrderRepository) GetByID(id string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var order models.Order
	err = r.db.Database.Collection(OrderCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&order)
	if err != nil {
		return nil, err
	}
	
	return &order, nil
}

// Update updates an order
func (r *OrderRepository) Update(order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(order.ID)
	if err != nil {
		return err
	}
	
	order.UpdatedAt = time.Now()
	
	_, err = r.db.Database.Collection(OrderCollection).ReplaceOne(
		ctx,
		bson.M{"_id": objectID},
		order,
	)
	
	return err
}

// Delete deletes an order
func (r *OrderRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.db.Database.Collection(OrderCollection).DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// Find finds orders based on filter
func (r *OrderRepository) Find(filter models.OrderFilter, page, limit int) ([]*models.Order, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Build query
	query := bson.M{}
	
	if filter.UserID != "" {
		query["userId"] = filter.UserID
	}
	
	if filter.Symbol != "" {
		query["symbol"] = filter.Symbol
	}
	
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	
	if filter.Direction != "" {
		query["direction"] = filter.Direction
	}
	
	if filter.ProductType != "" {
		query["productType"] = filter.ProductType
	}
	
	if filter.InstrumentType != "" {
		query["instrumentType"] = filter.InstrumentType
	}
	
	if filter.PortfolioID != "" {
		query["portfolioId"] = filter.PortfolioID
	}
	
	if filter.StrategyID != "" {
		query["strategyId"] = filter.StrategyID
	}
	
	if !filter.FromDate.IsZero() && !filter.ToDate.IsZero() {
		query["createdAt"] = bson.M{
			"$gte": filter.FromDate,
			"$lte": filter.ToDate,
		}
	} else if !filter.FromDate.IsZero() {
		query["createdAt"] = bson.M{"$gte": filter.FromDate}
	} else if !filter.ToDate.IsZero() {
		query["createdAt"] = bson.M{"$lte": filter.ToDate}
	}
	
	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}
	
	// Count total
	total, err := r.db.Database.Collection(OrderCollection).CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	
	// Set pagination
	skip := (page - 1) * limit
	
	// Find orders
	cursor, err := r.db.Database.Collection(OrderCollection).Find(
		ctx,
		query,
		options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "createdAt", Value: -1}}),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var orders []*models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, 0, err
	}
	
	return orders, int(total), nil
}

// PositionRepository provides methods for working with positions
type PositionRepository struct {
	db *MongoDB
}

// NewPositionRepository creates a new PositionRepository
func NewPositionRepository(db *MongoDB) *PositionRepository {
	return &PositionRepository{db: db}
}

// Create creates a new position
func (r *PositionRepository) Create(position *models.Position) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	position.CreatedAt = time.Now()
	position.UpdatedAt = time.Now()
	
	result, err := r.db.Database.Collection(PositionCollection).InsertOne(ctx, position)
	if err != nil {
		return "", err
	}
	
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// GetByID retrieves a position by ID
func (r *PositionRepository) GetByID(id string) (*models.Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var position models.Position
	err = r.db.Database.Collection(PositionCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&position)
	if err != nil {
		return nil, err
	}
	
	return &position, nil
}

// Update updates a position
func (r *PositionRepository) Update(position *models.Position) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(position.ID)
	if err != nil {
		return err
	}
	
	position.UpdatedAt = time.Now()
	
	_, err = r.db.Database.Collection(PositionCollection).ReplaceOne(
		ctx,
		bson.M{"_id": objectID},
		position,
	)
	
	return err
}

// Delete deletes a position
func (r *PositionRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.db.Database.Collection(PositionCollection).DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// Find finds positions based on filter
func (r *PositionRepository) Find(filter models.PositionFilter, page, limit int) ([]*models.Position, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Build query
	query := bson.M{}
	
	if filter.UserID != "" {
		query["userId"] = filter.UserID
	}
	
	if filter.Symbol != "" {
		query["symbol"] = filter.Symbol
	}
	
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	
	if filter.Direction != "" {
		query["direction"] = filter.Direction
	}
	
	if filter.ProductType != "" {
		query["productType"] = filter.ProductType
	}
	
	if filter.InstrumentType != "" {
		query["instrumentType"] = filter.InstrumentType
	}
	
	if filter.PortfolioID != "" {
		query["portfolioId"] = filter.PortfolioID
	}
	
	if filter.StrategyID != "" {
		query["strategyId"] = filter.StrategyID
	}
	
	if !filter.FromDate.IsZero() && !filter.ToDate.IsZero() {
		query["entryTime"] = bson.M{
			"$gte": filter.FromDate,
			"$lte": filter.ToDate,
		}
	} else if !filter.FromDate.IsZero() {
		query["entryTime"] = bson.M{"$gte": filter.FromDate}
	} else if !filter.ToDate.IsZero() {
		query["entryTime"] = bson.M{"$lte": filter.ToDate}
	}
	
	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}
	
	// Count total
	total, err := r.db.Database.Collection(PositionCollection).CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	
	// Set pagination
	skip := (page - 1) * limit
	
	// Find positions
	cursor, err := r.db.Database.Collection(PositionCollection).Find(
		ctx,
		query,
		options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "entryTime", Value: -1}}),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var positions []*models.Position
	if err = cursor.All(ctx, &positions); err != nil {
		return nil, 0, err
	}
	
	return positions, int(total), nil
}

// UserRepository provides methods for working with users
type UserRepository struct {
	db *MongoDB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *MongoDB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.PasswordChangedAt = time.Now()
	
	result, err := r.db.Database.Collection(UserCollection).InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var user models.User
	err = r.db.Database.Collection(UserCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	var user models.User
	err := r.db.Database.Collection(UserCollection).FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	var user models.User
	err := r.db.Database.Collection(UserCollection).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}
	
	user.UpdatedAt = time.Now()
	
	_, err = r.db.Database.Collection(UserCollection).ReplaceOne(
		ctx,
		bson.M{"_id": objectID},
		user,
	)
	
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.db.Database.Collection(UserCollection).DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// Find finds users based on filter
func (r *UserRepository) Find(filter models.UserFilter, page, limit int) ([]*models.User, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Build query
	query := bson.M{}
	
	if filter.Username != "" {
		query["username"] = bson.M{"$regex": filter.Username, "$options": "i"}
	}
	
	if filter.Email != "" {
		query["email"] = bson.M{"$regex": filter.Email, "$options": "i"}
	}
	
	if filter.Role != "" {
		query["role"] = filter.Role
	}
	
	if filter.Active != nil {
		query["active"] = *filter.Active
	}
	
	if !filter.FromDate.IsZero() && !filter.ToDate.IsZero() {
		query["createdAt"] = bson.M{
			"$gte": filter.FromDate,
			"$lte": filter.ToDate,
		}
	} else if !filter.FromDate.IsZero() {
		query["createdAt"] = bson.M{"$gte": filter.FromDate}
	} else if !filter.ToDate.IsZero() {
		query["createdAt"] = bson.M{"$lte": filter.ToDate}
	}
	
	// Count total
	total, err := r.db.Database.Collection(UserCollection).CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	
	// Set pagination
	skip := (page - 1) * limit
	
	// Find users
	cursor, err := r.db.Database.Collection(UserCollection).Find(
		ctx,
		query,
		options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "createdAt", Value: -1}}),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var users []*models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	
	return users, int(total), nil
}

// UserPreferenceRepository provides methods for working with user preferences
type UserPreferenceRepository struct {
	db *MongoDB
}

// NewUserPreferenceRepository creates a new UserPreferenceRepository
func NewUserPreferenceRepository(db *MongoDB) *UserPreferenceRepository {
	return &UserPreferenceRepository{db: db}
}

// GetByUserID retrieves user preferences by user ID
func (r *UserPreferenceRepository) GetByUserID(userID string) (*models.UserPreferences, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	var preferences models.UserPreferences
	err := r.db.Database.Collection(PreferenceCollection).FindOne(ctx, bson.M{"userId": userID}).Decode(&preferences)
	if err != nil {
		// If not found, create default preferences
		if err == mongo.ErrNoDocuments {
			preferences = *models.GetDefaultPreferences(userID)
			_, err = r.Create(&preferences)
			if err != nil {
				return nil, err
			}
			return &preferences, nil
		}
		return nil, err
	}
	
	return &preferences, nil
}

// Create creates new user preferences
func (r *UserPreferenceRepository) Create(preferences *models.UserPreferences) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	preferences.CreatedAt = time.Now()
	preferences.UpdatedAt = time.Now()
	
	result, err := r.db.Database.Collection(PreferenceCollection).InsertOne(ctx, preferences)
	if err != nil {
		return "", err
	}
	
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// Update updates user preferences
func (r *UserPreferenceRepository) Update(preferences *models.UserPreferences) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(preferences.ID)
	if err != nil {
		return err
	}
	
	preferences.UpdatedAt = time.Now()
	
	_, err = r.db.Database.Collection(PreferenceCollection).ReplaceOne(
		ctx,
		bson.M{"_id": objectID},
		preferences,
	)
	
	return err
}

// APIKeyRepository provides methods for working with API keys
type APIKeyRepository struct {
	db *MongoDB
}

// NewAPIKeyRepository creates a new APIKeyRepository
func NewAPIKeyRepository(db *MongoDB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// Create creates a new API key
func (r *APIKeyRepository) Create(apiKey *models.APIKey) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	apiKey.CreatedAt = time.Now()
	apiKey.UpdatedAt = time.Now()
	
	result, err := r.db.Database.Collection(APIKeyCollection).InsertOne(ctx, apiKey)
	if err != nil {
		return "", err
	}
	
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// GetByID retrieves an API key by ID
func (r *APIKeyRepository) GetByID(id string) (*models.APIKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var apiKey models.APIKey
	err = r.db.Database.Collection(APIKeyCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&apiKey)
	if err != nil {
		return nil, err
	}
	
	return &apiKey, nil
}

// GetByUserIDAndBroker retrieves API keys by user ID and broker
func (r *APIKeyRepository) GetByUserIDAndBroker(userID, broker string) (*models.APIKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	var apiKey models.APIKey
	err := r.db.Database.Collection(APIKeyCollection).FindOne(
		ctx,
		bson.M{
			"userId":     userID,
			"brokerName": broker,
			"isActive":   true,
		},
	).Decode(&apiKey)
	if err != nil {
		return nil, err
	}
	
	return &apiKey, nil
}

// Update updates an API key
func (r *APIKeyRepository) Update(apiKey *models.APIKey) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(apiKey.ID)
	if err != nil {
		return err
	}
	
	apiKey.UpdatedAt = time.Now()
	
	_, err = r.db.Database.Collection(APIKeyCollection).ReplaceOne(
		ctx,
		bson.M{"_id": objectID},
		apiKey,
	)
	
	return err
}

// Delete deletes an API key
func (r *APIKeyRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.db.Database.Collection(APIKeyCollection).DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// GetByUserID retrieves all API keys for a user
func (r *APIKeyRepository) GetByUserID(userID string) ([]*models.APIKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	cursor, err := r.db.Database.Collection(APIKeyCollection).Find(
		ctx,
		bson.M{"userId": userID},
		options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var apiKeys []*models.APIKey
	if err = cursor.All(ctx, &apiKeys); err != nil {
		return nil, err
	}
	
	return apiKeys, nil
}

// StrategyRepository provides methods for working with strategies
type StrategyRepository struct {
	db *MongoDB
}

// NewStrategyRepository creates a new StrategyRepository
func NewStrategyRepository(db *MongoDB) *StrategyRepository {
	return &StrategyRepository{db: db}
}

// Create creates a new strategy
func (r *StrategyRepository) Create(strategy *models.Strategy) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	strategy.CreatedAt = time.Now()
	strategy.UpdatedAt = time.Now()
	
	result, err := r.db.Database.Collection(StrategyCollection).InsertOne(ctx, strategy)
	if err != nil {
		return "", err
	}
	
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// GetByID retrieves a strategy by ID
func (r *StrategyRepository) GetByID(id string) (*models.Strategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var strategy models.Strategy
	err = r.db.Database.Collection(StrategyCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&strategy)
	if err != nil {
		return nil, err
	}
	
	return &strategy, nil
}

// Update updates a strategy
func (r *StrategyRepository) Update(strategy *models.Strategy) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(strategy.ID)
	if err != nil {
		return err
	}
	
	strategy.UpdatedAt = time.Now()
	
	_, err = r.db.Database.Collection(StrategyCollection).ReplaceOne(
		ctx,
		bson.M{"_id": objectID},
		strategy,
	)
	
	return err
}

// Delete deletes a strategy
func (r *StrategyRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.db.Database.Collection(StrategyCollection).DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// Find finds strategies based on filter
func (r *StrategyRepository) Find(filter models.StrategyFilter, page, limit int) ([]*models.Strategy, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Build query
	query := bson.M{}
	
	if filter.UserID != "" {
		query["userId"] = filter.UserID
	}
	
	if filter.Name != "" {
		query["name"] = bson.M{"$regex": filter.Name, "$options": "i"}
	}
	
	if filter.Type != "" {
		query["type"] = filter.Type
	}
	
	if filter.Tag != "" {
		query["tag"] = filter.Tag
	}
	
	if filter.Active != nil {
		query["active"] = *filter.Active
	}
	
	if filter.Symbol != "" {
		query["symbol"] = filter.Symbol
	}
	
	if filter.ProductType != "" {
		query["productType"] = filter.ProductType
	}
	
	if !filter.FromDate.IsZero() && !filter.ToDate.IsZero() {
		query["createdAt"] = bson.M{
			"$gte": filter.FromDate,
			"$lte": filter.ToDate,
		}
	} else if !filter.FromDate.IsZero() {
		query["createdAt"] = bson.M{"$gte": filter.FromDate}
	} else if !filter.ToDate.IsZero() {
		query["createdAt"] = bson.M{"$lte": filter.ToDate}
	}
	
	// Count total
	total, err := r.db.Database.Collection(StrategyCollection).CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	
	// Set pagination
	skip := (page - 1) * limit
	
	// Find strategies
	cursor, err := r.db.Database.Collection(StrategyCollection).Find(
		ctx,
		query,
		options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "createdAt", Value: -1}}),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var strategies []*models.Strategy
	if err = cursor.All(ctx, &strategies); err != nil {
		return nil, 0, err
	}
	
	return strategies, int(total), nil
}

// GetActiveStrategies retrieves all active strategies
func (r *StrategyRepository) GetActiveStrategies() ([]*models.Strategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	cursor, err := r.db.Database.Collection(StrategyCollection).Find(
		ctx,
		bson.M{"active": true},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var strategies []*models.Strategy
	if err = cursor.All(ctx, &strategies); err != nil {
		return nil, err
	}
	
	return strategies, nil
}

// PortfolioRepository provides methods for working with portfolios
type PortfolioRepository struct {
	db *MongoDB
}

// NewPortfolioRepository creates a new PortfolioRepository
func NewPortfolioRepository(db *MongoDB) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

// Create creates a new portfolio
func (r *PortfolioRepository) Create(portfolio *models.Portfolio) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	portfolio.CreatedAt = time.Now()
	portfolio.UpdatedAt = time.Now()
	
	result, err := r.db.Database.Collection(PortfolioCollection).InsertOne(ctx, portfolio)
	if err != nil {
		return "", err
	}
	
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// GetByID retrieves a portfolio by ID
func (r *PortfolioRepository) GetByID(id string) (*models.Portfolio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var portfolio models.Portfolio
	err = r.db.Database.Collection(PortfolioCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&portfolio)
	if err != nil {
		return nil, err
	}
	
	return &portfolio, nil
}

// Update updates a portfolio
func (r *PortfolioRepository) Update(portfolio *models.Portfolio) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(portfolio.ID)
	if err != nil {
		return err
	}
	
	portfolio.UpdatedAt = time.Now()
	
	_, err = r.db.Database.Collection(PortfolioCollection).ReplaceOne(
		ctx,
		bson.M{"_id": objectID},
		portfolio,
	)
	
	return err
}

// Delete deletes a portfolio
func (r *PortfolioRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.db.Database.Collection(PortfolioCollection).DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// Find finds portfolios based on filter
func (r *PortfolioRepository) Find(filter models.PortfolioFilter, page, limit int) ([]*models.Portfolio, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Build query
	query := bson.M{}
	
	if filter.UserID != "" {
		query["userId"] = filter.UserID
	}
	
	if filter.Name != "" {
		query["name"] = bson.M{"$regex": filter.Name, "$options": "i"}
	}
	
	if filter.StrategyID != "" {
		query["strategyId"] = filter.StrategyID
	}
	
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	
	if filter.Symbol != "" {
		query["symbol"] = filter.Symbol
	}
	
	if filter.ProductType != "" {
		query["productType"] = filter.ProductType
	}
	
	if !filter.FromDate.IsZero() && !filter.ToDate.IsZero() {
		query["createdAt"] = bson.M{
			"$gte": filter.FromDate,
			"$lte": filter.ToDate,
		}
	} else if !filter.FromDate.IsZero() {
		query["createdAt"] = bson.M{"$gte": filter.FromDate}
	} else if !filter.ToDate.IsZero() {
		query["createdAt"] = bson.M{"$lte": filter.ToDate}
	}
	
	// Count total
	total, err := r.db.Database.Collection(PortfolioCollection).CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	
	// Set pagination
	skip := (page - 1) * limit
	
	// Find portfolios
	cursor, err := r.db.Database.Collection(PortfolioCollection).Find(
		ctx,
		query,
		options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "createdAt", Value: -1}}),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var portfolios []*models.Portfolio
	if err = cursor.All(ctx, &portfolios); err != nil {
		return nil, 0, err
	}
	
	return portfolios, int(total), nil
}

// GetActivePortfolios retrieves all active portfolios
func (r *PortfolioRepository) GetActivePortfolios() ([]*models.Portfolio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	cursor, err := r.db.Database.Collection(PortfolioCollection).Find(
		ctx,
		bson.M{"status": models.PortfolioStatusActive},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var portfolios []*models.Portfolio
	if err = cursor.All(ctx, &portfolios); err != nil {
		return nil, err
	}
	
	return portfolios, nil
}

// GetPendingPortfolios retrieves all pending portfolios
func (r *PortfolioRepository) GetPendingPortfolios() ([]*models.Portfolio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	cursor, err := r.db.Database.Collection(PortfolioCollection).Find(
		ctx,
		bson.M{"status": models.PortfolioStatusPending},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var portfolios []*models.Portfolio
	if err = cursor.All(ctx, &portfolios); err != nil {
		return nil, err
	}
	
	return portfolios, nil
}
