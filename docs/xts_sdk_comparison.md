# XTS SDK Comparison Analysis

## Overview

This document analyzes the differences between the XTS PRO SDK and the XTS Client SDK to guide the Go implementation for both versions.

## Key Differences

### 1. API Parameters

**XTS Client SDK (GitHub):**
- The `place_order` function includes an additional parameter `apiOrderSource` that is not present in the XTS PRO SDK:
  ```python
  def place_order(self, exchangeSegment, exchangeInstrumentID, productType, orderType, 
                 orderSide, timeInForce, disclosedQuantity, orderQuantity, 
                 limitPrice, stopPrice, orderUniqueIdentifier, apiOrderSource, clientID=None):
  ```

**XTS PRO SDK:**
- The `place_order` function does not include the `apiOrderSource` parameter:
  ```python
  def place_order(self, exchangeSegment, exchangeInstrumentID, productType, orderType, 
                 orderSide, timeInForce, disclosedQuantity, orderQuantity, 
                 limitPrice, stopPrice, orderUniqueIdentifier, clientID=None):
  ```

### 2. Dealer API Support

**XTS Client SDK (GitHub):**
- Includes dealer-specific API endpoints:
  ```python
  "portfolio.dealerpositions": "interactive/portfolio/dealerpositions",
  "order.dealer.status": "/interactive/orders/dealerorderbook",
  "dealer.trades": "/interactive/orders/dealertradebook",
  ```
- Includes dealer-specific functions:
  ```python
  def get_dealer_orderbook(self, clientID=None):
  ```

**XTS PRO SDK:**
- Does not include dealer-specific endpoints or functions

### 3. Authentication Mechanism

Both SDKs use similar authentication mechanisms, but there might be subtle differences in token handling and session management.

### 4. Error Handling

Both SDKs use similar error handling approaches, but there might be differences in error codes and messages.

## Implementation Strategy for Go

### 1. Common Interface

Create a common interface for both XTS PRO and XTS Client implementations:

```go
type XTSClient interface {
    Login() (*models.Session, error)
    Logout() error
    PlaceOrder(order *models.Order) (*models.OrderResponse, error)
    ModifyOrder(order *models.ModifyOrder) (*models.OrderResponse, error)
    CancelOrder(orderID string, clientID string) (*models.OrderResponse, error)
    GetOrderBook(clientID string) (*models.OrderBook, error)
    GetPositions(clientID string) ([]models.Position, error)
    GetHoldings(clientID string) ([]models.Holding, error)
    // ... other common methods
}
```

### 2. Specific Implementations

#### XTS PRO Implementation

```go
type XTSProClient struct {
    // Common fields
    baseURL    string
    apiKey     string
    secretKey  string
    token      string
    httpClient *http.Client
}

func (c *XTSProClient) PlaceOrder(order *models.Order) (*models.OrderResponse, error) {
    // XTS PRO specific implementation
}
```

#### XTS Client Implementation

```go
type XTSClientImpl struct {
    // Common fields (same as XTSProClient)
    baseURL    string
    apiKey     string
    secretKey  string
    token      string
    httpClient *http.Client
}

func (c *XTSClientImpl) PlaceOrder(order *models.Order) (*models.OrderResponse, error) {
    // XTS Client specific implementation
    // Include apiOrderSource parameter
}

// Additional dealer-specific methods
func (c *XTSClientImpl) GetDealerOrderBook(clientID string) (*models.OrderBook, error) {
    // XTS Client specific implementation
}
```

### 3. Factory Pattern

Use a factory pattern to create the appropriate client based on configuration:

```go
func NewXTSClient(config *config.XTSConfig) (XTSClient, error) {
    if config.ClientType == "PRO" {
        return NewXTSProClient(config)
    } else if config.ClientType == "CLIENT" {
        return NewXTSClientImpl(config)
    }
    return nil, errors.New("invalid client type")
}
```

### 4. Models Adaptation

Extend the models to support both implementations:

```go
type Order struct {
    // Common fields
    ExchangeSegment      string
    ExchangeInstrumentID string
    ProductType          string
    OrderType            string
    OrderSide            string
    TimeInForce          string
    DisclosedQuantity    int
    OrderQuantity        int
    LimitPrice           float64
    StopPrice            float64
    OrderUniqueIdentifier string
    
    // XTS Client specific fields
    APIOrderSource       string `json:"apiOrderSource,omitempty"`
}
```

## Conclusion

The main differences between XTS PRO SDK and XTS Client SDK are:

1. The addition of the `apiOrderSource` parameter in the `place_order` function
2. Dealer API support in the XTS Client SDK
3. Potential differences in authentication and error handling

The Go implementation should use a common interface with specific implementations for each SDK type, using a factory pattern to create the appropriate client based on configuration. The models should be extended to support both implementations, with optional fields for XTS Client specific parameters.

This approach will allow for a clean, maintainable codebase that supports both XTS PRO and XTS Client users while minimizing code duplication.
