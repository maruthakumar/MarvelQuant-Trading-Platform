package repositories

import (
	"errors"
	"time"

	"github.com/trading-platform/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// PositionRepository defines the interface for position data operations
type PositionRepository interface {
	Create(position *models.Position) (*models.Position, error)
	GetByID(id string) (*models.Position, error)
	GetAll(filter models.PositionFilter, offset, limit int) ([]models.Position, int, error)
	Update(position *models.Position) (*models.Position, error)
	Delete(id string) error
}

// MongoPositionRepository implements PositionRepository using MongoDB
type MongoPositionRepository struct {
	collection *mongo.Collection
}

// NewMongoPositionRepository creates a new MongoPositionRepository
func NewMongoPositionRepository(db *mongo.Database) PositionRepository {
	return &MongoPositionRepository{
		collection: db.Collection("positions"),
	}
}

// Create adds a new position to the database
func (r *MongoPositionRepository) Create(position *models.Position) (*models.Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate a new ID if not provided
	if position.ID == "" {
		position.ID = primitive.NewObjectID().Hex()
	}

	// Insert the position
	_, err := r.collection.InsertOne(ctx, position)
	if err != nil {
		return nil, err
	}

	return position, nil
}

// GetByID retrieves a position by ID
func (r *MongoPositionRepository) GetByID(id string) (*models.Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var position models.Position
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&position)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("position not found")
		}
		return nil, err
	}

	return &position, nil
}

// GetAll retrieves positions with filtering and pagination
func (r *MongoPositionRepository) GetAll(filter models.PositionFilter, offset, limit int) ([]models.Position, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build the filter
	bsonFilter := bson.M{}
	if filter.UserID != "" {
		bsonFilter["userId"] = filter.UserID
	}
	if filter.Symbol != "" {
		bsonFilter["symbol"] = filter.Symbol
	}
	if filter.Status != "" {
		bsonFilter["status"] = filter.Status
	}
	if filter.Direction != "" {
		bsonFilter["direction"] = filter.Direction
	}
	if filter.ProductType != "" {
		bsonFilter["productType"] = filter.ProductType
	}
	if filter.InstrumentType != "" {
		bsonFilter["instrumentType"] = filter.InstrumentType
	}
	if filter.PortfolioID != "" {
		bsonFilter["portfolioId"] = filter.PortfolioID
	}
	if filter.StrategyID != "" {
		bsonFilter["strategyId"] = filter.StrategyID
	}
	if filter.OrderID != "" {
		bsonFilter["orderId"] = filter.OrderID
	}

	// Add date range filters if provided
	if !filter.FromDate.IsZero() || !filter.ToDate.IsZero() {
		dateFilter := bson.M{}
		if !filter.FromDate.IsZero() {
			dateFilter["$gte"] = filter.FromDate
		}
		if !filter.ToDate.IsZero() {
			dateFilter["$lte"] = filter.ToDate
		}
		bsonFilter["createdAt"] = dateFilter
	}

	// Add tags filter if provided
	if len(filter.Tags) > 0 {
		bsonFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	// Set up options for pagination and sorting
	findOptions := options.Find()
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"createdAt": -1}) // Sort by creation time, newest first

	// Execute the query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode the results
	var positions []models.Position
	if err := cursor.All(ctx, &positions); err != nil {
		return nil, 0, err
	}

	return positions, int(total), nil
}

// Update updates an existing position
func (r *MongoPositionRepository) Update(position *models.Position) (*models.Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the position
	filter := bson.M{"_id": position.ID}
	update := bson.M{"$set": position}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return position, nil
}

// Delete removes a position from the database
func (r *MongoPositionRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("position not found")
	}

	return nil
}
