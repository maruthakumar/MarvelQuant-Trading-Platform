package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trading-platform/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockCollection is a mock implementation of mongo.Collection
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(int64), args.Error(1)
}

// MockSingleResult is a mock implementation of mongo.SingleResult
type MockSingleResult struct {
	mock.Mock
}

func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockSingleResult) Err() error {
	args := m.Called()
	return args.Error(0)
}

// MockCursor is a mock implementation of mongo.Cursor
type MockCursor struct {
	mock.Mock
}

func (m *MockCursor) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockCursor) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockCursor) All(ctx context.Context, result interface{}) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}

func (m *MockCursor) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCursor) Err() error {
	args := m.Called()
	return args.Error(0)
}

func TestMongoOrderRepositoryCreate(t *testing.T) {
	// Create a mock collection
	mockCollection := new(MockCollection)
	
	// Create a sample order
	order := &models.Order{
		UserID:         "user123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		OrderType:      models.OrderTypeLimit,
		Direction:      models.OrderDirectionBuy,
		Quantity:       10,
		Price:          500.50,
		Status:         models.OrderStatusPending,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
	}
	
	// Set up the mock collection expectations
	mockCollection.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Order"), mock.Anything).Return(&mongo.InsertOneResult{InsertedID: "order123"}, nil)
	
	// Create the repository with the mock collection
	repo := &MongoOrderRepository{
		collection: mockCollection,
	}
	
	// Call the repository method
	createdOrder, err := repo.Create(order)
	
	// Check the result
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)
	assert.NotEmpty(t, createdOrder.ID)
	
	// Verify that the mock collection was called
	mockCollection.AssertExpectations(t)
}

func TestMongoOrderRepositoryGetByID(t *testing.T) {
	// Create a mock collection and single result
	mockCollection := new(MockCollection)
	mockSingleResult := new(MockSingleResult)
	
	// Create a sample order
	order := &models.Order{
		ID:             "order123",
		UserID:         "user123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		OrderType:      models.OrderTypeLimit,
		Direction:      models.OrderDirectionBuy,
		Quantity:       10,
		Price:          500.50,
		Status:         models.OrderStatusPending,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
	}
	
	// Set up the mock expectations
	mockSingleResult.On("Decode", mock.AnythingOfType("*models.Order")).Run(func(args mock.Arguments) {
		// Copy the order data to the provided pointer
		orderArg := args.Get(0).(*models.Order)
		*orderArg = *order
	}).Return(nil)
	
	mockCollection.On("FindOne", mock.Anything, bson.M{"_id": "order123"}, mock.Anything).Return(mockSingleResult)
	
	// Set up mock for not found case
	mockNotFoundResult := new(MockSingleResult)
	mockNotFoundResult.On("Decode", mock.Anything).Return(mongo.ErrNoDocuments)
	mockCollection.On("FindOne", mock.Anything, bson.M{"_id": "nonexistent"}, mock.Anything).Return(mockNotFoundResult)
	
	// Create the repository with the mock collection
	repo := &MongoOrderRepository{
		collection: mockCollection,
	}
	
	// Test successful retrieval
	retrievedOrder, err := repo.GetByID("order123")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedOrder)
	assert.Equal(t, order.ID, retrievedOrder.ID)
	assert.Equal(t, order.UserID, retrievedOrder.UserID)
	assert.Equal(t, order.Symbol, retrievedOrder.Symbol)
	
	// Test not found case
	retrievedOrder, err = repo.GetByID("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, retrievedOrder)
	assert.Contains(t, err.Error(), "order not found")
	
	// Verify that the mock collection was called
	mockCollection.AssertExpectations(t)
	mockSingleResult.AssertExpectations(t)
	mockNotFoundResult.AssertExpectations(t)
}

func TestMongoOrderRepositoryGetAll(t *testing.T) {
	// Create a mock collection and cursor
	mockCollection := new(MockCollection)
	mockCursor := new(MockCursor)
	
	// Create sample orders
	orders := []models.Order{
		{
			ID:             "order123",
			UserID:         "user123",
			Symbol:         "NIFTY",
			Exchange:       "NSE",
			OrderType:      models.OrderTypeLimit,
			Direction:      models.OrderDirectionBuy,
			Quantity:       10,
			Price:          500.50,
			Status:         models.OrderStatusPending,
			ProductType:    models.ProductTypeMIS,
			InstrumentType: models.InstrumentTypeOption,
		},
		{
			ID:             "order456",
			UserID:         "user123",
			Symbol:         "BANKNIFTY",
			Exchange:       "NSE",
			OrderType:      models.OrderTypeMarket,
			Direction:      models.OrderDirectionSell,
			Quantity:       5,
			Price:          0,
			Status:         models.OrderStatusExecuted,
			ProductType:    models.ProductTypeMIS,
			InstrumentType: models.InstrumentTypeFuture,
		},
	}
	
	// Set up the mock expectations
	mockCursor.On("All", mock.Anything, mock.AnythingOfType("*[]models.Order")).Run(func(args mock.Arguments) {
		// Copy the orders data to the provided slice
		ordersArg := args.Get(1).(*[]models.Order)
		*ordersArg = orders
	}).Return(nil)
	mockCursor.On("Close", mock.Anything).Return(nil)
	
	mockCollection.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(mockCursor, nil)
	mockCollection.On("CountDocuments", mock.Anything, mock.Anything, mock.Anything).Return(int64(2), nil)
	
	// Create the repository with the mock collection
	repo := &MongoOrderRepository{
		collection: mockCollection,
	}
	
	// Test retrieval with filter
	filter := models.OrderFilter{
		UserID: "user123",
		Symbol: "NIFTY",
	}
	retrievedOrders, total, err := repo.GetAll(filter, 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(retrievedOrders))
	assert.Equal(t, 2, total)
	
	// Verify that the mock collection was called
	mockCollection.AssertExpectations(t)
	mockCursor.AssertExpectations(t)
}

func TestMongoOrderRepositoryUpdate(t *testing.T) {
	// Create a mock collection
	mockCollection := new(MockCollection)
	
	// Create a sample order
	order := &models.Order{
		ID:             "order123",
		UserID:         "user123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		OrderType:      models.OrderTypeLimit,
		Direction:      models.OrderDirectionBuy,
		Quantity:       10,
		Price:          500.50,
		Status:         models.OrderStatusPending,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
		UpdatedAt:      time.Now(),
	}
	
	// Set up the mock collection expectations
	mockCollection.On("UpdateOne", mock.Anything, bson.M{"_id": "order123"}, bson.M{"$set": order}, mock.Anything).Return(&mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil)
	
	// Create the repository with the mock collection
	repo := &MongoOrderRepository{
		collection: mockCollection,
	}
	
	// Call the repository method
	updatedOrder, err := repo.Update(order)
	
	// Check the result
	assert.NoError(t, err)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, order.ID, updatedOrder.ID)
	
	// Verify that the mock collection was called
	mockCollection.AssertExpectations(t)
}

func TestMongoOrderRepositoryDelete(t *testing.T) {
	// Create a mock collection
	mockCollection := new(MockCollection)
	
	// Set up the mock collection expectations
	mockCollection.On("DeleteOne", mock.Anything, bson.M{"_id": "order123"}, mock.Anything).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
	mockCollection.On("DeleteOne", mock.Anything, bson.M{"_id": "nonexistent"}, mock.Anything).Return(&mongo.DeleteResult{DeletedCount: 0}, nil)
	
	// Create the repository with the mock collection
	repo := &MongoOrderRepository{
		collection: mockCollection,
	}
	
	// Test successful deletion
	err := repo.Delete("order123")
	assert.NoError(t, err)
	
	// Test not found case
	err = repo.Delete("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order not found")
	
	// Verify that the mock collection was called
	mockCollection.AssertExpectations(t)
}
