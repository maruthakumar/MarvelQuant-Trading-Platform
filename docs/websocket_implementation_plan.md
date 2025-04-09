# WebSocket Implementation Plan for Real-Time Market Data

## Overview

This document outlines the implementation plan for WebSocket functionality to provide real-time market data across all broker types in the Trading Platform. The implementation will follow a unified approach that abstracts away broker-specific differences while maintaining high performance and reliability.

## Architecture

The WebSocket implementation will follow a layered architecture:

1. **Broker-Specific WebSocket Clients**: Connect to each broker's WebSocket API
2. **Unified WebSocket Manager**: Manages connections and provides a consistent interface
3. **Data Normalization Layer**: Converts broker-specific data formats to a common format
4. **Event Distribution System**: Distributes market data events to subscribers
5. **Reconnection and Recovery**: Handles connection failures and data recovery

## Implementation Components

### 1. Broker-Specific WebSocket Clients

#### XTS WebSocket Client

```go
// xts_websocket.go
package xts

import (
    "github.com/gorilla/websocket"
    "github.com/trading-platform/backend/internal/broker/common"
)

type XTSWebSocketClient struct {
    conn           *websocket.Conn
    url            string
    token          string
    subscriptions  map[string]bool
    messageChannel chan []byte
    errorChannel   chan error
    done           chan struct{}
}

func NewXTSWebSocketClient(url, token string) (*XTSWebSocketClient, error) {
    // Implementation
}

func (c *XTSWebSocketClient) Connect() error {
    // Implementation
}

func (c *XTSWebSocketClient) Subscribe(symbols []string) error {
    // Implementation
}

func (c *XTSWebSocketClient) Unsubscribe(symbols []string) error {
    // Implementation
}

func (c *XTSWebSocketClient) Close() error {
    // Implementation
}

func (c *XTSWebSocketClient) ReadMessage() ([]byte, error) {
    // Implementation
}
```

#### Zerodha WebSocket Client

```go
// zerodha_websocket.go
package zerodha

import (
    "github.com/gorilla/websocket"
    "github.com/trading-platform/backend/internal/broker/common"
)

type ZerodhaWebSocketClient struct {
    conn           *websocket.Conn
    url            string
    apiKey         string
    accessToken    string
    subscriptions  map[string]bool
    messageChannel chan []byte
    errorChannel   chan error
    done           chan struct{}
}

func NewZerodhaWebSocketClient(url, apiKey, accessToken string) (*ZerodhaWebSocketClient, error) {
    // Implementation
}

// Similar methods as XTS WebSocket Client
```

### 2. Unified WebSocket Manager

```go
// websocket_manager.go
package websocket

import (
    "sync"
    "github.com/trading-platform/backend/internal/broker/common"
)

type WebSocketManager struct {
    clients       map[string]WebSocketClient
    subscriptions map[string]map[string]bool // clientID -> symbol -> bool
    dataChannel   chan common.MarketData
    errorChannel  chan error
    mu            sync.RWMutex
}

func NewWebSocketManager() *WebSocketManager {
    // Implementation
}

func (m *WebSocketManager) RegisterClient(clientID string, client WebSocketClient) error {
    // Implementation
}

func (m *WebSocketManager) Subscribe(clientID string, symbols []string) error {
    // Implementation
}

func (m *WebSocketManager) Unsubscribe(clientID string, symbols []string) error {
    // Implementation
}

func (m *WebSocketManager) GetDataChannel() <-chan common.MarketData {
    // Implementation
}

func (m *WebSocketManager) GetErrorChannel() <-chan error {
    // Implementation
}

func (m *WebSocketManager) Start() error {
    // Implementation
}

func (m *WebSocketManager) Stop() error {
    // Implementation
}
```

### 3. Data Normalization Layer

```go
// data_normalizer.go
package websocket

import (
    "encoding/json"
    "github.com/trading-platform/backend/internal/broker/common"
)

type DataNormalizer struct {
    brokerType string
}

func NewDataNormalizer(brokerType string) *DataNormalizer {
    // Implementation
}

func (n *DataNormalizer) Normalize(data []byte) (common.MarketData, error) {
    // Implementation based on broker type
    switch n.brokerType {
    case common.BrokerTypeXTSPro, common.BrokerTypeXTSClient:
        return n.normalizeXTS(data)
    case common.BrokerTypeZerodha:
        return n.normalizeZerodha(data)
    default:
        return common.MarketData{}, fmt.Errorf("unsupported broker type: %s", n.brokerType)
    }
}

func (n *DataNormalizer) normalizeXTS(data []byte) (common.MarketData, error) {
    // Implementation
}

func (n *DataNormalizer) normalizeZerodha(data []byte) (common.MarketData, error) {
    // Implementation
}
```

### 4. Event Distribution System

```go
// event_distributor.go
package websocket

import (
    "sync"
    "github.com/trading-platform/backend/internal/broker/common"
)

type EventDistributor struct {
    subscribers map[string][]chan common.MarketData // symbol -> channels
    mu          sync.RWMutex
}

func NewEventDistributor() *EventDistributor {
    // Implementation
}

func (d *EventDistributor) Subscribe(symbol string) (<-chan common.MarketData, error) {
    // Implementation
}

func (d *EventDistributor) Unsubscribe(symbol string, ch <-chan common.MarketData) error {
    // Implementation
}

func (d *EventDistributor) Distribute(data common.MarketData) {
    // Implementation
}
```

### 5. Reconnection and Recovery

```go
// reconnection.go
package websocket

import (
    "time"
    "github.com/trading-platform/backend/internal/broker/common"
)

type ReconnectionManager struct {
    client         WebSocketClient
    maxRetries     int
    retryInterval  time.Duration
    subscriptions  []string
    onReconnect    func() error
}

func NewReconnectionManager(client WebSocketClient, maxRetries int, retryInterval time.Duration) *ReconnectionManager {
    // Implementation
}

func (r *ReconnectionManager) SetSubscriptions(subscriptions []string) {
    // Implementation
}

func (r *ReconnectionManager) SetOnReconnect(callback func() error) {
    // Implementation
}

func (r *ReconnectionManager) HandleDisconnect() error {
    // Implementation with exponential backoff
}
```

## API Integration

The WebSocket functionality will be integrated with the existing API layer:

```go
// broker_manager.go (extension)
package api

// SubscribeToMarketData subscribes to real-time market data for the specified symbols
func (m *BrokerManager) SubscribeToMarketData(userID string, symbols []string) (<-chan common.MarketData, error) {
    // Implementation
}

// UnsubscribeFromMarketData unsubscribes from real-time market data for the specified symbols
func (m *BrokerManager) UnsubscribeFromMarketData(userID string, symbols []string, ch <-chan common.MarketData) error {
    // Implementation
}
```

## Testing Strategy

1. **Unit Tests**: Test individual components in isolation
   - Test data normalization for each broker type
   - Test event distribution with mock subscribers
   - Test reconnection logic with simulated disconnections

2. **Integration Tests**: Test the integration between components
   - Test WebSocket manager with mock WebSocket clients
   - Test end-to-end flow from connection to data distribution

3. **Mock WebSocket Server**: Create a mock WebSocket server for testing
   - Simulate broker-specific WebSocket responses
   - Test reconnection by simulating disconnections

## Performance Considerations

1. **Connection Pooling**: Reuse WebSocket connections when possible
2. **Efficient Data Structures**: Use efficient data structures for subscriptions and event distribution
3. **Buffered Channels**: Use buffered channels to prevent blocking
4. **Goroutine Management**: Carefully manage goroutines to prevent leaks
5. **Memory Usage**: Monitor and optimize memory usage, especially for high-frequency data

## Error Handling

1. **Connection Errors**: Handle connection failures with reconnection logic
2. **Data Parsing Errors**: Log and recover from data parsing errors
3. **Subscription Errors**: Handle subscription failures with retries
4. **Rate Limiting**: Implement rate limiting to prevent API throttling
5. **Logging**: Comprehensive logging for debugging and monitoring

## Implementation Timeline

1. **Week 1**: Implement broker-specific WebSocket clients
2. **Week 2**: Implement unified WebSocket manager and data normalization
3. **Week 3**: Implement event distribution and reconnection logic
4. **Week 4**: Integrate with API layer and implement tests
5. **Week 5**: Performance optimization and documentation

## Conclusion

This WebSocket implementation plan provides a comprehensive approach to real-time market data across all broker types. The implementation will follow a unified approach that abstracts away broker-specific differences while maintaining high performance and reliability.
