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

// StrategyRepository defines the interface for strategy data operations
type StrategyRepository interface {
	// Strategy CRUD operations
	Create(strategy *models.Strategy) (*models.Strategy, error)
	GetByID(id string) (*models.Strategy, error)
	GetByUser(userID string) ([]models.Strategy, error)
	GetByTag(tag string) ([]models.Strategy, error)
	Update(strategy *models.Strategy) (*models.Strategy, error)
	Delete(id string) error
	
	// Strategy schedule operations
	SaveSchedule(schedule *models.StrategySchedule) error
	GetSchedule(strategyID string) (*models.StrategySchedule, error)
	DeleteSchedule(strategyID string) error
}

// MongoStrategyRepository implements StrategyRepository using MongoDB
type MongoStrategyRepository struct {
	db *mongo.Database
}

// NewMongoStrategyRepository creates a new MongoStrategyRepository
func NewMongoStrategyRepository(db *mongo.Database) StrategyRepository {
	return &MongoStrategyRepository{
		db: db,
	}
}

// Create adds a new strategy to the database
func (r *MongoStrategyRepository) Create(strategy *models.Strategy) (*models.Strategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate a new ID if not provided
	if strategy.ID == "" {
		strategy.ID = primitive.NewObjectID().Hex()
	}

	// Insert the strategy
	_, err := r.db.Collection("strategies").InsertOne(ctx, strategy)
	if err != nil {
		return nil, err
	}

	return strategy, nil
}

// GetByID retrieves a strategy by ID
func (r *MongoStrategyRepository) GetByID(id string) (*models.Strategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var strategy models.Strategy
	filter := bson.M{"_id": id}

	err := r.db.Collection("strategies").FindOne(ctx, filter).Decode(&strategy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("strategy not found")
		}
		return nil, err
	}

	return &strategy, nil
}

// GetByUser retrieves all strategies for a user
func (r *MongoStrategyRepository) GetByUser(userID string) ([]models.Strategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var strategies []models.Strategy
	filter := bson.M{"userId": userID}

	cursor, err := r.db.Collection("strategies").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &strategies); err != nil {
		return nil, err
	}

	return strategies, nil
}

// GetByTag retrieves all strategies with a specific tag
func (r *MongoStrategyRepository) GetByTag(tag string) ([]models.Strategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var strategies []models.Strategy
	filter := bson.M{"tags": tag}

	cursor, err := r.db.Collection("strategies").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &strategies); err != nil {
		return nil, err
	}

	return strategies, nil
}

// Update updates an existing strategy
func (r *MongoStrategyRepository) Update(strategy *models.Strategy) (*models.Strategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the strategy
	filter := bson.M{"_id": strategy.ID}
	update := bson.M{"$set": strategy}

	_, err := r.db.Collection("strategies").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return strategy, nil
}

// Delete removes a strategy from the database
func (r *MongoStrategyRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := r.db.Collection("strategies").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("strategy not found")
	}

	return nil
}

// SaveSchedule saves a strategy schedule
func (r *MongoStrategyRepository) SaveSchedule(schedule *models.StrategySchedule) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if schedule exists
	var existingSchedule models.StrategySchedule
	filter := bson.M{"strategyId": schedule.StrategyID}
	err := r.db.Collection("strategy_schedules").FindOne(ctx, filter).Decode(&existingSchedule)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Insert new schedule
			_, err := r.db.Collection("strategy_schedules").InsertOne(ctx, schedule)
			return err
		}
		return err
	}

	// Update existing schedule
	update := bson.M{"$set": schedule}
	_, err = r.db.Collection("strategy_schedules").UpdateOne(ctx, filter, update)
	return err
}

// GetSchedule retrieves a strategy schedule
func (r *MongoStrategyRepository) GetSchedule(strategyID string) (*models.StrategySchedule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var schedule models.StrategySchedule
	filter := bson.M{"strategyId": strategyID}

	err := r.db.Collection("strategy_schedules").FindOne(ctx, filter).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}

	return &schedule, nil
}

// DeleteSchedule deletes a strategy schedule
func (r *MongoStrategyRepository) DeleteSchedule(strategyID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"strategyId": strategyID}
	result, err := r.db.Collection("strategy_schedules").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("schedule not found")
	}

	return nil
}
