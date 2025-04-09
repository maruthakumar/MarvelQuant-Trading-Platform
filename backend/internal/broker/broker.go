package broker

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// OrderType represents the type of order
type OrderType string

const (
	OrderTypeMarket     OrderType = "MARKET"
	OrderTypeLimit      OrderType = "LIMIT"
	OrderTypeStopLoss   OrderType = "SL"
	OrderTypeStopLimit  OrderType = "SL_LIMIT"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeBuy  TransactionType = "BUY"
	TransactionTypeSell TransactionType = "SELL"
)

// ProductType represents the type of product
type ProductType string

const (
	ProductTypeNRML ProductType = "NRML"
	ProductTypeMIS  ProductType = "MIS"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusOpen       OrderStatus = "OPEN"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
	OrderStatusRejected   OrderStatus = "REJECTED"
)

// Order represents an order in the system
type Order struct {
	ID              string          `json:"id"`
	BrokerOrderID   string          `json:"broker_order_id"`
	Symbol          string          `json:"symbol"`
	Exchange        string          `json:"exchange"`
	OrderType       OrderType       `json:"order_type"`
	TransactionType TransactionType `json:"transaction_type"`
	ProductType     ProductType     `json:"product_type"`
	Quantity        int             `json:"quantity"`
	Price           float64         `json:"price,omitempty"`
	TriggerPrice    float64         `json:"trigger_price,omitempty"`
	Status          OrderStatus     `json:"status"`
	Message         string          `json:"message,omitempty"`
	OrderTimestamp  time.Time       `json:"order_timestamp"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// Quote represents a market quote
type Quote struct {
	Symbol        string    `json:"symbol"`
	Exchange      string    `json:"exchange"`
	LastPrice     float64   `json:"last_price"`
	BidPrice      float64   `json:"bid_price"`
	AskPrice      float64   `json:"ask_price"`
	Volume        int       `json:"volume"`
	OpenInterest  int       `json:"open_interest"`
	Timestamp     time.Time `json:"timestamp"`
}

// Position represents a trading position
type Position struct {
	Symbol          string          `json:"symbol"`
	Exchange        string          `json:"exchange"`
	ProductType     ProductType     `json:"product_type"`
	Quantity        int             `json:"quantity"`
	AveragePrice    float64         `json:"average_price"`
	LastPrice       float64         `json:"last_price"`
	PnL             float64         `json:"pnl"`
	RealizedPnL     float64         `json:"realized_pnl"`
	UnrealizedPnL   float64         `json:"unrealized_pnl"`
	Timestamp       time.Time       `json:"timestamp"`
}

// BrokerConfig represents the configuration for a broker
type BrokerConfig struct {
	APIKey      string `json:"api_key"`
	APISecret   string `json:"api_secret"`
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	Endpoint    string `json:"endpoint"`
}

// OrderRequest represents a request to place an order
type OrderRequest struct {
	Symbol          string          `json:"symbol"`
	Exchange        string          `json:"exchange"`
	OrderType       OrderType       `json:"order_type"`
	TransactionType TransactionType `json:"transaction_type"`
	ProductType     ProductType     `json:"product_type"`
	Quantity        int             `json:"quantity"`
	Price           float64         `json:"price,omitempty"`
	TriggerPrice    float64         `json:"trigger_price,omitempty"`
}

// OrderResponse represents a response from placing an order
type OrderResponse struct {
	Success      bool   `json:"success"`
	OrderID      string `json:"order_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// Broker is an interface for broker operations
type Broker interface {
	// Initialize initializes the broker with configuration
	Initialize(config BrokerConfig) error
	
	// PlaceOrder places an order
	PlaceOrder(ctx context.Context, request OrderRequest) (*OrderResponse, error)
	
	// ModifyOrder modifies an existing order
	ModifyOrder(ctx context.Context, orderID string, request OrderRequest) (*OrderResponse, error)
	
	// CancelOrder cancels an order
	CancelOrder(ctx context.Context, orderID string) (*OrderResponse, error)
	
	// GetOrder gets an order by ID
	GetOrder(ctx context.Context, orderID string) (*Order, error)
	
	// GetOrders gets all orders
	GetOrders(ctx context.Context) ([]Order, error)
	
	// GetPositions gets all positions
	GetPositions(ctx context.Context) ([]Position, error)
	
	// GetQuote gets a quote for a symbol
	GetQuote(ctx context.Context, symbol, exchange string) (*Quote, error)
	
	// SubscribeQuotes subscribes to quotes for symbols
	SubscribeQuotes(ctx context.Context, symbols []string, exchange string) error
	
	// UnsubscribeQuotes unsubscribes from quotes for symbols
	UnsubscribeQuotes(ctx context.Context, symbols []string, exchange string) error
	
	// Close closes the broker connection
	Close() error
}

// BrokerFactory creates brokers
type BrokerFactory struct {
	brokers map[string]Broker
}

// NewBrokerFactory creates a new broker factory
func NewBrokerFactory() *BrokerFactory {
	return &BrokerFactory{
		brokers: make(map[string]Broker),
	}
}

// RegisterBroker registers a broker with the factory
func (f *BrokerFactory) RegisterBroker(name string, broker Broker) {
	f.brokers[name] = broker
}

// GetBroker gets a broker by name
func (f *BrokerFactory) GetBroker(name string) (Broker, error) {
	broker, ok := f.brokers[name]
	if !ok {
		return nil, fmt.Errorf("broker not found: %s", name)
	}
	return broker, nil
}

// ErrNotImplemented is returned when a method is not implemented
var ErrNotImplemented = errors.New("method not implemented")
