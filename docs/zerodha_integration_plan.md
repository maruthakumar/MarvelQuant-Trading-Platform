# Zerodha Integration Plan

## Overview

This document outlines the approach for integrating Zerodha's trading platform into our trading system using their official Go client (gokiteconnect).

## Zerodha Go Client Analysis

The official Zerodha Go client (gokiteconnect) provides a comprehensive implementation with the following key components:

1. **Client Structure**: A clean, well-organized client structure with proper separation of concerns
2. **Authentication**: Token-based authentication with session management
3. **API Endpoints**: Comprehensive coverage of Zerodha's trading API endpoints
4. **WebSocket Support**: Real-time market data and order updates via WebSocket
5. **Error Handling**: Robust error handling with custom error types
6. **Models**: Well-defined data models for various trading entities

## Integration Strategy

### 1. Adapter Pattern

To integrate Zerodha alongside XTS implementations, we'll use the adapter pattern:

```go
// ZerodhaAdapter implements the common broker interface using gokiteconnect
type ZerodhaAdapter struct {
    client *kiteconnect.Client
    ticker *kiteticker.Ticker
}

// Implement common broker interface methods
func (z *ZerodhaAdapter) Login(credentials *models.Credentials) (*models.Session, error) {
    // Use Zerodha client to authenticate
    // Convert response to common Session model
}

func (z *ZerodhaAdapter) PlaceOrder(order *models.Order) (*models.OrderResponse, error) {
    // Convert common Order model to Zerodha-specific order
    // Use Zerodha client to place order
    // Convert response to common OrderResponse model
}

// ... other interface methods
```

### 2. Common Interface

Extend our common broker interface to accommodate Zerodha-specific features:

```go
type BrokerClient interface {
    // Authentication
    Login(credentials *models.Credentials) (*models.Session, error)
    Logout() error
    
    // Order Management
    PlaceOrder(order *models.Order) (*models.OrderResponse, error)
    ModifyOrder(order *models.ModifyOrder) (*models.OrderResponse, error)
    CancelOrder(orderID string, clientID string) (*models.OrderResponse, error)
    GetOrderBook(clientID string) (*models.OrderBook, error)
    
    // Portfolio Management
    GetPositions(clientID string) ([]models.Position, error)
    GetHoldings(clientID string) ([]models.Holding, error)
    
    // Market Data
    GetQuote(symbols []string) (map[string]models.Quote, error)
    SubscribeToQuotes(symbols []string) (chan models.Quote, error)
    UnsubscribeFromQuotes(symbols []string) error
    
    // Broker-specific methods can be handled through type assertions
    // For example: if zerodhaClient, ok := client.(*ZerodhaAdapter); ok { zerodhaClient.PlaceGTT(...) }
}
```

### 3. Model Mapping

Create mapping functions between our common models and Zerodha-specific models:

```go
// Convert common Order to Zerodha Order
func convertToZerodhaOrder(order *models.Order) kiteconnect.OrderParams {
    return kiteconnect.OrderParams{
        Exchange:        mapExchange(order.ExchangeSegment),
        Tradingsymbol:   order.Symbol,
        TransactionType: mapTransactionType(order.OrderSide),
        Quantity:        order.OrderQuantity,
        Price:           order.LimitPrice,
        Product:         mapProduct(order.ProductType),
        OrderType:       mapOrderType(order.OrderType),
        Validity:        mapValidity(order.TimeInForce),
        // ... other fields
    }
}

// Convert Zerodha Order to common OrderResponse
func convertFromZerodhaOrder(order kiteconnect.Order) *models.OrderResponse {
    return &models.OrderResponse{
        OrderID:         order.OrderID,
        ExchangeOrderID: order.ExchangeOrderID,
        Status:          mapOrderStatus(order.Status),
        // ... other fields
    }
}
```

### 4. Configuration Management

Extend configuration to support Zerodha-specific settings:

```go
type ZerodhaConfig struct {
    APIKey      string
    APISecret   string
    RedirectURL string
    BaseURL     string // Optional, for custom endpoints
}

type BrokerConfig struct {
    BrokerType string // "XTS_PRO", "XTS_CLIENT", "ZERODHA"
    
    // Broker-specific configs
    XTSPro    *XTSProConfig
    XTSClient *XTSClientConfig
    Zerodha   *ZerodhaConfig
}
```

### 5. Factory Pattern

Use a factory pattern to create the appropriate client:

```go
func NewBrokerClient(config *BrokerConfig) (BrokerClient, error) {
    switch config.BrokerType {
    case "XTS_PRO":
        return NewXTSProClient(config.XTSPro)
    case "XTS_CLIENT":
        return NewXTSClientImpl(config.XTSClient)
    case "ZERODHA":
        return NewZerodhaAdapter(config.Zerodha)
    default:
        return nil, errors.New("unsupported broker type")
    }
}
```

## Implementation Plan

### Phase 1: Core Integration

1. **Create Adapter Structure**: Implement the ZerodhaAdapter struct
2. **Implement Authentication**: Login, logout, and session management
3. **Implement Order Management**: Place, modify, cancel orders
4. **Implement Portfolio Management**: Get positions and holdings
5. **Implement Market Data**: Get quotes and historical data

### Phase 2: Real-time Data

1. **Implement WebSocket Client**: Connect to Zerodha's WebSocket API
2. **Implement Market Data Streaming**: Subscribe to real-time quotes
3. **Implement Order Updates**: Handle real-time order status updates

### Phase 3: Advanced Features

1. **Implement GTT Orders**: Good Till Triggered orders
2. **Implement Basket Orders**: Multiple orders in a single request
3. **Implement Margins API**: Pre-order margin calculations

## Directory Structure

```
trading-platform/
├── backend/
│   ├── internal/
│   │   ├── broker/
│   │   │   ├── common/       # Common interfaces and models
│   │   │   ├── xts/          # XTS implementations
│   │   │   │   ├── pro/      # XTS PRO implementation
│   │   │   │   └── client/   # XTS Client implementation
│   │   │   └── zerodha/      # Zerodha implementation
│   │   │       ├── adapter.go    # Main adapter implementation
│   │   │       ├── mapper.go     # Model mapping functions
│   │   │       ├── websocket.go  # WebSocket client
│   │   │       └── config.go     # Zerodha-specific config
```

## Authentication Flow

Zerodha uses a slightly different authentication flow compared to XTS:

1. **Generate Login URL**: Create a URL for the user to visit
2. **User Authentication**: User logs in on Zerodha's website
3. **Receive Request Token**: Zerodha redirects to our app with a request token
4. **Exchange for Access Token**: Exchange request token for an access token
5. **Use Access Token**: Use the access token for API requests

This flow will need to be adapted to fit our common authentication interface.

## Error Handling

Zerodha's error handling will be mapped to our common error types:

```go
func mapZerodhaError(err error) error {
    if kiteErr, ok := err.(*kiteconnect.Error); ok {
        switch kiteErr.Code {
        case "TokenException":
            return errors.Wrap(err, "auth_error", "Authentication failed", kiteErr.Code)
        case "PermissionException":
            return errors.Wrap(err, "permission_error", "Permission denied", kiteErr.Code)
        // ... other error mappings
        default:
            return errors.Wrap(err, "zerodha_error", kiteErr.Message, kiteErr.Code)
        }
    }
    return err
}
```

## Testing Strategy

1. **Unit Tests**: Test adapter methods in isolation
2. **Integration Tests**: Test with Zerodha's sandbox environment
3. **Mock Tests**: Use mock responses for testing error scenarios

## Conclusion

The Zerodha integration will follow the adapter pattern to fit into our existing architecture while maintaining the unique features of the Zerodha platform. By using a common interface and factory pattern, we can provide a consistent experience across different brokers while allowing for broker-specific functionality when needed.
