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

// MultilegRepository defines the interface for multileg strategy data operations
type MultilegRepository interface {
	// Strategy CRUD operations
	CreateStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error)
	GetStrategyByID(id string) (*models.MultilegStrategy, error)
	GetStrategiesByUser(userID string) ([]models.MultilegStrategy, error)
	GetStrategiesByPortfolio(portfolioID string) ([]models.MultilegStrategy, error)
	UpdateStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error)
	DeleteStrategy(id string) error
}

// MongoMultilegRepository implements MultilegRepository using MongoDB
type MongoMultilegRepository struct {
	db *mongo.Database
}

// NewMongoMultilegRepository creates a new MongoMultilegRepository
func NewMongoMultilegRepository(db *mongo.Database) MultilegRepository {
	return &MongoMultilegRepository{
		db: db,
	}
}

// CreateStrategy adds a new multileg strategy to the database
func (r *MongoMultilegRepository) CreateStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate a new ID if not provided
	if strategy.ID == "" {
		strategy.ID = primitive.NewObjectID().Hex()
	}

	// Generate IDs for legs if not provided
	for i := range strategy.Legs {
		if strategy.Legs[i].ID == "" {
			strategy.Legs[i].ID = primitive.NewObjectID().Hex()
		}
	}

	// Insert the strategy
	_, err := r.db.Collection("multileg_strategies").InsertOne(ctx, strategy)
	if err != nil {
		return nil, err
	}

	return strategy, nil
}

// GetStrategyByID retrieves a multileg strategy by ID
func (r *MongoMultilegRepository) GetStrategyByID(id string) (*models.MultilegStrategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var strategy models.MultilegStrategy
	filter := bson.M{"_id": id}

	err := r.db.Collection("multileg_strategies").FindOne(ctx, filter).Decode(&strategy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("strategy not found")
		}
		return nil, err
	}

	return &strategy, nil
}

// GetStrategiesByUser retrieves all multileg strategies for a user
func (r *MongoMultilegRepository) GetStrategiesByUser(userID string) ([]models.MultilegStrategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var strategies []models.MultilegStrategy
	filter := bson.M{"userId": userID}

	cursor, err := r.db.Collection("multileg_strategies").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &strategies); err != nil {
		return nil, err
	}

	return strategies, nil
}

// GetStrategiesByPortfolio retrieves all multileg strategies for a portfolio
func (r *MongoMultilegRepository) GetStrategiesByPortfolio(portfolioID string) ([]models.MultilegStrategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var strategies []models.MultilegStrategy
	filter := bson.M{"portfolioId": portfolioID}

	cursor, err := r.db.Collection("multileg_strategies").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &strategies); err != nil {
		return nil, err
	}

	return strategies, nil
}

// UpdateStrategy updates an existing multileg strategy
func (r *MongoMultilegRepository) UpdateStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the strategy
	filter := bson.M{"_id": strategy.ID}
	update := bson.M{"$set": strategy}

	_, err := r.db.Collection("multileg_strategies").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return strategy, nil
}

// DeleteStrategy removes a multileg strategy from the database
func (r *MongoMultilegRepository) DeleteStrategy(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := r.db.Collection("multileg_strategies").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("strategy not found")
	}

	return nil
}
