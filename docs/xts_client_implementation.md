# XTS Client Adapter Implementation Documentation

## Overview

This document provides comprehensive documentation for the XTS Client adapter implementation, which enables integration with the XTS Client API in our trading platform. The implementation follows the common broker interface pattern while accommodating XTS Client-specific features.

## Key Components

### 1. Core Client Structure

The XTS Client adapter is implemented in the `client` package with the following key files:

- **client.go**: Core client structure and common methods
- **place_order.go**: Order placement implementation with apiOrderSource parameter
- **dealer.go**: Dealer-specific endpoints implementation
- **client_test.go**: Comprehensive unit tests

### 2. Key Features

#### Authentication

The adapter provides robust authentication handling:
- Login with API key and secret key
- Session token management
- Logout functionality
- Support for both investor and dealer clients

#### Order Management

Complete order management capabilities:
- Place orders with apiOrderSource parameter
- Retrieve order book
- Support for various order types (MARKET, LIMIT, etc.)
- Support for different product types (MIS, NRML, etc.)

#### Dealer-Specific Features

Comprehensive dealer functionality:
- Get dealer order book
- Get dealer trades
- Get dealer positions
- Place orders on behalf of clients

### 3. Integration with Common Interface

The XTS Client adapter implements the common `BrokerClient` interface, allowing it to be used interchangeably with other broker implementations through the factory pattern.

## Implementation Details

### 1. Client Initialization

```go
// NewXTSClientImpl creates a new XTS Client implementation
func NewXTSClientImpl(config *common.XTSClientConfig) (*XTSClientImpl, error) {
    if config == nil {
        return nil, errors.New("XTS Client configuration is required")
    }

    if config.APIKey == "" || config.SecretKey == "" {
        return nil, errors.New("API key and secret key are required")
    }

    baseURL := config.BaseURL
    if baseURL == "" {
        baseURL = "https://developers.symphonyfintech.in" // Default URL
    }

    return &XTSClientImpl{
        apiKey:     config.APIKey,
        secretKey:  config.SecretKey,
        source:     config.Source,
        baseURL:    baseURL,
        debug:      false,
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }, nil
}
```

### 2. Authentication

```go
// Login authenticates with the XTS Client API
func (c *XTSClientImpl) Login(credentials *common.Credentials) (*common.Session, error) {
    // Use credentials if provided, otherwise use the configured API key and secret key
    apiKey := c.apiKey
    secretKey := c.secretKey
    
    if credentials != nil && credentials.APIKey != "" && credentials.SecretKey != "" {
        apiKey = credentials.APIKey
        secretKey = credentials.SecretKey
    }

    // API call implementation...
}
```

### 3. Order Placement with apiOrderSource

```go
// PlaceOrder places a new order with the XTS Client API
func (c *XTSClientImpl) PlaceOrder(order *common.Order) (*common.OrderResponse, error) {
    // Validation and setup...
    
    // Add the apiOrderSource parameter which is specific to XTS Client
    if order.APIOrderSource != "" {
        params.Set("apiOrderSource", order.APIOrderSource)
    } else {
        // Default value if not provided
        params.Set("apiOrderSource", "WEBAPI")
    }
    
    // API call implementation...
}
```

### 4. Dealer-Specific Endpoints

```go
// GetDealerOrderBook retrieves the dealer order book for the specified client
func (c *XTSClientImpl) GetDealerOrderBook(clientID string) (*common.OrderBook, error) {
    if c.token == "" {
        return nil, errors.New("not logged in")
    }
    
    if c.isInvestor {
        return nil, errors.New("dealer endpoints are not available for investor clients")
    }
    
    // API call implementation...
}
```

## Error Handling

The implementation includes comprehensive error handling:

1. **Authentication Errors**: Proper handling of login failures and token validation
2. **API Errors**: Parsing and propagation of API error responses
3. **Validation Errors**: Input validation before making API calls
4. **Client Type Validation**: Checking if dealer endpoints are being accessed by investor clients

## Testing

The implementation includes thorough unit tests:

1. **Mock Server**: Using httptest package to mock API responses
2. **Success Cases**: Testing normal operation flows
3. **Error Cases**: Testing error handling and edge cases
4. **Validation**: Verifying input validation and parameter handling

## Factory Integration

The XTS Client adapter is integrated with the broker factory:

```go
func NewBrokerClient(config *common.BrokerConfig) (common.BrokerClient, error) {
    switch config.BrokerType {
    case common.BrokerTypeXTSPro:
        // XTS Pro implementation...
    case common.BrokerTypeXTSClient:
        if config.XTSClient == nil {
            return nil, errors.New("XTS Client configuration is required")
        }
        return client.NewXTSClientImpl(config.XTSClient)
    case common.BrokerTypeZerodha:
        // Zerodha implementation...
    default:
        return nil, errors.New("unsupported broker type")
    }
}
```

## Usage Examples

### Creating an XTS Client

```go
config := &common.BrokerConfig{
    BrokerType: common.BrokerTypeXTSClient,
    XTSClient: &common.XTSClientConfig{
        APIKey:    "your_api_key",
        SecretKey: "your_secret_key",
        Source:    "WEBAPI",
    },
}

client, err := factory.NewBrokerClient(config)
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}

// Login
session, err := client.Login(nil)
if err != nil {
    log.Fatalf("Failed to login: %v", err)
}
```

### Placing an Order

```go
order := &common.Order{
    ExchangeSegment:       "NSECM",
    ExchangeInstrumentID:  "123456",
    ProductType:           "MIS",
    OrderType:             "LIMIT",
    OrderSide:             "BUY",
    TimeInForce:           "DAY",
    OrderQuantity:         10,
    LimitPrice:            100.5,
    OrderUniqueIdentifier: "test123",
    APIOrderSource:        "WEBAPI", // XTS Client specific parameter
}

response, err := client.PlaceOrder(order)
if err != nil {
    log.Fatalf("Failed to place order: %v", err)
}
```

### Using Dealer Features

```go
// Get dealer order book
orderBook, err := xtsClient.GetDealerOrderBook("client123")
if err != nil {
    log.Fatalf("Failed to get dealer order book: %v", err)
}

// Get dealer positions
positions, err := xtsClient.GetDealerPositions("client123")
if err != nil {
    log.Fatalf("Failed to get dealer positions: %v", err)
}

// Place order on behalf of client
order.ClientID = "client123"
response, err := xtsClient.PlaceDealerOrder(order)
if err != nil {
    log.Fatalf("Failed to place dealer order: %v", err)
}
```

## Conclusion

The XTS Client adapter implementation provides a robust, well-tested integration with the XTS Client API. It follows the common broker interface pattern while accommodating XTS Client-specific features like the apiOrderSource parameter and dealer functionality. The implementation is ready for use in the trading platform and can be easily extended with additional features as needed.
