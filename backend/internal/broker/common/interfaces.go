// Package common provides interfaces and models shared across all broker implementations
package common

// BrokerClient defines the common interface for all broker implementations
type BrokerClient interface {
	// Authentication
	Login(credentials *Credentials) (*Session, error)
	Logout() error
	
	// Order Management
	PlaceOrder(order *Order) (*OrderResponse, error)
	ModifyOrder(order *ModifyOrder) (*OrderResponse, error)
	CancelOrder(orderID string, clientID string) (*OrderResponse, error)
	GetOrderBook(clientID string) (*OrderBook, error)
	
	// Portfolio Management
	GetPositions(clientID string) ([]Position, error)
	GetHoldings(clientID string) ([]Holding, error)
	
	// Market Data
	GetQuote(symbols []string) (map[string]Quote, error)
	SubscribeToQuotes(symbols []string) (chan Quote, error)
	UnsubscribeFromQuotes(symbols []string) error
}

// BrokerType represents the type of broker
type BrokerType string

// Broker types
const (
	BrokerTypeXTSPro    BrokerType = "XTS_PRO"
	BrokerTypeXTSClient BrokerType = "XTS_CLIENT"
	BrokerTypeZerodha   BrokerType = "ZERODHA"
)

// BrokerConfig contains configuration for broker clients
type BrokerConfig struct {
	BrokerType BrokerType
	XTSPro     *XTSProConfig
	XTSClient  *XTSClientConfig
	Zerodha    *ZerodhaConfig
}

// XTSProConfig contains configuration for XTS Pro
type XTSProConfig struct {
	APIKey     string
	SecretKey  string
	Source     string
	BaseURL    string
}

// XTSClientConfig contains configuration for XTS Client
type XTSClientConfig struct {
	APIKey     string
	SecretKey  string
	Source     string
	BaseURL    string
}

// ZerodhaConfig contains configuration for Zerodha
type ZerodhaConfig struct {
	APIKey      string
	APISecret   string
	RedirectURL string
	BaseURL     string
}
