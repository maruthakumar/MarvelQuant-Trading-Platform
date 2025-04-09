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

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	Create(order *models.Order) (*models.Order, error)
	GetByID(id string) (*models.Order, error)
	GetAll(filter models.OrderFilter, offset, limit int) ([]models.Order, int, error)
	Update(order *models.Order) (*models.Order, error)
	Delete(id string) error
}

// MongoOrderRepository implements OrderRepository using MongoDB
type MongoOrderRepository struct {
	collection *mongo.Collection
}

// NewMongoOrderRepository creates a new MongoOrderRepository
func NewMongoOrderRepository(db *mongo.Database) OrderRepository {
	return &MongoOrderRepository{
		collection: db.Collection("orders"),
	}
}

// Create adds a new order to the database
func (r *MongoOrderRepository) Create(order *models.Order) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate a new ID if not provided
	if order.ID == "" {
		order.ID = primitive.NewObjectID().Hex()
	}

	// Insert the order
	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetByID retrieves an order by ID
func (r *MongoOrderRepository) GetByID(id string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var order models.Order
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

// GetAll retrieves orders with filtering and pagination
func (r *MongoOrderRepository) GetAll(filter models.OrderFilter, offset, limit int) ([]models.Order, int, error) {
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
	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

// Update updates an existing order
func (r *MongoOrderRepository) Update(order *models.Order) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the order
	filter := bson.M{"_id": order.ID}
	update := bson.M{"$set": order}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// Delete removes an order from the database
func (r *MongoOrderRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("order not found")
	}

	return nil
}
