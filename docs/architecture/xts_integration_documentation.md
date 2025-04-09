# XTS Integration Documentation

## Overview

This document provides comprehensive documentation for the XTS trading platform integration implemented in Go. The implementation replaces the Python SDK with a native Go implementation to eliminate latency issues while maintaining all the functionality of the original Python SDK.

## Architecture

The XTS integration follows a modular architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────────┐
│                       XTS Go Integration                         │
│                                                                  │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────────┐ │
│  │                │  │                │  │                    │ │
│  │  REST Client   │  │  Market Data   │  │  Order Management  │ │
│  │  Module        │  │  WebSocket     │  │  WebSocket         │ │
│  │                │  │                │  │                    │ │
│  └────────────────┘  └────────────────┘  └────────────────────┘ │
│          │                   │                    │              │
│          ▼                   ▼                    ▼              │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────────┐ │
│  │                │  │                │  │                    │ │
│  │  Auth Service  │  │  Market Data   │  │  Order Service     │ │
│  │                │  │  Service       │  │                    │ │
│  └────────────────┘  └────────────────┘  └────────────────────┘ │
│          │                   │                    │              │
│          └───────────────────┼────────────────────┘              │
│                              │                                   │
│                              ▼                                   │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────────┐ │
│  │                │  │                │  │                    │ │
│  │  Models        │  │  Error         │  │  Recovery          │ │
│  │                │  │  Handling      │  │  Mechanisms        │ │
│  │                │  │                │  │                    │ │
│  └────────────────┘  └────────────────┘  └────────────────────┘ │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Components

### 1. Models (`/internal/xts/models`)

The models package defines all data structures used in the XTS integration:

- `Session`: Represents an authenticated session
- `Order`: Represents an order to be placed
- `OrderResponse`: Represents the response from placing an order
- `BracketOrder`: Represents a bracket order
- `CoverOrder`: Represents a cover order
- `ModifyOrder`: Represents an order modification request
- `Quote`: Represents a market quote
- `Position`: Represents a trading position
- `Instrument`: Represents a trading instrument
- `OrderBook`: Represents the order book

### 2. Configuration (`/internal/xts/config`)

The config package provides configuration management:

- `XTSConfig`: Configuration for XTS API connectivity
- `Routes`: Defines API endpoints for XTS

### 3. Error Handling (`/internal/xts/errors`)

The errors package provides comprehensive error handling:

- `XTSError`: Custom error type for XTS API errors
- Error checking functions: `IsAuthError`, `IsNetworkError`, etc.
- Standard error types for different categories of errors

### 4. Recovery Mechanisms (`/internal/xts/recovery`)

The recovery package implements resilience patterns:

- `RetryWithBackoff`: Retry mechanism with exponential backoff
- `CircuitBreaker`: Circuit breaker pattern for fault tolerance
- `RateLimiter`: Rate limiting to prevent API throttling

### 5. REST Client (`/internal/xts/rest`)

The rest package implements the HTTP client for XTS API:

- `Client`: Main client for REST API communication
- Authentication methods: `Login`, `Logout`
- Order methods: `PlaceOrder`, `ModifyOrder`, `CancelOrder`
- Market data methods: `GetQuotes`, `GetPositions`, etc.

### 6. WebSocket Clients (`/internal/xts/websocket`)

The websocket package implements real-time data streaming:

- `MarketDataClient`: Client for market data WebSocket
- `OrderClient`: Client for order updates WebSocket

### 7. Services (`/internal/xts/service`)

The service package provides high-level business logic:

- `MarketDataService`: Service for market data operations
- `OrderService`: Service for order management operations
- `AuthService`: Service for authentication operations

### 8. API Integration (`/internal/api`)

The api package integrates XTS with the backend gateway:

- `XTSController`: Controller for XTS API endpoints
- Route registration and request handling

## Usage Examples

### Initialization

```go
// Create configuration
cfg := config.NewXTSConfig()
cfg.BaseURL = "https://xts-api.trading.com"
cfg.APIKey = "your-api-key"
cfg.SecretKey = "your-secret-key"
cfg.Source = "WEB"

// Create REST client
restClient, err := rest.NewClient(cfg)
if err != nil {
    log.Fatalf("Failed to create REST client: %v", err)
}

// Create services
authService := service.NewAuthService(restClient, cfg)
marketService := service.NewMarketDataService(restClient, cfg)
orderService := service.NewOrderService(restClient, cfg)
```

### Authentication

```go
// Login
session, err := restClient.Login()
if err != nil {
    log.Fatalf("Login failed: %v", err)
}

// Start services
if err := marketService.Start(session.Token, session.UserID); err != nil {
    log.Fatalf("Failed to start market data service: %v", err)
}

if err := orderService.Start(session.Token, session.UserID); err != nil {
    log.Fatalf("Failed to start order service: %v", err)
}
```

### Placing an Order

```go
// Create order
order := &models.Order{
    ExchangeSegment:      "NSECM",
    ExchangeInstrumentID: "RELIANCE",
    ProductType:          models.ProductNRML,
    OrderType:            models.OrderTypeLimit,
    OrderSide:            models.TransactionTypeBuy,
    TimeInForce:          models.ValidityDay,
    DisclosedQuantity:    0,
    OrderQuantity:        100,
    LimitPrice:           2500.0,
    StopPrice:            0.0,
    OrderUniqueIdentifier: "test-order-123",
}

// Place order
response, err := orderService.PlaceOrder(order)
if err != nil {
    log.Printf("Order placement failed: %v", err)
    return
}

log.Printf("Order placed successfully: %s", response.OrderID)
```

### Subscribing to Market Data

```go
// Create instruments
instruments := []models.Instrument{
    {
        ExchangeSegment:      "NSECM",
        ExchangeInstrumentID: "RELIANCE",
    },
    {
        ExchangeSegment:      "NSECM",
        ExchangeInstrumentID: "INFY",
    },
}

// Subscribe to quotes
quoteChan, err := marketService.SubscribeToQuotes(instruments)
if err != nil {
    log.Printf("Subscription failed: %v", err)
    return
}

// Process quotes
go func() {
    for quote := range quoteChan {
        log.Printf("Received quote: %s %s - LTP: %.2f", 
            quote.ExchangeSegment, 
            quote.ExchangeInstrumentID, 
            quote.LastTradedPrice)
    }
}()
```

### Using Retry Mechanism

```go
// Set up retry with context
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Execute operation with retry
err = recovery.RetryWithBackoff(ctx, recovery.DefaultRetryConfig(), func() error {
    _, err := orderService.GetOrderBook("")
    return err
})

if err != nil {
    log.Printf("Operation failed after retries: %v", err)
    return
}
```

## Error Handling

The implementation includes comprehensive error handling:

1. **Custom Error Types**: Specific error types for different categories of errors
2. **Error Checking Functions**: Helper functions to check error types
3. **Detailed Error Information**: Errors include code, message, and description
4. **Context Preservation**: Errors maintain context through the call stack

Example:

```go
if err != nil {
    if errors.IsAuthError(err) {
        // Handle authentication error
    } else if errors.IsNetworkError(err) {
        // Handle network error
    } else if errors.IsOrderError(err) {
        // Handle order-related error
    } else {
        // Handle other errors
    }
}
```

## Resilience Patterns

The implementation includes several resilience patterns:

1. **Retry with Backoff**: Automatically retry operations with exponential backoff
2. **Circuit Breaker**: Prevent cascading failures by failing fast when the system is unhealthy
3. **Rate Limiting**: Prevent API throttling by limiting request rates
4. **Timeout Handling**: Proper handling of timeouts to prevent resource leaks
5. **Graceful Degradation**: Fallback mechanisms when services are unavailable

## Integration with Backend Gateway

The XTS integration is connected to the backend gateway through the `XTSController`:

1. **API Endpoints**: RESTful API endpoints for all XTS functionality
2. **Authentication**: Secure authentication and session management
3. **Request Validation**: Validation of incoming requests
4. **Error Handling**: Consistent error responses
5. **WebSocket Support**: Support for real-time data streaming

## Performance Considerations

The implementation is optimized for performance:

1. **Connection Pooling**: Reuse HTTP connections for better performance
2. **Efficient JSON Parsing**: Optimized JSON parsing for minimal overhead
3. **Memory Management**: Minimize allocations in hot paths
4. **Concurrency Control**: Appropriate synchronization primitives
5. **Buffering**: Buffer WebSocket messages to handle high throughput

## Testing

The implementation includes comprehensive tests:

1. **Unit Tests**: Test individual components in isolation
2. **Integration Tests**: Test interaction between components
3. **Mock Services**: Mock XTS services for testing
4. **Resilience Tests**: Test error handling and recovery mechanisms

## Deployment

The XTS integration can be deployed as part of the backend gateway or as a standalone service:

1. **Backend Integration**: Include as part of the main backend service
2. **Standalone Service**: Deploy as a separate microservice
3. **Containerization**: Package in Docker containers for easy deployment
4. **Configuration**: Configure through environment variables or config files

## Conclusion

This Go implementation of the XTS integration provides a high-performance, low-latency alternative to the Python SDK. It includes all the functionality of the original SDK while adding robust error handling, resilience patterns, and comprehensive testing. The modular architecture allows for easy maintenance and extension in the future.
