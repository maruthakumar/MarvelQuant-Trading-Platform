# Multi-Broker Integration Documentation

## Overview

This document provides comprehensive documentation for the multi-broker integration implementation in the Trading Platform. The implementation supports multiple broker types including XTS PRO, XTS Client, and Zerodha, with a unified API layer that abstracts away the differences between brokers.

## Architecture

The multi-broker integration follows a layered architecture:

1. **Common Interface Layer**: Defines the contract that all broker implementations must follow
2. **Broker-Specific Implementations**: Implements the common interface for each broker type
3. **Factory Layer**: Creates the appropriate broker client based on configuration
4. **Unified API Layer**: Provides a consistent API for interacting with any broker

### Directory Structure

```
/backend/internal/
├── broker/
│   ├── common/           # Common interfaces and models
│   ├── factory/          # Broker client factory
│   ├── integration/      # Integration tests
│   ├── xts/
│   │   ├── pro/          # XTS PRO implementation
│   │   └── client/       # XTS Client implementation
│   └── zerodha/          # Zerodha implementation
└── api/                  # Unified API layer
```

## Common Interface

The common interface (`BrokerClient`) defines the contract that all broker implementations must follow:

```go
type BrokerClient interface {
    Login(credentials *Credentials) (*Session, error)
    Logout() error
    PlaceOrder(order *Order) (*OrderResponse, error)
    ModifyOrder(order *ModifyOrder) (*OrderResponse, error)
    CancelOrder(orderID string, clientID string) (*OrderResponse, error)
    GetOrderBook(clientID string) (*OrderBook, error)
    GetPositions(clientID string) ([]Position, error)
    GetHoldings(clientID string) ([]Holding, error)
    GetQuote(symbols []string) (map[string]Quote, error)
    SubscribeToQuotes(symbols []string) (chan Quote, error)
    UnsubscribeFromQuotes(symbols []string) error
}
```

## Broker Implementations

### XTS PRO

The XTS PRO implementation is designed for users of the XTS PRO platform. It provides high-performance trading capabilities with a focus on low latency.

Key features:
- Direct integration with XTS PRO API
- Support for all order types and product types
- Real-time market data via WebSocket
- Comprehensive error handling and recovery mechanisms

### XTS Client

The XTS Client implementation is designed for users of the standard XTS Client platform. It extends the XTS PRO implementation with additional features specific to the XTS Client API.

Key features:
- All features of XTS PRO
- Support for dealer operations (placing orders on behalf of clients)
- Additional API source parameter for order placement
- Dealer-specific endpoints for managing client orders and positions

### Zerodha

The Zerodha implementation integrates with the official Zerodha Kite Connect API. It provides access to Zerodha's trading platform with a consistent interface matching the other broker implementations.

Key features:
- Integration with official Zerodha Kite Connect Go client
- Support for all Zerodha-specific order types and product types
- Mapping between common models and Zerodha-specific models
- Two-step authentication process

## Factory Layer

The factory layer provides a simple way to create the appropriate broker client based on configuration:

```go
func NewBrokerClient(config *BrokerConfig) (BrokerClient, error) {
    switch config.BrokerType {
    case BrokerTypeXTSPro:
        return xts.NewXTSProClient(config.XTSPro)
    case BrokerTypeXTSClient:
        return client.NewXTSClientImpl(config.XTSClient)
    case BrokerTypeZerodha:
        return zerodha.NewZerodhaAdapter(config.Zerodha)
    default:
        return nil, fmt.Errorf("unsupported broker type: %s", config.BrokerType)
    }
}
```

## Unified API Layer

The unified API layer (`BrokerManager`) provides a consistent API for interacting with any broker. It handles client registration, authentication, and provides methods for all trading operations.

Key features:
- Thread-safe client management
- User session tracking
- Consistent API for all broker operations
- Support for dealer-specific operations

Example usage:

```go
// Create a broker manager
manager := api.NewBrokerManager()

// Register a broker
config := &common.BrokerConfig{
    BrokerType: common.BrokerTypeXTSClient,
    XTSClient: &common.XTSClientConfig{
        APIKey:    "your_api_key",
        SecretKey: "your_secret_key",
        Source:    "WEBAPI",
    },
}
manager.RegisterBroker("client1", config)

// Login
credentials := &common.Credentials{
    UserID:   "your_user_id",
    Password: "your_password",
}
session, err := manager.Login("client1", credentials)

// Place an order
order := &common.Order{
    ExchangeSegment:       "NSECM",
    TradingSymbol:         "RELIANCE",
    OrderSide:             "BUY",
    OrderQuantity:         10,
    ProductType:           "MIS",
    OrderType:             "LIMIT",
    TimeInForce:           "DAY",
    LimitPrice:            2000.0,
    OrderUniqueIdentifier: "test123",
}
orderResponse, err := manager.PlaceOrder(session.UserID, order)
```

## Testing

The implementation includes comprehensive testing at multiple levels:

1. **Unit Tests**: Test individual components in isolation
2. **Integration Tests**: Test the integration between components
3. **Mock Tests**: Test the unified API layer with mock clients

Integration tests are designed to be skipped by default and only run when specific environment variables are set:

```go
// Skip if not running integration tests
if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
    t.Skip("Skipping integration tests. Set RUN_INTEGRATION_TESTS=1 to run.")
}
```

## Dealer Operations

The implementation includes special support for dealer operations, which allow a dealer to place orders on behalf of clients. This is particularly important for the XTS Client implementation.

The unified API layer provides methods for dealer operations:

```go
// Place an order on behalf of a client
orderResponse, err := manager.PlaceDealerOrder(dealerUserID, targetClientID, order)

// Get the order book for a client
orderBook, err := manager.GetDealerOrderBook(dealerUserID, targetClientID)

// Get the positions for a client
positions, err := manager.GetDealerPositions(dealerUserID, targetClientID)
```

## WebSocket Implementation

The implementation includes support for real-time market data via WebSocket. However, the WebSocket implementation is currently a placeholder and needs to be completed in a future update.

The interface is defined and ready to use:

```go
// Subscribe to real-time quotes
quoteChan, err := client.SubscribeToQuotes(symbols)

// Unsubscribe from real-time quotes
err := client.UnsubscribeFromQuotes(symbols)
```

## Error Handling

The implementation includes comprehensive error handling at all levels:

1. **Input Validation**: Validate all input parameters before making API calls
2. **API Error Handling**: Handle errors returned by the broker APIs
3. **Network Error Handling**: Handle network errors and timeouts
4. **Recovery Mechanisms**: Implement recovery mechanisms for connection failures

## Future Enhancements

1. **Complete WebSocket Implementation**: Implement real-time market data via WebSocket
2. **Add More Broker Types**: Add support for additional broker types
3. **Implement Rate Limiting**: Add rate limiting to prevent API throttling
4. **Add Caching**: Add caching for frequently accessed data
5. **Implement Circuit Breakers**: Add circuit breakers to prevent cascading failures

## Conclusion

The multi-broker integration implementation provides a flexible and extensible framework for integrating with multiple broker types. The unified API layer abstracts away the differences between brokers, providing a consistent interface for all trading operations.
